package api

import (
	"encoding/hex"
	"strconv"
	"sync/atomic"

	"github.com/arcology-network/common-lib/codec"
	common "github.com/arcology-network/common-lib/common"
	"github.com/arcology-network/concurrenturl"
	evmcommon "github.com/arcology-network/evm/common"
	"github.com/arcology-network/evm/core/vm"
	cumulativeu256 "github.com/arcology-network/vm-adaptor/api/commutative/u256"
	"github.com/arcology-network/vm-adaptor/api/parallel"
	runtime "github.com/arcology-network/vm-adaptor/api/runtime"
	execution "github.com/arcology-network/vm-adaptor/execution"

	ccurlcommon "github.com/arcology-network/concurrenturl/common"
	noncommutativeBytes "github.com/arcology-network/vm-adaptor/api/noncommutative/base"
	eucommon "github.com/arcology-network/vm-adaptor/common"
)

type API struct {
	logs  []eucommon.ILog
	depth uint8

	serialNums [4]uint64 // sub-process/container/element/uuid generator,

	schedule interface{}
	eu       *execution.EU
	reserved interface{}

	handlerDict map[[20]byte]eucommon.ApiCallHandler // APIs under the atomic namespace
	ccurl       *concurrenturl.ConcurrentUrl

	execResult *execution.Result

	filter eucommon.StateFilter
}

func NewAPI(ccurl *concurrenturl.ConcurrentUrl) *API {
	api := &API{
		eu:          nil,
		ccurl:       ccurl,
		handlerDict: make(map[[20]byte]eucommon.ApiCallHandler),
		depth:       0,
		execResult:  &execution.Result{},
		serialNums:  [4]uint64{},
	}
	api.filter = NewExportFilter(api)

	handlers := []eucommon.ApiCallHandler{
		parallel.NewParallelHandler(api),
		noncommutativeBytes.NewNoncommutativeBytesHandlers(api, nil),
		cumulativeu256.NewU256CumulativeHandlers(api),
		// cumulativei256.NewInt256CumulativeHandlers(api),
		runtime.NewRuntimeHandler(api),
	}

	for i, v := range handlers {
		if _, ok := api.handlerDict[(handlers)[i].Address()]; ok {
			panic("Error: Duplicate handler addresses found!!")
		}
		api.handlerDict[(handlers)[i].Address()] = v
	}

	api.ccurl.NewAccount(
		ccurlcommon.SYSTEM,
		hex.EncodeToString(codec.Bytes20(runtime.NewRuntimeHandler(api).Address()).Encode()),
	)
	return api
}

func (this *API) New(ccurl *concurrenturl.ConcurrentUrl, schedule interface{}) eucommon.EthApiRouter {
	api := NewAPI(ccurl)
	api.depth = this.depth + 1
	return api
}

func (this *API) CheckRuntimeConstrains() bool { // Execeeds the max recursion depth or the max sub processes
	return this.Depth() < eucommon.MAX_RECURSIION_DEPTH &&
		atomic.AddUint64(&eucommon.TotalSubProcesses, 1) <= eucommon.MAX_SUB_PROCESSES
}

func (this *API) StateFilter() eucommon.StateFilter { return this.filter }

func (this *API) IsLocal(txID uint32) bool         { return txID == ccurlcommon.SYSTEM } //A local tx
func (this *API) GetReserved() interface{}         { return this.reserved }
func (this *API) SetReserved(reserved interface{}) { this.reserved = reserved }

func (this *API) Depth() uint8                { return this.depth }
func (this *API) Coinbase() evmcommon.Address { return this.eu.VM().Context.Coinbase }
func (this *API) Origin() evmcommon.Address   { return this.eu.VM().TxContext.Origin }

func (this *API) SetSchedule(schedule interface{}) { this.schedule = schedule }
func (this *API) Schedule() interface{}            { return this.schedule }

func (this *API) VM() *vm.EVM {
	return common.IfThenDo1st(this.eu != nil, func() *vm.EVM { return this.eu.VM() }, nil)
}

func (this *API) GetEU() interface{}   { return this.eu }
func (this *API) SetEU(eu interface{}) { this.eu = eu.(*execution.EU) }

func (this *API) Ccurl() *concurrenturl.ConcurrentUrl { return this.ccurl }

func (this *API) GetSerialNum(idx int) uint64 {
	v := this.serialNums[idx]
	this.serialNums[idx]++
	return v
}

func (this *API) ElementUID() []byte {
	return []byte(hex.EncodeToString(this.eu.Message().TxHash[:8]) + "-" + strconv.Itoa(int(this.GetSerialNum(eucommon.ELEMENT_ID))))
}

// Generate an UUID based on transaction hash and the counter
func (this *API) UUID() []byte {
	id := codec.Bytes32(this.eu.Message().TxHash).UUID(this.GetSerialNum(eucommon.UUID))
	return id[:8]
}

func (this *API) AddLog(key, value string) {
	this.logs = append(this.logs, &execution.ExecutionLog{
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

func (this *API) Call(caller, callee [20]byte, input []byte, origin [20]byte, nonce uint64, blockhash evmcommon.Hash) (bool, []byte, bool, int64) {
	if handler, ok := this.handlerDict[callee]; ok {
		result, successful, fees := handler.Call(
			evmcommon.Address(codec.Bytes20(caller).Clone().(codec.Bytes20)),
			evmcommon.Address(codec.Bytes20(callee).Clone().(codec.Bytes20)),
			common.Clone(input),
			origin,
			nonce,
		)
		return true, result, successful, fees
	}
	return false, []byte{}, true, 0 // not an Arcology call, used 0 gas
}
