package tests

import (
	"math/big"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/arcology-network/common-lib/codec"
	"github.com/arcology-network/common-lib/common"
	evmcore "github.com/arcology-network/evm/core"
	"github.com/arcology-network/evm/crypto"
	eucommon "github.com/arcology-network/vm-adaptor/common"
	"github.com/arcology-network/vm-adaptor/compiler"
	tests "github.com/arcology-network/vm-adaptor/tests"
)

func TestSubcurrencyMint(t *testing.T) {
	currentPath, _ := os.Getwd()
	targetPath := path.Join((path.Dir(filepath.Dir(currentPath))), "concurrentlib/")

	err, eu, receipt := tests.DeployThenInvoke(targetPath, "examples/subcurrency/Subcurrency.sol", "0.8.19", "Coin", "", []byte{}, false)
	if err != nil {
		t.Error(err)
		return
	}

	coinAddress := receipt.ContractAddress

	var execResult *evmcore.ExecutionResult
	err, eu, execResult, receipt = tests.CallContract(eu, receipt.ContractAddress, crypto.Keccak256([]byte("getter()"))[:4], 0, false)

	code, err := compiler.CompileContracts(targetPath, "examples/subcurrency/subcurrency_test.sol", "0.8.19", "SubcurrencyCaller", false)
	if err != nil || len(code) == 0 {
		t.Error(err)
	}

	// data := crypto.Keccak256([]byte(funcName))[:4]
	// tests.CallContract(eu, receipt.ContractAddress, []byte{}, false)

	config := tests.MainTestConfig()
	config.Coinbase = &eucommon.Coinbase
	config.BlockNumber = new(big.Int).SetUint64(10000000)
	config.Time = new(big.Int).SetUint64(10000000)
	err, eu, execResult, receipt = tests.CallContract(eu, coinAddress, crypto.Keccak256([]byte("getter()"))[:4], 0, false)
	err, config, eu, receipt = tests.DepolyContract(eu, config, code, "", []byte{}, 2, false)
	if err != nil || receipt.Status != 1 {
		t.Error(err)
	}
	// err, eu, execResult, receipt = tests.CallContract(eu, coinAddress, crypto.Keccak256([]byte("getter()"))[:4], 0, false)

	addr := codec.Bytes32{}.Decode(common.PadLeft(coinAddress[:], 0, 32)).(codec.Bytes32) // Callee contract address
	funCall := crypto.Keccak256([]byte("call(address)"))[:4]
	funCall = append(funCall, addr[:]...)
	err, eu, execResult, receipt = tests.CallContract(eu, receipt.ContractAddress, funCall, 0, false)
	if receipt.Status != 1 {
		t.Error(execResult.Err)
	}
}
