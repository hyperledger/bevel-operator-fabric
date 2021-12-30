---
id: increase-storage
title: Increase storage
---

## Increase storage for the peer
```bash
kubectl hlf peer upgrade-storage \
    --name=peer1 --namespace=default \
    --storage-size=10Gi
```

## Increase storage for the orderer
```bash
kubectl hlf orderer upgrade-storage \
    --name=orderer1 --namespace=default \
    --storage-size=10Gi
```

## Increase storage for the certificate authority
```bash
kubectl hlf peer upgrade-storage \
    --name=peer1 --namespace=default \
    --storage-size=10Gi
```

## Increase storage for the CouchDB

```bash
kubectl hlf peer upgrade-storage \
    --name=peer1 --namespace=default \
    --storage-size=10Gi
```
