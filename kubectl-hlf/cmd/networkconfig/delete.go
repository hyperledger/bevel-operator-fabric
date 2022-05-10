package networkconfig

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
	networkConfigDeleteDesc = `
'delete' command deletes a Hyperledger Fabric Network Config tenant`
	networkConfigDeleteExample = `  kubectl hlf networkconfig delete --name org1-nc --namespace default`
)

type networkConfigDeleteCmd struct {
	out    io.Writer
	errOut io.Writer
	name   string
	ns     string
}

func newDeleteNetworkConfigCmd(out io.Writer, errOut io.Writer) *cobra.Command {
	c := &networkConfigDeleteCmd{out: out, errOut: errOut}

	cmd := &cobra.Command{
		Use:     "delete",
		Short:   "Delete a Network Config CRD",
		Long:    networkConfigDeleteDesc,
		Example: networkConfigDeleteExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := c.validate(); err != nil {
				return err
			}
			return c.run(args)
		},
	}

	f := cmd.Flags()
	f.StringVar(&c.name, "name", "", "Name of the Network Config to delete")
	f.StringVarP(&c.ns, "namespace", "n", helpers.DefaultNamespace, "Namespace scope for this request")
	return cmd
}

func (d *networkConfigDeleteCmd) validate() error {
	if d.name == "" {
		return errors.New("--name flag is required for CA deletion")
	}
	return nil
}

func (d *networkConfigDeleteCmd) run(args []string) error {
	oclient, err := helpers.GetKubeOperatorClient()
	if err != nil {
		return err
	}
	return deleteNetworkConfig(oclient, d)
}

func deleteNetworkConfig(client *operatorv1.Clientset, d *networkConfigDeleteCmd) error {
	networkConfig, err := client.HlfV1alpha1().FabricNetworkConfigs(d.ns).Get(context.Background(), d.name, metav1.GetOptions{})
	if err != nil {
		return err
	}
	if err := client.HlfV1alpha1().FabricNetworkConfigs(d.ns).Delete(context.Background(), d.name, metav1.DeleteOptions{}); err != nil {
		return err
	}

	fmt.Printf("Deleting NetworkConfig CA %s\n", networkConfig.ObjectMeta.Name)
	return nil
}
