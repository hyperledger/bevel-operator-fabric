package ca

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"github.com/kfsoftware/hlf-operator/controllers/ordnode"
	operatorv1alpha1 "github.com/kfsoftware/hlf-operator/pkg/client/clientset/versioned"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"log"
	"net"
	"net/url"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/Masterminds/sprig"
	pb "github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/ccpackager/lifecycle"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/resource"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/resource/genesisconfig"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"

	"github.com/kfsoftware/hlf-operator/controllers/certs"
	"github.com/kfsoftware/hlf-operator/controllers/testutils"
	"github.com/kfsoftware/hlf-operator/controllers/utils"
	"github.com/kfsoftware/hlf-operator/internal/github.com/hyperledger/fabric/common/policydsl"
	"github.com/lithammer/shortuuid/v3"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	k8sconfig "sigs.k8s.io/controller-runtime/pkg/client/config"

	hlfv1alpha1 "github.com/kfsoftware/hlf-operator/api/hlf.kungfusoftware.es/v1alpha1"
	"github.com/kfsoftware/hlf-operator/controllers/ordservice"
	"github.com/kfsoftware/hlf-operator/controllers/peer"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"path/filepath"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	// +kubebuilder:scaffold:imports
)

var cfg *rest.Config
var K8sClient client.Client
var ClientSet *kubernetes.Clientset
var testEnv *envtest.Environment
var kubeInt kubernetes.Interface
var dynamicClient dynamic.Interface
var k8sManager ctrl.Manager
var RestConfig *rest.Config

func getOrderers(releaseName string, ns string) []hlfv1alpha1.FabricOrdererNode {
	hlfClient, err := operatorv1alpha1.NewForConfig(RestConfig)
	Expect(err).ToNot(HaveOccurred())
	ctx := context.Background()
	ordNodesRes, err := hlfClient.HlfV1alpha1().FabricOrdererNodes(ns).List(
		ctx,
		v1.ListOptions{
			LabelSelector: fmt.Sprintf("release=%s", releaseName),
		},
	)
	Expect(err).ToNot(HaveOccurred())
	return ordNodesRes.Items
}

var _ = BeforeSuite(func(done Done) {
	logf.SetLogger(zap.New(zap.UseDevMode(true), zap.WriteTo(GinkgoWriter)))

	By("bootstrapping test environment")
	useExistingCluster := true
	testEnv = &envtest.Environment{
		CRDDirectoryPaths:  []string{filepath.Join("..", "config", "crd", "bases")},
		UseExistingCluster: &useExistingCluster,
	}

	var err error
	cfg, err = testEnv.Start()
	Expect(err).ToNot(HaveOccurred())
	Expect(cfg).ToNot(BeNil())

	err = hlfv1alpha1.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred())

	// +kubebuilder:scaffold:scheme

	k8sManager, err = ctrl.NewManager(
		cfg,
		ctrl.Options{
			Scheme:             scheme.Scheme,
			MetricsBindAddress: "0",
		},
	)
	Expect(err).ToNot(HaveOccurred())
	RestConfig = k8sManager.GetConfig()
	ClientSet, err = kubernetes.NewForConfig(RestConfig)
	Expect(err).NotTo(HaveOccurred())
	caChartPath, err := filepath.Abs("../../charts/hlf-ca")
	Expect(err).ToNot(HaveOccurred())

	caReconciler := &FabricCAReconciler{
		Client:    k8sManager.GetClient(),
		Log:       ctrl.Log.WithName("controllers").WithName("FabricCA"),
		Scheme:    nil,
		Config:    RestConfig,
		ClientSet: ClientSet,
		ChartPath: caChartPath,
	}
	err = caReconciler.SetupWithManager(k8sManager)
	Expect(err).ToNot(HaveOccurred())
	peerChartPath, err := filepath.Abs("../../charts/hlf-peer")
	Expect(err).ToNot(HaveOccurred())
	peerReconciler := &peer.FabricPeerReconciler{
		Client:    k8sManager.GetClient(),
		Log:       ctrl.Log.WithName("controllers").WithName("FabricPeer"),
		Scheme:    nil,
		Config:    RestConfig,
		ChartPath: peerChartPath,
	}
	err = peerReconciler.SetupWithManager(k8sManager)

	Expect(err).ToNot(HaveOccurred())
	ordChartPath, err := filepath.Abs("../../charts/hlf-ord")
	Expect(err).ToNot(HaveOccurred())
	ordReconciler := ordservice.FabricOrderingServiceReconciler{
		Client:    k8sManager.GetClient(),
		Log:       ctrl.Log.WithName("controllers").WithName("FabricOrderingService"),
		Scheme:    nil,
		ChartPath: ordChartPath,
		Config:    RestConfig,
	}
	err = ordReconciler.SetupWithManager(k8sManager)
	Expect(err).ToNot(HaveOccurred())

	ordNodeChartPath, err := filepath.Abs("../../charts/hlf-ordnode")
	Expect(err).ToNot(HaveOccurred())
	ordNodeReconciler := ordnode.FabricOrdererNodeReconciler{
		Client:    k8sManager.GetClient(),
		Log:       ctrl.Log.WithName("controllers").WithName("FabricOrdererNode"),
		Scheme:    nil,
		ChartPath: ordNodeChartPath,
		Config:    RestConfig,
	}
	err = ordNodeReconciler.SetupWithManager(k8sManager)
	Expect(err).ToNot(HaveOccurred())

	go func() {
		err = k8sManager.Start(ctrl.SetupSignalHandler())
		Expect(err).ToNot(HaveOccurred())
	}()
	mgrSyncCtx, mgrSyncCtxCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer mgrSyncCtxCancel()
	if synced := k8sManager.GetCache().WaitForCacheSync(mgrSyncCtx.Done()); !synced {
		fmt.Println("Failed to sync")
	}
	K8sClient = k8sManager.GetClient()
	Expect(K8sClient).ToNot(BeNil())

	restConfig := k8sconfig.GetConfigOrDie()
	Expect(restConfig).ToNot(BeNil())

	kubeInt = kubernetes.NewForConfigOrDie(restConfig)
	Expect(kubeInt).ToNot(BeNil())

	dynamicClient = dynamic.NewForConfigOrDie(restConfig)
	Expect(dynamicClient).ToNot(BeNil())

	close(done)
}, 60)

var _ = AfterSuite(func() {
	By("tearing down the test environment")
	err := testEnv.Stop()
	Expect(err).ToNot(HaveOccurred())
})

const (
	defTimeoutSecs  = "120s"
	peerTimeoutSecs = "240s"
	defInterval     = "1s"
	systemChannelID = "system-channel"
)

