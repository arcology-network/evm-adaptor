package evm_test

import (
	"bytes"
	"fmt"
	"math/big"
	"math/rand"
	"testing"
	"time"

	"github.com/arcology-network/concurrenturl/v2"
	urlcommon "github.com/arcology-network/concurrenturl/v2/common"
	urltype "github.com/arcology-network/concurrenturl/v2/type"
	commutative "github.com/arcology-network/concurrenturl/v2/type/commutative"
	noncommutative "github.com/arcology-network/concurrenturl/v2/type/noncommutative"
	evmcommon "github.com/arcology-network/evm/common"
	arbitrator "github.com/arcology-network/urlarbitrator-engine/go-wrapper"
	"github.com/arcology-network/vm-adaptor/evm"
)

func TestStateDBV2GetNonexistBalance(t *testing.T) {
	db := urlcommon.NewDataStore()
	meta, _ := commutative.NewMeta(urlcommon.NewPlatform().Eth10Account())
	db.Save(urlcommon.NewPlatform().Eth10Account(), meta)
	url := concurrenturl.NewConcurrentUrl(db)

	account := evmcommon.BytesToAddress([]byte{201, 202, 203, 204, 205})
	statedb := evm.NewStateDBV2(nil, db, url)
	statedb.Prepare(evmcommon.Hash{}, evmcommon.Hash{}, 1)
	statedb.CreateAccount(account)
	_, transitions := url.Export(true)
	t.Log("\n" + formatTransitions(transitions))
	url.Commit(transitions, []uint32{1})

	url = concurrenturl.NewConcurrentUrl(db)
	statedb = evm.NewStateDBV2(nil, db, url)
	statedb.Prepare(evmcommon.Hash{}, evmcommon.Hash{}, 2)
	balance := statedb.GetBalance(account)
	if balance == nil || balance.Cmp(new(big.Int)) != 0 {
		t.Fail()
	}
}

func TestStateDBV2GetNonexistCode(t *testing.T) {
	db := urlcommon.NewDataStore()
	meta, _ := commutative.NewMeta(urlcommon.NewPlatform().Eth10Account())
	db.Save(urlcommon.NewPlatform().Eth10Account(), meta)
	url := concurrenturl.NewConcurrentUrl(db)

	account := evmcommon.BytesToAddress([]byte{201, 202, 203, 204, 205})
	statedb := evm.NewStateDBV2(nil, db, url)
	statedb.Prepare(evmcommon.Hash{}, evmcommon.Hash{}, 1)
	statedb.CreateAccount(account)
	_, transitions := url.Export(true)
	t.Log("\n" + formatTransitions(transitions))
	url.Commit(transitions, []uint32{1})

	url = concurrenturl.NewConcurrentUrl(db)
	statedb = evm.NewStateDBV2(nil, db, url)
	statedb.Prepare(evmcommon.Hash{}, evmcommon.Hash{}, 2)
	code := statedb.GetCode(account)
	if code == nil || len(code) != 0 {
		t.Fail()
	}
}

func TestStateDBV2GetNonexistStorageState(t *testing.T) {
	db := urlcommon.NewDataStore()
	meta, _ := commutative.NewMeta(urlcommon.NewPlatform().Eth10Account())
	db.Save(urlcommon.NewPlatform().Eth10Account(), meta)
	url := concurrenturl.NewConcurrentUrl(db)

	account := evmcommon.BytesToAddress([]byte{201, 202, 203, 204, 205})
	statedb := evm.NewStateDBV2(nil, db, url)
	statedb.Prepare(evmcommon.Hash{}, evmcommon.Hash{}, 1)
	statedb.CreateAccount(account)
	_, transitions := url.Export(true)
	t.Log("\n" + formatTransitions(transitions))
	url.Commit(transitions, []uint32{1})

	url = concurrenturl.NewConcurrentUrl(db)
	statedb = evm.NewStateDBV2(nil, db, url)
	statedb.Prepare(evmcommon.Hash{}, evmcommon.Hash{}, 2)
	state := statedb.GetState(account, evmcommon.Hash{})
	if !bytes.Equal(state.Bytes(), evmcommon.Hash{}.Bytes()) {
		t.Fail()
	}
}

