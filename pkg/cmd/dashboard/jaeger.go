package dashboard

import (
	"errors"
	"fmt"

	"github.com/181192/ops-cli/pkg/cmd/cmdutils"
	"github.com/181192/ops-cli/pkg/kubernetes"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func jaegerDashboardCmd(cmd *cmdutils.Cmd) {
	var opts Options

	cmd.CobraCommand.Use = "jaeger"
	cmd.CobraCommand.Short = "Open Jaeger web UI"
	cmd.CobraCommand.Long = "Open Jaeger dashboard"
	cmd.CobraCommand.RunE = func(_ *cobra.Command, args []string) error {
		cmd.NameArg = cmdutils.GetNameArg(args)
		return doPortForwardJaeger(cmd, opts)
	}

	cmd.FlagSetGroup.InFlagSet("Dashboard", func(fs *pflag.FlagSet) {
		cmdutils.AddCommonFlagsForKubernetes(fs, &opts.KubeOptions)
	})
}

func doPortForwardJaeger(cmd *cmdutils.Cmd, opts Options) error {
	kubeConfig := opts.KubeOptions.KubeConfig
	kubeContext := opts.KubeOptions.KubeContext
	namespace := opts.KubeOptions.Namespace

	client, err := kubernetes.NewClient(kubeConfig, kubeContext)
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
		"http://localhost:%d", 16686, client, cmd.CobraCommand.OutOrStdout())
}
