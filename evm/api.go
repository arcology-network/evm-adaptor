package evm

import (
	"bytes"
	"fmt"
	"math/big"

	"github.com/arcology-network/common-lib/types"
	"github.com/arcology-network/concurrentlib"
	"github.com/arcology-network/concurrenturl/v2"
	urlcommon "github.com/arcology-network/concurrenturl/v2/common"
	"github.com/arcology-network/evm/common"
)

type (
	handlerTypeV2 func(*API, common.Address, common.Address, []byte, common.Address, uint64, common.Hash, common.Hash) ([]byte, bool)
)

var (
	sysapiLookupTable = map[[20]byte]map[[4]byte]handlerTypeV2{
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0x84}: {
			[4]byte{0xac, 0xaa, 0x8d, 0x70}: New,
			[4]byte{0x46, 0xf8, 0x1a, 0x87}: dynarrayCreateV2,
			[4]byte{0x43, 0x6a, 0x66, 0xe7}: dynarrayLengthV2,
			// [4]byte{0x47, 0xf3, 0x20, 0xc8}: dynarrayPushBackUint256V2,
			[4]byte{0xbc, 0x6d, 0xb0, 0xfe}: dynarrayPushBackAddressV2,
			[4]byte{0xa1, 0x1e, 0xc6, 0xda}: dynarrayPushBackBytesV2,
			// [4]byte{0x1c, 0x00, 0xd0, 0x65}: dynarrayTryPopFrontUint256V2,
			[4]byte{0x54, 0x34, 0x8f, 0xce}: dynarrayTryPopFrontAddressV2,
			[4]byte{0x7d, 0x50, 0x8c, 0xa0}: dynarrayTryPopFrontBytesV2,
			// [4]byte{0xfd, 0xf9, 0x38, 0x44}: dynarrayPopFrontUint256V2,
			[4]byte{0x8f, 0xb0, 0x3f, 0xc5}: dynarrayPopFrontAddressV2,
			[4]byte{0x77, 0x69, 0xf8, 0xeb}: dynarrayPopFrontBytesV2,
			// [4]byte{0x98, 0x97, 0xbc, 0x69}: dynarrayTryPopBackUint256V2,
			[4]byte{0xda, 0x96, 0xd8, 0x48}: dynarrayTryPopBackAddressV2,
			[4]byte{0x32, 0x2c, 0x59, 0xa5}: dynarrayTryPopBackBytesV2,
			// [4]byte{0xde, 0xec, 0x73, 0x67}: dynarrayPopBackUint256V2,
			[4]byte{0x8f, 0x19, 0xae, 0x0e}: dynarrayPopBackAddressV2,
			[4]byte{0x9d, 0xd6, 0x79, 0x3e}: dynarrayPopBackBytesV2,
			// [4]byte{0x66, 0xf4, 0x41, 0xdc}: dynarrayTryGetUint256V2,
			[4]byte{0xc6, 0x5f, 0xf6, 0x90}: dynarrayTryGetAddressV2,
			[4]byte{0x91, 0x76, 0x9b, 0x40}: dynarrayTryGetBytesV2,
			// [4]byte{0x8d, 0x20, 0x6a, 0xad}: dynarrayGetUint256V2,
			[4]byte{0x21, 0xec, 0x12, 0x53}: dynarrayGetAddressV2,
			[4]byte{0xbf, 0x2e, 0x89, 0x60}: dynarrayGetBytesV2,
		},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0xa0}: {
			[4]byte{0x25, 0x70, 0xd9, 0xd3}: uuidGenV2,
		},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0xa1}: {
			[4]byte{0x2d, 0xc7, 0x96, 0x88}: systemCreateDeferV2,
			[4]byte{0x06, 0xe3, 0x54, 0xdd}: systemCallDeferV2,
		},
	}
)

type API struct {
	logs         []ILog
	thash        common.Hash
	tindex       uint32
	dc           *types.DeferCall
	predecessors []common.Hash

	array     *concurrentlib.FixedLengthArray
	dynarray  *concurrentlib.DynamicArray
	deferCall *concurrentlib.DeferCall

	db  urlcommon.DatastoreInterface
	url *concurrenturl.ConcurrentUrl
}

func NewAPI(db urlcommon.DatastoreInterface, url *concurrenturl.ConcurrentUrl) *API {
	return &API{
		db:  db,
		url: url,
	}
}

// Implement KernelAPI interface.
func (api *API) AddLog(key, value string) {
	api.logs = append(api.logs, &types.ExecutingLog{
		Key:   key,
		Value: value,
	})
}

func (api *API) GetLogs() []ILog {
	return api.logs
}

func (api *API) ClearLogs() {
	api.logs = api.logs[:0]
}

func (api *API) IsKernelAPI(addr common.Address) bool {
	_, ok := sysapiLookupTable[[20]byte(addr)]
	return ok
}

func (api *API) Prepare(thash common.Hash, height *big.Int, tindex uint32) {
	api.thash = thash
	api.tindex = tindex
	api.dc = nil
	context := &txContext{height, tindex}
	api.array = concurrentlib.NewFixedLengthArray(api.url, context)
	api.dynarray = concurrentlib.NewDynamicArray(api.url, context)
	api.deferCall = concurrentlib.NewDeferCall(api.url, context)
}

func (api *API) Call(caller, callee common.Address, input []byte, origin common.Address, nonce uint64, blockhash common.Hash) ([]byte, bool) {
	for contract, handlers := range sysapiLookupTable {
		if !bytes.Equal(callee.Bytes(), contract[:]) {
			continue
		}

		for method, handler := range handlers {
			if !bytes.Equal(input[:4], method[:]) {
				continue
			}
			return handler(api, caller, callee, input, origin, nonce, api.thash, blockhash)
		}
	}
	panic("unexpected method got")
}

// For defer call.
func (api *API) SetDeferCall(contractAddress types.Address, deferID string) {
	sig := api.deferCall.GetSignature(contractAddress, deferID)
	if sig == "" {
		panic(fmt.Sprintf("unknown defer call on %s:%s", contractAddress, deferID))
	}

	api.dc = &types.DeferCall{
		DeferID:         deferID,
		ContractAddress: contractAddress,
		Signature:       sig,
	}
}

func (api *API) GetDeferCall() *types.DeferCall {
	return api.dc
}

func (api *API) SetPredecessors(predecessors []common.Hash) {
	api.predecessors = predecessors
}

func (api *API) IsInDeferCall() bool {
	return len(api.predecessors) > 0
}
