---
id: state-db
title: LevelDB / CouchDB
---

## Configuring LevelDB

In order to configure LevelDB, you need to set the following property in the CRD(Custom resource definition) of the peer:
```yaml
stateDb: leveldb
```

## Configuring CouchDB

You can configure couchdb by setting the following property in the yaml:

```yaml
stateDb: couchdb
```

And then you can configure also the username and password:
```yaml
couchdb:
  externalCouchDB: null
  password: couchdb
  user: couchdb
```

If you want to configure a custom image for CouchDB, you can set the `image`, `tag`, and `pullPolicy` properties under the `couchdb` property:
```yaml
couchdb:
    image: couchdb
    pullPolicy: IfNotPresent
    tag: 3.1.1
    user: couchdb
    password: couchdb
```
If you wish to use an external CouchDB instance, [check this page](./external-couchdb)
