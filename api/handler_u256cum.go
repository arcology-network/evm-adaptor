package api

import (
	"github.com/arcology-network/common-lib/codec"
	"github.com/arcology-network/common-lib/types"
	"github.com/holiman/uint256"

	"github.com/arcology-network/concurrenturl/commutative"
	ccinterfaces "github.com/arcology-network/concurrenturl/interfaces"
	abi "github.com/arcology-network/vm-adaptor/abi"
	intf "github.com/arcology-network/vm-adaptor/interface"
	evmcommon "github.com/ethereum/go-ethereum/common"

	"github.com/arcology-network/vm-adaptor/common"
	adaptorcommon "github.com/arcology-network/vm-adaptor/common"
)

// APIs under the concurrency namespace
type U256CumHandlers struct {
	api       intf.EthApiRouter
	connector *adaptorcommon.CcurlConnector
}

func NewU256CumulativeHandlers(api intf.EthApiRouter) *U256CumHandlers {
	return &U256CumHandlers{
		api:       api,
		connector: adaptorcommon.NewCCurlConnector("/container", api, api.Ccurl()),
	}
}

func (this *U256CumHandlers) Address() [20]byte {
	return common.CUMULATIVE_U256_HANDLER
}

// performed the changes on Delta only

func (this *U256CumHandlers) Call(caller, callee [20]byte, input []byte, origin [20]byte, nonce uint64) ([]byte, bool, int64) {
	signature := [4]byte{}
	copy(signature[:], input)

	switch signature {
	case [4]byte{0x1c, 0x64, 0x49, 0x9c}:
		return this.new(caller, input[4:])

	case [4]byte{0x59, 0xe0, 0x2d, 0xd7}: // 59 e0 2d d7
		return this.peek(caller, input[4:])

	case [4]byte{0x6d, 0x4c, 0xe6, 0x3c}:
		return this.get(caller, input[4:])

	case [4]byte{0xf8, 0x89, 0x79, 0x45}: // f8 89 79 45
		return this.min(caller, input[4:])

	case [4]byte{0x6a, 0xc5, 0xdb, 0x19}:
		return this.max(caller, input[4:])

	case [4]byte{0x10, 0x03, 0xe2, 0xd2}: // 10 03 e2 d2
		return this.add(caller, input[4:])

	case [4]byte{0x27, 0xee, 0x58, 0xa6}:
		return this.sub(caller, input[4:]) //27 ee 58 a6
	}
	return []byte{}, false, 0
}

func (this *U256CumHandlers) new(caller evmcommon.Address, input []byte) ([]byte, bool, int64) {
	txIndex := this.api.GetEU().(intf.EU).ID()
	if !this.connector.New(txIndex, types.Address(codec.Bytes20(caller).Hex())) { // A new container
		return []byte{}, false, 0
	}

	path := this.connector.Key(caller)
	key := path + string(this.api.ElementUID()) // Element ID

	// val, valErr := abi.Decode(input, 0, &uint256.Int{}, 1, 32)
	min, minErr := abi.Decode(input, 0, &uint256.Int{}, 1, 32)
	max, maxErr := abi.Decode(input, 1, &uint256.Int{}, 1, 32)
	if minErr != nil || maxErr != nil {
		return []byte{}, false, 0
	}

	newU256 := commutative.NewBoundedU256(min.(*uint256.Int), max.(*uint256.Int))
	if _, err := this.api.Ccurl().Write(txIndex, key, newU256); err != nil {
		return []byte{}, false, 0
	}
	return []byte{}, true, 0
}

func (this *U256CumHandlers) get(caller evmcommon.Address, input []byte) ([]byte, bool, int64) {
	path := this.connector.Key(caller) // Build container path
	if len(path) == 0 {
		return []byte{}, false, 0
	}

	if value, _, err := this.api.Ccurl().ReadAt(this.api.GetEU().(intf.EU).ID(), path, 0, new(commutative.U256)); value == nil || err != nil {
		return []byte{}, false, 0
	} else {
		updated := value.(uint256.Int)
		if encoded, err := abi.Encode(updated); err == nil { // Encode the result
			return encoded, true, 0
		}
	}
	return []byte{}, false, 0
}

