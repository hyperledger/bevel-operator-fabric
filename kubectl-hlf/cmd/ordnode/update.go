package ordnode

import (
	"context"
	"io"
	"strings"

	"github.com/kfsoftware/hlf-operator/kubectl-hlf/cmd/helpers"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type updateCmd struct {
	out         io.Writer
	errOut      io.Writer
	ordererOpts OrdererOptions
}

func (c *updateCmd) validate() error {
	return c.ordererOpts.Validate()
}
func (c *updateCmd) run() error {
	hlfClient, err := helpers.GetKubeOperatorClient()
	if err != nil {
		return err
	}
	log.Infof("updating name=%s namespace=%s", c.ordererOpts.Name, c.ordererOpts.NS)
	ctx := context.Background()
	ordnode, err := hlfClient.HlfV1alpha1().FabricOrdererNodes(c.ordererOpts.NS).Get(ctx, c.ordererOpts.Name, v1.GetOptions{})
	if err != nil {
		return err
	}
	var hostAliases []corev1.HostAlias
	for _, hostAlias := range c.ordererOpts.HostAliases {
		ipAndNames := strings.Split(hostAlias, ":")
		if len(ipAndNames) == 2 {
			aliases := strings.Split(ipAndNames[1], ",")
			if len(aliases) > 0 {
				hostAliases = append(hostAliases, corev1.HostAlias{IP: ipAndNames[0], Hostnames: aliases})
			} else {
				log.Warningf("ingnoring host-alias [%s]: must be in format <ip>:<alias1>,<alias2>...", hostAlias)
			}
		} else {
			log.Warningf("ingnoring host-alias [%s]: must be in format <ip>:<alias1>,<alias2>...", hostAlias)
		}
	}
	ordnode.Spec.HostAliases = hostAliases
	_, err = hlfClient.HlfV1alpha1().FabricOrdererNodes(c.ordererOpts.NS).Update(ctx, ordnode, v1.UpdateOptions{})
	if err != nil {
		return err
	}
	log.Infof("updated orderer node %s", c.ordererOpts.Name)
	return nil
}

func newUpdateOrdererCMD(out io.Writer, errOut io.Writer) *cobra.Command {
	c := updateCmd{out: out, errOut: errOut}
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update a Fabric orderer node",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := c.validate(); err != nil {
				return err
			}
			return c.run()
		},
	}
	f := cmd.Flags()
	f.StringVar(&c.ordererOpts.Name, "name", "", "Name of the Fabric orderer node to create")
	f.StringVarP(&c.ordererOpts.NS, "namespace", "n", helpers.DefaultNamespace, "Namespace scope for this request")
	f.StringArrayVarP(&c.ordererOpts.HostAliases, "host-aliases", "", []string{}, "Host aliases (e.g.: \"1.2.3.4:osn2.example.com,peer1.example.com\")")

	cmd.MarkPersistentFlagRequired("name")
	cmd.MarkPersistentFlagRequired("namespace")
	return cmd
}
