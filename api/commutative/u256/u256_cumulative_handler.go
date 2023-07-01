package u256

import (
	"encoding/hex"

	"github.com/arcology-network/common-lib/codec"
	"github.com/arcology-network/common-lib/types"
	"github.com/holiman/uint256"

	"github.com/arcology-network/concurrenturl/commutative"
	ccinterfaces "github.com/arcology-network/concurrenturl/interfaces"
	evmcommon "github.com/arcology-network/evm/common"
	abi "github.com/arcology-network/vm-adaptor/abi"
	apicommon "github.com/arcology-network/vm-adaptor/api/common"
	"github.com/arcology-network/vm-adaptor/common"
	eucommon "github.com/arcology-network/vm-adaptor/common"
)

// APIs under the concurrency namespace
type U256CumulativeHandlers struct {
	api       eucommon.EthApiRouter
	connector *apicommon.CcurlConnector
}

func NewU256CumulativeHandlers(api eucommon.EthApiRouter) *U256CumulativeHandlers {
	return &U256CumulativeHandlers{
		api:       api,
		connector: apicommon.NewCCurlConnector("/containers/", api, api.Ccurl()),
	}
}

func (this *U256CumulativeHandlers) Address() [20]byte {
	return common.CUMULATIVE_U256_HANDLER
}

func (this *U256CumulativeHandlers) Call(caller, callee [20]byte, input []byte, origin [20]byte, nonce uint64) ([]byte, bool, int64) {
	signature := [4]byte{}
	copy(signature[:], input)

	switch signature {
	case [4]byte{0x1c, 0x64, 0x49, 0x9c}:
		return this.new(caller, input[4:])

	case [4]byte{0xee, 0xb8, 0xa8, 0xd3}:
		return this.peek(caller, input[4:])

	case [4]byte{0x6d, 0x4c, 0xe6, 0x3c}:
		return this.get(caller, input[4:])

	case [4]byte{0x65, 0x83, 0x52, 0x0e}:
		return this.min(caller, input[4:])

	case [4]byte{0xaa, 0x11, 0xc9, 0x26}:
		return this.max(caller, input[4:])

	case [4]byte{0xaf, 0xc9, 0xfc, 0x46}:
		return this.add(caller, input[4:])

	case [4]byte{0xd8, 0x94, 0x8b, 0x09}:
		return this.sub(caller, input[4:])
	}
	return []byte{}, false, 0
}

func (this *U256CumulativeHandlers) new(caller evmcommon.Address, input []byte) ([]byte, bool, int64) {
	id := this.api.GenCcObjID()
	if !this.connector.New(types.Address(codec.Bytes20(caller).Hex()), hex.EncodeToString(id)) { // A new container
		return []byte{}, false, 0
	}

	path := this.connector.Key(types.Address(codec.Bytes20(caller).Hex()), hex.EncodeToString(id))
	key := path + string(this.api.GenCcElemUID()) // Element ID

	// val, valErr := abi.Decode(input, 0, &uint256.Int{}, 1, 32)
	min, minErr := abi.Decode(input, 0, &uint256.Int{}, 1, 32)
	max, maxErr := abi.Decode(input, 1, &uint256.Int{}, 1, 32)
	if minErr != nil || maxErr != nil {
		return []byte{}, false, 0
	}

	newU256 := commutative.NewU256(min.(*uint256.Int), max.(*uint256.Int))
	if _, err := this.api.Ccurl().Write(this.api.TxIndex(), key, newU256, true); err != nil {
		return []byte{}, false, 0
	}
	return id, true, 0
}

func (this *U256CumulativeHandlers) get(caller evmcommon.Address, input []byte) ([]byte, bool, int64) {
	path, err := this.buildPath(caller, input) // Build container path
	if len(path) == 0 || err != nil {
		return []byte{}, false, 0
	}

	if value, _, err := this.api.Ccurl().ReadAt(this.api.TxIndex(), path, 0); value == nil || err != nil {
		return []byte{}, false, 0
	} else {
		updated := value.(*uint256.Int)
		if encoded, err := abi.Encode(updated); err == nil { // Encode the result
			return encoded, true, 0
		}
	}
	return []byte{}, false, 0
}

