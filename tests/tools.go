package tests

import (
	"errors"
	"math"
	"math/big"

	"github.com/arcology-network/common-lib/cachedstorage"
	commontypes "github.com/arcology-network/common-lib/types"
	concurrenturl "github.com/arcology-network/concurrenturl"
	"github.com/arcology-network/concurrenturl/commutative"
	"github.com/arcology-network/concurrenturl/interfaces"
	ccurlstorage "github.com/arcology-network/concurrenturl/storage"
	evmcommon "github.com/arcology-network/evm/common"
	"github.com/arcology-network/evm/core"
	evmcore "github.com/arcology-network/evm/core"

	"github.com/arcology-network/evm/core/vm"
	"github.com/arcology-network/evm/crypto"
	"github.com/arcology-network/evm/params"

	ccapi "github.com/arcology-network/vm-adaptor/api"
	eucommon "github.com/arcology-network/vm-adaptor/common"
	"github.com/arcology-network/vm-adaptor/compiler"
	"github.com/arcology-network/vm-adaptor/eth"
	"github.com/arcology-network/vm-adaptor/execution"
)

// func Prepare(db interfaces.Datastore, height uint64, transitions []interfaces.Univalue, txs []uint32) (*execution.EU, *execution.Config) {
// 	url := concurrenturl.NewConcurrentUrl(db)
// 	if transitions != nil && len(transitions) != 0 {
// 		url.Import(transitions)
// 		url.Sort()
// 		url.Commit(txs)
// 	}

// 	api := ccapi.NewAPI(url)
// 	statedb := eth.NewImplStateDB(api)

// 	config := MainTestConfig()
// 	config.Coinbase = &eucommon.Coinbase
// 	config.BlockNumber = new(big.Int).SetUint64(height)
// 	config.Time = new(big.Int).SetUint64(height)

// 	return execution.NewEU(config.ChainConfig, *config.VMConfig, statedb, api), config
// }

// func Deploy(eu *execution.EU, config *execution.Config, owner evmcommon.Address, nonce uint64, code string, args ...[]byte) ([]interfaces.Univalue, *evmtypes.Receipt, error) {
// 	data := evmcommon.Hex2Bytes(code)
// 	for _, arg := range args {
// 		data = append(data, evmcommon.BytesToHash(arg).Bytes()...)
// 	}
// 	msg := core.NewMessage(owner, nil, nonce, new(big.Int).SetUint64(0), 1e15, new(big.Int).SetUint64(1), data, nil, false)
// 	receipt, _, err := eu.Run(evmcommon.BytesToHash([]byte{byte(nonce + 1), byte(nonce + 1), byte(nonce + 1)}), uint32(nonce+1), &msg, execution.NewEVMBlockContext(config), execution.NewEVMTxContext(msg))
// 	_, transitions := eu.Api().Ccurl().ExportAll()

// 	return transitions, receipt, err
// }

// func CallFunc(eu *execution.EU, config *execution.Config, from, to *evmcommon.Address, nonce uint64, checkNonce bool, function string, encodedArgs ...[]byte) ([]interfaces.Univalue, []interfaces.Univalue, *evmtypes.Receipt, error) {
// 	data := crypto.Keccak256([]byte(function))[:4]
// 	for _, arg := range encodedArgs {
// 		data = append(data, arg...)
// 	}
// 	msg := core.NewMessage(*from, to, nonce, new(big.Int).SetUint64(0), 1e15, new(big.Int).SetUint64(1), data, nil, checkNonce)
// 	receipt, _, err := eu.Run(evmcommon.BytesToHash([]byte{byte((nonce + 1) / 65536), byte((nonce + 1) / 256), byte((nonce + 1) % 256)}), uint32(nonce+1), &msg, execution.NewEVMBlockContext(config), execution.NewEVMTxContext(msg))
// 	accesses, transitions := eu.Api().Ccurl().ExportAll()
// 	return accesses, transitions, receipt, err
// }

