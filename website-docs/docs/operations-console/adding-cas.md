---
id: adding-cas
title: Adding Certificate Authorities
---

The steps to follow to add a Certificate Authority to the Fabric Operations console are:
- Export Certificate Authority to JSON format
- Enter the Fabric Operations Console UI
- Go to `Nodes`
- Click on `Import Certificate Authority`
- Select the JSON from the file system
- Click on `Add Certificate Authority`

## Export Certificate Authority to JSON

```bash
export CA_NAME=xxxxx
export CA_NS=default
kubectl hlf fop export ca --name="${CA_NAME}" --namespace="${CA_NS}" --out="${CA_NAME}.json"
```

## Enter the Fabric Operations Console UI

Open a browser and navigate to the URL you configured when creating the Fabric Operations Console.


## Go to `Nodes`

Click on `Nodes` at the sidenav to see the Certificate Authorities, Certificate Authorities and Ordering Services

![img_1.png](/img/img_1.png)

## Click on `Import Certificate Authority`

Click on `Import Certificate Authority` to open the dialog to import the Certificate Authority.


![img.png](/img/import_ca_table.png)

## Select the JSON from the file system

Click on `Add file` and select the JSON file you exported from the step `Export Certificate Authority to JSON`.

![img.png](/img/import_ord_service.png)

## Click on `Add Certificate Authority`

The last step is to click on `Add Certificate Authority` and the Certificate Authority will be imported to the console.


![img.png](/img/final_add_ca.png)

