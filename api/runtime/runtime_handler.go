package concurrency

import (
	"encoding/hex"

	evmcommon "github.com/arcology-network/evm/common"
	apicommon "github.com/arcology-network/vm-adaptor/api/common"
	"github.com/arcology-network/vm-adaptor/common"
	eucommon "github.com/arcology-network/vm-adaptor/common"
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

func (this *RuntimeHandler) Call(caller, callee [20]byte, input []byte, _ [20]byte, nonce uint64) ([]byte, bool, int64) {
	signature := [4]byte{}
	copy(signature[:], input)

	switch signature {
	case [4]byte{0xd8, 0x26, 0xf8, 0x8f}: // d8 26 f8 8f
		return this.reset(caller, input[4:])

	case [4]byte{0xbb, 0x07, 0xe8, 0x5d}: // bb 07 e8 5d
		return this.uuid(caller, callee, input[4:])
	}

	return []byte{}, false, 0
}

func (this *RuntimeHandler) reset(caller evmcommon.Address, input []byte) ([]byte, bool, int64) {
	if this.api.VM().ArcologyNetworkAPIs.IsInConstructor() {
		this.api.StateFilter().AddToIgnore(hex.EncodeToString(caller[:]))
		return []byte{}, true, 0
	}
	return []byte{}, false, 0
}

func (this *RuntimeHandler) uuid(caller, callee evmcommon.Address, input []byte) ([]byte, bool, int64) {
	return this.api.UUID(), true, 0
}
