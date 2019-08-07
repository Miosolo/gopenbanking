package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/op/go-logging"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

// return an log object.
var log = logging.MustGetLogger("CHAINCODE")

// Everything except the message has a custom color
// which is dependent on the log level. Many fields have a custom output
// formatting too, eg. the time returns the hour down to the milli second.
var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
)

// SimpleAsset implements a simple chaincode to manage an asset
type SimpleAsset struct {
}

// Init is called during chaincode instantiation to initialize any data.
// Note that chaincode upgrade also calls this function to reset or to migrate data.
// When calls function Init, you can set an original account and its value.
func (t *SimpleAsset) Init(stub shim.ChaincodeStubInterface) peer.Response {
	// Get the args from the transaction proposal
	function, args := stub.GetFunctionAndParameters()
	if function != "init" {
		log.Error("The first parameter needs to be a string: \"init\"")
		return shim.Error("The first parameter needs to be a string: \"init\"")
	}
	if len(args) != 2 {
		log.Error("Incorrect arguments. Expecting an account name and a balance value")
		return shim.Error("Incorrect arguments. Expecting an account name and a balance value")
	}

	// Set up any variables or assets here by calling stub.PutState()
	// We store the key and the value on the ledger
	err := stub.PutState(args[0], []byte(args[1]))
	if err != nil {
		log.Debug(fmt.Sprintf("not found orgAttribute"))
		return shim.Error(fmt.Sprintf("Failed to create asset: %s", args[0]))
	}
	return shim.Success([]byte(fmt.Sprintf("Success to create one account! Account: %s; value: %s", args[0], args[1])))

}

// Invoke is called per transaction on the chaincode. Each transaction is
// either a 'get' or a 'set' on the asset created by Init function. The Set
// method may create a new asset by specifying a new key-value pair.
func (t *SimpleAsset) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	// Extract the function and args from the transaction proposal
	fn, args := stub.GetFunctionAndParameters()

	var result string
	var err error

	if fn == "get" {
		result, err = get(stub, args)
	} else if fn == "add" {
		result, err = add(stub, args)
	} else if fn == "reduce" {
		result, err = reduce(stub, args)
	} else if fn == "create" {
		result, err = create(stub, args)
	} else if fn == "delete" {
		result, err = delete(stub, args)
	} else if fn == "transfer" {
		result, err = transfer(stub, args)
	} else if fn == "query" {
		result, err = query(stub, args)
	}

	if err != nil {
		log.Error(err.Error())
		return shim.Error(err.Error())
	}

	// Return the result as success payload
	return shim.Success([]byte(result))
}

// Get returns the value of the specified asset key
// When we need to query the remaining balance, we use this function.
func get(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 1 {
		return "", fmt.Errorf("Incorrect arguments. Expecting an account name.")
	}
	// get the account information from the database.
	value, err := stub.GetState(args[0])
	if err != nil {
		return "", fmt.Errorf("Failed to get asset: %s with error: %s", args[0], err)
	}
	if value == nil {
		return "", fmt.Errorf("Asset not found: %s", args[0])
	}

	//set the output string's color to be green.
	return fmt.Sprintf(" Account: %s; Balance: %s", args[0], string(value)), nil
}

// args[0] represents account, args[1] represents money.
// Add specific number of money to the specific account.
func add(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 2 {
		return "", fmt.Errorf("Incorrect arguments. Expecting an account name and a balance value.")
	}

	valueTemp, err := stub.GetState(args[0])
	if err != nil {
		return "", fmt.Errorf("Failed to get asset: %s with error: %s", args[0], err)
	}

	intArgs1, err := strconv.Atoi(args[1])
	if err != nil {
		return "", fmt.Errorf("Atoi fail! With Error: %s", err)
	}

	intValueTemp, err := strconv.Atoi(string(valueTemp))
	if err != nil {
		return "", fmt.Errorf("Atoi fail! With Error: %s", err)
	}

	err = stub.PutState(args[0], []byte(strconv.Itoa(intArgs1+intValueTemp)))
	if err != nil {
		return "", fmt.Errorf("Failed to set asset: %s with error: %s", args[0], err)
	}

	return fmt.Sprintf("Add is success! Account: %s; Remaining balance is: %d", args[0], intArgs1+intValueTemp), nil
}

