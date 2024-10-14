package nc

import (
	"bytes"
	"context"
	"fmt"
	"text/template"

	log "github.com/sirupsen/logrus"

	"github.com/Masterminds/sprig/v3"
	"github.com/kfsoftware/hlf-operator/kubectl-hlf/cmd/helpers"
	hlfv1alpha1 "github.com/kfsoftware/hlf-operator/pkg/apis/hlf.kungfusoftware.es/v1alpha1"
	operatorv1 "github.com/kfsoftware/hlf-operator/pkg/client/clientset/versioned"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
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
        enrollSecret: "{{ $ca.EnrollSecret }}"
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
	for _, peerOrg := range channel.Spec.ExternalPeerOrganizations {
		orgs = append(orgs, &Org{
			MSPID:     peerOrg.MSPID,
			CertAuths: []string{},
			Peers:     []string{},
			Orderers:  []string{},
		})
	}
	for _, ordOrg := range channel.Spec.OrdererOrganizations {
		var tlsCACert string
		if ordOrg.TLSCACert != "" {
			tlsCACert = ordOrg.TLSCACert
		} else {
			fabricCA, err := hlfClientSet.HlfV1alpha1().FabricCAs(ordOrg.CANamespace).Get(ctx, ordOrg.CAName, v1.GetOptions{})
			if err != nil {
				return nil, err
			}
			tlsCACert = fabricCA.Status.TLSCACert
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
				TLSCACert: tlsCACert,
			})
		}
		orgs = append(orgs, org)
	}
	for _, externalOrdOrg := range channel.Spec.ExternalOrdererOrganizations {
		org := &Org{
			MSPID:     externalOrdOrg.MSPID,
			CertAuths: []string{},
			Peers:     []string{},
			Orderers:  []string{},
		}
		for _, ordererEndpoint := range externalOrdOrg.OrdererEndpoints {
			ordererName := ordererEndpoint
			org.Orderers = append(org.Orderers, ordererName)
			ordererNodes = append(ordererNodes, &Orderer{
				URL:       fmt.Sprintf("grpcs://%s", ordererEndpoint),
				Name:      ordererName,
				TLSCACert: externalOrdOrg.TLSRootCert,
			})
		}
		orgs = append(orgs, org)
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

func GenerateNetworkConfigForChaincodeCommit(chCommit *hlfv1alpha1.FabricChaincodeCommit, kubeClientset *kubernetes.Clientset, hlfClientSet *operatorv1.Clientset, mspID string) (*NetworkConfigResponse, error) {
	tmpl, err := template.New("networkConfig").Funcs(sprig.HermeticTxtFuncMap()).Parse(tmplGoConfig)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	orgs := []*Org{}
	var peers []*Peer
	var ordererNodes []*Orderer
	var certAuths []*CA

	ctx := context.Background()

	org := &Org{
		MSPID:     chCommit.Spec.MSPID,
		CertAuths: []string{},
		Peers:     []string{},
		Orderers:  []string{},
	}

	for _, peer := range chCommit.Spec.Peers {
		fabricPeer, err := hlfClientSet.HlfV1alpha1().FabricPeers(peer.Namespace).Get(ctx, peer.Name, v1.GetOptions{})
		if err != nil {
			return nil, err
		}
		peerName := fmt.Sprintf("%s.%s", fabricPeer.Name, fabricPeer.Namespace)
		org.Peers = append(org.Peers, peerName)
		peerHost, err := helpers.GetPeerPublicURL(kubeClientset, *fabricPeer)
		if err != nil {
			return nil, err
		}
		peers = append(peers, &Peer{
			Name:      peerName,
			URL:       fmt.Sprintf("grpcs://%s", peerHost),
			TLSCACert: fabricPeer.Status.TlsCACert,
		})
	}

	for _, peer := range chCommit.Spec.ExternalPeers {
		peerName := peer.URL
		org.Peers = append(org.Peers, peerName)
		peers = append(peers, &Peer{
			Name:      peerName,
			URL:       peer.URL,
			TLSCACert: peer.TLSCACert,
		})
	}

	for _, orderer := range chCommit.Spec.Orderers {
		fabricOrderer, err := hlfClientSet.HlfV1alpha1().FabricOrdererNodes(orderer.Namespace).Get(ctx, orderer.Name, v1.GetOptions{})
		if err != nil {
			return nil, err
		}
		ordererName := fmt.Sprintf("%s.%s", fabricOrderer.Name, fabricOrderer.Namespace)
		org.Orderers = append(org.Orderers, ordererName)
		ordererHost, err := helpers.GetOrdererPublicURL(kubeClientset, *fabricOrderer)
		if err != nil {
			return nil, err
		}
		ordererNodes = append(ordererNodes, &Orderer{
			URL:       fmt.Sprintf("grpcs://%s", ordererHost),
			Name:      ordererName,
			TLSCACert: fabricOrderer.Status.TlsCert,
		})
	}

	for _, orderer := range chCommit.Spec.ExternalOrderers {
		ordererName := orderer.URL
		org.Orderers = append(org.Orderers, ordererName)
		ordererNodes = append(ordererNodes, &Orderer{
			URL:       orderer.URL,
			Name:      ordererName,
			TLSCACert: orderer.TLSCACert,
		})
	}

	orgs = append(orgs, org)
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

func GenerateNetworkConfigForChaincodeInstall(chInstall *hlfv1alpha1.FabricChaincodeInstall, kubeClientset *kubernetes.Clientset, hlfClientSet *operatorv1.Clientset, mspID string) (*NetworkConfigResponse, error) {
	tmpl, err := template.New("networkConfig").Funcs(sprig.HermeticTxtFuncMap()).Parse(tmplGoConfig)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	orgs := []*Org{}
	var peers []*Peer
	var certAuths []*CA

	ctx := context.Background()

	org := &Org{
		MSPID:     chInstall.Spec.MSPID,
		CertAuths: []string{},
		Peers:     []string{},
		Orderers:  []string{},
	}

	for _, peer := range chInstall.Spec.Peers {
		fabricPeer, err := hlfClientSet.HlfV1alpha1().FabricPeers(peer.Namespace).Get(ctx, peer.Name, v1.GetOptions{})
		if err != nil {
			return nil, err
		}
		peerName := fmt.Sprintf("%s.%s", fabricPeer.Name, fabricPeer.Namespace)
		org.Peers = append(org.Peers, peerName)
		peerHost, err := helpers.GetPeerPublicURL(kubeClientset, *fabricPeer)
		if err != nil {
			return nil, err
		}
		peers = append(peers, &Peer{
			Name:      peerName,
			URL:       fmt.Sprintf("grpcs://%s", peerHost),
			TLSCACert: fabricPeer.Status.TlsCACert,
		})
	}

	for _, peer := range chInstall.Spec.ExternalPeers {
		peerName := peer.URL
		org.Peers = append(org.Peers, peerName)
		peers = append(peers, &Peer{
			Name:      peerName,
			URL:       peer.URL,
			TLSCACert: peer.TLSCACert,
		})
	}
	orgs = append(orgs, org)
	err = tmpl.Execute(&buf, map[string]interface{}{
		"Peers":         peers,
		"Orderers":      []string{},
		"Organizations": orgs,
		"CertAuths":     certAuths,
		"Organization":  mspID,
		"Internal":      false,
	})
	if err != nil {
		return nil, err
	}
	log.Infof("Generated network config %s", buf.String())
	return &NetworkConfigResponse{
		NetworkConfig: buf.String(),
	}, nil
}

func GenerateNetworkConfigForChaincodeApprove(chInstall *hlfv1alpha1.FabricChaincodeApprove, kubeClientset *kubernetes.Clientset, hlfClientSet *operatorv1.Clientset, mspID string) (*NetworkConfigResponse, error) {
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
		MSPID:     mspID,
		CertAuths: []string{},
		Peers:     []string{},
		Orderers:  []string{},
	}

	for _, peer := range chInstall.Spec.Peers {
		fabricPeer, err := hlfClientSet.HlfV1alpha1().FabricPeers(peer.Namespace).Get(ctx, peer.Name, v1.GetOptions{})
		if err != nil {
			return nil, err
		}
		peerName := fmt.Sprintf("%s.%s", fabricPeer.Name, fabricPeer.Namespace)
		org.Peers = append(org.Peers, peerName)
		peerHost, err := helpers.GetPeerPublicURL(kubeClientset, *fabricPeer)
		if err != nil {
			return nil, err
		}
		peers = append(peers, &Peer{
			Name:      peerName,
			URL:       fmt.Sprintf("grpcs://%s", peerHost),
			TLSCACert: fabricPeer.Status.TlsCACert,
		})
	}

	for _, peer := range chInstall.Spec.ExternalPeers {
		peerName := peer.URL
		org.Peers = append(org.Peers, peerName)
		peers = append(peers, &Peer{
			Name:      peerName,
			URL:       peer.URL,
			TLSCACert: peer.TLSCACert,
		})
	}

	for _, orderer := range chInstall.Spec.Orderers {
		fabricOrderer, err := hlfClientSet.HlfV1alpha1().FabricOrdererNodes(orderer.Namespace).Get(ctx, orderer.Name, v1.GetOptions{})
		if err != nil {
			return nil, err
		}
		ordererName := fmt.Sprintf("%s.%s", fabricOrderer.Name, fabricOrderer.Namespace)
		ordererHost, err := helpers.GetOrdererPublicURL(kubeClientset, *fabricOrderer)
		if err != nil {
			return nil, err
		}
		ordererNodes = append(ordererNodes, &Orderer{
			Name:      ordererName,
			URL:       fmt.Sprintf("grpcs://%s", ordererHost),
			TLSCACert: fabricOrderer.Status.TlsCert,
		})
	}

	for _, orderer := range chInstall.Spec.ExternalOrderers {
		ordererName := orderer.URL
		ordererNodes = append(ordererNodes, &Orderer{
			Name:      ordererName,
			URL:       orderer.URL,
			TLSCACert: orderer.TLSCACert,
		})
	}

	orgs = append(orgs, org)

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
	log.Infof("Generated network config %s", buf.String())
	return &NetworkConfigResponse{
		NetworkConfig: buf.String(),
	}, nil
}

func GenerateNetworkConfigForFollower(chInstall *hlfv1alpha1.FabricFollowerChannel, kubeClientset *kubernetes.Clientset, hlfClientSet *operatorv1.Clientset, mspID string) (*NetworkConfigResponse, error) {
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
		MSPID:     chInstall.Spec.MSPID,
		CertAuths: []string{},
		Peers:     []string{},
		Orderers:  []string{},
	}
	for _, peer := range chInstall.Spec.PeersToJoin {
		fabricPeer, err := hlfClientSet.HlfV1alpha1().FabricPeers(peer.Namespace).Get(ctx, peer.Name, v1.GetOptions{})
		if err != nil {
			return nil, err
		}
		peerName := fmt.Sprintf("%s.%s", fabricPeer.Name, fabricPeer.Namespace)
		org.Peers = append(org.Peers, peerName)
		peerHost, err := helpers.GetPeerPublicURL(kubeClientset, *fabricPeer)
		if err != nil {
			return nil, err
		}
		peers = append(peers, &Peer{
			Name:      peerName,
			URL:       fmt.Sprintf("grpcs://%s", peerHost),
			TLSCACert: fabricPeer.Status.TlsCACert,
		})
	}
	for _, peer := range chInstall.Spec.ExternalPeersToJoin {
		peerName := peer.URL
		org.Peers = append(org.Peers, peerName)
		peers = append(peers, &Peer{
			Name:      peerName,
			URL:       peer.URL,
			TLSCACert: peer.TLSCACert,
		})
	}
	orgs = append(orgs, org)
	for _, orderer := range chInstall.Spec.Orderers {
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
