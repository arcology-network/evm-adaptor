package threading

import (
	"crypto/sha256"
	"math"
	"math/big"
	"strconv"
	"sync/atomic"

	"github.com/arcology-network/common-lib/codec"
	"github.com/arcology-network/common-lib/common"
	commonlibcommon "github.com/arcology-network/common-lib/common"
	"github.com/arcology-network/concurrenturl/indexer"
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

	case [4]byte{0xf8, 0xf0, 0xd9, 0x80}:
		return this.add(caller, callee, input[4:])

	case [4]byte{0x84, 0x67, 0x3c, 0xc9}:
		return this.length(input[4:])

	case [4]byte{0x3a, 0x27, 0x65, 0x23}: //3a 27 65 23
		return this.run(caller, callee, input[4:])

	case [4]byte{0x4d, 0xd4, 0x9a, 0xb4}: // 4d d4 9a b4
		return this.get(input[4:])

	case [4]byte{0x5e, 0x1d, 0x05, 0x4d}: // 5e 1d 05 4d
		return this.clear(input[4:])
	}

	return []byte{}, false, 0
}

func (this *ThreadingHandler) new(caller, callee evmcommon.Address, input []byte) ([]byte, bool, int64) {
	if this.api.Depth() >= eucommon.MAX_RECURSIION_DEPTH ||
		atomic.AddUint64(&eucommon.TotalSubProcesses, 1) > eucommon.MAX_SUB_PROCESSES {
		return []byte{}, false, 0 // Execeeds the max recursion depth or the max sub processes
	}

	threads, err := abi.DecodeTo(input, 0, uint8(1), 1, 32)
	if err != nil {
		return []byte{}, false, 0
	}

	id := strconv.Itoa(len(this.pools))
	this.pools[id] = execution.NewGeneration(uint32(len(this.pools)), 0, threads, []*execution.JobSequence{})
	return []byte(id), true, 0 // Create a new container
}

func (this *ThreadingHandler) add(caller, callee evmcommon.Address, input []byte) ([]byte, bool, int64) {
	if len(input) < 4 {
		return []byte{}, false, 0
	}

	id := this.ParseID(input)
	if len(id) == 0 || this.pools[id] == nil {
		return []byte{}, false, 0
	}

	// fmt.Println(input)
	gasLimit, err := abi.DecodeTo(input, 1, uint64(0), 1, 32)
	if err != nil {
		return []byte{}, false, 0
	}

	rawAddr, err := abi.DecodeTo(input, 2, [20]byte{}, 1, 32)
	if err != nil {
		return []byte{}, false, 0
	}
	calleeAddr := evmcommon.BytesToAddress(rawAddr[:]) // Callee contract

	funCall, err := abi.DecodeTo(input, 3, []byte{}, 2, math.MaxUint32)
	if err != nil {
		return []byte{}, false, 0
	}

	evmMsg := evmcore.NewMessage( // Build the message
		this.api.Origin(),
		&calleeAddr,
		0,
		new(big.Int).SetUint64(0), // Amount to transfer
		gasLimit,
		this.api.Message().GasPrice, // gas price
		funCall,
		nil,
		false, // Don't checking nonce
	)

	newJob := &execution.JobSequence{
		ID:        this.pools[id].Length(),
		ApiRouter: this.api,
	}

	stdMsg := &execution.StandardMessage{
		ID:     this.api.CCUID(),
		Native: &evmMsg,
		TxHash: sha256.Sum256(commonlibcommon.Flatten([][]byte{
			codec.Bytes32(this.api.TxHash()).Encode(),
			codec.Uint32((this.pools[id].BranchID())).Encode(),
			evmMsg.Data[:4],
			codec.Uint32(newJob.ID).Encode(),
		})),
	}

	newJob.StdMsgs = []*execution.StandardMessage{stdMsg}
	return []byte{}, this.pools[id].Add(newJob), 0
}

func (this *ThreadingHandler) clear(input []byte) ([]byte, bool, int64) {
	id := this.ParseID(input)
	if len(id) == 0 || this.pools[id] == nil {
		return []byte{}, false, 0
	}

	this.pools[id].Clear()
	return []byte{}, true, 0
}

func (this *ThreadingHandler) length(input []byte) ([]byte, bool, int64) {
	id := this.ParseID(input)
	if len(id) == 0 || this.pools[id] == nil {
		return []byte{}, false, 0
	}

	v, err := abi.Encode(this.pools[id].Length())
	return v, err == nil, 0
}

func (this *ThreadingHandler) run(caller, callee evmcommon.Address, input []byte) ([]byte, bool, int64) {
	id := this.ParseID(input)
	if len(id) == 0 || this.pools[id] == nil {
		return []byte{}, false, 0
	}

	preTransitions := indexer.Univalues(common.Clone(this.api.Ccurl().Export())).To(indexer.ITCTransition{})
	this.pools[id].Run(this.api, this.api.Ccurl().Snapshot(preTransitions))
	return []byte{}, true, 0
}

func (this *ThreadingHandler) get(input []byte) ([]byte, bool, int64) {
	id := this.ParseID(input)
	if len(id) == 0 || this.pools[id] == nil {
		return []byte{}, false, 0
	}

	if idx, err := abi.DecodeTo(input, 1, uint64(0), 1, 32); err == nil {
		if item := this.pools[id].At(idx); item != nil && item.Results[0].EvmResult != nil {
			return item.Results[0].EvmResult.ReturnData, item.Results[0].EvmResult.Err == nil, 0
		}
	}
	return []byte{}, false, 0
}

func (this *ThreadingHandler) error(input []byte) ([]byte, bool, int64) {
	// id := this.ParseID(input)
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
func (this *ThreadingHandler) ParseID(input []byte) string {
	id, _ := abi.DecodeTo(input, 0, []byte{}, 2, 32) // max 32 bytes                                                                          // container ID
	return string(id)                                // unique ID
}
