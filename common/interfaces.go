// KernelAPI provides system level function calls supported by arcology platform.
package common

import (
	"math/big"

	"github.com/arcology-network/concurrenturl"
	"github.com/arcology-network/evm/common"
	ethcommon "github.com/arcology-network/evm/common"
	"github.com/arcology-network/evm/consensus"
	"github.com/arcology-network/evm/core/types"
	"github.com/arcology-network/evm/core/vm"
)

type ConcurrentApiRouterInterface interface {
	Origin() ethcommon.Address
	Ccurl() *concurrenturl.ConcurrentUrl
	New(common.Hash, uint32, *concurrenturl.ConcurrentUrl, uint8) ConcurrentApiRouterInterface
	Coinbase() ethcommon.Address

	Depth() uint8
	SetEU(interface{})
	VM() *vm.EVM
	AddLog(key, value string)
	Call(caller, callee ethcommon.Address, input []byte, origin ethcommon.Address, nonce uint64, blockhash ethcommon.Hash) (bool, []byte, bool)
	Prepare(ethcommon.Hash, *big.Int, uint32)
	SetCallContext(interface{})

	TxIndex() uint32
	TxHash() [32]byte

	GenCtrnUID() []byte
	GenElemUID() uint64
}

type ILog interface {
	GetKey() string
	GetValue() string
}

type ChainContext interface {
	Engine() consensus.Engine                    // Engine retrieves the chain's consensus engine.
	GetHeader(common.Hash, uint64) *types.Header // GetHeader returns the hash corresponding to their hash.
}
