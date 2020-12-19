package ordservice

import (
	"github.com/spf13/cobra"
	"io"
)

// NewOrdServiceCmd creates a new root command to manage Ordering Services
func NewOrdServiceCmd(out io.Writer, errOut io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use: "ordservice",
	}
	cmd.AddCommand(newCreateOrderingServiceCmd(out, errOut))
	cmd.AddCommand(newOrderingServiceDeleteCmd(out, errOut))
	return cmd
}
