package chaincode

import (
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/common/policydsl"
	"io"
	"os"
	"time"

	"github.com/hyperledger/fabric-protos-go/common"
	pb "github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/kfsoftware/hlf-operator/kubectl-hlf/cmd/helpers"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type approveChaincodeCmd struct {
	configPath        string
	peer              string
	userName          string
	channelName       string
	packageID         string
	version           string
	name              string
	sequence          int64
	policy            string
	initRequired      bool
	collectionsConfig string
}

func (c *approveChaincodeCmd) validate() error {
	return nil
}
func (c *approveChaincodeCmd) run() error {
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
	peerName := peer.Name
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
	var sp *common.SignaturePolicyEnvelope
	if c.policy != "" {
		sp, err = policydsl.FromString(c.policy)
		if err != nil {
			return err
		}
	}
	var collectionConfigs []*pb.CollectionConfig

	if c.collectionsConfig != "" {
		collectionBytes, err := os.ReadFile(c.collectionsConfig)
		if err != nil {
			return err
		}
		collectionConfigs, err = helpers.GetCollectionConfigFromBytes(collectionBytes)
		if err != nil {
			return err
		}
	}
	if len(collectionConfigs) == 0 {
		collectionConfigs = nil
	}

	txID, err := resClient.LifecycleApproveCC(
		c.channelName,
		resmgmt.LifecycleApproveCCRequest{
			Name:              c.name,
			Version:           c.version,
			PackageID:         c.packageID,
			Sequence:          c.sequence,
			EndorsementPlugin: "escc",
			ValidationPlugin:  "vscc",
			SignaturePolicy:   sp,
			CollectionConfig:  collectionConfigs,
			InitRequired:      c.initRequired,
		},
		resmgmt.WithTargetEndpoints(peerName),
		resmgmt.WithTimeout(fab.ResMgmt, 20*time.Minute),
		resmgmt.WithTimeout(fab.PeerResponse, 20*time.Minute),
	)
	if err != nil {
		return err
	}
	log.Infof("Chaincode approved=%s", txID)
	return nil
}
func newChaincodeApproveCMD(out io.Writer, errOut io.Writer) *cobra.Command {
	c := &approveChaincodeCmd{}
	cmd := &cobra.Command{
		Use: "approveformyorg",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := c.validate(); err != nil {
				return err
			}
			return c.run()
		},
	}
	persistentFlags := cmd.PersistentFlags()
	persistentFlags.StringVarP(&c.configPath, "config", "", "", "Configuration file for the SDK")
	persistentFlags.StringVarP(&c.peer, "peer", "p", "", "Peer org to invoke the updates")
	persistentFlags.StringVarP(&c.userName, "user", "", "", "User name for the transaction")
	persistentFlags.StringVarP(&c.channelName, "channel", "", "", "Channel name")
	persistentFlags.StringVarP(&c.packageID, "package-id", "", "", "Package ID")
	persistentFlags.StringVarP(&c.version, "version", "", "1.0", "Version")
	persistentFlags.StringVarP(&c.name, "name", "", "", "Chaincode name")
	persistentFlags.Int64VarP(&c.sequence, "sequence", "", 1, "Sequence number")
	persistentFlags.StringVarP(&c.policy, "policy", "", "", "Policy")
	persistentFlags.BoolVarP(&c.initRequired, "init-required", "", false, "Init required")
	persistentFlags.StringVarP(&c.collectionsConfig, "collections-config", "", "", "Private data collections")

	cmd.MarkPersistentFlagRequired("user")
	cmd.MarkPersistentFlagRequired("peer")
	cmd.MarkPersistentFlagRequired("config")
	cmd.MarkPersistentFlagRequired("channelName")
	cmd.MarkPersistentFlagRequired("package-id")
	cmd.MarkPersistentFlagRequired("name")
	return cmd
}
