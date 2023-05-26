---
id: getting-started
title: Getting started
---


The gateway-api implementation has been tested with traefik and istio ingress proxies. But the following can be extended to other proxies as well.


## Setup

```bash
kubectl apply -f https://github.com/kubernetes-sigs/gateway-api/releases/download/v0.7.0/experimental-install.yaml
```

## Traefik implementation

The first step is to create a service for traefik with necessary RBAC.


```yaml
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: traefik-controller

---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: traefik

spec:
  replicas: 1
  selector:
    matchLabels:
      app: traefik-lb

  template:
    metadata:
      labels:
        app: traefik-lb

    spec:
      serviceAccountName: traefik-controller
      containers:
        - name: traefik
          image: traefik:v2.10
          args:
            - --entrypoints.web.address=:80
            - --entrypoints.websecure.address=:443
            - --experimental.kubernetesgateway
            - --providers.kubernetesgateway

          ports:
            - name: web
              containerPort: 80

            - name: websecure
              containerPort: 443

---
apiVersion: v1
kind: Service
metadata:
  name: traefik

spec:
  type: LoadBalancer
  selector:
    app: traefik-lb

  ports:
    - protocol: TCP
      port: 80
      targetPort: web
      name: web

    - protocol: TCP
      port: 443
      targetPort: websecure
      name: websecure
```

RBAC configuration:

```yaml
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: gateway-role
rules:
  - apiGroups:
      - ""
    resources:
      - namespaces
    verbs:
      - list
      - watch
  - apiGroups:
      - ""
    resources:
      - services
      - endpoints
      - secrets
    verbs:
      - get
      - list
      - watch
  - apiGroups:
      - gateway.networking.k8s.io
    resources:
      - gatewayclasses
      - gateways
      - httproutes
      - tcproutes
      - tlsroutes
    verbs:
      - get
      - list
      - watch
  - apiGroups:
      - gateway.networking.k8s.io
    resources:
      - gatewayclasses/status
      - gateways/status
      - httproutes/status
      - tcproutes/status
      - tlsroutes/status
    verbs:
      - update

---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: gateway-controller
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: gateway-role
subjects:
  - kind: ServiceAccount
    name: traefik-controller
    namespace: default
```

Create a gateway class:

```yaml
apiVersion: gateway.networking.k8s.io/v1alpha2
kind: GatewayClass
metadata:
  name: my-gateway-class #Name of the gateway class
spec:
  controllerName: traefik.io/gateway-controller
```

Create a Gateway resource:

```yaml
apiVersion: gateway.networking.k8s.io/v1alpha2
kind: Gateway
metadata:
  name: traefik-gateway #Name of the gateway
spec:
  gatewayClassName: my-gateway-class #Name of the gateway class to refer to
  listeners:
    - protocol: TLS 
      port: 443
      name: tcp
      tls:
        mode: Passthrough
      allowedRoutes:
        namespaces:
            from: Selector
            selector:
                matchLabels:
                    kubernetes.io/metadata.name: hlf #Namespace where the fabric resource is deployed (CA, orderer, peer etc)
```
For more info and configuration options refer to: [Traefik's Implementation](https://doc.traefik.io/traefik/routing/providers/kubernetes-gateway/) 


## Istio implementation

Install Istio using the minimal profile:

```bash
istioctl install --set profile=minimal -y
```
By default with this installation, a GatewayClass of name istio would be created.

Now, Create a Gateway Resource:

```yaml
apiVersion: gateway.networking.k8s.io/v1alpha2
kind: Gateway
metadata:
  name: istio-gateway #Name of the gateway
spec:
  gatewayClassName: istio #Name of the gateway class to refer to
  listeners:
    - protocol: TLS 
      port: 443
      name: tcp
      tls:
        mode: Passthrough
      allowedRoutes:
        namespaces:
            from: Selector
            selector:
                matchLabels:
                    kubernetes.io/metadata.name: hlf #Namespace where the fabric resource is deployed (CA, orderer, peer etc)
```

For more info and configuration options refer to: [Istio's Implementation](https://istio.io/latest/docs/tasks/traffic-management/ingress/gateway-api/)







