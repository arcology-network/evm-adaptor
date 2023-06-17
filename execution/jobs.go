package execution

import (
	"fmt"

	common "github.com/arcology-network/common-lib/common"
	indexer "github.com/arcology-network/concurrenturl/indexer"
	ccinterfaces "github.com/arcology-network/concurrenturl/interfaces"
	cceu "github.com/arcology-network/vm-adaptor"
	eucommon "github.com/arcology-network/vm-adaptor/common"
)

// APIs under the concurrency namespace
type Jobs struct {
	batchID         uint64
	maxThreads      uint8
	parentApiRouter eucommon.EthApiRouter
	jobs            []*Job // para jobs
}

func NewJobs(id int, maxThreads uint8, ethApiRouter eucommon.EthApiRouter, jobs []*Job) *Jobs {
	return &Jobs{
		batchID:         uint64(id),
		maxThreads:      maxThreads,
		parentApiRouter: ethApiRouter,
		jobs:            jobs,
	}
}

func NewJobsFromSequence(id int, maxThreads uint8, ethApiRouter eucommon.EthApiRouter, sequence *Sequence) *Jobs {
	if sequence == nil {
		return nil
	}

	this := &Jobs{
		batchID:         uint64(id),
		maxThreads:      maxThreads,
		parentApiRouter: ethApiRouter,
		jobs:            make([]*Job, len(sequence.Msgs)),
	}

	for i, msg := range sequence.Msgs {
		this.jobs[i] = NewJobFromNative(
			uint64(i),
			this.batchID, // batch ID
			msg.Native,
			ethApiRouter,
		)
	}
	return this
}

func (this *Jobs) Batch() uint64  { return uint64(this.batchID) }
func (this *Jobs) Length() uint64 { return uint64(len(this.jobs)) }

func (this *Jobs) At(idx uint64) *Job {
	return common.IfThenDo1st(idx < uint64(len(this.jobs)), func() *Job { return this.jobs[idx] }, nil)
}

func (this *Jobs) Add(job *Job) bool {
	if uint32(len(this.jobs)) > eucommon.MAX_NUMBER_THREADS {
		return false
	}

	this.jobs = append(this.jobs, job)
	return true
}

func (this *Jobs) Run(predecessors []*Result) []*Result {
	preTransitions := common.Concate(predecessors, func(v *Result) []ccinterfaces.Univalue { return v.Transitions })
	snapshotUrl := this.parentApiRouter.Ccurl().Snapshot(preTransitions)

	config := cceu.NewConfig().SetCoinbase(this.parentApiRouter.Coinbase())
	for i := 0; i < len(this.jobs); i++ {
		this.jobs[i].Result = this.jobs[i].Run(config, snapshotUrl)
	}

	results := common.Concate(this.jobs, func(job *Job) []*Result { return []*Result{job.Result} })
	if len(results) > 1 {
		results, _ = Results(results).DetectConflict() // Detect potential conflicts
	}

	// Run deferrred jobs
	subResults := this.RunSpawned(results)

	fmt.Println("Sub subResults 1 === ====================== ====================== ====================== =========================================")
	if len(subResults) > 0 {
		indexer.Univalues(Results((subResults[0])).Transitions()).SortByDefault().Print()
	}
	fmt.Println("Sub subResults 2 === ====================== ====================== ====================== =========================================")
	if len(subResults) > 1 {
		indexer.Univalues(Results((subResults[1])).Transitions()).SortByDefault().Print()
	}
	// Write the transitions back to the parent write cache
	catenated := append(results, common.Flatten(subResults)...)
	common.Foreach(catenated, func(v **Result) {
		(*v).WriteTo(this.parentApiRouter.TxIndex(), this.parentApiRouter.Ccurl().WriteCache()) // Merge the write cache to its parent
	})

	return catenated
}

func (this *Jobs) Clear() uint64 {
	length := len(this.jobs)
	this.jobs = this.jobs[:0]
	return uint64(length)
}

// Extract deferred calls if exist
func (this *Jobs) GetSpawned(results []*Result) ([]*Jobs, [][]*Result) {
	resultvec := common.GroupBy(results, func(v *Result) *[32]byte {
		return common.IfThenDo1st(v.Spawned != nil, func() *[32]byte { return &(v.Spawned.CallSig) }, nil)
	})

	subJobs := make([]*Jobs, len(resultvec))
	for i := 0; i < len(resultvec); i++ {
		seq := Results(resultvec[i]).ToSequence()
		subJobs[i] = NewJobsFromSequence(i, this.maxThreads, this.parentApiRouter, seq)
	}
	return common.Remove(&subJobs, nil), resultvec
}

func (this *Jobs) RunSpawned(results []*Result) [][]*Result {
	spawnedJobs, preResults := this.GetSpawned(results)
	if len(spawnedJobs) == 0 {
		return [][]*Result{}
	}

	// need to transfer the funds to a temp account
	spawnedResults := make([][]*Result, len(spawnedJobs))
	for i, jobs := range spawnedJobs {
		spawnedResults[i] = jobs.Run(preResults[i])
	}
	return spawnedResults
}
