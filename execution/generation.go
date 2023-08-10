package execution

import (
	"errors"

	common "github.com/arcology-network/common-lib/common"
	ccurlcommon "github.com/arcology-network/concurrenturl/common"
	indexer "github.com/arcology-network/concurrenturl/indexer"

	// evmeu "github.com/arcology-network/vm-adaptor"
	eucommon "github.com/arcology-network/vm-adaptor/common"
)

// APIs under the concurrency namespace
type Generation struct {
	ID         uint32
	numThreads uint8
	jobs       []*JobSequence // para jobs
}

func NewGeneration(id uint32, numThreads uint8, jobs []*JobSequence) *Generation {
	return &Generation{
		ID:         id,
		numThreads: numThreads,
		jobs:       jobs,
	}
}

// func (this *Generation) BranchID() uint32 { return this.branchID }
func (this *Generation) Length() uint64 { return uint64(len(this.jobs)) }

func (this *Generation) At(idx uint64) *JobSequence {
	return common.IfThenDo1st(idx < uint64(len(this.jobs)), func() *JobSequence { return this.jobs[idx] }, nil)
}

func (this *Generation) Add(job *JobSequence) bool {
	this.jobs = append(this.jobs, job)
	return true
}

func (this *Generation) Run(parentApiRouter eucommon.EthApiRouter) []*Result {
	preTransitions := indexer.Univalues(common.Clone(parentApiRouter.StateFilter().Raw())).To(indexer.ITCTransition{})

	snapshot := parentApiRouter.Ccurl().Snapshot(preTransitions)
	config := NewConfig().SetCoinbase(parentApiRouter.Coinbase())

	// t0 := time.Now()
	worker := func(start, end, idx int, args ...interface{}) {
		for i := start; i < end; i++ {
			this.jobs[i].Results = this.jobs[i].Run(config, snapshot)
		}
	}
	common.ParallelWorker(len(this.jobs), int(this.numThreads), worker)
	// fmt.Println(time.Since(t0))

	// Detect potential conflicts
	results := common.Concate(this.jobs, func(job *JobSequence) []*Result { return job.Results })
	dict, _ := Results(results).Detect().ToDict()

	for i := 0; i < len(results); i++ {
		if _, conflict := (*dict)[results[i].TxIndex]; conflict {
			results[i].Err = errors.New(ccurlcommon.WARN_ACCESS_CONFLICT)
		}
		results[i].ImmunizeGasTransition()
	}

	// this.ImmunizeGasTransition(this.Transitions)
	return results
}

func (this *Generation) Clear() uint64 {
	length := len(this.jobs)
	this.jobs = this.jobs[:0]
	return uint64(length)
}
