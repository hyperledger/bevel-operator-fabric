package chaincode

import (
	"crypto/sha256"
	"fmt"
	pb "github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/ccpackager/lifecycle"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"io"
	"io/ioutil"
	"strings"
)

type calculatePackageIDCMD struct {
	chaincodeLanguage string
	chaincodePath     string
	chaincodeLabel    string
}

func (c *calculatePackageIDCMD) validate() error {
	return nil
}
func (c *calculatePackageIDCMD) run(stdOut io.Writer, stdErr io.Writer) error {
	chLng, ok := pb.ChaincodeSpec_Type_value[strings.ToUpper(c.chaincodeLanguage)]
	if !ok {
		return errors.Errorf("Language %s not valid", c.chaincodeLanguage)
	}
	var pkg []byte
	var err error
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
	sum := sha256.Sum256(pkg)
	packageID := fmt.Sprintf("%s:%x", c.chaincodeLabel, sum)
	stdOut.Write([]byte(packageID))
	return nil
}
func newCalculatePackageIDCMD(stdOut io.Writer, stdErr io.Writer) *cobra.Command {
	c := &calculatePackageIDCMD{}
	cmd := &cobra.Command{
		Use: "calculatepackageid",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := c.validate(); err != nil {
				return err
			}
			return c.run(stdOut, stdErr)
		},
	}
	persistentFlags := cmd.PersistentFlags()
	persistentFlags.StringVarP(&c.chaincodeLanguage, "language", "l", "", "Chaincode language")
	persistentFlags.StringVarP(&c.chaincodePath, "path", "", "", "Chaincode path")
	persistentFlags.StringVarP(&c.chaincodeLabel, "label", "", "", "Chaincode label")
	cmd.MarkPersistentFlagRequired("path")
	cmd.MarkPersistentFlagRequired("label")
	cmd.MarkPersistentFlagRequired("language")
	return cmd
}
