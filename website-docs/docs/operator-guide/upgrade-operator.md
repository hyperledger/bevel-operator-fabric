---
id: upgrade-hlf-operator
title: Upgrade HLF operator
---

When there's a new release, the following resources may be added or modified:
- Custom resource definitions
- Cluster Role 
- Deployment

In order to upgrade the hlf-operator, you need to execute the following command:

```bash
export NEW_VERSION=1.7.0
helm install hlf-operator --version=$NEW_VERSION kfs/hlf-operator
```

If you specified a `values.yaml`, you'll need to pass the values to the upgrade command:

```bash
export NEW_VERSION=1.7.0
helm upgrade hlf-operator --values=values.yaml --version=$NEW_VERSION kfs/hlf-operator
```


After upgrading the operator, make sure it starts and there are no errors, in case there are and you don't know how to fix it, please, open an [issue in Github](https://github.com/hyperledger/bevel-operator-fabric/issues/new)
