package ca

import (
	"context"
	"fmt"
	"io"

	"github.com/kfsoftware/hlf-operator/api/hlf.kungfusoftware.es/v1alpha1"
	"github.com/kfsoftware/hlf-operator/kubectl-hlf/cmd/helpers"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Options struct {
	Name           string
	StorageClass   string
	Capacity       string
	NS             string
	Image          string
	Version        string
	EnrollID       string
	EnrollSecret   string
	Output         bool
	IngressGateway string
	IngressPort    int
	Hosts          []string
}

func (o Options) Validate() error {
	return nil
}

type createCmd struct {
	out    io.Writer
	errOut io.Writer
	caOpts Options
}

func (c *createCmd) validate() error {
	return c.caOpts.Validate()
}
func (c *createCmd) run(args []string) error {
	oclient, err := helpers.GetKubeOperatorClient()
	if err != nil {
		return err
	}
	identities := []v1alpha1.FabricCAIdentity{
		{
			Name:        c.caOpts.EnrollID,
			Pass:        c.caOpts.EnrollSecret,
			Type:        "client",
			Affiliation: "",
			Attrs: v1alpha1.FabricCAIdentityAttrs{
				RegistrarRoles: "*",
				DelegateRoles:  "*",
				Attributes:     "*",
				Revoker:        true,
				IntermediateCA: true,
				GenCRL:         true,
				AffiliationMgr: true,
			},
		},
	}
	ingressGateway := c.caOpts.IngressGateway
	ingressPort := c.caOpts.IngressPort
	istio := &v1alpha1.FabricIstio{
		Port:           ingressPort,
		Hosts:          []string{},
		IngressGateway: ingressGateway,
	}
	if len(c.caOpts.Hosts) > 0 {
		istio = &v1alpha1.FabricIstio{
			Port:           ingressPort,
			Hosts:          c.caOpts.Hosts,
			IngressGateway: ingressGateway,
		}
	}

	hosts := []string{
		"localhost",
		c.caOpts.Name,
		fmt.Sprintf("%s.%s", c.caOpts.Name, c.caOpts.NS),
	}
	hosts = append(hosts, c.caOpts.Hosts...)
	fabricCA := &v1alpha1.FabricCA{
		TypeMeta: v1.TypeMeta{
			Kind:       "FabricCA",
			APIVersion: v1alpha1.GroupVersion.String(),
		},
		ObjectMeta: v1.ObjectMeta{
			Name:      c.caOpts.Name,
			Namespace: c.caOpts.NS,
		},
		Spec: v1alpha1.FabricCASpec{
			Database: v1alpha1.FabricCADatabase{
				Type:       "sqlite3",
				Datasource: "fabric-ca-server.db",
			},
			Hosts: hosts,
			Service: v1alpha1.FabricCASpecService{
				ServiceType: "NodePort",
			},
			Image:        c.caOpts.Image,
			Version:      c.caOpts.Version,
			Debug:        false,
			Istio:        istio,
			CLRSizeLimit: 512000,
			TLS: v1alpha1.FabricCATLSConf{
				Subject: v1alpha1.FabricCASubject{
					CN: "ca",
					C:  "California",
					ST: "",
					O:  "Hyperledger",
					L:  "",
					OU: "Fabric",
				},
			},
			CA: v1alpha1.FabricCAItemConf{
				Name: "ca",
				CFG: v1alpha1.FabricCACFG{
					Identities: v1alpha1.FabricCACFGIdentities{
						AllowRemove: true,
					},
					Affiliations: v1alpha1.FabricCACFGAffilitions{
						AllowRemove: true,
					},
				},
				Subject: v1alpha1.FabricCASubject{
					CN: "ca",
					C:  "ES",
					ST: "Alicante",
					O:  "Kung Fu Software",
					L:  "Alicante",
					OU: "Tech",
				},
				CSR: v1alpha1.FabricCACSR{
					CN:    "ca",
					Hosts: []string{"localhost"},
					Names: []v1alpha1.FabricCANames{
						{C: "US", ST: "", O: "Hyperledger", L: "", OU: "North Carolina"},
					},
					CA: v1alpha1.FabricCACSRCA{
						Expiry:     "131400h",
						PathLength: 0,
					},
				},
				CRL: v1alpha1.FabricCACRL{
					Expiry: "24h",
				},
				Registry: v1alpha1.FabricCARegistry{
					MaxEnrollments: -1,
					Identities:     identities,
				},
				Intermediate: v1alpha1.FabricCAIntermediate{
					ParentServer: v1alpha1.FabricCAIntermediateParentServer{
						URL:    "",
						CAName: "",
					},
				},
				BCCSP: v1alpha1.FabricCABCCSP{
					Default: "SW",
					SW: v1alpha1.FabricCABCCSPSW{
						Hash:     "SHA2",
						Security: "256",
					},
				},
			},
			TLSCA: v1alpha1.FabricCAItemConf{
				Name: "tlsca",
				CFG: v1alpha1.FabricCACFG{
					Identities: v1alpha1.FabricCACFGIdentities{
						AllowRemove: true,
					},
					Affiliations: v1alpha1.FabricCACFGAffilitions{
						AllowRemove: true,
					},
				},
				Subject: v1alpha1.FabricCASubject{
					CN: "tlsca",
					C:  "ES",
					ST: "Alicante",
					O:  "Kung Fu Software",
					L:  "Alicante",
					OU: "Tech",
				},
				CSR: v1alpha1.FabricCACSR{
					CN:    "tlsca",
					Hosts: []string{"localhost"},
					Names: []v1alpha1.FabricCANames{
						{C: "US", ST: "", O: "Hyperledger", L: "", OU: "North Carolina"},
					},
					CA: v1alpha1.FabricCACSRCA{
						Expiry:     "131400h",
						PathLength: 0,
					},
				},
				CRL: v1alpha1.FabricCACRL{
					Expiry: "24h",
				},
				Registry: v1alpha1.FabricCARegistry{
					MaxEnrollments: -1,
					Identities:     identities,
				},
				Intermediate: v1alpha1.FabricCAIntermediate{
					ParentServer: v1alpha1.FabricCAIntermediateParentServer{
						URL:    "",
						CAName: "",
					},
				},
				BCCSP: v1alpha1.FabricCABCCSP{
					Default: "SW",
					SW: v1alpha1.FabricCABCCSPSW{
						Hash:     "SHA2",
						Security: "256",
					},
				},
			},
			Cors: v1alpha1.Cors{
				Enabled: false,
				Origins: []string{},
			},
			Resources: corev1.ResourceRequirements{
				Requests: corev1.ResourceList{},
				Limits:   corev1.ResourceList{},
			},
			Storage: v1alpha1.Storage{
				Size:         c.caOpts.Capacity,
				StorageClass: c.caOpts.StorageClass,
				AccessMode:   "ReadWriteOnce",
			},
			ServiceMonitor: nil,
			Metrics: v1alpha1.FabricCAMetrics{
				Provider: "prometheus",
				Statsd: &v1alpha1.FabricCAMetricsStatsd{
					Network:       "udp",
					Address:       "127.0.0.1:8125",
					WriteInterval: "10s",
					Prefix:        "server",
				},
			},
		},
	}
	if c.caOpts.Output {
		ot, err := helpers.MarshallWithoutStatus(fabricCA)
		if err != nil {
			return err
		}
		fmt.Println(string(ot))
	} else {
		ctx := context.Background()
		_, err = oclient.HlfV1alpha1().FabricCAs(c.caOpts.NS).Create(
			ctx,
			fabricCA,
			v1.CreateOptions{},
		)
		if err != nil {
			return err
		}
		log.Infof("Certificate authority %s created on namespace %s", fabricCA.Name, fabricCA.Namespace)
	}

	return nil
}

func newCreateCACmd(out io.Writer, errOut io.Writer) *cobra.Command {
	c := createCmd{out: out, errOut: errOut}
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a Fabric Certificate authority",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := c.validate(); err != nil {
				return err
			}
			return c.run(args)
		},
	}
	f := cmd.Flags()
	f.StringVar(&c.caOpts.Name, "name", "", "Name of the Certificate Authority tenant to create")
	f.StringVar(&c.caOpts.Capacity, "capacity", "", "Total raw capacity of Certificate Authority, e.g. 2Gi")
	f.StringVarP(&c.caOpts.NS, "namespace", "n", helpers.DefaultNamespace, "Namespace scope for this request")
	f.StringVarP(&c.caOpts.StorageClass, "storage-class", "s", helpers.DefaultStorageclass, "Storage class for this Certificate Authority tenant")
	f.StringVarP(&c.caOpts.Version, "version", "v", helpers.DefaultCAVersion, "Version of the Fabric CA")
	f.StringVarP(&c.caOpts.Image, "image", "i", helpers.DefaultCAImage, "Image of the Fabric CA")
	f.StringVarP(&c.caOpts.EnrollID, "enroll-id", "", "enroll", "Enroll ID of the CA")
	f.StringVarP(&c.caOpts.EnrollSecret, "enroll-pw", "", "enrollpw", "Enroll secret of the CA")
	f.BoolVarP(&c.caOpts.Output, "output", "o", false, "Output in yaml")
	f.StringArrayVarP(&c.caOpts.Hosts, "hosts", "", []string{}, "Hosts for Istio")
	f.StringVarP(&c.caOpts.IngressGateway, "istio-ingressgateway", "", "ingressgateway", "Istio ingress gateway name")
	f.IntVarP(&c.caOpts.IngressPort, "istio-port", "", 443, "Istio ingress port")
	return cmd
}
