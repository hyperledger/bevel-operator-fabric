# Fabric Peer

Properties

| Parameter | Description  |  Default | Required |
|---|---|---|---|
| `dockerSocketPath`  |  Socket to connect to the docker daemon | /var/run/docker.sock | No |
| `hosts`  |  CNAME for the generated certificate | [] | Yes |
| `operationHosts`  |  CNAME to generate the certificate for the operations service | [] | Yes |
| `operationIPs`  |  IPs to generate the certificate for the operations service | [] | Yes |
| `image` | Docker image name  |  null | Yes |
| `externalChaincodeBuilder` | Add file server to support external chaincode builder  |  null | Yes |
| `tag`  | Docker image version  | null  | Yes |
| `stateDb`  | State database for the Peer  | null  | Yes |
| `mspID`  | ID of the MSP the peer belongs to   | null  | Yes |
| `service.type`  | Kubernetes service type  | null  | Yes |
| `service.nodePortOperations`  | Port number for the operations service  | null  | Yes |
| `service.nodePortEvent`  | Port number for the events service  | null  | Yes |
| `service.nodePortRequest`  | Port number for the request peer service, required if ingress is not enabled for gossip  | null  | Yes |
| `secret.enrollment.component.cahost`  | CA server host  | null  | Yes |
| `secret.enrollment.component.caname`  | CA name  | null  | Yes |
| `secret.enrollment.component.caport`  | CA server port  | null  | Yes |
| `secret.enrollment.component.catls.cacert`  | CA server certificate  | null  | Yes |
| `secret.enrollment.component.csr.hosts`  | SANs for the generated certificate  | null  | No |
| `secret.enrollment.component.csr.cn`  | CN for the generated certificate  | null  | No |
| `secret.enrollment.component.enrollid`  | CA enroll username  | null  | Yes |
| `secret.enrollment.component.enrollsecret`  | CA enroll password  | null  | Yes |
| `secret.enrollment.tls.cahost`  | CA server port  | null  | Yes |
| `secret.enrollment.tls.caname`  | CA server port  | null  | Yes |
| `secret.enrollment.tls.caport`  | CA server port  | null  | Yes |
| `secret.enrollment.tls.catls.cacert`  | CA server port  | null  | Yes |
| `secret.enrollment.tls.csr.hosts`  | SANs for the generated certificate  | null  | No |
| `secret.enrollment.tls.csr.cn`  | CN for the generated certificate  | null  | No |
| `secret.enrollment.tls.enrollid`  | CA enroll username  | null  | Yes |
| `secret.enrollment.tls.enrollsecret`  | CA enroll password  | null  | Yes |

