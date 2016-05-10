package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"encoding/json"
)

const   STATE_TEMPLATE  	=  0
const   STATE_MANUFACTURE  	=  1
const   STATE_PRIVATE_OWNERSHIP =  2
const   STATE_LEASED_OUT 	=  3
const   STATE_BEING_SCRAPPED  	=  4

//==============================================================================================================================
//	 Structure Definitions 
//==============================================================================================================================
//	Chaincode - A blank struct for use with Shim (A HyperLedger included go file used for get/put state
//				and other HyperLedger functions)
//==============================================================================================================================
type  SimpleChaincode struct {
}

//==============================================================================================================================
//	Vehicle - Defines the structure for a car object. JSON on right tells it what JSON fields to map to
//			  that element when reading a JSON object into the struct e.g. JSON make -> Struct Make.
//==============================================================================================================================
type Vehicle struct {
	Make            string `json:"make"`
	Model           string `json:"model"`
	Reg             string `json:"reg"`
	VIN             int    `json:"VIN"`					
	Owner           string `json:"owner"`
	Scrapped        bool   `json:"scrapped"`
	Status          int    `json:"status"`
	Colour          string `json:"colour"`
	V5cID           string `json:"v5cID"`
	LeaseContractID string `json:"leaseContractID"`
}

//==============================================================================================================================
//	ECertResponse - Struct for storing the JSON response of retrieving an ECert. JSON OK -> Struct OK
//==============================================================================================================================
type ECertResponse struct {
	OK string `json:"OK"`
}					

//==============================================================================================================================
//	Init Function - Called when the user deploys the chaincode																	
//==============================================================================================================================
func (t *SimpleChaincode) Init(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	
	err := stub.PutState("Vehicle_Log_Address", []byte(args[0]))
	
															if err != nil { return nil, errors.New("Error storing vehicle log address") }
	
	return nil, nil
}

//==============================================================================================================================
//	 create_log - Invokes the function of event_code chaincode with the name 'chaincodeName' to log an
//					event.
//==============================================================================================================================
func (t *SimpleChaincode) create_log(stub *shim.ChaincodeStub, args []string) ([]byte, error) {	
																						
	chaincode_function := "create_vehicle_log"																																									
	chaincode_arguments := args

	vehicle_log_address, err := stub.GetState("Vehicle_Log_Address")
															if err != nil { return nil, errors.New("Error retrieving vehicle log address") }
	
	_, err = stub.InvokeChaincode(string(vehicle_log_address), chaincode_function, chaincode_arguments)
	
															if err != nil { return nil, errors.New("Failed to invoke vehicle_log_code") }
	
	return nil,nil
}


//==============================================================================================================================
//	 retrieve_v5c - Gets the state of the data at v5cID in the ledger then converts it from the stored 
//					JSON into the Vehicle struct for use in the contract. Returns the Vehcile struct.
//					Returns empty v if it errors.
//==============================================================================================================================
func (t *SimpleChaincode) retrieve_v5c(stub *shim.ChaincodeStub, v5cID string) (Vehicle, error) {
	
	var v Vehicle

	bytes, err := stub.GetState(v5cID)	;					
				
															if err != nil {	return v, errors.New("Error retrieving vehicle with v5cID = " + v5cID) }

	err = json.Unmarshal(bytes, &v)	;						

															if err != nil {	return v, errors.New("Corrupt vehicle record"+string(bytes))	}
	
	return v, nil
}

//==============================================================================================================================
// save_changes - Writes to the ledger the Vehicle struct passed in a JSON format. Uses the shim file's 
//				  method 'PutState'.
//==============================================================================================================================
func (t *SimpleChaincode) save_changes(stub *shim.ChaincodeStub, v Vehicle) (bool, error) {
	 
	bytes, err := json.Marshal(v)
	
																if err != nil { return false, errors.New("Error creating vehicle record") }

	err = stub.PutState(v.V5cID, bytes)
	
																if err != nil { return false, err }
	
	return true, nil
}

