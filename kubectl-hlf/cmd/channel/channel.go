package channel

import (
	"io"

	"github.com/spf13/cobra"
)

func NewChannelCmd(stdOut io.Writer, stdErr io.Writer) *cobra.Command {
	channelCmd := &cobra.Command{
		Use: "channel",
	}
	channelCmd.AddCommand(
		newUpdateChannelCMD(stdOut, stdErr),
		newInspectChannelCMD(stdOut, stdErr),
		newTopChannelCMD(stdOut, stdErr),
	)
	return channelCmd
}