func getClientForOrderer(updatedOrderer *hlfv1alpha1.FabricOrderingService, updatedCA *hlfv1alpha1.FabricCA) *resmgmt.Client {
	caurl := updatedCA.Status.URL
	caName := "ca"
	tlsCertString := updatedCA.Status.TlsCert
	caCert := updatedCA.Status.CACert
	enrollID := updatedCA.Spec.CA.Registry.Identities[0].Name
	enrollSecret := updatedCA.Spec.CA.Registry.Identities[0].Pass

	mspID := updatedOrderer.Spec.MspID
	userClient := "org1admin"
	userClientPW := "org1adminpw"
	registerParams := certs.RegisterUserRequest{
		TLSCert:      tlsCertString,
		URL:          caurl,
		Name:         caName,
		MSPID:        mspID,
		EnrollID:     enrollID,
		EnrollSecret: enrollSecret,
		User:         userClient,
		Secret:       userClientPW,
		Type:         "admin",
		Attributes:   nil,
	}
	_, err := certs.RegisterUser(registerParams)

	crt, pk, rootCrt, err := certs.EnrollUser(certs.EnrollUserRequest{
		TLSCert: tlsCertString,
		URL:     caurl,
		Name:    caName,
		MSPID:   mspID,
		User:    userClient,
		Secret:  userClientPW,
		Profile: "",
	})
	Expect(err).ToNot(HaveOccurred())
	Expect(crt).ToNot(BeNil())
	Expect(pk).ToNot(BeNil())
	Expect(rootCrt).ToNot(BeNil())

	certPem := string(utils.EncodeX509Certificate(crt))
	pkPem, err := utils.EncodePrivateKey(pk)
	Expect(err).ToNot(HaveOccurred())

	configYamlTpl := `
name: test-network-org1
version: 1.0.0
client:
  organization: OrdererMSP
  connection:
    timeout:
      peer:
        endorser: "300"
organizations:
  OrdererMSP:
    mspid: OrdererMSP
    cryptoPath: /tmp/cryptopath
    users:
      admin:
        key:
          pem: |
{{ .AdminKey | indent 12 }}

        cert:
          pem: |
{{ .AdminCert | indent 12 }}

    certificateAuthorities: []
    peers: []
    orderers:
    - "orderer"
certificateAuthorities: {}
peers: {}
orderers:
  "orderer":
    url: {{ .OrdUrl }}
    grpcOptions:
      allow-insecure: true
    tlsCACerts:
      pem: |
{{ .TlsCACrt | indent 8 }}

channels: {}
`
	tmpl, err := template.New("test").Funcs(sprig.HermeticTxtFuncMap()).Parse(configYamlTpl)
	var buf bytes.Buffer
	ordNodes := getOrderers(
		updatedOrderer.Name,
		updatedOrderer.Namespace,
	)
	ordNode := ordNodes[0]
	ordURL := ordNode.Status.URL
	Expect(err).ToNot(HaveOccurred())
	err = tmpl.Execute(&buf, map[string]interface{}{
		"AdminKey":  string(pkPem),
		"AdminCert": certPem,
		"TlsCACrt":  caCert,
		"OrdUrl":    ordURL,
	})
	configYaml := buf.Bytes()
	log.Print(string(configYaml))
	configBackend := config.FromRaw(configYaml, "yaml")
	sdk, err := fabsdk.New(configBackend)
	Expect(err).ToNot(HaveOccurred())
	org1AdminClientContext := sdk.Context(
		fabsdk.WithUser("admin"),
		fabsdk.WithOrg("OrdererMSP"),
	)
	ordClient, err := resmgmt.New(org1AdminClientContext)
	Expect(err).ToNot(HaveOccurred())
	return ordClient
}
func randomFabricCA(releaseName string, namespace string) *hlfv1alpha1.FabricCA {
	fabricCARegistry := hlfv1alpha1.FabricCARegistry{
		MaxEnrollments: -1,
		Identities: []hlfv1alpha1.FabricCAIdentity{
			{
				Name:        "enroll",
				Pass:        "enrollpw",
				Type:        "client",
				Affiliation: "",
				Attrs: hlfv1alpha1.FabricCAIdentityAttrs{
					RegistrarRoles: "*",
					DelegateRoles:  "*",
					Attributes:     "*",
					Revoker:        true,
					IntermediateCA: true,
					GenCRL:         true,
					AffiliationMgr: true,
				},
			},
		},
	}
	subject := hlfv1alpha1.FabricCASubject{
		CN: "ca",
		C:  "ES",
		ST: "Alicante",
		O:  "Kung Fu Software",
		L:  "Alicante",
		OU: "Tech",
	}
	tlsSubject := hlfv1alpha1.FabricCASubject{
		CN: "tlsca",
		C:  "ES",
		ST: "Alicante",
		O:  "Kung Fu Software",
		L:  "Alicante",
		OU: "Tech",
	}
	cabccsp := hlfv1alpha1.FabricCABCCSP{
		Default: "SW",
		SW: hlfv1alpha1.FabricCABCCSPSW{
			Hash:     "SHA2",
			Security: "256",
		},
	}
	cfg := hlfv1alpha1.FabricCACFG{
		Identities:   hlfv1alpha1.FabricCACFGIdentities{AllowRemove: true},
		Affiliations: hlfv1alpha1.FabricCACFGAffilitions{AllowRemove: true},
	}
	fabricCa := &hlfv1alpha1.FabricCA{
		TypeMeta: NewTypeMeta("FabricCA"),
		ObjectMeta: v1.ObjectMeta{
			Name:      releaseName,
			Namespace: namespace,
		},
		Spec: hlfv1alpha1.FabricCASpec{
			Database: hlfv1alpha1.FabricCADatabase{
				Type:       "sqlite3",
				Datasource: "fabric-ca-server.db",
			},
			Hosts: []string{
				"localhost",
				releaseName,
				fmt.Sprintf("%s.%s", releaseName, namespace),
			},
			Service: hlfv1alpha1.FabricCASpecService{
				ServiceType: "NodePort",
			},
			CLRSizeLimit: 512000,
			Image:        "hyperledger/fabric-ca",
			Version:      "1.4.9",
			Debug:        true,
			TLS:          hlfv1alpha1.FabricCATLSConf{Subject: subject},
			CA: hlfv1alpha1.FabricCAItemConf{
				Name:    "ca",
				Subject: subject,
				CFG:     cfg,
				CSR: hlfv1alpha1.FabricCACSR{
					CN:    "ca",
					Hosts: []string{"localhost"},
					Names: []hlfv1alpha1.FabricCANames{
						{
							C:  "US",
							ST: "",
							O:  "Hyperledger",
							L:  "",
							OU: "North Carolina",
						},
					},
					CA: hlfv1alpha1.FabricCACSRCA{
						Expiry:     "131400h",
						PathLength: 0,
					},
				},
				CRL:          hlfv1alpha1.FabricCACRL{Expiry: "24h"},
				Registry:     fabricCARegistry,
				Intermediate: hlfv1alpha1.FabricCAIntermediate{},
				BCCSP:        cabccsp,
			},
			TLSCA: hlfv1alpha1.FabricCAItemConf{
				Name:    "tlsca",
				Subject: tlsSubject,
				CFG:     cfg,
				CSR: hlfv1alpha1.FabricCACSR{
					CN:    "tlsca",
					Hosts: []string{"localhost"},
					Names: []hlfv1alpha1.FabricCANames{
						{
							C:  "US",
							ST: "",
							O:  "Hyperledger",
							L:  "",
							OU: "North Carolina",
						},
					},
					CA: hlfv1alpha1.FabricCACSRCA{
						Expiry:     "131400h",
						PathLength: 0,
					},
				},
				CRL:          hlfv1alpha1.FabricCACRL{Expiry: "24h"},
				Registry:     fabricCARegistry,
				Intermediate: hlfv1alpha1.FabricCAIntermediate{},
				BCCSP:        cabccsp,
			},
			Cors: hlfv1alpha1.Cors{
				Enabled: false,
				Origins: []string{},
			},
			Resources: hlfv1alpha1.Resources{
				Requests: hlfv1alpha1.Requests{
					CPU:    "10m",
					Memory: "256Mi",
				},
				Limits: hlfv1alpha1.RequestsLimit{
					CPU:    "2",
					Memory: "4Gi",
				},
			},
			Storage: hlfv1alpha1.Storage{
				Size:         "3Gi",
				StorageClass: "",
				AccessMode:   "ReadWriteOnce",
			},
			Metrics: hlfv1alpha1.FabricCAMetrics{
				Provider: "prometheus",
				Statsd: hlfv1alpha1.FabricCAMetricsStatsd{
					Network:       "udp",
					Address:       "127.0.0.1:8125",
					WriteInterval: "10s",
					Prefix:        "server",
				},
			},
		},
	}
	Expect(K8sClient.Create(context.Background(), fabricCa)).Should(Succeed())
	updatedCA := &hlfv1alpha1.FabricCA{}
	caKey := types.NamespacedName{Namespace: namespace, Name: releaseName}
	Eventually(
		func() bool {
			err := K8sClient.Get(context.Background(), caKey, updatedCA)
			if err != nil {
				return false
			}
			ctrl.Log.WithName("test").Info("after update", "updatedCA", updatedCA)
			return updatedCA.Status.Status == hlfv1alpha1.RunningStatus
		},
		defTimeoutSecs,
		defInterval,
	).Should(BeTrue(), "ca status should have been updated")
	return updatedCA
}
func getSDKForPeer(peer *hlfv1alpha1.FabricPeer, ca *hlfv1alpha1.FabricCA) *fabsdk.FabricSDK {
	userClient := "orgadmin"
	userClientPW := "orgadminpw"
	caName := "ca"
	enrollID := ca.Spec.CA.Registry.Identities[0].Name
	caURL := ca.Status.URL
	enrollSecret := ca.Spec.CA.Registry.Identities[0].Pass
	caCert := ca.Status.CACert
	tlsCertString := ca.Status.TlsCert
	mspID := peer.Spec.MspID
	registerParams := certs.RegisterUserRequest{
		TLSCert:      tlsCertString,
		URL:          caURL,
		Name:         caName,
		MSPID:        mspID,
		EnrollID:     enrollID,
		EnrollSecret: enrollSecret,
		User:         userClient,
		Secret:       userClientPW,
		Type:         "admin",
		Attributes:   nil,
	}
	_, err := certs.RegisterUser(registerParams)

	crt, pk, rootCrt, err := certs.EnrollUser(certs.EnrollUserRequest{
		TLSCert: tlsCertString,
		URL:     caURL,
		Name:    caName,
		MSPID:   mspID,
		User:    userClient,
		Secret:  userClientPW,
		Profile: "",
	})
	Expect(err).ToNot(HaveOccurred())
	Expect(crt).ToNot(BeNil())
	Expect(pk).ToNot(BeNil())
	Expect(rootCrt).ToNot(BeNil())
	certPem := string(utils.EncodeX509Certificate(crt))
	pkPem, err := utils.EncodePrivateKey(pk)
	Expect(err).ToNot(HaveOccurred())
	configYamlTpl := `
name: test-network-org1
version: 1.0.0
client:
  organization: {{.MSPID}}
  connection:
    timeout:
      peer:
        endorser: "300"
organizations:
  {{.MSPID}}:
    mspid: {{.MSPID}}
    cryptoPath: /tmp/cryptopath
    users:
      admin:
        key:
          pem: |
{{ .AdminKey | indent 12 }}
        cert:
          pem: |
{{ .AdminCert | indent 12 }}

    certificateAuthorities: []
    peers:
    - "peer"
certificateAuthorities: {}
peers:
  "peer":
    url: {{ .PeerUrl }}
    grpcOptions:
      hostnameOverride: ""
      ssl-target-name-override: ""
      allow-insecure: true
    tlsCACerts:
      pem: |
{{ .TlsCACrt | indent 8 }}

channels:
  _default:
    peers:
      "peer":
        endorsingPeer: true
        chaincodeQuery: true
        ledgerQuery: true
        eventSource: true


`
	tmpl, err := template.New("test").Funcs(sprig.HermeticTxtFuncMap()).Parse(configYamlTpl)
	var buf bytes.Buffer
	peerURL := peer.Status.URL
	Expect(err).ToNot(HaveOccurred())
	err = tmpl.Execute(&buf, map[string]interface{}{
		"MSPID":     peer.Spec.MspID,
		"AdminKey":  string(pkPem),
		"AdminCert": certPem,
		"TlsCACrt":  caCert,
		"PeerUrl":   peerURL,
	})
	configYaml := buf.Bytes()
	configBackend := config.FromRaw(configYaml, "yaml")
	sdk, err := fabsdk.New(
		configBackend,
	)
	Expect(err).ToNot(HaveOccurred())
	return sdk
}
func getSDKForPeerWithOrderer(peer *hlfv1alpha1.FabricPeer, ca *hlfv1alpha1.FabricCA, orderer *hlfv1alpha1.FabricOrderingService, ordererCA *hlfv1alpha1.FabricCA) *fabsdk.FabricSDK {
	userClient := "orgadmin"
	userClientPW := "orgadminpw"
	caName := "ca"
	enrollID := ca.Spec.CA.Registry.Identities[0].Name
	caURL := ca.Status.URL
	enrollSecret := ca.Spec.CA.Registry.Identities[0].Pass
	caCert := ca.Status.CACert
	tlsCertString := ca.Status.TlsCert
	mspID := peer.Spec.MspID
	registerParams := certs.RegisterUserRequest{
		TLSCert:      tlsCertString,
		URL:          caURL,
		Name:         caName,
		MSPID:        mspID,
		EnrollID:     enrollID,
		EnrollSecret: enrollSecret,
		User:         userClient,
		Secret:       userClientPW,
		Type:         "admin",
		Attributes:   nil,
	}
	_, err := certs.RegisterUser(registerParams)

	crt, pk, rootCrt, err := certs.EnrollUser(certs.EnrollUserRequest{
		TLSCert: tlsCertString,
		URL:     caURL,
		Name:    caName,
		MSPID:   mspID,
		User:    userClient,
		Secret:  userClientPW,
		Profile: "",
	})
	Expect(err).ToNot(HaveOccurred())
	Expect(crt).ToNot(BeNil())
	Expect(pk).ToNot(BeNil())
	Expect(rootCrt).ToNot(BeNil())
	certPem := string(utils.EncodeX509Certificate(crt))
	pkPem, err := utils.EncodePrivateKey(pk)
	Expect(err).ToNot(HaveOccurred())
	configYamlTpl := `
name: test-network-org1
version: 1.0.0
client:
  organization: {{.MSPID}}
  connection:
    timeout:
      peer:
        endorser: "300"
organizations:
  {{.MSPID}}:
    mspid: {{.MSPID}}
    cryptoPath: /tmp/cryptopath
    users:
      admin:
        key:
          pem: |
{{ .AdminKey | indent 12 }}
        cert:
          pem: |
{{ .AdminCert | indent 12 }}
    certificateAuthorities: []
    peers:
    - "peer"
certificateAuthorities: {}
peers:
  "peer":
    url: {{ .PeerUrl }}
    grpcOptions:
      allow-insecure: true
    tlsCACerts:
      pem: |
{{ .TlsCACrt | indent 8 }}
orderers:
  "orderer":
    url: {{ .OrdUrl }}
    grpcOptions:
      allow-insecure: true
    tlsCACerts:
      pem: |
{{ .OrdTlsCACrt | indent 8 }}

channels:
  _default:
    peers:
      "peer":
        endorsingPeer: true
        chaincodeQuery: true
        ledgerQuery: true
        eventSource: true
    orderers: ["orderer"]

`
	tmpl, err := template.New("test").Funcs(sprig.HermeticTxtFuncMap()).Parse(configYamlTpl)
	var buf bytes.Buffer
	peerURL := peer.Status.URL
	Expect(err).ToNot(HaveOccurred())
	ordNodes := getOrderers(
		orderer.Name,
		orderer.Namespace,
	)
	ordNode := ordNodes[0]
	err = tmpl.Execute(&buf, map[string]interface{}{
		"MSPID":       peer.Spec.MspID,
		"AdminKey":    string(pkPem),
		"AdminCert":   certPem,
		"TlsCACrt":    caCert,
		"PeerUrl":     peerURL,
		"OrdTlsCACrt": ordererCA.Status.CACert,
		"OrdUrl":      ordNode.Status.URL,
	})
	configYaml := buf.Bytes()
	configBackend := config.FromRaw(configYaml, "yaml")
	sdk, err := fabsdk.New(
		configBackend,
	)
	Expect(err).ToNot(HaveOccurred())
	return sdk
}
func getClientForPeer(peer *hlfv1alpha1.FabricPeer, ca *hlfv1alpha1.FabricCA) *resmgmt.Client {
	sdk := getSDKForPeer(peer, ca)
	org1AdminClientContext := sdk.Context(
		fabsdk.WithUser("admin"),
		fabsdk.WithOrg(peer.Spec.MspID),
	)
	resClient, err := resmgmt.New(org1AdminClientContext)
	Expect(err).ToNot(HaveOccurred())
	return resClient
}
func getClientForPeerWithOrderer(peer *hlfv1alpha1.FabricPeer, ca *hlfv1alpha1.FabricCA, orderer *hlfv1alpha1.FabricOrderingService, ordererCA *hlfv1alpha1.FabricCA) *resmgmt.Client {
	sdk := getSDKForPeerWithOrderer(peer, ca, orderer, ordererCA)
	org1AdminClientContext := sdk.Context(
		fabsdk.WithUser("admin"),
		fabsdk.WithOrg(peer.Spec.MspID),
	)
	resClient, err := resmgmt.New(org1AdminClientContext)
	Expect(err).ToNot(HaveOccurred())
	return resClient
}

