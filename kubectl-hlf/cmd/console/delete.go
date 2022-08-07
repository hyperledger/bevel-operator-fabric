package console

import (
	"context"
	"fmt"
	"github.com/kfsoftware/hlf-operator/kubectl-hlf/cmd/helpers"
	"github.com/spf13/cobra"
	"io"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type deleteConsoleCmd struct {
	name      string
	namespace string
}

func (c *deleteConsoleCmd) validate() error {
	if c.name == "" {
		return fmt.Errorf("--name is required")
	}
	if c.namespace == "" {
		return fmt.Errorf("--namespace is required")
	}
	return nil
}
func (c *deleteConsoleCmd) run() error {
	oclient, err := helpers.GetKubeOperatorClient()
	if err != nil {
		return err
	}
	ctx := context.Background()
	fabricOperationsConsole, err := oclient.HlfV1alpha1().FabricOperationsConsoles(c.namespace).Get(ctx, c.name, v1.GetOptions{})
	if err != nil {
		return err
	}
	err = oclient.HlfV1alpha1().FabricOperationsConsoles(c.namespace).Delete(
		ctx,
		fabricOperationsConsole.Name,
		v1.DeleteOptions{},
	)
	if err != nil {
		return err
	}
	fmt.Printf("Deleted operations console %s\n", fabricOperationsConsole.Name)
	return nil
}
func newDeleteConsoleCmd(out io.Writer, errOut io.Writer) *cobra.Command {
	c := &deleteConsoleCmd{}
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
	f.StringVar(&c.name, "name", "", "Name of the operations console")
	f.StringVar(&c.namespace, "namespace", "", "Namespace of the operations console")
	return cmd
}
