package mainchannel

import (
	"context"
	"fmt"
	"github.com/kfsoftware/hlf-operator/api/hlf.kungfusoftware.es/v1alpha1"
	"github.com/kfsoftware/hlf-operator/kubectl-hlf/cmd/helpers"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"io"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	EtcdRaftElectionTick         uint32
	EtcdRaftHeartbeatTick        uint32
	EtcdRaftMaxInflightBlocks    uint32
	EtcdRaftSnapshotIntervalSize uint32
	ChannelName                  string
	AdminPeerOrgs                []string
	AdminOrdererOrgs             []string
	PeerOrgs                     []string
	OrdererOrganizations         []string
	Consenters                   interface{}
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
	identities := map[string]v1alpha1.FabricMainChannelIdentity{}
	ordererOrganizations := []v1alpha1.FabricMainChannelOrdererOrganization{}
	adminPeerOrganizations := []v1alpha1.FabricMainChannelAdminPeerOrganizationSpec{}
	peerOrganizations := []v1alpha1.FabricMainChannelPeerOrganization{}
	adminOrdererOrganizations := []v1alpha1.FabricMainChannelAdminOrdererOrganizationSpec{}
	externalOrdererOrganizations := []v1alpha1.FabricMainChannelExternalOrdererOrganization{}
	consenters := []v1alpha1.FabricMainChannelConsenter{}
	//for _, orderer := range c.channelOpts.OrdererOrganizations {
	//	ordererOrganizations = append(ordererOrganizations, v1alpha1.FabricMainChannelOrdererOrganization{
	//		MSPID:                  "",
	//		CAName:                 "",
	//		CANamespace:            "",
	//		OrdererEndpoints:       nil,
	//		OrderersToJoin:         nil,
	//		ExternalOrderersToJoin: nil,
	//	})
	//}
	//for _, adminPeerOrg := range c.channelOpts.AdminPeerOrgs {
	//	adminPeerOrganizations = append(adminPeerOrganizations, v1alpha1.FabricMainChannelAdminPeerOrganizationSpec{
	//		MSPID: adminPeerOrg,
	//	})
	//}
	//for _, peerOrg := range c.channelOpts.PeerOrgs {
	//	peerOrganizations = append(peerOrganizations, v1alpha1.FabricMainChannelPeerOrganization{
	//		MSPID:       "",
	//		CAName:      "",
	//		CANamespace: "",
	//	})
	//}
	//for _, adminOrdererOrgMSPID := range c.channelOpts.AdminOrdererOrgs {
	//	adminOrdererOrganizations = append(adminOrdererOrganizations, v1alpha1.FabricMainChannelAdminOrdererOrganizationSpec{
	//		MSPID: adminOrdererOrgMSPID,
	//	})
	//}
	//
	//for _, consenter := range c.channelOpts.Consenters {
	//	consenters = append(consenters, v1alpha1.FabricMainChannelConsenter{
	//		Host:    "",
	//		Port:    0,
	//		TLSCert: "",
	//	})
	//}
	externalPeerOrganizations := []v1alpha1.FabricMainChannelExternalPeerOrganization{}
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
					ElectionTick:         c.channelOpts.EtcdRaftElectionTick,
					HeartbeatTick:        c.channelOpts.EtcdRaftHeartbeatTick,
					MaxInflightBlocks:    c.channelOpts.EtcdRaftMaxInflightBlocks,
					SnapshotIntervalSize: c.channelOpts.EtcdRaftSnapshotIntervalSize,
				},
			},
		},
		Capabilities: c.channelOpts.Capabilities,
		Policies:     nil,
	}
	fabricConsole := &v1alpha1.FabricMainChannel{
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
		ot, err := helpers.MarshallWithoutStatus(&fabricConsole)
		if err != nil {
			return err
		}
		fmt.Println(string(ot))
	} else {
		ctx := context.Background()
		_, err = oclient.HlfV1alpha1().FabricMainChannels().Create(
			ctx,
			fabricConsole,
			v1.CreateOptions{},
		)
		if err != nil {
			return err
		}
		log.Infof("Console %s created on namespace %s", fabricConsole.Name, fabricConsole.Namespace)
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
	f.BoolVarP(&c.channelOpts.Output, "output", "o", false, "Output in yaml")
	return cmd
}
