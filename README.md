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
## Deploy a Peer Organization

### Deploying a Certificate Authority
```bash
kubectl apply -f ./config/samples/hlf_v1alpha1_fabricca.yaml    
kubectl wait --timeout=180s --for=condition=Running fabriccas.hlf.kungfusoftware.es --all
kubectl get fabriccas.hlf.kungfusoftware.es org1-ca -o jsonpath='{.status.tls_cert}' > ${PWD}/tls-cert.pem
export CA_HOST=$(kubectl get fabriccas.hlf.kungfusoftware.es org1-ca -o jsonpath='{.status.host}')
export CA_PORT=$(kubectl get fabriccas.hlf.kungfusoftware.es org1-ca -o jsonpath='{.status.port}')
fabric-ca-client enroll -u "https://enroll:enrollpw@${CA_HOST}:${CA_PORT}" --caname ca --tls.certfiles ${PWD}/tls-cert.pem
fabric-ca-client register -u "https://enroll:enrollpw@${CA_HOST}:${CA_PORT}" \
    --caname ca --tls.certfiles ${PWD}/tls-cert.pem \
    --id.type admin --id.name admin-org1 --id.secret=adminpw
```
 
 ### Deploying a peer
 ```bash
kubectl get fabriccas.hlf.kungfusoftware.es org1-ca -o jsonpath='{.status.tls_cert}'
export CA_TLS_CRT=$(kubectl get fabriccas.hlf.kungfusoftware.es org1-ca -o jsonpath='{.status.tls_cert}')
fabric-ca-client register -u "https://admin:adminpw@${CA_HOST}:${CA_PORT}" \
    --caname ca --tls.certfiles "${PWD}/tls-cert.pem" \
    --id.type peer --id.name peer --id.secret=peerpw

# Set certificates of the Certificate Authority
yq w -i ./config/samples/hlf_v1alpha1_fabricpeer.yaml spec.secret.enrollment.component.cahost "${CA_HOST}" 
yq w -i ./config/samples/hlf_v1alpha1_fabricpeer.yaml spec.secret.enrollment.tls.cahost "${CA_HOST}"
yq w -i ./config/samples/hlf_v1alpha1_fabricpeer.yaml spec.secret.enrollment.component.caport "${CA_PORT}" 
yq w -i ./config/samples/hlf_v1alpha1_fabricpeer.yaml spec.secret.enrollment.tls.caport "${CA_PORT}" 
yq w -i ./config/samples/hlf_v1alpha1_fabricpeer.yaml spec.secret.enrollment.component.catls.cacert "$(echo ${CA_TLS_CRT} | base64)" 
yq w -i ./config/samples/hlf_v1alpha1_fabricpeer.yaml spec.secret.enrollment.component.catls.cacert "$(echo ${CA_TLS_CRT} | base64)" 
yq w -i ./config/samples/hlf_v1alpha1_fabricpeer.yaml spec.secret.enrollment.tls.catls.cacert "$(echo ${CA_TLS_CRT} | base64)"
kubectl apply -f ./config/samples/hlf_v1alpha1_fabricpeer.yaml
kubectl wait --timeout=180s --for=condition=Running fabricpeers.hlf.kungfusoftware.es --all

```

## Deploying an Ordering Service

### Deploying a certificate authority
```bash
kubectl apply -f ./config/samples/hlf_v1alpha1_fabricordererca.yaml    
kubectl wait --timeout=180s --for=condition=Running fabriccas.hlf.kungfusoftware.es --all
kubectl get fabriccas.hlf.kungfusoftware.es orderer-ca -o jsonpath='{.status.tls_cert}' > ${PWD}/tls-cert.pem
export CA_HOST=$(kubectl get fabriccas.hlf.kungfusoftware.es orderer-ca -o jsonpath='{.status.host}')
export CA_PORT=$(kubectl get fabriccas.hlf.kungfusoftware.es orderer-ca -o jsonpath='{.status.port}')
fabric-ca-client enroll -u "https://enroll:enrollpw@${CA_HOST}:${CA_PORT}" --caname ca --tls.certfiles ${PWD}/tls-cert.pem
fabric-ca-client register -u "https://enroll:enrollpw@${CA_HOST}:${CA_PORT}" \
    --caname ca --tls.certfiles ${PWD}/tls-cert.pem \
    --id.type admin --id.name admin-orderer --id.secret=adminpw
```

### Deploying the Ordering service
```bash

kubectl get fabriccas.hlf.kungfusoftware.es orderer-ca -o jsonpath='{.status.tls_cert}'
export CA_TLS_CRT=$(kubectl get fabriccas.hlf.kungfusoftware.es orderer-ca -o jsonpath='{.status.tls_cert}')
fabric-ca-client register -u "https://admin:adminpw@${CA_HOST}:${CA_PORT}" \
    --caname ca --tls.certfiles "${PWD}/tls-cert.pem" \
    --id.type orderer --id.name orderer --id.secret=ordererpw

yq w -i ./config/samples/hlf_v1alpha1_fabricorderer.yaml spec.enrollment.component.cahost "${CA_HOST}" 
yq w -i ./config/samples/hlf_v1alpha1_fabricorderer.yaml spec.enrollment.tls.cahost "${CA_HOST}"
yq w -i ./config/samples/hlf_v1alpha1_fabricorderer.yaml spec.enrollment.component.caport "${CA_PORT}" 
yq w -i ./config/samples/hlf_v1alpha1_fabricorderer.yaml spec.enrollment.tls.caport "${CA_PORT}" 
yq w -i ./config/samples/hlf_v1alpha1_fabricorderer.yaml spec.enrollment.component.catls.cacert "$(echo ${CA_TLS_CRT} | base64)" 
yq w -i ./config/samples/hlf_v1alpha1_fabricorderer.yaml spec.enrollment.tls.catls.cacert "$(echo ${CA_TLS_CRT} | base64)"

kubectl apply -f ./config/samples/hlf_v1alpha1_fabricorderer.yaml
kubectl wait --timeout=180s --for=condition=Running fabricorderingservices.hlf.kungfusoftware.es --all
```

At this point, you should have:
- Ordering service with 3 nodes and a CA
- Peer organization with a peer and a CA

From now on, you should:
- [Create a channel](https://hyperledger-fabric.readthedocs.io/en/release-2.2/create_channel/create_channel_overview.html)
- [Deploy a smart contract](https://hyperledger-fabric.readthedocs.io/en/release-2.2/deploy_chaincode.html)

If something went wrong or didn't work, please, open an issue.

### Cleanup the environment
```bash
kubectl delete -f ./config/samples/hlf_v1alpha1_fabricca.yaml
kubectl delete -f ./config/samples/hlf_v1alpha1_fabricpeer.yaml
kubectl delete -f ./config/samples/hlf_v1alpha1_fabricorderer.yaml
kubectl delete -f ./config/samples/hlf_v1alpha1_fabricordererca.yaml
```


## New release

```bash

make docker-build IMG=quay.io/kfsoftware/hlf-operator:1.0.1
make docker-push IMG=quay.io/kfsoftware/hlf-operator:1.0.0
```
