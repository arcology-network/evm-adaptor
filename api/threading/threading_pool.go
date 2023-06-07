package threading

import (
	"crypto/sha256"
	"errors"
	"math/big"
	"strings"

	"github.com/arcology-network/common-lib/codec"
	common "github.com/arcology-network/common-lib/common"
	concurrenturl "github.com/arcology-network/concurrenturl"
	arbitrator "github.com/arcology-network/concurrenturl/arbitrator"
	indexer "github.com/arcology-network/concurrenturl/indexer"
	ccinterfaces "github.com/arcology-network/concurrenturl/interfaces"
	evmcommon "github.com/arcology-network/evm/common"
	evmcoretypes "github.com/arcology-network/evm/core/types"
	interfaces "github.com/arcology-network/vm-adaptor/interfaces"

	cceu "github.com/arcology-network/vm-adaptor"

	ccurlinterfaces "github.com/arcology-network/concurrenturl/interfaces"
	eucommon "github.com/arcology-network/vm-adaptor/common"
	"github.com/arcology-network/vm-adaptor/eth"
	types "github.com/arcology-network/vm-adaptor/types"
)

// APIs under the concurrency namespace
type ThreadingPool struct {
	numThreads uint8
	jobs       []*types.Job // para jobs
}

func NewJobPool(numThreads uint8) *ThreadingPool {
	return &ThreadingPool{
		numThreads: numThreads,
		jobs:       []*types.Job{},
	}
}

func (this *ThreadingPool) Length() uint64 { return uint64(len(this.jobs)) }

func (this *ThreadingPool) At(idx uint64) *types.Job {
	return common.IfThenDo1st(idx < uint64(len(this.jobs)), func() *types.Job { return this.jobs[idx] }, nil)
}

func (this *ThreadingPool) Add(origin, calleeAddr evmcommon.Address, funCallData []byte) bool {
	if uint32(len(this.jobs)) >= eucommon.MAX_NUMBER_THREADS {
		return false
	}

	msg := evmcoretypes.NewMessage( // Build the message
		origin,
		&calleeAddr,
		0,
		new(big.Int).SetUint64(0), // Amount to transfer
		1e15,
		new(big.Int).SetUint64(1),
		funCallData,
		nil,
		false, // Don't checking nonce
	)

	this.jobs = append(this.jobs, &types.Job{Message: &msg})
	return true
}

func (this *ThreadingPool) ExportWriteCaches(jobs []*types.Job) []ccinterfaces.Univalue {
	infoVec := make([]ccinterfaces.Univalue, 0, len(jobs)*10) // Pre-allocation
	common.Foreach(jobs, func(job **types.Job) {
		if (**job).IsSuccessful() {
			infoVec = append(infoVec, (*(*job)).ApiRounter.Ccurl().Export()...)
		}
	})
	return infoVec
}

func (this *ThreadingPool) Run(parentApiRouter interfaces.EthApiRouter) bool {
	snapshotUrl := parentApiRouter.Ccurl().Snapshot()
	config := cceu.NewConfig().SetCoinbase(parentApiRouter.Coinbase()) // Share the same coinbase as the main thread

	for i := 0; i < len(this.jobs); i++ {
		ccurl := (&concurrenturl.ConcurrentUrl{}).New(
			indexer.NewWriteCache(snapshotUrl, parentApiRouter.Ccurl().Platform),
			parentApiRouter.Ccurl().Platform) // Init a write cache only since it doesn't need the importers

		// parentTxHash := parentApiRouter.TxHash()
		txHash := sha256.Sum256(append(codec.Bytes32(parentApiRouter.TxHash()).Encode(), codec.Uint64(i).Encode()...)) // A temp tx number, which will be replaced later
		this.jobs[i].ApiRounter = parentApiRouter.New(txHash, uint32(i), parentApiRouter.Depth(), ccurl)

		statedb := eth.NewImplStateDB(this.jobs[i].ApiRounter) // Eth state DB
		statedb.Prepare(txHash, [32]byte{}, i)                 // tx hash , block hash and tx index
		this.jobs[i].Run(config, statedb)
	}

	accesseVec := indexer.Univalues(this.ExportWriteCaches(this.jobs)).To(indexer.IPCAccess{})
	dict := arbitrator.Conflicts((&arbitrator.Arbitrator{}).Detect(accesseVec)).ToDict() // Detect potential conflicts
	common.Foreach(this.jobs, func(job **types.Job) {
		if _, conflict := (dict)[(**job).ApiRounter.TxIndex()]; conflict { // Label conflicts
			(**job).Err = errors.New("Error: Conflict in State accesses")
		}
	})

	executionResults := common.CopyIfDo(this.jobs,
		func(v *types.Job) bool { return true },
		func(v *types.Job) *types.ExecutionResult { return (*v).GetExecutionResult() })

	this.WriteBack(executionResults, parentApiRouter) // Merge back to the main write cache
	return true
}

// Merge all the transitions back to the main cache
func (this *ThreadingPool) WriteBack(executionResults []*types.ExecutionResult, targetApiRouter interfaces.EthApiRouter) {
	transitions := common.Concate(executionResults,
		func(v *types.ExecutionResult) []ccurlinterfaces.Univalue { return v.Transitions })

	common.RemoveIf(&transitions, func(v ccinterfaces.Univalue) bool {
		return strings.HasSuffix(*v.GetPath(), "/nonce") || common.IsPath(*v.GetPath()) // paths will be created as the elements inserted, but wow about empty paths
	})

	newPathTrans := common.MoveIf(&transitions, func(v ccinterfaces.Univalue) bool {
		return common.IsPath(*v.GetPath()) && !v.Preexist()
	})

	common.Foreach(newPathTrans, func(v *ccinterfaces.Univalue) {
		(*v).SetTx(targetApiRouter.TxIndex())              // use the parent tx index instead
		(*v).WriteTo(targetApiRouter.Ccurl().WriteCache()) // Write back to the parent writecache
	})

	common.Foreach(transitions, func(v *ccinterfaces.Univalue) {
		(*v).SetTx(targetApiRouter.TxIndex())              // use the parent tx index instead
		(*v).WriteTo(targetApiRouter.Ccurl().WriteCache()) // Write back to the parent writecache
	})
}

func (this *ThreadingPool) Clear() uint64 {
	length := len(this.jobs)
	this.jobs = this.jobs[:0]
	return uint64(length)
}
