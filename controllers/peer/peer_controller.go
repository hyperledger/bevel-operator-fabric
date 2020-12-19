package peer

import (
	"context"
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"github.com/operator-framework/operator-lib/status"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"helm.sh/helm/v3/pkg/cli"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/kubernetes/pkg/api/v1/pod"
	"os"
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"strings"
	"time"

	"github.com/go-logr/logr"
	hlfv1alpha1 "github.com/kfsoftware/hlf-operator/api/hlf.kungfusoftware.es/v1alpha1"
	"github.com/kfsoftware/hlf-operator/controllers/certs"
	"github.com/kfsoftware/hlf-operator/controllers/utils"
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
)

// FabricPeerReconciler reconciles a FabricPeer object
type FabricPeerReconciler struct {
	client.Client
	ChartPath string
	Log       logr.Logger
	Scheme    *runtime.Scheme
	Config    *rest.Config
}

func (r *FabricPeerReconciler) addFinalizer(reqLogger logr.Logger, m *hlfv1alpha1.FabricPeer) error {
	if len(m.GetFinalizers()) < 1 && m.GetDeletionTimestamp() == nil {
		reqLogger.Info("Adding Finalizer for the Memcached")
		m.SetFinalizers([]string{peerFinalizer})
		// Update CR
		err := r.Client.Update(context.TODO(), m)
		if err != nil {
			reqLogger.Error(err, "Failed to update Peer with finalizer")
			return err
		}
		reqLogger.Info(fmt.Sprintf("Finalizer for peer %s added", m.Name))
	}
	return nil
}

type Status struct {
	URL     string
	Status  hlfv1alpha1.DeploymentStatus
	TLSCert string
}

func GetPeerState(peer *hlfv1alpha1.FabricPeer, conf *action.Configuration, config *rest.Config, releaseName string, ns string, svc *corev1.Service) (*Status, error) {
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
	if ns == "" {
		ns = "default"
	}
	r := &Status{
		Status: hlfv1alpha1.PendingStatus,
	}
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
		}
	}
	if peer.Spec.ExternalEndpoint != "" {
		r.URL = fmt.Sprintf("grpcs://%s", peer.Spec.ExternalEndpoint)
	} else {
		for _, port := range svc.Spec.Ports {
			if port.Name == PeerPortName {
				r.URL = fmt.Sprintf("grpcs://%s:%d", k8sIP, port.NodePort)
			}
		}
	}
	tlsCrt, _, _, err := getExistingTLSCrypto(clientSet, releaseName, ns)
	if err != nil {
		return nil, err
	}
	r.TLSCert = string(utils.EncodeX509Certificate(tlsCrt))
	return r, nil
}

const peerFinalizer = "finalizer.peer.hlf.kungfusoftware.es"

const chartName = "hlf-peer"

// +kubebuilder:rbac:groups=hlf.kungfusoftware.es,resources=fabricpeers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=hlf.kungfusoftware.es,resources=fabricpeers/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=hlf.kungfusoftware.es,resources=fabricpeers/finalizers,verbs=get;update;patch

// +kubebuilder:rbac:groups=hlf.kungfusoftware.es,resources=fabricorderernodes,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=hlf.kungfusoftware.es,resources=fabricorderernodes/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=hlf.kungfusoftware.es,resources=fabricorderernodes/finalizers,verbs=get;update;patch

// +kubebuilder:rbac:groups=hlf.kungfusoftware.es,resources=fabricorderingservices,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=hlf.kungfusoftware.es,resources=fabricorderingservices/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=hlf.kungfusoftware.es,resources=fabricorderingservices/finalizers,verbs=get;update;patch

// +kubebuilder:rbac:groups=hlf.kungfusoftware.es,resources=fabriccas,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=hlf.kungfusoftware.es,resources=fabriccas/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=hlf.kungfusoftware.es,resources=fabriccas/finalizers,verbs=get;update;patch

// +kubebuilder:rbac:groups="",resources=nodes,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=serviceaccounts,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=clusterroles,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=clusterrolebindings,verbs=get;list;watch;create;update;patch;delete

// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=extensions,resources=deployments,verbs=get;list;watch;create;update;patch;delete

// +kubebuilder:rbac:groups=apps,resources=persistentvolumeclaims,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=persistentvolumeclaims,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=extensions,resources=persistentvolumeclaims,verbs=get;list;watch;create;update;patch;delete

