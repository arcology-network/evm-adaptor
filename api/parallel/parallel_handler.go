package u256

import (
	"encoding/hex"

	evmcommon "github.com/arcology-network/evm/common"
	apicommon "github.com/arcology-network/vm-adaptor/api/common"
)

// APIs under the concurrency namespace
type ParallelHandler struct {
	api  apicommon.ContextInfoInterface
	jobs [][]byte
}

func NewParallelHandler(api apicommon.ContextInfoInterface) *ParallelHandler {
	return &ParallelHandler{
		api: api,
	}
}

func (this *ParallelHandler) Address() [20]byte {
	return [20]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0x90}
}

func (this *ParallelHandler) Call(caller evmcommon.Address, input []byte, origin evmcommon.Address, nonce uint64) ([]byte, bool) {
	signature := [4]byte{}
	copy(signature[:], input)

	switch signature { // bf 22 6c 78
	case [4]byte{0xbf, 0x22, 0x6c, 0x78}:
		return this.newJob(caller, input[4:])

	case [4]byte{0xb6, 0xff, 0x8b, 0xd9}:
		return this.delJob(caller, input[4:])

	case [4]byte{0x02, 0xab, 0x4b, 0xbf}:
		return this.run(caller, input[4:])

	case [4]byte{0x64, 0x17, 0x43, 0x08}:
		return this.clear(caller, input[4:])
	}
	return this.unknow(caller, input[4:])
}

func (this *ParallelHandler) newJob(caller evmcommon.Address, input []byte) ([]byte, bool) {
	id := this.api.GenUUID()

	return id, true
}

func (this *ParallelHandler) delJob(caller evmcommon.Address, input []byte) ([]byte, bool) {
	id := this.api.GenUUID()
	return id, true
}

func (this *ParallelHandler) run(caller evmcommon.Address, input []byte) ([]byte, bool) {
	// id := this.api.GenUUID()
	// delta, err := abi.Decode(input, 1, &uint256.Int{}, 1, 32)
	return []byte{}, true
}

func (this *ParallelHandler) clear(caller evmcommon.Address, input []byte) ([]byte, bool) {
	// id := this.api.GenUUID()
	// delta, err := abi.Decode(input, 1, &uint256.Int{}, 1, 32)
	return []byte{}, true
}

func (this *ParallelHandler) unknow(caller evmcommon.Address, input []byte) ([]byte, bool) {
	this.api.AddLog("Unhandled function call in cumulative handler router", hex.EncodeToString(input))
	return []byte{}, false
}
