package tests

import (
	"math/big"
	"testing"

	cachedstorage "github.com/arcology-network/common-lib/cachedstorage"
	"github.com/arcology-network/concurrenturl/v2"
	urlcommon "github.com/arcology-network/concurrenturl/v2/common"
	curstorage "github.com/arcology-network/concurrenturl/v2/storage"
	"github.com/arcology-network/concurrenturl/v2/type/commutative"
	evmcommon "github.com/arcology-network/evm/common"
	adaptor "github.com/arcology-network/vm-adaptor/evm"
)

func TestNativeStorage(t *testing.T) {
	persistentDB := cachedstorage.NewDataStore()
	meta, _ := commutative.NewMeta(urlcommon.NewPlatform().Eth10Account())
	persistentDB.Inject(urlcommon.NewPlatform().Eth10Account(), meta)
	db := curstorage.NewTransientDB(persistentDB)

	url := concurrenturl.NewConcurrentUrl(db)
	api := adaptor.NewAPI(db, url)
	statedb := adaptor.NewStateDB(api, db, url)
	statedb.Prepare(evmcommon.Hash{}, evmcommon.Hash{}, 0)
	statedb.CreateAccount(coinbase)
	statedb.CreateAccount(owner)
	statedb.AddBalance(owner, new(big.Int).SetUint64(1e18))
	statedb.CreateAccount(user1)
	statedb.AddBalance(user1, new(big.Int).SetUint64(1e18))
	_, transitions := url.Export(true)

	// Deploy NativeStorage.
	eu, config := prepare(db, 10000000, transitions, []uint32{0})
	transitions, receipt := deploy(eu, config, owner, 0, NativeStorageCode)
	t.Log("\n" + FormatTransitions(transitions))
	t.Log(receipt)
	address := receipt.ContractAddress
	t.Log(address)

	// Call accessX.
	eu, config = prepare(db, 10000001, transitions, []uint32{1})
	acc, transitions, receipt := runEx(eu, config, &user1, &address, 1, true, "accessX()")
	t.Log("\n" + FormatTransitions(acc))
	t.Log(receipt)

	// Call accessY.
	eu, config = prepare(db, 10000002, transitions, []uint32{2})
	acc, _, receipt = runEx(eu, config, &user1, &address, 2, true, "accessY()")
	t.Log("\n" + FormatTransitions(acc))
	t.Log(receipt)
}
