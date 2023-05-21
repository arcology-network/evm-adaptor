// KernelAPI provides system level function calls supported by arcology platform.
package common

import (
	evmcommon "github.com/arcology-network/evm/common"
)

var (
	Coinbase = evmcommon.BytesToAddress([]byte("coinbase"))
	Owner    = evmcommon.BytesToAddress([]byte("owner"))
	Alice    = evmcommon.BytesToAddress([]byte("user1"))
	Bob      = evmcommon.BytesToAddress([]byte("user2"))
)
