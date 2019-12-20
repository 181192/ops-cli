package dashboard

import (
	"errors"
	"fmt"

	"github.com/181192/ops-cli/pkg/kubernetes"
	"github.com/spf13/cobra"
)

func promDashCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "prometheus",
		Short: "Open Prometheus web UI",
		Long:  `Open Prometheus dashboard`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := kubernetes.NewClient(kubeconfig, configContext)
			if err != nil {
				return fmt.Errorf("failed to create k8s client: %v", err)
			}

			if namespace == "" {
				namespace = "monitoring"
			}

			pl, err := client.PodsForSelector(namespace, "app=prometheus")
			if err != nil {
				return fmt.Errorf("not able to locate Prometheus pod: %v", err)
			}

			if len(pl.Items) < 1 {
				return errors.New("no Prometheus pods found")
			}

			// only use the first pod in the list
			return portForward(pl.Items[0].Name, namespace, "Prometheus",
				"http://localhost:%d", 9090, client, cmd.OutOrStdout())
		},
	}

	return cmd
}
