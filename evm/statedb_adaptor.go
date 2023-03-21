package evm

import (
	"math/big"

	uint256 "github.com/holiman/uint256"

	"github.com/arcology-network/concurrenturl/v2"
	urlcommon "github.com/arcology-network/concurrenturl/v2/common"
	commutative "github.com/arcology-network/concurrenturl/v2/type/commutative"
	noncommutative "github.com/arcology-network/concurrenturl/v2/type/noncommutative"
	evmcommon "github.com/arcology-network/evm/common"
	evmtypes "github.com/arcology-network/evm/core/types"
	"github.com/arcology-network/evm/crypto"
)

type ethState struct {
	refund uint64
	thash  evmcommon.Hash
	tid    uint32
	logs   map[evmcommon.Hash][]*evmtypes.Log

	kapi *API
	db   urlcommon.DatastoreInterface
	url  *concurrenturl.ConcurrentUrl
}

func NewStateDB(kapi *API, db urlcommon.DatastoreInterface, url *concurrenturl.ConcurrentUrl) StateDB {
	return &ethState{
		logs: make(map[evmcommon.Hash][]*evmtypes.Log),
		kapi: kapi,
		db:   db,
		url:  url,
	}
}

func (state *ethState) CreateAccount(addr evmcommon.Address) {
	createAccount(state.url, addr, state.tid)
}

func (state *ethState) SubBalance(addr evmcommon.Address, amount *big.Int) {
	state.AddBalance(addr, new(big.Int).Neg(amount))
}

func (state *ethState) AddBalance(addr evmcommon.Address, amount *big.Int) {
	if !accountExist(state.url, addr, state.tid) {
		createAccount(state.url, addr, state.tid)
	}

	if err := state.url.Write(state.tid, getBalancePath(state.url, addr), commutative.NewBalance(nil, amount)); err != nil {
		panic(err) //should not just panic
	}
}

func (state *ethState) GetBalance(addr evmcommon.Address) *big.Int {
	if !accountExist(state.url, addr, state.tid) {
		return new(big.Int)
	}

	if value, err := state.url.Read(state.tid, getBalancePath(state.url, addr)); err != nil {
		panic(err)
	} else {
		v, _, _ := value.(*commutative.Balance).Get(0, "", nil)
		return (*(v.(*commutative.Balance).Value().(*uint256.Int))).ToBig() // v.(*commutative.Balance).Value().(*big.Int)
	}
}

func (state *ethState) PeekBalance(addr evmcommon.Address) *big.Int {
	if !accountExist(state.url, addr, state.tid) {
		return new(big.Int)
	}

	if value, err := state.url.TryRead(urlcommon.SYSTEM, getBalancePath(state.url, addr)); err != nil {
		panic(err)
	} else {
		// return value.(*urltype.Univalue).Value().(*commutative.Balance).Value().(*big.Int)
		return value.(*commutative.Balance).Value().(*uint256.Int).ToBig()
	}
}

func (state *ethState) SetBalance(addr evmcommon.Address, amount *big.Int) {
	origin := state.PeekBalance(addr)
	state.AddBalance(addr, new(big.Int).Sub(amount, origin))
}

func (state *ethState) GetNonce(addr evmcommon.Address) uint64 {
	if !accountExist(state.url, addr, state.tid) {
		return 0
	}

	if value, err := state.url.TryRead(urlcommon.SYSTEM, getNoncePath(state.url, addr)); err != nil {
		panic(err)
	} else {
		return uint64(value.(*commutative.Int64).Value().(int64))
	}
}

func (state *ethState) SetNonce(addr evmcommon.Address, nonce uint64) {
	if !accountExist(state.url, addr, state.tid) {
		createAccount(state.url, addr, state.tid)
	}

	if err := state.url.Write(state.tid, getNoncePath(state.url, addr), commutative.NewInt64(0, 1)); err != nil {
		panic(err)
	}
}

func (state *ethState) GetCodeHash(addr evmcommon.Address) evmcommon.Hash {
	code := state.GetCode(addr)
	if len(code) == 0 {
		return evmcommon.Hash{}
	} else {
		return crypto.Keccak256Hash(code)
	}
}

