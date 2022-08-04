package operatorui

import (
	"github.com/spf13/cobra"
	"io"
)

func NewOperatorUICMD(out io.Writer, errOut io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use: "operatorui",
	}
	cmd.AddCommand(
		newCreateOperatorUICmd(out, errOut),
		newDeleteOperatorUICmd(out, errOut),
		newUpdateOperatorUICmd(out, errOut),
	)

	return cmd
}
