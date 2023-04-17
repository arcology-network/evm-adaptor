package multiprocess

import (
	"math/big"

	common "github.com/arcology-network/common-lib/common"
	"github.com/arcology-network/concurrenturl/v2"
	ccurlcommon "github.com/arcology-network/concurrenturl/v2/common"
	evmcommon "github.com/arcology-network/evm/common"

	types "github.com/arcology-network/evm/core/types"
	"github.com/arcology-network/evm/core/vm"
	"github.com/arcology-network/evm/params"

	cceu "github.com/arcology-network/vm-adaptor"
	eucommon "github.com/arcology-network/vm-adaptor/common"
	"github.com/arcology-network/vm-adaptor/eth"
)

type Job struct {
	sender      evmcommon.Address
	caller      evmcommon.Address
	callee      evmcommon.Address
	message     types.Message
	receipt     *types.Receipt
	result      interface{}
	prechkErr   error
	accesses    []ccurlcommon.UnivalueInterface
	transitions []ccurlcommon.UnivalueInterface
}

// APIs under the concurrency namespace
type JobManager struct {
	ccurl     *concurrenturl.ConcurrentUrl
	jobs      []Job
	threads   int
	apiRouter eucommon.ConcurrentApiRouterInterface
	// arbitrator ccur
}

func NewJobManager(apiRouter eucommon.ConcurrentApiRouterInterface) *JobManager {
	return &JobManager{
		jobs:      []Job{},
		apiRouter: apiRouter,
		threads:   16, // 16 threads by default
	}
}

func (this *JobManager) Add(calleeAddr evmcommon.Address, funCall []byte) int {
	this.apiRouter.From()
	this.jobs = append(this.jobs,
		Job{
			sender: this.apiRouter.From(),
			caller: evmcommon.Address{},
			callee: calleeAddr,
			message: types.NewMessage( // Build the message
				this.apiRouter.From(),
				&calleeAddr,
				0,
				new(big.Int).SetUint64(0), // Amount to transfer
				1e15,
				new(big.Int).SetUint64(1),
				funCall, //need to remove the wrapper first
				nil,
				false, // Stop checking nonce
			),
		},
	)
	return len(this.jobs)
}

func (this *JobManager) Start() {
	statedb := eth.NewImplStateDB(this.ccurl) // Eth state DB
	statedb.Prepare([32]byte{}, [32]byte{}, len(this.jobs))

	hasher := func(start, end, index int, args ...interface{}) {
		for i := start; i < end; i++ {
			statedb := eth.NewImplStateDB(this.ccurl)  // Eth state DB
			statedb.Prepare([32]byte{}, [32]byte{}, i) // tx hash , block hash and tx index

			eu := cceu.NewEU(
				params.MainnetChainConfig,
				vm.Config{},
				statedb,
				this.apiRouter.New(evmcommon.Hash{}, 0, this.ccurl), // Call function
			)

			config := cceu.NewConfig()
			// var result *core.ExecutionResult
			accesses, transitions, receipt, result, err :=
				eu.Run(evmcommon.Hash{}, i, &this.jobs[i].message, cceu.NewEVMBlockContext(config), cceu.NewEVMTxContext(this.jobs[i].message))

			this.jobs[i].accesses = accesses
			this.jobs[i].transitions = transitions
			this.jobs[i].receipt = receipt
			this.jobs[i].result = result
			this.jobs[i].prechkErr = err
		}
	}
	common.ParallelWorker(len(this.jobs), this.threads, hasher)

	total := 0
	for i := 0; i < len(this.jobs); i++ {
		total += len(this.jobs[i].accesses)
	}

	offset := 0
	accesses := make([]ccurlcommon.UnivalueInterface, total)
	for i := 0; i < len(this.jobs); i++ {
		copy(accesses[offset:], this.jobs[i].accesses)
		offset += len(this.jobs[i].accesses)
	}

	arbitrator := concurrenturl.NewArbitratorSlow()
	arbitrator.Detect(accesses)
}
