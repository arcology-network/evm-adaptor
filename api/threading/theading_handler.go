package threading

import (
	"math"

	"github.com/arcology-network/common-lib/codec"
	evmcommon "github.com/arcology-network/evm/common"
	"github.com/arcology-network/vm-adaptor/abi"
	eucommon "github.com/arcology-network/vm-adaptor/common"
)

// APIs under the concurrency namespace
type TheadingHandler struct {
	api       eucommon.ConcurrentApiRouterInterface
	jobQueues map[string]*Queue
}

func NewThreadingHandler(apiRounter eucommon.ConcurrentApiRouterInterface) *TheadingHandler {
	return &TheadingHandler{
		api:       apiRounter,
		jobQueues: map[string]*Queue{},
	}
}

func (this *TheadingHandler) Address() [20]byte {
	return [20]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0x90}
}

func (this *TheadingHandler) Call(caller, callee evmcommon.Address, input []byte, origin evmcommon.Address, nonce uint64) ([]byte, bool) {
	signature := [4]byte{}
	copy(signature[:], input)

	switch signature { // bf 22 6c 78

	case [4]byte{0x58, 0x16, 0xc4, 0x25}:
		return this.new(caller, callee, input[4:])

	case [4]byte{0x9b, 0xb5, 0x52, 0xd1}:
		return this.add(caller, callee, input[4:])

	case [4]byte{0x84, 0x67, 0x3c, 0xc9}:
		return this.length(input[4:])

	case [4]byte{0x3a, 0x27, 0x65, 0x23}: //3a 27 65 23
		return this.run(caller, callee, input[4:])

	case [4]byte{0x4d, 0xd4, 0x9a, 0xb4}: // 4d d4 9a b4
		return this.get(input[4:])

	case [4]byte{0x5e, 0x1d, 0x05, 0x4d}: // 5e 1d 05 4d
		return this.clear(input[4:])

		// case [4]byte{0xb4, 0x8f, 0xb6, 0xcf}:
		// 	return this.error(input[4:])

	}

	return []byte{}, false
}

func (this *TheadingHandler) new(caller, callee evmcommon.Address, input []byte) ([]byte, bool) {
	threads, err := abi.DecodeTo(input, 0, uint8(1), 1, 32)
	if err != nil {
		return []byte{}, false
	}

	id := this.api.GenCtrnUID()
	this.jobQueues[string(id)] = NewJobQueue(threads)
	return id, true // Create a new container
}

func (this *TheadingHandler) add(caller, callee evmcommon.Address, input []byte) ([]byte, bool) {
	if len(input) < 4 {
		return []byte{}, false
	}

	id := this.ParseID(input)
	if len(id) == 0 || this.jobQueues[id] == nil {
		return []byte{}, false
	}

	// fmt.Println(input)
	rawAddr, err := abi.DecodeTo(input, 1, [20]byte{}, 1, 32)
	if err != nil {
		return []byte{}, false
	}
	calleeAddr := evmcommon.BytesToAddress(rawAddr[:]) // Callee contract

	funCall, err := abi.DecodeTo(input, 2, []byte{}, 2, math.MaxUint32)
	if err != nil {
		return []byte{}, false
	}

	return []byte{}, this.jobQueues[id].Add(this.api.Origin(), calleeAddr, funCall)
}

func (this *TheadingHandler) clear(input []byte) ([]byte, bool) {
	id := this.ParseID(input)
	if len(id) == 0 || this.jobQueues[id] == nil {
		return []byte{}, false
	}

	this.jobQueues[id].Clear()
	return []byte{}, true
}

func (this *TheadingHandler) length(input []byte) ([]byte, bool) {
	id := this.ParseID(input)
	if len(id) == 0 || this.jobQueues[id] == nil {
		return []byte{}, false
	}

	v, err := abi.Encode(this.jobQueues[id].Length())
	return v, err == nil
}

func (this *TheadingHandler) error(input []byte) ([]byte, bool) {
	id := this.ParseID(input)
	if len(id) == 0 || this.jobQueues[id] == nil {
		return []byte{}, false
	}

	if idx, err := abi.DecodeTo(input, 1, uint64(0), 1, 32); err == nil {
		if item := this.jobQueues[id].At(idx); item != nil {
			buffer, err := abi.Encode(codec.String(item.prechkErr.Error() + item.prechkErr.Error()).Clone().(codec.String).ToBytes())
			return buffer, err == nil
		}
	}
	return []byte{}, false
}

func (this *TheadingHandler) run(caller, callee evmcommon.Address, input []byte) ([]byte, bool) {
	id := this.ParseID(input)
	if len(id) == 0 || this.jobQueues[id] == nil {
		return []byte{}, false
	}

	return []byte{}, this.jobQueues[id].Run(this.api)
}

func (this *TheadingHandler) get(input []byte) ([]byte, bool) {
	id := this.ParseID(input)
	if len(id) == 0 || this.jobQueues[id] == nil {
		return []byte{}, false
	}

	if idx, err := abi.DecodeTo(input, 1, uint64(0), 1, 32); err == nil {
		if item := this.jobQueues[id].At(idx); item != nil {
			return item.result.ReturnData, item.result.Err == nil
		}
	}
	return []byte{}, false
}

// Build the container path
func (this *TheadingHandler) ParseID(input []byte) string {
	id, _ := abi.DecodeTo(input, 0, []byte{}, 2, 32) // max 32 bytes                                                                          // container ID
	return string(id)                                // unique ID
}
