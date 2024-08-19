---
id: install-crd
title: Installing chaincodes with the FabricChaincodeInstall CRD
---

Certainly! I'll provide you with the complete, revised documentation page that incorporates all the information and improvements we've discussed. This will be in a format you can easily copy and paste.

# FabricChaincodeInstall CRD Documentation

## Overview

The `FabricChaincodeInstall` Custom Resource Definition (CRD) is used in the bevel-operator-fabric to install chaincodes on specific peers of an organization in a Hyperledger Fabric network. This CRD allows for fine-grained control over chaincode installation, specifying which peers should have the chaincode installed and the details of the chaincode package.

## Specification

### API Version and Kind

```yaml
apiVersion: hlf.kungfusoftware.es/v1alpha1
kind: FabricChaincodeInstall
```

### Metadata

The metadata section includes standard Kubernetes object metadata:

- `name`: A unique name for the FabricChaincodeInstall resource
- `namespace`: The Kubernetes namespace where the resource is created

### Spec

The `spec` section defines the desired state of the FabricChaincodeInstall resource:

#### Peers

- `peers`: A list of peers where the chaincode should be installed
  - `name`: The name of the peer
  - `namespace`: The namespace where the peer is deployed

#### External Peers

- `externalPeers`: A list of external peers (if any) where the chaincode should be installed

#### MSP ID

- `mspID`: The Membership Service Provider ID of the organization

#### HLF Identity

- `hlfIdentity`: Specifies the identity used for chaincode installation
  - `secretName`: Name of the Kubernetes secret containing the identity
  - `secretNamespace`: Namespace of the secret
  - `secretKey`: Key in the secret that contains the identity information

#### Chaincode Package

- `chaincodePackage`: Details of the chaincode to be installed
  - `name`: Name of the chaincode
  - `address`: Address where the chaincode is hosted
  - `type`: Type of the chaincode (e.g., 'ccaas' for Chaincode as a Service)
  - `dialTimeout`: Timeout for dialing the chaincode address
  - `tls`: TLS configuration for the chaincode
    - `required`: Boolean indicating if TLS is required

## Example Usage

```yaml
# FabricChaincodeInstall CRD Example with Field Descriptions

# API version of the CRD
apiVersion: hlf.kungfusoftware.es/v1alpha1
# Kind specifies that this is a FabricChaincodeInstall resource
kind: FabricChaincodeInstall
metadata:
  # Name of this FabricChaincodeInstall resource
  name: example-chaincode
  # Namespace where this resource will be created
  namespace: default
spec:
  # List of peers where the chaincode should be installed
  peers:
    # Each item in the list represents a peer
    - name: org1-peer0  # Name of the peer
      namespace: default  # Namespace where the peer is deployed
  # List of external peers (if any) where the chaincode should be installed
  # This is empty in this example
  externalPeers: []
  # Membership Service Provider ID of the organization
  mspID: Org1MSP
  # Identity used for chaincode installation
  hlfIdentity:
    # Name of the Kubernetes secret containing the identity
    secretName: org1-admin
    # Namespace where the secret is located
    secretNamespace: default
    # Key in the secret that contains the identity information
    secretKey: user.yaml
  # Details of the chaincode package to be installed
  chaincodePackage:
    # Name of the chaincode
    name: test
    # Address where the chaincode is hosted
    # Format: <service-name>.<namespace>:<port>
    address: 'example-chaincode.default:9999'
    # Type of the chaincode (e.g., 'ccaas' for Chaincode as a Service)
    type: 'ccaas'
    # Timeout for dialing the chaincode address
    dialTimeout: "10s"
    # TLS configuration for the chaincode
    tls:
      # Boolean indicating if TLS is required
      required: false
```

## Installation Process

When applying this CRD, the bevel-operator-fabric will perform the following steps:

1. Validate the CRD specification
2. Locate the specified peers within the cluster
3. Retrieve the HLF identity from the specified Kubernetes secret
4. Prepare the chaincode package based on the provided details
5. Connect to each specified peer
6. Install the chaincode package on each peer
7. Verify successful installation
8. Update the status of the FabricChaincodeInstall resource

## Notes

- Ensure that the specified peers are operational and accessible within the cluster
- The HLF identity used must have sufficient permissions to install chaincodes
- For external peers, additional configuration may be required to ensure connectivity
- The chaincode package must be available at the specified address before applying this CRD
- Adjust the `dialTimeout` as needed based on your network conditions
- Configure TLS settings appropriately for your environment

## Troubleshooting

If the chaincode installation fails, check the following:

- Peer accessibility and health
- Correct MSP ID
- Valid HLF identity and permissions
- Chaincode package availability and correctness
- Network connectivity to the chaincode address
- TLS configuration (if applicable)

Consult the bevel-operator-fabric logs for detailed error messages and installation status.