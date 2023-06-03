package concurrency

import (
	"github.com/arcology-network/common-lib/types"
	evmcommon "github.com/arcology-network/evm/common"
	"github.com/arcology-network/vm-adaptor/abi"
	interfaces "github.com/arcology-network/vm-adaptor/interfaces"
)

// APIs under the concurrency namespace
type ConcurrencyHandler struct {
	api interfaces.ApiRouter
}

func NewConcurrencyHandler(apiRounter interfaces.ApiRouter) *ConcurrencyHandler {
	return &ConcurrencyHandler{
		api: apiRounter,
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

	targetAddr, err := abi.DecodeTo(input, 0, [20]byte{}, 1, 32)
	if err != nil {
		return []byte{}, false
	}

	targetSignature, err := abi.DecodeTo(input, 1, []byte{}, 2, 4) // Function signature only, won't take any input argument
	if err != nil {
		return []byte{}, false
	}

	txHash := this.api.TxHash()
	this.api.SetDeferred(
		&types.DeferCall{
			DeferID:    string(txHash[:]),
			CallerAddr: types.Address(caller[:]),
			TargetAddr: types.Address(targetAddr[:]),
			TargetFunc: string(targetSignature),
		},
	)
	return []byte{}, true
}

func (this *ConcurrencyHandler) uuid(caller, callee evmcommon.Address, input []byte) ([]byte, bool) {
	return this.api.GenUUID(), true
}
