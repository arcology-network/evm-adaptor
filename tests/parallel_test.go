package tests

import (
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"testing"

	"github.com/arcology-network/concurrenturl/v2"
	evmcommon "github.com/arcology-network/evm/common"
	"github.com/arcology-network/evm/core/types"
	"github.com/arcology-network/evm/crypto"
	cceu "github.com/arcology-network/vm-adaptor"
	eucommon "github.com/arcology-network/vm-adaptor/common"
	compiler "github.com/arcology-network/vm-adaptor/compiler"
)

func TestParallelInvoke(t *testing.T) {
	eu, config, db, url := NewTestEU()

	// ================================== Compile the contract ==================================
	currentPath, _ := os.Getwd()
	project := filepath.Dir(currentPath)
	pyCompiler := project + "/compiler/compiler.py"

	code, err := compiler.CompileContracts(pyCompiler, project+"/api/parallel/parallel_test.sol", "ParallelInvokeTest")

	if err != nil || len(code) == 0 {
		t.Error("Error: Failed to generate the byte code")
	}
	// ================================== Deploy the contract ==================================
	msg := types.NewMessage(eucommon.User1, nil, 0, new(big.Int).SetUint64(0), 1e15, new(big.Int).SetUint64(1), evmcommon.Hex2Bytes(code), nil, false)     // Build the message
	_, transitions, receipt, _, err := eu.Run(evmcommon.BytesToHash([]byte{1, 1, 1}), 1, &msg, cceu.NewEVMBlockContext(config), cceu.NewEVMTxContext(msg)) // Execute it
	t.Log("\n" + eucommon.FormatTransitions(transitions))
	// ---------------

	contractAddress := receipt.ContractAddress
	if receipt.Status != 1 || err != nil {
		t.Error("Error: Deployment failed!!!", err)
	}
	fmt.Println(receipt.ContractAddress)

	// ================================== Call length() ==================================
	url = concurrenturl.NewConcurrentUrl(db)
	url.Import(transitions)
	url.PostImport()
	errs := url.Commit([]uint32{1})
	if len(errs) != 0 {
		t.Error(errs)
		return
	}

	_, transitions, receipt, _, err = eu.Run(evmcommon.BytesToHash([]byte{1, 1, 1}), 1, &msg, cceu.NewEVMBlockContext(config), cceu.NewEVMTxContext(msg))
	if err != nil {
		fmt.Print(err)
	}

	data := crypto.Keccak256([]byte("callPara()"))[:4]
	// data = append(data, evmcommon.BytesToHash(eucommon.User1.Bytes()).Bytes()...)
	// data = append(data, evmcommon.BytesToHash([]byte{0xcc}).Bytes()...)
	msg = types.NewMessage(eucommon.User1, &contractAddress, 1, new(big.Int).SetUint64(0), 1e15, new(big.Int).SetUint64(1), data, nil, false)

	_, transitions, receipt, execResult, err := eu.Run(evmcommon.BytesToHash([]byte{1, 1, 1}), 1, &msg, cceu.NewEVMBlockContext(config), cceu.NewEVMTxContext(msg))
	if err != nil {
		t.Error(err)
	}

	if execResult != nil && execResult.Err != nil {
		t.Error(execResult.Err)
	}

	// t.Log("\n" + FormatTransitions(transitions))
	// t.Log(receipt)
	// if receipt.Status != 1 {
	// 	t.Error("Error: Failed to calll length()!!!", err)
	// }

}
