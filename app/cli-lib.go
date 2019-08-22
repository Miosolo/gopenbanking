package app

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	clientmsp "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

var (
	// global map for orgName - domain reflection
	domainMap map[string]string
)

// Provider (app.Provider) contains the identity info & app running stubs
type Provider struct {
	channelID, orgID, orgUser, chaincodeID, // network parameters
	configPath, cryptoPath string // app config
	sdk *fabsdk.FabricSDK // SDK stub
}

// init domainMap with constants
func init() {
  domainMap = make(map[string]string)
	domainMap["ANZBank"] = "anz.italktoyou.cn"
	domainMap["CitiBank"] = "citi.italktoyou.cn"
	domainMap["Supervisor"] = "supervi.italktoyou.cn"
}

// New creates a new app.Provider instance & check the identity
func New(channelID, orgID, orgUser, chaincodeID, configPath, cryptoPath string) (p *Provider, err error) {
	// init app provider & its members
	ap := Provider{
		channelID:   channelID,
		orgID:       orgID,
		orgUser:     orgUser,
		chaincodeID: chaincodeID,
		configPath:  configPath,
		cryptoPath:  cryptoPath}

	// set env for YAML parsing
	os.Setenv("FABRIC_CRYPTOCONFIG_ROOT", os.Getenv("PWD")+"/"+cryptoPath)
	if domain, ok := domainMap[orgID]; ok {
		os.Setenv("FABRIC_CRYPTOCONFIG_USER",
			os.Getenv("FABRIC_CRYPTOCONFIG_ROOT")+"/peerOrganizations/"+
				domain+"/users/"+orgUser+"@"+domain+"/msp/")
		os.Setenv("FABRIC_ORG_ID", orgID)
	} else {
		return nil, errors.New("invalid organization")
	}

	// init the env
	configProvider := config.FromFile(configPath)
	ap.sdk, err = fabsdk.New(configProvider)
	if err != nil {
		return nil, fmt.Errorf("failed to create new SDK: %s", err)
	}

	// identify the org & role
	if err := ap.identify(); err != nil {
		return nil, fmt.Errorf("identify %s fail: %s", orgUser, err.Error())
	}

	return &ap, nil
}

// identify checks the user identity
func (ap Provider) identify() (err error) {
	mspClient, err := clientmsp.New(ap.sdk.Context(), clientmsp.WithOrg(ap.orgID))
	if err != nil {
		log.Printf("create msp client fail: %s\n", err.Error())
		return err
	}

	identity, err := mspClient.GetSigningIdentity(ap.orgUser)
	if err != nil {
		return err
	}

	log.Println("using identity: " + identity.Identifier().MSPID)
	return nil
}

// Invoke connects to the channel, makes up a transaction request,
// and handles the response
func (ap Provider) Invoke(ccFunction string, args []string) (resp string, err error) {
	channelProvider := ap.sdk.ChannelContext(ap.channelID,
		fabsdk.WithUser(ap.orgUser),
		fabsdk.WithOrg(ap.orgID))

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
		ChaincodeID: ap.chaincodeID,
		Fcn:         ccFunction,
		Args:        byteArgs,
	}

  var response channel.Response
	if ccFunction == "query" || ccFunction == "get" {
		response, err = channelClient.Query(request)
	} else {
		response, err = channelClient.Execute(request)
	}

	if err != nil {
		log.Println("operation fail: ", err.Error())
		return "", err
	}

	return string(response.Payload), nil
}
