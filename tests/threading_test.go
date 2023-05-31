package tests

import (
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"testing"

	"github.com/arcology-network/concurrenturl"
	evmcommon "github.com/arcology-network/evm/common"
	"github.com/arcology-network/evm/core/types"
	"github.com/arcology-network/evm/crypto"
	cceu "github.com/arcology-network/vm-adaptor"
	eucommon "github.com/arcology-network/vm-adaptor/common"
	compiler "github.com/arcology-network/vm-adaptor/compiler"
)

func TestThreadingBasic(t *testing.T) {
	eu, config, db, url, _ := NewTestEU()

	// ================================== Compile the contract ==================================
	currentPath, _ := os.Getwd()
	project := filepath.Dir(currentPath)
	pyCompiler := project + "/compiler/compiler.py"

	code, err := compiler.CompileContracts(pyCompiler, project+"/api/threading/threading_test.sol", "ThreadingTest")
	if err != nil || len(code) == 0 {
		t.Error("Error: Failed to generate the byte code")
	}
	// ================================== Deploy the contract ==================================
	msg := types.NewMessage(eucommon.Alice, nil, 0, new(big.Int).SetUint64(0), 1e15, new(big.Int).SetUint64(1), evmcommon.Hex2Bytes(code), nil, false) // Build the message
	receipt, _, err := eu.Run(evmcommon.BytesToHash([]byte{1, 1, 1}), 1, &msg, cceu.NewEVMBlockContext(config), cceu.NewEVMTxContext(msg))             // Execute it
	_, transitions := eu.Api().Ccurl().ExportAll()

	contractAddress := receipt.ContractAddress
	if receipt.Status != 1 || err != nil {
		t.Error("Error: Deployment failed!!!", err)
	}
	fmt.Println(receipt.ContractAddress)

	url = concurrenturl.NewConcurrentUrl(db)
	url.Import(transitions)
	url.Sort()
	url.Commit([]uint32{1})

	// ================================== Call() ==================================
	receipt, _, err = eu.Run(evmcommon.BytesToHash([]byte{1, 1, 1}), 1, &msg, cceu.NewEVMBlockContext(config), cceu.NewEVMTxContext(msg))
	// _, transitions = eu.Api().Ccurl().ExportAll()

	if err != nil {
		fmt.Print(err)
	}

	data := crypto.Keccak256([]byte("call()"))[:4]
	msg = types.NewMessage(eucommon.Alice, &contractAddress, 1, new(big.Int).SetUint64(0), 1e15, new(big.Int).SetUint64(1), data, nil, false)
	receipt, execResult, err := eu.Run(evmcommon.BytesToHash([]byte{1, 1, 1}), 1, &msg, cceu.NewEVMBlockContext(config), cceu.NewEVMTxContext(msg))
	_, transitions = eu.Api().Ccurl().ExportAll()

	if err != nil {
		t.Error(err)
	}

	if execResult != nil && execResult.Err != nil {
		t.Error(execResult.Err)
	}

	if receipt.Status != 1 || err != nil {
		t.Error("Error: Failed to call!!!", err)
	}
}

// func TestThreadingLib(t *testing.T) {
// 	eu, config, db, url, _ := NewTestEU()

// 	// ================================== Compile the contract ==================================
// 	currentPath, _ := os.Getwd()
// 	project := filepath.Dir(currentPath)
// 	pyCompiler := project + "/compiler/compiler.py"

// 	code, err := compiler.CompileContracts(pyCompiler, project+"/api/threading/threadingLib_test.sol", "ThreadingLibTest")
// 	if err != nil || len(code) == 0 {
// 		t.Error("Error: Failed to generate the byte code")
// 	}
// 	// ================================== Deploy the contract ==================================
// 	msg := types.NewMessage(eucommon.Alice, nil, 0, new(big.Int).SetUint64(0), 1e15, new(big.Int).SetUint64(1), evmcommon.Hex2Bytes(code), nil, false) // Build the message
// 	receipt, _, err := eu.Run(evmcommon.BytesToHash([]byte{1, 1, 1}), 1, &msg, cceu.NewEVMBlockContext(config), cceu.NewEVMTxContext(msg))             // Execute it
// 	_, transitions := eu.Api().Ccurl().ExportAll()

