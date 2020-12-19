package ca

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"github.com/operator-framework/operator-lib/status"
	"k8s.io/kubernetes/pkg/api/v1/pod"
	"log"
	"math/big"
	"net"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/go-logr/logr"
	hlfv1alpha1 "github.com/kfsoftware/hlf-operator/api/hlf.kungfusoftware.es/v1alpha1"
	"github.com/kfsoftware/hlf-operator/controllers/utils"
	"github.com/pkg/errors"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/storage/driver"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// FabricCAReconciler reconciles a FabricCA object
type FabricCAReconciler struct {
	client.Client
	ChartPath string
	Log       logr.Logger
	Scheme    *runtime.Scheme
	Config    *rest.Config
	ClientSet *kubernetes.Clientset
}

func parseECDSAPrivateKey(contents []byte) (*ecdsa.PrivateKey, error) {
	block, _ := pem.Decode(contents)
	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	ecdsaKey, ok := key.(*ecdsa.PrivateKey)
	if !ok {
		return nil, errors.New("private key is not of ECDSA type")
	}
	return ecdsaKey, nil
}
func parseX509Certificate(contents []byte) (*x509.Certificate, error) {
	block, _ := pem.Decode(contents)
	crt, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, err
	}
	return crt, nil
}

func getExistingTLSCrypto(client *kubernetes.Clientset, chartName string, namespace string) (*x509.Certificate, *ecdsa.PrivateKey, error) {
	secretName := fmt.Sprintf("%s--tls-cryptomaterial", chartName)

	secret, err := client.CoreV1().Secrets(namespace).Get(context.Background(), secretName, v1.GetOptions{})
	if err != nil {
		return nil, nil, err
	}
	tlsKeyData := secret.Data["tls.key"]
	tlsCrtData := secret.Data["tls.crt"]
	key, err := parseECDSAPrivateKey(tlsKeyData)
	if err != nil {
		return nil, nil, err
	}
	crt, err := parseX509Certificate(tlsCrtData)
	if err != nil {
		return nil, nil, err
	}
	return crt, key, nil
}

func getExistingSignCrypto(client *kubernetes.Clientset, chartName string, namespace string) (*x509.Certificate, *ecdsa.PrivateKey, error) {
	secretName := fmt.Sprintf("%s--msp-cryptomaterial", chartName)

	secret, err := client.CoreV1().Secrets(namespace).Get(context.Background(), secretName, v1.GetOptions{})
	if err != nil {
		return nil, nil, err
	}
	tlsKeyData := secret.Data["keyfile"]
	tlsCrtData := secret.Data["certfile"]
	key, err := parseECDSAPrivateKey(tlsKeyData)
	if err != nil {
		return nil, nil, err
	}
	crt, err := parseX509Certificate(tlsCrtData)
	if err != nil {
		return nil, nil, err
	}
	return crt, key, nil
}

func getExistingSignTLSCrypto(client *kubernetes.Clientset, chartName string, namespace string) (*x509.Certificate, *ecdsa.PrivateKey, error) {
	secretName := fmt.Sprintf("%s--msp-tls-cryptomaterial", chartName)

	secret, err := client.CoreV1().Secrets(namespace).Get(context.Background(), secretName, v1.GetOptions{})
	if err != nil {
		return nil, nil, err
	}
	tlsKeyData := secret.Data["keyfile"]
	tlsCrtData := secret.Data["certfile"]
	key, err := parseECDSAPrivateKey(tlsKeyData)
	if err != nil {
		return nil, nil, err
	}
	crt, err := parseX509Certificate(tlsCrtData)
	if err != nil {
		return nil, nil, err
	}
	return crt, key, nil
}

