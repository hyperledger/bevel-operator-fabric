package helpers

import (
	corev1 "k8s.io/api/core/v1"
)

const (
	// StaticYamlPath is path to the static yaml files
	//	StaticYamlPath = "https://raw.githubusercontent.com/nitisht/kubectl-minio/master/cmd/static/static.yaml"

	// ClusterRoleBindingName is the name for CRB
	ClusterRoleBindingName = "minio-operator-binding"

	// ClusterRoleName is the name for Cluster Role for operator
	ClusterRoleName = "minio-operator-role"

	// ContainerName is the name of operator container
	ContainerName = "minio-operator"

	// DeploymentName is the name of operator deployment
	DeploymentName = "minio-operator"

	// DefaultNamespace is the default namespace for all operations
	DefaultNamespace = "default"

	// DefaultStorageclass is the default storage class for MinIO Tenant PVC
	DefaultStorageclass = ""

	// DefaultServiceAccount is the service account for all
	DefaultServiceAccount = "minio-operator"

	// DefaultClusterDomain is the default domain of the Kubernetes cluster
	DefaultClusterDomain = "cluster.local"

	// DefaultSecretNameSuffix is the suffix added to tenant name to create the
	// credentials secret for this tenant
	DefaultSecretNameSuffix = "-creds-secret"

	// DefaultServiceNameSuffix is the suffix added to tenant name to create the
	// internal clusterIP service for this tenant
	DefaultServiceNameSuffix = "-internal-service"

	// MinIOPrometheusPath is the path where MinIO tenant exposes Prometheus metrics
	MinIOPrometheusPath = "/minio/prometheus/metrics"

	// MinIOPrometheusPort is the port where MinIO tenant exposes Prometheus metrics
	MinIOPrometheusPort = "9000"

	// MinIOMountPath is the path where MinIO related PVs are mounted in a container
	MinIOMountPath = "/export"

	// MinIOAccessMode is the default access mode to be used for PVC / PV binding
	MinIOAccessMode = "ReadWriteOnce"

	// DefaultImagePullPolicy specifies the policy to image pulls
	DefaultImagePullPolicy = corev1.PullIfNotPresent

	// DefaultOperatorImage is the default operator image to be used
	DefaultOperatorImage = "minio/k8s-operator:v3.0.28"

	// DefaultCAImage is the default MinIO image used while creating tenant
	DefaultCAImage   = "hyperledger/fabric-ca"
	DefaultCAVersion = "1.4.9"

	DefaultPeerImage   = "quay.io/kfsoftware/fabric-peer"
	DefaultPeerVersion = "amd64-2.2.0"

	DefaultOrdererImage = "hyperledger/fabric-orderer"
	DefaultOrdererVersion = "amd64-2.2.0"

	// DefaultKESImage is the default KES image used while creating tenant
	DefaultKESImage = "minio/kes:v0.11.0"

	// DefaultConsoleImage is the default console image used while creating tenant
	DefaultConsoleImage = "minio/console:v0.3.14"
)

// DeploymentReplicas is the number of replicas for MinIO Operator
var DeploymentReplicas int32 = 1

// KESReplicas is the number of replicas for MinIO KES
var KESReplicas int32 = 2

// ConsoleReplicas is the number of replicas for MinIO Console
var ConsoleReplicas int32 = 2
