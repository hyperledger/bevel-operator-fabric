package peer

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
	out      io.Writer
	errOut   io.Writer
	peerOpts Options
}

func (c *updateCmd) validate() error {
	return c.peerOpts.Validate()
}
func (c *updateCmd) run() error {
	hlfClient, err := helpers.GetKubeOperatorClient()
	if err != nil {
		return err
	}
	log.Infof("updating name=%s namespace=%s", c.peerOpts.Name, c.peerOpts.NS)
	ctx := context.Background()
	peer, err := hlfClient.HlfV1alpha1().FabricPeers(c.peerOpts.NS).Get(ctx, c.peerOpts.Name, v1.GetOptions{})
	if err != nil {
		return err
	}
	var hostAliases []corev1.HostAlias
	for _, hostAlias := range c.peerOpts.HostAliases {
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
	peer.Spec.HostAliases = hostAliases
	_, err = hlfClient.HlfV1alpha1().FabricPeers(c.peerOpts.NS).Update(ctx, peer, v1.UpdateOptions{})
	if err != nil {
		return err
	}
	log.Infof("updated peer %s", c.peerOpts.Name)
	return nil
}

func newUpdatePeerCMD(out io.Writer, errOut io.Writer) *cobra.Command {
	c := updateCmd{out: out, errOut: errOut}
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update a Fabric Peer",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := c.validate(); err != nil {
				return err
			}
			return c.run()
		},
	}
	f := cmd.Flags()
	f.StringVar(&c.peerOpts.Name, "name", "", "Name of the Fabric Peer to create")
	f.StringVarP(&c.peerOpts.NS, "namespace", "n", helpers.DefaultNamespace, "Namespace scope for this request")
	f.StringArrayVarP(&c.peerOpts.HostAliases, "host-aliases", "", []string{}, "Host aliases (e.g.: \"1.2.3.4:osn1.example.com,osn2.example.com\")")

	cmd.MarkPersistentFlagRequired("name")
	cmd.MarkPersistentFlagRequired("namespace")
	return cmd
}
