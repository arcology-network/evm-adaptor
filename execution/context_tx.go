package execution

import (
	"math/big"

	"github.com/arcology-network/evm/core"
	"github.com/arcology-network/evm/core/vm"
)

func NewEVMTxContext(msg core.Message) vm.TxContext {
	return vm.TxContext{
		Origin:   msg.From,
		GasPrice: new(big.Int).Set(msg.GasPrice),
	}
}
