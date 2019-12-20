package dashboard

import (
	"errors"
	"fmt"

	"github.com/181192/ops-cli/pkg/kubernetes"
	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
)

func grafanaDashCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "grafana",
		Short: "Open Grafana web UI",
		Long:  `Open Grafana dashboard`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := kubernetes.NewClient(kubeconfig, configContext)
			if err != nil {
				return fmt.Errorf("failed to create k8s client: %v", err)
			}

			if namespace == "" {
				namespace = "monitoring"
			}

			pl, err := client.PodsForSelector(namespace, "app=grafana")
			if err != nil {
				return fmt.Errorf("not able to locate Grafana pod: %v", err)
			}

			if len(pl.Items) < 1 {
				return errors.New("no Grafana pods found")
			}

			var pod v1.Pod

			for _, p := range pl.Items {
				if p.Status.Phase == "Running" {
					pod = p
				}
			}

			if pod.Name == "" {
				return errors.New("no running Grafana pods found")
			}

			// only use the first pod in the list
			return portForward(pod.Name, namespace, "Grafana",
				"http://localhost:%d", 3000, client, cmd.OutOrStdout())
		},
	}

	return cmd
}
