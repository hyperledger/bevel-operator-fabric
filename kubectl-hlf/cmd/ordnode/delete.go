package ordnode

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
'delete' command deletes a Hyperledger Fabric OrdererNode tenant`
	tenantDeleteExample = `  kubectl hlf ca delete --name org1-ca --namespace default`
)

type ordererNodeDeleteCmd struct {
	out    io.Writer
	errOut io.Writer
	name   string
	ns     string
}

func newOrdererNodeDeleteCmd(out io.Writer, errOut io.Writer) *cobra.Command {
	c := &ordererNodeDeleteCmd{out: out, errOut: errOut}

	cmd := &cobra.Command{
		Use:     "delete",
		Short:   "Delete an Orderer Node",
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
	f.StringVar(&c.name, "name", "", "Name of the Ordering Service to delete")
	f.StringVarP(&c.ns, "namespace", "n", helpers.DefaultNamespace, "Namespace scope for this request")
	return cmd
}

func (d *ordererNodeDeleteCmd) validate() error {
	if d.name == "" {
		return errors.New("--name flag is required for OrdererNode deletion")
	}
	return nil
}

func (d *ordererNodeDeleteCmd) run(args []string) error {
	oclient, err := helpers.GetKubeOperatorClient()
	if err != nil {
		return err
	}
	return deleteOrdererNode(oclient, d)
}

func deleteOrdererNode(client *operatorv1.Clientset, d *ordererNodeDeleteCmd) error {
	tenant, err := client.HlfV1alpha1().FabricOrdererNodes(d.ns).Get(context.Background(), d.name, metav1.GetOptions{})
	if err != nil {
		return err
	}
	if err := client.HlfV1alpha1().FabricOrdererNodes(d.ns).Delete(context.Background(), d.name, metav1.DeleteOptions{}); err != nil {
		return err
	}

	fmt.Printf("Deleting Fabric Ordering Service %s\n", tenant.ObjectMeta.Name)
	return nil
}
