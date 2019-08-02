package main

import (
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	clientmsp "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"log"
)

const (
	configPath "./config.yaml" // TODO: for multiple org & users, 2+ config?
	channelName "" // TODO
)

var (
	sdk	*fabsdk.FabricSDK
	orgID	string
	role string
)

//init() load the config file & initiate the sdk kit
func init() (err error) {
	configProvider := config.FromFile(configPath)
	sdk, err := fabsdk.New(configProvider)
	return err
}

//identify() inject the roles into the context
func identify() (err error) {
	mspClient, err := clientmsp.New(sdk.Context(), clientsmp.WithOrg(orgID))
	if err != nil {
		log.Println("create msp client fail: %s\n", err.Error())
		return err
	}

	identity, err := mspClient.GetSigningIdentity(role)
	if err != nil {
		log.Println("get admin identify fail: %s\n", err.Error())
		return err
	} else {
		fmt.Println("AdminIdentify is found:" + identity)
	}
}

//invoke() connects to the channel, makes up a transaction request,
//and handles the response
func invoke(fn string) (err error) {
	channelProvider := sdk.ChannelContext(channelName,
		fabsdk.WithUser(role),
		fabsdk.WithOrg(org))

	channelClient, err := channel.New(channelProvider)
	if err != nil {
		log.Println("create channel client fail: %s\n", err.Error())
		return err
	}

	var args [][]byte
	args = append(args, []byte("key1"))

	request := channel.Request{
		ChaincodeID: "", // TODO
		Fcn:         fn,
		Args:        args,
	}
	response, err := channelClient.Query(request)
	if err != nil {
		log.Println("operation fail: ", err.Error())
		return err
	} else {
		fmt.Printf("response is %s\n", response.Payload)
	}
}


func main() {

	// TODO: read params {org, role, fn} from the cli
	fn := "" // stub

	// init the env
	if err := init(); err != nil {
		log.Fatalf("create sdk fail: %s\n", err.Error())
	}

	// identify the org & role	
	if err = identify(); err != nil {
		log.Fatalf("identify user fail: %s\n", err.Error())
	}

	// invoke the smart contract
	if err = invoke(fn); err != nil {
		log.Fatalf("invoke chaincode fail: %s\n", err.Error())
	}	
}