package networkconfig

import (
	"context"
	"github.com/kfsoftware/hlf-operator/kubectl-hlf/cmd/helpers"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"io"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
)

const (
	networkConfigRefreshDesc = `
'refresh' command deletes a Hyperledger Fabric Network Config tenant`
	networkConfigRefreshExample = `  kubectl hlf networkconfig refresh --name org1-nc --namespace default`
)

type networkConfigRefreshCmd struct {
	out    io.Writer
	errOut io.Writer
	name   string
	ns     string
}

func newRefreshNetworkConfigCmd(out io.Writer, errOut io.Writer) *cobra.Command {
	c := &networkConfigRefreshCmd{out: out, errOut: errOut}

	cmd := &cobra.Command{
		Use:     "refresh",
		Short:   "Refresh a Network Config CRD",
		Long:    networkConfigRefreshDesc,
		Example: networkConfigRefreshExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := c.validate(); err != nil {
				return err
			}
			return c.run(args)
		},
	}

	f := cmd.Flags()
	f.StringVar(&c.name, "name", "", "Name of the Network Config to refresh")
	f.StringVarP(&c.ns, "namespace", "n", helpers.DefaultNamespace, "Namespace scope for this request")
	return cmd
}

func (d *networkConfigRefreshCmd) validate() error {
	if d.name == "" {
		return errors.New("--name flag is required for CA deletion")
	}
	return nil
}

func (d *networkConfigRefreshCmd) run(args []string) error {
	oclient, err := helpers.GetKubeOperatorClient()
	if err != nil {
		return err
	}
	networkConfig, err := oclient.HlfV1alpha1().FabricNetworkConfigs(d.ns).Get(context.Background(), d.name, metav1.GetOptions{})
	if err != nil {
		return err
	}
	if networkConfig.Annotations == nil {
        networkConfig.Annotations = make(map[string]string)
    }
	networkConfig.Annotations["reloader.hlf.kungfusoftware.es/time"] = time.Now().Format(time.RFC3339)
	_, err = oclient.HlfV1alpha1().FabricNetworkConfigs(d.ns).Update(context.Background(), networkConfig, metav1.UpdateOptions{})
	if err != nil {
		return err
	}
	return nil
}
