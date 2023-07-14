package parallel

import (
	"errors"
	"math"
	"math/big"
	"sync/atomic"

	"github.com/arcology-network/common-lib/common"
	evmcommon "github.com/arcology-network/evm/common"
	evmcore "github.com/arcology-network/evm/core"

	"github.com/arcology-network/vm-adaptor/abi"
	base "github.com/arcology-network/vm-adaptor/api/noncommutative/base"
	execution "github.com/arcology-network/vm-adaptor/execution"

	eucommon "github.com/arcology-network/vm-adaptor/common"
)

// APIs under the concurrency namespace
type ParallelHandler struct {
	*base.BytesHandlers
	erros   []error
	jobseqs []*execution.JobSequence
}

func NewParallelHandler(ethApiRouter eucommon.EthApiRouter) *ParallelHandler {
	handler := &ParallelHandler{
		erros:   []error{},
		jobseqs: []*execution.JobSequence{},
	}
	handler.BytesHandlers = base.NewNoncommutativeBytesHandlers(ethApiRouter, handler)
	return handler
}

func (this *ParallelHandler) Address() [20]byte { return eucommon.PARALLEL_HANDLER }

// func (this *ParallelHandler) Call(caller, callee [20]byte, input []byte, origin [20]byte, nonce uint64) ([]byte, bool, int64) {
// 	signature := [4]byte{}
// 	copy(signature[:], input)

// 	switch signature {
// 	case [4]byte{0x95, 0x07, 0xd3, 0x9a}: // 95 07 d3 9a
// 		return this.get(caller, input[4:])
// 	}
// 	return this.BytesHandlers.Call(caller, callee, input, origin, nonce)
// }

// func (this *ParallelHandler) get(caller evmcommon.Address, input []byte) ([]byte, bool, int64) {
// 	return []byte{}, false, 0
// }

func (this *ParallelHandler) Run(caller [20]byte, input []byte) ([]byte, bool, int64) {
	atomic.AddUint64(&eucommon.TotalSubProcesses, 1)
	if !this.Api().CheckRuntimeConstrains() {
		return []byte{}, false, 0
	}

	input, err := abi.DecodeTo(input, 0, []byte{}, 2, math.MaxInt64)
	if err != nil {
		return []byte{}, false, 0
	}

	numThreads, err := abi.DecodeTo(input, 0, uint64(1), 2, 32)
	if err != nil {
		return []byte{}, false, 0
	}
	threads := common.Min(common.Min(uint8(numThreads), 1), math.MaxUint8)

	path := this.Connector().Key(caller)
	length, successful, fee := this.Length(path)
	if !successful {
		return []byte{}, successful, fee
	}

	generation := execution.NewGeneration(0, threads, []*execution.JobSequence{})
	fees := make([]int64, length)
	this.erros = make([]error, length)
	this.jobseqs = make([]*execution.JobSequence, length)
	for i := uint64(0); i < length; i++ {
		data, successful, fee := this.Get(path, uint64(i))
		if fees[i] = fee; successful {
			this.jobseqs[i], this.erros[i] = this.toJobSeq(data)
		}
		generation.Add(this.jobseqs[i])
	}

	results := generation.Run(this.Api())

	if !this.Api().CheckRuntimeConstrains() {
		return []byte{}, false, fee
	}

	// execution.Results(results).Print()
	common.Foreach(results, func(v **execution.Result) { // Write the transitions back to the parent write cache
		(*v).WriteTo(uint32(this.Api().GetEU().(*execution.EU).Message().ID), this.Api().Ccurl().WriteCache()) // Merge the write cache to its parent
	})

	return []byte{}, true, common.Sum(fees, int64(0))
}

func (this *ParallelHandler) toJobSeq(input []byte) (*execution.JobSequence, error) {
	input, err := abi.DecodeTo(input, 0, []byte{}, 2, math.MaxInt64)
	if err != nil {
		return nil, errors.New("Error: Unrecognizable input format")
	}

	gasLimit, calleeAddr, funCall, err := abi.Parse3(input,
		uint64(0), 1, 32,
		[20]byte{}, 1, 32,
		[]byte{}, 2, math.MaxInt64)

	if err != nil {
		return nil, err
	}

	newJobSeq := &execution.JobSequence{
		ID:        this.BytesHandlers.Api().GetSerialNum(eucommon.SUB_PROCESS),
		ApiRouter: this.BytesHandlers.Api(),
	}

	addr := evmcommon.Address(calleeAddr)
	evmMsg := evmcore.NewMessage( // Build the message
		this.BytesHandlers.Api().Origin(),
		&addr,
		0,
		new(big.Int).SetUint64(0), // Amount to transfer
		gasLimit,
		this.BytesHandlers.Api().GetEU().(*execution.EU).Message().Native.GasPrice, // gas price
		funCall,
		nil,
		false, // Don't checking nonce
	)

	stdMsg := &execution.StandardMessage{
		ID:     newJobSeq.ID, // this is the problem !!!!
		Native: &evmMsg,
		TxHash: newJobSeq.DeriveNewHash(this.BytesHandlers.Api().GetEU().(*execution.EU).Message().TxHash),
	}

	newJobSeq.StdMsgs = []*execution.StandardMessage{stdMsg}
	return newJobSeq, nil
}
