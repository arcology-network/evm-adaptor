package evm

import (
	"encoding/hex"
	"math/big"

	"github.com/arcology-network/concurrenturl/v2"
	urlcommon "github.com/arcology-network/concurrenturl/v2/common"

	// urltype "github.com/arcology-network/concurrenturl/v2/type"
	commonlib "github.com/arcology-network/common-lib/common"
	commutative "github.com/arcology-network/concurrenturl/v2/type/commutative"
	noncommutative "github.com/arcology-network/concurrenturl/v2/type/noncommutative"
	evmcommon "github.com/arcology-network/evm/common"
	evmtypes "github.com/arcology-network/evm/core/types"
	"github.com/arcology-network/evm/crypto"
)

type ethStateV2 struct {
	refund uint64
	thash  evmcommon.Hash
	tid    uint32
	logs   map[evmcommon.Hash][]*evmtypes.Log

	kapi *APIV2
	db   urlcommon.DatastoreInterface
	url  *concurrenturl.ConcurrentUrl
}

func NewStateDBV2(kapi *APIV2, db urlcommon.DatastoreInterface, url *concurrenturl.ConcurrentUrl) StateDB {
	return &ethStateV2{
		logs: make(map[evmcommon.Hash][]*evmtypes.Log),
		kapi: kapi,
		db:   db,
		url:  url,
	}
}

func (state *ethStateV2) CreateAccount(addr evmcommon.Address) {
	createAccount(state.url, addr, state.tid)
}

func (state *ethStateV2) SubBalance(addr evmcommon.Address, amount *big.Int) {
	state.AddBalance(addr, new(big.Int).Neg(amount))
}

func (state *ethStateV2) AddBalance(addr evmcommon.Address, amount *big.Int) {
	if !accountExist(state.url, addr, state.tid) {
		createAccount(state.url, addr, state.tid)
	}

	if err := state.url.Write(state.tid, getBalancePath(state.url, addr), commutative.NewBalance(nil, amount)); err != nil {
		panic(err)
	}
}

func (state *ethStateV2) GetBalance(addr evmcommon.Address) *big.Int {
	if !accountExist(state.url, addr, state.tid) {
		return new(big.Int)
	}

	if value, err := state.url.Read(state.tid, getBalancePath(state.url, addr)); err != nil {
		panic(err)
	} else {
		v, _, _ := value.(*commutative.Balance).Get(0, "", nil)
		return v.(*commutative.Balance).Value().(*big.Int)
	}
}

func (state *ethStateV2) GetBalanceNoRecord(addr evmcommon.Address) *big.Int {
	if !accountExist(state.url, addr, state.tid) {
		return new(big.Int)
	}

	if value, err := state.url.TryRead(urlcommon.SYSTEM, getBalancePath(state.url, addr)); err != nil {
		panic(err)
	} else {
		// return value.(*urltype.Univalue).Value().(*commutative.Balance).Value().(*big.Int)
		return value.(*commutative.Balance).Value().(*big.Int)
	}
}

func (state *ethStateV2) SetBalance(addr evmcommon.Address, amount *big.Int) {
	origin := state.GetBalanceNoRecord(addr)
	state.AddBalance(addr, new(big.Int).Sub(amount, origin))
}

func (state *ethStateV2) GetNonce(addr evmcommon.Address) uint64 {
	if !accountExist(state.url, addr, state.tid) {
		return 0
	}

	if value, err := state.url.TryRead(urlcommon.SYSTEM, getNoncePath(state.url, addr)); err != nil {
		panic(err)
	} else {
		return uint64(value.(*commutative.Int64).Value().(int64))
	}
}

func (state *ethStateV2) SetNonce(addr evmcommon.Address, nonce uint64) {
	if !accountExist(state.url, addr, state.tid) {
		createAccount(state.url, addr, state.tid)
	}

	if err := state.url.Write(state.tid, getNoncePath(state.url, addr), commutative.NewInt64(0, 1)); err != nil {
		panic(err)
	}
}

