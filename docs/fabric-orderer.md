# Fabric orderer

| Parameter | Description  |  Default | Required |
|---|---|---|---|
| `genesis`  |  Base64 genesis block | null | Yes |
| `image` | Docker image name  |  null | Yes |
| `tag`  | Docker image version  | null  | Yes |
| `mspID`  | ID of the MSP of the orderer   | null  | Yes |
| `hosts`  | ID of the MSP of the orderer   | []  | Yes |
| `tls_root_cert`  | Root TLS certificate | null  | Yes |
| `tls_cert`  | TLS certificate | null  | Yes |
| `tls_key`  | Key certificate | null  | Yes |
| `service.type`  | Kubernetes service type  | null  | Yes |
| `service.nodePortOperations`  | Port number for the operations service  | null  | Yes |
| `service.nodePortRequest`  | Port number for the request orderer service  | null  | Yes |
| `secret.enrollment.component.cahost`  | CA server host  | null  | Yes |
| `secret.enrollment.component.caname`  | CA name  | null  | Yes |
| `secret.enrollment.component.caport`  | CA server port  | null  | Yes |
| `secret.enrollment.component.catls.cacert`  | CA server certificate  | null  | Yes |
| `secret.enrollment.component.csr.hosts`  | SANs for the generated certificate  | null  | No |
| `secret.enrollment.component.csr.cn`  | CN for the generated certificate  | null  | No |
| `secret.enrollment.component.enrollid`  | CA enroll username  | null  | Yes |
| `secret.enrollment.component.enrollsecret`  | CA enroll password  | null  | Yes |


