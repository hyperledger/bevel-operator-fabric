package console

import (
	"context"
	"fmt"
	"github.com/kfsoftware/hlf-operator/kubectl-hlf/cmd/helpers"
	"github.com/kfsoftware/hlf-operator/pkg/apis/hlf.kungfusoftware.es/v1alpha1"
	"github.com/sethvargo/go-password/password"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"io"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/api/networking/v1beta1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Options struct {
	Name                 string
	StorageClass         string
	Capacity             string
	NS                   string
	Image                string
	Version              string
	IngressGateway       string
	IngressPort          int
	Hosts                []string
	Output               bool
	InitialAdminPassword string
	InitialAdmin         string
	HostURL              string
	TLSSecretName        string
	ImagePullSecrets     []string
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
	oclient, err := helpers.GetKubeOperatorClient()
	if err != nil {
		return err
	}
	couchDBPassword, err := password.Generate(64, 10, 10, false, false)
	if err != nil {
		return err
	}
	hosts := []v1alpha1.IngressHost{}
	for _, host := range c.consoleOpts.Hosts {
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
	if c.consoleOpts.TLSSecretName != "" {
		ingress.TLS = []v1beta1.IngressTLS{
			{
				Hosts:      c.consoleOpts.Hosts,
				SecretName: c.consoleOpts.TLSSecretName,
			},
		}
	}
	var imagePullSecrets []corev1.LocalObjectReference
	if len(c.consoleOpts.ImagePullSecrets) > 0 {
		for _, v := range c.consoleOpts.ImagePullSecrets {
			imagePullSecrets = append(imagePullSecrets, corev1.LocalObjectReference{
				Name: v,
			})
		}
	}

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
				Scheme:   "couchdb",
				Username: c.consoleOpts.InitialAdmin,
				Password: c.consoleOpts.InitialAdminPassword,
			},
			Resources: corev1.ResourceRequirements{
				Limits:   nil,
				Requests: nil,
			},
			Image:            c.consoleOpts.Image,
			ImagePullSecrets: imagePullSecrets,
			Tag:              c.consoleOpts.Version,
			ImagePullPolicy:  "IfNotPresent",
			Tolerations:      []corev1.Toleration{},
			Replicas:         1,
			CouchDB: v1alpha1.FabricOperationsConsoleCouchDB{
				Image:    "couchdb",
				Tag:      "3.1.1",
				Username: "couchdb",
				Password: couchDBPassword,
				Storage: v1alpha1.Storage{
					Size:         c.consoleOpts.Capacity,
					StorageClass: c.consoleOpts.StorageClass,
					AccessMode:   "ReadWriteOnce",
				},
				Resources: &corev1.ResourceRequirements{
					Limits:   nil,
					Requests: nil,
				},
				ImagePullSecrets: []corev1.LocalObjectReference{},
				Affinity:         &corev1.Affinity{},
				Tolerations:      []corev1.Toleration{},
				ImagePullPolicy:  "IfNotPresent",
			},
			Env:      []corev1.EnvVar{},
			Affinity: &corev1.Affinity{},
			Port:     3000,
			Config:   "",
			Ingress:  ingress,
			HostURL:  c.consoleOpts.HostURL,
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
		ctx := context.Background()
		_, err = oclient.HlfV1alpha1().FabricOperationsConsoles(c.consoleOpts.NS).Create(
			ctx,
			fabricConsole,
			v1.CreateOptions{},
		)
		if err != nil {
			return err
		}
		log.Infof("Console %s created on namespace %s", fabricConsole.Name, fabricConsole.Namespace)
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
	f.StringVar(&c.consoleOpts.Capacity, "capacity", "1Gi", "Total raw capacity of Fabric Console in this zone, e.g. 16Ti")
	f.StringVarP(&c.consoleOpts.NS, "namespace", "n", helpers.DefaultNamespace, "Namespace scope for this request")
	f.StringVarP(&c.consoleOpts.StorageClass, "storage-class", "s", helpers.DefaultStorageclass, "Storage class for this Fabric Console")
	f.StringVarP(&c.consoleOpts.Image, "image", "", helpers.DefaultOperationsConsoleImage, "Image of the Fabric Console")
	f.StringVarP(&c.consoleOpts.Version, "version", "", helpers.DefaultOperationsConsoleVersion, "Version of the Fabric Console")
	f.StringVarP(&c.consoleOpts.TLSSecretName, "tls-secret-name", "", "", "TLS Secret name for serving the console in HTTPS")
	f.StringVarP(&c.consoleOpts.InitialAdmin, "admin-user", "", "", "User name of the console")
	f.StringVarP(&c.consoleOpts.InitialAdminPassword, "admin-pwd", "", "", "Admin password")
	f.StringVarP(&c.consoleOpts.IngressGateway, "istio-ingressgateway", "", "ingressgateway", "Istio ingress gateway name")
	f.StringVarP(&c.consoleOpts.HostURL, "host-url", "", "", "External host URL for the console")
	f.IntVarP(&c.consoleOpts.IngressPort, "istio-port", "", 443, "Istio ingress port")
	f.StringArrayVarP(&c.consoleOpts.Hosts, "hosts", "", []string{}, "External hosts")
	f.BoolVarP(&c.consoleOpts.Output, "output", "o", false, "Output in yaml")
	f.StringArrayVarP(&c.consoleOpts.ImagePullSecrets, "image-pull-secrets", "", []string{}, "Image Pull Secrets for the Peer Image")
	return cmd
}
