package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/op/go-logging"

	"github.com/hyperledger/fabric/core/chaincode/lib/cid"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

// return an log object.
var (
	log = logging.MustGetLogger("CHAINCODE")

	// Everything except the message has a custom color
	// which is dependent on the log level. Many fields have a custom output
	// formatting too, eg. the time returns the hour down to the milli second.
	format = logging.MustStringFormatter(
		`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`)
	paramLength      map[string]int
	paramLengthError map[string]string
)

// SimpleAsset implements a simple chaincode to manage an asset
type SimpleAsset struct {
}

func init() {
	paramLength = make(map[string]int)
	paramLengthError = make(map[string]string)
	// init the dict for parameter length check
	paramLength["add"] = 2
	paramLengthError["add"] = "Incorrect arguments. Expecting an account name and a balance value."
	paramLength["create"] = 2
	paramLengthError["create"] = "Incorrect arguments. Expecting an unique account name and an initial balance value."
	paramLength["delete"] = 1
	paramLengthError["delete"] = "Incorrect arguments. Expecting an account being deleted."
	paramLength["get"] = 1
	paramLengthError["get"] = "Incorrect arguments. Expecting an account name."
	paramLength["query"] = 2
	paramLengthError["query"] = "Incorrect arguments. Expecting an objectType and an account."
	paramLength["reduce"] = 2
	paramLengthError["reduce"] = "Incorrect arguments. Expecting an account name and a balance value."
	paramLength["rollback"] = 3
	paramLengthError["rollback"] = "Incorrect arguments. Expecting a debit account, credit account and a transaction id."
	paramLength["transfer"] = 3
	paramLengthError["transfer"] = "Incorrect arguments. Expecting a debit account, a credit account and a value"
}

// Init is called during chaincode instantiation to initialize any data.
// Note that chaincode upgrade also calls this function to reset or to migrate data.
// When calls function Init, you can set an original account and its value.
func (t *SimpleAsset) Init(stub shim.ChaincodeStubInterface) peer.Response {
	// Get the args from the transaction proposal
	function, args := stub.GetFunctionAndParameters()
	//the first argument is in the variable "function"
	if function != "init" {
		log.Error("The first parameter needs to be a string: \"init\"")
		return shim.Error("The first parameter needs to be a string: \"init\"")
	}
	if len(args) != 0 {
		log.Error("Incorrect arguments. Expecting an account name and a balance value")
		return shim.Error("Incorrect arguments. Expecting an account name and a balance value")
	}
	// get clientIdentity of the one who calls the chaincode
	client, err := cid.New(stub)
	if err != nil {
		return shim.Error(fmt.Sprintf("Get client identity failed! With error: %s", err))
	}
	// get clientMSPID of the one who calls the chaincode
	mspid, err := client.GetMSPID()
	if err != nil {
		return shim.Error(fmt.Sprintf("Get client MSPID failed! With error: %s", err))
	}

	if mspid == "ANZBankMSP" {
		return shim.Success([]byte(fmt.Sprintf("Success to initialize!")))
	} else {
		return shim.Error("You do not have authority to get access to this function!")
	}

}

