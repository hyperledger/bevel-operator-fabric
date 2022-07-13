package export

import (
	"archive/zip"
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/kfsoftware/hlf-operator/kubectl-hlf/cmd/helpers"
	"io"
)

type FabricOperationsCA struct {
	DisplayName   string `json:"display_name"`
	ApiUrl        string `json:"api_url"`
	OperationsUrl string `json:"operations_url"`
	CaUrl         string `json:"ca_url"`
	Type          string `json:"type"`
	CaName        string `json:"ca_name"`
	TlscaName     string `json:"tlsca_name"`
	TlsCert       string `json:"tls_cert"`
	Name          string `json:"name"`
}
type FabricUnused struct {
	Result FabricOrgResult `json:"result"`
}
type FabricOrgResult struct {
	CAName                    string `json:"CAName"`
	CAChain                   string `json:"CAChain"`
	IssuerPublicKey           string `json:"IssuerPublicKey"`
	IssuerRevocationPublicKey string `json:"IssuerRevocationPublicKey"`
	Version                   string `json:"Version"`
}
type FabricOperationsOrg struct {
	DisplayName   string        `json:"display_name"`
	MspId         string        `json:"msp_id"`
	Type          string        `json:"type"`
	Admins        []string      `json:"admins"`
	RootCerts     []string      `json:"root_certs"`
	TlsRootCerts  []string      `json:"tls_root_certs"`
	FabricNodeOus FabricNodeOus `json:"fabric_node_ous"`
	HostUrl       string        `json:"host_url,omitempty"`
	Name          string        `json:"name"`
}
type FabricNodeOus struct {
	AdminOuIdentifier   FabricOrgOUIdentifier `json:"admin_ou_identifier"`
	ClientOuIdentifier  FabricOrgOUIdentifier `json:"client_ou_identifier"`
	Enable              bool                  `json:"enable"`
	OrdererOuIdentifier FabricOrgOUIdentifier `json:"orderer_ou_identifier"`
	PeerOuIdentifier    FabricOrgOUIdentifier `json:"peer_ou_identifier"`
}
type FabricOrgOUIdentifier struct {
	Certificate                  string `json:"certificate"`
	OrganizationalUnitIdentifier string `json:"organizational_unit_identifier"`
}
type FabricOperationsPeer struct {
	DisplayName   string              `json:"display_name"`
	GrpcwpUrl     string              `json:"grpcwp_url"`
	ApiUrl        string              `json:"api_url"`
	OperationsUrl string              `json:"operations_url"`
	MspId         string              `json:"msp_id"`
	Name          string              `json:"name"`
	Type          string              `json:"type"`
	Msp           FabricOperationsMSP `json:"msp"`
	Pem           string              `json:"pem"`
	TlsCert       string              `json:"tls_cert"`
	TlsCaRootCert string              `json:"tls_ca_root_cert"`
}
type FabricOperationsMSP struct {
	Component FabricPeerComponentMSP `json:"component"`
	CA        FabricPeerMSPCA        `json:"ca"`
	TLSCA     FabricPeerMSPCA        `json:"tlsca"`
}
type FabricPeerComponentMSP struct {
	AdminCerts []interface{} `json:"admin_certs"`
	TlsCert    string        `json:"tls_cert"`
}
type FabricPeerMSPCA struct {
	RootCerts []string `json:"root_certs"`
}
type FeatureFlagsOrderer struct {
	OSNAdminFeatsEnabled bool `json:"osnadmin_feats_enabled"`
}
type FabricOperationsOrderer struct {
	DisplayName   string              `json:"display_name"`
	GrpcwpUrl     string              `json:"grpcwp_url"`
	ApiUrl        string              `json:"api_url"`
	OperationsUrl string              `json:"operations_url"`
	Type          string              `json:"type"`
	MspId         string              `json:"msp_id"`
	ClusterId     string              `json:"cluster_id"`
	ClusterName   string              `json:"cluster_name"`
	Name          string              `json:"name"`
	Msp           FabricOperationsMSP `json:"msp"`
	Pem           string              `json:"pem"`
	OSNAdminURL   string              `json:"osnadmin_url"`
	Systemless    bool                `json:"systemless"`
	TlsCert       string              `json:"tls_cert"`
	TlsCaRootCert string              `json:"tls_ca_root_cert"`
	FeatureFlags  FeatureFlagsOrderer `json:"feature_flags"`
}