func (state *ethStateV2) GetCodeHash(addr evmcommon.Address) evmcommon.Hash {
	code := state.GetCode(addr)
	if len(code) == 0 {
		return evmcommon.Hash{}
	} else {
		return crypto.Keccak256Hash(code)
	}
}

func (state *ethStateV2) GetCode(addr evmcommon.Address) []byte {
	if !accountExist(state.url, addr, state.tid) {
		return nil
	}
	if value, err := state.url.Read(state.tid, getCodePath(state.url, addr)); err != nil {
		panic(err)
	} else {
		return value.(*noncommutative.Bytes).Data()

	}
}

func (state *ethStateV2) SetCode(addr evmcommon.Address, code []byte) {
	if !accountExist(state.url, addr, state.tid) {
		createAccount(state.url, addr, state.tid)
	}

	if err := state.url.Write(state.tid, getCodePath(state.url, addr), noncommutative.NewBytes(code)); err != nil {
		panic(err)
	}
}

func (state *ethStateV2) GetCodeSize(addr evmcommon.Address) int {
	if state.kapi.IsKernelAPI(addr) {
		// FIXME!
		return 0xff
	}
	return len(state.GetCode(addr))
}

func (state *ethStateV2) AddRefund(amount uint64) {
	state.refund += amount
}

func (state *ethStateV2) SubRefund(amount uint64) {
	state.refund -= amount
}

func (state *ethStateV2) GetRefund() uint64 {
	return state.refund
}

func (state *ethStateV2) GetCommittedState(addr evmcommon.Address, key evmcommon.Hash) evmcommon.Hash {
	if value := state.db.Retrive(getStorageKeyPath(state.url, addr, key)); value == nil {
		return evmcommon.Hash{}
	} else {
		return evmcommon.BytesToHash(value.(*noncommutative.Bytes).Data())
	}
}

func (state *ethStateV2) GetState(addr evmcommon.Address, key evmcommon.Hash) evmcommon.Hash {
	if value, err := state.url.Read(state.tid, getStorageKeyPath(state.url, addr, key)); err != nil {
		panic(err)
	} else if value == nil {
		return evmcommon.Hash{}
	} else {
		return evmcommon.BytesToHash(value.(*noncommutative.Bytes).Data())
	}
}

func (state *ethStateV2) SetState(addr evmcommon.Address, key, value evmcommon.Hash) {
	if !accountExist(state.url, addr, state.tid) {
		createAccount(state.url, addr, state.tid)
	}

	if err := state.url.Write(state.tid, getStorageKeyPath(state.url, addr, key), noncommutative.NewBytes(value.Bytes())); err != nil {
		panic(err)
	}
}

func (state *ethStateV2) Suicide(addr evmcommon.Address) bool {
	return true
}

func (state *ethStateV2) HasSuicided(addr evmcommon.Address) bool {
	return false
}

func (state *ethStateV2) Exist(addr evmcommon.Address) bool {
	return accountExist(state.url, addr, state.tid)
}

func (state *ethStateV2) Empty(addr evmcommon.Address) bool {
	if !accountExist(state.url, addr, state.tid) {
		return true
	}

	if value, err := state.url.Read(urlcommon.SYSTEM, getBalancePath(state.url, addr)); err != nil {
		panic(err)
	} else {
		v, _, _ := value.(*commutative.Balance).Get(0, "", nil)
		if v.(*commutative.Balance).Value().(*big.Int).Cmp(new(big.Int).SetInt64(0)) != 0 {
			return false
		}
	}

	if value, err := state.url.Read(urlcommon.SYSTEM, getNoncePath(state.url, addr)); err != nil {
		panic(err)
	} else {
		if *value.(*noncommutative.Int64) != 0 {
			return false
		}
	}

	if value, err := state.url.Read(urlcommon.SYSTEM, getCodePath(state.url, addr)); err != nil {
		panic(err)
	} else {
		return len(value.(*noncommutative.Bytes).Data()) == 0
	}
}

