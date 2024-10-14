package approve

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/go-logr/logr"
	"github.com/hyperledger/fabric-protos-go/common"
	pb "github.com/hyperledger/fabric-protos-go/peer"
	fab2 "github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/kfsoftware/hlf-operator/controllers/utils"
	"github.com/kfsoftware/hlf-operator/internal/github.com/hyperledger/fabric/common/policydsl"
	"github.com/kfsoftware/hlf-operator/kubectl-hlf/cmd/helpers"
	hlfv1alpha1 "github.com/kfsoftware/hlf-operator/pkg/apis/hlf.kungfusoftware.es/v1alpha1"
	operatorv1 "github.com/kfsoftware/hlf-operator/pkg/client/clientset/versioned"
	"github.com/kfsoftware/hlf-operator/pkg/nc"
	"github.com/kfsoftware/hlf-operator/pkg/status"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/cryptosuite"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/cryptosuite/bccsp/sw"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	mspimpl "github.com/hyperledger/fabric-sdk-go/pkg/msp"
	"gopkg.in/yaml.v2"
	"k8s.io/client-go/kubernetes"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const chaincodeApproveFinalizer = "finalizer.chaincodeapprove.hlf.kungfusoftware.es"

type FabricChaincodeApproveReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
	Config *rest.Config
}

func (r *FabricChaincodeApproveReconciler) finalizeChaincodeApprove(reqLogger logr.Logger, m *hlfv1alpha1.FabricChaincodeApprove) error {
	// TODO: no need to do anything when finalizing
	reqLogger.Info("Successfully finalized ChaincodeApprove")
	return nil
}

func (r *FabricChaincodeApproveReconciler) addFinalizer(reqLogger logr.Logger, m *hlfv1alpha1.FabricChaincodeApprove) error {
	reqLogger.Info("Adding Finalizer for the ChaincodeApprove")
	controllerutil.AddFinalizer(m, chaincodeApproveFinalizer)

	// Update CR
	err := r.Update(context.TODO(), m)
	if err != nil {
		reqLogger.Error(err, "Failed to update ChaincodeApprove with finalizer")
		return err
	}
	return nil
}

// +kubebuilder:rbac:groups=hlf.kungfusoftware.es,resources=fabricchaincodeapproves,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=hlf.kungfusoftware.es,resources=fabricchaincodeapproves/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=hlf.kungfusoftware.es,resources=fabricchaincodeapproves/finalizers,verbs=get;update;patch

