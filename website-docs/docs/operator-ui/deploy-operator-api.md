---
id: deploy-operator-api
title: Deploy Operator API
---



## Create operator API
In order to create the operator API:

```bash
export API_URL=api-operator.<domain>
kubectl hlf operatorapi create --name=operator-api --namespace=default --hosts=$API_URL --ingress-class-name=istio
```

## Create operator API with authentication
```bash
export API_URL=api-operator.<domain>
export OIDC_ISSUER=https://<your_oidc_issuer>
export OIDC_JWKS=https://<oidc_jwks_url>
kubectl hlf operatorapi create --name=operator-api --namespace=default --hosts=$API_HOST --ingress-class-name=istio \
    --oidc-issuer="${OIDC_ISSUER}" --oidc-jwks="${OIDC_JWKS}"
```

## Create operator API with explorer

```bash
export API_URL=api-operator.<domain>
export HLF_SECRET_NAME="k8s-secret"
export HLF_MSPID="<your_mspid>"
export HLF_SECRET_KEY="<network_config_key_secret>" # e.g. networkConfig.yaml
export HLF_USER="<hlf_user>"
kubectl hlf operatorapi create --name=operator-api --namespace=default --hosts=$API_HOST --ingress-class-name=istio \
          --hlf-mspid="${HLF_MSPID}" --hlf-secret="${HLF_SECRET_NAME}" --hlf-secret-key="${HLF_SECRET_KEY}" \
          --hlf-user="${HLF_USER}"
```

## Update operator API

You can use the same commands with the same parameters, but instead of `create` use `update`

## Delete operator API
In order to delete the operator API:

```bash
kubectl hlf operatorapi delete --name=operator-api --namespace=default
```
