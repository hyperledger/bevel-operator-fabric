## Contributing


### Pre Requisites

- Goland 1.18
- Makefile


### Install controller-gen

```bash
go install  sigs.k8s.io/controller-tools/cmd/controller-gen@0.16.4
```


### Generate CRDs and Controller

```bash
make generate manifest install
```


### Deploy

```bash
# set the image name so that it always get redeployed
export IMAGE=kfsoftware/hlf-operator:1.11.0-support-chaincodes-$(date +%s%3N)
# build the binary 
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o hlf-operator ./main.go
 
# build the docker image
docker build -t $IMAGE --platform=linux/amd64 .

# import the images in all the nodes to avoid having to push the image to a registry 
k3d image import $IMAGE -c k8s-hlf

# deploy the new version of the operator 
make deploy IMG=$IMAGE

```
