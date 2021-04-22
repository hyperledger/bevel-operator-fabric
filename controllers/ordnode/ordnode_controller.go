package ordnode

import (
	"context"
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/go-logr/logr"
	hlfv1alpha1 "github.com/kfsoftware/hlf-operator/api/hlf.kungfusoftware.es/v1alpha1"
	"github.com/kfsoftware/hlf-operator/controllers/certs"
	"github.com/kfsoftware/hlf-operator/controllers/utils"
	"github.com/operator-framework/operator-lib/status"
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
	"k8s.io/kubernetes/pkg/api/v1/pod"
	"log"
	"os"
	"reflect"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"strings"
	"time"
)

// FabricOrdererNodeReconciler reconciles a FabricOrdererNode object
type FabricOrdererNodeReconciler struct {
	client.Client
	ChartPath string
	Log       logr.Logger
	Scheme    *runtime.Scheme
	Config    *rest.Config
}

const ordererNodeFinalizer = "finalizer.orderernode.hlf.kungfusoftware.es"

func (r *FabricOrdererNodeReconciler) finalizeOrderer(reqLogger logr.Logger, m *hlfv1alpha1.FabricOrdererNode) error {
	ns := m.Namespace
	if ns == "" {
		ns = "default"
	}
	cfg, err := newActionCfg(r.Log, r.Config, ns)
	if err != nil {
		return err
	}
	releaseName := m.Name
	reqLogger.Info("Successfully finalized orderer")
	cmd := action.NewUninstall(cfg)
	resp, err := cmd.Run(releaseName)
	if err != nil {
		if strings.Compare("Release not loaded", err.Error()) != 0 {
			return nil
		}
		return err
	}
	log.Printf("Release %s deleted=%s", releaseName, resp.Info)
	return nil
}

func (r *FabricOrdererNodeReconciler) addFinalizer(reqLogger logr.Logger, m *hlfv1alpha1.FabricOrdererNode) error {
	reqLogger.Info("Adding Finalizer for the Orderer")
	controllerutil.AddFinalizer(m, ordererNodeFinalizer)

	// Update CR
	err := r.Update(context.TODO(), m)
	if err != nil {
		reqLogger.Error(err, "Failed to update Orderer with finalizer")
		return err
	}
	return nil
}

