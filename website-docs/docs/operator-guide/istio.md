---
id: istio
title: Istio set up
---

Istio is a service mesh that provides a secure, high-performance networking platform for microservices and applications running on Kubernetes.

Node port solutions can work in the short term, but they are not long-term solutions, neither for production, as they require opening up multiple ports to the public internet.


The following diagram represents the architecture with Istio configured
![Istio](/img/istio_hlf.png)

As you can see, we can note the following:
- The Istio service mesh is running in the Kubernetes cluster
- The service only has one port exposed, which is the port of the Istio ingress gateway service.
- The ingress gateway routes the traffic to the peer, OSN or CA depending on the request.

## Installing istio


You can refer to your version of choice by going to this tutorial from the [istio docs](https://istio.io/latest/docs/setup/getting-started/) to get Istio installed in your Kubernetes cluster.


Alternatively, you can just execute this command to install the latest Istio version in your Kubernetes cluster:

```bash
curl -L https://istio.io/downloadIstio | sh - # download istioctl CLI

istioctl install --set profile=default -y # install Istio
```

## Locate the public IP or hostname of the ingress gateway

### Running on KinD/Minikube
```bash
PUBLIC_IP=$(kubectl get nodes -o jsonpath='{.items[0].status.addresses[?(@.type=="InternalIP")].address}')
# get node port

PORT=$(kubectl get svc istio-ingressgateway -n istio-system -o jsonpath='{.spec.ports[?(@.name=="https")].nodePort}')
```

### Running with load balancer IP

```bash
PUBLIC_IP=$(kubectl get svc istio-ingressgateway -n istio-system -o json | jq -r '.status.loadBalancer.ingress[0].ip')
PORT=443
```

### Running with load balancer hostname
```bash
PUBLIC_HOSTNAME=$(kubectl get svc istio-ingressgateway -n istio-system -o json | jq -r '.status.loadBalancer.ingress[0].hostname')
PORT=443
```


## Set up DNS

### Local DNS set up in Linux/Mac (for MiniKube and KinD)

Open up /etc/hosts
```bash
<PUBLIC_IP> peer0.org1.example.com
<PUBLIC_IP> ord1.ord-org.example.com
# and so on
```

### Set up DNS in your DNS provider

You will need to point the domain names you will use to the public IP of the ingress gateway, with either a A record, if you got a public IP, or a CNAME, if you got an ingress hostname


## Set up the network

### Deploying a Certificate Authority

```bash
kubectl hlf ca create --storage-class=standard --capacity=2Gi --name=org1-ca \
    --enroll-id=enroll --enroll-pw=enrollpw  
kubectl wait --timeout=180s --for=condition=Running fabriccas.hlf.kungfusoftware.es --all

# register user for the peers
kubectl hlf ca register --name=org1-ca --user=peer --secret=peerpw --type=peer \
 --enroll-id enroll --enroll-secret=enrollpw --mspid Org1MSP
```

### Create the peer

```bash
PEER1_DOMAIN=peer0.org1.example.com # domain for the peer
ISTIO_INGRESSGATEWAY=ingressgateway # name of the ingress gateway, in case there are many
ISTIO_GW_PORT=443 # port of the ingress gateway
kubectl hlf peer create --storage-class=standard --enroll-id=peer --mspid=Org1MSP \
        --enroll-pw=peerpw --capacity=5Gi --name=org1-peer0 --ca-name=org1-ca.default \
        --hosts=$PEER1_DOMAIN --istio-ingressgateway=$ISTIO_INGRESSGATEWAY --istio-port=$ISTIO_GW_PORT

kubectl wait --timeout=180s --for=condition=Running fabricpeers.hlf.kungfusoftware.es --all
```

If we inspect the virtual services and gateways of Istio, we must see a record per peer.

```bash
kubectl get virtualservices.networking.istio.io -A  # list all virtual services
kubectl get gateways.networking.istio.io -A  # list all gateways
```

To test that you can connect to the peer, you can use the following command to test directly from the command line(this test doesn't require DNS records to be set up):
```bash
echo "PUBLIC_IP=$PUBLIC_IP PORT=$PORT DOMAIN=$PEER1_DOMAIN"
openssl s_client -connect $PUBLIC_IP:$PORT -servername $PEER1_DOMAIN  -showcerts </dev/null
```


### Create the Certificate Authority for the orderer

```bash
kubectl hlf ca create --storage-class=standard --capacity=2Gi --name=ord-ca \
    --enroll-id=enroll --enroll-pw=enrollpw
kubectl wait --timeout=180s --for=condition=Running fabriccas.hlf.kungfusoftware.es --all
kubectl hlf ca register --name=ord-ca --user=orderer --secret=ordererpw \
    --type=orderer --enroll-id enroll --enroll-secret=enrollpw --mspid=OrdererMSP
```

### Deploying the Orderer nodes node

```bash
ORD1_DOMAIN=ord1.org1-node.example.com # domain for the orderer
ISTIO_INGRESSGATEWAY=ingressgateway # name of the ingress gateway, in case there are many
ISTIO_GW_PORT=443
kubectl hlf ordnode create --storage-class=standard --enroll-id=orderer --mspid=OrdererMSP \
    --enroll-pw=ordererpw --capacity=2Gi --name=ord-node1 --ca-name=ord-ca.default \
    --hosts=$ORD1_DOMAIN --istio-ingressgateway=$ISTIO_INGRESSGATEWAY --istio-port=$ISTIO_GW_PORT

kubectl wait --timeout=180s --for=condition=Running fabricorderernodes.hlf.kungfusoftware.es --all
```


### Testing orderer node connection
```bash
echo "PUBLIC_IP=$PUBLIC_IP PORT=$PORT DOMAIN=$ORD1_DOMAIN"
openssl s_client -connect $PUBLIC_IP:$PORT -servername $ORD1_DOMAIN  -showcerts </dev/null
```

### Create the Certificate Authority with Istio

```bash
CA_ORG2=ca.org2.example.com # domain for the orderer
ISTIO_INGRESSGATEWAY=ingressgateway # name of the ingress gateway, in case there are many
ISTIO_GW_PORT=443

kubectl hlf ca create --storage-class=standard --capacity=2Gi --name=org2-ca \
    --enroll-id=enroll --enroll-pw=enrollpw \
    --hosts=$CA_ORG2 --istio-ingressgateway=$ISTIO_INGRESSGATEWAY --istio-port=$ISTIO_GW_PORT

kubectl wait --timeout=180s --for=condition=Running fabriccas.hlf.kungfusoftware.es --all
```


### Testing CA node connection
```bash
echo "PUBLIC_IP=$PUBLIC_IP PORT=$PORT DOMAIN=$CA_ORG2"
openssl s_client -connect $PUBLIC_IP:$PORT -servername $CA_ORG2  -showcerts </dev/null
```
