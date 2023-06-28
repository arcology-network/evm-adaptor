package tests

import (
	"math/big"
	"os"
	"path/filepath"
	"testing"

	evmcommon "github.com/arcology-network/evm/common"
	"github.com/arcology-network/evm/core"
	ccEu "github.com/arcology-network/vm-adaptor"
	eucommon "github.com/arcology-network/vm-adaptor/common"
	"github.com/arcology-network/vm-adaptor/compiler"
)

func TestContractNoncommutativeInt256(t *testing.T) {
	eu, config, _, _, _ := NewTestEU()
	// ================================== Compile the contract ==================================
	currentPath, _ := os.Getwd()
	project := filepath.Dir(currentPath)
	targetPath := project + "/api/noncommutative"

	code, err := compiler.CompileContracts(targetPath, "int256/int256_test.sol", "0.8.19", "Int256Test", false)
	if err != nil || len(code) == 0 {
		t.Error("Error: Failed to generate the byte code")
	}
	// ================================== Deploy the contract ==================================
	msg := core.NewMessage(eucommon.Alice, nil, 0, new(big.Int).SetUint64(0), 1e15, new(big.Int).SetUint64(1), evmcommon.Hex2Bytes(code), nil, true) // Build the message
	receipt, _, err := eu.Run(evmcommon.BytesToHash([]byte{1, 1, 1}), 1, &msg, ccEu.NewEVMBlockContext(config), ccEu.NewEVMTxContext(msg))           // Execute it
	_, transitions := eu.Api().Ccurl().ExportAll()
	// ---------------

	// t.Log("\n" + FormatTransitions(accesses))
	t.Log("\n" + eucommon.FormatTransitions(transitions))
	// t.Log(receipt)
	// contractAddress := receipt.ContractAddress
	if receipt.Status != 1 || err != nil {
		t.Error("Error: Deployment failed!!!", err)
	}
}

func TestNoncommutativeInt256N(t *testing.T) {
	eu, config, _, _, _ := NewTestEU()

	// ================================== Compile the contract ==================================
	currentPath, _ := os.Getwd()
	project := filepath.Dir(currentPath)
	targetPath := project + "/api/noncommutative"

	code, err := compiler.CompileContracts(targetPath, "int256/int256N_test.sol", "0.8.19", "Int64NTest", false)
	if err != nil || len(code) == 0 {
		t.Error(err)
	}
	// ================================== Deploy the contract ==================================
	msg := core.NewMessage(eucommon.Alice, nil, 0, new(big.Int).SetUint64(0), 1e15, new(big.Int).SetUint64(1), evmcommon.Hex2Bytes(code), nil, true) // Build the message
	receipt, _, err := eu.Run(evmcommon.BytesToHash([]byte{1, 1, 1}), 1, &msg, ccEu.NewEVMBlockContext(config), ccEu.NewEVMTxContext(msg))           // Execute it
	_, transitions := eu.Api().Ccurl().ExportAll()
	// ---------------

	// t.Log("\n" + FormatTransitions(accesses))
	t.Log("\n" + eucommon.FormatTransitions(transitions))
	t.Log(receipt)
	// contractAddress := receipt.ContractAddress
	if receipt.Status != 1 || err != nil {
		t.Error("Error: Deployment failed!!!", err)
	}
}

/*
func TestCumulativeInt256(t *testing.T) {
	eu, config, _, _, _ := NewTestEU()

	// ================================== Compile the contract ==================================
	currentPath, _ := os.Getwd()
	project := filepath.Dir(currentPath)

	targetPath := project + "/api/commutative/"

	code, err := compiler.CompileContracts(targetPath+"int256", "int256Cumulative_test.sol", "0.8.19", "Int256CumulativeTest", false)
	// code, err := compiler.CompileContracts(pyCompiler, targetPath+"/int256/int256Cumulative_test.sol", "Int256CumulativeTest")

	if err != nil || len(code) == 0 {
		t.Error(err)
	}

	// ================================== Deploy the contract ==================================
	msg := core.NewMessage(eucommon.Alice, nil, 0, new(big.Int).SetUint64(0), 1e15, new(big.Int).SetUint64(1), evmcommon.Hex2Bytes(code), nil, true) // Build the message
	receipt, _, err := eu.Run(evmcommon.BytesToHash([]byte{1, 1, 1}), 1, &msg, ccEu.NewEVMBlockContext(config), ccEu.NewEVMTxContext(msg))           // Execute it
	_, transitions := eu.Api().Ccurl().ExportAll()
	// ---------------

	// t.Log("\n" + FormatTransitions(accesses))
	t.Log("\n" + eucommon.FormatTransitions(transitions))
	t.Log(receipt)
	// contractAddress := receipt.ContractAddress
	if receipt.Status != 1 || err != nil {
		t.Error("Error: Deployment failed!!!", err)
	}
}
*/

// func TestInt64Threading(t *testing.T) {
// 	eu, config, _, _, _ := NewTestEU()

// 	// ================================== Compile the contract ==================================
// 	currentPath, _ := os.Getwd()
// 	project := filepath.Dir(currentPath)
// 	pyCompiler := project + "/compiler/compiler.py"

// 	code, err := compiler.CompileContracts(pyCompiler, project+"/api/noncommutative/int256/int256_threading.sol", "ThreadingInt64")

// 	if err != nil || len(code) == 0 {
// 		t.Error("Error: ", err)
// 	}
// 	// ================================== Deploy the contract ==================================
// 	msg := types.NewMessage(eucommon.Alice, nil, 0, new(big.Int).SetUint64(0), 1e15, new(big.Int).SetUint64(1), evmcommon.Hex2Bytes(code), nil, false) // Build the message
// 	receipt, _, err := eu.Run(evmcommon.BytesToHash([]byte{1, 1, 1}), 1, &msg, evmeu.NewEVMBlockContext(config), evmeu.NewEVMTxContext(msg))       // Execute it
// 	// t.Log("\n" + eucommon.FormatTransitions(transitions))
// 	// indexer.Univalues(transitions).Print()
// 	if receipt.Status != 1 || err != nil {
// 		t.Error("Error: Deployment failed!!!", err)
// 	}
// 	fmt.Println(receipt.ContractAddress)

// 	// // ================================== CallBasic() ==================================
// 	data := crypto.Keccak256([]byte("call()"))[:4]
// 	contractAddress := receipt.ContractAddress
// 	msg = types.NewMessage(eucommon.Alice, &contractAddress, 1, new(big.Int).SetUint64(0), 1e15, new(big.Int).SetUint64(1), data, nil, false)
// 	receipt, execResult, err := eu.Run(evmcommon.BytesToHash([]byte{1, 1, 1}), 1, &msg, evmeu.NewEVMBlockContext(config), evmeu.NewEVMTxContext(msg))
// 	// indexer.Univalues(transitions).Print()

// 	if err != nil {
// 		t.Error(err)
// 	}

// 	if execResult != nil && execResult.Err != nil {
// 		t.Error(execResult.Err)
// 	}
// }
