package ordservice

import (
	"context"
	"github.com/go-logr/logr"
	hlfv1alpha1 "github.com/kfsoftware/hlf-operator/pkg/apis/hlf.kungfusoftware.es/v1alpha1"
	"github.com/kfsoftware/hlf-operator/pkg/status"
	log "github.com/sirupsen/logrus"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// FabricOrderingServiceReconciler reconciles a FabricOrderingService object
type FabricOrderingServiceReconciler struct {
	client.Client
	ChartPath string
	Log       logr.Logger
	Scheme    *runtime.Scheme
	Config    *rest.Config
}

const ordererFinalizer = "finalizer.orderer.hlf.kungfusoftware.es"

func (r *FabricOrderingServiceReconciler) finalizeOrderer(reqLogger logr.Logger, m *hlfv1alpha1.FabricOrderingService) error {
	ns := m.Namespace
	if ns == "" {
		ns = "default"
	}

	releaseName := m.Name
	reqLogger.Info("Successfully finalized orderer")

	log.Debugf("Release %s deleted", releaseName)
	return nil
}

func (r *FabricOrderingServiceReconciler) addFinalizer(reqLogger logr.Logger, m *hlfv1alpha1.FabricOrderingService) error {
	reqLogger.Info("Adding Finalizer for the Orderer")
	controllerutil.AddFinalizer(m, ordererFinalizer)

	// Update CR
	err := r.Update(context.TODO(), m)
	if err != nil {
		reqLogger.Error(err, "Failed to update Orderer with finalizer")
		return err
	}
	return nil
}

// +kubebuilder:rbac:groups=hlf.kungfusoftware.es,resources=fabricorderingservices,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=hlf.kungfusoftware.es,resources=fabricorderingservices/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=hlf.kungfusoftware.es,resources=fabricorderingservices/finalizers,verbs=get;update;patch
func (r *FabricOrderingServiceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	reqLogger := r.Log.WithValues("hlf", req.NamespacedName)
	reqLogger.Info("Reconciling FabricOrderingService")
	fabricOrderer := &hlfv1alpha1.FabricOrderingService{}
	fabricOrderer.Status.Status = hlfv1alpha1.PendingStatus
	fabricOrderer.Status.Conditions.SetCondition(status.Condition{
		Type:   "NOT_SUPPORTED",
		Status: "True",
	})
	if err := r.Status().Update(ctx, fabricOrderer); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *FabricOrderingServiceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&hlfv1alpha1.FabricOrderingService{}).
		Owns(&appsv1.Deployment{}).
		Complete(r)
}
