package ordservice

type FabricOrdChart struct {
	FullNameOverride string  `json:"fullnameOverride"`
	Image            Image   `json:"image"`
	Genesis          string  `json:"genesis"`
	Storage          Storage `json:"storage"`
	MspID            string  `json:"mspID"`
	Nodes            []Node  `json:"nodes"`
}
type Image struct {
	Repository string `json:"repository"`
	Tag        string `json:"tag"`
	PullPolicy string `json:"pullPolicy"`
}
type Storage struct {
	Size         string `json:"size"`
	AccessMode   string `json:"accessMode"`
	StorageClass string `json:"storageClass"`
}
type Service struct {
	Type            string `json:"type"`
	NodePortRequest int    `json:"nodePortRequest"`
}
type Node struct {
	SignKey      string   `json:"signKey"`
	SignCert     string   `json:"signCert"`
	SignRootCert string   `json:"signRootCert"`
	TLSCert      string   `json:"tlsCert"`
	TLSKey       string   `json:"tlsKey"`
	TLSRootCert  string   `json:"tlsRootCert"`
	Hosts        []string `json:"hosts"`
	Service      Service  `json:"service"`
}
