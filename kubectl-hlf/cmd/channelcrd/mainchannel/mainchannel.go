package mainchannel

import (
	"io"

	"github.com/spf13/cobra"
)

func NewChannelMainCmd(stdOut io.Writer, stdErr io.Writer) *cobra.Command {
	channelCmd := &cobra.Command{
		Use: "main",
	}
	channelCmd.AddCommand(
		newCreateMainChannelCmd(stdOut, stdErr),
		newUpdateMainChannelCmd(stdOut, stdErr),
		newDeleteMainChannelCmd(stdOut, stdErr),
	)
	return channelCmd
}
