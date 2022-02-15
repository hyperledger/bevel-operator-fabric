package channel

import (
	"bytes"
	"github.com/gogo/protobuf/proto"
	"github.com/hyperledger/fabric-config/configtx"
	"github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/resource"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"io"
)

type delAnchorPeerCmd struct {
	configPath  string
	channelName string
	userName    string
	peerHost    string
	peerPort    int
	mspID       string
}

func (c *delAnchorPeerCmd) validate() error {
	return nil
}

func (c *delAnchorPeerCmd) run() error {
	configBackend := config.FromFile(c.configPath)
	sdk, err := fabsdk.New(configBackend)
	if err != nil {
		return err
	}
	mspID := c.mspID
	org1AdminClientContext := sdk.Context(
		fabsdk.WithUser(c.userName),
		fabsdk.WithOrg(mspID),
	)
	resClient, err := resmgmt.New(org1AdminClientContext)
	if err != nil {
		return err
	}
	block, err := resClient.QueryConfigBlockFromOrderer(c.channelName)
	if err != nil {
		return err
	}
	cfgBlock, err := resource.ExtractConfigFromBlock(block)
	if err != nil {
		return err
	}
	cftxGen := configtx.New(cfgBlock)
	app := cftxGen.Application().Organization(mspID)
	err = app.RemoveAnchorPeer(configtx.Address{Host: c.peerHost, Port: c.peerPort})
	if err != nil {
		return err
	}
	configUpdateBytes, err := cftxGen.ComputeMarshaledUpdate(c.channelName)
	if err != nil {
		return err
	}
	configUpdate := &common.ConfigUpdate{}
	err = proto.Unmarshal(configUpdateBytes, configUpdate)
	if err != nil {
		return err
	}
	channelConfigBytes, err := CreateConfigUpdateEnvelope(c.channelName, configUpdate)
	if err != nil {
		return err
	}
	configUpdateReader := bytes.NewReader(channelConfigBytes)
	chResponse, err := resClient.SaveChannel(resmgmt.SaveChannelRequest{
		ChannelID:     c.channelName,
		ChannelConfig: configUpdateReader,
	})
	if err != nil {
		return err
	}
	log.Infof("anchor anchorPeers removed: %s", chResponse.TransactionID)
	return nil
}
func newDelAnchorPeerCMD(io.Writer, io.Writer) *cobra.Command {
	c := &delAnchorPeerCmd{}
	cmd := &cobra.Command{
		Use: "delanchorpeer",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := c.validate(); err != nil {
				return err
			}
			return c.run()
		},
	}
	persistentFlags := cmd.PersistentFlags()
	persistentFlags.StringVarP(&c.peerHost, "peer-host", "", "", "Peer host")
	persistentFlags.IntVarP(&c.peerPort, "peer-port", "", 0, "Peer port")
	persistentFlags.StringVarP(&c.channelName, "channel", "", "", "Channel name")
	persistentFlags.StringVarP(&c.configPath, "config", "", "", "Configuration file for the SDK")
	persistentFlags.StringVarP(&c.userName, "user", "", "", "User name for the transaction")
	persistentFlags.StringVarP(&c.mspID, "msp-id", "", "", "MSPID to remove the anchor peers from")

	cmd.MarkPersistentFlagRequired("channel")
	cmd.MarkPersistentFlagRequired("config")
	cmd.MarkPersistentFlagRequired("user")
	cmd.MarkPersistentFlagRequired("peer-host")
	cmd.MarkPersistentFlagRequired("peer-port")
	return cmd
}
