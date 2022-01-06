/*
This is for Training and evaluation purpose in Wipro, 2020 All Rights Reserved.

A Chain code implementation for Private Data Management in Hyperledger Fabric

*/

package main

import (
	"bytes"
	"encoding/json"
	"fmt"

	//"strings"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

var logger = shim.NewLogger("##################### eGadgets_Chaincode #####################")

// GadgetChaincode Chaincode implementation
//===================================================================================
type GadgetChaincode struct {
}

//defining structure for gadget with 5 properties
//===================================================================================
type gadget struct {
	GadgetId string `json:"GadgetId"`
	Make     string `json:"Make"`
	Model    string `json:"Model"`
	Color    string `json:"Color"`
	Owner    string `json:"Owner"`
}

type gadgetPrivateDetails struct {
	GadgetId string `json:"GadgetId"`
	Price    string `json:"Price"`
}

//===================================================================================
// Main method starts here
//===================================================================================
func main() {
	err := shim.Start(new(GadgetChaincode))
	if err != nil {
		fmt.Printf("Error starting eGadgets chaincode: %s", err)
	}
}

//===================================================================================
// Init initializes the chaincode
//===================================================================================
func (t *GadgetChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	logger.Info("##################### eGadgets_Chaincode_init #####################")
	return shim.Success(nil)
}

//===================================================================================
// Invoke
//===================================================================================
func (t *GadgetChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := stub.GetFunctionAndParameters()

	//Printing message
	fmt.Println("\n##################### invoke is running ##################### \n#####################invoked function name :  " + function + "#####################")

	//Handling functions

	if function == "initGadget" { //Initialising the Gadget
		return t.initGadget(stub, args)
	} else if function == "queryGadgetPrivateDetails" { //Fetch Private Details of Gadget from collectioneGadgetsPrivateDetails
		return t.queryGadgetPrivateDetails(stub, args)
	} else if function == "queryGadgetsDetails" { //fetch Gadget Details from collectioneGadgets
		return t.queryGadgetsDetails(stub, args)
	} else if function == "queryAllGadgetDetailsByRange" {
		return t.queryAllGadgetDetailsByRange(stub, args)
	}
	fmt.Println("invoke did not find function: " + function) //error
	return shim.Error("Received unknown function invocation")
}

//===================================================================================
// init
//===================================================================================
func (t *GadgetChaincode) initGadget(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var err error

	//args[] 		0		1		2		3		4					5
	//value 	G00001		HP		2020	Blue	Mr R K Sharma		34000
	if len(args) != 6 {
		return shim.Error("Expecting 6 arguments....")
	}
	fmt.Println("##################### initGadget is running #####################")

	//Checking the arguments zero or null values

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

	//Initialising the variables
	gadgetId := args[0]
	gadgetMake := args[1]
	gadgetModel := args[2]
	gadgetColor := args[3]
	gadgetOwner := args[4]
	gadgetPrice := args[5]

	//Check if Gadget already exists
	gadgetAsBytes, err := stub.GetPrivateData("collectioneGadgets", gadgetId)
	if err != nil {
		return shim.Error("Failed to get Gadget : " + err.Error())
	}
	if gadgetAsBytes != nil {
		fmt.Println("This Gadget already exists : " + gadgetId)
		return shim.Error("This Gadget already exists : " + gadgetId)
	}

	//Creating Object for Gadget and Marshal to JSON
	GadgetId := gadgetId
	gadget := &gadget{gadgetId, gadgetMake, gadgetModel, gadgetColor, gadgetOwner}
	gadgetJSONasBytes, err := json.Marshal(gadget)
	if err != nil {
		return shim.Error(err.Error())
	}

	// === Save Gadget details into state ===
	err = stub.PutPrivateData("collectioneGadgets", GadgetId, gadgetJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	// ==== Create private details for gadget with price, marshal to JSON, and save to state ====
	gadgetPrivateDetails := &gadgetPrivateDetails{gadgetId, gadgetPrice}

	gadgetPrivateDetailsJSONasBytes, err := json.Marshal(gadgetPrivateDetails)
	if err != nil {
		return shim.Error(err.Error())
	}
	// === Save Gadget private details into state ===
	err = stub.PutPrivateData("collectioneGadgetsPrivateDetails", GadgetId, gadgetPrivateDetailsJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("Gadget Details Saved as : ")

	fmt.Println(string(gadgetJSONasBytes), "Gadget Private Details : "+string(gadgetPrivateDetailsJSONasBytes))
	fmt.Println(" ******************** end of initGadget ******************** ")
	return shim.Success(nil)
}

//===================================================================================
// queryGadgetPrivateDetails
//===================================================================================
func (t *GadgetChaincode) queryGadgetPrivateDetails(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("##################### queryGadgetPrivateDetails started #####################")

	var gadgetId, jsonResp string
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting GadgetID of the Gadget to query")
	}

	gadgetId = args[0]
	gadgetAsBytes, err := stub.GetPrivateData("collectioneGadgetsPrivateDetails", gadgetId) //get the Gadget private details from chaincode state
	if err != nil {
		jsonResp = "Error : Not Authorised for private data " + err.Error()
		return shim.Error(jsonResp)
	} else if gadgetAsBytes == nil {
		jsonResp = "Error: Private details does not exist for GadgetId: " + gadgetId
		return shim.Error(jsonResp)
	}

	fmt.Println("Gadget Private Details are : ")
	fmt.Println(string(gadgetAsBytes))
	fmt.Println("  ******************** queryGadgetPrivateDetails ended ******************** ")

	return shim.Success(gadgetAsBytes)

}

//===================================================================================
// queryGadgetsDetails
//===================================================================================
func (t *GadgetChaincode) queryGadgetsDetails(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var gadgetId, jsonResp string
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting GadgetID for the Gadget to query")
	}

	gadgetId = args[0]
	gadgetAsBytes, err := stub.GetPrivateData("collectioneGadgets", gadgetId) //get the Gadget details from chaincode state
	if err != nil {
		jsonResp = "Error : Failed to get Gadget details for " + gadgetId + " : " + err.Error()
		return shim.Error(jsonResp)
	} else if gadgetAsBytes == nil {
		jsonResp = "Error : Gadget details does not exist for GadgetId: " + gadgetId + " : " + err.Error()
		return shim.Error(jsonResp)
	}

	fmt.Println("Gadget Details are : ")
	fmt.Println(string(gadgetAsBytes))
	fmt.Println("******************** queryGadgetDetails ended ******************** ")
	return shim.Success(gadgetAsBytes)

}

//===================================================================================
// queryAllGadgetDetailsByRange
//===================================================================================
func (t *GadgetChaincode) queryAllGadgetDetailsByRange(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) < 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	startKey := args[0]
	endKey := args[1]

	resultsIterator, err := stub.GetPrivateDataByRange("collectioneGadgets", startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}\n")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("******************** AllGadgetDetailsByRange Result:\n%s\n", buffer.String())
	fmt.Printf("\n******************** queryAllGadgetDetailsByRange ended ******************** ")
	return shim.Success(buffer.Bytes())
}