// compute Subject Key Identifier
func computeSKI(privKey *ecdsa.PrivateKey) []byte {
	// Marshall the public key
	raw := elliptic.Marshal(privKey.Curve, privKey.PublicKey.X, privKey.PublicKey.Y)

	// Hash it
	hash := sha256.Sum256(raw)
	return hash[:]
}
func CreateDefaultTLSCA(clientSet *kubernetes.Clientset, spec hlfv1alpha1.FabricCASpec) (*x509.Certificate, *ecdsa.PrivateKey, error) {
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		log.Fatalf("Failed to generate serial number: %v", err)
		return nil, nil, err
	}
	k8sIP, err := utils.GetPublicIPKubernetes(clientSet)
	if err != nil {
		return nil, nil, err
	}
	var dnsNames []string
	ips := []net.IP{net.ParseIP("127.0.0.1")}
	for _, host := range spec.Hosts {
		addr := net.ParseIP(host)
		if addr == nil {
			dnsNames = append(dnsNames, host)
		} else {
			ips = append(ips, addr)
		}
	}
	if !utils.Contains(spec.Hosts, k8sIP) {
		addr := net.ParseIP(k8sIP)
		if addr != nil {
			ips = append(ips, addr)
		}
	}
	caPrivKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, nil, err
	}

	x509Cert := &x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization:       []string{spec.TLS.Subject.O},
			Country:            []string{spec.TLS.Subject.C},
			Locality:           []string{spec.TLS.Subject.L},
			OrganizationalUnit: []string{spec.TLS.Subject.OU},
			StreetAddress:      []string{spec.TLS.Subject.ST},
		},
		NotBefore:             time.Now().AddDate(0, 0, -1),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		IsCA:                  true,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		BasicConstraintsValid: true,
		DNSNames:              dnsNames,
		IPAddresses:           ips,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		SubjectKeyId:          computeSKI(caPrivKey),
	}

	caBytes, err := x509.CreateCertificate(rand.Reader, x509Cert, x509Cert, &caPrivKey.PublicKey, caPrivKey)
	if err != nil {
		return nil, nil, err
	}
	crt, err := x509.ParseCertificate(caBytes)
	if err != nil {
		return nil, nil, err
	}
	return crt, caPrivKey, nil
}

func CreateDefaultCA(conf hlfv1alpha1.FabricCAItemConf) (*x509.Certificate, *ecdsa.PrivateKey, error) {
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		log.Fatalf("Failed to generate serial number: %v", err)
		return nil, nil, err
	}
	caPrivKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, nil, err
	}

	signCA := &x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization:       []string{conf.Subject.O},
			Country:            []string{conf.Subject.C},
			Locality:           []string{conf.Subject.L},
			OrganizationalUnit: []string{conf.Subject.OU},
			StreetAddress:      []string{conf.Subject.ST},
			CommonName:         conf.Subject.CN,
		},
		NotBefore:             time.Now().AddDate(0, 0, -1),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		IsCA:                  true,
		SubjectKeyId:          computeSKI(caPrivKey),
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign | x509.KeyUsageCRLSign | x509.KeyUsageKeyEncipherment,
		BasicConstraintsValid: true,
	}
	caBytes, err := x509.CreateCertificate(rand.Reader, signCA, signCA, &caPrivKey.PublicKey, caPrivKey)
	if err != nil {
		return nil, nil, err
	}
	crt, err := x509.ParseCertificate(caBytes)
	if err != nil {
		return nil, nil, err
	}
	return crt, caPrivKey, nil
}

func newActionCfg(log logr.Logger, clusterCfg *rest.Config, namespace string) (*action.Configuration, error) {
	err := os.Setenv("HELM_NAMESPACE", namespace)
	if err != nil {
		return nil, err
	}
	cfg := new(action.Configuration)
	ns := namespace
	err = cfg.Init(&genericclioptions.ConfigFlags{
		Namespace:   &ns,
		APIServer:   &clusterCfg.Host,
		CAFile:      &clusterCfg.CAFile,
		BearerToken: &clusterCfg.BearerToken,
	}, ns, "secret", actionLogger(log))
	return cfg, err
}

