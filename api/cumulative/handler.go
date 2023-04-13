package cumulative

import (
	"encoding/hex"
	"strconv"

	"github.com/arcology-network/common-lib/codec"
	"github.com/arcology-network/common-lib/types"
	"github.com/holiman/uint256"

	ccurlcommon "github.com/arcology-network/concurrenturl/v2/common"
	"github.com/arcology-network/concurrenturl/v2/type/commutative"
	evmcommon "github.com/arcology-network/evm/common"
	abi "github.com/arcology-network/vm-adaptor/abi"

	apicommon "github.com/arcology-network/vm-adaptor/api/common"
)

// APIs under the concurrency namespace
type Cumulative struct {
	api       apicommon.ContextInfoInterface
	connector *apicommon.CcurlConnector
	path      string
}

func NewCumulative(api apicommon.ContextInfoInterface) *Cumulative {
	return &Cumulative{
		api:       api,
		connector: apicommon.NewCCurlConnector("/storage/containers/", api.TxHash(), api.TxIndex(), api.Ccurl()),
	}
}

func (this *Cumulative) Address() [20]byte {
	return [20]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0x85}
}

func (this *Cumulative) Call(caller evmcommon.Address, input []byte, origin evmcommon.Address, nonce uint64) ([]byte, bool) {
	signature := [4]byte{}
	copy(signature[:], input)

	switch signature {
	case [4]byte{0x26, 0x4a, 0x7f, 0x20}:
		return this.New(caller, input[4:])

	case [4]byte{0x1e, 0x81, 0x3e, 0x5d}:
		return this.Get(caller, input[4:])

	case [4]byte{0x76, 0x18, 0xc3, 0x17}:
		return this.Add(caller, input[4:])

	case [4]byte{0x57, 0xf1, 0xa8, 0x7e}:
		return this.Sub(caller, input[4:])
	}
	return this.Unknow(caller, input[4:])
}

func (this *Cumulative) Unknow(caller evmcommon.Address, input []byte) ([]byte, bool) {
	this.api.AddLog("Unhandled function call in cumulative handler router", hex.EncodeToString(input))
	return []byte{}, false
}

func (this *Cumulative) New(caller evmcommon.Address, input []byte) ([]byte, bool) {
	id := this.api.GenUUID()
	if !this.connector.New(types.Address(codec.Bytes20(caller).Hex()), hex.EncodeToString(id), 0) { // A new container
		return []byte{}, false
	}

	txHash := this.api.TxHash()
	path := this.connector.Key(types.Address(codec.Bytes20(caller).Hex()), hex.EncodeToString(id))

	key := path + // Root
		hex.EncodeToString(txHash[:8]) + "-" + // Tx hash to avoid conflict
		strconv.Itoa(int(ccurlcommon.CommutativeUint256)) + "-" + // value type
		strconv.Itoa(int(this.api.SUID())) // Element ID

	val, valErr := abi.Decode(input, 0, &uint256.Int{}, 1, 32)
	min, minErr := abi.Decode(input, 1, &uint256.Int{}, 1, 32)
	max, maxErr := abi.Decode(input, 2, &uint256.Int{}, 1, 32)

	if valErr != nil || minErr != nil || maxErr != nil {
		return []byte{}, false
	}

	newU256 := commutative.NewU256(val.(*uint256.Int), min.(*uint256.Int), max.(*uint256.Int))
	if err := this.api.Ccurl().Write(this.api.TxIndex(), key, newU256); err != nil {
		return []byte{}, false
	}
	return id, true
}

func (this *Cumulative) Get(caller evmcommon.Address, input []byte) ([]byte, bool) {
	path, err := this.buildPath(caller, input) // Build container path
	if len(path) == 0 || err != nil {
		return []byte{}, false
	}

	if value, err := this.api.Ccurl().ReadAt(this.api.TxIndex(), path, 0); value == nil || err != nil {
		return []byte{}, false
	} else {

		updated := value.(*commutative.U256).Value().(*uint256.Int)
		if encoded, err := abi.Encode(updated); err == nil { // Encode the result
			return encoded, true
		}
	}
	return []byte{}, false
}

func (this *Cumulative) Add(caller evmcommon.Address, input []byte) ([]byte, bool) {
	path, err := this.buildPath(caller, input) // Build container path
	if len(path) == 0 || err != nil {
		return []byte{}, false
	}

	delta, err := abi.Decode(input, 1, &uint256.Int{}, 1, 32)
	if err != nil {
		return []byte{}, false
	}

	value := commutative.NewU256Delta(delta.(*uint256.Int), true)
	return []byte{}, this.api.Ccurl().WriteAt(this.api.TxIndex(), path, 0, value) == nil
}

func (this *Cumulative) Sub(caller evmcommon.Address, input []byte) ([]byte, bool) {
	path, err := this.buildPath(caller, input) // Build container path
	if len(path) == 0 || err != nil {
		return []byte{}, false
	}

	delta, err := abi.Decode(input, 1, &uint256.Int{}, 1, 32)
	if err != nil {
		return []byte{}, false
	}

	value := commutative.NewU256Delta(delta.(*uint256.Int), false)
	return []byte{}, this.api.Ccurl().WriteAt(this.api.TxIndex(), path, 0, value) == nil
}

// Build the container path
func (this *Cumulative) buildPath(caller evmcommon.Address, input []byte) (string, error) {
	id, err := abi.Decode(input, 0, []byte{}, 2, 32) // max 32 bytes
	if err != nil {
		return "", nil
	} // container ID
	return this.connector.Key(types.Address(codec.Bytes20(caller).Hex()), hex.EncodeToString(id.([]byte))), nil // unique ID
}
