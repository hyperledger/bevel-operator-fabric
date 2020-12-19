package inspect

import (
	"bytes"
	"fmt"
	"github.com/Masterminds/sprig/v3"
	"github.com/kfsoftware/hlf-operator/controllers/utils"
	"github.com/kfsoftware/hlf-operator/kubectl-hlf/cmd/helpers"
	"github.com/spf13/cobra"
	"io"
	"io/ioutil"
	"text/template"
)

const (
	createDesc = `
'inspect' command creates creates a configuration file ready to use for the go sdk`
	createExample = `  kubectl hlf inspect --output hlf-cfg.yaml`
)

type inspectCmd struct {
	fileOutput    string
	organizations []string
}

func (c *inspectCmd) validate() error {
	return nil
}

type OrderingService struct {
	MSPID        string
	OrderingNode string
}
type OrderingNode struct {
	TLSCert string
	URL     string
}
type Peer struct {
	URL     string
	TLSCert string
}

const tmplGoConfig = `
name: hlf-network
version: 1.0.0
client:
  organization: ""
organizations:
  {{ range $mspID, $org := .Organizations }}
  {{$mspID}}:
    mspid: {{$mspID}}
    cryptoPath: /tmp/cryptopath
    users: {}
    peers:
      {{- range $peer := $org.Peers }}
      - {{ $peer.Name }}
 	  {{- end }}
    orderers:
      {{- range $ordService := $org.OrderingServices }}
      {{- range $orderer := $ordService.Orderers }}
      - {{ $orderer.Name }}
 	  {{- end }}
 	  {{- end }}

    {{- end }}

orderers:
{{- range $ordService := .Orderers }}
{{- range $orderer := $ordService.Orderers }}
  "{{$orderer.Name}}":
    url: {{ $orderer.Status.URL }}
    grpcOptions:
      allow-insecure: false
    tlsCACerts:
      pem: |
{{ $orderer.Spec.TLSRootCert | indent 8 }}
{{- end }}
{{- end }}

peers:
  {{- range $peer := .Peers }}
  "{{$peer.Name}}":
    url: {{ $peer.Status.URL }}
    grpcOptions:
      hostnameOverride: ""
      ssl-target-name-override: ""
      allow-insecure: false
    tlsCACerts:
      pem: |
{{ $peer.Status.TlsCert | indent 8 }}
{{- end }}

channels: {}
`

func (c *inspectCmd) run(out io.Writer) error {
	oclient, err := helpers.GetKubeOperatorClient()
	if err != nil {
		return err
	}
	ns := ""
	certAuths, err := helpers.GetClusterCAs(oclient, ns)
	if err != nil {
		return err
	}
	ordOrgs, orderers, err := helpers.GetClusterOrderers(oclient, ns)
	if err != nil {
		return err
	}
	peerOrgs, peers, err := helpers.GetClusterPeers(oclient, ns)
	if err != nil {
		return err
	}
	filterByOrgs := len(c.organizations) > 0
	orgMap := map[string]*helpers.Organization{}
	for _, v := range ordOrgs {
		if filterByOrgs && utils.Contains(c.organizations, v.MspID) {
			orgMap[v.MspID] = v
		}
	}
	for _, v := range peerOrgs {
		if filterByOrgs && utils.Contains(c.organizations, v.MspID) {
			orgMap[v.MspID] = v
		}
	}
	tmpl, err := template.New("test").Funcs(sprig.HermeticTxtFuncMap()).Parse(tmplGoConfig)
	if err != nil {
		return err
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, map[string]interface{}{
		"Peers":         peers,
		"Orderers":      orderers,
		"Organizations": orgMap,
		"CertAuths":     certAuths,
	})
	if err != nil {
		return err
	}
	if c.fileOutput != "" {
		err = ioutil.WriteFile(c.fileOutput, buf.Bytes(), 0644)
		if err != nil {
			return err
		}
	} else {
		_, err = fmt.Fprint(out, buf.String())
		if err != nil {
			return err
		}
	}

	return nil
}

func NewInspectHLFConfig(out io.Writer) *cobra.Command {
	c := &inspectCmd{}
	cmd := &cobra.Command{
		Use:     "inspect",
		Short:   "Inspect the components deployed",
		Long:    createDesc,
		Example: createExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := c.validate(); err != nil {
				return err
			}
			return c.run(out)
		},
	}

	f := cmd.Flags()
	f.StringVar(&c.fileOutput, "output", "", "output file")
	f.StringArrayVarP(&c.organizations, "organizations", "o", []string{}, "organizations to export")

	return cmd
}
