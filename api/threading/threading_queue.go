package multiprocess

import (
	"crypto/sha256"
	"math/big"
	"strings"

	"github.com/arcology-network/common-lib/codec"
	common "github.com/arcology-network/common-lib/common"
	concurrenturl "github.com/arcology-network/concurrenturl"
	arbitrator "github.com/arcology-network/concurrenturl/arbitrator"
	indexer "github.com/arcology-network/concurrenturl/indexer"
	interfaces "github.com/arcology-network/concurrenturl/interfaces"
	ccurlstorage "github.com/arcology-network/concurrenturl/storage"
	evmcommon "github.com/arcology-network/evm/common"
	types "github.com/arcology-network/evm/core/types"

	cceu "github.com/arcology-network/vm-adaptor"

	apicommon "github.com/arcology-network/vm-adaptor/api/common"
	eucommon "github.com/arcology-network/vm-adaptor/common"
	"github.com/arcology-network/vm-adaptor/eth"
)

// APIs under the concurrency namespace
type Queue struct {
	jobs       []*Job
	arbitrator *arbitrator.Arbitrator
}

func NewJobQueue() *Queue {
	return &Queue{
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

func (this *Queue) Add(origin, calleeAddr evmcommon.Address, funCall []byte) int {
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
	return len(this.jobs) - 1
}

func (this *Queue) ExportWriteCaches(jobs []*Job) []interfaces.Univalue {
	infoVec := make([]interfaces.Univalue, 0, len(jobs)*10) // Pre-allocation
	common.Foreach(jobs, func(job **Job) {
		if (**job).prechkErr == nil && (**job).receipt.Status == 1 {
			infoVec = append(infoVec, (*(*job)).apiRounter.Ccurl().Export()...)
		}
	})
	return infoVec
}

func (this *Queue) FilterAccesses() []interfaces.Univalue {
	infoVec := this.ExportWriteCaches(this.jobs)

	accesseVec := indexer.Univalues(infoVec).To(indexer.IPCAccess{})
	indexer.Univalues(accesseVec).Print()
	return accesseVec
}

func (this *Queue) snapshot(mainApiRouter eucommon.ConcurrentApiRouterInterface) interfaces.Datastore {
	transitions := mainApiRouter.Ccurl().Export() // Get the all up-to-date transitions from the main thread
	indexer.Univalues(transitions).Print()

	mainProcessTrans := indexer.Univalues(common.Clone(transitions)).To(indexer.ITCTransition{})
	// indexer.Univalues(mainProcessTrans).Print()

	transientDB := ccurlstorage.NewTransientDB(mainApiRouter.Ccurl().WriteCache().Store()) // Should be the same as Importer().Store()
	snapshot := concurrenturl.NewConcurrentUrl(transientDB).Import(mainProcessTrans).Sort()
	return snapshot.Commit([]uint32{mainApiRouter.TxIndex()}).Importer().Store() // Commit these changes to the a transient DB
}

func (this *Queue) Run(threads uint8, parentApiRouter eucommon.ConcurrentApiRouterInterface) bool {
	if parentApiRouter.Depth() > apicommon.MAX_RECURSIION_DEPTH {
		return false //, errors.New("Error: Execeeds the max recursion depth")
	}
	snapshot := this.snapshot(parentApiRouter)
	config := cceu.NewConfig().SetCoinbase(parentApiRouter.Coinbase()) // Share the same coinbase as the main thread

	// t0 := time.Now()
	// executor := func(start, end, index int, args ...interface{}) {
	// for i := start; i < end; i++ {
	for i := 0; i < len(this.jobs); i++ {
		ccurl := (&concurrenturl.ConcurrentUrl{}).New(
			indexer.NewWriteCache(snapshot, parentApiRouter.Ccurl().Platform),
			parentApiRouter.Ccurl().Platform) // Init a write cache only since it doesn't need the importers

		txHash := sha256.Sum256(codec.Uint64(i).Encode()) // A temp tx number, which will be replaced later
		this.jobs[i].apiRounter = parentApiRouter.New(txHash, uint32(i), ccurl)

		statedb := eth.NewImplStateDB(this.jobs[i].apiRounter) // Eth state DB
		statedb.Prepare(txHash, [32]byte{}, i)                 // tx hash , block hash and tx index
		this.jobs[i].Run(config, statedb)
	}
	// }
	// common.ParallelWorker(len(this.jobs), int(threads), executor)
	// fmt.Println("Run: ", time.Since(t0))

	// t0 = time.Now()
	tx := (&arbitrator.Arbitrator{}).Detect(this.FilterAccesses()) // Detect potential conflicts

	this.LableConflicts(arbitrator.Conflicts(tx).TxIDs())
	this.WriteBack(parentApiRouter, this.jobs) // Merge back to the main write cache
	// fmt.Println("Commit: ", time.Since(t0))
	return true
}

// Merge all the transitions back to the main cache
func (this *Queue) LableConflicts(conflicTxs []uint32) {
	dict := common.MapFromArray(conflicTxs, true)
	common.Foreach(this.jobs, func(job **Job) {
		_, (**job).hasConflict = (*dict)[(**job).apiRounter.TxIndex()] // Label conflicts
	})
}

// Merge all the transitions back to the main cache
func (this *Queue) WriteBack(parentApiRouter eucommon.ConcurrentApiRouterInterface, jobs []*Job) {
	infoVec := this.ExportWriteCaches(this.jobs)
	transitions := []interfaces.Univalue(indexer.Univalues(infoVec).To(indexer.IPCTransition{}))

	common.RemoveIf(&transitions, func(v interfaces.Univalue) bool {
		return strings.HasSuffix(*v.GetPath(), "/nonce") || common.IsPath(*v.GetPath()) // paths will be created as the elements inserted, but wow about empty paths
	})

	common.Foreach(transitions, func(v *interfaces.Univalue) {
		(*v).SetTx(parentApiRouter.TxIndex())
		(*v).WriteTo(parentApiRouter.Ccurl().WriteCache()) // Write the path creation first
	})
}

func (this *Queue) Clear() uint64 {
	length := len(this.jobs)
	this.jobs = this.jobs[:0]
	return uint64(length)
}
