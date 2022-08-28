package nc

import (
	"bytes"
	"context"
	"fmt"
	"github.com/Masterminds/sprig/v3"
	hlfv1alpha1 "github.com/kfsoftware/hlf-operator/api/hlf.kungfusoftware.es/v1alpha1"
	"github.com/kfsoftware/hlf-operator/kubectl-hlf/cmd/helpers"
	operatorv1 "github.com/kfsoftware/hlf-operator/pkg/client/clientset/versioned"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"text/template"
)

type CA struct {
	Name         string
	URL          string
	EnrollID     string
	EnrollSecret string
	CAName       string
	TLSCert      string
}
type Org struct {
	MSPID     string
	CertAuths []string
	Peers     []string
	Orderers  []string
}

type Peer struct {
	Name      string
	URL       string
	TLSCACert string
}

type Orderer struct {
	URL       string
	Name      string
	TLSCACert string
}

const tmplGoConfig = `
name: hlf-network
version: 1.0.0
client:
  organization: "{{ .Organization }}"
{{- if not .Organizations }}
organizations: {}
{{- else }}
organizations:
  {{ range $org := .Organizations }}
  {{ $org.MSPID }}:
    mspid: {{ $org.MSPID }}
    cryptoPath: /tmp/cryptopath
    users: {}
{{- if not $org.CertAuths }}
    certificateAuthorities: []
{{- else }}
    certificateAuthorities: 
      {{- range $ca := $org.CertAuths }}
      - {{ $ca.Name }}
 	  {{- end }}
{{- end }}
{{- if not $org.Peers }}
    peers: []
{{- else }}
    peers:
      {{- range $peer := $org.Peers }}
      - {{ $peer }}
 	  {{- end }}
{{- end }}
{{- if not $org.Orderers }}
    orderers: []
{{- else }}
    orderers:
      {{- range $orderer := $org.Orderers }}
      - {{ $orderer }}
 	  {{- end }}

    {{- end }}
{{- end }}
{{- end }}

{{- if not .Orderers }}
{{- else }}
orderers:
{{- range $orderer := .Orderers }}
  {{$orderer.Name}}:
    url: {{ $orderer.URL }}
    grpcOptions:
      allow-insecure: false
    tlsCACerts:
      pem: |
{{ $orderer.TLSCACert | indent 8 }}
{{- end }}
{{- end }}

{{- if not .Peers }}
{{- else }}
peers:
  {{- range $peer := .Peers }}
  {{$peer.Name}}:
    url: {{ $peer.URL }}
    tlsCACerts:
      pem: |
{{ $peer.TLSCACert | indent 8 }}
{{- end }}
{{- end }}

{{- if not .CertAuths }}
{{- else }}
certificateAuthorities:
{{- range $ca := .CertAuths }}
  {{ $ca.Name }}:
    url: https://{{ $ca.URL }}
{{if $ca.EnrollID }}
    registrar:
        enrollId: {{ $ca.EnrollID }}
        enrollSecret: {{ $ca.EnrollSecret }}
{{ end }}
    caName: {{ $ca.CAName }}
    tlsCACerts:
      pem: 
       - |
{{ $ca.TLSCert | indent 12 }}

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

func GenerateNetworkConfig(channel *hlfv1alpha1.FabricMainChannel, kubeClientset *kubernetes.Clientset, hlfClientSet *operatorv1.Clientset, mspID string) (*NetworkConfigResponse, error) {
	tmpl, err := template.New("networkConfig").Funcs(sprig.HermeticTxtFuncMap()).Parse(tmplGoConfig)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	orgs := []*Org{}
	var peers []*Peer
	var certAuths []*CA
	var ordererNodes []*Orderer

	ctx := context.Background()
	for _, peerOrg := range channel.Spec.PeerOrganizations {
		orgs = append(orgs, &Org{
			MSPID:     peerOrg.MSPID,
			CertAuths: []string{},
			Peers:     []string{},
			Orderers:  []string{},
		})
	}
	for _, ordOrg := range channel.Spec.OrdererOrganizations {
		fabricCA, err := hlfClientSet.HlfV1alpha1().FabricCAs(ordOrg.CANamespace).Get(ctx, ordOrg.CAName, v1.GetOptions{})
		if err != nil {
			return nil, err
		}
		org := &Org{
			MSPID:     ordOrg.MSPID,
			CertAuths: []string{},
			Peers:     []string{},
			Orderers:  []string{},
		}
		for _, ordererEndpoint := range ordOrg.OrdererEndpoints {
			ordererName := ordererEndpoint
			org.Orderers = append(org.Orderers, ordererName)
			ordererNodes = append(ordererNodes, &Orderer{
				URL:       fmt.Sprintf("grpcs://%s", ordererEndpoint),
				Name:      ordererName,
				TLSCACert: fabricCA.Status.TLSCACert,
			})
		}
		orgs = append(orgs, org)
	}
	//for _, externalOrdOrg := range channel.Spec.ExternalOrdererOrganizations {
	//
	//}
	//for _, certAuth := range certAuths {
	//	tlsCACertPem := certAuth.Status.TLSCACert
	//	roots := x509.NewCertPool()
	//	ok := roots.AppendCertsFromPEM([]byte(tlsCACertPem))
	//	if !ok {
	//		panic("failed to parse root certificate")
	//	}
	//	for mspID, org := range orgMap {
	//		for _, peer := range org.Peers {
	//			block, _ := pem.Decode([]byte(peer.Status.TlsCert))
	//			if block == nil {
	//				continue
	//			}
	//			cert, err := x509.ParseCertificate(block.Bytes)
	//			if err != nil {
	//				continue
	//			}
	//			opts := x509.VerifyOptions{
	//				Roots:         roots,
	//				Intermediates: x509.NewCertPool(),
	//			}
	//
	//			if _, err := cert.Verify(opts); err == nil {
	//				orgMap[mspID].CertAuths = append(orgMap[mspID].CertAuths, certAuth)
	//			}
	//		}
	//	}
	//	for _, ord := range ordererNodes {
	//		block, _ := pem.Decode([]byte(ord.Status.TlsCert))
	//		if block == nil {
	//			continue
	//		}
	//		cert, err := x509.ParseCertificate(block.Bytes)
	//		if err != nil {
	//			continue
	//		}
	//		opts := x509.VerifyOptions{
	//			Roots:         roots,
	//			Intermediates: x509.NewCertPool(),
	//		}
	//		if _, err := cert.Verify(opts); err == nil {
	//			_, ok = orgMap[ord.Spec.MspID]
	//			if !ok {
	//				orgMap[ord.Spec.MspID] = &Organization{
	//					Type:         helpers.OrdererType,
	//					MspID:        "",
	//					OrdererNodes: []*helpers.ClusterOrdererNode{},
	//					Peers:        []*helpers.ClusterPeer{},
	//					CertAuths:    []*helpers.ClusterCA{certAuth},
	//				}
	//			} else {
	//				orgMap[ord.Spec.MspID].CertAuths = append(orgMap[ord.Spec.MspID].CertAuths, certAuth)
	//			}
	//		}
	//	}
	//
	//}
	//for _, ord := range ordererNodes {
	//	orgMap[ord.Spec.MspID].OrdererNodes = append(orgMap[ord.Spec.MspID].OrdererNodes, ord)
	//}
	//for _, peer := range clusterPeers {
	//	peers = append(peers, peer)
	//}
	err = tmpl.Execute(&buf, map[string]interface{}{
		"Peers":         peers,
		"Orderers":      ordererNodes,
		"Organizations": orgs,
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

func GenerateNetworkConfigForFollower(channel *hlfv1alpha1.FabricFollowerChannel, kubeClientset *kubernetes.Clientset, hlfClientSet *operatorv1.Clientset, mspID string) (*NetworkConfigResponse, error) {
	tmpl, err := template.New("networkConfig").Funcs(sprig.HermeticTxtFuncMap()).Parse(tmplGoConfig)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	orgs := []*Org{}
	var peers []*Peer
	var certAuths []*CA
	var ordererNodes []*Orderer

	ctx := context.Background()
	org := &Org{
		MSPID:     channel.Spec.MSPID,
		CertAuths: []string{},
		Peers:     []string{},
		Orderers:  []string{},
	}
	for _, peer := range channel.Spec.PeersToJoin {
		fabricPeer, err := hlfClientSet.HlfV1alpha1().FabricPeers(peer.Namespace).Get(ctx, peer.Name, v1.GetOptions{})
		if err != nil {
			return nil, err
		}
		peerName := fmt.Sprintf("%s.%s", fabricPeer.Name, fabricPeer.Namespace)
		org.Peers = append(org.Peers, peerName)
		peers = append(peers, &Peer{
			Name:      peerName,
			URL:       fmt.Sprintf("grpcs://%s:%d", fabricPeer.Spec.Istio.Hosts[0], fabricPeer.Spec.Istio.Port),
			TLSCACert: fabricPeer.Status.TlsCACert,
		})
	}
	orgs = append(orgs, org)
	for _, orderer := range channel.Spec.Orderers {
		ordererNodes = append(ordererNodes, &Orderer{
			URL:       orderer.URL,
			Name:      orderer.URL,
			TLSCACert: orderer.Certificate,
		})
	}
	err = tmpl.Execute(&buf, map[string]interface{}{
		"Peers":         peers,
		"Orderers":      ordererNodes,
		"Organizations": orgs,
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
