package tests

import (
	"math/big"
	"testing"

	cachedstorage "github.com/arcology-network/common-lib/cachedstorage"
	"github.com/arcology-network/concurrenturl/v2"
	urlcommon "github.com/arcology-network/concurrenturl/v2/common"
	curstorage "github.com/arcology-network/concurrenturl/v2/storage"
	"github.com/arcology-network/concurrenturl/v2/type/commutative"
	"github.com/arcology-network/evm/common"
	"github.com/arcology-network/evm/core/types"
	adaptor "github.com/arcology-network/vm-adaptor/evm"
)

func TestTransfer(t *testing.T) {
	config := MainConfig()
	persistentDB := cachedstorage.NewDataStore()
	meta, _ := commutative.NewMeta(urlcommon.NewPlatform().Eth10Account())
	persistentDB.Inject(urlcommon.NewPlatform().Eth10Account(), meta)
	db := curstorage.NewTransientDB(persistentDB)

	url := concurrenturl.NewConcurrentUrl(db)
	api := adaptor.NewAPIV2(db, url)
	statedb := adaptor.NewStateDBV2(api, db, url)
	statedb.Prepare(common.Hash{}, common.Hash{}, 0)
	statedb.CreateAccount(coinbase)
	statedb.CreateAccount(user1)
	statedb.AddBalance(user1, new(big.Int).SetUint64(1e18))
	_, transitions := url.Export(true)
	t.Log("\n" + formatTransitions(transitions))

	// Transfer.
	url = concurrenturl.NewConcurrentUrl(db)
	url.Import(transitions)
	url.PostImport()
	url.Commit([]uint32{0})
	api = adaptor.NewAPIV2(db, url)
	statedb = adaptor.NewStateDBV2(api, db, url)
	eu := adaptor.NewEUV2(config.ChainConfig, *config.VMConfig, config.Chain, statedb, api, db, url)

	config.Coinbase = &coinbase
	config.BlockNumber = new(big.Int).SetUint64(10000000)
	config.Time = new(big.Int).SetUint64(10000000)

	msg := types.NewMessage(user1, &user2, 0, new(big.Int).SetUint64(100), 1e15, new(big.Int).SetUint64(1), nil, nil, true)
	accesses, transitions, receipt := eu.Run(common.BytesToHash([]byte{1, 1, 1}), 1, &msg, adaptor.NewEVMBlockContextV2(config), adaptor.NewEVMTxContext(msg))
	t.Log("\n" + formatTransitions(accesses))
	t.Log("\n" + formatTransitions(transitions))
	t.Log(receipt)
}
