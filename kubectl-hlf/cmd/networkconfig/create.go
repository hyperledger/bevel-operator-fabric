package networkconfig

import (
	"context"
	"fmt"
	hlfv1alpha1 "github.com/kfsoftware/hlf-operator/api/hlf.kungfusoftware.es/v1alpha1"
	"github.com/kfsoftware/hlf-operator/kubectl-hlf/cmd/helpers"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"io"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type CreateOptions struct {
	Orgs       []string
	OutputPath string
	NS         string
	Name       string
	Internal   bool
	SecretName string
}

func (o CreateOptions) Validate() error {
	return nil
}

type createCmd struct {
	out    io.Writer
	errOut io.Writer
	opts   CreateOptions
}

func (c *createCmd) validate() error {
	return c.opts.Validate()
}
func (c *createCmd) run(args []string) error {
	oclient, err := helpers.GetKubeOperatorClient()
	if err != nil {
		return err
	}
	ctx := context.Background()
	secretName := fmt.Sprintf("%s-networkconfig", c.opts.Name)
	if c.opts.SecretName != "" {
		secretName = c.opts.SecretName
	}
	networkConfig := &hlfv1alpha1.FabricNetworkConfig{
		TypeMeta: v1.TypeMeta{
			Kind:       "FabricNetworkConfig",
			APIVersion: hlfv1alpha1.GroupVersion.String(),
		},
		ObjectMeta: v1.ObjectMeta{
			Name:      c.opts.Name,
			Namespace: c.opts.NS,
		},
		Spec: hlfv1alpha1.FabricNetworkConfigSpec{
			Organization:  "",
			Internal:      c.opts.Internal,
			Organizations: c.opts.Orgs,
			SecretName:    secretName,
		},
	}
	_, err = oclient.HlfV1alpha1().FabricNetworkConfigs(c.opts.NS).Create(
		ctx,
		networkConfig,
		v1.CreateOptions{},
	)
	if err != nil {
		return err
	}
	log.Infof("Certificate authority %s created on namespace %s", networkConfig.Name, networkConfig.Namespace)

	return nil
}

func newCreateNetworkConfigCmd(out io.Writer, errOut io.Writer) *cobra.Command {
	c := createCmd{out: out, errOut: errOut}
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates a network config",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := c.validate(); err != nil {
				return err
			}
			return c.run(args)
		},
	}
	f := cmd.Flags()
	f.StringSliceVarP(&c.opts.Orgs, "orgs", "o", []string{}, "Organizations to inspect")
	f.StringVar(&c.opts.Name, "name", "", "Name of the Network Config to create")
	f.StringVar(&c.opts.SecretName, "secret", "", "Secret name to store the network config")
	f.StringVarP(&c.opts.NS, "namespace", "n", helpers.DefaultNamespace, "Namespace scope for this request")
	f.BoolVarP(&c.opts.Internal, "internal", "i", false, "Use internal or external endpoints")
	f.StringVarP(&c.opts.OutputPath, "output-path", "", "", "Output path")

	return cmd
}
