package channel

import (
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/kfsoftware/hlf-operator/kubectl-hlf/cmd/helpers"
	"github.com/olekukonko/tablewriter"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"
)

type topChannelCmd struct {
	configPath  string
	peer        string
	channelName string
	userName    string
}

func (c *topChannelCmd) validate() error {
	return nil
}
func (c *topChannelCmd) run() error {
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
	configBackend := config.FromFile(c.configPath)
	sdk, err := fabsdk.New(configBackend)
	if err != nil {
		return err
	}
	org1ChannelContext := sdk.ChannelContext(
		c.channelName,
		fabsdk.WithUser(c.userName),
		fabsdk.WithOrg(mspID),
	)
	chCtx, err := org1ChannelContext()
	if err != nil {
		return err
	}
	discovery, err := chCtx.ChannelService().Discovery()
	if err != nil {
		return err
	}
	for {
		data := [][]string{}
		peers, err := discovery.GetPeers()
		if err != nil {
			log.Printf("Failed to get peers %v", err)
			return err
		}

		for _, peer := range peers {
			props := peer.Properties()
			ledgerHeight := props[fab.PropertyLedgerHeight]
			data = append(data, []string{
				peer.URL(), peer.MSPID(), fmt.Sprintf("L=%d", ledgerHeight),
			})
		}
		tableString := &strings.Builder{}
		table := tablewriter.NewWriter(tableString)
		table.SetHeader([]string{"URL", "MSP ID", "Height"})
		table.SetAutoWrapText(false)
		table.SetAutoFormatHeaders(true)
		table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
		table.SetAlignment(tablewriter.ALIGN_LEFT)
		table.SetCenterSeparator("")
		table.SetColumnSeparator("")
		table.SetRowSeparator("")
		table.SetTablePadding("\t")
		table.SetNoWhiteSpace(true)
		table.AppendBulk(data) // Add Bulk Data
		table.Render()
		cmd := exec.Command("clear") //Linux example, its tested
		cmd.Stdout = os.Stdout
		err = cmd.Run()
		if err != nil {
			return err
		}
		fmt.Printf("\r%s", tableString.String())
		time.Sleep(2 * time.Second)
	}
}
func newTopChannelCMD(io.Writer, io.Writer) *cobra.Command {
	c := &topChannelCmd{}
	cmd := &cobra.Command{
		Use: "top",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := c.validate(); err != nil {
				return err
			}
			return c.run()
		},
	}
	persistentFlags := cmd.PersistentFlags()
	persistentFlags.StringVarP(&c.peer, "peer", "p", "", "Name of the peer to invoke the updates")
	persistentFlags.StringVarP(&c.userName, "user", "", "", "User name for the transaction")
	persistentFlags.StringVarP(&c.channelName, "channel", "", "", "Channel name")
	persistentFlags.StringVarP(&c.configPath, "config", "", "", "Configuration file for the SDK")
	cmd.MarkPersistentFlagRequired("channel")
	cmd.MarkPersistentFlagRequired("user")
	cmd.MarkPersistentFlagRequired("peer")
	cmd.MarkPersistentFlagRequired("config")
	return cmd
}
