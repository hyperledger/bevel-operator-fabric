package peer

type RBAC struct {
	Ns string `json:"ns"`
}
type FabricPeerChart struct {
	ExternalChaincodeBuilder bool            `json:"externalChaincodeBuilder"`
	CouchdbUsername          string          `json:"couchdbUsername"`
	CouchdbPassword          string          `json:"couchdbPassword"`
	Image                    Image           `json:"image"`
	Rbac                     RBAC            `json:"rbac"`
	DockerSocketPath         string          `json:"dockerSocketPath"`
	Ingress                  Ingress         `json:"ingress"`
	Peer                     Peer            `json:"peer"`
	Cert                     string          `json:"cert"`
	Key                      string          `json:"key"`
	Hosts                    []string        `json:"hosts"`
	OperationHosts           []string        `json:"operationHosts"`
	TLS                      TLS             `json:"tls"`
	OPSTLS                   TLS             `json:"opsTLS"`
	Cacert                   string          `json:"cacert"`
	Tlsrootcert              string          `json:"tlsrootcert"`
	Resources                Resources       `json:"resources,omitempty"`
	NodeSelector             NodeSelector    `json:"nodeSelector,omitempty"`
	Tolerations              []interface{}   `json:"tolerations"`
	Affinity                 Affinity        `json:"affinity,omitempty"`
	ExternalHost             string          `json:"externalHost"`
	FullnameOverride         string          `json:"fullnameOverride"`
	HostAliases              []HostAliases   `json:"hostAliases"`
	Service                  Service         `json:"service"`
	Persistence              PeerPersistence `json:"persistence"`
	Logging                  Logging         `json:"logging"`
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
type Ingress struct {
	Enabled     bool          `json:"enabled"`
	Annotations Annotations   `json:"annotations"`
	Path        string        `json:"path"`
	Hosts       []string      `json:"hosts"`
	TLS         []interface{} `json:"tls"`
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
type HostAliases struct {
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
