package followerchannel

import (
	"bytes"
	"context"
	"fmt"
	"github.com/go-logr/logr"
	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric-config/configtx"
	"github.com/hyperledger/fabric-config/protolator"
	"github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/cryptosuite"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/cryptosuite/bccsp/sw"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/resource"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	mspimpl "github.com/hyperledger/fabric-sdk-go/pkg/msp"
	"github.com/hyperledger/fabric/protoutil"
	hlfv1alpha1 "github.com/kfsoftware/hlf-operator/api/hlf.kungfusoftware.es/v1alpha1"
	"github.com/kfsoftware/hlf-operator/controllers/utils"
	operatorv1 "github.com/kfsoftware/hlf-operator/pkg/client/clientset/versioned"
	"github.com/kfsoftware/hlf-operator/pkg/nc"
	"github.com/kfsoftware/hlf-operator/pkg/status"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"strings"
)

// FabricFollowerChannelReconciler reconciles a FabricFollowerChannel object
type FabricFollowerChannelReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
	Config *rest.Config
}

const mainChannelFinalizer = "finalizer.mainChannel.hlf.kungfusoftware.es"

func (r *FabricFollowerChannelReconciler) finalizeMainChannel(reqLogger logr.Logger, m *hlfv1alpha1.FabricFollowerChannel) error {
	ns := m.Namespace
	if ns == "" {
		ns = "default"
	}
	reqLogger.Info("Successfully finalized mainChannel")

	return nil
}

func (r *FabricFollowerChannelReconciler) addFinalizer(reqLogger logr.Logger, m *hlfv1alpha1.FabricFollowerChannel) error {
	reqLogger.Info("Adding Finalizer for the MainChannel")
	controllerutil.AddFinalizer(m, mainChannelFinalizer)

	// Update CR
	err := r.Update(context.TODO(), m)
	if err != nil {
		reqLogger.Error(err, "Failed to update MainChannel with finalizer")
		return err
	}
	return nil
}

