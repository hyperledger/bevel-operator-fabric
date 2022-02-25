package networkconfig

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/Masterminds/sprig/v3"
	"github.com/go-logr/logr"
	hlfv1alpha1 "github.com/kfsoftware/hlf-operator/api/hlf.kungfusoftware.es/v1alpha1"
	"github.com/kfsoftware/hlf-operator/controllers/utils"
	"github.com/kfsoftware/hlf-operator/kubectl-hlf/cmd/helpers"
	operatorv1 "github.com/kfsoftware/hlf-operator/pkg/client/clientset/versioned"
	"github.com/operator-framework/operator-lib/status"
	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"text/template"
)

// FabricNetworkConfigReconciler reconciles a FabricNetworkConfig object
type FabricNetworkConfigReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
	Config *rest.Config
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
  {{ range $mspID, $org := .Organizations }}
  {{$mspID}}:
    mspid: {{$mspID}}
    cryptoPath: /tmp/cryptopath
    users: {}
{{- if not $org.Peers }}
    peers: []
{{- else }}
    peers:
      {{- range $peer := $org.Peers }}
      - {{ $peer.Name }}
 	  {{- end }}
{{- end }}
{{- if not $org.OrderingServices }}
    orderers: []
{{- else }}
    orderers:
      {{- range $ordService := $org.OrderingServices }}
      {{- range $orderer := $ordService.Orderers }}
      - {{ $orderer.Name }}
 	  {{- end }}
 	  {{- end }}

    {{- end }}
{{- end }}
{{- end }}

{{- if not .Orderers }}
orderers: []
{{- else }}
orderers:
{{- range $ordService := .Orderers }}
{{- range $orderer := $ordService.Orderers }}
  {{$orderer.Name}}:
{{if $.Internal }}
    url: grpcs://{{ $orderer.PrivateURL }}
{{ else }}
    url: grpcs://{{ $orderer.PublicURL }}
{{ end }}
    grpcOptions:
      allow-insecure: false
    tlsCACerts:
      pem: |
{{ or $orderer.Status.TlsCACert $orderer.Status.TlsCert | indent 8 }}
{{- end }}
{{- end }}
{{- end }}

{{- if not .Peers }}
peers: []
{{- else }}
peers:
  {{- range $peer := .Peers }}
  {{$peer.Name}}:
{{if $.Internal }}
    url: grpcs://{{ $peer.PrivateURL }}
{{ else }}
    url: grpcs://{{ $peer.PublicURL }}
{{ end }}
    grpcOptions:
      hostnameOverride: ""
      ssl-target-name-override: ""
      allow-insecure: false
    tlsCACerts:
      pem: |
{{ $peer.Status.TlsCACert | indent 8 }}
{{- end }}
{{- end }}

{{- if not .CertAuths }}
certificateAuthorities: []
{{- else }}
certificateAuthorities:
{{- range $ca := .CertAuths }}
  
  {{ $ca.Name }}:
{{if $.Internal }}
    url: https://{{ $ca.PrivateURL }}
{{ else }}
    url: https://{{ $ca.PublicURL }}
{{ end }}
{{if $ca.EnrollID }}
    registrar:
        enrollId: {{ $ca.EnrollID }}
        enrollSecret: {{ $ca.EnrollPWD }}
{{ end }}
    caName: ca
    tlsCACerts:
      pem: 
       - |
{{ $ca.Status.TlsCert | indent 12 }}

{{- end }}
{{- end }}

channels:
  _default:
{{- if not .Orderers }}
    orderers: []
{{- else }}
    orderers:
{{- range $ordService := .Orderers }}
{{- range $orderer := $ordService.Orderers }}
      - {{$orderer.Name}}
{{- end }}
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

const networkConfigFinalizer = "finalizer.networkConfig.hlf.kungfusoftware.es"

func (r *FabricNetworkConfigReconciler) finalizeNetworkConfig(reqLogger logr.Logger, m *hlfv1alpha1.FabricNetworkConfig) error {
	ns := m.Namespace
	if ns == "" {
		ns = "default"
	}
	//releaseName := m.Name
	reqLogger.Info("Successfully finalized networkConfig")

	return nil
}

