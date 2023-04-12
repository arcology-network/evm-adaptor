package concurrentcontainer

import (
	"encoding/hex"
	"math"

	"github.com/arcology-network/common-lib/codec"
	"github.com/arcology-network/common-lib/types"

	commutative "github.com/arcology-network/concurrenturl/v2/type/commutative"
	"github.com/arcology-network/concurrenturl/v2/type/noncommutative"
	evmcommon "github.com/arcology-network/evm/common"
	abi "github.com/arcology-network/vm-adaptor/abi"

	apicommon "github.com/arcology-network/vm-adaptor/api/common"
	"github.com/holiman/uint256"
)

// APIs under the concurrency namespace
type Container struct {
	api       apicommon.ContextInfoInterface
	connector *apicommon.CcurlConnector
}

func NewContainer(api apicommon.ContextInfoInterface) *Container {
	return &Container{
		api:       api,
		connector: apicommon.NewCCurlConnector("/storage/containers/", api.TxHash(), api.TxIndex(), api.Ccurl()),
	}
}

func (this *Container) Address() [20]byte {
	return [20]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0x84}
}

func (this *Container) Call(caller evmcommon.Address, input []byte, origin evmcommon.Address, nonce uint64) ([]byte, bool) {
	signature := [4]byte{}
	copy(signature[:], input)

	switch signature {
	case [4]byte{0x5f, 0x7e, 0xdb, 0x9a}:
		return this.New(caller, input[4:])

	case [4]byte{0xb8, 0xeb, 0xfc, 0x10}:
		return this.Push(caller, input[4:], origin, nonce)

	case [4]byte{0xbe, 0x16, 0x6f, 0xcf}:
		return this.Length(caller, input[4:])

	case [4]byte{0xac, 0xc3, 0x86, 0x27}:
		return this.Get(caller, input[4:])

	case [4]byte{0x89, 0x0e, 0xee, 0xb1}:
		return this.Pop(caller, input[4:])

	case [4]byte{0x6a, 0x79, 0x6b, 0xfa}:
		return this.Set(caller, input[4:])
	}

	return this.Unknow(caller, input[4:])
}

func (this *Container) Unknow(caller evmcommon.Address, input []byte) ([]byte, bool) {
	this.api.AddLog("Unhandled function call", hex.EncodeToString(input))
	return []byte{}, false
}

func (this *Container) New(caller evmcommon.Address, input []byte) ([]byte, bool) {
	// elemType := int(input[31]) // Data type should only take one byte.
	id := this.api.GenUUID() // Generate a uuid for the container
	return id[:], this.connector.New(types.Address(codec.Bytes20(caller).Hex()), hex.EncodeToString(id), 0)
}

// Get the number of elements in the container
func (this *Container) Length(caller evmcommon.Address, input []byte) ([]byte, bool) {
	path := this.buildPath(caller, input) // Container path
	if meta, err := this.api.Ccurl().Read(this.api.TxIndex(), path); err == nil {
		if encoded, err := abi.Encode(uint256.NewInt(uint64(meta.(*commutative.Meta).Length()))); err == nil {
			return encoded, true
		}
	}
	return []byte{}, false
}

func (this *Container) Get(caller evmcommon.Address, input []byte) ([]byte, bool) {
	path := this.buildPath(caller, input) // Build container path
	idx, _ := abi.Decode(input, 1, uint64(0), 1, 32)

	if value, err := this.api.Ccurl().ReadAt(this.api.TxIndex(), path, idx.(uint64)); value == nil || err != nil {
		return []byte{}, false
	} else {
		if encoded, err := abi.Encode(value.(*noncommutative.Bytes).Data()); err == nil { // Encode the result
			offset := [32]byte{}
			offset[len(offset)-1] = uint8(len(offset))
			encoded = append(offset[:], encoded...)
			return encoded, true
		}
	}
	return []byte{}, false
}

func (this *Container) Set(caller evmcommon.Address, input []byte) ([]byte, bool) {
	path := this.buildPath(caller, input) // Build container path
	idx, err := abi.Decode(input, 1, uint64(0), 1, 32)
	if err != nil {
		return []byte{}, false
	}

	bytes, err := abi.Decode(input, 2, []byte{}, 2, math.MaxInt)
	if bytes == nil || err != nil {
		return []byte{}, false
	}

	value := noncommutative.NewBytes(bytes.([]byte))
	if err := this.api.Ccurl().WriteAt(this.api.TxIndex(), path, idx.(uint64), value); err == nil {
		return []byte{}, true
	}
	return []byte{}, false
}

// Push a new element into the container
func (this *Container) Push(caller evmcommon.Address, input []byte, origin evmcommon.Address, nonce uint64) ([]byte, bool) {
	path := this.buildPath(caller, input) // Container path

	buffer := codec.Bytes32(this.api.TxHash()).Clone()
	codec.Uint64(this.api.SUID()).EncodeToBuffer(buffer[len(buffer)-8:])
	key := path + hex.EncodeToString(buffer[:])

	value, err := abi.Decode(input, 1, []byte{}, 2, math.MaxInt)
	if value == nil || err != nil {
		return []byte{}, false
	}

	err = this.api.Ccurl().Write(this.api.TxIndex(), key, noncommutative.NewBytes(value.([]byte)))
	return []byte{}, err == nil
}

func (this *Container) Pop(caller evmcommon.Address, input []byte) ([]byte, bool) {
	path := this.buildPath(caller, input) // Build container path
	if value, err := this.api.Ccurl().PopBack(this.api.TxIndex(), path); err != nil {
		return []byte{}, false
	} else {
		if value != nil {
			encoded, err := abi.Encode([]byte(value.(*noncommutative.Bytes).Data()))

			offset := [32]byte{}
			offset[len(offset)-1] = uint8(len(offset))
			encoded = append(offset[:], encoded...)
			return encoded, err == nil
		}
	}
	return []byte{}, true
}

// Build the container path
func (this *Container) buildPath(caller evmcommon.Address, input []byte) string {
	id, _ := abi.Decode(input, 0, []byte{}, 2, 32)                                                         // max 32 bytes                                                                          // container ID
	return this.connector.Key(types.Address(codec.Bytes20(caller).Hex()), hex.EncodeToString(id.([]byte))) // unique ID
}
