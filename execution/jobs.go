package execution

import (
	common "github.com/arcology-network/common-lib/common"
	eucommon "github.com/arcology-network/vm-adaptor/common"
)

// APIs under the concurrency namespace
type Jobs struct {
	numThreads      uint8
	parentApiRouter eucommon.EthApiRouter
	jobs            []*Job // para jobs
	results         []*Result
}

func NewJobs(numThreads uint8, ethApiRouter eucommon.EthApiRouter, jobs []*Job) *Jobs {
	return &Jobs{
		numThreads:      numThreads,
		parentApiRouter: ethApiRouter,
		jobs:            jobs,
	}
}

func NewJobsFromSequence(numThreads uint8, ethApiRouter eucommon.EthApiRouter, sequence *Sequence) *Jobs {
	if sequence == nil {
		return nil
	}

	jobs := make([]*Job, len(sequence.Msgs))
	for i, msg := range sequence.Msgs {
		jobs[i] = &Job{
			Predecessors: sequence.Predecessors,
			Message:      msg.Native,
			ApiRouter:    ethApiRouter,
		}
	}

	return &Jobs{
		numThreads:      numThreads,
		parentApiRouter: ethApiRouter,
		jobs:            jobs,
		results:         []*Result{},
	}
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

func (this *Jobs) Run() []*Result {
	snapshotUrl := this.parentApiRouter.Ccurl().Snapshot()
	this.results = make([]*Result, len(this.jobs))
	for i := 0; i < len(this.jobs); i++ {
		this.results[i] = this.jobs[i].Run(uint32(i), this.parentApiRouter.Coinbase(), snapshotUrl)
	}

	if len(this.results) > 1 {
		this.results = Results(this.results).DetectConflict() // Detect potential conflicts
	}

	// Extract deferred calls if exist
	resultset := (&ResultSet{}).Categorize(this.results)
	subJobs := make([]*Jobs, len(resultset))
	for i := 0; i < len(resultset); i++ {
		seq := Results(resultset[i]).ToSequence()
		subJobs[i] = NewJobsFromSequence(this.numThreads, this.parentApiRouter, seq)
	}
	subJobs = common.Remove(&subJobs, nil)

	// Run deferrred jobs
	subResults := make([][]*Result, len(subJobs))
	for i, jobs := range subJobs {
		subResults[i] = jobs.Run() // need to transfer the funds into a temp account
	}

	// Write the transitions back to the parent write cache
	catenated := append(this.results, common.Flatten(subResults)...)
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
