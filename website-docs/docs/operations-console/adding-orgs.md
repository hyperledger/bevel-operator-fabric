---
id: adding-orgs
title: Adding Organizations
---

The steps to follow to add a Organization to the Fabric Operations console are:
- Export Organization` to JSON format
- Enter the Fabric Operations Console UI
- Go to `Nodes`
- Click on `Import MSP Definition`
- Select the JSON from the file system
- Click on `Import MSP Definition`

## Export Organization to JSON

```bash
export CA_NAME=xxxxx
export CA_NS=default
export MSP_ID=euipomsp
export CONSOLE_URL=https://console.hlf.example.org
kubectl hlf fop export org --msp-id=$MSP_ID --name=$CA_NAME --namespace=$CA_NS --out="${MSP_ID}_${CA_NAME}_${CA_NS}.json" --host-url="${CONSOLE_URL}"

```

## Enter the Fabric Operations Console UI

Open a browser and navigate to the URL you configured when creating the Fabric Operations Console.


## Go to `Organizations`

Click on `Organizations` at the sidenav to see the Available organizations

![img_4.png](/img/img_4.png)

## Click on `Import MSP Definition`

Click on `Import MSP Definition` to open the dialog to import the Organization.

![img_6.png](/img/img_6.png)

## Select the JSON from the file system

Click on `Add file` and select the JSON file you exported from the step `Export Organization to JSON`.

![img_7.png](/img/img_7.png)

## Click on `Add Organization`

The last step is to click on `Add Organization` and the Organization will be imported to the console.

![img_8.png](/img/img_8.png)
