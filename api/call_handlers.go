package api

import (
	"encoding/hex"
	"fmt"

	"github.com/arcology-network/common-lib/codec"
	"github.com/arcology-network/common-lib/types"
	"github.com/arcology-network/concurrenturl/v2/type/commutative"
	"github.com/arcology-network/evm/common"
	ethCommon "github.com/arcology-network/evm/common"
	abi "github.com/arcology-network/vm-adaptor/abi"
)

// APIs under the concurrency namespace
type ConcurrentContainer struct {
	api       *API
	connector *CCurlConnector
}

func NewConcurrentContainer(txHash ethCommon.Hash, txIndex uint32, api *API) *ConcurrentContainer {
	return &ConcurrentContainer{
		api:       api,
		connector: NewApiCCurlConnector(txHash, txIndex, api.ccurl),
	}
}

func (this *ConcurrentContainer) Call(caller common.Address, input []byte, origin common.Address, nonce uint64) ([]byte, bool) {
	signature := [4]byte{}
	copy(signature[:], input)

	fmt.Println(input)
	fmt.Println()
	Printer(input)
	fmt.Println()
	fmt.Println()

	switch signature {
	case [4]byte{0xac, 0xaa, 0x8d, 0x70}:
		return this.New(caller, input[4:], origin, nonce)

	// case [4]byte{0x58, 0x94, 0x13, 0x33}:
	// 	return this.Length(input[4:], origin, nonce)

	case [4]byte{0x9e, 0xc6, 0x69, 0x25}:
		return this.Push(caller, input[4:], origin, nonce)

	case [4]byte{0x05, 0x31, 0xab, 0xc6}:
		return this.Length(caller, input[4:])
	default:
		return this.Push(caller, input[4:], origin, nonce)
	}

}

func CopyLeftAligned(buffer []byte, encoded []byte) bool {
	if len(buffer) < len(encoded) {
		return false
	}
	copy(buffer[len(buffer)-len(encoded):], encoded)
	return true
}

func (this *ConcurrentContainer) New(caller common.Address, input []byte, origin common.Address, nonce uint64) ([]byte, bool) {
	elemType := int(input[31]) // Data type should only take one byte.

	buffer := [32]byte{}
	id := this.api.GenUUID() // Generate a uuid for the container
	if !CopyLeftAligned(buffer[:], id) {
		return []byte{}, true
	}
	return buffer[:], this.connector.New(types.Address(codec.Bytes20(caller).Hex()), hex.EncodeToString(id), 0, elemType)
}

// Push a new element into the container
func (this *ConcurrentContainer) Push(caller common.Address, input []byte, origin common.Address, nonce uint64) ([]byte, bool) {
	decoder := abi.NewDecoder()
	id, _ := decoder.Decode(input, 0, []byte{}, 2) // container ID
	fmt.Print(hex.EncodeToString(id.([]byte)))

	path := this.connector.buildContainerRootPath(types.Address(codec.Bytes20(caller).Hex()), hex.EncodeToString(id.([]byte))) // caller + container ID
	err := this.connector.ccurl.Write(this.api.txIndex, path, input[32:])
	return []byte{}, err == nil
}

// Get the number of elements in the container
func (this *ConcurrentContainer) Length(caller common.Address, input []byte) ([]byte, bool) {
	decoder := abi.NewDecoder()
	id, _ := decoder.Decode(input, 0, []byte{}, 2) // container ID
	fmt.Print(hex.EncodeToString(id.([]byte)))

	path := this.connector.buildContainerRootPath(types.Address(codec.Bytes20(caller).Hex()), hex.EncodeToString(id.([]byte))) // unique ID

	if meta, err := this.connector.ccurl.Read(this.api.txIndex, path); err == nil {
		keys := meta.(*commutative.Meta).PeekKeys()
		fmt.Println("Keys: ", keys)

		return []byte{}, true
	}

	return []byte{}, true
}

func Printer(input []byte) {
	fmt.Println(input[:4])
	input = input[4:]
	for i := int(0); i < len(input)/32; i++ {
		fmt.Println(input[i*32 : (i+1)*32])
	}
}
