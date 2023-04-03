package tests

import (
	"math/big"
	"testing"

	cachedstorage "github.com/arcology-network/common-lib/cachedstorage"
	"github.com/arcology-network/concurrenturl/v2"
	urlcmn "github.com/arcology-network/concurrenturl/v2/common"
	"github.com/arcology-network/concurrenturl/v2/type/commutative"
	evmcmn "github.com/arcology-network/evm/common"
	adaptor "github.com/arcology-network/vm-adaptor/eth"
)

func TestDynamicArrayBasic(t *testing.T) {
	db := cachedstorage.NewDataStore()
	meta, _ := commutative.NewMeta(urlcmn.NewPlatform().Eth10Account())
	db.Inject(urlcmn.NewPlatform().Eth10Account(), meta)

	url := concurrenturl.NewConcurrentUrl(db)
	api := adaptor.NewAPI(db, url)
	stateDB := adaptor.NewStateDB(api, db, url)
	stateDB.Prepare(evmcmn.Hash{}, evmcmn.Hash{}, 0)
	stateDB.CreateAccount(coinbase)
	stateDB.CreateAccount(owner)
	stateDB.AddBalance(owner, new(big.Int).SetUint64(1e18))
	_, transitions := url.Export(true)

	// Deploy DynamicArrayTest.
	eu, config := prepare(db, 10000000, transitions, []uint32{0})
	transitions, receipt := deploy(eu, config, owner, 0, darrayTestCode)
	if receipt.Status != 1 {
		t.Error("Error: Execution failed !!!")
	}

	t.Log("\n" + FormatTransitions(transitions))
	t.Log(receipt)
	// contractAddress := receipt.ContractAddress
	// t.Log(contractAddress)

	// Call push.
	// eu, config = prepare(db, 10000001, transitions, []uint32{1})
	// _, transitions, receipt = runEx(eu, config, &owner, &contractAddress, 1, false, "push(bytes)", []byte{32}, []byte{36}, []byte("abcdefghijklmnopqrstuvwxyz012345"), []byte("6789xxxxxxxxxxxxxxxxxxxxxxxxxxxx"))
	// t.Log("\n" + FormatTransitions(transitions))
	// t.Log(receipt)

	// // Call tryPop.
	// eu, config = prepare(db, 10000002, transitions, []uint32{2})
	// _, transitions, receipt = runEx(eu, config, &owner, &contractAddress, 2, false, "tryPop()")
	// t.Log("\n" + FormatTransitions(transitions))
	// t.Log(receipt)

	// // Call tryPop again.
	// eu, config = prepare(db, 10000003, transitions, []uint32{3})
	// _, transitions, receipt = runEx(eu, config, &owner, &contractAddress, 3, false, "tryPop()")
	// t.Log("\n" + FormatTransitions(transitions))
	// t.Log(receipt)

	// // Call push2.
	// eu, config = prepare(db, 10000004, transitions, []uint32{4})
	// _, transitions, receipt = runEx(eu, config, &owner, &contractAddress, 4, false, "push2(bytes)", []byte{32}, []byte{10}, []byte("0123456789xxxxxxxxxxxxxxxxxxxxxx"))
	// t.Log("\n" + FormatTransitions(transitions))
	// t.Log(receipt)

	// // Call tryPop2.
	// eu, config = prepare(db, 10000005, transitions, []uint32{5})
	// _, transitions, receipt = runEx(eu, config, &owner, &contractAddress, 5, false, "tryPop2()")
	// t.Log("\n" + FormatTransitions(transitions))
	// t.Log(receipt)
}
