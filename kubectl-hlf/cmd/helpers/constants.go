package helpers

const (

	// DefaultNamespace is the default namespace for all operations
	DefaultNamespace = "default"

	DefaultStorageclass = ""

	DefaultCAImage   = "hyperledger/fabric-ca"
	DefaultCAVersion = "1.5.7"

	DefaultPeerImage   = "hyperledger/fabric-peer"
	DefaultPeerVersion = "2.5.5"

	DefaultOperationsConsoleImage   = "ghcr.io/hyperledger-labs/fabric-console"
	DefaultOperationsConsoleVersion = "latest"

	DefaultOperationsOperatorUIImage   = "ghcr.io/kfsoftware/hlf-operator-ui"
	DefaultOperationsOperatorUIVersion = "0.0.16"

	DefaultOperationsOperatorAPIImage   = "ghcr.io/kfsoftware/hlf-operator-api"
	DefaultOperationsOperatorAPIVersion = "v0.0.16"

	DefaultFSServerImage   = "quay.io/kfsoftware/fs-peer"
	DefaultFSServerVersion = "amd64-2.2.0-0.0.1"

	DefaultCouchDBImage   = "couchdb"
	DefaultCouchDBVersion = "3.1.1"

	DefaultOrdererImage   = "hyperledger/fabric-orderer"
	DefaultOrdererVersion = "2.5.5"
)
