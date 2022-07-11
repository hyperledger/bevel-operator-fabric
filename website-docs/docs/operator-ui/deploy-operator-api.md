---
id: deploy-operator-api
title: Deploy Operator API
---



## Create operator API
In order to create

```bash
export API_URL=api-operator.hlf.kfs.es
kubectl hlf operatorapi create --name=operator-api --namespace=default --hosts=$API_URL --ingress-class-name=istio
```

## Delete operator API
In order to delete:
```bash
kubectl hlf operatorapi delete --name=operator-api --namespace=default
```