// Invoke is called per transaction on the chaincode. Each transaction is
// either a 'get' or a 'set' on the asset created by Init function. The Set
// method may create a new asset by specifying a new key-value pair.
func (t *SimpleAsset) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	// Extract the function and args from the transaction proposal
	fn, args := stub.GetFunctionAndParameters()

	// check the param length
	if expectedLen, ok := paramLength[fn]; !ok {
		return shim.Error("Undefined function")
	} else if expectedLen != len(args) {
		errStr, _ := paramLengthError[fn]
		return shim.Error(errStr)
	}

	var result string
	var err error

	// get clientIdentity of the one who calls the chaincode
	client, err := cid.New(stub)
	if err != nil {
		return shim.Error(fmt.Sprintf("Get client identity failed! With error: %s", err))
	}
	// get clientMSPID of the one who calls the chaincode
	mspid, err := client.GetMSPID()
	if err != nil {
		return shim.Error(fmt.Sprintf("Get client MSPID failed! With error: %s", err))
	}

	if mspid == "SuperviMSP" {
		// pass the params ASIS
		switch fn {
		case "query":
			result, err = query(stub, args)
		case "rollback":
			result, err = rollback(stub, args)
		default:
			return shim.Error("You do not have authority to get access to this function!")
		}
	} else {
		// add orgs to input
		fullAccount := func(account string) string {
			return account + "@" + mspid[:len(mspid)-3] // remove "MSP"
		}

		switch fn {
		case "get":
			result, err = get(stub, []string{fullAccount(args[0])})
		case "add":
			result, err = add(stub, []string{fullAccount(args[0]), args[1]})
		case "reduce":
			result, err = reduce(stub, []string{fullAccount(args[0]), args[1]})
		case "create":
			result, err = create(stub, []string{fullAccount(args[0]), args[1]})
		case "delete":
			result, err = delete(stub, []string{fullAccount(args[0])})
		case "transfer":
			result, err = transfer(stub, []string{fullAccount(args[0]), args[1], args[2]})
		case "query":
			result, err = query(stub, []string{args[0], fullAccount(args[1])})
		default:
			return shim.Error("You do not have authority to get access to this function!")
		}
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
	// get the account information from the database.
	value, err := stub.GetState(args[0])
	if err != nil {
		return "", fmt.Errorf("Failed to get asset: %s with error: %s", args[0], err)
	}
	if value == nil {
		return "", fmt.Errorf("Asset not found: %s", args[0])
	}

	return fmt.Sprintf(" Account: %s; Balance: %s", args[0], string(value)), nil
}

// args[0] represents account, args[1] represents money.
// Add specific number of money to the specific account.
func add(stub shim.ChaincodeStubInterface, args []string) (string, error) {
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

	return fmt.Sprintf("Add is success! Account: %s; Remaining balance is: %d", args[0], intValueTemp+intArgs1), nil

}

