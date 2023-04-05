package basic

import (
	"math/big"
	"os/exec"
	"testing"

	cachedstorage "github.com/arcology-network/common-lib/cachedstorage"
	"github.com/arcology-network/concurrenturl/v2"
	urlcommon "github.com/arcology-network/concurrenturl/v2/common"
	curstorage "github.com/arcology-network/concurrenturl/v2/storage"
	"github.com/arcology-network/concurrenturl/v2/type/commutative"
	"github.com/arcology-network/evm/common"
	"github.com/arcology-network/evm/core/types"
	"github.com/arcology-network/evm/crypto"
	ccEu "github.com/arcology-network/vm-adaptor"
	ccApi "github.com/arcology-network/vm-adaptor/api"
	eth "github.com/arcology-network/vm-adaptor/eth"
	tests "github.com/arcology-network/vm-adaptor/tests"
)

func TestApiInterfaces(t *testing.T) {
	config := tests.MainConfig()
	persistentDB := cachedstorage.NewDataStore()
	meta, _ := commutative.NewMeta(urlcommon.NewPlatform().Eth10Account())
	persistentDB.Inject(urlcommon.NewPlatform().Eth10Account(), meta)
	db := curstorage.NewTransientDB(persistentDB)

	url := concurrenturl.NewConcurrentUrl(db)
	//	api := ccApi.NewAPI(url)
	statedb := eth.NewImplStateDB(url)
	statedb.Prepare(common.Hash{}, common.Hash{}, 0)
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
	eu := ccEu.NewEU(config.ChainConfig, *config.VMConfig, config.Chain, statedb, api, db, url)

	config.Coinbase = &tests.Coinbase
	config.BlockNumber = new(big.Int).SetUint64(10000000)
	config.Time = new(big.Int).SetUint64(10000000)

	// ================================== Compile the contract ==================================
	_, err := exec.Command("python", "./compiler.py").Output() // capture the output of the Python script
	if err != nil {
		panic(err)
	}

	code, err := tests.BytecodeReader("./bytecode.txt") // Read the byte code
	if err != nil {
		t.Error("Error: ", err)
	}

	// ================================== Deploy the contract ==================================
	msg := types.NewMessage(tests.User1, nil, 0, new(big.Int).SetUint64(0), 1e15, new(big.Int).SetUint64(1), common.Hex2Bytes(code), nil, true)        // Build the message
	_, transitions, receipt, err := eu.Run(common.BytesToHash([]byte{1, 1, 1}), 1, &msg, ccEu.NewEVMBlockContextV2(config), ccEu.NewEVMTxContext(msg)) // Execute it
	// ---------------

	// t.Log("\n" + FormatTransitions(accesses))
	t.Log("\n" + tests.FormatTransitions(transitions))
	// t.Log(receipt)
	contractAddress := receipt.ContractAddress
	if receipt.Status != 1 || err != nil {
		t.Error("Error: Deployment failed!!!", err)
	}
	// return
	// ================================== Call length() ==================================
	url = concurrenturl.NewConcurrentUrl(db)
	url.Import(transitions)
	url.PostImport()
	errs := url.Commit([]uint32{1})
	if len(errs) != 0 {
		t.Error(errs)
		return
	}
	api = ccApi.NewAPI(url)
	statedb = eth.NewImplStateDB(url)
	eu = ccEu.NewEU(config.ChainConfig, *config.VMConfig, config.Chain, statedb, api, db, url)

	config.BlockNumber = new(big.Int).SetUint64(10000001)
	config.Time = new(big.Int).SetUint64(10000001)

	data := crypto.Keccak256([]byte("length()"))[:4]
	data = append(data, common.BytesToHash(tests.User1.Bytes()).Bytes()...)
	data = append(data, common.BytesToHash([]byte{0xcc}).Bytes()...)
	msg = types.NewMessage(tests.User1, &contractAddress, 1, new(big.Int).SetUint64(0), 1e15, new(big.Int).SetUint64(1), data, nil, true)
	_, transitions, receipt, err = eu.Run(common.BytesToHash([]byte{2, 2, 2}), 2, &msg, ccEu.NewEVMBlockContextV2(config), ccEu.NewEVMTxContext(msg))
	t.Log("\n" + tests.FormatTransitions(transitions))
	t.Log(receipt)
	if receipt.Status != 1 {
		t.Error("Error: Set failed!!!", err)
	}

	// Get.
	// url = concurrenturl.NewConcurrentUrl(db)
	// url.Import(transitions)
	// url.PostImport()
	// errs = url.Commit([]uint32{2})
	// if len(errs) != 0 {
	// 	t.Error(errs)
	// 	return
	// }
	// api = ccApi.NewAPI(url)
	// statedb = eth.NewImplStateDB(url)
	// eu = ccEu.NewEU(config.ChainConfig, *config.VMConfig, config.Chain, statedb, api, db, url)

	// config.BlockNumber = new(big.Int).SetUint64(10000002)
	// config.Time = new(big.Int).SetUint64(10000002)

	// data = crypto.Keccak256([]byte("getSum()"))[:4]
	// msg = types.NewMessage(tests.User1, &contractAddress, 2, new(big.Int).SetUint64(0), 1e15, new(big.Int).SetUint64(1), data, nil, true)
	// accesses, transitions, receipt, err := eu.Run(common.BytesToHash([]byte{3, 3, 3}), 3, &msg, ccEu.NewEVMBlockContextV2(config), ccEu.NewEVMTxContext(msg))
	// t.Log("\n" + tests.FormatTransitions(accesses))
	// t.Log("\n" + tests.FormatTransitions(transitions))
	// t.Log(receipt)

	// if receipt.Status != 1 {
	// 	t.Error("Error: Set failed!!!", err)
	// }
}
