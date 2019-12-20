package cmd

import (
	"fmt"

	"github.com/181192/ops-cli/pkg/kubernetes"
	"github.com/ghodss/yaml"
	"github.com/spf13/cobra"

	appsv1beta1 "k8s.io/api/apps/v1beta1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

var (
	labelSelector = ""
	applyManifest bool
	kubeconfig    string
	configContext string
	namespace     string
	label         string
	port          int
)

// manifestCmd represents the manifest command
var manifestCmd = &cobra.Command{
	Use:   "manifest",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("manifest called")
	},
}

func init() {
	rootCmd.AddCommand(manifestCmd)

	manifestCmd.AddCommand(manifestNamespacedTillerCmd())
	manifestCmd.PersistentFlags().StringVarP(&kubeconfig, "kubeconfig", "c", "", "Kubernetes configuration file")
	manifestCmd.PersistentFlags().StringVar(&configContext, "context", "", "The name of the kubeconfig context to use")
	manifestCmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", corev1.NamespaceAll, "Config namespace")
	manifestCmd.PersistentFlags().BoolVar(&applyManifest, "apply", false, "Apply kubernetes manifests to cluster")

}

func manifestNamespacedTillerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "namespaced-tiller",
		Short: "Generate manifests for namespaced-tiller",
		Long:  `Generate manifests for a given namespace, configure tiller and rbac rules.`,
		RunE: func(cmd *cobra.Command, args []string) error {

			ns := createNamespace(namespace)
			tiller := createTillerDeployment(namespace)
			nsYaml, err := yaml.Marshal(ns)
			if err != nil {
				return fmt.Errorf("failed create yaml %v", err)
			}

			tillerDeploymentYaml, err := yaml.Marshal(tiller.deployment)
			if err != nil {
				return fmt.Errorf("failed create yaml %v", err)
			}

			tillerServiceYaml, err := yaml.Marshal(tiller.service)
			if err != nil {
				return fmt.Errorf("failed create yaml %v", err)
			}

			tillerServiceAccountYaml, err := yaml.Marshal(tiller.serviceAccount)
			if err != nil {
				return fmt.Errorf("failed create yaml %v", err)
			}

			tillerClusterRoleBindingYaml, err := yaml.Marshal(tiller.clusterRoleBinding)
			if err != nil {
				return fmt.Errorf("failed create yaml %v", err)
			}
			fmt.Println("---")
			fmt.Println(string(nsYaml))
			fmt.Println("---")
			fmt.Println(string(tillerDeploymentYaml))
			fmt.Println("---")
			fmt.Println(string(tillerServiceYaml))
			fmt.Println("---")
			fmt.Println(string(tillerServiceAccountYaml))
			fmt.Println("---")
			fmt.Println(string(tillerClusterRoleBindingYaml))
			fmt.Println("---")

			if applyManifest {
				client, err := kubernetes.NewClient(kubeconfig, configContext)
				if err != nil {
					return fmt.Errorf("failed to create k8s client: %v", err)
				}

				clientset, err := kubernetes.NewClientset(client)
				if err != nil {
					return fmt.Errorf("failed to create k8s clientset: %v", err)
				}

				_, err = clientset.CoreV1().Namespaces().Update(ns)
				if err != nil {
					return fmt.Errorf("failed to create namespace %v", err)
				}

				_, err = clientset.CoreV1().Services(namespace).Update(tiller.service)
				if err != nil {
					return fmt.Errorf("failed to create service %v", err)
				}

				_, err = clientset.CoreV1().ServiceAccounts(namespace).Update(tiller.serviceAccount)
				if err != nil {
					return fmt.Errorf("failed to create serviceaccount %v", err)
				}

				_, err = clientset.AppsV1beta1().Deployments(namespace).Update(tiller.deployment)
				if err != nil {
					return fmt.Errorf("failed to create deployment %v", err)
				}

				_, err = clientset.RbacV1().ClusterRoleBindings().Update(tiller.clusterRoleBinding)
				if err != nil {
					return fmt.Errorf("failed to create clusterrolebinding %v", err)
				}
			}

			return nil
		},
	}

	return cmd
}

func createNamespace(namespace string) *corev1.Namespace {
	return &corev1.Namespace{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Namespace",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: namespace,
		},
		Spec: corev1.NamespaceSpec{},
	}
}

// TillerDeployment represent a tiller deployment
type TillerDeployment struct {
	deployment         *appsv1beta1.Deployment
	service            *corev1.Service
	serviceAccount     *corev1.ServiceAccount
	clusterRoleBinding *rbacv1.ClusterRoleBinding
}

