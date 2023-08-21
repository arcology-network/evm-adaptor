package tests

import (
	"math/big"

	"github.com/arcology-network/common-lib/cachedstorage"
	concurrenturl "github.com/arcology-network/concurrenturl"
	ccurlcommon "github.com/arcology-network/concurrenturl/common"
	"github.com/arcology-network/concurrenturl/commutative"
	ccurlstorage "github.com/arcology-network/concurrenturl/storage"
	evmcommon "github.com/arcology-network/evm/common"
	evmcore "github.com/arcology-network/evm/core"
	evmcoretypes "github.com/arcology-network/evm/core/types"
	ccapi "github.com/arcology-network/vm-adaptor/api"
	eucommon "github.com/arcology-network/vm-adaptor/common"
	"github.com/arcology-network/vm-adaptor/compiler"
	"github.com/arcology-network/vm-adaptor/eth"
	"github.com/arcology-network/vm-adaptor/execution"
)

type Contract struct {
	path    string
	name    string
	code    string
	eu      *execution.EU
	owner   [20]byte
	address evmcommon.Address

	execResult *evmcore.ExecutionResult
	receipt    *evmcoretypes.Receipt
	err        error
}

func NewContract(accounts []evmcommon.Address, owner [20]byte, targetPath, file, version, contractName string) (*Contract, error) {
	persistentDB := cachedstorage.NewDataStore()
	persistentDB.Inject(ccurlcommon.ETH10_ACCOUNT_PREFIX, commutative.NewPath())
	db := ccurlstorage.NewTransientDB(persistentDB)

	url := concurrenturl.NewConcurrentUrl(db)
	api := ccapi.NewAPI(url)

	statedb := eth.NewImplStateDB(api)
	statedb.PrepareFormer(evmcommon.Hash{}, evmcommon.Hash{}, 0)

	statedb.CreateAccount(eucommon.Coinbase)
	for i := 0; i < len(accounts); i++ {
		statedb.CreateAccount(accounts[i])
		statedb.AddBalance(accounts[i], new(big.Int).SetUint64(1e18))
	}

	_, transitions := api.StateFilter().ByType()
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

	if code, err := compiler.CompileContracts(targetPath, file, version, contractName, false); err == nil {
		return &Contract{
			name:  contractName,
			owner: owner,
			eu:    execution.NewEU(config.ChainConfig, *config.VMConfig, statedb, api),
			path:  targetPath + file,
			code:  code,
		}, nil
	} else {
		return nil, err
	}
}

// func (this *Contract) Deploy(msgSender [20]byte, nonce uint64) ([]byte, error) {
// 	data, err := this.Invoke(msgSender, evmcommon.Hex2Bytes(this.code))
// 	this.address = this.receipt.ContractAddress
// 	this.Apply()
// 	return data, err
// }

// func (this *Contract) Invoke(msgSender [20]byte, funcCall []byte) ([]byte, error) {
// 	msg := core.NewMessage(msgSender, &this.address, 10+1, new(big.Int).SetUint64(0), 1e15, new(big.Int).SetUint64(1), funcCall, nil, false)
// 	stdMsg := &execution.StandardMessage{
// 		ID:     1,
// 		TxHash: [32]byte{1, 1, 1},
// 		Native: &msg, // Build the message
// 		Source: commontypes.TX_SOURCE_LOCAL,
// 	}

// 	config := MainTestConfig()
// 	config.Coinbase = &eucommon.Coinbase
// 	config.BlockNumber = new(big.Int).SetUint64(10000000)
// 	config.Time = new(big.Int).SetUint64(10000000)

// 	this.receipt, this.execResult, this.err = this.eu.Run(stdMsg, execution.NewEVMBlockContext(config), execution.NewEVMTxContext(*stdMsg.Native)) // Execute it
// 	if this.err != nil || this.receipt.Status != 1 {
// 		return nil, errors.New("Error: Deployment failed!!!")
// 	}
// 	return this.execResult.ReturnData, nil
// }

// func (this *Contract) Apply() {
// 	_, transitionsFiltered := this.eu.Api().StateFilter().ByType()
// 	this.eu.Api().Ccurl().Import(transitionsFiltered)
// 	this.eu.Api().Ccurl().Sort()
// 	this.eu.Api().Ccurl().Commit([]uint32{1})
// }
