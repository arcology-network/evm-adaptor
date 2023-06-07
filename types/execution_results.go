package types

import "github.com/arcology-network/common-lib/codec"

type ExecutionResults []*ExecutionResult

func (this ExecutionResults) ToNewMessage() uint32 {
	return uint32((len(this) + 1) * codec.UINT32_LEN) // Header length
}

// func (this ExecutionResults) HeaderSize() uint32 {
// 	return uint32((len(this) + 1) * codec.UINT32_LEN) // Header length
// }

// func (this ExecutionResults) Size() uint32 {
// 	total := this.HeaderSize()
// 	for i := 0; i < len(this); i++ {
// 		total += (this)[i].Size()
// 	}
// 	return total
// }

// // Fill in the header info
// func (this ExecutionResults) FillHeader(buffer []byte) {
// 	codec.Uint32(len(this)).EncodeToBuffer(buffer)

// 	offset := uint32(0)
// 	for i := 0; i < len(this); i++ {
// 		codec.Uint32(offset).EncodeToBuffer(buffer[codec.UINT32_LEN*(i+1):])
// 		offset += (this)[i].Size()
// 	}
// }

// func (this ExecutionResults) GobEncode() ([]byte, error) {
// 	buffer := make([]byte, this.Size())
// 	this.FillHeader(buffer)

// 	offsets := make([]uint32, len(this)+1)
// 	offsets[0] = 0
// 	for i := 0; i < len(this); i++ {
// 		offsets[i+1] = offsets[i] + this[i].Size()
// 	}

// 	headerLen := this.HeaderSize()
// 	worker := func(start, end, index int, args ...interface{}) {
// 		for i := start; i < end; i++ {
// 			this[i].EncodeToBuffer(buffer[headerLen+offsets[i]:])
// 		}
// 	}
// 	common.ParallelWorker(len(this), 4, worker)
// 	return buffer, nil
// }

// func (this ExecutionResults) GobDecode(buffer []byte) error {
// 	bytesset := [][]byte(codec.Byteset{}.Decode(buffer).(codec.Byteset))
// 	euresults := make([]*ExecutionResult, len(bytesset))
// 	worker := func(start, end, index int, args ...interface{}) {
// 		for i := start; i < end; i++ {
// 			euresults[i] = (&ExecutionResult{}).Decode(bytesset[i])
// 		}
// 	}
// 	common.ParallelWorker(len(bytesset), 4, worker)
// 	this = euresults
// 	return nil
// }
