package dashboard

import (
	"fmt"
	"io"

	"github.com/181192/ops-cli/pkg/cmd/cmdutils"
	"github.com/181192/ops-cli/pkg/kubernetes"
	"github.com/181192/ops-cli/pkg/open"

	logger "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	v1 "k8s.io/api/core/v1"
)

// Options dashboard options
type Options struct {
	KubeOptions   cmdutils.KubernetesOpts
	Port          int
	LabelSelector string
	URL           string
}

// Command will create the `dashboard` commands
func Command(flagGrouping *cmdutils.FlagGrouping) *cobra.Command {
	dashboards := MakeDashboards()
	helpText := "Dashboards"
	helpTextLong := "Available dashboards:\n\n"

	for _, d := range dashboards {
		helpTextLong += fmt.Sprintf("  - %s\n", d.Name)
	}

	dashboardCmd := &cobra.Command{
		Use:     "dashboard",
		Aliases: []string{"d"},
		Short:   helpText,
		Long:    helpTextLong,
	}

	return generateDashboardCommands(flagGrouping, dashboardCmd, dashboards)
}

func generateDashboardCommands(flagGrouping *cmdutils.FlagGrouping, dashboardCmd *cobra.Command, dashboards []Dashboard) *cobra.Command {
	for _, dashboard := range dashboards {
		var opts Options

		cmdutils.AddResourceCmd(flagGrouping, dashboardCmd, func(cmd *cmdutils.Cmd) {
			cmd.CobraCommand.Use = dashboard.Name
			cmd.CobraCommand.Run = func(_ *cobra.Command, args []string) {
				for _, d := range dashboards {
					if d.Name == cmd.CobraCommand.Use {
						dashboard = d
						break
					}
				}

				if err := doPortForward(cmd.CobraCommand, &opts, args, &dashboard); err != nil {
					logger.Fatal(err)
				}
			}
			cmd.FlagSetGroup.InFlagSet("Dashboard", func(fs *pflag.FlagSet) {
				fs.IntVarP(&opts.Port, "port", "p", 0, "Target port to forward to")
				fs.StringVarP(&opts.LabelSelector, "label-selector", "l", "",
					"Selector (label query) to filter on, supports '=', '==', and '!='.(e.g. -l key1=value1,key2=value2)")
				fs.StringVarP(&opts.URL, "url", "u", "",
					"Relative URL to open (e.g. /metrics)")
				cmdutils.AddCommonFlagsForKubernetes(fs, &opts.KubeOptions, dashboard.Namespace)
			})
		})
	}

	return dashboardCmd
}

func doPortForward(cmd *cobra.Command, opts *Options, args []string, dashboard *Dashboard) error {

	kubeConfig := opts.KubeOptions.KubeConfig
	kubeContext := opts.KubeOptions.KubeContext
	overrideNamespace := opts.KubeOptions.Namespace
	overridePort := opts.Port
	overrideLabelSelector := opts.LabelSelector
	overrideURL := opts.URL

	client, err := kubernetes.NewClient(kubeConfig, kubeContext)
	if err != nil {
		return fmt.Errorf("failed to create k8s client: %v", err)
	}

	if overrideNamespace != "" {
		dashboard.Namespace = overrideNamespace
	}

	if overridePort != 0 {
		dashboard.Port = overridePort
	}

	if overrideLabelSelector != "" {
		dashboard.LabelSelector = overrideLabelSelector
	}

	if overrideURL != "" {
		dashboard.URL = overrideURL
	}

	pl, err := client.PodsForSelector(dashboard.Namespace, dashboard.LabelSelector)
	if err != nil {
		return fmt.Errorf("not able to locate %s pod in %s namespace using selector %s: %v",
			dashboard.Name, dashboard.Namespace, dashboard.LabelSelector, err)
	}

	if len(pl.Items) < 1 {
		return fmt.Errorf("no %s pods found in %s namespace using selector %s",
			dashboard.Name, dashboard.Namespace, dashboard.LabelSelector)
	}

	var pod v1.Pod

	for _, p := range pl.Items {
		if p.Status.Phase == "Running" {
			pod = p
		}
	}

	if pod.Name == "" {
		return fmt.Errorf("no running %s pods found in %s namespace using selector %s",
			dashboard.Name, dashboard.Namespace, dashboard.LabelSelector)
	}

	// only use the first pod in the list
	return portForward(pod.Name, dashboard.Namespace, dashboard.Name,
		"http://localhost:%d%s", dashboard.Port, dashboard.URL, client, cmd.OutOrStdout())
}

// portForward first tries to forward localhost:remotePort to podName:remotePort, falls back to dynamic local port
func portForward(podName, namespace, flavor, url string, remotePort int, urlSuffix string, client kubernetes.ExecClient, writer io.Writer) error {
	var err error
	for _, localPort := range []int{remotePort, 0} {
		fw, err := client.BuildPortForwarder(podName, namespace, localPort, remotePort)
		if err != nil {
			return fmt.Errorf("could not build port forwarder for %s: %v", flavor, err)
		}

		if err = kubernetes.RunPortForwarder(fw, func(fw *kubernetes.PortForward) error {
			logger.Infof("port-forward to %s pod in %s namespace ready\n", flavor, namespace)
			logger.Infof(fmt.Sprintf(url, fw.LocalPort, urlSuffix))
			open.Start(fmt.Sprintf(url, fw.LocalPort, urlSuffix))
			return nil
		}); err == nil {
			return nil
		}
	}

	return fmt.Errorf("failure running port forward process: %v", err)
}