func actionLogger(logger logr.Logger) func(format string, v ...interface{}) {
	return func(format string, v ...interface{}) {
		logger.Info(fmt.Sprintf(format, v...))
	}
}
func mapCRDItemConfToChart(conf hlfv1alpha1.FabricCAItemConf) FabricCAChartItemConf {
	names := []FabricCAChartNames{}
	for _, name := range conf.CSR.Names {
		names = append(names, FabricCAChartNames{
			C:  name.C,
			ST: name.ST,
			O:  name.O,
			L:  name.L,
			OU: name.OU,
		})
	}
	identities := []FabricCAChartIdentity{}
	for _, identity := range conf.Registry.Identities {
		identities = append(identities, FabricCAChartIdentity{
			Name:        identity.Name,
			Pass:        identity.Pass,
			Type:        identity.Type,
			Affiliation: identity.Affiliation,
			Attrs: FabricCAChartIdentityAttrs{
				RegistrarRoles: identity.Attrs.RegistrarRoles,
				DelegateRoles:  identity.Attrs.DelegateRoles,
				Attributes:     identity.Attrs.Attributes,
				Revoker:        identity.Attrs.Revoker,
				IntermediateCA: identity.Attrs.IntermediateCA,
				GenCRL:         identity.Attrs.GenCRL,
				AffiliationMgr: identity.Attrs.AffiliationMgr,
			},
		})
	}
	item := FabricCAChartItemConf{
		Name: conf.Name,
		CFG: FabricCAChartCFG{
			Identities:   FabricCAChartCFGIdentities{AllowRemove: conf.CFG.Identities.AllowRemove},
			Affiliations: FabricCAChartCFGAffilitions{AllowRemove: conf.CFG.Affiliations.AllowRemove},
		},
		CSR: FabricCAChartCSR{
			CN:    conf.CSR.CN,
			Hosts: conf.CSR.Hosts,
			Names: names,
			CA: FabricCAChartCSRCA{
				Expiry:     conf.CSR.CA.Expiry,
				PathLength: conf.CSR.CA.PathLength,
			},
		},
		CRL: FabricCAChartCRL{Expiry: conf.CRL.Expiry},
		Registry: FabricCAChartRegistry{
			MaxEnrollments: conf.Registry.MaxEnrollments,
			Identities:     identities,
		},
		Intermediate: FabricCAChartIntermediate{
			ParentServer: FabricCAChartIntermediateParentServer{
				URL:    conf.Intermediate.ParentServer.URL,
				CAName: conf.Intermediate.ParentServer.CAName,
			},
		},
		Affiliations: []Affiliation{},
		BCCSP: FabricCAChartBCCSP{
			Default: conf.BCCSP.Default,
			SW: FabricCAChartBCCSPSW{
				Hash:     conf.BCCSP.SW.Hash,
				Security: conf.BCCSP.SW.Security,
			},
		},
	}
	return item
}
func GetConfig(conf *hlfv1alpha1.FabricCA, client *kubernetes.Clientset, chartName string, namespace string) (*FabricCAChart, error) {
	spec := conf.Spec

	tlsCert, tlsKey, err := getExistingTLSCrypto(client, chartName, namespace)
	if err != nil {
		tlsCert, tlsKey, err = CreateDefaultTLSCA(client, spec)
		if err != nil {
			return nil, err
		}
	}
	signCert, signKey, err := getExistingSignCrypto(client, chartName, namespace)
	if err != nil {
		signCert, signKey, err = CreateDefaultCA(spec.CA)
		if err != nil {
			return nil, err
		}
	}
	caTLSSignCert, caTLSSignKey, err := getExistingSignTLSCrypto(client, chartName, namespace)
	if err != nil {
		caTLSSignCert, caTLSSignKey, err = CreateDefaultCA(spec.TLSCA)
		if err != nil {
			return nil, err
		}
	}
	tlsCRTEncoded := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: tlsCert.Raw,
	})
	tlsEncodedPK, err := x509.MarshalPKCS8PrivateKey(tlsKey)
	if err != nil {
		return nil, err
	}
	tlsPEMEncodedPK := pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: tlsEncodedPK,
	})

	signCRTEncoded := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: signCert.Raw,
	})
	signEncodedPK, err := x509.MarshalPKCS8PrivateKey(signKey)
	if err != nil {
		return nil, err
	}
	signPEMEncodedPK := pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: signEncodedPK,
	})

	caTLSSignCRTEncoded := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: caTLSSignCert.Raw,
	})
	caTLSSignEncodedPK, err := x509.MarshalPKCS8PrivateKey(caTLSSignKey)
	if err != nil {
		return nil, err
	}
	caTLSSignPEMEncodedPK := pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: caTLSSignEncodedPK,
	})
	var c = FabricCAChart{
		FullNameOverride: conf.Name,
		Image: Image{
			Repository: spec.Image,
			Tag:        spec.Version,
			PullPolicy: "IfNotPresent",
		},
		Service: Service{
			Type: spec.Service.ServiceType,
			Port: 7054,
		},
		Persistence: Persistence{
			Enabled:      true,
			Annotations:  map[string]string{},
			StorageClass: spec.Storage.StorageClass,
			AccessMode:   string(spec.Storage.AccessMode),
			Size:         spec.Storage.Size,
		},
		Msp: Msp{
			Keyfile:       string(signPEMEncodedPK),
			Certfile:      string(signCRTEncoded),
			TLSCAKeyfile:  string(caTLSSignPEMEncodedPK),
			TLSCACertfile: string(caTLSSignCRTEncoded),
			TlsKeyFile:    string(tlsPEMEncodedPK),
			TlsCertFile:   string(tlsCRTEncoded),
		},
		Database: Database{
			Type:       spec.Database.Type,
			Datasource: spec.Database.Datasource,
		},
		Resources: Resources{
			Requests: Requests{
				CPU:    spec.Resources.Requests.CPU,
				Memory: spec.Resources.Requests.Memory,
			},
			Limits: RequestsLimit{
				CPU:    spec.Resources.Limits.CPU,
				Memory: spec.Resources.Limits.Memory,
			},
		},
		NodeSelector: NodeSelector{},
		Tolerations:  nil,
		Affinity:     Affinity{},
		Debug:        spec.Debug,
		CLRSizeLimit: spec.CLRSizeLimit,
		Metrics: FabricCAChartMetrics{
			Provider: spec.Metrics.Provider,
			Statsd: FabricCAChartMetricsStatsd{
				Network:       spec.Metrics.Statsd.Network,
				Address:       spec.Metrics.Statsd.Address,
				WriteInterval: spec.Metrics.Statsd.WriteInterval,
				Prefix:        spec.Metrics.Statsd.Prefix,
			},
		},

		Ca:    mapCRDItemConfToChart(spec.CA),
		TLSCA: mapCRDItemConfToChart(spec.TLSCA),
		Cors: Cors{
			Enabled: spec.Cors.Enabled,
			Origins: spec.Cors.Origins,
		},
	}
	return &c, nil
}