type createPeerParams struct {
	MSPID string
}

func createPeer(releaseName string, namespace string, params createPeerParams, certauth *hlfv1alpha1.FabricCA) *hlfv1alpha1.FabricPeer {
	publicIP, err := utils.GetPublicIPKubernetes(ClientSet)
	Expect(err).ToNot(HaveOccurred())
	mspID := params.MSPID
	caHost := certauth.Status.Host
	caPort := certauth.Status.Port
	caName := "ca"
	caTLSCert := certauth.Status.TlsCert
	enrollID := certauth.Spec.CA.Registry.Identities[0].Name
	enrollSecret := certauth.Spec.CA.Registry.Identities[0].Pass
	caURL := certauth.Status.URL

	registerParams := certs.RegisterUserRequest{
		TLSCert:      caTLSCert,
		URL:          caURL,
		Name:         caName,
		MSPID:        mspID,
		EnrollID:     enrollID,
		EnrollSecret: enrollSecret,
		User:         "peer",
		Secret:       "peerpw",
		Type:         "peer",
		Attributes:   nil,
	}
	secret, err := certs.RegisterUser(registerParams)
	Expect(err).ToNot(HaveOccurred())
	Expect(secret).To(Equal(registerParams.Secret))
	peerEnrollID := "peer"
	peerEnrollSecret := "peerpw"
	hosts := []string{
		"127.0.0.1",
		publicIP,
	}
	fabricPeer := &hlfv1alpha1.FabricPeer{
		TypeMeta: NewTypeMeta("FabricPeer"),
		ObjectMeta: v1.ObjectMeta{
			Name:      releaseName,
			Namespace: namespace,
		},
		Spec: hlfv1alpha1.FabricPeerSpec{
			DockerSocketPath: "/var/run/docker.sock",
			Image:            "quay.io/kfsoftware/fabric-peer",
			Istio: hlfv1alpha1.FabricPeerIstio{
				Port: 443,
			},
			Gossip: hlfv1alpha1.FabricPeerSpecGossip{
				ExternalEndpoint:  "",
				Bootstrap:         "",
				Endpoint:          "",
				UseLeaderElection: true,
				OrgLeader:         false,
			},
			ExternalEndpoint:         "",
			Tag:                      "amd64-2.3.0",
			ExternalChaincodeBuilder: true,
			CouchDB: hlfv1alpha1.FabricPeerCouchDB{
				User:     "couchdb",
				Password: "couchdb",
			},
			MspID: mspID,
			Secret: hlfv1alpha1.Secret{
				Enrollment: hlfv1alpha1.Enrollment{
					Component: hlfv1alpha1.Component{
						Cahost: caHost,
						Caname: caName,
						Caport: caPort,
						Catls: hlfv1alpha1.Catls{
							Cacert: base64.StdEncoding.EncodeToString([]byte(caTLSCert)),
						},
						Enrollid:     peerEnrollID,
						Enrollsecret: peerEnrollSecret,
					},
					TLS: hlfv1alpha1.TLS{
						Cahost: caHost,
						Caname: caName,
						Caport: caPort,
						Catls: hlfv1alpha1.Catls{
							Cacert: base64.StdEncoding.EncodeToString([]byte(caTLSCert)),
						},
						Csr: hlfv1alpha1.Csr{
							Hosts: hosts,
							CN:    "",
						},
						Enrollid:     peerEnrollID,
						Enrollsecret: peerEnrollSecret,
					},
				},
			},
			Service: hlfv1alpha1.PeerService{
				Type: "NodePort",
			},
			StateDb: "leveldb",
			Storage: hlfv1alpha1.FabricPeerStorage{
				CouchDB: hlfv1alpha1.Storage{
					Size:         "1Gi",
					StorageClass: "",
					AccessMode:   "ReadWriteOnce",
				},
				Peer: hlfv1alpha1.Storage{
					Size:         "1Gi",
					StorageClass: "",
					AccessMode:   "ReadWriteOnce",
				},
				Chaincode: hlfv1alpha1.Storage{
					Size:         "1Gi",
					StorageClass: "",
					AccessMode:   "ReadWriteOnce",
				},
			},
			Discovery: hlfv1alpha1.FabricPeerDiscovery{
				Period:      "60s",
				TouchPeriod: "60s",
			},
			Logging: hlfv1alpha1.FabricPeerLogging{
				Level:    "info",
				Peer:     "info",
				Cauthdsl: "info",
				Gossip:   "info",
				Grpc:     "info",
				Ledger:   "info",
				Msp:      "info",
				Policies: "info",
			},
			Resources: hlfv1alpha1.FabricPeerResources{
				Peer: hlfv1alpha1.Resources{
					Requests: hlfv1alpha1.Requests{
						CPU:    "10m",
						Memory: "10M",
					},
					Limits: hlfv1alpha1.RequestsLimit{
						CPU:    "2",
						Memory: "4096M",
					},
				},
				CouchDB: hlfv1alpha1.Resources{
					Requests: hlfv1alpha1.Requests{
						CPU:    "10m",
						Memory: "10M",
					},
					Limits: hlfv1alpha1.RequestsLimit{
						CPU:    "2",
						Memory: "4096M",
					},
				},
				Chaincode: hlfv1alpha1.Resources{
					Requests: hlfv1alpha1.Requests{
						CPU:    "10m",
						Memory: "10M",
					},
					Limits: hlfv1alpha1.RequestsLimit{
						CPU:    "2",
						Memory: "4096M",
					},
				},
			},
			Hosts:          []string{},
			OperationHosts: []string{},
			OperationIPs:   []string{},
		},
	}
	Expect(K8sClient.Create(context.Background(), fabricPeer)).Should(Succeed())
	return fabricPeer
}

