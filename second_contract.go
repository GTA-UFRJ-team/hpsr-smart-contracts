package main

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"time"


	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

type  configurationTransaction struct {
	transactionType	            string `json:"transactionType"`
	name			            string `json:"name"`
	issuer			            string `json:"issuer"`
	configurationIdentifier		string `json:"configurationIdentifier"`
    versionIdentifier           string `json:"versionIdentifier"`
    description                 string `json:"description"`
    configuration               string `json:"configuration"`
}

type configurationRequestTransaction struct {
    transactionType	            string `json:"transactionType"`
	name			            string `json:"name"`
	issuer			            string `json:"issuer"`
	configurationIdentifier		string `json:"configurationIdentifier"`
    versionIdentifier           string `json:"versionIdentifier"`
    description                 string `json:"description"`
    nextRecipient               string `json:"nextRecipient"`
}
// ===================================================================================
// Main
// ===================================================================================
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Init initializes chaincode
// ===========================
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

// Invoke - Our entry point for Invocations
// ========================================
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()

	// Handle different functions
	if function == "initConfigurationTransaction"{ //create a configuration transaction
		return t.initInstructionTransaction(stub, args)
	} else if function == "initConfigurationRequestTransaction" {
		return t.initResponseTransaction(stub, args)
	} else if function == "getHistoryForTransaction" { //get history for transaction
		return t.getHistoryForContract(stub, args)
	}

	fmt.Println("invoke did not find func: " + function) //error
	return shim.Error("Received unknown function invocation")
}

func (t *SimpleChaincode) initConfigurationTransaction(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error

	if len(args) != 6 {
		return shim.Error("Incorrect number of arguments. Expecting 3.\n Usage: `{\"Args\":[\"[issuer]\",\"[transaction name]\",\"[configuration identifier]\",\"[version identifier]\",\"[description]\",\"[configuration]\"]}`")
	}

	if len(args[0]) <= 0 {
		return shim.Error("1st argument must be a non-empty string")
	}
	if len(args[1]) <= 0 {
		return shim.Error("2nd argument must be a non-empty string")
	}
	if len(args[2]) <= 0 {
		return shim.Error("3rd argument must be a non-empty string")
	}
	if len(args[3]) <= 0 {
		return shim.Error("4th argument must be a non-empty string")
	}
	if len(args[4]) <= 0 {
		return shim.Error("5th argument must be a non-empty string")
	}
	if len(args[5]) <= 0 {
		return shim.Error("6th argument must be a non-empty string")
	}

	var transactionType string = "configuration"
	issuer := strings.ToLower(args[1])
	name := strings.ToLower(args[0])
	configurationIdentifier := strings.ToLower(args[2])
    versionIdentifier := strings.ToLower(args[3])
    description := strings.ToLower(args[4])
    configuration := strings.ToLower(args[5])

	// ==== Check if transaction already exists ====
	contractAsBytes, err := stub.GetState(name)
	if err != nil {
		return shim.Error("Failed to get contract: " + err.Error())
	} else if contractAsBytes != nil {
		fmt.Println("This name already exists: " + name)
		return shim.Error("This contract already exists: " + name)
	}


    contractJSONasString := `{"transactionType": "` + transactionType +`","issuer": "` + issuer +`","name": "` + name +`","configurationIdentifier": "` + configurationIdentifier +`","versionIdentifier": "` +versionIdentifier +`","description": "` + description  + `", "configuration": "` + configuration  + `"}`
	contractJSONasBytes:= []byte(contractJSONasString)


	// === Save transaction to state ===
	err = stub.PutState(name, contractJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	// ==== Contract saved. Return success ====
	return shim.Success(nil)
}
func (t *SimpleChaincode) initConfigurationRequestTransaction(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error

	if len(args) != 6 {
		return shim.Error("Incorrect number of arguments. Expecting 3.\n Usage: `{\"Args\":[\"[issuer]\",\"[transaction name]\",\"[configuration identifier]\",\"[version identifier]\",\"[description]\",\"[next recipient]\"]}`")
	}

	if len(args[0]) <= 0 {
		return shim.Error("1st argument must be a non-empty string")
	}
	if len(args[1]) <= 0 {
		return shim.Error("2nd argument must be a non-empty string")
	}
	if len(args[2]) <= 0 {
		return shim.Error("3rd argument must be a non-empty string")
	}
	if len(args[3]) <= 0 {
		return shim.Error("4th argument must be a non-empty string")
	}
	if len(args[4]) <= 0 {
		return shim.Error("5th argument must be a non-empty string")
	}
	if len(args[5]) <= 0 {
		return shim.Error("6th argument must be a non-empty string")
	}

	var transactionType string = "configuration request"
	issuer := strings.ToLower(args[1])
	name := strings.ToLower(args[0])
	configurationIdentifier := strings.ToLower(args[2])
    versionIdentifier := strings.ToLower(args[3])
    description := strings.ToLower(args[4])
    nextRecipient := strings.ToLower(args[5])

	// ==== Check if transaction already exists ====
	contractAsBytes, err := stub.GetState(name)
	if err != nil {
		return shim.Error("Failed to get contract: " + err.Error())
	} else if contractAsBytes != nil {
		fmt.Println("This name already exists: " + name)
		return shim.Error("This contract already exists: " + name)
	}


    contractJSONasString := `{"transactionType": "` + transactionType +`","issuer": "` + issuer +`","name": "` + name +`","configurationIdentifier": "` + configurationIdentifier +`","versionIdentifier": "` +versionIdentifier +`","description": "` + description  + `", "nextRecipient": "` + nextRecipient  + `"}`
	contractJSONasBytes:= []byte(contractJSONasString)


	// === Save transaction to state ===
	err = stub.PutState(name, contractJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}


	// ==== Transaction saved. Return success ====
	return shim.Success(nil)
}


func (t *SimpleChaincode) getHistoryForTransaction(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	name := args[0]

	resultsIterator, err := stub.GetHistoryForKey(name)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing historic values for the marble
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"TxId\":")
		buffer.WriteString("\"")
		buffer.WriteString(response.TxId)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Value\":")
		// if it was a delete operation on given key, then we need to set the
		//corresponding value null. Else, we will write the response.Value
		//as-is (as the Value itself a JSON marble)
		if response.IsDelete {
			buffer.WriteString("null")
		} else {
			buffer.WriteString(string(response.Value))
		}

		buffer.WriteString(", \"Timestamp\":")
		buffer.WriteString("\"")
		buffer.WriteString(time.Unix(response.Timestamp.Seconds, int64(response.Timestamp.Nanos)).String())
		buffer.WriteString("\"")

		buffer.WriteString(", \"IsDelete\":")
		buffer.WriteString("\"")
		buffer.WriteString(strconv.FormatBool(response.IsDelete))
		buffer.WriteString("\"")

		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")


	return shim.Success(buffer.Bytes())
}


