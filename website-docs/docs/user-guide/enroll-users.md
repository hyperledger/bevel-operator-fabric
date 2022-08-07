---
id: enroll-users
title: Register & Enroll users
---


## Registering users

```bash
ENROLL_ID=enroll # enroll id for the CA, default `enroll`
ENROLL_SECRET=enrollpw # enroll secret for the CA, default `enrollpw`
USER_TYPE=peer # can be `peer`, `orderer`, `client` or `admin`
USER_NAME=peer
USER_SECRET=peerpw
MSP_ID=Org1MSP
kubectl hlf ca register --name=$CA_NAME --namespace=$CA_NAMESPACE \
    --user $USER --secret=$USER_SECRET --type=$USER_TYPE \
    --enroll-id=$ENROLL_ID --enroll-secret=$ENROLL_SECRET \
    --mspid $MSP_ID
```

## Registering users with attributes

```bash
ENROLL_ID=enroll # enroll id for the CA, default `enroll`
ENROLL_SECRET=enrollpw # enroll secret for the CA, default `enrollpw`
USER_TYPE=peer # can be `peer`, `orderer`, `client` or `admin`
USER_NAME=peer
USER_SECRET=peerpw
MSP_ID=Org1MSP
kubectl hlf ca register --name=$CA_NAME --namespace=$CA_NAMESPACE \
    --user $USER --secret=$USER_SECRET --type=$USER_TYPE \
    --enroll-id=$ENROLL_ID --enroll-secret=$ENROLL_SECRET \
    --mspid $MSP_ID --attributes="isAdmin=true,anotherAttribute=foo"
```

## Enrolling users in the TLS CA

```bash
CA_NAME=org1-ca
CA_NAMESPACE=default
CA_MSPID=Org1MSP
CA_TYPE=ca # can be `ca` or `tlsca`
kubectl hlf ca enroll --name=$CA_NAME --namespace=$CA_NAMESPACE \
    --user=admin --secret=adminpw --mspid $CA_MSPID \
    --ca-name $CA_TYPE  --output user.yaml 
```


## Enrolling users in the Sign CA

```bash
CA_NAME=org1-ca
CA_NAMESPACE=default
CA_MSPID=Org1MSP
CA_TYPE=tlsca # can be `ca` or `tlsca`
kubectl hlf ca enroll --name=$CA_NAME --namespace=$CA_NAMESPACE \
    --user=admin --secret=adminpw --mspid $CA_MSPID \
    --ca-name $CA_TYPE  --output user.yaml 
```


## Enrolling users in the Sign CA with attributes

```bash
CA_NAME=org1-ca
CA_NAMESPACE=default
CA_MSPID=Org1MSP
CA_TYPE=tlsca # can be `ca` or `tlsca`
kubectl hlf ca enroll --name=$CA_NAME --namespace=$CA_NAMESPACE \
    --user=admin --secret=adminpw --mspid $CA_MSPID \
    --ca-name $CA_TYPE  --output user.yaml --attributes="isAdmin,anotherAttribute:opt" # for optional attributes
```
