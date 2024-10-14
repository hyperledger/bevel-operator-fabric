package networkconfig

import (
	"context"
	"fmt"
	"github.com/kfsoftware/hlf-operator/kubectl-hlf/cmd/helpers"
	hlfv1alpha1 "github.com/kfsoftware/hlf-operator/pkg/apis/hlf.kungfusoftware.es/v1alpha1"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"io"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
)

type CreateOptions struct {
	Orgs       []string
	OutputPath string
	NS         string
	Identities []string
	Name       string
	Internal   bool
	SecretName string
	Channels   []string
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
	namespaces := []string{}
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
			Organization:     "",
			Internal:         c.opts.Internal,
			Organizations:    c.opts.Orgs,
			Namespaces:       namespaces,
			Channels:         c.opts.Channels,
			Identities:       identities,
			ExternalOrderers: []hlfv1alpha1.FabricNetworkConfigExternalOrderer{},
			ExternalPeers:    []hlfv1alpha1.FabricNetworkConfigExternalPeer{},
			SecretName:       secretName,
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
	log.Infof("Network Config %s created on namespace %s", networkConfig.Name, networkConfig.Namespace)

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
	f.StringSliceVarP(&c.opts.Channels, "channels", "c", []string{}, "Channels to inspect")
	f.StringVar(&c.opts.Name, "name", "", "Name of the Network Config to create")
	f.StringVar(&c.opts.SecretName, "secret", "", "Secret name to store the network config")
	f.StringVarP(&c.opts.NS, "namespace", "n", helpers.DefaultNamespace, "Namespace scope for this request")
	f.BoolVarP(&c.opts.Internal, "internal", "i", false, "Use internal or external endpoints")
	f.StringVarP(&c.opts.OutputPath, "output-path", "", "", "Output path")
	f.StringSliceVarP(&c.opts.Identities, "identities", "", []string{}, "Identities to add to the network config")

	return cmd
}
