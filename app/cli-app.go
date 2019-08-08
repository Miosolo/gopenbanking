package main

import (
	"flag"
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	clientmsp "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"log"
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
func invoke(sdk *fabsdk.FabricSDK, channelID, orgID,
	orgRole, chaincodeID, ccFunction string, args []string) (resp string, err error) {
	channelProvider := sdk.ChannelContext(channelID,
		fabsdk.WithUser(orgRole),
		fabsdk.WithOrg(orgID))

	channelClient, err := channel.New(channelProvider)
	if err != nil {
		log.Printf("create channel client fail: %s\n", err.Error())
		return err
	}

	var byteArgs [][]byte
	for _, arg := range(args) {
		byteArgs = append(byteArgs, []byte(arg))
	}

	request := channel.Request{
		ChaincodeID: chaincodeID,
		Fcn:         ccFunction,
		Args:        byteArgs,
	}

	response, err := channelClient.Query(request)
	if err != nil {
		log.Println("operation fail: ", err.Error())
		return "", err
	}

	return string(response.Payload), nil
}

func main() {
	// define the flags & parse the params
	channelID := flag.String("chan", "orgchannel", `Name of the channel, default "orgchannel"`)
	orgID := flag.String("org", "", "Name of your orgnization")
	orgRole := flag.String("role", "client", `Your role in this organization, default "client"`)
	chaincodeID := flag.String("cc", "", "ID of the chaincode instanciated")
	flag.Parse()

	// set env for YAML parsing
	os.Setenv("FABRIC_SDK_GO_PROJECT_PATH", "$PWD")
	os.Setenv("FABRIC_ORG_ID", *orgID)
	os.Setenv("CRYPTOCONFIG_FIXTURES_PATH", "crypto-config/"+*orgID)

	// init the env
	configProvider := config.FromFile(configPath)
	sdk, err := fabsdk.New(configProvider)
	defer sdk.Close()

	if err != nil {
		log.Fatalf("Failed to create new SDK: %s", err)
	}

	// identify the org & role
	if err := identify(sdk, *orgID, *orgRole); err != nil {
		log.Fatalf("identify user fail: %s\n", err.Error())
	}

	// print the instructions // TODO
	fmt.Println(`Functions and parameters of the ANZ-CITI Banking Network:
	- 
	- 
	- exit: terminate the loop and exit`)
	// start loop
	for true {
		// read the stdin input
		fmt.Printf("Enter the function & params: ")
		var fn string
		var args [3]string // max args: 3 due to the chaincode
		// inputCnt - 1 = args Count
		inputCnt, _ := fmt.Scanln(&fn, &args[0], &args[1], &args[2])

		if (inputCnt == 0) {
			continue
		} else if (fn == "exit") {
			fmt.Println("bye")
			return
		}

		// else, invoke the smart contract
		if response, err := invoke(sdk, *channelID, *orgID, *orgRole, *chaincodeID, fn, args[1:inputCnt]); err != nil {
			log.Println("invoke chaincode fail: %s\n", err.Error())
		} else {
			fmt.Println("Response: " + response)
		}
	}
}
