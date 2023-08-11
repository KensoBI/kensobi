package licensing

import (
	"github.com/grafana/grafana/pkg/api/dtos"
	contextmodel "github.com/grafana/grafana/pkg/services/contexthandler/model"
	"github.com/grafana/grafana/pkg/services/hooks"
	"github.com/grafana/grafana/pkg/services/navtree"
	"github.com/grafana/grafana/pkg/setting"
)

const (
	kensoBIEdition = "Kenso BI"
)

type KensoBILicensingService struct {
	Cfg          *setting.Cfg
	HooksService *hooks.HooksService
}

func (*KensoBILicensingService) Expiry() int64 {
	return 0
}

func (*KensoBILicensingService) Edition() string {
	return kensoBIEdition
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

func (*KensoBILicensingService) FeatureEnabled(feature string) bool {
	return false
}

func ProvideService(cfg *setting.Cfg, hooksService *hooks.HooksService) *KensoBILicensingService {
	l := &KensoBILicensingService{
		Cfg:          cfg,
		HooksService: hooksService,
	}
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
