package eu

import (
	"bytes"
	"fmt"
	"math"
	"math/big"

	common "github.com/arcology-network/common-lib/common"
	urlcommon "github.com/arcology-network/concurrenturl/v2/common"
	ethCommon "github.com/arcology-network/evm/common"
	"github.com/arcology-network/evm/core"
	"github.com/arcology-network/evm/core/types"
	"github.com/arcology-network/evm/core/vm"
	"github.com/arcology-network/evm/crypto"
	"github.com/arcology-network/evm/params"

	eucommon "github.com/arcology-network/vm-adaptor/common"
	eth "github.com/arcology-network/vm-adaptor/eth"
)

type EU struct {
	evm     *vm.EVM                               // Original ETH EVM
	statedb vm.StateDB                            // Arcology Implementation of Eth StateDB
	api     eucommon.ConcurrentApiRouterInterface // Arcology API calls
}

func NewEU(chainConfig *params.ChainConfig, vmConfig vm.Config, statedb vm.StateDB, api eucommon.ConcurrentApiRouterInterface) *EU {
	eu := &EU{
		evm:     vm.NewEVM(vm.BlockContext{BlockNumber: new(big.Int).SetUint64(100000000)}, vm.TxContext{}, statedb, chainConfig, vmConfig),
		statedb: statedb,
		api:     api,
	}

	eu.api.SetVM(eu.evm)
	eu.evm.SetApi(api)
	return eu
}

func (eu *EU) SetContext(statedb vm.StateDB, api eucommon.ConcurrentApiRouterInterface) {
	eu.api = api
	eu.statedb = statedb

	eu.evm.StateDB = eu.statedb
	eu.evm.SetApi(api)
}

func (eu *EU) Run(txHash ethCommon.Hash, txIndex int, msg *types.Message, blockContext vm.BlockContext, txContext vm.TxContext) (
	[]urlcommon.UnivalueInterface, []urlcommon.UnivalueInterface, *types.Receipt, *core.ExecutionResult, error) {
	eu.statedb.(*eth.ImplStateDB).Prepare(txHash, ethCommon.Hash{}, txIndex)
	eu.api.Prepare(txHash, blockContext.BlockNumber, uint32(txIndex))
	eu.evm.Context = blockContext
	eu.evm.TxContext = txContext

	gasPool := core.GasPool(math.MaxUint64)

	result, err := core.ApplyMessage(eu.evm, *msg, &gasPool) // Execute the transcation
	if err != nil {
		return []urlcommon.UnivalueInterface{}, []urlcommon.UnivalueInterface{}, nil, nil, err // Failed in Precheck before tx execution started
	}

	assertLog := GetAssertion(result.ReturnData) // Get the assertion
	if len(assertLog) > 0 {
		eu.api.AddLog("assert", assertLog)
	}

	// Create a new receipt
	receipt := types.NewReceipt(nil, result.Failed(), result.UsedGas)
	receipt.TxHash = txHash
	receipt.GasUsed = result.UsedGas

	// Check the newly created address
	if msg.To() == nil {
		userSpecifiedAddress := crypto.CreateAddress(eu.evm.Origin, msg.Nonce())
		receipt.ContractAddress = result.ContractAddress
		if !bytes.Equal(userSpecifiedAddress.Bytes(), result.ContractAddress.Bytes()) {
			eu.api.AddLog("ContractAddressWarning", fmt.Sprintf("user specified address = %v, inner address = %v", userSpecifiedAddress, result.ContractAddress))
		}
	}
	receipt.Logs = eu.statedb.(*eth.ImplStateDB).GetLogs(txHash)
	receipt.Bloom = types.CreateBloom(types.Receipts{receipt})
	accesses, transitions := eu.api.Ccurl().Export(false)

	if result.Failed() { // Failed
		accesses = accesses[:0]
		common.RemoveIf(&transitions, func(val urlcommon.UnivalueInterface) bool {
			path := val.GetPath()
			return len(*path) <= 5 || (*path)[len(*path)-5:] != "nonce" // Keep nonce transitions only, nonce needs to increment anyway.
		})
	}

	return accesses, transitions, receipt, result, err
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
