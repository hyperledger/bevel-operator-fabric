package org

import (
	"github.com/kfsoftware/hlf-operator/kubectl-hlf/cmd/helpers"
	"github.com/spf13/cobra"
	"io"
	"io/ioutil"
	"os"
	"path"
)

type InspectOptions struct {
	Orgs       []string
	OutputPath string
}

func (o InspectOptions) Validate() error {
	return nil
}

type inspectCmd struct {
	out    io.Writer
	errOut io.Writer
	caOpts InspectOptions
}

func (c *inspectCmd) validate() error {
	return c.caOpts.Validate()
}
func (c *inspectCmd) run(args []string) error {
	oclient, err := helpers.GetKubeOperatorClient()
	if err != nil {
		return err
	}
	baseOutputPath := c.caOpts.OutputPath
	peerOrgs, _, err := helpers.GetClusterPeers(oclient, "")
	if err != nil {
		return err
	}
	for _, peerOrg := range peerOrgs {
		firstPeer := peerOrg.Peers[0]
		certAuth, err := helpers.GetCertAuthByURL(
			oclient,
			firstPeer.Spec.Secret.Enrollment.Component.Cahost,
			firstPeer.Spec.Secret.Enrollment.Component.Caport,
		)
		if err != nil {
			return err
		}
		orgPath := path.Join(baseOutputPath, "peerOrganizations", peerOrg.MspID)
		mspPath := path.Join(orgPath, "msp")
		mspCaCerts := path.Join(mspPath, "cacerts")
		mspTLSCaCerts := path.Join(mspPath, "tlscacerts")

		err = os.MkdirAll(mspCaCerts, os.ModePerm)
		if err != nil {
			return err
		}
		err = os.MkdirAll(mspTLSCaCerts, os.ModePerm)
		if err != nil {
			return err
		}
		mspCACertPath := path.Join(mspCaCerts, "ca.pem")
		err = ioutil.WriteFile(mspCACertPath, []byte(certAuth.Status.CACert), os.ModePerm)
		if err != nil {
			return err
		}
		mspTLSCACertPath := path.Join(mspTLSCaCerts, "tlsca.pem")
		err = ioutil.WriteFile(mspTLSCACertPath, []byte(certAuth.Status.TLSCACert), os.ModePerm)
		if err != nil {
			return err
		}
		nodeOusContent := `
NodeOUs:
  Enable: true
  ClientOUIdentifier:
    Certificate: cacerts/ca.pem
    OrganizationalUnitIdentifier: client
  PeerOUIdentifier:
    Certificate: cacerts/ca.pem
    OrganizationalUnitIdentifier: peer
  AdminOUIdentifier:
    Certificate: cacerts/ca.pem
    OrganizationalUnitIdentifier: admin
  OrdererOUIdentifier:
    Certificate: cacerts/ca.pem
    OrganizationalUnitIdentifier: orderer
`
		nodeOusPath := path.Join(mspPath, "config.yaml")
		err = ioutil.WriteFile(nodeOusPath, []byte(nodeOusContent), os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}
func newOrgInspectCmd(out io.Writer, errOut io.Writer) *cobra.Command {
	c := inspectCmd{out: out, errOut: errOut}
	cmd := &cobra.Command{
		Use:   "inspect",
		Short: "Inpects and dumps the crypto material for the organization",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := c.validate(); err != nil {
				return err
			}
			return c.run(args)
		},
	}
	f := cmd.Flags()
	f.StringSliceVarP(&c.caOpts.Orgs, "orgs", "o", []string{}, "Organizations to inspect")
	f.StringVarP(&c.caOpts.OutputPath, "output-path", "", "", "Output path")

	return cmd
}
