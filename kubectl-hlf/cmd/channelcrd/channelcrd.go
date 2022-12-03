package channelcrd

import (
	"github.com/kfsoftware/hlf-operator/kubectl-hlf/cmd/channelcrd/follower"
	"github.com/kfsoftware/hlf-operator/kubectl-hlf/cmd/channelcrd/mainchannel"
	"io"

	"github.com/spf13/cobra"
)

func NewChannelCRDCmd(stdOut io.Writer, stdErr io.Writer) *cobra.Command {
	channelCmd := &cobra.Command{
		Use: "channelcrd",
	}
	channelCmd.AddCommand(
		mainchannel.NewChannelMainCmd(stdOut, stdErr),
		follower.NewChannelFollowerCmd(stdOut, stdErr),
	)
	return channelCmd
}
