package abi

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/arcology-network/common-lib/codec"
	"github.com/arcology-network/common-lib/common"
	"github.com/holiman/uint256"
)

func TestDecoder(t *testing.T) {
	raw := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 100, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 160, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 32, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 96, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 99, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 75, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 32, 25, 234, 68, 190, 137, 238, 206, 15, 212, 236, 116, 130, 4, 159, 71, 42, 17, 175, 25, 56, 75, 255, 179, 138, 136, 231, 123, 59, 29, 213, 76, 25, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 32, 25, 234, 68, 190, 137, 238, 206, 15, 212, 236, 116, 130, 4, 159, 71, 42, 17, 175, 25, 56, 75, 255, 179, 138, 136, 231, 123, 59, 29, 213, 76, 25}
	decorder := NewDecoder()

	// buffer, _ := decorder.At(decorder.raw, 0, uint32(0))
	Fields := codec.Hash32s{}.Decode(raw).(codec.Hash32s)
	for i := 0; i < len(Fields); i++ {
		fmt.Println(Fields[i])
	}

	buffer, _ := decorder.Decode(raw, 1, []byte{}, 2)
	fmt.Println(buffer)

	subbytes := buffer.([]byte)
	if len(buffer.([]byte)) != 75 || common.FindFirstIf(&subbytes, func(v byte) bool { return v != 65 }) != -1 {
		t.Error("Error; The array should be 75 byte long!")
	}

	buffer, _ = decorder.Decode(raw, 0, uint32(0), 1) //need to indicate here !!
	if buffer.(uint32) != 100 {
		t.Error("Error: Should be 100!")
	}

	// fmt.Println(reflect.TypeOf(uint256.NewInt(0)).Kind())

	n := reflect.TypeOf(uint256.NewInt(0)).String()

	if n == "*uint256.Int" {
		fmt.Println(n)
	}
	fmt.Println(n)
}
