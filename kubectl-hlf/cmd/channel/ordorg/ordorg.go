package ordorg

import (
	"io"

	"github.com/spf13/cobra"
)

func NewOrdOrgCmd(stdOut io.Writer, stdErr io.Writer) *cobra.Command {
	ordOrgCmd := &cobra.Command{
		Use: "ordorg",
	}
	ordOrgCmd.AddCommand(
		newAddOrgToChannelCMD(stdOut, stdErr),
		newRemoveOrgToChannelCMD(stdOut, stdErr),
	)
	return ordOrgCmd
}
