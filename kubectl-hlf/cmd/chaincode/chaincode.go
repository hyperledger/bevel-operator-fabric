package chaincode

import (
	"github.com/spf13/cobra"
	"io"
)

func NewChaincodeCmd(stdOut io.Writer, stdErr io.Writer) *cobra.Command {
	consortiumCmd := &cobra.Command{
		Use: "chaincode",
	}

	consortiumCmd.AddCommand(
		newChaincodeInstallCMD(stdOut, stdErr),
		newChaincodeQueryInstalledCMD(stdOut, stdErr),
		newChaincodeApproveCMD(stdOut, stdErr),
		newChaincodeCommitCMD(stdOut, stdErr),
		newQueryChaincodeCMD(stdOut, stdErr),
		newInvokeChaincodeCMD(stdOut, stdErr),
		newQueryCommittedCMD(stdOut, stdErr),
		newQueryApprovedCMD(stdOut, stdErr),
		newCalculatePackageIDCMD(stdOut, stdErr),
		newGetLatestInfoCMD(stdOut, stdErr),
	)
	return consortiumCmd
}
