# Fabric CA

Properties

| Parameter | Description  |  Default | Required | 
|---|---|---|---|
| `hosts`  |  CNAME for the generated certificate | [] | Yes |
| `image` | Docker image name  |  null | Yes |
| `version`  | Docker image version  | null  | Yes |
| `admin_user`  | Username for the admin user  | null  | Yes |
| `admin_password`  | Password for the admin user  | null  | Yes |
| `service.type`  | Kubernetes service type  | null  | Yes |
| `caName` | Default certificate authority name | ca | Yes | 