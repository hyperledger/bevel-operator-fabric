package identity

import "github.com/spf13/cobra"

func NewIdentityCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "identity",
		Short: "Manage HLF identities",
		Long:  `Manage HLF identities`,
	}
	cmd.AddCommand(newIdentityCreateCMD())
	cmd.AddCommand(newIdentityUpdateCMD())
	cmd.AddCommand(newIdentityDeleteCMD())
	return cmd
}
