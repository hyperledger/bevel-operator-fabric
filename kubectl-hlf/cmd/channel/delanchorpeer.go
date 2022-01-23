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
	//channelConfig, err := GetCurrentConfigFromPeer(resClient, channelID)
	//if err != nil {
	//	return err
	//}
	//modifiedConfig := &common.Config{}
	//modifiedConfigBytes, err := proto.Marshal(channelConfig)
	//if err != nil {
	//	return err
	//}
	//err = proto.Unmarshal(modifiedConfigBytes, modifiedConfig)
	//if err != nil {
	//	return err
	//}
	//configValues := modifiedConfig.ChannelGroup.Groups[channelconfig.ApplicationGroupKey].Groups[mspID].Values
	//anchorPeersBytes, ok := configValues[channelconfig.AnchorPeersKey]
	//anchorPeers := &peer.AnchorPeers{}
	//u, err := url.Parse(adminPeer.Status.URL)
	//if err != nil {
	//	return err
	//}
	//chunks := strings.Split(u.Host, ":")
	//peerHost := chunks[0]
	//port := chunks[1]
	//peerPort, err := strconv.Atoi(port)
	//if err != nil {
	//	return err
	//}
	//if ok {
	//	if err := proto.Unmarshal(anchorPeersBytes.Value, anchorPeers); err != nil {
	//		return err
	//	}
	//	var finalAnchorPeers []*peer.AnchorPeer
	//	for _, anchorPeer := range anchorPeers.AnchorPeers {
	//		if anchorPeer.Port != int32(peerPort) && anchorPeer.Host != peerHost {
	//			finalAnchorPeers = append(finalAnchorPeers, anchorPeer)
	//		} else {
	//			log.Infof("Removing anchor peer with host=%s port=%d", peerHost, peerPort)
	//		}
	//	}
	//	anchorPeers.AnchorPeers = uniqueAnchorPeers(finalAnchorPeers)
	//	log.Info("Anchor peers=%v", anchorPeers.AnchorPeers)
	//	newAnchorPeerBytes, err := proto.Marshal(anchorPeers)
	//	if err != nil {
	//		return err
	//	}
	//	configValues[channelconfig.AnchorPeersKey].Value = newAnchorPeerBytes
	//} else {
	//	log.Warn("No anchor peers configured for this organization in this channel")
	//	anchorPeers.AnchorPeers = []*peer.AnchorPeer{}
	//	newAnchorPeerBytes, err := proto.Marshal(anchorPeers)
	//	if err != nil {
	//		return err
	//	}
	//	configValues[channelconfig.AnchorPeersKey] = &common.ConfigValue{
	//		Version:   0,
	//		Value:     newAnchorPeerBytes,
	//		ModPolicy: "Admins",
	//	}
	//}
	//confUpdate, err := resmgmt.CalculateConfigUpdate(channelID, channelConfig, modifiedConfig)
	//if err != nil {
	//	return err
	//}
	//configEnvelopeBytes, err := GetConfigEnvelopeBytes(confUpdate)
	//if err != nil {
	//	return err
	//}
	//configReader := bytes.NewReader(configEnvelopeBytes)
	//txID, err := resClient.SaveChannel(resmgmt.SaveChannelRequest{
	//	ChannelID:     channelID,
	//	ChannelConfig: configReader,
	//})
	//if err != nil {
	//	return err
	//}
	//log.Infof("Anchor peer updated, txID=%s", string(txID.TransactionID))
	//return nil
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
