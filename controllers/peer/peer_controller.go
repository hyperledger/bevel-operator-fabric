package peer

import (
	"context"
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"helm.sh/helm/v3/pkg/release"
	"os"
	"reflect"
	"sort"
	"strings"
	"time"

	"github.com/kfsoftware/hlf-operator/controllers/hlfmetrics"
	"github.com/kfsoftware/hlf-operator/kubectl-hlf/cmd/helpers"

	operatorv1 "github.com/kfsoftware/hlf-operator/pkg/client/clientset/versioned"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller"

	"github.com/kfsoftware/hlf-operator/pkg/status"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"helm.sh/helm/v3/pkg/cli"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/kubernetes/pkg/api/v1/pod"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

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
	ChartPath                  string
	Log                        logr.Logger
	Scheme                     *runtime.Scheme
	Config                     *rest.Config
	AutoRenewCertificates      bool
	AutoRenewCertificatesDelta time.Duration
	Wait                       bool
	Timeout                    time.Duration
	MaxHistory                 int
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
	Status   hlfv1alpha1.DeploymentStatus
	TLSCert  string
	NodePort int
}

const (
	deploymentRestartTriggerAnnotation = "es.kungfusoftware.hlf.deployment-restart.timestamp"
)

func restartDeployment(config *rest.Config, deployment *appsv1.Deployment) error {
	clientSet, err := utils.GetClientKubeWithConf(config)
	if err != nil {
		return err
	}

	patchData := map[string]interface{}{}
	patchData["spec"] = map[string]interface{}{
		"template": map[string]interface{}{
			"metadata": map[string]interface{}{
				"annotations": map[string]interface{}{
					deploymentRestartTriggerAnnotation: time.Now().Format(time.Stamp),
				},
			},
		},
	}
	encodedData, err := json.Marshal(patchData)
	if err != nil {
		return err
	}
	_, err = clientSet.AppsV1().Deployments(deployment.Namespace).Patch(context.TODO(), deployment.Name, types.MergePatchType, encodedData, v1.PatchOptions{})
	if err != nil {
		return err
	}
	return nil
}
func GetPeerDeployment(conf *action.Configuration, config *rest.Config, releaseName string, ns string) (*appsv1.Deployment, error) {
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
	if ns == "" {
		ns = "default"
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
			return dep, nil
		}
	}
	return nil, errors.Errorf("Deployment not found")

}
func GetPeerState(conf *action.Configuration, config *rest.Config, releaseName string, ns string, svc *corev1.Service) (*hlfv1alpha1.FabricPeerStatus, error) {
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
	if ns == "" {
		ns = "default"
	}
	r := &hlfv1alpha1.FabricPeerStatus{
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
	for _, port := range svc.Spec.Ports {
		if port.Name == PeerPortName {
			r.NodePort = int(port.NodePort)
		}
	}

	tlsCrt, _, rootTlsCrt, err := getExistingTLSCrypto(clientSet, releaseName, ns)
	if err != nil {
		return nil, err
	}
	r.TlsCert = string(utils.EncodeX509Certificate(tlsCrt))
	r.TlsCACert = string(utils.EncodeX509Certificate(rootTlsCrt))
	hlfmetrics.UpdateCertificateExpiry(
		"peer",
		"tls",
		tlsCrt,
		releaseName,
		ns,
	)
	signCrt, _, rootSignCrt, err := getExistingSignCrypto(clientSet, releaseName, ns)
	if err != nil {
		return nil, err
	}
	r.SignCert = string(utils.EncodeX509Certificate(signCrt))
	r.SignCACert = string(utils.EncodeX509Certificate(rootSignCrt))
	hlfmetrics.UpdateCertificateExpiry(
		"peer",
		"sign",
		signCrt,
		releaseName,
		ns,
	)
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

// +kubebuilder:rbac:groups=monitoring.coreos.com,resources=servicemonitors,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=monitoring.coreos.com,resources=podmonitors,verbs=get;list;watch;create;update;patch;delete

// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=extensions,resources=deployments,verbs=get;list;watch;create;update;patch;delete

// +kubebuilder:rbac:groups=apps,resources=replicasets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=replicasets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=extensions,resources=replicasets,verbs=get;list;watch;create;update;patch;delete

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

func (r *FabricPeerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log.Infof("Reconciling FabricPeer %s/%s", req.Namespace, req.Name)
	defer func() {
		log.Infof("Reconciled FabricPeer %s/%s done", req.Namespace, req.Name)
	}()
	reqLogger := r.Log.WithValues("hlf", req.NamespacedName)
	fabricPeer := &hlfv1alpha1.FabricPeer{}
	releaseName := req.Name
	ns := req.Namespace
	cfg, err := newActionCfg(r.Log, r.Config, ns)
	if err != nil {
		r.setConditionStatus(ctx, fabricPeer, hlfv1alpha1.FailedStatus, false, err, false)
		return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricPeer)
	}
	err = r.Get(ctx, req.NamespacedName, fabricPeer)
	if err != nil {
		log.Debugf("Error getting the object %s error=%v", req.NamespacedName, err)
		if apierrors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			reqLogger.Info("Peer resource not found. Ignoring since object must be deleted.")

			return ctrl.Result{}, nil
		}
		reqLogger.Error(err, "Failed to get Peer.")
		r.setConditionStatus(ctx, fabricPeer, hlfv1alpha1.FailedStatus, false, err, false)
		return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricPeer)
	}

	isPeerMarkedToDelete := fabricPeer.GetDeletionTimestamp() != nil
	if isPeerMarkedToDelete {
		if utils.Contains(fabricPeer.GetFinalizers(), peerFinalizer) {
			if err := r.finalizePeer(reqLogger, fabricPeer); err != nil {
				r.setConditionStatus(ctx, fabricPeer, hlfv1alpha1.FailedStatus, false, err, false)
				return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricPeer)
			}
			controllerutil.RemoveFinalizer(fabricPeer, peerFinalizer)
			err := r.Update(ctx, fabricPeer)
			if err != nil {
				r.setConditionStatus(ctx, fabricPeer, hlfv1alpha1.FailedStatus, false, err, false)
				return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricPeer)
			}
		}
		return ctrl.Result{}, nil
	}
	if !utils.Contains(fabricPeer.GetFinalizers(), peerFinalizer) {
		if err := r.addFinalizer(reqLogger, fabricPeer); err != nil {
			r.setConditionStatus(ctx, fabricPeer, hlfv1alpha1.FailedStatus, false, err, false)
			return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricPeer)
		}
	}

	cmdStatus := action.NewStatus(cfg)
	exists := true
	helmStatus, err := cmdStatus.Run(releaseName)
	if err != nil {
		if errors.Is(err, driver.ErrReleaseNotFound) {
			// it doesn't exists
			exists = false
		} else {
			// it doesn't exist
			return ctrl.Result{}, err
		}
	}
	if exists && helmStatus.Info.Status == release.StatusPendingUpgrade {
		rollbackStatus := action.NewRollback(cfg)
		rollbackStatus.Version = helmStatus.Version - 1
		err = rollbackStatus.Run(releaseName)
		if err != nil {
			// it doesn't exist
			return ctrl.Result{}, err
		}
	} else if exists && helmStatus.Info.Status == release.StatusPendingRollback {
		historyAction := action.NewHistory(cfg)
		history, err := historyAction.Run(releaseName)
		if err != nil {
			return ctrl.Result{}, err
		}
		if len(history) > 0 {
			// find the last deployed revision
			// and rollback to it
			// sort history by revision number descending using raw go
			sort.Slice(history, func(i, j int) bool {
				return history[i].Version > history[j].Version
			})
			for _, historyItem := range history {
				if historyItem.Info.Status == release.StatusDeployed {
					rollbackStatus := action.NewRollback(cfg)
					rollbackStatus.Version = historyItem.Version
					err = rollbackStatus.Run(releaseName)
					if err != nil {
						// it doesn't exist
						return ctrl.Result{}, err
					}
					break
				}
			}
		}
	}
	log.Debugf("Release %s exists=%v", releaseName, exists)
	clientSet, err := utils.GetClientKubeWithConf(r.Config)
	if err != nil {
		r.setConditionStatus(ctx, fabricPeer, hlfv1alpha1.FailedStatus, false, err, false)
		return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricPeer)
	}
	svc, err := createPeerService(
		clientSet,
		chartName,
		fabricPeer,
	)
	if err != nil {
		r.setConditionStatus(ctx, fabricPeer, hlfv1alpha1.FailedStatus, false, err, false)
		return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricPeer)
	}
	reqLogger.Info(fmt.Sprintf("Service %s created", svc.Name))
	if exists {
		// update
		hlfClientSet, err := operatorv1.NewForConfig(r.Config)
		if err != nil {
			r.setConditionStatus(ctx, fabricPeer, hlfv1alpha1.FailedStatus, false, err, false)
			return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricPeer)
		}
		fabricPeerNode, err := hlfClientSet.HlfV1alpha1().FabricPeers(ns).Get(ctx, fabricPeer.Name, v1.GetOptions{})
		if err != nil {
			r.setConditionStatus(ctx, fabricPeer, hlfv1alpha1.FailedStatus, false, err, false)
			return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricPeer)
		}

		lastTimeCertsRenewed := fabricPeerNode.Status.LastCertificateUpdate
		certificatesNeedToBeRenewed := false
		if fabricPeerNode.Status.LastCertificateUpdate != nil && fabricPeerNode.Spec.UpdateCertificateTime != nil && fabricPeer.Spec.UpdateCertificateTime.Time.After(fabricPeerNode.Status.LastCertificateUpdate.Time) {
			certificatesNeedToBeRenewed = true
		}
		if lastTimeCertsRenewed == nil && fabricPeerNode.Spec.UpdateCertificateTime != nil {
			certificatesNeedToBeRenewed = true
		}
		// renew certificate 15 days before
		tlsCert, _, _, err := getExistingTLSCrypto(clientSet, releaseName, ns)
		if err != nil {
			return ctrl.Result{}, err
		}
		if r.AutoRenewCertificates && tlsCert.NotAfter.Before(time.Now().Add(r.AutoRenewCertificatesDelta)) {
			certificatesNeedToBeRenewed = true
		}
		requeueAfter := time.Second * 10
		log.Infof("Peer: Last time certs were updated: %v, they need to be renewed: %v", lastTimeCertsRenewed, certificatesNeedToBeRenewed)
		if certificatesNeedToBeRenewed {
			// must update the certificates and block until it's done
			// scale down to zero replicas
			// wait for the deployment to scale down
			// update the certs
			// scale up the peer
			log.Infof("Trying to upgrade certs")
			err := r.updateCerts(req, fabricPeer, clientSet, releaseName, svc, ctx, cfg, ns)
			if err != nil {
				log.Errorf("Error renewing certs: %v", err)
				r.setConditionStatus(ctx, fabricPeer, hlfv1alpha1.FailedStatus, false, err, false)
				return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricPeer)
			}
			now := v1.NewTime(time.Now().Add(time.Minute * 5)) // to avoid the upgrade certs to be triggered again
			lastTimeCertsRenewed = &now
			log.Infof("Peer certs updated, last time updated: %v", lastTimeCertsRenewed)
			requeueAfter = time.Minute * 5
		} else {
			c, err := GetConfig(fabricPeer, clientSet, releaseName, req.Namespace, svc, false)
			if err != nil {
				r.setConditionStatus(ctx, fabricPeer, hlfv1alpha1.FailedStatus, false, err, false)
				return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricPeer)
			}
			if helmStatus.Info.Status != release.StatusPendingUpgrade {
				err = r.upgradeChart(cfg, err, ns, releaseName, c)
				if err != nil {
					r.setConditionStatus(ctx, fabricPeer, hlfv1alpha1.FailedStatus, false, err, false)
					return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricPeer)
				}
			}
			requeueAfter = time.Minute * 60
		}
		s, err := GetPeerState(cfg, r.Config, releaseName, ns, svc)
		if err != nil {
			r.setConditionStatus(ctx, fabricPeer, hlfv1alpha1.FailedStatus, false, err, false)
			return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricPeer)
		}

		fPeer := fabricPeer.DeepCopy()
		fPeer.Status.Status = s.Status
		fPeer.Status.Message = ""
		fPeer.Status.TlsCert = s.TlsCert
		fPeer.Status.TlsCACert = s.TlsCACert
		fPeer.Status.SignCert = s.SignCert
		fPeer.Status.SignCACert = s.SignCACert
		fPeer.Status.NodePort = s.NodePort
		fPeer.Status.LastCertificateUpdate = lastTimeCertsRenewed
		fPeer.Status.Conditions.SetCondition(status.Condition{
			Type:   status.ConditionType(s.Status),
			Status: "True",
		})
		if !reflect.DeepEqual(fPeer.Status, fabricPeer.Status) {
			if err := r.Status().Update(ctx, fPeer); err != nil {
				log.Errorf("Error updating the status: %v", err)
				r.setConditionStatus(ctx, fabricPeer, hlfv1alpha1.FailedStatus, false, err, false)
				return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricPeer)
			}
		}
		log.Infof("Peer %s in %s status requeueAfter %v", fPeer.Name, string(s.Status), requeueAfter)
		switch s.Status {
		case hlfv1alpha1.PendingStatus:
			log.Infof("Peer %s in %s status", fPeer.Name, string(s.Status))
			return ctrl.Result{
				RequeueAfter: 10 * time.Second,
			}, nil
		case hlfv1alpha1.FailedStatus:
			log.Infof("Peer %s in %s status", fPeer.Name, string(s.Status))
			return ctrl.Result{
				RequeueAfter: requeueAfter,
			}, nil
		case hlfv1alpha1.RunningStatus:
			return ctrl.Result{
				RequeueAfter: requeueAfter,
			}, nil
		case hlfv1alpha1.UpdatingCertificates:
			return ctrl.Result{
				RequeueAfter: requeueAfter,
			}, nil
		default:
			return ctrl.Result{
				RequeueAfter: 2 * time.Second,
			}, nil
		}
	} else {
		cmd := action.NewInstall(cfg)
		cmd.Wait = r.Wait
		cmd.Timeout = r.Timeout
		name, chart, err := cmd.NameAndChart([]string{releaseName, r.ChartPath})
		if err != nil {
			r.setConditionStatus(ctx, fabricPeer, hlfv1alpha1.FailedStatus, false, err, false)
			return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricPeer)
		}

		cmd.ReleaseName = name
		ch, err := loader.Load(chart)
		if err != nil {
			r.setConditionStatus(ctx, fabricPeer, hlfv1alpha1.FailedStatus, false, err, false)
			return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricPeer)
		}
		c, err := GetConfig(
			fabricPeer,
			clientSet,
			name,
			req.Namespace,
			svc,
			false,
		)
		if err != nil {
			r.setConditionStatus(ctx, fabricPeer, hlfv1alpha1.FailedStatus, false, err, false)
			return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricPeer)
		}
		var inInterface map[string]interface{}
		inrec, err := json.Marshal(c)
		if err != nil {
			r.setConditionStatus(ctx, fabricPeer, hlfv1alpha1.FailedStatus, false, err, false)
			return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricPeer)
		}
		err = json.Unmarshal(inrec, &inInterface)
		if err != nil {
			r.setConditionStatus(ctx, fabricPeer, hlfv1alpha1.FailedStatus, false, err, false)
			return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricPeer)
		}
		release, err := cmd.Run(ch, inInterface)
		if err != nil {
			reqLogger.Error(err, "Failed to install chart")
			r.setConditionStatus(ctx, fabricPeer, hlfv1alpha1.FailedStatus, false, err, false)
			return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricPeer)
		}
		log.Infof("Chart installed %s", release.Name)
		err = r.Get(ctx, req.NamespacedName, fabricPeer)
		if err != nil {
			r.setConditionStatus(ctx, fabricPeer, hlfv1alpha1.FailedStatus, false, err, false)
			return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricPeer)
		}
		fabricPeer.Status.Status = hlfv1alpha1.PendingStatus
		fabricPeer.Status.Conditions.SetCondition(status.Condition{
			Type:               "DEPLOYED",
			Status:             "True",
			LastTransitionTime: v1.Time{},
		})
		if err := r.Status().Update(ctx, fabricPeer); err != nil {
			r.setConditionStatus(ctx, fabricPeer, hlfv1alpha1.FailedStatus, false, err, false)
			return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricPeer)
		}
		return ctrl.Result{
			Requeue:      false,
			RequeueAfter: 10 * time.Second,
		}, nil
	}
}

