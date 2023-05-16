package u256

import (
	"github.com/arcology-network/common-lib/common"
	"github.com/arcology-network/concurrenturl/commutative"
	abi "github.com/arcology-network/vm-adaptor/abi"
	"github.com/holiman/uint256"
)

// APIs under the concurrency namespace
type U256Codec struct{}

func (this *U256Codec) new(input []byte) (interface{}, error) {
	min, minErr := abi.DecodeTo(input, 0, &uint256.Int{}, 1, 32)
	max, maxErr := abi.DecodeTo(input, 1, &uint256.Int{}, 1, 32)
	if minErr != nil || maxErr != nil {
		return nil, common.IfThen(minErr != nil, minErr, maxErr)
	}
	return commutative.NewU256(min, max), nil
}

func (this *U256Codec) delta(input []byte) (interface{}, error) {
	delta, err := abi.DecodeTo(input, 1, &uint256.Int{}, 1, 32)
	if err != nil {
		return nil, err
	}
	return commutative.NewU256Delta(delta, false), nil
}

func (this *U256Codec) encode(value interface{}) ([]byte, error) {
	updated := value.(*uint256.Int)
	if encoded, err := abi.Encode(updated); err == nil { // Encode the result
		return encoded, err
	}
	return []byte{}, nil
}
