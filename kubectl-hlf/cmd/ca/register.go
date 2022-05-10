package ca

import (
	"io"

	"github.com/kfsoftware/hlf-operator/controllers/certs"
	"github.com/kfsoftware/hlf-operator/kubectl-hlf/cmd/helpers"
	"github.com/spf13/cobra"
)

type RegisterOptions struct {
	Name         string
	NS           string
	User         string
	Secret       string
	Type         string
	MspID        string
	EnrollID     string
	EnrollSecret string
}

func (o RegisterOptions) Validate() error {
	return nil
}

type registerCmd struct {
	out    io.Writer
	errOut io.Writer
	caOpts RegisterOptions
}

func (c *registerCmd) validate() error {
	return c.caOpts.Validate()
}
func (c *registerCmd) run(args []string) error {
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

	url, err := helpers.GetURLForCA(certAuth)
	if err != nil {
		return err
	}
	_, err = certs.RegisterUser(certs.RegisterUserRequest{
		TLSCert:      certAuth.Status.TlsCert,
		URL:          url,
		Name:         "",
		MSPID:        c.caOpts.MspID,
		EnrollID:     c.caOpts.EnrollID,
		EnrollSecret: c.caOpts.EnrollSecret,
		User:         c.caOpts.User,
		Secret:       c.caOpts.Secret,
		Type:         c.caOpts.Type,
		Attributes:   nil,
	})
	if err != nil {
		return err
	}
	return nil
}
func newCARegisterCmd(out io.Writer, errOut io.Writer) *cobra.Command {
	c := registerCmd{out: out, errOut: errOut}
	cmd := &cobra.Command{
		Use:   "register",
		Short: "Create a Fabric Certificate authority",
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
	f.StringVarP(&c.caOpts.EnrollID, "enroll-id", "", "", "Enroll ID to register new users")
	f.StringVarP(&c.caOpts.EnrollSecret, "enroll-secret", "", "", "Enroll secret to register new users")
	f.StringVarP(&c.caOpts.User, "user", "", "", "Username for the new user")
	f.StringVarP(&c.caOpts.Secret, "secret", "", "", "Password for the new user")
	f.StringVarP(&c.caOpts.Type, "type", "", "", "Type of the identity to create (peer/client/orderer/admin)")
	f.StringVarP(&c.caOpts.MspID, "mspid", "", "", "MSP ID of the organization")

	return cmd
}
