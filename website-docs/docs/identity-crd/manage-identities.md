---
id: manage-identities
title: Managing identities with CRDs
---

`FabricIdentity` controller uses the internal communication (port 7054) to the Fabric CA that's by default enabled when the Fabric CA is deployed with the operator.

## Create a HLF identity

Use the create command to create a new HLF identity.

```bash
kubectl hlf identity create --name <name> --namespace <namespace> \
    --ca-name <ca-name> --ca-namespace <ca-namespace> \
    --ca <ca> --mspid <mspid> --enroll-id <enroll-id> --enroll-secret <enroll-secret>

```

Arguments:

- --name: Name of the external chaincode.
- --namespace: Namespace of the external chaincode.
- --ca-name: Name of the CA (Certificate Authority).
- --ca-namespace: Namespace of the CA.
- --ca: CA name.
- --mspid: MSP ID.
- --enroll-id: Enroll ID.
- --enroll-secret: Enroll Secret.

## Update HLF Identity

Use the update command to update an existing HLF identity.


```bash
kubectl hlf identity update --name <name> --namespace <namespace> \
    --ca-name <ca-name> --ca-namespace <ca-namespace> --ca <ca> \
    --mspid <mspid> --enroll-id <enroll-id> --enroll-secret <enroll-secret>
```

Arguments:

- --name: Name of the external chaincode.
- --namespace: Namespace of the external chaincode.
- --ca-name: Name of the CA (Certificate Authority).
- --ca-namespace: Namespace of the CA.
- --ca: CA name.
- --mspid: MSP ID.
- --enroll-id: Enroll ID.
- --enroll-secret: Enroll Secret.


## Delete HLF Identity

Use the delete command to delete an existing HLF identity.

```bash
kubectl hlf identity delete --name <name> --namespace <namespace>
```

Arguments:

- --name: Name of the identity.
- --namespace: Namespace of the identity.
