package tests

import (
	"math"
	"math/big"

	"github.com/arcology/evm/common"
	"github.com/arcology/evm/consensus"
	"github.com/arcology/evm/core/types"
	"github.com/arcology/evm/core/vm"
	"github.com/arcology/evm/params"
	adaptor "github.com/arcology/vm-adaptor/evm"
)

var (
	coinbase = common.BytesToAddress([]byte("coinbase"))
	owner    = common.BytesToAddress([]byte("owner"))
	user1    = common.BytesToAddress([]byte("user1"))
	user2    = common.BytesToAddress([]byte("user2"))
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

func MainConfig() *adaptor.Config {
	vmConfig := vm.Config{}
	cfg := &adaptor.Config{
		ChainConfig: params.MainnetChainConfig,
		VMConfig:    &vmConfig,
		BlockNumber: new(big.Int).SetUint64(10000000),
		ParentHash:  common.Hash{},
		Time:        new(big.Int).SetUint64(10000000),
		Coinbase:    &coinbase,
		GasLimit:    math.MaxUint64,
		Difficulty:  new(big.Int).SetUint64(10000000),
	}
	cfg.Chain = new(fakeChain)
	return cfg
}
