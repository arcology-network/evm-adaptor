package concurrentcontainer

import (
	"math"

	"github.com/arcology-network/common-lib/codec"
	"github.com/arcology-network/common-lib/types"

	"github.com/arcology-network/concurrenturl/interfaces"
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
	handler   interface{}
}

func NewHandler(api eucommon.EthApiRouter, handler interface{}) *BytesHandlers {
	return &BytesHandlers{
		api:       api,
		connector: apicommon.NewCCurlConnector("/container", api, api.Ccurl()),
		handler:   handler,
	}
}

func (this *BytesHandlers) Address() [20]byte {
	return common.BYTES_HANDLER
}

func (this *BytesHandlers) Connector() *apicommon.CcurlConnector { return this.connector }

func (this *BytesHandlers) Call(caller, callee [20]byte, input []byte, origin [20]byte, nonce uint64) ([]byte, bool, int64) {
	signature := [4]byte{}
	copy(signature[:], input)

	switch signature {
	case [4]byte{0xcd, 0xbf, 0x60, 0x8d}:
		return this.New(caller, input[4:])

	case [4]byte{0x59, 0xe0, 0x2d, 0xd7}:
		return this.PeekLength(caller, input[4:])

	case [4]byte{0x3b, 0x3d, 0xca, 0x76}: // 3b 3d ca 76
		return this.rand(caller, input[4:])

	case [4]byte{0x7d, 0xac, 0xda, 0x03}: // 7d ac da 03
		return this.push(caller, input[4:], origin, nonce)

	case [4]byte{0x1f, 0x7b, 0x6d, 0x32}: // 1f 7b 6d 32
		return this.length(caller, input[4:])

	case [4]byte{0x95, 0x07, 0xd3, 0x9a}: // 95 07 d3 9a
		return this.get(caller, input[4:])

	case [4]byte{0xe9, 0x1e, 0xa7, 0x22}:
		return this.find(caller, input[4:])

	case [4]byte{0xa4, 0xec, 0xe5, 0x2c}: // a4 ec e5 2c
		return this.pop(caller, input[4:])

	case [4]byte{0x20, 0xba, 0x5b, 0x60}: // 20 ba 5b 60
		return this.insert(caller, input[4:])

	case [4]byte{0x52, 0xef, 0xea, 0x6e}: // 52 ef ea 6e
		return this.clear(caller, input[4:])

	case [4]byte{0x8b, 0x28, 0x29, 0x47}: // 8b 28 29 47
		return this.set(caller, input[4:])
	}

	if this.handler != nil {
		return this.handler.(interface {
			Run([20]byte, []byte) ([]byte, bool, int64)
		}).Run(caller, input[4:])
	}

	return []byte{}, false, 0 // unknown
}

func (this *BytesHandlers) Api() eucommon.EthApiRouter { return this.api }

func (this *BytesHandlers) New(caller evmcommon.Address, input []byte) ([]byte, bool, int64) {
	connected := this.connector.New(
		uint32(this.api.GetEU().(*execution.EU).Message().ID), // Tx ID for conflict detection
		types.Address(codec.Bytes20(caller).Hex()),            // Main contract address
	)
	return caller[:], connected, 0 // Create a new container
}

// get the number of elements in the container
func (this *BytesHandlers) length(caller evmcommon.Address, input []byte) ([]byte, bool, int64) {
	path := this.connector.Key(caller)
	if length, successful, _ := this.Length(path); successful {
		if encoded, err := abi.Encode(uint256.NewInt(length)); err == nil {
			return encoded, true, 0
		}
	}
	return []byte{}, false, 0
}

// get the intial length of the container
func (this *BytesHandlers) PeekLength(caller evmcommon.Address, input []byte) ([]byte, bool, int64) {
	path := this.connector.Key(caller) // BytesHandlers path
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

func (this *BytesHandlers) get(caller evmcommon.Address, input []byte) ([]byte, bool, int64) {
	path := this.connector.Key(caller) // Build container path
	if len(path) == 0 {
		return []byte{}, false, 0
	}

	idx, err := abi.DecodeTo(input, 0, uint64(0), 1, 32)
	if err != nil {
		return []byte{}, false, 0
	}

	values, successful, _ := this.Get(path, idx)
	if len(values) > 0 && successful {
		if encoded, err := abi.Encode(values); err == nil { // Encode the result
			offset := [32]byte{}
			offset[len(offset)-1] = uint8(len(offset))
			encoded = append(offset[:], encoded...)
			return encoded, true, 0
		}
	}
	return []byte{}, false, 0
}

// push a new element into the container
func (this *BytesHandlers) find(caller evmcommon.Address, input []byte) ([]byte, bool, int64) {
	path := this.connector.Key(caller) // BytesHandlers path
	if len(path) == 0 {
		return []byte{}, false, 0
	}

	key, err := abi.DecodeTo(input, 0, []byte{}, 2, 32)
	if err != nil {
		return []byte{}, false, 0
	}

	if err == nil {
		return this.Insert(path, (key))
	}

	return []byte{}, false, 0
}

func (this *BytesHandlers) set(caller evmcommon.Address, input []byte) ([]byte, bool, int64) {
	path := this.connector.Key(caller) // Build container path
	idx, bytes, err := abi.Parse2(input, uint64(0), 1, 32, []byte{}, 2, math.MaxInt)

	if err != nil {
		return []byte{}, false, 0
	}

	if successful, fee := this.Set(path, idx, bytes); successful {
		return []byte{}, true, fee
	}
	return []byte{}, false, 0
}

func (this *BytesHandlers) rand(caller evmcommon.Address, input []byte) ([]byte, bool, int64) {
	randNum := this.api.ElementUID()
	return randNum, true, 0
}

// push a new element into the container
func (this *BytesHandlers) push(caller evmcommon.Address, input []byte, origin evmcommon.Address, nonce uint64) ([]byte, bool, int64) {
	path := this.connector.Key(caller) // BytesHandlers path
	return this.Push(path, input)
}

// push a new element into the container
func (this *BytesHandlers) insert(caller evmcommon.Address, input []byte) ([]byte, bool, int64) {
	path := this.connector.Key(caller) // BytesHandlers path
	if len(path) == 0 {
		return []byte{}, false, 0
	}

	key, value, err := abi.Parse2(input,
		[]byte{}, 2, math.MaxInt,
		[]byte{}, 2, math.MaxInt,
	)

	if err == nil {
		return this.Insert(path+string(key), value)
	}

	return []byte{}, false, 0
}

func (this *BytesHandlers) pop(caller evmcommon.Address, _ []byte) ([]byte, bool, int64) {
	path := this.connector.Key(caller) // Build container path
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

func (this *BytesHandlers) clear(caller evmcommon.Address, input []byte) ([]byte, bool, int64) {
	for {
		if _, ok, _ := this.pop(caller, nil); !ok {
			break
		}
	}
	return []byte{}, true, 0
}
