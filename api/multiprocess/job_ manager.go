package multiprocess

import (
	"errors"
	"math/big"

	common "github.com/arcology-network/common-lib/common"
	concurrenturl "github.com/arcology-network/concurrenturl"
	ccurlcommon "github.com/arcology-network/concurrenturl/common"
	indexer "github.com/arcology-network/concurrenturl/indexer"
	ccurlstorage "github.com/arcology-network/concurrenturl/storage"
	evmcommon "github.com/arcology-network/evm/common"
	"github.com/arcology-network/evm/core"
	types "github.com/arcology-network/evm/core/types"
	"github.com/arcology-network/evm/core/vm"
	"github.com/arcology-network/evm/params"

	cceu "github.com/arcology-network/vm-adaptor"
	eucommon "github.com/arcology-network/vm-adaptor/common"
	"github.com/arcology-network/vm-adaptor/eth"
)

type Job struct {
	sender      evmcommon.Address
	caller      evmcommon.Address
	callee      evmcommon.Address
	message     types.Message
	receipt     *types.Receipt
	result      *core.ExecutionResult
	prechkErr   error
	accesses    []ccurlcommon.UnivalueInterface
	transitions []ccurlcommon.UnivalueInterface
	ccurl       *concurrenturl.ConcurrentUrl
}

// APIs under the concurrency namespace
type JobManager struct {
	jobs       []Job
	threads    int
	apiRouter  eucommon.ConcurrentApiRouterInterface
	arbitrator *indexer.ArbitratorSlow
}

func NewJobManager(apiRouter eucommon.ConcurrentApiRouterInterface) *JobManager {
	return &JobManager{
		jobs:       []Job{},
		apiRouter:  apiRouter,
		threads:    16, // 16 threads by default
		arbitrator: indexer.NewArbitratorSlow(),
	}
}

func (this *JobManager) Length() uint64 { return uint64(len(this.jobs)) }

func (this *JobManager) At(idx uint64) ([]byte, error) {
	if idx >= uint64(len(this.jobs)) {
		return []byte{}, errors.New("Access out of range")
	}

	if this.jobs[idx].result != nil {
		return this.jobs[idx].result.ReturnData, this.jobs[idx].result.Err
	}
	return []byte{}, this.jobs[idx].prechkErr
}

func (this *JobManager) Snapshot(current *concurrenturl.ConcurrentUrl) (*concurrenturl.ConcurrentUrl, error) {
	_, transitions := current.ExportAll()

	snapshotUrl := concurrenturl.NewConcurrentUrl(ccurlstorage.NewTransientDB(*(current.Store())))
	snapshotUrl.Import(transitions)
	snapshotUrl.PostImport()

	if errs := snapshotUrl.Commit(nil); len(errs) != 0 { // Commit all
		return nil, errors.New("Error: Failed to import transitions")
	}
	return snapshotUrl, nil
}

func (this *JobManager) Add(calleeAddr evmcommon.Address, funCall []byte) int {
	sender := this.apiRouter.From()
	this.jobs = append(this.jobs,
		Job{
			sender: sender,
			caller: evmcommon.Address{},
			callee: calleeAddr,
			message: types.NewMessage( // Build the message
				sender,
				&calleeAddr,
				0,
				new(big.Int).SetUint64(0), // Amount to transfer
				1e15,
				new(big.Int).SetUint64(1),
				funCall, //Need to a deepcopy
				nil,
				false, // Stop checking nonce
			),
		},
	)
	return len(this.jobs) - 1
}

