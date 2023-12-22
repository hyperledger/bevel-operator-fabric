package mainchannel

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"fmt"
	"github.com/go-logr/logr"
	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric-config/configtx"
	"github.com/hyperledger/fabric-config/configtx/membership"
	"github.com/hyperledger/fabric-config/configtx/orderer"
	"github.com/hyperledger/fabric-config/protolator"
	"github.com/hyperledger/fabric-protos-go/common"
	cb "github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric-protos-go/orderer/smartbft"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	fab2 "github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
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
	"github.com/kfsoftware/hlf-operator/kubectl-hlf/cmd/helpers"
	"github.com/kfsoftware/hlf-operator/kubectl-hlf/cmd/helpers/osnadmin"
	operatorv1 "github.com/kfsoftware/hlf-operator/pkg/client/clientset/versioned"
	"github.com/kfsoftware/hlf-operator/pkg/nc"
	"github.com/kfsoftware/hlf-operator/pkg/status"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"net"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"strconv"
	"strings"
	"time"
)

// FabricMainChannelReconciler reconciles a FabricMainChannel object
type FabricMainChannelReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
	Config *rest.Config
}

const mainChannelFinalizer = "finalizer.mainChannel.hlf.kungfusoftware.es"

func (r *FabricMainChannelReconciler) finalizeMainChannel(reqLogger logr.Logger, m *hlfv1alpha1.FabricMainChannel) error {
	ns := m.Namespace
	if ns == "" {
		ns = "default"
	}
	reqLogger.Info("Successfully finalized mainChannel")

	return nil
}

