package eth

import (
	"bytes"
	"fmt"
	"math/big"
	"testing"
	"time"

	cachedstorage "github.com/arcology-network/common-lib/cachedstorage"
	"github.com/arcology-network/concurrenturl/v2"
	urlcommon "github.com/arcology-network/concurrenturl/v2/common"
	commutative "github.com/arcology-network/concurrenturl/v2/type/commutative"
	evmcommon "github.com/arcology-network/evm/common"
	eth "github.com/arcology-network/vm-adaptor/eth"
)

func TestStateDBV2GetNonexistBalance(t *testing.T) {
	db := cachedstorage.NewDataStore()
	meta, _ := commutative.NewMeta(urlcommon.NewPlatform().Eth10Account())
	db.Inject(urlcommon.NewPlatform().Eth10Account(), meta)
	url := concurrenturl.NewConcurrentUrl(db)

	account := evmcommon.BytesToAddress([]byte{201, 202, 203, 204, 205})
	statedb := eth.NewStateDB(nil, db, url)
	statedb.Prepare(evmcommon.Hash{}, evmcommon.Hash{}, 1)
	statedb.CreateAccount(account)
	_, transitions := url.Export(true)
	fmt.Println("\n" + FormatTransitions(transitions))
	url.Import(transitions)
	url.PostImport()
	url.Commit([]uint32{1})

	url = concurrenturl.NewConcurrentUrl(db)
	statedb = eth.NewStateDB(nil, db, url)
	statedb.Prepare(evmcommon.Hash{}, evmcommon.Hash{}, 2)
	balance := statedb.GetBalance(account)
	if balance == nil || balance.Cmp(new(big.Int)) != 0 {
		t.Fail()
	}
}

func TestStateDBV2GetNonexistCode(t *testing.T) {
	db := cachedstorage.NewDataStore()
	meta, _ := commutative.NewMeta(urlcommon.NewPlatform().Eth10Account())
	db.Inject(urlcommon.NewPlatform().Eth10Account(), meta)
	url := concurrenturl.NewConcurrentUrl(db)

	account := evmcommon.BytesToAddress([]byte{201, 202, 203, 204, 205})
	statedb := eth.NewStateDB(nil, db, url)
	statedb.Prepare(evmcommon.Hash{}, evmcommon.Hash{}, 1)
	statedb.CreateAccount(account)
	_, transitions := url.Export(true)
	fmt.Println("\n" + FormatTransitions(transitions))
	url.Import(transitions)
	url.PostImport()
	url.Commit([]uint32{1})

	url = concurrenturl.NewConcurrentUrl(db)
	statedb = NewStateDB(nil, db, url)
	statedb.Prepare(evmcommon.Hash{}, evmcommon.Hash{}, 2)
	code := statedb.GetCode(account)
	if code == nil || len(code) != 0 {
		t.Fail()
	}
}

func TestStateDBV2GetNonexistStorageState(t *testing.T) {
	db := cachedstorage.NewDataStore()
	meta, _ := commutative.NewMeta(urlcommon.NewPlatform().Eth10Account())
	db.Inject(urlcommon.NewPlatform().Eth10Account(), meta)
	url := concurrenturl.NewConcurrentUrl(db)

	account := evmcommon.BytesToAddress([]byte{201, 202, 203, 204, 205})
	statedb := NewStateDB(nil, db, url)
	statedb.Prepare(evmcommon.Hash{}, evmcommon.Hash{}, 1)
	statedb.CreateAccount(account)
	_, transitions := url.Export(true)
	fmt.Println("\n" + FormatTransitions(transitions))
	url.Import(transitions)
	url.PostImport()
	url.Commit([]uint32{1})

	url = concurrenturl.NewConcurrentUrl(db)
	statedb = NewStateDB(nil, db, url)
	statedb.Prepare(evmcommon.Hash{}, evmcommon.Hash{}, 2)
	state := statedb.GetState(account, evmcommon.Hash{})
	if !bytes.Equal(state.Bytes(), evmcommon.Hash{}.Bytes()) {
		t.Fail()
	}
}

