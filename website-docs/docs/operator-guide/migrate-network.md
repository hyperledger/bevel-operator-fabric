---
id: migrate-network
title: Migrate network
---

This document is a walkthrough of the steps required to migrate a network from another method to the HLF operator.

## Peer migration

The best migration is to perform migration at all. For the peers this is possible by spinning up a new set of peers with new domains.

In order to do this, we must create the certificate authority with the Keys and Certificates that were used to create the previous peer certificates.

## Ordering service migration

For the ordering service, a new set of orderer nodes must be created with new domains. In order to do this, we must create the certificate authority with the Keys and Certificates that were used to create the previous peer certificates.

For example, if we have an existing ordering service with the following URLs:
- orderer1.myorg.com:7050
- orderer2.myorg.com:7050
- orderer3.myorg.com:7050
- orderer4.myorg.com:7050
- orderer5.myorg.com:7050


We must create with the HLF operator the following orderer nodes:
- orderer6.myorg.com:7050
- orderer7.myorg.com:7050
- orderer8.myorg.com:7050
- orderer9.myorg.com:7050
- orderer10.myorg.com:7050


After these nodes are created, we must join them to the channel, when the orderer nodes are joined to the channel and we're confident enough about the new set of orderers, progressively, we must update the channel configuration to include the new orderers as consenters, and remove the old ones.
