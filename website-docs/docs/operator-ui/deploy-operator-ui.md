---
id: deploy-operator-ui
title: Deploy Operator UI
---


## Create operator UI
In order to create

```bash
export HOST=operator-ui.hlf.kfs.es
export API_URL="http://api-operator.hlf.kfs.es/graphql"
kubectl hlf operatorui create --name=operator-ui --namespace=default --hosts=$HOST --ingress-class-name=istio --api-url=$API_URL
```

## Delete operator UI
In order to delete:
```bash
kubectl hlf operatorui delete --name=operator-ui --namespace=default
```

