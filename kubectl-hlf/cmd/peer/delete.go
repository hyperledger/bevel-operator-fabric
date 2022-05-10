package peer

import (
	"context"
	"fmt"
	"github.com/kfsoftware/hlf-operator/kubectl-hlf/cmd/helpers"
	operatorv1 "github.com/kfsoftware/hlf-operator/pkg/client/clientset/versioned"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"io"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	tenantDeleteDesc = `
'delete' command deletes a Hyperledger Fabric Peer tenant`
	tenantDeleteExample = `  kubectl hlf ca delete --name org1-ca --namespace default`
)

type peerDeleteCmd struct {
	out    io.Writer
	errOut io.Writer
	name   string
	ns     string
}

func newPeerDeleteCmd(out io.Writer, errOut io.Writer) *cobra.Command {
	c := &peerDeleteCmd{out: out, errOut: errOut}

	cmd := &cobra.Command{
		Use:     "delete",
		Short:   "Delete a Peer authority",
		Long:    tenantDeleteDesc,
		Example: tenantDeleteExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := c.validate(); err != nil {
				return err
			}
			return c.run(args)
		},
	}

	f := cmd.Flags()
	f.StringVar(&c.name, "name", "", "Name of the Peer to delete")
	f.StringVarP(&c.ns, "namespace", "n", helpers.DefaultNamespace, "Namespace scope for this request")
	return cmd
}

func (d *peerDeleteCmd) validate() error {
	if d.name == "" {
		return errors.New("--name flag is required for Peer deletion")
	}
	return nil
}

func (d *peerDeleteCmd) run(args []string) error {
	oclient, err := helpers.GetKubeOperatorClient()
	if err != nil {
		return err
	}
	return deletePeer(oclient, d)
}

func deletePeer(client *operatorv1.Clientset, d *peerDeleteCmd) error {
	tenant, err := client.HlfV1alpha1().FabricPeers(d.ns).Get(context.Background(), d.name, metav1.GetOptions{})
	if err != nil {
		return err
	}
	if err := client.HlfV1alpha1().FabricPeers(d.ns).Delete(context.Background(), d.name, metav1.DeleteOptions{}); err != nil {
		return err
	}

	fmt.Printf("Deleting Fabric Peer %s\n", tenant.ObjectMeta.Name)
	return nil
}