func (r *FabricChaincodeApproveReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	reqLogger := r.Log.WithValues("hlf", req.NamespacedName)
	reqLogger.Info("Reconciling ChaincodeApprove")
	fabricChaincodeApprove := &hlfv1alpha1.FabricChaincodeApprove{}

	err := r.Get(ctx, req.NamespacedName, fabricChaincodeApprove)
	if err != nil {
		if apierrors.IsNotFound(err) {
			reqLogger.Info("FabricChaincodeApprove resource not found. Ignoring since object must be deleted.")
			return ctrl.Result{}, nil
		}
		reqLogger.Error(err, "Failed to get FabricChaincodeApprove.")
		return ctrl.Result{}, err
	}

	// Check if the FabricChaincodeApprove instance is marked to be deleted
	isMarkedToBeDeleted := fabricChaincodeApprove.GetDeletionTimestamp() != nil
	if isMarkedToBeDeleted {
		if utils.Contains(fabricChaincodeApprove.GetFinalizers(), chaincodeApproveFinalizer) {
			if err := r.finalizeChaincodeApprove(reqLogger, fabricChaincodeApprove); err != nil {
				return ctrl.Result{}, err
			}
			controllerutil.RemoveFinalizer(fabricChaincodeApprove, chaincodeApproveFinalizer)
			err := r.Update(ctx, fabricChaincodeApprove)
			if err != nil {
				return ctrl.Result{}, err
			}
		}
		return ctrl.Result{}, nil
	}

	// Add finalizer for this CR
	if !utils.Contains(fabricChaincodeApprove.GetFinalizers(), chaincodeApproveFinalizer) {
		if err := r.addFinalizer(reqLogger, fabricChaincodeApprove); err != nil {
			return ctrl.Result{}, err
		}
	}

	// TODO: Implement the logic for approving the chaincode
	// This should include:
	// 1. Getting the necessary clients (Kubernetes, HLF)
	// 2. Generating the network config
	// 3. Getting the resource management client
	// 4. Approving the chaincode
	// 5. Updating the status of the FabricChaincodeApprove resource

	// Example of how to update the status (you'll need to implement the actual logic):
	clientSet, err := utils.GetClientKubeWithConf(r.Config)
	if err != nil {
		r.setConditionStatus(ctx, fabricChaincodeApprove, hlfv1alpha1.FailedStatus, false, err, false)
		return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricChaincodeApprove)
	}
	hlfClientSet, err := operatorv1.NewForConfig(r.Config)
	if err != nil {
		r.setConditionStatus(ctx, fabricChaincodeApprove, hlfv1alpha1.FailedStatus, false, err, false)
		return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricChaincodeApprove)
	}
	ncResponse, err := nc.GenerateNetworkConfigForChaincodeApprove(fabricChaincodeApprove, clientSet, hlfClientSet, fabricChaincodeApprove.Spec.MSPID)
	if err != nil {
		r.setConditionStatus(ctx, fabricChaincodeApprove, hlfv1alpha1.FailedStatus, false, err, false)
		return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricChaincodeApprove)
	}

	networkConfig := ncResponse.NetworkConfig
	resClient, sdk, err := getResmgmtBasedOnIdentity(ctx, fabricChaincodeApprove, networkConfig, clientSet, fabricChaincodeApprove.Spec.MSPID)
	if err != nil {
		r.setConditionStatus(ctx, fabricChaincodeApprove, hlfv1alpha1.FailedStatus, false, errors.Wrapf(err, "failed to get resmgmt"), false)
		return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricChaincodeApprove)
	}
	defer sdk.Close()

	var sp *common.SignaturePolicyEnvelope
	if fabricChaincodeApprove.Spec.EndorsementPolicy != "" {
		sp, err = policydsl.FromString(fabricChaincodeApprove.Spec.EndorsementPolicy)
		if err != nil {
			r.setConditionStatus(ctx, fabricChaincodeApprove, hlfv1alpha1.FailedStatus, false, err, false)
			return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricChaincodeApprove)
		}
	}
	var collectionConfigs []*pb.CollectionConfig

	if len(fabricChaincodeApprove.Spec.PrivateDataCollections) > 0 {
		collectionBytes, err := json.Marshal(fabricChaincodeApprove.Spec.PrivateDataCollections)
		if err != nil {
			r.setConditionStatus(ctx, fabricChaincodeApprove, hlfv1alpha1.FailedStatus, false, err, false)
			return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricChaincodeApprove)
		}
		collectionConfigs, err = helpers.GetCollectionConfigFromBytes(collectionBytes)
		if err != nil {
			r.setConditionStatus(ctx, fabricChaincodeApprove, hlfv1alpha1.FailedStatus, false, err, false)
			return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricChaincodeApprove)
		}
	}
	if len(collectionConfigs) == 0 {
		collectionConfigs = nil
	}
	// get peerName of the first peer, either from peers or externalPeers
	var peerTarget string
	if len(fabricChaincodeApprove.Spec.Peers) > 0 {
		peerTarget = fmt.Sprintf("%s.%s", fabricChaincodeApprove.Spec.Peers[0].Name, fabricChaincodeApprove.Spec.Peers[0].Namespace)
	} else if len(fabricChaincodeApprove.Spec.ExternalPeers) > 0 {
		peerTarget = fabricChaincodeApprove.Spec.ExternalPeers[0].URL
	}
	if peerTarget == "" {
		r.setConditionStatus(ctx, fabricChaincodeApprove, hlfv1alpha1.FailedStatus, false, errors.New("peerTarget is empty"), false)
		return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricChaincodeApprove)
	}
	if fabricChaincodeApprove.Spec.Sequence > 1 {
		info, err := resClient.LifecycleQueryCommittedCC(
			fabricChaincodeApprove.Spec.ChannelName,
			resmgmt.LifecycleQueryCommittedCCRequest{
				Name: fabricChaincodeApprove.Spec.ChaincodeName,
			},
			resmgmt.WithTargetEndpoints(peerTarget),
		)

		if err != nil {
			r.setConditionStatus(ctx, fabricChaincodeApprove, hlfv1alpha1.FailedStatus, false, errors.Wrapf(err, "failed to query committed chaincode"), false)
			return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricChaincodeApprove)
		}
		log.Infof("info: %+v", info)
		lastSequence := info[0].Sequence
		if fabricChaincodeApprove.Spec.Sequence <= lastSequence {
			log.Infof("Sequence %d already committed", fabricChaincodeApprove.Spec.Sequence)
			fabricChaincodeApprove.Status.Status = hlfv1alpha1.RunningStatus
			fabricChaincodeApprove.Status.Message = "Chaincode already committed"
			fabricChaincodeApprove.Status.Conditions.SetCondition(status.Condition{
				Type:   status.ConditionType(hlfv1alpha1.RunningStatus),
				Status: corev1.ConditionTrue,
			})
			return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricChaincodeApprove)
		}
	}
	approveCCRequest := resmgmt.LifecycleApproveCCRequest{
		Name:              fabricChaincodeApprove.Spec.ChaincodeName,
		Version:           fabricChaincodeApprove.Spec.Version,
		PackageID:         fabricChaincodeApprove.Spec.PackageID,
		Sequence:          fabricChaincodeApprove.Spec.Sequence,
		EndorsementPlugin: "escc",
		ValidationPlugin:  "vscc",
		SignaturePolicy:   sp,
		CollectionConfig:  collectionConfigs,
		InitRequired:      fabricChaincodeApprove.Spec.InitRequired,
	}
	mustApprove := true
	// get current approved chaincode
	currentApprovedCC, err := resClient.LifecycleQueryApprovedCC(
		fabricChaincodeApprove.Spec.ChannelName,
		resmgmt.LifecycleQueryApprovedCCRequest{
			Name:     fabricChaincodeApprove.Spec.ChaincodeName,
			Sequence: fabricChaincodeApprove.Spec.Sequence,
		},
		resmgmt.WithTargetEndpoints(peerTarget),
	)
	if err == nil {
		mustApprove = currentApprovedCC.PackageID != fabricChaincodeApprove.Spec.PackageID || currentApprovedCC.Sequence != fabricChaincodeApprove.Spec.Sequence
	}

	log.Infof("currentApprovedCC: %+v", currentApprovedCC)
	log.Infof("approveCCRequest: %+v", approveCCRequest)

	log.Infof("mustApprove: %t", mustApprove)
	// compare currentApprovedCC with approveCCRequest and decide if we need to approve again
	if !mustApprove {
		r.setConditionStatus(ctx, fabricChaincodeApprove, hlfv1alpha1.RunningStatus, false, errors.New("chaincode already approved"), false)
		return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricChaincodeApprove)
	}

	txID, err := resClient.LifecycleApproveCC(
		fabricChaincodeApprove.Spec.ChannelName,
		approveCCRequest,
		resmgmt.WithTargetEndpoints(peerTarget),
		resmgmt.WithTimeout(fab2.ResMgmt, 20*time.Minute),
		resmgmt.WithTimeout(fab2.PeerResponse, 20*time.Minute),
	)
	if err != nil && (!strings.Contains(err.Error(), "attempted to redefine uncommitted") && !strings.Contains(err.Error(), "attempted to redefine the current committed")) {
		r.setConditionStatus(ctx, fabricChaincodeApprove, hlfv1alpha1.FailedStatus, false, errors.Wrapf(err, "failed to approve chaincode"), false)
		return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricChaincodeApprove)
	}
	r.Log.Info(fmt.Sprintf("ChaincodeApprove %s approved: %s", fabricChaincodeApprove.Name, txID))
	fabricChaincodeApprove.Status.Status = hlfv1alpha1.RunningStatus
	fabricChaincodeApprove.Status.Message = "Chaincode approved"
	if txID != "" {
		fabricChaincodeApprove.Status.TransactionID = string(txID)
	}
	fabricChaincodeApprove.Status.Conditions.SetCondition(status.Condition{
		Type:   status.ConditionType(hlfv1alpha1.RunningStatus),
		Status: corev1.ConditionTrue,
	})
	return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricChaincodeApprove)
}

