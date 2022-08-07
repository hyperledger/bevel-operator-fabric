---
id: getting-started
title: Getting started
---

## Deploying Operations Console

This guide intends to showcase the installation of [Fabric Operations Console](https://github.com/hyperledger-labs/fabric-operations-console).

## How to deploy the Fabric Operations console

:::caution

Since the Fabric Operations Console connects directly with the peers/orderers and CAs, the console needs to be served via HTTPS.

Make sure you use [cert-manager](https://cert-manager.io/docs/) to generate the certificates and then specify the generated secret while creating the Fabric Operations Console.

:::

### Generate a certificate for TLS

This step is critical since the Operations Console need a secure communication in order to connect with the CAs, Peers and Orderer nodes

```bash
export CONSOLE_PASSWORD="admin"
export TLS_SECRET_NAME="console-operator-tls"
kubectl hlf console create --name=console --namespace=default --version="latest" --image="ghcr.io/hyperledger-labs/fabric-console" \
      --admin-user="admin" --admin-pwd="$CONSOLE_PASSWORD" --tls-secret-name="$TLS_SECRET_NAME"
```

