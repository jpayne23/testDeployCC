package main

import (
	"errors"
	"fmt"
	"strconv"
	"github.com/openblockchain/obc-peer/openchain/chaincode/shim"
)

type Vehicle struct {
	
}
func (t *Vehicle) init(stub *shim.ChaincodeStub, args []string) ([]byte, error) {

	var make string
	var model string
	var reg string
	var VIN int
	var owner string
	var scrapped bool
	var status int
	var colour string
	var v5cID string
	var leaseContractID string
	var err error

	make = "UNDEFINED"
	model = "UNDEFINED"
	reg = "UNDEFINED"
	VIN = 0
	colour = "UNDEFINED"
	owner = "DVLA"
	status = 0
	v5cID = args[0]
	leaseContractID = "UNDEFINED"
	scrapped = false
	
	err = stub.PutState("Make", []byte(make))
	err = stub.PutState("Model", []byte(model))
	err = stub.PutState("Reg", []byte(reg))
	err = stub.PutState("Owner", []byte(owner))
	err = stub.PutState("Status", []byte(strconv.Itoa(status)))
	err = stub.PutState("Scrapped", []byte(strconv.FormatBool(scrapped)))
	err = stub.PutState("Colour", []byte(colour))
	err = stub.PutState("VIN", []byte(strconv.Itoa(VIN)))
	err = stub.PutState("V5cID", []byte(v5cID))
	err = stub.PutState("Lease_Contract_ID", []byte(leaseContractID))
	
	if err != nil {
                return nil, nil
        }

	return nil, nil
}

func (t *Vehicle) authority_to_manufacturer(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	
	statusBytes, err := stub.GetState("Status")
	status, err := strconv.Atoi(string(statusBytes))
	
	if err != nil {
		return nil, errors.New("Failed to get state")
	}
	
	ownerBytes, err := stub.GetState("Owner")
	owner := string(ownerBytes)
	
	if err != nil {
		return nil, errors.New("Failed to get state")
	}
	
	scrappedBytes, err := stub.GetState("Scrapped")
	scrapped, err := strconv.ParseBool(string(scrappedBytes))
	
	if err != nil {
		return nil, errors.New("Failed to get state")
	}
	
	if status == 0 && owner == "DVLA" && args[0] == "Toyota" && !scrapped {
		owner = args[0]
		status = 1
	} else {
		return nil, errors.New("Permission denied")
	}
	
	err = stub.PutState("Status", []byte(strconv.Itoa(status)))
	err = stub.PutState("Owner", []byte(owner))
	
	return nil, nil
	
}

func (t *Vehicle) manufacturer_to_private(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	
	statusBytes, err := stub.GetState("Status")
	status, err := strconv.Atoi(string(statusBytes))
	
	if err != nil {
		return nil, errors.New("Failed to get state")
	}
	
	makeBytes, err := stub.GetState("Make")
	make := string(makeBytes)
	
	if err != nil {
		return nil, errors.New("Failed to get state")
	}
	
	modelBytes, err := stub.GetState("Model")
	model := string(modelBytes)
	
	if err != nil {
		return nil, errors.New("Failed to get state")
	}
	
	regBytes, err := stub.GetState("Reg")
	reg := string(regBytes)
	
	if err != nil {
		return nil, errors.New("Failed to get state")
	}
	
	colourBytes, err := stub.GetState("Colour")
	colour := string(colourBytes)
	
	if err != nil {
		return nil, errors.New("Failed to get state")
	}
	
	vinBytes, err := stub.GetState("VIN")
	vin, err := strconv.Atoi(string(vinBytes))
	
	if err != nil {
		return nil, errors.New("Failed to get state")
	}
	
	ownerBytes, err := stub.GetState("Owner")
	owner := string(ownerBytes)
	
	if err != nil {
		return nil, errors.New("Failed to get state")
	}
	
	scrappedBytes, err := stub.GetState("Scrapped")
	scrapped, err := strconv.ParseBool(string(scrappedBytes))
	
	if err != nil {
		return nil, errors.New("Failed to get state")
	}
	
	if make == "UNDEFINED" || model == "UNDEFINED" || reg == "UNDEFINED" || colour == "UNDEFINED" || vin == 0 {
		return nil, errors.New("Car not fully defined")
	}
	
	if status == 1 && owner == "Toyota" && !scrapped {
		owner = args[0]
		status = 2
	} else {
		return nil, errors.New("Permission denied")
	}
	
	err = stub.PutState("Status", []byte(strconv.Itoa(status)))
	err = stub.PutState("Owner", []byte(owner))
	
	return nil, nil
	
}

