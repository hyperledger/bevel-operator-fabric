package nc

import (
	"bytes"
	"crypto/x509"
	"encoding/pem"
	"github.com/Masterminds/sprig/v3"
	"github.com/kfsoftware/hlf-operator/kubectl-hlf/cmd/helpers"
	operatorv1 "github.com/kfsoftware/hlf-operator/pkg/client/clientset/versioned"
	"k8s.io/client-go/kubernetes"
	"text/template"
)

const tmplGoConfig = `
name: hlf-network
version: 1.0.0
client:
  organization: "{{ .Organization }}"
{{- if not .Organizations }}
organizations: {}
{{- else }}
organizations:
  {{ range $mspID, $org := .Organizations }}
  {{$mspID}}:
    mspid: {{$mspID}}
    cryptoPath: /tmp/cryptopath
    users: {}
{{- if not $org.CertAuths }}
    certificateAuthorities: []
{{- else }}
    certificateAuthorities: 
      {{- range $ca := $org.CertAuths }}
      - {{ $ca.Name }}-sign
      - {{ $ca.Name }}-tls
 	  {{- end }}
{{- end }}
{{- if not $org.Peers }}
    peers: []
{{- else }}
    peers:
      {{- range $peer := $org.Peers }}
      - {{ $peer.Name }}
 	  {{- end }}
{{- end }}
{{- if not $org.OrdererNodes }}
    orderers: []
{{- else }}
    orderers:
      {{- range $orderer := $org.OrdererNodes }}
      - {{ $orderer.Name }}
 	  {{- end }}

    {{- end }}
{{- end }}
{{- end }}

{{- if not .Orderers }}
orderers: []
{{- else }}
orderers:
{{- range $orderer := .Orderers }}
  {{$orderer.Name}}:
{{if $.Internal }}
    url: grpcs://{{ $orderer.PrivateURL }}
{{ else }}
    url: grpcs://{{ $orderer.PublicURL }}
{{ end }}
    grpcOptions:
      allow-insecure: false
    tlsCACerts:
      pem: |
{{ or $orderer.Status.TlsCACert $orderer.Status.TlsCert | indent 8 }}
{{- end }}
{{- end }}

{{- if not .Peers }}
peers: []
{{- else }}
peers:
  {{- range $peer := .Peers }}
  {{$peer.Name}}:
{{if $.Internal }}
    url: grpcs://{{ $peer.PrivateURL }}
{{ else }}
    url: grpcs://{{ $peer.PublicURL }}
{{ end }}
    grpcOptions:
      hostnameOverride: ""
      ssl-target-name-override: ""
      allow-insecure: false
    tlsCACerts:
      pem: |
{{ $peer.Status.TlsCACert | indent 8 }}
{{- end }}
{{- end }}

{{- if not .CertAuths }}
certificateAuthorities: []
{{- else }}
certificateAuthorities:
{{- range $ca := .CertAuths }}
  
  {{ $ca.Name }}-tls:
{{if $.Internal }}
    url: https://{{ $ca.PrivateURL }}
{{ else }}
    url: https://{{ $ca.PublicURL }}
{{ end }}
{{if $ca.EnrollID }}
    registrar:
        enrollId: {{ $ca.EnrollID }}
        enrollSecret: {{ $ca.EnrollPWD }}
{{ end }}
    caName: tlsca
    tlsCACerts:
      pem: 
       - |
{{ $ca.Status.TlsCert | indent 12 }}
  
  {{ $ca.Name }}-sign:
{{if $.Internal }}
    url: https://{{ $ca.PrivateURL }}
{{ else }}
    url: https://{{ $ca.PublicURL }}
{{ end }}
{{if $ca.EnrollID }}
    registrar:
        enrollId: {{ $ca.EnrollID }}
        enrollSecret: {{ $ca.EnrollPWD }}
{{ end }}
    caName: ca
    tlsCACerts:
      pem: 
       - |
{{ $ca.Status.TlsCert | indent 12 }}

{{- end }}
{{- end }}

channels:
  _default:
{{- if not .Orderers }}
    orderers: []
{{- else }}
    orderers:
{{- range $orderer := .Orderers }}
      - {{$orderer.Name}}
{{- end }}
{{- end }}
{{- if not .Peers }}
    peers: {}
{{- else }}
    peers:
{{- range $peer := .Peers }}
       {{$peer.Name}}:
        discover: true
        endorsingPeer: true
        chaincodeQuery: true
        ledgerQuery: true
        eventSource: true
{{- end }}
{{- end }}

`

