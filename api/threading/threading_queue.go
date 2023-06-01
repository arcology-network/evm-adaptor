package threading

import (
	"crypto/sha256"
	"math/big"
	"strings"

	"github.com/arcology-network/common-lib/codec"
	common "github.com/arcology-network/common-lib/common"
	concurrenturl "github.com/arcology-network/concurrenturl"
	arbitrator "github.com/arcology-network/concurrenturl/arbitrator"
	indexer "github.com/arcology-network/concurrenturl/indexer"
	ccinterfaces "github.com/arcology-network/concurrenturl/interfaces"
	evmcommon "github.com/arcology-network/evm/common"
	types "github.com/arcology-network/evm/core/types"
	interfaces "github.com/arcology-network/vm-adaptor/interfaces"

	cceu "github.com/arcology-network/vm-adaptor"

	eucommon "github.com/arcology-network/vm-adaptor/common"
	"github.com/arcology-network/vm-adaptor/eth"
)

// APIs under the concurrency namespace
type Queue struct {
	numThreads uint8
	jobs       []*Job
	arbitrator *arbitrator.Arbitrator
}

func NewJobQueue(numThreads uint8) *Queue {
	return &Queue{
		numThreads: numThreads,
		jobs:       []*Job{},
		arbitrator: &arbitrator.Arbitrator{},
	}
}

func (this *Queue) Length() uint64 { return uint64(len(this.jobs)) }

func (this *Queue) At(idx uint64) *Job {
	return common.IfThenDo1st(idx < uint64(len(this.jobs)), func() *Job { return this.jobs[idx] }, nil)
}

func (this *Queue) Del(idx uint64) {
	common.IfThenDo(idx < uint64(len(this.jobs)), func() { this.jobs[idx] = nil }, func() {})
	common.RemoveIf(&this.jobs, func(job *Job) bool { return job == nil })
}

func (this *Queue) Add(origin, calleeAddr evmcommon.Address, funCall []byte) bool {
	if uint32(len(this.jobs)) >= eucommon.MAX_NUMBER_THREADS {
		return false
	}

	this.jobs = append(this.jobs,
		&Job{
			sender: origin,
			caller: evmcommon.Address{},
			callee: calleeAddr,
			message: types.NewMessage( // Build the message
				origin,
				&calleeAddr,
				0,
				new(big.Int).SetUint64(0), // Amount to transfer
				1e15,
				new(big.Int).SetUint64(1),
				funCall,
				nil,
				false, // Don't checking nonce
			),
		},
	)
	return true
}

func (this *Queue) ExportWriteCaches(jobs []*Job) []ccinterfaces.Univalue {
	infoVec := make([]ccinterfaces.Univalue, 0, len(jobs)*10) // Pre-allocation
	common.Foreach(jobs, func(job **Job) {
		if (**job).prechkErr == nil && (**job).receipt.Status == 1 {
			infoVec = append(infoVec, (*(*job)).apiRounter.Ccurl().Export()...)
		}
	})
	return infoVec
}

func (this *Queue) Run(parentApiRouter interfaces.ApiRouter) bool {
	snapshotUrl := parentApiRouter.Ccurl().Snapshot()
	// snapshotUrl := this.snapshot(parentApiRouter)
	config := cceu.NewConfig().SetCoinbase(parentApiRouter.Coinbase()) // Share the same coinbase as the main thread

	// t0 := time.Now()
	// executor := func(start, end, index int, args ...interface{}) {
	// for i := start; i < end; i++ {
	for i := 0; i < len(this.jobs); i++ {
		ccurl := (&concurrenturl.ConcurrentUrl{}).New(
			indexer.NewWriteCache(snapshotUrl, parentApiRouter.Ccurl().Platform),
			parentApiRouter.Ccurl().Platform) // Init a write cache only since it doesn't need the importers

		// parentTxHash := parentApiRouter.TxHash()
		txHash := sha256.Sum256(append(codec.Bytes32(parentApiRouter.TxHash()).Encode(), codec.Uint64(i).Encode()...)) // A temp tx number, which will be replaced later
		this.jobs[i].apiRounter = parentApiRouter.New(txHash, uint32(i), parentApiRouter.Depth(), ccurl)

		statedb := eth.NewImplStateDB(this.jobs[i].apiRounter) // Eth state DB
		statedb.Prepare(txHash, [32]byte{}, i)                 // tx hash , block hash and tx index
		this.jobs[i].Run(config, statedb)
	}
	// }
	// common.ParallelWorker(len(this.jobs), int(this.numThreads), executor)
	// fmt.Println("Run: ", time.Since(t0))

	// t0 = time.Now()
	accesseVec := indexer.Univalues(this.ExportWriteCaches(this.jobs)).To(indexer.IPCAccess{})
	tx := (&arbitrator.Arbitrator{}).Detect(accesseVec) // Detect potential conflicts

	this.WriteBack(arbitrator.Conflicts(tx).TxIDs(), parentApiRouter, this.jobs) // Merge back to the main write cache
	// fmt.Println("Commit: ", time.Since(t0))
	return true
}

// Merge all the transitions back to the main cache
func (this *Queue) WriteBack(conflicTxs []uint32, parentApiRouter interfaces.ApiRouter, jobs []*Job) {
	dict := common.MapFromArray(conflicTxs, true)
	common.Foreach(this.jobs, func(job **Job) {
		_, (**job).hasConflict = (*dict)[(**job).apiRounter.TxIndex()] // Label conflicts
	})

	transitionVec := this.ExportWriteCaches(this.jobs)
	transitions := []ccinterfaces.Univalue(indexer.Univalues(transitionVec).To(indexer.IPCTransition{}))

	common.RemoveIf(&transitions, func(v ccinterfaces.Univalue) bool {
		return strings.HasSuffix(*v.GetPath(), "/nonce") || common.IsPath(*v.GetPath()) // paths will be created as the elements inserted, but wow about empty paths
	})

	newPathTrans := common.MoveIf(&transitions, func(v ccinterfaces.Univalue) bool {
		return common.IsPath(*v.GetPath()) && !v.Preexist()
	})

	common.Foreach(newPathTrans, func(v *ccinterfaces.Univalue) {
		(*v).SetTx(parentApiRouter.TxIndex())              // use the parent tx index instead
		(*v).WriteTo(parentApiRouter.Ccurl().WriteCache()) // Write back to the parent writecache
	})

	common.Foreach(transitions, func(v *ccinterfaces.Univalue) {
		(*v).SetTx(parentApiRouter.TxIndex())              // use the parent tx index instead
		(*v).WriteTo(parentApiRouter.Ccurl().WriteCache()) // Write back to the parent writecache
	})
}

func (this *Queue) Clear() uint64 {
	length := len(this.jobs)
	this.jobs = this.jobs[:0]
	return uint64(length)
}
