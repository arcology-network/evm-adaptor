// KernelAPI provides system level function calls supported by arcology platform.
package interfaces

import (
	"math/big"

	"github.com/arcology-network/concurrenturl/univalue"
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
	Origin() evmcommon.Address
	// WriteCache() *StorageCommitter.ConcurrentUrl
	WriteCache() interface{}
	// DataReader() interface{}
	SetReadOnlyDataSource(interface{})
	New(interface{}, interface{}) EthApiRouter
	Coinbase() evmcommon.Address

	StateFilter() StateFilter

	SetEU(interface{})
	GetEU() interface{}
	VM() interface{} //*vm.EVM
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
	Raw() []*univalue.Univalue
	ByType() ([]*univalue.Univalue, []*univalue.Univalue)
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

type EU interface {
	GasPrice() *big.Int
	Message() interface{}
	VM() interface{} //*vm.EVM
	ID() uint32
	TxHash() [32]byte
	Origin() [20]byte
	Coinbase() [20]byte
}

type JobSequence interface {
	GetID() uint32
	New(uint32, EthApiRouter) JobSequence
	DeriveNewHash([32]byte) [32]byte
	AppendMsg(interface{})
}

type Generation interface {
	New(uint32, uint8, []JobSequence) Generation
	Add(JobSequence) bool
	Run(EthApiRouter) []*univalue.Univalue
	JobSeqs() []JobSequence
	JobT() JobSequence
}