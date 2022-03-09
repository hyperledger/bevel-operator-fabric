package ordnode

import corev1 "k8s.io/api/core/v1"

type fabricOrdChart struct {
	Istio                       Istio               `json:"istio"`
	AdminIstio                  Istio               `json:"adminIstio"`
	Replicas                    int                 `json:"replicas"`
	Genesis                     string              `json:"genesis"`
	ChannelParticipationEnabled bool                `json:"channelParticipationEnabled"`
	BootstrapMethod             string              `json:"bootstrapMethod"`
	Admin                       admin               `json:"admin"`
	Cacert                      string              `json:"cacert"`
	Tlsrootcert                 string              `json:"tlsrootcert"`
	AdminCert                   string              `json:"adminCert"`
	Cert                        string              `json:"cert"`
	Key                         string              `json:"key"`
	TLS                         tls                 `json:"tls"`
	Tolerations                 []corev1.Toleration `json:"tolerations,omitempty"`
	Resources                   Resources           `json:"resources,omitempty"`
	FullnameOverride            string              `json:"fullnameOverride"`
	HostAliases                 []HostAlias         `json:"hostAliases"`
	Service                     service             `json:"service"`
	Image                       image               `json:"image"`
	Persistence                 persistence         `json:"persistence"`
	Ord                         ord                 `json:"ord"`
	Clientcerts                 clientcerts         `json:"clientcerts"`
	Hosts                       []string            `json:"hosts"`
	Logging                     Logging             `json:"logging"`
	ServiceMonitor              ServiceMonitor      `json:"serviceMonitor"`
	EnvVars                     []corev1.EnvVar     `json:"envVars"`
}
type Resources struct {
	Limits   Limits   `json:"limits"`
	Requests Requests `json:"requests"`
}
type Limits struct {
	CPU    string `json:"cpu"`
	Memory string `json:"memory"`
}
type Requests struct {
	CPU    string `json:"cpu"`
	Memory string `json:"memory"`
}
type ServiceMonitor struct {
	Enabled           bool              `json:"enabled"`
	Labels            map[string]string `json:"labels"`
	Interval          string            `json:"interval"`
	ScrapeTimeout     string            `json:"scrapeTimeout"`
	Scheme            string            `json:"scheme"`
	Relabelings       []interface{}     `json:"relabelings"`
	TargetLabels      []interface{}     `json:"targetLabels"`
	MetricRelabelings []interface{}     `json:"metricRelabelings"`
	SampleLimit       int               `json:"sampleLimit"`
}

type Logging struct {
	Spec string `json:"spec"`
}

type admin struct {
	Cert          string `json:"cert"`
	Key           string `json:"key"`
	RootCAs       string `json:"rootCAs"`
	ClientRootCAs string `json:"clientRootCAs"`
}
type tls struct {
	Cert string `json:"cert"`
	Key  string `json:"key"`
}
type HostAlias struct {
	IP        string   `json:"ip"`
	Hostnames []string `json:"hostnames"`
}
type service struct {
	Type               string `json:"type"`
	Port               int    `json:"port"`
	NodePort           int    `json:"nodePort"`
	PortOperations     int    `json:"portOperations"`
	NodePortOperations int    `json:"nodePortOperations"`
}
type image struct {
	Repository string `json:"repository"`
	Tag        string `json:"tag"`
	PullPolicy string `json:"pullPolicy"`
}
type annotations struct {
}
type persistence struct {
	Enabled      bool        `json:"enabled"`
	Annotations  annotations `json:"annotations"`
	StorageClass string      `json:"storageClass"`
	AccessMode   string      `json:"accessMode"`
	Size         string      `json:"size"`
}
type ordServer struct {
	Enabled bool `json:"enabled"`
}
type ordClient struct {
	Enabled bool `json:"enabled"`
}
type tlsConfiguration struct {
	Server ordServer `json:"server"`
	Client ordClient `json:"client"`
}
type ord struct {
	Type  string           `json:"type"`
	MspID string           `json:"mspID"`
	TLS   tlsConfiguration `json:"tls"`
}
type clientcerts struct {
	CertPem string `json:"cert.pem,omitempty"`
}

type Istio struct {
	Port           int      `json:"port"`
	Hosts          []string `json:"hosts"`
	IngressGateway string   `json:"ingressGateway"`
}
