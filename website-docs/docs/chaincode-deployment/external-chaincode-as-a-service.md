---
id: external-chaincode-as-a-service
title: External chaincode as a service
---

Fabric 2.0 enabled the ability to deploy chaincode as a service. This is a chaincode that runs on a remote machine, container, or baremetal, and the peer will connect to it.

HLF-operator supports the installation of chaincode as a service, using the `ccaas` external builder that Hyperledger Fabric supports out-of-the-box since the version `2.4.1`

![External chaincode as a service](/img/external_chaincode_as_a_service.png)


## Step 1: Prepare a Docker image for the chaincode

To deploy the chaincode in Kubernetes, we need to prepare a Docker image that has the chaincode in it.

You can use the following code as a baseline:
[https://github.com/hyperledger/fabric-samples/tree/main/chaincode/fabcar/external](https://github.com/hyperledger/fabric-samples/tree/main/chaincode/fabcar/external)

It contains a `Dockerfile` and a `fabcar.go`, if you build that folder:

```bash
docker build -t <username>/chaincode .
```

And then push it to Docker Hub:

```bash
docker push <username>/chaincode
```

You will be able to deploy the chaincode in the step 3

## Step 2: Install the chaincode

In this step we will install the chaincode using the `ccaas` external builder, telling the address of the chaincode service where the chaincode will be deployed.

```bash
# remove the code.tar.gz asset-transfer-basic-external.tgz if they exist
rm code.tar.gz asset-transfer-basic-external.tgz
export CHAINCODE_NAME=asset
export CHAINCODE_LABEL=asset
cat << METADATA-EOF > "metadata.json"
{
    "type": "ccaas",
    "label": "${CHAINCODE_LABEL}"
}
METADATA-EOF

cat > "connection.json" <<CONN_EOF
{
  "address": "${CHAINCODE_NAME}:7052",
  "dial_timeout": "10s",
  "tls_required": false
}
CONN_EOF

tar cfz code.tar.gz connection.json
tar cfz asset-transfer-basic-external.tgz metadata.json code.tar.gz
export PACKAGE_ID=$(kubectl hlf chaincode calculatepackageid --path=asset-transfer-basic-external.tgz --language=node --label=$CHAINCODE_LABEL)
echo "PACKAGE_ID=$PACKAGE_ID"

kubectl hlf chaincode install --path=./asset-transfer-basic-external.tgz \
    --config=org1.yaml --language=golang --label=$CHAINCODE_LABEL --user=admin --peer=org1-peer0.default

# this can take 3-4 minutes
```

## Step 3: Deploy the chaincode

In this last step, we will deploy the chaincode using the `ccaas` external builder, specifying the name of the chaincode, which must be the same as we specified in the previous step.

And we will also specify the image which you have pushed in the first step.

```bash
kubectl hlf externalchaincode sync --image=<username>/chaincode:latest \
    --name=$CHAINCODE_NAME \
    --namespace=default \
    --package-id=$PACKAGE_ID \
    --tls-required=false \
    --replicas=1
```