func TestStateDBV2(t *testing.T) {
	db := cachedstorage.NewDataStore()
	meta, _ := commutative.NewMeta(urlcommon.NewPlatform().Eth10Account())
	db.Inject(urlcommon.NewPlatform().Eth10Account(), meta)
	url := concurrenturl.NewConcurrentUrl(db)

	account := evmcommon.BytesToAddress([]byte{201, 202, 203, 204, 205})
	statedb := NewStateDB(nil, db, url)
	statedb.Prepare(evmcommon.Hash{}, evmcommon.Hash{}, 1)
	statedb.CreateAccount(account)
	statedb.AddBalance(account, new(big.Int).SetUint64(100))
	statedb.SubBalance(account, new(big.Int).SetUint64(10)) // Balance 90

	accesses, transitions := url.Export(true)
	fmt.Println("\n" + FormatTransitions(accesses))
	fmt.Println("\n" + FormatTransitions(transitions))
	url.Import(transitions)
	url.PostImport()
	url.Commit([]uint32{1})

	url1 := concurrenturl.NewConcurrentUrl(db)
	statedb1 := NewStateDB(nil, db, url1)
	statedb1.Prepare(evmcommon.Hash{}, evmcommon.Hash{}, 2)

	url2 := concurrenturl.NewConcurrentUrl(db)
	statedb2 := NewStateDB(nil, db, url2)
	statedb2.Prepare(evmcommon.Hash{}, evmcommon.Hash{}, 3)

	statedb1.AddBalance(account, new(big.Int).SetUint64(200)) // + 200 + 90
	statedb2.AddBalance(account, new(big.Int).SetUint64(300)) // + 300 + 90

	ar1, t1 := url1.Export(true)
	fmt.Println("\n" + FormatTransitions(ar1))
	fmt.Println("\n" + FormatTransitions(t1))

	ar2, t2 := url2.Export(true)
	fmt.Println("\n" + FormatTransitions(ar2))
	fmt.Println("\n" + FormatTransitions(t2))

	txs, groups, flags := detectConflict(append(ar1, ar2...))
	fmt.Println(txs)
	fmt.Println(groups)
	fmt.Println(flags)

	url.Import(append(t1, t2...))
	url.PostImport()
	url.Commit([]uint32{2, 3})
	url = concurrenturl.NewConcurrentUrl(db)
	statedb = NewStateDB(nil, db, url)
	balance := statedb.GetBalance(account)
	fmt.Println(balance)

	if balance.Int64() != 590 {
		t.Error("Error: Expected: 590 ,", "Actual: ", balance.Int64())
	}
}

func TestStateDBV2BalanceReadWriteConflict(t *testing.T) {
	db := cachedstorage.NewDataStore()
	meta, _ := commutative.NewMeta(urlcommon.NewPlatform().Eth10Account())
	db.Inject(urlcommon.NewPlatform().Eth10Account(), meta)
	url := concurrenturl.NewConcurrentUrl(db)

	account := evmcommon.BytesToAddress([]byte{201, 202, 203, 204, 205})
	statedb := NewStateDB(nil, db, url)
	statedb.Prepare(evmcommon.Hash{}, evmcommon.Hash{}, 1)
	statedb.CreateAccount(account)
	statedb.AddBalance(account, new(big.Int).SetUint64(100))

	_, transitions := url.Export(true)
	url.Import(transitions)
	url.PostImport()
	url.Commit([]uint32{1}) // Write balance 100 to the storage

	url1 := concurrenturl.NewConcurrentUrl(db)
	statedb1 := NewStateDB(nil, db, url1)
	statedb1.Prepare(evmcommon.Hash{}, evmcommon.Hash{}, 2)

	url2 := concurrenturl.NewConcurrentUrl(db)
	statedb2 := NewStateDB(nil, db, url2)
	statedb2.Prepare(evmcommon.Hash{}, evmcommon.Hash{}, 3)

	url3 := concurrenturl.NewConcurrentUrl(db)
	statedb3 := NewStateDB(nil, db, url3)
	statedb3.Prepare(evmcommon.Hash{}, evmcommon.Hash{}, 4)

	if b := statedb1.GetBalance(account); b.Uint64() != 100 { // read the balance, should conflict with other txs.
		t.Fail()
	}

	// 100 + 300 - 50 == 350
	statedb2.AddBalance(account, new(big.Int).SetUint64(300))
	statedb3.SubBalance(account, new(big.Int).SetUint64(50))

	access1, tx1 := url1.Export(true)
	fmt.Println("Access Records: ", FormatTransitions(access1))
	fmt.Println("Transition    : ", FormatTransitions(tx1))
	fmt.Println()

	access2, tx2 := url2.Export(true)
	fmt.Println("Access Records: ", FormatTransitions(access2))
	fmt.Println("Transition    : ", FormatTransitions(tx2))
	fmt.Println()

	access3, tx3 := url3.Export(true)
	fmt.Println("Access Records: ", FormatTransitions(access3))
	fmt.Println("Transition    : ", FormatTransitions(tx3))
	fmt.Println()

	txs, groups, flags := detectConflict(append(append(access1, access2...), access3...))
	fmt.Println(txs)
	fmt.Println(groups)
	fmt.Println(flags)

	if len(flags) != 2 {
		t.Error("Error: There should be two conflicting TXs")
	}
}

