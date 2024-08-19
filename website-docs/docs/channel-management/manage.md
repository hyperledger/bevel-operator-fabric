---
id: manage
title: Manage the channel
---

## Add peer organization to the channel


You can add more organizations by updating the `peerOrganizations` or `externalPeerOrganizations` property in the `FabricMainChannel` CRD.

If the organization is not in the cluster, you need to add the organization to the `externalPeerOrganizations` property, with the `mspID`, `signRootCert` and `tlsRootCert`.

```yaml
  externalPeerOrganizations:
    - mspID: <MSP_ID>
      signRootCert: |
        <SIGN_ROOT_CRT_PEM>
      tlsRootCert: |
        <TLS_ROOT_CRT_PEM>
```

If the organization is in the cluster, you need to add the organization to the `peerOrganizations` property, with the `mspID`, `signRootCert` and `tlsRootCert`.

```yaml
  peerOrganizations:
    - caName: <CA_NAME>
      caNamespace: <CA_NS>
      mspID: <MSP_ID>
```



## Add orderer organization to the channel


You can add more organizations by updating the `peerOrganizations` or `externalPeerOrganizations` property in the `FabricMainChannel` CRD.

If the organization is not in the cluster, you need to add the organization to the `externalPeerOrganizations` property, with the `mspID`, `signRootCert` and `tlsRootCert`.

```yaml
  externalOrdererOrganizations:
    - mspID: <MSP_ID>
      signRootCert: |
        <SIGN_ROOT_CRT_PEM>
      tlsRootCert: |
        <TLS_ROOT_CRT_PEM>
      ordererEndpoints: # orderer endpoints for the organization in the channel configuration
        - <ORDERER0_ENDPOINT>
```

If the organization is in the cluster, you need to add the organization to the `peerOrganizations` property, with the `mspID`, `signRootCert` and `tlsRootCert`.

```yaml
  ordererOrganizations:
    - caName: <CA_NAME>
      caNamespace: <CA_NS>
      externalOrderersToJoin:
        - host: <ADMIN_ORDERER_HOST>
          port: <ADMIN_ORDERER_PORT>
      mspID: <MSP_ID>
      ordererEndpoints: # orderer endpoints for the organization in the channel configuration
        - <ORDERER0_ENDPOINT>
      orderersToJoin: []
```


