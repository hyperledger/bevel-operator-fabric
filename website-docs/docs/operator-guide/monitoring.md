---
id: monitoring
title: Monitoring
---

The CRDs for the orderer, peer, and certificate authority have an optional parameter to create the service monitors to scrape the metrics automatically if Prometheus Operator is installed on the cluster.


```yaml
  serviceMonitor:
    enabled: true
    interval: 10s
    labels: {}
    sampleLimit: 0
    scrapeTimeout: 10s
```

There are some dashboards available in the Github repository for Grafana available at https://github.com/hyperledger/bevel-operator-fabric/dashboards.
