package transfer

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
	ccEu "github.com/arcology-network/vm-adaptor"
	ccApi "github.com/arcology-network/vm-adaptor/api"
	eth "github.com/arcology-network/vm-adaptor/eth"
	tests "github.com/arcology-network/vm-adaptor/tests"
)

func TestTransfer(t *testing.T) {
	config := tests.MainConfig()
	persistentDB := cachedstorage.NewDataStore()
	meta, _ := commutative.NewMeta(urlcommon.NewPlatform().Eth10Account())
	persistentDB.Inject(urlcommon.NewPlatform().Eth10Account(), meta)
	db := curstorage.NewTransientDB(persistentDB)

	url := concurrenturl.NewConcurrentUrl(db)
	// api := ccApi.NewAPI(url)
	statedb := eth.NewImplStateDB(url)
	statedb.Prepare(common.Hash{}, common.Hash{}, 0)
	statedb.CreateAccount(tests.Coinbase)
	statedb.CreateAccount(tests.User1)
	statedb.AddBalance(tests.User1, new(big.Int).SetUint64(1e18))
	_, transitions := url.Export(true)
	t.Log("\n" + tests.FormatTransitions(transitions))

	// Transfer.
	url = concurrenturl.NewConcurrentUrl(db)
	url.Import(transitions)
	url.PostImport()
	url.Commit([]uint32{0})
	api = ccApi.NewAPI(url)
	statedb = eth.NewImplStateDB(url)
	eu := ccEu.NewEU(config.ChainConfig, *config.VMConfig, config.Chain, statedb, api, url)

	config.Coinbase = &tests.Coinbase
	config.BlockNumber = new(big.Int).SetUint64(10000000)
	config.Time = new(big.Int).SetUint64(10000000)

	msg := types.NewMessage(tests.User1, &tests.User2, 0, new(big.Int).SetUint64(100), 1e15, new(big.Int).SetUint64(1), nil, nil, true)
	accesses, transitions, receipt, err := eu.Run(common.BytesToHash([]byte{1, 1, 1}), 1, &msg, ccEu.NewEVMBlockContextV2(config), ccEu.NewEVMTxContext(msg))
	t.Log("\n" + tests.FormatTransitions(accesses))
	t.Log("\n" + tests.FormatTransitions(transitions))
	t.Log(receipt)

	if receipt.Status != 1 {
		t.Error("Error: Transfer failed !!!", err)
	}
}
