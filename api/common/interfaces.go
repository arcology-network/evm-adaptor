package common

import (
	"github.com/arcology-network/concurrenturl/v2"
	evmCommon "github.com/arcology-network/evm/common"
)

type ConcurrencyHandlerInterface interface {
	Address() [20]byte
	Call(evmCommon.Address, []byte, evmCommon.Address, uint64) ([]byte, bool)
}

type ContextInfoInterface interface {
	TxIndex() uint32
	TxHash() [32]byte
	Ccurl() *concurrenturl.ConcurrentUrl
	GenUUID() []byte
	SUID() uint64
	AddLog(key, value string)
}