// func PrintInput(input []byte) {
// 	fmt.Println(input)
// 	fmt.Println()
// 	fmt.Println(input[:4])
// 	input = input[4:]
// 	for i := int(0); i < len(input)/32; i++ {
// 		fmt.Println(input[i*32 : (i+1)*32])
// 	}
// 	fmt.Println()
// }

// func Run(eu *execution.EU, config *execution.Config, from, to *evmcommon.Address, nonce uint64, checkNonce bool, function string, args ...[]byte) ([]interfaces.Univalue, *evmtypes.Receipt) {
// 	data := crypto.Keccak256([]byte(function))[:4]
// 	for _, arg := range args {
// 		data = append(data, evmcommon.BytesToHash(arg).Bytes()...)
// 	}
// 	msg := core.NewMessage(*from, to, nonce, new(big.Int).SetUint64(0), 1e15, new(big.Int).SetUint64(1), data, nil, checkNonce)
// 	receipt, _, _ := eu.Run(evmcommon.BytesToHash([]byte{byte((nonce + 1) / 65536), byte((nonce + 1) / 256), byte((nonce + 1) % 256)}), uint32(nonce+1), &msg, execution.NewEVMBlockContext(config), execution.NewEVMTxContext(msg))
// 	_, transitions := eu.Api().Ccurl().ExportAll()

// 	return transitions, receipt
// }

// func RunEx(eu *execution.EU, config *execution.Config, from, to *evmcommon.Address, nonce uint64, checkNonce bool, function string, args ...[]byte) ([]interfaces.Univalue, []interfaces.Univalue, *evmtypes.Receipt) {
// 	data := crypto.Keccak256([]byte(function))[:4]
// 	for _, arg := range args {
// 		data = append(data, evmcommon.BytesToHash(arg).Bytes()...)
// 	}
// 	msg := core.NewMessage(*from, to, nonce, new(big.Int).SetUint64(0), 1e15, new(big.Int).SetUint64(1), data, nil, checkNonce)
// 	receipt, _, _ := eu.Run(evmcommon.BytesToHash([]byte{byte((nonce + 1) / 65536), byte((nonce + 1) / 256), byte((nonce + 1) % 256)}), uint32(nonce+1), &msg, execution.NewEVMBlockContext(config), execution.NewEVMTxContext(msg))
// 	accesses, transitions := eu.Api().Ccurl().ExportAll()

// 	return accesses, transitions, receipt
// }

func MainTestConfig() *execution.Config {
	vmConfig := vm.Config{}
	cfg := &execution.Config{
		ChainConfig: params.MainnetChainConfig,
		VMConfig:    &vmConfig,
		BlockNumber: big.NewInt(0),
		ParentHash:  evmcommon.Hash{},
		Time:        big.NewInt(0),
		Coinbase:    &eucommon.Coinbase,
		GasLimit:    math.MaxUint64, // Should come from the message
		Difficulty:  big.NewInt(0),
	}
	cfg.Chain = new(execution.DummyChain)
	return cfg
}

func NewTestEU() (*execution.EU, *execution.Config, interfaces.Datastore, *concurrenturl.ConcurrentUrl, []interfaces.Univalue) {
	persistentDB := cachedstorage.NewDataStore()
	persistentDB.Inject((&concurrenturl.Platform{}).Eth10Account(), commutative.NewPath())
	db := ccurlstorage.NewTransientDB(persistentDB)

	url := concurrenturl.NewConcurrentUrl(db)
	api := ccapi.NewAPI(url)

	statedb := eth.NewImplStateDB(api)
	statedb.PrepareFormer(evmcommon.Hash{}, evmcommon.Hash{}, 0)
	statedb.CreateAccount(eucommon.Coinbase)
	statedb.CreateAccount(eucommon.Alice)
	statedb.AddBalance(eucommon.Alice, new(big.Int).SetUint64(1e18))

	statedb.CreateAccount(eucommon.ATOMIC_HANDLER)
	// statedb.AddBalance(eucommon.ATOMIC_HANDLER, new(big.Int).SetUint64(1e18))

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

	return execution.NewEU(config.ChainConfig, *config.VMConfig, statedb, api), config, db, url, transitions
}

