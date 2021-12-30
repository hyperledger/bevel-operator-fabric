package consenter

import (
	"io"

	"github.com/spf13/cobra"
)

func NewConsenterCmd(stdOut io.Writer, stdErr io.Writer) *cobra.Command {
	consenterCmd := &cobra.Command{
		Use: "consenter",
	}
	consenterCmd.AddCommand(
		newAddConsenterCMD(stdOut, stdErr),
		newDelConsenterCMD(stdOut, stdErr),
		newReplaceConsenterCMD(stdOut, stdErr),
	)
	return consenterCmd
}
