package main

import (
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	mspclient "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"time"
)

func main() {
	//configFilePath := os.Args[1]
	configFilePath := "connection-org.yaml"
	channelName := "mychannel"
	mspID := "Org1MSP"
	chaincodeName := "mycc"

	enrollID := randomString(10)
	registerEnrollUser(configFilePath, enrollID, mspID)

	invokeCC(configFilePath, channelName, enrollID, mspID, chaincodeName, "CreateCar")
	//invokeCCgw(configFilePath, channelName, enrollID, mspID, chaincodeName, "CreateCar")
	queryCC(configFilePath, channelName, enrollID, mspID, chaincodeName, "QueryAllCars")
	//queryCCgw(configFilePath, channelName, enrollID, mspID, chaincodeName, "QueryAllCars")
}

func registerEnrollUser(configFilePath, enrollID, mspID string) {
	log.Info("Registering User : ", enrollID)
	sdk, err := fabsdk.New(config.FromFile(configFilePath))
	ctx := sdk.Context()
	caClient, err := mspclient.New(ctx) //mspclient.WithCAInstance("hq-guild-ca.fabric"),
	//mspclient.WithOrg(mspID),

	if err != nil {
		log.Error("Failed to create msp client: %s\n", err)
	}
	if caClient != nil {
		log.Info("ca client created")
	}
	enrollmentSecret, err := caClient.Register(&mspclient.RegistrationRequest{
		Name:           enrollID,
		Type:           "client",
		MaxEnrollments: -1,
		Affiliation:    "",
		Attributes:     nil,
		Secret:         enrollID,
	})
	if err != nil {
		log.Error(err)
	}
	err = caClient.Enroll(
		enrollID,
		mspclient.WithSecret(enrollmentSecret),
		mspclient.WithProfile("tls"),
	)
	if err != nil {
		log.Error(errors.WithMessage(err, "failed to register identity"))
	}

	wallet, err := gateway.NewFileSystemWallet(fmt.Sprintf("wallet/%s", mspID))

	signingIdentity, err := caClient.GetSigningIdentity(enrollID)
	key, err := signingIdentity.PrivateKey().Bytes()
	identity := gateway.NewX509Identity(mspID, string(signingIdentity.EnrollmentCertificate()), string(key))

	err = wallet.Put(enrollID, identity)
	if err != nil {
		log.Error(err)
	}

}
func invokeCCgw(configFilePath, channelName, userName, mspID, chaincodeName, fcn string) {

	configBackend := config.FromFile(configFilePath)
	sdk, err := fabsdk.New(configBackend)
	if err != nil {
		log.Error(err)
	}

	gw, err := gateway.Connect(
		gateway.WithSDK(sdk),
		gateway.WithUser(userName),
	)
	if err != nil {
		log.Error("Failed to create new Gateway: %s", err)
	}
	defer gw.Close()
	nw, err := gw.GetNetwork(channelName)
	if err != nil {
		log.Error("Failed to get network: %s", err)
	}

	contract := nw.GetContract(chaincodeName)

	resp, err := contract.SubmitTransaction(fcn, userName, "a", "b", "1", "ewdscwds")

	if err != nil {
		log.Error("Failed submit transaction: %s", err)
	}
	log.Info(resp)

}
func invokeCC(configFilePath, channelName, userName, mspID, chaincodeName, fcn string) {

	configBackend := config.FromFile(configFilePath)
	sdk, err := fabsdk.New(configBackend)
	if err != nil {
		log.Error(err)
	}

	chContext := sdk.ChannelContext(
		channelName,
		fabsdk.WithUser(userName),
		fabsdk.WithOrg(mspID),
	)

	ch, err := channel.New(chContext)
	if err != nil {
		log.Error(err)
	}

	var args [][]byte

	inputArgs := []string{userName, "23", "234", "2324", "234"}
	for _, arg := range inputArgs {
		args = append(args, []byte(arg))
	}
	response, err := ch.Execute(
		channel.Request{
			ChaincodeID:     chaincodeName,
			Fcn:             fcn,
			Args:            args,
			TransientMap:    nil,
			InvocationChain: nil,
			IsInit:          false,
		},
	)

	if err != nil {
		log.Error(err)
	}

	log.Infof("txid=%s", response.TransactionID)
}

func queryCC(configFilePath, channelName, userName, mspID, chaincodeName, fcn string) {

	configBackend := config.FromFile(configFilePath)
	sdk, err := fabsdk.New(configBackend)
	if err != nil {
		log.Error(err)
	}
	log.Println(sdk)
	chContext := sdk.ChannelContext(
		channelName,
		fabsdk.WithUser(userName),
		fabsdk.WithOrg(mspID),
	)

	ch, err := channel.New(chContext)
	if err != nil {
		log.Error(err)
	}

	response, err := ch.Query(
		channel.Request{
			ChaincodeID:     chaincodeName,
			Fcn:             fcn,
			Args:            nil,
			TransientMap:    nil,
			InvocationChain: nil,
			IsInit:          false,
		},
	)

	if err != nil {
		log.Error(err)
	}
	log.Infof("response=%s", response.Payload)
}

func queryCCgw(configFilePath, channelName, userName, mspID, chaincodeName, fcn string) {

	configBackend := config.FromFile(configFilePath)
	sdk, err := fabsdk.New(configBackend)
	if err != nil {
		log.Error(err)
	}
	gw, err := gateway.Connect(
		gateway.WithSDK(sdk),
		gateway.WithUser(userName),
	)

	if err != nil {
		log.Error("Failed to create new Gateway: %s", err)
	}
	defer gw.Close()
	nw, err := gw.GetNetwork(channelName)
	if err != nil {
		log.Error("Failed to get network: %s", err)
	}

	contract := nw.GetContract(chaincodeName)

	resp, err := contract.EvaluateTransaction(fcn)

	if err != nil {
		log.Error("Failed submit transaction: %s", err)
	}
	log.Info(string(resp))

}

func randomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, length)
	rand.Read(b)
	return fmt.Sprintf("%x", b)[:length]
}
