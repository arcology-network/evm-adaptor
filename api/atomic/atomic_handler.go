package concurrency

import (
	evmcommon "github.com/arcology-network/evm/common"
	"github.com/arcology-network/vm-adaptor/common"
	eucommon "github.com/arcology-network/vm-adaptor/common"
	"github.com/arcology-network/vm-adaptor/execution"
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

func (this *AtomicHandler) Call(caller, callee [20]byte, input []byte, origin [20]byte, nonce uint64) ([]byte, bool, int64) {
	signature := [4]byte{}
	copy(signature[:], input)

	switch signature {
	case [4]byte{0xb6, 0x26, 0x54, 0xfb}:
		return this.singleton(origin, input[4:])

	case [4]byte{0xbb, 0x07, 0xe8, 0x5d}: // bb 07 e8 5d
		return this.uuid(caller, callee, input[4:])
	}

	return []byte{}, false, 0
}

func (this *AtomicHandler) singleton(origin evmcommon.Address, input []byte) ([]byte, bool, int64) {
	schedule := this.api.Schedule().(*execution.Schedule)
	if schedule != nil {
		schedule.IsLast(this.api.TxHash(), this.api.Message())
	}
	return []byte{}, false, 0

}

func (this *AtomicHandler) uuid(caller, callee evmcommon.Address, input []byte) ([]byte, bool, int64) {
	return this.api.GenUUID(), true, 0
}