//==============================================================================================================================
//	 Router Functions
//==============================================================================================================================
//	Run - Called on chaincode invoke. Takes a function name passed and calls that function. Converts some
//		  initial arguments passed to other things for use in the called function e.g. name -> ecert
//==============================================================================================================================
func (t *SimpleChaincode) Invoke(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	
	
	if function == "create_vehicle" { return t.create_vehicle(stub, args[0], args[1])
	} else { 																				// If the function is not a create then there must be a car so we need to retrieve the car.
		
		argPos := 2
		
		if function == "scrap_vehicle" {																// If its a scrap vehicle then only two arguments are passed (no update value) all others have three arguments and the v5cID is expected in the last argument
			argPos = 1
		}
		
		v, err := t.retrieve_v5c(stub, args[argPos])
		
																							if err != nil { return nil, err }
																		
		if strings.Contains(function, "update") == false           && 
		   function 							!= "scrap_vehicle"    { 									// If the function is not an update or a scrappage it must be a transfer so we need to get the ecert of the recipient.
			
				if 		   function == "authority_to_manufacturer" { return t.authority_to_manufacturer(stub, v, args[0], args[1])
				} else if  function == "manufacturer_to_private"   { return t.manufacturer_to_private(stub, v, args[0], args[1])
				} else if  function == "private_to_private" 	   { return t.private_to_private(stub, v, args[0], args[1])
				} else if  function == "private_to_lease_company"  { return t.private_to_lease_company(stub, v, args[0], args[1])
				} else if  function == "lease_company_to_private"  { return t.lease_company_to_private(stub, v, args[0], args[1])
				} else if  function == "private_to_scrap_merchant" { return t.private_to_scrap_merchant(stub, v, args[0], args[1])
				}
			
		} else if function == "update_make"  	    { return t.update_make(stub, v, args[0], args[1])
		} else if function == "update_model"        { return t.update_model(stub, v, args[0], args[1])
		} else if function == "update_registration" { return t.update_registration(stub, v, args[0], args[1])
		} else if function == "update_vin" 			{ return t.update_vin(stub, v, args[0], args[1])
		} else if function == "update_colour" 		{ return t.update_colour(stub, v, args[0], args[1])
		} else if function == "scrap_vehicle" 		{ return t.scrap_vehicle(stub, v, args[0]) }
		
																						return nil, errors.New("Function of that name doesn't exist.")
			
	}
	
}
//=================================================================================================================================	
//	Query - Called on chaincode query. Takes a function name passed and calls that function. Passes the
//  		initial arguments passed are passed on to the called function.
//=================================================================================================================================	
func (t *SimpleChaincode) Query(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	
	if len(args) != 2 { return nil, errors.New("Incorrect number of arguments passed") }
																			
	v, err := t.retrieve_v5c(stub, args[1])
																							if err != nil { return nil, err }
																							
	if function == "get_all" { return t.get_all(stub, v, args[1], args[0]) }
	
																							return nil, errors.New("Received unknown function invocation")
}

//=================================================================================================================================
//	 Create Function
//=================================================================================================================================									
//	 Create Vehicle - Creates the initial JSON for the vehcile and then saves it to the ledger.									
//=================================================================================================================================
func (t *SimpleChaincode) create_vehicle(stub *shim.ChaincodeStub, caller_name string, v5cID string) ([]byte, error) {								

	var v Vehicle																																										
		
	v5c_ID         := "\"v5cID\":\""+v5cID+"\", "							// Variables to define the JSON
	vin            := "\"VIN\":0, "
	make           := "\"Make\":\"UNDEFINED\", "
	model          := "\"Model\":\"UNDEFINED\", "
	reg            := "\"Reg\":\"UNDEFINED\", "
	owner          := "\"Owner\":\""+caller_name+"\", "
	colour         := "\"Colour\":\"UNDEFINED\", "
	leaseContract  := "\"LeaseContractID\":\"UNDEFINED\", "
	status         := "\"Status\":0, "
	scrapped       := "\"Scrapped\":false"
	
	vehicle_json := "{"+v5c_ID+vin+make+model+reg+owner+colour+leaseContract+status+scrapped+"}" 	// Concatenates the variables to create the total JSON object
	
	//matched, err := regexp.Match("^[A-z][A-z][0-9]{7}", []byte(v5cID))  				// matched = true if the v5cID passed fits format of two letters followed by seven digits
	
																		//if err != nil { return nil, errors.New("Invalid v5cID") }
	
	//if 				v5c_ID  == "" 	 || 
	//				matched == false    {
																		//return nil, errors.New("Invalid v5cID provided")
	//}

	err := json.Unmarshal([]byte(vehicle_json), &v)							// Convert the JSON defined above into a vehicle object for go
	
																		if err != nil { return nil, errors.New("Invalid JSON object") }

	//record, err := stub.GetState(v.V5cID) 								// If not an error then a record exists so cant create a new car with this V5cID as it must be unique
	
																		//if record != nil { return nil, errors.New("Vehicle already exists") }
	
	_, err  = t.save_changes(stub, v)									
			
																		if err != nil { return nil, err }
																		
	_, err  = t.create_log(stub,[]string{ "Create",								// Type of event 
											"Create V5C",		// Event text
											v.V5cID, caller_name})	// Which car and who caused it
	
																		if err != nil { return nil, err }																	
	
	return nil, nil

}

