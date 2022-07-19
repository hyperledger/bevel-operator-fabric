---
id: k8s-builder
title: Kubernetes chaincode builder
---

## Architecture

In order to build the chaincode, you have to use an image that has a Kubernetes builder bundled with it.


The following picture illustrates how the Kubernetes builder bundled with the previous images works:

![](/img/kubernetes_builder_chaincode.png)

## Creating peers with external builder

To use the Kubernetes builder, you have to create a peer with the following command:

```bash
export PEER_IMAGE=quay.io/kfsoftware/fabric-peer
export PEER_VERSION=2.4.1-v0.0.3
export ORG=org1
export MSP_ORG=Org1MSP
export CA_NAME=ca-org1.default
export PEER_SECRET=peerpw

kubectl hlf peer create --image=$PEER_IMAGE --version=$PEER_VERSION --storage-class=$STORAGE_CLASS --enroll-id=peer --mspid=$MSP_ORG \
        --enroll-pw=$PEER_SECRET --capacity=5Gi --name=$ORG-peer1 --ca-name=$CA_NAME --k8s-builder=true
```

If the peer is already created, you need to change the following properties:
```yaml
external_chaincode_builder: true
```
Add the following item to the `externalBuilders` property in the `spec` section:
```yaml
  - name: k8s-builder
    path: /builders/golang
    propagateEnvironment:
    - CHAINCODE_SHARED_DIR
    - FILE_SERVER_BASE_IP
    - KUBERNETES_SERVICE_HOST
    - KUBERNETES_SERVICE_PORT
    - K8SCC_CFGFILE
    - TMPDIR
    - LD_LIBRARY_PATH
    - LIBPATH
    - PATH
    - EXTERNAL_BUILDER_HTTP_PROXY
    - EXTERNAL_BUILDER_HTTPS_PROXY
    - EXTERNAL_BUILDER_NO_PROXY
    - EXTERNAL_BUILDER_PEER_URL
```

And finally, update the image and version to the following:

```yaml
  image: quay.io/kfsoftware/fabric-peer
  tag: 2.4.1-v0.0.3
```

## Install chaincode

Install chaincode using the Kubernetes builder:
```bash
kubectl hlf chaincode install --path=./fixtures/chaincodes/fabcar/go \
    --config=org1.yaml --language=golang --label=fabcar --user=admin --peer=org1-peer0.default

# this can take 3-4 minutes
```