// args[0] represents account, args[1] represents money.
// Reduce specific number of money to the specific account.
func reduce(stub shim.ChaincodeStubInterface, args []string) (string, error) {
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

// The function of this module is to create an account of ledger
// args[0] means the account ID
// args[1] means the account initial value.
func create(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	var name []byte
	name, err := stub.GetState(args[0])
	if name != nil {
		return "", fmt.Errorf(fmt.Sprintf("The account has already existed!"))
	}
	if err != nil {
		return "", fmt.Errorf(fmt.Sprintf("Failed to get access to asset: %s; With error: %s", args[0], err))
	}

	// Set up any variables or assets here by calling stub.PutState()
	// We store the key and the value on the ledger
	err = stub.PutState(args[0], []byte(args[1]))
	if err != nil {
		return "", fmt.Errorf(fmt.Sprintf("Failed to create asset: %s; With Error: %s", args[0], err))
	}

	return fmt.Sprintf("Create account: %s  is success!", args[0]), nil

}

// delete an account of ledger.
// args[0] represents the account ID.
func delete(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	// delete the account.
	err := stub.DelState(args[0])
	if err != nil {
		return "", fmt.Errorf("Failed to delete asset: %s with error: %s", args[0], err)
	}

	return fmt.Sprintf("Delete is success! Account: %s", args[0]), nil
}

// args[0] represents the debit account
// args[1] represents the credit account
// args[2] represents the money.
// transfer the money from the debit account to the credit account.
func transfer(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	//reduce money from the debit account.
	var argsD []string = make([]string, 2)
	argsD[0] = args[0]
	argsD[1] = args[2]
	_, err := reduce(stub, argsD)
	if err != nil {
		return "", fmt.Errorf(fmt.Sprintf("Reduce debit account failed!"))
	}

	//add money to the cebit account.
	var argsC []string = make([]string, 2)
	argsC[0] = args[1]
	argsC[1] = args[2]
	_, err = add(stub, argsC)
	if err != nil {
		return "", fmt.Errorf(fmt.Sprintf("Add credit account failed!"))
	}
	// store the transfer record into the database
	// "out" means the money go out from one's account,
	// so the organization of the key-value pair is:
	// Key is a composite key, its sequence is ["out"debit account] [credit account] [uuid] [time]
	// value is the amount of money been transfered.
	msg, err := createHistoryKey(stub, args, "out")
	if err != nil {
		return "", fmt.Errorf("Create history records failed! with error: %s", err)
	}
	log.Info(msg)
	// store the transfer record into the database
	// "in" means the money go into one's account,
	// so the organization of the key-value pair is:
	// Key is a composite key, its sequence is ["in"credit account] [debit account] [uuid] [time]
	// value is the amount of money been transfered.
	msg, err = createHistoryKey(stub, args, "in")
	if err != nil {
		return "", fmt.Errorf("Create history records failed! with error: %s", err)
	}
	log.Info(msg)

	return fmt.Sprintf("Transfer is success!"), nil
}

// create history transferring records
// "out" means the money go out from one's account,
// "in" means the money go into one's account,
// both "out" and "in" is tags, they emphasize on going out or in records
func createHistoryKey(stub shim.ChaincodeStubInterface, args []string, first string) (string, error) {
	// get the time of the transaction been finished.
	FormatTime, err := stub.GetTxTimestamp()
	if err != nil {
		return "", fmt.Errorf(fmt.Sprintf("Get transaction timestamp failed!"))
	}
	tm := time.Unix(FormatTime.Seconds, 0)

	// if we need to create an "out" record
	// the organization of the key-value pair is:
	// Key is a composite key, its sequence is ["out"debit account] [credit account] [uuid] [time]
	// value is the amount of money been transfered.
	if first == "out" {
		historyKey, err := stub.CreateCompositeKey(first, []string{
			args[0], "->", args[1],
			"\t", stub.GetTxID(),
			"\t", tm.Format("Mon Jan 2 15:04:05 +0800 UTC 2006"),
		})
		if err != nil {
			return "", fmt.Errorf("Create historyKey failed! With error: %s", err)
		}

		err = stub.PutState(historyKey, []byte(args[2]))
		if err != nil {
			return "", fmt.Errorf("Store transfer information failed! With error: %s", err)
		}

	} else if first == "in" {
		// so the organization of the key-value pair is:
		// Key is a composite key, its sequence is ["in"credit account] [debit account] [uuid] [time]
		// value is the amount of money been transfered.
		historyKey, err := stub.CreateCompositeKey(first, []string{
			args[1], "<-", args[0],
			"\t", stub.GetTxID(),
			"\t", tm.Format("Mon Jan 2 15:04:05 +0800 UTC 2006"),
		})
		if err != nil {
			return "", fmt.Errorf("Create historyKey failed! With error: %s", err)
		}

		err = stub.PutState(historyKey, []byte(args[2]))
		if err != nil {
			return "", fmt.Errorf("Store transfer information failed! With error: %s", err)
		}
	}

	return fmt.Sprintf("Insert records success!"), nil
}

// query for the transferring history.
// args[0] represents the objectType, that is, "in" or "out"
// the variable "objectType" will store with the first argument of the composite key as one string.
// for example, if we store "Yongmao", "Songyue", "1", "10:01:10" with objectType "in",
// actually the string will be: inYongmao Songyue 1 10:01:10,
// every space is the seperator of each string.
// args[1] represents the account name
func query(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	var PCKey []string = make([]string, 1)
	PCKey[0] = args[1]
	// intend to get the record of transferring
	it, err := stub.GetStateByPartialCompositeKey(args[0], PCKey)
	if err != nil {
		return "", fmt.Errorf(fmt.Sprintf("Cannot get by partial composite key!"))
	}

	defer it.Close()
	// result contains all the appropriate results
	result := ""
	header := "AccountAccociation | ID | Time | Amount"
	if args[0] != "in" && args[0] != "out" {
		return "", fmt.Errorf(fmt.Sprintf("You have typed a wrong objectType!"))
	}

	for it.HasNext() {
		item, err := it.Next()
		if err != nil {
			return "", fmt.Errorf(fmt.Sprintf("Get next of iterator failed!"))
		}
		log.Info(fmt.Sprintf("%s %s", item.GetKey(), item.GetValue()))
		result = result + fmt.Sprintf("%s\t%s\n", item.GetKey()[len(args[0]):], item.GetValue()) // omit "in" / "out"
	}

	if result == "" {
		return "", fmt.Errorf("Do not have any records!")
	} else {
		return fmt.Sprintf("Query success! The entries obeys <%s>:\n%s", header, result), nil
	}
}

// the supervisor can rollback the transferring operation
// args[0] represents debit account in transferring record
// args[1] represents credit account in transferring record
// args[2] represents transaction id in transferring record
func rollback(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	// get satisfied out record
	var PCKeyOut []string = make([]string, 1)
	PCKeyOut[0] = args[0]
	itOut, err := stub.GetStateByPartialCompositeKey("out", PCKeyOut)
	if err != nil {
		return "", fmt.Errorf(fmt.Sprintf("Cannot get by partial composite key when get \"in\" record!"))
	}
	//get money value and delete "out" record
	defer itOut.Close()
	var money []byte
	if itOut.HasNext() == false {
		return "", fmt.Errorf(fmt.Sprintf("Database do not have such records! Please check you arguments!"))
	}
	for itOut.HasNext() {
		item, err := itOut.Next()
		if err != nil {
			return "", fmt.Errorf(fmt.Sprintf("Get next of iteratorOut failed!"))
		}
		log.Info(fmt.Sprintf("%s %s", item.GetKey(), item.GetValue()))
		// get attribute from composite key
		_, attrArray, err := stub.SplitCompositeKey(item.GetKey())
		if err != nil {
			return "", fmt.Errorf(fmt.Sprintf("Split composite key failed!"))
		}
		// compare the input hash code with the hash code stored in database
		IsThisOne := strings.Compare(attrArray[4], args[2])
		if IsThisOne == 0 {
			money = item.GetValue()
			stub.DelState(item.GetKey())
			break
		}
	}
	// delete "in" record
	var PCKeyIn []string = make([]string, 1)
	PCKeyIn[0] = args[1]
	itIn, err := stub.GetStateByPartialCompositeKey("in", PCKeyIn)
	if err != nil {
		return "", fmt.Errorf(fmt.Sprintf("Cannot get by partial composite key when get \"in\" record!"))
	}

	defer itIn.Close()
	if itIn.HasNext() == false {
		return "", fmt.Errorf(fmt.Sprintf("Database do not have such records! Please check you arguments!"))
	}
	for itIn.HasNext() {
		item, err := itIn.Next()
		if err != nil {
			return "", fmt.Errorf(fmt.Sprintf("Get next of iteratorIn failed!"))
		}
		log.Info(fmt.Sprintf("%s %s", item.GetKey(), item.GetValue()))
		// get attribute from composite key
		_, attrArray, err := stub.SplitCompositeKey(item.GetKey())
		if err != nil {
			return "", fmt.Errorf(fmt.Sprintf("Split composite key failed!"))
		}
		// compare the input hash code with the hash code stored in database
		IsThisOne := strings.Compare(attrArray[4], args[2])
		if IsThisOne == 0 {
			stub.DelState(item.GetKey())
			break
		}
	}

	// Then we should put money back into debit account.
	//reduce money from the debit account.
	var argsD []string = make([]string, 2)
	argsD[0] = args[1]
	argsD[1] = string(money)
	_, err = reduce(stub, argsD)
	if err != nil {
		return "", fmt.Errorf(fmt.Sprintf("Reduce debit account failed! With error: %s", err))
	}

	//add money to the cebit account.
	var argsC []string = make([]string, 2)
	argsC[0] = args[0]
	argsC[1] = string(money)
	_, err = add(stub, argsC)
	if err != nil {
		return "", fmt.Errorf(fmt.Sprintf("Add cebit account failed! With error: %s", err))
	}

	return fmt.Sprintf("rollback Success!"), nil
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
		fmt.Printf("Error starting SimpleAsset chaincode: %s\n", err)
	}
}
