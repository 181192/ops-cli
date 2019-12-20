package dashboard

import (
	"errors"
	"fmt"
	"io"

	"github.com/181192/ops-cli/pkg/kubernetes"
	"github.com/181192/ops-cli/pkg/open"

	"github.com/spf13/cobra"

	v1 "k8s.io/api/core/v1"
)

var (
	labelSelector = ""

	kubeconfig    string
	configContext string
	namespace     string
	label         string
	port          int
)

// dashboardCmd represents the dashboard command
var dashboardCmd = &cobra.Command{
	Use:   "dashboard",
	Short: "Access to various web UIs",
	Long:  `Access to various web UIs`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if label == "" {
			return errors.New("no labels given to filter")
		}

		client, err := kubernetes.NewClient(kubeconfig, configContext)
		if err != nil {
			return fmt.Errorf("failed to create k8s client: %v", err)
		}

		pl, err := client.PodsForSelector(namespace, label)
		if err != nil {
			return fmt.Errorf("not able to locate pod: %v", err)
		}

		if len(pl.Items) < 1 {
			return errors.New("no pods found by " + label)
		}

		// only use the first pod in the list
		return portForward(pl.Items[0].Name, namespace, "Pod",
			"http://localhost:%d", port, client, cmd.OutOrStdout())
	},
}

// Command will create the `dashboard` commands
func Command() *cobra.Command {
	dashboardCmd.AddCommand(kialiDashCmd())
	dashboardCmd.AddCommand(promDashCmd())
	dashboardCmd.AddCommand(grafanaDashCmd())
	dashboardCmd.AddCommand(jaegerDashCmd())
	dashboardCmd.AddCommand(alertmanagerDashCmd())

	dashboardCmd.PersistentFlags().StringVarP(&kubeconfig, "kubeconfig", "c", "", "Kubernetes configuration file")
	dashboardCmd.PersistentFlags().StringVar(&configContext, "context", "", "The name of the kubeconfig context to use")
	dashboardCmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", v1.NamespaceAll, "Config namespace")

	dashboardCmd.MarkFlagRequired("label")
	dashboardCmd.Flags().StringVarP(&label, "label", "l", "", "key=value")

	dashboardCmd.MarkFlagRequired("port")
	dashboardCmd.Flags().IntVarP(&port, "port", "p", 0, "container port")

	return dashboardCmd

}

// portForward first tries to forward localhost:remotePort to podName:remotePort, falls back to dynamic local port
func portForward(podName, namespace, flavor, url string, remotePort int, client kubernetes.ExecClient, writer io.Writer) error {
	var err error
	for _, localPort := range []int{remotePort, 0} {
		fw, err := client.BuildPortForwarder(podName, namespace, localPort, remotePort)
		if err != nil {
			return fmt.Errorf("could not build port forwarder for %s: %v", flavor, err)
		}

		if err = kubernetes.RunPortForwarder(fw, func(fw *kubernetes.PortForward) error {
			fmt.Printf("port-forward to %s pod ready\n", flavor)
			fmt.Println(fmt.Sprintf(url, fw.LocalPort))
			open.Start(fmt.Sprintf(url, fw.LocalPort))
			return nil
		}); err == nil {
			return nil
		}
	}

	return fmt.Errorf("failure running port forward process: %v", err)
}
