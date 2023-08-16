package execution

import (
	"errors"

	common "github.com/arcology-network/common-lib/common"

	// evmeu "github.com/arcology-network/vm-adaptor"
	ccurlcommon "github.com/arcology-network/concurrenturl/common"
	indexer "github.com/arcology-network/concurrenturl/indexer"
	"github.com/arcology-network/concurrenturl/interfaces"
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
func (this *Generation) Length() uint64       { return uint64(len(this.jobs)) }
func (this *Generation) Jobs() []*JobSequence { return this.jobs }

func (this *Generation) At(idx uint64) *JobSequence {
	return common.IfThenDo1st(idx < uint64(len(this.jobs)), func() *JobSequence { return this.jobs[idx] }, nil)
}

func (this *Generation) Add(job *JobSequence) bool {
	this.jobs = append(this.jobs, job)
	return true
}

func (this *Generation) Run(parentApiRouter eucommon.EthApiRouter) []interfaces.Univalue {
	config := NewConfig().SetCoinbase(parentApiRouter.Coinbase())

	// t0 := time.Now()
	worker := func(start, end, idx int, args ...interface{}) {
		for i := start; i < end; i++ {
			this.jobs[i].Results = this.jobs[i].Run(config, parentApiRouter)
		}
	}
	common.ParallelWorker(len(this.jobs), int(this.numThreads), worker)
	// fmt.Println(time.Since(t0))

	_, groupDict, _ := JobSequences(this.jobs).Detect().ToDict()
	JobSequences(this.jobs).ProcessConflicts(groupDict, errors.New(ccurlcommon.WARN_ACCESS_CONFLICT))

	//no filtering here !!!

	transitions := common.ConcateDo(this.jobs,
		func(v *JobSequence) uint64 {
			return uint64(len((*v).TransitionBuffer))
		},

		func(v *JobSequence) []interfaces.Univalue {
			return (*v).TransitionBuffer
		},
	)
	return indexer.Univalues(transitions).To(indexer.ITCTransition{})
}

func (this *Generation) Clear() uint64 {
	length := len(this.jobs)
	this.jobs = this.jobs[:0]
	return uint64(length)
}
