package execution

import (
	"crypto/sha256"

	"github.com/arcology-network/common-lib/codec"
	"github.com/arcology-network/common-lib/common"
	"github.com/arcology-network/concurrenturl"
	"github.com/arcology-network/concurrenturl/commutative"
	indexer "github.com/arcology-network/concurrenturl/indexer"
	"github.com/arcology-network/concurrenturl/interfaces"

	ccurlinterfaces "github.com/arcology-network/concurrenturl/interfaces"
	evmcommon "github.com/arcology-network/evm/common"
	"github.com/arcology-network/evm/core/vm"
	evmparams "github.com/arcology-network/evm/params"
	eucommon "github.com/arcology-network/vm-adaptor/common"
	"github.com/arcology-network/vm-adaptor/eth"
	"github.com/holiman/uint256"
)

type JobSequence struct {
	ID               uint32 // group id
	PreTxs           []uint32
	StdMsgs          []*StandardMessage
	Results          []*Result
	ApiRouter        eucommon.EthApiRouter
	RecordBuffer     []ccurlinterfaces.Univalue
	TransitionBuffer []ccurlinterfaces.Univalue
	immunedBuffer    []ccurlinterfaces.Univalue
}

func (this *JobSequence) DeriveNewHash(seed [32]byte) [32]byte {
	return sha256.Sum256(common.Flatten([][]byte{
		codec.Bytes32(seed).Encode(),
		codec.Uint32(this.ID).Encode(),
	}))
}

func (this *JobSequence) Length() int { return len(this.StdMsgs) }

func (this *JobSequence) Run(config *Config, mainApi eucommon.EthApiRouter) []*Result { //
	results := make([]*Result, len(this.StdMsgs))
	this.ApiRouter = mainApi.New((&concurrenturl.ConcurrentUrl{}).New(indexer.NewWriteCache(mainApi.Ccurl().WriteCache())), this.ApiRouter.Schedule())

	for i, msg := range this.StdMsgs {
		pendingApi := this.ApiRouter.New((&concurrenturl.ConcurrentUrl{}).New(indexer.NewWriteCache(this.ApiRouter.Ccurl().WriteCache())), this.ApiRouter.Schedule())
		pendingApi.DecrementDepth()

		results[i] = this.execute(msg, config, pendingApi)                         // What happens if it fails
		this.ApiRouter.Ccurl().WriteCache().AddTransitions(results[i].transitions) // merge transitions to the main cache here !!!
	}

	this.RecordBuffer = indexer.Univalues(this.ApiRouter.Ccurl().Export()).To(indexer.IPCAccess{})

	this.immunedBuffer = common.Concate(this.Results,
		func(v *Result) []interfaces.Univalue {
			return (*v).immunedTransitions
		},
	)

	this.TransitionBuffer = append(this.TransitionBuffer, indexer.Univalues(this.ApiRouter.Ccurl().Export()).To(indexer.ITCTransition{})...)
	this.TransitionBuffer = append(this.TransitionBuffer, this.immunedBuffer...)
	return results
}

func (this *JobSequence) FlagConflict(dict *map[uint32]uint64, err error) {
	for i := 0; i < len(this.Results); i++ {
		this.Results[i].Err = err // Flag the transitions for the WriteTo().
	}

	this.RecordBuffer = this.RecordBuffer[:0]
	this.TransitionBuffer = this.immunedBuffer
}

func (this *JobSequence) execute(stdMsg *StandardMessage, config *Config, api eucommon.EthApiRouter) *Result { //
	statedb := eth.NewImplStateDB(api)                                  // Eth state DB
	statedb.PrepareFormer(stdMsg.TxHash, [32]byte{}, uint32(stdMsg.ID)) // tx hash , block hash and tx index

	eu := NewEU(
		config.ChainConfig,
		vm.Config{},
		statedb,
		api, // Tx hash, tx id and url
	)

	// var prechkErr error
	receipt, evmResult, prechkErr :=
		eu.Run(
			stdMsg,
			NewEVMBlockContext(config),
			NewEVMTxContext(*stdMsg.Native),
		)

	return (&Result{
		TxIndex:          uint32(stdMsg.ID),
		TxHash:           common.IfThenDo1st(receipt != nil, func() evmcommon.Hash { return receipt.TxHash }, evmcommon.Hash{}),
		rawStateAccesses: api.StateFilter().Raw(), // Transitions + Accesses
		Err:              common.IfThenDo1st(prechkErr == nil, func() error { return evmResult.Err }, prechkErr),
		From:             stdMsg.Native.From,
		Coinbase:         *config.Coinbase,
		Receipt:          receipt,
		EvmResult:        evmResult,
	}).Postprocess()
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
