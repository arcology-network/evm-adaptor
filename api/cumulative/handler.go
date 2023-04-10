package cumulative

import (
	"encoding/hex"
	"fmt"
	"math"
	"reflect"

	"github.com/arcology-network/common-lib/codec"
	"github.com/arcology-network/common-lib/types"

	"github.com/arcology-network/concurrenturl/v2/type/noncommutative"
	evmcommon "github.com/arcology-network/evm/common"
	abi "github.com/arcology-network/vm-adaptor/abi"

	apicommon "github.com/arcology-network/vm-adaptor/api/common"
)

// APIs under the concurrency namespace
type Cumulative struct {
	api       apicommon.ContextInfoInterface
	connector *apicommon.CCurlPathBuilder
}

func NewCumulative(api apicommon.ContextInfoInterface) *Cumulative {
	return &Cumulative{
		api:       api,
		connector: apicommon.NewCCurlPathBuilder(api.TxHash(), api.TxIndex(), api.Ccurl()),
	}
}

func (this *Cumulative) Address() [20]byte {
	return [20]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0x85}
}

func (this *Cumulative) Call(caller evmcommon.Address, input []byte, origin evmcommon.Address, nonce uint64) ([]byte, bool) {
	signature := [4]byte{}
	copy(signature[:], input)

	switch signature {
	case [4]byte{0x26, 0x4a, 0x7f, 0x20}:
		return this.New(caller, input[4:])

	case [4]byte{0x1e, 0x81, 0x3e, 0x5d}:
		return this.Get(caller, input[4:])

	case [4]byte{0x76, 0x18, 0xc3, 0x17}:
		return this.Add(caller, input[4:])

	case [4]byte{0x57, 0xf1, 0xa8, 0x7e}:
		return this.Sub(caller, input[4:])
	}
	return this.Unknow(caller, input[4:])
}

func (this *Cumulative) Unknow(caller evmcommon.Address, input []byte) ([]byte, bool) {
	this.api.AddLog("Unhandled function call in cumulative handler router", hex.EncodeToString(input))
	return []byte{}, false
}

func (this *Cumulative) New(caller evmcommon.Address, input []byte) ([]byte, bool) {
	// elemType := int(input[31]) // Data type should only take one byte.
	id := this.api.GenUUID() // Generate a uuid for the container
	return id[:], this.connector.New(types.Address(codec.Bytes20(caller).Hex()), hex.EncodeToString(id), 0)
}

func (this *Cumulative) Get(caller evmcommon.Address, input []byte) ([]byte, bool) {
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

func (this *Cumulative) Add(caller evmcommon.Address, input []byte) ([]byte, bool) {
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
	if err := this.api.Ccurl().WriteAt(this.api.TxIndex(), path, idx.(uint64), value); err == nil {
		return []byte{}, true
	}
	return []byte{}, false
}

func (this *Cumulative) Sub(caller evmcommon.Address, input []byte) ([]byte, bool) {
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
	if err := this.api.Ccurl().WriteAt(this.api.TxIndex(), path, idx.(uint64), value); err == nil {
		return []byte{}, true
	}
	return []byte{}, false
}

// Build the container path
func (this *Cumulative) buildPath(caller evmcommon.Address, input []byte) string {
	id, _ := abi.Decode(input, 0, []byte{}, 2, 32)                                                                            // max 32 bytes                                                                          // container ID
	return this.connector.BuildContainerRootPath(types.Address(codec.Bytes20(caller).Hex()), hex.EncodeToString(id.([]byte))) // unique ID
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
