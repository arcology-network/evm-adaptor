package types

import (
	"github.com/arcology-network/common-lib/codec"
)

func (this *ExecutionResult) HeaderSize() uint32 {
	return 8 * codec.UINT32_LEN
}

// func (this *ExecutionResult) Size() uint32 {
// 	return this.HeaderSize() +
// 		uint32(len(this.H)) +
// 		codec.UINT32_LEN +
// 		codec.Byteset(this.Transitions).Size() +
// 		codec.Bytes(this.TransitTypes).Size() +
// 		this.DC.Size() +
// 		codec.UINT64_LEN +
// 		codec.UINT64_LEN
// }

// func (this *ExecutionResult) Encode() []byte {
// 	buffer := make([]byte, this.Size())
// 	this.EncodeToBuffer(buffer)
// 	return buffer
// }

// func (this *ExecutionResult) EncodeToBuffer(buffer []byte) int {
// 	if this == nil {
// 		return 0
// 	}

// 	offset := codec.Encoder{}.FillHeader(
// 		buffer,
// 		[]uint32{
// 			codec.String(this.H).Size(),
// 			codec.Uint32(this.ID).Size(),
// 			codec.Byteset(this.Transitions).Size(),
// 			codec.Bytes(this.TransitTypes).Size(),
// 			this.DC.Size(),
// 			codec.UINT64_LEN,
// 			codec.UINT64_LEN,
// 		},
// 	)

// 	offset += codec.String(this.H).EncodeToBuffer(buffer[offset:])
// 	offset += codec.Uint32(this.ID).EncodeToBuffer(buffer[offset:])
// 	offset += codec.Byteset(this.Transitions).EncodeToBuffer(buffer[offset:])
// 	offset += codec.Bytes(this.TransitTypes).EncodeToBuffer(buffer[offset:])
// 	offset += this.DC.EncodeToBuffer(buffer[offset:])
// 	offset += codec.Uint64(this.Status).EncodeToBuffer(buffer[offset:])
// 	offset += codec.Uint64(this.GasUsed).EncodeToBuffer(buffer[offset:])

// 	return offset
// }

// func (this *ExecutionResult) Decode(buffer []byte) *ExecutionResult {
// 	fields := [][]byte(codec.Byteset{}.Decode(buffer).(codec.Byteset))

// 	this.H = string(fields[0])
// 	this.ID = uint32(codec.Uint32(0).Decode(fields[1]).(codec.Uint32))

// 	this.Transitions = [][]byte(codec.Byteset{}.Decode(fields[2]).(codec.Byteset))
// 	this.TransitTypes = []byte(codec.Bytes{}.Decode(fields[3]).(codec.Bytes))

// 	if len(fields[4]) > 0 {
// 		this.DC = (&DeferredCall{}).Decode(fields[4])
// 	}
// 	this.Status = uint64(codec.Uint64(0).Decode(fields[5]).(codec.Uint64))
// 	this.GasUsed = uint64(codec.Uint64(0).Decode(fields[6]).(codec.Uint64))
// 	return this
// }

// func (this *ExecutionResult) GobEncode() ([]byte, error) {
// 	return this.Encode(), nil
// }

// func (this *ExecutionResult) GobDecode(buffer []byte) error {
// 	this.Decode(buffer)
// 	return nil
// }
