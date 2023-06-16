package eu

import (
	"math"
	"math/big"

	"github.com/arcology-network/evm/common"
	evmcommon "github.com/arcology-network/evm/common"
	"github.com/arcology-network/evm/consensus"
	"github.com/arcology-network/evm/core/types"
	"github.com/arcology-network/evm/core/vm"
	"github.com/arcology-network/evm/params"
	eucommon "github.com/arcology-network/vm-adaptor/common"
)

// DummyChain implements the ChainContext interface.
type DummyChain struct{}

func (chain *DummyChain) GetHeader(evmcommon.Hash, uint64) *types.Header { return &types.Header{} }
func (chain *DummyChain) Engine() consensus.Engine                       { return nil }

// Config contains all the static settings used in Schedule.
type Config struct {
	ChainConfig *params.ChainConfig
	VMConfig    *vm.Config
	BlockNumber *big.Int    // types.Header.Number
	ParentHash  common.Hash // types.Header.ParentHash
	Time        *big.Int    // types.Header.Time
	Chain       eucommon.ChainContext
	Coinbase    *evmcommon.Address
	GasLimit    uint64   // types.Header.GasLimit
	Difficulty  *big.Int // types.Header.Difficulty
}

func (this *Config) SetCoinbase(coinbase evmcommon.Address) *Config {
	this.Coinbase = &coinbase
	return this
}

func NewConfig() *Config {
	cfg := &Config{
		ChainConfig: params.MainnetChainConfig,
		VMConfig:    &vm.Config{},
		BlockNumber: big.NewInt(0),
		ParentHash:  evmcommon.Hash{},
		Time:        big.NewInt(0),
		Coinbase:    &evmcommon.Address{},
		GasLimit:    math.MaxUint64,
		Difficulty:  big.NewInt(0),
	}
	cfg.Chain = new(DummyChain)
	return cfg
}
