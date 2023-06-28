package tests

import (
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"testing"

	concurrenturl "github.com/arcology-network/concurrenturl"
	evmcommon "github.com/arcology-network/evm/common"
	"github.com/arcology-network/evm/core"
	"github.com/arcology-network/evm/crypto"
	evmeu "github.com/arcology-network/vm-adaptor"
	eucommon "github.com/arcology-network/vm-adaptor/common"
	"github.com/arcology-network/vm-adaptor/compiler"
)

func TestBase(t *testing.T) {
	eu, config, db, url, _ := NewTestEU()

	// ================================== Compile the contract ==================================
	currentPath, _ := os.Getwd()
	project := filepath.Dir(currentPath)

	targetPath := project + "/api/"

	code, err := compiler.CompileContracts(targetPath+"noncommutative/", "base/base_test.sol", "0.8.19", "BaseTest", false)
	if err != nil || len(code) == 0 {
		t.Error("Error: Failed to generate the byte code")
	}

	// ================================== Deploy the contract ==================================
	msg := core.NewMessage(eucommon.Alice, nil, 0, new(big.Int).SetUint64(0), 1e15, new(big.Int).SetUint64(1), evmcommon.Hex2Bytes(code), nil, true) // Build the message
	receipt, _, err := eu.Run(evmcommon.BytesToHash([]byte{1, 1, 1}), 1, &msg, evmeu.NewEVMBlockContext(config), evmeu.NewEVMTxContext(msg))         // Execute it
	_, transitions := eu.Api().Ccurl().ExportAll()

	t.Log("\n" + eucommon.FormatTransitions(transitions))
	// t.Log(receipt)

	if receipt.Status != 1 || err != nil {
		t.Error("Error: Deployment failed!!!", err)
	}

	contractAddress := receipt.ContractAddress
	url = concurrenturl.NewConcurrentUrl(db)
	url.Import(transitions)
	url.Sort()
	url.Commit([]uint32{1})

	// ================================== Call() ==================================
	receipt, _, err = eu.Run(evmcommon.BytesToHash([]byte{1, 1, 1}), 1, &msg, evmeu.NewEVMBlockContext(config), evmeu.NewEVMTxContext(msg))
	// _, transitions = eu.Api().Ccurl().ExportAll()

	if err != nil {
		fmt.Print(err)
	}

	data := crypto.Keccak256([]byte("call()"))[:4]
	msg = core.NewMessage(eucommon.Alice, &contractAddress, 1, new(big.Int).SetUint64(0), 1e15, new(big.Int).SetUint64(1), data, nil, false)
	receipt, execResult, err := eu.Run(evmcommon.BytesToHash([]byte{1, 1, 1}), 1, &msg, evmeu.NewEVMBlockContext(config), evmeu.NewEVMTxContext(msg))
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