func (t *Vehicle) private_to_private(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	
	statusBytes, err := stub.GetState("Status")
	status, err := strconv.Atoi(string(statusBytes))
	
	if err != nil {
		return nil, errors.New("Failed to get state")
	}
	
	ownerBytes, err := stub.GetState("Owner")
	owner := string(ownerBytes)
	
	if err != nil {
		return nil, errors.New("Failed to get state")
	}
	
	scrappedBytes, err := stub.GetState("Scrapped")
	scrapped, err := strconv.ParseBool(string(scrappedBytes))
	
	if err != nil {
		return nil, errors.New("Failed to get state")
	}
	
	if status == 2 && !scrapped {
		owner = args[0]
	} else {
		return nil, errors.New("Permission denied")
	}
	
	err = stub.PutState("Owner", []byte(owner))
	
	return nil, nil
	
}

func (t *Vehicle) private_to_lease_company(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	
	statusBytes, err := stub.GetState("Status")
	status, err := strconv.Atoi(string(statusBytes))
	
	if err != nil {
		return nil, errors.New("Failed to get state")
	}
	
	ownerBytes, err := stub.GetState("Owner")
	owner := string(ownerBytes)
	
	if err != nil {
		return nil, errors.New("Failed to get state")
	}
	
	scrappedBytes, err := stub.GetState("Scrapped")
	scrapped, err := strconv.ParseBool(string(scrappedBytes))
	
	if err != nil {
		return nil, errors.New("Failed to get state")
	}
	
	if status == 2 && args[0] == "LeaseCan" && !scrapped {
		owner = args[0]
	} else {
		return nil, errors.New("Permission denied")
	}
	
	err = stub.PutState("Owner", []byte(owner))
	
	return nil, nil
	
}

func (t *Vehicle) lease_company_to_private(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	
	statusBytes, err := stub.GetState("Status")
	status, err := strconv.Atoi(string(statusBytes))
	
	if err != nil {
		return nil, errors.New("Failed to get state")
	}
	
	ownerBytes, err := stub.GetState("Owner")
	owner := string(ownerBytes)
	
	if err != nil {
		return nil, errors.New("Failed to get state")
	}
	
	scrappedBytes, err := stub.GetState("Scrapped")
	scrapped, err := strconv.ParseBool(string(scrappedBytes))
	
	if err != nil {
		return nil, errors.New("Failed to get state")
	}
	
	if status == 2 && owner == "LeaseCan" && !scrapped {
		owner = args[0]
	} else {
		return nil, errors.New("Permission denied")
	}
	
	err = stub.PutState("Owner", []byte(owner))
	
	return nil, nil
	
}

func (t *Vehicle) private_to_scrap_merchant(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	
	statusBytes, err := stub.GetState("Status")
	status, err := strconv.Atoi(string(statusBytes))
	
	if err != nil {
		return nil, errors.New("Failed to get state")
	}
	
	ownerBytes, err := stub.GetState("Owner")
	owner := string(ownerBytes)
	
	if err != nil {
		return nil, errors.New("Failed to get state")
	}
	
	scrappedBytes, err := stub.GetState("Scrapped")
	scrapped, err := strconv.ParseBool(string(scrappedBytes))
	
	if err != nil {
		return nil, errors.New("Failed to get state")
	}
	
	if status == 2 && args[0] == "Cray Bros (London) Ltd" && !scrapped {
		owner = args[0]
		status = 4
	} else {
		return nil, errors.New("Permission denied")
	}
	
	err = stub.PutState("Status", []byte(strconv.Itoa(status)))
	err = stub.PutState("Owner", []byte(owner))
	
	return nil, nil
	
}

