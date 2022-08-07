package console

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/api/networking/v1beta1"
)

type FabricOperationsConsoleChart struct {
	Replicas         int                           `json:"replicaCount"`
	Image            Image                         `json:"image"`
	ImagePullSecrets []corev1.LocalObjectReference `json:"imagePullSecrets"`
	PodAnnotations   map[string]string             `json:"podAnnotations"`
	Ingress          Ingress                       `json:"ingress"`
	Resources        corev1.ResourceRequirements   `json:"resources"`
	Tolerations      []corev1.Toleration           `json:"tolerations"`
	Affinity         *corev1.Affinity              `json:"affinity"`
	Port             int                           `json:"port"`
	HostUrl          string                        `json:"hostUrl"`
	Auth             Auth                          `json:"auth"`
	CouchDB          CouchDB                       `json:"couchdb"`
}

type CouchDB struct {
	External    CouchDBExternal              `json:"external"`
	Image       string                       `json:"image"`
	Tag         string                       `json:"tag"`
	PullPolicy  corev1.PullPolicy            `json:"pullPolicy"`
	Username    string                       `json:"username"`
	Password    string                       `json:"password"`
	Persistence Persistence                  `json:"persistence"`
	Resources   *corev1.ResourceRequirements `json:"resources"`
}
type Persistence struct {
	Enabled      bool                              `json:"enabled"`
	Annotations  map[string]string                 `json:"annotations"`
	StorageClass string                            `json:"storageClass"`
	AccessMode   corev1.PersistentVolumeAccessMode `json:"accessMode"`
	Size         string                            `json:"size"`
}
type CouchDBExternal struct {
	Enabled bool   `json:"enabled"`
	Host    string `json:"host"`
	Port    int    `json:"port"`
}
type Auth struct {
	Scheme   string `json:"scheme"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Image struct {
	Repository string            `json:"repository"`
	Tag        string            `json:"tag"`
	PullPolicy corev1.PullPolicy `json:"pullPolicy"`
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