type createOrdererParams struct {
	MSPID string
}

func createOrderer(releaseName string, namespace string, params createOrdererParams, certauth *hlfv1alpha1.FabricCA) *hlfv1alpha1.FabricOrderingService {
	mspID := params.MSPID
	By("create a fabric orderer")
	caHost := certauth.Status.Host
	caPort := certauth.Status.Port
	caName := "ca"
	caTLSCert := certauth.Status.TlsCert
	enrollID := certauth.Spec.CA.Registry.Identities[0].Name
	enrollSecret := certauth.Spec.CA.Registry.Identities[0].Pass
	caURL := fmt.Sprintf("https://%s:%d", caHost, caPort)
	ordEnrollID := "orderer"
	ordEnrollSecret := "ordererpw"
	ordType := "orderer"
	certs.RegisterUser(certs.RegisterUserRequest{
		TLSCert:      caTLSCert,
		URL:          caURL,
		Name:         caName,
		MSPID:        mspID,
		EnrollID:     enrollID,
		EnrollSecret: enrollSecret,
		User:         ordEnrollID,
		Secret:       ordEnrollSecret,
		Type:         ordType,
		Attributes:   nil,
	})

	fabricOrderer := &hlfv1alpha1.FabricOrderingService{
		TypeMeta: NewTypeMeta("FabricOrderingService"),
		ObjectMeta: v1.ObjectMeta{
			Name:      releaseName,
			Namespace: namespace,
		},
		Spec: hlfv1alpha1.FabricOrderingServiceSpec{
			Storage: hlfv1alpha1.Storage{
				Size:         "30Gi",
				StorageClass: "standard",
				AccessMode:   "ReadWriteOnce",
			},
			SystemChannel: hlfv1alpha1.OrdererSystemChannel{
				Name: systemChannelID,
				Config: hlfv1alpha1.ChannelConfig{
					BatchTimeout:            "2s",
					MaxMessageCount:         500,
					AbsoluteMaxBytes:        10 * 1024 * 1024,
					PreferredMaxBytes:       2 * 1024 * 1024,
					OrdererCapabilities:     hlfv1alpha1.OrdererCapabilities{V2_0: true},
					ApplicationCapabilities: hlfv1alpha1.ApplicationCapabilities{V2_0: true},
					ChannelCapabilities:     hlfv1alpha1.ChannelCapabilities{V2_0: true},
					SnapshotIntervalSize:    19,
					TickInterval:            "500ms",
					ElectionTick:            10,
					HeartbeatTick:           1,
					MaxInflightBlocks:       5,
				},
			},
			Nodes: []hlfv1alpha1.OrdererNode{
				{
					ID:   "orderer0",
					Host: "",
					Port: 0,
					Enrollment: hlfv1alpha1.OrdererNodeEnrollment{
						TLS: hlfv1alpha1.OrdererNodeEnrollmentTLS{
							Csr: hlfv1alpha1.Csr{
								Hosts: []string{"orderer0.example.com"},
							},
						},
					},
				},
			},
			Image: "hyperledger/fabric-orderer",
			Tag:   "amd64-2.3.0",
			MspID: mspID,
			Enrollment: hlfv1alpha1.OrdererEnrollment{
				Component: hlfv1alpha1.Component{
					Cahost: caHost,
					Caname: caName,
					Caport: caPort,
					Catls: hlfv1alpha1.Catls{
						Cacert: base64.StdEncoding.EncodeToString([]byte(caTLSCert)),
					},
					Enrollid:     enrollID,
					Enrollsecret: enrollSecret,
				},
				TLS: hlfv1alpha1.TLS{
					Cahost: caHost,
					Caname: caName,
					Caport: caPort,
					Catls: hlfv1alpha1.Catls{
						Cacert: base64.StdEncoding.EncodeToString([]byte(caTLSCert)),
					},
					Enrollid:     enrollID,
					Enrollsecret: enrollSecret,
					Csr: hlfv1alpha1.Csr{
						Hosts: []string{},
						CN:    "",
					},
				},
			},
			Service: hlfv1alpha1.OrdererService{
				Type: "NodePort",
			},
		},
	}
	Expect(K8sClient.Create(context.Background(), fabricOrderer)).Should(Succeed())
	return fabricOrderer
}

