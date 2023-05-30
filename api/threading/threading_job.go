package multiprocess

import (
	common "github.com/arcology-network/common-lib/common"
	"github.com/arcology-network/concurrenturl/commutative"
	interfaces "github.com/arcology-network/concurrenturl/interfaces"
	evmcommon "github.com/arcology-network/evm/common"

	"github.com/arcology-network/evm/core"
	ethtypes "github.com/arcology-network/evm/core/types"
	ethparams "github.com/arcology-network/evm/params"
	eucommon "github.com/arcology-network/vm-adaptor/common"
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
	apiRounter  eucommon.ConcurrentApiRouterInterface
	hasConflict bool // Arcology detected errors
}

func (this *Job) CalcualteRefund() uint64 { //
	amount := uint64(0)
	for _, v := range *this.apiRounter.Ccurl().WriteCache().Cache() {
		typed := v.Value().(interfaces.Type)
		amount += common.IfThen(
			!v.Preexist(),
			(uint64(typed.Size())/32)*uint64(v.Writes())*ethparams.SstoreSetGas,
			(uint64(typed.Size())/32)*uint64(v.Writes()),
		)
	}
	return amount
}

func (this *Job) RefundTo(payer, recipent interfaces.Univalue, amount uint64) (uint64, error) {
	// amount := uint64(this.receipt.GasUsed)
	credit := commutative.NewU256Delta(uint256.NewInt(amount), true).(*commutative.U256)
	if _, _, _, _, err := recipent.Value().(interfaces.Type).Set(credit, nil); err != nil {
		return 0, err
	}
	recipent.IncrementDeltaWrites(1)

	debit := commutative.NewU256Delta(uint256.NewInt(amount), false).(*commutative.U256)
	if _, _, _, _, err := payer.Value().(interfaces.Type).Set(debit, nil); err != nil {
		return 0, err
	}
	payer.IncrementDeltaWrites(1)

	return amount, nil
}