type Status struct {
	URL       string
	Port      int
	Host      string
	Status    hlfv1alpha1.DeploymentStatus
	TlsCert   string
	CACert    string
	TLSCACert string
}

func GetServiceName(releaseName string) string {
	return releaseName
}
func GetDeploymentName(releaseName string) string {
	return releaseName
}
func GetCAState(clientSet *kubernetes.Clientset, releaseName string, ns string) (*Status, error) {
	ctx := context.Background()
	k8sIP, err := utils.GetPublicIPKubernetes(clientSet)
	if err != nil {
		return nil, err
	}
	r := &Status{
		Status: hlfv1alpha1.PendingStatus,
	}
	depName := GetDeploymentName(releaseName)
	dep, err := clientSet.AppsV1().Deployments(ns).Get(ctx, depName, v1.GetOptions{})
	if err != nil {
		return nil, err
	}
	pods, err := clientSet.CoreV1().Pods(ns).List(ctx, v1.ListOptions{
		LabelSelector: fmt.Sprintf("release=%s", releaseName),
	})
	if err != nil {
		return nil, err
	}
	if len(pods.Items) > 0 {
		for _, item := range pods.Items {
			if pod.IsPodReadyConditionTrue(item.Status) {
				r.Status = hlfv1alpha1.RunningStatus
			} else {
				switch item.Status.Phase {
				case corev1.PodPending:
					r.Status = hlfv1alpha1.PendingStatus
				case corev1.PodFailed:
					r.Status = hlfv1alpha1.FailedStatus
				case corev1.PodSucceeded:
				case corev1.PodRunning:
				case corev1.PodUnknown:
					r.Status = hlfv1alpha1.UnknownStatus
				}
			}
		}
	} else {
		if dep.Status.ReadyReplicas == *dep.Spec.Replicas {
			r.Status = hlfv1alpha1.RunningStatus
		} else {
			r.Status = hlfv1alpha1.PendingStatus
		}
	}
	svcName := GetServiceName(releaseName)
	svc, err := clientSet.CoreV1().Services(ns).Get(ctx, svcName, v1.GetOptions{})
	if err != nil {
		return nil, err
	}
	nodePort := svc.Spec.Ports[0].NodePort
	r.URL = fmt.Sprintf("https://%s:%d", k8sIP, nodePort)
	r.Port = int(nodePort)
	r.Host = k8sIP
	tlsCrt, _, err := getExistingTLSCrypto(clientSet, releaseName, ns)
	if err != nil {
		return nil, err
	}
	r.TlsCert = string(utils.EncodeX509Certificate(tlsCrt))
	signCrt, _, err := getExistingSignCrypto(clientSet, releaseName, ns)
	if err != nil {
		return nil, err
	}
	r.CACert = string(utils.EncodeX509Certificate(signCrt))
	tlsCACrt, _, err := getExistingSignTLSCrypto(clientSet, releaseName, ns)
	if err != nil {
		return nil, err
	}
	r.TLSCACert = string(utils.EncodeX509Certificate(tlsCACrt))
	return r, nil
}

