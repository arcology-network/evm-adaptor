package eu

import (
	"math/big"

	"github.com/arcology-network/evm/core/types"
	"github.com/arcology-network/evm/core/vm"
)

func NewEVMTxContext(msg types.Message) vm.TxContext {
	return vm.TxContext{
		Origin:   msg.From(),
		GasPrice: new(big.Int).Set(msg.GasPrice()),
	}
}
