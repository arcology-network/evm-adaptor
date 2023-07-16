package tests

import (
	"math/big"
	"os"
	"path/filepath"
	"testing"

	commontype "github.com/arcology-network/common-lib/types"
	evmcommon "github.com/arcology-network/evm/common"
	"github.com/arcology-network/evm/core"
	eucommon "github.com/arcology-network/vm-adaptor/common"
	"github.com/arcology-network/vm-adaptor/compiler"
	execution "github.com/arcology-network/vm-adaptor/execution"
)

func TestContainerPair(t *testing.T) {
	eu, config, _, _, _ := NewTestEU()

	// ================================== Compile the contract ==================================
	currentPath, _ := os.Getwd()
	project := filepath.Dir(currentPath)
	targetPath := project + "/api/"

	code, err := compiler.CompileContracts(targetPath, "combo/bytes_bool_test.sol", "0.8.19", "PairTest", false)
	if err != nil || len(code) == 0 {
		t.Error("Error: Failed to generate the byte code")
	}

	// ================================== Deploy the contract ==================================
	msg := core.NewMessage(eucommon.Alice, nil, 0, new(big.Int).SetUint64(0), 1e15, new(big.Int).SetUint64(1), evmcommon.Hex2Bytes(code), nil, true)
	stdMsg := &execution.StandardMessage{
		ID:     1,
		TxHash: [32]byte{1, 1, 1},
		Native: &msg, // Build the message
		Source: commontype.TX_SOURCE_LOCAL,
	}

	receipt, _, err := eu.Run(stdMsg, execution.NewEVMBlockContext(config), execution.NewEVMTxContext(*stdMsg.Native)) // Execute it
	// _, transitions := eu.Api().StateFilter().ByType()

	// msg := core.NewMessage(eucommon.Alice, nil, 0, new(big.Int).SetUint64(0), 1e15, new(big.Int).SetUint64(1), evmcommon.Hex2Bytes(code), nil, true) // Build the message
	// receipt, _, err := eu.Run(evmcommon.BytesToHash([]byte{1, 1, 1}), 1, &msg, execution.NewEVMBlockContext(config), execution.NewEVMTxContext(msg)) // Execute it
	// _, transitions := eu.Api().StateFilter().ByType()

	// ---------------

	// t.Log("\n" + FormatTransitions(accesses))
	//t.Log("\n" + eucommon.FormatTransitions(transitions))
	// t.Log(receipt)
	// contractAddress := receipt.ContractAddress
	if receipt.Status != 1 || err != nil {
		t.Error("Error: Deployment failed!!!", err)
	}
}

func TestSetTest(t *testing.T) {
	currentPath, _ := os.Getwd()
	targetPath := filepath.Dir(currentPath) + "/api/"

	err, _ := InvokeTestContract(targetPath, "combo/set_test.sol", "0.8.19", "SetTest", "", []byte{}, false)
	if err != nil {
		t.Error(err)
	}
}

func TestMapTest(t *testing.T) {
	currentPath, _ := os.Getwd()
	targetPath := filepath.Dir(currentPath) + "/api/"

	err, _ := InvokeTestContract(targetPath, "combo/map_test.sol", "0.8.19", "MapTest", "", []byte{}, false)
	if err != nil {
		t.Error(err)
	}
}
