package main

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/golang-collections/go-datastructures/queue"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

type  instructionTransaction struct {
	transactionType		string `json:"transactionType"` 
	name			    string `json:"name"`
	issuer			    string `json:"issuer"`
	instruction		    string `json:"instruction"`
}

type responseTransaction struct {
	transactionType		string `json:"transactionType"`
	name			    string `json:"name"`
	issuer			    string `json:"issuer"`
	instructionId		string `json:"instructionId"`
	result			    string `json:"result"`
	errorString		    string `json:"errorString"`

}
var q queue.Queue
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
	if function == "initInstructionTransaction" { //create an instruction transaction
		return t.initInstructionTransaction(stub, args)
	} else if function == "initResponseTransaction" {
		return t.initResponseTransaction(stub, args)
	} else if function == "getPendingInstructionTransaction" { //get pending instruction transaction
		return t.getPendingInstructionTransaction(stub, args)
	} else if function == "getHistoryForTransaction" { //get history of values for a transaction
		return t.getHistoryForTransaction(stub, args)
	}

	fmt.Println("invoke did not find func: " + function) //error
	return shim.Error("Received unknown function invocation")
}

func (t *SimpleChaincode) initInstructionTransaction(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3.\n Usage: `{\"Args\":[\"[issuer]\",\"[instruction]\"],\"[contract name]\"}`")
	}

	// ==== Input sanitation ====
    //	fmt.Println("- start init instruction transaction")
	if len(args[0]) <= 0 {
		return shim.Error("1st argument must be a non-empty string")
	}
	if len(args[1]) <= 0 {
		return shim.Error("2nd argument must be a non-empty string")
	}
	if len(args[2]) <= 0 {
		return shim.Error("3rd argument must be a non-empty string")
	}

	var transactionType string = "instruction"
	issuer := strings.ToLower(args[1])
	name := strings.ToLower(args[0])
	instruction := strings.ToLower(args[2])

	// ==== Check if transaction already exists ====
	contractAsBytes, err := stub.GetState(name)
	if err != nil {
		return shim.Error("Failed to get contract: " + err.Error())
	} else if contractAsBytes != nil {
		fmt.Println("This name already exists: " + name)
		return shim.Error("This contract already exists: " + name)
	}


	contractJSONasString := `{"transactionType": "` + transactionType +`","issuer": "` + issuer +`","name": "` + name +`","instruction": "` + instruction +`"}`
	contractJSONasBytes:= []byte(contractJSONasString)


	// === Save transaction to state ===
	err = stub.PutState(name, contractJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

    //Get Transaction ID
	var transactionId string = stub.GetTxID()

    //Put transaction in queue of pending transactions identified by the transaction ID
	q.Put(transactionId)


	// ==== Transaction saved. Return success ====
	return shim.Success(nil)
}

func (t *SimpleChaincode) initResponseTransaction(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error

	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments. Expecting 5")
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

	var transactionType string = "response"
	name := strings.ToLower(args[0])
	issuer := strings.ToLower(args[1])
	instructionId := args[2]
	result := strings.ToLower (args[3])
	errorString := strings.ToLower (args[4])

	// ==== Check if transaction already exists ====
	contractAsBytes, err := stub.GetState(name)
	if err != nil {
		return shim.Error("Failed to get contract: " + err.Error())
	} else if contractAsBytes != nil {
		fmt.Println("This name already exists: " + name)
		return shim.Error("This contract already exists: " + name)
	}
	contractJSONasString := `{"transactionType": "` + transactionType +`","issuer": "` + issuer +`","name": "` + name +`","instructionId": "` + instructionId +`","result":"` + result +`","errorString": "` + errorString  + `"}`
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

	// buffer is a JSON array containing historic values for the the transaction
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

func (t *SimpleChaincode) getPendingInstructionTransaction (stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) < 0 {
		return shim.Error("Incorrect number of arguments. Expected 0 arguments")
	}

    //check if queue is not empty
	if (q.Len() == 0){
		return shim.Error("No instruction transactions to be processed!\n")
	}
    //get a transaction id from the queue
	results, err := q.Get(1)
	if err != nil {
		return shim.Error(err.Error())
	}
	var buffer bytes.Buffer
	buffer.WriteString(results[0].(string))
	return shim.Success(buffer.Bytes())
}

