package abi

import (
	"encoding/binary"
	"errors"
	"reflect"

	"github.com/arcology-network/common-lib/common"
	"github.com/holiman/uint256"
)

func Parse2[T0, T1 any](input []byte,
	_v0 T0, _depth0 uint8, _len0 int,
	_v1 T1, _depth1 uint8, _len1 int) (T0, T1, error) {

	decodedv0, err := DecodeTo(input, 0, _v0, _depth0, _len0)
	if err != nil {
		return _v0, _v1, errors.New("Error: Failed to decode the first")
	}

	decodedv1, err := DecodeTo(input, 1, _v1, _depth1, _len1)
	if err != nil {
		return _v0, _v1, errors.New("Error: Failed to decode the second")
	}
	return decodedv0, decodedv1, nil
}

func Parse3[T0, T1, T2 any](input []byte,
	_v0 T0, _depth0 uint8, _len0 int,
	_v1 T1, _depth1 uint8, _len1 int,
	_v2 T2, _depth2 uint8, _len2 int) (T0, T1, T2, error) {

	decodedv0, err := DecodeTo(input, 0, _v0, _depth0, _len0)
	if err != nil {
		return _v0, _v1, _v2, errors.New("Error: Failed to decode v0")
	}

	decodedv1, err := DecodeTo(input, 1, _v1, _depth1, _len1)
	if err != nil {
		return _v0, _v1, _v2, errors.New("Error: Failed to decode v1")
	}

	decodedv2, err := DecodeTo(input, 2, _v2, _depth2, _len2)
	if err != nil {
		return _v0, _v1, _v2, errors.New("Error: Failed to parse v2")
	}
	return decodedv0, decodedv1, decodedv2, nil
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
			return raw[idx*32 : common.Min(idx*32+length, len(raw))], nil
		}
		depth--

		sub := raw[idx*32+32-4 : idx*32+32]
		offset := binary.BigEndian.Uint32(sub)
		return next(raw, offset, depth, maxLength)
	}

	return raw, errors.New("Error: Unknown type")
}

func next(raw []byte, offset uint32, depth uint8, maxLength int) (interface{}, error) {
	if len(raw) <= int(offset) {
		return nil, errors.New("Error: Access out of range")
	}

	length, _ := Decode(raw[offset:], 0, uint32(0), depth, maxLength)

	if offset+length.(uint32)+32 > uint32(len(raw)) {
		return nil, errors.New("Error: Access out of range")
	}

	sub := raw[offset+32 : offset+length.(uint32)+32]
	return Decode(sub, 0, []byte{}, depth, maxLength)
}
