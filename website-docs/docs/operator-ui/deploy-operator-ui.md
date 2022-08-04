---
id: deploy-operator-ui
title: Deploy Operator UI
---


## Create operator UI
In order to create the operator UI:

```bash
export HOST=operator-ui.<domain>
export API_URL="http://api-operator.<domain>/graphql"
kubectl hlf operatorui create --name=operator-ui --namespace=default --hosts=$HOST --ingress-class-name=istio --api-url=$API_URL
```

## Create operator UI with authentication

```bash
```

## Delete operator UI
In order to delete the operator UI:

```bash
kubectl hlf operatorui delete --name=operator-ui --namespace=default
```

