---
id: deploy-operator-api
title: Deploy Operator API
---



## Create operator API
In order to create the operator API:

```bash
export API_URL=api-operator.<domain>
kubectl hlf operatorapi create --name=operator-api --namespace=default --hosts=$API_URL --ingress-class-name=istio --network-config=./network-config.yaml
```

## Delete operator API
In order to delete the operator API:

```bash
kubectl hlf operatorapi delete --name=operator-api --namespace=default
```
