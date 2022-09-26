---
id: followerchannel
title: Follower channel
---


Second, we create the main channel CRD and apply it.
```bash
kubectl hlf channelcrd follower create \
    --channel-name=demo \
    --mspid=Org1MSP \
    --name="demo-org1msp" \
    --orderer-certificates="./orderer-cert.pem" \
    --orderer-urls="grpcs://ord-node1.default:7050" \
    --anchor-peers="org1-peer0:7051" \
    --peers="org1-peer0.default" \
    --secret-name=wallet \
    --secret-ns=default \
    --secret-key="peer-org1.yaml"
```
