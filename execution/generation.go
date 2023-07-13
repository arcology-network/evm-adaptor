package execution

import (
	"errors"
	"fmt"

	common "github.com/arcology-network/common-lib/common"
	ccurlcommon "github.com/arcology-network/concurrenturl/common"
	indexer "github.com/arcology-network/concurrenturl/indexer"

	// evmeu "github.com/arcology-network/vm-adaptor"
	eucommon "github.com/arcology-network/vm-adaptor/common"
)

// APIs under the concurrency namespace
type Generation struct {
	ID         uint32
	maxThreads uint8
	jobs       []*JobSequence // para jobs
}

func NewGeneration(id uint32, maxThreads uint8, jobs []*JobSequence) *Generation {
	return &Generation{
		ID:         id,
		maxThreads: maxThreads,
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
	// common.ParallelForeach(this.jobs, this.maxThreads, func(job **JobSequence) *JobSequence {
	// 	(**job).Results = (**job).Run(config, snapshot)
	// 	return (*job)
	// })

	for i := 0; i < len(this.jobs); i++ {
		this.jobs[i].Results = this.jobs[i].Run(config, snapshot)
	}

	// Detect potential conflicts
	results := common.Concate(this.jobs, func(job *JobSequence) []*Result { return job.Results })

	conflicts := Results(results).Detect()
	dict := conflicts.ToDict()

	if len(*dict) > 0 {
		Results(results).Detect()
		fmt.Println("Warning: Conflict detected!", len(*dict))
		conflicts.Print()
	}

	for i := 0; i < len(results); i++ {
		if _, conflict := (*dict)[results[i].TxIndex]; conflict {
			results[i].Err = errors.New(ccurlcommon.ERR_ACCESS_CONFLICT)
		}
	}

	// common.Foreach(results, func(v **Result) { // Write the transitions back to the parent write cache
	// 	(*v).WriteTo(uint32(parentApiRouter.GetEU().(*EU).Message().ID), parentApiRouter.Ccurl().WriteCache()) // Merge the write cache to its parent
	// })

	return results
}

func (this *Generation) Clear() uint64 {
	length := len(this.jobs)
	this.jobs = this.jobs[:0]
	return uint64(length)
}