func (this *U256CumulativeHandlers) peek(caller evmcommon.Address, input []byte) ([]byte, bool, int64) { // Get the initial value
	path, err := this.buildPath(caller, input) // Build container path
	if len(path) == 0 || err != nil {
		return []byte{}, false, 0
	}

	// Peek the initial value
	value, _, err := this.api.Ccurl().DoAt(this.api.TxIndex(), path, 0, func(v interface{}) interface{} {
		return []interface{}{uint32(0), uint32(0), uint32(0), v.(ccinterfaces.Univalue).Value()}
	})

	if value != nil && err == nil {
		if initv := value.(*commutative.U256).Value().(*codec.Uint256); initv != nil {
			if encoded, err := abi.Encode((*uint256.Int)(initv)); err == nil { // Encode the result
				return encoded, true, 0
			}
		}
	}
	return []byte{}, false, 0
}

func (this *U256CumulativeHandlers) add(caller evmcommon.Address, input []byte) ([]byte, bool, int64) {
	return this.set(caller, input, true)
}

func (this *U256CumulativeHandlers) sub(caller evmcommon.Address, input []byte) ([]byte, bool, int64) {
	return this.set(caller, input, false)
}

func (this *U256CumulativeHandlers) set(caller evmcommon.Address, input []byte, isPositive bool) ([]byte, bool, int64) {
	path, err := this.buildPath(caller, input) // Build container path
	if len(path) == 0 || err != nil {
		return []byte{}, false, 0
	}

	delta, err := abi.Decode(input, 1, &uint256.Int{}, 1, 32)
	if err != nil {
		return []byte{}, false, 0
	}

	value := commutative.NewU256Delta(delta.(*uint256.Int), isPositive)

	_, v := this.api.Ccurl().WriteAt(this.api.TxIndex(), path, 0, value, true)
	return []byte{}, v == nil, 0
}

func (this *U256CumulativeHandlers) min(caller evmcommon.Address, input []byte) ([]byte, bool, int64) {
	path, err := this.buildPath(caller, input) // Build container path
	if len(path) == 0 || err != nil {
		return []byte{}, false, 0
	}

	value, _, err := this.api.Ccurl().DoAt(this.api.TxIndex(), path, 0, func(v interface{}) interface{} {
		return []interface{}{uint32(1), uint32(0), uint32(0), v.(ccinterfaces.Univalue).Value()}
	})

	if value != nil && err == nil {
		minv := value.(*commutative.U256).Min().(*codec.Uint256)
		if encoded, err := abi.Encode((*uint256.Int)(minv)); err == nil { // Encode the result
			return encoded, true, 0
		}
	}
	return []byte{}, false, 0
}

func (this *U256CumulativeHandlers) max(caller evmcommon.Address, input []byte) ([]byte, bool, int64) {
	path, err := this.buildPath(caller, input) // Build container path
	if len(path) == 0 || err != nil {
		return []byte{}, false, 0
	}

	value, _, err := this.api.Ccurl().DoAt(this.api.TxIndex(), path, 0, func(v interface{}) interface{} {
		return []interface{}{uint32(1), uint32(0), uint32(0), v.(ccinterfaces.Univalue).Value()}
	})

	if value != nil && err == nil {
		minv := value.(*commutative.U256).Max().(*codec.Uint256)
		if encoded, err := abi.Encode((*uint256.Int)(minv)); err == nil { // Encode the result
			return encoded, true, 0
		}
	}
	return []byte{}, false, 0
}

// Build the container path
func (this *U256CumulativeHandlers) buildPath(caller evmcommon.Address, input []byte) (string, error) {
	id, err := abi.Decode(input, 0, []byte{}, 2, 32) // max 32 bytes
	if err != nil {
		return "", nil
	} // container ID
	return this.connector.Key(types.Address(codec.Bytes20(caller).Hex()), hex.EncodeToString(id.([]byte))), nil // unique ID
}