func TestStateDBV2NonceWrite(t *testing.T) {
	db := cachedstorage.NewDataStore()
	meta, _ := commutative.NewMeta(urlcommon.NewPlatform().Eth10Account())
	db.Inject(urlcommon.NewPlatform().Eth10Account(), meta)
	url := concurrenturl.NewConcurrentUrl(db)

	account := evmcommon.BytesToAddress([]byte{201, 202, 203, 204, 205})
	statedb := NewStateDB(nil, db, url)
	statedb.Prepare(evmcommon.Hash{}, evmcommon.Hash{}, 1)
	statedb.CreateAccount(account)

	_, transitions := url.Export(true)
	url.Import(transitions)
	url.PostImport()
	url.Commit([]uint32{1})

	url1 := concurrenturl.NewConcurrentUrl(db)
	statedb1 := NewStateDB(nil, db, url1)
	statedb1.Prepare(evmcommon.Hash{}, evmcommon.Hash{}, 2)
	url2 := concurrenturl.NewConcurrentUrl(db)
	statedb2 := NewStateDB(nil, db, url2)
	statedb2.Prepare(evmcommon.Hash{}, evmcommon.Hash{}, 3)

	statedb1.SetNonce(account, 10)
	statedb2.SetNonce(account, 20)

	ar1, t1 := url1.Export(true)
	ar2, t2 := url2.Export(true)
	fmt.Println("\n" + FormatTransitions(ar1))
	fmt.Println("\n" + FormatTransitions(t1))
	fmt.Println("\n" + FormatTransitions(ar2))
	fmt.Println("\n" + FormatTransitions(t2))
	txs, groups, flags := detectConflict(append(ar1, ar2...))
	fmt.Println(txs)
	fmt.Println(groups)
	fmt.Println(flags)

	url.Import(append(t1, t2...))
	url.PostImport()
	url.Commit([]uint32{2, 3})
	url = concurrenturl.NewConcurrentUrl(db)
	statedb = NewStateDB(nil, db, url)
	nonce := statedb.GetNonce(account)

	if len(flags) != 0 {
		t.Error("Here should be no conflict")
	}

	if nonce != 2 {
		t.Error("Nonce should be equal to 2")
	}
}

func TestExport1(b *testing.T) {
	db := cachedstorage.NewDataStore()
	meta, _ := commutative.NewMeta(urlcommon.NewPlatform().Eth10Account())
	db.Inject(urlcommon.NewPlatform().Eth10Account(), meta)

	begin := time.Now()
	var transitions []urlcommon.UnivalueInterface
	for i := 0; i < 1000; i++ {
		url := concurrenturl.NewConcurrentUrl(db)
		statedb := NewStateDB(nil, db, url)
		statedb.Prepare(evmcommon.Hash{}, evmcommon.Hash{}, 1)
		for j := 0; j < 10; j++ {
			// acc := evmcommon.BytesToAddress([]byte(RandStringRunes(20)))
			acc := evmcommon.BytesToAddress([]byte{byte(j)})
			statedb.CreateAccount(acc)
			statedb.AddBalance(acc, new(big.Int).SetUint64(1e18))
		}
		_, ts := url.Export(true)
		transitions = append(transitions, ts...)
	}
	b.Log(time.Duration(time.Since(begin)))
	b.Log(len(transitions))
	b.Log("\n" + FormatTransitions(transitions[:9]))

	// begin = time.Now()
	// url := concurrenturl.NewConcurrentUrl(db)
	// url.Commit(transitions, []uint32{1})
	// b.Log(time.Duration(time.Since(begin)))
}

func BenchmarkExport(b *testing.B) {
	db := cachedstorage.NewDataStore()
	meta, _ := commutative.NewMeta(urlcommon.NewPlatform().Eth10Account())
	db.Inject(urlcommon.NewPlatform().Eth10Account(), meta)
	url := concurrenturl.NewConcurrentUrl(db)

	statedb := NewStateDB(nil, db, url)
	statedb.Prepare(evmcommon.Hash{}, evmcommon.Hash{}, 1)

	begin := time.Now()
	for i := 0; i < 1000; i++ {
		acc := evmcommon.BytesToAddress([]byte(RandStringRunes(20)))
		statedb.CreateAccount(acc)
		statedb.AddBalance(acc, new(big.Int).SetUint64(1e18))
	}
	b.Log(time.Duration(time.Since(begin)))
	b.ResetTimer()

	var transitions []urlcommon.UnivalueInterface
	for i := 0; i < b.N; i++ {
		begin = time.Now()
		_, transitions = url.Export(true)
		b.Log(len(transitions))
		// b.Log(time.Duration(time.Since(begin)))
	}

	b.Log("\n" + FormatTransitions(transitions[:9]))
}
