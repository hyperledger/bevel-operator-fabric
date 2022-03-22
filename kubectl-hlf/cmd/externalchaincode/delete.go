package externalchaincode

import (
	"context"
	"fmt"
	"github.com/kfsoftware/hlf-operator/kubectl-hlf/cmd/helpers"
	"github.com/spf13/cobra"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type deleteExternalChaincodeCmd struct {
	name      string
	namespace string
}

func (c *deleteExternalChaincodeCmd) validate() error {
	if c.name == "" {
		return fmt.Errorf("--name is required")
	}
	if c.namespace == "" {
		return fmt.Errorf("--namespace is required")
	}
	return nil
}
func (c *deleteExternalChaincodeCmd) run() error {
	oclient, err := helpers.GetKubeOperatorClient()
	if err != nil {
		return err
	}
	ctx := context.Background()
	fabricChaincode, err := oclient.HlfV1alpha1().FabricChaincodes(c.namespace).Get(ctx, c.name, v1.GetOptions{})
	if err != nil {
		return err
	}
	err = oclient.HlfV1alpha1().FabricChaincodes(c.namespace).Delete(
		ctx,
		fabricChaincode.Name,
		v1.DeleteOptions{},
	)
	if err != nil {
		return err
	}
	fmt.Printf("Deleted external chaincode %s\n", fabricChaincode.Name)
	return nil
}
func newExternalChaincodeDeleteCmd() *cobra.Command {
	c := &deleteExternalChaincodeCmd{}
	cmd := &cobra.Command{
		Use: "delete",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := c.validate(); err != nil {
				return err
			}
			return c.run()
		},
	}
	f := cmd.Flags()
	f.StringVar(&c.name, "name", "", "Name of the external chaincode")
	f.StringVar(&c.namespace, "namespace", "", "Namespace of the external chaincode")
	return cmd
}