func (t *Vehicle) update_vin(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	
	vinBytes, err := stub.GetState("VIN")
	vin, err := strconv.Atoi(string(vinBytes))
	
	if err != nil {
		return nil, errors.New("Failed to get state")
	}
	
	ownerBytes, err := stub.GetState("Owner")
	owner := string(ownerBytes)
	
	if err != nil {
		return nil, errors.New("Failed to get state")
	}
	
	statusBytes, err := stub.GetState("Status")
	status, err := strconv.Atoi(string(statusBytes))
	
	if err != nil {
		return nil, errors.New("Failed to get state")
	}
	
	newVin, err := strconv.Atoi(args[0])
	
	if err != nil || len(args[0]) != 15 {
		return nil, errors.New("Invalid value passed for new VIN")
	}
	
	if status == 1 && owner == "Toyota" && vin == 0 {
		vin = newVin
	} else {
		return nil, errors.New("Permission denied")
	}
	
	
	err = stub.PutState("VIN", []byte(strconv.Itoa(vin)))
	
	return nil, nil
	
}

func (t *Vehicle) update_registration(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	
	regBytes, err := stub.GetState("Reg")
	
	reg := string(regBytes)
	
	if err != nil {
		return nil, errors.New("Failed to get state")
	}
	
	statusBytes, err := stub.GetState("Status")
	status, err := strconv.Atoi(string(statusBytes))
	
	if err != nil {
		return nil, errors.New("Failed to get state")
	}
	
	if status == 1 || status == 2 {
		reg = args[0]
	} else {
		return nil, errors.New("Permission denied")
	}
	
	err = stub.PutState("Reg", []byte(reg))
	
	return nil, nil
	
}

func (t *Vehicle) update_colour(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	
	colourBytes, err := stub.GetState("Colour")
	
	colour := string(colourBytes)
	
	if err != nil {
		return nil, errors.New("Failed to get state")
	}
	
	statusBytes, err := stub.GetState("Status")
	status, err := strconv.Atoi(string(statusBytes))
	
	if err != nil {
		return nil, errors.New("Failed to get state")
	}
	
	if status == 1 || status == 2 {
		colour = args[0]
	} else {
		return nil, errors.New("Permission denied")
	}
	
	err = stub.PutState("Colour", []byte(colour))
	
	return nil, nil
	
}

func (t *Vehicle) update_make(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	
	makeBytes, err := stub.GetState("Make")
	
	make := string(makeBytes)
	
	if err != nil {
		return nil, errors.New("Failed to get state")
	}
	
	statusBytes, err := stub.GetState("Status")
	status, err := strconv.Atoi(string(statusBytes))
	
	if err != nil {
		return nil, errors.New("Failed to get state")
	}
	
	ownerBytes, err := stub.GetState("Owner")
	owner := string(ownerBytes)
	
	if err != nil {
		return nil, errors.New("Failed to get state")
	}
	
	if status == 1 && owner == "Toyota" {
		make = args[0]
	} else {
		return nil, errors.New("Permission denied")
	}
	
	err = stub.PutState("Make", []byte(make))
	
	return nil, nil
}

func (t *Vehicle) update_model(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	
	modelBytes, err := stub.GetState("Model")
	
	model := string(modelBytes)
	
	if err != nil {
		return nil, errors.New("Failed to get state")
	}
	
	statusBytes, err := stub.GetState("Status")
	status, err := strconv.Atoi(string(statusBytes))
	
	if err != nil {
		return nil, errors.New("Failed to get state")
	}
	
	ownerBytes, err := stub.GetState("Owner")
	owner := string(ownerBytes)
	
	if err != nil {
		return nil, errors.New("Failed to get state")
	}
	
	if status == 1 && owner == "Toyota" {
		model = args[0]
	} else {
		return nil, errors.New("Permission denied")
	}
	
	err = stub.PutState("Model", []byte(model))
	
	return nil, nil
}

func (t *Vehicle) scrap_merchant_to_scrap(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
		
	statusBytes, err := stub.GetState("Status")
	status, err := strconv.Atoi(string(statusBytes))
	
	if err != nil {
		return nil, errors.New("Failed to get state")
	}
	
	scrappedBytes, err := stub.GetState("Scrapped")
	scrapped, err := strconv.ParseBool(string(scrappedBytes))
	
	if err != nil {
		return nil, errors.New("Failed to get state")
	}
	
	ownerBytes, err := stub.GetState("Owner")
	owner := string(ownerBytes)
	
	if err != nil {
		return nil, errors.New("Failed to get state")
	}
	
	if(!scrapped && owner == "Cray Bros (London) Ltd" && status == 4) {
		scrapped = true
	} else {
		return nil, errors.New("Permission denied")
	}

	err = stub.PutState("Scrapped", []byte(strconv.FormatBool(scrapped)))
	
	return nil, nil
}


