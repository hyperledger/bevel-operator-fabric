---
id: configuration
title: Configuration (Affinity, NodeSelector, Tolerations)
---

## Set Affinity

### Set affinity for the FabricCA

```bash
export CA_NAME=org1-ca
export CA_NS=default
cat <<EOT > affinity-patch.yaml
spec:
  affinity:
    nodeAffinity:
      requiredDuringSchedulingIgnoredDuringExecution:
        nodeSelectorTerms:
        - matchExpressions:
          - key: kubernetes.io/e2e-az-name
            operator: In
            values:
            - e2e-az1
            - e2e-az2
      preferredDuringSchedulingIgnoredDuringExecution:
      - weight: 1
        preference:
          matchExpressions:
          - key: another-node-label-key
            operator: In
            values:
            - another-node-label-value
EOT

kubectl patch fabriccas.hlf.kungfusoftware.es $CA_NAME --namespace=$CA_NS --patch="$(cat affinity-patch.yaml)" --type=merge

```


### Set affinity for the FabricPeer

```bash
export PEER_NAME=org1-peer0
export PEER_NS=default
cat <<EOT > affinity-patch.yaml
spec:
  affinity:
    nodeAffinity:
      requiredDuringSchedulingIgnoredDuringExecution:
        nodeSelectorTerms:
        - matchExpressions:
          - key: kubernetes.io/e2e-az-name
            operator: In
            values:
            - e2e-az1
            - e2e-az2
      preferredDuringSchedulingIgnoredDuringExecution:
      - weight: 1
        preference:
          matchExpressions:
          - key: another-node-label-key
            operator: In
            values:
            - another-node-label-value
EOT

kubectl patch fabricpeers.hlf.kungfusoftware.es $PEER_NAME --namespace=$PEER_NS --patch="$(cat affinity-patch.yaml)" --type=merge
```


### Set affinity for the FabricOrdererNode

```bash
export ORDERER_NAME=org1-peer0
export ORDERER_NS=default
cat <<EOT > affinity-patch.yaml
spec:
  affinity:
    nodeAffinity:
      requiredDuringSchedulingIgnoredDuringExecution:
        nodeSelectorTerms:
        - matchExpressions:
          - key: kubernetes.io/e2e-az-name
            operator: In
            values:
            - e2e-az1
            - e2e-az2
      preferredDuringSchedulingIgnoredDuringExecution:
      - weight: 1
        preference:
          matchExpressions:
          - key: another-node-label-key
            operator: In
            values:
            - another-node-label-value
EOT

kubectl patch fabricorderernodes.hlf.kungfusoftware.es $ORDERER_NAME --namespace=$ORDERER_NS --patch="$(cat affinity-patch.yaml)" --type=merge
```


## Set tolerations


### Set tolerations for the FabricCA

```bash
export CA_NAME=org1-ca
export CA_NS=default
cat <<EOT > tolerations-patch.yaml
spec:
  tolerations:
    - effect: NoSchedule
      key: kubernetes.azure.com/scalesetpriority
      operator: Equal
      value: spot
EOT

kubectl patch fabriccas.hlf.kungfusoftware.es $CA_NAME --namespace=$CA_NS --patch="$(cat tolerations-patch.yaml)" --type=merge

```


### Set tolerations for the FabricPeer

```bash
export PEER_NAME=org1-peer0
export PEER_NS=default
cat <<EOT > tolerations-patch.yaml
spec:
  tolerations:
    - effect: NoSchedule
      key: kubernetes.azure.com/scalesetpriority
      operator: Equal
      value: spot
EOT

kubectl patch fabricpeers.hlf.kungfusoftware.es $PEER_NAME --namespace=$PEER_NS --patch="$(cat tolerations-patch.yaml)" --type=merge
```


### Set tolerations for the FabricOrdererNode

```bash
export ORDERER_NAME=org1-peer0
export ORDERER_NS=default
cat <<EOT > tolerations-patch.yaml
spec:
  tolerations:
    - effect: NoSchedule
      key: kubernetes.azure.com/scalesetpriority
      operator: Equal
      value: spot
EOT

kubectl patch fabricorderernodes.hlf.kungfusoftware.es $ORDERER_NAME --namespace=$ORDERER_NS --patch="$(cat tolerations-patch.yaml)" --type=merge
```


## Set Node Selector

### Set nodeselector for the FabricCA

```bash
export CA_NAME=org1-ca
export CA_NS=default
cat <<EOT > nodeselector-patch.yaml
spec:
  nodeSelector:
    disktype: ssd
EOT

kubectl patch fabriccas.hlf.kungfusoftware.es $CA_NAME --namespace=$CA_NS --patch="$(cat nodeselector-patch.yaml)" --type=merge

```


### Set nodeselector for the FabricPeer

```bash
export PEER_NAME=org1-peer0
export PEER_NS=default
cat <<EOT > nodeselector-patch.yaml
spec:
  nodeSelector:
    disktype: ssd
EOT

kubectl patch fabricpeers.hlf.kungfusoftware.es $PEER_NAME --namespace=$PEER_NS --patch="$(cat nodeselector-patch.yaml)" --type=merge
```


### Set nodeselector for the FabricOrdererNode

```bash
export ORDERER_NAME=org1-peer0
export ORDERER_NS=default
cat <<EOT > nodeselector-patch.yaml
spec:
  nodeSelector:
    disktype: ssd
EOT

kubectl patch fabricorderernodes.hlf.kungfusoftware.es $ORDERER_NAME --namespace=$ORDERER_NS --patch="$(cat nodeselector-patch.yaml)" --type=merge
```

