package ssh

import (
	"fmt"

	cmdUtil "github.com/181192/ops-cli/pkg/util"

	"github.com/spf13/cobra"
)

// nodeSSHCmd represents the node ssh command
func nodeSSHCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "node",
		Short: "ssh into a Kubernetes node",
		Long:  `Enter into any Kubernetes node by creating a new pod tolerated to the specified node.`,
		Run: func(cmd *cobra.Command, args []string) {
			nodeName := args[0]

			podName := "k-nsenter-" + nodeName

			spec := fmt.Sprintf(`{
				"spec": {
					"hostPID": true,
					"hostNetwork": true,
					"nodeSelector": {
						"kubernetes.io/hostname": %q,
					},
					"tolerations": [{
							"operator": "Exists"
					}],
					"containers": [
						{
							"name": "nsenter",
							"image": "alexeiled/nsenter:2.34",
							"command": [
								"/nsenter", "--all", "--target=1", "--", "su", "-"
							],
							"stdin": true,
							"tty": true,
							"securityContext": {
								"privileged": true
							},
							"resources": {
								"requests": {
									"cpu": "10m"
								}
							}
						}
					]
				}
			}`, nodeName)

			cmdUtil.ExecuteCmd(cmd, "kubectl", []string{"run", podName, "--restart=Never", "-it", "--rm", "--image overriden", "--overrides " + spec})
		},
	}

	return cmd
}
