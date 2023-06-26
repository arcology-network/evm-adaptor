package int64

import (
	"encoding/hex"

	"github.com/arcology-network/common-lib/codec"
	"github.com/arcology-network/common-lib/types"
	"github.com/holiman/uint256"

	"github.com/arcology-network/concurrenturl/commutative"
	evmcommon "github.com/arcology-network/evm/common"
	abi "github.com/arcology-network/vm-adaptor/abi"
	"github.com/arcology-network/vm-adaptor/common"
	eucommon "github.com/arcology-network/vm-adaptor/common"

	apicommon "github.com/arcology-network/vm-adaptor/api/common"
)

// APIs under the concurrency namespace
type Int256CumulativeHandlers struct {
	api       eucommon.EthApiRouter
	connector *apicommon.CcurlConnector
}

func NewInt256CumulativeHandlers(api eucommon.EthApiRouter) *Int256CumulativeHandlers {
	return &Int256CumulativeHandlers{
		api:       api,
		connector: apicommon.NewCCurlConnector("/containers/", api, api.Ccurl()),
	}
}

func (this *Int256CumulativeHandlers) Address() [20]byte {
	return common.CUMULATIVE_I256_HANDLER
}

func (this *Int256CumulativeHandlers) Call(caller, callee [20]byte, input []byte, origin [20]byte, nonce uint64) ([]byte, bool, int64) {
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

func (this *Int256CumulativeHandlers) Unknow(caller evmcommon.Address, input []byte) ([]byte, bool, int64) {
	this.api.AddLog("Unhandled function call in cumulative handler router", hex.EncodeToString(input))
	return []byte{}, false, 0
}

func (this *Int256CumulativeHandlers) new(caller evmcommon.Address, input []byte) ([]byte, bool, int64) {
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
	if _, err := this.api.Ccurl().Write(this.api.TxIndex(), key, newU256); err != nil {
		return []byte{}, false, 0
	}
	return id, true, 0
}

func (this *Int256CumulativeHandlers) get(caller evmcommon.Address, input []byte) ([]byte, bool, int64) {
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

func (this *Int256CumulativeHandlers) add(caller evmcommon.Address, input []byte) ([]byte, bool, int64) {
	path, err := this.buildPath(caller, input) // Build container path
	if len(path) == 0 || err != nil {
		return []byte{}, false, 0
	}

	delta, err := abi.Decode(input, 1, &uint256.Int{}, 1, 32)
	if err != nil {
		return []byte{}, false, 0
	}

	value := commutative.NewU256Delta(delta.(*uint256.Int), true)
	_, err = this.api.Ccurl().WriteAt(this.api.TxIndex(), path, 0, value)
	return []byte{}, err == nil, 0
}

func (this *Int256CumulativeHandlers) sub(caller evmcommon.Address, input []byte) ([]byte, bool, int64) {
	path, err := this.buildPath(caller, input) // Build container path
	if len(path) == 0 || err != nil {
		return []byte{}, false, 0
	}

	delta, err := abi.Decode(input, 1, &uint256.Int{}, 1, 32)
	if err != nil {
		return []byte{}, false, 0
	}

	value := commutative.NewU256Delta(delta.(*uint256.Int), false)
	_, err = this.api.Ccurl().WriteAt(this.api.TxIndex(), path, 0, value)
	return []byte{}, err == nil, 0
}

func (this *Int256CumulativeHandlers) set(caller evmcommon.Address, input []byte) ([]byte, bool, int64) {
	path, err := this.buildPath(caller, input) // Build container path
	if len(path) == 0 || err != nil {
		return []byte{}, false, 0
	}

	delta, err := abi.DecodeTo(input, 1, &uint256.Int{}, 1, 32)
	if err != nil {
		return []byte{}, false, 0
	}

	sign, err := abi.DecodeTo(input, 1, bool(true), 1, 32)
	if err != nil {
		return []byte{}, false, 0
	}

	value := commutative.NewU256Delta(delta, sign)
	_, err = this.api.Ccurl().WriteAt(this.api.TxIndex(), path, 0, value)
	return []byte{}, err == nil, 0
}

// Build the container path
func (this *Int256CumulativeHandlers) buildPath(caller evmcommon.Address, input []byte) (string, error) {
	id, err := abi.Decode(input, 0, []byte{}, 2, 32) // max 32 bytes
	if err != nil {
		return "", nil
	} // container ID
	return this.connector.Key(types.Address(codec.Bytes20(caller).Hex()), hex.EncodeToString(id.([]byte))), nil // unique ID
}
