```bash
mkdir -p crypto-config/OrdererMSP/msp/tlscacerts  
mkdir -p crypto-config/OrdererMSP/msp/cacerts  
mkdir -p crypto-config/OrdererMSP/msp/keystore
mkdir -p crypto-config/OrdererMSP/msp/signcerts



kubectl apply -f orderer0.yaml -f ./orderer1.yaml -f ./orderer2.yaml -f ./orderer3.yaml

kubectl get fabricorderernodes.hlf.kungfusoftware.es  orderer0  -o jsonpath='{.status.url}'  
kubectl get fabricorderernodes.hlf.kungfusoftware.es  orderer1  -o jsonpath='{.status.url}'  
kubectl get fabricorderernodes.hlf.kungfusoftware.es  orderer2  -o jsonpath='{.status.url}'  
kubectl get fabricorderernodes.hlf.kungfusoftware.es  orderer3  -o jsonpath='{.status.url}'

kubectl get fabricorderernodes.hlf.kungfusoftware.es  orderer0  -o jsonpath='{.status.adminPort}'  
kubectl get fabricorderernodes.hlf.kungfusoftware.es  orderer1  -o jsonpath='{.status.adminPort}'  
kubectl get fabricorderernodes.hlf.kungfusoftware.es  orderer2  -o jsonpath='{.status.adminPort}'  
kubectl get fabricorderernodes.hlf.kungfusoftware.es  orderer3  -o jsonpath='{.status.adminPort}'

kubectl get fabricorderernodes.hlf.kungfusoftware.es  orderer0  -o jsonpath='{.status.tlsCert}' > crypto-config/OrdererMSP/orderer0-server.pem 
kubectl get fabricorderernodes.hlf.kungfusoftware.es  orderer1  -o jsonpath='{.status.tlsCert}' > crypto-config/OrdererMSP/orderer1-server.pem 
kubectl get fabricorderernodes.hlf.kungfusoftware.es  orderer2  -o jsonpath='{.status.tlsCert}' > crypto-config/OrdererMSP/orderer2-server.pem 
kubectl get fabricorderernodes.hlf.kungfusoftware.es  orderer3  -o jsonpath='{.status.tlsCert}' > crypto-config/OrdererMSP/orderer2-server.pem 


kubectl get service  orderer0  -o jsonpath='{.spec.ports[1].nodePort}'  
kubectl get service  orderer1  -o jsonpath='{.spec.ports[1].nodePort}'  
kubectl get service  orderer2  -o jsonpath='{.spec.ports[1].nodePort}'  
kubectl get service  orderer3  -o jsonpath='{.spec.ports[1].nodePort}'  

export FABRIC_CFG_PATH=$PWD
configtxgen -profile SampleAppChannelEtcdRaft -outputBlock genesis_block.pb -channelID channel1

kubectl hlf ca register --name=orderer-ca --namespace=default --user=admin --secret=adminpw \    
    --type=admin --enroll-id enroll --enroll-secret=enrollpw --mspid=OrdererMSP 

kubectl hlf ca enroll --name=orderer-ca --namespace=default --user=admin --secret=adminpw --mspid OrdererMSP \
        --ca-name tlsca  --output admintls-ordservice.yaml
        

export OSN_TLS_CA_ROOT_CERT=./crypto-config/OrdererMSP/msp/tlscacerts/tlscacert.pem
export ADMIN_TLS_SIGN_CERT=./admin-crt.pem
export ADMIN_TLS_PRIVATE_KEY=./admin-pk.pem

osnadmin channel join --channelID channel1  --config-block genesis_block.pb -o 172.24.0.2:31732 --ca-file $OSN_TLS_CA_ROOT_CERT --client-cert $ADMIN_TLS_SIGN_CERT --client-key $ADMIN_TLS_PRIVATE_KEY

osnadmin channel join --channelID channel1  --config-block genesis_block.pb -o 172.24.0.2:30569 --ca-file $OSN_TLS_CA_ROOT_CERT --client-cert $ADMIN_TLS_SIGN_CERT --client-key $ADMIN_TLS_PRIVATE_KEY
osnadmin channel join --channelID channel1  --config-block genesis_block.pb -o 172.24.0.2:31805 --ca-file $OSN_TLS_CA_ROOT_CERT --client-cert $ADMIN_TLS_SIGN_CERT --client-key $ADMIN_TLS_PRIVATE_KEY
osnadmin channel join --channelID channel1  --config-block genesis_block.pb -o 172.24.0.2:30816 --ca-file $OSN_TLS_CA_ROOT_CERT --client-cert $ADMIN_TLS_SIGN_CERT --client-key $ADMIN_TLS_PRIVATE_KEY
osnadmin channel join --channelID channel1  --config-block genesis_block.pb -o 172.24.0.2:31185 --ca-file $OSN_TLS_CA_ROOT_CERT --client-cert $ADMIN_TLS_SIGN_CERT --client-key $ADMIN_TLS_PRIVATE_KEY

osnadmin channel list --channelID channel1 -o 172.24.0.2:30569 --ca-file $OSN_TLS_CA_ROOT_CERT --client-cert $ADMIN_TLS_SIGN_CERT --client-key $ADMIN_TLS_PRIVATE_KEY
osnadmin channel list --channelID channel1 -o 172.24.0.2:31805 --ca-file $OSN_TLS_CA_ROOT_CERT --client-cert $ADMIN_TLS_SIGN_CERT --client-key $ADMIN_TLS_PRIVATE_KEY
osnadmin channel list --channelID channel1 -o 172.24.0.2:30816 --ca-file $OSN_TLS_CA_ROOT_CERT --client-cert $ADMIN_TLS_SIGN_CERT --client-key $ADMIN_TLS_PRIVATE_KEY
osnadmin channel list --channelID channel1 -o 172.24.0.2:31185 --ca-file $OSN_TLS_CA_ROOT_CERT --client-cert $ADMIN_TLS_SIGN_CERT --client-key $ADMIN_TLS_PRIVATE_KEY

# add an additional orderer node
kubectl apply -f ./orderer4.yaml

kubectl get fabricorderernodes.hlf.kungfusoftware.es  orderer4  -o jsonpath='{.status.url}'  
kubectl get fabricorderernodes.hlf.kungfusoftware.es  orderer4  -o jsonpath='{.status.tlsCert}' > crypto-config/OrdererMSP/orderer4-server.pem 
kubectl get fabricorderernodes.hlf.kungfusoftware.es  orderer4  -o jsonpath='{.status.tlsCert}' > crypto-config/OrdererMSP/orderer4-server.pem | base64 --wrap=0
kubectl get service  orderer4  -o jsonpath='{.spec.ports[1].nodePort}'  
osnadmin channel join --channelID channel1  --config-block genesis_block.pb -o 172.24.0.2:31765 --ca-file $OSN_TLS_CA_ROOT_CERT --client-cert $ADMIN_TLS_SIGN_CERT --client-key $ADMIN_TLS_PRIVATE_KEY
osnadmin channel list --channelID channel1 -o 172.24.0.2:31765 --ca-file $OSN_TLS_CA_ROOT_CERT --client-cert $ADMIN_TLS_SIGN_CERT --client-key $ADMIN_TLS_PRIVATE_KEY


# update channel to include consenter
export CH_NAME=channel1
peer channel fetch config config_block.pb -o 172.24.0.2:30503 -c channel1 --tls --cafile $OSN_TLS_CA_ROOT_CERT
configtxlator proto_decode --input config_block.pb --type common.Block --output config_block.json
jq '.data.data[0].payload.data.config' config_block.json > config.json
cp config.json modified_config.json
configtxlator proto_encode --input config.json --type common.Config --output config.pb
# update the modified_config.json
configtxlator proto_encode --input modified_config.json --type common.Config --output modified_config.pb
configtxlator compute_update --channel_id $CH_NAME --original config.pb --updated modified_config.pb --output config_update.pb

configtxlator proto_decode --input config_update.pb --type common.ConfigUpdate --output config_update.json
echo '{"payload":{"header":{"channel_header":{"channel_id":"'$CH_NAME'", "type":2}},"data":{"config_update":'$(cat config_update.json)'}}}' | jq . > config_update_in_envelope.json
configtxlator proto_encode --input config_update_in_envelope.json --type common.Envelope --output config_update_in_envelope.pb
export CORE_PEER_LOCALMSPID=OrdererMSP
export CORE_PEER_MSPCONFIGPATH=$PWD/crypto-config/OrdererMSP/msp
peer channel update -f config_update_in_envelope.pb -c $CH_NAME -o 172.24.0.2:30503 --tls --cafile $OSN_TLS_CA_ROOT_CERT


# remove an orderer node
rm config_block.pb config_block.json config.json config.pb modified_config.json modified_config.pb config_update.json config_update.pb config_update_in_envelope.json config_update_in_envelope.pb

peer channel fetch config config_block.pb -o 172.24.0.2:30503 -c channel1 --tls --cafile $OSN_TLS_CA_ROOT_CERT
configtxlator proto_decode --input config_block.pb --type common.Block --output config_block.json
jq '.data.data[0].payload.data.config' config_block.json > config.json
cp config.json modified_config.json
configtxlator proto_encode --input config.json --type common.Config --output config.pb
# update the modified_config.json
configtxlator proto_encode --input modified_config.json --type common.Config --output modified_config.pb
configtxlator compute_update --channel_id $CH_NAME --original config.pb --updated modified_config.pb --output config_update.pb

configtxlator proto_decode --input config_update.pb --type common.ConfigUpdate --output config_update.json
echo '{"payload":{"header":{"channel_header":{"channel_id":"'$CH_NAME'", "type":2}},"data":{"config_update":'$(cat config_update.json)'}}}' | jq . > config_update_in_envelope.json
configtxlator proto_encode --input config_update_in_envelope.json --type common.Envelope --output config_update_in_envelope.pb
export CORE_PEER_LOCALMSPID=OrdererMSP
export CORE_PEER_MSPCONFIGPATH=$PWD/crypto-config/OrdererMSP/msp
peer channel update -f config_update_in_envelope.pb -c $CH_NAME -o 172.24.0.2:30503 --tls --cafile $OSN_TLS_CA_ROOT_CERT

```