type Organization struct {
	Type         helpers.OrganizationType
	MspID        string
	OrdererNodes []*helpers.ClusterOrdererNode
	Peers        []*helpers.ClusterPeer
	CertAuths    []*helpers.ClusterCA
}
type NetworkConfigResponse struct {
	NetworkConfig string
}

func GenerateNetworkConfig(kubeClientset *kubernetes.Clientset, hlfClientSet *operatorv1.Clientset, mspID string) (*NetworkConfigResponse, error) {
	tmpl, err := template.New("networkConfig").Funcs(sprig.HermeticTxtFuncMap()).Parse(tmplGoConfig)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	certAuths, err := helpers.GetClusterCAs(kubeClientset, hlfClientSet, "")
	if err != nil {
		return nil, err
	}
	ordOrgs, _, err := helpers.GetClusterOrderers(kubeClientset, hlfClientSet, "")
	if err != nil {
		return nil, err
	}
	ordererNodes, err := helpers.GetClusterOrdererNodes(kubeClientset, hlfClientSet, "")
	if err != nil {
		return nil, err
	}
	peerOrgs, clusterPeers, err := helpers.GetClusterPeers(kubeClientset, hlfClientSet, "")
	if err != nil {
		return nil, err
	}
	orgMap := map[string]*Organization{}
	for _, v := range ordOrgs {
		orgMap[v.MspID] = &Organization{
			Type:         v.Type,
			MspID:        v.MspID,
			OrdererNodes: v.OrdererNodes,
			Peers:        v.Peers,
			CertAuths:    []*helpers.ClusterCA{},
		}
	}
	for _, v := range peerOrgs {
		orgMap[v.MspID] = &Organization{
			Type:         v.Type,
			MspID:        v.MspID,
			OrdererNodes: v.OrdererNodes,
			Peers:        v.Peers,
			CertAuths:    []*helpers.ClusterCA{},
		}
	}
	for _, certAuth := range certAuths {
		tlsCACertPem := certAuth.Status.TLSCACert
		roots := x509.NewCertPool()
		ok := roots.AppendCertsFromPEM([]byte(tlsCACertPem))
		if !ok {
			panic("failed to parse root certificate")
		}
		for mspID, org := range orgMap {
			for _, peer := range org.Peers {
				block, _ := pem.Decode([]byte(peer.Status.TlsCert))
				if block == nil {
					continue
				}
				cert, err := x509.ParseCertificate(block.Bytes)
				if err != nil {
					continue
				}
				opts := x509.VerifyOptions{
					Roots:         roots,
					Intermediates: x509.NewCertPool(),
				}

				if _, err := cert.Verify(opts); err == nil {
					orgMap[mspID].CertAuths = append(orgMap[mspID].CertAuths, certAuth)
				}
			}
		}
		for _, ord := range ordererNodes {
			block, _ := pem.Decode([]byte(ord.Status.TlsCert))
			if block == nil {
				continue
			}
			cert, err := x509.ParseCertificate(block.Bytes)
			if err != nil {
				continue
			}
			opts := x509.VerifyOptions{
				Roots:         roots,
				Intermediates: x509.NewCertPool(),
			}
			if _, err := cert.Verify(opts); err == nil {
				orgMap[ord.Spec.MspID].CertAuths = append(orgMap[ord.Spec.MspID].CertAuths, certAuth)
			}
		}

	}
	for _, ord := range ordererNodes {
		orgMap[ord.Spec.MspID].OrdererNodes = append(orgMap[ord.Spec.MspID].OrdererNodes, ord)
	}
	var peers []*helpers.ClusterPeer
	for _, peer := range clusterPeers {
		peers = append(peers, peer)
	}
	err = tmpl.Execute(&buf, map[string]interface{}{
		"Peers":         peers,
		"Orderers":      ordererNodes,
		"Organizations": orgMap,
		"CertAuths":     certAuths,
		"Organization":  mspID,
		"Internal":      false,
	})
	if err != nil {
		return nil, err
	}
	return &NetworkConfigResponse{
		NetworkConfig: buf.String(),
	}, nil
}
