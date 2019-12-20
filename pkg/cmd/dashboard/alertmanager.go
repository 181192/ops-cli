package dashboard

import (
	"errors"
	"fmt"

	"github.com/181192/ops-cli/pkg/kubernetes"
	"github.com/spf13/cobra"
)

func alertmanagerDashCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "alertmanager",
		Short: "Open Alertmanager web UI",
		Long:  `Open Alertmanager dashboard`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := kubernetes.NewClient(kubeconfig, configContext)
			if err != nil {
				return fmt.Errorf("failed to create k8s client: %v", err)
			}

			if namespace == "" {
				namespace = "monitoring"
			}

			pl, err := client.PodsForSelector(namespace, "app=alertmanager")
			if err != nil {
				return fmt.Errorf("not able to locate Alertmanager pod: %v", err)
			}

			if len(pl.Items) < 1 {
				return errors.New("no Alertmanager pods found")
			}

			// only use the first pod in the list
			return portForward(pl.Items[0].Name, namespace, "Alertmanager",
				"http://localhost:%d", 9093, client, cmd.OutOrStdout())
		},
	}

	return cmd
}
