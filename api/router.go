package api

import (
	"math/big"

	"github.com/arcology-network/common-lib/codec"
	"github.com/arcology-network/common-lib/types"

	"github.com/arcology-network/concurrenturl/v2"
	"github.com/arcology-network/evm/common"
	euCommon "github.com/arcology-network/vm-adaptor/common"
)

var (
	apiNamespaces = map[[20]byte]bool{
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0x84}: true,
	}
)

type API struct {
	logs         []euCommon.ILog
	txHash       common.Hash // Tx hash
	txIndex      uint32      // Tx index in the block
	dc           *types.DeferCall
	predecessors []common.Hash

	seed uint64 // for uuid generation

	// dynarray  *concurrentlib.DynamicArray
	// deferCall *concurrentlib.DeferCall

	// db  urlcommon.DatastoreInterface
	concurrencyHandler *ConcurrentContainer // APIs under the concurrency namespace
	ccurl              *concurrenturl.ConcurrentUrl
}

func NewAPI(ccurl *concurrenturl.ConcurrentUrl) *API {
	return &API{
		ccurl: ccurl,
	}
}

func (this *API) Prepare(txHash common.Hash, height *big.Int, txIndex uint32) {
	this.txHash = txHash
	this.txIndex = txIndex
	this.dc = nil

	this.concurrencyHandler = NewConcurrentContainer(txHash, txIndex, this)
}

// Generate an UUID based on transaction hash and the counter
func (this *API) GenUUID() []byte {
	this.seed++
	id := codec.Hash32(this.txHash).UUID(this.seed)
	return id[:]
}

// Implement KernelAPI interface.
func (this *API) AddLog(key, value string) {
	this.logs = append(this.logs, &types.ExecutingLog{
		Key:   key,
		Value: value,
	})
}

func (this *API) GetLogs() []euCommon.ILog {
	return this.logs
}

func (this *API) ClearLogs() {
	this.logs = this.logs[:0]
}

func (this *API) IsKernelAPI(callee common.Address) bool {
	_, ok := apiNamespaces[[20]byte(callee)]
	return ok
}

func (this *API) Call(caller, callee common.Address, input []byte, origin common.Address, nonce uint64, blockhash common.Hash) ([]byte, bool) {
	_, ok := apiNamespaces[callee]
	if !ok {
		panic("Should never enter here")
	}

	if callee == [20]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0x84} {
		return this.concurrencyHandler.Call(caller, input, origin, nonce)
	}

	panic("unexpected method got")

	return []byte{}, false
}

// For defer call.
func (this *API) SetDeferCall(contractAddress types.Address, deferID string) {
	// sig := this.deferCall.GetSignature(contractAddress, deferID)
	// if sig == "" {
	// 	panic(fmt.Sprintf("unknown defer call on %s:%s", contractAddress, deferID))
	// }

	this.dc = &types.DeferCall{
		DeferID:         deferID,
		ContractAddress: contractAddress,
		// Signature:       sig,
	}
}

func (this *API) GetDeferCall() *types.DeferCall {
	return this.dc
}

func (this *API) SetPredecessors(predecessors []common.Hash) {
	this.predecessors = predecessors
}

func (this *API) IsInDeferCall() bool {
	return len(this.predecessors) > 0
}
