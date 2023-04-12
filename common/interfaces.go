// KernelAPI provides system level function calls supported by arcology platform.
package common

import (
	"math/big"

	ethcommon "github.com/arcology-network/evm/common"
)

type ConcurrentApiInterface interface {
	AddLog(key, value string)
	Call(caller, callee ethcommon.Address, input []byte, origin ethcommon.Address, nonce uint64, blockhash ethcommon.Hash) (bool, []byte, bool)
	Prepare(ethcommon.Hash, *big.Int, uint32)
}

type ILog interface {
	GetKey() string
	GetValue() string
}
