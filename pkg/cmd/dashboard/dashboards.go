package dashboard

// Dashboard represents input to kubectl port-forward
type Dashboard struct {
	Name          string
	Namespace     string
	Port          int
	LabelSelector string
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

	return dashboards
}
