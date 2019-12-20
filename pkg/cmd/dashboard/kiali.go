package dashboard

import (
	"errors"
	"fmt"

	"github.com/181192/ops-cli/pkg/kubernetes"
	"github.com/spf13/cobra"
)

func kialiDashCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kiali",
		Short: "Open Kiali web UI",
		Long:  `Open Istio's Kiali dashboard`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := kubernetes.NewClient(kubeconfig, configContext)
			if err != nil {
				return fmt.Errorf("failed to create k8s client: %v", err)
			}

			if namespace == "" {
				namespace = "istio-system"
			}

			pl, err := client.PodsForSelector(namespace, "app=kiali")
			if err != nil {
				return fmt.Errorf("not able to locate Kiali pod: %v", err)
			}

			if len(pl.Items) < 1 {
				return errors.New("no Kiali pods found")
			}

			// only use the first pod in the list
			return portForward(pl.Items[0].Name, namespace, "Kiali",
				"http://localhost:%d/kiali", 20001, client, cmd.OutOrStdout())
		},
	}

	return cmd
}
