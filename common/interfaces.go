// KernelAPI provides system level function calls supported by arcology platform.
package common

import (
	"math/big"

	"github.com/arcology-network/concurrenturl"
	"github.com/arcology-network/concurrenturl/interfaces"
	"github.com/arcology-network/evm/common"

	evmcommon "github.com/arcology-network/evm/common"
	"github.com/arcology-network/evm/consensus"
	"github.com/arcology-network/evm/core"
	"github.com/arcology-network/evm/core/types"
	"github.com/arcology-network/evm/core/vm"
)

type EthApiRouter interface {
	Origin() evmcommon.Address
	Ccurl() *concurrenturl.ConcurrentUrl
	New(common.Hash, uint32, *concurrenturl.ConcurrentUrl, interface{}) EthApiRouter
	Coinbase() evmcommon.Address

	SetEU(interface{})
	GetEU() interface{}

	Message() *core.Message
	VM() *vm.EVM

	GetReserved() interface{}
	SetReserved(interface{})

	Depth() uint8
	AddLog(key, value string)
	Call(caller, callee [20]byte, input []byte, origin [20]byte, nonce uint64, blockhash evmcommon.Hash) (bool, []byte, bool, int64)
	SetRuntimeContext([32]byte, uint32, *big.Int)

	TxIndex() uint32
	TxHash() [32]byte
	Schedule() interface{}

	CCUID() uint64
	CCElemID() uint64

	GenCcObjID() []byte
	GenCcElemUID() []byte
	GenUUID() []byte
}

type ILog interface {
	GetKey() string
	GetValue() string
}

type ChainContext interface {
	Engine() consensus.Engine                    // Engine retrieves the chain's consensus engine.
	GetHeader(common.Hash, uint64) *types.Header // GetHeader returns the hash corresponding to their hash.
}

type ApiCallHandler interface {
	Address() [20]byte
	Call([20]byte, [20]byte, []byte, [20]byte, uint64) ([]byte, bool, int64)
}

type StateConflict interface {
	Detect() []uint32
}

type SnapshotMaker interface {
	Import([]interfaces.Univalue)
	Make([]uint32) interface{}
	Clear()
}

// type LocalSnapshotMaker struct {
// 	transitions []interfaces.Univalue
// }

// func (this *LocalSnapshotMaker) Import(values []interfaces.Univalue) { this.transitions = values }
// func (this *LocalSnapshotMaker) Clear()                              { this.transitions = this.transitions[:0] }

// func (this *LocalSnapshotMaker) Make(ccurl *concurrenturl.ConcurrentUrl, _ []uint32) interface{} {
// 	return ccurl.Snapshot(this.transitions)
// }
