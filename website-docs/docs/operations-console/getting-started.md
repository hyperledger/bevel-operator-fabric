---
id: getting-started
title: Getting started
---

## Deploying Operations Console

This guide intends to showcase the installation of [Fabric Operations Console](https://github.com/hyperledger-labs/fabric-operations-console).

## How to deploy the Fabric Operations console

:::caution

Since the

:::

### Generate a certificate for TLS

This step is critical since the Operations Console need a secure communication in order to connect with the CAs, Peers and Orderer nodes

```bash
CONSOLE_PASSWORD="admin"
kubectl hlf console create --name=console --namespace=default --version="latest" --image="ghcr.io/hyperledger-labs/fabric-console" --admin-user="admin" --admin-pwd="$CONSOLE_PASSWORD"
```

