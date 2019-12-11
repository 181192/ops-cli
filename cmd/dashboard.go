package cmd

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

func init() {
	RootCmd.AddCommand(dashboardCmd)

	dashboardCmd.AddCommand(kialiDashCmd())
	dashboardCmd.AddCommand(promDashCmd())
	dashboardCmd.AddCommand(grafanaDashCmd())
	dashboardCmd.AddCommand(jaegerDashCmd())

	dashboardCmd.PersistentFlags().StringVarP(&kubeconfig, "kubeconfig", "c", "", "Kubernetes configuration file")
	dashboardCmd.PersistentFlags().StringVar(&configContext, "context", "", "The name of the kubeconfig context to use")
	dashboardCmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", v1.NamespaceAll, "Config namespace")

	dashboardCmd.MarkFlagRequired("label")
	dashboardCmd.Flags().StringVarP(&label, "label", "l", "", "key=value")

	dashboardCmd.MarkFlagRequired("port")
	dashboardCmd.Flags().IntVarP(&port, "port", "p", 0, "container port")
}

func grafanaDashCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "grafana",
		Short: "Open Grafana web UI",
		Long:  `Open Grafana dashboard`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := kubernetes.NewClient(kubeconfig, configContext)
			if err != nil {
				return fmt.Errorf("failed to create k8s client: %v", err)
			}

			if namespace == "" {
				namespace = "monitoring"
			}

			pl, err := client.PodsForSelector(namespace, "app=grafana")
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
				"http://localhost:%d", 3000, client, cmd.OutOrStdout())
		},
	}

	return cmd
}

func promDashCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "prometheus",
		Short: "Open Prometheus web UI",
		Long:  `Open Prometheus dashboard`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := kubernetes.NewClient(kubeconfig, configContext)
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
				"http://localhost:%d", 9090, client, cmd.OutOrStdout())
		},
	}

	return cmd
}

func kialiDashCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kiali",
		Short: "Open Kiali web UI",
		Long:  `Open Istio's Kiali dashboard`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := kubernetes.NewClient(kubeconfig, configContext)
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
				"http://localhost:%d/kiali", 20001, client, cmd.OutOrStdout())
		},
	}

	return cmd
}

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
