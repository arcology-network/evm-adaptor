package execution

import (
	"fmt"

	"github.com/arcology-network/common-lib/codec"
	common "github.com/arcology-network/common-lib/common"
	indexer "github.com/arcology-network/concurrenturl/indexer"
	eucommon "github.com/arcology-network/vm-adaptor/common"
)

// APIs under the concurrency namespace
type Jobs struct {
	id              int
	numThreads      uint8
	parentApiRouter eucommon.EthApiRouter
	jobs            []*Job // para jobs
	results         []*Result
}

func NewJobs(id int, numThreads uint8, ethApiRouter eucommon.EthApiRouter, jobs []*Job) *Jobs {
	return &Jobs{
		id:              id,
		numThreads:      numThreads,
		parentApiRouter: ethApiRouter,
		jobs:            jobs,
	}
}

func NewJobsFromSequence(id int, numThreads uint8, ethApiRouter eucommon.EthApiRouter, sequence *Sequence) *Jobs {
	if sequence == nil {
		return nil
	}

	this := &Jobs{
		id:              id,
		numThreads:      numThreads,
		parentApiRouter: ethApiRouter,
		jobs:            make([]*Job, len(sequence.Msgs)),
		results:         []*Result{},
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

func (this *Jobs) Run() []*Result {
	snapshotUrl := this.parentApiRouter.Ccurl().Snapshot()
	this.results = make([]*Result, len(this.jobs))
	for i := 0; i < len(this.jobs); i++ {
		this.results[i] = this.jobs[i].Run(this.parentApiRouter.Coinbase(), snapshotUrl)
	}

	if len(this.results) > 1 {
		this.results, _ = Results(this.results).DetectConflict() // Detect potential conflicts
	}

	// Extract deferred calls if exist
	ResultDict := (&ResultDict{}).Categorize(this.results)
	subJobs := make([]*Jobs, len(ResultDict))
	for i := 0; i < len(ResultDict); i++ {
		seq := Results(ResultDict[i]).ToSequence()
		subJobs[i] = NewJobsFromSequence(i, this.numThreads, this.parentApiRouter, seq)
	}
	subJobs = common.Remove(&subJobs, nil)

	// Run deferrred jobs
	subResults := make([][]*Result, len(subJobs))
	for i, jobs := range subJobs {
		subResults[i] = jobs.Run() // need to transfer the funds into a temp account
	}

	fmt.Println("Sub subResults 1 === ====================== ====================== ====================== =========================================")
	if len(subResults) > 0 {
		indexer.Univalues(Results((subResults[0])).Transitions()).SortByDefault().Print()
	}

	fmt.Println("Sub subResults 2 === ====================== ====================== ====================== =========================================")
	if len(subResults) > 1 {
		indexer.Univalues(Results((subResults[1])).Transitions()).SortByDefault().Print()
	}

	// Write the transitions back to the parent write cache
	catenated := append(this.results, common.Flatten(subResults)...)
	common.Foreach(catenated, func(v **Result) {
		(*v).WriteTo(this.parentApiRouter.TxIndex(), this.parentApiRouter.Ccurl().WriteCache()) // Merge the write cache to its parent
	})

	// fmt.Println("Sub Catenated === ====================== ====================== ====================== =========================================")
	// indexer.Univalues(Results(common.Flatten(subResults)).Transitions()).SortByDefault().Print()

	// fmt.Println("Current Catenated === ====================== ====================== ====================== =========================================")
	// indexer.Univalues(Results(this.results).Transitions()).SortByDefault().Print()

	return catenated
}

func (this *Jobs) Clear() uint64 {
	length := len(this.jobs)
	this.jobs = this.jobs[:0]
	return uint64(length)
}
