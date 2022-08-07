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
export HOST=operator-ui.<domain>
export API_URL="http://api-operator.<domain>/graphql"
export OIDC_AUTHORITY="<url_authority>" # without the /.well-known/openid-configuration
export OIDC_CLIENT_ID="<client_id>" # OIDC Client ID for the Operator UI
export OIDC_SCOPE="profile email" # OIDC Scope for the Operator UI
kubectl hlf operatorui create --name=operator-ui --namespace=default --hosts=$HOST --ingress-class-name=istio --api-url=$API_URL \
      --oidc-authority="${OIDC_AUTHORITY}" --oidc-client-id="${OIDC_CLIENT_ID}" --oidc-scope="${OIDC_SCOPE}"         
```

## Update operator API

You can use the same commands with the same parameters, but instead of `create` use `update`

## Delete operator UI
In order to delete the operator UI:

```bash
kubectl hlf operatorui delete --name=operator-ui --namespace=default
```

