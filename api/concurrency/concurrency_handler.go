package concurrency

import (
	evmcommon "github.com/arcology-network/evm/common"
	interfaces "github.com/arcology-network/vm-adaptor/interfaces"
	types "github.com/arcology-network/vm-adaptor/types"
)

// APIs under the concurrency namespace
type ConcurrencyHandler struct {
	api interfaces.EthApiRouter
}

func NewConcurrencyHandler(ethApiRounter interfaces.EthApiRouter) *ConcurrencyHandler {
	return &ConcurrencyHandler{
		api: ethApiRounter,
	}
}

func (this *ConcurrencyHandler) Address() [20]byte {
	return [20]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0xa0}
}

func (this *ConcurrencyHandler) Call(caller, callee evmcommon.Address, input []byte, origin evmcommon.Address, nonce uint64) ([]byte, bool) {
	signature := [4]byte{}
	copy(signature[:], input)

	switch signature {
	case [4]byte{0x03, 0x74, 0xa0, 0x4d}: // 03 74 a0 4d 5
		return this.deferred(caller, callee, input[4:])

	case [4]byte{0xbb, 0x07, 0xe8, 0x5d}: // bb 07 e8 5d
		return this.uuid(caller, callee, input[4:])
	}

	return []byte{}, false
}

func (this *ConcurrencyHandler) deferred(caller, callee evmcommon.Address, input []byte) ([]byte, bool) {
	if this.api.VM().ArcologyNetworkAPIs.Depth() > 2 {
		return []byte{}, false
	}

	this.api.SetReserved(types.NewDeferredCall(callee, this.api)) // System address pays for it
	return []byte{}, true
}

func (this *ConcurrencyHandler) uuid(caller, callee evmcommon.Address, input []byte) ([]byte, bool) {
	return this.api.GenUUID(), true
}
