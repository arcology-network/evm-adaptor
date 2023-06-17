package execution

import (
	"crypto/sha256"
	"math/big"

	"github.com/arcology-network/common-lib/codec"
	common "github.com/arcology-network/common-lib/common"
	"github.com/arcology-network/concurrenturl"
	"github.com/arcology-network/concurrenturl/commutative"
	"github.com/arcology-network/concurrenturl/indexer"
	cceu "github.com/arcology-network/vm-adaptor"

	ccurlinterfaces "github.com/arcology-network/concurrenturl/interfaces"
	evmcommon "github.com/arcology-network/evm/common"
	evmcore "github.com/arcology-network/evm/core"
	evmTypes "github.com/arcology-network/evm/core/types"
	"github.com/arcology-network/evm/core/vm"
	evmparams "github.com/arcology-network/evm/params"
	eucommon "github.com/arcology-network/vm-adaptor/common"
	"github.com/arcology-network/vm-adaptor/eth"
	"github.com/holiman/uint256"
)

type Job struct {
	ID           uint64
	TxHash       [32]byte
	Predecessors [][32]byte
	EvmMsg       *evmcore.Message
	ApiRouter    eucommon.EthApiRouter
	Receipt      *evmTypes.Receipt
	EvmResult    *evmcore.ExecutionResult
	Result       *Result
}

func NewJob(jobID, batchID uint64, from, to evmcommon.Address, funCallData []byte, gaslimit uint64, parentApiRouter eucommon.EthApiRouter) *Job {
	msg := evmcore.NewMessage( // Build the message
		from,
		&to,
		0,
		new(big.Int).SetUint64(0), // Amount to transfer
		gaslimit,
		parentApiRouter.Message().GasPrice, // gas price
		funCallData,
		nil,
		false, // Don't checking nonce
	)
	return NewJobFromNative(jobID, batchID, &msg, parentApiRouter)
}

func NewJobFromNative(jobID, batchID uint64, nativeMsg *evmcore.Message, parentApiRouter eucommon.EthApiRouter) *Job {
	job := &Job{
		ID:        jobID,
		EvmMsg:    nativeMsg,
		ApiRouter: parentApiRouter,
	}
	job.TxHash = job.Hash(batchID)
	return job
}

func (this *Job) Hash(batchID uint64) [32]byte {
	return sha256.Sum256(common.Flatten([][]byte{
		codec.Bytes32(this.ApiRouter.TxHash()).Encode(),
		codec.Uint32((batchID)).Encode(),
		this.EvmMsg.Data[:4],
		codec.Uint32(this.ID).Encode(),
	}))
}

func (this *Job) CaptureStates(snapshotUrl ccurlinterfaces.Datastore) eucommon.EthApiRouter {
	ccurl := (&concurrenturl.ConcurrentUrl{}).New(
		indexer.NewWriteCache(snapshotUrl, this.ApiRouter.Ccurl().Platform),
		this.ApiRouter.Ccurl().Platform) // Init a write cache only since it doesn't need the importers

	return this.ApiRouter.New(this.TxHash, uint32(this.ID), this.ApiRouter.Depth(), ccurl)
}

func (this *Job) Run(config *cceu.Config, snapshotUrl ccurlinterfaces.Datastore) *Result { //
	this.ApiRouter = this.CaptureStates(snapshotUrl)
	statedb := eth.NewImplStateDB(this.ApiRouter)                // Eth state DB
	statedb.PrepareFormer(this.TxHash, [32]byte{}, int(this.ID)) // tx hash , block hash and tx index

	eu := cceu.NewEU(
		config.ChainConfig,
		vm.Config{},
		statedb,
		this.ApiRouter, // Tx hash, tx id and url
	)

	var prechkErr error
	this.Receipt, this.EvmResult, prechkErr =
		eu.Run(
			this.TxHash,
			int(this.ID),
			this.EvmMsg,
			cceu.NewEVMBlockContext(config),
			cceu.NewEVMTxContext(*this.EvmMsg),
		)

	// Do gas transfer
	if prechkErr == nil && this.EvmResult != nil && this.EvmResult.Err == nil && this.ApiRouter.GetReserved() != nil {
		deferred := this.ApiRouter.GetReserved().(*StandardMessage)
		if this.EvmMsg.GasLimit-this.EvmResult.UsedGas >= deferred.Native.GasLimit {
			eu.VM().Context.Transfer(
				eu.VM().StateDB,
				this.EvmMsg.From,
				eucommon.ATOMIC_HANDLER,
				big.NewInt(int64(deferred.Native.GasLimit)),
			)
		}
	}
	transitions := this.ApiRouter.Ccurl().Export()
	indexer.Univalues(transitions).Print()

	return &Result{
		TxIndex: uint32(this.ID),
		TxHash:  common.IfThenDo1st(this.Receipt != nil, func() evmcommon.Hash { return this.Receipt.TxHash }, evmcommon.Hash{}),
		Spawned: common.IfThenDo1st(this.ApiRouter.GetReserved() != nil,
			func() *StandardMessage {
				return this.ApiRouter.GetReserved().(*StandardMessage)
			},
			nil),
		Transitions: transitions, // Transitions + Accesses
		Err:         common.IfThenDo1st(prechkErr == nil, func() error { return this.EvmResult.Err }, prechkErr),
		GasUsed:     common.IfThenDo1st(this.Receipt != nil, func() uint64 { return this.Receipt.GasUsed }, 0),
	}
}

func (this *Job) CalcualteRefund() uint64 {
	amount := uint64(0)
	for _, v := range *this.ApiRouter.Ccurl().WriteCache().Cache() {
		typed := v.Value().(ccurlinterfaces.Type)
		amount += common.IfThen(
			!v.Preexist(),
			(uint64(typed.Size())/32)*uint64(v.Writes())*evmparams.SstoreSetGas,
			(uint64(typed.Size())/32)*uint64(v.Writes()),
		)
	}
	return amount
}

func (this *Job) RefundTo(payer, recipent ccurlinterfaces.Univalue, amount uint64) (uint64, error) {
	// amount := uint64(this.receipt.GasUsed)
	credit := commutative.NewU256Delta(uint256.NewInt(amount), true).(*commutative.U256)
	if _, _, _, _, err := recipent.Value().(ccurlinterfaces.Type).Set(credit, nil); err != nil {
		return 0, err
	}
	recipent.IncrementDeltaWrites(1)

	debit := commutative.NewU256Delta(uint256.NewInt(amount), false).(*commutative.U256)
	if _, _, _, _, err := payer.Value().(ccurlinterfaces.Type).Set(debit, nil); err != nil {
		return 0, err
	}
	payer.IncrementDeltaWrites(1)

	return amount, nil
}
