package abi

import (
	"encoding/binary"
	"errors"
	"reflect"

	"github.com/arcology-network/common-lib/common"
	"github.com/holiman/uint256"
)

func CanDecodeTo[T any](raw []byte, idx int, initv T, depth uint8, maxLength int, validate func(v T) bool) bool {
	if v, err := DecodeTo(raw, idx, initv, depth, maxLength); err == nil {
		if validate != nil {
			return validate(v)
		}
		return true
	}
	return false
}

func DecodeTo[T any](raw []byte, idx int, initv T, depth uint8, maxLength int) (T, error) {
	v, err := Decode(raw, idx, initv, depth, maxLength)
	if err == nil {
		if reflect.TypeOf(v) == reflect.TypeOf(initv) {
			return v.(T), nil
		}
	}
	return initv, err
}

func Decode(raw []byte, idx int, initv interface{}, depth uint8, maxLength int) (interface{}, error) {
	if depth < 1 {
		return nil, errors.New("Error: Can be 0 deep!!")
	}

	if idx >= len(raw)/32 {
		return raw, nil
	}

	if idx*32+32 > len(raw) {
		return nil, errors.New("Error: Access out of range")
	}

	switch initv.(type) {
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

	case *uint256.Int:
		var v uint256.Int
		v.SetBytes(raw[idx*32 : idx*32+32])
		return &v, nil

	case uint256.Int:
		var v uint256.Int
		v.SetBytes(raw[idx*32 : idx*32+32])
		return v, nil

	case [20]uint8:
		var v [20]byte
		copy(v[:], raw[idx*32+12:idx*32+32])
		return v, nil

	case [32]uint8:
		var v [32]byte
		copy(v[:], raw[idx*32:idx*32+32])
		return v, nil

	case []uint8:
		if depth == 1 {
			length := common.Min(len(raw), maxLength)
			return raw[idx*32 : idx*32+length], nil
		}
		depth--

		sub := raw[idx*32+32-4 : idx*32+32]
		offset := binary.BigEndian.Uint32(sub)
		return next(raw, offset, depth, maxLength)
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
