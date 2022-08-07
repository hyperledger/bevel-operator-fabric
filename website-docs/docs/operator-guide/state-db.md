---
id: state-db
title: LevelDB / CouchDB
---
:::caution

Once you set the state database of the peer, you cannot change it since the structure is different between using **LevelDB** and **CouchDB** you can find more information in [the official HLF docs](https://hyperledger-fabric.readthedocs.io/en/release-2.3/couchdb_as_state_database.html)

:::


## Configuring LevelDB

In order to configure LevelDB, you need to set the following property in the CRD(Custom resource definition) of the peer:
```yaml
stateDb: leveldb
```

## Configuring CouchDB

You can configure the world state to be CouchDB by setting the property `stateDb` in the CRD of the peer:

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
If you wish to use an external CouchDB instance, [check this page](../couchdb/external-couchdb)

