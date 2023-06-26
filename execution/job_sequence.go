package execution

import (
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

type JobSequence struct {
	ID        uint64
	StdMsgs   []*StandardMessage
	Results   []*Result
	ApiRouter eucommon.EthApiRouter
}

func (this *JobSequence) Run(config *cceu.Config, snapshotUrl ccurlinterfaces.Datastore) []*Result { //
	results := make([]*Result, len(this.StdMsgs))

	for i, msg := range this.StdMsgs {
		results[i] = this.execute(msg, config, snapshotUrl) // What happens if it fails
		transitions := indexer.Univalues(common.Clone(results[i].Transitions)).To(indexer.ITCTransition{})
		if i < len(this.StdMsgs)-1 {
			snapshotUrl = this.ApiRouter.Ccurl().Snapshot(transitions)
		}
	}
	return results
}

func (this *JobSequence) execute(stdMsg *StandardMessage, config *cceu.Config, snapshotUrl ccurlinterfaces.Datastore) *Result { //
	ccurl := (&concurrenturl.ConcurrentUrl{}).New(
		indexer.NewWriteCache(snapshotUrl, this.ApiRouter.Ccurl().Platform),
		this.ApiRouter.Ccurl().Platform) // Init a write cache only since it doesn't need the importers

	this.ApiRouter = this.ApiRouter.New(stdMsg.TxHash, uint32(stdMsg.ID), this.ApiRouter.Depth(), ccurl, this.ApiRouter.Schedule())
	statedb := eth.NewImplStateDB(this.ApiRouter)                       // Eth state DB
	statedb.PrepareFormer(stdMsg.TxHash, [32]byte{}, uint32(stdMsg.ID)) // tx hash , block hash and tx index

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
			uint32(stdMsg.ID),
			stdMsg.Native,
			cceu.NewEVMBlockContext(config),
			cceu.NewEVMTxContext(*stdMsg.Native),
		)

	// Do gas transfer
	// if prechkErr == nil && evmResult != nil && evmResult.Err == nil && this.ApiRouter.GetReserved() != nil {
	// 	deferred := this.ApiRouter.GetReserved().(*StandardMessage)
	// 	if stdMsg.Native.GasLimit-evmResult.UsedGas >= deferred.Native.GasLimit {
	// 		eu.VM().Context.Transfer(
	// 			eu.VM().StateDB,
	// 			stdMsg.Native.From,
	// 			eucommon.ATOMIC_HANDLER,
	// 			big.NewInt(int64(deferred.Native.GasLimit)),
	// 		)
	// 	}
	// }
	indexer.Univalues(this.ApiRouter.Ccurl().Export()).Print()

	return &Result{
		TxIndex:     uint32(stdMsg.ID),
		TxHash:      common.IfThenDo1st(receipt != nil, func() evmcommon.Hash { return receipt.TxHash }, evmcommon.Hash{}),
		Transitions: this.ApiRouter.Ccurl().Export(), // Transitions + Accesses
		Err:         common.IfThenDo1st(prechkErr == nil, func() error { return evmResult.Err }, prechkErr),

		Receipt:   receipt,
		EvmResult: evmResult,
	}
}

func (this *JobSequence) CalcualteRefund() uint64 {
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

func (this *JobSequence) RefundTo(payer, recipent ccurlinterfaces.Univalue, amount uint64) (uint64, error) {
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
