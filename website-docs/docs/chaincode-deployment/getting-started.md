---
id: getting-started
title: Getting started
---

The chaincode is a program that handles business logic agreed to by members of the network, as a smart contract. Chaincode can run in a variety of different platforms:
- Docker
- Kubernetes (not neccesarily on Docker, it can be ContainerD)
- On baremetal

Since this documentation is for the HLF Operator, we will focus on deploying the chaincode in Kubernetes.

If you want to know how to develop chaincodes remotely from your machine, you can read the [developing-chaincode](../chaincode-development/getting-started.md) section.


In order to know how to deploy the chaincode, you have to options:
- [Deploy using the external chaincode as a service](./external-chaincode-as-a-service.md)
- [Deploy using the Kubernetes builder](./k8s-builder.md)

