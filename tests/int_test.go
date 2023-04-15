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

func TestContractInt(t *testing.T) {
	eu, config, _, _ := NewTestEU()
	// ================================== Compile the contract ==================================
	currentPath, _ := os.Getwd()
	project := filepath.Dir(currentPath)
	pyCompiler := project + "/compiler/compiler.py"
	targetPath := project + "/api/types/"
	baseFile := targetPath + "base/Base.sol"

	if err := common.CopyFile(baseFile, targetPath+"int/Base.sol"); err != nil {
		t.Error(err)
	}

	code, err := compiler.CompileContracts(pyCompiler, targetPath+"int/int_test.sol", "IntTest")
	os.Remove(targetPath + "/int/Base.sol")

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

	// ================================== Call length() ==================================
	// url = concurrenturl.NewConcurrentUrl(db)
	// url.Import(transitions)
	// url.PostImport()
	// errs := url.Commit([]uint32{1})
	// if len(errs) != 0 {
	// 	t.Error(errs)
	// 	return
	// }
	// api = ccApi.NewAPI(url)
	// statedb = eth.NewImplStateDB(url)
	// eu = ccEu.NewEU(config.ChainConfig, *config.VMConfig, config.Chain, statedb, api)

	// config.BlockNumber = new(big.Int).SetUint64(10000001)
	// config.Time = new(big.Int).SetUint64(10000001)

	// data := crypto.Keccak256([]byte("length()"))[:4]
	// data = append(data, evmcommon.BytesToHash(User1.Bytes()).Bytes()...)
	// data = append(data, evmcommon.BytesToHash([]byte{0xcc}).Bytes()...)
	// msg = types.NewMessage(eucommon.User1, &contractAddress, 1, new(big.Int).SetUint64(0), 1e15, new(big.Int).SetUint64(1), data, nil, true)
	// _, transitions, receipt, err = eu.Run(evmcommon.BytesToHash([]byte{2, 2, 2}), 2, &msg, ccEu.NewEVMBlockContext(config), ccEu.NewEVMTxContext(msg))
	// t.Log("\n" + FormatTransitions(transitions))
	// t.Log(receipt)
	// if receipt.Status != 1 {
	// 	t.Error("Error: Failed to calll length()!!!", err)
	// }

}