func (state *ethState) GetCode(addr evmcommon.Address) []byte {
	if !accountExist(state.url, addr, state.tid) {
		return nil
	}
	if value, err := state.url.Read(state.tid, getCodePath(state.url, addr)); err != nil {
		panic(err)
	} else {
		return value.(*noncommutative.Bytes).Data()

	}
}

func (state *ethState) SetCode(addr evmcommon.Address, code []byte) {
	if !accountExist(state.url, addr, state.tid) {
		createAccount(state.url, addr, state.tid)
	}

	if err := state.url.Write(state.tid, getCodePath(state.url, addr), noncommutative.NewBytes(code)); err != nil {
		panic(err)
	}
}

func (state *ethState) GetCodeSize(addr evmcommon.Address) int {
	if state.kapi.IsKernelAPI(addr) {
		// FIXME!
		return 0xff
	}
	return len(state.GetCode(addr))
}

func (state *ethState) AddRefund(amount uint64) {
	state.refund += amount
}

func (state *ethState) SubRefund(amount uint64) {
	state.refund -= amount
}

func (state *ethState) GetRefund() uint64 {
	return state.refund
}

func (state *ethState) GetCommittedState(addr evmcommon.Address, key evmcommon.Hash) evmcommon.Hash {
	if value, _ := state.db.Retrive(getStorageKeyPath(state.url, addr, key)); value == nil {
		return evmcommon.Hash{}
	} else {
		return evmcommon.BytesToHash(value.(*noncommutative.Bytes).Data())
	}
}

func (state *ethState) GetState(addr evmcommon.Address, key evmcommon.Hash) evmcommon.Hash {
	if value, err := state.url.Read(state.tid, getStorageKeyPath(state.url, addr, key)); err != nil {
		panic(err)
	} else if value == nil {
		return evmcommon.Hash{}
	} else {
		return evmcommon.BytesToHash(value.(*noncommutative.Bytes).Data())
	}
}

func (state *ethState) SetState(addr evmcommon.Address, key, value evmcommon.Hash) {
	if !accountExist(state.url, addr, state.tid) {
		createAccount(state.url, addr, state.tid)
	}

	if err := state.url.Write(state.tid, getStorageKeyPath(state.url, addr, key), noncommutative.NewBytes(value.Bytes())); err != nil {
		panic(err)
	}
}

func (state *ethState) Suicide(addr evmcommon.Address) bool {
	return true
}

func (state *ethState) HasSuicided(addr evmcommon.Address) bool {
	return false
}

func (state *ethState) Exist(addr evmcommon.Address) bool {
	return accountExist(state.url, addr, state.tid)
}

func (state *ethState) Empty(addr evmcommon.Address) bool {
	if !accountExist(state.url, addr, state.tid) {
		return true
	}

	if value, err := state.url.TryRead(urlcommon.SYSTEM, getBalancePath(state.url, addr)); err != nil {
		panic(err)
	} else {
		if value.(*commutative.Balance).Value().(*big.Int).Cmp(new(big.Int).SetInt64(0)) != 0 {
			return false
		}
	}

	if value, err := state.url.TryRead(urlcommon.SYSTEM, getNoncePath(state.url, addr)); err != nil {
		panic(err)
	} else {
		if value.(*commutative.Int64).Value().(int64) != 0 {
			return false
		}
	}

	if value, err := state.url.TryRead(urlcommon.SYSTEM, getCodePath(state.url, addr)); err != nil {
		panic(err)
	} else {
		return len(value.(*noncommutative.Bytes).Data()) == 0
	}
}

func (state *ethState) RevertToSnapshot(id int) {
	// Do nothing.
}

func (state *ethState) Snapshot() int {
	return 0
}

func (state *ethState) AddLog(log *evmtypes.Log) {
	state.logs[state.thash] = append(state.logs[state.thash], log)
}

func (state *ethState) AddPreimage(hash evmcommon.Hash, preimage []byte) {
	// Do nothing.
}

func (state *ethState) ForEachStorage(addr evmcommon.Address, f func(evmcommon.Hash, evmcommon.Hash) bool) error {
	return nil
}

func (state *ethState) PrepareAccessList(sender evmcommon.Address, dest *evmcommon.Address, precompiles []evmcommon.Address, txAccesses evmtypes.AccessList) {
	// Do nothing.
}

func (state *ethState) AddressInAccessList(addr evmcommon.Address) bool {
	return true
}

