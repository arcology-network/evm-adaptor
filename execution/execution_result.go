package execution

import (
	// "github.com/arcology-network/common-lib/codec"

	"errors"
	"fmt"

	common "github.com/arcology-network/common-lib/common"
	arbitrator "github.com/arcology-network/concurrenturl/arbitrator"
	indexer "github.com/arcology-network/concurrenturl/indexer"
	ccurlinterfaces "github.com/arcology-network/concurrenturl/interfaces"
)

type Result struct {
	GroupID     uint32
	TxIndex     uint32
	TxHash      [32]byte
	Deferred    *StandardMessage
	Transitions []ccurlinterfaces.Univalue
	Err         error
	GasUsed     uint64
}

func (this *Result) WriteTo(newTxIdx uint32, targetCache *indexer.WriteCache) {
	transitions := []ccurlinterfaces.Univalue(indexer.Univalues(common.Clone(this.Transitions)).To(TransitionFilter{status: 0}))

	newPathTrans := common.MoveIf(&transitions, func(v ccurlinterfaces.Univalue) bool {
		return common.IsPath(*v.GetPath()) && !v.Preexist() // Move new path creation transitions
	})

	common.Foreach(newPathTrans, func(v *ccurlinterfaces.Univalue) {
		(*v).SetTx(newTxIdx)      // use the parent tx index instead
		(*v).WriteTo(targetCache) // Write back to the parent writecache
	})

	common.Foreach(transitions, func(v *ccurlinterfaces.Univalue) {
		(*v).SetTx(newTxIdx)      // use the parent tx index instead
		(*v).WriteTo(targetCache) // Write back to the parent writecache
	})
}

type Results []*Result

func (this Results) Transitions() []ccurlinterfaces.Univalue {
	all := []ccurlinterfaces.Univalue{}
	common.Foreach(this, func(v **Result) {
		all = append(all, (**v).Transitions...)
	})
	return all
}

func (this Results) SetGroupIDs(GroupID uint32) {
	common.Foreach(this, func(v **Result) {
		(**v).GroupID = GroupID
	})
}

func (this Results) DetectConflict() ([]*Result, int) {
	groupIDs := []uint32{}
	accesseVec := []ccurlinterfaces.Univalue{}
	for _, v := range this {
		if v.Err == nil {
			groupIDs = append(groupIDs, common.Fill(make([]uint32, len(v.Transitions)), v.GroupID)...)
			accesseVec = append(accesseVec, indexer.Univalues(common.Clone(v.Transitions)).To(indexer.IPCAccess{})...)
		}
	}

	conflicInfo := arbitrator.Conflicts((&arbitrator.Arbitrator{}).Detect(groupIDs, accesseVec))
	dict := conflicInfo.ToDict()

	if len(dict) > 0 {
		fmt.Println("Conflict")
	}

	totalConflict := 0
	for i := 0; i < len(this); i++ {
		if _, conflict := (dict)[this[i].TxIndex]; conflict {
			this[i].Err = errors.New("Error: Conflicts detected in state accesses")
			totalConflict++
		}
	}
	return this, totalConflict
}

func (this Results) ToSequence() *Sequence {
	if this[0].Deferred == nil {
		return nil
	}

	predecessors := make([][32]byte, 0, len(this))
	common.Foreach(this, func(v **Result) { predecessors = append(predecessors, (**v).Deferred.TxHash) })
	return NewSequence([32]byte{}, predecessors, []*StandardMessage{this[0].Deferred}, true)
}

// This works with the deferred execution
type ResultDict map[[32]byte][]*Result

func (this *ResultDict) Categorize(results []*Result) [][]*Result {
	if len(results) == 1 {
		return [][]*Result{results}
	}

	for _, v := range results {
		if v.Deferred == nil {
			continue
		}

		vec := (*this)[v.Deferred.GroupBy]
		if vec != nil {
			vec = []*Result{}
		}
		(*this)[v.Deferred.GroupBy] = append(vec, v)
	}
	return common.MapValues((map[[32]byte][]*Result)(*this))
}
