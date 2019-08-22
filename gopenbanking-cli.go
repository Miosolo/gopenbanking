package main

import (
	"flag"
	"fmt"

	"github.com/Miosolo/gopenbanking/app"
)

// provides an interactive cli interface to multi-org users
func main() {
	// define the flags & parse the params
	channelID := flag.String("chan", "orgschannel", `Name of the channel`)
	orgID := flag.String("org", "", "Name of your orgnization")
	orgUser := flag.String("user", "", `Your User ID in this organization`)
	chaincodeID := flag.String("cc", "cc_gopenbanking", "ID of the chaincode instanciated")
	configPath := flag.String("conf", "app/config.yaml", "path of app configeration config.yaml")
	cryptoPath := flag.String("crypto", "crypto-config", "path of crypto-config")
	flag.Parse()

	ap, err := app.New(*channelID, *orgID, *orgUser, *chaincodeID, *configPath, *cryptoPath)
	if err != nil {
		fmt.Println("Cannot start up the app: " + err.Error())
		return
	}

	// print the instructions
	if *orgID == "Supervisor" {
		fmt.Println(`==========INSTRUCTIONS==========
  - "rollback" + debit account + credit account + transaction ID
  - "query" + "in" / "out" + account
  - "exit": terminate the loop and exit
================================`)
	} else { // banks
		fmt.Println(`==========INSTRUCTIONS==========
Functions and parameters of the ANZ-CITI Banking Network:
  - "get" + account
  - "add" + account + value
  - "reduce" + account + value
  - "create" + account + inititial value
  - "delete" + account
  - "tranfer" + this side's account + full account (this side/ other sides) + tranfer amount
  - "query" + "in" / "out" + account
  - "exit": terminate the loop and exit
================================`)
	}

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
		if response, err := ap.Invoke(fn, args[0:inputCnt-1]); err != nil {
			fmt.Println("Invoking chaincode failed: " + err.Error())
		} else {
			fmt.Println("Response: " + response)
		}
	}
}
