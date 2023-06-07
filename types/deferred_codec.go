package types

// func (this *DeferredCall) Encode() []byte {
// 	buffer := make([]byte, this.Size())
// 	this.EncodeToBuffer(buffer)
// 	return buffer
// }

// func (this *DeferredCall) HeaderSize() uint32 {
// 	return 4 * codec.UINT32_LEN
// }

// func (this *DeferredCall) Size() uint32 {
// 	if this == nil {
// 		return 0
// 	}

// 	return this.HeaderSize() +
// 		uint32(len(this.TxHash)+len(this.groupBy)+len(this.callData))
// }

// func (this *DeferredCall) EncodeToBuffer(buffer []byte) int {
// 	if this == nil {
// 		return 0
// 	}

// 	offset := codec.Encoder{}.FillHeader(
// 		buffer,
// 		[]uint32{
// 			codec.Bytes32(this.TxHash).Size(),
// 			codec.Bytes32(this.groupBy).Size(),
// 			codec.Bytes(this.callData).Size(),
// 		},
// 	)

// 	offset += codec.Bytes32(this.TxHash).EncodeToBuffer(buffer[offset:])
// 	offset += codec.Bytes32(this.groupBy).EncodeToBuffer(buffer[offset:])
// 	offset += codec.Bytes(this.callData).EncodeToBuffer(buffer[offset:])
// 	return offset
// }

// func (this *DeferredCall) Decode(data []byte) *DeferredCall {
// 	buffers := [][]byte(codec.Byteset{}.Decode(data).(codec.Byteset))
// 	this.TxHash = codec.Bytes32{}.Decode(buffers[0]).(codec.Bytes32)
// 	this.groupBy = codec.Bytes32{}.Decode(buffers[1]).(codec.Bytes32)
// 	this.callData = codec.Bytes{}.Decode(buffers[2]).(codec.Bytes)
// 	return this
// }
