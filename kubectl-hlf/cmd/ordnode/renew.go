package ordnode

import (
	"context"
	"github.com/kfsoftware/hlf-operator/kubectl-hlf/cmd/helpers"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"io"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
)

type renewChannelCmd struct {
	name      string
	namespace string
}

func (c *renewChannelCmd) validate() error {
	if c.namespace == "" {
		return errors.Errorf("--namespace is required")
	}
	if c.name == "" {
		return errors.Errorf("--name is required")
	}
	return nil
}
func (c *renewChannelCmd) run() error {
	hlfClient, err := helpers.GetKubeOperatorClient()
	if err != nil {
		return err
	}
	log.Infof("name=%s namespace=%s", c.name, c.namespace)
	now := v1.NewTime(time.Now())
	ctx := context.Background()
	ordererNode, err := hlfClient.HlfV1alpha1().FabricOrdererNodes(c.namespace).Get(ctx, c.name, v1.GetOptions{})
	if err != nil {
		return err
	}
	ordererNode.Spec.UpdateCertificateTime = &now
	_, err = hlfClient.HlfV1alpha1().FabricOrdererNodes(c.namespace).Update(ctx, ordererNode, v1.UpdateOptions{})
	if err != nil {
		return err
	}
	log.Infof("Renewed certificate for orderer node %s", c.name)
	return nil
}

func newRenewChannelCMD(io.Writer, io.Writer) *cobra.Command {
	c := &renewChannelCmd{}
	cmd := &cobra.Command{
		Use: "renew",
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
	cmd.MarkPersistentFlagRequired("name")
	cmd.MarkPersistentFlagRequired("namespace")
	return cmd
}