func verifyCADeployment(fabricCA *hlfv1alpha1.FabricCA) {
	caurl := fabricCA.Status.URL
	caName := "ca"
	tlsCertString := fabricCA.Status.TlsCert
	enrollID := fabricCA.Spec.CA.Registry.Identities[0].Name
	enrollSecret := fabricCA.Spec.CA.Registry.Identities[0].Pass
	By("enroll the admin user")
	mspID := "Org1MSP"
	tlsCert, tlsKey, tlsRootCert, err := certs.EnrollUser(certs.EnrollUserRequest{
		TLSCert: tlsCertString,
		URL:     caurl,
		Name:    caName,
		MSPID:   mspID,
		User:    enrollID,
		Secret:  enrollSecret,
	})
	Expect(err).ToNot(HaveOccurred())
	Expect(tlsCert).ToNot(BeNil())
	Expect(tlsKey).ToNot(BeNil())
	Expect(tlsRootCert).ToNot(BeNil())

	registerParams := certs.RegisterUserRequest{
		TLSCert:      tlsCertString,
		URL:          caurl,
		Name:         caName,
		MSPID:        mspID,
		EnrollID:     enrollID,
		EnrollSecret: enrollSecret,
		User:         "peer",
		Secret:       "peerpw",
		Type:         "peer",
		Attributes:   nil,
	}
	secret, err := certs.RegisterUser(registerParams)
	Expect(err).ToNot(HaveOccurred())
	Expect(secret).To(Equal(registerParams.Secret))
}

