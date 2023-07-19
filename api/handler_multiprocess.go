package api

import (
	"math"
	"math/big"
	"sync/atomic"

	"github.com/arcology-network/common-lib/common"
	evmcommon "github.com/arcology-network/evm/common"
	evmcore "github.com/arcology-network/evm/core"

	"github.com/arcology-network/vm-adaptor/abi"
	execution "github.com/arcology-network/vm-adaptor/execution"

	eucommon "github.com/arcology-network/vm-adaptor/common"
)

// APIs under the concurrency namespace
type MultiprocessHandlers struct {
	*BaseHandlers
	erros   []error
	jobseqs []*execution.JobSequence
}

func NewMultiprocessHandlers(ethApiRouter eucommon.EthApiRouter) *MultiprocessHandlers {
	handler := &MultiprocessHandlers{
		erros:   []error{},
		jobseqs: []*execution.JobSequence{},
	}
	handler.BaseHandlers = NewBaseHandlers(ethApiRouter, handler)
	return handler
}

func (this *MultiprocessHandlers) Address() [20]byte { return eucommon.PARALLEL_HANDLER }

func (this *MultiprocessHandlers) Run(caller [20]byte, input []byte) ([]byte, bool, int64) {
	if atomic.AddUint64(&eucommon.TotalSubProcesses, 1); !this.Api().CheckRuntimeConstrains() {
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
		funCall, successful, fee := this.GetByIndex(path, uint64(i))
		if fees[i] = fee; successful {
			this.jobseqs[i], this.erros[i] = this.toJobSeq(funCall)
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

func (this *MultiprocessHandlers) toJobSeq(input []byte) (*execution.JobSequence, error) {
	gasLimit, calleeAddr, funCall, err := abi.Parse3(input,
		uint64(0), 1, 32,
		[20]byte{}, 1, 32,
		[]byte{}, 2, math.MaxInt64)

	if err != nil {
		return nil, err
	}

	newJobSeq := &execution.JobSequence{
		ID:        this.BaseHandlers.Api().GetSerialNum(eucommon.SUB_PROCESS),
		ApiRouter: this.BaseHandlers.Api(),
	}

	addr := evmcommon.Address(calleeAddr)
	evmMsg := evmcore.NewMessage( // Build the message
		this.BaseHandlers.Api().Origin(),
		&addr,
		0,
		new(big.Int).SetUint64(0), // Amount to transfer
		gasLimit,
		this.BaseHandlers.Api().GetEU().(*execution.EU).Message().Native.GasPrice, // gas price
		funCall,
		nil,
		false, // Don't checking nonce
	)

	stdMsg := &execution.StandardMessage{
		ID:     newJobSeq.ID, // this is the problem !!!!
		Native: &evmMsg,
		TxHash: newJobSeq.DeriveNewHash(this.BaseHandlers.Api().GetEU().(*execution.EU).Message().TxHash),
	}

	newJobSeq.StdMsgs = []*execution.StandardMessage{stdMsg}
	return newJobSeq, nil
}
