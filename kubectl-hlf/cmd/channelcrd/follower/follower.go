package follower

import (
	"io"

	"github.com/spf13/cobra"
)

func NewChannelFollowerCmd(stdOut io.Writer, stdErr io.Writer) *cobra.Command {
	channelCmd := &cobra.Command{
		Use: "follower",
	}
	channelCmd.AddCommand(
		newCreateFollowerChannelCmd(stdOut, stdErr),
		newUpdateFollowerChannelCmd(stdOut, stdErr),
		newDeleteFollowerChannelCmd(stdOut, stdErr),
	)
	return channelCmd
}
