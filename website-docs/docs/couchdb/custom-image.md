---
id: custom-image
title: Using custom CouchDB image
---


If you want to use a custom image for the CouchDB instance, you can use the `image` and `tag` parameter in the `fabricpeer` CRD.
```yaml
apiVersion: hlf.kungfusoftware.es/v1alpha1
kind: FabricPeer
metadata:
  name: org1-peer0
  namespace: default
spec:
# ...more props
  couchdb:
    password: couchdb
    user: couchdb
    image: <yourownreg>
    tag: <yourtag>
# ...more props
```
