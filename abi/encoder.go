package abi

import (
	"encoding/binary"
	"errors"

	"github.com/holiman/uint256"
)

func Encode(typed interface{}) ([]byte, error) {
	buffer := [32]byte{}

	switch typed.(type) {
	case bool:
		if typed.(bool) {
			buffer[31] = 1
		}
		return buffer[:], nil

	case uint8:
		buffer[31] = typed.(uint8)
		return buffer[:], nil

	case uint16:
		binary.BigEndian.PutUint16(buffer[32-2:], typed.(uint16))
		return buffer[:], nil

	case uint32:
		binary.BigEndian.PutUint32(buffer[32-4:], typed.(uint32))
		return buffer[:], nil

	case uint64:
		binary.BigEndian.PutUint64(buffer[32-8:], typed.(uint64))
		return buffer[:], nil

	case *uint256.Int:
		bytes := typed.(*uint256.Int).Bytes32()
		return bytes[:], nil

	case []uint8:
		binary.BigEndian.PutUint32(buffer[32-4:], uint32(len(typed.([]byte))))
		if len(typed.([]byte))%32 == 0 {
			return append(buffer[:], typed.([]byte)...), nil
		}

		body := make([]byte, len(typed.([]byte))/32*32+32)
		copy(body, typed.([]byte))
		return append(buffer[:], body...), nil
	}

	return []byte{}, errors.New("Error: Unsupported data type")
}

func AddOffset(sections [][]byte) []byte {
	encoded := []byte{}

	// sumLength := 0
	// for i := 0; i < len(sections); i++ {
	// 	sumLength
	// }

	// offset := [32]byte{}
	// offset[len(offset)-1] = uint8(len(offset))
	// encoded = append(offset[:], encoded...)
	return encoded
}
