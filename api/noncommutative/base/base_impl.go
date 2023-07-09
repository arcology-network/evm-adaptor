package concurrentcontainer

import (
	"math"

	"github.com/arcology-network/concurrenturl/noncommutative"
	abi "github.com/arcology-network/vm-adaptor/abi"
	"github.com/arcology-network/vm-adaptor/execution"
)

// // get the number of elements in the container
func (this *BytesHandlers) Length(path string) (uint64, bool, int64) {
	if len(path) == 0 {
		return 0, false, 0
	}

	if path, _ := this.api.Ccurl().Read(uint32(this.api.GetEU().(*execution.EU).Message().ID), path); path != nil {
		return uint64(len(path.([]string))), true, 0
	}
	return 0, false, 0
}

// // get the number of elements in the container
func (this *BytesHandlers) Get(path string, idx uint64) ([]byte, bool, int64) {
	if value, _, err := this.api.Ccurl().ReadAt(uint32(this.api.GetEU().(*execution.EU).Message().ID), path, idx); err == nil && value != nil {
		return value.([]byte), true, 0
	}
	return []byte{}, false, 0
}

func (this *BytesHandlers) Set(path string, idx uint64, bytes []byte) (bool, int64) {
	if len(path) > 0 {
		value := noncommutative.NewBytes(bytes)
		if _, err := this.api.Ccurl().WriteAt(uint32(this.api.GetEU().(*execution.EU).Message().ID), path, idx, value, true); err == nil {
			return true, 0
		}
	}
	return false, 0
}

func (this *BytesHandlers) Push(path string, input []byte) ([]byte, bool, int64) {
	if len(path) == 0 {
		return []byte{}, false, 0
	}

	// if this.deploymentAddr != caller {
	// 	fmt.Println(caller[:])
	// 	fmt.Println(this.deploymentAddr[:])
	// 	panic("Mismatch !!!!!")
	// }

	value, err := abi.DecodeTo(input, 0, []byte{}, 2, math.MaxInt)
	if value == nil || err != nil {
		return []byte{}, false, 0
	}

	key := path + string(this.api.ElementUID())
	_, err = this.api.Ccurl().Write(uint32(this.api.GetEU().(*execution.EU).Message().ID), key, noncommutative.NewBytes(value), true)
	return []byte{}, err == nil, 0
}

// func (this *BytesHandlers) run(caller, callee evmcommon.Address, input []byte) ([]byte, bool, int64) {
// 	path := this.Connector().Key(caller)

// 	length, successful, fee := this.Length(path)
// 	if !successful {
// 		return []byte{}, successful, fee
// 	}

// 	generation := execution.NewGeneration(0, 4, []*execution.JobSequence{})

// 	fees := make([]int64, length)
// 	erros := make([]error, length)
// 	jobseqs := make([]*execution.JobSequence, length)
// 	for i := uint64(0); i < length; i++ {
// 		data, successful, fee := this.Get(path, uint64(i))
// 		if fees[i] = fee; successful {
// 			jobseqs[i], erros[i] = this.toJobSeq(data)
// 		}
// 		generation.Add(jobseqs[i])
// 	}

// 	generation.Run(this.Api())
// 	return []byte{}, true, common.Sum(fees, int64(0))
// }

// func (this *BytesHandlers) toJobSeq(input []byte) (*execution.JobSequence, error) {
// 	input, err := abi.DecodeTo(input, 0, []byte{}, 2, math.MaxInt64)
// 	if err != nil {
// 		return nil, errors.New("Error: Unrecognizable input format")
// 	}

// 	gasLimit, calleeAddr, funCall, err := abi.Parse3(input,
// 		uint64(0), 1, 32,
// 		[20]byte{}, 1, 32,
// 		[]byte{}, 2, math.MaxInt64)

// 	if err != nil {
// 		return nil, err
// 	}

// 	newJobSeq := &execution.JobSequence{
// 		ID:        this.Api().GetSerialNum(eucommon.SUB_PROCESS),
// 		ApiRouter: this.Api(),
// 	}

// 	addr := evmcommon.Address(calleeAddr)
// 	evmMsg := evmcore.NewMessage( // Build the message
// 		this.Api().Origin(),
// 		&addr,
// 		0,
// 		new(big.Int).SetUint64(0), // Amount to transfer
// 		gasLimit,
// 		this.Api().GetEU().(*execution.EU).Message().Native.GasPrice, // gas price
// 		funCall,
// 		nil,
// 		false, // Don't checking nonce
// 	)

// 	stdMsg := &execution.StandardMessage{
// 		ID:     newJobSeq.ID, // this is the problem !!!!
// 		Native: &evmMsg,
// 		TxHash: newJobSeq.DeriveNewHash(this.Api().GetEU().(*execution.EU).Message().TxHash),
// 	}

// 	newJobSeq.StdMsgs = []*execution.StandardMessage{stdMsg}
// 	return newJobSeq, nil
// }
