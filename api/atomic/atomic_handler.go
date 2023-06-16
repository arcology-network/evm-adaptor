package concurrency

import (
	"crypto/sha256"
	"math"
	"math/big"

	commonlibcommon "github.com/arcology-network/common-lib/common"
	commonlibtypes "github.com/arcology-network/common-lib/types"
	evmcommon "github.com/arcology-network/evm/common"
	evmcoretypes "github.com/arcology-network/evm/core/types"
	"github.com/arcology-network/vm-adaptor/abi"
	"github.com/arcology-network/vm-adaptor/common"
	eucommon "github.com/arcology-network/vm-adaptor/common"
	execution "github.com/arcology-network/vm-adaptor/execution"
	"golang.org/x/crypto/sha3"
)

// APIs under the concurrency namespace
type AtomicHandler struct {
	api eucommon.EthApiRouter
}

func NewAtomicHandler(ethApiRouter eucommon.EthApiRouter) *AtomicHandler {
	return &AtomicHandler{
		api: ethApiRouter,
	}
}

func (this *AtomicHandler) Address() [20]byte {
	return common.ATOMIC_HANDLER
}

func (this *AtomicHandler) Call(caller, callee evmcommon.Address, input []byte, origin evmcommon.Address, nonce uint64) ([]byte, bool) {
	signature := [4]byte{}
	copy(signature[:], input)

	switch signature {
	case [4]byte{0x92, 0x8a, 0x5d, 0x96}:
		return this.deferred(origin, input[4:])

	case [4]byte{0xbb, 0x07, 0xe8, 0x5d}: // bb 07 e8 5d
		return this.uuid(caller, callee, input[4:])
	}

	return []byte{}, false
}

func (this *AtomicHandler) deferred(origin evmcommon.Address, input []byte) ([]byte, bool) {
	if this.api.VM().ArcologyNetworkAPIs.Depth() > 4 {
		return []byte{}, false
	}

	gasLimit, err := abi.DecodeTo(input, 0, uint64(0), 1, 32)
	if err != nil {
		return []byte{}, false
	}

	calleeAddr, err := abi.DecodeTo(input, 1, [20]byte{}, 1, 20)
	if err != nil {
		return []byte{}, false
	}

	funCallData, err := abi.DecodeTo(input, 2, []byte{}, 2, math.MaxUint32)
	if err != nil {
		return []byte{}, false
	}

	hierarchy := this.api.VM().ArcologyNetworkAPIs.CallHierarchy()
	groupBy := sha3.Sum256(commonlibcommon.Flatten(commonlibcommon.Reverse[[]byte](&hierarchy)))

	addr := evmcommon.Address(calleeAddr)
	evmMsg := evmcoretypes.NewMessage(
		common.ATOMIC_HANDLER, // From the system account
		&addr,
		0,
		big.NewInt(0),
		gasLimit,
		big.NewInt(1),
		funCallData,
		nil,
		false,
	)

	msg := &execution.StandardMessage{
		TxHash:  sha256.Sum256(funCallData),
		GroupBy: groupBy,
		Native:  &evmMsg,
		Source:  commonlibtypes.TX_SOURCE_DEFERRED,
	}

	this.api.SetReserved(msg) // System address pays for it
	return []byte{}, true
}

func (this *AtomicHandler) uuid(caller, callee evmcommon.Address, input []byte) ([]byte, bool) {
	return this.api.GenUUID(), true
}
