package concurrency

import (
	"strings"

	"github.com/arcology-network/concurrenturl/noncommutative"
	evmcommon "github.com/arcology-network/evm/common"
	"github.com/arcology-network/vm-adaptor/abi"
	apicommon "github.com/arcology-network/vm-adaptor/api/common"
	"github.com/arcology-network/vm-adaptor/common"
	eucommon "github.com/arcology-network/vm-adaptor/common"
	execution "github.com/arcology-network/vm-adaptor/execution"
)

// APIs under the concurrency namespace
type RuntimeHandler struct {
	api       eucommon.EthApiRouter
	connector *apicommon.CcurlConnector
}

func NewRuntimeHandler(ethApiRouter eucommon.EthApiRouter) *RuntimeHandler {
	return &RuntimeHandler{
		api:       ethApiRouter,
		connector: apicommon.NewCCurlConnector("/native/local/", ethApiRouter, ethApiRouter.Ccurl()),
	}
}

func (this *RuntimeHandler) Address() [20]byte {
	return common.RUNTIME_HANDLER
}

func (this *RuntimeHandler) Call(caller, callee [20]byte, input []byte, origin [20]byte, nonce uint64) ([]byte, bool, int64) {
	signature := [4]byte{}
	copy(signature[:], input)

	switch signature {
	case [4]byte{0xb6, 0x26, 0x54, 0xfb}:
		return this.singleton(origin, input[4:])

	case [4]byte{0xbb, 0x07, 0xe8, 0x5d}: // bb 07 e8 5d
		return this.uuid(caller, callee, input[4:])

	case [4]byte{0x4c, 0x45, 0xdc, 0x4a}: // bb 07 e8 5d
		return this.localize(caller, callee, input[4:])
	}

	return []byte{}, false, 0
}

func (this *RuntimeHandler) singleton(origin evmcommon.Address, input []byte) ([]byte, bool, int64) {
	// schedule := this.api.Schedule()
	// if schedule != nil {
	// 	schedule.IsLast(this.api.TxHash(), this.api.Message())
	// }
	return []byte{}, false, 0
}

func (this *RuntimeHandler) uuid(caller, callee evmcommon.Address, input []byte) ([]byte, bool, int64) {
	return this.api.UUID(), true, 0
}

func (this *RuntimeHandler) localize(caller, callee evmcommon.Address, input []byte) ([]byte, bool, int64) {
	if this.api.GetEU().(*execution.EU).Message().Native.To != nil {
		return []byte{}, false, 0 // Only in constructor
	}

	_, err := abi.DecodeTo(input, 0, uint64(0), 1, 32) // max 32 bytes
	if err != nil {
		return []byte{}, false, 0
	}

	path := this.connector.Key(caller) // unique ID
	path = strings.TrimSuffix(path, "/")

	value := noncommutative.NewBytes([]byte{})
	if _, err := this.api.Ccurl().Write(uint32(this.api.GetEU().(*execution.EU).Message().ID), path, value, false); err == nil {
		return []byte{}, true, 0
	}
	return []byte{}, false, 0
}
