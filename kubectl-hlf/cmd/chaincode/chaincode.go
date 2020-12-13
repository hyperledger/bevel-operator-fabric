package chaincode

import (
	"github.com/spf13/cobra"
	"io"
)

func NewChaincodeCmd(stdOut io.Writer, stdErr io.Writer) *cobra.Command {
	consortiumCmd := &cobra.Command{
		Use: "chaincode",
	}

	consortiumCmd.AddCommand(newChaincodeInstallCMD(stdOut, stdErr))
	consortiumCmd.AddCommand(newChaincodeQueryInstalledCMD(stdOut, stdErr))
	consortiumCmd.AddCommand(newChaincodeApproveCMD(stdOut, stdErr))
	consortiumCmd.AddCommand(newChaincodeCommitCMD(stdOut, stdErr))
	consortiumCmd.AddCommand(newQueryChaincodeCMD(stdOut, stdErr))
	consortiumCmd.AddCommand(newInvokeChaincodeCMD(stdOut, stdErr))
	return consortiumCmd
}