func (r *FabricPeerReconciler) updateCerts(req ctrl.Request, fPeer *hlfv1alpha1.FabricPeer, clientSet *kubernetes.Clientset, releaseName string, svc *corev1.Service, ctx context.Context, cfg *action.Configuration, ns string) error {
	log.Infof("Trying to upgrade certs")
	r.setConditionStatus(ctx, fPeer, hlfv1alpha1.UpdatingCertificates, false, nil, false)
	config, err := GetConfig(fPeer, clientSet, releaseName, req.Namespace, svc, true)
	if err != nil {
		log.Errorf("Error getting the config: %v", err)
		return err
	}
	err = r.upgradeChart(cfg, err, ns, releaseName, config)
	if err != nil {
		return err
	}
	dep, err := GetPeerDeployment(
		cfg,
		r.Config,
		releaseName,
		req.Namespace,
	)
	if err != nil {
		return err
	}
	err = restartDeployment(
		r.Config,
		dep,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *FabricPeerReconciler) upgradeChart(
	cfg *action.Configuration,
	err error,
	ns string,
	releaseName string,
	c *FabricPeerChart,
) error {
	inrec, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	var inInterface map[string]interface{}
	err = json.Unmarshal(inrec, &inInterface)
	if err != nil {
		return err
	}
	cmd := action.NewUpgrade(cfg)
	err = os.Setenv("HELM_NAMESPACE", ns)
	if err != nil {
		return err
	}
	settings := cli.New()
	chartPath, err := cmd.LocateChart(r.ChartPath, settings)
	ch, err := loader.Load(chartPath)
	if err != nil {
		return err
	}
	cmd.Wait = r.Wait
	cmd.MaxHistory = r.MaxHistory
	cmd.Timeout = r.Timeout
	log.Infof("Upgrading chart %s", inrec)
	release, err := cmd.Run(releaseName, ch, inInterface)
	if err != nil {
		return err
	}
	log.Infof("Chart upgraded %s", release.Name)
	return nil
}

func (r *FabricPeerReconciler) setConditionStatus(ctx context.Context, p *hlfv1alpha1.FabricPeer, conditionType hlfv1alpha1.DeploymentStatus, statusFlag bool, err error, statusUnknown bool) (update bool) {
	statusStr := func() corev1.ConditionStatus {
		if statusUnknown {
			return corev1.ConditionUnknown
		}
		if statusFlag {
			return corev1.ConditionTrue
		} else {
			return corev1.ConditionFalse
		}
	}
	if p.Status.Status != conditionType {
		depCopy := client.MergeFrom(p.DeepCopy())
		p.Status.Status = conditionType
		err = r.Status().Patch(ctx, p, depCopy)
		if err != nil {
			log.Warnf("Failed to update status to %s: %v", conditionType, err)
		}
	}
	if err != nil {
		p.Status.Message = err.Error()
	}
	condition := func() status.Condition {
		if err != nil {
			return status.Condition{
				Type:    status.ConditionType(conditionType),
				Status:  statusStr(),
				Reason:  status.ConditionReason(err.Error()),
				Message: err.Error(),
			}
		}
		return status.Condition{
			Type:   status.ConditionType(conditionType),
			Status: statusStr(),
		}
	}
	return p.Status.Conditions.SetCondition(condition())
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

var (
	ErrClientK8s = errors.New("k8sAPIClientError")
)

func (r *FabricPeerReconciler) updateCRStatusOrFailReconcile(ctx context.Context, log logr.Logger, p *hlfv1alpha1.FabricPeer) (
	reconcile.Result, error) {
	if err := r.Status().Update(ctx, p); err != nil {
		log.Error(err, fmt.Sprintf("%v failed to update the application status", ErrClientK8s))
		return reconcile.Result{}, err
	}
	return reconcile.Result{
		RequeueAfter: 1 * time.Minute,
	}, nil
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

func ReenrollSignCryptoMaterial(
	conf *hlfv1alpha1.FabricPeer,
	caName string,
	caurl string,
	enrollID string,
	tlsCertString string,
	signCertPem string,
	privateKey *ecdsa.PrivateKey,
) (*x509.Certificate, *ecdsa.PrivateKey, *x509.Certificate, error) {
	signCert, signRootCert, err := certs.ReenrollUser(
		certs.ReenrollUserRequest{
			TLSCert:  tlsCertString,
			URL:      caurl,
			Name:     caName,
			EnrollID: enrollID,
			MSPID:    conf.Spec.MspID,
		},
		signCertPem,
		privateKey,
	)
	if err != nil {
		return nil, nil, nil, err
	}
	return signCert, privateKey, signRootCert, nil
}

func ReenrollTLSCryptoMaterial(
	conf *hlfv1alpha1.FabricPeer,
	caName string,
	caurl string,
	enrollID string,
	tlsCertString string,
	hosts []string,
	tlsCertPem string,
	tlsKey *ecdsa.PrivateKey,
) (*x509.Certificate, *ecdsa.PrivateKey, *x509.Certificate, error) {
	tlsCert, tlsRootCert, err := certs.ReenrollUser(
		certs.ReenrollUserRequest{
			TLSCert:    tlsCertString,
			URL:        caurl,
			Name:       caName,
			MSPID:      conf.Spec.MspID,
			Hosts:      hosts,
			EnrollID:   enrollID,
			CN:         "",
			Profile:    "tls",
			Attributes: nil,
		},
		tlsCertPem,
		tlsKey,
	)
	if err != nil {
		return nil, nil, nil, err
	}
	return tlsCert, tlsKey, tlsRootCert, nil
}

func getCertBytesFromCATLS(client *kubernetes.Clientset, caTls hlfv1alpha1.Catls) ([]byte, error) {
	var signCertBytes []byte
	var err error
	if caTls.Cacert != "" {
		signCertBytes, err = base64.StdEncoding.DecodeString(caTls.Cacert)
		if err != nil {
			return nil, err
		}
	} else if caTls.SecretRef != nil {
		secret, err := client.CoreV1().Secrets(caTls.SecretRef.Namespace).Get(context.Background(), caTls.SecretRef.Name, v1.GetOptions{})
		if err != nil {
			return nil, err
		}
		signCertBytes = secret.Data[caTls.SecretRef.Key]
	} else {
		return nil, errors.New("invalid ca tls")
	}
	return signCertBytes, nil
}

func GetConfig(
	conf *hlfv1alpha1.FabricPeer,
	client *kubernetes.Clientset,
	chartName string,
	namespace string,
	svc *corev1.Service,
	refreshCerts bool,
) (*FabricPeerChart, error) {
	spec := conf.Spec
	tlsParams := conf.Spec.Secret.Enrollment.TLS
	tlsCAUrl := fmt.Sprintf("https://%s:%d", tlsParams.Cahost, tlsParams.Caport)
	ingressHosts := spec.Hosts
	var hosts []string
	hosts = append(hosts, tlsParams.Csr.Hosts...)
	hosts = append(hosts, ingressHosts...)
	var tlsCert, tlsRootCert, tlsOpsCert, signCert, signRootCert *x509.Certificate
	var tlsKey, tlsOpsKey, signKey *ecdsa.PrivateKey
	var err error
	ctx := context.Background()
	if tlsParams.External != nil {
		secret, err := client.CoreV1().Secrets(tlsParams.External.SecretNamespace).Get(ctx, tlsParams.External.SecretName, v1.GetOptions{})
		if err != nil {
			return nil, errors.Wrapf(err, "failed to get secret %s", tlsParams.External.SecretName)
		}
		tlsCert, err = utils.ParseX509Certificate(secret.Data[tlsParams.External.CertificateKey])
		if err != nil {
			return nil, errors.Wrapf(err, "failed to parse tls certificate")
		}
		tlsRootCert, err = utils.ParseX509Certificate(secret.Data[tlsParams.External.RootCertificateKey])
		if err != nil {
			return nil, errors.Wrapf(err, "failed to parse tls root certificate")
		}
		tlsKey, err = utils.ParseECDSAPrivateKey(secret.Data[tlsParams.External.PrivateKeyKey])
		if err != nil {
			return nil, errors.Wrapf(err, "failed to parse tls private key")
		}
	} else if refreshCerts {
		cacert, err := getCertBytesFromCATLS(client, tlsParams.Catls)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to decode tls ca cert")
		}
		tlsCert, tlsKey, tlsRootCert, err = getExistingTLSCrypto(client, chartName, namespace)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to get existing tls crypto material")
		}
		if tlsCert.NotAfter.Before(time.Now()) {
			log.Infof("Enrolling tls crypto material for %s", chartName)
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
			log.Infof("Successfully enrolled tls crypto material for %s", chartName)
		} else {
			tlsCert, tlsKey, tlsRootCert, err = ReenrollTLSCryptoMaterial(
				conf,
				tlsParams.Caname,
				tlsCAUrl,
				tlsParams.Enrollid,
				string(cacert),
				hosts,
				string(utils.EncodeX509Certificate(tlsCert)),
				tlsKey,
			)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to reenroll tls crypto material")
			}
			log.Infof("Successfully reenrolled tls crypto material for %s", chartName)
		}
		log.Infof("Successfully reenrolled tls crypto material for %s", chartName)
	} else {
		tlsCert, tlsKey, tlsRootCert, err = getExistingTLSCrypto(client, chartName, namespace)
		if err != nil {
			cacert, err := getCertBytesFromCATLS(client, tlsParams.Catls)
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
	}
	if refreshCerts {
		cacert, err := getCertBytesFromCATLS(client, tlsParams.Catls)
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
			hosts,
		)
		if err != nil {
			return nil, err
		}
	} else {
		tlsOpsCert, tlsOpsKey, _, err = getExistingTLSOPSCrypto(client, chartName, namespace)
		if err != nil {
			cacert, err := getCertBytesFromCATLS(client, tlsParams.Catls)
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
				hosts,
			)
			if err != nil {
				return nil, err
			}
		}
	}
	signParams := conf.Spec.Secret.Enrollment.Component
	caUrl := fmt.Sprintf("https://%s:%d", signParams.Cahost, signParams.Caport)
	if signParams.External != nil {
		secret, err := client.CoreV1().Secrets(signParams.External.SecretNamespace).Get(ctx, signParams.External.SecretName, v1.GetOptions{})
		if err != nil {
			return nil, errors.Wrapf(err, "failed to get secret %s", signParams.External.SecretName)
		}
		signCert, err = utils.ParseX509Certificate(secret.Data[signParams.External.CertificateKey])
		if err != nil {
			return nil, errors.Wrapf(err, "failed to parse sign certificate")
		}
		signRootCert, err = utils.ParseX509Certificate(secret.Data[signParams.External.RootCertificateKey])
		if err != nil {
			return nil, errors.Wrapf(err, "failed to parse sign root certificate")
		}
		signKey, err = utils.ParseECDSAPrivateKey(secret.Data[signParams.External.PrivateKeyKey])
		if err != nil {
			return nil, errors.Wrapf(err, "failed to parse sign private key")
		}
	} else if refreshCerts {
		cacert, err := getCertBytesFromCATLS(client, signParams.Catls)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to decode sign ca cert")
		}
		signCert, signKey, signRootCert, err = getExistingSignCrypto(client, chartName, namespace)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to get existing sign crypto material")
		}
		signCertPem := utils.EncodeX509Certificate(signCert)
		if signCert.NotAfter.Before(time.Now()) {
			log.Infof("Renewing certificates using enroll")
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
			log.Infof("Enrolled sign crypto material")
		} else {
			log.Infof("Renewing certificates using reenroll")
			signCert, signKey, signRootCert, err = ReenrollSignCryptoMaterial(
				conf,
				signParams.Caname,
				caUrl,
				signParams.Enrollid,
				string(cacert),
				string(signCertPem),
				signKey,
			)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to reenroll sign crypto material")
			}
			log.Infof("Reenrolled sign crypto material")
		}
		log.Infof("Reenrolled sign crypto material")
	} else {
		signCert, signKey, signRootCert, err = getExistingSignCrypto(client, chartName, namespace)
		if err != nil {
			cacert, err := getCertBytesFromCATLS(client, signParams.Catls)
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
	externalBuilders := []ExternalBuilder{}
	for _, builder := range spec.ExternalBuilders {
		externalBuilders = append(externalBuilders, ExternalBuilder{
			Name:                 builder.Name,
			Path:                 builder.Path,
			PropagateEnvironment: builder.PropagateEnvironment,
		})
	}
	imagePullPolicy := spec.ImagePullPolicy
	if imagePullPolicy == "" {
		imagePullPolicy = hlfv1alpha1.DefaultImagePullPolicy
	}
	var hostAliases []HostAlias
	for _, hostAlias := range spec.HostAliases {
		hostAliases = append(hostAliases, HostAlias{
			IP:        hostAlias.IP,
			Hostnames: hostAlias.Hostnames,
		})
	}
	stateDb := "goleveldb"
	switch spec.StateDb {
	case hlfv1alpha1.StateDBCouchDB:
		stateDb = "CouchDB"
	case hlfv1alpha1.StateDBLevelDB:
		stateDb = "goleveldb"
	case hlfv1alpha1.StateDBPostgres:
		stateDb = "pg"
	default:
		stateDb = "goleveldb"
	}
	var monitor ServiceMonitor
	if spec.ServiceMonitor != nil && spec.ServiceMonitor.Enabled {
		monitor = ServiceMonitor{
			Enabled:           spec.ServiceMonitor.Enabled,
			Labels:            spec.ServiceMonitor.Labels,
			Interval:          spec.ServiceMonitor.Interval,
			ScrapeTimeout:     spec.ServiceMonitor.ScrapeTimeout,
			Scheme:            "http",
			Relabelings:       []interface{}{},
			TargetLabels:      []interface{}{},
			MetricRelabelings: []interface{}{},
			SampleLimit:       spec.ServiceMonitor.SampleLimit,
		}
	} else {
		monitor = ServiceMonitor{Enabled: false}
	}
	var istio Istio
	if spec.Istio != nil {
		gateway := spec.Istio.IngressGateway
		if gateway == "" {
			gateway = "ingressgateway"
		}
		istio = Istio{
			Port:           spec.Istio.Port,
			Hosts:          spec.Istio.Hosts,
			IngressGateway: gateway,
		}
	} else {
		istio = Istio{
			Port:           0,
			Hosts:          []string{},
			IngressGateway: "",
		}
	}
	traefik := Traefik{}
	if spec.Traefik != nil {
		var middlewares []TraefikMiddleware
		if spec.Traefik.Middlewares != nil {
			for _, middleware := range spec.Traefik.Middlewares {
				middlewares = append(middlewares, TraefikMiddleware{
					Name:      middleware.Name,
					Namespace: middleware.Namespace,
				})
			}
		}
		traefik = Traefik{
			Entrypoints: spec.Traefik.Entrypoints,
			Middlewares: middlewares,
			Hosts:       spec.Traefik.Hosts,
		}
	}
	var gatewayApi GatewayApi
	if spec.GatewayApi != nil {
		gatewayApiName := spec.GatewayApi.GatewayName
		gatewayApiNamespace := spec.GatewayApi.GatewayNamespace
		if gatewayApiName == "" {
			gatewayApiName = "hlf-gateway"
		}
		if gatewayApiNamespace == "" {
			gatewayApiName = "default"
		}
		gatewayApi = GatewayApi{
			Port:             spec.GatewayApi.Port,
			Hosts:            spec.GatewayApi.Hosts,
			GatewayName:      gatewayApiName,
			GatewayNamespace: gatewayApiNamespace,
		}
	} else {
		gatewayApi = GatewayApi{
			Port:             443,
			Hosts:            []string{},
			GatewayName:      "",
			GatewayNamespace: "",
		}
	}
	exporter := CouchDBExporter{
		Enabled:    false,
		Image:      "",
		Tag:        "",
		PullPolicy: "",
	}
	if spec.CouchDBExporter != nil && spec.CouchDBExporter.Enabled {
		exporter.Enabled = spec.CouchDBExporter.Enabled
		if spec.CouchDBExporter.Image != "" {
			exporter.Image = spec.CouchDBExporter.Image
		}
		if spec.CouchDBExporter.Tag != "" {
			exporter.Tag = spec.CouchDBExporter.Tag
		}
		if spec.CouchDBExporter.ImagePullPolicy != "" {
			exporter.PullPolicy = string(spec.CouchDBExporter.ImagePullPolicy)
		}
	}

	couchDB := CouchDB{}
	if spec.CouchDB.ExternalCouchDB != nil && spec.CouchDB.ExternalCouchDB.Enabled {
		couchDB.External = CouchDBExternal{
			Enabled: true,
			Host:    spec.CouchDB.ExternalCouchDB.Host,
			Port:    spec.CouchDB.ExternalCouchDB.Port,
		}
	}
	if spec.CouchDB.Image != "" && spec.CouchDB.Tag != "" {
		couchDB.Image = spec.CouchDB.Image
		couchDB.Tag = spec.CouchDB.Tag
	} else {
		couchDB.Image = helpers.DefaultCouchDBImage
		couchDB.Tag = helpers.DefaultCouchDBVersion
	}
	if spec.CouchDB.PullPolicy != "" {
		couchDB.PullPolicy = string(spec.CouchDB.PullPolicy)
	} else {
		couchDB.PullPolicy = string(hlfv1alpha1.DefaultImagePullPolicy)
	}

	fsServer := FSServer{
		Image:      helpers.DefaultFSServerImage,
		Tag:        helpers.DefaultFSServerVersion,
		PullPolicy: string(hlfv1alpha1.DefaultImagePullPolicy),
	}
	if spec.FSServer != nil && spec.FSServer.Image != "" && spec.FSServer.Tag != "" {
		fsServer.Image = spec.FSServer.Image
		fsServer.Tag = spec.FSServer.Tag
		fsServer.PullPolicy = string(spec.FSServer.PullPolicy)
	} else {
		fsServer.Image = helpers.DefaultFSServerImage
		fsServer.Tag = helpers.DefaultFSServerVersion
		fsServer.PullPolicy = string(hlfv1alpha1.DefaultImagePullPolicy)
	}
	proxy := GRPCProxy{
		Enabled:          false,
		Image:            "",
		Tag:              "",
		PullPolicy:       "",
		ImagePullSecrets: nil,
		Istio:            Istio{},
	}
	if spec.GRPCProxy != nil && spec.GRPCProxy.Enabled {
		proxy = GRPCProxy{
			Enabled:          spec.GRPCProxy.Enabled,
			Image:            spec.GRPCProxy.Image,
			Tag:              spec.GRPCProxy.Tag,
			PullPolicy:       spec.GRPCProxy.Image,
			ImagePullSecrets: spec.GRPCProxy.ImagePullSecrets,
			Istio: Istio{
				Port:           spec.GRPCProxy.Istio.Port,
				Hosts:          spec.GRPCProxy.Istio.Hosts,
				IngressGateway: spec.GRPCProxy.Istio.IngressGateway,
			},
		}
	}
	peerResources := PeerResources{
		Peer:            spec.Resources.Peer,
		CouchDB:         spec.Resources.CouchDB,
		Chaincode:       spec.Resources.Chaincode,
		CouchDBExporter: spec.Resources.CouchDBExporter,
		Proxy:           spec.Resources.Proxy,
	}
	var c = FabricPeerChart{
		DeliveryClientaddressOverrides: spec.DeliveryClientaddressOverrides,
		Volumes:                        spec.Volumes,
		PeerVolumeMounts:               spec.PeerVolumeMounts,
		PodLabels:                      spec.PodLabels,
		PodAnnotations:                 spec.PodAnnotations,
		EnvVars:                        spec.Env,
		Replicas:                       spec.Replicas,
		ImagePullSecrets:               spec.ImagePullSecrets,
		GatewayApi:                     gatewayApi,
		Istio:                          istio,
		Traefik:                        traefik,
		Image: Image{
			Repository: spec.Image,
			Tag:        spec.Tag,
			PullPolicy: string(imagePullPolicy),
		},
		Proxy:            proxy,
		ServiceMonitor:   monitor,
		ExternalBuilders: externalBuilders,
		DockerSocketPath: spec.DockerSocketPath,
		CouchDBExporter:  exporter,
		CouchDB:          couchDB,
		FSServer:         fsServer,
		Peer: Peer{
			DatabaseType: stateDb,
			MspID:        spec.MspID,
			Gossip: Gossip{
				Bootstrap:              spec.Gossip.Bootstrap,
				Endpoint:               gossipEndpoint,
				ExternalEndpoint:       gossipExternalEndpoint,
				OrgLeader:              spec.Gossip.OrgLeader,
				UseLeaderElection:      spec.Gossip.UseLeaderElection,
				ReconnectInterval:      spec.Gossip.ReconnectInterval,
				AliveExpirationTimeout: spec.Gossip.AliveExpirationTimeout,
				AliveTimeInterval:      spec.Gossip.AliveTimeInterval,
				ResponseWaitTime:       spec.Gossip.ResponseWaitTime,
			},
			TLS: TLSAuth{
				Server: Server{Enabled: true},
				Client: Client{Enabled: false},
			},
		},
		ExternalChaincodeBuilder: conf.Spec.ExternalChaincodeBuilder,
		CouchdbPassword:          conf.Spec.CouchDB.Password,
		CouchdbUsername:          conf.Spec.CouchDB.User,
		Rbac:                     RBAC{Ns: namespace},
		Cert:                     string(signCRTEncoded),
		Key:                      string(signPEMEncodedPK),
		Hosts:                    ingressHosts,
		TLS: TLS{
			Cert: string(tlsCRTEncoded),
			Key:  string(tlsPEMEncodedPK),
		},
		OPSTLS: TLS{
			Cert: string(tlsOpsCRTEncoded),
			Key:  string(tlsOpsPEMEncodedPK),
		},
		Cacert:           string(signRootCRTEncoded),
		IntCacert:        ``,
		Tlsrootcert:      string(tlsRootCRTEncoded),
		Resources:        peerResources,
		NodeSelector:     spec.NodeSelector,
		Tolerations:      spec.Tolerations,
		Affinity:         spec.Affinity,
		ExternalHost:     externalEndpoint,
		FullnameOverride: conf.Name,
		HostAliases:      hostAliases,
		Service: Service{
			Type: string(spec.Service.Type),
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

func (r *FabricPeerReconciler) SetupWithManager(mgr ctrl.Manager, maxConcurrentReconciles int) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&hlfv1alpha1.FabricPeer{}).
		Owns(&appsv1.Deployment{}).
		WithOptions(controller.Options{
			MaxConcurrentReconciles: maxConcurrentReconciles,
		}).
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
	cmd.Wait = r.Wait
	cmd.Timeout = r.Timeout
	resp, err := cmd.Run(releaseName)
	if err != nil {
		if strings.Compare("Release not loaded", err.Error()) != 0 {
			return nil
		}
		log.Errorf("Failed to uninstall release %s %v", releaseName, err)
		return err
	}
	log.Infof("Release %s deleted=%s", releaseName, resp.Info)
	return nil
}

const PeerPortName = "peer"
const PeerPortName443 = "peer-443"
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
	labels := map[string]string{
		"app":     chartName,
		"release": releaseName,
	}
	serviceSpec := corev1.ServiceSpec{
		Type: peer.Spec.Service.Type,
		Ports: []corev1.ServicePort{
			{
				Name:     PeerPortName443,
				Protocol: "TCP",
				Port:     443,
				TargetPort: intstr.IntOrString{
					Type:   intstr.Int,
					IntVal: 7051,
				},
			},
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
		Selector: labels,
	}

	if exists {
		// update the service
		svc.Spec = serviceSpec
		return clientSet.CoreV1().Services(ns).Update(ctx, svc, v1.UpdateOptions{})
	}

	svc = &apiv1.Service{
		ObjectMeta: v1.ObjectMeta{
			Name:      svcName,
			Namespace: ns,
			Labels:    labels,
		},
		Spec:   serviceSpec,
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