// +kubebuilder:rbac:groups=hlf.kungfusoftware.es,resources=fabricorderernodes,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=hlf.kungfusoftware.es,resources=fabricorderernodes/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=hlf.kungfusoftware.es,resources=fabricorderernodes/finalizers,verbs=get;update;patch
func (r *FabricOrdererNodeReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	reqLogger := r.Log.WithValues("hlf", req.NamespacedName)
	fabricOrdererNode := &hlfv1alpha1.FabricOrdererNode{}
	releaseName := req.Name
	ns := req.Namespace
	cfg, err := newActionCfg(r.Log, r.Config, ns)
	if err != nil {
		return ctrl.Result{}, err
	}
	err = r.Get(ctx, req.NamespacedName, fabricOrdererNode)
	if err != nil {
		log.Printf("Error getting the object %s error=%v", req.NamespacedName, err)
		if apierrors.IsNotFound(err) {
			reqLogger.Info("Orderer resource not found. Ignoring since object must be deleted.")
			return ctrl.Result{}, nil
		}
		reqLogger.Error(err, "Failed to get Orderer.")
		return ctrl.Result{}, err
	}
	isMemcachedMarkedToBeDeleted := fabricOrdererNode.GetDeletionTimestamp() != nil
	if isMemcachedMarkedToBeDeleted {
		if utils.Contains(fabricOrdererNode.GetFinalizers(), ordererNodeFinalizer) {
			if err := r.finalizeOrderer(reqLogger, fabricOrdererNode); err != nil {
				return ctrl.Result{}, err
			}
			controllerutil.RemoveFinalizer(fabricOrdererNode, ordererNodeFinalizer)
			err := r.Update(ctx, fabricOrdererNode)
			if err != nil {
				return ctrl.Result{}, err
			}
		}
		return ctrl.Result{}, nil
	}
	if !utils.Contains(fabricOrdererNode.GetFinalizers(), ordererNodeFinalizer) {
		if err := r.addFinalizer(reqLogger, fabricOrdererNode); err != nil {
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
		s, err := GetOrdererState(cfg, r.Config, releaseName, ns)
		if err != nil {
			log.Printf("Failed to get orderer state=%v", err)
			return ctrl.Result{}, err
		}
		fOrderer := fabricOrdererNode.DeepCopy()
		fOrderer.Status.Status = s.Status
		fOrderer.Status.URL = s.URL
		fOrderer.Status.Host = s.Host
		fOrderer.Status.Port = s.Port
		fOrderer.Status.TlsCert = s.TlsCert
		fOrderer.Status.TlsAdminCert = s.TlsAdminCert
		fOrderer.Status.AdminPort = s.AdminPort
		fOrderer.Status.OperationsPort = s.OperationsPort
		fOrderer.Status.Conditions.SetCondition(status.Condition{
			Type:   status.ConditionType(s.Status),
			Status: "True",
		})
		if reflect.DeepEqual(fOrderer.Status, fabricOrdererNode.Status) {
			log.Printf("Status hasn't changed, skipping update")
		} else {
			if err := r.Status().Update(ctx, fOrderer); err != nil {
				log.Printf("Error updating the status: %v", err)
				return ctrl.Result{}, err
			}
		}
		if s.Status == hlfv1alpha1.RunningStatus {
			return ctrl.Result{
				//RequeueAfter: 120 * time.Second,
			}, nil
		} else {
			return ctrl.Result{
				RequeueAfter: 2 * time.Second,
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
		c, err := getConfig(fabricOrdererNode, clientSet, releaseName, req.Namespace)
		if err != nil {
			reqLogger.Error(err, "Failed to get config for orderer %s/%s", req.Namespace, req.Name)
			return ctrl.Result{}, err
		}
		var inInterface map[string]interface{}
		inrec, err := json.Marshal(c)
		if err != nil {
			return ctrl.Result{}, err
		}
		log.Println(string(inrec))
		err = json.Unmarshal(inrec, &inInterface)
		if err != nil {
			log.Printf("Failed to unmarshall JSON %v", err)
			return ctrl.Result{}, err
		}
		if fabricOrdererNode.Spec.Genesis == "" && fabricOrdererNode.Spec.BootstrapMethod != "none" {
			waitForGenesis := 2 * time.Second
			log.Printf("Waiting %v since bootstrapMethod is %s", waitForGenesis, fabricOrdererNode.Spec.BootstrapMethod)
			return ctrl.Result{
				RequeueAfter: waitForGenesis,
			}, err
		}
		release, err := cmd.Run(ch, inInterface)
		if err != nil {
			log.Printf("Failed to run release %v", err)
			return ctrl.Result{}, err
		}
		log.Printf("Chart installed %s", release.Name)
		fabricOrdererNode.Status.Status = hlfv1alpha1.PendingStatus
		fabricOrdererNode.Status.Conditions.SetCondition(status.Condition{
			Type:               "DEPLOYED",
			Status:             "True",
			LastTransitionTime: v1.Time{},
		})
		if err := r.Status().Update(ctx, fabricOrdererNode); err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{
			Requeue:      false,
			RequeueAfter: 10 * time.Second,
		}, nil
	}
}

func (r *FabricOrdererNodeReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&hlfv1alpha1.FabricOrdererNode{}).
		Owns(&appsv1.Deployment{}).
		Complete(r)
}

func getExistingTLSAdminCrypto(client *kubernetes.Clientset, chartName string, namespace string) (*x509.Certificate, *ecdsa.PrivateKey, *x509.Certificate, *x509.Certificate, error) {
	secretName := fmt.Sprintf("%s-admin", chartName)
	secret, err := client.CoreV1().Secrets(namespace).Get(context.Background(), secretName, v1.GetOptions{})
	if err != nil {
		return nil, nil, nil, nil, err
	}
	// cacert.pem
	tlsKeyData := secret.Data["tls.key"]
	tlsCrtData := secret.Data["tls.crt"]
	rootTLSCrtData := secret.Data["cacert.crt"]
	clientRootCrtData := secret.Data["clientcacert.crt"]
	key, err := utils.ParseECDSAPrivateKey(tlsKeyData)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	crt, err := utils.ParseX509Certificate(tlsCrtData)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	rootCrt, err := utils.ParseX509Certificate(rootTLSCrtData)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	clientRootCrt, err := utils.ParseX509Certificate(clientRootCrtData)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	return crt, key, rootCrt, clientRootCrt, nil
}

func getExistingTLSCrypto(client *kubernetes.Clientset, chartName string, namespace string) (*x509.Certificate, *ecdsa.PrivateKey, *x509.Certificate, error) {
	secretName := fmt.Sprintf("%s-tls", chartName)
	tlsRootSecretName := fmt.Sprintf("%s-tlsrootcert", chartName)
	secret, err := client.CoreV1().Secrets(namespace).Get(context.Background(), secretName, v1.GetOptions{})
	if err != nil {
		return nil, nil, nil, err
	}
	rootCertSecret, err := client.CoreV1().Secrets(namespace).Get(context.Background(), tlsRootSecretName, v1.GetOptions{})
	if err != nil {
		return nil, nil, nil, err
	}
	// cacert.pem
	tlsKeyData := secret.Data["tls.key"]
	tlsCrtData := secret.Data["tls.crt"]
	rootTLSCrtData := rootCertSecret.Data["cacert.pem"]
	key, err := utils.ParseECDSAPrivateKey(tlsKeyData)
	if err != nil {
		return nil, nil, nil, err
	}
	crt, err := utils.ParseX509Certificate(tlsCrtData)
	if err != nil {
		return nil, nil, nil, err
	}
	rootCrt, err := utils.ParseX509Certificate(rootTLSCrtData)
	if err != nil {
		return nil, nil, nil, err
	}
	return crt, key, rootCrt, nil
}

func getExistingSignCrypto(client *kubernetes.Clientset, chartName string, namespace string) (*x509.Certificate, *ecdsa.PrivateKey, *x509.Certificate, error) {
	secretCrtName := fmt.Sprintf("%s-idcert", chartName)
	secretKeyName := fmt.Sprintf("%s-idkey", chartName)
	secretRootCrtName := fmt.Sprintf("%s-cacert", chartName)

	secretCrt, err := client.CoreV1().Secrets(namespace).Get(context.Background(), secretCrtName, v1.GetOptions{})
	if err != nil {
		return nil, nil, nil, err
	}
	secretKey, err := client.CoreV1().Secrets(namespace).Get(context.Background(), secretKeyName, v1.GetOptions{})
	if err != nil {
		return nil, nil, nil, err
	}
	secretRootCrt, err := client.CoreV1().Secrets(namespace).Get(context.Background(), secretRootCrtName, v1.GetOptions{})
	if err != nil {
		return nil, nil, nil, err
	}
	signCrtData := secretCrt.Data["cert.pem"]
	signKeyData := secretKey.Data["key.pem"]
	signRootCrtData := secretRootCrt.Data["cacert.pem"]
	crt, err := utils.ParseX509Certificate(signCrtData)
	if err != nil {
		return nil, nil, nil, err
	}
	rootCrt, err := utils.ParseX509Certificate(signRootCrtData)
	if err != nil {
		return nil, nil, nil, err
	}
	key, err := utils.ParseECDSAPrivateKey(signKeyData)
	if err != nil {
		return nil, nil, nil, err
	}
	return crt, key, rootCrt, nil
}

func CreateTLSCryptoMaterial(conf *hlfv1alpha1.FabricOrdererNode, caName string, caurl string, enrollID string, enrollSecret string, tlsCertString string, hosts []string) (*x509.Certificate, *ecdsa.PrivateKey, *x509.Certificate, error) {
	tlsCert, tlsKey, tlsRootCert, err := certs.EnrollUser(certs.EnrollUserRequest{
		TLSCert:    tlsCertString,
		URL:        caurl,
		Name:       caName,
		MSPID:      conf.Spec.MspID,
		User:       enrollID,
		Secret:     enrollSecret,
		Hosts:      hosts,
		CN:         "",
		Profile:    "tls",
		Attributes: nil,
	})
	if err != nil {
		return nil, nil, nil, err
	}
	return tlsCert, tlsKey, tlsRootCert, nil
}

func CreateTLSAdminCryptoMaterial(conf *hlfv1alpha1.FabricOrdererNode, caName string, caurl string, enrollID string, enrollSecret string, tlsCertString string, hosts []string) (*x509.Certificate, *ecdsa.PrivateKey, *x509.Certificate, *x509.Certificate, error) {
	tlsCert, tlsKey, tlsRootCert, err := certs.EnrollUser(
		certs.EnrollUserRequest{
			TLSCert:    tlsCertString,
			URL:        caurl,
			Name:       caName,
			MSPID:      conf.Spec.MspID,
			User:       enrollID,
			Secret:     enrollSecret,
			Hosts:      hosts,
			CN:         "",
			Profile:    "tls",
			Attributes: nil,
		},
	)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	return tlsCert, tlsKey, tlsRootCert, tlsRootCert, nil
}

func CreateSignCryptoMaterial(conf *hlfv1alpha1.FabricOrdererNode, caName string, caurl string, enrollID string, enrollSecret string, tlsCertString string) (*x509.Certificate, *ecdsa.PrivateKey, *x509.Certificate, error) {
	tlsCert, tlsKey, tlsRootCert, err := certs.EnrollUser(certs.EnrollUserRequest{
		TLSCert: tlsCertString,
		URL:     caurl,
		Name:    caName,
		MSPID:   conf.Spec.MspID,
		User:    enrollID,
		Secret:  enrollSecret,
	})
	if err != nil {
		return nil, nil, nil, err
	}
	return tlsCert, tlsKey, tlsRootCert, nil
}

func getConfig(conf *hlfv1alpha1.FabricOrdererNode, client *kubernetes.Clientset, chartName string, namespace string) (*fabricOrdChart, error) {
	spec := conf.Spec
	tlsParams := conf.Spec.Secret.Enrollment.TLS
	tlsCAUrl := fmt.Sprintf("https://%s:%d", tlsParams.Cahost, tlsParams.Caport)
	tlsHosts := []string{}
	ingressHosts := []string{}
	tlsHosts = append(tlsHosts, tlsParams.Csr.Hosts...)
	tlsCert, tlsKey, tlsRootCert, err := getExistingTLSCrypto(client, chartName, namespace)
	if err != nil {
		cacert, err := base64.StdEncoding.DecodeString(tlsParams.Catls.Cacert)
		if err != nil {
			return nil, err
		}
		tlsCert, tlsKey, tlsRootCert, err = CreateTLSCryptoMaterial(
			conf,
			tlsParams.Caname,
			tlsCAUrl,
			tlsParams.Enrollid,
			tlsParams.Enrollsecret,
			string(cacert),
			tlsHosts,
		)
		if err != nil {
			return nil, err
		}
	}

	adminCert, adminKey, adminRootCert, adminClientRootCert, err := getExistingTLSAdminCrypto(client, chartName, namespace)
	if err != nil {
		cacert, err := base64.StdEncoding.DecodeString(tlsParams.Catls.Cacert)
		if err != nil {
			return nil, err
		}
		adminCert, adminKey, adminRootCert, adminClientRootCert, err = CreateTLSAdminCryptoMaterial(
			conf,
			tlsParams.Caname,
			tlsCAUrl,
			tlsParams.Enrollid,
			tlsParams.Enrollsecret,
			string(cacert),
			tlsHosts,
		)
		if err != nil {
			return nil, err
		}
	}
	signParams := conf.Spec.Secret.Enrollment.Component
	caUrl := fmt.Sprintf("https://%s:%d", signParams.Cahost, signParams.Caport)
	signCert, signKey, signRootCert, err := getExistingSignCrypto(client, chartName, namespace)
	if err != nil {
		cacert, err := base64.StdEncoding.DecodeString(signParams.Catls.Cacert)
		if err != nil {
			return nil, err
		}
		signCert, signKey, signRootCert, err = CreateSignCryptoMaterial(
			conf,
			signParams.Caname,
			caUrl,
			signParams.Enrollid,
			signParams.Enrollsecret,
			string(cacert),
		)
		if err != nil {
			return nil, err
		}
	}
	tlsCRTEncoded := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: tlsCert.Raw,
	})
	tlsRootCRTEncoded := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: tlsRootCert.Raw,
	})
	tlsEncodedPK, err := utils.EncodePrivateKey(tlsKey)
	if err != nil {
		return nil, err
	}

	adminCRTEncoded := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: adminCert.Raw,
	})
	adminRootCRTEncoded := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: adminRootCert.Raw,
	})
	adminClientRootCRTEncoded := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: adminClientRootCert.Raw,
	})
	adminEncodedPK, err := utils.EncodePrivateKey(adminKey)
	if err != nil {
		return nil, err
	}

	signCRTEncoded := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: signCert.Raw,
	})
	signRootCRTEncoded := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: signRootCert.Raw,
	})
	signEncodedPK, err := utils.EncodePrivateKey(signKey)
	if err != nil {
		return nil, err
	}
	fabricOrdChart := fabricOrdChart{
		Ingress: ingress{
			Enabled: false,
		},
		Genesis:                     spec.Genesis,
		Cacert:                      string(signRootCRTEncoded),
		Tlsrootcert:                 string(tlsRootCRTEncoded),
		ChannelParticipationEnabled: spec.ChannelParticipationEnabled,
		BootstrapMethod:             string(spec.BootstrapMethod),
		Admin: admin{
			Cert:          string(adminCRTEncoded),
			Key:           string(adminEncodedPK),
			RootCAs:       string(adminRootCRTEncoded),
			ClientRootCAs: string(adminClientRootCRTEncoded),
		},
		Cert: string(signCRTEncoded),
		Key:  string(signEncodedPK),
		TLS: tls{
			Cert: string(tlsCRTEncoded),
			Key:  string(tlsEncodedPK),
		},
		FullnameOverride: conf.Name,
		HostAliases:      nil,
		Service: service{
			Type:               spec.Service.Type,
			Port:               7050,
			PortOperations:     9443,
			NodePort:           spec.Service.NodePortRequest,
			NodePortOperations: spec.Service.NodePortOperations,
		},
		Image: image{
			Repository: spec.Image,
			Tag:        spec.Tag,
			PullPolicy: string(spec.PullPolicy),
		},
		Persistence: persistence{
			Enabled:      true,
			Annotations:  annotations{},
			StorageClass: spec.Storage.StorageClass,
			AccessMode:   string(spec.Storage.AccessMode),
			Size:         spec.Storage.Size,
		},
		Ord: ord{
			Type:  "etcdraft",
			MspID: spec.MspID,
			TLS: tlsConfiguration{
				Server: ordServer{
					Enabled: true,
				},
				Client: ordClient{
					Enabled: false,
				},
			},
		},
		Clientcerts: clientcerts{},
		Hosts:       ingressHosts,
	}

	return &fabricOrdChart, nil
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

