package helpers

const (

	// DefaultNamespace is the default namespace for all operations
	DefaultNamespace = "default"

	DefaultStorageclass = ""

	DefaultCAImage   = "hyperledger/fabric-ca"
	DefaultCAVersion = "1.4.9"

	DefaultPeerImage   = "quay.io/kfsoftware/fabric-peer"
	DefaultPeerVersion = "2.3.0-v0.0.2"

	DefaultFSServerImage   = "quay.io/kfsoftware/fs-peer"
	DefaultFSServerVersion = "amd64-2.2.0-0.0.1"

	DefaultCouchDBImage   = "couchdb"
	DefaultCouchDBVersion = "3.1.1"

	DefaultOrdererImage   = "hyperledger/fabric-orderer"
	DefaultOrdererVersion = "amd64-2.3.0"
)