func appendFile(filename string, data []byte, zipw *zip.Writer) error {
	wr, err := zipw.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create entry for %s in zip file: %s", filename, err)
	}
	dataReader := bytes.NewReader(data)
	if _, err := io.Copy(wr, dataReader); err != nil {
		return fmt.Errorf("failed to write %s to zip: %s", filename, err)
	}
	return nil
}

func mapFabricOperationsCA(clusterCA *helpers.ClusterCA) *FabricOperationsCA {
	caURL := fmt.Sprintf("https://%s", clusterCA.PublicURL)
	displayName := fmt.Sprintf("%s_%s", clusterCA.Object.Name, clusterCA.Object.Namespace)
	if len(displayName) >= 30 {
		displayName = displayName[0:29]
	}
	internalOperationsURL := fmt.Sprintf("http://%s.%s:%d", clusterCA.Object.Name, clusterCA.Object.Namespace, 9443)
	ca := &FabricOperationsCA{
		DisplayName:   displayName,
		ApiUrl:        caURL,
		OperationsUrl: internalOperationsURL,
		CaUrl:         caURL,
		Type:          "fabric-ca",
		CaName:        "ca",
		TlscaName:     "tlsca",
		TlsCert:       clusterCA.Object.Status.TlsCert,
		Name:          displayName,
	}
	return ca
}

func mapFabricOperationsPeer(clusterPeer *helpers.ClusterPeer) (*FabricOperationsPeer, error) {
	apiURL := fmt.Sprintf("grpcs://%s", clusterPeer.PublicURL)
	displayName := fmt.Sprintf("%s_%s", clusterPeer.Object.Name, clusterPeer.Object.Namespace)
	if len(displayName) >= 30 {
		displayName = displayName[0:29]
	}
	if clusterPeer.Spec.GRPCProxy == nil {
		return nil, fmt.Errorf("grpc proxy not configured for peer %s", clusterPeer.Object.Name)
	}
	grpcwpURL := fmt.Sprintf("https://%s:%d", clusterPeer.Spec.GRPCProxy.Istio.Hosts[0], clusterPeer.Spec.GRPCProxy.Istio.Port)
	internalOperationsURL := fmt.Sprintf("http://%s.%s:%d", clusterPeer.ObjectMeta.Name, clusterPeer.ObjectMeta.Namespace, 9443)
	fabricOperationsPeer := &FabricOperationsPeer{
		DisplayName:   displayName,
		GrpcwpUrl:     grpcwpURL,
		ApiUrl:        apiURL,
		OperationsUrl: internalOperationsURL,
		MspId:         clusterPeer.Spec.MspID,
		Name:          displayName,
		Type:          "fabric-peer",
		Msp: FabricOperationsMSP{
			Component: FabricPeerComponentMSP{
				AdminCerts: []interface{}{},
				TlsCert:    base64.StdEncoding.EncodeToString([]byte(clusterPeer.Object.Status.TlsCert)),
			},
			CA: FabricPeerMSPCA{
				RootCerts: []string{
					base64.StdEncoding.EncodeToString([]byte(clusterPeer.Object.Status.SignCACert)),
				},
			},
			TLSCA: FabricPeerMSPCA{
				RootCerts: []string{
					base64.StdEncoding.EncodeToString([]byte(clusterPeer.Object.Status.TlsCACert)),
				},
			},
		},
		Pem:           base64.StdEncoding.EncodeToString([]byte(clusterPeer.Status.SignCert)),
		TlsCert:       base64.StdEncoding.EncodeToString([]byte(clusterPeer.Status.TlsCert)),
		TlsCaRootCert: base64.StdEncoding.EncodeToString([]byte(clusterPeer.Status.TlsCACert)),
	}
	return fabricOperationsPeer, nil
}

type MapFabricOperationsOrderer struct {
	ClusterID   string
	ClusterName string
	OSNAdminURL string
}

