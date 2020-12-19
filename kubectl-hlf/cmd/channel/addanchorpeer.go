package channel

import (
	"bytes"
	"fmt"
	"github.com/gogo/protobuf/proto"
	"github.com/hyperledger/fabric-config/protolator"
	"github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/resource"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/kfsoftware/hlf-operator/internal/github.com/hyperledger/fabric/common/channelconfig"
	"github.com/kfsoftware/hlf-operator/kubectl-hlf/cmd/helpers"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"io"
	"net/url"
	"strconv"
	"strings"
)

type addAnchorPeerCmd struct {
	configPath  string
	peer        string
	channelName string
	userName    string
}

func (c *addAnchorPeerCmd) validate() error {
	return nil
}

func GetConfigEnvelopeBytes(configUpdate *common.ConfigUpdate) ([]byte, error) {
	var buf bytes.Buffer
	if err := protolator.DeepMarshalJSON(&buf, configUpdate); err != nil {
		return nil, err
	}

	channelConfigBytes, err := proto.Marshal(configUpdate)
	if err != nil {
		return nil, err
	}
	configUpdateEnvelope := &common.ConfigUpdateEnvelope{
		ConfigUpdate: channelConfigBytes,
		Signatures:   nil,
	}
	configUpdateEnvelopeBytes, err := proto.Marshal(configUpdateEnvelope)
	if err != nil {
		return nil, err
	}
	payload := &common.Payload{
		Data: configUpdateEnvelopeBytes,
	}
	payloadBytes, err := proto.Marshal(payload)
	if err != nil {
		return nil, err
	}
	configEnvelope := &common.Envelope{
		Payload: payloadBytes,
	}

	return proto.Marshal(configEnvelope)
}

func GetCurrentConfigFromPeer(resClient *resmgmt.Client, channelID string) (*common.Config, error) {
	block, err := resClient.QueryConfigBlockFromOrderer(channelID)
	if err != nil {
		return nil, err
	}
	channelConfig, err := resource.ExtractConfigFromBlock(block)
	if err != nil {
		return nil, err
	}
	return channelConfig, nil
}
func uniqueAnchorPeers(anchorPeers []*peer.AnchorPeer) []*peer.AnchorPeer {
	keys := make(map[string]bool)
	var list []*peer.AnchorPeer
	for _, entry := range anchorPeers {
		key := fmt.Sprintf("%s:%d", entry.Host, entry.Port)
		if _, value := keys[key]; !value {
			keys[key] = true
			list = append(list, entry)
		}
	}
	return list
}
func (c *addAnchorPeerCmd) run() error {
	oclient, err := helpers.GetKubeOperatorClient()
	if err != nil {
		return err
	}
	configBackend := config.FromFile(c.configPath)
	sdk, err := fabsdk.New(configBackend)
	if err != nil {
		return err
	}
	adminPeer, err := helpers.GetPeerByFullName(oclient, c.peer)
	if err != nil {
		return err
	}
	mspID := adminPeer.Spec.MspID
	org1AdminClientContext := sdk.Context(
		fabsdk.WithUser(c.userName),
		fabsdk.WithOrg(mspID),
	)
	resClient, err := resmgmt.New(org1AdminClientContext)
	if err != nil {
		return err
	}

	channelID := c.channelName
	channelConfig, err := GetCurrentConfigFromPeer(resClient, channelID)
	if err != nil {
		return err
	}
	modifiedConfig := &common.Config{}
	modifiedConfigBytes, err := proto.Marshal(channelConfig)
	if err != nil {
		return err
	}
	err = proto.Unmarshal(modifiedConfigBytes, modifiedConfig)
	if err != nil {
		return err
	}
	configValues := modifiedConfig.ChannelGroup.Groups[channelconfig.ApplicationGroupKey].Groups[mspID].Values
	anchorPeersBytes, ok := configValues[channelconfig.AnchorPeersKey]
	anchorPeers := &peer.AnchorPeers{}
	u, err := url.Parse(adminPeer.Status.URL)
	if err != nil {
		return err
	}
	chunks := strings.Split(u.Host, ":")
	peerHost := chunks[0]
	port := chunks[1]
	peerPort, err := strconv.Atoi(port)
	if err != nil {
		return err
	}
	if ok {
		if err := proto.Unmarshal(anchorPeersBytes.Value, anchorPeers); err != nil {
			return err
		}
		anchorPeers.AnchorPeers = append(anchorPeers.AnchorPeers, &peer.AnchorPeer{
			Host: peerHost,
			Port: int32(peerPort),
		})
		anchorPeers.AnchorPeers = uniqueAnchorPeers(anchorPeers.AnchorPeers)
		newAnchorPeerBytes, err := proto.Marshal(anchorPeers)
		if err != nil {
			return err
		}
		configValues[channelconfig.AnchorPeersKey].Value = newAnchorPeerBytes
	} else {
		anchorPeers.AnchorPeers = []*peer.AnchorPeer{
			{
				Host: peerHost,
				Port: int32(peerPort),
			},
		}
		newAnchorPeerBytes, err := proto.Marshal(anchorPeers)
		if err != nil {
			return err
		}
		configValues[channelconfig.AnchorPeersKey] = &common.ConfigValue{
			Version:   0,
			Value:     newAnchorPeerBytes,
			ModPolicy: "Admins",
		}
	}
	confUpdate, err := resmgmt.CalculateConfigUpdate(channelID, channelConfig, modifiedConfig)
	if err != nil {
		return err
	}
	configEnvelopeBytes, err := GetConfigEnvelopeBytes(confUpdate)
	if err != nil {
		return err
	}
	configReader := bytes.NewReader(configEnvelopeBytes)
	txID, err := resClient.SaveChannel(resmgmt.SaveChannelRequest{
		ChannelID:     channelID,
		ChannelConfig: configReader,
	})
	if err != nil {
		return err
	}
	log.Infof("Anchor peer updated, txID=%s", string(txID.TransactionID))
	return nil
}
func newAddAnchorPeerCMD(io.Writer, io.Writer) *cobra.Command {
	c := &addAnchorPeerCmd{}
	cmd := &cobra.Command{
		Use: "addanchorpeer",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := c.validate(); err != nil {
				return err
			}
			return c.run()
		},
	}
	persistentFlags := cmd.PersistentFlags()
	persistentFlags.StringVarP(&c.peer, "peer", "", "", "Admin org to invoke the updates")
	persistentFlags.StringVarP(&c.channelName, "channel", "", "", "Channel name")
	persistentFlags.StringVarP(&c.configPath, "config", "", "", "Configuration file for the SDK")
	persistentFlags.StringVarP(&c.userName, "user", "", "", "User name for the transaction")
	cmd.MarkPersistentFlagRequired("channel")
	cmd.MarkPersistentFlagRequired("config")
	cmd.MarkPersistentFlagRequired("user")
	cmd.MarkPersistentFlagRequired("peer")
	return cmd
}
