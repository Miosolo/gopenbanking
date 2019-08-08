package main

import (
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	clientmsp "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"log"
	"flag"
	"os"
)

const (
	configPath = "./config.yaml" // TODO: for multiple org & users, 2+ config?
)

//identify() inject the roles into the context
func identify(sdk *fabsdk.FabricSDK, orgName, orgRole string) (err error) {
	mspClient, err := clientmsp.New(sdk.Context(), clientmsp.WithOrg(orgName))
	if err != nil {
		log.Printf("create msp client fail: %s\n", err.Error())
		return err
	}

	identity, err := mspClient.GetSigningIdentity(orgRole)
	if err != nil {
		log.Printf("get identify fail: %s\n", err.Error())
		return err
	}
	fmt.Println("Identify is found: " + identity.Identifier().MSPID)
	return nil
}

//invoke() connects to the channel, makes up a transaction request,
//and handles the response
func invoke(sdk *fabsdk.FabricSDK, channelID, orgID, orgRole, ccFunction string) (err error) {
	channelProvider := sdk.ChannelContext(channelID,
		fabsdk.WithUser(orgRole),
		fabsdk.WithOrg(orgID))

	channelClient, err := channel.New(channelProvider)
	if err != nil {
		log.Printf("create channel client fail: %s\n", err.Error())
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
	}
	
	fmt.Printf("response is %s\n", response.Payload)
	return nil
}

func main() {
	// define the flags & parse the params
	channelID := flag.String("chan", "orgchannel", `Name of the channel, default "orgchannel"`)
	orgName := flag.String("org", "", "Name of your orgnization")
	orgRole := flag.String("role", "client", `Your role in this organization, default "client"`)
	chaincodeID := flag.String("cc", "", "ID of the chaincode instanciated")
	ccFunction := flag.String("fn", "", "The function of smart contract to call")
	flag.Parse()

	// TODO: set env
	os.Setenv("FABRIC_SDK_GO_PROJECT_PATH", "")

	// init the env
	configProvider := config.FromFile(configPath)
	sdk, err := fabsdk.New(configProvider)
	defer sdk.Close()

	if err != nil {
		log.Fatalf("Failed to create new SDK: %s", err)
	}

	// identify the org & role	
	if err := identify(); err != nil {
		log.Fatalf("identify user fail: %s\n", err.Error())
	}

	// invoke the smart contract
	if err := invoke(fn); err != nil {
		log.Fatalf("invoke chaincode fail: %s\n", err.Error())
	}	
}