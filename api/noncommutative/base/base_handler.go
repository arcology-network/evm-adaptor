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
		return this.New(caller, input[4:])

	case [4]byte{0x59, 0xe0, 0x2d, 0xd7}:
		return this.Peek(caller, input[4:])

	case [4]byte{0x7d, 0xac, 0xda, 0x03}: // 7d ac da 03
		return this.Push(caller, input[4:], origin, nonce)

	case [4]byte{0x1f, 0x7b, 0x6d, 0x32}: // 1f 7b 6d 32
		return this.Length(caller, input[4:])

	case [4]byte{0x95, 0x07, 0xd3, 0x9a}: // 95 07 d3 9a
		return this.Get(caller, input[4:])

	case [4]byte{0xa4, 0xec, 0xe5, 0x2c}: // a4 ec e5 2c
		return this.Pop(caller, input[4:])

	case [4]byte{0x5e, 0x1d, 0x05, 0x4d}: // 5e 1d 05 4d
		return this.Clear(caller, input[4:])

	case [4]byte{0x8b, 0x28, 0x29, 0x47}: // 8b 28 29 47
		return this.Set(caller, input[4:])
	}

	return []byte{}, false, 0 // unknown
}

func (this *BytesHandlers) Api() eucommon.EthApiRouter { return this.api }

func (this *BytesHandlers) New(caller evmcommon.Address, input []byte) ([]byte, bool, int64) {
	// id := this.api.UUID() // Generate a uuid for the container
	id := []byte(this.GetAddress())
	return id[:], this.connector.New(
		uint32(this.api.GetEU().(*execution.EU).Message().ID),
		types.Address(codec.Bytes20(caller).Hex()),
		hex.EncodeToString(id),
	), 0 // Create a new container
}

// Get the number of elements in the container
func (this *BytesHandlers) Length(caller evmcommon.Address, input []byte) ([]byte, bool, int64) {
	path := this.buildPath(caller) // BytesHandlers path
	if len(path) == 0 {
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
func (this *BytesHandlers) Peek(caller evmcommon.Address, input []byte) ([]byte, bool, int64) {
	path := this.buildPath(caller) // BytesHandlers path
	if len(path) == 0 {
		return []byte{}, false, 0
	}

	typedv, fees := this.api.Ccurl().PeekCommitted(path)
	if typedv != nil {
		type measurable interface{ Length() int }
		numKeys := uint64(typedv.(interfaces.Type).Value().(measurable).Length())
		if encoded, err := abi.Encode(uint256.NewInt(numKeys)); err == nil {
			return encoded, true, int64(fees)
		}
	}
	return []byte{}, false, int64(fees)
}

func (this *BytesHandlers) Get(caller evmcommon.Address, input []byte) ([]byte, bool, int64) {
	path := this.buildPath(caller) // Build container path
	if len(path) == 0 {
		return []byte{}, false, 0
	}

	idx, err := abi.DecodeTo(input, 0, uint64(0), 1, 32)
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

func (this *BytesHandlers) Set(caller evmcommon.Address, input []byte) ([]byte, bool, int64) {
	path := this.buildPath(caller) // Build container path
	if len(path) == 0 {
		return []byte{}, false, 0
	}

	idx, err := abi.DecodeTo(input, 0, uint64(0), 1, 32)
	if err != nil {
		return []byte{}, false, 0
	}

	bytes, err := abi.Decode(input, 1, []byte{}, 2, math.MaxInt)
	if bytes == nil || err != nil {
		return []byte{}, false, 0
	}

	value := noncommutative.NewBytes(bytes.([]byte))
	if _, err := this.api.Ccurl().WriteAt(uint32(this.api.GetEU().(*execution.EU).Message().ID), path, idx, value, true); err == nil {
		return []byte{}, true, 0
	}
	return []byte{}, false, 0
}

// Push a new element into the container
func (this *BytesHandlers) Push(caller evmcommon.Address, input []byte, origin evmcommon.Address, nonce uint64) ([]byte, bool, int64) {
	path := this.buildPath(caller) // BytesHandlers path
	if len(path) == 0 {
		return []byte{}, false, 0
	}

	value, err := abi.Decode(input, 0, []byte{}, 2, math.MaxInt)
	if value == nil || err != nil {
		return []byte{}, false, 0
	}

	key := path + string(this.api.ElementUID())
	_, err = this.api.Ccurl().Write(uint32(this.api.GetEU().(*execution.EU).Message().ID), key, noncommutative.NewBytes(value.([]byte)), true)
	return []byte{}, err == nil, 0
}

func (this *BytesHandlers) Pop(caller evmcommon.Address, input []byte) ([]byte, bool, int64) {
	path := this.buildPath(caller) // Build container path
	if len(path) == 0 {
		return []byte{}, false, 0
	}

	if value, _, err := this.api.Ccurl().PopBack(uint32(this.api.GetEU().(*execution.EU).Message().ID), path, true); err != nil {
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

func (this *BytesHandlers) Clear(caller evmcommon.Address, input []byte) ([]byte, bool, int64) {
	for {
		if _, ok, _ := this.Pop(caller, input); !ok {
			break
		}
	}
	return []byte{}, true, 0
}

// Build the container path
func (this *BytesHandlers) buildPath(caller evmcommon.Address) string {
	id := []byte(this.GetAddress())
	return this.connector.Key(types.Address(codec.Bytes20(caller).Hex()), hex.EncodeToString(id)) // unique ID
}

func (this *BytesHandlers) GetAddress() string {
	return string(this.api.VM().ArcologyNetworkAPIs.CallContext.Contract.CodeAddr[:]) // unique ID
}
