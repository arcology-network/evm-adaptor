package tests

import (
	"fmt"
	"math"
	"math/big"

	"github.com/arcology-network/common-lib/cachedstorage"
	concurrenturl "github.com/arcology-network/concurrenturl"
	"github.com/arcology-network/concurrenturl/commutative"
	"github.com/arcology-network/concurrenturl/interfaces"
	ccurlstorage "github.com/arcology-network/concurrenturl/storage"
	evmcommon "github.com/arcology-network/evm/common"

	evmtypes "github.com/arcology-network/evm/core/types"
	"github.com/arcology-network/evm/core/vm"
	"github.com/arcology-network/evm/crypto"
	"github.com/arcology-network/evm/params"
	cceu "github.com/arcology-network/vm-adaptor"

	ccapi "github.com/arcology-network/vm-adaptor/api"
	eucommon "github.com/arcology-network/vm-adaptor/common"
	"github.com/arcology-network/vm-adaptor/eth"
	cceueth "github.com/arcology-network/vm-adaptor/eth"
)

func Prepare(db interfaces.Datastore, height uint64, transitions []interfaces.Univalue, txs []uint32) (*cceu.EU, *cceu.Config) {
	url := concurrenturl.NewConcurrentUrl(db)
	if transitions != nil && len(transitions) != 0 {
		url.Import(transitions)
		url.Sort()
		url.Commit(txs)
	}

	api := ccapi.NewAPI(url)
	statedb := cceueth.NewImplStateDB(api)

	config := MainTestConfig()
	config.Coinbase = &eucommon.Coinbase
	config.BlockNumber = new(big.Int).SetUint64(height)
	config.Time = new(big.Int).SetUint64(height)

	return cceu.NewEU(config.ChainConfig, *config.VMConfig, statedb, api), config
}

func Deploy(eu *cceu.EU, config *cceu.Config, owner evmcommon.Address, nonce uint64, code string, args ...[]byte) ([]interfaces.Univalue, *evmtypes.Receipt, error) {
	data := evmcommon.Hex2Bytes(code)
	for _, arg := range args {
		data = append(data, evmcommon.BytesToHash(arg).Bytes()...)
	}
	msg := evmtypes.NewMessage(owner, nil, nonce, new(big.Int).SetUint64(0), 1e15, new(big.Int).SetUint64(1), data, nil, false)
	receipt, _, err := eu.Run(evmcommon.BytesToHash([]byte{byte(nonce + 1), byte(nonce + 1), byte(nonce + 1)}), int(nonce+1), &msg, cceu.NewEVMBlockContext(config), cceu.NewEVMTxContext(msg))
	_, transitions := eu.Api().Ccurl().ExportAll()

	return transitions, receipt, err
}

func CallFunc(eu *cceu.EU, config *cceu.Config, from, to *evmcommon.Address, nonce uint64, checkNonce bool, function string, encodedArgs ...[]byte) ([]interfaces.Univalue, []interfaces.Univalue, *evmtypes.Receipt, error) {
	data := crypto.Keccak256([]byte(function))[:4]
	for _, arg := range encodedArgs {
		data = append(data, arg...)
	}
	msg := evmtypes.NewMessage(*from, to, nonce, new(big.Int).SetUint64(0), 1e15, new(big.Int).SetUint64(1), data, nil, checkNonce)
	receipt, _, err := eu.Run(evmcommon.BytesToHash([]byte{byte((nonce + 1) / 65536), byte((nonce + 1) / 256), byte((nonce + 1) % 256)}), int(nonce+1), &msg, cceu.NewEVMBlockContext(config), cceu.NewEVMTxContext(msg))
	accesses, transitions := eu.Api().Ccurl().ExportAll()
	return accesses, transitions, receipt, err
}

func PrintInput(input []byte) {
	fmt.Println(input)
	fmt.Println()
	fmt.Println(input[:4])
	input = input[4:]
	for i := int(0); i < len(input)/32; i++ {
		fmt.Println(input[i*32 : (i+1)*32])
	}
	fmt.Println()
}

func MainTestConfig() *cceu.Config {
	vmConfig := vm.Config{}
	cfg := &cceu.Config{
		ChainConfig: params.MainnetChainConfig,
		VMConfig:    &vmConfig,
		BlockNumber: big.NewInt(0),
		ParentHash:  evmcommon.Hash{},
		Time:        big.NewInt(0),
		Coinbase:    &eucommon.Coinbase,
		GasLimit:    math.MaxUint64, // Should come from the message
		Difficulty:  big.NewInt(0),
	}
	cfg.Chain = new(cceu.DummyChain)
	return cfg
}

func NewTestEU() (*cceu.EU, *cceu.Config, interfaces.Datastore, *concurrenturl.ConcurrentUrl, []interfaces.Univalue) {
	persistentDB := cachedstorage.NewDataStore()
	persistentDB.Inject((&concurrenturl.Platform{}).Eth10Account(), commutative.NewPath())
	db := ccurlstorage.NewTransientDB(persistentDB)

	url := concurrenturl.NewConcurrentUrl(db)
	api := ccapi.NewAPI(url)

	statedb := eth.NewImplStateDB(api)
	statedb.Prepare(evmcommon.Hash{}, evmcommon.Hash{}, 0)
	statedb.CreateAccount(eucommon.Coinbase)
	statedb.CreateAccount(eucommon.Alice)
	statedb.AddBalance(eucommon.Alice, new(big.Int).SetUint64(1e18))
	// transitions := url.Export()
	_, transitions := url.ExportAll()
	// fmt.Println("\n" + eucommon.FormatTransitions(transitions))

	// Deploy.
	url = concurrenturl.NewConcurrentUrl(db)
	url.Import(transitions)
	url.Sort()
	url.Commit([]uint32{0})
	api = ccapi.NewAPI(url)
	statedb = eth.NewImplStateDB(api)

	config := MainTestConfig()
	config.Coinbase = &eucommon.Coinbase
	config.BlockNumber = new(big.Int).SetUint64(10000000)
	config.Time = new(big.Int).SetUint64(10000000)

	return cceu.NewEU(config.ChainConfig, *config.VMConfig, statedb, api), config, db, url, transitions
}

// func NewTestEUwithUrl(db interfaces.Datastore, ccurl *concurrenturl.ConcurrentUrl) (*cceu.EU, *cceu.Config, interfaces.Datastore, *concurrenturl.ConcurrentUrl) {
// 	// Deploy.
// 	ccurl = concurrenturl.NewConcurrentUrl(db)
// 	ccurl.Import(transitions)
// 	ccurl.PostImport()
// 	ccurl.Commit([]uint32{0})
// 	api := ccapi.NewAPI(url)
// 	statedb = eth.NewImplStateDB(url)

// 	config := MainTestConfig()
// 	config.Coinbase = &eucommon.Coinbase
// 	config.BlockNumber = new(big.Int).SetUint64(10000000)
// 	config.Time = new(big.Int).SetUint64(10000000)

// 	return cceu.NewEU(config.ChainConfig, *config.VMConfig, statedb, api), config, db, url
// }
