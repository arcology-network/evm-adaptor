package multiprocess

import (
	"crypto/sha256"
	"math/big"
	"strings"

	"github.com/arcology-network/common-lib/codec"
	common "github.com/arcology-network/common-lib/common"
	concurrenturl "github.com/arcology-network/concurrenturl"
	ccurlcommon "github.com/arcology-network/concurrenturl/common"
	indexer "github.com/arcology-network/concurrenturl/indexer"
	ccurlstorage "github.com/arcology-network/concurrenturl/storage"
	"github.com/arcology-network/concurrenturl/univalue"
	evmcommon "github.com/arcology-network/evm/common"
	types "github.com/arcology-network/evm/core/types"
	"github.com/arcology-network/evm/core/vm"
	"github.com/arcology-network/evm/params"

	cceu "github.com/arcology-network/vm-adaptor"

	apicommon "github.com/arcology-network/vm-adaptor/api/common"
	eucommon "github.com/arcology-network/vm-adaptor/common"
	"github.com/arcology-network/vm-adaptor/eth"
)

// APIs under the concurrency namespace
type Queue struct {
	jobs       []*Job
	arbitrator *indexer.Arbitrator
}

func NewJobQueue() *Queue {
	return &Queue{
		jobs:       []*Job{},
		arbitrator: &indexer.Arbitrator{},
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

func (this *Queue) FilteredAccesses() []ccurlcommon.UnivalueInterface {
	accesseVec := make([]ccurlcommon.UnivalueInterface, 0, len(this.jobs)*10) // Pre-allocation
	common.Foreach(this.jobs, func(job **Job) {
		if (**job).prechkErr == nil && (**job).receipt.Status == 1 {
			accesseVec = append(accesseVec, (*(*job)).FilteredAccesses()...)
		}
	})
	univalue.Univalues(accesseVec).Print()
	return accesseVec
}

func (this *Queue) snapshot(mainApiRouter eucommon.ConcurrentApiRouterInterface) ccurlcommon.DatastoreInterface {
	transitions := mainApiRouter.Ccurl().Export() // Get the all up-to-date transitions from the main thread
	univalue.Univalues(transitions).Print()

	mainProcessTrans := univalue.Univalues(common.Clone(transitions)).To(
		univalue.RemoveReadOnly,
		univalue.DelNonExist,
		univalue.CloneValue,
	)

	univalue.Univalues(mainProcessTrans).Print()

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
		eu := cceu.NewEU(
			params.MainnetChainConfig,
			vm.Config{},
			statedb,
			this.jobs[i].apiRounter, // Tx hash, tx id and url
		)

		this.jobs[i].receipt, this.jobs[i].result, this.jobs[i].prechkErr =
			eu.Run(eu.Api().TxHash(), int(eu.Api().TxIndex()), &this.jobs[i].message, cceu.NewEVMBlockContext(config), cceu.NewEVMTxContext(this.jobs[i].message))
	}
	// }
	// common.ParallelWorker(len(this.jobs), int(threads), executor)
	// fmt.Println("Run: ", time.Since(t0))

	// t0 = time.Now()
	tx := (&indexer.Arbitrator{}).Detect(this.FilteredAccesses()) // Detect potential conflicts
	this.LableJobs(indexer.Conflicts(tx).TxIDs())
	this.WriteBack(parentApiRouter, this.jobs) // Merge back to the main write cache
	// fmt.Println("Commit: ", time.Since(t0))
	return true
}

// Merge all the transitions back to the main cache
func (this *Queue) LableJobs(conflicTxs []uint32) {
	dict := common.MapFromArray(conflicTxs, true)
	common.Foreach(this.jobs, func(job **Job) {
		_, (**job).hasConflict = (*dict)[(**job).apiRounter.TxIndex()] // Label conflicts
	})
}

// Merge all the transitions back to the main cache
func (this *Queue) WriteBack(parentApiRouter eucommon.ConcurrentApiRouterInterface, jobs []*Job) {
	for i := 0; i < len(jobs); i++ { // transitt
		transitions := this.jobs[i].FilteredTransitions()
		common.RemoveIf(&transitions, func(v ccurlcommon.UnivalueInterface) bool {
			return strings.HasSuffix(*v.GetPath(), "/nonce") || ccurlcommon.IsPath(*v.GetPath())
		})

		common.Foreach(transitions, func(v *ccurlcommon.UnivalueInterface) {
			(*v).SetTx(parentApiRouter.TxIndex())
			(*v).WriteTo(parentApiRouter.Ccurl().WriteCache()) // Write the path creation first
		})
	}
}

func (this *Queue) Clear() uint64 {
	length := len(this.jobs)
	this.jobs = this.jobs[:0]
	return uint64(length)
}
