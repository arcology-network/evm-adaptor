// KernelAPI provides system level function calls supported by arcology platform.
package interfaces

import (
	"github.com/ethereum/go-ethereum/common"

	evmcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/core/types"
)

type ApiCallHandler interface {
	Address() [20]byte
	Call([20]byte, [20]byte, []byte, [20]byte, uint64) ([]byte, bool, int64)
}

type EthApiRouter interface {
	GetDeployer() evmcommon.Address
	SetDeployer(evmcommon.Address)

	GetEU() interface{}
	SetEU(interface{})

	GetSchedule() interface{}
	SetSchedule(interface{})

	Origin() evmcommon.Address

	AuxDict() map[string]interface{}
	WriteCachePool() interface{}
	WriteCache() interface{}
	SetReadOnlyDataSource(interface{})
	New(interface{}, interface{}, evmcommon.Address, interface{}) EthApiRouter
	Coinbase() evmcommon.Address

	VM() interface{} //*vm.EVM

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

type ILog interface {
	GetByKey() string
	GetValue() string
}

type ChainContext interface {
	Engine() consensus.Engine                    // Engine retrieves the chain's consensus engine.
	GetHeader(common.Hash, uint64) *types.Header // GetHeader returns the hash corresponding to their hash.
}