// 	contractAddress := receipt.ContractAddress
// 	if receipt.Status != 1 || err != nil {
// 		t.Error("Error: Deployment failed!!!", err)
// 	}
// 	fmt.Println(receipt.ContractAddress)

// 	url = concurrenturl.NewConcurrentUrl(db)
// 	url.Import(transitions)
// 	url.Sort()
// 	url.Commit([]uint32{1})

// 	// ================================== Call() ==================================
// 	receipt, _, err = eu.Run(evmcommon.BytesToHash([]byte{1, 1, 1}), 1, &msg, cceu.NewEVMBlockContext(config), cceu.NewEVMTxContext(msg))
// 	_, transitions = eu.Api().Ccurl().ExportAll()

// 	if err != nil {
// 		fmt.Print(err)
// 	}

// 	data := crypto.Keccak256([]byte("call()"))[:4]
// 	msg = types.NewMessage(eucommon.Alice, &contractAddress, 1, new(big.Int).SetUint64(0), 1e15, new(big.Int).SetUint64(1), data, nil, false)
// 	receipt, execResult, err := eu.Run(evmcommon.BytesToHash([]byte{1, 1, 1}), 1, &msg, cceu.NewEVMBlockContext(config), cceu.NewEVMTxContext(msg))
// 	_, transitions = eu.Api().Ccurl().ExportAll()

// 	if err != nil {
// 		t.Error(err)
// 	}

// 	if execResult != nil && execResult.Err != nil {
// 		t.Error(execResult.Err)
// 	}

// 	if receipt.Status != 1 || err != nil {
// 		t.Error("Error: Failed to call!!!", err)
// 	}
// }

func BenchmarkThreadingReverseString(t *testing.B) {
	eu, config, db, url, _ := NewTestEU()

	// ================================== Compile the contract ==================================
	currentPath, _ := os.Getwd()
	project := filepath.Dir(currentPath)
	pyCompiler := project + "/compiler/compiler.py"

	code, err := compiler.CompileContracts(pyCompiler, project+"/api/threading/threading_test.sol", "Benchmarking")
	if err != nil || len(code) == 0 {
		t.Error("Error: Failed to generate the byte code")
	}
	// ================================== Deploy the contract ==================================
	msg := types.NewMessage(eucommon.Alice, nil, 0, new(big.Int).SetUint64(0), 1e15, new(big.Int).SetUint64(1), evmcommon.Hex2Bytes(code), nil, false) // Build the message
	receipt, _, err := eu.Run(evmcommon.BytesToHash([]byte{1, 1, 1}), 1, &msg, cceu.NewEVMBlockContext(config), cceu.NewEVMTxContext(msg))             // Execute it
	_, transitions := eu.Api().Ccurl().ExportAll()

	// t.Log("\n" + eucommon.FormatTransitions(transitions))

	contractAddress := receipt.ContractAddress
	if receipt.Status != 1 || err != nil {
		t.Error("Error: Deployment failed!!!", err)
	}
	// fmt.Println(receipt.ContractAddress)

	url = concurrenturl.NewConcurrentUrl(db)
	url.Import(transitions)
	url.Sort()
	url.Commit([]uint32{1})

	// ================================== Test SHA1() ==================================
	receipt, _, err = eu.Run(evmcommon.BytesToHash([]byte{1, 1, 1}), 1, &msg, cceu.NewEVMBlockContext(config), cceu.NewEVMTxContext(msg))
	_, transitions = eu.Api().Ccurl().ExportAll()

	if err != nil {
		fmt.Print(err)
	}

	data := crypto.Keccak256([]byte("testReverseString10k()"))[:4]
	msg = types.NewMessage(eucommon.Alice, &contractAddress, 1, new(big.Int).SetUint64(0), 1e15, new(big.Int).SetUint64(1), data, nil, false)
	receipt, execResult, err := eu.Run(evmcommon.BytesToHash([]byte{1, 1, 1}), 1, &msg, cceu.NewEVMBlockContext(config), cceu.NewEVMTxContext(msg))
	_, transitions = eu.Api().Ccurl().ExportAll()

	fmt.Println()
	if err != nil {
		t.Error(err)
	}

	if execResult != nil && execResult.Err != nil {
		t.Error(execResult.Err)
	}

	if receipt.Status != 1 || err != nil {
		t.Error("Error: Failed to call!!!", err)
	}
}