func (r *FabricMainChannelReconciler) addFinalizer(reqLogger logr.Logger, m *hlfv1alpha1.FabricMainChannel) error {
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

// +kubebuilder:rbac:groups=hlf.kungfusoftware.es,resources=fabricmainchannels,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=hlf.kungfusoftware.es,resources=fabricmainchannels/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=hlf.kungfusoftware.es,resources=fabricmainchannels/finalizers,verbs=get;update;patch
func (r *FabricMainChannelReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	reqLogger := r.Log.WithValues("hlf", req.NamespacedName)
	fabricMainChannel := &hlfv1alpha1.FabricMainChannel{}

	err := r.Get(ctx, req.NamespacedName, fabricMainChannel)
	if err != nil {
		log.Debugf("Error getting the object %s error=%v", req.NamespacedName, err)
		if apierrors.IsNotFound(err) {
			reqLogger.Info("MainChannel resource not found. Ignoring since object must be deleted.")
			return ctrl.Result{}, nil
		}
		reqLogger.Error(err, "Failed to get MainChannel.")
		return ctrl.Result{}, err
	}
	markedToBeDeleted := fabricMainChannel.GetDeletionTimestamp() != nil
	if markedToBeDeleted {
		if utils.Contains(fabricMainChannel.GetFinalizers(), mainChannelFinalizer) {
			if err := r.finalizeMainChannel(reqLogger, fabricMainChannel); err != nil {
				return ctrl.Result{}, err
			}
			controllerutil.RemoveFinalizer(fabricMainChannel, mainChannelFinalizer)
			err := r.Update(ctx, fabricMainChannel)
			if err != nil {
				return ctrl.Result{}, err
			}
		}
		return ctrl.Result{}, nil
	}
	if !utils.Contains(fabricMainChannel.GetFinalizers(), mainChannelFinalizer) {
		if err := r.addFinalizer(reqLogger, fabricMainChannel); err != nil {
			return ctrl.Result{}, err
		}
	}
	clientSet, err := utils.GetClientKubeWithConf(r.Config)
	if err != nil {
		r.setConditionStatus(ctx, fabricMainChannel, hlfv1alpha1.FailedStatus, false, err, false)
		return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricMainChannel)
	}
	hlfClientSet, err := operatorv1.NewForConfig(r.Config)
	if err != nil {
		r.setConditionStatus(ctx, fabricMainChannel, hlfv1alpha1.FailedStatus, false, err, false)
		return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricMainChannel)
	}
	ncResponse, err := nc.GenerateNetworkConfig(fabricMainChannel, clientSet, hlfClientSet, "")
	if err != nil {
		r.setConditionStatus(ctx, fabricMainChannel, hlfv1alpha1.FailedStatus, false, errors.Wrapf(err, "failed to generate network config"), false)
		return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricMainChannel)
	}
	log.Infof("Generated network config: %s", ncResponse.NetworkConfig)
	configBackend := config.FromRaw([]byte(ncResponse.NetworkConfig), "yaml")
	sdk, err := fabsdk.New(configBackend)
	if err != nil {
		r.setConditionStatus(ctx, fabricMainChannel, hlfv1alpha1.FailedStatus, false, err, false)
		return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricMainChannel)
	}
	defer sdk.Close()
	firstAdminOrgMSPID := fabricMainChannel.Spec.AdminPeerOrganizations[0].MSPID
	idConfig, ok := fabricMainChannel.Spec.Identities[firstAdminOrgMSPID]
	if !ok {
		r.setConditionStatus(ctx, fabricMainChannel, hlfv1alpha1.FailedStatus, false, fmt.Errorf("identity not found for MSPID %s", firstAdminOrgMSPID), false)
		return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricMainChannel)
	}
	secret, err := clientSet.CoreV1().Secrets(idConfig.SecretNamespace).Get(ctx, idConfig.SecretName, v1.GetOptions{})
	if err != nil {
		r.setConditionStatus(ctx, fabricMainChannel, hlfv1alpha1.FailedStatus, false, err, false)
		return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricMainChannel)
	}
	secretData, ok := secret.Data[idConfig.SecretKey]
	if !ok {
		r.setConditionStatus(ctx, fabricMainChannel, hlfv1alpha1.FailedStatus, false, fmt.Errorf("secret key %s not found", idConfig.SecretKey), false)
		return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricMainChannel)
	}
	id := &identity{}
	err = yaml.Unmarshal(secretData, id)
	if err != nil {
		r.setConditionStatus(ctx, fabricMainChannel, hlfv1alpha1.FailedStatus, false, err, false)
		return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricMainChannel)
	}
	sdkConfig, err := sdk.Config()
	if err != nil {
		r.setConditionStatus(ctx, fabricMainChannel, hlfv1alpha1.FailedStatus, false, err, false)
		return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricMainChannel)
	}
	cryptoConfig := cryptosuite.ConfigFromBackend(sdkConfig)
	cryptoSuite, err := sw.GetSuiteByConfig(cryptoConfig)
	if err != nil {
		r.setConditionStatus(ctx, fabricMainChannel, hlfv1alpha1.FailedStatus, false, err, false)
		return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricMainChannel)
	}
	userStore := mspimpl.NewMemoryUserStore()
	endpointConfig, err := fab.ConfigFromBackend(sdkConfig)
	if err != nil {
		r.setConditionStatus(ctx, fabricMainChannel, hlfv1alpha1.FailedStatus, false, err, false)
		return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricMainChannel)
	}
	identityManager, err := mspimpl.NewIdentityManager(firstAdminOrgMSPID, userStore, cryptoSuite, endpointConfig)
	if err != nil {
		r.setConditionStatus(ctx, fabricMainChannel, hlfv1alpha1.FailedStatus, false, err, false)
		return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricMainChannel)
	}
	signingIdentity, err := identityManager.CreateSigningIdentity(
		msp.WithPrivateKey([]byte(id.Key.Pem)),
		msp.WithCert([]byte(id.Cert.Pem)),
	)
	if err != nil {
		r.setConditionStatus(ctx, fabricMainChannel, hlfv1alpha1.FailedStatus, false, err, false)
		return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricMainChannel)
	}
	sdkContext := sdk.Context(
		fabsdk.WithIdentity(signingIdentity),
		fabsdk.WithOrg(firstAdminOrgMSPID),
	)
	resClient, err := resmgmt.New(sdkContext)
	if err != nil {
		r.setConditionStatus(ctx, fabricMainChannel, hlfv1alpha1.FailedStatus, false, err, false)
		return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricMainChannel)
	}
	resmgmtOptions := []resmgmt.RequestOption{
		resmgmt.WithTimeout(fab2.ResMgmt, 30*time.Second),
	}
	for _, ordOrg := range fabricMainChannel.Spec.OrdererOrganizations {
		for _, endpoint := range ordOrg.OrdererEndpoints {
			resmgmtOptions = append(resmgmtOptions, resmgmt.WithOrdererEndpoint(endpoint))
		}
	}
	var blockBytes []byte

	channelBlock, err := resClient.QueryConfigBlockFromOrderer(fabricMainChannel.Spec.Name, resmgmtOptions...)
	if err == nil {
		log.Infof("Channel %s already exists", fabricMainChannel.Spec.Name)
		blockBytes, err = proto.Marshal(channelBlock)
		if err != nil {
			r.setConditionStatus(ctx, fabricMainChannel, hlfv1alpha1.FailedStatus, false, err, false)
			return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricMainChannel)
		}
	} else {
		log.Infof("Channel %s does not exist, creating it: %v", fabricMainChannel.Spec.Name, err)
		channelConfig, err := r.mapToConfigTX(fabricMainChannel)
		if err != nil {
			r.setConditionStatus(ctx, fabricMainChannel, hlfv1alpha1.FailedStatus, false, err, false)
			return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricMainChannel)
		}
		block, err := configtx.NewApplicationChannelGenesisBlock(channelConfig, fabricMainChannel.Spec.Name)
		if err != nil {
			r.setConditionStatus(ctx, fabricMainChannel, hlfv1alpha1.FailedStatus, false, err, false)
			return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricMainChannel)
		}
		blockBytes, err = proto.Marshal(block)
		if err != nil {
			r.setConditionStatus(ctx, fabricMainChannel, hlfv1alpha1.FailedStatus, false, err, false)
			return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricMainChannel)
		}
	}

	// join orderers
	for _, ordererOrg := range fabricMainChannel.Spec.OrdererOrganizations {
		var tlsCACert string
		if ordererOrg.CAName != "" && ordererOrg.CANamespace != "" {
			certAuth, err := helpers.GetCertAuthByName(
				clientSet,
				hlfClientSet,
				ordererOrg.CAName,
				ordererOrg.CANamespace,
			)
			if err != nil {
				r.setConditionStatus(ctx, fabricMainChannel, hlfv1alpha1.FailedStatus, false, err, false)
				return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricMainChannel)
			}
			tlsCACert = certAuth.Status.TLSCACert

		} else if ordererOrg.TLSCACert != "" && ordererOrg.SignCACert != "" {
			tlsCACert = ordererOrg.TLSCACert
		}
		certPool := x509.NewCertPool()
		ok := certPool.AppendCertsFromPEM([]byte(tlsCACert))
		if !ok {
			r.setConditionStatus(ctx, fabricMainChannel, hlfv1alpha1.FailedStatus, false, fmt.Errorf("couldn't append certs from org %s", ordererOrg.MSPID), false)
			return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricMainChannel)
		}
		idConfig, ok := fabricMainChannel.Spec.Identities[fmt.Sprintf("%s-tls", ordererOrg.MSPID)]
		if !ok {
			log.Infof("Identity for MSPID %s not found, trying with normal identity", fmt.Sprintf("%s-tls", ordererOrg.MSPID))
			// try with normal identity
			idConfig, ok = fabricMainChannel.Spec.Identities[ordererOrg.MSPID]
			if !ok {
				r.setConditionStatus(ctx, fabricMainChannel, hlfv1alpha1.FailedStatus, false, fmt.Errorf("identity not found for MSPID %s", ordererOrg.MSPID), false)
				return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricMainChannel)
			}
		}
		secret, err := clientSet.CoreV1().Secrets(idConfig.SecretNamespace).Get(ctx, idConfig.SecretName, v1.GetOptions{})
		if err != nil {
			r.setConditionStatus(ctx, fabricMainChannel, hlfv1alpha1.FailedStatus, false, err, false)
			return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricMainChannel)
		}
		id := &identity{}
		secretData, ok := secret.Data[idConfig.SecretKey]
		if !ok {
			r.setConditionStatus(ctx, fabricMainChannel, hlfv1alpha1.FailedStatus, false, fmt.Errorf("secret key %s not found", idConfig.SecretKey), false)
			return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricMainChannel)
		}
		err = yaml.Unmarshal(secretData, id)
		if err != nil {
			r.setConditionStatus(ctx, fabricMainChannel, hlfv1alpha1.FailedStatus, false, err, false)
			return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricMainChannel)
		}
		tlsClientCert, err := tls.X509KeyPair(
			[]byte(id.Cert.Pem),
			[]byte(id.Key.Pem),
		)
		for _, cc := range ordererOrg.ExternalOrderersToJoin {
			osnUrl := fmt.Sprintf("https://%s:%d", cc.Host, cc.AdminPort)
			log.Infof("Trying to join orderer %s to channel %s", osnUrl, fabricMainChannel.Spec.Name)
			chResponse, err := osnadmin.Join(osnUrl, blockBytes, certPool, tlsClientCert)
			if err != nil {
				r.setConditionStatus(ctx, fabricMainChannel, hlfv1alpha1.FailedStatus, false, err, false)
				return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricMainChannel)
			}
			defer chResponse.Body.Close()
			if chResponse.StatusCode == 405 {
				log.Infof("Orderer %s already joined to channel %s", osnUrl, fabricMainChannel.Spec.Name)
				continue
			}
			responseData, err := ioutil.ReadAll(chResponse.Body)
			if err != nil {
				r.setConditionStatus(ctx, fabricMainChannel, hlfv1alpha1.FailedStatus, false, err, false)
				return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricMainChannel)
			}
			log.Infof("Orderer %s joined Status code=%d", osnUrl, chResponse.StatusCode)

			if chResponse.StatusCode != 201 {
				r.setConditionStatus(
					ctx,
					fabricMainChannel,
					hlfv1alpha1.FailedStatus,
					false,
					fmt.Errorf(
						"response from orderer %s trying to join to the channel %s: %d, response: %s",
						osnUrl,
						fabricMainChannel.Spec.Name,
						chResponse.StatusCode,
						string(responseData),
					),
					false,
				)
				return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricMainChannel)
			}
			chInfo := &osnadmin.ChannelInfo{}
			err = json.Unmarshal(responseData, chInfo)
			if err != nil {
				r.setConditionStatus(ctx, fabricMainChannel, hlfv1alpha1.FailedStatus, false, err, false)
				return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricMainChannel)
			}
		}

		for _, cc := range ordererOrg.OrderersToJoin {
			ordererNode, err := hlfClientSet.HlfV1alpha1().FabricOrdererNodes(cc.Namespace).Get(ctx, cc.Name, v1.GetOptions{})
			if err != nil {
				r.setConditionStatus(ctx, fabricMainChannel, hlfv1alpha1.FailedStatus, false, err, false)
				return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricMainChannel)
			}
			adminHost, adminPort, err := helpers.GetOrdererAdminHostAndPort(clientSet, ordererNode.Spec, ordererNode.Status)
			if err != nil {
				r.setConditionStatus(ctx, fabricMainChannel, hlfv1alpha1.FailedStatus, false, err, false)
				return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricMainChannel)
			}
			osnUrl := fmt.Sprintf("https://%s:%d", adminHost, adminPort)
			log.Infof("Trying to join orderer %s to channel %s", osnUrl, fabricMainChannel.Spec.Name)
			chResponse, err := osnadmin.Join(osnUrl, blockBytes, certPool, tlsClientCert)
			if err != nil {
				r.setConditionStatus(ctx, fabricMainChannel, hlfv1alpha1.FailedStatus, false, err, false)
				return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricMainChannel)
			}
			defer chResponse.Body.Close()
			if chResponse.StatusCode == 405 {
				log.Infof("Orderer %s already joined to channel %s", osnUrl, fabricMainChannel.Spec.Name)
				continue
			}
			responseData, err := ioutil.ReadAll(chResponse.Body)
			if err != nil {
				r.setConditionStatus(ctx, fabricMainChannel, hlfv1alpha1.FailedStatus, false, err, false)
				return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricMainChannel)
			}
			log.Infof("Orderer %s.%s joined Status code=%d", cc.Name, cc.Namespace, chResponse.StatusCode)
			if chResponse.StatusCode != 201 {
				r.setConditionStatus(
					ctx,
					fabricMainChannel,
					hlfv1alpha1.FailedStatus,
					false,
					fmt.Errorf(
						"response from orderer %s trying to join to the channel %s: %d, response: %s",
						osnUrl,
						fabricMainChannel.Spec.Name,
						chResponse.StatusCode,
						string(responseData),
					),
					false,
				)
				return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricMainChannel)
			}
			chInfo := &osnadmin.ChannelInfo{}
			err = json.Unmarshal(responseData, chInfo)
			if err != nil {
				r.setConditionStatus(ctx, fabricMainChannel, hlfv1alpha1.FailedStatus, false, err, false)
				return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricMainChannel)
			}
		}
	}

	r.Log.Info("Fetching block from orderer")
	var ordererChannelBlock *common.Block
	attemptsLeft := 5
	for {
		ordererChannelBlock, err = resClient.QueryConfigBlockFromOrderer(fabricMainChannel.Spec.Name, resmgmtOptions...)
		if err == nil || attemptsLeft == 0 {
			break
		}
		if err != nil {
			attemptsLeft--
		}
		r.Log.Info(fmt.Sprintf("Failed to get block %v, attempts left %d", err, attemptsLeft))
		time.Sleep(1500 * time.Millisecond)
	}

	if err != nil {
		r.Log.Info(fmt.Sprintf("Failed to get block %v", err))
		r.setConditionStatus(ctx, fabricMainChannel, hlfv1alpha1.FailedStatus, false, errors.Wrapf(err, "failed to get block from channel %s", fabricMainChannel.Spec.Name), false)
		return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricMainChannel)
	}
	r.Log.Info(fmt.Sprintf("Block from channel %s fetched from orderer", fabricMainChannel.Spec.Name))
	cfgBlock, err := resource.ExtractConfigFromBlock(ordererChannelBlock)
	if err != nil {
		r.setConditionStatus(ctx, fabricMainChannel, hlfv1alpha1.FailedStatus, false, errors.Wrapf(err, "failed to extract config from channel block"), false)
		return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricMainChannel)
	}
	currentConfigTx := configtx.New(cfgBlock)
	newConfigTx, err := r.mapToConfigTX(fabricMainChannel)
	if err != nil {
		r.setConditionStatus(ctx, fabricMainChannel, hlfv1alpha1.FailedStatus, false, errors.Wrapf(err, "error mapping channel to configtx channel"), false)
		return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricMainChannel)
	}
	var buf2 bytes.Buffer
	err = protolator.DeepMarshalJSON(&buf2, cfgBlock)
	if err != nil {
		r.setConditionStatus(ctx, fabricMainChannel, hlfv1alpha1.FailedStatus, false, errors.Wrapf(err, "error converting block to JSON"), false)
		return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricMainChannel)
	}
	log.Debug(fmt.Sprintf("Config block main channel: %s", buf2.String()))
	log.Debug(fmt.Sprintf("ConfigTX: %v", newConfigTx))
	err = updateApplicationChannelConfigTx(currentConfigTx, newConfigTx)
	if err != nil {
		r.setConditionStatus(ctx, fabricMainChannel, hlfv1alpha1.FailedStatus, false, errors.Wrapf(err, "failed to update application channel config"), false)
		return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricMainChannel)
	}
	configUpdate, err := resmgmt.CalculateConfigUpdate(fabricMainChannel.Spec.Name, cfgBlock, currentConfigTx.UpdatedConfig())
	if err != nil {
		if !strings.Contains(err.Error(), "no differences detected between original and updated config") {
			r.setConditionStatus(ctx, fabricMainChannel, hlfv1alpha1.FailedStatus, false, errors.Wrapf(err, "error calculating config update"), false)
			return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricMainChannel)
		}
		log.Infof("No differences detected between original and updated config")
	} else {
		channelConfigBytes, err := CreateConfigUpdateEnvelope(fabricMainChannel.Spec.Name, configUpdate)
		if err != nil {
			r.setConditionStatus(ctx, fabricMainChannel, hlfv1alpha1.FailedStatus, false, errors.Wrapf(err, "error creating config update envelope"), false)
			return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricMainChannel)
		}
		var configSignatures []*common.ConfigSignature
		for _, adminPeer := range fabricMainChannel.Spec.AdminPeerOrganizations {
			configUpdateReader := bytes.NewReader(channelConfigBytes)
			idConfig, ok := fabricMainChannel.Spec.Identities[adminPeer.MSPID]
			if !ok {
				r.setConditionStatus(ctx, fabricMainChannel, hlfv1alpha1.FailedStatus, false, fmt.Errorf("identity not found for MSPID %s", adminPeer.MSPID), false)
				return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricMainChannel)
			}
			secret, err := clientSet.CoreV1().Secrets(idConfig.SecretNamespace).Get(ctx, idConfig.SecretName, v1.GetOptions{})
			if err != nil {
				r.setConditionStatus(ctx, fabricMainChannel, hlfv1alpha1.FailedStatus, false, err, false)
				return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricMainChannel)
			}
			secretData, ok := secret.Data[idConfig.SecretKey]
			if !ok {
				r.setConditionStatus(ctx, fabricMainChannel, hlfv1alpha1.FailedStatus, false, fmt.Errorf("secret key %s not found", idConfig.SecretKey), false)
				return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricMainChannel)
			}
			id := &identity{}
			err = yaml.Unmarshal(secretData, id)
			if err != nil {
				r.setConditionStatus(ctx, fabricMainChannel, hlfv1alpha1.FailedStatus, false, err, false)
				return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricMainChannel)
			}
			sdkConfig, err := sdk.Config()
			if err != nil {
				r.setConditionStatus(ctx, fabricMainChannel, hlfv1alpha1.FailedStatus, false, err, false)
				return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricMainChannel)
			}
			cryptoConfig := cryptosuite.ConfigFromBackend(sdkConfig)
			cryptoSuite, err := sw.GetSuiteByConfig(cryptoConfig)
			if err != nil {
				r.setConditionStatus(ctx, fabricMainChannel, hlfv1alpha1.FailedStatus, false, err, false)
				return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricMainChannel)
			}
			userStore := mspimpl.NewMemoryUserStore()
			endpointConfig, err := fab.ConfigFromBackend(sdkConfig)
			if err != nil {
				r.setConditionStatus(ctx, fabricMainChannel, hlfv1alpha1.FailedStatus, false, err, false)
				return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricMainChannel)
			}
			identityManager, err := mspimpl.NewIdentityManager(adminPeer.MSPID, userStore, cryptoSuite, endpointConfig)
			if err != nil {
				r.setConditionStatus(ctx, fabricMainChannel, hlfv1alpha1.FailedStatus, false, err, false)
				return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricMainChannel)
			}
			signingIdentity, err := identityManager.CreateSigningIdentity(
				msp.WithPrivateKey([]byte(id.Key.Pem)),
				msp.WithCert([]byte(id.Cert.Pem)),
			)
			if err != nil {
				r.setConditionStatus(ctx, fabricMainChannel, hlfv1alpha1.FailedStatus, false, err, false)
				return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricMainChannel)
			}

			sdkContext := sdk.Context(
				fabsdk.WithIdentity(signingIdentity),
				fabsdk.WithOrg(adminPeer.MSPID),
			)
			resClient, err := resmgmt.New(sdkContext)
			if err != nil {
				r.setConditionStatus(ctx, fabricMainChannel, hlfv1alpha1.FailedStatus, false, err, false)
				return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricMainChannel)
			}
			signature, err := resClient.CreateConfigSignatureFromReader(signingIdentity, configUpdateReader)
			if err != nil {
				r.setConditionStatus(ctx, fabricMainChannel, hlfv1alpha1.FailedStatus, false, err, false)
				return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricMainChannel)
			}
			configSignatures = append(configSignatures, signature)
		}
		configUpdateReader := bytes.NewReader(channelConfigBytes)
		saveChannelOpts := []resmgmt.RequestOption{
			resmgmt.WithConfigSignatures(configSignatures...),
		}
		saveChannelOpts = append(saveChannelOpts, resmgmtOptions...)
		saveChannelResponse, err := resClient.SaveChannel(
			resmgmt.SaveChannelRequest{
				ChannelID:         fabricMainChannel.Spec.Name,
				ChannelConfig:     configUpdateReader,
				SigningIdentities: []msp.SigningIdentity{},
			},
			saveChannelOpts...,
		)
		if err != nil {
			r.setConditionStatus(ctx, fabricMainChannel, hlfv1alpha1.FailedStatus, false, errors.Wrapf(err, "error saving application configuration"), false)
			return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricMainChannel)
		}
		log.Infof("Application configuration updated with transaction ID: %s", saveChannelResponse.TransactionID)
		currentConfigTx := configtx.New(cfgBlock)
		newConfigTx, err := r.mapToConfigTX(fabricMainChannel)
		if err != nil {
			r.setConditionStatus(ctx, fabricMainChannel, hlfv1alpha1.FailedStatus, false, errors.Wrapf(err, "error mapping channel to configtx channel"), false)
			return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricMainChannel)
		}
		err = updateOrdererChannelConfigTx(currentConfigTx, newConfigTx)
		if err != nil {
			r.setConditionStatus(ctx, fabricMainChannel, hlfv1alpha1.FailedStatus, false, errors.Wrapf(err, "failed to update application channel config"), false)
			return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricMainChannel)
		}
		configUpdate, err := resmgmt.CalculateConfigUpdate(fabricMainChannel.Spec.Name, cfgBlock, currentConfigTx.UpdatedConfig())
		if err != nil {
			if !strings.Contains(err.Error(), "no differences detected between original and updated config") {
				r.setConditionStatus(ctx, fabricMainChannel, hlfv1alpha1.FailedStatus, false, errors.Wrapf(err, "error calculating config update"), false)
				return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricMainChannel)
			}
			log.Infof("No differences detected between original and updated config")
		} else {
			channelConfigBytes, err := CreateConfigUpdateEnvelope(fabricMainChannel.Spec.Name, configUpdate)
			if err != nil {
				r.setConditionStatus(ctx, fabricMainChannel, hlfv1alpha1.FailedStatus, false, errors.Wrapf(err, "error creating config update envelope"), false)
				return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricMainChannel)
			}
			var buf2 bytes.Buffer
			err = protolator.DeepMarshalJSON(&buf2, cfgBlock)
			if err != nil {
				r.setConditionStatus(ctx, fabricMainChannel, hlfv1alpha1.FailedStatus, false, errors.Wrapf(err, "error converting block to JSON"), false)
				return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricMainChannel)
			}
			configSignatures = []*cb.ConfigSignature{}
			for _, adminOrderer := range fabricMainChannel.Spec.AdminOrdererOrganizations {
				configUpdateReader := bytes.NewReader(channelConfigBytes)
				identityName := fmt.Sprintf("%s-sign", adminOrderer.MSPID)
				idConfig, ok := fabricMainChannel.Spec.Identities[identityName]
				if !ok {
					r.setConditionStatus(ctx, fabricMainChannel, hlfv1alpha1.FailedStatus, false, fmt.Errorf("identity not found for MSPID %s", identityName), false)
					return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricMainChannel)
				}
				secret, err := clientSet.CoreV1().Secrets(idConfig.SecretNamespace).Get(ctx, idConfig.SecretName, v1.GetOptions{})
				if err != nil {
					r.setConditionStatus(ctx, fabricMainChannel, hlfv1alpha1.FailedStatus, false, err, false)
					return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricMainChannel)
				}
				secretData, ok := secret.Data[idConfig.SecretKey]
				if !ok {
					r.setConditionStatus(ctx, fabricMainChannel, hlfv1alpha1.FailedStatus, false, fmt.Errorf("secret key %s not found", idConfig.SecretKey), false)
					return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricMainChannel)
				}
				id := &identity{}
				err = yaml.Unmarshal(secretData, id)
				if err != nil {
					r.setConditionStatus(ctx, fabricMainChannel, hlfv1alpha1.FailedStatus, false, err, false)
					return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricMainChannel)
				}
				sdkConfig, err := sdk.Config()
				if err != nil {
					r.setConditionStatus(ctx, fabricMainChannel, hlfv1alpha1.FailedStatus, false, err, false)
					return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricMainChannel)
				}
				cryptoConfig := cryptosuite.ConfigFromBackend(sdkConfig)
				cryptoSuite, err := sw.GetSuiteByConfig(cryptoConfig)
				if err != nil {
					r.setConditionStatus(ctx, fabricMainChannel, hlfv1alpha1.FailedStatus, false, err, false)
					return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricMainChannel)
				}
				userStore := mspimpl.NewMemoryUserStore()
				endpointConfig, err := fab.ConfigFromBackend(sdkConfig)
				if err != nil {
					r.setConditionStatus(ctx, fabricMainChannel, hlfv1alpha1.FailedStatus, false, err, false)
					return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricMainChannel)
				}
				identityManager, err := mspimpl.NewIdentityManager(adminOrderer.MSPID, userStore, cryptoSuite, endpointConfig)
				if err != nil {
					r.setConditionStatus(ctx, fabricMainChannel, hlfv1alpha1.FailedStatus, false, err, false)
					return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricMainChannel)
				}
				signingIdentity, err := identityManager.CreateSigningIdentity(
					msp.WithPrivateKey([]byte(id.Key.Pem)),
					msp.WithCert([]byte(id.Cert.Pem)),
				)
				if err != nil {
					r.setConditionStatus(ctx, fabricMainChannel, hlfv1alpha1.FailedStatus, false, err, false)
					return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricMainChannel)
				}

				sdkContext := sdk.Context(
					fabsdk.WithIdentity(signingIdentity),
					fabsdk.WithOrg(adminOrderer.MSPID),
				)
				resClient, err := resmgmt.New(sdkContext)
				if err != nil {
					r.setConditionStatus(ctx, fabricMainChannel, hlfv1alpha1.FailedStatus, false, err, false)
					return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricMainChannel)
				}
				signature, err := resClient.CreateConfigSignatureFromReader(signingIdentity, configUpdateReader)
				if err != nil {
					r.setConditionStatus(ctx, fabricMainChannel, hlfv1alpha1.FailedStatus, false, err, false)
					return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricMainChannel)
				}
				configSignatures = append(configSignatures, signature)
			}
			configUpdateReader = bytes.NewReader(channelConfigBytes)
			saveChannelOpts = []resmgmt.RequestOption{
				resmgmt.WithConfigSignatures(configSignatures...),
			}
			saveChannelOpts = append(saveChannelOpts, resmgmtOptions...)
			saveChannelResponse, err = resClient.SaveChannel(
				resmgmt.SaveChannelRequest{
					ChannelID:         fabricMainChannel.Spec.Name,
					ChannelConfig:     configUpdateReader,
					SigningIdentities: []msp.SigningIdentity{},
				},
				saveChannelOpts...,
			)
			if err != nil {
				r.setConditionStatus(ctx, fabricMainChannel, hlfv1alpha1.FailedStatus, false, errors.Wrapf(err, "error saving orderer configuration"), false)
				return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricMainChannel)
			}
			log.Infof("Orderer configuration updated with transaction ID: %s", saveChannelResponse.TransactionID)
		}
	}
	r.Log.Info(fmt.Sprintf("fetching block every 1 second waiting for orderers to reconcile %s", fabricMainChannel.Name))
	ordererChannelCh := make(chan *common.Block, 1)
	go func() {
		for {
			ordererChannelBlock, err = resClient.QueryConfigBlockFromOrderer(fabricMainChannel.Spec.Name, resmgmtOptions...)
			if err != nil {
				log.Errorf("error querying orderer channel: %v", err)
				time.Sleep(1 * time.Second)
			} else {
				log.Infof("orderer channel fetched")
				ordererChannelCh <- ordererChannelBlock
				break
			}
		}
	}()
	select {
	case res := <-ordererChannelCh:
		ordererChannelBlock = res
	case <-time.After(12 * time.Second):
		err = errors.New("timeout querying orderer channel")
		r.Log.Error(err, "error querying orderer channel")
		r.setConditionStatus(ctx, fabricMainChannel, hlfv1alpha1.FailedStatus, false, err, false)
		return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricMainChannel)
	}
	cmnConfig, err := resource.ExtractConfigFromBlock(ordererChannelBlock)
	if err != nil {
		r.setConditionStatus(ctx, fabricMainChannel, hlfv1alpha1.FailedStatus, false, errors.Wrapf(err, "error extracting the config from block"), false)
		return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricMainChannel)
	}
	var buf bytes.Buffer
	err = protolator.DeepMarshalJSON(&buf, cmnConfig)
	if err != nil {
		r.setConditionStatus(ctx, fabricMainChannel, hlfv1alpha1.FailedStatus, false, errors.Wrapf(err, "error converting block to JSON"), false)
		return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricMainChannel)
	}
	configMapName := fmt.Sprintf("%s-config", fabricMainChannel.ObjectMeta.Name)
	createConfigMap := false
	configMapNamespace := "default"
	configMap, err := clientSet.CoreV1().ConfigMaps(configMapNamespace).Get(ctx, configMapName, v1.GetOptions{})
	if err != nil {
		if apierrors.IsNotFound(err) {
			reqLogger.Info(fmt.Sprintf("ConfigMap %s not found, creating it", configMapName))
			createConfigMap = true
		} else {
			r.setConditionStatus(ctx, fabricMainChannel, hlfv1alpha1.FailedStatus, false, errors.Wrapf(err, "error getting configmap"), false)
			return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricMainChannel)
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
			r.setConditionStatus(ctx, fabricMainChannel, hlfv1alpha1.FailedStatus, false, errors.Wrapf(err, "error creating config map"), false)
			return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricMainChannel)
		}
	} else {
		configMap.Data["channel.json"] = buf.String()
		_, err = clientSet.CoreV1().ConfigMaps(configMapNamespace).Update(ctx, configMap, v1.UpdateOptions{})
		if err != nil {
			r.setConditionStatus(ctx, fabricMainChannel, hlfv1alpha1.FailedStatus, false, errors.Wrapf(err, "error updating config map"), false)
			return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricMainChannel)
		}
	}
	fabricMainChannel.Status.Status = hlfv1alpha1.RunningStatus
	fabricMainChannel.Status.Message = "Channel setup completed"

	fabricMainChannel.Status.Conditions.SetCondition(status.Condition{
		Type:   status.ConditionType(fabricMainChannel.Status.Status),
		Status: "True",
	})
	if err := r.Status().Update(ctx, fabricMainChannel); err != nil {
		r.setConditionStatus(ctx, fabricMainChannel, hlfv1alpha1.FailedStatus, false, err, false)
		return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricMainChannel)
	}
	r.setConditionStatus(ctx, fabricMainChannel, hlfv1alpha1.RunningStatus, true, nil, false)
	return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricMainChannel)
}