func (t *Vehicle) Run(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {

	// Handle different functions
	if function == "init" {
		return t.init(stub, args)
	} else if function == "authority_to_manufacturer" {
		return t.authority_to_manufacturer(stub, args)
	} else if function == "manufacturer_to_private" {
		return t.manufacturer_to_private(stub, args)
	} else if function == "private_to_private" {
		return t.private_to_private(stub, args)
	} else if function == "private_to_lease_company" {
		return t.private_to_lease_company(stub, args)
	} else if function == "lease_company_to_private" {
		return t.lease_company_to_private(stub, args)
	} else if function == "private_to_scrap_merchant" {
		return t.private_to_scrap_merchant(stub, args)
	} else if function == "update_registration" {
		return t.update_registration(stub, args)
	} else if function == "update_make" {
		return t.update_make(stub, args)
	} else if function == "update_model" {
		return t.update_model(stub, args)
	} else if function == "update_vin" {
		return t.update_vin(stub, args)
	} else if function == "update_colour" {
		return t.update_colour(stub, args)
	} else if function == "scrap_merchant_to_scrap" {
		return t.scrap_merchant_to_scrap(stub, args)
	}
	
	return nil, nil
}


func (t *Vehicle) Query(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	
	if function == "get_attribute" {
		return t.get_attribute(stub, args)
	} else {
		return t.get_all(stub, args)
	}
}

func (t *Vehicle) get_attribute(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	state, err := stub.GetState(args[0])
	
	if err != nil {
		return nil, errors.New("Failed to get state")
	}
	
	jsonResp := "{\""+args[0]+"\":\"" + string(state) + "\"}"
	fmt.Printf("Query Response:%s\n", jsonResp)
	return []byte(jsonResp), nil
}

func (t *Vehicle) get_all(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	
	make, err := stub.GetState("Make")
	
	if err != nil {
		return nil, errors.New("Failed to get state")
	}
	
	model, err := stub.GetState("Model")
	
	if err != nil {
		return nil, errors.New("Failed to get state")
	}
	
	reg, err := stub.GetState("Reg")
	
	if err != nil {
		return nil, errors.New("Failed to get state")
	}
	
	owner, err := stub.GetState("Owner")
	
	if err != nil {
		return nil, errors.New("Failed to get state")
	}
	
	scrapped, err := stub.GetState("Scrapped")
	
	if err != nil {
		return nil, errors.New("Failed to get state")
	}
	
	status, err := stub.GetState("Status")
	
	if err != nil {
		return nil, errors.New("Failed to get state")
	}
	
	vin, err := stub.GetState("VIN")
	
	if err != nil {
		return nil, errors.New("Failed to get state")
	}

	v5cID, err := stub.GetState("V5cID")

        if err != nil {
                return nil, errors.New("Failed to get state")
        }

	
	colour, err := stub.GetState("Colour")
	
	if err != nil {
		return nil, errors.New("Failed to get state")
	}
	
	var jsonResp string
	
	if args[0] == string(owner) || args[0] == "DVLA" {

		jsonResp = "{\"VIN\":\"" + string(vin) + "\",\"make\":\"" + string(make) + "\",\"model\":\"" + string(model) + "\",\"reg\":\"" + string(reg) + "\", \"owner\":\"" + string(owner) + "\",\"colour\":\"" + string(colour) + "\", \"scrapped\":\"" + string(scrapped) + "\",\"status\":\"" + string(status) + "\",\"v5cID\":\"" + string(v5cID) + "\"}"
	} else {
		jsonResp = "{\"Error\":\"NO\"}"	
	}

	fmt.Printf("Query Response:%s\n", jsonResp)
	return []byte(jsonResp), nil
}

func main() {
	err := shim.Start(new(Vehicle))
	if err != nil {
		fmt.Printf("Error starting Vehicle: %s", err)
	}
}
