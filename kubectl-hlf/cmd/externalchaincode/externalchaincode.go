package externalchaincode

import (
	"io"

	"github.com/spf13/cobra"
)

func NewExternalChaincodeCmd(stdOut io.Writer, stdErr io.Writer) *cobra.Command {
	externalChaincodeCmd := &cobra.Command{
		Use: "externalchaincode",
	}
	externalChaincodeCmd.AddCommand(
		newExternalChaincodeCreateCmd(),
		newExternalChaincodeUpdateCmd(),
		newExternalChaincodeDeleteCmd(),
		newExternalChaincodeSyncCmd(),
	)
	return externalChaincodeCmd
}
