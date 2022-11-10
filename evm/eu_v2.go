package evm

import (
	"bytes"
	"fmt"
	"math"
	"math/big"

	"github.com/arcology-network/concurrenturl/v2"
	urlcommon "github.com/arcology-network/concurrenturl/v2/common"
	"github.com/arcology-network/evm/common"
	"github.com/arcology-network/evm/core"
	"github.com/arcology-network/evm/core/types"
	"github.com/arcology-network/evm/core/vm"
	"github.com/arcology-network/evm/crypto"
	"github.com/arcology-network/evm/params"
)

type EUV2 struct {
	evm     *vm.EVM
	statedb vm.StateDB
	api     *APIV2
	db      urlcommon.DatastoreInterface
	url     *concurrenturl.ConcurrentUrl
}

func NewEUV2(chainConfig *params.ChainConfig, vmConfig vm.Config, chainContext core.ChainContext, statedb vm.StateDB, api *APIV2, db urlcommon.DatastoreInterface, url *concurrenturl.ConcurrentUrl) *EUV2 {
	return &EUV2{
		evm:     vm.NewEVMEx(vm.BlockContext{BlockNumber: new(big.Int).SetUint64(100000000)}, vm.TxContext{}, statedb, chainConfig, vmConfig, api),
		statedb: statedb,
		api:     api,
		db:      db,
		url:     url,
	}
}

func (eu *EUV2) SetContext(statedb vm.StateDB, api *APIV2, db urlcommon.DatastoreInterface, url *concurrenturl.ConcurrentUrl) {
	eu.api = api
	eu.statedb = statedb
	eu.db = db
	eu.url = url

	eu.evm.StateDB = eu.statedb
	eu.evm.SetApi(api)
}

func (eu *EUV2) Run(thash common.Hash, tindex int, msg *types.Message, blockContext vm.BlockContext, txContext vm.TxContext) ([]urlcommon.UnivalueInterface, []urlcommon.UnivalueInterface, *types.Receipt) {
	eu.statedb.(*ethStateV2).Prepare(thash, common.Hash{}, tindex)
	eu.api.Prepare(thash, blockContext.BlockNumber, uint32(tindex))
	eu.evm.Context = blockContext
	eu.evm.TxContext = txContext

	gp := core.GasPool(math.MaxUint64)

	result, err := core.ApplyMessage(eu.evm, *msg, &gp)
	if err != nil {
		result = &core.ExecutionResult{
			Err: err,
		}
	}

	assertLog := GetAssert(result.ReturnData)
	if len(assertLog) > 0 {
		eu.api.AddLog("assert", assertLog)
	}

	receipt := types.NewReceipt(nil, result.Failed(), result.UsedGas)
	receipt.TxHash = thash
	receipt.GasUsed = result.UsedGas
	if msg.To() == nil {
		userSpecifiedAddress := crypto.CreateAddress(eu.evm.Origin, msg.Nonce())
		receipt.ContractAddress = result.ContractAddress
		if !bytes.Equal(userSpecifiedAddress.Bytes(), result.ContractAddress.Bytes()) {
			eu.api.AddLog("ContractAddressWarning", fmt.Sprintf("user specified address = %v, inner address = %v", userSpecifiedAddress, result.ContractAddress))
		}
	}
	receipt.Logs = eu.statedb.(*ethStateV2).GetLogs(thash)
	receipt.Bloom = types.CreateBloom(types.Receipts{receipt})

	if !result.Failed() {
		accesses, transitions := eu.url.Export(false)
		_, nonceTransitions := ExportOnConfliction(eu.db, tindex, msg.From())
		return accesses, append(transitions, nonceTransitions...), receipt
	} else {
		accesses, transitions := ExportOnFailure(eu.db, tindex, msg.From(), blockContext.Coinbase, receipt.GasUsed, msg.GasPrice())
		return accesses, transitions, receipt
	}
}

func (eu *EUV2) RunEx(thash common.Hash, tindex int, msg *types.Message, blockContext vm.BlockContext, txContext vm.TxContext) ([][]byte, [][]byte, [][]byte, *types.Receipt, []byte) {
	eu.statedb.(*ethStateV2).Prepare(thash, common.Hash{}, tindex)
	eu.api.Prepare(thash, blockContext.BlockNumber, uint32(tindex))
	eu.evm.Context = blockContext
	eu.evm.TxContext = txContext

	gp := core.GasPool(math.MaxUint64)

	result, err := core.ApplyMessage(eu.evm, *msg, &gp)

	if err != nil {
		fmt.Printf("core.ApplyMessage err:%v\n", err)
		result = &core.ExecutionResult{
			Err: err,
		}
	}

	if result.Err != nil {
		fmt.Printf("result.Err err: %v\n", result.Err)
	}

	assertLog := GetAssert(result.ReturnData)
	if len(assertLog) > 0 {
		eu.api.AddLog("assert", assertLog)
	}

	receipt := types.NewReceipt(nil, result.Failed(), result.UsedGas)
	receipt.TxHash = thash
	receipt.GasUsed = result.UsedGas
	if msg.To() == nil {
		userSpecifiedAddress := crypto.CreateAddress(eu.evm.Origin, msg.Nonce())
		receipt.ContractAddress = result.ContractAddress
		if !bytes.Equal(userSpecifiedAddress.Bytes(), result.ContractAddress.Bytes()) {
			eu.api.AddLog("ContractAddressWarning", fmt.Sprintf("user specified address = %v, inner address = %v", userSpecifiedAddress, result.ContractAddress))
		}
	}
	receipt.Logs = eu.statedb.(*ethStateV2).GetLogs(thash)
	receipt.Bloom = types.CreateBloom(types.Receipts{receipt})

	var accesses, transitions [][]byte
	if !result.Failed() {
		accesses, transitions = eu.url.ExportEncoded(nil)
	} else {
		accesses, transitions = ExportOnFailureEx(eu.db, tindex, msg.From(), blockContext.Coinbase, receipt.GasUsed, msg.GasPrice())
	}

	_, nonceTransitions := ExportOnConflictionEx(eu.db, tindex, msg.From())
	return accesses, transitions, nonceTransitions, receipt, result.ReturnData
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
