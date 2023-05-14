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
	From() ethcommon.Address

	SetEU(interface{})
	VM() *vm.EVM
	AddLog(key, value string)
	Call(caller, callee ethcommon.Address, input []byte, origin ethcommon.Address, nonce uint64, blockhash ethcommon.Hash) (bool, []byte, bool)
	Prepare(ethcommon.Hash, *big.Int, uint32)
	Ccurl() *concurrenturl.ConcurrentUrl

	TxIndex() uint32
	TxHash() [32]byte

	GenUUID() []byte
	SUID() uint64

	New(common.Hash, uint32, *concurrenturl.ConcurrentUrl) ConcurrentApiRouterInterface
	Coinbase() ethcommon.Address
}

type ILog interface {
	GetKey() string
	GetValue() string
}

type ChainContext interface {
	Engine() consensus.Engine                    // Engine retrieves the chain's consensus engine.
	GetHeader(common.Hash, uint64) *types.Header // GetHeader returns the hash corresponding to their hash.
}
