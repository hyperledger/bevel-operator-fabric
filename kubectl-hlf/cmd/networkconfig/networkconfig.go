package networkconfig

import (
	"github.com/spf13/cobra"
	"io"
)

func NewNetworkConfigCmd(out io.Writer, errOut io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use: "networkconfig",
	}
	cmd.AddCommand(
		newCreateNetworkConfigCmd(out, errOut),
		newDeleteNetworkConfigCmd(out, errOut),
		newExportNetworkConfigCmd(out, errOut),
		newRefreshNetworkConfigCmd(out, errOut),
		newUpdateNetworkConfigCmd(out, errOut),
	)
	return cmd
}
