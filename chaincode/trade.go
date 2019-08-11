package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type Trade struct {
}

type contract struct {
	// ObjectType   string `json:"docType"` //docType is used to distinguish the various types of objects in state database
	contractId   string  `json:"contractId"`
	buyerName    string  `json:"buyerName"`
	sellerName   string  `json:"sellerName"`
	commodity    string  `json:"commodity"`
	weight       float64 `json:"weight"`
	price        float64 `json:"price"`
	grade        string  `json:"grade"`
	deliveryDate string  `json:"deliveryDate"`
	status       string  `json:"status"`
}

type contractSecret struct {
	contractId     string `json:"contractId"`
	additionalInfo string `json:"additionalInfo"`
}

func (t *Trade) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("Initializing Trade Workflow")
	return shim.Success(nil)
}

func (t *Trade) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	// fmt.Println("ex02 Invoke")
	function, args := stub.GetFunctionAndParameters()
	if function == "generateSalesContract" {
		return t.generateSalesContract(stub, args)
	} else if function == "readContract" {
		return t.readContract(stub, args)
	}
	// } else if function == "query" {
	// 	// the old "Query" is now implemtned in invoke
	// 	return t.query(stub, args)
	// }

	return shim.Error("Invalid invoke function name. Expecting \"invoke\" \"delete\" \"query\"")
}

// Transaction makes payment of X units from A to B
func (t *Trade) generateSalesContract(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	// if len(args) != 3 {
	// 	return shim.Error("Incorrect number of arguments. Expecting 3")
	// }

	// All Validations remaining
	// Make event emitter
	// ACLs remaining
	// Composite key

	contractId := args[0]

	contractBytes, err := stub.GetState(contractId)
	if err != nil {
		return shim.Error("Failed to get contract: " + err.Error())
	} else if contractBytes != nil {
		fmt.Println("This contract already exists: " + contractId)
		return shim.Error("This contract already exists: " + contractId)
	}

	buyerName := args[1]
	sellerName := args[2]
	commodity := args[3]
	weight, err1 := strconv.ParseFloat(args[4], 32)
	price, err2 := strconv.ParseFloat(args[5], 32)
	if err1 != nil || err2 != nil {
		return shim.Error("Error parsing the values")
	}
	grade := args[6]
	deliveryDate := args[7]
	additionalInfo := args[8]
	status := "SellerCreated"

	contract := &contract{contractId, buyerName, sellerName, commodity, weight, price, grade, deliveryDate, status}
	contractBytes, err = json.Marshal(contract)

	if err != nil {
		return shim.Error(err.Error())
	}

	contractSecret := &contractSecret{contractId, additionalInfo}
	contractSecretBytes, err := json.Marshal(contractSecret)

	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(contractId, contractBytes)

	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutPrivateData("secret", contractId, contractSecretBytes)

	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (t *Trade) readContract(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting contractId of the contract to query")
	}

	contractId := args[0]
	valAsbytes, err := stub.GetState(contractId)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + contractId + "\"}"
		return shim.Error(jsonResp)
	} else if valAsbytes == nil {
		jsonResp := "{\"Error\":\"contract does not exist: " + contractId + "\"}"
		return shim.Error(jsonResp)
	}

	return shim.Success(valAsbytes)
}

// func (t *Trade) updateContract(stub shim.ChaincodeStubInterface, args []string) pb.Response {

// 	if len(args) != 1 {
// 		return shim.Error("Incorrect number of arguments. Expecting contractId of the contract to query")
// 	}

// 	contractId := args[0]
// 	valAsbytes, err := stub.GetState(contractId)
// 	if err != nil {
// 		jsonResp := "{\"Error\":\"Failed to get state for " + contractId + "\"}"
// 		return shim.Error(jsonResp)
// 	} else if valAsbytes == nil {
// 		jsonResp := "{\"Error\":\"contract does not exist: " + contractId + "\"}"
// 		return shim.Error(jsonResp)
// 	}

// 	return shim.Success(valAsbytes)
// }

func main() {
	err := shim.Start(new(Trade))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
