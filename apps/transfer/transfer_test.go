package transfer

import (
	"math/big"
	"testing"

	cachedstorage "github.com/arcology-network/common-lib/cachedstorage"
	"github.com/arcology-network/concurrenturl"
	"github.com/arcology-network/concurrenturl/commutative"
	curstorage "github.com/arcology-network/concurrenturl/storage"
	"github.com/arcology-network/evm/common"
	"github.com/arcology-network/evm/core"
	ccEu "github.com/arcology-network/vm-adaptor"
	ccApi "github.com/arcology-network/vm-adaptor/api"
	vmcommon "github.com/arcology-network/vm-adaptor/common"
	eth "github.com/arcology-network/vm-adaptor/eth"
	tests "github.com/arcology-network/vm-adaptor/tests"
)

func TestTransfer(t *testing.T) {
	coinbase := vmcommon.Coinbase
	// owner := vmcommon.Owner
	user1 := vmcommon.Alice
	user2 := vmcommon.Bob

	config := tests.MainTestConfig()
	persistentDB := cachedstorage.NewDataStore()
	meta := commutative.NewPath()
	persistentDB.Inject((&concurrenturl.Platform{}).Eth10Account(), meta)
	db := curstorage.NewTransientDB(persistentDB)

	url := concurrenturl.NewConcurrentUrl(db)
	api := ccApi.NewAPI(url)
	statedb := eth.NewImplStateDB(api)
	statedb.PrepareFormer(common.Hash{}, common.Hash{}, 0)
	statedb.CreateAccount(coinbase)
	statedb.CreateAccount(user1)
	statedb.AddBalance(user1, new(big.Int).SetUint64(1e18))
	_, transitions := url.ExportAll()
	t.Log("\n" + vmcommon.FormatTransitions(transitions))

	// Transfer.
	url = concurrenturl.NewConcurrentUrl(db)
	url.Import(transitions)
	url.Sort()
	url.Commit([]uint32{0})
	api = ccApi.NewAPI(url)
	statedb = eth.NewImplStateDB(api)
	eu := ccEu.NewEU(config.ChainConfig, *config.VMConfig, statedb, api)

	config.Coinbase = &coinbase
	config.BlockNumber = new(big.Int).SetUint64(10000000)
	config.Time = new(big.Int).SetUint64(10000000)

	msg := core.NewMessage(user1, &user2, 0, new(big.Int).SetUint64(100), 1e15, new(big.Int).SetUint64(1), nil, nil, true)
	receipt, _, err := eu.Run(common.BytesToHash([]byte{1, 1, 1}), 1, &msg, ccEu.NewEVMBlockContext(config), ccEu.NewEVMTxContext(msg))
	accesses, transitions := eu.Api().StateFilter().ByType()

	t.Log("\n" + vmcommon.FormatTransitions(accesses))
	t.Log("\n" + vmcommon.FormatTransitions(transitions))
	t.Log(receipt)

	if receipt.Status != 1 {
		t.Error("Error: Transfer failed !!!", err)
	}
}
