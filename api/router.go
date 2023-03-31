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

	seed   uint64 // for uuid generation
	serial uint64

	// dynarray  *concurrentlib.DynamicArray
	// deferCall *concurrentlib.DeferCall
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

func (this *API) Serial() uint64 {
	this.serial++
	return this.serial
}

// Generate an UUID based on transaction hash and the counter
func (this *API) GenUUID() []byte {
	this.seed++
	id := codec.Hash32(this.txHash).UUID(this.seed)
	return id[:]
}

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

func (this *API) Call(caller, callee common.Address, input []byte, origin common.Address, nonce uint64, blockhash common.Hash) (bool, []byte, bool) {
	if callee == [20]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0x84} {
		result, successful := this.concurrencyHandler.Call(caller, input, origin, nonce)
		return true, result, successful
	}
	return false, []byte{}, false
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
