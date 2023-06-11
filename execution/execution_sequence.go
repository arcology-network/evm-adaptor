package execution

import (
	"crypto/sha256"
)

type Sequence struct {
	ID              [32]byte
	Predecessors    [][32]byte
	PredecessorHash [32]byte
	Msgs            []*StandardMessage
	Parallel        bool
}

func NewSequence(ID [32]byte, predecessors [][32]byte, msgs []*StandardMessage, parallel bool) *Sequence {
	buffer := make([]byte, len(msgs)*32)
	for i, v := range msgs {
		copy(buffer[i*32:], v.TxHash[:])
	}

	return &Sequence{
		ID:              sha256.Sum256(buffer),
		Predecessors:    predecessors,
		PredecessorHash: sha256.Sum256(buffer),
		Msgs:            msgs,
		Parallel:        parallel,
	}
}

// func (this ExecutingSequences) Encode() ([]byte, error) {
// 	if this == nil {
// 		return []byte{}, nil
// 	}

// 	data := make([][]byte, len(this))
// 	worker := func(start, end, idx int, args ...interface{}) {
// 		executingSequences := args[0].([]interface{})[0].(ExecutingSequences)
// 		data := args[0].([]interface{})[1].([][]byte)
// 		for i := start; i < end; i++ {
// 			standardMessages := StandardMessages(executingSequences[i].Msgs)
// 			standardMessagesData, err := standardMessages.Encode()
// 			if err != nil {
// 				standardMessagesData = []byte{}
// 			}

// 			tmpData := [][]byte{
// 				standardMessagesData,
// 				codec.Bools([]bool{executingSequences[i].Parallel}).Encode(),
// 				executingSequences[i].SequenceId[:],
// 				codec.Uint32s(executingSequences[i].Txids).Encode(),
// 			}
// 			data[i] = codec.Byteset(tmpData).Encode()
// 		}
// 	}
// 	common.ParallelWorker(len(this), concurrency, worker, this, data)
// 	return codec.Byteset(data).Encode(), nil
// }

// func (this *ExecutingSequences) Decode(data []byte) ([]*Sequence, error) {
// 	fields := codec.Byteset{}.Decode(data).(codec.Byteset)
// 	v := ExecutingSequences(make([]*Sequence, len(fields)))
// 	this = &v

// 	worker := func(start, end, idx int, args ...interface{}) {
// 		datas := args[0].([]interface{})[0].(codec.Byteset)
// 		executingSequences := args[0].([]interface{})[1].(ExecutingSequences)

// 		for i := start; i < end; i++ {
// 			executingSequence := new(Sequence)

// 			datafields := codec.Byteset{}.Decode(datas[i]).(codec.Byteset)
// 			msgResults, err := new(StandardMessages).Decode(datafields[0])
// 			if err != nil {
// 				msgResults = StandardMessages{}
// 			}
// 			executingSequence.Msgs = msgResults
// 			parallels := new(encoding.Bools).Decode(datafields[1])
// 			if len(parallels) > 0 {
// 				executingSequence.Parallel = parallels[0]
// 			}
// 			executingSequence.SequenceId = ethCommon.BytesToHash(datafields[2])
// 			executingSequence.Txids = new(encoding.Uint32s).Decode(datafields[3])
// 			executingSequences[i] = executingSequence

// 		}
// 	}
// 	common.ParallelWorker(len(fields), concurrency, worker, fields, *this)
// 	return ([]*Sequence)(*this), nil
// }
