package multiprocess

import (
	"encoding/hex"
	"errors"
	"math"

	"github.com/arcology-network/common-lib/codec"
	common "github.com/arcology-network/common-lib/common"
	evmcommon "github.com/arcology-network/evm/common"
	"github.com/arcology-network/vm-adaptor/abi"
	eucommon "github.com/arcology-network/vm-adaptor/common"
)

// APIs under the concurrency namespace
type TheadingHandler struct {
	api      eucommon.ConcurrentApiRouterInterface
	jobQueue *Queue
}

func NewMultiprocessHandler(apiRounter eucommon.ConcurrentApiRouterInterface) *TheadingHandler {
	return &TheadingHandler{
		api:      apiRounter,
		jobQueue: NewJobQueue(apiRounter),
	}
}

func (this *TheadingHandler) Address() [20]byte {
	return [20]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0x90}
}

func (this *TheadingHandler) Call(caller, callee evmcommon.Address, input []byte, origin evmcommon.Address, nonce uint64) ([]byte, bool) {
	signature := [4]byte{}
	copy(signature[:], input)

	switch signature { // bf 22 6c 78
	case [4]byte{0x0b, 0x2a, 0xcb, 0x3f}: // 0b 2a cb 3f
		return this.add(caller, callee, input[4:])

	case [4]byte{0x1f, 0x7b, 0x6d, 0x32}:
		return this.length()

	case [4]byte{0x1c, 0x82, 0xed, 0x4c}: // 1c 82 ed 4c
		return this.del(caller, callee, input[4:])

	case [4]byte{0xc0, 0x40, 0x62, 0x26}:
		return this.run(caller, callee, input[4:])

	case [4]byte{0x64, 0xf1, 0xbd, 0x63}:
		return this.get(input[4:])

	case [4]byte{0x52, 0xef, 0xea, 0x6e}:
		return this.clear()

	case [4]byte{0xb4, 0x8f, 0xb6, 0xcf}:
		return this.error(input[4:])

	}
	return this.unknow(caller, callee, input)
}

func (this *TheadingHandler) add(caller, callee evmcommon.Address, input []byte) ([]byte, bool) {
	if len(input) < 4 {
		return []byte(errors.New("Error: Invalid input").Error()), false
	}

	// fmt.Println(input)
	rawAddr, err := abi.DecodeTo(input, 0, [20]byte{}, 1, 32)
	if err != nil {
		return []byte(err.Error()), false
	}
	calleeAddr := evmcommon.BytesToAddress(rawAddr[:]) // Callee contract

	funCall, err := abi.DecodeTo(input, 1, []byte{}, 2, math.MaxUint32)
	if err != nil {
		return []byte(err.Error()), false
	}

	jobID := this.jobQueue.Add(calleeAddr, funCall)

	if buffer, err := abi.Encode(uint64(jobID)); err != nil {
		return []byte{}, false
	} else {
		return buffer, true
	}
}

func (this *TheadingHandler) clear() ([]byte, bool) {
	buffer, err := abi.Encode(codec.Uint64(this.jobQueue.Clear()))
	return buffer, err == nil
}

func (this *TheadingHandler) length() ([]byte, bool) {
	if v, err := abi.Encode(this.jobQueue.Length()); err == nil {
		return v, true
	}
	return []byte{}, false
}

func (this *TheadingHandler) error(input []byte) ([]byte, bool) {
	if idx, err := abi.DecodeTo(input, 1, uint64(0), 1, 32); err == nil {
		if item := this.jobQueue.At(idx); item != nil {
			buffer, err := abi.Encode(codec.String(item.prechkErr.Error() + item.prechkErr.Error()).Clone().(codec.String).ToBytes())
			return buffer, err == nil
		}
	}
	return []byte{}, false
}

func (this *TheadingHandler) peek(input []byte) ([]byte, bool) {
	if idx, err := abi.DecodeTo(input, 1, uint64(0), 1, 32); err == nil {
		if item := this.jobQueue.At(idx); item != nil {
			return item.message.Data(), true
		}
	}
	return []byte{}, false
}

func (this *TheadingHandler) run(caller, callee evmcommon.Address, input []byte) ([]byte, bool) {
	if threads, err := abi.DecodeTo(input, 1, uint64(0), 1, 32); err == nil {
		return []byte{}, this.jobQueue.Run(uint8(common.Min(common.Max(threads, 1), math.MaxUint8)))
	}

	return []byte{}, this.jobQueue.Run(1)
}

func (this *TheadingHandler) del(caller, callee evmcommon.Address, input []byte) ([]byte, bool) {
	if idx, err := abi.DecodeTo(input, 0, uint64(0), 1, 32); err == nil {
		this.jobQueue.Del(idx)
		return []byte{}, true
	}
	return []byte{}, false
}

func (this *TheadingHandler) get(input []byte) ([]byte, bool) {
	if idx, err := abi.DecodeTo(input, 1, uint64(0), 1, 32); err == nil {
		if item := this.jobQueue.At(idx); item != nil {
			return item.result.ReturnData, item.result.Err == nil
		}
	}
	return []byte{}, false
}

func (this *TheadingHandler) unknow(caller, callee evmcommon.Address, input []byte) ([]byte, bool) {
	this.api.AddLog("Unhandled function call in cumulative handler router", hex.EncodeToString(input))
	return []byte{}, false
}
