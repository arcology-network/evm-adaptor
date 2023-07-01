package execution

import (
	"errors"

	common "github.com/arcology-network/common-lib/common"
	ccurlcommon "github.com/arcology-network/concurrenturl/common"
	ccinterfaces "github.com/arcology-network/concurrenturl/interfaces"
	evmeu "github.com/arcology-network/vm-adaptor"
	eucommon "github.com/arcology-network/vm-adaptor/common"
)

// APIs under the concurrency namespace
type Generation struct {
	ID         uint32
	branchID   uint32
	maxThreads uint8
	jobs       []*JobSequence // para jobs
}

func NewGeneration(id uint32, branchID uint32, maxThreads uint8, jobs []*JobSequence) *Generation {
	return &Generation{
		ID:         id,
		branchID:   branchID,
		maxThreads: maxThreads,
		jobs:       jobs,
	}
}

func (this *Generation) BranchID() uint32 { return this.branchID }
func (this *Generation) Length() uint64   { return uint64(len(this.jobs)) }

func (this *Generation) At(idx uint64) *JobSequence {
	return common.IfThenDo1st(idx < uint64(len(this.jobs)), func() *JobSequence { return this.jobs[idx] }, nil)
}

func (this *Generation) Add(job *JobSequence) bool {
	this.jobs = append(this.jobs, job)
	return true
}

func (this *Generation) Run(parentApiRouter eucommon.EthApiRouter, snapshot ccinterfaces.Datastore) []*Result {
	config := evmeu.NewConfig().SetCoinbase(parentApiRouter.Coinbase())

	// common.ParallelForeach(this.jobs, this.maxThreads, func(job **JobSequence) *JobSequence {
	// 	(**job).Results = (**job).Run(config, snapshot)
	// 	return (*job)
	// })

	for i := 0; i < len(this.jobs); i++ {
		this.jobs[i].Results = this.jobs[i].Run(config, snapshot)
	}

	// Detect potential conflicts
	results := common.Concate(this.jobs, func(job *JobSequence) []*Result { return job.Results })
	dict := Results(results).Detect()

	for i := 0; i < len(results); i++ {
		if _, conflict := (*dict)[results[i].TxIndex]; conflict {
			results[i].Err = errors.New(ccurlcommon.ERR_ACCESS_CONFLICT)
		}
	}

	common.Foreach(results, func(v **Result) {
		(*v).WriteTo(parentApiRouter.TxIndex(), parentApiRouter.Ccurl().WriteCache()) // Merge the write cache to its parent
	})

	// indexer.Univalues(results).To()

	// fmt.Println("Sub subResults 1 === ====================== ====================== ====================== =========================================")
	// if len(subResults) > 0 {
	// 	indexer.Univalues(Results((subResults[0])).Transitions()).SortByDefault().Print()
	// }
	// fmt.Println("Sub subResults 2 === ====================== ====================== ====================== =========================================")
	// if len(subResults) > 1 {
	// 	indexer.Univalues(Results((subResults[1])).Transitions()).SortByDefault().Print()
	// }
	// Write the transitions back to the parent write cache
	return results
}

func (this *Generation) Clear() uint64 {
	length := len(this.jobs)
	this.jobs = this.jobs[:0]
	return uint64(length)
}

// func (this *Generation) RunSpawned(parentApiRouter eucommon.EthApiRouter, results []*Result) [][]*Result {
// 	grouped := common.GroupBy(results, func(v *Result) *[32]byte {
// 		return common.IfThenDo1st(v.Spawned != nil, func() *[32]byte { return &(v.Spawned.CallSig) }, nil)
// 	})

// 	spawnedJobs := make([]*Generation, 0, len(grouped))
// 	for i := 0; i < len(grouped); i++ {
// 		seq := Results(grouped[i]).ToSequence()
// 		if job := NewGenerationFromSequence(uint32(i), this.maxThreads, parentApiRouter, seq); job != nil {
// 			spawnedJobs = append(spawnedJobs, job)
// 		}
// 	}

// 	if len(spawnedJobs) == 0 {
// 		return [][]*Result{}
// 	}

// 	spawnedResults := make([][]*Result, len(spawnedJobs))
// 	for i, jobs := range spawnedJobs {
// 		preTransitions := common.Concate(grouped[i], func(v *Result) []ccinterfaces.Univalue { return v.Transitions })

// 		var maker eucommon.LocalSnapshotMaker
// 		maker.Import(preTransitions)
// 		snapshot := maker.Make(parentApiRouter.Ccurl(), nil).(interfaces.Datastore)
// 		spawnedResults[i] = jobs.Run(parentApiRouter, snapshot)
// 	}
// 	return spawnedResults
// }
