package operatorapi

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/api/networking/v1beta1"
)

type Image struct {
	Repository string            `json:"repository"`
	PullPolicy corev1.PullPolicy `json:"pullPolicy"`
	Tag        string            `json:"tag"`
}

type ServiceAccount struct {
	Create      bool              `json:"create"`
	Annotations map[string]string `json:"annotations"`
	Name        string            `json:"name"`
}

type Service struct {
	Type string `json:"type"`
	Port int    `json:"port"`
}

type Autoscaling struct {
	Enabled                        bool `json:"enabled"`
	MinReplicas                    int  `json:"minReplicas"`
	MaxReplicas                    int  `json:"maxReplicas"`
	TargetCPUUtilizationPercentage int  `json:"targetCPUUtilizationPercentage"`
}

type HLFConfig struct {
	MspID         string           `json:"mspID"`
	User          string           `json:"user"`
	NetworkConfig HLFNetworkConfig `json:"networkConfig"`
}

type HLFNetworkConfig struct {
	SecretName string `json:"secretName"`
	Key        string `json:"key"`
}

type Ingress struct {
	Enabled     bool                 `json:"enabled"`
	ClassName   string               `json:"className"`
	Annotations map[string]string    `json:"annotations"`
	TLS         []v1beta1.IngressTLS `json:"tls"`
	Hosts       []IngressHost        `json:"hosts"`
}

type IngressHost struct {
	Host  string        `json:"host"`
	Paths []IngressPath `json:"paths"`
}

type IngressPath struct {
	Path     string `json:"path"`
	PathType string `json:"pathType"`
}

type HLFOperatorAPIChart struct {
	PodLabels        map[string]string             `json:"podLabels"`
	ReplicaCount     int                           `json:"replicaCount"`
	LogoURL          string                        `json:"logoUrl"`
	Image            Image                         `json:"image"`
	Hlf              HLFConfig                     `json:"hlf,omitempty"`
	ImagePullSecrets []corev1.LocalObjectReference `json:"imagePullSecrets"`
	ServiceAccount   ServiceAccount                `json:"serviceAccount"`
	PodAnnotations   map[string]string             `json:"podAnnotations"`
	Service          Service                       `json:"service"`
	Ingress          Ingress                       `json:"ingress"`
	Resources        *corev1.ResourceRequirements  `json:"resources"`
	Autoscaling      Autoscaling                   `json:"autoscaling"`
	Tolerations      []corev1.Toleration           `json:"tolerations,omitempty"`
	Affinity         *corev1.Affinity              `json:"affinity"`
	Auth             Auth                          `json:"auth"`
}

type Auth struct {
	OIDCJWKS      string `json:"oidcJWKS"`
	OIDCIssuer    string `json:"oidcIssuer"`
	OIDCAuthority string `json:"oidcAuthority"`
	OIDCClientId  string `json:"oidcClientId"`
	OIDCScope     string `json:"oidcScope"`
}
