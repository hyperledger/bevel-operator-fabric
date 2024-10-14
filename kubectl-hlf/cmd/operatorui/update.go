package operatorui

import (
	"context"
	"fmt"
	"github.com/kfsoftware/hlf-operator/kubectl-hlf/cmd/helpers"
	"github.com/kfsoftware/hlf-operator/pkg/apis/hlf.kungfusoftware.es/v1alpha1"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"io"
	"k8s.io/api/networking/v1beta1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type updateCmd struct {
	out    io.Writer
	errOut io.Writer
	uiOpts Options
}

func (c *updateCmd) validate() error {
	return c.uiOpts.Validate()
}
func (c *updateCmd) run() error {
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
	ctx := context.Background()
	fabricOperatorUI, err := oclient.HlfV1alpha1().FabricOperatorUIs(c.uiOpts.NS).Get(ctx, c.uiOpts.Name, v1.GetOptions{})
	if err != nil {
		return errors.Wrapf(err, "failed to get Fabric Operator UI %s", c.uiOpts.Name)
	}
	fabricOperatorUI.Spec.Image = c.uiOpts.Image
	fabricOperatorUI.Spec.Tag = c.uiOpts.Version
	fabricOperatorUI.Spec.Auth = &v1alpha1.FabricOperatorUIAuth{
		OIDCAuthority: c.uiOpts.OIDCAuthority,
		OIDCClientId:  c.uiOpts.OIDCClientId,
		OIDCScope:     c.uiOpts.OIDCScope,
	}
	fabricOperatorUI.Spec.LogoURL = c.uiOpts.LogoURL
	fabricOperatorUI.Spec.Ingress = ingress
	fabricOperatorUI.Spec.APIURL = c.uiOpts.APIURL
	fabricOperatorUI.Spec.Replicas = c.uiOpts.Replicas
	if c.uiOpts.Output {
		ot, err := helpers.MarshallWithoutStatus(&fabricOperatorUI)
		if err != nil {
			return err
		}
		fmt.Println(string(ot))
	} else {
		ctx := context.Background()
		_, err = oclient.HlfV1alpha1().FabricOperatorUIs(c.uiOpts.NS).Update(
			ctx,
			fabricOperatorUI,
			v1.UpdateOptions{},
		)
		if err != nil {
			return err
		}
		log.Infof("Operator UI %s updated on namespace %s", fabricOperatorUI.Name, fabricOperatorUI.Namespace)
	}
	return nil
}
func newUpdateOperatorUICmd(out io.Writer, errOut io.Writer) *cobra.Command {
	c := updateCmd{out: out, errOut: errOut}
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update a Operator UI",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := c.validate(); err != nil {
				return err
			}
			return c.run()
		},
	}
	f := cmd.Flags()
	f.StringVar(&c.uiOpts.Name, "name", "", "Name of the Operator UI to update")
	f.StringVarP(&c.uiOpts.NS, "namespace", "n", helpers.DefaultNamespace, "Namespace scope for this request")
	f.StringVarP(&c.uiOpts.Image, "image", "", helpers.DefaultOperationsOperatorUIImage, "Image of the Operator UI")
	f.StringVarP(&c.uiOpts.Version, "version", "", helpers.DefaultOperationsOperatorUIVersion, "Version of the Operator UI")
	f.StringVarP(&c.uiOpts.IngresClassName, "ingress-class-name", "", "istio", "Ingress class name")
	f.StringVarP(&c.uiOpts.TLSSecretName, "tls-secret-name", "", "", "TLS Secret for the Operator UI")
	f.StringVarP(&c.uiOpts.APIURL, "api-url", "", "", "API URL for the Operator UI")
	f.IntVarP(&c.uiOpts.Replicas, "replicas", "", 1, "Number of replicas of the Operator UI")
	f.StringArrayVarP(&c.uiOpts.Hosts, "hosts", "", []string{}, "External hosts")
	f.BoolVarP(&c.uiOpts.Output, "output", "o", false, "Output in yaml")
	f.StringVarP(&c.uiOpts.LogoURL, "logo-url", "", "", "Logo URL for the Operator UI")
	f.StringVarP(&c.uiOpts.OIDCAuthority, "oidc-authority", "", "", "OIDC Authority for the Operator UI")
	f.StringVarP(&c.uiOpts.OIDCClientId, "oidc-client-id", "", "", "OIDC Client ID for the Operator UI")
	f.StringVarP(&c.uiOpts.OIDCScope, "oidc-scope", "", "", "OIDC Scope for the Operator UI")
	return cmd
}
