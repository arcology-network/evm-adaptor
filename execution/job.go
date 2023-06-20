package execution

import (
	"math/big"

	common "github.com/arcology-network/common-lib/common"
	"github.com/arcology-network/concurrenturl"
	"github.com/arcology-network/concurrenturl/commutative"
	"github.com/arcology-network/concurrenturl/indexer"
	cceu "github.com/arcology-network/vm-adaptor"

	ccurlinterfaces "github.com/arcology-network/concurrenturl/interfaces"
	evmcommon "github.com/arcology-network/evm/common"
	"github.com/arcology-network/evm/core/vm"
	evmparams "github.com/arcology-network/evm/params"
	eucommon "github.com/arcology-network/vm-adaptor/common"
	"github.com/arcology-network/vm-adaptor/eth"
	"github.com/holiman/uint256"
)

type Job struct {
	JobID     uint64
	BranchID  uint64
	StdMsgs   []*StandardMessage
	Results   []*Result
	ApiRouter eucommon.EthApiRouter
}

func (this *Job) CaptureStates(stdMsg *StandardMessage, snapshotUrl ccurlinterfaces.Datastore) eucommon.EthApiRouter {
	ccurl := (&concurrenturl.ConcurrentUrl{}).New(
		indexer.NewWriteCache(snapshotUrl, this.ApiRouter.Ccurl().Platform),
		this.ApiRouter.Ccurl().Platform) // Init a write cache only since it doesn't need the importers

	return this.ApiRouter.New(stdMsg.TxHash, uint32(stdMsg.ID), this.ApiRouter.Depth(), ccurl)
}

func (this *Job) Run(config *cceu.Config, snapshotUrl ccurlinterfaces.Datastore) []*Result { //
	results := make([]*Result, len(this.StdMsgs))
	for i, msg := range this.StdMsgs {
		results[i] = this.runMsg(msg, config, snapshotUrl)
	}
	return results
}

func (this *Job) runMsg(stdMsg *StandardMessage, config *cceu.Config, snapshotUrl ccurlinterfaces.Datastore) *Result { //
	this.ApiRouter = this.CaptureStates(stdMsg, snapshotUrl)
	statedb := eth.NewImplStateDB(this.ApiRouter)                    // Eth state DB
	statedb.PrepareFormer(stdMsg.TxHash, [32]byte{}, int(stdMsg.ID)) // tx hash , block hash and tx index

	eu := cceu.NewEU(
		config.ChainConfig,
		vm.Config{},
		statedb,
		this.ApiRouter, // Tx hash, tx id and url
	)

	// var prechkErr error
	receipt, evmResult, prechkErr :=
		eu.Run(
			stdMsg.TxHash,
			int(stdMsg.ID),
			stdMsg.Native,
			cceu.NewEVMBlockContext(config),
			cceu.NewEVMTxContext(*stdMsg.Native),
		)

	// Do gas transfer
	if prechkErr == nil && evmResult != nil && evmResult.Err == nil && this.ApiRouter.GetReserved() != nil {
		deferred := this.ApiRouter.GetReserved().(*StandardMessage)
		if stdMsg.Native.GasLimit-evmResult.UsedGas >= deferred.Native.GasLimit {
			eu.VM().Context.Transfer(
				eu.VM().StateDB,
				stdMsg.Native.From,
				eucommon.ATOMIC_HANDLER,
				big.NewInt(int64(deferred.Native.GasLimit)),
			)
		}
	}
	transitions := this.ApiRouter.Ccurl().Export()
	indexer.Univalues(transitions).Print()

	return &Result{
		TxIndex: uint32(stdMsg.ID),
		TxHash:  common.IfThenDo1st(receipt != nil, func() evmcommon.Hash { return receipt.TxHash }, evmcommon.Hash{}),
		Spawned: common.IfThenDo1st(this.ApiRouter.GetReserved() != nil,
			func() *StandardMessage {
				return this.ApiRouter.GetReserved().(*StandardMessage)
			},
			nil),
		Transitions: transitions, // Transitions + Accesses
		Err:         common.IfThenDo1st(prechkErr == nil, func() error { return evmResult.Err }, prechkErr),
		// GasUsed:     common.IfThenDo1st(this.Receipt != nil, func() uint64 { return this.Receipt.GasUsed }, 0),

		Receipt:   receipt,
		EvmResult: evmResult,
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
