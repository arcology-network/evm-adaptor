package abi

import (
	"encoding/binary"
	"errors"
	"reflect"

	"github.com/holiman/uint256"
)

type Decoder struct{}

func NewDecoder() *Decoder {
	return &Decoder{}
}

func (this *Decoder) Decode(raw []byte, idx int, typed interface{}, depth uint8) (interface{}, error) {
	if idx >= len(raw)/32 {
		return nil, errors.New("Error: Out of range access!!")
	}

	switch reflect.TypeOf(typed).String() {
	case "bool":
		return raw[len(raw[idx*32:idx*32+32])-1] == 1, nil
	case "uint8":
		return uint8(raw[idx*32+32-1]), nil
	case "uint16":
		return binary.BigEndian.Uint16(raw[idx*32+32-2 : idx*32+32]), nil
	case "uint32":
		return binary.BigEndian.Uint32(raw[idx*32+32-4 : idx*32+32]), nil
	case "uint64":
		return binary.BigEndian.Uint64(raw[idx*32+32-8 : idx*32+32]), nil
	case "[]uint8":
		if depth == 1 {
			return raw, nil
		}
		depth--

		offset := binary.BigEndian.Uint32(raw[idx*32+32-4 : idx*32+32])
		return this.next(raw, offset, depth)

	case "*uint256.Int":
		var v uint256.Int
		v.SetBytes(raw[idx*32 : idx*32+32])
		return v, nil
	}
	return raw, nil
}

func (this *Decoder) next(raw []byte, offset uint32, depth uint8) (interface{}, error) {
	length, _ := this.Decode(raw[offset:], 0, uint32(0), depth)
	sub := raw[offset+32 : offset+length.(uint32)+32]

	return this.Decode(sub, 0, []byte{}, depth)
}
