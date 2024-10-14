package operatorui

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
	Name            string
	NS              string
	Image           string
	Version         string
	Hosts           []string
	Output          bool
	TLSSecretName   string
	APIURL          string
	IngresClassName string
	LogoURL         string
	OIDCAuthority   string
	OIDCClientId    string
	OIDCScope       string
	Replicas        int
}

func (o Options) Validate() error {
	if o.APIURL == "" {
		return fmt.Errorf("--api-url is required")
	}
	if o.Replicas < 1 {
		return fmt.Errorf("--replicas must be greater than 0")
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
		ClassName: "",
		Annotations: map[string]string{
			"kubernetes.io/ingress.class": c.uiOpts.IngresClassName,
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
			LogoURL: c.uiOpts.LogoURL,
			Auth: &v1alpha1.FabricOperatorUIAuth{
				OIDCAuthority: c.uiOpts.OIDCAuthority,
				OIDCClientId:  c.uiOpts.OIDCClientId,
				OIDCScope:     c.uiOpts.OIDCScope,
			},
			APIURL:           c.uiOpts.APIURL,
			Image:            c.uiOpts.Image,
			Tag:              c.uiOpts.Version,
			ImagePullPolicy:  "Always",
			Tolerations:      []corev1.Toleration{},
			Replicas:         c.uiOpts.Replicas,
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
	f.StringVarP(&c.uiOpts.NS, "namespace", "n", helpers.DefaultNamespace, "Namespace scope for this request")
	f.StringVarP(&c.uiOpts.Image, "image", "", helpers.DefaultOperationsOperatorUIImage, "Image of the Operator UI")
	f.StringVarP(&c.uiOpts.Version, "version", "", helpers.DefaultOperationsOperatorUIVersion, "Version of the Operator UI")
	f.StringVarP(&c.uiOpts.IngresClassName, "ingress-class-name", "", "istio", "Ingress class name")
	f.StringVarP(&c.uiOpts.TLSSecretName, "tls-secret-name", "", "", "TLS Secret for the Operator UI")
	f.StringVarP(&c.uiOpts.APIURL, "api-url", "", "", "API URL for the Operator UI")
	f.StringArrayVarP(&c.uiOpts.Hosts, "hosts", "", []string{}, "External hosts")
	f.BoolVarP(&c.uiOpts.Output, "output", "o", false, "Output in yaml")
	f.IntVarP(&c.uiOpts.Replicas, "replicas", "", 1, "Number of replicas of the Operator UI")
	f.StringVarP(&c.uiOpts.LogoURL, "logo-url", "", "", "Logo URL for the Operator UI")
	f.StringVarP(&c.uiOpts.OIDCAuthority, "oidc-authority", "", "", "OIDC Authority for the Operator UI")
	f.StringVarP(&c.uiOpts.OIDCClientId, "oidc-client-id", "", "", "OIDC Client ID for the Operator UI")
	f.StringVarP(&c.uiOpts.OIDCScope, "oidc-scope", "", "", "OIDC Scope for the Operator UI")
	return cmd
}
