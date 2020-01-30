package dashboard

import (
	"fmt"
	"io"

	"github.com/181192/ops-cli/pkg/cmd/cmdutils"
	"github.com/181192/ops-cli/pkg/kubernetes"
	"github.com/181192/ops-cli/pkg/open"

	"github.com/spf13/cobra"
)

// Options dashboard options
type Options struct {
	KubeOptions   cmdutils.KubernetesOpts
	Label         string
	Port          string
	labelSelector string
}

// dashboardCmd represents the dashboard command
var dashboardCmd = &cobra.Command{
	Use:   "dashboard",
	Short: "Access to various web UIs",
	Long:  `Access to various web UIs`,
}

// Command will create the `dashboard` commands
func Command(flagGrouping *cmdutils.FlagGrouping) *cobra.Command {

	cmdutils.AddResourceCmd(flagGrouping, dashboardCmd, kialiDashboardCmd)
	cmdutils.AddResourceCmd(flagGrouping, dashboardCmd, prometheusDashboardCmd)
	cmdutils.AddResourceCmd(flagGrouping, dashboardCmd, grafanaDashboardCmd)
	cmdutils.AddResourceCmd(flagGrouping, dashboardCmd, jaegerDashboardCmd)
	cmdutils.AddResourceCmd(flagGrouping, dashboardCmd, alertmanagerDashboardCmd)

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
