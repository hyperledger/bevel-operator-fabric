---
id: getting-started
title: Getting started
---


The HLF Operator UI provides a graphical interface for a convenient blockchain-as-a-service user experience. Once the operator is set up it's very easy for teams to create, clone, watch, edit and delete their own Certificate Authorities, Peers and Orderer nodes.


## Why another UI for Fabric



The HLF Operator UI consists of two components:
- hlf-operator-ui
- hlf-operator-api

![hlf operator ui](/img/hlf_operator_ui.png)

## HLF Operator UI

The goal of the HLF Operator UI is to provides a graphical interface for a convenient blockchain-as-a-service user experience by:
- Creating peers
- Creating CAs
- Creating orderers
- Renewing certificates

Apart from this, it will be an explorer for the network config provided to the HLF Operator API to see:
- Channels
- Channel details
  + Height
  + Organizations
  + Peers + height

## HLF Operator API

It provides access to the data consumed by the HLF Operator UI:
- Channels
- Peers
- Orderer Nodes
- Certificate authorities

## Deploying

You need to deploy the HLF Operator UI and the HLF Operator API separately:

- [Deploy the HLF Operator UI](./deploy-operator-ui.md)
- [Deploy the HLF Operator API](./deploy-operator-api.md)
