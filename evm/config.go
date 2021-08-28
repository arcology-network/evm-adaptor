package evm

import (
	"math/big"

	"github.com/arcology/evm/common"
	"github.com/arcology/evm/consensus"
	"github.com/arcology/evm/core/types"
	"github.com/arcology/evm/core/vm"
	"github.com/arcology/evm/params"
)

// ChainContext supports retrieving headers and consensus parameters from the
// current blockchain to be used during transaction processing.
type ChainContext interface {
	// Engine retrieves the chain's consensus engine.
	Engine() consensus.Engine

	// GetHeader returns the hash corresponding to their hash.
	GetHeader(common.Hash, uint64) *types.Header
}

// Config contains all the static settings used in Schedule.
type Config struct {
	ChainConfig *params.ChainConfig
	VMConfig    *vm.Config
	BlockNumber *big.Int    // types.Header.Number
	ParentHash  common.Hash // types.Header.ParentHash
	Time        *big.Int    // types.Header.Time
	Chain       ChainContext
	Coinbase    *common.Address
	GasLimit    uint64   // types.Header.GasLimit
	Difficulty  *big.Int // types.Header.Difficulty
}

func NewEVMBlockContextV2(cfg *Config) vm.BlockContext {
	return vm.BlockContext{
		CanTransfer: CanTransferV2,
		Transfer:    Transfer,
		GetHash:     GetHashFn(cfg.BlockNumber, cfg.ParentHash, cfg.Chain),
		Coinbase:    *cfg.Coinbase,
		BlockNumber: new(big.Int).Set(cfg.BlockNumber),
		Time:        new(big.Int).Set(cfg.Time),
		Difficulty:  new(big.Int).Set(cfg.Difficulty),
		GasLimit:    cfg.GasLimit,
	}
}

func NewEVMTxContext(msg types.Message) vm.TxContext {
	return vm.TxContext{
		Origin:   msg.From(),
		GasPrice: new(big.Int).Set(msg.GasPrice()),
	}
}

// GetHashFn returns a GetHashFunc which retrieves header hashes by number
func GetHashFn(blockNumber *big.Int, parentHash common.Hash, chain ChainContext) func(n uint64) common.Hash {
	// var cache map[uint64]common.Hash

	return func(n uint64) common.Hash {
		// If there's no hash cache yet, make one
		// if cache == nil {
		// 	cache = map[uint64]common.Hash{
		// 		blockNumber.Uint64() - 1: parentHash,
		// 	}
		// }
		// // Try to fulfill the request from the cache
		// if hash, ok := cache[n]; ok {
		// 	return hash
		// }
		// // Not cached, iterate the blocks and cache the hashes
		// for header := chain.GetHeader(parentHash, blockNumber.Uint64()-1); header != nil; header = chain.GetHeader(header.ParentHash, header.Number.Uint64()-1) {
		// 	cache[header.Number.Uint64()-1] = header.ParentHash
		// 	if n == header.Number.Uint64()-1 {
		// 		return header.ParentHash
		// 	}
		// }
		return common.Hash{}
	}
}

// CanTransfer checks whether there are enough funds in the address' account to make a transfer.
// This does not take the necessary gas in to account to make the transfer valid.

func CanTransferV2(db vm.StateDB, addr common.Address, amount *big.Int) bool {
	return db.(*ethStateV2).GetBalanceNoRecord(addr).Cmp(amount) >= 0
}

// Transfer subtracts amount from sender and adds amount to recipient using the given Db
func Transfer(db vm.StateDB, sender, recipient common.Address, amount *big.Int) {
	db.SubBalance(sender, amount)
	db.AddBalance(recipient, amount)
}
