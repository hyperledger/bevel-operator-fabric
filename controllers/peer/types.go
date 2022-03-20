package peer

import corev1 "k8s.io/api/core/v1"

type RBAC struct {
	Ns string `json:"ns"`
}
type CouchDB struct {
	External   CouchDBExternal `json:"external"`
	Image      string          `json:"image"`
	Tag        string          `json:"tag"`
	PullPolicy string          `json:"pullPolicy"`
}
type CouchDBExternal struct {
	Enabled bool   `json:"enabled"`
	Host    string `json:"host"`
	Port    int    `json:"port"`
}
type FSServer struct {
	Image      string `json:"image"`
	Tag        string `json:"tag"`
	PullPolicy string `json:"pullPolicy"`
}

type FabricPeerChart struct {
	FSServer                 FSServer            `json:"fsServer"`
	Istio                    Istio               `json:"istio"`
	Replicas                 int                 `json:"replicas"`
	ExternalChaincodeBuilder bool                `json:"externalChaincodeBuilder"`
	CouchdbUsername          string              `json:"couchdbUsername"`
	CouchdbPassword          string              `json:"couchdbPassword"`
	Image                    Image               `json:"image"`
	CouchDB                  CouchDB             `json:"couchdb"`
	Rbac                     RBAC                `json:"rbac"`
	DockerSocketPath         string              `json:"dockerSocketPath"`
	Peer                     Peer                `json:"peer"`
	Cert                     string              `json:"cert"`
	Key                      string              `json:"key"`
	Hosts                    []string            `json:"hosts"`
	TLS                      TLS                 `json:"tls"`
	OPSTLS                   TLS                 `json:"opsTLS"`
	Cacert                   string              `json:"cacert"`
	IntCacert                string              `json:"intCAcert"`
	Tlsrootcert              string              `json:"tlsrootcert"`
	Resources                PeerResources       `json:"resources,omitempty"`
	NodeSelector             NodeSelector        `json:"nodeSelector,omitempty"`
	Tolerations              []corev1.Toleration `json:"tolerations,omitempty"`
	Affinity                 Affinity            `json:"affinity,omitempty"`
	ExternalHost             string              `json:"externalHost"`
	FullnameOverride         string              `json:"fullnameOverride"`
	CouchDBExporter          CouchDBExporter     `json:"couchdbExporter"`
	HostAliases              []HostAlias         `json:"hostAliases"`
	Service                  Service             `json:"service"`
	Persistence              PeerPersistence     `json:"persistence"`
	Logging                  Logging             `json:"logging"`
	ExternalBuilders         []ExternalBuilder   `json:"externalBuilders"`
	ServiceMonitor           ServiceMonitor      `json:"serviceMonitor"`
	EnvVars                  []corev1.EnvVar     `json:"envVars"`
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

type ExternalBuilder struct {
	Name                 string   `json:"name"`
	Path                 string   `json:"path"`
	PropagateEnvironment []string `json:"propagateEnvironment"`
}

type Istio struct {
	Port           int      `json:"port"`
	Hosts          []string `json:"hosts"`
	IngressGateway string   `json:"ingressGateway"`
}
type PeerResources struct {
	Peer            Resources  `json:"peer"`
	CouchDB         Resources  `json:"couchdb"`
	Chaincode       Resources  `json:"chaincode"`
	CouchDBExporter *Resources `json:"couchdbExporter,omitempty"`
}
type CouchDBExporter struct {
	Enabled    bool   `json:"enabled"`
	Image      string `json:"image"`
	Tag        string `json:"tag"`
	PullPolicy string `json:"pullPolicy"`
}
type PeerPersistence struct {
	Peer      Persistence `json:"peer"`
	CouchDB   Persistence `json:"couchdb"`
	Chaincode Persistence `json:"chaincode"`
}
type Image struct {
	Repository string `json:"repository"`
	Tag        string `json:"tag"`
	PullPolicy string `json:"pullPolicy"`
}
type Annotations struct {
}
type Gossip struct {
	Bootstrap         string `json:"bootstrap"`
	Endpoint          string `json:"endpoint"`
	ExternalEndpoint  string `json:"externalEndpoint"`
	OrgLeader         bool   `json:"orgLeader"`
	UseLeaderElection bool   `json:"useLeaderElection"`
}
type Server struct {
	Enabled bool `json:"enabled"`
}
type Client struct {
	Enabled bool `json:"enabled"`
}
type TLSAuth struct {
	Server Server `json:"server"`
	Client Client `json:"client"`
}
type Peer struct {
	DatabaseType    string  `json:"databaseType"`
	CouchdbInstance string  `json:"couchdbInstance"`
	MspID           string  `json:"mspID"`
	Gossip          Gossip  `json:"gossip"`
	TLS             TLSAuth `json:"tls"`
}
type TLS struct {
	Cert string `json:"cert"`
	Key  string `json:"key"`
}
type Limits struct {
	CPU    string `json:"cpu"`
	Memory string `json:"memory"`
}
type Requests struct {
	CPU    string `json:"cpu"`
	Memory string `json:"memory"`
}
type Resources struct {
	Limits   Limits   `json:"limits"`
	Requests Requests `json:"requests"`
}
type NodeSelector struct {
}
type Affinity struct {
}
type HostAlias struct {
	IP        string   `json:"ip"`
	Hostnames []string `json:"hostnames"`
}
type Service struct {
	Type               string `json:"type"`
	PortRequest        int    `json:"portRequest"`
	PortEvent          int    `json:"portEvent"`
	PortOperations     int    `json:"portOperations"`
	NodePortOperations int    `json:"nodePortOperations,omitempty"`
	NodePortEvent      int    `json:"nodePortEvent,omitempty"`
	NodePortRequest    int    `json:"nodePortRequest,omitempty"`
}
type Persistence struct {
	Enabled      bool        `json:"enabled"`
	Annotations  Annotations `json:"annotations"`
	StorageClass string      `json:"storageClass"`
	AccessMode   string      `json:"accessMode"`
	Size         string      `json:"size"`
}
type Logging struct {
	Level    string `json:"level"`
	Peer     string `json:"peer"`
	Cauthdsl string `json:"cauthdsl"`
	Gossip   string `json:"gossip"`
	Grpc     string `json:"grpc"`
	Ledger   string `json:"ledger"`
	Msp      string `json:"msp"`
	Policies string `json:"policies"`
}
