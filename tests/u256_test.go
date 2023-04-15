package tests

import (
	"math/big"
	"os"
	"path/filepath"
	"testing"

	common "github.com/arcology-network/common-lib/common"
	evmcommon "github.com/arcology-network/evm/common"
	"github.com/arcology-network/evm/core/types"
	ccEu "github.com/arcology-network/vm-adaptor"
	eucommon "github.com/arcology-network/vm-adaptor/common"
	compiler "github.com/arcology-network/vm-adaptor/compiler"
)

func TestU256Dynamic(t *testing.T) {
	eu, config, _, _ := NewTestEU()

	// ================================== Compile the contract ==================================
	currentPath, _ := os.Getwd()
	project := filepath.Dir(currentPath)
	pyCompiler := project + "/compiler/compiler.py"
	targetPath := project + "/api/types/"
	baseFile := targetPath + "base/Base.sol"

	if err := common.CopyFile(baseFile, targetPath+"/u256/Base.sol"); err != nil {
		t.Error(err)
	}

	code, err := compiler.CompileContracts(pyCompiler, targetPath+"/u256/u256_test.sol", "U256DynamicTest")
	os.Remove(targetPath + "/u256/Base.sol")

	if err != nil || len(code) == 0 {
		t.Error("Error: Failed to generate the byte code")
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

func TestU256Fixed(t *testing.T) {
	eu, config, _, _ := NewTestEU()

	// ================================== Compile the contract ==================================
	currentPath, _ := os.Getwd()
	project := filepath.Dir(currentPath)
	pyCompiler := project + "/compiler/compiler.py"
	targetPath := project + "/api/types/"
	baseFile := targetPath + "base/Base.sol"

	if err := common.CopyFile(baseFile, targetPath+"/u256/Base.sol"); err != nil {
		t.Error(err)
	}

	code, err := compiler.CompileContracts(pyCompiler, targetPath+"/u256/u256Fixed_test.sol", "U256FixedTest")
	os.Remove(targetPath + "/u256/Base.sol")

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

func TestCumulativeU256(t *testing.T) {
	eu, config, _, _ := NewTestEU()

	// ================================== Compile the contract ==================================
	currentPath, _ := os.Getwd()
	project := filepath.Dir(currentPath)
	pyCompiler := project + "/compiler/compiler.py"
	targetPath := project + "/api/types/"
	baseFile := targetPath + "base/Base.sol"

	if err := common.CopyFile(baseFile, targetPath+"/u256/Base.sol"); err != nil {
		t.Error(err)
	}

	code, err := compiler.CompileContracts(pyCompiler, targetPath+"/u256/u256Cumulative_test.sol", "CumulativeU256Test")
	os.Remove(targetPath + "/u256/Base.sol")

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
