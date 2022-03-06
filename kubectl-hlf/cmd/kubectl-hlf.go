package cmd

import (
	"github.com/kfsoftware/hlf-operator/kubectl-hlf/cmd/ca"
	"github.com/kfsoftware/hlf-operator/kubectl-hlf/cmd/chaincode"
	"github.com/kfsoftware/hlf-operator/kubectl-hlf/cmd/channel"
	"github.com/kfsoftware/hlf-operator/kubectl-hlf/cmd/externalchaincode"
	"github.com/kfsoftware/hlf-operator/kubectl-hlf/cmd/inspect"
	"github.com/kfsoftware/hlf-operator/kubectl-hlf/cmd/networkconfig"
	"github.com/kfsoftware/hlf-operator/kubectl-hlf/cmd/ordnode"
	"github.com/kfsoftware/hlf-operator/kubectl-hlf/cmd/org"
	"github.com/kfsoftware/hlf-operator/kubectl-hlf/cmd/peer"
	"github.com/kfsoftware/hlf-operator/kubectl-hlf/cmd/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	// Workaround for authentication plugins https://krew.sigs.k8s.io/docs/developer-guide/develop/best-practices/#auth-plugins
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

const (
	hlfDesc = `
kubectl plugin to manage HLF operator CRDs.`
)

// NewCmdHLF creates a new root command for kubectl-hlf
func NewCmdHLF() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "hlf",
		Short:        "manage HLF operator CRDs",
		Long:         hlfDesc,
		SilenceUsage: true,
	}
	logrus.SetLevel(logrus.DebugLevel)
	cmd.AddCommand(
		inspect.NewInspectHLFConfig(cmd.OutOrStdout()),
		channel.NewChannelCmd(cmd.OutOrStdout(), cmd.ErrOrStderr()),
		ca.NewCACmd(cmd.OutOrStdout(), cmd.ErrOrStderr()),
		peer.NewPeerCmd(cmd.OutOrStdout(), cmd.ErrOrStderr()),
		ordnode.NewOrdNodeCmd(cmd.OutOrStdout(), cmd.ErrOrStderr()),
		chaincode.NewChaincodeCmd(cmd.OutOrStdout(), cmd.ErrOrStderr()),
		org.NewOrgCmd(cmd.OutOrStdout(), cmd.ErrOrStderr()),
		utils.NewUtilsCMD(cmd.OutOrStdout(), cmd.ErrOrStderr()),
		networkconfig.NewNetworkConfigCmd(cmd.OutOrStdout(), cmd.ErrOrStderr()),
		externalchaincode.NewExternalChaincodeCmd(cmd.OutOrStdout(), cmd.ErrOrStderr()),
)
	return cmd
}
