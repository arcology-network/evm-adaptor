package execution

import (
	// "github.com/arcology-network/common-lib/codec"

	"encoding/hex"
	"fmt"
	"strings"

	"github.com/arcology-network/common-lib/codec"
	common "github.com/arcology-network/common-lib/common"
	indexer "github.com/arcology-network/concurrenturl/indexer"
	ccurlinterfaces "github.com/arcology-network/concurrenturl/interfaces"
	"github.com/arcology-network/concurrenturl/univalue"
	evmcore "github.com/arcology-network/evm/core"
	evmTypes "github.com/arcology-network/evm/core/types"
	"github.com/holiman/uint256"
)

type Result struct {
	GroupID            uint32 // == Group ID
	TxIndex            uint32
	TxHash             [32]byte
	From               [20]byte
	Coinbase           [20]byte
	rawStateAccesses   []ccurlinterfaces.Univalue
	immunedTransitions []ccurlinterfaces.Univalue
	transitions        []ccurlinterfaces.Univalue
	Receipt            *evmTypes.Receipt
	EvmResult          *evmcore.ExecutionResult
	Err                error
}

func (this *Result) BreakdownBalanceTransition(balanceTransition ccurlinterfaces.Univalue, gasDelta *uint256.Int, isCredit bool) ccurlinterfaces.Univalue {
	if delta := (*uint256.Int)(balanceTransition.Value().(ccurlinterfaces.Type).Delta().(*codec.Uint256)); delta.Cmp(gasDelta) >= 0 {
		transfer := delta.Sub(delta, (*uint256.Int)(gasDelta))                                  // balance - gas
		(balanceTransition).Value().(ccurlinterfaces.Type).SetDelta((*codec.Uint256)(transfer)) // Set the transfer, Won't change the initial value.
		(balanceTransition).Value().(ccurlinterfaces.Type).SetDeltaSign(false)
		//
		newGasTransition := balanceTransition.Clone().(ccurlinterfaces.Univalue)
		newGasTransition.Value().(ccurlinterfaces.Type).SetDelta((*codec.Uint256)(gasDelta))
		newGasTransition.Value().(ccurlinterfaces.Type).SetDeltaSign(isCredit)
		newGasTransition.GetUnimeta().(*univalue.Unimeta).SetPersistent(true)
		return newGasTransition
	}
	return nil
}

func (this *Result) Postprocess() *Result {
	_, senderBalance := common.FindFirstIf(this.rawStateAccesses, func(v ccurlinterfaces.Univalue) bool {
		return v != nil && strings.HasSuffix(*v.GetPath(), "/balance") && strings.Contains(*v.GetPath(), hex.EncodeToString(this.From[:]))
	})

	if senderGasTransition := this.BreakdownBalanceTransition(*senderBalance, uint256.NewInt(this.Receipt.GasUsed), false); senderGasTransition != nil {
		this.immunedTransitions = append(this.immunedTransitions, senderGasTransition)
	}

	_, coinbaseBalance := common.FindFirstIf(this.rawStateAccesses, func(v ccurlinterfaces.Univalue) bool {
		return v != nil && strings.HasSuffix(*v.GetPath(), "/balance") || strings.Contains(*v.GetPath(), hex.EncodeToString(this.Coinbase[:]))
	})

	if coinbaseGasTransition := this.BreakdownBalanceTransition(*coinbaseBalance, uint256.NewInt(this.Receipt.GasUsed), true); coinbaseGasTransition != nil {
		this.immunedTransitions = append(this.immunedTransitions, coinbaseGasTransition)
	}

	common.Foreach(this.rawStateAccesses, func(v *ccurlinterfaces.Univalue) {
		if v != nil {
			return
		}

		path := (*v).GetPath()
		if strings.HasSuffix(*path, "/nonce") && strings.Contains(*path, hex.EncodeToString(this.From[:])) {
			(*v).GetUnimeta().(*univalue.Unimeta).SetPersistent(true)
		}
	})

	if this.Err == nil {
		this.transitions = indexer.Univalues(this.rawStateAccesses).To(indexer.ITCTransition{Err: this.Err})
		this.transitions = common.MoveIf(&this.transitions, func(v ccurlinterfaces.Univalue) bool { return !v.Persistent() })
	}
	return this
}

func (this *Result) Print() {
	// fmt.Println("GroupID: ", this.GroupID)
	fmt.Println("TxIndex: ", this.TxIndex)
	fmt.Println("TxHash: ", this.TxHash)
	fmt.Println()
	indexer.Univalues(this.rawStateAccesses).Print()
	fmt.Println("Error: ", this.Err)
}

type Results []*Result

func (this Results) Print() {
	fmt.Println("Execution Results: ")
	for _, v := range this {
		v.Print()
		fmt.Println()
	}
}