// +kubebuilder:rbac:groups=hlf.kungfusoftware.es,resources=fabricfollowerchannels,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=hlf.kungfusoftware.es,resources=fabricfollowerchannels/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=hlf.kungfusoftware.es,resources=fabricfollowerchannels/finalizers,verbs=get;update;patch
func (r *FabricFollowerChannelReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	reqLogger := r.Log.WithValues("hlf", req.NamespacedName)
	fabricFollowerChannel := &hlfv1alpha1.FabricFollowerChannel{}

	err := r.Get(ctx, req.NamespacedName, fabricFollowerChannel)
	if err != nil {
		log.Debugf("Error getting the object %s error=%v", req.NamespacedName, err)
		if apierrors.IsNotFound(err) {
			reqLogger.Info("MainChannel resource not found. Ignoring since object must be deleted.")
			return ctrl.Result{}, nil
		}
		reqLogger.Error(err, "Failed to get MainChannel.")
		return ctrl.Result{}, err
	}
	markedToBeDeleted := fabricFollowerChannel.GetDeletionTimestamp() != nil
	if markedToBeDeleted {
		if utils.Contains(fabricFollowerChannel.GetFinalizers(), mainChannelFinalizer) {
			if err := r.finalizeMainChannel(reqLogger, fabricFollowerChannel); err != nil {
				return ctrl.Result{}, err
			}
			controllerutil.RemoveFinalizer(fabricFollowerChannel, mainChannelFinalizer)
			err := r.Update(ctx, fabricFollowerChannel)
			if err != nil {
				return ctrl.Result{}, err
			}
		}
		return ctrl.Result{}, nil
	}
	if !utils.Contains(fabricFollowerChannel.GetFinalizers(), mainChannelFinalizer) {
		if err := r.addFinalizer(reqLogger, fabricFollowerChannel); err != nil {
			return ctrl.Result{}, err
		}
	}
	clientSet, err := utils.GetClientKubeWithConf(r.Config)
	if err != nil {
		r.setConditionStatus(ctx, fabricFollowerChannel, hlfv1alpha1.FailedStatus, false, err, false)
		return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricFollowerChannel)
	}
	hlfClientSet, err := operatorv1.NewForConfig(r.Config)
	if err != nil {
		r.setConditionStatus(ctx, fabricFollowerChannel, hlfv1alpha1.FailedStatus, false, err, false)
		return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricFollowerChannel)
	}

	// join peers
	mspID := fabricFollowerChannel.Spec.MSPID

	ncResponse, err := nc.GenerateNetworkConfigForFollower(fabricFollowerChannel, clientSet, hlfClientSet, mspID)
	if err != nil {
		r.setConditionStatus(ctx, fabricFollowerChannel, hlfv1alpha1.FailedStatus, false, errors.Wrapf(err, "failed to generate network config"), false)
		return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricFollowerChannel)
	}
	log.Infof("Generated network config: %s", ncResponse.NetworkConfig)
	configBackend := config.FromRaw([]byte(ncResponse.NetworkConfig), "yaml")
	sdk, err := fabsdk.New(configBackend)
	if err != nil {
		r.setConditionStatus(ctx, fabricFollowerChannel, hlfv1alpha1.FailedStatus, false, err, false)
		return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricFollowerChannel)
	}
	defer sdk.Close()
	idConfig := fabricFollowerChannel.Spec.HLFIdentity
	secret, err := clientSet.CoreV1().Secrets(idConfig.SecretNamespace).Get(ctx, idConfig.SecretName, v1.GetOptions{})
	if err != nil {
		r.setConditionStatus(ctx, fabricFollowerChannel, hlfv1alpha1.FailedStatus, false, err, false)
		return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricFollowerChannel)
	}
	secretData, ok := secret.Data[idConfig.SecretKey]
	if !ok {
		r.setConditionStatus(ctx, fabricFollowerChannel, hlfv1alpha1.FailedStatus, false, fmt.Errorf("secret key %s not found", idConfig.SecretKey), false)
		return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricFollowerChannel)
	}
	id := &identity{}
	err = yaml.Unmarshal(secretData, id)
	if err != nil {
		r.setConditionStatus(ctx, fabricFollowerChannel, hlfv1alpha1.FailedStatus, false, err, false)
		return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricFollowerChannel)
	}
	sdkConfig, err := sdk.Config()
	if err != nil {
		r.setConditionStatus(ctx, fabricFollowerChannel, hlfv1alpha1.FailedStatus, false, err, false)
		return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricFollowerChannel)
	}
	cryptoConfig := cryptosuite.ConfigFromBackend(sdkConfig)
	cryptoSuite, err := sw.GetSuiteByConfig(cryptoConfig)
	if err != nil {
		r.setConditionStatus(ctx, fabricFollowerChannel, hlfv1alpha1.FailedStatus, false, err, false)
		return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricFollowerChannel)
	}
	userStore := mspimpl.NewMemoryUserStore()
	endpointConfig, err := fab.ConfigFromBackend(sdkConfig)
	if err != nil {
		r.setConditionStatus(ctx, fabricFollowerChannel, hlfv1alpha1.FailedStatus, false, err, false)
		return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricFollowerChannel)
	}
	identityManager, err := mspimpl.NewIdentityManager(mspID, userStore, cryptoSuite, endpointConfig)
	if err != nil {
		r.setConditionStatus(ctx, fabricFollowerChannel, hlfv1alpha1.FailedStatus, false, err, false)
		return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricFollowerChannel)
	}
	signingIdentity, err := identityManager.CreateSigningIdentity(
		msp.WithPrivateKey([]byte(id.Key.Pem)),
		msp.WithCert([]byte(id.Cert.Pem)),
	)
	if err != nil {
		r.setConditionStatus(ctx, fabricFollowerChannel, hlfv1alpha1.FailedStatus, false, err, false)
		return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricFollowerChannel)
	}
	sdkContext := sdk.Context(
		fabsdk.WithIdentity(signingIdentity),
		fabsdk.WithOrg(mspID),
	)
	resClient, err := resmgmt.New(sdkContext)
	if err != nil {
		r.setConditionStatus(ctx, fabricFollowerChannel, hlfv1alpha1.FailedStatus, false, err, false)
		return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricFollowerChannel)
	}
	for _, peer := range fabricFollowerChannel.Spec.PeersToJoin {
		r.Log.Info(fmt.Sprintf("Joining peer %s namespace %s", peer.Name, peer.Namespace))
		err = resClient.JoinChannel(
			fabricFollowerChannel.Spec.Name,
			resmgmt.WithTargetEndpoints(fmt.Sprintf("%s.%s", peer.Name, peer.Namespace)),
		)
		if err != nil {
			if strings.Contains(err.Error(), "already exists") {
				r.Log.Info(fmt.Sprintf("Peer %s already joined channel %s", peer.Name, fabricFollowerChannel.Spec.Name))
				continue
			}
			r.setConditionStatus(ctx, fabricFollowerChannel, hlfv1alpha1.FailedStatus, false, err, false)
			return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricFollowerChannel)
		}
	}
	for _, peer := range fabricFollowerChannel.Spec.ExternalPeersToJoin {
		r.Log.Info(fmt.Sprintf("Joining peer %s", peer.URL))
		err = resClient.JoinChannel(
			fabricFollowerChannel.Spec.Name,
			resmgmt.WithTargetEndpoints(peer.URL),
		)
		if err != nil {
			if strings.Contains(err.Error(), "already exists") {
				r.Log.Info(fmt.Sprintf("Peer %s already joined channel %s", peer.URL, fabricFollowerChannel.Spec.Name))
				continue
			}
			r.setConditionStatus(ctx, fabricFollowerChannel, hlfv1alpha1.FailedStatus, false, err, false)
			return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricFollowerChannel)
		}
	}

	// set anchor peers
	block, err := resClient.QueryConfigBlockFromOrderer(fabricFollowerChannel.Spec.Name)
	if err != nil {
		r.Log.Info(fmt.Sprintf("Failed to get block %v", err))
		r.setConditionStatus(ctx, fabricFollowerChannel, hlfv1alpha1.FailedStatus, false, errors.Wrapf(err, "failed to get block from channel %s", fabricFollowerChannel.Spec.Name), false)
		return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricFollowerChannel)
	}
	cfgBlock, err := resource.ExtractConfigFromBlock(block)
	if err != nil {
		r.setConditionStatus(ctx, fabricFollowerChannel, hlfv1alpha1.FailedStatus, false, err, false)
		return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricFollowerChannel)
	}
	var buf2 bytes.Buffer
	err = protolator.DeepMarshalJSON(&buf2, cfgBlock)
	if err != nil {
		r.setConditionStatus(ctx, fabricFollowerChannel, hlfv1alpha1.FailedStatus, false, errors.Wrapf(err, "error converting block to JSON"), false)
		return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricFollowerChannel)
	}
	log.Infof("Config block: %s", buf2.Bytes())
	cftxGen := configtx.New(cfgBlock)
	app := cftxGen.Application().Organization(mspID)
	anchorPeers, err := app.AnchorPeers()
	if err != nil {
		r.setConditionStatus(ctx, fabricFollowerChannel, hlfv1alpha1.FailedStatus, false, err, false)
		return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricFollowerChannel)
	}
	r.Log.Info(fmt.Sprintf("Old anchor peers %v", anchorPeers))

	for _, anchorPeer := range anchorPeers {
		err = app.RemoveAnchorPeer(configtx.Address{
			Host: anchorPeer.Host,
			Port: anchorPeer.Port,
		})
		if err != nil {
			r.setConditionStatus(ctx, fabricFollowerChannel, hlfv1alpha1.FailedStatus, false, err, false)
			return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricFollowerChannel)
		}
	}
	r.Log.Info(fmt.Sprintf("New anchor peers %v", anchorPeers))

	for _, anchorPeer := range fabricFollowerChannel.Spec.AnchorPeers {
		err = app.AddAnchorPeer(configtx.Address{
			Host: anchorPeer.Host,
			Port: anchorPeer.Port,
		})
		if err != nil {
			r.setConditionStatus(ctx, fabricFollowerChannel, hlfv1alpha1.FailedStatus, false, err, false)
			return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricFollowerChannel)
		}
	}
	configUpdateBytes, err := cftxGen.ComputeMarshaledUpdate(fabricFollowerChannel.Spec.Name)
	if err != nil {
		if !strings.Contains(err.Error(), "no differences detected between original and updated config") {
			r.setConditionStatus(ctx, fabricFollowerChannel, hlfv1alpha1.FailedStatus, false, errors.Wrapf(err, "error calculating config update"), false)
			return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricFollowerChannel)
		}
		r.Log.Info("No differences detected between original and updated config")
	} else {
		configUpdate := &common.ConfigUpdate{}
		err = proto.Unmarshal(configUpdateBytes, configUpdate)
		if err != nil {
			r.setConditionStatus(ctx, fabricFollowerChannel, hlfv1alpha1.FailedStatus, false, err, false)
			return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricFollowerChannel)
		}
		channelConfigBytes, err := CreateConfigUpdateEnvelope(fabricFollowerChannel.Spec.Name, configUpdate)
		if err != nil {
			r.setConditionStatus(ctx, fabricFollowerChannel, hlfv1alpha1.FailedStatus, false, err, false)
			return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricFollowerChannel)
		}
		configUpdateReader := bytes.NewReader(channelConfigBytes)
		chResponse, err := resClient.SaveChannel(resmgmt.SaveChannelRequest{
			ChannelID:     fabricFollowerChannel.Spec.Name,
			ChannelConfig: configUpdateReader,
		})
		if err != nil {
			r.setConditionStatus(ctx, fabricFollowerChannel, hlfv1alpha1.FailedStatus, false, err, false)
			return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricFollowerChannel)
		}
		log.Infof("anchor anchorPeers added: %s", chResponse.TransactionID)
	}

	// update config map with the configuration
	ordererChannelBlock, err := resClient.QueryConfigBlockFromOrderer(fabricFollowerChannel.Spec.Name)
	if err != nil {
		r.setConditionStatus(ctx, fabricFollowerChannel, hlfv1alpha1.FailedStatus, false, errors.Wrapf(err, "error fetching block from orderer"), false)
		return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricFollowerChannel)
	}
	cmnConfig, err := resource.ExtractConfigFromBlock(ordererChannelBlock)
	if err != nil {
		r.setConditionStatus(ctx, fabricFollowerChannel, hlfv1alpha1.FailedStatus, false, errors.Wrapf(err, "error extracting the config from block"), false)
		return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricFollowerChannel)
	}
	var buf bytes.Buffer
	err = protolator.DeepMarshalJSON(&buf, cmnConfig)
	if err != nil {
		r.setConditionStatus(ctx, fabricFollowerChannel, hlfv1alpha1.FailedStatus, false, errors.Wrapf(err, "error converting block to JSON"), false)
		return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricFollowerChannel)
	}
	configMapNamespace := "default"
	configMapName := fmt.Sprintf("%s-follower-config", fabricFollowerChannel.ObjectMeta.Name)
	createConfigMap := false
	configMap, err := clientSet.CoreV1().ConfigMaps(configMapNamespace).Get(ctx, configMapName, v1.GetOptions{})
	if err != nil {
		if apierrors.IsNotFound(err) {
			reqLogger.Info(fmt.Sprintf("ConfigMap %s not found, creating it", configMapName))
			createConfigMap = true
		} else {
			r.setConditionStatus(ctx, fabricFollowerChannel, hlfv1alpha1.FailedStatus, false, errors.Wrapf(err, "error getting configmap"), false)
			return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricFollowerChannel)
		}
	}
	if createConfigMap {
		_, err = clientSet.CoreV1().ConfigMaps(configMapNamespace).Create(ctx, &corev1.ConfigMap{
			TypeMeta: v1.TypeMeta{},
			ObjectMeta: v1.ObjectMeta{
				Name:      configMapName,
				Namespace: configMapNamespace,
			},
			Data: map[string]string{
				"channel.json": buf.String(),
			},
		}, v1.CreateOptions{})
		if err != nil {
			r.setConditionStatus(ctx, fabricFollowerChannel, hlfv1alpha1.FailedStatus, false, errors.Wrapf(err, "error creating config map"), false)
			return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricFollowerChannel)
		}
	} else {
		configMap.Data["channel.json"] = buf.String()
		_, err = clientSet.CoreV1().ConfigMaps(configMapNamespace).Update(ctx, configMap, v1.UpdateOptions{})
		if err != nil {
			r.setConditionStatus(ctx, fabricFollowerChannel, hlfv1alpha1.FailedStatus, false, errors.Wrapf(err, "error updating config map"), false)
			return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricFollowerChannel)
		}
	}

	fabricFollowerChannel.Status.Status = hlfv1alpha1.RunningStatus
	fabricFollowerChannel.Status.Message = "Peers and anchor peers completed"
	fabricFollowerChannel.Status.Conditions.SetCondition(status.Condition{
		Type:   status.ConditionType(fabricFollowerChannel.Status.Status),
		Status: "True",
	})
	if err := r.Status().Update(ctx, fabricFollowerChannel); err != nil {
		r.setConditionStatus(ctx, fabricFollowerChannel, hlfv1alpha1.FailedStatus, false, err, false)
		return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricFollowerChannel)
	}
	return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricFollowerChannel)
}

