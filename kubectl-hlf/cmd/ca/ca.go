package ca

import (
	"github.com/spf13/cobra"
	"io"
)

func NewCACmd(out io.Writer, errOut io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use: "ca",
	}
	cmd.AddCommand(newCreateCACmd(out, errOut))
	cmd.AddCommand(newCADeleteCmd(out, errOut))
	cmd.AddCommand(newCARegisterCmd(out, errOut))
	cmd.AddCommand(newCAEnrollCmd(out, errOut))
	return cmd
}
