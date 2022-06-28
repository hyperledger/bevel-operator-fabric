package console

import (
	"fmt"
	"github.com/kfsoftware/hlf-operator/api/hlf.kungfusoftware.es/v1alpha1"
	"github.com/kfsoftware/hlf-operator/kubectl-hlf/cmd/helpers"
	"github.com/spf13/cobra"
	"io"
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
	MspID          string
	IngressGateway string
	IngressPort    int
	Hosts          []string
	Output         bool
}

func (o Options) Validate() error {
	return nil
}

type createCmd struct {
	out         io.Writer
	errOut      io.Writer
	consoleOpts Options
}

func (c *createCmd) validate() error {
	return c.consoleOpts.Validate()
}
func (c *createCmd) run() error {
	//oclient, err := helpers.GetKubeOperatorClient()
	//if err != nil {
	//	return err
	//}
	fabricConsole := &v1alpha1.FabricOperationsConsole{
		TypeMeta: v1.TypeMeta{
			Kind:       "FabricOperationsConsole",
			APIVersion: v1alpha1.GroupVersion.String(),
		},
		ObjectMeta: v1.ObjectMeta{
			Name:      c.consoleOpts.Name,
			Namespace: c.consoleOpts.NS,
		},
		Spec: v1alpha1.FabricOperationsConsoleSpec{
			Auth: v1alpha1.FabricOperationsConsoleAuth{
				Scheme:   "",
				Username: "",
				Password: "",
			},
			Resources: corev1.ResourceRequirements{
				Limits:   nil,
				Requests: nil,
			},
			Image:           "ghcr.io/hyperledger-labs/fabric-console",
			Tag:             "latest",
			ImagePullPolicy: "Always",
			Tolerations:     []corev1.Toleration{},
			Replicas:        1,
			CouchDB: v1alpha1.FabricOperationsConsoleCouchDB{
				Image:    "",
				Tag:      "",
				Username: "",
				Password: "",
				Storage: v1alpha1.Storage{
					Size:         "",
					StorageClass: "",
					AccessMode:   "",
				},
				Resources: &corev1.ResourceRequirements{
					Limits:   nil,
					Requests: nil,
				},
				ImagePullSecrets: nil,
				Affinity: &corev1.Affinity{
					NodeAffinity: &corev1.NodeAffinity{
						RequiredDuringSchedulingIgnoredDuringExecution: &corev1.NodeSelector{
							NodeSelectorTerms: nil,
						},
						PreferredDuringSchedulingIgnoredDuringExecution: nil,
					},
					PodAffinity: &corev1.PodAffinity{
						RequiredDuringSchedulingIgnoredDuringExecution:  nil,
						PreferredDuringSchedulingIgnoredDuringExecution: nil,
					},
					PodAntiAffinity: &corev1.PodAntiAffinity{
						RequiredDuringSchedulingIgnoredDuringExecution:  nil,
						PreferredDuringSchedulingIgnoredDuringExecution: nil,
					},
				},
				Tolerations:     nil,
				ImagePullPolicy: "",
			},
			Env:              nil,
			ImagePullSecrets: nil,
			Affinity:         nil,
			Port:             0,
			Config:           "",
			Ingress: v1alpha1.Ingress{
				Enabled:     false,
				ClassName:   "",
				Annotations: nil,
				TLS:         nil,
				Hosts:       nil,
			},
			HostURL: "",
		},
		Status: v1alpha1.FabricOperationsConsoleStatus{},
	}
	if c.consoleOpts.Output {
		ot, err := helpers.MarshallWithoutStatus(&fabricConsole)
		if err != nil {
			return err
		}
		fmt.Println(string(ot))
	} else {
		//ctx := context.Background()
		//_, err = oclient.HlfV1alpha1().FabricConsoles(c.consoleOpts.NS).Create(
		//	ctx,
		//	fabricConsole,
		//	v1.CreateOptions{},
		//)
		//if err != nil {
		//	return err
		//}
		//log.Infof("Console %s created on namespace %s", fabricConsole.Name, fabricConsole.Namespace)
	}
	return nil
}
func newCreateConsoleCmd(out io.Writer, errOut io.Writer) *cobra.Command {
	c := createCmd{out: out, errOut: errOut}
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a Fabric Console",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := c.validate(); err != nil {
				return err
			}
			return c.run()
		},
	}
	f := cmd.Flags()
	f.StringVar(&c.consoleOpts.Name, "name", "", "Name of the Fabric Console to create")
	f.StringVar(&c.consoleOpts.Capacity, "capacity", "5Gi", "Total raw capacity of Fabric Console in this zone, e.g. 16Ti")
	f.StringVarP(&c.consoleOpts.NS, "namespace", "n", helpers.DefaultNamespace, "Namespace scope for this request")
	f.StringVarP(&c.consoleOpts.StorageClass, "storage-class", "s", helpers.DefaultStorageclass, "Storage class for this Fabric Console")
	f.StringVarP(&c.consoleOpts.Image, "image", "", helpers.DefaultOperationsConsoleImage, "Version of the Fabric Console")
	f.StringVarP(&c.consoleOpts.Version, "version", "", helpers.DefaultOperationsConsoleVersion, "Version of the Fabric Console")
	f.StringVarP(&c.consoleOpts.MspID, "mspid", "", "", "MSP ID of the organization")
	f.StringVarP(&c.consoleOpts.IngressGateway, "istio-ingressgateway", "", "ingressgateway", "Istio ingress gateway name")
	f.IntVarP(&c.consoleOpts.IngressPort, "istio-port", "", 443, "Istio ingress port")
	f.StringArrayVarP(&c.consoleOpts.Hosts, "hosts", "", []string{}, "External hosts")
	f.BoolVarP(&c.consoleOpts.Output, "output", "o", false, "Output in yaml")
	return cmd
}
