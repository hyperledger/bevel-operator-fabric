package ordnode

type FabricOrdChart struct {
	Genesis          string        `json:"genesis"`
	Ingress          Ingress       `json:"ingress"`
	Cacert           string        `json:"cacert"`
	Tlsrootcert      string        `json:"tlsrootcert"`
	AdminCert        string        `json:"adminCert"`
	Cert             string        `json:"cert"`
	Key              string        `json:"key"`
	TLS              TLS           `json:"tls"`
	FullnameOverride string        `json:"fullnameOverride"`
	HostAliases      []HostAliases `json:"hostAliases"`
	Service          Service       `json:"service"`
	Image            Image         `json:"image"`
	Persistence      Persistence   `json:"persistence"`
	Ord              Ord           `json:"ord"`
	Clientcerts      Clientcerts   `json:"clientcerts"`
	Hosts            []string      `json:"hosts"`
}
type Ingress struct {
	Enabled     bool          `json:"enabled"`
	Annotations Annotations   `json:"annotations"`
	Path        string        `json:"path"`
	Hosts       []string      `json:"hosts"`
	TLS         []interface{} `json:"tls"`
}
type TLS struct {
	Cert string `json:"cert"`
	Key  string `json:"key"`
}
type HostAliases struct {
	IP        string   `json:"ip"`
	Hostnames []string `json:"hostnames"`
}
type Service struct {
	Type               string `json:"type"`
	Port               int    `json:"port"`
	NodePort           int    `json:"nodePort"`
	PortOperations     int    `json:"portOperations"`
	NodePortOperations int    `json:"nodePortOperations"`
}
type Image struct {
	Repository string `json:"repository"`
	Tag        string `json:"tag"`
	PullPolicy string `json:"pullPolicy"`
}
type Annotations struct {
}
type Persistence struct {
	Enabled      bool        `json:"enabled"`
	Annotations  Annotations `json:"annotations"`
	StorageClass string      `json:"storageClass"`
	AccessMode   string      `json:"accessMode"`
	Size         string      `json:"size"`
}
type Server struct {
	Enabled bool `json:"enabled"`
}
type Client struct {
	Enabled bool `json:"enabled"`
}
type TLSConfiguration struct {
	Server Server `json:"server"`
	Client Client `json:"client"`
}
type Ord struct {
	Type  string           `json:"type"`
	MspID string           `json:"mspID"`
	TLS   TLSConfiguration `json:"tls"`
}
type Clientcerts struct {
	CertPem string `json:"cert.pem,omitempty"`
}