var (
	ErrClientK8s = errors.New("k8sAPIClientError")
)

func (r *FabricMainChannelReconciler) updateCRStatusOrFailReconcile(ctx context.Context, log logr.Logger, p *hlfv1alpha1.FabricMainChannel) (
	reconcile.Result, error) {
	if err := r.Status().Update(ctx, p); err != nil {
		log.Error(err, fmt.Sprintf("%v failed to update the application status", ErrClientK8s))
		return reconcile.Result{}, err
	}
	return reconcile.Result{}, nil
}

func (r *FabricMainChannelReconciler) setConditionStatus(ctx context.Context, p *hlfv1alpha1.FabricMainChannel, conditionType hlfv1alpha1.DeploymentStatus, statusFlag bool, err error, statusUnknown bool) (update bool) {
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

func (r *FabricMainChannelReconciler) SetupWithManager(mgr ctrl.Manager) error {
	managedBy := ctrl.NewControllerManagedBy(mgr)
	return managedBy.
		For(&hlfv1alpha1.FabricMainChannel{}).
		Owns(&corev1.Secret{}).
		Complete(r)
}

func (r *FabricMainChannelReconciler) mapToConfigTX(channel *hlfv1alpha1.FabricMainChannel) (configtx.Channel, error) {
	consenters := []orderer.Consenter{}
	for _, consenter := range channel.Spec.Consenters {
		tlsCert, err := utils.ParseX509Certificate([]byte(consenter.TLSCert))
		if err != nil {
			return configtx.Channel{}, err
		}
		channelConsenter := orderer.Consenter{
			Address: orderer.EtcdAddress{
				Host: consenter.Host,
				Port: consenter.Port,
			},
			ClientTLSCert: tlsCert,
			ServerTLSCert: tlsCert,
		}
		consenters = append(consenters, channelConsenter)
	}
	clientSet, err := utils.GetClientKubeWithConf(r.Config)
	if err != nil {
		return configtx.Channel{}, err
	}
	hlfClientSet, err := operatorv1.NewForConfig(r.Config)
	if err != nil {
		return configtx.Channel{}, err
	}
	ordererOrgs := []configtx.Organization{}
	for _, ordererOrg := range channel.Spec.OrdererOrganizations {
		var tlsCACert *x509.Certificate
		var caCert *x509.Certificate
		if ordererOrg.CAName != "" && ordererOrg.CANamespace != "" {
			certAuth, err := helpers.GetCertAuthByName(
				clientSet,
				hlfClientSet,
				ordererOrg.CAName,
				ordererOrg.CANamespace,
			)
			if err != nil {
				return configtx.Channel{}, err
			}
			tlsCACert, err = utils.ParseX509Certificate([]byte(certAuth.Status.TLSCACert))
			if err != nil {
				return configtx.Channel{}, err
			}
			caCert, err = utils.ParseX509Certificate([]byte(certAuth.Status.CACert))
			if err != nil {
				return configtx.Channel{}, err
			}
		} else if ordererOrg.TLSCACert != "" && ordererOrg.SignCACert != "" {
			tlsCACert, err = utils.ParseX509Certificate([]byte(ordererOrg.TLSCACert))
			if err != nil {
				return configtx.Channel{}, err
			}
			caCert, err = utils.ParseX509Certificate([]byte(ordererOrg.SignCACert))
			if err != nil {
				return configtx.Channel{}, err
			}
		}
		ordererOrgs = append(ordererOrgs, r.mapOrdererOrg(ordererOrg.MSPID, ordererOrg.OrdererEndpoints, caCert, tlsCACert))
	}
	for _, ordererOrg := range channel.Spec.ExternalOrdererOrganizations {
		tlsCACert, err := utils.ParseX509Certificate([]byte(ordererOrg.TLSRootCert))
		if err != nil {
			return configtx.Channel{}, err
		}
		caCert, err := utils.ParseX509Certificate([]byte(ordererOrg.SignRootCert))
		if err != nil {
			return configtx.Channel{}, err
		}
		ordererOrgs = append(ordererOrgs, r.mapOrdererOrg(ordererOrg.MSPID, ordererOrg.OrdererEndpoints, caCert, tlsCACert))
	}
	etcdRaftOptions := orderer.EtcdRaftOptions{
		TickInterval:         "500ms",
		ElectionTick:         10,
		HeartbeatTick:        1,
		MaxInflightBlocks:    5,
		SnapshotIntervalSize: 16 * 1024 * 1024, // 16 MB
	}
	if channel.Spec.ChannelConfig != nil &&
		channel.Spec.ChannelConfig.Orderer != nil &&
		channel.Spec.ChannelConfig.Orderer.EtcdRaft != nil &&
		channel.Spec.ChannelConfig.Orderer.EtcdRaft.Options != nil {
		etcdRaftOptions.TickInterval = channel.Spec.ChannelConfig.Orderer.EtcdRaft.Options.TickInterval
		etcdRaftOptions.ElectionTick = channel.Spec.ChannelConfig.Orderer.EtcdRaft.Options.ElectionTick
		etcdRaftOptions.HeartbeatTick = channel.Spec.ChannelConfig.Orderer.EtcdRaft.Options.HeartbeatTick
		etcdRaftOptions.MaxInflightBlocks = channel.Spec.ChannelConfig.Orderer.EtcdRaft.Options.MaxInflightBlocks
		etcdRaftOptions.SnapshotIntervalSize = channel.Spec.ChannelConfig.Orderer.EtcdRaft.Options.SnapshotIntervalSize
	}
	ordererAdminRule := "MAJORITY Admins"
	if channel.Spec.AdminOrdererOrganizations != nil {
		ordererAdminRule = "OR("
		for idx, adminOrg := range channel.Spec.AdminOrdererOrganizations {
			ordererAdminRule += "'" + adminOrg.MSPID + ".admin'"
			if idx < len(channel.Spec.AdminOrdererOrganizations)-1 {
				ordererAdminRule += ","
			}
		}
		ordererAdminRule += ")"
	}
	adminOrdererPolicies := map[string]configtx.Policy{
		"Readers": {
			Type: "ImplicitMeta",
			Rule: "ANY Readers",
		},
		"Writers": {
			Type: "ImplicitMeta",
			Rule: "ANY Writers",
		},
		"Admins": {
			Type: "Signature",
			Rule: ordererAdminRule,
		},
		"BlockValidation": {
			Type: "ImplicitMeta",
			Rule: "ANY Writers",
		},
	}
	ordConfigtx := configtx.Orderer{
		OrdererType:   "etcdraft",
		Organizations: ordererOrgs,
		SmartBFT: &smartbft.Options{
			RequestBatchMaxCount:      100,
			RequestBatchMaxInterval:   "50ms",
			RequestForwardTimeout:     "2s",
			RequestComplainTimeout:    "20s",
			RequestAutoRemoveTimeout:  "3m0s",
			ViewChangeResendInterval:  "5s",
			ViewChangeTimeout:         "20s",
			LeaderHeartbeatTimeout:    "1m0s",
			CollectTimeout:            "1s",
			RequestBatchMaxBytes:      10485760,
			IncomingMessageBufferSize: 200,
			RequestPoolSize:           100000,
			LeaderHeartbeatCount:      10,
		},
		EtcdRaft: orderer.EtcdRaft{
			Consenters: consenters,
			Options:    etcdRaftOptions,
		},
		Policies:     adminOrdererPolicies,
		Capabilities: []string{"V2_0"},
		BatchSize: orderer.BatchSize{
			MaxMessageCount:   100,
			AbsoluteMaxBytes:  1024 * 1024,
			PreferredMaxBytes: 512 * 1024,
		},
		BatchTimeout: 2 * time.Second,
		State:        "STATE_NORMAL",
	}
	if channel.Spec.ChannelConfig != nil {
		if channel.Spec.ChannelConfig.Orderer != nil {
			if channel.Spec.ChannelConfig.Orderer.BatchTimeout != "" {
				batchTimeout, err := time.ParseDuration(channel.Spec.ChannelConfig.Orderer.BatchTimeout)
				if err != nil {
					return configtx.Channel{}, err
				}
				ordConfigtx.BatchTimeout = batchTimeout
			}
			if channel.Spec.ChannelConfig.Orderer.BatchSize != nil {
				ordConfigtx.BatchSize.MaxMessageCount = uint32(channel.Spec.ChannelConfig.Orderer.BatchSize.MaxMessageCount)
				ordConfigtx.BatchSize.AbsoluteMaxBytes = uint32(channel.Spec.ChannelConfig.Orderer.BatchSize.AbsoluteMaxBytes)
				ordConfigtx.BatchSize.PreferredMaxBytes = uint32(channel.Spec.ChannelConfig.Orderer.BatchSize.PreferredMaxBytes)
			}
		}
	}
	peerOrgs := []configtx.Organization{}
	for _, peerOrg := range channel.Spec.PeerOrganizations {
		certAuth, err := helpers.GetCertAuthByName(
			clientSet,
			hlfClientSet,
			peerOrg.CAName,
			peerOrg.CANamespace,
		)
		if err != nil {
			return configtx.Channel{}, err
		}
		tlsCACert, err := utils.ParseX509Certificate([]byte(certAuth.Status.TLSCACert))
		if err != nil {
			return configtx.Channel{}, err
		}
		caCert, err := utils.ParseX509Certificate([]byte(certAuth.Status.CACert))
		if err != nil {
			return configtx.Channel{}, err
		}
		peerOrgs = append(peerOrgs, r.mapPeerOrg(peerOrg.MSPID, caCert, tlsCACert))
	}
	for _, peerOrg := range channel.Spec.ExternalPeerOrganizations {
		tlsCACert, err := utils.ParseX509Certificate([]byte(peerOrg.TLSRootCert))
		if err != nil {
			return configtx.Channel{}, err
		}
		caCert, err := utils.ParseX509Certificate([]byte(peerOrg.SignRootCert))
		if err != nil {
			return configtx.Channel{}, err
		}
		peerOrgs = append(peerOrgs, r.mapPeerOrg(peerOrg.MSPID, caCert, tlsCACert))
	}
	var adminAppPolicy string
	if len(channel.Spec.AdminPeerOrganizations) == 0 {
		adminAppPolicy = "MAJORITY Admins"
	} else {
		adminAppPolicy = "OR("
		for idx, adminPeerOrg := range channel.Spec.AdminPeerOrganizations {
			adminAppPolicy += "'" + adminPeerOrg.MSPID + ".admin'"
			if idx < len(channel.Spec.AdminPeerOrganizations)-1 {
				adminAppPolicy += ","
			}
		}
		adminAppPolicy += ")"
	}
	policies := map[string]configtx.Policy{
		"Readers": {
			Type: "ImplicitMeta",
			Rule: "ANY Readers",
		},
		"Writers": {
			Type: "ImplicitMeta",
			Rule: "ANY Writers",
		},
		"Admins": {
			Type: "Signature",
			Rule: adminAppPolicy,
		},
		"Endorsement": {
			Type: "ImplicitMeta",
			Rule: "MAJORITY Endorsement",
		},
		"LifecycleEndorsement": {
			Type: "ImplicitMeta",
			Rule: "MAJORITY Endorsement",
		},
	}
	application := configtx.Application{
		Organizations: peerOrgs,
		Capabilities:  []string{"V2_0"},
		Policies:      policies,
		ACLs:          defaultACLs(),
	}
	channelConfig := configtx.Channel{
		Orderer:      ordConfigtx,
		Application:  application,
		Capabilities: []string{"V2_0"},
		Policies: map[string]configtx.Policy{
			"Readers": {
				Type: "ImplicitMeta",
				Rule: "ANY Readers",
			},
			"Writers": {
				Type: "ImplicitMeta",
				Rule: "ANY Writers",
			},
			"Admins": {
				Type: "ImplicitMeta",
				Rule: "MAJORITY Admins",
			},
		},
	}
	return channelConfig, nil
}

func (r *FabricMainChannelReconciler) mapOrdererOrg(mspID string, ordererEndpoints []string, caCert *x509.Certificate, tlsCACert *x509.Certificate) configtx.Organization {
	return configtx.Organization{
		Name: mspID,
		Policies: map[string]configtx.Policy{
			"Admins": {
				Type: "Signature",
				Rule: fmt.Sprintf("OR('%s.admin')", mspID),
			},
			"Readers": {
				Type: "Signature",
				Rule: fmt.Sprintf("OR('%s.member')", mspID),
			},
			"Writers": {
				Type: "Signature",
				Rule: fmt.Sprintf("OR('%s.member')", mspID),
			},
			"Endorsement": {
				Type: "Signature",
				Rule: fmt.Sprintf("OR('%s.member')", mspID),
			},
		},
		MSP: configtx.MSP{
			Name:         mspID,
			RootCerts:    []*x509.Certificate{caCert},
			TLSRootCerts: []*x509.Certificate{tlsCACert},
			NodeOUs: membership.NodeOUs{
				Enable: true,
				ClientOUIdentifier: membership.OUIdentifier{
					Certificate:                  caCert,
					OrganizationalUnitIdentifier: "client",
				},
				PeerOUIdentifier: membership.OUIdentifier{
					Certificate:                  caCert,
					OrganizationalUnitIdentifier: "peer",
				},
				AdminOUIdentifier: membership.OUIdentifier{
					Certificate:                  caCert,
					OrganizationalUnitIdentifier: "admin",
				},
				OrdererOUIdentifier: membership.OUIdentifier{
					Certificate:                  caCert,
					OrganizationalUnitIdentifier: "orderer",
				},
			},
			Admins:                        []*x509.Certificate{},
			IntermediateCerts:             []*x509.Certificate{},
			RevocationList:                []*pkix.CertificateList{},
			OrganizationalUnitIdentifiers: []membership.OUIdentifier{},
			CryptoConfig:                  membership.CryptoConfig{},
			TLSIntermediateCerts:          []*x509.Certificate{},
		},
		AnchorPeers:      []configtx.Address{},
		OrdererEndpoints: ordererEndpoints,
		ModPolicy:        "",
	}
}

func (r *FabricMainChannelReconciler) mapPeerOrg(mspID string, caCert *x509.Certificate, tlsCACert *x509.Certificate) configtx.Organization {
	return configtx.Organization{
		Name: mspID,
		Policies: map[string]configtx.Policy{
			"Admins": {
				Type: "Signature",
				Rule: fmt.Sprintf("OR('%s.admin')", mspID),
			},
			"Readers": {
				Type: "Signature",
				Rule: fmt.Sprintf("OR('%s.member')", mspID),
			},
			"Writers": {
				Type: "Signature",
				Rule: fmt.Sprintf("OR('%s.member')", mspID),
			},
			"Endorsement": {
				Type: "Signature",
				Rule: fmt.Sprintf("OR('%s.member')", mspID),
			},
		},
		MSP: configtx.MSP{
			Name:         mspID,
			RootCerts:    []*x509.Certificate{caCert},
			TLSRootCerts: []*x509.Certificate{tlsCACert},
			NodeOUs: membership.NodeOUs{
				Enable: true,
				ClientOUIdentifier: membership.OUIdentifier{
					Certificate:                  caCert,
					OrganizationalUnitIdentifier: "client",
				},
				PeerOUIdentifier: membership.OUIdentifier{
					Certificate:                  caCert,
					OrganizationalUnitIdentifier: "peer",
				},
				AdminOUIdentifier: membership.OUIdentifier{
					Certificate:                  caCert,
					OrganizationalUnitIdentifier: "admin",
				},
				OrdererOUIdentifier: membership.OUIdentifier{
					Certificate:                  caCert,
					OrganizationalUnitIdentifier: "orderer",
				},
			},
			Admins:                        []*x509.Certificate{},
			IntermediateCerts:             []*x509.Certificate{},
			RevocationList:                []*pkix.CertificateList{},
			OrganizationalUnitIdentifiers: []membership.OUIdentifier{},
			CryptoConfig:                  membership.CryptoConfig{},
			TLSIntermediateCerts:          []*x509.Certificate{},
		},
		AnchorPeers:      []configtx.Address{},
		OrdererEndpoints: []string{},
		ModPolicy:        "",
	}
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
	configUpdateEnvelope := &cb.ConfigUpdateEnvelope{}
	configUpdateEnvelope.ConfigUpdate = configUpdateData
	envelope, err := protoutil.CreateSignedEnvelope(cb.HeaderType_CONFIG_UPDATE, channelID, nil, configUpdateEnvelope, 0, 0)
	if err != nil {
		return nil, err
	}
	envelopeData, err := proto.Marshal(envelope)
	if err != nil {
		return nil, err
	}
	return envelopeData, nil
}

func updateApplicationChannelConfigTx(currentConfigTX configtx.ConfigTx, newConfigTx configtx.Channel) error {
	err := currentConfigTX.Application().SetPolicies(
		newConfigTx.Application.Policies,
	)
	if err != nil {
		return errors.Wrapf(err, "failed to set application")
	}
	app, err := currentConfigTX.Application().Configuration()
	if err != nil {
		return errors.Wrapf(err, "failed to get application configuration")
	}
	log.Infof("Current organizations %v", app.Organizations)
	log.Infof("New organizations %v", newConfigTx.Application.Organizations)
	for _, channelPeerOrg := range app.Organizations {
		deleted := true
		for _, organization := range newConfigTx.Application.Organizations {
			if organization.Name == channelPeerOrg.Name {
				deleted = false
				break
			}
		}
		if deleted {
			log.Infof("Removing organization %s", channelPeerOrg.Name)
			currentConfigTX.Application().RemoveOrganization(channelPeerOrg.Name)
		}
	}
	for _, organization := range newConfigTx.Application.Organizations {
		found := false
		for _, channelPeerOrg := range app.Organizations {
			if channelPeerOrg.Name == organization.Name {
				found = true
				break
			}
		}
		if !found {
			log.Infof("Adding organization %s", organization.Name)
			err = currentConfigTX.Application().SetOrganization(organization)
			if err != nil {
				return errors.Wrapf(err, "failed to set organization %s", organization.Name)
			}
		}
	}
	err = currentConfigTX.Application().SetPolicies(
		newConfigTx.Application.Policies,
	)
	if err != nil {
		return errors.Wrap(err, "failed to set application policies")
	}
	err = currentConfigTX.Application().SetACLs(
		newConfigTx.Application.ACLs,
	)
	if err != nil {
		return errors.Wrapf(err, "failed to set ACLs")
	}
	err = currentConfigTX.Orderer().SetBatchTimeout(newConfigTx.Orderer.BatchTimeout)
	if err != nil {
		return errors.Wrapf(err, "failed to set batch timeout")
	}
	return nil
}
func updateOrdererChannelConfigTx(currentConfigTX configtx.ConfigTx, newConfigTx configtx.Channel) error {
	err := currentConfigTX.Orderer().SetPolicies(
		newConfigTx.Orderer.Policies,
	)
	if err != nil {
		return errors.Wrapf(err, "failed to set application")
	}
	ord, err := currentConfigTX.Orderer().Configuration()
	if err != nil {
		return errors.Wrapf(err, "failed to get application configuration")
	}
	log.Infof("Current orderer organizations %v", ord.Organizations)
	log.Infof("New orderer organizations %v", newConfigTx.Orderer.Organizations)
	for _, channelOrdOrg := range ord.Organizations {
		deleted := true
		for _, organization := range newConfigTx.Orderer.Organizations {
			if organization.Name == channelOrdOrg.Name {
				deleted = false
				break
			}
		}
		if deleted {
			log.Infof("Removing organization %s", channelOrdOrg.Name)
			currentConfigTX.Orderer().RemoveOrganization(channelOrdOrg.Name)
		}
	}
	for _, organization := range newConfigTx.Orderer.Organizations {
		found := false
		for _, channelPeerOrg := range ord.Organizations {
			if channelPeerOrg.Name == organization.Name {
				found = true
				break
			}
		}
		if found {
			ordConfig, err := currentConfigTX.Orderer().Organization(organization.Name).Configuration()
			if err != nil {
				return errors.Wrapf(err, "failed to get orderer organization configuration")
			}
			// remove all previous endpoints
			for _, endpoint := range ordConfig.OrdererEndpoints {
				// extract host and port for endpoint
				host, portStr, err := net.SplitHostPort(endpoint)
				if err != nil {
					return errors.Wrapf(err, "failed to split host and port for endpoint %s", endpoint)
				}
				port, err := strconv.Atoi(portStr)
				if err != nil {
					return errors.Wrapf(err, "failed to convert port %s to int", portStr)
				}
				err = currentConfigTX.Orderer().Organization(organization.Name).RemoveEndpoint(
					configtx.Address{
						Host: host,
						Port: port,
					},
				)
				if err != nil {
					return errors.Wrapf(err, "failed to remove endpoint %s", endpoint)
				}
			}
			// add endpoints
			for _, endpoint := range organization.OrdererEndpoints {
				host, portStr, err := net.SplitHostPort(endpoint)
				if err != nil {
					return errors.Wrapf(err, "failed to split host and port for endpoint %s", endpoint)
				}
				port, err := strconv.Atoi(portStr)
				if err != nil {
					return errors.Wrapf(err, "failed to convert port %s to int", portStr)
				}
				err = currentConfigTX.Orderer().Organization(organization.Name).SetEndpoint(configtx.Address{
					Host: host,
					Port: port,
				})
				if err != nil {
					return errors.Wrapf(err, "failed to add endpoint %s", endpoint)
				}
			}
		} else {
			log.Infof("Adding organization %s", organization.Name)
			err = currentConfigTX.Orderer().SetOrganization(organization)
			if err != nil {
				return errors.Wrapf(err, "failed to set organization %s", organization.Name)
			}

		}
	}
	for _, consenter := range ord.EtcdRaft.Consenters {
		deleted := true
		for _, newConsenter := range newConfigTx.Orderer.EtcdRaft.Consenters {
			if newConsenter.Address.Host == consenter.Address.Host && newConsenter.Address.Port == consenter.Address.Port {
				deleted = false
				break
			}
		}
		if deleted {
			log.Infof("Removing consenter %s:%d", consenter.Address.Host, consenter.Address.Port)
			err = currentConfigTX.Orderer().RemoveConsenter(consenter)
			if err != nil {
				return errors.Wrapf(err, "failed to remove consenter %s:%d", consenter.Address.Host, consenter.Address.Port)
			}
		}
	}
	for _, newConsenter := range newConfigTx.Orderer.EtcdRaft.Consenters {
		found := false
		for _, consenter := range ord.EtcdRaft.Consenters {
			if newConsenter.Address.Host == consenter.Address.Host && newConsenter.Address.Port == consenter.Address.Port {
				found = true
				break
			}
		}
		if !found {
			log.Infof("Adding consenter %s:%d", newConsenter.Address.Host, newConsenter.Address.Port)
			err = currentConfigTX.Orderer().AddConsenter(newConsenter)
			if err != nil {
				return errors.Wrapf(err, "failed to add consenter %s:%d", newConsenter.Address.Host, newConsenter.Address.Port)
			}
		}
	}

	err = currentConfigTX.Orderer().BatchSize().SetMaxMessageCount(
		newConfigTx.Orderer.BatchSize.MaxMessageCount,
	)
	if err != nil {
		return errors.Wrapf(err, "failed to set max message count")
	}
	err = currentConfigTX.Orderer().BatchSize().SetAbsoluteMaxBytes(
		newConfigTx.Orderer.BatchSize.AbsoluteMaxBytes,
	)
	if err != nil {
		return errors.Wrapf(err, "failed to set absolute max bytes")
	}
	err = currentConfigTX.Orderer().BatchSize().SetPreferredMaxBytes(
		newConfigTx.Orderer.BatchSize.PreferredMaxBytes,
	)
	if err != nil {
		return errors.Wrapf(err, "failed to set preferred max bytes")
	}
	err = currentConfigTX.Orderer().SetPolicies(
		newConfigTx.Orderer.Policies,
	)
	if err != nil {
		return errors.Wrap(err, "failed to set application policies")
	}

	return nil
}

func defaultACLs() map[string]string {
	return map[string]string{
		"_lifecycle/CheckCommitReadiness": "/Channel/Application/Writers",

		//  ACL policy for _lifecycle's "CommitChaincodeDefinition" function
		"_lifecycle/CommitChaincodeDefinition": "/Channel/Application/Writers",

		//  ACL policy for _lifecycle's "QueryChaincodeDefinition" function
		"_lifecycle/QueryChaincodeDefinition": "/Channel/Application/Writers",

		//  ACL policy for _lifecycle's "QueryChaincodeDefinitions" function
		"_lifecycle/QueryChaincodeDefinitions": "/Channel/Application/Writers",

		// ---Lifecycle System Chaincode (lscc) function to policy mapping for access control---//

		//  ACL policy for lscc's "getid" function
		"lscc/ChaincodeExists": "/Channel/Application/Readers",

		//  ACL policy for lscc's "getdepspec" function
		"lscc/GetDeploymentSpec": "/Channel/Application/Readers",

		//  ACL policy for lscc's "getccdata" function
		"lscc/GetChaincodeData": "/Channel/Application/Readers",

		//  ACL Policy for lscc's "getchaincodes" function
		"lscc/GetInstantiatedChaincodes": "/Channel/Application/Readers",

		// ---Query System Chaincode (qscc) function to policy mapping for access control---//

		//  ACL policy for qscc's "GetChainInfo" function
		"qscc/GetChainInfo": "/Channel/Application/Readers",

		//  ACL policy for qscc's "GetBlockByNumber" function
		"qscc/GetBlockByNumber": "/Channel/Application/Readers",

		//  ACL policy for qscc's  "GetBlockByHash" function
		"qscc/GetBlockByHash": "/Channel/Application/Readers",

		//  ACL policy for qscc's "GetTransactionByID" function
		"qscc/GetTransactionByID": "/Channel/Application/Readers",

		//  ACL policy for qscc's "GetBlockByTxID" function
		"qscc/GetBlockByTxID": "/Channel/Application/Readers",

		// ---Configuration System Chaincode (cscc) function to policy mapping for access control---//

		//  ACL policy for cscc's "GetConfigBlock" function
		"cscc/GetConfigBlock": "/Channel/Application/Readers",

		//  ACL policy for cscc's "GetChannelConfig" function
		"cscc/GetChannelConfig": "/Channel/Application/Readers",

		// ---Miscellaneous peer function to policy mapping for access control---//

		//  ACL policy for invoking chaincodes on peer
		"peer/Propose": "/Channel/Application/Writers",

		//  ACL policy for chaincode to chaincode invocation
		"peer/ChaincodeToChaincode": "/Channel/Application/Writers",

		// ---Events resource to policy mapping for access control// // // ---//

		//  ACL policy for sending block events
		"event/Block": "/Channel/Application/Readers",

		//  ACL policy for sending filtered block events
		"event/FilteredBlock": "/Channel/Application/Readers",
	}
}
