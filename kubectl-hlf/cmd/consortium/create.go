package consortium

import (
	"bytes"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/resource"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/kfsoftware/hlf-operator/controllers/testutils"
	"github.com/kfsoftware/hlf-operator/controllers/utils"
	"github.com/kfsoftware/hlf-operator/kubectl-hlf/cmd/helpers"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"io"
)

type createConsortiumCmd struct {
	configPath      string
	consortiumName  string
	peers           []string
	ordererOrg      string
	systemChannelID string
	user            string
}

func (c *createConsortiumCmd) validate() error {
	return nil
}

func (c *createConsortiumCmd) run() error {
	oclient, err := helpers.GetKubeOperatorClient()
	if err != nil {
		return err
	}
	configBackend := config.FromFile(c.configPath)
	sdk, err := fabsdk.New(configBackend)
	if err != nil {
		return err
	}
	ordService, err := helpers.GetOrderingServiceByFullName(oclient, c.ordererOrg)
	if err != nil {
		return err
	}
	mspID := ordService.Spec.MspID

	org1AdminClientContext := sdk.Context(
		fabsdk.WithUser(c.user),
		fabsdk.WithOrg(mspID),
	)
	ordClient, err := resmgmt.New(org1AdminClientContext)
	if err != nil {
		return err
	}
	block, err := ordClient.QueryConfigBlockFromOrderer(c.systemChannelID)
	if err != nil {
		return err
	}
	systemChannelConfig, err := resource.ExtractConfigFromBlock(block)
	if err != nil {
		return err
	}
	ns := ""
	_, peers, err := helpers.GetClusterPeers(oclient, ns)
	if err != nil {
		return err
	}
	var peerOrgs []testutils.PeerOrganization
	for _, peer := range peers {
		if !utils.Contains(c.peers, peer.Name) {
			continue
		}
		certAuth, err := helpers.GetCertAuthByURL(
			oclient,
			peer.Spec.Secret.Enrollment.Component.Cahost,
			peer.Spec.Secret.Enrollment.Component.Caport,
		)
		if err != nil {
			return err
		}
		var nodes []testutils.PeerNode
		peerOrgs = append(peerOrgs, testutils.PeerOrganization{
			RootCert:    certAuth.Status.CACert,
			TLSRootCert: certAuth.Status.TLSCACert,
			MspID:       peer.Spec.MspID,
			Peers:       nodes,
		})
	}
	if len(peerOrgs) == 0 {
		return errors.Errorf("No peer orgs specified")
	}
	modifiedConfig, err := testutils.AddConsortiumToConfig(
		systemChannelConfig,
		testutils.AddConsortiumRequest{
			Name:          c.consortiumName,
			Organizations: peerOrgs,
		},
	)
	if err != nil {
		return err
	}
	confUpdate, err := resmgmt.CalculateConfigUpdate(
		c.systemChannelID,
		systemChannelConfig,
		modifiedConfig,
	)
	if err != nil {
		return err
	}
	configEnvelopeBytes, err := testutils.GetConfigEnvelopeBytes(confUpdate)
	if err != nil {
		return err
	}
	configReader := bytes.NewReader(configEnvelopeBytes)
	saveResponse, err := ordClient.SaveChannel(resmgmt.SaveChannelRequest{
		ChannelID:     c.systemChannelID,
		ChannelConfig: configReader,
	})
	if err != nil {
		return err
	}
	log.Infof("Channel updated, txID=%s", saveResponse.TransactionID)
	return nil
}
func NewCreateConsortiumCMD(io.Writer, io.Writer) *cobra.Command {
	c := &createConsortiumCmd{}
	cmd := &cobra.Command{
		Use: "create",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := c.validate(); err != nil {
				return err
			}
			return c.run()
		},
	}
	persistentFlags := cmd.PersistentFlags()
	persistentFlags.StringVarP(&c.consortiumName, "name", "n", "", "Name of the new consortium")
	persistentFlags.StringVarP(&c.ordererOrg, "orderer-org", "", "", "Ordering service name")
	persistentFlags.StringVarP(&c.systemChannelID, "system-channel-id", "", "", "System channel name")
	persistentFlags.StringVarP(&c.configPath, "config", "", "", "Configuration file for the SDK")
	persistentFlags.StringVarP(&c.user, "user", "", "", "User used to issue the transaction")
	persistentFlags.StringSliceVarP(&c.peers, "peers", "p", nil, "Organizations belonging to the consortium")
	cmd.MarkPersistentFlagRequired("name")
	return cmd
}