//=================================================================================================================================
//	 Transfer Functions
//=================================================================================================================================
//	 authority_to_manufacturer
//=================================================================================================================================
func (t *SimpleChaincode) authority_to_manufacturer(stub *shim.ChaincodeStub, v Vehicle, caller_name string, recipient_name string) ([]byte, error) {
	
	if     	v.Status			== STATE_TEMPLATE      	   &&
			v.Owner				== caller_name 		   &&
			v.Scrapped			== false 				  {		// If the roles and users are ok 
	
					v.Owner  = recipient_name		// then make the owner the new owner
					v.Status = STATE_MANUFACTURE			// and mark it in the state of manufacture
	
	} else {									// Otherwise if there is an error
	
															return nil, errors.New("Permission Denied")
	
	}
	
	_, err := t.save_changes(stub, v)						// Write new state

															if err != nil {	return nil, err	}
	
											// Log the Event
														
	return nil, nil									// We are Done
	
}

//=================================================================================================================================
//	 manufacturer_to_private
//=================================================================================================================================
func (t *SimpleChaincode) manufacturer_to_private(stub *shim.ChaincodeStub, v Vehicle, caller_name string, recipient_name string) ([]byte, error) {
	
	if 		v.Make 	 == "UNDEFINED" || 					
			v.Model  == "UNDEFINED" || 
			v.Reg 	 == "UNDEFINED" || 
			v.Colour == "UNDEFINED" || 
			v.VIN == 0 				   {					//If any part of the car is undefined it has not bene fully manufacturered so cannot be sent
			
															return nil, errors.New("Car not fully defined")
	}
	
	if 		v.Status       == STATE_MANUFACTURE    && 
			v.Owner  == caller_name 	       &&  
			v.Scrapped     == false 				  {
			
					v.Owner = recipient_name
					v.Status = STATE_PRIVATE_OWNERSHIP
					
	} else {
															return nil, errors.New("Permission denied")
	}
	
	_, err := t.save_changes(stub, v)
	
															if err != nil { return nil, err }
	
	return nil, nil
	
}

//=================================================================================================================================
//	 private_to_private
//=================================================================================================================================
func (t *SimpleChaincode) private_to_private(stub *shim.ChaincodeStub, v Vehicle, caller_name string, recipient_name string) ([]byte, error) {
	
	if 		v.Status       == STATE_PRIVATE_OWNERSHIP &&
			v.Owner  == caller_name 		  &&
			v.Scrapped     == false 					 {
			
					v.Owner = recipient_name
					
	} else {
		
															return nil, errors.New("Permission denied")
	
	}
	
	_, err := t.save_changes(stub, v)
	
															if err != nil { return nil, err }
	
	return nil, nil
	
}

//=================================================================================================================================
//	 private_to_lease_company
//=================================================================================================================================
func (t *SimpleChaincode) private_to_lease_company(stub *shim.ChaincodeStub, v Vehicle, caller_name string, recipient_name string) ([]byte, error) {
	
	if 		v.Status       == STATE_PRIVATE_OWNERSHIP && 
			v.Owner  == caller_name 		  && 
			v.Scrapped     == false						 {
		
					v.Owner = recipient_name
					
	} else {
															return nil, errors.New("Permission denied")
	}
	
	_, err := t.save_changes(stub, v)
															if err != nil { return nil, err }
	
	return nil, nil
	
}

//=================================================================================================================================
//	 lease_company_to_private
//=================================================================================================================================
func (t *SimpleChaincode) lease_company_to_private(stub *shim.ChaincodeStub, v Vehicle, caller_name string, recipient_name string) ([]byte, error) {
	
	if		v.Status       == STATE_PRIVATE_OWNERSHIP &&
			v.Owner  == caller_name 		  &&  
			v.Scrapped	   == false					     {
		
				v.Owner = recipient_name
	
	} else {
															return nil, errors.New("Permission denied")
	}
	
	_, err := t.save_changes(stub, v)
															if err != nil { return nil, err }
	
	return nil, nil
	
}

//=================================================================================================================================
//	 private_to_scrap_merchant
//=================================================================================================================================
func (t *SimpleChaincode) private_to_scrap_merchant(stub *shim.ChaincodeStub, v Vehicle, caller_name string, recipient_name string) ([]byte, error) {
	
	if		v.Status       == STATE_PRIVATE_OWNERSHIP &&
			v.Owner  == caller_name 		  && 
			v.Scrapped 	   == false 					 {
			
					v.Owner = recipient_name
					v.Status = STATE_BEING_SCRAPPED
	
	} else {
		
															return nil, errors.New("Permission denied")
	
	}
	
	_, err := t.save_changes(stub, v)
	
															if err != nil { return nil, err }
	
	return nil, nil
	
}

