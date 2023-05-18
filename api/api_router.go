package api

import (
	"math/big"

	"github.com/arcology-network/common-lib/codec"
	common "github.com/arcology-network/common-lib/common"
	"github.com/arcology-network/common-lib/types"

	"github.com/arcology-network/concurrenturl"
	evmcommon "github.com/arcology-network/evm/common"
	"github.com/arcology-network/evm/core/vm"
	cceu "github.com/arcology-network/vm-adaptor"
	apicommon "github.com/arcology-network/vm-adaptor/api/common"
	cumulativei256 "github.com/arcology-network/vm-adaptor/api/commutative/int256"
	cumulativeu256 "github.com/arcology-network/vm-adaptor/api/commutative/u256"
	noncommutativeBytes "github.com/arcology-network/vm-adaptor/api/noncommutative/base"
	threading "github.com/arcology-network/vm-adaptor/api/threading"
	eucommon "github.com/arcology-network/vm-adaptor/common"
)

type API struct {
	logs         []eucommon.ILog
	txHash       evmcommon.Hash // Tx hash
	txIndex      uint32         // Tx index in the block
	dc           *types.DeferCall
	predecessors []evmcommon.Hash

	seed   uint64 // for uuid generation
	serial uint64
	// deferCall *concurrentlib.DeferCall

	eu          *cceu.EU
	handlerDict map[[20]byte]apicommon.ConcurrentApiHandlerInterface // APIs under the concurrency namespace
	ccurl       *concurrenturl.ConcurrentUrl
}

func NewAPI(ccurl *concurrenturl.ConcurrentUrl) *API {
	api := &API{
		eu:          nil,
		ccurl:       ccurl,
		handlerDict: make(map[[20]byte]apicommon.ConcurrentApiHandlerInterface),
	}

	handlers := []apicommon.ConcurrentApiHandlerInterface{
		noncommutativeBytes.NewNoncommutativeBytesHandlers(api),
		cumulativeu256.NewU256CumulativeHandlers(api),
		cumulativei256.NewInt256CumulativeHandlers(api),
		threading.NewThreadingHandler(api),
	}

	for i, v := range handlers {
		if _, ok := api.handlerDict[(handlers)[i].Address()]; ok {
			panic("Error: Duplicate handler addresses found!!")
		}
		api.handlerDict[(handlers)[i].Address()] = v
	}
	return api
}

func (this *API) New(txHash evmcommon.Hash, txIndex uint32, ccurl *concurrenturl.ConcurrentUrl) eucommon.ConcurrentApiRouterInterface {
	api := NewAPI(ccurl)
	api.txHash = txHash
	api.txIndex = txIndex
	// api.SetEU(this.eu)
	return api
}

func (this *API) Coinbase() evmcommon.Address { return this.eu.VM().Context.Coinbase }
func (this *API) Origin() evmcommon.Address   { return this.eu.VM().TxContext.Origin }
func (this *API) VM() *vm.EVM                 { return this.eu.VM() }
func (this *API) SetEU(eu interface{})        { this.eu = eu.(*cceu.EU) }

func (this *API) TxHash() [32]byte                    { return this.txHash }
func (this *API) TxIndex() uint32                     { return this.txIndex }
func (this *API) Ccurl() *concurrenturl.ConcurrentUrl { return this.ccurl }

func (this *API) Prepare(txHash evmcommon.Hash, height *big.Int, txIndex uint32) {
	this.txHash = txHash
	this.txIndex = txIndex
	this.dc = nil
}

func (this *API) GenElemUID() uint64 {
	this.serial++
	return this.serial
}

// Generate an UUID based on transaction hash and the counter
func (this *API) GenCtrnUID() []byte {
	this.seed++
	id := codec.Bytes32(this.txHash).UUID(this.seed)
	return id[:8]
}

func (this *API) AddLog(key, value string) {
	this.logs = append(this.logs, &types.ExecutingLog{
		Key:   key,
		Value: value,
	})
}

func (this *API) GetLogs() []eucommon.ILog {
	return this.logs
}

func (this *API) ClearLogs() {
	this.logs = this.logs[:0]
}

func (this *API) Call(caller, callee evmcommon.Address, input []byte, origin evmcommon.Address, nonce uint64, blockhash evmcommon.Hash) (bool, []byte, bool) {
	if handler, ok := this.handlerDict[callee]; ok {
		result, successful := handler.Call(
			evmcommon.Address(codec.Bytes20(caller).Clone().(codec.Bytes20)),
			evmcommon.Address(codec.Bytes20(callee).Clone().(codec.Bytes20)),
			common.Clone(input),
			origin,
			nonce,
		)
		return true, result, successful
	}
	return false, []byte{}, false
}

// For defer call.
func (this *API) SetDeferCall(contractAddress types.Address, deferID string) {
	// sig := this.deferCall.GetSignature(contractAddress, deferID)
	// if sig == "" {
	// 	panic(fmt.Sprintf("unknown defer call on %s:%s", contractAddress, deferID))
	// }

	this.dc = &types.DeferCall{
		DeferID:         deferID,
		ContractAddress: contractAddress,
		// Signature:       sig,
	}
}

func (this *API) GetDeferCall() *types.DeferCall {
	return this.dc
}

func (this *API) SetPredecessors(predecessors []evmcommon.Hash) {
	this.predecessors = predecessors
}

func (this *API) IsInDeferCall() bool {
	return len(this.predecessors) > 0
}
