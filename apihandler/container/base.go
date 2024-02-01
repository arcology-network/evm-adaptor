package api

import (
	"encoding/hex"
	"math"
	"sort"

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
	api         intf.EthApiRouter
	pathBuilder *adaptorcommon.PathBuilder
	caller      [20]byte
	args        []interface{}
}

func NewBaseHandlers(api intf.EthApiRouter, args ...interface{}) *BaseHandlers {
	return &BaseHandlers{
		api:         api,
		pathBuilder: adaptorcommon.NewPathBuilder("/container", api),
		args:        args,
	}
}

func (this *BaseHandlers) Address() [20]byte                     { return common.BYTES_HANDLER }
func (this *BaseHandlers) Connector() *adaptorcommon.PathBuilder { return this.pathBuilder }

func (this *BaseHandlers) Call(caller, callee [20]byte, input []byte, origin [20]byte, nonce uint64) ([]byte, bool, int64) {
	signature := [4]byte{}
	copy(signature[:], input)

	switch signature {
	case [4]byte{0xcd, 0xbf, 0x60, 0x8d}:
		return this.new(caller, input[4:]) // Create a new container

	case [4]byte{0x59, 0xe0, 0x2d, 0xd7}:
		return this.peekLength(caller, input[4:]) // Get the initial length of the container, it remains the same in the same block.

	case [4]byte{0xf1, 0x06, 0x84, 0x54}:
		return this.pid(caller, input[4:]) // Get the pesudo process ID.

	case [4]byte{0x1f, 0x7b, 0x6d, 0x32}:
		return this.length(caller, input[4:]) // Get the current number of elements in the container.

	case [4]byte{0x6a, 0x3a, 0x16, 0xbd}:
		return this.indexByKey(caller, input[4:]) // Get the index of the element by its key.

	case [4]byte{0xb7, 0xc5, 0x64, 0x6c}:
		return this.keyByIndex(caller, input[4:]) // Get the key of the element by its index.

	case [4]byte{0x8e, 0x7c, 0xb6, 0xe1}:
		return this.getByIndex(caller, input[4:]) // Get the element by its index.

	case [4]byte{0xaf, 0x4b, 0xaa, 0x7d}:
		return this.setByIndex(caller, input[4:]) // Set the element by its index.

	case [4]byte{0x7f, 0xed, 0x84, 0xf2}:
		return this.getByKey(caller, input[4:]) // Get the element by its key.

	case [4]byte{0xc2, 0x78, 0xb7, 0x99}:
		return this.setByKey(caller, input[4:]) // Set the element by its key.

	case [4]byte{0x90, 0xd2, 0x44, 0xd8}:
		return this.delByIndex(caller, input[4:]) // Delete the element by its index.

	case [4]byte{0x37, 0x79, 0xc0, 0x34}:
		return this.delByKey(caller, input[4:]) // Delete the element by its key.

	case [4]byte{0x52, 0xef, 0xea, 0x6e}:
		return this.clear(caller, input[4:]) // Clear the container.
	}

	// Custom function call. The base handler may have a custom function to call..
	if len(this.args) > 0 {
		customFun := this.args[0].(func([20]byte, []byte, ...interface{}) ([]byte, bool, int64))
		customFun(caller, input[4:], this.args[1:]...)
	}

	return []byte{}, false, 0 // unknown
}

func (this *BaseHandlers) Api() intf.EthApiRouter { return this.api }

// Create a new container. This function is called when the constructor of the base contract is called in the concurrentlib.
func (this *BaseHandlers) new(caller evmcommon.Address, input []byte) ([]byte, bool, int64) {
	connected := this.pathBuilder.New(
		this.api.GetEU().(interface{ ID() uint32 }).ID(), // Tx ID for conflict detection
		types.Address(codec.Bytes20(caller).Hex()),       // Main contract address
	)
	this.caller = caller
	return caller[:], connected, 0 // Create a new container
}

