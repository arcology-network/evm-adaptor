package api

import (
	"encoding/hex"
	"math"

	"github.com/arcology-network/common-lib/codec"
	"github.com/arcology-network/common-lib/types"

	"github.com/arcology-network/concurrenturl/commutative"
	"github.com/arcology-network/concurrenturl/interfaces"
	"github.com/arcology-network/eu/cache"
	abi "github.com/arcology-network/vm-adaptor/abi"
	"github.com/arcology-network/vm-adaptor/common"
	adaptorcommon "github.com/arcology-network/vm-adaptor/common"
	intf "github.com/arcology-network/vm-adaptor/interface"
	evmcommon "github.com/ethereum/go-ethereum/common"

	"github.com/holiman/uint256"
)

// APIs under the concurrency namespace
type BaseHandlers struct {
	api       intf.EthApiRouter
	connector *adaptorcommon.BuiltinPathMaker
	args      []interface{}
}

func NewBaseHandlers(api intf.EthApiRouter, args ...interface{}) *BaseHandlers {
	return &BaseHandlers{
		api:       api,
		connector: adaptorcommon.NewBuiltinPathMaker("/container", api),
		args:      args,
	}
}

func (this *BaseHandlers) Address() [20]byte                          { return common.BYTES_HANDLER }
func (this *BaseHandlers) Connector() *adaptorcommon.BuiltinPathMaker { return this.connector }

func (this *BaseHandlers) Call(caller, callee [20]byte, input []byte, origin [20]byte, nonce uint64) ([]byte, bool, int64) {
	signature := [4]byte{}
	copy(signature[:], input)

	switch signature {
	case [4]byte{0xcd, 0xbf, 0x60, 0x8d}:
		return this.New(caller, input[4:])

	case [4]byte{0x59, 0xe0, 0x2d, 0xd7}:
		return this.peekLength(caller, input[4:])

	case [4]byte{0xf1, 0x06, 0x84, 0x54}: // f1 06 84 54
		return this.pid(caller, input[4:])

	case [4]byte{0x1f, 0x7b, 0x6d, 0x32}: // 1f 7b 6d 32
		return this.length(caller, input[4:])

	case [4]byte{0x6a, 0x3a, 0x16, 0xbd}: // =========  6a 3a 16 bd
		return this.indexByKey(caller, input[4:])

	case [4]byte{0xb7, 0xc5, 0x64, 0x6c}: // b7c5646c
		return this.keyByIndex(caller, input[4:])

	case [4]byte{0x8e, 0x7c, 0xb6, 0xe1}: // 8e 7c b6 e1
		return this.getByIndex(caller, input[4:])

	case [4]byte{0xaf, 0x4b, 0xaa, 0x7d}: // af 4b aa 7d
		return this.setByIndex(caller, input[4:])

	case [4]byte{0x7f, 0xed, 0x84, 0xf2}: //7f ed 84 f2
		return this.getByKey(caller, input[4:])

	case [4]byte{0xc2, 0x78, 0xb7, 0x99}: // c2 78 b7 99
		return this.setByKey(caller, input[4:])

	case [4]byte{0x90, 0xd2, 0x44, 0xd8}: //  90 d2 44 d8
		return this.delByIndex(caller, input[4:])

	case [4]byte{0x37, 0x79, 0xc0, 0x34}:
		return this.delByKey(caller, input[4:])

	case [4]byte{0x52, 0xef, 0xea, 0x6e}:
		return this.clear(caller, input[4:])
	}

	if len(this.args) > 0 {
		// this.args[0].(interface {
		// 	Run([20]byte, []byte, ...interface{}) ([]byte, bool, int64)
		// }).Run(caller, input[4:], this.args[1:])

		customFun := this.args[0].(func([20]byte, []byte, ...interface{}) ([]byte, bool, int64))
		customFun(caller, input[4:], this.args[1:]...)

		// return this.args.(interface {
		// 	Run([20]byte, []byte) ([]byte, bool, int64)
		// }).Run(caller, input[4:]) //more variables
	}

	return []byte{}, false, 0 // unknown
}

func (this *BaseHandlers) Api() intf.EthApiRouter { return this.api }

func (this *BaseHandlers) New(caller evmcommon.Address, input []byte) ([]byte, bool, int64) {
	connected := this.connector.New(
		this.api.GetEU().(intf.EU).ID(),            // Tx ID for conflict detection
		types.Address(codec.Bytes20(caller).Hex()), // Main contract address
	)
	return caller[:], connected, 0 // Create a new container
}

func (this *BaseHandlers) pid(caller evmcommon.Address, input []byte) ([]byte, bool, int64) {
	pidNum := this.api.Pid()
	return []byte(hex.EncodeToString(pidNum[:])), true, 0
}

// getByIndex the number of elements in the container
func (this *BaseHandlers) length(caller evmcommon.Address, input []byte) ([]byte, bool, int64) {
	path := this.connector.Key(caller)
	if length, successful, _ := this.Length(path); successful {
		if encoded, err := abi.Encode(uint256.NewInt(length)); err == nil {
			return encoded, true, 0
		}
	}
	return []byte{}, false, 0
}

