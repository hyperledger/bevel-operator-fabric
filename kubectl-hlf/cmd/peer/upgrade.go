package peer

import (
	"context"
	"github.com/kfsoftware/hlf-operator/kubectl-hlf/cmd/helpers"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"io"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type upgradePeerCmd struct {
	name      string
	namespace string
	image     string
	version   string
}

func (c *upgradePeerCmd) validate() error {
	if c.name == "" {
		return errors.Errorf("--name is required")
	}
	if c.namespace == "" {
		return errors.Errorf("--namespace is required")
	}
	if c.image == "" {
		return errors.Errorf("--image is required")
	}
	if c.version == "" {
		return errors.Errorf("--version is required")
	}
	return nil
}
func (c *upgradePeerCmd) run() error {
	hlfClient, err := helpers.GetKubeOperatorClient()
	if err != nil {
		return err
	}
	log.Infof("upgrading name=%s namespace=%s", c.name, c.namespace)
	ctx := context.Background()
	peer, err := hlfClient.HlfV1alpha1().FabricPeers(c.namespace).Get(ctx, c.name, v1.GetOptions{})
	if err != nil {
		return err
	}
	peer.Spec.Image = c.image
	peer.Spec.Tag = c.version
	_, err = hlfClient.HlfV1alpha1().FabricPeers(c.namespace).Update(ctx, peer, v1.UpdateOptions{})
	if err != nil {
		return err
	}
	log.Infof("Upgraded peer %s", c.name)
	return nil
}

func newUpgradePeerCMD(io.Writer, io.Writer) *cobra.Command {
	c := &upgradePeerCmd{}
	cmd := &cobra.Command{
		Use: "upgrade",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := c.validate(); err != nil {
				return err
			}
			return c.run()
		},
	}
	persistentFlags := cmd.PersistentFlags()
	persistentFlags.StringVarP(&c.name, "name", "", "", "Peer Service name")
	persistentFlags.StringVarP(&c.namespace, "namespace", "", "default", "Namespace scope for this request")
	persistentFlags.StringVarP(&c.image, "image", "", helpers.DefaultPeerImage, "Version of the Fabric Peer")
	persistentFlags.StringVarP(&c.version, "version", "", helpers.DefaultPeerVersion, "Version of the Fabric Peer")

	cmd.MarkPersistentFlagRequired("name")
	cmd.MarkPersistentFlagRequired("namespace")
	cmd.MarkPersistentFlagRequired("image")
	cmd.MarkPersistentFlagRequired("version")
	return cmd
}
