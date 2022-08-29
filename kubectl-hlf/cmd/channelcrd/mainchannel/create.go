package mainchannel

import (
	"context"
	"fmt"
	"github.com/kfsoftware/hlf-operator/api/hlf.kungfusoftware.es/v1alpha1"
	"github.com/kfsoftware/hlf-operator/kubectl-hlf/cmd/helpers"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"io"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Options struct {
	Name                 string
	StorageClass         string
	Capacity             string
	NS                   string
	Image                string
	Version              string
	IngressGateway       string
	IngressPort          int
	Hosts                []string
	Output               bool
	InitialAdminPassword string
	InitialAdmin         string
	HostURL              string
	TLSSecretName        string
}

func (o Options) Validate() error {
	return nil
}

type createCmd struct {
	out         io.Writer
	errOut      io.Writer
	channelOpts Options
}

func (c *createCmd) validate() error {
	return c.channelOpts.Validate()
}
func (c *createCmd) run() error {
	oclient, err := helpers.GetKubeOperatorClient()
	if err != nil {
		return err
	}
	fabricConsole := &v1alpha1.FabricMainChannel{
		TypeMeta: v1.TypeMeta{
			Kind:       "FabricMainChannel",
			APIVersion: v1alpha1.GroupVersion.String(),
		},
		ObjectMeta: v1.ObjectMeta{
			Name:      c.channelOpts.Name,
			Namespace: c.channelOpts.NS,
		},
		Spec: v1alpha1.FabricMainChannelSpec{
			Name:                         "",
			Identities:                   nil,
			AdminPeerOrganizations:       nil,
			PeerOrganizations:            nil,
			ExternalPeerOrganizations:    nil,
			ChannelConfig:                nil,
			AdminOrdererOrganizations:    nil,
			OrdererOrganizations:         nil,
			ExternalOrdererOrganizations: nil,
			Consenters:                   nil,
		},
	}
	if c.channelOpts.Output {
		ot, err := helpers.MarshallWithoutStatus(&fabricConsole)
		if err != nil {
			return err
		}
		fmt.Println(string(ot))
	} else {
		ctx := context.Background()
		_, err = oclient.HlfV1alpha1().FabricMainChannels(c.channelOpts.NS).Create(
			ctx,
			fabricConsole,
			v1.CreateOptions{},
		)
		if err != nil {
			return err
		}
		log.Infof("Console %s created on namespace %s", fabricConsole.Name, fabricConsole.Namespace)
	}
	return nil
}
func newCreateMainChannelCmd(out io.Writer, errOut io.Writer) *cobra.Command {
	c := createCmd{out: out, errOut: errOut}
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a main channel",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := c.validate(); err != nil {
				return err
			}
			return c.run()
		},
	}
	f := cmd.Flags()
	f.StringVar(&c.channelOpts.Name, "name", "", "Name of the Fabric Console to create")
	f.BoolVarP(&c.channelOpts.Output, "output", "o", false, "Output in yaml")
	return cmd
}
