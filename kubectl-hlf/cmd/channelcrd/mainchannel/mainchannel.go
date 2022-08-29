package mainchannel

import (
	"io"

	"github.com/spf13/cobra"
)

func NewChannelMainCmd(stdOut io.Writer, stdErr io.Writer) *cobra.Command {
	channelCmd := &cobra.Command{
		Use: "main",
	}
	channelCmd.AddCommand()
	return channelCmd
}
