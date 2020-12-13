package consortium

import (
	"github.com/spf13/cobra"
	"io"
)

func NewConsortiumCmd(stdOut io.Writer, stdErr io.Writer) *cobra.Command {
	consortiumCmd := &cobra.Command{
		Use: "consortiums",
	}

	consortiumCmd.AddCommand(NewCreateConsortiumCMD(stdOut, stdErr))
	return consortiumCmd
}
