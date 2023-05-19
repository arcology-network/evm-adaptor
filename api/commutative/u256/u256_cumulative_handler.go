package u256

import (
	"encoding/hex"
	"strconv"

	"github.com/arcology-network/common-lib/codec"
	"github.com/arcology-network/common-lib/types"
	"github.com/holiman/uint256"

	"github.com/arcology-network/concurrenturl/commutative"
	evmcommon "github.com/arcology-network/evm/common"
	abi "github.com/arcology-network/vm-adaptor/abi"

	apicommon "github.com/arcology-network/vm-adaptor/api/common"
	eucommon "github.com/arcology-network/vm-adaptor/common"
)

// APIs under the concurrency namespace
type U256CumulativeHandlers struct {
	api       eucommon.ConcurrentApiRouterInterface
	connector *apicommon.CcurlConnector
}

func NewU256CumulativeHandlers(api eucommon.ConcurrentApiRouterInterface) *U256CumulativeHandlers {
	return &U256CumulativeHandlers{
		api:       api,
		connector: apicommon.NewCCurlConnector("/containers/", api, api.Ccurl()),
	}
}

func (this *U256CumulativeHandlers) Address() [20]byte {
	return [20]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0x85}
}

func (this *U256CumulativeHandlers) Call(caller, callee evmcommon.Address, input []byte, origin evmcommon.Address, nonce uint64) ([]byte, bool) {
	signature := [4]byte{}
	copy(signature[:], input)

	switch signature {
	case [4]byte{0x1c, 0x64, 0x49, 0x9c}:
		return this.new(caller, input[4:])

	case [4]byte{0x6d, 0x4c, 0xe6, 0x3c}:
		return this.get(caller, input[4:])

	// case [4]byte{0xc9, 0xef, 0xba, 0xb9}:
	// 	return this.set(caller, input[4:])

	case [4]byte{0xaf, 0xc9, 0xfc, 0x46}:
		return this.add(caller, input[4:])

	case [4]byte{0xd8, 0x94, 0x8b, 0x09}:
		return this.sub(caller, input[4:])
	}
	return this.Unknow(caller, input)
}

func (this *U256CumulativeHandlers) Unknow(caller evmcommon.Address, input []byte) ([]byte, bool) {
	this.api.AddLog("Unhandled function call in cumulative handler router", hex.EncodeToString(input))
	return []byte{}, false
}

func (this *U256CumulativeHandlers) new(caller evmcommon.Address, input []byte) ([]byte, bool) {
	id := this.api.GenCtrnUID()
	if !this.connector.New(types.Address(codec.Bytes20(caller).Hex()), hex.EncodeToString(id)) { // A new container
		return []byte{}, false
	}

	txHash := this.api.TxHash()
	path := this.connector.Key(types.Address(codec.Bytes20(caller).Hex()), hex.EncodeToString(id))

	key := path +
		hex.EncodeToString(txHash[:8]) + "-" + // Tx hash to avoid conflict
		strconv.Itoa(int(this.api.GenElemUID())) // Element ID

	// val, valErr := abi.Decode(input, 0, &uint256.Int{}, 1, 32)
	min, minErr := abi.Decode(input, 0, &uint256.Int{}, 1, 32)
	max, maxErr := abi.Decode(input, 1, &uint256.Int{}, 1, 32)
	if minErr != nil || maxErr != nil {
		return []byte{}, false
	}

	newU256 := commutative.NewU256(min.(*uint256.Int), max.(*uint256.Int))
	if err := this.api.Ccurl().Write(this.api.TxIndex(), key, newU256); err != nil {
		return []byte{}, false
	}
	return id, true
}

func (this *U256CumulativeHandlers) get(caller evmcommon.Address, input []byte) ([]byte, bool) {
	path, err := this.buildPath(caller, input) // Build container path
	if len(path) == 0 || err != nil {
		return []byte{}, false
	}

	if value, err := this.api.Ccurl().ReadAt(this.api.TxIndex(), path, 0); value == nil || err != nil {
		return []byte{}, false
	} else {

		updated := value.(*uint256.Int)
		if encoded, err := abi.Encode(updated); err == nil { // Encode the result
			return encoded, true
		}
	}
	return []byte{}, false
}

func (this *U256CumulativeHandlers) add(caller evmcommon.Address, input []byte) ([]byte, bool) {
	return this.set(caller, input, true)
}

func (this *U256CumulativeHandlers) sub(caller evmcommon.Address, input []byte) ([]byte, bool) {
	return this.set(caller, input, false)
}

func (this *U256CumulativeHandlers) set(caller evmcommon.Address, input []byte, isPositive bool) ([]byte, bool) {
	path, err := this.buildPath(caller, input) // Build container path
	if len(path) == 0 || err != nil {
		return []byte{}, false
	}

	delta, err := abi.Decode(input, 1, &uint256.Int{}, 1, 32)
	if err != nil {
		return []byte{}, false
	}

	value := commutative.NewU256Delta(delta.(*uint256.Int), isPositive)
	return []byte{}, this.api.Ccurl().WriteAt(this.api.TxIndex(), path, 0, value) == nil
}

// Build the container path
func (this *U256CumulativeHandlers) buildPath(caller evmcommon.Address, input []byte) (string, error) {
	id, err := abi.Decode(input, 0, []byte{}, 2, 32) // max 32 bytes
	if err != nil {
		return "", nil
	} // container ID
	return this.connector.Key(types.Address(codec.Bytes20(caller).Hex()), hex.EncodeToString(id.([]byte))), nil // unique ID
}
