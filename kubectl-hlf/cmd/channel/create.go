package channel

import (
	"bytes"
	"fmt"
	"io"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/resource"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/resource/genesisconfig"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/kfsoftware/hlf-operator/controllers/testutils"
	"github.com/kfsoftware/hlf-operator/controllers/utils"
	"github.com/kfsoftware/hlf-operator/kubectl-hlf/cmd/helpers"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type createChannelCmd struct {
	configPath     string
	organizations  []string
	adminOrg       string
	ordererOrg     string
	channelName    string
	consortiumName string
	userName       string
	ordererNodes   []string
}

func (c *createChannelCmd) validate() error {
	return nil
}
func (c *createChannelCmd) run() error {
	oclient, err := helpers.GetKubeOperatorClient()
	if err != nil {
		return err
	}
	clientSet, err := helpers.GetKubeClient()
	if err != nil {
		return err
	}
	configBackend := config.FromFile(c.configPath)
	sdk, err := fabsdk.New(configBackend)
	if err != nil {
		return err
	}
	ordService, err := helpers.GetOrderingServiceByFullName(clientSet, oclient, c.ordererOrg)
	if err != nil {
		return err
	}
	adminPeer, err := helpers.GetPeerByFullName(clientSet, oclient, c.adminOrg)
	if err != nil {
		return err
	}
	org1AdminClientContext := sdk.Context(
		fabsdk.WithUser(c.userName),
		fabsdk.WithOrg(adminPeer.Spec.MspID),
	)
	resClient, err := resmgmt.New(org1AdminClientContext)
	if err != nil {
		return err
	}
	ns := ""
	_, peers, err := helpers.GetClusterPeers(clientSet, oclient, ns)
	if err != nil {
		return err
	}
	var peerOrgs []testutils.PeerOrganization
	for _, peer := range peers {
		if !utils.Contains(c.organizations, peer.Name) {
			continue
		}
		certAuth, err := helpers.GetCertAuthByURL(
			clientSet,
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
	certAuth, err := helpers.GetCertAuthByURL(
		clientSet,
		oclient,
		ordService.Spec.Enrollment.Component.Cahost,
		ordService.Spec.Enrollment.Component.Caport,
	)
	if err != nil {
		return err
	}
	ordOrganization := testutils.OrdererOrganization{
		Nodes:        []testutils.OrdererNode{},
		RootTLSCert:  certAuth.Status.TLSCACert,
		RootSignCert: certAuth.Status.CACert,
		MspID:        c.ordererOrg,
	}

	profileConfig, err := testutils.GetChannelProfileConfig(
		ordOrganization,
		peerOrgs,
		c.consortiumName,
		fmt.Sprintf(`OR('%s.admin')`, adminPeer.Spec.MspID),
	)
	if err != nil {
		return err
	}
	var baseProfile *genesisconfig.Profile
	channelTx, err := resource.CreateChannelCreateTx(
		profileConfig,
		baseProfile,
		c.channelName,
	)
	if err != nil {
		return err
	}
	channelConfig := bytes.NewReader(channelTx)
	saveResponse, err := resClient.SaveChannel(resmgmt.SaveChannelRequest{
		ChannelID:     c.channelName,
		ChannelConfig: channelConfig,
	})
	if err != nil {
		return err
	}
	log.Infof("Channel created, txID=%s", saveResponse.TransactionID)
	return nil
}
func newCreateChannelCMD(io.Writer, io.Writer) *cobra.Command {
	c := &createChannelCmd{}
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
	persistentFlags.StringVarP(&c.adminOrg, "admin-org", "", "", "Admin org to invoke the updates")
	persistentFlags.StringVarP(&c.ordererOrg, "ordering-service", "", "", "Orderer Service name")
	persistentFlags.StringVarP(&c.channelName, "name", "", "", "Channel name")
	persistentFlags.StringVarP(&c.configPath, "config", "", "", "Configuration file for the SDK")
	persistentFlags.StringVarP(&c.userName, "user", "", "", "User name for the transaction")
	persistentFlags.StringVarP(&c.consortiumName, "consortium", "", "", "Consortium name")
	persistentFlags.StringSliceVarP(&c.organizations, "organizations", "p", []string{}, "Organizations belonging to the consortium")
	persistentFlags.StringSliceVarP(&c.ordererNodes, "orderer-nodes", "o", []string{}, "Consenter orderer nodes")
	cmd.MarkPersistentFlagRequired("channel-id")
	cmd.MarkPersistentFlagRequired("organizations")
	cmd.MarkPersistentFlagRequired("config")
	cmd.MarkPersistentFlagRequired("peer-org")
	cmd.MarkPersistentFlagRequired("consortium")
	cmd.MarkPersistentFlagRequired("orderer-org")
	return cmd
}