var _ = Describe("Fabric Controllers", func() {
	FabricNamespace := ""
	BeforeEach(func() {
		FabricNamespace = "hlf-operator-" + getRandomChannelID()
		testNamespace := &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: FabricNamespace,
			},
		}
		Expect(K8sClient.Create(context.Background(), testNamespace)).Should(Succeed())
	})
	AfterEach(func() {
		Expect(ClientSet.CoreV1().Namespaces().Delete(context.Background(), FabricNamespace, v1.DeleteOptions{})).Should(Succeed())
	})
	Specify("create a new Fabric CA instance", func() {
		By("create the CA object")
		objName := "org1-ca"
		fabricCARegistry := hlfv1alpha1.FabricCARegistry{
			MaxEnrollments: -1,
			Identities: []hlfv1alpha1.FabricCAIdentity{
				{
					Name:        "enroll",
					Pass:        "enrollpw",
					Type:        "client",
					Affiliation: "",
					Attrs: hlfv1alpha1.FabricCAIdentityAttrs{
						RegistrarRoles: "*",
						DelegateRoles:  "*",
						Attributes:     "*",
						Revoker:        true,
						IntermediateCA: true,
						GenCRL:         true,
						AffiliationMgr: true,
					},
				},
			},
		}
		subject := hlfv1alpha1.FabricCASubject{
			CN: "ca",
			C:  "ES",
			ST: "Alicante",
			O:  "Kung Fu Software",
			L:  "Alicante",
			OU: "Tech",
		}
		tlsSubject := hlfv1alpha1.FabricCASubject{
			CN: "tlsca",
			C:  "ES",
			ST: "Alicante",
			O:  "Kung Fu Software",
			L:  "Alicante",
			OU: "Tech",
		}
		cabccsp := hlfv1alpha1.FabricCABCCSP{
			Default: "SW",
			SW: hlfv1alpha1.FabricCABCCSPSW{
				Hash:     "SHA2",
				Security: "256",
			},
		}
		cfg := hlfv1alpha1.FabricCACFG{
			Identities:   hlfv1alpha1.FabricCACFGIdentities{AllowRemove: true},
			Affiliations: hlfv1alpha1.FabricCACFGAffilitions{AllowRemove: true},
		}
		fabricCa := &hlfv1alpha1.FabricCA{
			TypeMeta: NewTypeMeta("FabricCA"),
			ObjectMeta: v1.ObjectMeta{
				Name:      objName,
				Namespace: FabricNamespace,
			},
			Spec: hlfv1alpha1.FabricCASpec{
				Database: hlfv1alpha1.FabricCADatabase{
					Type:       "sqlite3",
					Datasource: "fabric-ca-server.db",
				},
				Hosts: []string{
					"localhost",
					objName,
					fmt.Sprintf("%s.%s", objName, FabricNamespace),
				},
				Service: hlfv1alpha1.FabricCASpecService{
					ServiceType: "NodePort",
				},
				CLRSizeLimit: 512000,
				Image:        "hyperledger/fabric-ca",
				Version:      "1.4.9",
				Debug:        true,
				TLS:          hlfv1alpha1.FabricCATLSConf{Subject: subject},
				CA: hlfv1alpha1.FabricCAItemConf{
					Name:    "ca",
					Subject: subject,
					CFG:     cfg,
					CSR: hlfv1alpha1.FabricCACSR{
						CN:    "ca",
						Hosts: []string{"localhost"},
						Names: []hlfv1alpha1.FabricCANames{
							{
								C:  "US",
								ST: "",
								O:  "Hyperledger",
								L:  "",
								OU: "North Carolina",
							},
						},
						CA: hlfv1alpha1.FabricCACSRCA{
							Expiry:     "131400h",
							PathLength: 0,
						},
					},
					CRL:          hlfv1alpha1.FabricCACRL{Expiry: "24h"},
					Registry:     fabricCARegistry,
					Intermediate: hlfv1alpha1.FabricCAIntermediate{},
					BCCSP:        cabccsp,
				},
				TLSCA: hlfv1alpha1.FabricCAItemConf{
					Name:    "tlsca",
					Subject: tlsSubject,
					CFG:     cfg,
					CSR: hlfv1alpha1.FabricCACSR{
						CN:    "tlsca",
						Hosts: []string{"localhost"},
						Names: []hlfv1alpha1.FabricCANames{
							{
								C:  "US",
								ST: "",
								O:  "Hyperledger",
								L:  "",
								OU: "North Carolina",
							},
						},
						CA: hlfv1alpha1.FabricCACSRCA{
							Expiry:     "131400h",
							PathLength: 0,
						},
					},
					CRL:          hlfv1alpha1.FabricCACRL{Expiry: "24h"},
					Registry:     fabricCARegistry,
					Intermediate: hlfv1alpha1.FabricCAIntermediate{},
					BCCSP:        cabccsp,
				},
				Cors: hlfv1alpha1.Cors{
					Enabled: false,
					Origins: []string{},
				},
				Resources: hlfv1alpha1.Resources{
					Requests: hlfv1alpha1.Requests{
						CPU:    "10m",
						Memory: "256Mi",
					},
					Limits: hlfv1alpha1.RequestsLimit{
						CPU:    "2",
						Memory: "4Gi",
					},
				},
				Storage: hlfv1alpha1.Storage{
					Size:         "3Gi",
					StorageClass: "",
					AccessMode:   "ReadWriteOnce",
				},
				Metrics: hlfv1alpha1.FabricCAMetrics{
					Provider: "prometheus",
					Statsd: hlfv1alpha1.FabricCAMetricsStatsd{
						Network:       "udp",
						Address:       "127.0.0.1:8125",
						WriteInterval: "10s",
						Prefix:        "server",
					},
				},
			},
		}
		Expect(K8sClient.Create(context.Background(), fabricCa)).Should(Succeed())
		updatedCA := &hlfv1alpha1.FabricCA{}
		caKey := types.NamespacedName{Namespace: FabricNamespace, Name: objName}
		Eventually(
			func() bool {
				err := K8sClient.Get(context.Background(), caKey, updatedCA)
				if err != nil {
					return false
				}
				ctrl.Log.WithName("test").Info("after update", "updatedCA", updatedCA)
				return updatedCA.Status.Status == hlfv1alpha1.RunningStatus
			},
			defTimeoutSecs,
			defInterval,
		).Should(BeTrue(), "ca status should have been updated")
		verifyCADeployment(updatedCA)
		Expect(K8sClient.Delete(context.Background(), updatedCA)).Should(Succeed())
		Eventually(
			func() bool {
				depName := GetDeploymentName(objName)
				_, err := ClientSet.AppsV1().Deployments(FabricNamespace).Get(context.Background(), depName, v1.GetOptions{})
				if err != nil {
					return apierrors.IsNotFound(err)
				} else {
					return false
				}
			},
			defTimeoutSecs,
			defInterval,
		).Should(BeTrue(), "ca deployment should have been deleted")

	})
	Specify("create a new Fabric Peer instance", func() {
		releaseNameCA := "org1-ca"
		releaseNamePeer := "org1-peer"
		mspID := "Org1MSP"
		By("create a fabric ca")
		updatedCA := randomFabricCA(releaseNameCA, FabricNamespace)
		Expect(updatedCA).ToNot(BeNil())
		By("create a fabric peer with leveldb")
		params := createPeerParams{
			MSPID: mspID,
		}
		createPeer(
			releaseNamePeer,
			FabricNamespace,
			params,
			updatedCA,
		)
		peer := &hlfv1alpha1.FabricPeer{}
		peerKey := types.NamespacedName{Namespace: FabricNamespace, Name: releaseNamePeer}
		Eventually(
			func() bool {
				err := K8sClient.Get(context.Background(), peerKey, peer)
				if err != nil {
					return false
				}
				ctrl.Log.WithName("test").Info("after update", "peer", peer)
				return peer.Status.Status == hlfv1alpha1.RunningStatus
			},
			peerTimeoutSecs,
			defInterval,
		).Should(BeTrue(), "peer status should have been updated")
		Expect(peer.Status.URL).ToNot(BeEmpty())
		Expect(peer.Status.TlsCert).ToNot(BeEmpty())

		By("test the peer API")
		Eventually(
			func() bool {
				resClient := getClientForPeer(peer, updatedCA)
				channelResponse, err := resClient.QueryChannels(
					resmgmt.WithTargetEndpoints("peer"),
				)
				if err != nil {
					return false
				}
				return len(channelResponse.Channels) == 0
			},
			peerTimeoutSecs,
			defInterval,
		).Should(BeTrue(), "peer channels should be 0")
		resClient := getClientForPeer(peer, updatedCA)
		By("Installing a new chaincode")
		pkgLabel := "fabcar"
		packageBytes, err := lifecycle.NewCCPackage(&lifecycle.Descriptor{
			Path:  "../../fixtures/chaincodes/fabcar/javascript",
			Type:  pb.ChaincodeSpec_NODE,
			Label: pkgLabel,
		})
		Expect(err).ToNot(HaveOccurred())
		responses, err := resClient.LifecycleInstallCC(
			resmgmt.LifecycleInstallCCRequest{
				Label:   pkgLabel,
				Package: packageBytes,
			},
			resmgmt.WithTimeout(fab.ResMgmt, 20*time.Minute),
			resmgmt.WithTimeout(fab.PeerResponse, 20*time.Minute),
		)
		Expect(err).ToNot(HaveOccurred())
		Expect(responses).To(HaveLen(1))
		installedRes, err := resClient.LifecycleQueryInstalledCC(
			resmgmt.WithTargetEndpoints("peer"),
		)
		Expect(err).ToNot(HaveOccurred())
		Expect(installedRes).To(HaveLen(1))
	})
	Specify("create a new Fabric Orderer instance", func() {
		releaseNameOrdCA := "org1-ca"
		releaseNameOrd := "org1-orderer"
		By("create a fabric ca")
		ordererCA := randomFabricCA(releaseNameOrdCA, FabricNamespace)
		Expect(ordererCA).ToNot(BeNil())
		By("create a fabric orderer")
		ordererMSPID := "OrdererMSP"
		ordParams := createOrdererParams{
			MSPID: ordererMSPID,
		}
		createOrderer(
			releaseNameOrd,
			FabricNamespace,
			ordParams,
			ordererCA,
		)
		orderer := &hlfv1alpha1.FabricOrderingService{}
		ordererKey := types.NamespacedName{
			Namespace: FabricNamespace,
			Name:      releaseNameOrd,
		}
		Eventually(
			func() bool {
				err := K8sClient.Get(context.Background(), ordererKey, orderer)
				if err != nil {
					return false
				}
				ctrl.Log.WithName("test").Info("after update", "orderer", orderer)
				return orderer.Status.Status == hlfv1alpha1.RunningStatus
			},
			peerTimeoutSecs,
			defInterval,
		).Should(BeTrue(), "peer status should have been updated")

		By("create a fabric peer")
		releaseNamePeer := "org1-peer0"
		releaseNamePeerCA := "org1-peer0-ca"
		peerCA := randomFabricCA(releaseNamePeerCA, FabricNamespace)
		Expect(peerCA).ToNot(BeNil())
		peerMSPID := "Org1MSP"
		peerParams := createPeerParams{
			MSPID: peerMSPID,
		}
		createPeer(
			releaseNamePeer,
			FabricNamespace,
			peerParams,
			peerCA,
		)
		peer := &hlfv1alpha1.FabricPeer{}
		peerKey := types.NamespacedName{
			Namespace: FabricNamespace,
			Name:      releaseNamePeer,
		}
		Eventually(
			func() bool {
				err := K8sClient.Get(context.Background(), peerKey, peer)
				if err != nil {
					return false
				}
				ctrl.Log.WithName("test").Info("after update", "peer", peer)
				return peer.Status.Status == hlfv1alpha1.RunningStatus
			},
			peerTimeoutSecs,
			defInterval,
		).Should(BeTrue(), "peer status should have been updated")
		Expect(peer.Status.URL).ToNot(BeEmpty())
		Expect(peer.Status.TlsCert).ToNot(BeEmpty())

		By("add the peer to the consortium")
		ordClient := getClientForOrderer(orderer, ordererCA)
		block, err := ordClient.QueryConfigBlockFromOrderer(systemChannelID)
		Expect(err).ToNot(HaveOccurred())
		systemChannelConfig, err := resource.ExtractConfigFromBlock(block)
		Expect(err).ToNot(HaveOccurred())
		log.Print(systemChannelConfig)
		u, err := url.Parse(peer.Status.URL)
		Expect(err).ToNot(HaveOccurred())
		host, portStr, err := net.SplitHostPort(u.Host)
		Expect(err).ToNot(HaveOccurred())
		port, err := strconv.Atoi(portStr)
		Expect(err).ToNot(HaveOccurred())
		consortiumName := "SampleConsortium"

		modifiedConfig, err := testutils.AddConsortiumToConfig(
			systemChannelConfig,
			testutils.AddConsortiumRequest{
				Name: consortiumName,
				Organizations: []testutils.PeerOrganization{
					{
						RootCert:    peerCA.Status.CACert,
						TLSRootCert: peerCA.Status.CACert,
						MspID:       peer.Spec.MspID,
						Peers: []testutils.PeerNode{
							{
								Host: host,
								Port: port,
							},
						},
					},
				},
			},
		)
		Expect(err).ToNot(HaveOccurred())
		confUpdate, err := resmgmt.CalculateConfigUpdate(
			systemChannelID,
			systemChannelConfig,
			modifiedConfig,
		)
		Expect(err).ToNot(HaveOccurred())
		configEnvelopeBytes, err := testutils.GetConfigEnvelopeBytes(confUpdate)
		Expect(err).ToNot(HaveOccurred())
		configReader := bytes.NewReader(configEnvelopeBytes)
		saveResponse, err := ordClient.SaveChannel(resmgmt.SaveChannelRequest{
			ChannelID:     systemChannelID,
			ChannelConfig: configReader,
		})
		Expect(err).ToNot(HaveOccurred())
		Expect(saveResponse).ToNot(BeNil())

		By("create a channel")

		channelID := getRandomChannelID()
		ordNodes := getOrderers(
			releaseNameOrd,
			FabricNamespace,
		)
		nodes := []testutils.OrdererNode{}
		for _, item := range ordNodes {
			nodes = append(nodes, testutils.OrdererNode{
				Port:    item.Status.Port,
				Host:    item.Status.Host,
				TLSCert: item.Spec.TLSCert,
			})
		}
		orgOrganization := testutils.OrdererOrganization{
			Nodes:        nodes,
			RootTLSCert:  ordererCA.Status.TLSCACert,
			MspID:        orderer.Spec.MspID,
			RootSignCert: ordererCA.Status.CACert,
		}
		u, err = url.Parse(peer.Status.URL)
		Expect(err).ToNot(HaveOccurred())
		host, portStr, err = net.SplitHostPort(u.Host)
		Expect(err).ToNot(HaveOccurred())
		port, err = strconv.Atoi(portStr)
		Expect(err).ToNot(HaveOccurred())
		peerOrgs := []testutils.PeerOrganization{
			{
				RootCert:    peerCA.Status.CACert,
				TLSRootCert: peerCA.Status.CACert,
				MspID:       peer.Spec.MspID,
				Peers: []testutils.PeerNode{
					{
						Host: host,
						Port: port,
					},
				},
			},
		}
		profileConfig, err := testutils.GetChannelProfileConfig(
			orgOrganization,
			peerOrgs,
			consortiumName,
			fmt.Sprintf("OR('%s.admin')", peerOrgs[0].MspID),
		)
		Expect(err).ToNot(HaveOccurred())
		var baseProfile *genesisconfig.Profile
		channelTx, err := resource.CreateChannelCreateTx(
			profileConfig,
			baseProfile,
			channelID,
		)
		Expect(err).ToNot(HaveOccurred())
		channelConfig := bytes.NewReader(channelTx)
		resClient := getClientForPeerWithOrderer(
			peer,
			peerCA,
			orderer,
			ordererCA,
		)
		createChannelResponse, err := resClient.SaveChannel(resmgmt.SaveChannelRequest{
			ChannelID:     channelID,
			ChannelConfig: channelConfig,
		})
		Expect(err).ToNot(HaveOccurred())
		Expect(createChannelResponse).ToNot(BeNil())

		By("join the peer to the channel")
		time.Sleep(2 * time.Second) // wait for the transaction to be committed
		err = resClient.JoinChannel(
			channelID,
			resmgmt.WithTargetEndpoints("peer"),
			resmgmt.WithOrdererEndpoint("orderer"),
		)
		Expect(err).ToNot(HaveOccurred())
		By("install a chaincode for the org")
		pkgLabel := "fabcar"
		packageBytes, err := lifecycle.NewCCPackage(&lifecycle.Descriptor{
			Path:  "../../fixtures/chaincodes/fabcar/javascript",
			Type:  pb.ChaincodeSpec_NODE,
			Label: pkgLabel,
		})
		Expect(err).ToNot(HaveOccurred())
		_, err = resClient.LifecycleInstallCC(
			resmgmt.LifecycleInstallCCRequest{
				Label:   pkgLabel,
				Package: packageBytes,
			},
			resmgmt.WithTimeout(fab.ResMgmt, 20*time.Minute),
			resmgmt.WithTimeout(fab.PeerResponse, 20*time.Minute),
		)
		By("approve a chaincode in the peer")
		ccName := "fabcar"
		version := "1.0"
		sequence := 1
		Expect(err).ToNot(HaveOccurred())
		packageID := lifecycle.ComputePackageID(pkgLabel, packageBytes)
		sp, err := policydsl.FromString("OR('Org1MSP.peer')")
		Expect(err).ToNot(HaveOccurred())
		m, err := resClient.LifecycleQueryCommittedCC(channelID, resmgmt.LifecycleQueryCommittedCCRequest{})
		Expect(err).ToNot(HaveOccurred())
		Expect(m).To(HaveLen(0))
		_, err = resClient.LifecycleApproveCC(
			channelID,
			resmgmt.LifecycleApproveCCRequest{
				Name:              ccName,
				Version:           version,
				PackageID:         packageID,
				Sequence:          int64(sequence),
				EndorsementPlugin: "escc",
				ValidationPlugin:  "vscc",
				SignaturePolicy:   sp,
				CollectionConfig:  nil,
				InitRequired:      false,
			},
		)
		Expect(err).ToNot(HaveOccurred())

		_, err = resClient.LifecycleCommitCC(
			channelID,
			resmgmt.LifecycleCommitCCRequest{
				Name:              ccName,
				Version:           version,
				Sequence:          int64(sequence),
				EndorsementPlugin: "escc",
				ValidationPlugin:  "vscc",
				SignaturePolicy:   sp,
				CollectionConfig:  nil,
				InitRequired:      false,
			},
		)
		Expect(err).ToNot(HaveOccurred())
		time.Sleep(1 * time.Second)
		sdk := getSDKForPeerWithOrderer(peer, peerCA, orderer, ordererCA)
		channelCtx := sdk.ChannelContext(channelID, fabsdk.WithUser("admin"), fabsdk.WithOrg(peer.Spec.MspID))
		channelClient, err := channel.New(channelCtx)
		Expect(err).ToNot(HaveOccurred())
		exRes, err := channelClient.Execute(
			channel.Request{
				ChaincodeID:     ccName,
				Fcn:             "initLedger",
				Args:            [][]byte{},
				TransientMap:    nil,
				InvocationChain: nil,
				IsInit:          false,
			},
		)
		Expect(err).ToNot(HaveOccurred())
		Expect(exRes.TransactionID).ToNot(BeEmpty())
		time.Sleep(1 * time.Second)
		qr, err := channelClient.Query(
			channel.Request{
				ChaincodeID:     ccName,
				Fcn:             "queryAllCars",
				Args:            [][]byte{},
				TransientMap:    nil,
				InvocationChain: nil,
				IsInit:          false,
			},
		)
		Expect(err).ToNot(HaveOccurred())
		Expect(qr.TransactionID).ToNot(BeEmpty())

	})
})

func NewTypeMeta(kind string) metav1.TypeMeta {
	return metav1.TypeMeta{
		Kind:       kind,
		APIVersion: "hlf.kungfusoftware.es/v1alpha1",
	}
}

func getRandomChannelID() string {
	alphabet := "23456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
	channelID := shortuuid.NewWithAlphabet(alphabet)

	for i := 1; i <= 9; i++ {
		channelID = strings.Replace(channelID, strconv.Itoa(i), "a", -1)
	}
	return strings.ToLower(channelID)
}
