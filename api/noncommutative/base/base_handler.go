package concurrentcontainer

import (
	"encoding/hex"
	"math"
	"strconv"

	"github.com/arcology-network/common-lib/codec"
	"github.com/arcology-network/common-lib/types"

	"github.com/arcology-network/concurrenturl/noncommutative"
	evmcommon "github.com/arcology-network/evm/common"
	abi "github.com/arcology-network/vm-adaptor/abi"

	apicommon "github.com/arcology-network/vm-adaptor/api/common"
	eucommon "github.com/arcology-network/vm-adaptor/common"
	"github.com/holiman/uint256"
)

// APIs under the concurrency namespace
type NoncommutativeBytesHandlers struct {
	api       eucommon.ConcurrentApiRouterInterface
	connector *apicommon.CcurlConnector
}

func NewNoncommutativeBytesHandlers(api eucommon.ConcurrentApiRouterInterface) *NoncommutativeBytesHandlers {
	return &NoncommutativeBytesHandlers{
		api:       api,
		connector: apicommon.NewCCurlConnector("/containers/", api, api.Ccurl()),
	}
}

func (this *NoncommutativeBytesHandlers) Address() [20]byte {
	return [20]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0x84}
}

func (this *NoncommutativeBytesHandlers) Call(caller, callee evmcommon.Address, input []byte, origin evmcommon.Address, nonce uint64) ([]byte, bool) {
	signature := [4]byte{}
	copy(signature[:], input)

	switch signature {
	case [4]byte{0xcd, 0xbf, 0x60, 0x8d}: //cd bf 60 8d
		return this.new(caller, input[4:])

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

	return this.unknow(caller, input)
}

func (this *NoncommutativeBytesHandlers) unknow(caller evmcommon.Address, input []byte) ([]byte, bool) {
	this.api.AddLog("Unhandled function call", hex.EncodeToString(input))
	return []byte{}, false
}

func (this *NoncommutativeBytesHandlers) new(caller evmcommon.Address, input []byte) ([]byte, bool) {
	id := this.api.GenCtrnUID()                                                                          // Generate a uuid for the container
	return id[:], this.connector.New(types.Address(codec.Bytes20(caller).Hex()), hex.EncodeToString(id)) // Create a new container
}

// Get the number of elements in the container
func (this *NoncommutativeBytesHandlers) length(caller evmcommon.Address, input []byte) ([]byte, bool) {
	path := this.buildPath(caller, input) // NoncommutativeBytesHandlers path
	if path, err := this.api.Ccurl().Read(this.api.TxIndex(), path); err == nil && path != nil {
		if encoded, err := abi.Encode(uint256.NewInt(uint64(len(path.([]string))))); err == nil {
			return encoded, true
		}
	}
	return []byte{}, false
}

func (this *NoncommutativeBytesHandlers) get(caller evmcommon.Address, input []byte) ([]byte, bool) {
	idx, err := abi.DecodeTo(input, 1, uint64(0), 1, 32)
	if err != nil {
		return []byte{}, false
	}

	path := this.buildPath(caller, input) // Build container path
	if value, err := this.api.Ccurl().ReadAt(this.api.TxIndex(), path, idx); value == nil || err != nil {
		return []byte{}, false
	} else {
		if encoded, err := abi.Encode(value.([]byte)); err == nil { // Encode the result
			offset := [32]byte{}
			offset[len(offset)-1] = uint8(len(offset))
			encoded = append(offset[:], encoded...)
			return encoded, true
		}
	}
	return []byte{}, false
}

func (this *NoncommutativeBytesHandlers) set(caller evmcommon.Address, input []byte) ([]byte, bool) {
	idx, err := abi.DecodeTo(input, 1, uint64(0), 1, 32)
	if err != nil {
		return []byte{}, false
	}

	bytes, err := abi.Decode(input, 2, []byte{}, 2, math.MaxInt)
	if bytes == nil || err != nil {
		return []byte{}, false
	}

	path := this.buildPath(caller, input) // Build container path
	value := noncommutative.NewBytes(bytes.([]byte))
	if err := this.api.Ccurl().WriteAt(this.api.TxIndex(), path, idx, value); err == nil {
		return []byte{}, true
	}
	return []byte{}, false
}

// Push a new element into the container
func (this *NoncommutativeBytesHandlers) push(caller evmcommon.Address, input []byte, origin evmcommon.Address, nonce uint64) ([]byte, bool) {
	path := this.buildPath(caller, input) // NoncommutativeBytesHandlers path

	txHash := this.api.TxHash()
	key := path + hex.EncodeToString(txHash[:8]) + "-" + strconv.Itoa(int(this.api.GenElemUID()))

	value, err := abi.Decode(input, 1, []byte{}, 2, math.MaxInt)
	if value == nil || err != nil {
		return []byte{}, false
	}

	err = this.api.Ccurl().Write(this.api.TxIndex(), key, noncommutative.NewBytes(value.([]byte)))
	return []byte{}, err == nil
}

func (this *NoncommutativeBytesHandlers) pop(caller evmcommon.Address, input []byte) ([]byte, bool) {
	path := this.buildPath(caller, input) // Build container path
	if value, err := this.api.Ccurl().PopBack(this.api.TxIndex(), path); err != nil {
		return []byte{}, false
	} else {
		if value != nil {
			encoded, err := abi.Encode(value.([]byte))

			offset := [32]byte{}
			offset[len(offset)-1] = uint8(len(offset))
			encoded = append(offset[:], encoded...)
			return encoded, err == nil
		}
	}
	return []byte{}, true
}

func (this *NoncommutativeBytesHandlers) clear(caller evmcommon.Address, input []byte) ([]byte, bool) {
	for {
		if _, ok := this.pop(caller, input); !ok {
			break
		}
	}
	return []byte{}, true
}

// Build the container path
func (this *NoncommutativeBytesHandlers) buildPath(caller evmcommon.Address, input []byte) string {
	id, _ := abi.Decode(input, 0, []byte{}, 2, 32)                                                         // max 32 bytes                                                                          // container ID
	return this.connector.Key(types.Address(codec.Bytes20(caller).Hex()), hex.EncodeToString(id.([]byte))) // unique ID
}
