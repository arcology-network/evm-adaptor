package eu

import (
	"bytes"
	"fmt"
	"math"
	"math/big"

	ethCommon "github.com/arcology-network/evm/common"
	"github.com/arcology-network/evm/core"
	evmcore "github.com/arcology-network/evm/core"
	"github.com/arcology-network/evm/core/types"
	evmcoretypes "github.com/arcology-network/evm/core/types"
	"github.com/arcology-network/evm/core/vm"
	corevm "github.com/arcology-network/evm/core/vm"
	"github.com/arcology-network/evm/crypto"
	"github.com/arcology-network/evm/params"
	eucommon "github.com/arcology-network/vm-adaptor/common"
	eth "github.com/arcology-network/vm-adaptor/eth"
)

type EU struct {
	msg         *evmcoretypes.Message
	evm         *vm.EVM               // Original ETH EVM
	statedb     vm.StateDB            // Arcology Implementation of Eth StateDB
	api         eucommon.EthApiRouter // Arcology API calls
	CallContext *corevm.ScopeContext  // Arcology API calls
}

func NewEU(chainConfig *params.ChainConfig, vmConfig vm.Config, statedb vm.StateDB, api eucommon.EthApiRouter) *EU {
	eu := &EU{
		evm:     vm.NewEVM(vm.BlockContext{BlockNumber: new(big.Int).SetUint64(100000000)}, vm.TxContext{}, statedb, chainConfig, vmConfig),
		statedb: statedb,
		api:     api,
	}

	eu.api.SetEU(eu)
	eu.evm.ArcologyNetworkAPIs.APIs = (api)
	return eu
}

func (this *EU) Message() *evmcoretypes.Message { return this.msg }
func (this *EU) VM() *vm.EVM                    { return this.evm }
func (this *EU) Statedb() vm.StateDB            { return this.statedb }
func (this *EU) Api() eucommon.EthApiRouter     { return this.api }

// func (this *EU) Depth() uint8                               { return this.depth }

func (this *EU) SetContext(statedb vm.StateDB, api eucommon.EthApiRouter) {
	this.api = api
	this.statedb = statedb

	this.evm.StateDB = this.statedb
	this.evm.ArcologyNetworkAPIs.APIs = (api)
}

func (this *EU) Run(txHash ethCommon.Hash, txIndex int, msg *evmcoretypes.Message, blockContext vm.BlockContext, txContext vm.TxContext) (*types.Receipt, *evmcore.ExecutionResult, error) {
	this.statedb.(*eth.ImplStateDB).Prepare(txHash, ethCommon.Hash{}, txIndex)
	this.api.SetContext(txHash, blockContext.BlockNumber, uint32(txIndex))
	this.evm.Context = blockContext
	this.evm.TxContext = txContext
	this.msg = msg

	gasPool := core.GasPool(math.MaxUint64)

	result, err := core.ApplyMessage(this.evm, *msg, &gasPool) // Execute the transcation
	if err != nil {
		return nil, nil, err // Failed in Precheck before tx execution started
	}

	assertLog := GetAssertion(result.ReturnData) // Get the assertion
	if len(assertLog) > 0 {
		this.api.AddLog("assert", assertLog)
	}

	// Create a new receipt
	receipt := types.NewReceipt(nil, result.Failed(), result.UsedGas)
	receipt.TxHash = txHash
	receipt.GasUsed = result.UsedGas

	// Check the newly created address
	if msg.To() == nil {
		userSpecifiedAddress := crypto.CreateAddress(this.evm.Origin, msg.Nonce())
		receipt.ContractAddress = result.ContractAddress
		if !bytes.Equal(userSpecifiedAddress.Bytes(), result.ContractAddress.Bytes()) {
			this.api.AddLog("ContractAddressWarning", fmt.Sprintf("user specified address = %v, inner address = %v", userSpecifiedAddress, result.ContractAddress))
		}
	}
	receipt.Logs = this.statedb.(*eth.ImplStateDB).GetLogs(txHash)
	receipt.Bloom = types.CreateBloom(types.Receipts{receipt})
	// accesses, transitions := this.api.Ccurl().ExportAll()

	// if result.Failed() { // Failed
	// 	accesses = accesses[:0]
	// 	common.RemoveIf(&transitions, func(val interfaces.Univalue) bool {
	// 		path := val.GetPath()
	// 		return len(*path) <= 5 || (*path)[len(*path)-5:] != "nonce" // Keep nonce transitions only, nonce needs to increment anyway.
	// 	})
	// }

	return receipt, result, err
}

// Get the assertion info from the execution result
func GetAssertion(ret []byte) string {
	offset := 4 + 32 + 32
	pattern := []byte{8, 195, 121, 160}
	if ret != nil && len(ret) > offset {
		starts := ret[:4]
		if string(pattern) == string(starts) {
			remains := ret[offset:]
			end := 0
			for i := range remains {
				if remains[i] == 0 {
					end = i
					break
				}
			}
			return string(remains[:end])
		}
	}
	return ""
}
