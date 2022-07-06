---
id: adding-cas
title: Adding Certificate Authorities
---



```bash
export PEER_NAME=xxxxx
export PEER_NS=default
kubectl hlf fop export peer --name=$PEER_NAME --namespace=$PEER_NS --out="${PEER_NAME}_${PEER_NS}.json"
```

