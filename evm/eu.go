package evm

import (
	"bytes"
	"fmt"
	"math"
	"math/big"

	common "github.com/arcology-network/common-lib/common"
	"github.com/arcology-network/concurrenturl/v2"
	urlcommon "github.com/arcology-network/concurrenturl/v2/common"
	ethCommon "github.com/arcology-network/evm/common"
	"github.com/arcology-network/evm/core"
	"github.com/arcology-network/evm/core/types"
	"github.com/arcology-network/evm/core/vm"
	"github.com/arcology-network/evm/crypto"
	"github.com/arcology-network/evm/params"
)

type EU struct {
	evm     *vm.EVM // Original ETH EVM
	statedb vm.StateDB
	api     *API // Arcology APIs
	db      urlcommon.DatastoreInterface
	url     *concurrenturl.ConcurrentUrl
	uid     uint64 // Number of UIDs genearted
}

func NewEU(chainConfig *params.ChainConfig, vmConfig vm.Config, chainContext core.ChainContext, statedb vm.StateDB, api *API, db urlcommon.DatastoreInterface, url *concurrenturl.ConcurrentUrl) *EU {
	eu := &EU{
		evm:     vm.NewEVM(vm.BlockContext{BlockNumber: new(big.Int).SetUint64(100000000)}, vm.TxContext{}, statedb, chainConfig, vmConfig),
		statedb: statedb,
		api:     api,
		db:      db,
		url:     url,
	}

	eu.evm.SetApi(api)
	return eu
}

func (eu *EU) SetContext(statedb vm.StateDB, api *API, db urlcommon.DatastoreInterface, url *concurrenturl.ConcurrentUrl) {
	eu.api = api
	eu.statedb = statedb
	eu.db = db
	eu.url = url

	eu.evm.StateDB = eu.statedb
	eu.evm.SetApi(api)
}

func (eu *EU) Run(thash ethCommon.Hash, tindex int, msg *types.Message, blockContext vm.BlockContext, txContext vm.TxContext) ([]urlcommon.UnivalueInterface, []urlcommon.UnivalueInterface, *types.Receipt) {
	eu.statedb.(*ethState).Prepare(thash, ethCommon.Hash{}, tindex)
	eu.api.Prepare(thash, blockContext.BlockNumber, uint32(tindex))
	eu.evm.Context = blockContext
	eu.evm.TxContext = txContext

	gasPool := core.GasPool(math.MaxUint64)

	result, err := core.ApplyMessage(eu.evm, *msg, &gasPool) // Execute the transcation
	if err != nil {
		result = &core.ExecutionResult{
			Err: err,
		}
	}

	assertLog := GetAssert(result.ReturnData)
	if len(assertLog) > 0 {
		eu.api.AddLog("assert", assertLog)
	}

	// Create a new receipt
	receipt := types.NewReceipt(nil, result.Failed(), result.UsedGas)
	receipt.TxHash = thash
	receipt.GasUsed = result.UsedGas

	// Check the newly created address
	if msg.To() == nil {
		userSpecifiedAddress := crypto.CreateAddress(eu.evm.Origin, msg.Nonce())
		receipt.ContractAddress = result.ContractAddress
		if !bytes.Equal(userSpecifiedAddress.Bytes(), result.ContractAddress.Bytes()) {
			eu.api.AddLog("ContractAddressWarning", fmt.Sprintf("user specified address = %v, inner address = %v", userSpecifiedAddress, result.ContractAddress))
		}
	}
	receipt.Logs = eu.statedb.(*ethState).GetLogs(thash)
	receipt.Bloom = types.CreateBloom(types.Receipts{receipt})
	accesses, transitions := eu.url.Export(false)

	if result.Failed() { // Failed
		accesses = accesses[:0]
		common.RemoveIf(&transitions, func(val urlcommon.UnivalueInterface) bool {
			path := val.GetPath()
			return len(*path) <= 5 || (*path)[len(*path)-5:] != "nonce" // Remove all except the nonce transitions, nonce needs to increment anyway.
		})
	}

	return accesses, transitions, receipt
}

func GetAssert(ret []byte) string {
	startIdx := 4 + 32 + 32
	pattern := []byte{8, 195, 121, 160}
	if ret != nil || len(ret) > startIdx {
		starts := ret[:4]
		if string(pattern) == string(starts) {
			remains := ret[startIdx:]
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