func (r *FabricNetworkConfigReconciler) addFinalizer(reqLogger logr.Logger, m *hlfv1alpha1.FabricNetworkConfig) error {
	reqLogger.Info("Adding Finalizer for the NetworkConfig")
	controllerutil.AddFinalizer(m, networkConfigFinalizer)

	// Update CR
	err := r.Update(context.TODO(), m)
	if err != nil {
		reqLogger.Error(err, "Failed to update NetworkConfig with finalizer")
		return err
	}
	return nil
}

// +kubebuilder:rbac:groups=hlf.kungfusoftware.es,resources=fabricnetworkconfigs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=hlf.kungfusoftware.es,resources=fabricnetworkconfigs/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=hlf.kungfusoftware.es,resources=fabricnetworkconfigs/finalizers,verbs=get;update;patch
func (r *FabricNetworkConfigReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	reqLogger := r.Log.WithValues("hlf", req.NamespacedName)
	fabricNetworkConfig := &hlfv1alpha1.FabricNetworkConfig{}
	//releaseName := req.Name

	err := r.Get(ctx, req.NamespacedName, fabricNetworkConfig)
	if err != nil {
		log.Debugf("Error getting the object %s error=%v", req.NamespacedName, err)
		if apierrors.IsNotFound(err) {
			reqLogger.Info("NetworkConfig resource not found. Ignoring since object must be deleted.")
			return ctrl.Result{}, nil
		}
		reqLogger.Error(err, "Failed to get NetworkConfig.")
		return ctrl.Result{}, err
	}
	isMemcachedMarkedToBeDeleted := fabricNetworkConfig.GetDeletionTimestamp() != nil
	if isMemcachedMarkedToBeDeleted {
		if utils.Contains(fabricNetworkConfig.GetFinalizers(), networkConfigFinalizer) {
			if err := r.finalizeNetworkConfig(reqLogger, fabricNetworkConfig); err != nil {
				return ctrl.Result{}, err
			}
			controllerutil.RemoveFinalizer(fabricNetworkConfig, networkConfigFinalizer)
			err := r.Update(ctx, fabricNetworkConfig)
			if err != nil {
				return ctrl.Result{}, err
			}
		}
		return ctrl.Result{}, nil
	}
	if !utils.Contains(fabricNetworkConfig.GetFinalizers(), networkConfigFinalizer) {
		if err := r.addFinalizer(reqLogger, fabricNetworkConfig); err != nil {
			return ctrl.Result{}, err
		}
	}
	tmpl, err := template.New("networkConfig").Funcs(sprig.HermeticTxtFuncMap()).Parse(tmplGoConfig)
	if err != nil {
		r.setConditionStatus(ctx, fabricNetworkConfig, hlfv1alpha1.FailedStatus, false, err, false)
		return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricNetworkConfig)
	}
	hlfClientSet, err := operatorv1.NewForConfig(r.Config)
	if err != nil {
		r.setConditionStatus(ctx, fabricNetworkConfig, hlfv1alpha1.FailedStatus, false, err, false)
		return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricNetworkConfig)
	}
	kubeClientset, err := kubernetes.NewForConfig(r.Config)
	if err != nil {
		r.setConditionStatus(ctx, fabricNetworkConfig, hlfv1alpha1.FailedStatus, false, err, false)
		return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricNetworkConfig)
	}
	var buf bytes.Buffer
	certAuths, err := helpers.GetClusterCAs(kubeClientset, hlfClientSet, "")
	if err != nil {
		r.setConditionStatus(ctx, fabricNetworkConfig, hlfv1alpha1.FailedStatus, false, err, false)
		return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricNetworkConfig)
	}
	ordOrgs, orderers, err := helpers.GetClusterOrderers(kubeClientset, hlfClientSet, "")
	if err != nil {
		r.setConditionStatus(ctx, fabricNetworkConfig, hlfv1alpha1.FailedStatus, false, err, false)
		return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricNetworkConfig)
	}
	peerOrgs, clusterPeers, err := helpers.GetClusterPeers(kubeClientset, hlfClientSet, "")
	if err != nil {
		r.setConditionStatus(ctx, fabricNetworkConfig, hlfv1alpha1.FailedStatus, false, err, false)
		return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricNetworkConfig)
	}
	orgMap := map[string]*helpers.Organization{}
	filterByOrgs := len(fabricNetworkConfig.Spec.Organizations) > 0
	for _, v := range ordOrgs {
		if (filterByOrgs && utils.Contains(fabricNetworkConfig.Spec.Organizations, v.MspID)) || !filterByOrgs {
			orgMap[v.MspID] = v
		}
	}
	for _, v := range peerOrgs {
		if (filterByOrgs && utils.Contains(fabricNetworkConfig.Spec.Organizations, v.MspID)) || !filterByOrgs {
			orgMap[v.MspID] = v
		}
	}
	var peers []*helpers.ClusterPeer
	for _, peer := range clusterPeers {
		if (filterByOrgs && utils.Contains(fabricNetworkConfig.Spec.Organizations, peer.MSPID)) || !filterByOrgs {
			peers = append(peers, peer)
		}
	}
	err = tmpl.Execute(&buf, map[string]interface{}{
		"Peers":         peers,
		"Orderers":      orderers,
		"Organizations": orgMap,
		"CertAuths":     certAuths,
		"Organization":  fabricNetworkConfig.Spec.Organization,
		"Internal":      fabricNetworkConfig.Spec.Internal,
	})
	if err != nil {
		r.setConditionStatus(ctx, fabricNetworkConfig, hlfv1alpha1.FailedStatus, false, err, false)
		return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricNetworkConfig)
	}
	ns := req.Namespace
	if ns == "" {
		ns = "default"
	}
	secretName := fabricNetworkConfig.Spec.SecretName
	secretData := map[string][]byte{
		"config.yaml": buf.Bytes(),
	}
	secret, err := kubeClientset.CoreV1().Secrets(ns).Get(ctx, secretName, metav1.GetOptions{})
	if err != nil {
		if apierrors.IsNotFound(err) {
			// creating secret
			secret, err = kubeClientset.CoreV1().Secrets(ns).Create(
				ctx,
				&corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Name:      secretName,
						Namespace: ns,
					},
					Data: secretData,
				},
				metav1.CreateOptions{},
			)
			if err != nil {
				r.setConditionStatus(ctx, fabricNetworkConfig, hlfv1alpha1.FailedStatus, false, err, false)
				return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricNetworkConfig)
			}
			return ctrl.Result{}, nil
		}
		r.setConditionStatus(ctx, fabricNetworkConfig, hlfv1alpha1.FailedStatus, false, err, false)
		return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricNetworkConfig)
	}
	secret.Data = secretData
	secret, err = kubeClientset.CoreV1().Secrets(ns).Update(
		ctx,
		secret,
		metav1.UpdateOptions{},
	)
	if err != nil {
		r.setConditionStatus(ctx, fabricNetworkConfig, hlfv1alpha1.FailedStatus, false, err, false)
		return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricNetworkConfig)
	}
	return ctrl.Result{}, nil
}