func createTillerDeployment(namespace string) *TillerDeployment {
	deployment := &appsv1beta1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "extensions/v1beta1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "tiller-deploy",
			Namespace: namespace,
			Labels: map[string]string{
				"app":     "helm",
				"name":    "tiller",
				"version": "v2.15.1",
			},
		},
		Spec: appsv1beta1.DeploymentSpec{
			Replicas: ptrint32(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app":     "helm",
					"name":    "tiller",
					"version": "v2.15.1",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app":     "helm",
						"name":    "tiller",
						"version": "v2.15.1",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						corev1.Container{
							Name:  "tiller",
							Image: "gcr.io/kubernetes-helm/tiller:v2.15.1",
							Ports: []corev1.ContainerPort{
								corev1.ContainerPort{
									Name:          "tiller",
									ContainerPort: 44134,
									Protocol:      corev1.Protocol("TCP"),
								},
								corev1.ContainerPort{
									Name:          "http",
									ContainerPort: 44135,
									Protocol:      corev1.Protocol("TCP"),
								},
							},
							Env: []corev1.EnvVar{
								corev1.EnvVar{
									Name: "TILLER_NAMESPACE",
									ValueFrom: &corev1.EnvVarSource{
										FieldRef: &corev1.ObjectFieldSelector{
											FieldPath: "metadata.namespace",
										},
									},
								},
								corev1.EnvVar{
									Name:  "TILLER_HISTORY_MAX",
									Value: "10",
								},
							},
							Resources: corev1.ResourceRequirements{
								Limits: corev1.ResourceList{
									"cpu":    *resource.NewQuantity(55, resource.DecimalSI),
									"memory": *resource.NewQuantity(157286400, resource.BinarySI),
								},
								Requests: corev1.ResourceList{
									"cpu":    *resource.NewQuantity(45, resource.DecimalSI),
									"memory": *resource.NewQuantity(157286400, resource.BinarySI),
								},
							},
							LivenessProbe: &corev1.Probe{
								Handler: corev1.Handler{
									HTTPGet: &corev1.HTTPGetAction{
										Path: "/liveness",
										Port: intstr.IntOrString{
											Type:   intstr.Type(0),
											IntVal: 44135,
										},
										Scheme: corev1.URIScheme("HTTP"),
									},
								},
								InitialDelaySeconds: 1,
								TimeoutSeconds:      1,
								PeriodSeconds:       10,
								SuccessThreshold:    1,
								FailureThreshold:    3,
							},
							ReadinessProbe: &corev1.Probe{
								Handler: corev1.Handler{
									HTTPGet: &corev1.HTTPGetAction{
										Path: "/readiness",
										Port: intstr.IntOrString{
											Type:   intstr.Type(0),
											IntVal: 44135,
										},
										Scheme: corev1.URIScheme("HTTP"),
									},
								},
								InitialDelaySeconds: 1,
								TimeoutSeconds:      1,
								PeriodSeconds:       10,
								SuccessThreshold:    1,
								FailureThreshold:    3,
							},
							TerminationMessagePath:   "/dev/termination-log",
							TerminationMessagePolicy: corev1.TerminationMessagePolicy("File"),
							ImagePullPolicy:          corev1.PullPolicy("IfNotPresent"),
						},
					},
					TerminationGracePeriodSeconds: ptrint64(30),
					ServiceAccountName:            "tiller",
					DeprecatedServiceAccount:      "tiller",
					AutomountServiceAccountToken:  ptrbool(true),
				},
			},
			Strategy: appsv1beta1.DeploymentStrategy{
				Type: appsv1beta1.DeploymentStrategyType("RollingUpdate"),
				RollingUpdate: &appsv1beta1.RollingUpdateDeployment{
					MaxUnavailable: &intstr.IntOrString{
						Type:   intstr.Type(1),
						IntVal: 0,
						StrVal: "25%",
					},
					MaxSurge: &intstr.IntOrString{
						Type:   intstr.Type(1),
						IntVal: 0,
						StrVal: "25%",
					},
				},
			},
			MinReadySeconds:         0,
			RevisionHistoryLimit:    ptrint32(10),
			ProgressDeadlineSeconds: ptrint32(600),
		},
	}

	service := &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "tiller-deploy",
			Namespace: namespace,
			Labels: map[string]string{
				"app":     "helm",
				"name":    "tiller",
				"version": "v2.15.1",
			},
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				corev1.ServicePort{
					Name:     "tiller",
					Protocol: corev1.Protocol("TCP"),
					Port:     44134,
					TargetPort: intstr.IntOrString{
						Type:   intstr.Type(1),
						StrVal: "tiller",
					},
				},
			},
			Selector: map[string]string{
				"app":     "helm",
				"name":    "tiller",
				"version": "v2.15.1",
			},
		},
	}

	serviceAccount := &corev1.ServiceAccount{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ServiceAccount",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "tiller",
			Namespace: namespace,
		},
	}

	clusterRoleBinding := &rbacv1.ClusterRoleBinding{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ClusterRoleBinding",
			APIVersion: "rbac.authorization.k8s.io/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "tiller-" + namespace,
		},
		Subjects: []rbacv1.Subject{
			rbacv1.Subject{
				Kind:      "ServiceAccount",
				Name:      "tiller",
				Namespace: namespace,
			},
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "ClusterRole",
			Name:     "cluster-admin",
		},
	}

	return &TillerDeployment{
		deployment:         deployment,
		service:            service,
		serviceAccount:     serviceAccount,
		clusterRoleBinding: clusterRoleBinding,
	}
}

func ptrint32(p int32) *int32 {
	return &p
}

func ptrint64(p int64) *int64 {
	return &p
}

func ptrbool(p bool) *bool {
	return &p
}
