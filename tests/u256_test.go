package tests

import (
	"math/big"
	"os"
	"path/filepath"
	"testing"

	cachedstorage "github.com/arcology-network/common-lib/cachedstorage"
	common "github.com/arcology-network/common-lib/common"
	"github.com/arcology-network/concurrenturl/v2"
	ccurlcommon "github.com/arcology-network/concurrenturl/v2/common"
	ccurlstorage "github.com/arcology-network/concurrenturl/v2/storage"
	"github.com/arcology-network/concurrenturl/v2/type/commutative"
	evmcommon "github.com/arcology-network/evm/common"
	"github.com/arcology-network/evm/core/types"
	ccEu "github.com/arcology-network/vm-adaptor"
	ccApi "github.com/arcology-network/vm-adaptor/api"
	compiler "github.com/arcology-network/vm-adaptor/compiler"
	eth "github.com/arcology-network/vm-adaptor/eth"
)

func TestU256Dynamic(t *testing.T) {
	config := MainConfig()
	persistentDB := cachedstorage.NewDataStore()
	meta, _ := commutative.NewMeta(ccurlcommon.NewPlatform().Eth10Account())
	persistentDB.Inject(ccurlcommon.NewPlatform().Eth10Account(), meta)
	db := ccurlstorage.NewTransientDB(persistentDB)

	url := concurrenturl.NewConcurrentUrl(db)
	statedb := eth.NewImplStateDB(url)
	statedb.Prepare(evmcommon.Hash{}, evmcommon.Hash{}, 0)
	statedb.CreateAccount(Coinbase)
	statedb.CreateAccount(User1)
	statedb.AddBalance(User1, new(big.Int).SetUint64(1e18))
	_, transitions := url.Export(true)
	t.Log("\n" + FormatTransitions(transitions))

	// Deploy.
	url = concurrenturl.NewConcurrentUrl(db)
	url.Import(transitions)
	url.PostImport()
	url.Commit([]uint32{0})
	api := ccApi.NewAPI(url)
	statedb = eth.NewImplStateDB(url)
	eu := ccEu.NewEU(config.ChainConfig, *config.VMConfig, config.Chain, statedb, api, url)

	config.Coinbase = &Coinbase
	config.BlockNumber = new(big.Int).SetUint64(10000000)
	config.Time = new(big.Int).SetUint64(10000000)

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
	msg := types.NewMessage(User1, nil, 0, new(big.Int).SetUint64(0), 1e15, new(big.Int).SetUint64(1), evmcommon.Hex2Bytes(code), nil, true)              // Build the message
	_, transitions, receipt, err := eu.Run(evmcommon.BytesToHash([]byte{1, 1, 1}), 1, &msg, ccEu.NewEVMBlockContextV2(config), ccEu.NewEVMTxContext(msg)) // Execute it
	// ---------------

	// t.Log("\n" + FormatTransitions(accesses))
	t.Log("\n" + FormatTransitions(transitions))
	t.Log(receipt)
	// contractAddress := receipt.ContractAddress
	if receipt.Status != 1 || err != nil {
		t.Error("Error: Deployment failed!!!", err)
	}
}

func TestU256Fixed(t *testing.T) {
	config := MainConfig()
	persistentDB := cachedstorage.NewDataStore()
	meta, _ := commutative.NewMeta(ccurlcommon.NewPlatform().Eth10Account())
	persistentDB.Inject(ccurlcommon.NewPlatform().Eth10Account(), meta)
	db := ccurlstorage.NewTransientDB(persistentDB)

	url := concurrenturl.NewConcurrentUrl(db)
	statedb := eth.NewImplStateDB(url)
	statedb.Prepare(evmcommon.Hash{}, evmcommon.Hash{}, 0)
	statedb.CreateAccount(Coinbase)
	statedb.CreateAccount(User1)
	statedb.AddBalance(User1, new(big.Int).SetUint64(1e18))
	_, transitions := url.Export(true)
	t.Log("\n" + FormatTransitions(transitions))

	// Deploy.
	url = concurrenturl.NewConcurrentUrl(db)
	url.Import(transitions)
	url.PostImport()
	url.Commit([]uint32{0})
	api := ccApi.NewAPI(url)
	statedb = eth.NewImplStateDB(url)
	eu := ccEu.NewEU(config.ChainConfig, *config.VMConfig, config.Chain, statedb, api, url)

	config.Coinbase = &Coinbase
	config.BlockNumber = new(big.Int).SetUint64(10000000)
	config.Time = new(big.Int).SetUint64(10000000)

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
	msg := types.NewMessage(User1, nil, 0, new(big.Int).SetUint64(0), 1e15, new(big.Int).SetUint64(1), evmcommon.Hex2Bytes(code), nil, true)              // Build the message
	_, transitions, receipt, err := eu.Run(evmcommon.BytesToHash([]byte{1, 1, 1}), 1, &msg, ccEu.NewEVMBlockContextV2(config), ccEu.NewEVMTxContext(msg)) // Execute it
	// ---------------

	// t.Log("\n" + FormatTransitions(accesses))
	t.Log("\n" + FormatTransitions(transitions))
	t.Log(receipt)
	// contractAddress := receipt.ContractAddress
	if receipt.Status != 1 || err != nil {
		t.Error("Error: Deployment failed!!!", err)
	}
}

func TestCumulativeU256(t *testing.T) {
	config := MainConfig()
	persistentDB := cachedstorage.NewDataStore()
	meta, _ := commutative.NewMeta(ccurlcommon.NewPlatform().Eth10Account())
	persistentDB.Inject(ccurlcommon.NewPlatform().Eth10Account(), meta)
	db := ccurlstorage.NewTransientDB(persistentDB)

	url := concurrenturl.NewConcurrentUrl(db)
	statedb := eth.NewImplStateDB(url)
	statedb.Prepare(evmcommon.Hash{}, evmcommon.Hash{}, 0)
	statedb.CreateAccount(Coinbase)
	statedb.CreateAccount(User1)
	statedb.AddBalance(User1, new(big.Int).SetUint64(1e18))
	_, transitions := url.Export(true)
	t.Log("\n" + FormatTransitions(transitions))

	// Deploy.
	url = concurrenturl.NewConcurrentUrl(db)
	url.Import(transitions)
	url.PostImport()
	url.Commit([]uint32{0})
	api := ccApi.NewAPI(url)
	statedb = eth.NewImplStateDB(url)
	eu := ccEu.NewEU(config.ChainConfig, *config.VMConfig, config.Chain, statedb, api, url)

	config.Coinbase = &Coinbase
	config.BlockNumber = new(big.Int).SetUint64(10000000)
	config.Time = new(big.Int).SetUint64(10000000)

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
	msg := types.NewMessage(User1, nil, 0, new(big.Int).SetUint64(0), 1e15, new(big.Int).SetUint64(1), evmcommon.Hex2Bytes(code), nil, true)              // Build the message
	_, transitions, receipt, err := eu.Run(evmcommon.BytesToHash([]byte{1, 1, 1}), 1, &msg, ccEu.NewEVMBlockContextV2(config), ccEu.NewEVMTxContext(msg)) // Execute it
	// ---------------

	// t.Log("\n" + FormatTransitions(accesses))
	t.Log("\n" + FormatTransitions(transitions))
	t.Log(receipt)
	// contractAddress := receipt.ContractAddress
	if receipt.Status != 1 || err != nil {
		t.Error("Error: Deployment failed!!!", err)
	}

}
