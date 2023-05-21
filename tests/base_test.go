package tests

import (
	"math/big"
	"os"
	"path/filepath"
	"testing"

	"github.com/arcology-network/concurrenturl"
	evmcommon "github.com/arcology-network/evm/common"
	"github.com/arcology-network/evm/core/types"
	"github.com/arcology-network/evm/crypto"
	cceu "github.com/arcology-network/vm-adaptor"
	ccapi "github.com/arcology-network/vm-adaptor/api"
	eucommon "github.com/arcology-network/vm-adaptor/common"
	compiler "github.com/arcology-network/vm-adaptor/compiler"
	eth "github.com/arcology-network/vm-adaptor/eth"
)

func TestBase(t *testing.T) {
	eu, config, db, url, _ := NewTestEU()

	// ================================== Compile the contract ==================================
	currentPath, _ := os.Getwd()
	project := filepath.Dir(currentPath)
	pyCompiler := project + "/compiler/compiler.py"
	targetPath := project + "/api/noncommutative/"

	code, err := compiler.CompileContracts(pyCompiler, targetPath+"/base/base_test.sol", "BaseTest")
	if err != nil || len(code) == 0 {
		t.Error("Error: Failed to generate the byte code")
	}

	// ================================== Deploy the contract ==================================
	msg := types.NewMessage(eucommon.Alice, nil, 0, new(big.Int).SetUint64(0), 1e15, new(big.Int).SetUint64(1), evmcommon.Hex2Bytes(code), nil, true) // Build the message
	receipt, _, err := eu.Run(evmcommon.BytesToHash([]byte{1, 1, 1}), 1, &msg, cceu.NewEVMBlockContext(config), cceu.NewEVMTxContext(msg))            // Execute it
	_, transitions := eu.Api().Ccurl().ExportAll()

	t.Log("\n" + eucommon.FormatTransitions(transitions))
	// t.Log(receipt)
	contractAddress := receipt.ContractAddress
	if receipt.Status != 1 || err != nil {
		t.Error("Error: Deployment failed!!!", err)
	}

	// ================================== Call length() ==================================
	url = concurrenturl.NewConcurrentUrl(db)
	url.Import(transitions)
	url.Sort()
	url.Commit([]uint32{1})

	statedb := eth.NewImplStateDB(url)
	eu = cceu.NewEU(config.ChainConfig, *config.VMConfig, statedb, ccapi.NewAPI(url))

	config.BlockNumber = new(big.Int).SetUint64(10000001)
	config.Time = new(big.Int).SetUint64(10000001)

	data := crypto.Keccak256([]byte("length()"))[:4]
	data = append(data, evmcommon.BytesToHash(eucommon.Alice.Bytes()).Bytes()...)
	data = append(data, evmcommon.BytesToHash([]byte{0xcc}).Bytes()...)
	msg = types.NewMessage(eucommon.Alice, &contractAddress, 1, new(big.Int).SetUint64(0), 1e15, new(big.Int).SetUint64(1), data, nil, true)
	receipt, _, err = eu.Run(evmcommon.BytesToHash([]byte{2, 2, 2}), 2, &msg, cceu.NewEVMBlockContext(config), cceu.NewEVMTxContext(msg))
	_, transitions = eu.Api().Ccurl().ExportAll()

	t.Log("\n" + eucommon.FormatTransitions(transitions))
	t.Log(receipt)
	if receipt.Status != 1 {
		t.Error("Error: Failed to call length()!!!", err)
	}
}