const caFinalizer = "finalizer.ca.hlf.kungfusoftware.es"

func (r *FabricCAReconciler) finalizeCA(reqLogger logr.Logger, m *hlfv1alpha1.FabricCA) error {
	ns := m.Namespace
	if ns == "" {
		ns = "default"
	}
	cfg, err := newActionCfg(r.Log, r.Config, ns)
	if err != nil {
		return err
	}
	releaseName := m.Name
	reqLogger.Info("Successfully finalized ca")
	cmd := action.NewUninstall(cfg)
	resp, err := cmd.Run(releaseName)
	if err != nil {
		if strings.Compare("Release not loaded", err.Error()) != 0 {
			return nil
		}
		log.Printf("Failed to uninstall release %s %v", releaseName, err)
		return err
	}
	log.Printf("Release %s deleted=%s", releaseName, resp.Info)
	return nil
}

func (r *FabricCAReconciler) addFinalizer(reqLogger logr.Logger, m *hlfv1alpha1.FabricCA) error {
	reqLogger.Info("Adding Finalizer for the CA")
	controllerutil.AddFinalizer(m, caFinalizer)

	// Update CR
	err := r.Update(context.TODO(), m)
	if err != nil {
		reqLogger.Error(err, "Failed to update CA with finalizer")
		return err
	}
	return nil
}

