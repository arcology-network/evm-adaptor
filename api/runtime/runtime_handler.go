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

	// Contract := this.api.VM().ArcologyNetworkAPIs.CallContext.Contract
	// fmt.Print(Contract)

	switch signature {
	case [4]byte{0x4e, 0xd3, 0x88, 0x5e}: // 4e d3 88 5e
		return this.set(caller, input[4:])

	case [4]byte{0xbb, 0x07, 0xe8, 0x5d}: // bb 07 e8 5d
		return this.uuid(caller, callee, input[4:])
	}

	return []byte{}, false, 0
}

func (this *RuntimeHandler) set(caller evmcommon.Address, input []byte) ([]byte, bool, int64) {
	if this.api.VM().ArcologyNetworkAPIs.IsInConstructor() {
		this.api.StateFilter().AddToIgnore(hex.EncodeToString(caller[:]))
		return []byte{}, true, 0
	}
	return []byte{}, false, 0
}

func (this *RuntimeHandler) uuid(caller, callee evmcommon.Address, input []byte) ([]byte, bool, int64) {
	return this.api.UUID(), true, 0
}

// func (this *RuntimeHandler) localize(caller, callee evmcommon.Address, input []byte) ([]byte, bool, int64) {
// 	if this.api.GetEU().(*execution.EU).Message().Native.To != nil {
// 		return []byte{}, false, 0 // Only in constructor
// 	}

// 	_, err := abi.DecodeTo(input, 0, uint64(0), 1, 32) // max 32 bytes
// 	if err != nil {
// 		return []byte{}, false, 0
// 	}

// 	path := this.connector.Key(caller) // unique ID
// 	path = strings.TrimSuffix(path, "/")

// 	value := noncommutative.NewBytes([]byte{})
// 	if _, err := this.api.Ccurl().Write(uint32(this.api.GetEU().(*execution.EU).Message().ID), path, value, false); err == nil {
// 		return []byte{}, true, 0
// 	}
// 	return []byte{}, false, 0
// }
