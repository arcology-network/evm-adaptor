package execution

import (
	// "github.com/arcology-network/common-lib/codec"

	"fmt"

	common "github.com/arcology-network/common-lib/common"
	arbitrator "github.com/arcology-network/concurrenturl/arbitrator"
	indexer "github.com/arcology-network/concurrenturl/indexer"
	ccurlinterfaces "github.com/arcology-network/concurrenturl/interfaces"
	evmcore "github.com/arcology-network/evm/core"
	evmTypes "github.com/arcology-network/evm/core/types"
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

func (this *Result) WriteTo(newTxIdx uint32, targetCache *indexer.WriteCache) {
	transitions := []ccurlinterfaces.Univalue(indexer.Univalues(common.Clone(this.Transitions)).To(
		TransitionFilter{
			Err:      this.Err,
			Sender:   this.From,
			Coinbase: *this.Config.Coinbase,
		},
	))

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