func (this *U256CumHandlers) peek(caller evmcommon.Address, input []byte) ([]byte, bool, int64) { // Get the initial value
	path := this.connector.Key(caller) // Build container path
	if len(path) == 0 {
		return []byte{}, false, 0
	}

	// Peek the initial value
	value, _, err := this.api.Ccurl().DoAt(this.api.GetEU().(intf.EU).ID(), path, 0, func(v interface{}) (uint32, uint32, uint32, interface{}) {
		return uint32(0), uint32(0), uint32(0), v.(ccinterfaces.Univalue).Value()
	}, new(commutative.U256))

	if value != nil && err == nil {
		initv := value.(*commutative.U256).Value().(uint256.Int)
		// if ; initv != nil {
		if encoded, err := abi.Encode((*uint256.Int)(&initv)); err == nil { // Encode the result
			return encoded, true, 0
		}
		// }
	}
	return []byte{}, false, 0
}

func (this *U256CumHandlers) add(caller evmcommon.Address, input []byte) ([]byte, bool, int64) {
	return this.set(caller, input, true)
}

func (this *U256CumHandlers) sub(caller evmcommon.Address, input []byte) ([]byte, bool, int64) {
	return this.set(caller, input, false)
}

func (this *U256CumHandlers) set(caller evmcommon.Address, input []byte, isPositive bool) ([]byte, bool, int64) {
	path := this.connector.Key(caller) // Build container path
	if len(path) == 0 {
		return []byte{}, false, 0
	}

	delta, err := abi.Decode(input, 0, &uint256.Int{}, 1, 32)
	if err != nil {
		return []byte{}, false, 0
	}

	value := commutative.NewU256Delta(delta.(*uint256.Int), isPositive)

	_, err = this.api.Ccurl().WriteAt(this.api.GetEU().(intf.EU).ID(), path, 0, value)
	return []byte{}, err == nil, 0
}

func (this *U256CumHandlers) min(caller evmcommon.Address, input []byte) ([]byte, bool, int64) {
	path := this.connector.Key(caller) // Build container path
	if len(path) == 0 {
		return []byte{}, false, 0
	}

	value, _, err := this.api.Ccurl().DoAt(this.api.GetEU().(intf.EU).ID(), path, 0, func(v interface{}) (uint32, uint32, uint32, interface{}) {
		return uint32(1), uint32(0), uint32(0), v.(ccinterfaces.Univalue).Value()
	}, new(commutative.U256))

	if value != nil && err == nil {
		minv := value.(*commutative.U256).Min().(uint256.Int)
		if encoded, err := abi.Encode((*uint256.Int)(&minv)); err == nil { // Encode the result
			return encoded, true, 0
		}
	}
	return []byte{}, false, 0
}

func (this *U256CumHandlers) max(caller evmcommon.Address, input []byte) ([]byte, bool, int64) {
	path := this.connector.Key(caller) // Build container path
	if len(path) == 0 {
		return []byte{}, false, 0
	}

	value, _, err := this.api.Ccurl().DoAt(this.api.GetEU().(intf.EU).ID(), path, 0, func(v interface{}) (uint32, uint32, uint32, interface{}) {
		return uint32(1), uint32(0), uint32(0), v.(ccinterfaces.Univalue).Value()
	}, new(commutative.U256))

	if value != nil && err == nil {
		minv := value.(*commutative.U256).Max().(uint256.Int)
		if encoded, err := abi.Encode((*uint256.Int)(&minv)); err == nil { // Encode the result
			return encoded, true, 0
		}
	}
	return []byte{}, false, 0
}
