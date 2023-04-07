package abi

import (
	"encoding/binary"
	"errors"

	"github.com/arcology-network/common-lib/common"
	"github.com/holiman/uint256"
)

func Decode(raw []byte, idx int, typed interface{}, depth uint8, maxLength int) (interface{}, error) {
	if depth < 1 {
		return nil, errors.New("Error: Can be 0 deep!!")
	}

	if idx >= len(raw)/32 {
		return raw, nil
	}

	if idx*32+32 > len(raw) {
		return nil, errors.New("Error: Access out of range")
	}

	switch typed.(type) {
	case bool:
		return raw[len(raw[idx*32:idx*32+32])-1] == 1, nil
	case uint8:
		return uint8(raw[idx*32+32-1]), nil
	case uint16:
		return binary.BigEndian.Uint16(raw[idx*32+32-2 : idx*32+32]), nil
	case uint32:
		return binary.BigEndian.Uint32(raw[idx*32+32-4 : idx*32+32]), nil
	case uint64:
		return binary.BigEndian.Uint64(raw[idx*32+32-8 : idx*32+32]), nil
	case []uint8:
		if depth == 1 {
			length := common.Min(len(raw), maxLength)
			return raw[idx*32 : idx*32+length], nil
		}
		depth--

		sub := raw[idx*32+32-4 : idx*32+32]
		offset := binary.BigEndian.Uint32(sub)
		return next(raw, offset, depth, maxLength)

	case *uint256.Int:
		var v uint256.Int
		v.SetBytes(raw[idx*32 : idx*32+32])
		return v, nil
	}
	return raw, errors.New("Error: Unknown type")
}

func next(raw []byte, offset uint32, depth uint8, maxLength int) (interface{}, error) {
	length, _ := Decode(raw[offset:], 0, uint32(0), depth, maxLength)

	if offset+length.(uint32)+32 > uint32(len(raw)) {
		return nil, errors.New("Error: Access out of range")
	}

	sub := raw[offset+32 : offset+length.(uint32)+32]
	return Decode(sub, 0, []byte{}, depth, maxLength)
}
