package ordnode

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/kfsoftware/hlf-operator/api/hlf.kungfusoftware.es/v1alpha1"
	"github.com/kfsoftware/hlf-operator/controllers/utils"
	"github.com/kfsoftware/hlf-operator/kubectl-hlf/cmd/helpers"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"io"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type OrdererOptions struct {
	Name           string
	StorageClass   string
	Capacity       string
	NS             string
	Image          string
	Version        string
	MspID          string
	EnrollID       string
	EnrollPW       string
	CAName         string
	Hosts          []string
	Output         bool
	IngressGateway string
	IngressPort    int
	AdminHosts     []string
}

func (o OrdererOptions) Validate() error {
	return nil
}

type createCmd struct {
	out         io.Writer
	errOut      io.Writer
	ordererOpts OrdererOptions
}

func (c *createCmd) validate() error {
	return c.ordererOpts.Validate()
}
func (c *createCmd) run(args []string) error {
	oclient, err := helpers.GetKubeOperatorClient()
	if err != nil {
		return err
	}
	clientSet, err := helpers.GetKubeClient()
	if err != nil {
		return err
	}
	certAuth, err := helpers.GetCertAuthByFullName(clientSet, oclient, c.ordererOpts.CAName)
	if err != nil {
		return err
	}
	k8sIP, err := utils.GetPublicIPKubernetes(clientSet)
	if err != nil {
		return err
	}

	csrHosts := []string{
		"127.0.0.1",
		"localhost",
	}
	csrHosts = append(csrHosts, k8sIP)
	csrHosts = append(csrHosts, c.ordererOpts.Name)
	csrHosts = append(csrHosts, fmt.Sprintf("%s.%s", c.ordererOpts.Name, c.ordererOpts.NS))
	ingressGateway := c.ordererOpts.IngressGateway
	ingressPort := c.ordererOpts.IngressPort
	istio := &v1alpha1.FabricIstio{
		Port:           ingressPort,
		Hosts:          []string{},
		IngressGateway: ingressGateway,
	}
	if len(c.ordererOpts.Hosts) > 0 {
		istio = &v1alpha1.FabricIstio{
			Port:           ingressPort,
			Hosts:          c.ordererOpts.Hosts,
			IngressGateway: ingressGateway,
		}
		csrHosts = append(csrHosts, c.ordererOpts.Hosts...)
	}
	adminIstio := &v1alpha1.FabricIstio{
		Port:           ingressPort,
		Hosts:          []string{},
		IngressGateway: ingressGateway,
	}
	if len(c.ordererOpts.AdminHosts) > 0 {
		adminIstio = &v1alpha1.FabricIstio{
			Port:           ingressPort,
			Hosts:          c.ordererOpts.AdminHosts,
			IngressGateway: ingressGateway,
		}
	}
	fabricOrderer := &v1alpha1.FabricOrdererNode{
		TypeMeta: v1.TypeMeta{
			Kind:       "FabricOrdererNode",
			APIVersion: v1alpha1.GroupVersion.String(),
		},
		ObjectMeta: v1.ObjectMeta{
			Name:      c.ordererOpts.Name,
			Namespace: c.ordererOpts.NS,
		},
		Spec: v1alpha1.FabricOrdererNodeSpec{
			ServiceMonitor:              nil,
			HostAliases:                 []corev1.HostAlias{},
			Resources:                   corev1.ResourceRequirements{},
			Replicas:                    1,
			Image:                       c.ordererOpts.Image,
			Tag:                         c.ordererOpts.Version,
			PullPolicy:                  corev1.PullAlways,
			MspID:                       c.ordererOpts.MspID,
			Genesis:                     "",
			BootstrapMethod:             v1alpha1.BootstrapMethodNone,
			ChannelParticipationEnabled: true,
			Storage: v1alpha1.Storage{
				Size:         c.ordererOpts.Capacity,
				StorageClass: c.ordererOpts.StorageClass,
				AccessMode:   "ReadWriteOnce",
			},
			Service: v1alpha1.OrdererNodeService{
				Type: "NodePort",
			},
			Secret: &v1alpha1.Secret{
				Enrollment: v1alpha1.Enrollment{
					Component: v1alpha1.Component{
						Cahost: k8sIP,
						Caname: certAuth.Spec.CA.Name,
						Caport: certAuth.Status.NodePort,
						Catls: v1alpha1.Catls{
							Cacert: base64.StdEncoding.EncodeToString([]byte(certAuth.Status.TlsCert)),
						},
						Enrollid:     c.ordererOpts.EnrollID,
						Enrollsecret: c.ordererOpts.EnrollPW,
					},
					TLS: v1alpha1.TLS{
						Cahost: k8sIP,
						Caname: certAuth.Spec.TLSCA.Name,
						Caport: certAuth.Status.NodePort,
						Catls: v1alpha1.Catls{
							Cacert: base64.StdEncoding.EncodeToString([]byte(certAuth.Status.TlsCert)),
						},
						Csr: v1alpha1.Csr{
							Hosts: csrHosts,
							CN:    "",
						},
						Enrollid:     c.ordererOpts.EnrollID,
						Enrollsecret: c.ordererOpts.EnrollPW,
					},
				},
			},
			Istio:      istio,
			AdminIstio: adminIstio,
		},
	}
	if c.ordererOpts.Output {
		ot, err := helpers.MarshallWithoutStatus(&fabricOrderer)
		if err != nil {
			return err
		}
		fmt.Println(string(ot))
	} else {
		ctx := context.Background()
		ordService, err := oclient.HlfV1alpha1().FabricOrdererNodes(c.ordererOpts.NS).Create(
			ctx,
			fabricOrderer,
			v1.CreateOptions{},
		)
		if err != nil {
			return err
		}
		log.Infof("Ordering service %s created on namespace %s", ordService.Name, ordService.Namespace)
	}
	return nil
}
func newCreateOrdererNodeCmd(out io.Writer, errOut io.Writer) *cobra.Command {
	c := createCmd{out: out, errOut: errOut}
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a Fabric Ordering Service Node(OSN)",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := c.validate(); err != nil {
				return err
			}
			return c.run(args)
		},
	}
	f := cmd.Flags()
	f.StringVar(&c.ordererOpts.Name, "name", "", "Name of the Fabric Orderer to create")
	f.StringVar(&c.ordererOpts.CAName, "ca-name", "", "CA name to enroll this user")
	f.StringVar(&c.ordererOpts.EnrollID, "enroll-id", "", "Enroll ID of the CA")
	f.StringVar(&c.ordererOpts.EnrollPW, "enroll-pw", "", "Enroll secret of the CA")
	f.StringVar(&c.ordererOpts.Capacity, "capacity", "5Gi", "Total raw capacity of Fabric Orderer in this zone, e.g. 16Ti")
	f.StringVarP(&c.ordererOpts.NS, "namespace", "n", helpers.DefaultNamespace, "Namespace scope for this request")
	f.StringVarP(&c.ordererOpts.StorageClass, "storage-class", "s", helpers.DefaultStorageclass, "Storage class for this Fabric Orderer")
	f.StringVarP(&c.ordererOpts.Image, "image", "", helpers.DefaultOrdererImage, "Version of the Fabric Orderer")
	f.StringVarP(&c.ordererOpts.Version, "version", "", helpers.DefaultOrdererVersion, "Version of the Fabric Orderer")
	f.StringVarP(&c.ordererOpts.IngressGateway, "istio-ingressgateway", "", "ingressgateway", "Istio ingress gateway name")
	f.IntVarP(&c.ordererOpts.IngressPort, "istio-port", "", 443, "Istio ingress port")
	f.StringVarP(&c.ordererOpts.MspID, "mspid", "", "", "MSP ID of the organization")
	f.StringArrayVarP(&c.ordererOpts.Hosts, "hosts", "", []string{}, "Hosts")
	f.StringArrayVarP(&c.ordererOpts.AdminHosts, "admin-hosts", "", []string{}, "Hosts for the admin API(introduced in v2.3)")
	f.BoolVarP(&c.ordererOpts.Output, "output", "o", false, "Output in yaml")
	return cmd
}