//=================================================================================================================================
//	 Update Functions
//=================================================================================================================================
//	 update_vin
//=================================================================================================================================
func (t *SimpleChaincode) update_vin(stub *shim.ChaincodeStub, v Vehicle, caller_name string, new_value string) ([]byte, error) {
	
	new_vin, err := strconv.Atoi(string(new_value)) 		                // will return an error if the new vin contains non numerical chars
	
															if err != nil || len(string(new_value)) != 15 { return nil, errors.New("Invalid value passed for new VIN") }
	
	if 		v.Status       == STATE_MANUFACTURE    && 
			v.Owner  == caller_name 	       && 
			v.VIN          == 0 		       &&			// Can't change the VIN after its initial assignment
			v.Scrapped     == false 				  {
			
					v.VIN = new_vin					// Update to the new value
	} else {
	
															return nil, errors.New("Permission denied")
		
	}
	
	_, err  = t.save_changes(stub, v)						// Save the changes in the blockchain
	
															if err != nil { return nil, err } 
	
	return nil, nil
	
}


//=================================================================================================================================
//	 update_registration
//=================================================================================================================================
func (t *SimpleChaincode) update_registration(stub *shim.ChaincodeStub, v Vehicle, caller_name string, new_value string) ([]byte, error) {

	
	if		v.Owner  == caller_name 	       && 
			v.Scrapped     == false                   {
			
					v.Reg = new_value
	
	} else {
															return nil, errors.New("Permission denied")
	}
	
	_, err := t.save_changes(stub, v)
	
															if err != nil { return nil, err }
	
	return nil, nil
	
}

//=================================================================================================================================
//	 update_colour
//=================================================================================================================================
func (t *SimpleChaincode) update_colour(stub *shim.ChaincodeStub, v Vehicle, caller_name string, new_value string) ([]byte, error) {
	
	if 		v.Owner  == caller_name 	       && 
			v.Scrapped     == false  				  {
			
					v.Colour = new_value
	} else {
	
															return nil, errors.New("Permission denied")
	}
	
	_, err := t.save_changes(stub, v)
	
															if err != nil { return nil, err }
	
	return nil, nil
	
}

//=================================================================================================================================
//	 update_make
//=================================================================================================================================
func (t *SimpleChaincode) update_make(stub *shim.ChaincodeStub, v Vehicle, caller_name string, new_value string) ([]byte, error) {
	
	if 		v.Status       == STATE_MANUFACTURE    &&
			v.Owner  == caller_name 	       && 
			v.Scrapped     == false 				  {
			
					v.Make = new_value
	} else {
	
															return nil, errors.New("Permission denied")
	
	}
	
	_, err := t.save_changes(stub, v)
	
															if err != nil { return nil, err }
	
	return nil, nil
	
}

//=================================================================================================================================
//	 update_model
//=================================================================================================================================
func (t *SimpleChaincode) update_model(stub *shim.ChaincodeStub, v Vehicle, caller_name string, new_value string) ([]byte, error) {
	
	if 		v.Status       == STATE_MANUFACTURE    &&
			v.Owner  == caller_name          &&  
			v.Scrapped     == false 				  {
			
					v.Model = new_value
					
	} else {
															return nil, errors.New("Permission denied")
	}
	
	_, err := t.save_changes(stub, v)
	
															if err != nil { return nil, err }
	
	return nil, nil
	
}

//=================================================================================================================================
//	 scrap_vehicle
//=================================================================================================================================
func (t *SimpleChaincode) scrap_vehicle(stub *shim.ChaincodeStub, v Vehicle, caller_name string) ([]byte, error) {

	if		v.Status       == STATE_BEING_SCRAPPED && 
			v.Owner  == caller_name 	       && 
			v.Scrapped     == false 				  {
		
					v.Scrapped = true
				
	} else {
		return nil, errors.New("Permission denied")
	}
	
	_, err := t.save_changes(stub, v)
	
															if err != nil { return nil, err }
	
	return nil, nil
	
}

//=================================================================================================================================
//	 Read Functions
//=================================================================================================================================
//	 get_all
//=================================================================================================================================
func (t *SimpleChaincode) get_all(stub *shim.ChaincodeStub, v Vehicle, current_owner string, caller_name string) ([]byte, error) {
	
	bytes, err := json.Marshal(v)
	
																if err != nil { return nil, errors.New("Invalid vehicle object") }
	
	return bytes, nil

}

//=================================================================================================================================
//	 Main - main - Starts up the chaincode
//=================================================================================================================================
func main() {

	err := shim.Start(new(SimpleChaincode))
	
															if err != nil { fmt.Printf("Error starting Chaincode: %s", err) }
}
