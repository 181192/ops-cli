package dashboard

import (
	"github.com/181192/ops-cli/pkg/config"
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Dashboard represents input to kubectl port-forward
type Dashboard struct {
	Name          string `json:"name"`
	Namespace     string `json:"namespace"`
	Port          int    `json:"port"`
	LabelSelector string `json:"labelSelector"`
	URL           string `json:"url,omitempty"`
}

// MakeDashboards returns list of all dashboards
func MakeDashboards() []Dashboard {
	dashboards := []Dashboard{}

	dashboards = append(dashboards, Dashboard{
		Name:          "alertmanager",
		Namespace:     "monitoring",
		Port:          9093,
		LabelSelector: "app=alertmanager",
	})

	dashboards = append(dashboards, Dashboard{
		Name:          "grafana",
		Namespace:     "monitoring",
		Port:          3000,
		LabelSelector: "app.kubernetes.io/name=grafana",
	})

	dashboards = append(dashboards, Dashboard{
		Name:          "jaeger",
		Namespace:     "tracing",
		Port:          16686,
		LabelSelector: "app=jaeger,app.kubernetes.io/name=jaeger-operator-jaeger",
	})

	dashboards = append(dashboards, Dashboard{
		Name:          "kiali",
		Namespace:     "istio-system",
		Port:          20001,
		LabelSelector: "app=kiali",
	})

	dashboards = append(dashboards, Dashboard{
		Name:          "prometheus",
		Namespace:     "monitoring",
		Port:          9090,
		LabelSelector: "app=prometheus",
	})

	config.InitConfig()

	var userDashboards []Dashboard
	if err := viper.UnmarshalKey("dashboards", &userDashboards); err != nil {
		logger.Error(err)
	}

	dashboards = append(dashboards, userDashboards...)

	return unique(dashboards)
}

func unique(dashboards []Dashboard) []Dashboard {
	var uniques []Dashboard

DashboardLoop:
	for _, dashboard := range dashboards {
		for i, u := range uniques {
			if dashboard.Name == u.Name {
				uniques[i] = dashboard
				continue DashboardLoop
			}
		}
		uniques = append(uniques, dashboard)
	}
	return uniques
}
