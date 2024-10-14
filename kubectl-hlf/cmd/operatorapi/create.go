package operatorapi

import (
	"context"
	"fmt"
	"github.com/kfsoftware/hlf-operator/kubectl-hlf/cmd/helpers"
	"github.com/kfsoftware/hlf-operator/pkg/apis/hlf.kungfusoftware.es/v1alpha1"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"io"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/api/networking/v1beta1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Options struct {
	Name             string
	NS               string
	Image            string
	Version          string
	Hosts            []string
	Output           bool
	TLSSecretName    string
	IngressClassName string
	HLFSecretName    string
	MSPID            string
	User             string
	HLFKey           string
	OIDCIssuer       string
	OIDCJWKS         string
	Replicas         int
}

func (o Options) Validate() error {
	if o.Replicas < 1 {
		return fmt.Errorf("replicas must be greater than 0")
	}
	return nil
}

type createCmd struct {
	out     io.Writer
	errOut  io.Writer
	apiOpts Options
}

func (c *createCmd) validate() error {
	return c.apiOpts.Validate()
}
func (c *createCmd) run() error {
	oclient, err := helpers.GetKubeOperatorClient()
	if err != nil {
		return err
	}
	hosts := []v1alpha1.IngressHost{}
	for _, host := range c.apiOpts.Hosts {
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
		ClassName: "",
		Annotations: map[string]string{
			"kubernetes.io/ingress.class": c.apiOpts.IngressClassName,
		},
		TLS:   []v1beta1.IngressTLS{},
		Hosts: hosts,
	}
	if c.apiOpts.TLSSecretName != "" {
		ingress.TLS = []v1beta1.IngressTLS{
			{
				Hosts:      c.apiOpts.Hosts,
				SecretName: c.apiOpts.TLSSecretName,
			},
		}
	}
	fabricAPI := &v1alpha1.FabricOperatorAPI{
		TypeMeta: v1.TypeMeta{
			Kind:       "FabricOperatorAPI",
			APIVersion: v1alpha1.GroupVersion.String(),
		},
		ObjectMeta: v1.ObjectMeta{
			Name:      c.apiOpts.Name,
			Namespace: c.apiOpts.NS,
		},
		Spec: v1alpha1.FabricOperatorAPISpec{
			PodLabels:       map[string]string{},
			Image:           c.apiOpts.Image,
			Tag:             c.apiOpts.Version,
			ImagePullPolicy: "Always",
			Auth: &v1alpha1.FabricOperatorAPIAuth{
				OIDCJWKS:   c.apiOpts.OIDCJWKS,
				OIDCIssuer: c.apiOpts.OIDCIssuer,
			},
			Istio:    v1alpha1.FabricIstio{},
			Ingress:  ingress,
			Replicas: c.apiOpts.Replicas,
			HLFConfig: v1alpha1.FabricOperatorAPIHLFConfig{
				MSPID: c.apiOpts.MSPID,
				User:  c.apiOpts.User,
				NetworkConfig: v1alpha1.FabricOperatorAPINetworkConfig{
					SecretName: c.apiOpts.HLFSecretName,
					Key:        c.apiOpts.HLFKey,
				},
			},
			Tolerations:      []corev1.Toleration{},
			ImagePullSecrets: []corev1.LocalObjectReference{},
			Env:              []corev1.EnvVar{},
			Affinity:         &corev1.Affinity{},
			Resources: &corev1.ResourceRequirements{
				Limits:   nil,
				Requests: nil,
			},
		},
		Status: v1alpha1.FabricOperatorAPIStatus{},
	}
	if c.apiOpts.Output {
		ot, err := helpers.MarshallWithoutStatus(&fabricAPI)
		if err != nil {
			return err
		}
		fmt.Println(string(ot))
	} else {
		ctx := context.Background()
		_, err = oclient.HlfV1alpha1().FabricOperatorAPIs(c.apiOpts.NS).Create(
			ctx,
			fabricAPI,
			v1.CreateOptions{},
		)
		if err != nil {
			return err
		}
		log.Infof("Operator API %s created on namespace %s", fabricAPI.Name, fabricAPI.Namespace)
	}
	return nil
}
func newCreateOperatorAPICmd(out io.Writer, errOut io.Writer) *cobra.Command {
	c := createCmd{out: out, errOut: errOut}
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a Operator API",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := c.validate(); err != nil {
				return err
			}
			return c.run()
		},
	}
	f := cmd.Flags()
	f.StringVar(&c.apiOpts.Name, "name", "", "Name of the Operator API to create")
	f.StringVarP(&c.apiOpts.NS, "namespace", "n", helpers.DefaultNamespace, "Namespace scope for this request")
	f.StringVarP(&c.apiOpts.Image, "image", "", helpers.DefaultOperationsOperatorAPIImage, "Image of the Operator API")
	f.StringVarP(&c.apiOpts.Version, "version", "", helpers.DefaultOperationsOperatorAPIVersion, "Version of the Operator API")
	f.StringVarP(&c.apiOpts.TLSSecretName, "tls-secret-name", "", "", "TLS Secret for the Operator API")
	f.StringVarP(&c.apiOpts.IngressClassName, "ingress-class-name", "", "istio", "Ingress class name")
	f.StringVarP(&c.apiOpts.MSPID, "hlf-mspid", "", "", "HLF Network Config MSPID")
	f.StringVarP(&c.apiOpts.User, "hlf-user", "", "", "HLF Network Config User")
	f.StringVarP(&c.apiOpts.HLFSecretName, "hlf-secret", "", "", "HLF Network Config Secret name")
	f.StringVarP(&c.apiOpts.HLFKey, "hlf-secret-key", "", "", "HLF Network Config Secret key")
	f.StringVarP(&c.apiOpts.OIDCJWKS, "oidc-jwks", "", "", "OIDC JWKS URL")
	f.StringVarP(&c.apiOpts.OIDCIssuer, "oidc-issuer", "", "", "OIDC Issuer URL")
	f.IntVarP(&c.apiOpts.Replicas, "replicas", "", 1, "Number of replicas of the Operator UI")
	f.StringArrayVarP(&c.apiOpts.Hosts, "hosts", "", []string{}, "External hosts")
	f.BoolVarP(&c.apiOpts.Output, "output", "o", false, "Output in yaml")
	return cmd
}