func (state *ethStateV2) RevertToSnapshot(id int) {
	// Do nothing.
}

func (state *ethStateV2) Snapshot() int {
	return 0
}

func (state *ethStateV2) AddLog(log *evmtypes.Log) {
	state.logs[state.thash] = append(state.logs[state.thash], log)
}

func (state *ethStateV2) AddPreimage(hash evmcommon.Hash, preimage []byte) {
	// Do nothing.
}

func (state *ethStateV2) ForEachStorage(addr evmcommon.Address, f func(evmcommon.Hash, evmcommon.Hash) bool) error {
	return nil
}

func (state *ethStateV2) PrepareAccessList(sender evmcommon.Address, dest *evmcommon.Address, precompiles []evmcommon.Address, txAccesses evmtypes.AccessList) {
	// Do nothing.
}

func (state *ethStateV2) AddressInAccessList(addr evmcommon.Address) bool {
	return true
}

func (state *ethStateV2) SlotInAccessList(addr evmcommon.Address, slot evmcommon.Hash) (addressOk bool, slotOk bool) {
	return true, true
}

func (state *ethStateV2) AddAddressToAccessList(addr evmcommon.Address) {
	// Do nothing.
}

func (state *ethStateV2) AddSlotToAccessList(addr evmcommon.Address, slot evmcommon.Hash) {
	// Do nothing.
}

func (state *ethStateV2) Prepare(thash, bhash evmcommon.Hash, ti int) {
	state.refund = 0
	state.thash = thash
	state.tid = uint32(ti)
	state.logs = make(map[evmcommon.Hash][]*evmtypes.Log)
}

func (state *ethStateV2) GetLogs(hash evmcommon.Hash) []*evmtypes.Log {
	return state.logs[hash]
}

func (state *ethStateV2) Copy() StateDB {
	return &ethStateV2{
		logs: make(map[evmcommon.Hash][]*evmtypes.Log),
		kapi: state.kapi,
		db:   state.db,
		url:  state.url,
	}
}

func (state *ethStateV2) Set(eac EthAccountCache, esc EthStorageCache) {
	// TODO
}

func GetAccountPathSize(url *concurrenturl.ConcurrentUrl) int {
	return len(url.Platform.Eth10()) + 2*evmcommon.AddressLength
}

func ExportOnFailure(
	db urlcommon.DatastoreInterface,
	txIndex int,
	from, coinbase evmcommon.Address,
	gasUsed uint64, gasPrice *big.Int,
) ([]urlcommon.UnivalueInterface, []urlcommon.UnivalueInterface) {
	url := concurrenturl.NewConcurrentUrl(db)
	state := NewStateDBV2(nil, db, url)
	state.Prepare(evmcommon.Hash{}, evmcommon.Hash{}, txIndex)
	state.AddBalance(coinbase, new(big.Int).Mul(new(big.Int).SetUint64(gasUsed), gasPrice))
	state.SubBalance(from, new(big.Int).Mul(new(big.Int).SetUint64(gasUsed), gasPrice))
	state.SetNonce(from, 0)
	return url.Export(true)
}

func ExportOnConfliction(
	db urlcommon.DatastoreInterface,
	txIndex int,
	from evmcommon.Address,
) ([]urlcommon.UnivalueInterface, []urlcommon.UnivalueInterface) {
	url := concurrenturl.NewConcurrentUrl(db)
	state := NewStateDBV2(nil, db, url)
	state.Prepare(evmcommon.Hash{}, evmcommon.Hash{}, txIndex)
	state.SetNonce(from, 0)
	return url.Export(true)
}

