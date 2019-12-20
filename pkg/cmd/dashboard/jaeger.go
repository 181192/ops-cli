package dashboard

import (
	"errors"
	"fmt"

	"github.com/181192/ops-cli/pkg/kubernetes"
	"github.com/spf13/cobra"
)

func jaegerDashCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "jaeger",
		Short: "Open Jaeger web UI",
		Long:  `Open Istio's Jaeger dashboard`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := kubernetes.NewClient(kubeconfig, configContext)
			if err != nil {
				return fmt.Errorf("failed to create k8s client: %v", err)
			}

			if namespace == "" {
				namespace = "istio-system"
			}

			pl, err := client.PodsForSelector(namespace, "app=jaeger")
			if err != nil {
				return fmt.Errorf("not able to locate Jaeger pod: %v", err)
			}

			if len(pl.Items) < 1 {
				return errors.New("no Jaeger pods found")
			}
			// only use the first pod in the list
			return portForward(pl.Items[0].Name, namespace, "Jaeger",
				"http://localhost:%d", 16686, client, cmd.OutOrStdout())
		},
	}

	return cmd
}
