package execution

import (
	// "github.com/arcology-network/common-lib/codec"

	"encoding/hex"
	"fmt"
	"strings"

	"github.com/arcology-network/common-lib/codec"
	common "github.com/arcology-network/common-lib/common"
	arbitrator "github.com/arcology-network/concurrenturl/arbitrator"
	indexer "github.com/arcology-network/concurrenturl/indexer"
	ccurlinterfaces "github.com/arcology-network/concurrenturl/interfaces"
	"github.com/arcology-network/concurrenturl/univalue"
	evmcore "github.com/arcology-network/evm/core"
	evmTypes "github.com/arcology-network/evm/core/types"
	"github.com/holiman/uint256"
)

type Result struct {
	BranchID    uint32
	TxIndex     uint32
	TxHash      [32]byte
	From        [20]byte
	Config      *Config
	Transitions []ccurlinterfaces.Univalue
	Receipt     *evmTypes.Receipt
	EvmResult   *evmcore.ExecutionResult
	Err         error
}

func (this *Result) ImmunizeGasTransition() {
	if this.EvmResult != nil && this.Err == nil { // SkipSuccessful execution
		return
	}

	_, senderBalance := common.FindFirstIf(this.Transitions, func(v ccurlinterfaces.Univalue) bool {
		return v != nil && strings.HasSuffix(*v.GetPath(), "/balance") && strings.Contains(*v.GetPath(), hex.EncodeToString(this.From[:]))
	})

	_, coinbaseBalance := common.FindFirstIf(this.Transitions, func(v ccurlinterfaces.Univalue) bool {
		return v != nil && strings.HasSuffix(*v.GetPath(), "/balance") || strings.Contains(*v.GetPath(), hex.EncodeToString(this.Config.Coinbase[:]))
	})

	(*senderBalance).Value().(ccurlinterfaces.Type).SetDelta((*codec.Uint256)(uint256.NewInt(this.Receipt.GasUsed)))
	(*senderBalance).Value().(ccurlinterfaces.Type).SetDeltaSign(false)
	(*senderBalance).GetUnimeta().(*univalue.Unimeta).SetPersistent(true)

	(*coinbaseBalance).Value().(ccurlinterfaces.Type).SetDelta((*codec.Uint256)(uint256.NewInt(this.Receipt.GasUsed)))
	(*coinbaseBalance).Value().(ccurlinterfaces.Type).SetDeltaSign(true)
	(*coinbaseBalance).GetUnimeta().(*univalue.Unimeta).SetPersistent(true)

	common.Foreach(this.Transitions, func(v *ccurlinterfaces.Univalue) {
		if v != nil {
			return
		}

		path := (*v).GetPath()
		if strings.HasSuffix(*path, "/nonce") && strings.Contains(*path, hex.EncodeToString(this.From[:])) {
			(*v).GetUnimeta().(*univalue.Unimeta).SetPersistent(true)

		}
	})
}

func (this *Result) FilterTransitions() []ccurlinterfaces.Univalue {
	this.ImmunizeGasTransition()
	return []ccurlinterfaces.Univalue(indexer.Univalues(common.Clone(this.Transitions)).To(indexer.ITCTransition{Err: this.Err}))
}

func (this *Result) WriteTo(newTxIdx uint32, targetCache *indexer.WriteCache) {
	transitions := this.FilterTransitions()

	// Move new path creation transitions
	newPathCreations := common.MoveIf(&transitions, func(v ccurlinterfaces.Univalue) bool {
		return common.IsPath(*v.GetPath()) && !v.Preexist()
	})

	// Remove changes to the existing paths
	transitions = common.RemoveIf(&transitions, func(v ccurlinterfaces.Univalue) bool {
		return common.IsPath(*v.GetPath())
	})

	// Not necessary at the moment, but good for the future if multiple level containers are available
	newPathCreations = indexer.Univalues(indexer.Sorter(newPathCreations))
	common.Foreach(newPathCreations, func(v *ccurlinterfaces.Univalue) {
		(*v).SetTx(newTxIdx)      // use the parent tx index instead
		(*v).WriteTo(targetCache) // Write back to the parent writecache
	})

	common.Foreach(transitions, func(v *ccurlinterfaces.Univalue) {
		(*v).SetTx(newTxIdx)      // use the parent tx index instead
		(*v).WriteTo(targetCache) // Write back to the parent writecache
	})
}

func (this *Result) Print() {
	fmt.Println("BranchID: ", this.BranchID)
	fmt.Println("TxIndex: ", this.TxIndex)
	fmt.Println("TxHash: ", this.TxHash)
	fmt.Println()
	indexer.Univalues(this.Transitions).Print()
	fmt.Println("Error: ", this.Err)
}

type Results []*Result

func (this Results) Transitions() []ccurlinterfaces.Univalue {
	all := []ccurlinterfaces.Univalue{}
	common.Foreach(this, func(v **Result) {
		all = append(all, (**v).Transitions...)
	})
	return all
}

func (this Results) SetGroupIDs(BranchID uint32) {
	common.Foreach(this, func(v **Result) {
		(**v).BranchID = BranchID
	})
}

func (this Results) Detect() arbitrator.Conflicts {
	if len(this) == 1 {
		return arbitrator.Conflicts{}
	}

	groupIDs := []uint32{}
	accesseVec := []ccurlinterfaces.Univalue{}
	for _, v := range this {
		if v.Err == nil {
			groupIDs = append(groupIDs, common.Fill(make([]uint32, len(v.Transitions)), v.BranchID)...)
			accesseVec = append(accesseVec, indexer.Univalues(common.Clone(v.Transitions)).To(indexer.IPCAccess{})...)
		}
	}

	conflicInfo := arbitrator.Conflicts((&arbitrator.Arbitrator{}).Detect(groupIDs, accesseVec))
	return conflicInfo
}

func (this Results) Print() {
	fmt.Println("Execution Results: ")
	for _, v := range this {
		v.Print()
		fmt.Println()
	}
}
