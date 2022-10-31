package chaincode

import (
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/kfsoftware/hlf-operator/kubectl-hlf/cmd/helpers"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"io"
)

type invokeChaincodeCmd struct {
	configPath string
	peer       string
	userName   string
	channel    string
	chaincode  string
	fcn        string
	args       []string
}

func (c *invokeChaincodeCmd) validate() error {
	return nil
}
func (c *invokeChaincodeCmd) run(out io.Writer) error {
	oclient, err := helpers.GetKubeOperatorClient()
	if err != nil {
		return err
	}
	clientSet, err := helpers.GetKubeClient()
	if err != nil {
		return err
	}
	peer, err := helpers.GetPeerByFullName(clientSet, oclient, c.peer)
	if err != nil {
		return err
	}
	mspID := peer.Spec.MspID
	configBackend := config.FromFile(c.configPath)
	sdk, err := fabsdk.New(configBackend)
	if err != nil {
		return err
	}
	chContext := sdk.ChannelContext(
		c.channel,
		fabsdk.WithUser(c.userName),
		fabsdk.WithOrg(mspID),
	)
	ch, err := channel.New(chContext)
	if err != nil {
		return err
	}
	var args [][]byte
	for _, arg := range c.args {
		args = append(args, []byte(arg))
	}

	response, err := ch.Execute(
		channel.Request{
			ChaincodeID:     c.chaincode,
			Fcn:             c.fcn,
			Args:            args,
			TransientMap:    nil,
			InvocationChain: nil,
			IsInit:          false,
		},
	)
	if err != nil {
		return err
	}
	_, err = fmt.Fprint(out, string(response.Payload))
	if err != nil {
		return err
	}
	log.Infof("txid=%s", response.TransactionID)
	return nil
}
func newInvokeChaincodeCMD(out io.Writer, errOut io.Writer) *cobra.Command {
	c := &invokeChaincodeCmd{}
	cmd := &cobra.Command{
		Use: "invoke",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := c.validate(); err != nil {
				return err
			}
			return c.run(out)
		},
	}
	persistentFlags := cmd.PersistentFlags()
	persistentFlags.StringVarP(&c.peer, "peer", "p", "", "Peer org to invoke the updates")
	persistentFlags.StringVarP(&c.userName, "user", "", "", "User name for the transaction")
	persistentFlags.StringVarP(&c.configPath, "config", "", "", "Configuration file for the SDK")
	persistentFlags.StringVarP(&c.channel, "channel", "", "", "Channel name")
	persistentFlags.StringVarP(&c.chaincode, "chaincode", "", "", "Chaincode label")
	persistentFlags.StringVarP(&c.fcn, "fcn", "", "", "Function name")
	persistentFlags.StringArrayVarP(&c.args, "args", "a", []string{}, "Function arguments")
	cmd.MarkPersistentFlagRequired("user")
	cmd.MarkPersistentFlagRequired("peer")
	cmd.MarkPersistentFlagRequired("config")
	cmd.MarkPersistentFlagRequired("chaincode")
	cmd.MarkPersistentFlagRequired("fcn")
	return cmd
}
