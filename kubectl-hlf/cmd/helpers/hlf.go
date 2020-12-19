package helpers

import (
	"context"
	"fmt"
	hlfv1alpha1 "github.com/kfsoftware/hlf-operator/api/hlf.kungfusoftware.es/v1alpha1"
	operatorv1 "github.com/kfsoftware/hlf-operator/pkg/client/clientset/versioned"
	"github.com/pkg/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type OrganizationType = string

const (
	PeerType    = "PEER"
	OrdererType = "ORDERER"
)

type Organization struct {
	Type             OrganizationType
	MspID            string
	OrderingServices []*ClusterOrderingService
	Peers            []*ClusterPeer
}

type ClusterCA struct {
	Object hlfv1alpha1.FabricCA
	Spec   hlfv1alpha1.FabricCASpec
	Status hlfv1alpha1.FabricCAStatus
	Name   string
}
type ClusterOrderingService struct {
	Name     string
	Object   hlfv1alpha1.FabricOrderingService
	Spec     hlfv1alpha1.FabricOrderingServiceSpec
	Status   hlfv1alpha1.FabricOrderingServiceStatus
	Orderers []*ClusterOrdererNode
}
type ClusterOrdererNode struct {
	Name   string
	Spec   hlfv1alpha1.FabricOrdererNodeSpec
	Status hlfv1alpha1.FabricOrdererNodeStatus
}

type ClusterPeer struct {
	Name      string
	Spec      hlfv1alpha1.FabricPeerSpec
	Status    hlfv1alpha1.FabricPeerStatus
	TLSCACert string
	RootCert  string
	Identity  Identity
}
type Identity struct {
	Key  string
	Cert string
}

func GetClusterCAs(oclient *operatorv1.Clientset, ns string) ([]*ClusterCA, error) {
	ctx := context.Background()
	certAuthsRes, err := oclient.HlfV1alpha1().FabricCAs(ns).List(ctx, v1.ListOptions{})
	if err != nil {
		return nil, err
	}
	var certAuths []*ClusterCA
	for _, certAuth := range certAuthsRes.Items {
		certauthName := fmt.Sprintf("%s.%s", certAuth.Name, certAuth.Namespace)
		certAuths = append(certAuths, &ClusterCA{
			Object: certAuth,
			Spec:   certAuth.Spec,
			Status: certAuth.Status,
			Name:   certauthName,
		})
	}
	return certAuths, nil
}

func GetClusterOrderers(
	oclient *operatorv1.Clientset,
	ns string,
) ([]*Organization, []*ClusterOrderingService, error) {
	ctx := context.Background()
	orderingServices, err := oclient.HlfV1alpha1().FabricOrderingServices(ns).List(ctx, v1.ListOptions{})
	if err != nil {
		return nil, nil, err
	}
	var orderers []*ClusterOrderingService
	for _, ordService := range orderingServices.Items {
		ordNodesRes, err := oclient.HlfV1alpha1().FabricOrdererNodes(ns).List(
			ctx,
			v1.ListOptions{
				LabelSelector: fmt.Sprintf("release=%s", ordService.Name),
			},
		)
		if err != nil {
			return nil, nil, err
		}
		orderingService := &ClusterOrderingService{
			Name:     ordService.FullName(),
			Object:   ordService,
			Spec:     ordService.Spec,
			Status:   ordService.Status,
			Orderers: []*ClusterOrdererNode{},
		}
		orderers = append(orderers, orderingService)
		for _, ordNode := range ordNodesRes.Items {
			orderingService.Orderers = append(
				orderingService.Orderers,
				&ClusterOrdererNode{
					Name:   ordNode.FullName(),
					Spec:   ordNode.Spec,
					Status: ordNode.Status,
				},
			)
		}
	}
	if len(orderers) == 0 {
		return nil, nil, nil
	}
	var organizations []*Organization
	for _, ord := range orderers {
		org := &Organization{
			Type:             OrdererType,
			MspID:            ord.Spec.MspID,
			OrderingServices: []*ClusterOrderingService{ord},
			Peers:            []*ClusterPeer{},
		}
		organizations = append(organizations, org)
	}
	return organizations, orderers, nil
}
func GetCertAuthByURL(oclient *operatorv1.Clientset, host string, port int) (*ClusterCA, error) {
	certAuths, err := GetClusterCAs(oclient, "")
	if err != nil {
		return nil, err
	}
	for _, certAuth := range certAuths {
		if
		// if host and port is specified by NodePort
		certAuth.Status.Host == host && certAuth.Status.Port == port ||
			// if host and port is specified by kubernetes DNS
			certAuth.Name == host && certAuth.Status.Port == 7054 {
			return certAuth, nil
		}

	}
	return nil, errors.Errorf("CA with host=%s port=%d not found", host, port)
}

func GetCertAuthByName(oclient *operatorv1.Clientset, name string, ns string) (*ClusterCA, error) {
	certAuths, err := GetClusterCAs(oclient, "")
	if err != nil {
		return nil, err
	}
	for _, certAuth := range certAuths {
		if certAuth.Object.Name == name && certAuth.Object.Namespace == ns {
			return certAuth, nil
		}

	}
	return nil, errors.Errorf("CA with name=%s not found", name)
}

func GetOrderingServiceByFullName(oclient *operatorv1.Clientset, name string) (*ClusterOrderingService, error) {
	_, ordServices, err := GetClusterOrderers(oclient, "")
	if err != nil {
		return nil, err
	}
	for _, ordService := range ordServices {
		if ordService.Name == name {
			return ordService, nil
		}

	}
	return nil, errors.Errorf("Ordering Service with name=%s not found", name)
}
func GetPeerByFullName(oclient *operatorv1.Clientset, name string) (*ClusterPeer, error) {
	_, peers, err := GetClusterPeers(oclient, "")
	if err != nil {
		return nil, err
	}
	for _, peer := range peers {
		if peer.Name == name {
			return peer, nil
		}

	}
	return nil, errors.Errorf("Peer with name=%s not found", name)
}
func GetCertAuthByFullName(oclient *operatorv1.Clientset, name string) (*ClusterCA, error) {
	certAuths, err := GetClusterCAs(oclient, "")
	if err != nil {
		return nil, err
	}
	for _, certAuth := range certAuths {
		if certAuth.Name == name {
			return certAuth, nil
		}

	}
	return nil, errors.Errorf("CA with name=%s not found", name)
}
func GetClusterPeers(oclient *operatorv1.Clientset, ns string) ([]*Organization, []*ClusterPeer, error) {
	ctx := context.Background()

	peerResponse, err := oclient.HlfV1alpha1().FabricPeers(ns).List(ctx, v1.ListOptions{})
	if err != nil {
		return nil, nil, err
	}
	var peers []*ClusterPeer
	for _, peer := range peerResponse.Items {

		peers = append(
			peers,
			&ClusterPeer{
				Name:     peer.FullName(),
				Spec:     peer.Spec,
				Status:   peer.Status,
				Identity: Identity{},
			},
		)
	}
	orgMap := map[string]*Organization{}
	for _, peer := range peers {
		mspID := peer.Spec.MspID
		org, ok := orgMap[mspID]
		if !ok {
			orgMap[mspID] = &Organization{
				Type:             PeerType,
				MspID:            mspID,
				OrderingServices: []*ClusterOrderingService{},
				Peers:            []*ClusterPeer{},
			}
			org = orgMap[mspID]
		}
		org.Peers = append(org.Peers, peer)
	}
	var organizations []*Organization
	for _, org := range orgMap {
		organizations = append(organizations, org)
	}
	return organizations, peers, nil
}
