package helpers

import (
	"context"
	"fmt"
	"github.com/kfsoftware/hlf-operator/controllers/utils"
	"k8s.io/client-go/kubernetes"
	"strings"

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
	Type         OrganizationType
	MspID        string
	OrdererNodes []*ClusterOrdererNode
	Peers        []*ClusterPeer
}

type ClusterCA struct {
	Object     hlfv1alpha1.FabricCA
	Spec       hlfv1alpha1.FabricCASpec
	Status     hlfv1alpha1.FabricCAStatus
	Name       string
	PublicURL  string
	PrivateURL string
	EnrollID   string
	EnrollPWD  string
	Item       hlfv1alpha1.FabricCA
}

func (c ClusterCA) GetFullName() string {
	return fmt.Sprintf("%s.%s", c.Object.Name, c.Object.Namespace)
}

type ClusterOrderingService struct {
	MSPID    string
	Name     string
	Object   hlfv1alpha1.FabricOrderingService
	Spec     hlfv1alpha1.FabricOrderingServiceSpec
	Status   hlfv1alpha1.FabricOrderingServiceStatus
	Orderers []*ClusterOrdererNode
}

type ClusterOrdererNode struct {
	ObjectMeta v1.ObjectMeta
	Name       string
	PublicURL  string
	PrivateURL string
	Spec       hlfv1alpha1.FabricOrdererNodeSpec
	Status     hlfv1alpha1.FabricOrdererNodeStatus
	Item       hlfv1alpha1.FabricOrdererNode
}

type ClusterPeer struct {
	Name       string
	Spec       hlfv1alpha1.FabricPeerSpec
	Status     hlfv1alpha1.FabricPeerStatus
	PublicURL  string
	PrivateURL string
	TLSCACert  string
	RootCert   string
	Identity   Identity
	MSPID      string
	ObjectMeta v1.ObjectMeta
}
type Identity struct {
	Key  string
	Cert string
}

func GetClusterCAs(clientSet *kubernetes.Clientset, oclient *operatorv1.Clientset, ns string) ([]*ClusterCA, error) {
	ctx := context.Background()
	certAuthsRes, err := oclient.HlfV1alpha1().FabricCAs(ns).List(ctx, v1.ListOptions{})
	if err != nil {
		return nil, err
	}
	var certAuths []*ClusterCA
	for _, certAuth := range certAuthsRes.Items {
		certauthName := fmt.Sprintf("%s.%s", certAuth.Name, certAuth.Namespace)
		privateURL := GetCAPrivateURL(certAuth)
		publicURL, err := GetCAPublicURL(clientSet, certAuth)
		if err != nil {
			return nil, err
		}
		certAuthIdentities := certAuth.Spec.CA.Registry.Identities
		var enrollId string
		var enrollPwd string
		if len(certAuthIdentities) > 0 {
			enrollId = certAuthIdentities[0].Name
			enrollPwd = certAuthIdentities[0].Pass
		}
		certAuths = append(certAuths, &ClusterCA{
			Object:     certAuth,
			Spec:       certAuth.Spec,
			Status:     certAuth.Status,
			Name:       certauthName,
			PublicURL:  publicURL,
			PrivateURL: privateURL,
			EnrollID:   enrollId,
			EnrollPWD:  enrollPwd,
			Item:       certAuth,
		})
	}
	return certAuths, nil
}

