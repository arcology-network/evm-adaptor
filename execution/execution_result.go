package execution

import (
	// "github.com/arcology-network/common-lib/codec"

	"errors"

	common "github.com/arcology-network/common-lib/common"
	arbitrator "github.com/arcology-network/concurrenturl/arbitrator"
	indexer "github.com/arcology-network/concurrenturl/indexer"
	ccurlinterfaces "github.com/arcology-network/concurrenturl/interfaces"
)

type Result struct {
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

func (this Results) DetectConflict() []*Result {
	accesseVec := common.Concate(this, func(v *Result) []ccurlinterfaces.Univalue {
		return common.IfThen(
			v.Err == nil,
			indexer.Univalues(common.Clone(v.Transitions)).To(indexer.IPCAccess{}),
			[]ccurlinterfaces.Univalue{})
	})
	dict := arbitrator.Conflicts((&arbitrator.Arbitrator{}).Detect(accesseVec)).ToDict()

	common.Foreach(this, func(result **Result) {
		if _, conflict := (dict)[(**result).TxIndex]; conflict { // Label conflicts
			(**result).Err = errors.New("Error: Conflicts detected in state accesses")
		}
	})
	return this
}

func (this Results) ToSequence() *Sequence {
	if this[0].Deferred == nil {
		return nil
	}

	// to := evmcommon.Address(this[0].Deferred.Addr)
	// evmMsg := evmcoretypes.NewMessage(
	// 	this[0].Deferred.From, // From the system account
	// 	&to,
	// 	0,
	// 	big.NewInt(0),
	// 	1e15,
	// 	big.NewInt(1),
	// 	this[0].Deferred.FuncCallData,
	// 	nil,
	// 	false,
	// )

	predecessors := make([][32]byte, 0, len(this))
	common.Foreach(this, func(v **Result) { predecessors = append(predecessors, (**v).Deferred.TxHash) })

	// msg := &StandardMessage{
	// 	TxHash: sha256.Sum256(this[0].Deferred.FuncCallData),
	// 	Native: &evmMsg,
	// 	Source: commontypes.TX_SOURCE_DEFERRED,
	// }
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
