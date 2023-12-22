// KernelAPI provides system level function calls supported by arcology platform.
package common

import (
	"github.com/arcology-network/concurrenturl"
	"github.com/arcology-network/concurrenturl/interfaces"
	"github.com/ethereum/go-ethereum/common"

	evmcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
)

type ApiCallHandler interface {
	Address() [20]byte
	Call([20]byte, [20]byte, []byte, [20]byte, uint64) ([]byte, bool, int64)
}

type EthApiRouter interface {
	Origin() evmcommon.Address
	Ccurl() *concurrenturl.ConcurrentUrl
	SetCcurl(*concurrenturl.ConcurrentUrl)
	New(*concurrenturl.ConcurrentUrl, interface{}) EthApiRouter
	Coinbase() evmcommon.Address

	StateFilter() StateFilter

	SetEU(interface{})
	GetEU() interface{}
	VM() *vm.EVM
	Schedule() interface{}

	CheckRuntimeConstrains() bool

	DecrementDepth() uint8
	Depth() uint8
	AddLog(key, value string)
	Call(caller, callee [20]byte, input []byte, origin [20]byte, nonce uint64, blockhash evmcommon.Hash) (bool, []byte, bool, int64)

	GetSerialNum(int) uint64
	Pid() [32]byte
	UUID() []byte
	ElementUID() []byte
}

type StateFilter interface {
	Raw() []interfaces.Univalue
	ByType() ([]interfaces.Univalue, []interfaces.Univalue)
	AddToAutoReversion(addr string)
	RemoveByAddress(string)
}

type ILog interface {
	GetByKey() string
	GetValue() string
}

type ChainContext interface {
	Engine() consensus.Engine                    // Engine retrieves the chain's consensus engine.
	GetHeader(common.Hash, uint64) *types.Header // GetHeader returns the hash corresponding to their hash.
}

type MessageReader interface {
	Message() interface {
		ID() uint32
		TxHash() [32]byte
	}
}

type EUInterface interface {
	Message() *StandardMessage
	VM() *vm.EVM
	ID() uint32
	TxHash() [32]byte
	Origin() [20]byte
	Coinbase() [20]byte
}
