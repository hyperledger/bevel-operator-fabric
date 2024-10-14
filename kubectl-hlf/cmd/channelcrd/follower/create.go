package follower

import (
	"context"
	"fmt"
	"github.com/kfsoftware/hlf-operator/kubectl-hlf/cmd/helpers"
	"github.com/kfsoftware/hlf-operator/pkg/apis/hlf.kungfusoftware.es/v1alpha1"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"io"
	"io/ioutil"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net"
	"strconv"
	"strings"
)

type CreateOptions struct {
	Name                string
	Output              bool
	MSPID               string
	AnchorPeers         []string
	SecretName          string
	SecretNamespace     string
	SecretKey           string
	ChannelName         string
	Peers               []string
	OrdererCertificates []string
	OrdererURLs         []string
}

func (o CreateOptions) Validate() error {
	if o.Name == "" {
		return fmt.Errorf("--name is required")
	}
	if o.MSPID == "" {
		return fmt.Errorf("--msp-id is required")
	}
	if o.SecretName == "" {
		return fmt.Errorf("--secret-name is required")
	}
	if o.SecretNamespace == "" {
		return fmt.Errorf("--secret-namespace is required")
	}
	if o.SecretKey == "" {
		return fmt.Errorf("--secret-key is required")
	}
	if len(o.AnchorPeers) == 0 {
		return fmt.Errorf("--anchor-peers is required")
	}
	if len(o.OrdererURLs) == 0 {
		return fmt.Errorf("--orderer-urls is required")
	}
	if len(o.OrdererCertificates) == 0 {
		return fmt.Errorf("--orderer-certificates is required")
	}
	if len(o.Peers) == 0 {
		return fmt.Errorf("--peers is required")
	}
	return nil
}

type createCmd struct {
	out         io.Writer
	errOut      io.Writer
	channelOpts CreateOptions
}

func (c *createCmd) validate() error {
	return c.channelOpts.Validate()
}
func (c *createCmd) run() error {
	oclient, err := helpers.GetKubeOperatorClient()
	if err != nil {
		return err
	}
	orderers := []v1alpha1.FabricFollowerChannelOrderer{}
	for idx, orderer := range c.channelOpts.OrdererURLs {
		if len(c.channelOpts.OrdererCertificates)-1 < idx {
			return fmt.Errorf("orderer certificate not found for orderer %s", orderer)
		}
		ordererCrtFile := c.channelOpts.OrdererCertificates[idx]
		ordererCertificate, err := ioutil.ReadFile(ordererCrtFile)
		if err != nil {
			return fmt.Errorf("error reading orderer certificate file %s: %s", ordererCrtFile, err)
		}
		orderers = append(orderers, v1alpha1.FabricFollowerChannelOrderer{
			URL:         orderer,
			Certificate: string(ordererCertificate),
		})
	}
	peers := []v1alpha1.FabricFollowerChannelPeer{}
	for _, peer := range c.channelOpts.Peers {
		chunks := strings.Split(peer, ".")
		if len(chunks) != 2 {
			return fmt.Errorf("invalid peer format: %s", peer)
		}
		name := chunks[0]
		namespace := chunks[1]
		fabricPeer, err := oclient.HlfV1alpha1().FabricPeers(namespace).Get(context.TODO(), name, v1.GetOptions{})
		if err != nil {
			return err
		}
		peers = append(peers, v1alpha1.FabricFollowerChannelPeer{
			Name:      fabricPeer.Name,
			Namespace: fabricPeer.Namespace,
		})
	}
	anchorPeers := []v1alpha1.FabricFollowerChannelAnchorPeer{}
	for _, anchorPeer := range c.channelOpts.AnchorPeers {
		host, port, err := net.SplitHostPort(anchorPeer)
		if err != nil {
			return err
		}
		portNumber, err := strconv.Atoi(port)
		if err != nil {
			return err
		}
		anchorPeers = append(anchorPeers, v1alpha1.FabricFollowerChannelAnchorPeer{
			Host: host,
			Port: portNumber,
		})
	}
	identity := v1alpha1.HLFIdentity{
		SecretName:      c.channelOpts.SecretName,
		SecretNamespace: c.channelOpts.SecretNamespace,
		SecretKey:       c.channelOpts.SecretKey,
	}
	fabricFollowerChannel := &v1alpha1.FabricFollowerChannel{
		TypeMeta: v1.TypeMeta{
			Kind:       "FabricFollowerChannel",
			APIVersion: v1alpha1.GroupVersion.String(),
		},
		ObjectMeta: v1.ObjectMeta{
			Name: c.channelOpts.Name,
		},
		Spec: v1alpha1.FabricFollowerChannelSpec{
			Name:                c.channelOpts.ChannelName,
			MSPID:               c.channelOpts.MSPID,
			Orderers:            orderers,
			PeersToJoin:         peers,
			ExternalPeersToJoin: []v1alpha1.FabricFollowerChannelExternalPeer{},
			AnchorPeers:         anchorPeers,
			HLFIdentity:         identity,
		},
	}
	if c.channelOpts.Output {
		ot, err := helpers.MarshallWithoutStatus(&fabricFollowerChannel)
		if err != nil {
			return err
		}
		fmt.Println(string(ot))
	} else {
		ctx := context.Background()
		_, err = oclient.HlfV1alpha1().FabricFollowerChannels().Create(
			ctx,
			fabricFollowerChannel,
			v1.CreateOptions{},
		)
		if err != nil {
			return err
		}
		log.Infof("Channel %s created on namespace %s", fabricFollowerChannel.Name, fabricFollowerChannel.Namespace)
	}
	return nil
}

func newCreateFollowerChannelCmd(out io.Writer, errOut io.Writer) *cobra.Command {
	c := createCmd{out: out, errOut: errOut}
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a follower channel",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := c.validate(); err != nil {
				return err
			}
			return c.run()
		},
	}
	f := cmd.Flags()
	f.StringVar(&c.channelOpts.Name, "name", "", "Name of the Fabric Console to create")
	f.StringVar(&c.channelOpts.ChannelName, "channel-name", "", "Name of the channel to join")
	f.StringVar(&c.channelOpts.MSPID, "mspid", "", "MSPID of the channel")
	f.StringArrayVar(&c.channelOpts.AnchorPeers, "anchor-peers", []string{}, "Anchor peers of the channel")
	f.StringArrayVar(&c.channelOpts.OrdererURLs, "orderer-urls", []string{}, "Orderer URLs of the channel, e.g grpcs://<host>:<port>")
	f.StringArrayVar(&c.channelOpts.OrdererCertificates, "orderer-certificates", []string{}, "Orderer certificates of the channel")
	f.StringArrayVar(&c.channelOpts.Peers, "peers", []string{}, "Peers of the channel")
	f.StringVar(&c.channelOpts.SecretName, "secret-name", "", "Name of the secret containing the certificate to join and set the anchor peers")
	f.StringVar(&c.channelOpts.SecretNamespace, "secret-ns", "", "Namespace of the secret containing the certificate to join and set the anchor peers")
	f.StringVar(&c.channelOpts.SecretKey, "secret-key", "", "Key of the secret containing the certificate to join and set the anchor peers")
	f.BoolVarP(&c.channelOpts.Output, "output", "o", false, "Output in yaml")
	return cmd
}