func TestStateDBV2(t *testing.T) {
	db := urlcommon.NewDataStore()
	meta, _ := commutative.NewMeta(urlcommon.NewPlatform().Eth10Account())
	db.Save(urlcommon.NewPlatform().Eth10Account(), meta)
	url := concurrenturl.NewConcurrentUrl(db)

	account := evmcommon.BytesToAddress([]byte{201, 202, 203, 204, 205})
	statedb := evm.NewStateDBV2(nil, db, url)
	statedb.Prepare(evmcommon.Hash{}, evmcommon.Hash{}, 1)
	statedb.CreateAccount(account)
	statedb.AddBalance(account, new(big.Int).SetUint64(100))
	statedb.SubBalance(account, new(big.Int).SetUint64(10))

	accesses, transitions := url.Export(true)
	t.Log("\n" + formatTransitions(accesses))
	t.Log("\n" + formatTransitions(transitions))
	url.Commit(transitions, []uint32{1})

	url1 := concurrenturl.NewConcurrentUrl(db)
	statedb1 := evm.NewStateDBV2(nil, db, url1)
	statedb1.Prepare(evmcommon.Hash{}, evmcommon.Hash{}, 2)
	url2 := concurrenturl.NewConcurrentUrl(db)
	statedb2 := evm.NewStateDBV2(nil, db, url2)
	statedb2.Prepare(evmcommon.Hash{}, evmcommon.Hash{}, 3)

	statedb1.AddBalance(account, new(big.Int).SetUint64(200))
	statedb2.AddBalance(account, new(big.Int).SetUint64(300))

	ar1, t1 := url1.Export(true)
	ar2, t2 := url2.Export(true)
	t.Log("\n" + formatTransitions(ar1))
	t.Log("\n" + formatTransitions(t1))
	t.Log("\n" + formatTransitions(ar2))
	t.Log("\n" + formatTransitions(t2))
	txs, groups, flags := detectConflict(append(ar1, ar2...))
	t.Log(txs)
	t.Log(groups)
	t.Log(flags)

	url.Commit(append(t1, t2...), []uint32{2, 3})
	url = concurrenturl.NewConcurrentUrl(db)
	statedb = evm.NewStateDBV2(nil, db, url)
	balance := statedb.GetBalance(account)
	t.Log(balance)
}

func TestStateDBV2BalanceReadWriteConflict(t *testing.T) {
	db := urlcommon.NewDataStore()
	meta, _ := commutative.NewMeta(urlcommon.NewPlatform().Eth10Account())
	db.Save(urlcommon.NewPlatform().Eth10Account(), meta)
	url := concurrenturl.NewConcurrentUrl(db)

	account := evmcommon.BytesToAddress([]byte{201, 202, 203, 204, 205})
	statedb := evm.NewStateDBV2(nil, db, url)
	statedb.Prepare(evmcommon.Hash{}, evmcommon.Hash{}, 1)
	statedb.CreateAccount(account)
	statedb.AddBalance(account, new(big.Int).SetUint64(100))

	_, transitions := url.Export(true)
	url.Commit(transitions, []uint32{1})

	url1 := concurrenturl.NewConcurrentUrl(db)
	statedb1 := evm.NewStateDBV2(nil, db, url1)
	statedb1.Prepare(evmcommon.Hash{}, evmcommon.Hash{}, 2)
	url2 := concurrenturl.NewConcurrentUrl(db)
	statedb2 := evm.NewStateDBV2(nil, db, url2)
	statedb2.Prepare(evmcommon.Hash{}, evmcommon.Hash{}, 3)
	url3 := concurrenturl.NewConcurrentUrl(db)
	statedb3 := evm.NewStateDBV2(nil, db, url3)
	statedb3.Prepare(evmcommon.Hash{}, evmcommon.Hash{}, 4)

	if b := statedb1.GetBalance(account); b.Uint64() != 100 {
		t.Fail()
	}
	statedb2.AddBalance(account, new(big.Int).SetUint64(300))
	statedb3.SubBalance(account, new(big.Int).SetUint64(50))

	ar1, t1 := url1.Export(true)
	ar2, t2 := url2.Export(true)
	ar3, t3 := url3.Export(true)
	t.Log("\n" + formatTransitions(ar1))
	t.Log("\n" + formatTransitions(t1))
	t.Log("\n" + formatTransitions(ar2))
	t.Log("\n" + formatTransitions(t2))
	t.Log("\n" + formatTransitions(ar3))
	t.Log("\n" + formatTransitions(t3))
	txs, groups, flags := detectConflict(append(append(ar1, ar2...), ar3...))
	t.Log(txs)
	t.Log(groups)
	t.Log(flags)
}

