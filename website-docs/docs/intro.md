---
id: intro
title: Introduction
sidebar_label: Introduction
slug: /
---

## What's HLF Operator?
HLF Operator is a Kubernetes Operator built with the [operator sdk](https://sdk.operatorframework.io/) to manage the Hyperledger Fabric components:
- Peer
- Ordering service nodes(OSN)
- Certificate authorities

You can see more info in https://google.es


## Why another tool to manage Hyperledger Fabric networks?
There are some alternatives such as:
- [Cello](https://github.com/hyperledger/cello)
- [Workflow based on Helm Charts and ArgoCD workflows](https://github.com/hyfen-nl/PIVT)

These tools are much complex, since it requires a deep knowledge in Hyperledger Fabric in order to get the most from these tools.

Instead, what if we could we get the simplicity of Kubernetes and the power from Hyperledger Fabric? This is when this operator comes in. With CRDs(Custom resource definition) for the Peer, Certificate Authority and Ordering Services we can set up a fully network.
