package ordnode

import (
	"context"
	"github.com/kfsoftware/hlf-operator/kubectl-hlf/cmd/helpers"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"io"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type upgradeOrdererCmd struct {
	name      string
	namespace string
	image     string
	version   string
}

func (c *upgradeOrdererCmd) validate() error {
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
func (c *upgradeOrdererCmd) run() error {
	hlfClient, err := helpers.GetKubeOperatorClient()
	if err != nil {
		return err
	}
	log.Infof("name=%s namespace=%s", c.name, c.namespace)
	ctx := context.Background()
	ordererNode, err := hlfClient.HlfV1alpha1().FabricOrdererNodes(c.namespace).Get(ctx, c.name, v1.GetOptions{})
	if err != nil {
		return err
	}
	ordererNode.Spec.Image = c.image
	ordererNode.Spec.Tag = c.version
	_, err = hlfClient.HlfV1alpha1().FabricOrdererNodes(c.namespace).Update(ctx, ordererNode, v1.UpdateOptions{})
	if err != nil {
		return err
	}
	log.Infof("Upgraded orderer node %s", c.name)
	return nil
}

func newUpgradeOrdererCMD(io.Writer, io.Writer) *cobra.Command {
	c := &upgradeOrdererCmd{}
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
	persistentFlags.StringVarP(&c.name, "name", "", "", "Orderer Service name")
	persistentFlags.StringVarP(&c.namespace, "namespace", "", "default", "Namespace scope for this request")
	persistentFlags.StringVarP(&c.image, "image", "", helpers.DefaultOrdererImage, "Version of the Fabric Orderer Node")
	persistentFlags.StringVarP(&c.version, "version", "", helpers.DefaultOrdererVersion, "Version of the Fabric Orderer Node")
	cmd.MarkPersistentFlagRequired("name")
	cmd.MarkPersistentFlagRequired("namespace")
	cmd.MarkPersistentFlagRequired("image")
	cmd.MarkPersistentFlagRequired("version")
	return cmd
}
