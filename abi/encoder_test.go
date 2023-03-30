package abi

import (
	"fmt"
	"testing"

	"github.com/holiman/uint256"
)

func TestEncoder(t *testing.T) {
	var err error
	encoder := NewEncoder()
	buffer, _ := encoder.Encode(uint64(99))
	fmt.Println(buffer)
	if buffer[31] != 99 {
		t.Error("Error: Should be 99!")
	}

	buffer, _ = encoder.Encode(uint64(256))
	fmt.Println(buffer)
	if buffer[30] != 1 || buffer[31] != 0 || err != nil {
		t.Error("Error: Wrong value!")
	}

	data := [31]byte{12, 13, 14, 1, 5, 16, 17, 18}
	buffer, err = encoder.Encode(data[:])
	fmt.Println(buffer)
	if len(buffer)%32 != 0 || err != nil {
		t.Error("Error: Should be 32-byte!")
	}

	data2 := [66]byte{12, 13, 14, 1, 5, 16, 17, 18, 19, 21}
	buffer, err = encoder.Encode(data2[:])
	fmt.Println(buffer)

	if len(buffer)%32 != 0 || err != nil {
		t.Error("Error: Should be 32-byte!")
	}

	u256 := uint256.NewInt(999)
	buffer, err = encoder.Encode(u256)
	if len(buffer)%32 != 0 || err != nil {
		t.Error("Error: Should be 32-byte!")
	}
}
