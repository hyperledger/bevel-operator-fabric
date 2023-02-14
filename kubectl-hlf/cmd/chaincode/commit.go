package chaincode

import (
	"github.com/hyperledger/fabric-protos-go/common"
	pb "github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/common/policydsl"
	"github.com/kfsoftware/hlf-operator/kubectl-hlf/cmd/helpers"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"io"
	"io/ioutil"
	"time"
)

type commitChaincodeCmd struct {
	configPath        string
	userName          string
	channelName       string
	version           string
	name              string
	sequence          int64
	policy            string
	initRequired      bool
	collectionsConfig string
	mspID             string
}

func (c *commitChaincodeCmd) validate() error {
	return nil
}
func (c *commitChaincodeCmd) run() error {
	configBackend := config.FromFile(c.configPath)
	sdk, err := fabsdk.New(configBackend)
	if err != nil {
		return err
	}
	org1AdminClientContext := sdk.Context(
		fabsdk.WithUser(c.userName),
		fabsdk.WithOrg(c.mspID),
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
		collectionBytes, err := ioutil.ReadFile(c.collectionsConfig)
		if err != nil {
			return err
		}
		collectionConfigs, err = helpers.GetCollectionConfigFromBytes(collectionBytes)
		if err != nil {
			return err
		}
	}
	txID, err := resClient.LifecycleCommitCC(
		c.channelName,
		resmgmt.LifecycleCommitCCRequest{
			Name:              c.name,
			Version:           c.version,
			Sequence:          c.sequence,
			EndorsementPlugin: "escc",
			ValidationPlugin:  "vscc",
			SignaturePolicy:   sp,
			CollectionConfig:  collectionConfigs,
			InitRequired:      c.initRequired,
		},
		resmgmt.WithTimeout(fab.ResMgmt, 20*time.Minute),
		resmgmt.WithTimeout(fab.PeerResponse, 20*time.Minute),
	)
	if err != nil {
		return err
	}
	log.Infof("Chaincode commited=%s", txID)
	return nil
}
func newChaincodeCommitCMD(out io.Writer, errOut io.Writer) *cobra.Command {
	c := &commitChaincodeCmd{}
	cmd := &cobra.Command{
		Use: "commit",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := c.validate(); err != nil {
				return err
			}
			return c.run()
		},
	}
	persistentFlags := cmd.PersistentFlags()
	persistentFlags.StringVarP(&c.configPath, "config", "", "", "Configuration file for the SDK")
	persistentFlags.StringVarP(&c.userName, "user", "", "", "User name for the transaction")
	persistentFlags.StringVarP(&c.channelName, "channel", "", "", "Channel name")
	persistentFlags.StringVarP(&c.version, "version", "", "1.0", "Version")
	persistentFlags.StringVarP(&c.name, "name", "", "", "Chaincode name")
	persistentFlags.StringVarP(&c.mspID, "mspid", "", "", "MSP ID of the organization")
	persistentFlags.Int64VarP(&c.sequence, "sequence", "", 1, "Sequence number")
	persistentFlags.StringVarP(&c.policy, "policy", "", "", "Policy")
	persistentFlags.BoolVarP(&c.initRequired, "init-required", "", false, "Init required")
	persistentFlags.StringVarP(&c.collectionsConfig, "collections-config", "", "", "Private data collections")

	cmd.MarkPersistentFlagRequired("user")
	cmd.MarkPersistentFlagRequired("mspid")
	cmd.MarkPersistentFlagRequired("config")
	cmd.MarkPersistentFlagRequired("channelName")
	cmd.MarkPersistentFlagRequired("name")
	return cmd
}
