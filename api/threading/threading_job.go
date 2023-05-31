package threading

import (
	"strings"

	common "github.com/arcology-network/common-lib/common"
	"github.com/arcology-network/concurrenturl/commutative"
	"github.com/arcology-network/concurrenturl/indexer"
	evmcommon "github.com/arcology-network/evm/common"
	cceu "github.com/arcology-network/vm-adaptor"

	ccinterfaces "github.com/arcology-network/concurrenturl/interfaces"
	"github.com/arcology-network/evm/core"
	ethtypes "github.com/arcology-network/evm/core/types"
	"github.com/arcology-network/evm/core/vm"
	ethparams "github.com/arcology-network/evm/params"
	"github.com/arcology-network/vm-adaptor/eth"
	interfaces "github.com/arcology-network/vm-adaptor/interfaces"
	"github.com/holiman/uint256"
)

type Job struct {
	sender      evmcommon.Address
	caller      evmcommon.Address
	callee      evmcommon.Address
	message     ethtypes.Message
	receipt     *ethtypes.Receipt
	result      *core.ExecutionResult
	prechkErr   error
	apiRounter  interfaces.ApiRouter
	hasConflict bool // Arcology detected errors
}

func (this *Job) Run(config *cceu.Config, statedb *eth.ImplStateDB) { //
	eu := cceu.NewEU(
		config.ChainConfig,
		vm.Config{},
		statedb,
		this.apiRounter, // Tx hash, tx id and url
	)

	this.receipt, this.result, this.prechkErr =
		eu.Run(eu.Api().TxHash(), int(eu.Api().TxIndex()), &this.message, cceu.NewEVMBlockContext(config), cceu.NewEVMTxContext(this.message))
}

func (this *Job) GetAccessInfo() []ccinterfaces.Univalue {
	if this.isSuccessful() || this.hasConflict {
		return []ccinterfaces.Univalue{}
	}
	all := this.apiRounter.Ccurl().Export()
	return indexer.Univalues(all).To(indexer.ITCAccess{})
}

func (this *Job) GetTransitions() []ccinterfaces.Univalue {
	all := []ccinterfaces.Univalue(indexer.Univalues(this.apiRounter.Ccurl().Export()).To(indexer.IPCTransition{}))
	return common.RemoveIf(
		&all,
		func(v ccinterfaces.Univalue) bool { // Nonce or committed path
			return strings.HasSuffix(*v.GetPath(), "/nonce") || (common.IsPath(*v.GetPath()) && v.Preexist())
		},
	)
}

func (this *Job) isSuccessful() bool {
	return this.prechkErr == nil && this.receipt.Status == 1
}

func (this *Job) CalcualteRefund() uint64 {
	amount := uint64(0)
	for _, v := range *this.apiRounter.Ccurl().WriteCache().Cache() {
		typed := v.Value().(ccinterfaces.Type)
		amount += common.IfThen(
			!v.Preexist(),
			(uint64(typed.Size())/32)*uint64(v.Writes())*ethparams.SstoreSetGas,
			(uint64(typed.Size())/32)*uint64(v.Writes()),
		)
	}
	return amount
}

func (this *Job) RefundTo(payer, recipent ccinterfaces.Univalue, amount uint64) (uint64, error) {
	// amount := uint64(this.receipt.GasUsed)
	credit := commutative.NewU256Delta(uint256.NewInt(amount), true).(*commutative.U256)
	if _, _, _, _, err := recipent.Value().(ccinterfaces.Type).Set(credit, nil); err != nil {
		return 0, err
	}
	recipent.IncrementDeltaWrites(1)

	debit := commutative.NewU256Delta(uint256.NewInt(amount), false).(*commutative.U256)
	if _, _, _, _, err := payer.Value().(ccinterfaces.Type).Set(debit, nil); err != nil {
		return 0, err
	}
	payer.IncrementDeltaWrites(1)

	return amount, nil
}
