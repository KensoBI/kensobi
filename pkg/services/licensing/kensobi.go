package licensing

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"os"
	"strings"
	"time"

	"github.com/grafana/grafana/pkg/api/dtos"
	contextmodel "github.com/grafana/grafana/pkg/services/contexthandler/model"
	"github.com/grafana/grafana/pkg/services/hooks"
	"github.com/grafana/grafana/pkg/services/navtree"
	"github.com/grafana/grafana/pkg/setting"
)

const (
	KensoBIEdition = "KensoBI"
	OssEdition     = "oss"
	publicKey      = `-----BEGIN PUBLIC KEY-----
MIIBCgKCAQEA2gA2lRtZoQbovt03x2mwtKXRNJY+PcX7vZXYQTLeQWMYMBHG+40I
TQ1mkZfGqTEAXX5zwqZP6UBcvg+vGkP5VzxFy3SZrftd5c5XN+CnD3Zcvar4muI/
qEC0SoyW3u5r4HEz/zIpgEmhLjUnu3hcsfN94GDtc17kyO2XiVsxpl20dUiBT4TD
QCJM+bazE1IpnP7nXTfu+F4wXL6m0iRCHUuVphtOEmUtIhgQ0+1xiszF2utIYr8V
5CuOH3lzfKKwdYUIG4oGhhPb3PRcz6rXSKwLFQDc24uUNE9KcgLkSR8qIW+GYEF+
rabOjpOMw2Ajx2ojcgt2kW4++JfxC/oADQIDAQAB
-----END PUBLIC KEY-----
`
)

type licenseContent struct {
	UserName string `json:"user_name"`
	Expiry   string `json:"expiry"`
}

type KensoBILicensingService struct {
	Cfg          *setting.Cfg
	HooksService *hooks.HooksService
	license      *licenseContent
}

// Expiry returns the unix epoch timestamp when the license expires, or 0 if no valid license is provided
func (k *KensoBILicensingService) Expiry() int64 {
	if k.license == nil {
		return 0
	}

	if k.license.Expiry != "never" {
		expiryDate, err := time.Parse("2006-01-02", k.license.Expiry)
		if err != nil {
			k.Cfg.Logger.Error("Error parsing expiry date")
			return 0
		}
		if time.Now().After(expiryDate) {
			k.Cfg.Logger.Error("License expired")
			return 0
		}
		return expiryDate.Unix()
	}

	//"never" - return epoch time in seconds 100 years from now
	return time.Now().AddDate(100, 0, 0).Unix()
}

func (k *KensoBILicensingService) HasLicense() bool {
	return k.Expiry() > 0
}

func (k *KensoBILicensingService) Edition() string {
	if k.HasLicense() {
		return KensoBIEdition
	}
	return OssEdition
}

func (*KensoBILicensingService) StateInfo() string {
	return ""
}

func (*KensoBILicensingService) ContentDeliveryPrefix() string {
	return "grafana-oss"
}

func (l *KensoBILicensingService) LicenseURL(showAdminLicensingPage bool) string {
	if showAdminLicensingPage {
		return l.Cfg.AppSubURL + "/admin/kenso-license"
	}

	return "https://kensobi.com/"
}

func (*KensoBILicensingService) EnabledFeatures() map[string]bool {
	return map[string]bool{}
}

func (k *KensoBILicensingService) FeatureEnabled(feature string) bool {
	k.Cfg.Logger.Info("KensoBI feature enabled", "feature", feature)
	return false
}

func (k *KensoBILicensingService) onInvalidated() {
	//TODO restart whole application or stop all kenso plugins
	k.Cfg.Logger.Info("KensoBI license invalidated")
}

func (k *KensoBILicensingService) checkLicense() bool {
	defer func() {
		if r := recover(); r != nil {
			k.Cfg.Logger.Error("Failed to check license", "error", r)
		}
	}()

	readLicense := func() *licenseContent {
		path := k.Cfg.EnterpriseLicensePath
		if path == "" {
			return nil
		}
		licensePEM, err := os.ReadFile(path)
		if err != nil {
			k.Cfg.Logger.Info("No license file found")
			return nil
		}

		block, _ := pem.Decode([]byte(publicKey))
		if block == nil {
			k.Cfg.Logger.Error("Error decoding public key")
			return nil
		}
		pubKey, err := x509.ParsePKCS1PublicKey(block.Bytes)
		if err != nil {
			k.Cfg.Logger.Error("Error parsing public key")
			return nil
		}
		license := string(licensePEM)
		split := strings.Split(license, ".")
		if len(split) != 2 {
			k.Cfg.Logger.Error("Error parsing license")
			return nil
		}
		license = split[0]
		signature := split[1]

		//decode signature
		signatureBytes, err := base64.StdEncoding.DecodeString(signature)
		licenseBytes, err := base64.StdEncoding.DecodeString(license)
		if err != nil {
			k.Cfg.Logger.Error("Error decoding license")
			return nil
		}

		//verify signature
		hashed := sha256.Sum256(licenseBytes)
		err = rsa.VerifyPKCS1v15(pubKey, crypto.SHA256, hashed[:], signatureBytes)
		if err != nil {
			k.Cfg.Logger.Error("Error verifying signature")
			return nil
		}

		//parse license
		var licenseObj licenseContent
		err = json.Unmarshal(licenseBytes, &licenseObj)
		if err != nil {
			k.Cfg.Logger.Error("Error parsing license")
			return nil
		}

		return &licenseObj
	}

	prevState := k.HasLicense()
	k.license = readLicense()
	isValid := k.HasLicense()
	if !isValid {
		k.license = nil
	}

	if prevState != isValid && prevState {
		k.onInvalidated()
	}

	return isValid
}

func (k *KensoBILicensingService) startChecker() {
	go func() {
		for {
			if !k.checkLicense() {
				break
			}
			time.Sleep(1 * time.Hour)
		}
	}()
}

func ProvideService(cfg *setting.Cfg, hooksService *hooks.HooksService) *KensoBILicensingService {
	l := &KensoBILicensingService{
		Cfg:          cfg,
		HooksService: hooksService,
		license:      nil,
	}
	l.startChecker()
	l.HooksService.AddIndexDataHook(func(indexData *dtos.IndexViewData, req *contextmodel.ReqContext) {
		if !req.IsGrafanaAdmin {
			return
		}

		var adminNodeID string

		if cfg.IsFeatureToggleEnabled("topnav") {
			adminNodeID = navtree.NavIDCfg
		} else {
			adminNodeID = navtree.NavIDAdmin
		}

		if adminNode := indexData.NavTree.FindById(adminNodeID); adminNode != nil {
			if req.IsGrafanaAdmin {
				adminNode.Children = append(adminNode.Children, &navtree.NavLink{
					Text: "Stats",
					Id:   "stats",
					Url:  l.Cfg.AppSubURL + "/admin/stats",
					Icon: "unlock",
				})
				adminNode.Children = append(adminNode.Children, &navtree.NavLink{
					Text: "KensoBI License",
					Id:   "kenso-license",
					Url:  l.Cfg.AppSubURL + "/admin/kenso-license",
					Icon: "unlock",
				})
			}
		}
	})

	return l
}
