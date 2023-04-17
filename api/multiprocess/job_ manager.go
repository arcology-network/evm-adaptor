package multiprocess

import (
	"math/big"

	concurrenturl "github.com/arcology-network/concurrenturl/v2"
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
	jobs       []Job
	threads    int
	apiRouter  eucommon.ConcurrentApiRouterInterface
	arbitrator *concurrenturl.ArbitratorSlow
}

func NewJobManager(apiRouter eucommon.ConcurrentApiRouterInterface) *JobManager {
	return &JobManager{
		jobs:       []Job{},
		apiRouter:  apiRouter,
		threads:    16, // 16 threads by default
		arbitrator: concurrenturl.NewArbitratorSlow(),
	}
}

func (this *JobManager) Add(calleeAddr evmcommon.Address, funCall []byte) int {
	sender := this.apiRouter.From()
	this.jobs = append(this.jobs,
		Job{
			sender: sender,
			caller: evmcommon.Address{},
			callee: calleeAddr,
			message: types.NewMessage( // Build the message
				sender,
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
	return len(this.jobs) - 1
}

func (this *JobManager) Start() {
	// processor := func(start, end, index int, args ...interface{}) {
	for i := 0; i < len(this.jobs); i++ {
		// for i := start; i < end; i++ {
		statedb := eth.NewImplStateDB(this.apiRouter.Ccurl()) // Eth state DB
		statedb.Prepare([32]byte{}, [32]byte{}, i)            // tx hash , block hash and tx index

		eu := cceu.NewEU(
			params.MainnetChainConfig,
			vm.Config{},
			statedb,
			this.apiRouter.New(evmcommon.Hash{}, 0, this.apiRouter.Ccurl()), // Call the function
		)

		config := cceu.NewConfig()
		this.jobs[i].accesses, this.jobs[i].transitions, this.jobs[i].receipt, this.jobs[i].result, this.jobs[i].prechkErr =
			eu.Run(evmcommon.Hash{}, i, &this.jobs[i].message, cceu.NewEVMBlockContext(config), cceu.NewEVMTxContext(this.jobs[i].message))
	}
	// }
	// common.ParallelWorker(len(this.jobs), this.threads, processor)

	// Copy all the transition to a single array for conflict detection
	total := 0
	for i := 0; i < len(this.jobs); i++ {
		total += len(this.jobs[i].accesses)
	}
	accesses := make([]ccurlcommon.UnivalueInterface, total) // Pre-allocation for better performance

	offset := 0
	for i := 0; i < len(this.jobs); i++ {
		copy(accesses[offset:], this.jobs[i].accesses)
		offset += len(this.jobs[i].accesses)
	}

	// Detect conflicts
	arbitrator := concurrenturl.NewArbitratorSlow()
	arbitrator.Detect(accesses)
}