func mapFabricOperationsOrderer(clusterOrdererNode *helpers.ClusterOrdererNode, opts MapFabricOperationsOrderer) (*FabricOperationsOrderer, error) {
	apiURL := fmt.Sprintf("grpcs://%s", clusterOrdererNode.PublicURL)
	displayName := fmt.Sprintf("%s_%s", clusterOrdererNode.ObjectMeta.Name, clusterOrdererNode.ObjectMeta.Namespace)
	if len(displayName) >= 30 {
		displayName = displayName[0:29]
	}
	if clusterOrdererNode.Spec.GRPCProxy == nil {
		return nil, fmt.Errorf("grpc proxy not configured for peer %s", clusterOrdererNode.ObjectMeta.Name)
	}
	grpcwpURL := fmt.Sprintf("https://%s:%d", clusterOrdererNode.Spec.GRPCProxy.Istio.Hosts[0], clusterOrdererNode.Spec.GRPCProxy.Istio.Port)
	internalOperationsURL := fmt.Sprintf("http://%s.%s:%d", clusterOrdererNode.ObjectMeta.Name, clusterOrdererNode.ObjectMeta.Namespace, 9443)

	fabricOperationsPeer := &FabricOperationsOrderer{
		DisplayName:   displayName,
		GrpcwpUrl:     grpcwpURL,
		ApiUrl:        apiURL,
		OperationsUrl: internalOperationsURL,
		MspId:         clusterOrdererNode.Spec.MspID,
		FeatureFlags:  FeatureFlagsOrderer{OSNAdminFeatsEnabled: true},
		Name:          displayName,
		OSNAdminURL:   opts.OSNAdminURL,
		Systemless:    true,
		Type:          "fabric-orderer",
		ClusterId:     opts.ClusterID,
		ClusterName:   opts.ClusterName,
		Msp: FabricOperationsMSP{
			Component: FabricPeerComponentMSP{
				AdminCerts: []interface{}{},
				TlsCert:    base64.StdEncoding.EncodeToString([]byte(clusterOrdererNode.Status.TlsCert)),
			},
			CA: FabricPeerMSPCA{
				RootCerts: []string{
					base64.StdEncoding.EncodeToString([]byte(clusterOrdererNode.Status.SignCACert)),
				},
			},
			TLSCA: FabricPeerMSPCA{
				RootCerts: []string{
					base64.StdEncoding.EncodeToString([]byte(clusterOrdererNode.Status.TlsCACert)),
				},
			},
		},
		Pem:           base64.StdEncoding.EncodeToString([]byte(clusterOrdererNode.Status.SignCert)),
		TlsCert:       base64.StdEncoding.EncodeToString([]byte(clusterOrdererNode.Status.TlsCert)),
		TlsCaRootCert: base64.StdEncoding.EncodeToString([]byte(clusterOrdererNode.Status.TlsCACert)),
	}
	return fabricOperationsPeer, nil
}

type MapFabricOperationsOrg struct {
	MSPID   string
	HostURL string
}

func mapFabricOperationsOrg(clusterCA *helpers.ClusterCA, opts MapFabricOperationsOrg) (*FabricOperationsOrg, error) {
	displayName := fmt.Sprintf("%s_%s", clusterCA.Object.Name, clusterCA.Object.Namespace)
	if len(displayName) >= 30 {
		displayName = displayName[0:29]
	}
	fabricOperationsPeer := &FabricOperationsOrg{
		DisplayName:  displayName,
		MspId:        opts.MSPID,
		Type:         "msp",
		HostUrl:      opts.HostURL,
		Admins:       []string{},
		RootCerts:    []string{base64.StdEncoding.EncodeToString([]byte(clusterCA.Object.Status.CACert))},
		TlsRootCerts: []string{base64.StdEncoding.EncodeToString([]byte(clusterCA.Object.Status.TLSCACert))},
		FabricNodeOus: FabricNodeOus{
			AdminOuIdentifier: FabricOrgOUIdentifier{
				Certificate:                  base64.StdEncoding.EncodeToString([]byte(clusterCA.Object.Status.CACert)),
				OrganizationalUnitIdentifier: "admin",
			},
			ClientOuIdentifier: FabricOrgOUIdentifier{
				Certificate:                  base64.StdEncoding.EncodeToString([]byte(clusterCA.Object.Status.CACert)),
				OrganizationalUnitIdentifier: "client",
			},
			Enable: true,
			OrdererOuIdentifier: FabricOrgOUIdentifier{
				Certificate:                  base64.StdEncoding.EncodeToString([]byte(clusterCA.Object.Status.CACert)),
				OrganizationalUnitIdentifier: "orderer",
			},
			PeerOuIdentifier: FabricOrgOUIdentifier{
				Certificate:                  base64.StdEncoding.EncodeToString([]byte(clusterCA.Object.Status.CACert)),
				OrganizationalUnitIdentifier: "peer",
			},
		},
		Name: displayName,
	}
	return fabricOperationsPeer, nil
}
