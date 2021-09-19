package utils

import (
	"github.com/spf13/cobra"
	"io"
)

func NewUtilsCMD(out io.Writer, errOut io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use: "utils",
	}
	cmd.AddCommand(newAddUserCmd(out, errOut))
	return cmd
}