func InvokeTestContract(targetPath, file, version, contractName, funcName string, inputData []byte, checkNonce bool) (error, *execution.EU) {
	code, err := compiler.CompileContracts(targetPath, file, version, contractName, false)

	eu, config, _, _, _ := NewTestEU()
	if err != nil || len(code) == 0 {
		return err, eu
	}

	// ================================== Deploy the contract ==================================
	msg := core.NewMessage(eucommon.Alice, nil, 0, new(big.Int).SetUint64(0), 1e15, new(big.Int).SetUint64(1), evmcommon.Hex2Bytes(code), nil, false)
	stdMsg := &execution.StandardMessage{
		ID:     1,
		TxHash: [32]byte{1, 1, 1},
		Native: &msg, // Build the message
		Source: commontypes.TX_SOURCE_LOCAL,
	}

	receipt, _, err := eu.Run(stdMsg, execution.NewEVMBlockContext(config), execution.NewEVMTxContext(*stdMsg.Native)) // Execute it
	_, transitions := eu.Api().Ccurl().ExportAll()

	// msg := core.NewMessage(eucommon.Alice, nil, 0, new(big.Int).SetUint64(0), 1e15, new(big.Int).SetUint64(1), evmcommon.Hex2Bytes(code), nil, checkNonce) // Build the message
	// receipt, _, err := eu.Run(evmcommon.BytesToHash([]byte{1, 1, 1}), 1, &msg, execution.NewEVMBlockContext(config), execution.NewEVMTxContext(msg))       // Execute it
	// _, transitions := eu.Api().Ccurl().ExportAll()

	eu.Api().Ccurl().Import(transitions)
	eu.Api().Ccurl().Sort()
	eu.Api().Ccurl().Commit([]uint32{1})

	// ---------------
	contractAddress := receipt.ContractAddress
	if receipt.Status != 1 || err != nil {
		errmsg := ""
		if err != nil {
			errmsg = err.Error()
		}
		return errors.New("Error: Deployment failed!!!" + errmsg), eu
	}

	if len(funcName) == 0 {
		return err, eu
	}

	// ================================== CallBasic() ==================================
	// eu, config, _, _, _ = NewTestEU()
	data := crypto.Keccak256([]byte(funcName))[:4]
	data = append(data, inputData...)

	msg = core.NewMessage(eucommon.Alice, &contractAddress, 10, new(big.Int).SetUint64(0), 1e15, new(big.Int).SetUint64(1), data, nil, false)
	stdMsg = &execution.StandardMessage{
		ID:     1,
		TxHash: [32]byte{1, 1, 1},
		Native: &msg, // Build the message
		Source: commontypes.TX_SOURCE_LOCAL,
	}

	var execResult *evmcore.ExecutionResult
	receipt, execResult, err = eu.Run(stdMsg, execution.NewEVMBlockContext(config), execution.NewEVMTxContext(*stdMsg.Native)) // Execute it
	// _, transitions := eu.Api().Ccurl().ExportAll()

	// msg = core.NewMessage(eucommon.Alice, &contractAddress, 1, new(big.Int).SetUint64(0), 1e15, new(big.Int).SetUint64(1), data, nil, false)
	// receipt, execResult, _ := eu.Run(evmcommon.BytesToHash([]byte{1, 1, 1}), 1, &msg, execution.NewEVMBlockContext(config), execution.NewEVMTxContext(msg))
	// _, transitions = eu.Api().Ccurl().ExportAll()

	if receipt.Status != 1 {
		return execResult.Err, eu
	}

	if execResult != nil && execResult.Err != nil {
		return execResult.Err, eu
	}
	return nil, eu
}