func (r *FabricChaincodeApproveReconciler) setConditionStatus(ctx context.Context, p *hlfv1alpha1.FabricChaincodeApprove, conditionType hlfv1alpha1.DeploymentStatus, statusFlag bool, err error, statusUnknown bool) (update bool) {
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

var (
	ErrClientK8s = errors.New("k8sAPIClientError")
)

func (r *FabricChaincodeApproveReconciler) updateCRStatusOrFailReconcile(ctx context.Context, log logr.Logger, p *hlfv1alpha1.FabricChaincodeApprove) (
	reconcile.Result, error) {
	if err := r.Status().Update(ctx, p); err != nil {
		log.Error(err, fmt.Sprintf("%v failed to update the application status", ErrClientK8s))
		return reconcile.Result{
			Requeue:      false,
			RequeueAfter: 0,
		}, nil
	}
	if p.Status.Status == hlfv1alpha1.FailedStatus {
		return reconcile.Result{
			RequeueAfter: 1 * time.Minute,
		}, nil
	}
	r.Log.Info(fmt.Sprintf("Requeueing after 1 minute for %s", p.Name))
	return reconcile.Result{
		Requeue:      false,
		RequeueAfter: 0,
	}, nil
}

func (r *FabricChaincodeApproveReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&hlfv1alpha1.FabricChaincodeApprove{}).
		Complete(r)
}

