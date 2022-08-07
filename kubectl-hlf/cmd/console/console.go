package console

import (
	"github.com/spf13/cobra"
	"io"
)

func NewConsoleCmd(out io.Writer, errOut io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use: "console",
	}
	cmd.AddCommand(
		newCreateConsoleCmd(out, errOut),
		newDeleteConsoleCmd(out, errOut),
	)

	return cmd
}