func Reconcile(
	req ctrl.Request,
	r *FabricCAReconciler,
	cfg *action.Configuration,
) (ctrl.Result, error) {
	ctx := context.Background()
	reqLogger := r.Log.WithValues("hlf", req.NamespacedName)
	ns := req.Namespace
	hlf := &hlfv1alpha1.FabricCA{}
	releaseName := req.Name
	err := r.Get(ctx, req.NamespacedName, hlf)
	if err != nil {
		log.Printf("Error getting the object %s error=%v", req.NamespacedName, err)
		if apierrors.IsNotFound(err) {
			reqLogger.Info("CA resource not found. Ignoring since object must be deleted.")
			return ctrl.Result{}, nil
		}
		reqLogger.Error(err, "Failed to get CA.")
		return ctrl.Result{}, err
	}
	isMemcachedMarkedToBeDeleted := hlf.GetDeletionTimestamp() != nil
	if isMemcachedMarkedToBeDeleted {
		if utils.Contains(hlf.GetFinalizers(), caFinalizer) {
			// Run finalization logic for caFinalizer. If the
			// finalization logic fails, don't remove the finalizer so
			// that we can retry during the next reconciliation.
			if err := r.finalizeCA(reqLogger, hlf); err != nil {
				return ctrl.Result{}, err
			}

			// Remove caFinalizer. Once all finalizers have been
			// removed, the object will be deleted.
			controllerutil.RemoveFinalizer(hlf, caFinalizer)
			err := r.Update(ctx, hlf)
			if err != nil {
				return ctrl.Result{}, err
			}
		}
		return ctrl.Result{}, nil
	}
	if !utils.Contains(hlf.GetFinalizers(), caFinalizer) {
		if err := r.addFinalizer(reqLogger, hlf); err != nil {
			return ctrl.Result{}, err
		}
	}

	cmdStatus := action.NewStatus(cfg)
	exists := true
	_, err = cmdStatus.Run(releaseName)
	if err != nil {
		if errors.Is(err, driver.ErrReleaseNotFound) {
			// it doesn't exists
			exists = false
		} else {
			// it doesnt exist
			return ctrl.Result{}, err
		}
	}
	log.Printf("Release %s exists=%v", releaseName, exists)
	clientSet, err := utils.GetClientKubeWithConf(r.Config)
	if err != nil {
		return ctrl.Result{}, err
	}

	if exists {
		// update
		s, err := GetCAState(r.ClientSet, releaseName, ns)
		if err != nil {
			return ctrl.Result{}, err
		}
		fca := hlf.DeepCopy()
		fca.Status.Status = s.Status
		fca.Status.URL = s.URL
		fca.Status.TlsCert = s.TlsCert
		fca.Status.TLSCACert = s.TLSCACert
		fca.Status.CACert = s.CACert
		fca.Status.Host = s.Host
		fca.Status.Port = s.Port
		fca.Status.Conditions.SetCondition(status.Condition{
			Type:               status.ConditionType(s.Status),
			Status:             "True",
			LastTransitionTime: v1.Time{},
		})
		if reflect.DeepEqual(fca.Status, hlf.Status) {
			log.Printf("Status hasn't changed, skipping CA update")
		} else {
			// TODO: DON'T UPGRADE IF NOT NEEDED TO
			//cmdGet := action.NewGet(cfg)
			//rel, err := cmdGet.Run(releaseName)
			//if err != nil {
			//	return ctrl.Result{}, err
			//}
			//ns := rel.Namespace
			//if ns == "" {
			//	ns = "default"
			//}
			//c, err := GetConfig(hlf, clientSet, releaseName, ns)
			//if err != nil {
			//	return ctrl.Result{}, err
			//}
			//inrec, err := json.Marshal(c)
			//if err != nil {
			//	return ctrl.Result{}, err
			//}
			//var inInterface map[string]interface{}
			//err = json.Unmarshal(inrec, &inInterface)
			//if err != nil {
			//	return ctrl.Result{}, err
			//}
			//cmd := action.NewUpgrade(cfg)
			//settings := cli.New()
			//chartPath, err := cmd.LocateChart(chartPath, settings)
			//ch, err := loader.Load(chartPath)
			//if err != nil {
			//	return ctrl.Result{}, err
			//}
			//release, err := cmd.Run(releaseName, ch, inInterface)
			//if err != nil {
			//	return ctrl.Result{}, err
			//}
			//log.Printf("Chart upgraded %s", release.Name)

			if err := r.Status().Update(ctx, fca); err != nil {
				log.Printf("Error updating the status: %v", err)
				return ctrl.Result{}, err
			}
			if err := r.Status().Update(ctx, hlf); err != nil {
				return ctrl.Result{}, err
			}
		}

		if s.Status == hlfv1alpha1.RunningStatus {
			return ctrl.Result{
				Requeue:      false,
				RequeueAfter: 0,
			}, nil
		} else {
			return ctrl.Result{
				Requeue:      false,
				RequeueAfter: 5 * time.Second,
			}, nil
		}
	} else {
		cmd := action.NewInstall(cfg)
		name, chart, err := cmd.NameAndChart([]string{releaseName, r.ChartPath})
		if err != nil {
			return ctrl.Result{}, err
		}

		cmd.ReleaseName = name
		ch, err := loader.Load(chart)
		if err != nil {
			return ctrl.Result{}, err
		}
		c, err := GetConfig(hlf, clientSet, name, req.Namespace)
		if err != nil {
			reqLogger.Error(err, "Failed to get config")
			return ctrl.Result{}, err
		}
		var inInterface map[string]interface{}
		inrec, err := json.Marshal(c)
		if err != nil {
			reqLogger.Error(err, "Failed to marshall helm values")
			return ctrl.Result{}, err
		}
		log.Println(string(inrec))
		err = json.Unmarshal(inrec, &inInterface)
		if err != nil {
			return ctrl.Result{}, err
		}
		release, err := cmd.Run(ch, inInterface)
		if err != nil {
			reqLogger.Error(err, "Failed to install helm chart")
			return ctrl.Result{}, err
		}
		log.Printf("Chart installed %s", release.Name)
		hlf.Status.Status = hlfv1alpha1.PendingStatus
		hlf.Status.Conditions.SetCondition(status.Condition{
			Type:               "DEPLOYED",
			Status:             "True",
			LastTransitionTime: v1.Time{},
		})
		if err := r.Status().Update(ctx, hlf); err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{
			Requeue:      false,
			RequeueAfter: 10 * time.Second,
		}, nil
	}
}

// +kubebuilder:rbac:groups=hlf.kungfusoftware.es,resources=fabriccas,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=hlf.kungfusoftware.es,resources=fabriccas/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=hlf.kungfusoftware.es,resources=fabriccas/finalizers,verbs=update
func (r *FabricCAReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {

	ns := req.Namespace
	cfg, err := newActionCfg(r.Log, r.Config, ns)
	if err != nil {
		return ctrl.Result{}, err
	}
	return Reconcile(
		req,
		r,
		cfg,
	)

}

func (r *FabricCAReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&hlfv1alpha1.FabricCA{}).
		Owns(&appsv1.Deployment{}).
		Complete(r)
}
