package ordnode

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/kfsoftware/hlf-operator/api/hlf.kungfusoftware.es/v1alpha1"
	"github.com/kfsoftware/hlf-operator/controllers/utils"
	"github.com/kfsoftware/hlf-operator/kubectl-hlf/cmd/helpers"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"io"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
	"strconv"
	"strings"
)

type OrdererOptions struct {
	Name              string
	StorageClass      string
	Capacity          string
	NS                string
	Image             string
	Version           string
	MspID             string
	EnrollID          string
	EnrollPW          string
	CAName            string
	SystemChannelName string
	NumOrderers       int
	Hosts             []string
	Output            bool
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
	c.ordererOpts.Image = helpers.DefaultOrdererImage
	return c.ordererOpts.Validate()
}
func (c *createCmd) run(args []string) error {
	oclient, err := helpers.GetKubeOperatorClient()
	if err != nil {
		return err
	}
	certAuth, err := helpers.GetCertAuthByFullName(oclient, c.ordererOpts.CAName)
	if err != nil {
		return err
	}
	clientSet, err := helpers.GetKubeClient()
	if err != nil {
		return err
	}
	k8sIP, err := utils.GetPublicIPKubernetes(clientSet)
	if err != nil {
		return err
	}
	includeHosts := len(c.ordererOpts.Hosts) > 0
	if includeHosts && len(c.ordererOpts.Hosts) != c.ordererOpts.NumOrderers {
		return errors.Errorf("there are %d orderers but %d hosts", c.ordererOpts.NumOrderers, len(c.ordererOpts.Hosts))
	}
	var nodes []v1alpha1.OrdererNode
	for i := 0; i < c.ordererOpts.NumOrderers; i++ {
		hosts := []string{}
		host := ""
		port := 0
		if includeHosts {
			ordHost := c.ordererOpts.Hosts[i]
			chunks := strings.Split(ordHost, ":")
			host = chunks[0]
			if len(chunks) < 2 {
				return errors.Errorf("Host %s doesn't have port", ordHost)
			}
			port, err = strconv.Atoi(chunks[1])
			if err != nil {
				return err
			}
			hosts = append(hosts, host)
		}
		nodes = append(nodes, v1alpha1.OrdererNode{
			ID:   fmt.Sprintf("orderer%d", i),
			Host: host,
			Port: port,
			Enrollment: v1alpha1.OrdererNodeEnrollment{
				TLS: v1alpha1.OrdererNodeEnrollmentTLS{
					Csr: v1alpha1.Csr{
						Hosts: hosts,
					},
				},
			},
		})
	}
	if len(nodes) == 0 {
		return errors.Errorf("Orderers are empty")
	}

	fabricOrderer := &v1alpha1.FabricOrdererNode{
		TypeMeta: v1.TypeMeta{
			Kind:       "FabricOrderingService",
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
							Hosts: []string{
								"127.0.0.1",
								"localhost",
								k8sIP,
							},
							CN: "",
						},
						Enrollid:     c.ordererOpts.EnrollID,
						Enrollsecret: c.ordererOpts.EnrollPW,
					},
				},
			},
			Istio:      nil,
			AdminIstio: nil,
		},
	}
	if c.ordererOpts.Output {
		ot, err := yaml.Marshal(&fabricOrderer)
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
	f.StringVar(&c.ordererOpts.Name, "name", "", "name of the Fabric Orderer to create")
	f.StringVar(&c.ordererOpts.CAName, "ca-name", "", "name of the Fabric Orderer")
	f.StringVar(&c.ordererOpts.EnrollID, "enroll-id", "", "user to enroll orderer certificates")
	f.StringVar(&c.ordererOpts.EnrollPW, "enroll-pw", "", "password to enroll orderer certificates")
	f.StringVar(&c.ordererOpts.Capacity, "capacity", "5Gi", "total raw capacity of Fabric Orderer in this zone, e.g. 16Ti")
	f.StringVarP(&c.ordererOpts.NS, "namespace", "n", helpers.DefaultNamespace, "namespace scope for this request")
	f.StringVarP(&c.ordererOpts.StorageClass, "storage-class", "s", helpers.DefaultStorageclass, "storage class for this Fabric Orderer")
	f.StringVarP(&c.ordererOpts.Image, "image", "", helpers.DefaultOrdererImage, "version of the Fabric Orderer")
	f.StringVarP(&c.ordererOpts.Version, "version", "", helpers.DefaultOrdererVersion, "version of the Fabric Orderer")
	f.StringVarP(&c.ordererOpts.MspID, "mspid", "", "", "MSP ID of the organization")
	f.StringVarP(&c.ordererOpts.SystemChannelName, "system-channel", "", "system-channel", "System channel name")
	f.IntVarP(&c.ordererOpts.NumOrderers, "num-orderers", "", 3, "Orderers nodes to create")
	f.StringArrayVarP(&c.ordererOpts.Hosts, "hosts", "", []string{}, "Hosts")
	f.BoolVarP(&c.ordererOpts.Output, "output", "o", false, "output in yaml")
	return cmd
}
