package mainchannel

import (
	"context"
	"fmt"
	"github.com/kfsoftware/hlf-operator/api/hlf.kungfusoftware.es/v1alpha1"
	"github.com/kfsoftware/hlf-operator/controllers/utils"
	"github.com/kfsoftware/hlf-operator/kubectl-hlf/cmd/helpers"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"io"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
)

type Options struct {
	Name                         string
	Output                       bool
	Capabilities                 []string
	BatchTimeout                 string
	MaxMessageCount              int
	AbsoluteMaxBytes             int
	PreferredMaxBytes            int
	EtcdRaftTickInterval         string
	EtcdRaftElectionTick         int
	EtcdRaftHeartbeatTick        int
	EtcdRaftMaxInflightBlocks    int
	EtcdRaftSnapshotIntervalSize int
	ChannelName                  string
	AdminPeerOrgs                []string
	AdminOrdererOrgs             []string
	PeerOrgs                     []string
	OrdererOrgs                  []string
	Consenters                   []string
	SecretName                   string
	SecretNS                     string
	Identities                   []string
}

func (o Options) Validate() error {
	return nil
}

type createCmd struct {
	out         io.Writer
	errOut      io.Writer
	channelOpts Options
}

func (c *createCmd) validate() error {
	return c.channelOpts.Validate()
}
func (c *createCmd) run() error {
	oclient, err := helpers.GetKubeOperatorClient()
	if err != nil {
		return err
	}
	clientSet, err := helpers.GetKubeClient()
	if err != nil {
		return err
	}
	identities := map[string]v1alpha1.FabricMainChannelIdentity{}
	for _, identity := range c.channelOpts.Identities {
		chunks := strings.Split(identity, ";")
		if len(chunks) != 2 {
			return fmt.Errorf("invalid identity %s, example format <msp_id>;<secret_key>", identity)
		}
		mspID := chunks[0]
		secretKey := chunks[1]
		identities[mspID] = v1alpha1.FabricMainChannelIdentity{
			SecretName:      c.channelOpts.SecretName,
			SecretNamespace: c.channelOpts.SecretNS,
			SecretKey:       secretKey,
		}

	}
	ordererOrganizations := []v1alpha1.FabricMainChannelOrdererOrganization{}
	adminPeerOrganizations := []v1alpha1.FabricMainChannelAdminPeerOrganizationSpec{}
	peerOrganizations := []v1alpha1.FabricMainChannelPeerOrganization{}
	adminOrdererOrganizations := []v1alpha1.FabricMainChannelAdminOrdererOrganizationSpec{}
	externalOrdererOrganizations := []v1alpha1.FabricMainChannelExternalOrdererOrganization{}
	consenters := []v1alpha1.FabricMainChannelConsenter{}
	ns := ""
	orderers, err := helpers.GetClusterOrdererNodes(clientSet, oclient, ns)
	if err != nil {
		return err
	}
	ordererMap := map[string][]*helpers.ClusterOrdererNode{}
	for _, orderer := range orderers {
		if !utils.Contains(c.channelOpts.OrdererOrgs, orderer.Spec.MspID) {
			continue
		}
		tlsCert, err := utils.ParseX509Certificate([]byte(orderer.Status.TlsCert))
		if err != nil {
			return err
		}
		consenterHost, consenterPort, err := helpers.GetOrdererHostAndPort(
			clientSet,
			orderer.Spec,
			orderer.Status,
		)
		if err != nil {
			return err
		}
		consenters = append(consenters, v1alpha1.FabricMainChannelConsenter{
			Host:    consenterHost,
			Port:    consenterPort,
			TLSCert: string(utils.EncodeX509Certificate(tlsCert)),
		})
		_, ok := ordererMap[orderer.Spec.MspID]
		if !ok {
			ordererMap[orderer.Spec.MspID] = []*helpers.ClusterOrdererNode{}
		}
		ordererMap[orderer.Spec.MspID] = append(ordererMap[orderer.Spec.MspID], orderer)
	}
	for mspID, nodes := range ordererMap {
		node := nodes[0]
		var ordererEndpoints []string
		for _, ordererNode := range nodes {
			ordererHost, ordererPort, err := helpers.GetOrdererHostAndPort(clientSet, ordererNode.Spec, ordererNode.Status)
			if err != nil {
				return err
			}
			ordererEndpoints = append(ordererEndpoints, fmt.Sprintf("%s:%d", ordererHost, ordererPort))
		}
		tlsCACert := node.Status.TlsCACert
		signCACert := node.Status.SignCACert
		ordererNodes := []v1alpha1.FabricMainChannelExternalOrdererNode{}
		for _, ordererNode := range nodes {
			adminOrdererHost, adminOrdererPort, err := helpers.GetOrdererAdminHostAndPort(clientSet, ordererNode.Spec, ordererNode.Status)
			if err != nil {
				return err
			}
			ordererNodes = append(ordererNodes, v1alpha1.FabricMainChannelExternalOrdererNode{
				Host:      adminOrdererHost,
				AdminPort: adminOrdererPort,
			})
		}
		ordererOrganizations = append(ordererOrganizations, v1alpha1.FabricMainChannelOrdererOrganization{
			MSPID:                  mspID,
			TLSCACert:              tlsCACert,
			SignCACert:             signCACert,
			OrdererEndpoints:       ordererEndpoints,
			OrderersToJoin:         []v1alpha1.FabricMainChannelOrdererNode{},
			ExternalOrderersToJoin: ordererNodes,
		})
	}
	for _, adminPeerOrg := range c.channelOpts.AdminPeerOrgs {
		adminPeerOrganizations = append(adminPeerOrganizations, v1alpha1.FabricMainChannelAdminPeerOrganizationSpec{
			MSPID: adminPeerOrg,
		})
	}
	externalPeerOrganizations := []v1alpha1.FabricMainChannelExternalPeerOrganization{}
	peerOrgs, _, err := helpers.GetClusterPeers(clientSet, oclient, ns)
	if err != nil {
		return err
	}
	for _, peerOrg := range peerOrgs {
		if len(peerOrg.Peers) == 0 {
			return fmt.Errorf("no peers found for organization %s", peerOrg.MspID)
		}
		if !utils.Contains(c.channelOpts.PeerOrgs, peerOrg.MspID) {
			continue
		}
		firstPeer := peerOrg.Peers[0]
		externalPeerOrganizations = append(externalPeerOrganizations, v1alpha1.FabricMainChannelExternalPeerOrganization{
			MSPID:        peerOrg.MspID,
			TLSRootCert:  firstPeer.Status.TlsCACert,
			SignRootCert: firstPeer.Status.SignCACert,
		})
	}
	for _, adminOrdererOrgMSPID := range c.channelOpts.AdminOrdererOrgs {
		adminOrdererOrganizations = append(adminOrdererOrganizations, v1alpha1.FabricMainChannelAdminOrdererOrganizationSpec{
			MSPID: adminOrdererOrgMSPID,
		})
	}

	channelConfig := &v1alpha1.FabricMainChannelConfig{
		Application: &v1alpha1.FabricMainChannelApplicationConfig{
			Capabilities: c.channelOpts.Capabilities,
			Policies:     nil,
			ACLs:         nil,
		},
		Orderer: &v1alpha1.FabricMainChannelOrdererConfig{
			OrdererType:  "etcdraft",
			Capabilities: c.channelOpts.Capabilities,
			Policies:     nil,
			BatchTimeout: c.channelOpts.BatchTimeout,
			BatchSize: &v1alpha1.FabricMainChannelOrdererBatchSize{
				MaxMessageCount:   c.channelOpts.MaxMessageCount,
				AbsoluteMaxBytes:  c.channelOpts.AbsoluteMaxBytes,
				PreferredMaxBytes: c.channelOpts.PreferredMaxBytes,
			},
			State: "STATE_NORMAL",
			EtcdRaft: &v1alpha1.FabricMainChannelEtcdRaft{
				Options: &v1alpha1.FabricMainChannelEtcdRaftOptions{
					TickInterval:         c.channelOpts.EtcdRaftTickInterval,
					ElectionTick:         uint32(c.channelOpts.EtcdRaftElectionTick),
					HeartbeatTick:        uint32(c.channelOpts.EtcdRaftHeartbeatTick),
					MaxInflightBlocks:    uint32(c.channelOpts.EtcdRaftMaxInflightBlocks),
					SnapshotIntervalSize: uint32(c.channelOpts.EtcdRaftSnapshotIntervalSize),
				},
			},
		},
		Capabilities: c.channelOpts.Capabilities,
		Policies:     nil,
	}
	fabricMainChannel := &v1alpha1.FabricMainChannel{
		TypeMeta: v1.TypeMeta{
			Kind:       "FabricMainChannel",
			APIVersion: v1alpha1.GroupVersion.String(),
		},
		ObjectMeta: v1.ObjectMeta{
			Name: c.channelOpts.Name,
		},
		Spec: v1alpha1.FabricMainChannelSpec{
			Name:                         c.channelOpts.ChannelName,
			Identities:                   identities,
			AdminPeerOrganizations:       adminPeerOrganizations,
			PeerOrganizations:            peerOrganizations,
			ExternalPeerOrganizations:    externalPeerOrganizations,
			ChannelConfig:                channelConfig,
			AdminOrdererOrganizations:    adminOrdererOrganizations,
			OrdererOrganizations:         ordererOrganizations,
			ExternalOrdererOrganizations: externalOrdererOrganizations,
			Consenters:                   consenters,
		},
	}
	if c.channelOpts.Output {
		ot, err := helpers.MarshallWithoutStatus(&fabricMainChannel)
		if err != nil {
			return err
		}
		fmt.Println(string(ot))
	} else {
		ctx := context.Background()
		_, err = oclient.HlfV1alpha1().FabricMainChannels().Create(
			ctx,
			fabricMainChannel,
			v1.CreateOptions{},
		)
		if err != nil {
			return err
		}
		log.Infof("MainChannel %s created on namespace %s", fabricMainChannel.Name, fabricMainChannel.Namespace)
	}
	return nil
}
func newCreateMainChannelCmd(out io.Writer, errOut io.Writer) *cobra.Command {
	c := createCmd{out: out, errOut: errOut}
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a main channel",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := c.validate(); err != nil {
				return err
			}
			return c.run()
		},
	}
	f := cmd.Flags()
	f.StringVar(&c.channelOpts.Name, "name", "", "Name of the Fabric Console to create")
	f.StringSliceVar(&c.channelOpts.Capabilities, "capabilities", []string{"V2_0"}, "Capabilities of the channel")
	f.StringSliceVar(&c.channelOpts.AdminPeerOrgs, "admin-peer-orgs", []string{}, "MSP IDs of the admin peer organizations")
	f.StringSliceVar(&c.channelOpts.AdminOrdererOrgs, "admin-orderer-orgs", []string{}, "MSP IDs of the admin orderer organizations")
	f.StringSliceVar(&c.channelOpts.OrdererOrgs, "orderer-orgs", []string{}, "MSP IDs of the orderer organizations")
	f.StringSliceVar(&c.channelOpts.PeerOrgs, "peer-orgs", []string{}, "MSP IDs of the peer organizations")

	f.StringVar(&c.channelOpts.ChannelName, "channel-name", "mychannel", "Name of the channel")
	f.StringVar(&c.channelOpts.BatchTimeout, "batch-timeout", "2s", "Batch timeout")
	f.IntVar(&c.channelOpts.MaxMessageCount, "max-message-count", 10, "Max message count")
	f.IntVar(&c.channelOpts.AbsoluteMaxBytes, "absolute-max-bytes", 1048576, "Absolute max bytes")
	f.IntVar(&c.channelOpts.PreferredMaxBytes, "preferred-max-bytes", 524288, "Preferred max bytes")
	f.StringVar(&c.channelOpts.EtcdRaftTickInterval, "etcd-raft-tick-interval", "500ms", "Etcd raft tick interval")
	f.IntVar(&c.channelOpts.EtcdRaftElectionTick, "etcd-raft-election-tick", 10, "Etcd raft election tick")
	f.IntVar(&c.channelOpts.EtcdRaftHeartbeatTick, "etcd-raft-heartbeat-tick", 1, "Etcd raft heartbeat tick")
	f.IntVar(&c.channelOpts.EtcdRaftMaxInflightBlocks, "etcd-raft-max-inflight-blocks", 5, "Etcd raft max inflight blocks")
	f.IntVar(&c.channelOpts.EtcdRaftSnapshotIntervalSize, "etcd-raft-snapshot-interval-size", 16777216, "Etcd raft snapshot interval size")

	f.StringVar(&c.channelOpts.SecretName, "secret-name", "", "Secret name for the identities")
	f.StringVar(&c.channelOpts.SecretNS, "secret-ns", "", "Secret namespace for the identities")
	f.StringSliceVar(&c.channelOpts.Identities, "identities", []string{}, "Identity map")

	f.BoolVarP(&c.channelOpts.Output, "output", "o", false, "Output in yaml")
	return cmd
}
