package commit

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
	hlfv1alpha1 "github.com/kfsoftware/hlf-operator/api/hlf.kungfusoftware.es/v1alpha1"
	"github.com/kfsoftware/hlf-operator/controllers/utils"
	"github.com/kfsoftware/hlf-operator/internal/github.com/hyperledger/fabric/common/policydsl"
	"github.com/kfsoftware/hlf-operator/kubectl-hlf/cmd/helpers"
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
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"k8s.io/client-go/kubernetes"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/cryptosuite"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/cryptosuite/bccsp/sw"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab"
	mspimpl "github.com/hyperledger/fabric-sdk-go/pkg/msp"
	"gopkg.in/yaml.v2"
)

const chaincodeCommitFinalizer = "finalizer.chaincodeCommit.hlf.kungfusoftware.es"

type FabricChaincodeCommitReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
	Config *rest.Config
}

func (r *FabricChaincodeCommitReconciler) finalizeChaincodeCommit(reqLogger logr.Logger, m *hlfv1alpha1.FabricChaincodeCommit) error {
	reqLogger.Info("Successfully finalized ChaincodeCommit")
	return nil
}

func (r *FabricChaincodeCommitReconciler) addFinalizer(reqLogger logr.Logger, m *hlfv1alpha1.FabricChaincodeCommit) error {
	reqLogger.Info("Adding Finalizer for the ChaincodeCommit")
	controllerutil.AddFinalizer(m, chaincodeCommitFinalizer)

	// Update CR
	err := r.Update(context.TODO(), m)
	if err != nil {
		reqLogger.Error(err, "Failed to update ChaincodeCommit with finalizer")
		return err
	}
	return nil
}

// +kubebuilder:rbac:groups=hlf.kungfusoftware.es,resources=fabricchaincodecommits,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=hlf.kungfusoftware.es,resources=fabricchaincodecommits/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=hlf.kungfusoftware.es,resources=fabricchaincodecommits/finalizers,verbs=get;update;patch

