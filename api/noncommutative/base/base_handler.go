package concurrentcontainer

import (
	"encoding/hex"
	"math"

	"github.com/arcology-network/common-lib/codec"
	"github.com/arcology-network/common-lib/types"

	"github.com/arcology-network/concurrenturl/interfaces"
	"github.com/arcology-network/concurrenturl/noncommutative"
	evmcommon "github.com/arcology-network/evm/common"
	abi "github.com/arcology-network/vm-adaptor/abi"
	"github.com/arcology-network/vm-adaptor/common"
	eucommon "github.com/arcology-network/vm-adaptor/common"
	"github.com/arcology-network/vm-adaptor/execution"

	apicommon "github.com/arcology-network/vm-adaptor/api/common"
	"github.com/holiman/uint256"
)

// APIs under the concurrency namespace
type BytesHandlers struct {
	api       eucommon.EthApiRouter
	connector *apicommon.CcurlConnector
}

func NewNoncommutativeBytesHandlers(api eucommon.EthApiRouter) *BytesHandlers {
	return &BytesHandlers{
		api:       api,
		connector: apicommon.NewCCurlConnector("/containers/", api, api.Ccurl()),
	}
}

func (this *BytesHandlers) Address() [20]byte {
	return common.BYTES_HANDLER
}

func (this *BytesHandlers) Call(caller, callee [20]byte, input []byte, origin [20]byte, nonce uint64) ([]byte, bool, int64) {
	signature := [4]byte{}
	copy(signature[:], input)

	switch signature {
	case [4]byte{0xcd, 0xbf, 0x60, 0x8d}: //cd bf 60 8d
		return this.new(caller, input[4:])

	case [4]byte{0xee, 0xb8, 0xa8, 0xd3}:
		return this.peek(caller, input[4:])

	case [4]byte{0xe7, 0x71, 0xee, 0x0d}: // e771ee0d
		return this.push(caller, input[4:], origin, nonce)

	case [4]byte{0x84, 0x67, 0x3c, 0xc9}: // 84 67 3c c9
		return this.length(caller, input[4:])

	case [4]byte{0x4d, 0xd4, 0x9a, 0xb4}: // 4d d4 9a b4
		return this.get(caller, input[4:])

	case [4]byte{0xa4, 0xec, 0xe5, 0x2c}: // a4 ec e5 2c
		return this.pop(caller, input[4:])

	case [4]byte{0x5e, 0x1d, 0x05, 0x4d}: // 5e 1d 05 4d
		return this.clear(caller, input[4:])

	case [4]byte{0x4c, 0x51, 0xa8, 0x8f}: // 4c51a88f
		return this.set(caller, input[4:])
	}

	return []byte{}, false, 0 // unknown
}

func (this *BytesHandlers) new(caller evmcommon.Address, input []byte) ([]byte, bool, int64) {
	id := this.api.UUID() // Generate a uuid for the container
	return id[:], this.connector.New(
		uint32(this.api.GetEU().(*execution.EU).Message().ID),
		types.Address(codec.Bytes20(caller).Hex()),
		hex.EncodeToString(id),
	), 0 // Create a new container
}

// Get the number of elements in the container
func (this *BytesHandlers) length(caller evmcommon.Address, input []byte) ([]byte, bool, int64) {
	path, err := this.buildPath(caller, input) // BytesHandlers path
	if len(path) == 0 || err != nil {
		return []byte{}, false, 0
	}

	if path, _ := this.api.Ccurl().Read(uint32(this.api.GetEU().(*execution.EU).Message().ID), path); path != nil {
		if encoded, err := abi.Encode(uint256.NewInt(uint64(len(path.([]string))))); err == nil {
			return encoded, true, 0
		}
	}
	return []byte{}, false, 0
}

// Get the intial length of the container
func (this *BytesHandlers) peek(caller evmcommon.Address, input []byte) ([]byte, bool, int64) {
	path, err := this.buildPath(caller, input) // BytesHandlers path
	if len(path) == 0 || err != nil {
		return []byte{}, false, 0
	}

	typedv, fees := this.api.Ccurl().PeekCommitted(path)
	if typedv != nil && err == nil {
		type measurable interface{ Length() int }
		numKeys := uint64(typedv.(interfaces.Type).Value().(measurable).Length())
		if encoded, err := abi.Encode(uint256.NewInt(numKeys)); err == nil {
			return encoded, true, int64(fees)
		}
	}
	return []byte{}, false, int64(fees)
}

