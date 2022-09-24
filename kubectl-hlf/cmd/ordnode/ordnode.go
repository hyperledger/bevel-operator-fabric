package ordnode

import (
	"io"

	"github.com/spf13/cobra"
)

// NewOrdNodeCmd creates a new root command to manage Ordering Services
func NewOrdNodeCmd(out io.Writer, errOut io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use: "ordnode",
	}
	cmd.AddCommand(
		newCreateOrdererNodeCmd(out, errOut),
		newOrdererNodeDeleteCmd(out, errOut),
		newJoinChannelCMD(out, errOut),
		newRenewChannelCMD(out, errOut),
		newRemoveChannelCMD(out, errOut),
		newUpgradeOrdererCMD(out, errOut),
		newUpdateOrdererCMD(out, errOut),
	)
	return cmd
}
