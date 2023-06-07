package types

import (
	"math/big"

	common "github.com/arcology-network/common-lib/common"
	evmcommon "github.com/arcology-network/evm/common"
	evmcoretypes "github.com/arcology-network/evm/core/types"
)

// func NewMessage(
// 	from common.Address,
// 	to *common.Address,
// 	nonce uint64,
// 	amount *big.Int,
// 	gasLimit uint64,
// 	gasPrice *big.Int,
// 	data []byte,
// 	accessList AccessList,
// 	checkNonce bool) Message {
//	return Message{
//		from:       from,
//		to:         to,
//		nonce:      nonce,
//		amount:     amount,
//		gasLimit:   gasLimit,
//		gasPrice:   gasPrice,
//		data:       data,
//		accessList: accessList,
//		checkNonce: checkNonce,
//	}
// }

const (
	TX_FROM_REMOTE = iota
	TX_FROM_LOCAL
	TX_FROM_BLOCK
	TX_FROM_DEFERRED
)

type DeferredMessageBuilder struct {
	dict map[[32]byte]*StandardMessage
}

func (this *DeferredMessageBuilder) NewDeferredMessageBuilder() *DeferredMessageBuilder {
	return &DeferredMessageBuilder{
		make(map[[32]byte]*StandardMessage),
	}
}

func (this *DeferredMessageBuilder) ToStandardMessages() []*StandardMessage {
	return common.MapValues(this.dict)
}

func (this *DeferredMessageBuilder) Add(results []*ExecutionResult) {
	for _, result := range results {
		message := this.dict[result.Deferred.Signature]
		if message != nil {
			message.Predecessors = append(message.Predecessors, result.TxHash)
		}

		var to evmcommon.Address
		to.SetBytes(result.Deferred.Signature[:20])

		evmMsg := evmcoretypes.NewMessage(
			evmcommon.Address(result.Deferred.From), // From the system account
			&to,
			0,
			big.NewInt(0),
			1e15,
			big.NewInt(1),
			result.Deferred.FuncCallData,
			nil,
			false,
		)

		this.dict[result.Deferred.Signature] = &StandardMessage{
			TxHash:       result.Deferred.Signature,
			Predecessors: [][32]byte{result.TxHash},
			Native:       &evmMsg,
			Source:       TX_FROM_DEFERRED,
		}
	}
}

// type evmcoretypes.Message struct {
// 	groupBy   [32]byte
// 	callerTxs [][32]byte // Predecessors
// 	msg       *evmcoretypes.Message
// }

// func (this *evmcoretypes.Message) BuildEthMsg() *evmcoretypes.Message {
// 	from := sha3.Sum256(codec.Bytes32s(this.callerTxs).Flatten())
// 	to := evmcommon.Address{}
// 	to.SetBytes(this.msg.Data()[4:])

// 	msg := evmcoretypes.NewMessage(
// 		evmcommon.Address(codec.Bytes20{}.Decode(from[:]).(codec.Bytes20)),
// 		&to,
// 		0,
// 		big.NewInt(0),
// 		1e15,
// 		big.NewInt(1),
// 		this.msg.Data(),
// 		nil,
// 		false,
// 	)
// 	return &msg
// }
