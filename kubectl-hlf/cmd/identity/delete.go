package identity

import (
	"context"
	"fmt"
	"github.com/kfsoftware/hlf-operator/kubectl-hlf/cmd/helpers"
	"github.com/spf13/cobra"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type deleteIdentityCmd struct {
	name      string
	namespace string
}

func (c *deleteIdentityCmd) validate() error {
	if c.name == "" {
		return fmt.Errorf("--name is required")
	}
	if c.namespace == "" {
		return fmt.Errorf("--namespace is required")
	}
	return nil
}
func (c *deleteIdentityCmd) run() error {
	oclient, err := helpers.GetKubeOperatorClient()
	if err != nil {
		return err
	}
	ctx := context.Background()
	fabricIdentity, err := oclient.HlfV1alpha1().FabricIdentities(c.namespace).Get(ctx, c.name, v1.GetOptions{})
	if err != nil {
		return err
	}
	err = oclient.HlfV1alpha1().FabricIdentities(c.namespace).Delete(
		ctx,
		fabricIdentity.Name,
		v1.DeleteOptions{},
	)
	if err != nil {
		return err
	}
	fmt.Printf("Deleted identity %s\n", fabricIdentity.Name)
	return nil
}
func newIdentityDeleteCMD() *cobra.Command {
	c := &deleteIdentityCmd{}
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete HLF identity",
		Long:  `Delete HLF identity`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := c.validate(); err != nil {
				return err
			}
			return c.run()
		},
	}
	f := cmd.Flags()
	f.StringVar(&c.name, "name", "", "Name of the identity")
	f.StringVar(&c.namespace, "namespace", "", "Namespace of the identity")
	return cmd
}
