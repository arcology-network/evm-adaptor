package types

import (
	// "github.com/arcology-network/common-lib/codec"
	interfaces "github.com/arcology-network/concurrenturl/interfaces"
)

type ExecutionResult struct {
	TxHash      [32]byte
	Deferred    *DeferredCall
	Transitions []interfaces.Univalue
	Err         error
	GasUsed     uint64
}

// type NewExecutionResult struct {
// 	TxHash      [32]byte
// 	TxID        uint32
// 	Deferred    *DeferredCall
// 	Transitions []interfaces.Univalue
// 	Status      uint64
// 	GasUsed     uint64
// }

// func NewExecutionResult() *ExecutionResult {
// 	return uint32((len(this) + 1) * codec.UINT32_LEN) // Header length
// }
