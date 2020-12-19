package ordnode

type fabricOrdChart struct {
	Genesis          string        `json:"genesis"`
	Ingress          ingress       `json:"ingress"`
	Cacert           string        `json:"cacert"`
	Tlsrootcert      string        `json:"tlsrootcert"`
	AdminCert        string        `json:"adminCert"`
	Cert             string        `json:"cert"`
	Key              string        `json:"key"`
	TLS              tls           `json:"tls"`
	FullnameOverride string        `json:"fullnameOverride"`
	HostAliases      []hostAliases `json:"hostAliases"`
	Service          service       `json:"service"`
	Image            image         `json:"image"`
	Persistence      persistence   `json:"persistence"`
	Ord              ord           `json:"ord"`
	Clientcerts      clientcerts   `json:"clientcerts"`
	Hosts            []string      `json:"hosts"`
}
type ingress struct {
	Enabled     bool          `json:"enabled"`
	Annotations annotations   `json:"annotations"`
	Path        string        `json:"path"`
	Hosts       []string      `json:"hosts"`
	TLS         []interface{} `json:"tls"`
}
type tls struct {
	Cert string `json:"cert"`
	Key  string `json:"key"`
}
type hostAliases struct {
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
