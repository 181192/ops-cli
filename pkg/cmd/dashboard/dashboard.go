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
}

// Command will create the `dashboard` commands
func Command(flagGrouping *cmdutils.FlagGrouping) *cobra.Command {
	var opts Options
	dashboards := MakeDashboards()
	validArgs := []string{}
	helpText := "Available dashboards:\n\n"

	for _, d := range dashboards {
		helpText += fmt.Sprintf("  - %s\n", d.Name)
		validArgs = append(validArgs, d.Name)
	}

	dashboardCmd := &cobra.Command{
		Use:       "dashboard",
		Aliases:   []string{"d"},
		Short:     helpText,
		Long:      helpText,
		ValidArgs: validArgs,
		Run: func(cmd *cobra.Command, args []string) {
			if err := doPortForward(cmd, opts, args, dashboards); err != nil {
				logger.Fatal(err)
			}
		},
	}

	flagSetGroup := flagGrouping.New(dashboardCmd)

	flagSetGroup.InFlagSet("Dashboard", func(fs *pflag.FlagSet) {
		fs.IntVarP(&opts.Port, "port", "p", 0, "Target port to forward to")
		fs.StringVarP(&opts.LabelSelector, "label-selector", "l", "",
			"Selector (label query) to filter on, supports '=', '==', and '!='.(e.g. -l key1=value1,key2=value2)")
		cmdutils.AddCommonFlagsForKubernetes(fs, &opts.KubeOptions)
	})

	flagSetGroup.AddTo(dashboardCmd)

	return dashboardCmd
}

func doPortForward(cmd *cobra.Command, opts Options, args []string, dashboards []Dashboard) error {
	var dashboard *Dashboard

	if len(args) == 0 {
		return cmd.Help()
	}

	if len(args) > 1 {
		return fmt.Errorf("only one argument is allowed to be used as a name")
	}

	if len(args) == 1 {
		for _, d := range dashboards {
			if d.Name == args[0] {
				dashboard = &d
				break
			}
		}
	}

	if dashboard == nil {
		return fmt.Errorf("cannot get dashboard for: %s", args[0])
	}

	kubeConfig := opts.KubeOptions.KubeConfig
	kubeContext := opts.KubeOptions.KubeContext
	overrideNamespace := opts.KubeOptions.Namespace
	overridePort := opts.Port
	overrideLabelSelector := opts.LabelSelector

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
		"http://localhost:%d", dashboard.Port, client, cmd.OutOrStdout())
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
			logger.Infof("port-forward to %s pod in %s namespace ready\n", flavor, namespace)
			logger.Infof(fmt.Sprintf(url, fw.LocalPort))
			open.Start(fmt.Sprintf(url, fw.LocalPort))
			return nil
		}); err == nil {
			return nil
		}
	}

	return fmt.Errorf("failure running port forward process: %v", err)
}
