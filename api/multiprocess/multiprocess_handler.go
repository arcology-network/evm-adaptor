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
type MultiprocessHandler struct {
	api        eucommon.ConcurrentApiRouterInterface
	jobManager *JobManager
}

func NewParallelHandler(apiRounter eucommon.ConcurrentApiRouterInterface) *MultiprocessHandler {
	return &MultiprocessHandler{
		api:        apiRounter,
		jobManager: NewJobManager(apiRounter),
	}
}

func (this *MultiprocessHandler) Address() [20]byte {
	return [20]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0x90}
}

func (this *MultiprocessHandler) Call(caller, callee evmcommon.Address, input []byte, origin evmcommon.Address, nonce uint64) ([]byte, bool) {
	signature := [4]byte{}
	copy(signature[:], input)

	switch signature { // bf 22 6c 78
	case [4]byte{0xa4, 0x62, 0x12, 0x2d}: // a4 62 12 2d
		return this.addJob(caller, callee, input[4:])

	case [4]byte{0xe0, 0x88, 0x6f, 0x90}:
		return this.at(input[4:])

	case [4]byte{0xb6, 0xff, 0x8b, 0xd9}:
		return this.delJob(caller, callee, input[4:])

	case [4]byte{0x1f, 0x7b, 0x6d, 0x32}:
		return this.length()

	case [4]byte{0xc0, 0x40, 0x62, 0x26}:
		return this.run(caller, callee, input[4:])

	case [4]byte{0xb4, 0x8f, 0xb6, 0xcf}:
		return this.error(input[4:])

	case [4]byte{0x64, 0x17, 0x43, 0x08}:
		return this.clear()
	}
	return this.unknow(caller, callee, input)
}

func (this *MultiprocessHandler) clear() ([]byte, bool) {
	buffer := [32]byte{}
	codec.Uint64(this.jobManager.Clear()).EncodeToBuffer(buffer[len(buffer)-8:])
	return buffer[:], true
}

func (this *MultiprocessHandler) length() ([]byte, bool) {
	if v, err := abi.Encode(this.jobManager.Length()); err == nil {
		return v, true
	}
	return []byte{}, false
}

func (this *MultiprocessHandler) error(input []byte) ([]byte, bool) {
	if idx, err := abi.DecodeTo(input, 1, uint64(0), 1, 32); err == nil {
		if _, err := this.jobManager.At(idx); err != nil {
			return codec.String(err.Error()).ToBytes(), err == nil
		}
	}
	return []byte{}, false
}

func (this *MultiprocessHandler) at(input []byte) ([]byte, bool) {
	if idx, err := abi.DecodeTo(input, 1, uint64(0), 1, 32); err == nil {
		data, err := this.jobManager.At(idx)
		return data, err == nil
	}
	return []byte{}, false
}

func (this *MultiprocessHandler) unknow(caller, callee evmcommon.Address, input []byte) ([]byte, bool) {
	this.api.AddLog("Unhandled function call in cumulative handler router", hex.EncodeToString(input))
	return []byte{}, false
}

func (this *MultiprocessHandler) run(caller, callee evmcommon.Address, input []byte) ([]byte, bool) {
	if threads, err := abi.DecodeTo(input, 1, uint64(0), 1, 32); err == nil {
		return []byte{}, this.jobManager.Run(uint8(common.Max(common.Min(threads, 1), math.MaxUint8)))
	}

	return []byte{}, this.jobManager.Run(1)
}

func (this *MultiprocessHandler) delJob(caller, callee evmcommon.Address, input []byte) ([]byte, bool) {
	if idx, err := abi.DecodeTo(input, 1, uint64(0), 1, 32); err != nil {
		if _, err := this.jobManager.At(idx); err != nil {
			return codec.String(err.Error()).ToBytes(), err == nil
		}
	}

	id := this.api.GenUUID()
	return id, true
}

func (this *MultiprocessHandler) addJob(caller, callee evmcommon.Address, input []byte) ([]byte, bool) {
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

	jobID := this.jobManager.Add(calleeAddr, funCall)

	if buffer, err := abi.Encode(uint64(jobID)); err != nil {
		return []byte(err.Error()), false
	} else {
		return buffer, true
	}
}
