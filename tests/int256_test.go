package tests

import (
	"math/big"
	"os"
	"path/filepath"
	"testing"

	"github.com/arcology-network/common-lib/common"
	evmcommon "github.com/arcology-network/evm/common"
	"github.com/arcology-network/evm/core/types"
	ccEu "github.com/arcology-network/vm-adaptor"
	eucommon "github.com/arcology-network/vm-adaptor/common"
	compiler "github.com/arcology-network/vm-adaptor/compiler"
)

func TestContractNoncommutativeInt256(t *testing.T) {
	eu, config, _, _, _ := NewTestEU()
	// ================================== Compile the contract ==================================
	currentPath, _ := os.Getwd()
	project := filepath.Dir(currentPath)
	pyCompiler := project + "/compiler/compiler.py"
	targetPath := project + "/api/noncommutative/"
	baseFile := targetPath + "base/Base.sol"

	if err := common.CopyFile(baseFile, targetPath+"int256/Base.sol"); err != nil {
		t.Error(err)
	}

	if err := common.CopyFile(project+"/api/threading/Threading.sol", targetPath+"/int256/Threading.sol"); err != nil {
		t.Error(err)
	}

	code, err := compiler.CompileContracts(pyCompiler, targetPath+"int256/int256_test.sol", "Int256Test")

	if err != nil || len(code) == 0 {
		t.Error("Error: Failed to generate the byte code")
	}
	// ================================== Deploy the contract ==================================
	msg := types.NewMessage(eucommon.User1, nil, 0, new(big.Int).SetUint64(0), 1e15, new(big.Int).SetUint64(1), evmcommon.Hex2Bytes(code), nil, true)      // Build the message
	_, transitions, receipt, _, err := eu.Run(evmcommon.BytesToHash([]byte{1, 1, 1}), 1, &msg, ccEu.NewEVMBlockContext(config), ccEu.NewEVMTxContext(msg)) // Execute it
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
	pyCompiler := project + "/compiler/compiler.py"
	targetPath := project + "/api/noncommutative/"
	baseFile := targetPath + "base/Base.sol"

	if err := common.CopyFile(baseFile, targetPath+"/int256/Base.sol"); err != nil {
		t.Error(err)
	}

	if err := common.CopyFile(project+"/api/threading/Threading.sol", targetPath+"/int256/Threading.sol"); err != nil {
		t.Error(err)
	}

	if err := common.CopyFile(project+"/api/threading/Threading.sol", targetPath+"/int256/Threading.sol"); err != nil {
		t.Error(err)
	}

	code, err := compiler.CompileContracts(pyCompiler, targetPath+"/int256/int256N_test.sol", "Int64NTest")

	if err != nil || len(code) == 0 {
		t.Error(err)
	}
	// ================================== Deploy the contract ==================================
	msg := types.NewMessage(eucommon.User1, nil, 0, new(big.Int).SetUint64(0), 1e15, new(big.Int).SetUint64(1), evmcommon.Hex2Bytes(code), nil, true)      // Build the message
	_, transitions, receipt, _, err := eu.Run(evmcommon.BytesToHash([]byte{1, 1, 1}), 1, &msg, ccEu.NewEVMBlockContext(config), ccEu.NewEVMTxContext(msg)) // Execute it
	// ---------------

	// t.Log("\n" + FormatTransitions(accesses))
	t.Log("\n" + eucommon.FormatTransitions(transitions))
	t.Log(receipt)
	// contractAddress := receipt.ContractAddress
	if receipt.Status != 1 || err != nil {
		t.Error("Error: Deployment failed!!!", err)
	}
}

func TestCumulativeInt256(t *testing.T) {
	eu, config, _, _, _ := NewTestEU()

	// ================================== Compile the contract ==================================
	currentPath, _ := os.Getwd()
	project := filepath.Dir(currentPath)
	pyCompiler := project + "/compiler/compiler.py"
	targetPath := project + "/api/commutative/"

	code, err := compiler.CompileContracts(pyCompiler, targetPath+"/int256/int256Cumulative_test.sol", "Int256CumulativeTest")

	if err != nil || len(code) == 0 {
		t.Error(err)
	}

	// ================================== Deploy the contract ==================================
	msg := types.NewMessage(eucommon.User1, nil, 0, new(big.Int).SetUint64(0), 1e15, new(big.Int).SetUint64(1), evmcommon.Hex2Bytes(code), nil, true)      // Build the message
	_, transitions, receipt, _, err := eu.Run(evmcommon.BytesToHash([]byte{1, 1, 1}), 1, &msg, ccEu.NewEVMBlockContext(config), ccEu.NewEVMTxContext(msg)) // Execute it
	// ---------------

	// t.Log("\n" + FormatTransitions(accesses))
	t.Log("\n" + eucommon.FormatTransitions(transitions))
	t.Log(receipt)
	// contractAddress := receipt.ContractAddress
	if receipt.Status != 1 || err != nil {
		t.Error("Error: Deployment failed!!!", err)
	}
}

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
// 	msg := types.NewMessage(eucommon.User1, nil, 0, new(big.Int).SetUint64(0), 1e15, new(big.Int).SetUint64(1), evmcommon.Hex2Bytes(code), nil, false) // Build the message
// 	_, _, receipt, _, err := eu.Run(evmcommon.BytesToHash([]byte{1, 1, 1}), 1, &msg, cceu.NewEVMBlockContext(config), cceu.NewEVMTxContext(msg))       // Execute it
// 	// t.Log("\n" + eucommon.FormatTransitions(transitions))
// 	// univalue.Univalues(transitions).Print()
// 	if receipt.Status != 1 || err != nil {
// 		t.Error("Error: Deployment failed!!!", err)
// 	}
// 	fmt.Println(receipt.ContractAddress)

// 	// // ================================== CallBasic() ==================================
// 	data := crypto.Keccak256([]byte("call()"))[:4]
// 	contractAddress := receipt.ContractAddress
// 	msg = types.NewMessage(eucommon.User1, &contractAddress, 1, new(big.Int).SetUint64(0), 1e15, new(big.Int).SetUint64(1), data, nil, false)
// 	_, _, receipt, execResult, err := eu.Run(evmcommon.BytesToHash([]byte{1, 1, 1}), 1, &msg, cceu.NewEVMBlockContext(config), cceu.NewEVMTxContext(msg))
// 	// univalue.Univalues(transitions).Print()

// 	if err != nil {
// 		t.Error(err)
// 	}

// 	if execResult != nil && execResult.Err != nil {
// 		t.Error(execResult.Err)
// 	}
// }