func (state *ethState) SlotInAccessList(addr evmcommon.Address, slot evmcommon.Hash) (addressOk bool, slotOk bool) {
	return true, true
}

func (state *ethState) AddAddressToAccessList(addr evmcommon.Address) {
	// Do nothing.
}

func (state *ethState) AddSlotToAccessList(addr evmcommon.Address, slot evmcommon.Hash) {
	// Do nothing.
}

func (state *ethState) Prepare(thash, bhash evmcommon.Hash, ti int) {
	state.refund = 0
	state.thash = thash
	state.tid = uint32(ti)
	state.logs = make(map[evmcommon.Hash][]*evmtypes.Log)
}

func (state *ethState) GetLogs(hash evmcommon.Hash) []*evmtypes.Log {
	return state.logs[hash]
}

func (state *ethState) Copy() StateDB {
	return &ethState{
		logs: make(map[evmcommon.Hash][]*evmtypes.Log),
		kapi: state.kapi,
		db:   state.db,
		url:  state.url,
	}
}

func (state *ethState) Set(eac EthAccountCache, esc EthStorageCache) {
	// TODO
}

func GetAccountPathSize(url *concurrenturl.ConcurrentUrl) int {
	return len(url.Platform.Eth10()) + 2*evmcommon.AddressLength
}

func ExportOnFailure(
	db urlcommon.DatastoreInterface,
	txIndex int,
	from evmcommon.Address,
	coinbase evmcommon.Address,
	gasUsed uint64,
	gasPrice *big.Int,
) ([]urlcommon.UnivalueInterface, []urlcommon.UnivalueInterface) {
	url := concurrenturl.NewConcurrentUrl(db)
	state := NewStateDB(nil, db, url)
	state.Prepare(evmcommon.Hash{}, evmcommon.Hash{}, txIndex)

	// Chage the tranasction fees
	state.AddBalance(coinbase, new(big.Int).Mul(new(big.Int).SetUint64(gasUsed), gasPrice))
	state.SubBalance(from, new(big.Int).Mul(new(big.Int).SetUint64(gasUsed), gasPrice))
	//	state.SetNonce(from, 0) don't increase if the transcation failes.
	return url.Export(false)
}

// func ExportOnFailureEx(
// 	db urlcommon.DatastoreInterface,
// 	txIndex int,
// 	from, coinbase evmcommon.Address,
// 	gasUsed uint64, gasPrice *big.Int,
// ) ([][]byte, [][]byte) {
// 	url := concurrenturl.NewConcurrentUrl(db)
// 	state := NewStateDB(nil, db, url)
// 	state.Prepare(evmcommon.Hash{}, evmcommon.Hash{}, txIndex)
// 	state.AddBalance(coinbase, new(big.Int).Mul(new(big.Int).SetUint64(gasUsed), gasPrice))
// 	state.SubBalance(from, new(big.Int).Mul(new(big.Int).SetUint64(gasUsed), gasPrice))
// 	// state.SetNonce(from, 0)
// 	return url.ExportEncoded(nil)
// }

// func ExportOnConfliction(
// 	db urlcommon.DatastoreInterface,
// 	txIndex int,
// 	from evmcommon.Address,
// ) ([]urlcommon.UnivalueInterface, []urlcommon.UnivalueInterface) {
// 	url := concurrenturl.NewConcurrentUrl(db)
// 	state := NewStateDB(nil, db, url)
// 	state.Prepare(evmcommon.Hash{}, evmcommon.Hash{}, txIndex)
// 	state.SetNonce(from, 0)
// 	return url.Export(true)
// }

// func ExportOnConflictionEx(
// 	db urlcommon.DatastoreInterface,
// 	txIndex int,
// 	from evmcommon.Address,
// ) ([][]byte, [][]byte) {
// 	url := concurrenturl.NewConcurrentUrl(db)
// 	state := NewStateDB(nil, db, url)
// 	state.Prepare(evmcommon.Hash{}, evmcommon.Hash{}, txIndex)
// 	state.SetNonce(from, 0)
// 	return url.ExportEncoded(func(accesses, transitions []urlcommon.UnivalueInterface) ([]urlcommon.UnivalueInterface, []urlcommon.UnivalueInterface) {
// 		for _, t := range transitions {
// 			t.SetTransitionType(urlcommon.INVARIATE_TRANSITIONS)
// 		}
// 		return accesses, transitions
// 	})
// }