func (r *FabricChaincodeCommitReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	reqLogger := r.Log.WithValues("hlf", req.NamespacedName)
	fabricChaincodeCommit := &hlfv1alpha1.FabricChaincodeCommit{}

	err := r.Get(ctx, req.NamespacedName, fabricChaincodeCommit)
	if err != nil {
		if apierrors.IsNotFound(err) {
			reqLogger.Info("FabricChaincodeCommit resource not found. Ignoring since object must be deleted.")
			return ctrl.Result{}, nil
		}
		reqLogger.Error(err, "Failed to get FabricChaincodeCommit.")
		return ctrl.Result{}, err
	}

	isMarkedToBeDeleted := fabricChaincodeCommit.GetDeletionTimestamp() != nil
	if isMarkedToBeDeleted {
		if utils.Contains(fabricChaincodeCommit.GetFinalizers(), chaincodeCommitFinalizer) {
			if err := r.finalizeChaincodeCommit(reqLogger, fabricChaincodeCommit); err != nil {
				return ctrl.Result{}, err
			}
			controllerutil.RemoveFinalizer(fabricChaincodeCommit, chaincodeCommitFinalizer)
			err := r.Update(ctx, fabricChaincodeCommit)
			if err != nil {
				return ctrl.Result{}, err
			}
		}
		return ctrl.Result{}, nil
	}

	if !utils.Contains(fabricChaincodeCommit.GetFinalizers(), chaincodeCommitFinalizer) {
		if err := r.addFinalizer(reqLogger, fabricChaincodeCommit); err != nil {
			return ctrl.Result{}, err
		}
	}

	clientSet, err := utils.GetClientKubeWithConf(r.Config)
	if err != nil {
		r.setConditionStatus(ctx, fabricChaincodeCommit, hlfv1alpha1.FailedStatus, false, err, false)
		return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricChaincodeCommit)
	}
	hlfClientSet, err := operatorv1.NewForConfig(r.Config)
	if err != nil {
		r.setConditionStatus(ctx, fabricChaincodeCommit, hlfv1alpha1.FailedStatus, false, err, false)
		return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricChaincodeCommit)
	}
	ncResponse, err := nc.GenerateNetworkConfigForChaincodeCommit(fabricChaincodeCommit, clientSet, hlfClientSet, fabricChaincodeCommit.Spec.MSPID)
	if err != nil {
		r.setConditionStatus(ctx, fabricChaincodeCommit, hlfv1alpha1.FailedStatus, false, err, false)
		return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricChaincodeCommit)
	}

	networkConfig := ncResponse.NetworkConfig
	resClient, sdk, err := getResmgmtBasedOnIdentity(ctx, fabricChaincodeCommit, networkConfig, clientSet, fabricChaincodeCommit.Spec.MSPID)
	if err != nil {
		r.setConditionStatus(ctx, fabricChaincodeCommit, hlfv1alpha1.FailedStatus, false, errors.Wrapf(err, "failed to get resmgmt"), false)
		return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricChaincodeCommit)
	}
	defer sdk.Close()

	var sp *common.SignaturePolicyEnvelope
	if fabricChaincodeCommit.Spec.EndorsementPolicy != "" {
		sp, err = policydsl.FromString(fabricChaincodeCommit.Spec.EndorsementPolicy)
		if err != nil {
			r.setConditionStatus(ctx, fabricChaincodeCommit, hlfv1alpha1.FailedStatus, false, err, false)
			return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricChaincodeCommit)
		}
	}
	var collectionConfigs []*pb.CollectionConfig

	if len(fabricChaincodeCommit.Spec.PrivateDataCollections) > 0 {
		collectionBytes, err := json.Marshal(fabricChaincodeCommit.Spec.PrivateDataCollections)
		if err != nil {
			r.setConditionStatus(ctx, fabricChaincodeCommit, hlfv1alpha1.FailedStatus, false, err, false)
			return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricChaincodeCommit)
		}
		collectionConfigs, err = helpers.GetCollectionConfigFromBytes(collectionBytes)
		if err != nil {
			r.setConditionStatus(ctx, fabricChaincodeCommit, hlfv1alpha1.FailedStatus, false, err, false)
			return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricChaincodeCommit)
		}
	}
	if len(collectionConfigs) == 0 {
		collectionConfigs = nil
	}

	txID, err := resClient.LifecycleCommitCC(
		fabricChaincodeCommit.Spec.ChannelName,
		resmgmt.LifecycleCommitCCRequest{
			Name:              fabricChaincodeCommit.Spec.ChaincodeName,
			Version:           fabricChaincodeCommit.Spec.Version,
			Sequence:          fabricChaincodeCommit.Spec.Sequence,
			EndorsementPlugin: "escc",
			ValidationPlugin:  "vscc",
			SignaturePolicy:   sp,
			CollectionConfig:  collectionConfigs,
			InitRequired:      fabricChaincodeCommit.Spec.InitRequired,
		},
		resmgmt.WithTimeout(fab2.ResMgmt, 20*time.Minute),
		resmgmt.WithTimeout(fab2.PeerResponse, 20*time.Minute),
	)
	if err != nil && !strings.Contains(err.Error(), "new definition must be sequence") {
		r.setConditionStatus(ctx, fabricChaincodeCommit, hlfv1alpha1.FailedStatus, false, err, false)
		return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricChaincodeCommit)
	}

	fabricChaincodeCommit.Status.Status = hlfv1alpha1.RunningStatus
	fabricChaincodeCommit.Status.Message = "Chaincode committed"
	if txID != "" {
		fabricChaincodeCommit.Status.TransactionID = string(txID)
	}
	fabricChaincodeCommit.Status.Conditions.SetCondition(status.Condition{
		Type:   status.ConditionType(hlfv1alpha1.RunningStatus),
		Status: corev1.ConditionTrue,
	})
	return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricChaincodeCommit)
}

func (r *FabricChaincodeCommitReconciler) setConditionStatus(ctx context.Context, p *hlfv1alpha1.FabricChaincodeCommit, conditionType hlfv1alpha1.DeploymentStatus, statusFlag bool, err error, statusUnknown bool) (update bool) {
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

func (r *FabricChaincodeCommitReconciler) updateCRStatusOrFailReconcile(ctx context.Context, log logr.Logger, p *hlfv1alpha1.FabricChaincodeCommit) (
	reconcile.Result, error) {
	if err := r.Status().Update(ctx, p); err != nil {
		log.Error(err, fmt.Sprintf("%v failed to update the application status", ErrClientK8s))
		return reconcile.Result{}, err
	}
	if p.Status.Status == hlfv1alpha1.FailedStatus {
		return reconcile.Result{
			RequeueAfter: 1 * time.Minute,
		}, nil
	}
	return reconcile.Result{}, nil
}

func (r *FabricChaincodeCommitReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&hlfv1alpha1.FabricChaincodeCommit{}).
		Complete(r)
}

type identity struct {
	Cert Pem `json:"cert"`
	Key  Pem `json:"key"`
}

type Pem struct {
	Pem string
}

func getResmgmtBasedOnIdentity(ctx context.Context, chInstall *hlfv1alpha1.FabricChaincodeCommit, networkConfig string, clientSet *kubernetes.Clientset, mspID string) (*resmgmt.Client, *fabsdk.FabricSDK, error) {
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
