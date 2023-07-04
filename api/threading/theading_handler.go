package threading

import (
	"math"
	"math/big"

	evmcommon "github.com/arcology-network/evm/common"
	evmcore "github.com/arcology-network/evm/core"
	"github.com/arcology-network/vm-adaptor/abi"

	eucommon "github.com/arcology-network/vm-adaptor/common"
	execution "github.com/arcology-network/vm-adaptor/execution"
)

// APIs under the concurrency namespace
type ThreadingHandler struct {
	api   eucommon.EthApiRouter
	pools map[string]*execution.Generation
}

func NewThreadingHandler(ethApiRouter eucommon.EthApiRouter) *ThreadingHandler {
	return &ThreadingHandler{
		api:   ethApiRouter,
		pools: map[string]*execution.Generation{},
	}
}

func (this *ThreadingHandler) Address() [20]byte {
	return eucommon.THREADING_HANDLER
}

func (this *ThreadingHandler) Call(caller, callee [20]byte, input []byte, origin [20]byte, nonce uint64) ([]byte, bool, int64) {
	signature := [4]byte{}
	copy(signature[:], input)

	switch signature { // bf 22 6c 78

	case [4]byte{0x58, 0x16, 0xc4, 0x25}:
		return this.new(caller, callee, input[4:])

	case [4]byte{0x55, 0x46, 0x3c, 0x05}:
		return this.add(caller, callee, input[4:])

	case [4]byte{0x1f, 0x7b, 0x6d, 0x32}: //
		return this.length(input[4:])

	case [4]byte{0xc0, 0x40, 0x62, 0x26}: //
		return this.run(caller, callee, input[4:])

	case [4]byte{0x95, 0x07, 0xd3, 0x9a}: //
		return this.get(input[4:])

	case [4]byte{0x52, 0xef, 0xea, 0x6e}:
		return this.clear(input[4:])
	}

	return []byte{}, false, 0
}

func (this *ThreadingHandler) new(caller, callee evmcommon.Address, input []byte) ([]byte, bool, int64) {
	// if this.api.Depth() >= eucommon.MAX_RECURSIION_DEPTH ||
	// 	atomic.AddUint64(&eucommon.TotalProcesses, 1) > eucommon.MAX_SUB_PROCESSES {
	// 	return []byte{}, false, 0 // Execeeds the max recursion depth or the max sub processes
	// }

	threads, err := abi.DecodeTo(input, 0, uint8(1), 1, 32)
	if err != nil {
		return []byte{}, false, 0
	}

	address := this.GetAddress()
	this.pools[string(address)] = execution.NewGeneration(uint32(len(this.pools)), threads, []*execution.JobSequence{})
	return []byte(address), true, 0 // Create a new container
}

func (this *ThreadingHandler) add(caller, callee evmcommon.Address, input []byte) ([]byte, bool, int64) {
	if len(input) < 4 {
		return []byte{}, false, 0
	}

	id := this.GetAddress()
	address := this.api.VM().ArcologyNetworkAPIs.CallContext.Contract.CodeAddr[:]
	id = string(address)

	if len(id) == 0 || this.pools[id] == nil {
		return []byte{}, false, 0
	}

	// fmt.Println(input)
	gasLimit, err := abi.DecodeTo(input, 0, uint64(0), 1, 32)
	if err != nil {
		return []byte{}, false, 0
	}

	rawAddr, err := abi.DecodeTo(input, 1, [20]byte{}, 1, 32)
	if err != nil {
		return []byte{}, false, 0
	}
	calleeAddr := evmcommon.BytesToAddress(rawAddr[:]) // Callee contract

	funCall, err := abi.DecodeTo(input, 2, []byte{}, 2, math.MaxUint32)
	if err != nil {
		return []byte{}, false, 0
	}

	newJob := &execution.JobSequence{
		ID:        this.api.GetSerialNum(eucommon.SUB_PROCESS),
		ApiRouter: this.api,
	}

	// txHash := newJob.DeriveNewHash(this.api.GetEU().(*execution.EU).Message().TxHash)
	evmMsg := evmcore.NewMessage( // Build the message
		this.api.Origin(),
		&calleeAddr,
		0,
		new(big.Int).SetUint64(0), // Amount to transfer
		gasLimit,
		this.api.GetEU().(*execution.EU).Message().Native.GasPrice, // gas price
		funCall,
		nil,
		false, // Don't checking nonce
	)

	stdMsg := &execution.StandardMessage{
		ID:     newJob.ID, // this is the problem !!!!
		Native: &evmMsg,
		TxHash: newJob.DeriveNewHash(this.api.GetEU().(*execution.EU).Message().TxHash),
	}

	newJob.StdMsgs = []*execution.StandardMessage{stdMsg}
	return []byte{}, this.pools[id].Add(newJob), 0
}

func (this *ThreadingHandler) clear(input []byte) ([]byte, bool, int64) {
	id := this.GetAddress()
	if len(id) == 0 || this.pools[id] == nil {
		return []byte{}, false, 0
	}

	this.pools[id].Clear()
	return []byte{}, true, 0
}

func (this *ThreadingHandler) length(input []byte) ([]byte, bool, int64) {
	id := this.GetAddress()
	if len(id) == 0 || this.pools[id] == nil {
		return []byte{}, false, 0
	}

	v, err := abi.Encode(this.pools[id].Length())
	return v, err == nil, 0
}

func (this *ThreadingHandler) run(caller, callee evmcommon.Address, input []byte) ([]byte, bool, int64) {
	id := this.GetAddress()
	if len(id) == 0 || this.pools[id] == nil {
		return []byte{}, false, 0
	}

	this.pools[id].Run(this.api)
	return []byte{}, true, 0
}

func (this *ThreadingHandler) get(input []byte) ([]byte, bool, int64) {
	id := this.GetAddress()
	if len(id) == 0 || this.pools[id] == nil {
		return []byte{}, false, 0
	}

	if idx, err := abi.DecodeTo(input, 0, uint64(0), 1, 32); err == nil {
		if item := this.pools[id].At(idx); item != nil && item.Results[0].EvmResult != nil {
			return item.Results[0].EvmResult.ReturnData, item.Results[0].EvmResult.Err == nil, 0
		}
	}
	return []byte{}, false, 0
}

func (this *ThreadingHandler) error(input []byte) ([]byte, bool, int64) {
	// id := this.GetAddress()
	// if len(id) == 0 || this.pools[id] == nil {
	// 	return []byte{}, false
	// }

	// if idx, err := abi.DecodeTo(input, 1, uint64(0), 1, 32); err == nil {
	// 	if item := this.pools[id].At(idx); item != nil {
	// 		buffer, err := abi.Encode(item.Err.Error())
	// 		return buffer, err == nil
	// 	}
	// }
	return []byte{}, false, 0
}

// Build the container path
// func (this *ThreadingHandler) GetAddress(input []byte) string {
// 	id, _ := abi.DecodeTo(input, 0, []byte{}, 2, 32) // max 32 bytes                                                                          // container ID
// 	return string(id)                                // unique ID
// }

// Build the container path
func (this *ThreadingHandler) GetAddress() string {
	// id, _ := abi.DecodeTo(input, 0, []byte{}, 2, 32) // max 32 bytes                                                                          // container ID
	return string(this.api.VM().ArcologyNetworkAPIs.CallContext.Contract.CodeAddr[:]) // unique ID
}