// +kubebuilder:rbac:groups=apps,resources=pods,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=pods,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=extensions,resources=pods,verbs=get;list;watch;create;update;patch;delete

// +kubebuilder:rbac:groups=apps,resources=services,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=services,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=extensions,resources=services,verbs=get;list;watch;create;update;patch;delete

// +kubebuilder:rbac:groups=apps,resources=secrets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=extensions,resources=secrets,verbs=get;list;watch;create;update;patch;delete

// +kubebuilder:rbac:groups=apps,resources=pods/log,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=pods/log,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=extensions,resources=pods/log,verbs=get;list;watch;create;update;patch;delete

// +kubebuilder:rbac:groups=apps,resources=pods/status,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=pods/status,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=extensions,resources=pods/status,verbs=get;list;watch;create;update;patch;delete

// +kubebuilder:rbac:groups=apps,resources=configmaps,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=extensions,resources=configmaps,verbs=get;list;watch;create;update;patch;delete

// +kubebuilder:rbac:groups=networking.istio.io,resources=gateways,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=networking.istio.io,resources=virtualservices,verbs=get;list;watch;create;update;patch;delete

func (r *FabricPeerReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	reqLogger := r.Log.WithValues("hlf", req.NamespacedName)
	fabricPeer := &hlfv1alpha1.FabricPeer{}
	releaseName := req.Name
	ns := req.Namespace
	cfg, err := newActionCfg(r.Log, r.Config, ns)
	if err != nil {
		return ctrl.Result{}, err
	}
	err = r.Get(ctx, req.NamespacedName, fabricPeer)
	if err != nil {
		log.Printf("Error getting the object %s error=%v", req.NamespacedName, err)
		if apierrors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			reqLogger.Info("Peer resource not found. Ignoring since object must be deleted.")

			return ctrl.Result{}, nil
		}
		reqLogger.Error(err, "Failed to get Peer.")
		return ctrl.Result{}, err
	}

	isPeerMarkedToDelete := fabricPeer.GetDeletionTimestamp() != nil
	if isPeerMarkedToDelete {
		if utils.Contains(fabricPeer.GetFinalizers(), peerFinalizer) {
			if err := r.finalizePeer(reqLogger, fabricPeer); err != nil {
				return ctrl.Result{}, err
			}
			controllerutil.RemoveFinalizer(fabricPeer, peerFinalizer)
			err := r.Update(ctx, fabricPeer)
			if err != nil {
				return ctrl.Result{}, err
			}
		}
		return ctrl.Result{}, nil
	}
	if !utils.Contains(fabricPeer.GetFinalizers(), peerFinalizer) {
		if err := r.addFinalizer(reqLogger, fabricPeer); err != nil {
			return ctrl.Result{}, err
		}
	}

	cmdStatus := action.NewStatus(cfg)
	exists := true
	_, err = cmdStatus.Run(releaseName)
	if err != nil {
		if errors.Is(err, driver.ErrReleaseNotFound) {
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
	svc, err := createPeerService(
		clientSet,
		chartName,
		fabricPeer,
	)
	if err != nil {
		return ctrl.Result{}, err
	}

	reqLogger.Info(fmt.Sprintf("Service %s created", svc.Name))
	if exists {
		// update
		s, err := GetPeerState(
			fabricPeer,
			cfg,
			r.Config,
			releaseName,
			ns,
			svc,
		)
		if err != nil {
			return ctrl.Result{}, err
		}
		fPeer := fabricPeer.DeepCopy()
		fPeer.Status.Status = s.Status
		fPeer.Status.URL = s.URL
		fPeer.Status.TlsCert = s.TLSCert
		fPeer.Status.Conditions.SetCondition(status.Condition{
			Type:   status.ConditionType(s.Status),
			Status: "True",
		})
		if reflect.DeepEqual(fPeer.Status, fabricPeer.Status) {
			log.Printf("Status hasn't changed, skipping update")
			c, err := GetConfig(fabricPeer, clientSet, releaseName, req.Namespace, svc)
			if err != nil {
				return ctrl.Result{}, err
			}
			inrec, err := json.Marshal(c)
			if err != nil {
				return ctrl.Result{}, err
			}
			var inInterface map[string]interface{}
			err = json.Unmarshal(inrec, &inInterface)
			if err != nil {
				return ctrl.Result{}, err
			}
			cmd := action.NewUpgrade(cfg)
			err = os.Setenv("HELM_NAMESPACE", ns)
			if err != nil {
				return ctrl.Result{}, err
			}
			settings := cli.New()
			chartPath, err := cmd.LocateChart(r.ChartPath, settings)
			ch, err := loader.Load(chartPath)
			if err != nil {
				return ctrl.Result{}, err
			}
			release, err := cmd.Run(releaseName, ch, inInterface)
			if err != nil {
				return ctrl.Result{}, err
			}
			log.Infof("Chart upgraded %s", release.Name)
		} else {
			if err := r.Status().Update(ctx, fPeer); err != nil {
				log.Printf("Error updating the status: %v", err)
				return ctrl.Result{}, err
			}
		}
		log.Infof("Peer %s in %s status", fPeer.Name, string(s.Status))
		if s.Status == hlfv1alpha1.RunningStatus {
			return ctrl.Result{}, nil
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
		c, err := GetConfig(
			fabricPeer,
			clientSet,
			name,
			req.Namespace,
			svc,
		)
		if err != nil {
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
			return ctrl.Result{}, err
		}
		release, err := cmd.Run(ch, inInterface)
		if err != nil {
			reqLogger.Error(err, "Failed to install chart")
			return ctrl.Result{}, err
		}
		log.Printf("Chart installed %s", release.Name)
		fabricPeer.Status.Status = hlfv1alpha1.PendingStatus
		fabricPeer.Status.Conditions.SetCondition(status.Condition{
			Type:               "DEPLOYED",
			Status:             "True",
			LastTransitionTime: v1.Time{},
		})
		if err := r.Status().Update(ctx, fabricPeer); err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{
			Requeue:      false,
			RequeueAfter: 10 * time.Second,
		}, nil
	}
}

func getExistingTLSOPSCrypto(client *kubernetes.Clientset, chartName string, namespace string) (*x509.Certificate, *ecdsa.PrivateKey, *x509.Certificate, error) {
	secretName := fmt.Sprintf("%s-tls-ops", chartName)
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

func CreateTLSCryptoMaterial(conf *hlfv1alpha1.FabricPeer, caName string, caurl string, enrollID string, enrollSecret string, tlsCertString string, hosts []string) (*x509.Certificate, *ecdsa.PrivateKey, *x509.Certificate, error) {
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

func CreateTLSOPSCryptoMaterial(conf *hlfv1alpha1.FabricPeer, caName string, caurl string, enrollID string, enrollSecret string, tlsCertString string, hosts []string) (*x509.Certificate, *ecdsa.PrivateKey, *x509.Certificate, error) {
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
		return nil, nil, nil, err
	}
	return tlsCert, tlsKey, tlsRootCert, nil
}

func CreateSignCryptoMaterial(conf *hlfv1alpha1.FabricPeer, caName string, caurl string, enrollID string, enrollSecret string, tlsCertString string) (*x509.Certificate, *ecdsa.PrivateKey, *x509.Certificate, error) {
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

func GetConfig(conf *hlfv1alpha1.FabricPeer, client *kubernetes.Clientset, chartName string, namespace string, svc *corev1.Service) (*FabricPeerChart, error) {
	spec := conf.Spec
	tlsParams := conf.Spec.Secret.Enrollment.TLS
	tlsCAUrl := fmt.Sprintf("https://%s:%d", tlsParams.Cahost, tlsParams.Caport)
	ingressHosts := spec.Hosts
	operationHosts := spec.OperationHosts
	var hosts []string
	hosts = append(hosts, tlsParams.Csr.Hosts...)
	hosts = append(hosts, ingressHosts...)
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
			hosts,
		)
		if err != nil {
			return nil, err
		}
	}
	tlsOpsCert, tlsOpsKey, _, err := getExistingTLSOPSCrypto(client, chartName, namespace)
	if err != nil {
		cacert, err := base64.StdEncoding.DecodeString(tlsParams.Catls.Cacert)
		if err != nil {
			return nil, err
		}
		tlsOpsCert, tlsOpsKey, _, err = CreateTLSOPSCryptoMaterial(
			conf,
			tlsParams.Caname,
			tlsCAUrl,
			tlsParams.Enrollid,
			tlsParams.Enrollsecret,
			string(cacert),
			spec.OperationIPs,
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
	tlsEncodedPK, err := x509.MarshalPKCS8PrivateKey(tlsKey)
	if err != nil {
		return nil, err
	}
	tlsPEMEncodedPK := pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: tlsEncodedPK,
	})
	tlsOpsCRTEncoded := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: tlsOpsCert.Raw,
	})
	tlsOpsEncodedPK, err := x509.MarshalPKCS8PrivateKey(tlsOpsKey)
	if err != nil {
		return nil, err
	}
	tlsOpsPEMEncodedPK := pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: tlsOpsEncodedPK,
	})

	signCRTEncoded := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: signCert.Raw,
	})
	signRootCRTEncoded := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: signRootCert.Raw,
	})
	signEncodedPK, err := x509.MarshalPKCS8PrivateKey(signKey)
	if err != nil {
		return nil, err
	}
	signPEMEncodedPK := pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: signEncodedPK,
	})
	var externalEndpoint string
	if spec.ExternalEndpoint != "" {
		externalEndpoint = spec.ExternalEndpoint
	} else {
		requestNodePort, err := getRequestNodePort(svc)
		if err != nil {
			return nil, err
		}
		publicIP, err := utils.GetPublicIPKubernetes(client)
		if err != nil {
			return nil, err
		}
		externalEndpoint = fmt.Sprintf("%s:%d", publicIP, requestNodePort)
	}

	gossipExternalEndpoint := spec.Gossip.ExternalEndpoint
	if gossipExternalEndpoint == "" {
		gossipExternalEndpoint = externalEndpoint
	}
	gossipEndpoint := spec.Gossip.Endpoint
	if gossipEndpoint == "" {
		gossipEndpoint = externalEndpoint
	}
	var c = FabricPeerChart{
		Image: Image{
			Repository: spec.Image,
			Tag:        spec.Tag,
			PullPolicy: "Always",
		},
		DockerSocketPath: spec.DockerSocketPath,
		Peer: Peer{
			DatabaseType: string(spec.StateDb),
			MspID:        spec.MspID,
			Gossip: Gossip{
				Bootstrap:         spec.Gossip.Bootstrap,
				Endpoint:          gossipEndpoint,
				ExternalEndpoint:  gossipExternalEndpoint,
				OrgLeader:         spec.Gossip.OrgLeader,
				UseLeaderElection: spec.Gossip.UseLeaderElection,
			},
			TLS: TLSAuth{
				Server: Server{Enabled: true},
				Client: Client{Enabled: false},
			},
		},
		ExternalChaincodeBuilder: conf.Spec.ExternalChaincodeBuilder,
		CouchdbPassword:          conf.Spec.CouchDB.User,
		CouchdbUsername:          conf.Spec.CouchDB.Password,
		Rbac:                     RBAC{Ns: namespace},
		Cert:                     string(signCRTEncoded),
		Key:                      string(signPEMEncodedPK),
		Hosts:                    ingressHosts,
		OperationHosts:           operationHosts,
		TLS: TLS{
			Cert: string(tlsCRTEncoded),
			Key:  string(tlsPEMEncodedPK),
		},
		OPSTLS: TLS{
			Cert: string(tlsOpsCRTEncoded),
			Key:  string(tlsOpsPEMEncodedPK),
		},
		Cacert:      string(signRootCRTEncoded),
		Tlsrootcert: string(tlsRootCRTEncoded),
		Resources: Resources{
			Limits: Limits{
				CPU:    "100m",
				Memory: "128Mi",
			},
			Requests: Requests{
				CPU:    "100m",
				Memory: "128Mi",
			},
		},
		NodeSelector:     NodeSelector{},
		Tolerations:      nil,
		Affinity:         Affinity{},
		ExternalHost:     externalEndpoint,
		FullnameOverride: conf.Name,
		HostAliases:      nil,
		Service: Service{
			Type: spec.Service.Type,
		},
		Persistence: PeerPersistence{
			Peer: Persistence{
				Enabled:      true,
				Annotations:  Annotations{},
				StorageClass: spec.Storage.Peer.StorageClass,
				AccessMode:   string(spec.Storage.Peer.AccessMode),
				Size:         spec.Storage.Peer.Size,
			},
			CouchDB: Persistence{
				Enabled:      true,
				Annotations:  Annotations{},
				StorageClass: spec.Storage.CouchDB.StorageClass,
				AccessMode:   string(spec.Storage.CouchDB.AccessMode),
				Size:         spec.Storage.CouchDB.Size,
			},
			Chaincode: Persistence{
				Enabled:      true,
				Annotations:  Annotations{},
				StorageClass: spec.Storage.Chaincode.StorageClass,
				AccessMode:   string(spec.Storage.Chaincode.AccessMode),
				Size:         spec.Storage.Chaincode.Size,
			},
		},
		Logging: Logging{
			Level:    conf.Spec.Logging.Level,
			Peer:     conf.Spec.Logging.Peer,
			Cauthdsl: conf.Spec.Logging.Cauthdsl,
			Gossip:   conf.Spec.Logging.Gossip,
			Grpc:     conf.Spec.Logging.Grpc,
			Ledger:   conf.Spec.Logging.Ledger,
			Msp:      conf.Spec.Logging.Msp,
			Policies: conf.Spec.Logging.Policies,
		},
	}
	return &c, nil
}