func (this *BaseHandlers) pid(caller evmcommon.Address, input []byte) ([]byte, bool, int64) {
	pidNum := this.api.Pid()
	return []byte(hex.EncodeToString(pidNum[:])), true, 0
}

// getByIndex the number of elements in the container
func (this *BaseHandlers) length(caller evmcommon.Address, input []byte) ([]byte, bool, int64) {
	path := this.pathBuilder.Key(caller)
	if length, successful, _ := this.Length(path); successful {
		if encoded, err := abi.Encode(uint256.NewInt(length)); err == nil {
			return encoded, true, 0
		}
	}
	return []byte{}, false, 0
}

// peekLength the initial length of the container, which would remain the same in the same block.
func (this *BaseHandlers) peekLength(caller evmcommon.Address, input []byte) ([]byte, bool, int64) {
	path := this.pathBuilder.Key(caller) // BaseHandlers path
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

// getByIndex the element by its index
func (this *BaseHandlers) getByIndex(caller evmcommon.Address, input []byte) ([]byte, bool, int64) {
	path := this.pathBuilder.Key(caller) // Build container path
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
	path := this.pathBuilder.Key(caller) // Build container path
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

// Set the element by its index. This function is different from the SetByKey function in that if
// the index is out of range, the function will return false.
func (this *BaseHandlers) setByIndex(caller evmcommon.Address, input []byte) ([]byte, bool, int64) {
	path := this.pathBuilder.Key(caller) // Build container path

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

// Push a new element into the container. If the key does not exist, it will be created and the value will be set.
func (this *BaseHandlers) setByKey(caller evmcommon.Address, input []byte) ([]byte, bool, int64) {
	path := this.pathBuilder.Key(caller) // BaseHandlers path
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
	path := this.pathBuilder.Key(caller) // BaseHandlers path
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
	path := this.pathBuilder.Key(caller) // BaseHandlers path
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
	path := this.pathBuilder.Key(caller) // Build container path
	idx, err := abi.DecodeTo(input, 0, uint64(0), 1, 32)
	if err == nil {
		if successful, fee := this.SetByIndex(path, idx, nil); successful {
			return []byte{}, true, fee
		}
	}
	return []byte{}, false, 0
}

func (this *BaseHandlers) delByKey(caller evmcommon.Address, input []byte) ([]byte, bool, int64) {
	path := this.pathBuilder.Key(caller) // Build container path

	key, err := abi.DecodeTo(input, 0, []byte{}, 2, math.MaxInt)
	if err == nil {
		str := hex.EncodeToString(key)
		if successful, fee := this.SetByKey(path+str, nil); successful {
			return []byte{}, true, fee
		}
	}
	return []byte{}, false, 0
}

func (this *BaseHandlers) max(caller evmcommon.Address, input []byte) ([]byte, bool, int64) {
	entries, _, _ := this.Export(this.pathBuilder.Key(caller))
	sort.Slice(entries, func(i, j int) bool {
		return string(entries[i]) < string(entries[j])
	})
	return []byte{}, false, 0
}

func (this *BaseHandlers) min(caller evmcommon.Address, input []byte) ([]byte, bool, int64) {
	entries, _, _ := this.Export(this.pathBuilder.Key(caller))
	sort.Slice(entries, func(i, j int) bool {
		return string(entries[i]) > string(entries[j])
	})
	return []byte{}, false, 0
}

func (this *BaseHandlers) clear(caller evmcommon.Address, input []byte) ([]byte, bool, int64) {
	path := this.pathBuilder.Key(caller) // Build container path
	if len(path) == 0 {
		return []byte{}, false, 0
	}

	for {
		if _, _, err := this.api.WriteCache().(*cache.WriteCache).PopBack(this.api.GetEU().(interface{ ID() uint32 }).ID(), path, nil); err != nil {
			break
		}
	}
	return []byte{}, true, 0
}
