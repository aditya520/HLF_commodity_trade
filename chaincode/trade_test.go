package main

import (
	"fmt"
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

var testLog = shim.NewLogger("trade_cc_test")

func checkInit(t *testing.T, stub *shim.MockStub, args [][]byte) {
	res := stub.MockInit("1", args)
	if res.Status != shim.OK {
		fmt.Println("Init failed", string(res.Message))
		t.FailNow()
	}
}

func checkState(t *testing.T, stub *shim.MockStub, name string, value string) {
	bytes := stub.State[name]
	if bytes == nil {
		testLog.Info("State", name, "failed to get value")
		t.FailNow()
	}
	if string(bytes) != value {
		testLog.Info("State value", name, "was", string(bytes), "and not", value, "as expected")
		t.FailNow()
	} else {
		testLog.Info("State value", name, "is", string(bytes), "as expected")
	}
}

func checkNoState(t *testing.T, stub *shim.MockStub, name string) {
	bytes := stub.State[name]
	if bytes != nil {
		testLog.Info("State", name, "should be absent; found value")
		t.FailNow()
	} else {
		testLog.Info("State", name, "is absent as it should be")
	}
}

func checkQueryOneArg(t *testing.T, stub *shim.MockStub, function string, argument string, value string) {
	res := stub.MockInvoke("1", [][]byte{[]byte(function), []byte(argument)})
	if res.Status != shim.OK {
		testLog.Info("Query", function, "failed", string(res.Message))
		t.FailNow()
	}
	if res.Payload == nil {
		testLog.Info("Query", function, "failed to get value")
		t.FailNow()
	}
	payload := string(res.Payload)
	if payload != value {
		testLog.Info("Query value", function, "was", payload, "and not", value, "as expected")
		t.FailNow()
	} else {
		testLog.Info("Query value", function, "is", payload, "as expected")
	}
}

func checkBadQuery(t *testing.T, stub *shim.MockStub, function string, name string) {
	res := stub.MockInvoke("1", [][]byte{[]byte(function), []byte(name)})
	if res.Status == shim.OK {
		testLog.Info("Query", function, "unexpectedly succeeded")
		t.FailNow()
	} else {
		testLog.Info("Query", function, "failed as espected, with message: ", res.Message)

	}
}

func checkInvoke(t *testing.T, stub *shim.MockStub, functionAndArgs []string) {
	functionAndArgsAsBytes := lib.ParseStringSliceToByteSlice(functionAndArgs)
	res := stub.MockInvoke("1", functionAndArgsAsBytes)
	if res.Status != shim.OK {
		testLog.Info("Invoke", functionAndArgs, "failed", string(res.Message))
		t.FailNow()
	} else {
		testLog.Info("Invoke", functionAndArgs, "successful", string(res.Message))
	}
}

func checkBadInvoke(t *testing.T, stub *shim.MockStub, functionAndArgs []string) {
	functionAndArgsAsBytes := lib.ParseStringSliceToByteSlice(functionAndArgs)
	res := stub.MockInvoke("1", functionAndArgsAsBytes)
	if res.Status == shim.OK {
		testLog.Info("Invoke", functionAndArgs, "unexpectedly succeeded")
		t.FailNow()
	} else {
		testLog.Info("Invoke", functionAndArgs, "failed as espected, with message: "+res.Message)
	}
}

func TestExample02_Init(t *testing.T) {
	scc := new(Trade)
	stub := shim.NewMockStub("trade_test", scc)
	checkInit(t, stub, [][]byte{})
}

func TestExample02_Invoke(t *testing.T) {
	scc := new(Trade)
	stub := shim.NewMockStub("trade_test", scc)

	// checkInit(t, stub, [][]byte{})

	// checkInvoke(t, stub, [][]byte{[]byte("generateSalesContract"), []byte("3"), []byte("Buyer1"), []byte("Seller1"), []byte("Commodity"), []byte("10"), []byte("2000"), []byte("02/12/2015")})
	// checkQuery(t, stub, "1")
	// checkQuery(t, stub, "B", "801")

	// // Invoke B->A for 234
	// checkInvoke(t, stub, [][]byte{[]byte("invoke"), []byte("B"), []byte("A"), []byte("234")})
	// checkQuery(t, stub, "A", "678")
	// checkQuery(t, stub, "B", "567")
	// checkQuery(t, stub, "A", "678")
	// checkQuery(t, stub, "B", "567")
}

// func TestExample02_Query(t *testing.T) {
// 	scc := new(SimpleChaincode)
// 	stub := shim.NewMockStub("ex02", scc)

// 	// Init A=345 B=456
// 	checkInit(t, stub, [][]byte{[]byte("init"), []byte("A"), []byte("345"), []byte("B"), []byte("456")})

// 	// Query A
// 	checkQuery(t, stub, "A", "345")

// 	// Query B
// 	checkQuery(t, stub, "B", "456")
// }
