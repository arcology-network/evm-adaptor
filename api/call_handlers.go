package api

import (
	"encoding/hex"
	"fmt"
	"math"
	"reflect"

	"github.com/arcology-network/common-lib/codec"
	common "github.com/arcology-network/common-lib/common"
	"github.com/arcology-network/common-lib/types"
	"github.com/arcology-network/concurrenturl/v2/type/commutative"
	"github.com/arcology-network/concurrenturl/v2/type/noncommutative"
	ethCommon "github.com/arcology-network/evm/common"
	evmCommon "github.com/arcology-network/evm/common"
	abi "github.com/arcology-network/vm-adaptor/abi"
	"github.com/holiman/uint256"
)

// APIs under the concurrency namespace
type ConcurrentContainer struct {
	api       *API
	connector *CCurlPathBuilder
}

func NewConcurrentContainer(txHash ethCommon.Hash, txIndex uint32, api *API) *ConcurrentContainer {
	return &ConcurrentContainer{
		api:       api,
		connector: NewCCurlPathBuilder(txHash, txIndex, api.ccurl),
	}
}

func (this *ConcurrentContainer) Call(caller evmCommon.Address, input []byte, origin evmCommon.Address, nonce uint64) ([]byte, bool) {
	signature := [4]byte{}
	copy(signature[:], input)

	switch signature {
	case [4]byte{0xac, 0xaa, 0x8d, 0x70}:
		return this.New(caller, input[4:])

	case [4]byte{0xe1, 0xbb, 0x50, 0x1d}:
		return this.Push(caller, input[4:], origin, nonce)

	case [4]byte{0x05, 0x31, 0xab, 0xc6}:
		return this.Length(caller, input[4:])

	case [4]byte{0x31, 0xfe, 0x88, 0xd0}:
		return this.Get(caller, input[4:])

	case [4]byte{0xdb, 0x33, 0xc5, 0x7a}:
		print(input)
		return this.Set(caller, input[4:])

	default:
		this.api.AddLog("Unhandled function call", hex.EncodeToString(input))
	}
	return []byte{}, true
}

func (this *ConcurrentContainer) New(caller evmCommon.Address, input []byte) ([]byte, bool) {
	elemType := int(input[31]) // Data type should only take one byte.
	id := this.api.GenUUID()   // Generate a uuid for the container
	return id[:], this.connector.New(types.Address(codec.Bytes20(caller).Hex()), hex.EncodeToString(id), 0, elemType)
}

// Push a new element into the container
func (this *ConcurrentContainer) Push(caller evmCommon.Address, input []byte, origin evmCommon.Address, nonce uint64) ([]byte, bool) {
	path := this.buildPath(caller, input) // Container path

	buffer := common.DeepCopy(this.api.txHash[:])
	codec.Uint64(this.api.Serial()).EncodeToBuffer(buffer[len(buffer)-8:])
	key := path + hex.EncodeToString(buffer)

	value, _ := abi.Decode(input, 1, []byte{}, 2, math.MaxInt)
	err := this.connector.ccurl.Write(this.api.txIndex, key, noncommutative.NewBytes(value.([]byte)))
	return []byte{}, err == nil
}

// Get the number of elements in the container
func (this *ConcurrentContainer) Length(caller evmCommon.Address, input []byte) ([]byte, bool) {
	path := this.buildPath(caller, input) // Container path

	if meta, err := this.connector.ccurl.Read(this.api.txIndex, path); err == nil {
		keys := meta.(*commutative.Meta).KeyView()
		if encoded, err := abi.Encode(uint256.NewInt(uint64(len(keys)))); err == nil {
			return encoded, true
		}
	}
	return []byte{}, false
}

func (this *ConcurrentContainer) Get(caller evmCommon.Address, input []byte) ([]byte, bool) {
	path := this.buildPath(caller, input) // Build container path
	idx, _ := abi.Decode(input, 1, uint64(0), 1, 32)

	if value, err := this.connector.ccurl.ReadAt(this.api.txIndex, path, idx.(uint64)); err != nil {
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

func (this *ConcurrentContainer) Set(caller evmCommon.Address, input []byte) ([]byte, bool) {
	path := this.buildPath(caller, input) // Build container path
	idx, err := abi.Decode(input, 1, uint64(0), 1, 32)
	if err != nil {
		return []byte{}, false
	}

	bytes, err := abi.Decode(input, 2, []byte{}, 2, math.MaxInt)
	if bytes == nil || err != nil {
		return []byte{}, false
	}

	if reflect.TypeOf(bytes).String() != "[]uint8" { // Check the value data type
		return []byte{}, false
	}

	value := noncommutative.NewBytes(bytes.([]byte))
	if err := this.connector.ccurl.WriteAt(this.api.txIndex, path, idx.(uint64), value); err == nil {
		return []byte{}, true
	}
	return []byte{}, false
}

// Build the container path
func (this *ConcurrentContainer) buildPath(caller evmCommon.Address, input []byte) string {
	id, _ := abi.Decode(input, 0, []byte{}, 2, 32)                                                                            // max 32 bytes                                                                          // container ID
	return this.connector.buildContainerRootPath(types.Address(codec.Bytes20(caller).Hex()), hex.EncodeToString(id.([]byte))) // unique ID
}

func print(input []byte) {
	fmt.Println(input)
	fmt.Println()
	fmt.Println(input[:4])
	input = input[4:]
	for i := int(0); i < len(input)/32; i++ {
		fmt.Println(input[i*32 : (i+1)*32])
	}
	fmt.Println()
}
