package dashboard

import (
	"errors"
	"fmt"

	"github.com/181192/ops-cli/pkg/cmd/cmdutils"
	"github.com/181192/ops-cli/pkg/kubernetes"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func prometheusDashboardCmd(cmd *cmdutils.Cmd) {
	var opts Options

	cmd.CobraCommand.Use = "prometheus"
	cmd.CobraCommand.Short = "Open Prometheus web UI"
	cmd.CobraCommand.Long = "Open Prometheus dashboard"
	cmd.CobraCommand.RunE = func(_ *cobra.Command, args []string) error {
		cmd.NameArg = cmdutils.GetNameArg(args)
		return doPortForwardPrometheus(cmd, opts)
	}

	cmd.FlagSetGroup.InFlagSet("Dashboard", func(fs *pflag.FlagSet) {
		cmdutils.AddCommonFlagsForKubernetes(fs, &opts.KubeOptions)
	})
}

func doPortForwardPrometheus(cmd *cmdutils.Cmd, opts Options) error {
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

	pl, err := client.PodsForSelector(namespace, "app=prometheus")
	if err != nil {
		return fmt.Errorf("not able to locate Prometheus pod: %v", err)
	}

	if len(pl.Items) < 1 {
		return errors.New("no Prometheus pods found")
	}

	// only use the first pod in the list
	return portForward(pl.Items[0].Name, namespace, "Prometheus",
		"http://localhost:%d", 9090, client, cmd.CobraCommand.OutOrStdout())
}