// args[0] represents account, args[1] represents money.
// Reduce specific number of money to the specific account.
func reduce(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 2 {
		return "", fmt.Errorf("Incorrect arguments. Expecting an account name and a balance value.")
	}
	// Get the account from the worldstate database.
	valueTemp, err := stub.GetState(args[0])
	if err != nil {
		return "", fmt.Errorf("Failed to get asset: %s with error: %s", args[0], err)
	}
	// change the argument into integer.
	intArgs1, err := strconv.Atoi(args[1])
	if err != nil {
		return "", fmt.Errorf("Atoi fail! With Error: %s", err)
	}
	//
	intValueTemp, err := strconv.Atoi(string(valueTemp))
	if err != nil {
		return "", fmt.Errorf("Atoi fail! With Error: %s", err)
	}

	if intArgs1 > intValueTemp {
		return "", fmt.Errorf("The balance in %s's account is not enough to reduce!", args[0])
	}

	err = stub.PutState(args[0], []byte(strconv.Itoa(intValueTemp-intArgs1)))
	if err != nil {
		return "", fmt.Errorf("Failed to set asset: %s;  With Error: %s", args[0], err)
	}

	return fmt.Sprintf("Reduce is success! Account: %s; Remaining balance is: %d", args[0], intValueTemp-intArgs1), nil
}

// create an account of ledger, args[0] means the account ID, args[1] means the account initial value.
func create(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 2 {
		return "", fmt.Errorf("Incorrect arguments. Expecting an unique account name and an initial balance value.")
	}

	// Set up any variables or assets here by calling stub.PutState()
	// We store the key and the value on the ledger
	err := stub.PutState(args[0], []byte(args[1]))
	if err != nil {
		return "", fmt.Errorf(fmt.Sprintf("Failed to create asset: %s; With Error: %s", args[0], err))
	}

	return fmt.Sprintf("Create account: %s  is success!", args[0]), nil
}

// delete an account of ledger. args[0] represents the account ID.
func delete(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 1 {
		return "", fmt.Errorf("Incorrect arguments. Expecting an account being deleted.")
	}
	// delete the account.
	err := stub.DelState(args[0])
	if err != nil {
		return "", fmt.Errorf("Failed to delete asset: %s with error: %s", args[0], err)
	}

	return "Delete" + args[0] + "is success!", nil
}

// args[0] represents the debit account, args[1] represents the credit account, args[2] represents the money.
// transfer the money from the debit account to the credit account.
func transfer(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 3 {
		return "", fmt.Errorf("Incorrect arguments. Expecting a debit account, a credit account and a value")
	}

	//reduce money from the debit account.
	var argsD []string = make([]string, 2)
	argsD[0] = args[0]
	args[1] = args[2]
	reduce(stub, argsD)

	//add money to the cebit account.
	var argsC []string = make([]string, 2)
	argsC[0] = args[1]
	argsC[1] = args[2]
	add(stub, argsC)

	FormatTime, err := stub.GetTxTimestamp()
	tm := time.Unix(FormatTime.Seconds, 0)
	historyKey, err := stub.CreateCompositeKey("history", []string{
		args[0],
		args[1],
		stub.GetTxID(),
		tm.Format("2019-08-06 08:08:08 PM"),
	})
	if err != nil {
		return "", fmt.Errorf("Create historyKey failed! With error: %s", err)
	}

	err = stub.PutState(historyKey, []byte(args[2]))
	if err != nil {
		return "", fmt.Errorf("Store transfer information failed! With error: %s", err)
	}

	return fmt.Sprintf("Transfer is success!"), nil
}

// query for the history
func query(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 1 {
		return "", fmt.Errorf("Incorrect arguments. Expecting an account")
	}
	var PCKey []string
	PCKey[0] = args[0]
	it, err := stub.GetStateByPartialCompositeKey("history", PCKey)
	if err != nil {
		return "", fmt.Errorf(fmt.Sprintf("Cannot get by partial composite key!"))
	}

	defer it.Close()
	for it.HasNext() {
		item, err := it.Next()
		if err != nil {
			return "", fmt.Errorf(fmt.Sprintf("Get next of iterator failed!"))
		}
		log.Info(fmt.Sprintf("%s %s", item.GetKey(), item.GetValue()))
	}
	return fmt.Sprintf("Query success!"), nil
}

// main function starts up the chaincode in the container during instantiate
func main() {

	// For demo purposes, create two backend for os.Stderr.
	backend1 := logging.NewLogBackend(os.Stderr, "", 0)
	backend2 := logging.NewLogBackend(os.Stderr, "", 0)

	// For messages written to backend2 we want to add some additional
	// information to the output, including the used log level and the name of
	// the function.
	backend2Formatter := logging.NewBackendFormatter(backend2, format)

	// Only errors and more severe messages should be sent to backend1
	backend1Leveled := logging.AddModuleLevel(backend1)
	backend1Leveled.SetLevel(logging.ERROR, "")

	// Set the backends to be used.
	logging.SetBackend(backend1Leveled, backend2Formatter)
	if err := shim.Start(new(SimpleAsset)); err != nil {
		fmt.Printf("Error starting SimpleAsset chaincode: %s", err)
	}
}
