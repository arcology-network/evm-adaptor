package abi

import (
	"bytes"
	"fmt"
	"math"
	"reflect"
	"testing"

	"github.com/arcology-network/common-lib/codec"
	"github.com/arcology-network/common-lib/common"
	"github.com/holiman/uint256"
)

func TestDecoder(t *testing.T) {
	raw := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 100, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 160, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 32, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 96, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 99, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 75, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 65, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 32, 25, 234, 68, 190, 137, 238, 206, 15, 212, 236, 116, 130, 4, 159, 71, 42, 17, 175, 25, 56, 75, 255, 179, 138, 136, 231, 123, 59, 29, 213, 76, 25, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 32, 25, 234, 68, 190, 137, 238, 206, 15, 212, 236, 116, 130, 4, 159, 71, 42, 17, 175, 25, 56, 75, 255, 179, 138, 136, 231, 123, 59, 29, 213, 76, 25}
	// buffer, _ := decorder.At(decorder.raw, 0, uint32(0))
	Fields := codec.Hash32s{}.Decode(raw).(codec.Hash32s)
	for i := 0; i < len(Fields); i++ {
		fmt.Println(Fields[i])
	}

	buffer, _ := Decode(raw, 1, []byte{}, 2, math.MaxInt)
	fmt.Println(buffer)

	subbytes := buffer.([]byte)
	idx, _ := common.FindFirstIf(&subbytes, func(v byte) bool { return v != 65 })
	if len(buffer.([]byte)) != 75 || idx != -1 {
		t.Error("Error; The array should be 75 byte long!")
	}

	buffer, _ = Decode(raw, 0, uint32(0), 1, math.MaxInt) //need to indicate here !!
	if buffer.(uint32) != 100 {
		t.Error("Error: Should be 100!")
	}
	fmt.Println(reflect.TypeOf(uint256.NewInt(0)).Kind())

	raw = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 64, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 32, 25, 234, 68, 190, 137, 238, 206, 15, 212, 236, 116, 130, 4, 159, 71, 42, 17, 175, 25, 56, 75, 255, 179, 138, 136, 231, 123, 59, 29, 213, 76, 25}
	buffer, _ = Decode(raw, 1, []byte{}, 1, 32)
	fmt.Println()
	fmt.Println(buffer)

	// field 0
	raw = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 96, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 160, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 32, 25, 234, 68, 190, 137, 238, 206, 15, 212, 236, 116, 130, 4, 159, 71, 42, 17, 175, 25, 56, 75, 255, 179, 138, 136, 231, 123, 59, 29, 213, 76, 25, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 5, 170, 170, 170, 170, 170, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	buffer, _ = Decode(raw, 0, []byte{}, 2, 32)
	fmt.Println()
	fmt.Println(buffer)

	if !bytes.Equal(buffer.([]byte), []byte{25, 234, 68, 190, 137, 238, 206, 15, 212, 236, 116, 130, 4, 159, 71, 42, 17, 175, 25, 56, 75, 255, 179, 138, 136, 231, 123, 59, 29, 213, 76, 25}) {
		t.Error("Error: Should be equal!")
	}

	// field 1
	buffer, _ = Decode(raw, 1, []byte{}, 1, 32)
	if len(buffer.([]byte)) != 32 || int(buffer.([]byte)[31]) != 1 {
		t.Error("Error: Wrong length")
	}

	// field 2
	buffer, _ = Decode(raw, 2, []byte{}, 2, 32)
	if len(buffer.([]byte)) != 5 {
		t.Error("Error: Wrong length")
	}

	buffer32, _ := Decode(raw, 0, [32]byte{}, 2, math.MaxInt)
	if buffer32.([32]byte)[31] != 96 {
		t.Error("Error: Wrong [32]byte length")
	}
}