func GetClusterOrderers(
	clientSet *kubernetes.Clientset,
	oclient *operatorv1.Clientset,
	ns string,
) ([]*Organization, []*ClusterOrderingService, error) {
	ctx := context.Background()
	ordererNodes, err := oclient.HlfV1alpha1().FabricOrdererNodes(ns).List(ctx, v1.ListOptions{})
	if err != nil {
		return nil, nil, err
	}
	orderingServices, err := oclient.HlfV1alpha1().FabricOrderingServices(ns).List(ctx, v1.ListOptions{})
	if err != nil {
		return nil, nil, err
	}
	var orderers []*ClusterOrderingService
	if len(ordererNodes.Items) > 0 {
		orderingService := &ClusterOrderingService{
			Name:     ordererNodes.Items[0].FullName(),
			MSPID:    ordererNodes.Items[0].Spec.MspID,
			Orderers: []*ClusterOrdererNode{},
		}
		orderers = append(orderers, orderingService)
		for _, ordNode := range ordererNodes.Items {
			publicURL, err := GetOrdererPublicURL(clientSet, ordNode)
			if err != nil {
				return nil, nil, err
			}
			privateURL := GetOrdererPrivateURL(ordNode)
			orderingService.Orderers = append(
				orderingService.Orderers,
				&ClusterOrdererNode{
					Name:       ordNode.FullName(),
					ObjectMeta: ordNode.ObjectMeta,
					Spec:       ordNode.Spec,
					Status:     ordNode.Status,
					PublicURL:  publicURL,
					PrivateURL: privateURL,
				},
			)
		}
	}
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
			MSPID:    ordService.Spec.MspID,
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
			Type:         OrdererType,
			MspID:        ord.MSPID,
			OrdererNodes: []*ClusterOrdererNode{},
			Peers:        []*ClusterPeer{},
		}
		organizations = append(organizations, org)
	}
	return organizations, orderers, nil
}

