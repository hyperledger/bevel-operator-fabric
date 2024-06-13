package ca

import (
	"github.com/kfsoftware/hlf-operator/controllers/certs"
	"github.com/kfsoftware/hlf-operator/internal/github.com/hyperledger/fabric-ca/api"
	"github.com/kfsoftware/hlf-operator/kubectl-hlf/cmd/helpers"
	"github.com/spf13/cobra"
	"io"
)

type RevokeOptions struct {
	Name         string
	NS           string
	MspID        string
	EnrollID     string
	EnrollSecret string
	CAURL        string

	RevName   string
	RevSerial string
	RevAKI    string
	RevReason string
	RevCAName string
	RevGenCRL bool
}

func (o RevokeOptions) Validate() error {
	return nil
}

type revokeCmd struct {
	out    io.Writer
	errOut io.Writer
	caOpts RevokeOptions
}

func (c *revokeCmd) validate() error {
	return c.caOpts.Validate()
}
func (c *revokeCmd) run(args []string) error {
	oclient, err := helpers.GetKubeOperatorClient()
	if err != nil {
		return err
	}
	clientSet, err := helpers.GetKubeClient()
	if err != nil {
		return err
	}
	certAuth, err := helpers.GetCertAuthByName(clientSet, oclient, c.caOpts.Name, c.caOpts.NS)
	if err != nil {
		return err
	}
	var url string
	if c.caOpts.CAURL != "" {
		url = c.caOpts.CAURL
	} else {
		url, err = helpers.GetURLForCA(certAuth)
		if err != nil {
			return err
		}
	}
	err = certs.RevokeUser(certs.RevokeUserRequest{
		TLSCert:      certAuth.Status.TlsCert,
		URL:          url,
		Name:         "",
		MSPID:        c.caOpts.MspID,
		EnrollID:     c.caOpts.EnrollID,
		EnrollSecret: c.caOpts.EnrollSecret,
		RevocationRequest: &api.RevocationRequest{
			Name:   c.caOpts.RevName,
			Serial: c.caOpts.RevSerial,
			AKI:    c.caOpts.RevAKI,
			Reason: c.caOpts.RevReason,
			CAName: c.caOpts.RevCAName,
			GenCRL: c.caOpts.RevGenCRL,
		},
	})
	if err != nil {
		return err
	}
	return nil
}
func newCARevokeCmd(out io.Writer, errOut io.Writer) *cobra.Command {
	c := revokeCmd{out: out, errOut: errOut}
	cmd := &cobra.Command{
		Use:   "revoke",
		Short: "Revoke a user from the Fabric CA",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := c.validate(); err != nil {
				return err
			}
			return c.run(args)
		},
	}
	f := cmd.Flags()
	f.StringVar(&c.caOpts.Name, "name", "", "Name of the Certificate Authority in the cluster, e.g ca.default")
	f.StringVarP(&c.caOpts.NS, "namespace", "n", helpers.DefaultNamespace, "Namespace scope for this request")
	f.StringVarP(&c.caOpts.EnrollID, "enroll-id", "", "", "Enroll ID to revoke new users")
	f.StringVarP(&c.caOpts.EnrollSecret, "enroll-secret", "", "", "Enroll secret to revoke new users")
	f.StringVarP(&c.caOpts.MspID, "mspid", "", "", "MSP ID of the organization")
	f.StringVarP(&c.caOpts.CAURL, "ca-url", "", "", "Fabric CA URL")

	f.StringVarP(&c.caOpts.RevName, "rev-name", "", "", "Name of the user to revoke")
	f.StringVarP(&c.caOpts.RevSerial, "rev-serial", "", "", "Serial number of the certificate to revoke")
	f.StringVarP(&c.caOpts.RevAKI, "rev-aki", "", "", "Authority Key Identifier of the certificate to revoke")
	f.StringVarP(&c.caOpts.RevReason, "rev-reason", "", "", "Reason for revocation")
	f.StringVarP(&c.caOpts.RevCAName, "rev-ca-name", "", "", "Name of the CA to revoke the user from")
	f.BoolVarP(&c.caOpts.RevGenCRL, "rev-gen-crl", "", false, "Generate CRL after revocation")
	return cmd
}
