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

	ccurlcommon "github.com/arcology-network/concurrenturl/common"
	ccapi "github.com/arcology-network/vm-adaptor/api"
	eucommon "github.com/arcology-network/vm-adaptor/common"
	"github.com/arcology-network/vm-adaptor/compiler"
	"github.com/arcology-network/vm-adaptor/eth"
	"github.com/arcology-network/vm-adaptor/execution"
)

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
	persistentDB.Inject(ccurlcommon.ETH10_ACCOUNT_PREFIX, commutative.NewPath())
	db := ccurlstorage.NewTransientDB(persistentDB)

	url := concurrenturl.NewConcurrentUrl(db)
	api := ccapi.NewAPI(url)

	statedb := eth.NewImplStateDB(api)
	statedb.PrepareFormer(evmcommon.Hash{}, evmcommon.Hash{}, 0)
	statedb.CreateAccount(eucommon.Coinbase)
	statedb.CreateAccount(eucommon.Alice)
	statedb.AddBalance(eucommon.Alice, new(big.Int).SetUint64(1e18))

	statedb.CreateAccount(eucommon.RUNTIME_HANDLER)
	// statedb.AddBalance(eucommon.RUNTIME_HANDLER, new(big.Int).SetUint64(1e18))

	_, transitions := api.StateFilter().ByType()
	// indexer.Univalues(transitionsFiltered).Print()

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

	// _, transitions := eu.Api().Ccurl().ExportAll()
	// indexer.Univalues(transitions).Print()
	// fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>  Ignore addresses Removed  <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<")
	// fmt.Println()
	_, transitionsFiltered := eu.Api().StateFilter().ByType()

	eu.Api().Ccurl().Import(transitionsFiltered)
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
	// _, transitions := eu.Api().StateFilter().ByType()

	// msg = core.NewMessage(eucommon.Alice, &contractAddress, 1, new(big.Int).SetUint64(0), 1e15, new(big.Int).SetUint64(1), data, nil, false)
	// receipt, execResult, _ := eu.Run(evmcommon.BytesToHash([]byte{1, 1, 1}), 1, &msg, execution.NewEVMBlockContext(config), execution.NewEVMTxContext(msg))
	// _, transitions = eu.Api().StateFilter().ByType()

	if receipt.Status != 1 {
		return execResult.Err, eu
	}

	if execResult != nil && execResult.Err != nil {
		return execResult.Err, eu
	}
	return nil, eu
}
