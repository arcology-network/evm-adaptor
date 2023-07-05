package threading

import (
	"errors"
	"math"
	"math/big"

	evmcommon "github.com/arcology-network/evm/common"
	evmcore "github.com/arcology-network/evm/core"

	"github.com/arcology-network/vm-adaptor/abi"
	base "github.com/arcology-network/vm-adaptor/api/noncommutative/base"

	eucommon "github.com/arcology-network/vm-adaptor/common"
	execution "github.com/arcology-network/vm-adaptor/execution"
)

// APIs under the concurrency namespace
type ParallelHandler struct {
	*base.BytesHandlers
}

func NewParallelHandler(ethApiRouter eucommon.EthApiRouter) *ParallelHandler {
	return &ParallelHandler{
		base.NewNoncommutativeBytesHandlers(ethApiRouter),
	}
}

func (this *ParallelHandler) Address() [20]byte {
	return eucommon.THREADING_HANDLER
}

func (this *ParallelHandler) Call(caller, callee [20]byte, input []byte, origin [20]byte, nonce uint64) ([]byte, bool, int64) {
	if returned, called, fees := this.BytesHandlers.Call(caller, callee, input, origin, nonce); called {
		return returned, called, fees
	}

	signature := [4]byte{}
	copy(signature[:], input)

	switch signature { // bf 22 6c 78
	case [4]byte{0xc0, 0x40, 0x62, 0x26}: //
		return this.run(caller, callee, input[4:])
	}

	return []byte{}, false, 0
}

func (this *ParallelHandler) run(caller, callee evmcommon.Address, input []byte) ([]byte, bool, int64) {
	// id := this.GetAddress()
	// if len(id) == 0 || this.pools[id] == nil {
	// 	return []byte{}, false, 0
	// }

	// this.pools[id].Run(this.BytesHandlers.Api())

	// for i := 0; i < this.BytesHandlers.Length(caller); i++ {

	// }

	return []byte{}, true, 0
}

func (this *ParallelHandler) parse(input []byte) (*execution.JobSequence, error) {

	// id := this.GetAddress()
	// address := this.api.VM().ArcologyNetworkAPIs.CallContext.Contract.CodeAddr[:]
	// id = string(address)

	// if len(id) == 0 || this.pools[id] == nil {
	// 	return []byte{}, false, 0
	// }
	// fmt.Println(input)
	gasLimit, err := abi.DecodeTo(input, 0, uint64(0), 1, 32)
	if err != nil {
		return nil, errors.New("Error: Failed to part gas limit")
	}

	rawAddr, err := abi.DecodeTo(input, 1, [20]byte{}, 1, 32)
	if err != nil {
		return nil, errors.New("Error: Failed to parse callee Addr")
	}
	calleeAddr := evmcommon.BytesToAddress(rawAddr[:]) // Callee contract

	funCall, err := abi.DecodeTo(input, 2, []byte{}, 2, math.MaxUint32)
	if err != nil {
		return nil, errors.New("Error: Failed to parse callee function call data")
	}

	newSeq := &execution.JobSequence{
		ID:        this.BytesHandlers.Api().GetSerialNum(eucommon.SUB_PROCESS),
		ApiRouter: this.BytesHandlers.Api(),
	}

	evmMsg := evmcore.NewMessage( // Build the message
		this.BytesHandlers.Api().Origin(),
		&calleeAddr,
		0,
		new(big.Int).SetUint64(0), // Amount to transfer
		gasLimit,
		this.BytesHandlers.Api().GetEU().(*execution.EU).Message().Native.GasPrice, // gas price
		funCall,
		nil,
		false, // Don't checking nonce
	)

	stdMsg := &execution.StandardMessage{
		ID:     0, //Can ONly be 1 msg at the most
		Native: &evmMsg,
		TxHash: newSeq.DeriveNewHash(this.BytesHandlers.Api().GetEU().(*execution.EU).Message().TxHash),
	}

	newSeq.StdMsgs = []*execution.StandardMessage{stdMsg}
	return newSeq, nil
}
