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

## Roadmap

- [ ] More parametrization on the Peer
- [ ] More parametrization on the Fabric CA
- [ ] More parametrization on the Fabric Ordering services

## Ideas for the future

- [ ] Install chaincode in peer using Custom Resource Definitions
- [ ] Manage channel configuration using Custom Resource Definitions

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
# when kubectl-hlf is in the krew-index
kubectl krew install hlf 
# now
kubectl krew install --manifest=krew-plugin.yaml
```

## Deploy a Peer Organization

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

kubectl hlf peer create --storage-class=standard --enroll-id=peer --mspid=Org1MSP \
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

### Deploying the Ordering service

```bash
kubectl hlf ordservice create  --storage-class=standard --enroll-id=orderer --mspid=OrdererMSP \
    --enroll-pw=ordererpw --capacity=2Gi --name=ordservice --ca-name=ord-ca.default \
    --system-channel testchainid --num-orderers=1
kubectl wait --timeout=180s --for=condition=Running fabricorderingservices.hlf.kungfusoftware.es --all
```

## Preparing a connection string for the ordering service
```bash
kubectl hlf inspect --output ordservice.yaml -o OrdererMSP
kubectl hlf ca register --name=ord-ca --user=admin --secret=adminpw \
    --type=admin --enroll-id enroll --enroll-secret=enrollpw --mspid=Ord2MSP

kubectl hlf ca enroll --name=ord-ca --user=admin --secret=adminpw --mspid Ord2MSP \
        --ca-name ca  --output admin-ordservice.yaml 
## add user from admin-ordservice.yaml to ordservice.yaml

```

## Create a consortium
```bash
kubectl hlf consortiums create --name=Default --system-channel-id="testchainid" \
    --config=ordservice.yaml --orderer-org=ordservice.default --user=admin \
    -p=org1-peer0.default
```
## Preparing a connection string for the peer
```bash
kubectl hlf ca register --name=org1-ca --user=admin --secret=adminpw --type=admin \
 --enroll-id enroll --enroll-secret=enrollpw --mspid Org1MSP  

kubectl hlf ca enroll --name=org1-ca --user=admin --secret=adminpw --mspid Org1MSP \
        --ca-name ca  --output peer-org1.yaml

kubectl hlf inspect --output org1.yaml -o Org1MSP -o OrdererMSP

## add user key and cert to org1.yaml from admin-ordservice.yaml
```

## Create a channel
```bash

kubectl hlf channel create --name=ch1 --config=org1.yaml \
    --admin-org=org1-peer0.default --user=admin \
    -p=org1-peer0.default --ordering-service=ordservice.default \
    --consortium=Default
```

## Add anchor peer
```bash
kubectl hlf channel addanchorpeer --channel=ch1 --config=org1.yaml \
    --user=admin --peer=org1-peer0.default 

```
## Join channel
```bash
kubectl hlf channel join --name=ch1 --config=org1.yaml \
    --user=admin -p=org1-peer0.default

```


## See ledger height
In case of error, you may need to add the following to the org1.yaml configuration file:
```yaml
channels:
  _default:
    peers:
      "org1-peer0.default":
          endorsingPeer: true
          chaincodeQuery: true
          ledgerQuery: true
          eventSource: true
```
```bash
kubectl hlf channel top --channel=ch1 --config=org1.yaml \
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
kubectl hlf chaincode approveformyorg --config=org1.yaml --user=admin --peer=org1-peer0.default \
    --package-id=fabcar:db8d009f7e2e9fa4a40ddfd6b7e603d3177b126d18cdbeabcf8481f9a6de519f \
    --version "1.0" --sequence 1 --name=fabcar \
    --policy="OR('Org1MSP.member')" --channel=ch1
```

## Commit chaincode
```bash
kubectl hlf chaincode commit --config=org1.yaml --user=admin --peer=org1-peer0.default \
    --version "1.0" --sequence 1 --name=fabcar \
    --policy="OR('Org1MSP.member')" --channel=ch1
```


## Invoke a transaction in the ledger
```bash
kubectl hlf chaincode invoke --config=org1.yaml \
    --user=admin --peer=org1-peer0.default \
    --chaincode=fabcar --channel=ch1 \
    --fcn=initLedger -a '[]'
```

## Query the ledger
```bash
kubectl hlf chaincode query --config=org1.yaml \
    --user=admin --peer=org1-peer0.default \
    --chaincode=fabcar --channel=ch1 \
    --fcn=QueryAllCars -a '[]'
```

At this point, you should have:

- Ordering service with 3 nodes and a CA
- Peer organization with a peer and a CA
- A channel **ch1**
- A chaincode install in peer0
- A chaincode approved and committed

If something went wrong or didn't work, please, open an issue.

### Cleanup the environment

```bash
kubectl delete fabricorderingservices.hlf.kungfusoftware.es --all-namespaces --all
kubectl delete fabricpeers.hlf.kungfusoftware.es --all-namespaces --all
kubectl delete fabriccas.hlf.kungfusoftware.es --all-namespaces --all
```