func TestStateDBV2NonceWrite(t *testing.T) {
	db := urlcommon.NewDataStore()
	meta, _ := commutative.NewMeta(urlcommon.NewPlatform().Eth10Account())
	db.Save(urlcommon.NewPlatform().Eth10Account(), meta)
	url := concurrenturl.NewConcurrentUrl(db)

	account := evmcommon.BytesToAddress([]byte{201, 202, 203, 204, 205})
	statedb := evm.NewStateDBV2(nil, db, url)
	statedb.Prepare(evmcommon.Hash{}, evmcommon.Hash{}, 1)
	statedb.CreateAccount(account)

	_, transitions := url.Export(true)
	url.Commit(transitions, []uint32{1})

	url1 := concurrenturl.NewConcurrentUrl(db)
	statedb1 := evm.NewStateDBV2(nil, db, url1)
	statedb1.Prepare(evmcommon.Hash{}, evmcommon.Hash{}, 2)
	url2 := concurrenturl.NewConcurrentUrl(db)
	statedb2 := evm.NewStateDBV2(nil, db, url2)
	statedb2.Prepare(evmcommon.Hash{}, evmcommon.Hash{}, 3)

	statedb1.SetNonce(account, 10)
	statedb2.SetNonce(account, 20)

	ar1, t1 := url1.Export(true)
	ar2, t2 := url2.Export(true)
	t.Log("\n" + formatTransitions(ar1))
	t.Log("\n" + formatTransitions(t1))
	t.Log("\n" + formatTransitions(ar2))
	t.Log("\n" + formatTransitions(t2))
	txs, groups, flags := detectConflict(append(ar1, ar2...))
	t.Log(txs)
	t.Log(groups)
	t.Log(flags)

	url.Commit(append(t1, t2...), []uint32{2, 3})
	url = concurrenturl.NewConcurrentUrl(db)
	statedb = evm.NewStateDBV2(nil, db, url)
	nonce := statedb.GetNonce(account)
	t.Log(nonce)
}

func TestExportOnFailure(t *testing.T) {
	db := urlcommon.NewDataStore()
	meta, _ := commutative.NewMeta(urlcommon.NewPlatform().Eth10Account())
	db.Save(urlcommon.NewPlatform().Eth10Account(), meta)

	url := concurrenturl.NewConcurrentUrl(db)
	account := evmcommon.BytesToAddress([]byte{0xcc})
	coinbase := evmcommon.BytesToAddress([]byte{0xdd})
	statedb := evm.NewStateDBV2(nil, db, url)
	statedb.Prepare(evmcommon.Hash{}, evmcommon.Hash{}, 1)
	statedb.CreateAccount(account)
	statedb.CreateAccount(coinbase)
	_, transitions := url.Export(true)
	url.Commit(transitions, []uint32{1})

	accesses, transitions := evm.ExportOnFailure(
		db,
		2,
		account,
		coinbase,
		1024,
		new(big.Int).SetUint64(32),
	)
	t.Log("\n" + formatTransitions(accesses))
	t.Log("\n" + formatTransitions(transitions))
}

func TestExportOnConfliction(t *testing.T) {
	db := urlcommon.NewDataStore()
	meta, _ := commutative.NewMeta(urlcommon.NewPlatform().Eth10Account())
	db.Save(urlcommon.NewPlatform().Eth10Account(), meta)

	url := concurrenturl.NewConcurrentUrl(db)
	account := evmcommon.BytesToAddress([]byte{0xcc})
	statedb := evm.NewStateDBV2(nil, db, url)
	statedb.Prepare(evmcommon.Hash{}, evmcommon.Hash{}, 1)
	statedb.CreateAccount(account)
	_, transitions := url.Export(true)
	url.Commit(transitions, []uint32{1})

	accesses, transitions := evm.ExportOnConfliction(
		db,
		2,
		evmcommon.BytesToAddress([]byte{0xcc}),
	)
	t.Log("\n" + formatTransitions(accesses))
	t.Log("\n" + formatTransitions(transitions))
}