func (r *FabricPeerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&hlfv1alpha1.FabricPeer{}).
		Owns(&appsv1.Deployment{}).
		Complete(r)
}
func getServiceName(peer *hlfv1alpha1.FabricPeer) string {
	return peer.Name
}
func (r *FabricPeerReconciler) finalizePeer(reqLogger logr.Logger, peer *hlfv1alpha1.FabricPeer) error {
	ns := peer.Namespace
	if ns == "" {
		ns = "default"
	}
	svcName := getServiceName(peer)
	clientSet, err := utils.GetClientKubeWithConf(r.Config)
	if err != nil {
		return err
	}
	ctx := context.Background()
	err = clientSet.CoreV1().Services(ns).Delete(ctx, svcName, v1.DeleteOptions{})
	if err != nil {
		if apierrors.IsNotFound(err) {
			reqLogger.Info(fmt.Sprintf("Service %s couldn't be found", svcName))
		} else {
			reqLogger.Info(fmt.Sprintf("Service %s couldn't be deleted: %v", svcName, err))
		}
	}
	cfg, err := newActionCfg(r.Log, r.Config, ns)
	if err != nil {
		return err
	}
	releaseName := peer.Name
	reqLogger.Info("Successfully finalized peer")
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

const PeerPortName = "peer"
const ChaincodePortName = "chaincode"
const EventPortName = "event"
const OperationsPortName = "operations"

func getRequestNodePort(svc *corev1.Service) (int, error) {
	for _, port := range svc.Spec.Ports {
		if port.Name == PeerPortName {
			return int(port.NodePort), nil
		}
	}
	return 0, errors.Errorf("")
}
func getReleaseName(peer *hlfv1alpha1.FabricPeer) string {
	return peer.Name
}
func getNamespace(peer *hlfv1alpha1.FabricPeer) string {
	ns := peer.Namespace
	if ns == "" {
		ns = "default"
	}
	return ns
}
func createPeerService(
	clientSet *kubernetes.Clientset,
	chartName string,
	peer *hlfv1alpha1.FabricPeer,
) (*apiv1.Service, error) {
	releaseName := getReleaseName(peer)
	ns := getNamespace(peer)
	ctx := context.Background()
	svcName := releaseName
	svc, err := clientSet.CoreV1().Services(ns).Get(
		ctx,
		svcName,
		v1.GetOptions{},
	)
	exists := true
	if err != nil {
		if apierrors.IsNotFound(err) {
			exists = false
		} else {
			return nil, err
		}
	}
	if exists {
		return svc, nil
	}
	svc = &apiv1.Service{
		ObjectMeta: v1.ObjectMeta{
			Name:      svcName,
			Namespace: ns,
			Labels:    map[string]string{},
		},
		Spec: corev1.ServiceSpec{
			Type: apiv1.ServiceType(peer.Spec.Service.Type),
			Ports: []corev1.ServicePort{
				{
					Name:     PeerPortName,
					Protocol: "TCP",
					Port:     7051,
					TargetPort: intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: 7051,
					},
				},
				{
					Name:     ChaincodePortName,
					Protocol: "TCP",
					Port:     7052,
					TargetPort: intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: 7052,
					},
				},
				{
					Name:     EventPortName,
					Protocol: "TCP",
					Port:     7053,
					TargetPort: intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: 7053,
					},
				},
				{
					Name:     OperationsPortName,
					Protocol: "TCP",
					Port:     9443,
					TargetPort: intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: 9443,
					},
				},
			},
			Selector: map[string]string{
				"app":     chartName,
				"release": releaseName,
			},
		},
		Status: corev1.ServiceStatus{},
	}
	return clientSet.CoreV1().Services(ns).Create(ctx, svc, v1.CreateOptions{})
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
