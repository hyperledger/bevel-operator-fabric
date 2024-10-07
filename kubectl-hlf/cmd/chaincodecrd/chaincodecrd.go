package chaincodecrd

import (
	"io"

	"github.com/spf13/cobra"
)

func NewChaincodeCRDCmd(out io.Writer, errOut io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "chaincodecrd",
		Short: "Manage chaincode CRDs",
		Long:  "Manage chaincode CRDs for installation, approval, and commitment",
	}

	cmd.AddCommand(newInstallCmd(out, errOut))
	cmd.AddCommand(newApproveCmd(out, errOut))
	cmd.AddCommand(newCommitCmd(out, errOut))

	return cmd
}