var (
	ErrClientK8s = errors.New("k8sAPIClientError")
)

func (r *FabricNetworkConfigReconciler) updateCRStatusOrFailReconcile(ctx context.Context, log logr.Logger, p *hlfv1alpha1.FabricNetworkConfig) (
	reconcile.Result, error) {
	if err := r.Status().Update(ctx, p); err != nil {
		log.Error(err, fmt.Sprintf("%v failed to update the application status", ErrClientK8s))
		return reconcile.Result{}, err
	}
	return reconcile.Result{}, nil
}

func (r *FabricNetworkConfigReconciler) setConditionStatus(ctx context.Context, p *hlfv1alpha1.FabricNetworkConfig, conditionType hlfv1alpha1.DeploymentStatus, statusFlag bool, err error, statusUnknown bool) (update bool) {
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

func (r *FabricNetworkConfigReconciler) SetupWithManager(mgr ctrl.Manager) error {
	managedBy := ctrl.NewControllerManagedBy(mgr)
	return managedBy.
		For(&hlfv1alpha1.FabricNetworkConfig{}).
		Owns(&corev1.Secret{}).
		Complete(r)
}

func actionLogger(logger logr.Logger) func(format string, v ...interface{}) {
	return func(format string, v ...interface{}) {
		logger.Info(fmt.Sprintf(format, v...))
	}
}
