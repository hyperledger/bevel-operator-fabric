package operatorapi

import (
	"github.com/spf13/cobra"
	"io"
)

func NewOperatorAPICMD(out io.Writer, errOut io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use: "operatorapi",
	}
	cmd.AddCommand(
		newCreateOperatorAPICmd(out, errOut),
		newDeleteOperatorAPICmd(out, errOut),
		newUpdateOperatorAPICmd(out, errOut),
	)

	return cmd
}
