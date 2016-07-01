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

	_, err := stub.InvokeChaincode("6198f2344e2d949a24df601c8826315d030d3111b329fce1fee72ed994f2fcd88da0ab395e1cc30450a081e3e9f5ccf91bf0c717a093271b40feba619fdb5ab3","add_number",args)

	if err != nil { return nil, errors.New("Unable to invoke chaincode") }

	return nil, nil
}

//==============================================================================================================================
//	get_number - Queries the numbers chaincode and calls get_number, chaincode name currently hardcoded.
//				 Returns the current number value stored in the world state
//==============================================================================================================================
func (t *SimpleChaincode) get_number(stub *shim.ChaincodeStub, args []string) ([]byte, error) {

	res, err := stub.QueryChaincode("6198f2344e2d949a24df601c8826315d030d3111b329fce1fee72ed994f2fcd88da0ab395e1cc30450a081e3e9f5ccf91bf0c717a093271b40feba619fdb5ab3","get_number",args)
		
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