func ExportOnFailureEx(
	db urlcommon.DatastoreInterface,
	txIndex int,
	from, coinbase evmcommon.Address,
	gasUsed uint64, gasPrice *big.Int,
) ([][]byte, [][]byte) {
	url := concurrenturl.NewConcurrentUrl(db)
	state := NewStateDBV2(nil, db, url)
	state.Prepare(evmcommon.Hash{}, evmcommon.Hash{}, txIndex)
	state.AddBalance(coinbase, new(big.Int).Mul(new(big.Int).SetUint64(gasUsed), gasPrice))
	state.SubBalance(from, new(big.Int).Mul(new(big.Int).SetUint64(gasUsed), gasPrice))
	state.SetNonce(from, 0)
	return url.ExportEncoded()
}

func ExportOnConflictionEx(
	db urlcommon.DatastoreInterface,
	txIndex int,
	from evmcommon.Address,
) ([][]byte, [][]byte) {
	url := concurrenturl.NewConcurrentUrl(db)
	state := NewStateDBV2(nil, db, url)
	state.Prepare(evmcommon.Hash{}, evmcommon.Hash{}, txIndex)
	state.SetNonce(from, 0)
	return url.ExportEncoded()
}

func addressToHex(addr evmcommon.Address) string {
	var accHex [2 * evmcommon.AddressLength]byte
	hex.Encode(accHex[:], addr[:])
	return string(accHex[:])
}

func getAccountRootPath(url *concurrenturl.ConcurrentUrl, account evmcommon.Address) string {
	var accHex [2 * evmcommon.AddressLength]byte
	hex.Encode(accHex[:], account[:])
	return commonlib.StrCat(url.Platform.Eth10Account(), string(accHex[:]), "/")
}

func getStorageRootPath(url *concurrenturl.ConcurrentUrl, account evmcommon.Address) string {
	var accHex [2 * evmcommon.AddressLength]byte
	hex.Encode(accHex[:], account[:])
	return commonlib.StrCat(url.Platform.Eth10Account(), string(accHex[:]), "/storage/native/")
}

func getStorageKeyPath(url *concurrenturl.ConcurrentUrl, account evmcommon.Address, key evmcommon.Hash) string {
	var accHex [2 * evmcommon.AddressLength]byte
	hex.Encode(accHex[:], account[:])
	return commonlib.StrCat(url.Platform.Eth10Account(), string(accHex[:]), "/storage/native/", key.Hex())
}

func getBalancePath(url *concurrenturl.ConcurrentUrl, account evmcommon.Address) string {
	var accHex [2 * evmcommon.AddressLength]byte
	hex.Encode(accHex[:], account[:])
	return commonlib.StrCat(url.Platform.Eth10Account(), string(accHex[:]), "/balance")
}

func getNoncePath(url *concurrenturl.ConcurrentUrl, account evmcommon.Address) string {
	var accHex [2 * evmcommon.AddressLength]byte
	hex.Encode(accHex[:], account[:])
	return commonlib.StrCat(url.Platform.Eth10Account(), string(accHex[:]), "/nonce")
}

func getCodePath(url *concurrenturl.ConcurrentUrl, account evmcommon.Address) string {
	var accHex [2 * evmcommon.AddressLength]byte
	hex.Encode(accHex[:], account[:])
	return commonlib.StrCat(url.Platform.Eth10Account(), string(accHex[:]), "/code")
}

func accountExist(url *concurrenturl.ConcurrentUrl, account evmcommon.Address, tid uint32) bool {
	return url.IfExists(getAccountRootPath(url, account))
}

func createAccount(url *concurrenturl.ConcurrentUrl, account evmcommon.Address, tid uint32) {
	if err := url.CreateAccount(tid, url.Platform.Eth10(), addressToHex(account)); err != nil {
		panic(err)
	}

	if err := url.Write(tid, getBalancePath(url, account), commutative.NewBalance(new(big.Int), new(big.Int))); err != nil {
		panic(err)
	}
	if err := url.Write(tid, getNoncePath(url, account), commutative.NewInt64(0, 0)); err != nil {
		panic(err)
	}
	// if err := url.Write(tid, getCodePath(url, account), noncommutative.NewBytes(nil)); err != nil {
	// 	panic(err)
	// }
}