type identity struct {
	Cert Pem `json:"cert"`
	Key  Pem `json:"key"`
}

type Pem struct {
	Pem string
}

func getResmgmtBasedOnIdentity(ctx context.Context, chInstall *hlfv1alpha1.FabricChaincodeApprove, networkConfig string, clientSet *kubernetes.Clientset, mspID string) (*resmgmt.Client, *fabsdk.FabricSDK, error) {
	configBackend := config.FromRaw([]byte(networkConfig), "yaml")
	sdk, err := fabsdk.New(configBackend)
	if err != nil {
		return nil, nil, err
	}
	idConfig := chInstall.Spec.HLFIdentity
	secret, err := clientSet.CoreV1().Secrets(idConfig.SecretNamespace).Get(ctx, idConfig.SecretName, v1.GetOptions{})
	if err != nil {
		return nil, nil, err
	}
	secretData, ok := secret.Data[idConfig.SecretKey]
	if !ok {

		return nil, nil, err
	}
	id := &identity{}
	err = yaml.Unmarshal(secretData, id)
	if err != nil {
		return nil, nil, err
	}
	sdkConfig, err := sdk.Config()
	if err != nil {
		return nil, nil, err
	}
	cryptoConfig := cryptosuite.ConfigFromBackend(sdkConfig)
	cryptoSuite, err := sw.GetSuiteByConfig(cryptoConfig)
	if err != nil {
		return nil, nil, err
	}
	userStore := mspimpl.NewMemoryUserStore()
	endpointConfig, err := fab.ConfigFromBackend(sdkConfig)
	if err != nil {
		return nil, nil, err
	}
	identityManager, err := mspimpl.NewIdentityManager(mspID, userStore, cryptoSuite, endpointConfig)
	if err != nil {
		return nil, nil, err
	}
	signingIdentity, err := identityManager.CreateSigningIdentity(
		msp.WithPrivateKey([]byte(id.Key.Pem)),
		msp.WithCert([]byte(id.Cert.Pem)),
	)
	if err != nil {
		return nil, nil, err
	}
	sdkContext := sdk.Context(
		fabsdk.WithIdentity(signingIdentity),
		fabsdk.WithOrg(mspID),
	)
	resClient, err := resmgmt.New(sdkContext)
	if err != nil {
		return nil, nil, err
	}
	return resClient, sdk, nil
}
