package multiprocess

import (
	"encoding/hex"
	"strings"

	common "github.com/arcology-network/common-lib/common"
	ccurlcommon "github.com/arcology-network/concurrenturl/common"
	"github.com/arcology-network/concurrenturl/commutative"
	"github.com/arcology-network/concurrenturl/univalue"
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
	for _, v := range *this.apiRounter.Ccurl().WriteCache().Cache() {
		typed := v.Value().(ccurlcommon.TypeInterface)
		return (uint64(typed.Size()) / 32) * uint64(v.Reads())
		common.IfThen(
			!v.Preexist(),
			(uint64(typed.Size())/32)*uint64(v.Writes())*ethparams.SstoreSetGas,
			(uint64(typed.Size())/32)*uint64(v.Writes()),
		)
	}
}

func (this *Job) RefundTo(payer, recipent ccurlcommon.UnivalueInterface, amount uint64) uint64 {
	// amount := uint64(this.receipt.GasUsed)
	credit := commutative.NewU256Delta(uint256.NewInt(amount), true).(*commutative.U256)
	recipent.Value().(ccurlcommon.TypeInterface).Set(credit, nil)
	recipent.IncrementDelta(1)

	debit := commutative.NewU256Delta(uint256.NewInt(amount), false).(*commutative.U256)
	payer.Value().(ccurlcommon.TypeInterface).Set(debit, nil)
	payer.IncrementDelta(1)
	return amount
}

func (this *Job) FilteredAccesses() []ccurlcommon.UnivalueInterface {
	if this.prechkErr != nil && this.receipt.Status != 1 {
		return []ccurlcommon.UnivalueInterface{}
	}
	return univalue.Univalues((this.apiRounter.Ccurl().Export())).To(univalue.AccessCodecFilterSet()...)
}

func (this *Job) FilteredTransitions() []ccurlcommon.UnivalueInterface {
	if this.prechkErr != nil && this.receipt.Status != 1 { // Failed transaction
		return common.CopyIf(this.apiRounter.Ccurl().Export(), func(v ccurlcommon.UnivalueInterface) bool {
			return v != nil &&
				(strings.HasSuffix(*v.GetPath(), "/nonce") ||
					strings.HasSuffix(*v.GetPath(), "/balance"))
		})
	}

	if this.hasConflict { // Transaction has conflicts, refund some gas to the sender
		_, recipent := common.FindFirstIf(this.apiRounter.Ccurl().Export(), func(v ccurlcommon.UnivalueInterface) bool {
			return strings.HasSuffix(*v.GetPath(), "/balance") && strings.Index(*v.GetPath(), hex.EncodeToString(this.sender[:])) > 0
		})

		_, payer := common.FindFirstIf(this.apiRounter.Ccurl().Export(), func(v ccurlcommon.UnivalueInterface) bool {
			coinbase := this.apiRounter.Coinbase()
			return strings.HasSuffix(*v.GetPath(), "/balance") && strings.Index(*v.GetPath(), hex.EncodeToString(coinbase[:])) > 0
		})

		this.RefundTo(*payer, *recipent, 0)
	}

	return univalue.Univalues(this.apiRounter.Ccurl().Export()).To(
		univalue.RemoveReadOnly,
		univalue.DelNonExist,
	)
}
