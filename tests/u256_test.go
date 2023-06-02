package tests

import (
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"testing"

	common "github.com/arcology-network/common-lib/common"
	evmcommon "github.com/arcology-network/evm/common"
	"github.com/arcology-network/evm/core/types"
	"github.com/arcology-network/evm/crypto"
	ccEu "github.com/arcology-network/vm-adaptor"
	cceu "github.com/arcology-network/vm-adaptor"
	eucommon "github.com/arcology-network/vm-adaptor/common"
	"github.com/arcology-network/vm-adaptor/compilers"
)

func TestNoncommutativeU256Dynamic(t *testing.T) {
	eu, config, _, _, _ := NewTestEU()

	// ================================== Compile the contract ==================================
	currentPath, _ := os.Getwd()
	project := filepath.Dir(currentPath)
	// pyCompiler := project + "/compiler/compiler.py"
	targetPath := project + "/api/noncommutative/"

	if err := common.CopyFile(targetPath+"base/Base.sol", targetPath+"/u256/Base.sol"); err != nil {
		t.Error(err)
	}

	if err := common.CopyFile(project+"/api/threading/Threading.sol", targetPath+"/u256/Threading.sol"); err != nil {
		t.Error(err)
	}

	code, err := compilers.CompileContracts(targetPath+"/u256", "u256_test.sol", "0.5.0", "U256DynamicTest", false)
	// code, err := compiler.CompileContracts(pyCompiler, targetPath+"/u256/u256_test.sol", "U256DynamicTest")

	if err != nil || len(code) == 0 {
		t.Error("Error: Failed to generate the byte code")
	}

	// ================================== Deploy the contract ==================================
	msg := types.NewMessage(eucommon.Alice, nil, 0, new(big.Int).SetUint64(0), 1e15, new(big.Int).SetUint64(1), evmcommon.Hex2Bytes(code), nil, true) // Build the message
	receipt, _, err := eu.Run(evmcommon.BytesToHash([]byte{1, 1, 1}), 1, &msg, ccEu.NewEVMBlockContext(config), ccEu.NewEVMTxContext(msg))            // Execute it
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

func TestNoncommutative256N(t *testing.T) {
	eu, config, _, _, _ := NewTestEU()

	// ================================== Compile the contract ==================================
	currentPath, _ := os.Getwd()
	project := filepath.Dir(currentPath)
	// pyCompiler := project + "/compiler/compiler.py"
	targetPath := project + "/api/noncommutative/"
	baseFile := targetPath + "base/Base.sol"

	if err := common.CopyFile(baseFile, targetPath+"/u256/Base.sol"); err != nil {
		t.Error(err)
	}

	if err := common.CopyFile(project+"/api/threading/Threading.sol", targetPath+"/u256/Threading.sol"); err != nil {
		t.Error(err)
	}

	code, err := compilers.CompileContracts(targetPath+"/u256", "u256N_test.sol", "0.5.0", "U256NTest", false)
	// code, err := compiler.CompileContracts(pyCompiler, targetPath+"/u256/u256N_test.sol", "U256NTest")

	if err != nil || len(code) == 0 {
		t.Error(err)
	}
	// ================================== Deploy the contract ==================================
	msg := types.NewMessage(eucommon.Alice, nil, 0, new(big.Int).SetUint64(0), 1e15, new(big.Int).SetUint64(1), evmcommon.Hex2Bytes(code), nil, true) // Build the message
	receipt, _, err := eu.Run(evmcommon.BytesToHash([]byte{1, 1, 1}), 1, &msg, ccEu.NewEVMBlockContext(config), ccEu.NewEVMTxContext(msg))            // Execute it
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

func TestCumulativeU256Case1(t *testing.T) {
	eu, config, _, _, _ := NewTestEU()

	// ================================== Compile the contract ==================================
	currentPath, _ := os.Getwd()
	project := filepath.Dir(currentPath)
	// pyCompiler := project + "/compiler/compiler.py"
	targetPath := project + "/api/commutative/"

	if err := common.CopyFile(project+"/api/threading/Threading.sol", targetPath+"/u256/Threading.sol"); err != nil {
		t.Error(err)
	}

	code, err := compilers.CompileContracts(targetPath+"/u256", "u256Cumulative_test.sol", "0.5.0", "ThreadingCumulativeU256", false)
	// code, err := compiler.CompileContracts(pyCompiler, targetPath+"/u256/u256Cumulative_test.sol", "ThreadingCumulativeU256")

	if err != nil || len(code) == 0 {
		t.Error(err)
	}

	// ================================== Deploy the contract ==================================
	msg := types.NewMessage(eucommon.Alice, nil, 0, new(big.Int).SetUint64(0), 1e15, new(big.Int).SetUint64(1), evmcommon.Hex2Bytes(code), nil, true) // Build the message
	receipt, _, err := eu.Run(evmcommon.BytesToHash([]byte{1, 1, 1}), 1, &msg, ccEu.NewEVMBlockContext(config), ccEu.NewEVMTxContext(msg))            // Execute it

	_, transitions := eu.Api().Ccurl().ExportAll()
	eu.Api().Ccurl().Import(transitions)
	eu.Api().Ccurl().Sort()
	eu.Api().Ccurl().Commit([]uint32{1})

	// ---------------
	t.Log(receipt)
	contractAddress := receipt.ContractAddress
	if receipt.Status != 1 || err != nil {
		t.Error("Error: Deployment failed!!!", err)
	}

	// ================================== CallBasic() ==================================
	receipt, _, err = eu.Run(evmcommon.BytesToHash([]byte{1, 1, 1}), 1, &msg, cceu.NewEVMBlockContext(config), cceu.NewEVMTxContext(msg))
	_, transitions = eu.Api().Ccurl().ExportAll()
	if err != nil {
		fmt.Print(err)
	}

	data := crypto.Keccak256([]byte("call()"))[:4]
	msg = types.NewMessage(eucommon.Alice, &contractAddress, 1, new(big.Int).SetUint64(0), 1e15, new(big.Int).SetUint64(1), data, nil, false)
	receipt, execResult, err := eu.Run(evmcommon.BytesToHash([]byte{1, 1, 1}), 1, &msg, cceu.NewEVMBlockContext(config), cceu.NewEVMTxContext(msg))
	_, transitions = eu.Api().Ccurl().ExportAll()

	if receipt.Status != 1 || err != nil {
		t.Error(err)
	}

	if execResult != nil && execResult.Err != nil {
		t.Error(execResult.Err)
	}
}

func TestCumulativeU256Case2(t *testing.T) {
	eu, config, _, _, _ := NewTestEU()

	// ================================== Compile the contract ==================================
	currentPath, _ := os.Getwd()
	project := filepath.Dir(currentPath)
	// pyCompiler := project + "/compiler/compiler.py"
	targetPath := project + "/api/commutative/"

	code, err := compilers.CompileContracts(targetPath+"/u256", "u256Cumulative_test.sol", "0.5.0", "ThreadingCumulativeU256", false)
	// code, err := compiler.CompileContracts(pyCompiler, targetPath+"/u256/u256Cumulative_test.sol", "ThreadingCumulativeU256")

	if err != nil || len(code) == 0 {
		t.Error(err)
	}

	// ================================== Deploy the contract ==================================
	msg := types.NewMessage(eucommon.Alice, nil, 0, new(big.Int).SetUint64(0), 1e15, new(big.Int).SetUint64(1), evmcommon.Hex2Bytes(code), nil, true) // Build the message
	receipt, _, err := eu.Run(evmcommon.BytesToHash([]byte{1, 1, 1}), 1, &msg, ccEu.NewEVMBlockContext(config), ccEu.NewEVMTxContext(msg))            // Execute it

	_, transitions := eu.Api().Ccurl().ExportAll()
	eu.Api().Ccurl().Import(transitions)
	eu.Api().Ccurl().Sort()
	eu.Api().Ccurl().Commit([]uint32{1})

	// ---------------
	t.Log(receipt)
	contractAddress := receipt.ContractAddress
	if receipt.Status != 1 || err != nil {
		t.Error("Error: Deployment failed!!!", err)
	}

	// ================================== CallBasic() ==================================
	receipt, _, err = eu.Run(evmcommon.BytesToHash([]byte{1, 1, 1}), 1, &msg, cceu.NewEVMBlockContext(config), cceu.NewEVMTxContext(msg))
	_, transitions = eu.Api().Ccurl().ExportAll()
	if err != nil {
		fmt.Print(err)
	}

	data := crypto.Keccak256([]byte("call1()"))[:4]
	msg = types.NewMessage(eucommon.Alice, &contractAddress, 1, new(big.Int).SetUint64(0), 1e15, new(big.Int).SetUint64(1), data, nil, false)
	receipt, execResult, err := eu.Run(evmcommon.BytesToHash([]byte{1, 1, 1}), 1, &msg, cceu.NewEVMBlockContext(config), cceu.NewEVMTxContext(msg))
	_, transitions = eu.Api().Ccurl().ExportAll()

	if receipt.Status != 1 || err != nil {
		t.Error(err)
	}

	if execResult != nil && execResult.Err != nil {
		t.Error(execResult.Err)
	}
}

func TestCumulativeU256ThreadingMulti(t *testing.T) {
	eu, config, _, _, _ := NewTestEU()

	// ================================== Compile the contract ==================================
	currentPath, _ := os.Getwd()
	project := filepath.Dir(currentPath)
	// pyCompiler := project + "/compiler/compiler.py"
	targetPath := project + "/api/commutative/"

	code, err := compilers.CompileContracts(targetPath+"/u256", "u256Cumulative_test.sol", "0.5.0", "ThreadingCumulativeU256Multi", false)
	// code, err := compiler.CompileContracts(pyCompiler, targetPath+"/u256/u256Cumulative_test.sol", "ThreadingCumulativeU256Multi")

	if err != nil || len(code) == 0 {
		t.Error(err)
	}

	// ================================== Deploy the contract ==================================
	msg := types.NewMessage(eucommon.Alice, nil, 0, new(big.Int).SetUint64(0), 1e15, new(big.Int).SetUint64(1), evmcommon.Hex2Bytes(code), nil, true) // Build the message
	receipt, _, err := eu.Run(evmcommon.BytesToHash([]byte{1, 1, 1}), 1, &msg, ccEu.NewEVMBlockContext(config), ccEu.NewEVMTxContext(msg))            // Execute it

	_, transitions := eu.Api().Ccurl().ExportAll()
	eu.Api().Ccurl().Import(transitions)
	eu.Api().Ccurl().Sort()
	eu.Api().Ccurl().Commit([]uint32{1})

	// ---------------
	t.Log(receipt)
	contractAddress := receipt.ContractAddress
	if receipt.Status != 1 || err != nil {
		t.Error("Error: Deployment failed!!!", err)
	}

	// ================================== CallBasic() ==================================
	receipt, _, err = eu.Run(evmcommon.BytesToHash([]byte{1, 1, 1}), 1, &msg, cceu.NewEVMBlockContext(config), cceu.NewEVMTxContext(msg))
	_, transitions = eu.Api().Ccurl().ExportAll()
	if err != nil {
		fmt.Print(err)
	}

	data := crypto.Keccak256([]byte("testCase1()"))[:4]
	msg = types.NewMessage(eucommon.Alice, &contractAddress, 1, new(big.Int).SetUint64(0), 1e15, new(big.Int).SetUint64(1), data, nil, false)
	receipt, execResult, err := eu.Run(evmcommon.BytesToHash([]byte{1, 1, 1}), 1, &msg, cceu.NewEVMBlockContext(config), cceu.NewEVMTxContext(msg))
	_, transitions = eu.Api().Ccurl().ExportAll()

	if receipt.Status != 1 || err != nil {
		t.Error(err)
	}

	if execResult != nil && execResult.Err != nil {
		t.Error(execResult.Err)
	}
}

func TestU256Threading(t *testing.T) {
	eu, config, _, _, _ := NewTestEU()

	// ================================== Compile the contract ==================================
	currentPath, _ := os.Getwd()
	project := filepath.Dir(currentPath)
	// pyCompiler := project + "/compiler/compiler.py"
	targetPath := project + "/api/noncommutative/"

	code, err := compilers.CompileContracts(targetPath+"/u256", "u256_test.sol", "0.5.0", "U256ThreadingTest", false)
	// code, err := compiler.CompileContracts(pyCompiler, targetPath+"/u256/u256_test.sol", "U256ThreadingTest")

	if err != nil || len(code) == 0 {
		t.Error(err)
	}

	// ================================== Deploy the contract ==================================
	msg := types.NewMessage(eucommon.Alice, nil, 0, new(big.Int).SetUint64(0), 1e15, new(big.Int).SetUint64(1), evmcommon.Hex2Bytes(code), nil, true) // Build the message
	receipt, _, err := eu.Run(evmcommon.BytesToHash([]byte{1, 1, 1}), 1, &msg, ccEu.NewEVMBlockContext(config), ccEu.NewEVMTxContext(msg))            // Execute it
	_, transitions := eu.Api().Ccurl().ExportAll()
	// ---------------

	// t.Log("\n" + FormatTransitions(accesses))
	t.Log("\n" + eucommon.FormatTransitions(transitions))
	t.Log(receipt)
	contractAddress := receipt.ContractAddress
	if receipt.Status != 1 || err != nil {
		t.Error("Error: Deployment failed!!!", err)
	}

	// ================================== CallBasic() ==================================
	receipt, _, err = eu.Run(evmcommon.BytesToHash([]byte{1, 1, 1}), 1, &msg, cceu.NewEVMBlockContext(config), cceu.NewEVMTxContext(msg))
	_, transitions = eu.Api().Ccurl().ExportAll()

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

func TestArrayThreading(t *testing.T) {
	eu, config, _, _, _ := NewTestEU()

	// ================================== Compile the contract ==================================
	currentPath, _ := os.Getwd()
	project := filepath.Dir(currentPath)
	// pyCompiler := project + "/compiler/compiler.py"
	targetPath := project + "/api/noncommutative/"

	code, err := compilers.CompileContracts(targetPath+"/u256", "u256_test.sol", "0.5.0", "ArrayThreadingTest", false)
	// code, err := compiler.CompileContracts(pyCompiler, targetPath+"/u256/u256_test.sol", "ArrayThreadingTest")

	if err != nil || len(code) == 0 {
		t.Error(err)
	}

	// ================================== Deploy the contract ==================================
	msg := types.NewMessage(eucommon.Alice, nil, 0, new(big.Int).SetUint64(0), 1e15, new(big.Int).SetUint64(1), evmcommon.Hex2Bytes(code), nil, true) // Build the message
	receipt, _, err := eu.Run(evmcommon.BytesToHash([]byte{1, 1, 1}), 1, &msg, ccEu.NewEVMBlockContext(config), ccEu.NewEVMTxContext(msg))            // Execute it
	_, transitions := eu.Api().Ccurl().ExportAll()

	// ---------------

	// t.Log("\n" + FormatTransitions(accesses))
	t.Log("\n" + eucommon.FormatTransitions(transitions))
	t.Log(receipt)
	contractAddress := receipt.ContractAddress
	if receipt.Status != 1 || err != nil {
		t.Error("Error: Deployment failed!!!", err)
	}

	// ================================== CallBasic() ==================================
	receipt, _, err = eu.Run(evmcommon.BytesToHash([]byte{1, 1, 1}), 1, &msg, cceu.NewEVMBlockContext(config), cceu.NewEVMTxContext(msg))
	_, transitions = eu.Api().Ccurl().ExportAll()

	if err != nil {
		fmt.Print(err)
	}

	data := crypto.Keccak256([]byte("call()"))[:4]
	msg = types.NewMessage(eucommon.Alice, &contractAddress, 1, new(big.Int).SetUint64(0), 1e15, new(big.Int).SetUint64(1), data, nil, false)
	receipt, execResult, err := eu.Run(evmcommon.BytesToHash([]byte{1, 1, 1}), 1, &msg, cceu.NewEVMBlockContext(config), cceu.NewEVMTxContext(msg))
	_, transitions = eu.Api().Ccurl().ExportAll()

	if receipt.Status != 1 || err != nil {
		t.Error(err)
	}

	if execResult != nil && execResult.Err != nil {
		t.Error(execResult.Err)
	}
}