var (
	ErrClientK8s = errors.New("k8sAPIClientError")
)

func (r *FabricFollowerChannelReconciler) updateCRStatusOrFailReconcile(ctx context.Context, log logr.Logger, p *hlfv1alpha1.FabricFollowerChannel) (
	reconcile.Result, error) {
	if err := r.Status().Update(ctx, p); err != nil {
		log.Error(err, fmt.Sprintf("%v failed to update the application status", ErrClientK8s))
		return reconcile.Result{}, err
	}
	return reconcile.Result{}, nil
}

func (r *FabricFollowerChannelReconciler) setConditionStatus(ctx context.Context, p *hlfv1alpha1.FabricFollowerChannel, conditionType hlfv1alpha1.DeploymentStatus, statusFlag bool, err error, statusUnknown bool) (update bool) {
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

func (r *FabricFollowerChannelReconciler) SetupWithManager(mgr ctrl.Manager) error {
	managedBy := ctrl.NewControllerManagedBy(mgr)
	return managedBy.
		For(&hlfv1alpha1.FabricFollowerChannel{}).
		Owns(&corev1.Secret{}).
		Complete(r)
}

type identity struct {
	Cert Pem `json:"cert"`
	Key  Pem `json:"key"`
}

type Pem struct {
	Pem string
}

func CreateConfigUpdateEnvelope(channelID string, configUpdate *common.ConfigUpdate) ([]byte, error) {
	configUpdate.ChannelId = channelID
	configUpdateData, err := proto.Marshal(configUpdate)
	if err != nil {
		return nil, err
	}
	configUpdateEnvelope := &common.ConfigUpdateEnvelope{}
	configUpdateEnvelope.ConfigUpdate = configUpdateData
	envelope, err := protoutil.CreateSignedEnvelope(common.HeaderType_CONFIG_UPDATE, channelID, nil, configUpdateEnvelope, 0, 0)
	if err != nil {
		return nil, err
	}
	envelopeData, err := proto.Marshal(envelope)
	if err != nil {
		return nil, err
	}
	return envelopeData, nil
}
