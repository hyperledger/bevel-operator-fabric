---
id: getting-started
title: Getting started
---
[![Build Status](https://img.shields.io/travis/kfsoftware/hlf-operator/main.svg?label=E2E%20testing)](https://travis-ci.org/kfsoftware/hlf-operator)

# Hyperledger Fabric Operator

## Features

- [x] Create certificates authorities (CA)
- [x] Create peers
- [x] Create ordering services
- [x] Create resources without manual provisioning of cryptographic material
- [x] Domain routing with SNI using Istio
- [x] Run chaincode as external chaincode in Kubernetes
- [x] Support Hyperledger Fabric 2.3+
- [x] Managed genesis for Ordering services
- [x] E2E testing including the execution of chaincodes in KIND
- [x] Renewal of certificates

## Stay Up-to-Date

`hlf-operator` is currently in stable. Watch **releases** of this repository to be notified for future updates:

![hlf-operator-star-github](https://user-images.githubusercontent.com/6862893/123808402-022aa800-d8f1-11eb-8df4-8a9552f126a2.gif)

## Sponsor

|||
|-|-|
|![kfs logo](https://avatars.githubusercontent.com/u/74511895?s=200&v=4)|If you want to design and deploy a secure Blockchain network based on the latest version of  Hyperledger Fabric, feel free to contact dviejo@kungfusoftware.es or visit [https://kfs.es/blockchain](https://kfs.es/blockchain)|


## Getting started

### Requirements

- Fabric CA client
- YQ binary to replace values in YAML (for the getting started)
- KubeCTL
- Kubernetes 1.15+
- Istio

### Install Istio

```bash
kubectl apply -f ./hack/istio-operator/crds/*
helm template ./hack/istio-operator/ \
  --set hub=docker.io/istio \
  --set tag=1.8.0 \
  --set operatorNamespace=istio-operator \
  --set watchedNamespaces=istio-system | kubectl apply -f -

kubectl create ns istio-system
kubectl apply -n istio-system -f ./hack/istio-operator.yaml
```

### Installing the operator

```bash
helm install hlf-operator ./chart/hlf-operator
```

### Installing the Kubectl HLF Plugin


```bash
kubectl krew install hlf
```

## Deploy a Peer Organization

### Setup versions
```bash
export PEER_IMAGE=quay.io/kfsoftware/fabric-peer
export PEER_VERSION=2.4.1-v0.0.3

export ORDERER_IMAGE=hyperledger/fabric-orderer
export ORDERER_VERSION=2.4.1

```

### Deploying a Certificate Authority

```bash
kubectl hlf ca create --storage-class=standard --capacity=2Gi --name=org1-ca \
    --enroll-id=enroll --enroll-pw=enrollpw  
kubectl wait --timeout=180s --for=condition=Running fabriccas.hlf.kungfusoftware.es --all

# register user for the peers
kubectl hlf ca register --name=org1-ca --user=peer --secret=peerpw --type=peer \
 --enroll-id enroll --enroll-secret=enrollpw --mspid Org1MSP
```

### Deploying a peer

```bash

kubectl hlf peer create --image=$PEER_IMAGE --version=$PEER_VERSION --storage-class=standard --enroll-id=peer --mspid=Org1MSP \
        --enroll-pw=peerpw --capacity=5Gi --name=org1-peer0 --ca-name=org1-ca.default
kubectl wait --timeout=180s --for=condition=Running fabricpeers.hlf.kungfusoftware.es --all
```

## Deploying an Ordering Service

### Deploying a certificate authority

```bash
kubectl hlf ca create --storage-class=standard --capacity=2Gi --name=ord-ca \
    --enroll-id=enroll --enroll-pw=enrollpw
kubectl wait --timeout=180s --for=condition=Running fabriccas.hlf.kungfusoftware.es --all
kubectl hlf ca register --name=ord-ca --user=orderer --secret=ordererpw \
    --type=orderer --enroll-id enroll --enroll-secret=enrollpw --mspid=OrdererMSP

```

### Deploying the Orderer nodes node

```bash
kubectl hlf ordnode create --image=$ORDERER_IMAGE --version=$ORDERER_VERSION \
    --storage-class=standard --enroll-id=orderer --mspid=OrdererMSP \
    --enroll-pw=ordererpw --capacity=2Gi --name=ord-node1 --ca-name=ord-ca.default
kubectl wait --timeout=180s --for=condition=Running fabricorderernodes.hlf.kungfusoftware.es --all
```

## Preparing a connection string for the ordering service
```bash
kubectl hlf inspect --output ordservice.yaml -o OrdererMSP
kubectl hlf ca register --name=ord-ca --user=admin --secret=adminpw \
    --type=admin --enroll-id enroll --enroll-secret=enrollpw --mspid=OrdererMSP

kubectl hlf ca enroll --name=ord-ca --user=admin --secret=adminpw --mspid OrdererMSP \
        --ca-name ca  --output admin-ordservice.yaml 
## add user from admin-ordservice.yaml to ordservice.yaml
kubectl hlf utils adduser --userPath=admin-ordservice.yaml --config=ordservice.yaml --username=admin --mspid=OrdererMSP
```

## Create a channel
```bash
kubectl hlf channel generate --output=demo.block --name=demo --organizations Org1MSP --ordererOrganizations OrdererMSP

# enroll using the TLS CA
kubectl hlf ca enroll --name=ord-ca --namespace=default --user=admin --secret=adminpw --mspid OrdererMSP \
        --ca-name tlsca  --output admin-tls-ordservice.yaml 

kubectl hlf ordnode join --block=demo.block --name=ord-node1 --namespace=default --identity=admin-tls-ordservice.yaml

```


## Preparing a connection string for the peer
```bash
kubectl hlf ca register --name=org1-ca --user=admin --secret=adminpw --type=admin \
 --enroll-id enroll --enroll-secret=enrollpw --mspid Org1MSP  

kubectl hlf ca enroll --name=org1-ca --user=admin --secret=adminpw --mspid Org1MSP \
        --ca-name ca  --output peer-org1.yaml

kubectl hlf inspect --output org1.yaml -o Org1MSP -o OrdererMSP

## add user key and cert to org1.yaml from admin-ordservice.yaml
kubectl hlf utils adduser --userPath=peer-org1.yaml --config=org1.yaml --username=admin --mspid=Org1MSP
```


## Join channel
```bash
kubectl hlf channel join --name=demo --config=org1.yaml \
    --user=admin -p=org1-peer0.default

```
## Inspect the channel
```bash
kubectl hlf channel inspect --channel=demo --config=org1.yaml \
    --user=admin -p=org1-peer0.default > demo.json
```
## Add anchor peer
```bash
kubectl hlf channel addanchorpeer --channel=demo --config=org1.yaml \
    --user=admin --peer=org1-peer0.default 

```

## See ledger height
In case of error, you may need to add the following to the org1.yaml configuration file:
```yaml
channels:
  _default:
    orderers:
      - ord-node1.default
    peers:
      "org1-peer0.default":
        endorsingPeer: true
        chaincodeQuery: true
        ledgerQuery: true
        eventSource: true

```
```bash
kubectl hlf channel top --channel=demo --config=org1.yaml \
    --user=admin -p=org1-peer0.default
```

## Install a chaincode
```bash
kubectl hlf chaincode install --path=./fixtures/chaincodes/fabcar/go \
    --config=org1.yaml --language=golang --label=fabcar --user=admin --peer=org1-peer0.default

# this can take 3-4 minutes
```

## Query chaincodes installed
```bash
kubectl hlf chaincode queryinstalled --config=org1.yaml --user=admin --peer=org1-peer0.default
```

## Approve chaincode
```bash
PACKAGE_ID=fabcar:0c616be7eebace4b3c2aa0890944875f695653dbf80bef7d95f3eed6667b5f40 # replace it with the package id of your 
kubectl hlf chaincode approveformyorg --config=org1.yaml --user=admin --peer=org1-peer0.default \
    --package-id=$PACKAGE_ID \
    --version "1.0" --sequence 1 --name=fabcar \
    --policy="OR('Org1MSP.member')" --channel=demo
```

## Commit chaincode
```bash
kubectl hlf chaincode commit --config=org1.yaml --user=admin --mspid=Org1MSP \
    --version "1.0" --sequence 1 --name=fabcar \
    --policy="OR('Org1MSP.member')" --channel=demo
```


## Invoke a transaction in the ledger
```bash
kubectl hlf chaincode invoke --config=org1.yaml \
    --user=admin --peer=org1-peer0.default \
    --chaincode=fabcar --channel=demo \
    --fcn=initLedger -a '[]'
```

## Query the ledger
```bash
kubectl hlf chaincode query --config=org1.yaml \
    --user=admin --peer=org1-peer0.default \
    --chaincode=fabcar --channel=demo \
    --fcn=QueryAllCars -a '[]'
```

At this point, you should have:

- Ordering service with 3 nodes and a CA
- Peer organization with a peer and a CA
- A channel **demo**
- A chaincode install in peer0
- A chaincode approved and committed

If something went wrong or didn't work, please, open an issue.

## Cleanup the environment

```bash
kubectl delete fabricorderernodes.hlf.kungfusoftware.es --all-namespaces --all
kubectl delete fabricpeers.hlf.kungfusoftware.es --all-namespaces --all
kubectl delete fabriccas.hlf.kungfusoftware.es --all-namespaces --all
```

## Troubleshooting

### Chaincode installation/build error
Chaincode installation/build can fail due to unsupported local kubertenes version such as [minikube](https://github.com/kubernetes/minikube).

```shell
$ kubectl hlf chaincode install --path=./fixtures/chaincodes/fabcar/go \
        --config=org1.yaml --language=golang --label=fabcar --user=admin --peer=org1-peer0.default
        
Error: Transaction processing for endorser [192.168.49.2:31278]: Chaincode status Code: (500) UNKNOWN. 
Description: failed to invoke backing implementation of 'InstallChaincode': could not build chaincode: 
external builder failed: external builder failed to build: external builder 'my-golang-builder' failed:
exit status 1
```

If your purpose is to test the hlf-operator please consider to switch to [kind](https://github.com/kubernetes-sigs/kind) that is tested and supported.
