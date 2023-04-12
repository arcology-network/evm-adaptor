package api

import (
	"math/big"
	"os"
	"path/filepath"
	"testing"

	cachedstorage "github.com/arcology-network/common-lib/cachedstorage"
	"github.com/arcology-network/concurrenturl/v2"
	ccurlcommon "github.com/arcology-network/concurrenturl/v2/common"
	ccurlstorage "github.com/arcology-network/concurrenturl/v2/storage"
	"github.com/arcology-network/concurrenturl/v2/type/commutative"
	evmcommon "github.com/arcology-network/evm/common"
	"github.com/arcology-network/evm/core/types"
	ccEu "github.com/arcology-network/vm-adaptor"
	ccApi "github.com/arcology-network/vm-adaptor/api"
	eth "github.com/arcology-network/vm-adaptor/eth"
	tests "github.com/arcology-network/vm-adaptor/tests"
)

func TestCumulativeU256(t *testing.T) {
	config := tests.MainConfig()
	persistentDB := cachedstorage.NewDataStore()
	meta, _ := commutative.NewMeta(ccurlcommon.NewPlatform().Eth10Account())
	persistentDB.Inject(ccurlcommon.NewPlatform().Eth10Account(), meta)
	db := ccurlstorage.NewTransientDB(persistentDB)

	url := concurrenturl.NewConcurrentUrl(db)
	statedb := eth.NewImplStateDB(url)
	statedb.Prepare(evmcommon.Hash{}, evmcommon.Hash{}, 0)
	statedb.CreateAccount(tests.Coinbase)
	statedb.CreateAccount(tests.User1)
	statedb.AddBalance(tests.User1, new(big.Int).SetUint64(1e18))
	_, transitions := url.Export(true)
	t.Log("\n" + tests.FormatTransitions(transitions))

	// Deploy.
	url = concurrenturl.NewConcurrentUrl(db)
	url.Import(transitions)
	url.PostImport()
	url.Commit([]uint32{0})
	api := ccApi.NewAPI(url)
	statedb = eth.NewImplStateDB(url)
	eu := ccEu.NewEU(config.ChainConfig, *config.VMConfig, config.Chain, statedb, api, url)

	config.Coinbase = &tests.Coinbase
	config.BlockNumber = new(big.Int).SetUint64(10000000)
	config.Time = new(big.Int).SetUint64(10000000)

	// ================================== Compile the contract ==================================
	currentPath, _ := os.Getwd()
	compiler := filepath.Dir(filepath.Dir(filepath.Dir(currentPath))) + "/tests/compiler.py"
	code, err := tests.CompileContracts(compiler, "./u256_test.sol", "CumulativeU256Test")
	if err != nil || len(code) == 0 {
		t.Error("Error: Failed to generate the byte code")
	}

	// ================================== Deploy the contract ==================================
	msg := types.NewMessage(tests.User1, nil, 0, new(big.Int).SetUint64(0), 1e15, new(big.Int).SetUint64(1), evmcommon.Hex2Bytes(code), nil, true)        // Build the message
	_, transitions, receipt, err := eu.Run(evmcommon.BytesToHash([]byte{1, 1, 1}), 1, &msg, ccEu.NewEVMBlockContextV2(config), ccEu.NewEVMTxContext(msg)) // Execute it
	// ---------------

	// t.Log("\n" + FormatTransitions(accesses))
	t.Log("\n" + tests.FormatTransitions(transitions))
	t.Log(receipt)
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
	// eu = ccEu.NewEU(config.ChainConfig, *config.VMConfig, config.Chain, statedb, api, url)

	// config.BlockNumber = new(big.Int).SetUint64(10000001)
	// config.Time = new(big.Int).SetUint64(10000001)

	// data := crypto.Keccak256([]byte("length()"))[:4]
	// data = append(data, evmcommon.BytesToHash(tests.User1.Bytes()).Bytes()...)
	// data = append(data, evmcommon.BytesToHash([]byte{0xcc}).Bytes()...)
	// msg = types.NewMessage(tests.User1, &contractAddress, 1, new(big.Int).SetUint64(0), 1e15, new(big.Int).SetUint64(1), data, nil, true)
	// _, transitions, receipt, err = eu.Run(evmcommon.BytesToHash([]byte{2, 2, 2}), 2, &msg, ccEu.NewEVMBlockContextV2(config), ccEu.NewEVMTxContext(msg))
	// t.Log("\n" + tests.FormatTransitions(transitions))
	// t.Log(receipt)
	// if receipt.Status != 1 {
	// 	t.Error("Error: Failed to calll length()!!!", err)
	// }

}