func BenchmarkExport(b *testing.B) {
	db := urlcommon.NewDataStore()
	meta, _ := commutative.NewMeta(urlcommon.NewPlatform().Eth10Account())
	db.Save(urlcommon.NewPlatform().Eth10Account(), meta)
	url := concurrenturl.NewConcurrentUrl(db)

	statedb := evm.NewStateDBV2(nil, db, url)
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

	b.Log("\n" + formatTransitions(transitions[:9]))
}

func TestExport1(b *testing.T) {
	db := urlcommon.NewDataStore()
	meta, _ := commutative.NewMeta(urlcommon.NewPlatform().Eth10Account())
	db.Save(urlcommon.NewPlatform().Eth10Account(), meta)

	begin := time.Now()
	var transitions []urlcommon.UnivalueInterface
	for i := 0; i < 1000; i++ {
		url := concurrenturl.NewConcurrentUrl(db)
		statedb := evm.NewStateDBV2(nil, db, url)
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
	b.Log("\n" + formatTransitions(transitions[:9]))

	// begin = time.Now()
	// url := concurrenturl.NewConcurrentUrl(db)
	// url.Commit(transitions, []uint32{1})
	// b.Log(time.Duration(time.Since(begin)))
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func formatValue(value interface{}) string {
	switch value.(type) {
	case *commutative.Meta:
		meta := value.(*commutative.Meta)
		var str string
		str += "{"
		for i, k := range meta.GetKeys() {
			str += k
			if i != len(meta.GetKeys())-1 {
				str += ", "
			}
		}
		str += "}"
		if len(meta.GetAdded()) != 0 {
			str += " + {"
			for i, k := range meta.GetAdded() {
				str += k
				if i != len(meta.GetAdded())-1 {
					str += ", "
				}
			}
			str += "}"
		}
		if len(meta.GetRemoved()) != 0 {
			str += " - {"
			for i, k := range meta.GetRemoved() {
				str += k
				if i != len(meta.GetRemoved())-1 {
					str += ", "
				}
			}
			str += "}"
		}
		return str
	case *noncommutative.Int64:
		return fmt.Sprintf(" = %v", int64(*value.(*noncommutative.Int64)))
	case *noncommutative.Bytes:
		return fmt.Sprintf(" = %v", value.(*noncommutative.Bytes).Data())
	case *commutative.Balance:
		v := value.(*commutative.Balance).Value()
		d := value.(*commutative.Balance).GetDelta()
		return fmt.Sprintf(" = %v + %v", v.(*big.Int).Uint64(), d.Int64())
	case *commutative.Int64:
		v := value.(*commutative.Int64).Value()
		d := value.(*commutative.Int64).GetDelta()
		return fmt.Sprintf(" = %v + %v", v, d)
	}
	return ""
}

func formatTransitions(transitions []urlcommon.UnivalueInterface) string {
	var str string
	for _, t := range transitions {
		str += fmt.Sprintf("[%v:%v,%v,%v,%v]%s%s\n", t.(*urltype.Univalue).GetTx(), t.(*urltype.Univalue).Reads(), t.(*urltype.Univalue).Writes(), t.(*urltype.Univalue).Preexist(), t.(*urltype.Univalue).Composite(), t.(*urltype.Univalue).GetPath(), formatValue(t.(*urltype.Univalue).Value()))
	}
	return str
}

func detectConflict(transitions []urlcommon.UnivalueInterface) ([]uint32, []uint32, []bool) {
	length := len(transitions)
	txs := make([]uint32, length)
	paths := make([]string, length)
	reads := make([]uint32, length)
	writes := make([]uint32, length)
	addOrDelete := make([]bool, length)
	composite := make([]bool, length)
	for i, t := range transitions {
		txs[i] = t.(*urltype.Univalue).GetTx()
		paths[i] = t.(*urltype.Univalue).GetPath()
		reads[i] = t.(*urltype.Univalue).Reads()
		writes[i] = t.(*urltype.Univalue).Writes()
		addOrDelete[i] = t.(*urltype.Univalue).IfAddOrDelete()
		composite[i] = t.(*urltype.Univalue).Composite()
	}

	engine := arbitrator.Start()
	_, buf := arbitrator.Insert(engine, txs, paths, reads, writes, addOrDelete, composite)
	txs, groups, flags := arbitrator.Detect(engine, uint32(length))
	arbitrator.Clear(engine, buf)
	return txs, groups, flags
}
