package chaincode

import (
	"github.com/hyperledger/fabric-protos-go/common"
	pb "github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/common/policydsl"
	"github.com/kfsoftware/hlf-operator/kubectl-hlf/cmd/helpers"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"io"
	"io/ioutil"
)

type queryCheckCommitReadiness struct {
	configPath        string
	peer              string
	userName          string
	channelName       string
	chaincodeName     string
	sequence          int64
	policy            string
	initRequired      bool
	collectionsConfig string
	version           string
}

func (c *queryCheckCommitReadiness) validate() error {
	return nil
}
func (c *queryCheckCommitReadiness) run(out io.Writer) error {
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
	var sp *common.SignaturePolicyEnvelope
	if c.policy != "" {
		sp, err = policydsl.FromString(c.policy)
		if err != nil {
			return err
		}
	}
	if len(collectionConfigs) == 0 {
		collectionConfigs = nil
	}
	chaincode, err := resClient.LifecycleCheckCCCommitReadiness(
		c.channelName,
		resmgmt.LifecycleCheckCCCommitReadinessRequest{
			Name:              c.chaincodeName,
			Version:           c.version,
			Sequence:          c.sequence,
			EndorsementPlugin: "escc",
			ValidationPlugin:  "vscc",
			SignaturePolicy:   sp,
			CollectionConfig:  collectionConfigs,
			InitRequired:      c.initRequired,
		},
		resmgmt.WithTargetEndpoints(peerName),
	)
	data := [][]string{}
	for mspID, approved := range chaincode.Approvals {
		isApproved := "false"
		if approved {
			isApproved = "true"
		}
		data = append(data, []string{mspID, isApproved})
	}
	table := tablewriter.NewWriter(out)
	table.SetHeader([]string{"MSP ID", "Approved"})
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetTablePadding("\t") // pad with tabs
	table.SetNoWhiteSpace(true)
	table.AppendBulk(data) // Add Bulk Data
	table.Render()
	return nil
}
func newCheckCommitReadiness(out io.Writer, errOut io.Writer) *cobra.Command {
	c := &queryCheckCommitReadiness{}
	cmd := &cobra.Command{
		Use: "checkcommitreadiness",
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
	persistentFlags.StringVarP(&c.channelName, "channel", "C", "", "Channel name")
	persistentFlags.StringVarP(&c.chaincodeName, "chaincode", "c", "", "Chaincode label")
	persistentFlags.StringVarP(&c.version, "version", "", "1.0", "Version")
	persistentFlags.Int64VarP(&c.sequence, "sequence", "", 1, "Sequence number")
	persistentFlags.StringVarP(&c.policy, "policy", "", "", "Policy")
	persistentFlags.BoolVarP(&c.initRequired, "init-required", "", false, "Init required")
	persistentFlags.StringVarP(&c.collectionsConfig, "collections-config", "", "", "Private data collections")

	cmd.MarkPersistentFlagRequired("user")
	cmd.MarkPersistentFlagRequired("peer")
	cmd.MarkPersistentFlagRequired("config")
	cmd.MarkPersistentFlagRequired("channel")

	return cmd
}
