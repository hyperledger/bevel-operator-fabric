package chaincode

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/kfsoftware/hlf-operator/kubectl-hlf/cmd/helpers"
	"github.com/spf13/cobra"
	"io"
)

type queryChaincodeCmd struct {
	configPath string
	peer       string
	userName   string
	channel    string
	chaincode  string
	fcn        string
	args       []string
	transient  string
	mspID      string
}

func (c *queryChaincodeCmd) validate() error {
	return nil
}
func (c *queryChaincodeCmd) run(out io.Writer) error {
	var mspID string
	if c.mspID != "" {
		mspID = c.mspID
	} else {
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
		mspID = peer.Spec.MspID
	}
	peerName := c.peer
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
	var transientMap map[string][]byte
	if c.transient != "" {
		err = json.Unmarshal([]byte(c.transient), &transientMap)
		if err != nil {
			return err
		}
	}

	response, err := ch.Query(
		channel.Request{
			ChaincodeID:     c.chaincode,
			Fcn:             c.fcn,
			Args:            args,
			TransientMap:    transientMap,
			InvocationChain: nil,
			IsInit:          false,
		},
		channel.WithTargetEndpoints(peerName),
	)
	if err != nil {
		return err
	}
	_, err = fmt.Fprint(out, string(response.Payload))
	if err != nil {
		return err
	}
	return nil
}

func newQueryChaincodeCMD(out io.Writer, errOut io.Writer) *cobra.Command {
	c := &queryChaincodeCmd{}
	cmd := &cobra.Command{
		Use: "query",
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
	persistentFlags.StringVarP(&c.transient, "transient", "t", "", "Transient map")
	persistentFlags.StringVarP(&c.mspID, "mspID", "", "", "MSP ID")
	cmd.MarkPersistentFlagRequired("user")
	cmd.MarkPersistentFlagRequired("peer")
	cmd.MarkPersistentFlagRequired("config")
	cmd.MarkPersistentFlagRequired("chaincode")
	cmd.MarkPersistentFlagRequired("fcn")
	return cmd
}
