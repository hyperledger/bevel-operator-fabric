package peer

import (
	"github.com/spf13/cobra"
	"io"
)

func NewPeerCmd(out io.Writer, errOut io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use: "peer",
	}
	cmd.AddCommand(newCreatePeerCmd(out, errOut))
	cmd.AddCommand(newPeerDeleteCmd(out, errOut))
	return cmd
}
