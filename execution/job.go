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
	evmcoretypes "github.com/arcology-network/evm/core/types"
	"github.com/arcology-network/evm/core/vm"
	evmparams "github.com/arcology-network/evm/params"
	eucommon "github.com/arcology-network/vm-adaptor/common"
	"github.com/arcology-network/vm-adaptor/eth"
	"github.com/holiman/uint256"
)

type Job struct {
	ID           uint32
	Prefix       []byte
	Predecessors [][32]byte
	Message      *evmcoretypes.Message
	ApiRouter    eucommon.EthApiRouter
	Receipt      *evmcoretypes.Receipt
	EvmResult    *evmcore.ExecutionResult
}

func NewJob(ID uint32, from, to evmcommon.Address, funCallData []byte, gaslimit uint64, parentApiRouter eucommon.EthApiRouter) *Job {
	msg := evmcoretypes.NewMessage( // Build the message
		from,
		&to,
		0,
		new(big.Int).SetUint64(0), // Amount to transfer
		gaslimit,
		parentApiRouter.Message().GasPrice(),
		funCallData,
		nil,
		false, // Don't checking nonce
	)

	return &Job{
		ID:        ID,
		Message:   &msg,
		ApiRouter: parentApiRouter,
	}
}

func (this *Job) Run(coinbase [20]byte, snapshotUrl ccurlinterfaces.Datastore) *Result { //
	this.ApiRouter = this.CaptureStates(this.Prefix, snapshotUrl)
	statedb := eth.NewImplStateDB(this.ApiRouter)                      // Eth state DB
	statedb.Prepare(this.ApiRouter.TxHash(), [32]byte{}, int(this.ID)) // tx hash , block hash and tx index

	config := cceu.NewConfig().SetCoinbase(coinbase) // Share the same coinbase as the main thread
	eu := cceu.NewEU(
		config.ChainConfig,
		vm.Config{},
		statedb,
		this.ApiRouter, // Tx hash, tx id and url
	)

	var prechkErr error
	this.Receipt, this.EvmResult, prechkErr =
		eu.Run(
			eu.Api().TxHash(),
			int(eu.Api().TxIndex()),
			this.Message,
			cceu.NewEVMBlockContext(config),
			cceu.NewEVMTxContext(*this.Message),
		)

	transitions := this.ApiRouter.Ccurl().Export()
	transitions = common.RemoveIf(&transitions, func(v ccurlinterfaces.Univalue) bool {
		return v.GetTx() == ccurlinterfaces.SYSTEM // remove temp
	})

	return &Result{
		TxHash:      common.IfThenDo1st(this.Receipt != nil, func() evmcommon.Hash { return this.Receipt.TxHash }, evmcommon.Hash{}),
		Deferred:    common.IfThenDo1st(this.ApiRouter.GetReserved() != nil, func() *DeferredCall { return this.ApiRouter.GetReserved().(*DeferredCall) }, nil),
		Transitions: transitions, // Transitions + Accesses
		Err:         common.IfThenDo1st(prechkErr == nil, func() error { return this.EvmResult.Err }, prechkErr),
		GasUsed:     common.IfThenDo1st(this.Receipt != nil, func() uint64 { return this.Receipt.GasUsed }, 0),
	}
}

func (this *Job) CaptureStates(prefix []byte, snapshotUrl ccurlinterfaces.Datastore) eucommon.EthApiRouter {
	ccurl := (&concurrenturl.ConcurrentUrl{}).New(
		indexer.NewWriteCache(snapshotUrl, this.ApiRouter.Ccurl().Platform),
		this.ApiRouter.Ccurl().Platform) // Init a write cache only since it doesn't need the importers

	// parentTxHash := parentApiRouter.TxHash()
	txHash := sha256.Sum256(common.Flatten([][]byte{
		codec.Bytes32(this.ApiRouter.TxHash()).Encode(),
		prefix,
		codec.Uint32(this.ID).Encode(),
	}))

	return this.ApiRouter.New(txHash, this.ID, this.ApiRouter.Depth(), ccurl)
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
