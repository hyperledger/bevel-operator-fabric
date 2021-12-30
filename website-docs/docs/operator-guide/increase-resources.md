---
id: increase-resources
title: Increase resources
---

## Increase resources for the peer

To increase the storage for the orderer node, you can modify the `resources` section in the fabricpeers object.

Note that there are 5 different types of resources that can be increased:
- peer
- couchdb
- chaincode
- couchdbExporter

```yaml
resources:
  peer:
    limits:
      cpu: "2"
      memory: 2Gi
    requests:
      cpu: 10m
      memory: 256Mi
  chaincode:
    limits:
      cpu: 500m
      memory: 256Mi
    requests:
      cpu: 10m
      memory: 256Mi
  couchdb:
    limits:
      cpu: "2"
      memory: 2Gi
    requests:
      cpu: 10m
      memory: 256Mi
  couchdbExporter:
    limits:
      cpu: 500m
      memory: 256Mi
    requests:
      cpu: 10m
      memory: 256Mi
```


### Peer
These resources are the ones used for the fabric-peer container.

### CouchDB
These resources are the ones used for the fabric-couchdb container.


### CouchDB Exporter
These resources are the ones used for the fabric-couchdb-exporter container in case it's enabled with the following property:

```yaml
  couchDBexporter:
    enabled: true
    image: gesellix/couchdb-prometheus-exporter
    imagePullPolicy: IfNotPresent
    tag: v30.0.0
```

### Chaincode

This is used in case externalBuilder is enabled, in which case the chaincode container is created, this container is used to store the chaincode build output.


## Increase storage for the orderer

To increase the storage for the orderer node, you can modify the `resources` section in the fabricorderernode object

```yaml
resources:
  limits:
    cpu: "2"
    memory: 2Gi
  requests:
    cpu: 10m
    memory: 256Mi
```

## Increase storage for the certificate authority

To increase the storage for the certificate authority, you can modify the `resources` section in the fabriccas object

```yaml
resources:
  limits:
    cpu: "2"
    memory: 2Gi
  requests:
    cpu: 10m
    memory: 256Mi
```
