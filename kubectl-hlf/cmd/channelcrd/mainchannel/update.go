package mainchannel

import (
	"context"
	"fmt"
	"github.com/kfsoftware/hlf-operator/kubectl-hlf/cmd/helpers"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"io"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type updateCmd struct {
	out         io.Writer
	errOut      io.Writer
	channelOpts Options
}

func (c *updateCmd) validate() error {
	return c.channelOpts.Validate()
}
func (c *updateCmd) run() error {
	oclient, err := helpers.GetKubeOperatorClient()
	if err != nil {
		return err
	}
	fabricMainChannel, err := oclient.HlfV1alpha1().FabricMainChannels().Get(context.TODO(), c.channelOpts.Name, v1.GetOptions{})
	if err != nil {
		return err
	}
	fabricMainChannelSpec, err := c.channelOpts.mapToFabricMainChannel()
	if err != nil {
		return err
	}
	fabricMainChannel.Spec = *fabricMainChannelSpec
	if c.channelOpts.Output {
		ot, err := helpers.MarshallWithoutStatus(&fabricMainChannel)
		if err != nil {
			return err
		}
		fmt.Println(string(ot))
	} else {
		ctx := context.Background()
		_, err = oclient.HlfV1alpha1().FabricMainChannels().Update(
			ctx,
			fabricMainChannel,
			v1.UpdateOptions{},
		)
		if err != nil {
			return err
		}
		log.Infof("MainChannel %s updated on namespace %s", fabricMainChannel.Name, fabricMainChannel.Namespace)
	}
	return nil
}
func newUpdateMainChannelCmd(out io.Writer, errOut io.Writer) *cobra.Command {
	c := updateCmd{out: out, errOut: errOut}
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update a main channel",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := c.validate(); err != nil {
				return err
			}
			return c.run()
		},
	}
	f := cmd.Flags()
	f.StringVar(&c.channelOpts.Name, "name", "", "Name of the Fabric Console to update")
	f.StringSliceVar(&c.channelOpts.Capabilities, "capabilities", []string{"V2_0"}, "Capabilities of the channel")
	f.StringSliceVar(&c.channelOpts.AdminPeerOrgs, "admin-peer-orgs", []string{}, "MSP IDs of the admin peer organizations")
	f.StringSliceVar(&c.channelOpts.AdminOrdererOrgs, "admin-orderer-orgs", []string{}, "MSP IDs of the admin orderer organizations")
	f.StringSliceVar(&c.channelOpts.OrdererOrgs, "orderer-orgs", []string{}, "MSP IDs of the orderer organizations")
	f.StringSliceVar(&c.channelOpts.PeerOrgs, "peer-orgs", []string{}, "MSP IDs of the peer organizations")
	f.StringSliceVar(&c.channelOpts.Consenters, "consenters", []string{}, "Consenters of the channel")
	f.StringSliceVar(&c.channelOpts.ConsenterCertificates, "consenter-certificates", []string{}, "Consenter certificates of the channel")

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
