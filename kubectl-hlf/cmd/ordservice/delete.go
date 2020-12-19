package ordservice

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
'delete' command deletes a Hyperledger Fabric OrderingService tenant`
	tenantDeleteExample = `  kubectl hlf ca delete --name org1-ca --namespace default`
)

type ordererDeleteCmd struct {
	out    io.Writer
	errOut io.Writer
	name   string
	ns     string
}

func newOrderingServiceDeleteCmd(out io.Writer, errOut io.Writer) *cobra.Command {
	c := &ordererDeleteCmd{out: out, errOut: errOut}

	cmd := &cobra.Command{
		Use:     "delete",
		Short:   "Delete an Ordering service",
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
	f.StringVar(&c.name, "name", "", "name of the Ordering Service to delete")
	f.StringVarP(&c.ns, "namespace", "n", helpers.DefaultNamespace, "namespace scope for this request")
	return cmd
}

func (d *ordererDeleteCmd) validate() error {
	if d.name == "" {
		return errors.New("--name flag is required for OrderingService deletion")
	}
	return nil
}

func (d *ordererDeleteCmd) run(args []string) error {
	oclient, err := helpers.GetKubeOperatorClient()
	if err != nil {
		return err
	}
	return deleteOrderingService(oclient, d)
}

func deleteOrderingService(client *operatorv1.Clientset, d *ordererDeleteCmd) error {
	tenant, err := client.HlfV1alpha1().FabricOrderingServices(d.ns).Get(context.Background(), d.name, metav1.GetOptions{})
	if err != nil {
		return err
	}
	if err := client.HlfV1alpha1().FabricOrderingServices(d.ns).Delete(context.Background(), d.name, metav1.DeleteOptions{}); err != nil {
		return err
	}

	fmt.Printf("Deleting Fabric Ordering Service %s\n", tenant.ObjectMeta.Name)
	return nil
}