func GetOrdererState(conf *action.Configuration, config *rest.Config, releaseName string, ns string) (*hlfv1alpha1.FabricOrdererNodeStatus, error) {
	ctx := context.Background()
	cmd := action.NewGet(conf)
	rel, err := cmd.Run(releaseName)
	if err != nil {
		return nil, err
	}
	clientSet, err := utils.GetClientKubeWithConf(config)
	if err != nil {
		return nil, err
	}
	k8sIP, err := utils.GetPublicIPKubernetes(clientSet)
	if err != nil {
		return nil, err
	}
	r := &hlfv1alpha1.FabricOrdererNodeStatus{
		Status: hlfv1alpha1.RunningStatus,
	}
	tlsCrt, _, _, err := getExistingTLSCrypto(clientSet, releaseName, ns)
	if err != nil {
		return nil, err
	}
	r.TlsCert = string(utils.EncodeX509Certificate(tlsCrt))
	tlsAdminCrt, _, _, _, err := getExistingTLSAdminCrypto(clientSet, releaseName, ns)
	if err != nil {
		return nil, err
	}
	r.TlsAdminCert = string(utils.EncodeX509Certificate(tlsAdminCrt))
	objects := utils.ParseK8sYaml([]byte(rel.Manifest))
	for _, object := range objects {
		kind := object.GetObjectKind().GroupVersionKind().Kind
		if kind == "Deployment" {
			depSpec := object.(*appsv1.Deployment)
			dep, err := clientSet.AppsV1().Deployments(ns).Get(ctx, depSpec.Name, v1.GetOptions{})
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
						case corev1.PodSucceeded:
						case corev1.PodRunning:
							r.Status = hlfv1alpha1.RunningStatus
						case corev1.PodFailed:
							r.Status = hlfv1alpha1.FailedStatus
						case corev1.PodUnknown:
							r.Status = hlfv1alpha1.UnknownStatus
						}
					}
				}
			} else {
				if dep.Status.ReadyReplicas == *depSpec.Spec.Replicas {
					r.Status = hlfv1alpha1.RunningStatus
				} else {
					r.Status = hlfv1alpha1.PendingStatus
				}
			}
		} else if kind == "Service" {
			svcSpec := object.(*corev1.Service)
			svc, err := clientSet.CoreV1().Services(ns).Get(ctx, svcSpec.Name, v1.GetOptions{})
			if err != nil {
				return nil, err
			}
			for _, port := range svc.Spec.Ports {
				if port.Name == "grpc" {
					r.URL = fmt.Sprintf("grpcs://%s:%d", k8sIP, port.NodePort)
					r.Port = int(port.NodePort)
					r.Host = k8sIP
				} else if port.Name == "admin" {
					r.AdminPort = int(port.NodePort)
				} else if port.Name == "operations" {
					r.OperationsPort = int(port.NodePort)
				}
			}
		}
	}
	return r, nil
}
