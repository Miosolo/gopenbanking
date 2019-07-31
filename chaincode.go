package main

import (
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

// SimpleAsset implements a simple chaincode to manage an asset
type SimpleAsset struct {
}

// Init is called during chaincode instantiation to initialize any
// data. Note that chaincode upgrade also calls this function to reset
// or to migrate data.
func (t *SimpleAsset) Init(stub shim.ChaincodeStubInterface) peer.Response {
	// Get the args from the transaction proposal
	args := stub.GetStringArgs()
	if len(args) != 2 {
		return shim.Error("Incorrect arguments. Expecting a key and a value")
	}

	// Set up any variables or assets here by calling stub.PutState()

	// We store the key and the value on the ledger
	err := stub.PutState(args[0], []byte(args[1]))
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to create asset: %s", args[0]))
	}
	return shim.Success(nil)
}

// Invoke is called per transaction on the chaincode. Each transaction is
// either a 'get' or a 'set' on the asset created by Init function. The Set
// method may create a new asset by specifying a new key-value pair.
func (t *SimpleAsset) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	// Extract the function and args from the transaction proposal
	fn, args := stub.GetFunctionAndParameters()

	var result string
	var err error
	if fn == "set" {
		result, err = set(stub, args)
	} else if fn == "get" { // assume 'get' even if fn is nil
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
	}
	if err != nil {
		return shim.Error(err.Error())
	}

	// Return the result as success payload
	return shim.Success([]byte(result))
}

// Set stores the asset (both key and value) on the ledger. If the key exists,
// it will override the value with the new one
func set(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 2 {
		return "", fmt.Errorf("Incorrect arguments. Expecting an account name and a balance value.")
	}
	// set the account and its balance value to specific values.
	err := stub.PutState(args[0], []byte(args[1]))
	if err != nil {
		return "", fmt.Errorf("Failed to set asset: %s", args[0])
	}
	return "Account:" + args[0] + "Balance: " + string(value), nil
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
	return "Account:" + args[0] + "Balance: " + string(value), nil
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
	intValueTemp, err := strconv.Atoi(string(valueTemp))
	err = stub.PutState(args[0], int32TobyteArray(int32(intArgs1+intValueTemp)))
	if err != nil {
		return "", fmt.Errorf("Failed to set asset: %s", args[0])
	}
	return "Add is success!" + "Account: " + args[0] + "Remaining balance is:" + string(int32TobyteArray(int32(intArgs1+intValueTemp))), nil
}

// args[0] represents account, args[1] represents money.
// Reduce specific number of money to the specific account.
func reduce(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 2 {
		return "", fmt.Errorf("Incorrect arguments. Expecting an account name and a balance value.")
	}
	valueTemp, err := stub.GetState(args[0])
	if err != nil {
		return "", fmt.Errorf("Failed to get asset: %s with error: %s", args[0], err)
	}
	intArgs1, err := strconv.Atoi(args[1])
	intValueTemp, err := strconv.Atoi(string(valueTemp))
	if intArgs1 > intValueTemp {
		return "", fmt.Errorf("The balance in %s's account is not enough to reduce!", args[0])
	}
	err = stub.PutState(args[0], int32TobyteArray(int32(intValueTemp-intArgs1)))
	if err != nil {
		return "", fmt.Errorf("Failed to set asset: %s", args[0])
	}
	return "Reduce is success!" + "Account: " + args[0] + " Remaining balance is:" + string(int32TobyteArray(int32(intArgs1-intValueTemp))), nil
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
		return "", fmt.Errorf(fmt.Sprintf("Failed to create asset: %s", args[0]))
	}
	return "Create account: " + args[0] + "is success!", nil
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

// it is used to convert an int32 variable into a byte array.
func int32TobyteArray(core int32) []byte {
	var result []byte

	result[3] = uint8(core)
	result[2] = uint8(core >> 8)
	result[1] = uint8(core >> 16)
	result[0] = uint8(core >> 24)

	return result
}

// args[0] represents the debit account, args[1] represents the credit account, args[2] represents the money.
// transfer the money from the debit account to the credit account.
func transfer(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 3 {
		return "", fmt.Errorf("Incorrect arguments. Expecting a debit account, a credit account and a value")
	}
	//get the remaining balance from the debit account.
	valueTempD, err := stub.GetState(args[0])
	if err != nil {
		return "", fmt.Errorf("Failed to get asset: %s with error: %s", args[0], err)
	}
	intValueTempD, err := strconv.Atoi(string(valueTempD))
	// convert the money into an integer variable.
	intArgs2, err := strconv.Atoi(args[2])
	if intValueTempD < intArgs2 {
		return "", fmt.Errorf("The balance in %s's account is not enough to reduce!", args[0])
	}
	//Then we really reduce the balance from the debit account and put it into the credit account.
	err = stub.PutState(args[0], int32TobyteArray(int32(intValueTempD-intArgs2)))
	//get the remaining balance from the cebit account.
	valueTempC, err := stub.GetState(args[0])
	if err != nil {
		return "", fmt.Errorf("Failed to get asset: %s with error: %s", args[0], err)
	}
	intValueTempC, err := strconv.Atoi(string(valueTempC))
	// add the money to the credit account.
	err = stub.PutState(args[0], int32TobyteArray(int32(intArgs2+intValueTempC)))
	if err != nil {
		return "", fmt.Errorf("Failed to set asset: %s", args[0])
	}
	return "Transfer is success", nil
}

// main function starts up the chaincode in the container during instantiate
func main() {
	if err := shim.Start(new(SimpleAsset)); err != nil {
		fmt.Printf("Error starting SimpleAsset chaincode: %s", err)
	}
}
