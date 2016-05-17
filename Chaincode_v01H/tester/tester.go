package main

import (
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"errors"
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
//	Init Function - Called when the user deploys the chaincode															
//==============================================================================================================================
func (t *SimpleChaincode) Init(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
														
	return nil, nil
}

//==============================================================================================================================
//	Invoke Functions
//==============================================================================================================================
//	add_number - Invokes the numbers chaincode and calls the function add_number, chaincode name currently hardcoded
//
//==============================================================================================================================
func (t *SimpleChaincode) add_number(stub *shim.ChaincodeStub, args []string) ([]byte, error) {

	_, err := stub.InvokeChaincode("eeec8f203c1249eb1a9f797e42cffd6d4f388abcbb0cefd99c50f58d417f1f6bbdacf8c0c58a343943c50d557e756f93e08a35e66a9769467fb56efd2c2bac70","add_number",args)

	if err != nil { return nil, errors.New("Unable to invoke chaincode") }

	return nil, nil
}

//==============================================================================================================================
//	get_number - Queries the numbers chaincode and calls get_number, chaincode name currently hardcoded.
//				 Returns the current number value stored in the world state
//==============================================================================================================================
func (t *SimpleChaincode) get_number(stub *shim.ChaincodeStub, args []string) ([]byte, error) {

	res, err := stub.QueryChaincode("eeec8f203c1249eb1a9f797e42cffd6d4f388abcbb0cefd99c50f58d417f1f6bbdacf8c0c58a343943c50d557e756f93e08a35e66a9769467fb56efd2c2bac70","get_number",args)
		
																			if err != nil { return nil, errors.New("Unable to query chaincode") }														

	return res, nil
	
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
	
	// Handle different functions
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

