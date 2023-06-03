---
id: implementation
title: Implementation
---

With the GatewayClass and Gateway resources of your respective proxy setup, let's create the fabric resources like CA, peer and orderer.

Note that this setup is similiar to the original setup given in the docs where we are using coredns to resolve the ip addresses. The gateway api implementation also works externally by making the gateway-api service a LoadBalancer.

The first step is to get the address of the gateway which needs to be resolved for the fabric resources.

For istio:

```bash
export INGRESS_HOST=$(kubectl get gateways.gateway.networking.k8s.io gateway -n istio-ingress -ojsonpath='{.status.addresses[*].value}')
```

For traefik, the ingress host is the ClusterIP of the traefik-service which is deployed earlier in the setup

```bash
export INGRESS_HOST=$(kubectl get svc traefik -n gateway-api -o jsonpath='{.status.loadBalancer.ingress[0].ip}')
```

```bash
kubectl apply -f - <<EOF
kind: ConfigMap
apiVersion: v1
metadata:
  name: coredns
  namespace: kube-system
data:
  Corefile: |
    .:53 {
        errors
        health {
           lameduck 5s
        }
        rewrite name regex (.*)\.localho\.st host.ingress.internal
        hosts {
          ${INGRESS_HOST} host.ingress.internal
          fallthrough
        }
        ready
        kubernetes cluster.local in-addr.arpa ip6.arpa {
           pods insecure
           fallthrough in-addr.arpa ip6.arpa
           ttl 30
        }
        prometheus :9153
        forward . /etc/resolv.conf {
           max_concurrent 1000
        }
        cache 30
        loop
        reload
        loadbalance
    }
EOF
```

## Setup

```bash
export PEER_IMAGE=hyperledger/fabric-peer
export PEER_VERSION=2.4.6

export ORDERER_IMAGE=hyperledger/fabric-orderer
export ORDERER_VERSION=2.4.6

export CA_IMAGE=hyperledger/fabric-ca
export CA_VERSION=1.5.6-beta2

export NAMESPACE=hlf
export GATEWAYNAME=gateway  
export GATEWAYNAMESPACE=istio-ingress  

```

Watch out for the following configuration:

--gateway-api-hosts : The hosts that are configured to be used with gateway-api
--gateway-api-name : The name of the gateway (Name of the 'Gateway' Resource created earlier)
--gateway-api-namespace : The namespace where the 'Gateway' resource is deployed


### Create CA

```bash
kubectl hlf ca create  --image=$CA_IMAGE --version=$CA_VERSION --storage-class=standard --capacity=1Gi --name=org1-ca     --enroll-id=enroll --enroll-pw=enrollpw --gateway-api-hosts=org1-ca.localho.st --gateway-api-name $GATEWAYNAME --gateway-api-namespace $GATEWAYNAMESPACE -n $NAMESPACE
```

Make sure the CA is reachable and gives a response

```bash
curl -k https://org1-ca.localho.st:443/cainfo
```


### Create Peers

```bash
 kubectl hlf peer create --statedb=couchdb --image=$PEER_IMAGE --version=$PEER_VERSION --storage-class=standard --enroll-id=peer --mspid=Org1MSP \
        --enroll-pw=peerpw --capacity=5Gi --name=org1-peer0 --ca-name=org1-ca.$NAMESPACE \
        --gateway-api-hosts=peer0-org1.localho.st --gateway-api-name $GATEWAYNAME --gateway-api-namespace $GATEWAYNAMESPACE -n $NAMESPACE
```

Make sure the Peer is reachable and gives a response

```bash
openssl s_client -connect peer0-org1.localho.st:443
```

### Create Ordering Node

```bash
kubectl hlf ordnode create --image=$ORDERER_IMAGE --version=$ORDERER_VERSION \
    --storage-class=standard --enroll-id=orderer --mspid=OrdererMSP \
    --enroll-pw=ordererpw --capacity=2Gi --name=ord-node1 --ca-name=ord-ca.$NAMESPACE \
    --gateway-api-hosts=orderer0-ord.localho.st --gateway-api-name $GATEWAYNAME --gateway-api-namespace $GATEWAYNAMESPACE -n $NAMESPACE --admin-gateway-api-hosts orderer0-ord-admin.localho.st
```

Make sure the Orderer is reachable and gives a response

```bash
openssl s_client -connect orderer0-ord.localho.st:443
```







