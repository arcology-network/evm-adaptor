package api

import (
	"math"
	"sync/atomic"

	"github.com/arcology-network/common-lib/common"
	"github.com/arcology-network/concurrenturl/univalue"
	"github.com/arcology-network/eu/cache"
	evmcommon "github.com/ethereum/go-ethereum/common"
	evmcore "github.com/ethereum/go-ethereum/core"
	"github.com/holiman/uint256"

	"github.com/arcology-network/vm-adaptor/abi"

	eu "github.com/arcology-network/eu"
	eucommon "github.com/arcology-network/eu/common"
	adaptorcommon "github.com/arcology-network/vm-adaptor/common"
	adaptorintf "github.com/arcology-network/vm-adaptor/interface"

	basecontainer "github.com/arcology-network/vm-adaptor/apihandler/container"
)

// APIs under the concurrency namespace
type MultiprocessHandler struct {
	*basecontainer.BaseHandlers
	erros   []error
	jobseqs []*eu.JobSequence
}

func NewMultiprocessHandler(ethApiRouter adaptorintf.EthApiRouter) *MultiprocessHandler {
	handler := &MultiprocessHandler{
		erros:   []error{},
		jobseqs: common.To[*eu.JobSequence, *eu.JobSequence]([]*eu.JobSequence{}),
	}
	handler.BaseHandlers = basecontainer.NewBaseHandlers(ethApiRouter, handler.Run, &eu.Generation{})
	return handler
}

func (this *MultiprocessHandler) Address() [20]byte { return adaptorcommon.MULTIPROCESS_HANDLER }

func (this *MultiprocessHandler) Run(caller [20]byte, input []byte, args ...interface{}) ([]byte, bool, int64) {
	if atomic.AddUint64(&eucommon.TotalSubProcesses, 1); !this.Api().CheckRuntimeConstrains() {
		return []byte{}, false, 0
	}

	input, err := abi.DecodeTo(input, 0, []byte{}, 2, math.MaxInt64)
	if err != nil {
		return []byte{}, false, 0
	}

	numThreads, err := abi.DecodeTo(input, 0, uint64(1), 1, 8)
	if err != nil {
		return []byte{}, false, 0
	}
	threads := common.Min(common.Max(uint8(numThreads), 1), math.MaxUint8) // [1, 255]

	path := this.Connector().Key(caller)
	length, successful, fee := this.Length(path)
	length = common.Min(eucommon.MAX_VM_INSTANCES, length)
	if !successful {
		return []byte{}, successful, fee
	}

	generation := args[0].(*eu.Generation).New(0, threads, args[0].(*eu.Generation).JobSeqs()[:0])
	fees := make([]int64, length)
	this.erros = make([]error, length)

	this.jobseqs = common.Resize(this.jobseqs, int(length))
	for i := uint64(0); i < length; i++ {
		funCall, successful, fee := this.GetByIndex(path, uint64(i))
		if fees[i] = fee; successful {
			this.jobseqs[i], this.erros[i] = this.toJobSeq(funCall, generation.JobT())
		}
		generation.Add(this.jobseqs[i]) // Add the job sequence to the generation regardless of the error
	}
	transitions := generation.Run(this.Api()) // Run the generation

	// Sub processes may have been spawned during the execution, recheck it.
	if !this.Api().CheckRuntimeConstrains() {
		return []byte{}, false, fee
	}

	// Unify tx IDs c
	mainTxID := uint32(this.Api().GetEU().(adaptorintf.EU).ID())
	common.Foreach(transitions, func(v **univalue.Univalue, _ int) { (*v).SetTx(mainTxID) })

	this.Api().WriteCache().(*cache.WriteCache).AddTransitions(transitions) // Merge the write cache to the main cache
	return []byte{}, true, common.Sum[int64](fees)
}

// toJobSeq converts the input byte slice into a JobSequence object.
// For multiprocessor, a job sequence only contains one message.
// To keep the same structure with the transaction level processing,
// the message is wrapped
func (this *MultiprocessHandler) toJobSeq(input []byte, T *eu.JobSequence) (*eu.JobSequence, error) {
	gasLimit, value, calleeAddr, funCall, err := abi.Parse4(input,
		uint64(0), 1, 32,
		uint256.NewInt(0), 1, 32,
		[20]byte{}, 1, 32,
		[]byte{}, 2, math.MaxInt64)

	if err != nil {
		return nil, err
	}

	transfer := value.ToBig()
	addr := evmcommon.Address(calleeAddr)
	evmMsg := evmcore.NewMessage( // Build the message
		this.BaseHandlers.Api().Origin(),
		&addr,
		0,
		transfer, // Amount to transfer
		gasLimit,
		this.BaseHandlers.Api().GetEU().(adaptorintf.EU).GasPrice(), // gas price
		funCall,
		nil,
		false, // Don't checking nonce
	)

	// newJobSeq creates a new job sequence using the TYPE INFO.
	newJobSeq := T.New(
		uint32(this.BaseHandlers.Api().GetSerialNum(eucommon.SUB_PROCESS)),
		this.BaseHandlers.Api(),
	)

	newJobSeq.AppendMsg(&eucommon.StandardMessage{
		ID:     uint64(newJobSeq.GetID()),
		Native: &evmMsg,
		TxHash: newJobSeq.DeriveNewHash(this.BaseHandlers.Api().GetEU().(adaptorintf.EU).TxHash()),
	})
	return newJobSeq, nil
}
