package networkconfig

import (
	"context"
	"fmt"
	"github.com/kfsoftware/hlf-operator/kubectl-hlf/cmd/helpers"
	hlfv1alpha1 "github.com/kfsoftware/hlf-operator/pkg/apis/hlf.kungfusoftware.es/v1alpha1"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"io"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
)

type UpdateOptions struct {
	Orgs       []string
	OutputPath string
	NS         string
	SecretName string
	Name       string
	Internal   bool
	Identities []string
}

func (o UpdateOptions) Validate() error {
	return nil
}

type updateCmd struct {
	out    io.Writer
	errOut io.Writer
	opts   UpdateOptions
}

func (c *updateCmd) validate() error {
	return c.opts.Validate()
}
func (c *updateCmd) run(args []string) error {
	oclient, err := helpers.GetKubeOperatorClient()
	if err != nil {
		return err
	}
	networkConfig, err := oclient.HlfV1alpha1().FabricNetworkConfigs(c.opts.NS).Get(context.Background(), c.opts.Name, v1.GetOptions{})
	if err != nil {
		return err
	}
	ctx := context.Background()
	identities := []hlfv1alpha1.FabricNetworkConfigIdentity{}
	for _, identity := range c.opts.Identities {
		chunks := strings.Split(identity, ".")
		if len(chunks) != 2 {
			return fmt.Errorf("identity %s is not valid, must be in format <name>.<ns>", identity)
		}
		name := chunks[0]
		ns := chunks[1]
		_, err = oclient.HlfV1alpha1().FabricIdentities(ns).Get(
			ctx,
			name,
			v1.GetOptions{},
		)
		if err != nil {
			return errors.Wrapf(err, "error getting identity %s on namespace %s", name, ns)
		}
		identities = append(identities, hlfv1alpha1.FabricNetworkConfigIdentity{
			Name:      name,
			Namespace: ns,
		})
	}
	networkConfig.Spec.Identities = identities
	networkConfig.Spec.Internal = c.opts.Internal
	networkConfig.Spec.Organizations = c.opts.Orgs
	secretName := fmt.Sprintf("%s-networkconfig", c.opts.Name)
	if c.opts.SecretName != "" {
		secretName = c.opts.SecretName
	}
	networkConfig.Spec.SecretName = secretName
	_, err = oclient.HlfV1alpha1().FabricNetworkConfigs(c.opts.NS).Update(context.Background(), networkConfig, v1.UpdateOptions{})
	if err != nil {
		return err
	}
	return nil
}

func newUpdateNetworkConfigCmd(out io.Writer, errOut io.Writer) *cobra.Command {
	c := updateCmd{out: out, errOut: errOut}
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Updates a network config",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := c.validate(); err != nil {
				return err
			}
			return c.run(args)
		},
	}
	f := cmd.Flags()
	f.StringSliceVarP(&c.opts.Orgs, "orgs", "o", []string{}, "Organizations to inspect")
	f.StringVar(&c.opts.Name, "name", "", "Name of the Network Config to update")
	f.StringVar(&c.opts.SecretName, "secret", "", "Secret name to store the network config")
	f.StringVarP(&c.opts.NS, "namespace", "n", helpers.DefaultNamespace, "Namespace scope for this request")
	f.BoolVarP(&c.opts.Internal, "internal", "i", false, "Use internal or external endpoints")
	f.StringVarP(&c.opts.OutputPath, "output-path", "", "", "Output path")
	f.StringSliceVarP(&c.opts.Identities, "identities", "", []string{}, "Identities to use")
	return cmd
}
