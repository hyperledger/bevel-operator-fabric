---
id: network-config
title: Get a network config
---

Generating a network config is one of the most common operations once you have a network up and running.

You can get a network config with the following command:

```bash
kubectl hlf inspect --output networkConfig.yaml -o OrdererMSP -o Org1MSP
```

Network config will look like this:

```yaml

```

```bash

kubectl hlf ca register --name=ca-kfsmsp --namespace=hlf --user=admin --secret=adminpw --mspid KFSMSP

kubectl hlf ca register --name=ca-kfsmsp --namespace=hlf --user=admin --secret=adminpw --type=admin \
 --enroll-id enroll --enroll-secret=enrollpw --mspid KFSMSP
```