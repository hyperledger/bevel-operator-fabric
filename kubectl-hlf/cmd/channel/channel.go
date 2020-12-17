package channel

import (
	"github.com/spf13/cobra"
	"io"
)

func NewChannelCmd(stdOut io.Writer, stdErr io.Writer) *cobra.Command {
	consortiumCmd := &cobra.Command{
		Use: "channel",
	}
	consortiumCmd.AddCommand(newCreateChannelCMD(stdOut, stdErr))
	consortiumCmd.AddCommand(newJoinChannelCMD(stdOut, stdErr))
	consortiumCmd.AddCommand(newAddAnchorPeerCMD(stdOut, stdErr))
	consortiumCmd.AddCommand(newInspectChannelCMD(stdOut, stdErr))
	consortiumCmd.AddCommand(newTopChannelCMD(stdOut, stdErr))
	consortiumCmd.AddCommand(newAddOrgToChannelCMD(stdOut, stdErr))
	return consortiumCmd
}