func (this *JobManager) Run() {
	snapshotUrl, err := this.Snapshot(this.apiRouter.Ccurl())
	if err != nil {
		return
	}

	processor := func(start, end, index int, args ...interface{}) {
		for i := 0; i < len(this.jobs); i++ {
			// for i := start; i < end; i++ {
			this.jobs[i].ccurl = concurrenturl.NewConcurrentUrl(ccurlstorage.NewTransientDB(*snapshotUrl.Store()))
			statedb := eth.NewImplStateDB(this.jobs[i].ccurl) // Eth state DB
			statedb.Prepare([32]byte{}, [32]byte{}, i)        // tx hash , block hash and tx index

			eu := cceu.NewEU(
				params.MainnetChainConfig,
				vm.Config{},
				statedb,
				this.apiRouter.New(evmcommon.Hash{}, 0, this.apiRouter.Ccurl()), // Call the function
			)

			config := cceu.NewConfig()
			this.jobs[i].accesses, this.jobs[i].transitions, this.jobs[i].receipt, this.jobs[i].result, this.jobs[i].prechkErr =
				eu.Run(evmcommon.Hash{}, i, &this.jobs[i].message, cceu.NewEVMBlockContext(config), cceu.NewEVMTxContext(this.jobs[i].message))
		}
	}
	common.ParallelWorker(len(this.jobs), this.threads, processor)

	// Detect potential conflicts
	arbitrator := indexer.NewArbitratorSlow()
	// accesses := common.ConcateFrom(this.jobs, func(v Job) []ccurlcommon.UnivalueInterface { return v.accesses })

	accesseVec := []ccurlcommon.UnivalueInterface{}
	common.Foreach(this.jobs, func(job *Job) { accesseVec = append(accesseVec, job.accesses...) })
	arbitrator.Detect(accesseVec)
}

// Merge all the transitions
func (this *JobManager) Commit(getter func(v interface{}) []ccurlcommon.UnivalueInterface) {
	for i := 0; i < len(this.jobs); i++ {
		this.apiRouter.Ccurl().Indexer().MergeFrom(this.jobs[i].ccurl.Indexer())
	}
}

func (this *JobManager) Clear() uint64 {
	length := len(this.jobs)
	this.jobs = this.jobs[:0]
	return uint64(length)
}

// func MainTestConfig() *cceu.Config {
// 	vmConfig := vm.Config{}
// 	cfg := &cceu.Config{
// 		ChainConfig: params.MainnetChainConfig,
// 		VMConfig:    &vmConfig,
// 		BlockNumber: big.NewInt(0),
// 		ParentHash:  evmcommon.Hash{},
// 		Time:        big.NewInt(0),
// 		Coinbase:    &eucommon.Coinbase,
// 		GasLimit:    math.MaxUint64, // Should come from the message
// 		Difficulty:  big.NewInt(0),
// 	}
// 	cfg.Chain = new(cceu.DummyChain)
// 	return cfg
// }

// func NewTestEU() (*cceu.EU, *cceu.Config, ccurlcommon.DatastoreInterface, *concurrenturl.ConcurrentUrl) {
// 	persistentDB := cachedstorage.NewDataStore()
// 	persistentDB.Inject((&concurrenturl.Platform{}).Eth10Account(), commutative.NewPath())
// 	db := ccurlstorage.NewTransientDB(persistentDB)

// 	url := concurrenturl.NewConcurrentUrl(db)
// 	statedb := eth.NewImplStateDB(url)
// 	statedb.Prepare(evmcommon.Hash{}, evmcommon.Hash{}, 0)
// 	statedb.CreateAccount(eucommon.Coinbase)
// 	statedb.CreateAccount(eucommon.User1)
// 	statedb.AddBalance(eucommon.User1, new(big.Int).SetUint64(1e18))
// 	_, transitions := url.ExportAll()
// 	fmt.Println("\n" + eucommon.FormatTransitions(transitions))

// 	// Deploy.
// 	url = concurrenturl.NewConcurrentUrl(db)
// 	url.Import(transitions)
// 	url.PostImport()
// 	url.Commit([]uint32{0})
// 	api := ccapi.NewAPI(url)
// 	statedb = eth.NewImplStateDB(url)

// 	config := MainTestConfig()
// 	config.Coinbase = &eucommon.Coinbase
// 	config.BlockNumber = new(big.Int).SetUint64(10000000)
// 	config.Time = new(big.Int).SetUint64(10000000)

// 	return cceu.NewEU(config.ChainConfig, *config.VMConfig, statedb, api), config, db, url
// }
