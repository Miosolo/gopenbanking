package main

import (
  "flag"
  "fmt"
  "log"
  "os"

  "github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
  clientmsp "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
  "github.com/hyperledger/fabric-sdk-go/pkg/core/config"
  "github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

const (
  configPath = "./config.yaml"
)

//identify() checks the user identity
func identify(sdk *fabsdk.FabricSDK, orgID, orgUser string) (err error) {
  mspClient, err := clientmsp.New(sdk.Context(), clientmsp.WithOrg(orgID))
  if err != nil {
    log.Printf("create msp client fail: %s\n", err.Error())
    return err
  }

  identity, err := mspClient.GetSigningIdentity(orgUser)
  if err != nil {
    return err
  }

  log.Println("Identity is found: " + identity.Identifier().MSPID)
  return nil
}

//invoke() connects to the channel, makes up a transaction request,
//and handles the response
func invoke(sdk *fabsdk.FabricSDK, channelID, orgID,
  orgUser, chaincodeID, ccFunction string, args []string) (resp string, err error) {
  channelProvider := sdk.ChannelContext(channelID,
    fabsdk.WithUser(orgUser),
    fabsdk.WithOrg(orgID))

  channelClient, err := channel.New(channelProvider)
  if err != nil {
    log.Printf("create channel client fail: %s\n", err.Error())
    return "", err
  }

  var byteArgs [][]byte
  for _, arg := range args {
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
  channelID := flag.String("chan", "orgchannel", `Name of the channel`)
  orgID := flag.String("org", "ANZBank", "Name of your orgnization")
  orgUser := flag.String("user", "Admin", `Your User ID in this organization`)
  chaincodeID := flag.String("cc", "", "ID of the chaincode instanciated")
  flag.Parse()

  // set env for YAML parsing
  os.Setenv("FABRIC_CRYPTOCONFIG_PATH", os.Getenv("PWD")+"/../crypto-config/")
  os.Setenv("FABRIC_ORG_ID", *orgID)
  os.Setenv("FABRIC_ORG_USER", *orgUser)

  // init the env
  configProvider := config.FromFile(configPath)
  sdk, err := fabsdk.New(configProvider)
  defer sdk.Close()

  if err != nil {
    log.Fatalf("Failed to create new SDK: %s", err)
  }

  // identify the org & role
  if err := identify(sdk, *orgID, *orgUser); err != nil {
    log.Fatalf("identify %s fail: %s\n", *orgUser, err.Error())
  }

  // print the instructions
  fmt.Println(`==========INSTRUCTIONS==========
Functions and parameters of the ANZ-CITI Banking Network:
  - (Bank) "get" + account
  - (Bank) "add" + account + value
  - (Bank) "reduce" + account + value
  - (Bank) "create" + account + inititial value
  - (Bank) "delete" + account
  - (Bank) "tranfer" + debit account + credit account + tranfer value
  - (Any) "query" + object type + account name
  - (Any) "exit": terminate the loop and exit
=================================`)
  // start loop
  for true {
    // read the stdin input
    fmt.Printf("Enter the function & params: ")
    var fn string
    var args [3]string // max args: 3 due to the chaincode
    // inputCnt - 1 = args Count
    inputCnt, _ := fmt.Scanln(&fn, &args[0], &args[1], &args[2])

    if inputCnt == 0 {
      continue
    } else if fn == "exit" {
      fmt.Println("bye")
      return
    }

    // else, invoke the smart contract
    if response, err := invoke(sdk, *channelID, *orgID, *orgUser, *chaincodeID, fn, args[1:inputCnt]); err != nil {
      log.Printf("invoke chaincode fail: %s\n", err.Error())
    } else {
      fmt.Println("Response: " + response)
    }
  }
}