func GetClusterOrdererNodes(
	clientSet *kubernetes.Clientset,
	oclient *operatorv1.Clientset,
	ns string,
) ([]*ClusterOrdererNode, error) {
	ctx := context.Background()
	ordererNodeList, err := oclient.HlfV1alpha1().FabricOrdererNodes(ns).List(ctx, v1.ListOptions{})
	if err != nil {
		return nil, err
	}
	var ordererNodes []*ClusterOrdererNode

	for _, ordNode := range ordererNodeList.Items {
		publicURL, err := GetOrdererPublicURL(clientSet, ordNode)
		if err != nil {
			return nil, err
		}
		privateURL := GetOrdererPrivateURL(ordNode)

		ordererNodes = append(
			ordererNodes,
			&ClusterOrdererNode{
				Name:       ordNode.FullName(),
				PublicURL:  publicURL,
				PrivateURL: privateURL,
				Spec:       ordNode.Spec,
				Status:     ordNode.Status,
				Item:       ordNode,
			},
		)
	}
	return ordererNodes, nil
}
func GetCertAuthByURL(clientSet *kubernetes.Clientset, oclient *operatorv1.Clientset, host string, port int) (*ClusterCA, error) {
	cahost := host
	ns := ""
	if strings.Contains(cahost, ".") && len(strings.Split(cahost, ".")) == 2 {
		chunks := strings.Split(cahost, ".")
		cahost = chunks[0]
		ns = chunks[1]
	}
	certAuths, err := GetClusterCAs(clientSet, oclient, ns)
	if err != nil {
		return nil, err
	}
	for _, certAuth := range certAuths {
		if // if host and port is specified by kubernetes DNS
		utils.Contains(certAuth.Spec.Hosts, host) || certAuth.Item.Name == cahost || (certAuth.Status.NodePort != 7054 && certAuth.Status.NodePort == port) {
			return certAuth, nil
		}

	}
	return nil, errors.Errorf("CA with host=%s port=%d not found", host, port)
}
func GetURLForCA(certAuth *ClusterCA) (string, error) {
	var host string
	var port int
	if len(certAuth.Spec.Istio.Hosts) > 0 {
		host = certAuth.Spec.Istio.Hosts[0]
		port = certAuth.Spec.Istio.Port
	} else {
		client, err := GetKubeClient()
		if err != nil {
			return "", err
		}
		host, err = utils.GetPublicIPKubernetes(client)
		if err != nil {
			return "", err
		}
		port = certAuth.Status.NodePort
	}
	return fmt.Sprintf("https://%s:%d", host, port), nil
}
func GetCertAuthByName(clientSet *kubernetes.Clientset, oclient *operatorv1.Clientset, name string, ns string) (*ClusterCA, error) {
	certAuths, err := GetClusterCAs(clientSet, oclient, "")
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

func GetOrderingServiceByFullName(clientSet *kubernetes.Clientset, oclient *operatorv1.Clientset, name string) (*ClusterOrderingService, error) {
	_, ordServices, err := GetClusterOrderers(clientSet, oclient, "")
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
func GetPeerByFullName(clientSet *kubernetes.Clientset, oclient *operatorv1.Clientset, name string) (*ClusterPeer, error) {
	_, peers, err := GetClusterPeers(clientSet, oclient, "")
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

type HostPort struct {
	Host string
	Port int
}

func GetOrdererPublicURL(clientset *kubernetes.Clientset, node hlfv1alpha1.FabricOrdererNode) (string, error) {
	hostPort, err := GetOrdererHostPort(clientset, node)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s:%d", hostPort.Host, hostPort.Port), nil
}
func GetOrdererHostAndPort(clientset *kubernetes.Clientset, nodeSpec hlfv1alpha1.FabricOrdererNodeSpec, nodeStatus hlfv1alpha1.FabricOrdererNodeStatus) (string, int, error) {
	hostName, err := utils.GetPublicIPKubernetes(clientset)
	if err != nil {
		return "", 0, err
	}
	ordererPort := nodeStatus.NodePort
	if len(nodeSpec.Istio.Hosts) > 0 {
		hostName = nodeSpec.Istio.Hosts[0]
		ordererPort = nodeSpec.Istio.Port
	}
	return hostName, ordererPort, nil
}
func GetPeerHostAndPort(clientset *kubernetes.Clientset, nodeSpec hlfv1alpha1.FabricPeerSpec, nodeStatus hlfv1alpha1.FabricPeerStatus) (string, int, error) {
	hostName, err := utils.GetPublicIPKubernetes(clientset)
	if err != nil {
		return "", 0, err
	}
	ordererPort := nodeStatus.NodePort
	if len(nodeSpec.Istio.Hosts) > 0 {
		hostName = nodeSpec.Istio.Hosts[0]
		ordererPort = nodeSpec.Istio.Port
	}
	return hostName, ordererPort, nil
}
func GetOrdererAdminHostAndPort(clientset *kubernetes.Clientset, nodeSpec hlfv1alpha1.FabricOrdererNodeSpec, nodeStatus hlfv1alpha1.FabricOrdererNodeStatus) (string, int, error) {
	hostName, err := utils.GetPublicIPKubernetes(clientset)
	if err != nil {
		return "", 0, err
	}
	ordererPort := nodeStatus.AdminPort
	if len(nodeSpec.AdminIstio.Hosts) > 0 {
		hostName = nodeSpec.AdminIstio.Hosts[0]
		ordererPort = nodeSpec.AdminIstio.Port
	}
	return hostName, ordererPort, nil
}
func GetOrdererPrivateURL(node hlfv1alpha1.FabricOrdererNode) string {
	return fmt.Sprintf("%s.%s:%s", node.Name, node.Namespace, "7050")
}

func GetCAPublicURL(clientset *kubernetes.Clientset, node hlfv1alpha1.FabricCA) (string, error) {
	hostPort, err := GetCAHostPort(clientset, node)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s:%d", hostPort.Host, hostPort.Port), nil
}
func GetCAPrivateURL(node hlfv1alpha1.FabricCA) string {
	return fmt.Sprintf("%s.%s:%s", node.Name, node.Namespace, "7054")
}

func GetPeerPublicURL(clientset *kubernetes.Clientset, node hlfv1alpha1.FabricPeer) (string, error) {
	hostPort, err := GetPeerHostPort(clientset, node)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s:%d", hostPort.Host, hostPort.Port), nil
}
func GetPeerHostPort(clientset *kubernetes.Clientset, node hlfv1alpha1.FabricPeer) (*HostPort, error) {
	k8sIP, err := utils.GetPublicIPKubernetes(clientset)
	if err != nil {
		return nil, err
	}
	if node.Spec.Istio != nil && len(node.Spec.Istio.Hosts) > 0 {
		return &HostPort{
			Host: node.Spec.Istio.Hosts[0],
			Port: node.Spec.Istio.Port,
		}, nil
	}
	return &HostPort{
		Host: k8sIP,
		Port: node.Status.NodePort,
	}, nil
}
func GetPeerPrivateURL(node hlfv1alpha1.FabricPeer) string {
	return fmt.Sprintf("%s.%s:%s", node.Name, node.Namespace, "7051")
}
func GetOrdererHostPort(clientset *kubernetes.Clientset, node hlfv1alpha1.FabricOrdererNode) (*HostPort, error) {
	k8sIP, err := utils.GetPublicIPKubernetes(clientset)
	if err != nil {
		return nil, err
	}
	if node.Spec.Istio != nil && len(node.Spec.Istio.Hosts) > 0 {
		return &HostPort{
			Host: node.Spec.Istio.Hosts[0],
			Port: node.Spec.Istio.Port,
		}, nil
	}
	return &HostPort{
		Host: k8sIP,
		Port: node.Status.NodePort,
	}, nil
}

func GetCAHostPort(clientset *kubernetes.Clientset, node hlfv1alpha1.FabricCA) (*HostPort, error) {
	k8sIP, err := utils.GetPublicIPKubernetes(clientset)
	if err != nil {
		return nil, err
	}
	if node.Spec.Istio != nil && len(node.Spec.Istio.Hosts) > 0 {
		return &HostPort{
			Host: node.Spec.Istio.Hosts[0],
			Port: node.Spec.Istio.Port,
		}, nil
	}
	return &HostPort{
		Host: k8sIP,
		Port: node.Status.NodePort,
	}, nil
}
func GetOrdererNodeByFullName(clientSet *kubernetes.Clientset, oclient *operatorv1.Clientset, name string) (*ClusterOrdererNode, error) {
	ordererNodes, err := GetClusterOrdererNodes(clientSet, oclient, "")
	if err != nil {
		return nil, err
	}
	for _, ordNode := range ordererNodes {
		if ordNode.Name == name {
			return ordNode, nil
		}
	}
	return nil, errors.Errorf("Orderer Node with name=%s not found", name)
}
func GetCertAuthByFullName(clientSet *kubernetes.Clientset, oclient *operatorv1.Clientset, name string) (*ClusterCA, error) {
	certAuths, err := GetClusterCAs(clientSet, oclient, "")
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
func GetClusterPeers(
	clientSet *kubernetes.Clientset,
	oclient *operatorv1.Clientset, ns string) ([]*Organization, []*ClusterPeer, error) {
	ctx := context.Background()

	peerResponse, err := oclient.HlfV1alpha1().FabricPeers(ns).List(ctx, v1.ListOptions{})
	if err != nil {
		return nil, nil, err
	}
	var peers []*ClusterPeer
	for _, peer := range peerResponse.Items {
		publicURL, err := GetPeerPublicURL(clientSet, peer)
		if err != nil {
			return nil, nil, err
		}
		privateURL := GetPeerPrivateURL(peer)
		peers = append(
			peers,
			&ClusterPeer{
				Name:       peer.FullName(),
				ObjectMeta: peer.ObjectMeta,
				Spec:       peer.Spec,
				Status:     peer.Status,
				Identity:   Identity{},
				PublicURL:  publicURL,
				PrivateURL: privateURL,
				MSPID:      peer.Spec.MspID,
			},
		)
	}
	orgMap := map[string]*Organization{}
	for _, peer := range peers {
		mspID := peer.Spec.MspID
		org, ok := orgMap[mspID]
		if !ok {
			orgMap[mspID] = &Organization{
				Type:         PeerType,
				MspID:        mspID,
				OrdererNodes: []*ClusterOrdererNode{},
				Peers:        []*ClusterPeer{},
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
