---
id: network-config
title: Get a network config
---

Generating a network config is one of the most common operations once you have a network up and running.

### Generate network config

You can get a network config with the following command:

```bash
kubectl hlf inspect --output networkConfig.yaml -o OrdererMSP -o Org1MSP
```

Network config will look like this:

```yaml
name: hlf-network
version: 1.0.0
client:
  organization: ""
... rest of your network config ...
```

In order to have users in your network, first you need to register and enroll them:


### Setup env variables
```bash
CA_NAME=ca-org1
CA_NAMESPACE=default
MSP_ID=Org1MSP

USER_NAME=admin
USER_PWD=adminpw
USER_TYPE=admin # it can be client, admin, peer, orderer
ENROLL_USER=enroll # username of the enroll user of the CA
ENROLL_PWD=enrollpw # username of the enroll password of the CA

USER_CA_TYPE=ca # it can be ca, tlsca
```
### Register a user

```bash
kubectl hlf ca register --name=$CA_NAME \
 --namespace=$CA_NAMESPACE --mspid=$MSP_ID \
 --user=$USER_NAME --secret=$USER_PWD --type=$USER_TYPE \
 --enroll-id=$ENROLL_USER --enroll-secret=$ENROLL_PWD
```
If it has been already registered, the following error will appear
```log
Error: failed to register user: failed to register user: Response from server: Error Code: 74 - Identity 'admin' is already registered
```

### Enroll a user
```bash
kubectl hlf ca enroll --name=$CA_NAME --namespace=$CA_NAMESPACE \
    --user=$USER_NAME --secret=$USER_PWD --mspid $MSP_ID \
    --ca-name=$USER_CA_TYPE  --output user.yaml
```

### Utility: add user to network config

```bash
kubectl hlf inspect --output org1.yaml -o Org1MSP -o OrdererMSP

## add user key and cert to org1.yaml from peer-org1.yaml

kubectl hlf utils adduser --userPath=user.yaml \
  --config=org1.yaml --username=admin --mspid=$MSP_ID
```
