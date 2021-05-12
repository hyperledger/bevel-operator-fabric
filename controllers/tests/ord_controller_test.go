package tests

import (
	"context"
	log "github.com/sirupsen/logrus"

	hlfv1alpha1 "github.com/kfsoftware/hlf-operator/api/hlf.kungfusoftware.es/v1alpha1"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	// +kubebuilder:scaffold:imports
)

var _ = Describe("Fabric Orderer Controller", func() {
	FabricNamespace := ""
	BeforeEach(func() {
		FabricNamespace = "hlf-operator-" + getRandomChannelID()
		testNamespace := &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: FabricNamespace,
			},
		}
		log.Infof("Creating namespace %s", FabricNamespace)
		Expect(K8sClient.Create(context.Background(), testNamespace)).Should(Succeed())
	})
	Specify("create a new Fabric Orderer with channel participation", func() {
		releaseNameOrdCA := "org1-ca"
		releaseNameOrd := "org1-orderer"
		By("create a fabric ca")
		ordererCA := randomFabricCA(releaseNameOrdCA, FabricNamespace)
		Expect(ordererCA).ToNot(BeNil())
		By("create a fabric orderer")
		ordererMSPID := "OrdererMSP"
		ordParams := createOrdererParams{
			MSPID: ordererMSPID,
		}
		createOrdererNode(
			releaseNameOrd,
			FabricNamespace,
			ordParams,
			ordererCA,
		)
		orderer := &hlfv1alpha1.FabricOrdererNode{}
		ordererKey := types.NamespacedName{
			Namespace: FabricNamespace,
			Name:      releaseNameOrd,
		}
		Eventually(
			func() bool {
				err := K8sClient.Get(context.Background(), ordererKey, orderer)
				if err != nil {
					return false
				}
				ctrl.Log.WithName("test").Info("after update", "orderer", orderer)
				return orderer.Status.Status == hlfv1alpha1.RunningStatus
			},
			peerTimeoutSecs,
			defInterval,
		).Should(BeTrue(), "peer status should have been updated")
	})

})