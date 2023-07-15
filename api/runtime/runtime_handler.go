package concurrency

import (
	"encoding/hex"
	"math"

	evmcommon "github.com/arcology-network/evm/common"
	"github.com/arcology-network/vm-adaptor/abi"
	apicommon "github.com/arcology-network/vm-adaptor/api/common"
	"github.com/arcology-network/vm-adaptor/common"
	eucommon "github.com/arcology-network/vm-adaptor/common"
	"github.com/arcology-network/vm-adaptor/execution"
)

// APIs under the concurrency namespace
type RuntimeHandler struct {
	api       eucommon.EthApiRouter
	connector *apicommon.CcurlConnector
}

func NewHandler(ethApiRouter eucommon.EthApiRouter) *RuntimeHandler {
	return &RuntimeHandler{
		api:       ethApiRouter,
		connector: apicommon.NewCCurlConnector("/native/local/", ethApiRouter, ethApiRouter.Ccurl()),
	}
}

func (this *RuntimeHandler) Address() [20]byte {
	return common.RUNTIME_HANDLER
}

func (this *RuntimeHandler) Call(caller, callee [20]byte, input []byte, _ [20]byte, nonce uint64) ([]byte, bool, int64) {
	signature := [4]byte{}
	copy(signature[:], input)

	switch signature {
	// case [4]byte{0xd3, 0x01, 0xe8, 0xfe}: // d3 01 e8 fe
	// 	return this.undo(caller, input[4:])

	case [4]byte{0x79, 0xfc, 0x09, 0xa2}: // 79 fc 09 a2
		return this.exists(caller, input[4:])

	case [4]byte{0x64, 0x23, 0xdb, 0x34}: // d3 01 e8 fe
		return this.reset(caller, input[4:])

	case [4]byte{0xbb, 0x07, 0xe8, 0x5d}: // bb 07 e8 5d
		return this.uuid(caller, callee, input[4:])
	}

	return []byte{}, false, 0
}

// func (this *RuntimeHandler) undo(caller evmcommon.Address, input []byte) ([]byte, bool, int64) {
// 	if this.api.VM().ArcologyNetworkAPIs.IsInConstructor() {
// 		this.api.StateFilter().AddToAutoReversion(hex.EncodeToString(caller[:]))
// 		return []byte{}, true, 0
// 	}
// 	return []byte{}, true, 0
// }

func (this *RuntimeHandler) exists(caller evmcommon.Address, input []byte) ([]byte, bool, int64) {
	path := this.connector.Key(caller) // Build container path
	if len(path) == 0 {
		return []byte{}, false, 0
	}

	if key, err := abi.DecodeTo(input, 0, []byte{}, 1, math.MaxInt64); err != nil {
		v, _ := this.api.Ccurl().Read(uint32(this.api.GetEU().(*execution.EU).Message().ID), path+hex.EncodeToString(key))
		if encoded, err := abi.Encode(v != nil); err == nil {
			return encoded, true, 0
		}
	}
	return []byte{}, false, 0
}

func (this *RuntimeHandler) reset(caller evmcommon.Address, input []byte) ([]byte, bool, int64) {
	this.api.StateFilter().RemoveByAddress(hex.EncodeToString(caller[:]))
	return []byte{}, true, 0
}

func (this *RuntimeHandler) uuid(caller, callee evmcommon.Address, input []byte) ([]byte, bool, int64) {
	return this.api.ElementUID(), true, 0
}
