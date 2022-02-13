# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).


## [1.5.0]

### Features
- Improved documentation for how to choose the World State database.
- Create network config in the cluster using a CRD
- Add an example client application for Node.JS
- Configure the `image`, `tag`, and `pullPolicy` for the CouchDB container
- Configure the `image`, `tag`, and `pullPolicy` File Server container (used to build chaincodes in Kubernetes).

### Kubectl plugin
- Allow passing empty signature policy to the `approve` and `commit` chaincode command takes the default signature policy of the channel.
- Add `networkconfig` commands to create and update the network config CRD that generates secret in Kubernetes so the application can use it.

## [1.4.0]

- Renewal of certificates
- Improved documentation for Istio
- Support for Prometheus metrics using Kube Prometheus Operator