// getByIndex the intial length of the container
func (this *BaseHandlers) peekLength(caller evmcommon.Address, input []byte) ([]byte, bool, int64) {
	path := this.connector.Key(caller) // BaseHandlers path
	if len(path) == 0 {
		return []byte{}, false, 0
	}

	typedv, fees := this.api.WriteCache().(*cache.WriteCache).PeekCommitted(path, new(commutative.Path))
	if typedv != nil {
		type measurable interface{ Length() int }
		numKeys := uint64(typedv.(interfaces.Type).Value().(measurable).Length())
		if encoded, err := abi.Encode(uint256.NewInt(numKeys)); err == nil {
			return encoded, true, int64(fees)
		}
	}
	return []byte{}, false, int64(fees)
}

func (this *BaseHandlers) getByIndex(caller evmcommon.Address, input []byte) ([]byte, bool, int64) {
	path := this.connector.Key(caller) // Build container path
	if len(path) == 0 {
		return []byte{}, false, 0
	}

	idx, err := abi.DecodeTo(input, 0, uint64(0), 1, 32)
	if err != nil {
		return []byte{}, false, 0
	}

	values, successful, _ := this.GetByIndex(path, idx)
	if len(values) > 0 && successful {
		return values, true, 0
	}
	return []byte{}, false, 0
}

func (this *BaseHandlers) getByKey(caller evmcommon.Address, input []byte) ([]byte, bool, int64) {
	path := this.connector.Key(caller) // Build container path
	if len(path) == 0 {
		return []byte{}, false, 0
	}

	if key, err := abi.DecodeTo(input, 0, []byte{}, 2, math.MaxInt); err == nil && len(key) > 0 {
		str := hex.EncodeToString(key)
		bytes, successful, _ := this.GetByKey(path + str)
		if len(bytes) > 0 && successful {
			return bytes, true, 0
		}
	}
	return []byte{}, false, 0
}

func (this *BaseHandlers) setByIndex(caller evmcommon.Address, input []byte) ([]byte, bool, int64) {
	path := this.connector.Key(caller) // Build container path

	idx, bytes, err := abi.Parse2(input,
		uint64(0), 1, 32,
		[]byte{}, 2, math.MaxInt,
	)

	if err != nil {
		return []byte{}, false, 0
	}

	if successful, fee := this.SetByIndex(path, idx, bytes); successful {
		return []byte{}, true, fee
	}
	return []byte{}, false, 0
}

// push a new element into the container
func (this *BaseHandlers) setByKey(caller evmcommon.Address, input []byte) ([]byte, bool, int64) {
	path := this.connector.Key(caller) // BaseHandlers path
	if len(path) == 0 {
		return []byte{}, false, 0
	}

	key, value, err := abi.Parse2(input,
		[]byte{}, 2, math.MaxInt,
		[]byte{}, 2, math.MaxInt,
	)

	if err == nil {
		str := hex.EncodeToString(key)
		successful, _ := this.SetByKey(path+str, value)
		return []byte{}, successful, 0
	}

	return []byte{}, false, 0
}

func (this *BaseHandlers) indexByKey(caller evmcommon.Address, input []byte) ([]byte, bool, int64) {
	path := this.connector.Key(caller) // BaseHandlers path
	if len(path) == 0 {
		return []byte{}, false, 0
	}

	if key, err := abi.DecodeTo(input, 0, []byte{}, 2, math.MaxInt); err == nil {
		index, _ := this.IndexOf(path, hex.EncodeToString(key))
		if encoded, err := abi.Encode(index); index != math.MaxUint64 && err == nil { // Encode the result
			return encoded, true, 0
		}
	}
	return []byte{}, false, 0
}

// 4223b5c2
func (this *BaseHandlers) keyByIndex(caller evmcommon.Address, input []byte) ([]byte, bool, int64) {
	path := this.connector.Key(caller) // BaseHandlers path
	if len(path) == 0 {
		return []byte{}, false, 0
	}

	if index, err := abi.DecodeTo(input, 0, uint64(0), 1, 32); err == nil {
		key, _ := this.KeyAt(path, index)
		v, _ := hex.DecodeString(key)
		return v, true, 0
	}
	return []byte{}, false, 0
}

func (this *BaseHandlers) delByIndex(caller evmcommon.Address, input []byte) ([]byte, bool, int64) {
	path := this.connector.Key(caller) // Build container path
	idx, err := abi.DecodeTo(input, 0, uint64(0), 1, 32)
	if err == nil {
		if successful, fee := this.SetByIndex(path, idx, nil); successful {
			return []byte{}, true, fee
		}
	}
	return []byte{}, false, 0
}

func (this *BaseHandlers) delByKey(caller evmcommon.Address, input []byte) ([]byte, bool, int64) {
	path := this.connector.Key(caller) // Build container path

	key, err := abi.DecodeTo(input, 0, []byte{}, 2, math.MaxInt)
	if err == nil {
		str := hex.EncodeToString(key)
		if successful, fee := this.SetByKey(path+str, nil); successful {
			return []byte{}, true, fee
		}
	}
	return []byte{}, false, 0
}

func (this *BaseHandlers) clear(caller evmcommon.Address, input []byte) ([]byte, bool, int64) {
	path := this.connector.Key(caller) // Build container path
	if len(path) == 0 {
		return []byte{}, false, 0
	}

	for {
		if _, _, err := this.api.WriteCache().(*cache.WriteCache).PopBack(this.api.GetEU().(intf.EU).ID(), path, nil); err != nil {
			break
		}
	}
	return []byte{}, true, 0
}
