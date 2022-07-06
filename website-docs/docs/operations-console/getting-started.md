---
id: getting-started
title: Getting started
---

## Deploying Operations Console

This guide intends to showcase the installation of [Fabric Operations Console](https://github.com/hyperledger-labs/fabric-operations-console).

## How to deploy the Fabric Operations console

:::caution

Once you set the state database of the peer, you cannot change it since the structure is different between using **LevelDB** and **CouchDB** you can find more information in [the official HLF docs](https://hyperledger-fabric.readthedocs.io/en/release-2.3/couchdb_as_state_database.html)

:::

### Generate a certificate for TLS

This step is critical since the Operations Console need a secure communication in order to connect with the CAs, Peers and Orderer nodes

```bash
CONSOLE_PASSWORD="admin"
kubectl hlf console create --name=console --namespace=default --version="latest" --image="ghcr.io/hyperledger-labs/fabric-console" --admin-user="admin" --admin-pwd="$CONSOLE_PASSWORD"
```

