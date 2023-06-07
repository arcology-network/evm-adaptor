package api

import (
	"encoding/hex"
	"math/big"
	"strconv"

	"github.com/arcology-network/common-lib/codec"
	common "github.com/arcology-network/common-lib/common"
	commontypes "github.com/arcology-network/common-lib/types"
	"github.com/arcology-network/concurrenturl"
	evmcommon "github.com/arcology-network/evm/common"
	evmcoretypes "github.com/arcology-network/evm/core/types"
	"github.com/arcology-network/evm/core/vm"
	cceu "github.com/arcology-network/vm-adaptor"
	cumulativei256 "github.com/arcology-network/vm-adaptor/api/commutative/int256"
	cumulativeu256 "github.com/arcology-network/vm-adaptor/api/commutative/u256"
	"github.com/arcology-network/vm-adaptor/api/concurrency"

	noncommutativeBytes "github.com/arcology-network/vm-adaptor/api/noncommutative/base"
	threading "github.com/arcology-network/vm-adaptor/api/threading"
	interfaces "github.com/arcology-network/vm-adaptor/interfaces"
	types "github.com/arcology-network/vm-adaptor/types"
)

type API struct {
	logs    []interfaces.ILog
	txHash  evmcommon.Hash // Tx hash
	txIndex uint32         // Tx index in the block

	uuid     uint64
	ccUID    uint64 // for uuid generation
	ccElemID uint64
	depth    uint8

	reserved interface{}
	eu       *cceu.EU

	handlerDict map[[20]byte]interfaces.ApiCallHandler // APIs under the concurrency namespace
	ccurl       *concurrenturl.ConcurrentUrl

	execResult *types.ExecutionResult
}

func NewAPI(ccurl *concurrenturl.ConcurrentUrl) *API {
	api := &API{
		eu:          nil,
		ccurl:       ccurl,
		handlerDict: make(map[[20]byte]interfaces.ApiCallHandler),
		depth:       0,
		execResult:  &types.ExecutionResult{},
	}

	handlers := []interfaces.ApiCallHandler{
		noncommutativeBytes.NewNoncommutativeBytesHandlers(api),
		cumulativeu256.NewU256CumulativeHandlers(api),
		cumulativei256.NewInt256CumulativeHandlers(api),
		threading.NewThreadingHandler(api),
		concurrency.NewConcurrencyHandler(api),
	}

	for i, v := range handlers {
		if _, ok := api.handlerDict[(handlers)[i].Address()]; ok {
			panic("Error: Duplicate handler addresses found!!")
		}
		api.handlerDict[(handlers)[i].Address()] = v
	}
	return api
}

func (this *API) New(txHash evmcommon.Hash, txIndex uint32, parentDepth uint8, ccurl *concurrenturl.ConcurrentUrl) interfaces.EthApiRouter {
	api := NewAPI(ccurl)

	api.txHash = txHash
	api.txIndex = txIndex

	api.uuid = 0
	api.ccUID = 0
	api.ccElemID = 0

	api.depth = parentDepth + 1
	return api
}

func (this *API) GetReserved() interface{}         { return this.reserved }
func (this *API) SetReserved(reserved interface{}) { this.reserved = reserved }

func (this *API) Depth() uint8                { return this.depth }
func (this *API) Coinbase() evmcommon.Address { return this.eu.VM().Context.Coinbase }
func (this *API) Origin() evmcommon.Address   { return this.eu.VM().TxContext.Origin }

func (this *API) Message() *evmcoretypes.Message { return this.eu.Message() }

func (this *API) VM() *vm.EVM { return this.eu.VM() }

func (this *API) GetEU() interface{}   { return this.eu }
func (this *API) SetEU(eu interface{}) { this.eu = eu.(*cceu.EU) }

func (this *API) TxHash() [32]byte                    { return this.txHash }
func (this *API) TxIndex() uint32                     { return this.txIndex }
func (this *API) Ccurl() *concurrenturl.ConcurrentUrl { return this.ccurl }

func (this *API) SetContext(txHash evmcommon.Hash, height *big.Int, txIndex uint32) {
	this.txHash = txHash
	this.txIndex = txIndex
}

func (this *API) GenUUID() []byte {
	this.uuid++
	id := codec.Bytes32(this.txHash).UUID(this.uuid)
	return id[:]
}

func (this *API) GenCcElemUID() []byte {
	this.ccElemID++
	return []byte(hex.EncodeToString(this.txHash[:8]) + "-" + strconv.Itoa(int(this.ccElemID)))
}

// Generate an UUID based on transaction hash and the counter
func (this *API) GenCcObjID() []byte {
	this.ccUID++
	id := codec.Bytes32(this.txHash).UUID(this.ccUID)
	return id[:8]
}

func (this *API) AddLog(key, value string) {
	this.logs = append(this.logs, &commontypes.ExecutingLog{
		Key:   key,
		Value: value,
	})
}

func (this *API) GetLogs() []interfaces.ILog {
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