func (this *BytesHandlers) get(caller evmcommon.Address, input []byte) ([]byte, bool, int64) {
	path, err := this.buildPath(caller, input) // Build container path
	if len(path) == 0 || err != nil {
		return []byte{}, false, 0
	}

	idx, err := abi.DecodeTo(input, 1, uint64(0), 1, 32)
	if err != nil {
		return []byte{}, false, 0
	}

	if value, _, err := this.api.Ccurl().ReadAt(uint32(this.api.GetEU().(*execution.EU).Message().ID), path, idx); err != nil && value == nil {
		return []byte{}, false, 0
	} else {
		if encoded, err := abi.Encode(value.([]byte)); err == nil { // Encode the result
			offset := [32]byte{}
			offset[len(offset)-1] = uint8(len(offset))
			encoded = append(offset[:], encoded...)
			return encoded, true, 0
		}
	}
	return []byte{}, false, 0
}

func (this *BytesHandlers) set(caller evmcommon.Address, input []byte) ([]byte, bool, int64) {
	path, err := this.buildPath(caller, input) // Build container path
	if len(path) == 0 || err != nil {
		return []byte{}, false, 0
	}

	idx, err := abi.DecodeTo(input, 1, uint64(0), 1, 32)
	if err != nil {
		return []byte{}, false, 0
	}

	bytes, err := abi.Decode(input, 2, []byte{}, 2, math.MaxInt)
	if bytes == nil || err != nil {
		return []byte{}, false, 0
	}

	value := noncommutative.NewBytes(bytes.([]byte))
	if _, err := this.api.Ccurl().WriteAt(uint32(this.api.GetEU().(*execution.EU).Message().ID), path, idx, value); err == nil {
		return []byte{}, true, 0
	}
	return []byte{}, false, 0
}

// Push a new element into the container
func (this *BytesHandlers) push(caller evmcommon.Address, input []byte, origin evmcommon.Address, nonce uint64) ([]byte, bool, int64) {
	path, err := this.buildPath(caller, input) // BytesHandlers path
	if len(path) == 0 || err != nil {
		return []byte{}, false, 0
	}

	key := path + string(this.api.ElementUID())
	value, err := abi.Decode(input, 1, []byte{}, 2, math.MaxInt)
	if value == nil || err != nil {
		return []byte{}, false, 0
	}

	_, err = this.api.Ccurl().Write(uint32(this.api.GetEU().(*execution.EU).Message().ID), key, noncommutative.NewBytes(value.([]byte)))
	return []byte{}, err == nil, 0
}

func (this *BytesHandlers) pop(caller evmcommon.Address, input []byte) ([]byte, bool, int64) {
	path, err := this.buildPath(caller, input) // Build container path
	if len(path) == 0 || err != nil {
		return []byte{}, false, 0
	}

	if value, _, err := this.api.Ccurl().PopBack(uint32(this.api.GetEU().(*execution.EU).Message().ID), path); err != nil {
		return []byte{}, false, 0
	} else {
		if value != nil {
			encoded, err := abi.Encode(value.([]byte))

			offset := [32]byte{}
			offset[len(offset)-1] = uint8(len(offset))
			encoded = append(offset[:], encoded...)
			return encoded, err == nil, 0
		}
	}
	return []byte{}, true, 0
}

func (this *BytesHandlers) clear(caller evmcommon.Address, input []byte) ([]byte, bool, int64) {
	for {
		if _, ok, _ := this.pop(caller, input); !ok {
			break
		}
	}
	return []byte{}, true, 0
}

// Build the container path
func (this *BytesHandlers) buildPath(caller evmcommon.Address, input []byte) (string, error) {
	id, err := abi.Decode(input, 0, []byte{}, 2, 32)                                                            // max 32 bytes                                                                          // container ID
	return this.connector.Key(types.Address(codec.Bytes20(caller).Hex()), hex.EncodeToString(id.([]byte))), err // unique ID
}
