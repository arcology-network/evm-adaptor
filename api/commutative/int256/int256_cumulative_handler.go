package int64

import (
	"encoding/hex"

	"github.com/arcology-network/common-lib/codec"
	"github.com/arcology-network/common-lib/types"
	"github.com/holiman/uint256"

	"github.com/arcology-network/concurrenturl/commutative"
	evmcommon "github.com/arcology-network/evm/common"
	abi "github.com/arcology-network/vm-adaptor/abi"
	interfaces "github.com/arcology-network/vm-adaptor/interfaces"

	apicommon "github.com/arcology-network/vm-adaptor/api/common"
)

// APIs under the concurrency namespace
type Int256CumulativeHandlers struct {
	api       interfaces.ApiRouter
	connector *apicommon.CcurlConnector
}

func NewInt256CumulativeHandlers(api interfaces.ApiRouter) *Int256CumulativeHandlers {
	return &Int256CumulativeHandlers{
		api:       api,
		connector: apicommon.NewCCurlConnector("/containers/", api, api.Ccurl()),
	}
}

func (this *Int256CumulativeHandlers) Address() [20]byte {
	return [20]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0x86}
}

func (this *Int256CumulativeHandlers) Call(caller, callee evmcommon.Address, input []byte, origin evmcommon.Address, nonce uint64) ([]byte, bool) {
	signature := [4]byte{}
	copy(signature[:], input)

	switch signature {
	case [4]byte{0x90, 0x54, 0xce, 0x5f}:
		return this.new(caller, input[4:])

	case [4]byte{0x6d, 0x4c, 0xe6, 0x3c}:
		return this.get(caller, input[4:])

	case [4]byte{0xa4, 0xc6, 0xa7, 0x68}:
		return this.add(caller, input[4:])

	case [4]byte{0xc8, 0xda, 0xaa, 0xab}:
		return this.sub(caller, input[4:])
	}
	return this.Unknow(caller, input)
}

func (this *Int256CumulativeHandlers) Unknow(caller evmcommon.Address, input []byte) ([]byte, bool) {
	this.api.AddLog("Unhandled function call in cumulative handler router", hex.EncodeToString(input))
	return []byte{}, false
}

func (this *Int256CumulativeHandlers) new(caller evmcommon.Address, input []byte) ([]byte, bool) {
	id := this.api.GenCcObjID()
	if !this.connector.New(types.Address(codec.Bytes20(caller).Hex()), hex.EncodeToString(id)) { // A new container
		return []byte{}, false
	}

	txHash := this.api.TxHash()
	path := this.connector.Key(types.Address(codec.Bytes20(caller).Hex()), hex.EncodeToString(id))

	key := path +
		hex.EncodeToString(txHash[:8]) + "-" + // Tx hash to avoid conflict
		string(this.api.GenCcElemUID()) // Element ID

	// val, valErr := abi.Decode(input, 0, &uint256.Int{}, 1, 32)
	min, minErr := abi.Decode(input, 0, &uint256.Int{}, 1, 32)
	max, maxErr := abi.Decode(input, 1, &uint256.Int{}, 1, 32)
	if minErr != nil || maxErr != nil {
		return []byte{}, false
	}

	newU256 := commutative.NewU256(min.(*uint256.Int), max.(*uint256.Int))
	if _, err := this.api.Ccurl().Write(this.api.TxIndex(), key, newU256); err != nil {
		return []byte{}, false
	}
	return id, true
}

func (this *Int256CumulativeHandlers) get(caller evmcommon.Address, input []byte) ([]byte, bool) {
	path, err := this.buildPath(caller, input) // Build container path
	if len(path) == 0 || err != nil {
		return []byte{}, false
	}

	if value, _, err := this.api.Ccurl().ReadAt(this.api.TxIndex(), path, 0); value == nil || err != nil {
		return []byte{}, false
	} else {

		updated := value.(*uint256.Int)
		if encoded, err := abi.Encode(updated); err == nil { // Encode the result
			return encoded, true
		}
	}
	return []byte{}, false
}

func (this *Int256CumulativeHandlers) add(caller evmcommon.Address, input []byte) ([]byte, bool) {
	path, err := this.buildPath(caller, input) // Build container path
	if len(path) == 0 || err != nil {
		return []byte{}, false
	}

	delta, err := abi.Decode(input, 1, &uint256.Int{}, 1, 32)
	if err != nil {
		return []byte{}, false
	}

	value := commutative.NewU256Delta(delta.(*uint256.Int), true)
	_, err = this.api.Ccurl().WriteAt(this.api.TxIndex(), path, 0, value)
	return []byte{}, err == nil
}

func (this *Int256CumulativeHandlers) sub(caller evmcommon.Address, input []byte) ([]byte, bool) {
	path, err := this.buildPath(caller, input) // Build container path
	if len(path) == 0 || err != nil {
		return []byte{}, false
	}

	delta, err := abi.Decode(input, 1, &uint256.Int{}, 1, 32)
	if err != nil {
		return []byte{}, false
	}

	value := commutative.NewU256Delta(delta.(*uint256.Int), false)
	_, err = this.api.Ccurl().WriteAt(this.api.TxIndex(), path, 0, value)
	return []byte{}, err == nil
}

func (this *Int256CumulativeHandlers) set(caller evmcommon.Address, input []byte) ([]byte, bool) {
	path, err := this.buildPath(caller, input) // Build container path
	if len(path) == 0 || err != nil {
		return []byte{}, false
	}

	delta, err := abi.DecodeTo(input, 1, &uint256.Int{}, 1, 32)
	if err != nil {
		return []byte{}, false
	}

	sign, err := abi.DecodeTo(input, 1, bool(true), 1, 32)
	if err != nil {
		return []byte{}, false
	}

	value := commutative.NewU256Delta(delta, sign)
	_, err = this.api.Ccurl().WriteAt(this.api.TxIndex(), path, 0, value)
	return []byte{}, err == nil
}

// Build the container path
func (this *Int256CumulativeHandlers) buildPath(caller evmcommon.Address, input []byte) (string, error) {
	id, err := abi.Decode(input, 0, []byte{}, 2, 32) // max 32 bytes
	if err != nil {
		return "", nil
	} // container ID
	return this.connector.Key(types.Address(codec.Bytes20(caller).Hex()), hex.EncodeToString(id.([]byte))), nil // unique ID
}
