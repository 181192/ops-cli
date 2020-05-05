package dashboard

import (
	"errors"
	"fmt"

	"github.com/181192/ops-cli/pkg/cmd/cmdutils"
	"github.com/181192/ops-cli/pkg/kubernetes"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	v1 "k8s.io/api/core/v1"
)

func grafanaDashboardCmd(cmd *cmdutils.Cmd) {
	var opts Options

	cmd.CobraCommand.Use = "grafana"
	cmd.CobraCommand.Short = "Open Grafana web UI"
	cmd.CobraCommand.Long = "Open Grafana dashboard"
	cmd.CobraCommand.RunE = func(_ *cobra.Command, args []string) error {
		cmd.NameArg = cmdutils.GetNameArg(args)
		return doPortForwardGrafana(cmd, opts)
	}

	cmd.FlagSetGroup.InFlagSet("Dashboard", func(fs *pflag.FlagSet) {
		cmdutils.AddCommonFlagsForKubernetes(fs, &opts.KubeOptions)
	})
}

func doPortForwardGrafana(cmd *cmdutils.Cmd, opts Options) error {
	kubeConfig := opts.KubeOptions.KubeConfig
	kubeContext := opts.KubeOptions.KubeContext
	namespace := opts.KubeOptions.Namespace

	client, err := kubernetes.NewClient(kubeConfig, kubeContext)
	if err != nil {
		return fmt.Errorf("failed to create k8s client: %v", err)
	}

	if namespace == "" {
		namespace = "monitoring"
	}

	pl, err := client.PodsForSelector(namespace, "app.kubernetes.io/name=grafana")
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
		"http://localhost:%d", 3000, client, cmd.CobraCommand.OutOrStdout())
}
