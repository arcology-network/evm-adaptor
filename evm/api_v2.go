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
	handlerTypeV2 func(*APIV2, common.Address, common.Address, []byte, common.Address, uint64, common.Hash, common.Hash) ([]byte, bool)
)

var (
	sysapiLookupTable = map[[20]byte]map[[4]byte]handlerTypeV2{
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0x80}: {
			[4]byte{0xf0, 0x2e, 0x3a, 0xff}: arrayCreateV2,
			[4]byte{0xce, 0x86, 0x99, 0x03}: arraySizeV2,
			[4]byte{0xbf, 0x39, 0x13, 0x86}: arraySetUint256V2,
			[4]byte{0x42, 0x45, 0x0b, 0x03}: arrayGetUint256V2,
			[4]byte{0xd7, 0x33, 0xe7, 0x67}: arraySetAddressV2,
			[4]byte{0x3f, 0xb7, 0xde, 0x0c}: arrayGetAddressV2,
			[4]byte{0xe9, 0x7e, 0x40, 0x65}: arraySetBytesV2,
			[4]byte{0x6e, 0xee, 0x17, 0x2d}: arrayGetBytesV2,
		},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0x81}: {
			[4]byte{0xf0, 0x2e, 0x3a, 0xff}: hashmapCreateV2,
			[4]byte{0x21, 0xec, 0x12, 0x53}: hashmapUint256AddressGetV2,
			[4]byte{0xce, 0x58, 0x62, 0x1e}: hashmapUint256AddressSetV2,
			[4]byte{0x8d, 0x20, 0x6a, 0xad}: hashmapUint256Uint256GetV2,
			[4]byte{0x36, 0xf3, 0xc7, 0x7d}: hashmapUint256Uint256SetV2,
			[4]byte{0xbf, 0x2e, 0x89, 0x60}: hashmapUint256BytesGetV2,
			[4]byte{0x66, 0x84, 0x5e, 0xad}: hashmapUint256BytesSetV2,
			[4]byte{0xe5, 0xc2, 0xbe, 0x85}: hashmapAddressAddressGetV2,
			[4]byte{0xd6, 0xf5, 0x2d, 0xbe}: hashmapAddressAddressSetV2,
			[4]byte{0xc4, 0x1e, 0xb8, 0x5a}: hashmapAddressUint256GetV2,
			[4]byte{0x4f, 0x7c, 0x4f, 0x4c}: hashmapAddressUint256SetV2,
			[4]byte{0xcb, 0xb2, 0x89, 0xeb}: hashmapAddressBytesGetV2,
			[4]byte{0x4d, 0x57, 0x80, 0x65}: hashmapAddressBytesSetV2,
			[4]byte{0x39, 0xc2, 0x3c, 0x41}: hashmapDeleteKeyUint256V2,
			[4]byte{0x5d, 0x26, 0x1f, 0xfa}: hashmapDeleteKeyAddressV2,
		},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0x82}: {
			[4]byte{0x46, 0xf8, 0x1a, 0x87}: queueCreateV2, // create(string,uint256)
			[4]byte{0xce, 0x86, 0x99, 0x03}: queueSizeV2,
			[4]byte{0xa0, 0xaa, 0x9f, 0x29}: queuePushUint256V2,
			[4]byte{0xf6, 0x1f, 0xe1, 0x44}: queuePopUint256V2,
		},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0x83}: {
			[4]byte{0xfb, 0x99, 0xd9, 0x25}: varCreateV2,
			[4]byte{0x8a, 0x42, 0xeb, 0xe9}: varSetUint256V2,
			[4]byte{0x0b, 0xb6, 0x87, 0xe3}: varGetUint256V2,
			[4]byte{0xa8, 0x15, 0xff, 0x15}: varSetAddressV2,
			[4]byte{0xbf, 0x40, 0xfa, 0xc1}: varGetAddressV2,
			[4]byte{0x2b, 0x29, 0xc0, 0xfa}: varSetBytesV2,
			[4]byte{0xd8, 0xde, 0x89, 0x9d}: varGetBytesV2,
		},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0x84}: {
			[4]byte{0x46, 0xf8, 0x1a, 0x87}: dynarrayCreateV2,
			[4]byte{0x43, 0x6a, 0x66, 0xe7}: dynarrayLengthV2,
			[4]byte{0x47, 0xf3, 0x20, 0xc8}: dynarrayPushBackUint256V2,
			[4]byte{0xbc, 0x6d, 0xb0, 0xfe}: dynarrayPushBackAddressV2,
			[4]byte{0xa1, 0x1e, 0xc6, 0xda}: dynarrayPushBackBytesV2,
			[4]byte{0x1c, 0x00, 0xd0, 0x65}: dynarrayTryPopFrontUint256V2,
			[4]byte{0x54, 0x34, 0x8f, 0xce}: dynarrayTryPopFrontAddressV2,
			[4]byte{0x7d, 0x50, 0x8c, 0xa0}: dynarrayTryPopFrontBytesV2,
			[4]byte{0xfd, 0xf9, 0x38, 0x44}: dynarrayPopFrontUint256V2,
			[4]byte{0x8f, 0xb0, 0x3f, 0xc5}: dynarrayPopFrontAddressV2,
			[4]byte{0x77, 0x69, 0xf8, 0xeb}: dynarrayPopFrontBytesV2,
			[4]byte{0x98, 0x97, 0xbc, 0x69}: dynarrayTryPopBackUint256V2,
			[4]byte{0xda, 0x96, 0xd8, 0x48}: dynarrayTryPopBackAddressV2,
			[4]byte{0x32, 0x2c, 0x59, 0xa5}: dynarrayTryPopBackBytesV2,
			[4]byte{0xde, 0xec, 0x73, 0x67}: dynarrayPopBackUint256V2,
			[4]byte{0x8f, 0x19, 0xae, 0x0e}: dynarrayPopBackAddressV2,
			[4]byte{0x9d, 0xd6, 0x79, 0x3e}: dynarrayPopBackBytesV2,
			[4]byte{0x66, 0xf4, 0x41, 0xdc}: dynarrayTryGetUint256V2,
			[4]byte{0xc6, 0x5f, 0xf6, 0x90}: dynarrayTryGetAddressV2,
			[4]byte{0x91, 0x76, 0x9b, 0x40}: dynarrayTryGetBytesV2,
			[4]byte{0x8d, 0x20, 0x6a, 0xad}: dynarrayGetUint256V2,
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

type txContext struct {
	height *big.Int
	index  uint32
}

func (c *txContext) GetIndex() uint32 {
	return c.index
}

func (c *txContext) GetHeight() *big.Int {
	return c.height
}

type APIV2 struct {
	logs         []ILog
	thash        common.Hash
	tindex       uint32
	dc           *types.DeferCall
	predecessors []common.Hash

	array     *concurrentlib.FixedLengthArray
	sortedMap *concurrentlib.SortedMap
	queue     *concurrentlib.Queue
	dynarray  *concurrentlib.DynamicArray
	deferCall *concurrentlib.DeferCall

	db  urlcommon.DatastoreInterface
	url *concurrenturl.ConcurrentUrl
}

func NewAPIV2(db urlcommon.DatastoreInterface, url *concurrenturl.ConcurrentUrl) *APIV2 {
	return &APIV2{
		db:  db,
		url: url,
	}
}

// Implement KernelAPI interface.
func (api *APIV2) AddLog(key, value string) {
	api.logs = append(api.logs, &types.ExecutingLog{
		Key:   key,
		Value: value,
	})
}

func (api *APIV2) GetLogs() []ILog {
	return api.logs
}

func (api *APIV2) ClearLogs() {
	api.logs = api.logs[:0]
}

func (api *APIV2) IsKernelAPI(addr common.Address) bool {
	_, ok := sysapiLookupTable[[20]byte(addr)]
	return ok
}

func (api *APIV2) Prepare(thash common.Hash, height *big.Int, tindex uint32) {
	api.thash = thash
	api.tindex = tindex
	api.dc = nil
	context := &txContext{height, tindex}
	api.array = concurrentlib.NewFixedLengthArray(api.url, context)
	api.sortedMap = concurrentlib.NewSortedMap(api.url, context)
	api.queue = concurrentlib.NewQueue(api.url, context)
	api.dynarray = concurrentlib.NewDynamicArray(api.url, context)
	api.deferCall = concurrentlib.NewDeferCall(api.url, context)
}

func (api *APIV2) Call(caller, callee common.Address, input []byte, origin common.Address, nonce uint64, blockhash common.Hash) ([]byte, bool) {
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
func (api *APIV2) SetDeferCall(contractAddress types.Address, deferID string) {
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

func (api *APIV2) GetDeferCall() *types.DeferCall {
	return api.dc
}

func (api *APIV2) SetPredecessors(predecessors []common.Hash) {
	api.predecessors = predecessors
}

func (api *APIV2) IsInDeferCall() bool {
	return len(api.predecessors) > 0
}
