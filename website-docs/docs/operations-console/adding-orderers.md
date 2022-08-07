---
id: adding-orderers
title: Adding Orderer nodes
---



The steps to follow to add a Ordering Service to the Fabric Operations console are:
- Export Ordering Service to JSON format
- Enter the Fabric Operations Console UI
- Go to `Nodes`
- Click on `Import Ordering services`
- Select the JSON from the file system
- Click on `Add Ordering services`

## Export Ordering Service to JSON

```bash
export ORDERER_NAME=orderer0-ordmsp068wi-5vph
export ORDERER_NS=default
kubectl hlf fop export orderer --cluster-id=orderermsp1 --cluster-name="Cluster 1" --name=$ORDERER_NAME --namespace=$ORDERER_NS --out="${ORDERER_NAME}_${ORDERER_NS}.json"
```

## Enter the Fabric Operations Console UI

Open a browser and navigate to the URL you configured when creating the Fabric Operations Console.


## Go to `Nodes`

Click on `Nodes` at the sidenav to see the Peers, Certificate Authorities and Ordering Services

![img_1.png](/img/img_1.png)

## Click on `Import Ordering services`

Click on `Import Ordering services` to open the dialog to import the Ordering Service.

![img.png](/img/ordering_service_import.png)

## Select the JSON from the file system

Click on `Add file` and select the JSON file you exported from the step `Export ordering services to JSON`.

![img.png](/img/select_json_ordering_service.png)

## Click on `Add Ordering services`

The last step is to set `Ordering service location` to Kubernetes and to click on `Add Ordering services` and the Ordering Service will be imported to the console.

![img.png](/img/add_ordering_service.png)



