package execution

import (
	"errors"

	common "github.com/arcology-network/common-lib/common"
	"github.com/arcology-network/concurrenturl/interfaces"
	ccinterfaces "github.com/arcology-network/concurrenturl/interfaces"
	cceu "github.com/arcology-network/vm-adaptor"
	eucommon "github.com/arcology-network/vm-adaptor/common"
)

// APIs under the concurrency namespace
type ParallelJobs struct {
	branchID   uint64
	maxThreads uint8
	preTxs     []uint32
	jobs       []*Job // para jobs
}

func NewParallelJobs(id int, maxThreads uint8, ethApiRouter eucommon.EthApiRouter, jobs []*Job) *ParallelJobs {
	return &ParallelJobs{
		branchID:   uint64(id),
		maxThreads: maxThreads,
		jobs:       jobs,
	}
}

func NewParallelJobsFromSequence(id int, maxThreads uint8, ethApiRouter eucommon.EthApiRouter, sequence *Sequence) *ParallelJobs {
	if sequence == nil {
		return nil
	}

	this := &ParallelJobs{
		branchID:   uint64(id),
		maxThreads: maxThreads,
		jobs:       make([]*Job, len(sequence.Msgs)),
	}

	for i, stdMsg := range sequence.Msgs {
		this.jobs[i] = &Job{
			uint64(i),
			this.branchID,
			[]*StandardMessage{stdMsg},
			make([]*Result, 1),
			ethApiRouter,
		}
	}
	return this
}

func (this *ParallelJobs) BranchID() uint64 { return uint64(this.branchID) }
func (this *ParallelJobs) Length() uint64   { return uint64(len(this.jobs)) }

func (this *ParallelJobs) At(idx uint64) *Job {
	return common.IfThenDo1st(idx < uint64(len(this.jobs)), func() *Job { return this.jobs[idx] }, nil)
}

func (this *ParallelJobs) Add(job *Job) bool {
	if uint32(len(this.jobs)) > eucommon.MAX_NUMBER_THREADS {
		return false
	}

	this.jobs = append(this.jobs, job)
	return true
}

func (this *ParallelJobs) Run(parentApiRouter eucommon.EthApiRouter, snapshot ccinterfaces.Datastore) []*Result {
	config := cceu.NewConfig().SetCoinbase(parentApiRouter.Coinbase())

	common.ParallelForeach(this.jobs, this.maxThreads, func(job **Job) *Job {
		(**job).Results = (**job).Run(config, snapshot)
		return (*job)
	})

	// Detect potential conflicts
	results := common.Concate(this.jobs, func(job *Job) []*Result { return job.Results })
	dict := Results(results).Detect()
	for i := 0; i < len(results); i++ {
		if _, conflict := (*dict)[results[i].TxIndex]; conflict {
			results[i].Err = errors.New("Error: Conflicts detected in state accesses")
		}
	}

	// Run deferrred jobs
	subResults := this.RunSpawned(parentApiRouter, results)
	catenated := append(results, common.Flatten(subResults)...)
	common.Foreach(catenated, func(v **Result) {
		(*v).WriteTo(parentApiRouter.TxIndex(), parentApiRouter.Ccurl().WriteCache()) // Merge the write cache to its parent
	})

	// fmt.Println("Sub subResults 1 === ====================== ====================== ====================== =========================================")
	// if len(subResults) > 0 {
	// 	indexer.Univalues(Results((subResults[0])).Transitions()).SortByDefault().Print()
	// }
	// fmt.Println("Sub subResults 2 === ====================== ====================== ====================== =========================================")
	// if len(subResults) > 1 {
	// 	indexer.Univalues(Results((subResults[1])).Transitions()).SortByDefault().Print()
	// }
	// Write the transitions back to the parent write cache
	return catenated
}

func (this *ParallelJobs) Clear() uint64 {
	length := len(this.jobs)
	this.jobs = this.jobs[:0]
	return uint64(length)
}

func (this *ParallelJobs) RunSpawned(parentApiRouter eucommon.EthApiRouter, results []*Result) [][]*Result {
	grouped := common.GroupBy(results, func(v *Result) *[32]byte {
		return common.IfThenDo1st(v.Spawned != nil, func() *[32]byte { return &(v.Spawned.CallSig) }, nil)
	})

	spawnedJobs := make([]*ParallelJobs, 0, len(grouped))
	for i := 0; i < len(grouped); i++ {
		seq := Results(grouped[i]).ToSequence()
		if job := NewParallelJobsFromSequence(i, this.maxThreads, parentApiRouter, seq); job != nil {
			spawnedJobs = append(spawnedJobs, job)
		}
	}

	if len(spawnedJobs) == 0 {
		return [][]*Result{}
	}

	spawnedResults := make([][]*Result, len(spawnedJobs))
	for i, jobs := range spawnedJobs {
		preTransitions := common.Concate(grouped[i], func(v *Result) []ccinterfaces.Univalue { return v.Transitions })

		var maker eucommon.LocalSnapshotMaker
		maker.Import(preTransitions)
		snapshot := maker.Make(parentApiRouter.Ccurl(), nil).(interfaces.Datastore)
		spawnedResults[i] = jobs.Run(parentApiRouter, snapshot)
	}
	return spawnedResults
}
