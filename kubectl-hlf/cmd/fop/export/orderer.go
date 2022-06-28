package export

import "github.com/spf13/cobra"

func newExportOrdererCMD() *cobra.Command {
	return &cobra.Command{
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		Use: "orderer",
	}
}
