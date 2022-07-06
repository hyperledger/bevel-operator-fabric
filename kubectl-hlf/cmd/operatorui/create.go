package operatorui

import (
	"context"
	"fmt"
	"github.com/kfsoftware/hlf-operator/api/hlf.kungfusoftware.es/v1alpha1"
	"github.com/kfsoftware/hlf-operator/kubectl-hlf/cmd/helpers"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"io"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/api/networking/v1beta1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Options struct {
	Name           string
	StorageClass   string
	Capacity       string
	NS             string
	Image          string
	Version        string
	IngressGateway string
	IngressPort    int
	Hosts          []string
	Output         bool
	TLSSecretName  string
	APIURL         string
}

func (o Options) Validate() error {
	if o.APIURL == "" {
		return fmt.Errorf("--api-url is required")
	}
	return nil
}

type createCmd struct {
	out    io.Writer
	errOut io.Writer
	uiOpts Options
}

func (c *createCmd) validate() error {
	return c.uiOpts.Validate()
}
func (c *createCmd) run() error {
	oclient, err := helpers.GetKubeOperatorClient()
	if err != nil {
		return err
	}
	hosts := []v1alpha1.IngressHost{}
	for _, host := range c.uiOpts.Hosts {
		hosts = append(hosts, v1alpha1.IngressHost{
			Paths: []v1alpha1.IngressPath{
				{
					Path:     "/",
					PathType: "Prefix",
				},
			},
			Host: host,
		})
	}
	ingress := v1alpha1.Ingress{
		Enabled:   true,
		ClassName: "istio",
		Annotations: map[string]string{
			"kubernetes.io/ingress.class": "istio",
		},
		TLS:   []v1beta1.IngressTLS{},
		Hosts: hosts,
	}
	if c.uiOpts.TLSSecretName != "" {
		ingress.TLS = []v1beta1.IngressTLS{
			{
				Hosts:      c.uiOpts.Hosts,
				SecretName: c.uiOpts.TLSSecretName,
			},
		}
	}
	fabricOperatorUI := &v1alpha1.FabricOperatorUI{
		TypeMeta: v1.TypeMeta{
			Kind:       "FabricOperatorUI",
			APIVersion: v1alpha1.GroupVersion.String(),
		},
		ObjectMeta: v1.ObjectMeta{
			Name:      c.uiOpts.Name,
			Namespace: c.uiOpts.NS,
		},
		Spec: v1alpha1.FabricOperatorUISpec{
			Resources: &corev1.ResourceRequirements{
				Limits:   nil,
				Requests: nil,
			},
			APIURL:           c.uiOpts.APIURL,
			Image:            c.uiOpts.Image,
			Tag:              c.uiOpts.Version,
			ImagePullPolicy:  "Always",
			Tolerations:      []corev1.Toleration{},
			Replicas:         1,
			Env:              []corev1.EnvVar{},
			ImagePullSecrets: []corev1.LocalObjectReference{},
			Affinity:         &corev1.Affinity{},
			Ingress:          ingress,
		},
		Status: v1alpha1.FabricOperatorUIStatus{},
	}
	if c.uiOpts.Output {
		ot, err := helpers.MarshallWithoutStatus(&fabricOperatorUI)
		if err != nil {
			return err
		}
		fmt.Println(string(ot))
	} else {
		ctx := context.Background()
		_, err = oclient.HlfV1alpha1().FabricOperatorUIs(c.uiOpts.NS).Create(
			ctx,
			fabricOperatorUI,
			v1.CreateOptions{},
		)
		if err != nil {
			return err
		}
		log.Infof("Operator UI %s created on namespace %s", fabricOperatorUI.Name, fabricOperatorUI.Namespace)
	}
	return nil
}
func newCreateOperatorUICmd(out io.Writer, errOut io.Writer) *cobra.Command {
	c := createCmd{out: out, errOut: errOut}
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a Operator UI",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := c.validate(); err != nil {
				return err
			}
			return c.run()
		},
	}
	f := cmd.Flags()
	f.StringVar(&c.uiOpts.Name, "name", "", "Name of the Operator UI to create")
	f.StringVar(&c.uiOpts.Capacity, "capacity", "1Gi", "Total raw capacity of Operator UI in this zone, e.g. 16Ti")
	f.StringVarP(&c.uiOpts.NS, "namespace", "n", helpers.DefaultNamespace, "Namespace scope for this request")
	f.StringVarP(&c.uiOpts.StorageClass, "storage-class", "s", helpers.DefaultStorageclass, "Storage class for this Operator UI")
	f.StringVarP(&c.uiOpts.Image, "image", "", helpers.DefaultOperationsOperatorUIImage, "Image of the Operator UI")
	f.StringVarP(&c.uiOpts.Version, "version", "", helpers.DefaultOperationsOperatorUIVersion, "Version of the Operator UI")
	f.StringVarP(&c.uiOpts.IngressGateway, "istio-ingressgateway", "", "ingressgateway", "Istio ingress gateway name")
	f.StringVarP(&c.uiOpts.TLSSecretName, "tls-secret-name", "", "", "TLS Secret for the Operator UI")
	f.StringVarP(&c.uiOpts.APIURL, "api-url", "", "", "API URL for the Operator UI")
	f.IntVarP(&c.uiOpts.IngressPort, "istio-port", "", 443, "Istio ingress port")
	f.StringArrayVarP(&c.uiOpts.Hosts, "hosts", "", []string{}, "External hosts")
	f.BoolVarP(&c.uiOpts.Output, "output", "o", false, "Output in yaml")
	return cmd
}
