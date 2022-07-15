---
id: architecture
title: Architecture
---



## Problem

When developing chaincodes, we need to have an HLF network deployed, usually in our development environment. It is hard for new developers to get started since they need to know how to deploy an HLF network.

This tool aims to ease the onboarding of new developers that are not familiar with the internals of the HLF network but are interested in developing chaincodes.

The only requirement is that the developer has access to a working HLF network. This network can be set up by Administrators that are used to perform these operations, so the developer doesn't need to.

With this tool, instead of installing the chaincode in the peers, approving the chaincode definition, and finally, committing the chaincode, the developer can have the chaincode started in its machine. With one command, it can install, approve and commit the chaincode in one go. SupposeThen, ifhe developer needs to modify the chaincode logic. In that case, all it needs to do is restart the chaincode program running in its machine, just like any other application the developer is used to developing.



## Solution

The chaincode needs to be hosted in the developer machine, by doing this the developer can modify the chaincode logic without having to re-deploy the chaincode program.

But the developer machine is not accessible by the peer, which is deployed in another location, for this reason, we need a tunnel to be able to specify an address for the peer to connect and then, the tunnel will forward the traffic to the developer machine.

[`hlf-cc-dev`](https://github.com/kfsoftware/hlf-cc-dev) is a tool that helps to deploy the chaincode program using a simple CLI interface.

For more inforamtion you can check the following diagram.

![img.png](/img/arch_chaincode_dev.png)
