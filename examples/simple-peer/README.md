```bash
kubectl get fabriccas.hlf.kungfusoftware.es org1-ca -o jsonpath='{.status.tls_cert}' | base64 --wrap=0

kubectl apply -f ./org1-ca.yaml

kubectl apply -f ./peer0.yaml

kubectl get fabriccas.hlf.kungfusoftware.es org1-ca -o jsonpath='{.status.tlsca_cert}' > crypto-config/Org1MSP/msp/tlscacerts/tlsca.pem
kubectl get fabriccas.hlf.kungfusoftware.es org1-ca -o jsonpath='{.status.ca_cert}' > crypto-config/Org1MSP/msp/cacerts/cacert.pem

kubectl get fabricpeers.hlf.kungfusoftware.es org1-peer0 -o jsonpath='{.status.url}'
export FABRIC_CFG_PATH=$PWD
mkdir -p crypto-config/Org1MSP/msp/tlscacerts  
mkdir -p crypto-config/Org1MSP/msp/cacerts  
mkdir -p crypto-config/Org1MSP/msp/keystore
mkdir -p crypto-config/Org1MSP/msp/signcerts

configtxgen -printOrg Org1MSP

kubectl hlf ca register --name=org1-ca --user=admin --secret=adminpw --type=admin \
 --enroll-id enroll --enroll-secret=enrollpw --mspid Org1MSP  

kubectl hlf ca enroll --name=org1-ca --user=admin --secret=adminpw --mspid Org1MSP \
        --ca-name ca  --output peer-org1.yaml

kubectl hlf inspect --output org1.yaml -o Org1MSP -o OrdererMSP

kubectl hlf channel join --name=channel1 --config=org1.yaml \
    --user=admin -p=org1-peer0.default

kubectl hlf channel addanchorpeer --channel=channel1 --config=org1.yaml \
    --user=admin --peer=org1-peer0.default 


kubectl get fabricpeers.hlf.kungfusoftware.es org1-peer0 -o jsonpath='{.status.tls_cert}' 
kubectl get fabricpeers.hlf.kungfusoftware.es org1-peer1 -o jsonpath='{.status.tls_cert}' 


kubectl hlf channel join --name=channel1 --config=org1.yaml \
    --user=admin -p=org1-peer1.default

kubectl hlf chaincode install --path=../../fixtures/chaincodes/fabcar/go \
    --config=org1.yaml --language=golang --label=fabcar --user=admin --peer=org1-peer0.default

kubectl hlf chaincode install --path=../../fixtures/chaincodes/fabcar/go \
    --config=org1.yaml --language=golang --label=fabcar --user=admin --peer=org1-peer1.default

kubectl hlf chaincode queryinstalled --config=org1.yaml --user=admin --peer=org1-peer0.default

kubectl hlf chaincode queryinstalled --config=org1.yaml --user=admin --peer=org1-peer1.default

kubectl hlf chaincode approveformyorg --config=org1.yaml --user=admin --peer=org1-peer0.default \
    --package-id=fabcar:0c616be7eebace4b3c2aa0890944875f695653dbf80bef7d95f3eed6667b5f40 \
    --version "1.0" --sequence 1 --name=fabcar --collections-config=./pdc.json \
    --policy="OR('Org1MSP.member')" --channel=channel1


kubectl hlf chaincode commit --config=org1.yaml --user=admin --peer=org1-peer0.default \
    --version "1.0" --sequence 1 --name=fabcar --collections-config=./pdc.json  \
    --policy="OR('Org1MSP.member')" --channel=channel1

kubectl hlf chaincode invoke --config=org1.yaml \
    --user=admin --peer=org1-peer0.default \
    --chaincode=fabcar --channel=channel1 \
    --fcn=initLedger -a '[]'

kubectl hlf chaincode query --config=org1.yaml \
    --user=admin --peer=org1-peer0.default \
    --chaincode=fabcar --channel=channel1 \
    --fcn=QueryAllCars -a '[]'

```