package dashboard

import (
	"errors"
	"fmt"

	"github.com/181192/ops-cli/pkg/cmd/cmdutils"
	"github.com/181192/ops-cli/pkg/kubernetes"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func alertmanagerDashboardCmd(cmd *cmdutils.Cmd) {
	var opts Options

	cmd.CobraCommand.Use = "alertmanager"
	cmd.CobraCommand.Short = "Open Alertmanager web UI"
	cmd.CobraCommand.Long = "Open Alertmanager dashboard"
	cmd.CobraCommand.RunE = func(_ *cobra.Command, args []string) error {
		cmd.NameArg = cmdutils.GetNameArg(args)
		return doPortForwardAlertmanager(cmd, opts)
	}

	cmd.FlagSetGroup.InFlagSet("Dashboard", func(fs *pflag.FlagSet) {
		cmdutils.AddCommonFlagsForKubernetes(fs, &opts.KubeOptions)
	})
}

func doPortForwardAlertmanager(cmd *cmdutils.Cmd, opts Options) error {
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

	pl, err := client.PodsForSelector(namespace, "app=alertmanager")
	if err != nil {
		return fmt.Errorf("not able to locate Alertmanager pod: %v", err)
	}

	if len(pl.Items) < 1 {
		return errors.New("no Alertmanager pods found")
	}

	// only use the first pod in the list
	return portForward(pl.Items[0].Name, namespace, "Alertmanager",
		"http://localhost:%d", 9093, client, cmd.CobraCommand.OutOrStdout())
}
