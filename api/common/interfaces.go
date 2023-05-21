package common

import (
	evmCommon "github.com/arcology-network/evm/common"
)

type ConcurrentApiHandlerInterface interface {
	Address() [20]byte
	Call(evmCommon.Address, evmCommon.Address, []byte, evmCommon.Address, uint64) ([]byte, bool)
}

const (
	MAX_RECURSIION_DEPTH = uint8(3)
)
