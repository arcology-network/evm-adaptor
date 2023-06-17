package execution

import (
	"fmt"

	"github.com/arcology-network/common-lib/codec"
	common "github.com/arcology-network/common-lib/common"
	indexer "github.com/arcology-network/concurrenturl/indexer"
	ccinterfaces "github.com/arcology-network/concurrenturl/interfaces"
	eucommon "github.com/arcology-network/vm-adaptor/common"
)

// APIs under the concurrency namespace
type Jobs struct {
	id              int
	maxThreads      uint8
	parentApiRouter eucommon.EthApiRouter
	jobs            []*Job // para jobs
}

func NewJobs(id int, maxThreads uint8, ethApiRouter eucommon.EthApiRouter, jobs []*Job) *Jobs {
	return &Jobs{
		id:              id,
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
		id:              id,
		maxThreads:      maxThreads,
		parentApiRouter: ethApiRouter,
		jobs:            make([]*Job, len(sequence.Msgs)),
		// results:         []*Result{},
	}

	for i, msg := range sequence.Msgs {
		this.jobs[i] = NewJobFromNative(
			i,
			this.Prefix(), // Parent prefix for uid generation
			msg.Native,
			ethApiRouter,
		)
	}
	return this
}

func (this *Jobs) Prefix() []byte {
	hash := this.parentApiRouter.TxHash()
	return append(hash[:], codec.Uint32(this.id).Encode()...)
}

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

	// this.results = make([]*Result, len(this.jobs))
	for i := 0; i < len(this.jobs); i++ {
		this.jobs[i].Result = this.jobs[i].Run(this.parentApiRouter.Coinbase(), snapshotUrl)
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
	resultDict := (&ResultDict{}).Categorize(results)
	subJobs := make([]*Jobs, len(resultDict))
	for i := 0; i < len(resultDict); i++ {
		seq := Results(resultDict[i]).ToSequence()
		subJobs[i] = NewJobsFromSequence(i, this.maxThreads, this.parentApiRouter, seq)
	}
	return common.Remove(&subJobs, nil), resultDict
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
