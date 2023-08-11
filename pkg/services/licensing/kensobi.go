package licensing

import (
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
)

type KensoBILicensingService struct {
	Cfg          *setting.Cfg
	HooksService *hooks.HooksService
	hasLicense   bool
}

func (k *KensoBILicensingService) Expiry() int64 {
	if !k.hasLicense {
		return 0
	}
	//unix epoch timestamp of NOW + 1 year
	return time.Now().AddDate(1, 0, 0).Unix() //TODO: change to real expiry
}

func (k *KensoBILicensingService) Edition() string {
	if k.hasLicense {
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

func (k *KensoBILicensingService) checkLicense() {
	defer func() {
		if r := recover(); r != nil {
			k.Cfg.Logger.Error("Failed to check license", "error", r)
		}
	}()

	isValid := func() bool {
		path := k.Cfg.EnterpriseLicensePath
		if path == "" {
			return false
		}
		//TODO check license with public key
		return true
	}()

	prevState := k.hasLicense

	//check license	here
	k.hasLicense = isValid

	if prevState != k.hasLicense && prevState {
		k.onInvalidated()
	}
}

func (k *KensoBILicensingService) startChecker() {
	go func() {
		for {
			k.checkLicense()
			time.Sleep(1 * time.Hour)
		}
	}()
}

func ProvideService(cfg *setting.Cfg, hooksService *hooks.HooksService) *KensoBILicensingService {
	l := &KensoBILicensingService{
		Cfg:          cfg,
		HooksService: hooksService,
		hasLicense:   false,
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
