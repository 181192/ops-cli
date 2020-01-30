package dashboard

import (
	"errors"
	"fmt"

	"github.com/181192/ops-cli/pkg/cmd/cmdutils"
	"github.com/181192/ops-cli/pkg/kubernetes"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func kialiDashboardCmd(cmd *cmdutils.Cmd) {
	var opts Options

	cmd.CobraCommand.Use = "kiali"
	cmd.CobraCommand.Short = "Open Kiali web UI"
	cmd.CobraCommand.Long = "Open Istio's Kiali dashboard"
	cmd.CobraCommand.RunE = func(_ *cobra.Command, args []string) error {
		cmd.NameArg = cmdutils.GetNameArg(args)
		return doPortForwardKiali(cmd, opts)
	}

	cmd.FlagSetGroup.InFlagSet("Dashboard", func(fs *pflag.FlagSet) {
		// fs.StringVarP(&opts.Label, "label", "l", "", "Label selector (key=value)")
		// fs.StringVarP(&opts.Port, "port", "", "", "Container port to forward in")
		// _ = cobra.MarkFlagRequired(fs, "label")
		// _ = cobra.MarkFlagRequired(fs, "port")
		cmdutils.AddCommonFlagsForKubernetes(fs, &opts.KubeOptions)
	})
}

func doPortForwardKiali(cmd *cmdutils.Cmd, opts Options) error {
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

	pl, err := client.PodsForSelector(namespace, "app=kiali")
	if err != nil {
		return fmt.Errorf("not able to locate Kiali pod: %v", err)
	}

	if len(pl.Items) < 1 {
		return errors.New("no Kiali pods found")
	}

	// only use the first pod in the list
	return portForward(pl.Items[0].Name, namespace, "Kiali",
		"http://localhost:%d/kiali", 20001, client, cmd.CobraCommand.OutOrStdout())
}
