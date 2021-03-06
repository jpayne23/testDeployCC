package main

import (
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"errors"
	"strconv"
)

//==============================================================================================================================
//	 Structure Definitions 
//==============================================================================================================================
//	Chaincode - A blank struct for use with Shim (A HyperLedger included go file used for get/put state
//				and other HyperLedger functions)
//==============================================================================================================================
type  SimpleChaincode struct {
}

//==============================================================================================================================
//	Init Function - Called when the user deploys the chaincode sets up original number value, passed as an argument																
//==============================================================================================================================
func (t *SimpleChaincode) Init(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {

	//Args
	//				0
	//		start Number value

																
	err := stub.PutState("Number", []byte(args[0]))
	
															if err != nil { return nil, errors.New("INIT: Error putting state") }

	return nil, nil
}

//==============================================================================================================================
//	Invoke Functions
//==============================================================================================================================
//	add_number - Retrieves the current number value stored in the world state and adds a number passed by the invoker to it
//				and updates Number to the new value in the world state
//==============================================================================================================================
func (t *SimpleChaincode) add_number(stub *shim.ChaincodeStub, args []string) ([]byte, error) {


	//Args
	//				0
	//			Value to add

	adder, _ := strconv.Atoi(args[0])


	bytes, err := stub.GetState("Number")
		
															if err != nil { return nil, errors.New("Unable to get number") }
	
	number, _ := strconv.Atoi(string(bytes))
															
	
	newNumber := number + adder
	
	toPut := strconv.Itoa(newNumber)
	

	err = stub.PutState("Number", []byte(toPut))

															if err != nil { return nil, errors.New("Unable to put the state") }

	return nil, nil
}

//==============================================================================================================================
//	Query Functions
//==============================================================================================================================
//	get_number - Retrieves the current number value stored in the world state and returns it
//
//==============================================================================================================================
func (t *SimpleChaincode) get_number(stub *shim.ChaincodeStub, args []string) ([]byte, error) {

	bytes, err := stub.GetState("Number")
		
																			if err != nil { return nil, errors.New("Unable to get number") }

	return bytes, nil
	
}

//==============================================================================================================================
//	 Router Functions
//==============================================================================================================================
//	Invoke - Called on chaincode invoke. Takes a function name passed and calls that function.
//==============================================================================================================================
func (t *SimpleChaincode) Invoke(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {

	// Handle different functions
	if function == "add_number" {
		
		return t.add_number(stub, args) 
	}
	
	return nil, errors.New("Function of that name not found")
}
//==============================================================================================================================
//	Query - Called on chaincode query. Takes a function name passed and calls that function.
//==============================================================================================================================
func (t *SimpleChaincode) Query(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	
	if function == "get_number" { 
			return t.get_number(stub, args) 		
	}
	
	return nil, errors.New("Function of that name not found")
}

//=================================================================================================================================
//	 Main - main - Starts up the chaincode
//=================================================================================================================================
func main() {

	err := shim.Start(new(SimpleChaincode))
	
															if err != nil { fmt.Printf("Error starting SimpleChaincode: %s", err) }
}

