package tests

import (
	"math"
	"math/big"

	"github.com/arcology-network/evm/common"
	"github.com/arcology-network/evm/consensus"
	"github.com/arcology-network/evm/core/types"
	"github.com/arcology-network/evm/core/vm"
	"github.com/arcology-network/evm/params"
	eu "github.com/arcology-network/vm-adaptor"
)

var (
	Coinbase = common.BytesToAddress([]byte("coinbase"))
	Owner    = common.BytesToAddress([]byte("owner"))
	User1    = common.BytesToAddress([]byte("user1"))
	User2    = common.BytesToAddress([]byte("user2"))
)

// fakeChain implements the ChainContext interface.
type fakeChain struct {
}

func (chain *fakeChain) GetHeader(common.Hash, uint64) *types.Header {
	return &types.Header{}
}

func (chain *fakeChain) Engine() consensus.Engine {
	return nil
}

func MainConfig() *eu.Config {
	vmConfig := vm.Config{}
	cfg := &eu.Config{
		ChainConfig: params.MainnetChainConfig,
		VMConfig:    &vmConfig,
		BlockNumber: new(big.Int).SetUint64(10000000),
		ParentHash:  common.Hash{},
		Time:        new(big.Int).SetUint64(10000000),
		Coinbase:    &Coinbase,
		GasLimit:    math.MaxUint64,
		Difficulty:  new(big.Int).SetUint64(10000000),
	}
	cfg.Chain = new(fakeChain)
	return cfg
}
