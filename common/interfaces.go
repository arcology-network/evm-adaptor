// KernelAPI provides system level function calls supported by arcology platform.
package common

import (
	"github.com/arcology-network/concurrenturl"
	"github.com/arcology-network/concurrenturl/interfaces"
	"github.com/arcology-network/evm/common"

	evmcommon "github.com/arcology-network/evm/common"
	"github.com/arcology-network/evm/consensus"
	"github.com/arcology-network/evm/core/types"
	"github.com/arcology-network/evm/core/vm"
)

type ApiCallHandler interface {
	Address() [20]byte
	Call([20]byte, [20]byte, []byte, [20]byte, uint64) ([]byte, bool, int64)
}

type EthApiRouter interface {
	Origin() evmcommon.Address
	Ccurl() *concurrenturl.ConcurrentUrl
	New(*concurrenturl.ConcurrentUrl, interface{}) EthApiRouter
	Coinbase() evmcommon.Address

	StateFilter() StateFilter

	SetEU(interface{})
	GetEU() interface{}
	VM() *vm.EVM
	Schedule() interface{}

	GetReserved() interface{}
	SetReserved(interface{})

	CheckRuntimeConstrains() bool

	Depth() uint8
	AddLog(key, value string)
	Call(caller, callee [20]byte, input []byte, origin [20]byte, nonce uint64, blockhash evmcommon.Hash) (bool, []byte, bool, int64)

	GetSerialNum(int) uint64
	UUID() []byte
	ElementUID() []byte
}

type StateFilter interface {
	Raw() []interfaces.Univalue
	ByType() ([]interfaces.Univalue, []interfaces.Univalue)
	AddToIgnore(addr string)
}

type ILog interface {
	GetKey() string
	GetValue() string
}

type ChainContext interface {
	Engine() consensus.Engine                    // Engine retrieves the chain's consensus engine.
	GetHeader(common.Hash, uint64) *types.Header // GetHeader returns the hash corresponding to their hash.
}
