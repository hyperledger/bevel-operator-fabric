package channel

import (
	"bytes"
	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric-protos-go/common"
	mspclient "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"

	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"io"
	"io/ioutil"
)

type signUpdateChannelCmd struct {
	configPath  string
	channelName string
	userName    string
	file        string
	output      string
	mspID       string
}

func (c *signUpdateChannelCmd) validate() error {
	return nil
}

func (c *signUpdateChannelCmd) run() error {
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
	updateEnvelopeBytes, err := ioutil.ReadFile(c.file)
	if err != nil {
		return err
	}
	envelope := &common.Envelope{}
	err = proto.Unmarshal(updateEnvelopeBytes, envelope)
	if err != nil {
		return err
	}
	configUpdateReader := bytes.NewReader(updateEnvelopeBytes)
	mspClient, err := mspclient.New(org1AdminClientContext, mspclient.WithOrg(c.mspID))
	if err != nil {
		return err
	}
	usr, err := mspClient.GetSigningIdentity(c.userName)
	if err != nil {
		return err
	}
	signature, err := resClient.CreateConfigSignatureFromReader(usr, configUpdateReader)
	if err != nil {
		return err
	}
	signatureBytes, err := proto.Marshal(signature)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(c.output, signatureBytes, 0777)
	if err != nil {
		return err
	}
	log.Infof("channel signed output: %s", c.output)
	return nil
}
func newSignUpdateChannelCMD(io.Writer, io.Writer) *cobra.Command {
	c := &signUpdateChannelCmd{}
	cmd := &cobra.Command{
		Use: "signupdate",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := c.validate(); err != nil {
				return err
			}
			return c.run()
		},
	}
	persistentFlags := cmd.PersistentFlags()
	persistentFlags.StringVarP(&c.mspID, "mspid", "", "", "Organization to use to submit the channel update")
	persistentFlags.StringVarP(&c.channelName, "channel", "", "", "Channel name")
	persistentFlags.StringVarP(&c.configPath, "config", "", "", "Configuration file for the SDK")
	persistentFlags.StringVarP(&c.userName, "user", "", "", "User name for the transaction")
	persistentFlags.StringVarP(&c.file, "file", "f", "", "Config update file")
	persistentFlags.StringVarP(&c.output, "output", "o", "", "Output signature file")
	cmd.MarkPersistentFlagRequired("mspid")
	cmd.MarkPersistentFlagRequired("channel")
	cmd.MarkPersistentFlagRequired("config")
	cmd.MarkPersistentFlagRequired("user")
	cmd.MarkPersistentFlagRequired("file")
	cmd.MarkPersistentFlagRequired("output")
	return cmd
}
