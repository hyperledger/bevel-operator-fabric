package cmd

import (
	"github.com/kfsoftware/hlf-operator/kubectl-hlf/cmd/ca"
	"github.com/kfsoftware/hlf-operator/kubectl-hlf/cmd/chaincode"
	"github.com/kfsoftware/hlf-operator/kubectl-hlf/cmd/channel"
	"github.com/kfsoftware/hlf-operator/kubectl-hlf/cmd/consortium"
	"github.com/kfsoftware/hlf-operator/kubectl-hlf/cmd/inspect"
	"github.com/kfsoftware/hlf-operator/kubectl-hlf/cmd/ordservice"
	"github.com/kfsoftware/hlf-operator/kubectl-hlf/cmd/peer"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"

	// Workaround for authentication plugins https://krew.sigs.k8s.io/docs/developer-guide/develop/best-practices/#auth-plugins
	_ "k8s.io/client-go/plugin/pkg/client/auth"

)

const (
	hlfDesc = `
kubectl plugin to manage HLF operator CRDs.`
)

// NewCmdHLF creates a new root command for kubectl-hlf
func NewCmdHLF(streams genericclioptions.IOStreams) *cobra.Command {
	cmd := &cobra.Command{
		Use:          "hlf",
		Short:        "manage HLF operator CRDs",
		Long:         hlfDesc,
		SilenceUsage: true,
	}
	cmd.AddCommand(inspect.NewInspectHLFConfig(cmd.OutOrStdout()))
	cmd.AddCommand(consortium.NewConsortiumCmd(cmd.OutOrStdout(), cmd.ErrOrStderr()))
	cmd.AddCommand(channel.NewChannelCmd(cmd.OutOrStdout(), cmd.ErrOrStderr()))
	cmd.AddCommand(ca.NewCACmd(cmd.OutOrStdout(), cmd.ErrOrStderr()))
	cmd.AddCommand(peer.NewPeerCmd(cmd.OutOrStdout(), cmd.ErrOrStderr()))
	cmd.AddCommand(ordservice.NewOrdServiceCmd(cmd.OutOrStdout(), cmd.ErrOrStderr()))
	cmd.AddCommand(chaincode.NewChaincodeCmd(cmd.OutOrStdout(), cmd.ErrOrStderr()))
	return cmd
}
