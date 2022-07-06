---
id: getting-started
title: Getting started
---

## Deploying Operations Console

This guide intends to showcase the installation of [Fabric Operations Console](https://github.com/hyperledger-labs/fabric-operations-console).

## How to deploy the Fabric Operations console



```bash
CONSOLE_PASSWORD="admin"
kubectl hlf console create --name=console --namespace=default --version="latest" --image="ghcr.io/hyperledger-labs/fabric-console" --admin-user="admin" --admin-pwd="$CONSOLE_PASSWORD"
```
