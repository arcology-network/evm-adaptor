// KernelAPI provides system level function calls supported by arcology platform.
package common

import (
	"math/big"

	"github.com/arcology-network/concurrenturl"
	"github.com/arcology-network/evm/common"

	evmCommon "github.com/arcology-network/evm/common"
	"github.com/arcology-network/evm/consensus"
	"github.com/arcology-network/evm/core/types"
	evmcoretypes "github.com/arcology-network/evm/core/types"
	"github.com/arcology-network/evm/core/vm"
)

type EthApiRouter interface {
	Origin() evmCommon.Address
	Ccurl() *concurrenturl.ConcurrentUrl
	New(common.Hash, uint32, uint8, *concurrenturl.ConcurrentUrl) EthApiRouter
	Coinbase() evmCommon.Address

	SetEU(interface{})
	GetEU() interface{}

	Message() *evmcoretypes.Message
	VM() *vm.EVM

	GetReserved() interface{}
	SetReserved(interface{})

	Depth() uint8
	AddLog(key, value string)
	Call(caller, callee evmCommon.Address, input []byte, origin evmCommon.Address, nonce uint64, blockhash evmCommon.Hash) (bool, []byte, bool)
	SetContext(evmCommon.Hash, *big.Int, uint32)

	TxIndex() uint32
	TxHash() [32]byte

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
	Call(evmCommon.Address, evmCommon.Address, []byte, evmCommon.Address, uint64) ([]byte, bool)
}
