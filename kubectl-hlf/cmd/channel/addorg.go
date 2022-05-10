package channel

import (
	"bytes"
	"github.com/kfsoftware/hlf-operator/internal/github.com/hyperledger/fabric/common/channelconfig"
	"io"
	"io/ioutil"

	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric-config/protolator"
	"github.com/hyperledger/fabric-protos-go/common"
	cb "github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/kfsoftware/hlf-operator/internal/github.com/hyperledger/fabric/sdkinternal/configtxgen/encoder"
	"github.com/kfsoftware/hlf-operator/internal/github.com/hyperledger/fabric/sdkinternal/configtxgen/genesisconfig"
	"github.com/kfsoftware/hlf-operator/kubectl-hlf/cmd/helpers"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

type addOrgCmd struct {
	configPath  string
	orgPath     string
	peer        string
	channelName string
	userName    string
	mspID       string
	dryRun      bool
}

func (c *addOrgCmd) validate() error {
	return nil
}
func (c *addOrgCmd) run(out io.Writer) error {
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
	org1AdminClientContext := sdk.Context(
		fabsdk.WithUser(c.userName),
		fabsdk.WithOrg(mspID),
	)
	resClient, err := resmgmt.New(org1AdminClientContext)
	if err != nil {
		return err
	}
	channelID := c.channelName
	channelConfig, err := helpers.GetCurrentConfigFromPeer(resClient, channelID)
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
	orgBytes, err := ioutil.ReadFile(c.orgPath)
	if err != nil {
		return err
	}
	topLevel := &genesisconfig.TopLevel{}
	err = yaml.Unmarshal(orgBytes, topLevel)
	if err != nil {
		return err
	}
	var orgConfig *cb.ConfigGroup
	for _, org := range topLevel.Organizations {
		if org.Name == c.mspID {
			orgConfig, err = encoder.NewApplicationOrgGroup(org)
			if err != nil {
				return err
			}
		}
	}
	if orgConfig == nil {
		return errors.Errorf("msp ID %s not found", c.mspID)
	}
	modifiedConfig.ChannelGroup.Groups[channelconfig.ApplicationGroupKey].Groups[c.mspID] = orgConfig
	confUpdate, err := resmgmt.CalculateConfigUpdate(channelID, channelConfig, modifiedConfig)
	if err != nil {
		return err
	}
	if c.dryRun {
		err = protolator.DeepMarshalJSON(out, confUpdate)
		if err != nil {
			return err
		}
	} else {
		configEnvelopeBytes, err := helpers.GetConfigEnvelopeBytes(confUpdate)
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
		log.Infof("Channel updated, txID=%s", string(txID.TransactionID))
	}
	return nil
}
func newAddOrgToChannelCMD(out io.Writer, errOut io.Writer) *cobra.Command {
	c := &addOrgCmd{}
	cmd := &cobra.Command{
		Use: "addorg",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := c.validate(); err != nil {
				return err
			}
			return c.run(out)
		},
	}
	persistentFlags := cmd.PersistentFlags()
	persistentFlags.StringVarP(&c.peer, "peer", "", "", "Name of the peer to invoke the updates")
	persistentFlags.StringVarP(&c.channelName, "name", "", "", "Channel name")
	persistentFlags.StringVarP(&c.configPath, "config", "", "", "Configuration file for the SDK")
	persistentFlags.StringVarP(&c.userName, "user", "", "", "User name for the transaction")
	persistentFlags.StringVarP(&c.mspID, "msp-id", "", "", "MSP ID for the new organization")
	persistentFlags.StringVarP(&c.orgPath, "org-config", "", "", "JSON with the crypto material for the new org")
	persistentFlags.BoolVarP(&c.dryRun, "dry-run", "", false, "Output configuration as JSON and not update")
	cmd.MarkPersistentFlagRequired("name")
	cmd.MarkPersistentFlagRequired("config")
	cmd.MarkPersistentFlagRequired("peer")
	cmd.MarkPersistentFlagRequired("org-config")
	cmd.MarkPersistentFlagRequired("msp-id")
	cmd.MarkPersistentFlagRequired("user")
	return cmd
}
