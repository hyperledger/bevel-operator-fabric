package chaincodecrd

import (
	"io"

	"github.com/spf13/cobra"
)

func newApproveCmd(out io.Writer, errOut io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "approve",
		Short: "Approve a chaincode definition",
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement chaincode approval logic
			return nil
		},
	}
	return cmd
}
