---
id: external-couchdb
title: Using external CouchDB
---


If you want to use an external CouchDB, you can use the `externalCouchDB` parameter in the `fabricpeer` CRD.
```yaml
apiVersion: hlf.kungfusoftware.es/v1alpha1
kind: FabricPeer
metadata:
  creationTimestamp: null
  name: org1-peer0
  namespace: default
spec:
# ...more props
  couchdb:
    password: couchdb
    user: couchdb
    externalCouchDB:
      enabled: true # "true" to use external couchdb, "false" to have it as a sidecar
      host: <EXTERNAL_COUCHDB_HOST>
      port: <EXTERNAL_COUCHDB_PORT>
# ...more props
```

You can change this property at anytime on your peers, but it's recommended to use the sidecar CouchDB, since it makes it easier to create and destroy peers.
