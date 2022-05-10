package chaincode

import (
	pb "github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/ccpackager/lifecycle"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/kfsoftware/hlf-operator/kubectl-hlf/cmd/helpers"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"io"
	"io/ioutil"
	"strings"
	"time"
)

type installChaincodeCmd struct {
	configPath        string
	peer              string
	chaincodeLanguage string
	chaincodePath     string
	userName          string
	chaincodeLabel    string
}

func (c *installChaincodeCmd) validate() error {
	return nil
}
func (c *installChaincodeCmd) run() error {
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
	chLng, ok := pb.ChaincodeSpec_Type_value[strings.ToUpper(c.chaincodeLanguage)]
	if !ok {
		return errors.Errorf("Language %s not valid", c.chaincodeLanguage)
	}
	var pkg []byte
	if strings.HasSuffix(c.chaincodePath, ".tar.gz") || strings.HasSuffix(c.chaincodePath, ".tgz") {
		pkg, err = ioutil.ReadFile(c.chaincodePath)
		if err != nil {
			return err
		}
	} else {
		pkg, err = lifecycle.NewCCPackage(&lifecycle.Descriptor{
			Path:  c.chaincodePath,
			Type:  pb.ChaincodeSpec_Type(chLng),
			Label: c.chaincodeLabel,
		})
		if err != nil {
			return err
		}
	}
	packageID := lifecycle.ComputePackageID(c.chaincodeLabel, pkg)
	responses, err := resClient.LifecycleInstallCC(
		resmgmt.LifecycleInstallCCRequest{
			Label:   c.chaincodeLabel,
			Package: pkg,
		},
		resmgmt.WithTargetEndpoints(peerName),
		resmgmt.WithTimeout(fab.ResMgmt, 20*time.Minute),
		resmgmt.WithTimeout(fab.PeerResponse, 20*time.Minute),
	)
	if err != nil {
		return err
	}
	for _, res := range responses {
		log.Infof("Package id=%s Status=%d", res.PackageID, res.Status)
	}
	log.Infof("Chaincode installed %s", packageID)
	return nil
}
func newChaincodeInstallCMD(io.Writer, io.Writer) *cobra.Command {
	c := &installChaincodeCmd{}
	cmd := &cobra.Command{
		Use: "install",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := c.validate(); err != nil {
				return err
			}
			return c.run()
		},
	}
	persistentFlags := cmd.PersistentFlags()
	persistentFlags.StringVarP(&c.peer, "peer", "p", "", "Peer org to invoke the updates")
	persistentFlags.StringVarP(&c.userName, "user", "", "", "User name for the transaction")
	persistentFlags.StringVarP(&c.chaincodeLanguage, "language", "l", "", "Chaincode language")
	persistentFlags.StringVarP(&c.chaincodePath, "path", "", "", "Chaincode path")
	persistentFlags.StringVarP(&c.chaincodeLabel, "label", "", "", "Chaincode label")
	persistentFlags.StringVarP(&c.configPath, "config", "", "", "Configuration file for the SDK")
	cmd.MarkPersistentFlagRequired("path")
	cmd.MarkPersistentFlagRequired("label")
	cmd.MarkPersistentFlagRequired("language")
	cmd.MarkPersistentFlagRequired("user")
	cmd.MarkPersistentFlagRequired("peer")
	cmd.MarkPersistentFlagRequired("config")
	return cmd
}