## Add organization
```bash

# update channel to include organization
export CH_NAME=channel1
peer channel fetch config config_block.pb -o 172.24.0.2:30503 -c channel1 --tls --cafile $OSN_TLS_CA_ROOT_CERT
configtxlator proto_decode --input config_block.pb --type common.Block --output config_block.json
jq '.data.data[0].payload.data.config' config_block.json > config.json
cp config.json modified_config.json
configtxlator proto_encode --input config.json --type common.Config --output config.pb
# update the modified_config.json
jq -s '.[0] * {"channel_group":{"groups":{"Application":{"groups": {"Org1MSP":.[1]}}}}}' config.json ../simple-peer/org1.json > modified_config.json


configtxlator proto_encode --input modified_config.json --type common.Config --output modified_config.pb
configtxlator compute_update --channel_id $CH_NAME --original config.pb --updated modified_config.pb --output config_update.pb

configtxlator proto_decode --input config_update.pb --type common.ConfigUpdate --output config_update.json
echo '{"payload":{"header":{"channel_header":{"channel_id":"'$CH_NAME'", "type":2}},"data":{"config_update":'$(cat config_update.json)'}}}' | jq . > config_update_in_envelope.json
configtxlator proto_encode --input config_update_in_envelope.json --type common.Envelope --output config_update_in_envelope.pb
export CORE_PEER_LOCALMSPID=OrdererMSP
export CORE_PEER_MSPCONFIGPATH=$PWD/crypto-config/OrdererMSP/msp
peer channel update -f config_update_in_envelope.pb -c $CH_NAME -o 172.24.0.2:30503 --tls --cafile $OSN_TLS_CA_ROOT_CERT

```

## Delete orderers
```bash
kubectl delete -f orderer0.yaml -f ./orderer1.yaml -f ./orderer2.yaml -f ./orderer3.yaml -f orderer4.yaml

```