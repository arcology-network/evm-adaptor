package eth

import (
	"math/big"

	uint256 "github.com/holiman/uint256"

	codec "github.com/arcology-network/common-lib/codec"
	"github.com/arcology-network/concurrenturl"
	commutative "github.com/arcology-network/concurrenturl/commutative"
	noncommutative "github.com/arcology-network/concurrenturl/noncommutative"
	evmcommon "github.com/arcology-network/evm/common"
	evmtypes "github.com/arcology-network/evm/core/types"
	"github.com/arcology-network/evm/crypto"
)

// Arcology implementation of Eth ImplStateDB interfaces.
type ImplStateDB struct {
	refund uint64
	txHash evmcommon.Hash
	tid    uint32 // tx id
	logs   map[evmcommon.Hash][]*evmtypes.Log

	url *concurrenturl.ConcurrentUrl
}

func NewImplStateDB(url *concurrenturl.ConcurrentUrl) *ImplStateDB {
	return &ImplStateDB{
		logs: make(map[evmcommon.Hash][]*evmtypes.Log),
		url:  url,
	}
}

func (state *ImplStateDB) CreateAccount(addr evmcommon.Address) {
	createAccount(state.url, addr, state.tid)
}

func (state *ImplStateDB) SubBalance(addr evmcommon.Address, amount *big.Int) {
	state.AddBalance(addr, new(big.Int).Neg(amount))
}

func (state *ImplStateDB) AddBalance(addr evmcommon.Address, amount *big.Int) {
	if !accountExist(state.url, addr, state.tid) {
		createAccount(state.url, addr, state.tid)
	}

	if delta, ok := commutative.NewU256DeltaFromBigInt(amount); ok {
		if err := state.url.Write(state.tid, getBalancePath(state.url, addr), delta); err != nil {
			panic(err) //should not just panic
		}
		return
	}
	panic("Error: Failed to call AddBalance()")
}

func (state *ImplStateDB) GetBalance(addr evmcommon.Address) *big.Int {
	if !accountExist(state.url, addr, state.tid) {
		return new(big.Int)
	}

	if value, err := state.url.Read(state.tid, getBalancePath(state.url, addr)); err != nil {
		panic(err)
	} else {
		return (*(value.(*uint256.Int))).ToBig() // v.(*commutative.U256).Value().(*big.Int)
	}
}

func (state *ImplStateDB) PeekBalance(addr evmcommon.Address) *big.Int {
	if !accountExist(state.url, addr, state.tid) {
		return new(big.Int)
	}

	if value, err := state.url.Peek(getBalancePath(state.url, addr)); err != nil {
		panic(err)
	} else {
		// return value.(*urltype.Univalue).Value().(*commutative.U256).Value().(*big.Int)
		v, _, _ := value.(*commutative.U256).Get()
		return v.(*uint256.Int).ToBig()
	}
}

func (state *ImplStateDB) SetBalance(addr evmcommon.Address, amount *big.Int) {
	origin := state.PeekBalance(addr)
	state.AddBalance(addr, new(big.Int).Sub(amount, origin))
}

func (state *ImplStateDB) GetNonce(addr evmcommon.Address) uint64 {
	if !accountExist(state.url, addr, state.tid) {
		return 0
	}

	if value, err := state.url.Peek(getNoncePath(state.url, addr)); err != nil {
		panic(err)
	} else {
		v, _, _ := value.(*commutative.Uint64).Get()
		return v.(uint64)
	}
}

func (state *ImplStateDB) SetNonce(addr evmcommon.Address, nonce uint64) {
	if !accountExist(state.url, addr, state.tid) {
		createAccount(state.url, addr, state.tid)
	}

	if err := state.url.Write(state.tid, getNoncePath(state.url, addr), commutative.NewUint64Delta(1)); err != nil {
		panic(err)
	}
}

func (state *ImplStateDB) GetCodeHash(addr evmcommon.Address) evmcommon.Hash {
	code := state.GetCode(addr)
	if len(code) == 0 {
		return evmcommon.Hash{}
	} else {
		return crypto.Keccak256Hash(code)
	}
}

func (state *ImplStateDB) GetCode(addr evmcommon.Address) []byte {
	if !accountExist(state.url, addr, state.tid) {
		return nil
	}
	if value, err := state.url.Read(state.tid, getCodePath(state.url, addr)); err != nil {
		panic(err)
	} else {
		return value.([]byte)
	}
}

func (this *ImplStateDB) SetCode(addr evmcommon.Address, code []byte) {
	if !accountExist(this.url, addr, this.tid) {
		createAccount(this.url, addr, this.tid)
	}

	if err := this.url.Write(this.tid, getCodePath(this.url, addr), noncommutative.NewBytes(code)); err != nil {
		panic(err)
	}
}

func (this *ImplStateDB) GetCodeSize(addr evmcommon.Address) int {
	return len(this.GetCode(addr))
}

func (this *ImplStateDB) AddRefund(amount uint64) {
	this.refund += amount
}

func (this *ImplStateDB) SubRefund(amount uint64) {
	this.refund -= amount
}

func (this *ImplStateDB) GetRefund() uint64 {
	return this.refund
}

// Get from DB directly, bypassing ccurl since it make have some temporary states
func (this *ImplStateDB) GetCommittedState(addr evmcommon.Address, key evmcommon.Hash) evmcommon.Hash {
	if value, _ := this.url.ReadCommitted(this.tid, getStorageKeyPath(this.url, addr, key)); value == nil {
		return evmcommon.Hash{}
	} else {
		v, _, _ := value.(*noncommutative.Bytes).Get()
		return evmcommon.BytesToHash(v.([]byte))
	}
}

func (state *ImplStateDB) GetState(addr evmcommon.Address, key evmcommon.Hash) evmcommon.Hash {
	if value, err := state.url.Read(state.tid, getStorageKeyPath(state.url, addr, key)); err != nil {
		panic(err)
	} else if value == nil {
		return evmcommon.Hash{}
	} else {
		return evmcommon.BytesToHash(value.([]byte))
	}
}

func (this *ImplStateDB) SetState(addr evmcommon.Address, key, value evmcommon.Hash) {
	if !accountExist(this.url, addr, this.tid) {
		createAccount(this.url, addr, this.tid)
	}

	if err := this.url.Write(this.tid, getStorageKeyPath(this.url, addr, key), noncommutative.NewBytes(value.Bytes())); err != nil {
		panic(err)
	}
}

func (state *ImplStateDB) Suicide(addr evmcommon.Address) bool {
	return true
}

func (state *ImplStateDB) HasSuicided(addr evmcommon.Address) bool {
	return false
}

func (state *ImplStateDB) Exist(addr evmcommon.Address) bool {
	return accountExist(state.url, addr, state.tid)
}

func (state *ImplStateDB) Empty(addr evmcommon.Address) bool {
	if !accountExist(state.url, addr, state.tid) {
		return true
	}

	if value, err := state.url.Peek(getBalancePath(state.url, addr)); err != nil {
		panic(err)
	} else {
		v, _, _ := value.(*commutative.U256).Get()
		if v.(*uint256.Int).Cmp(commutative.U256_ZERO) != 0 {
			return false
		}
	}

	if value, err := state.url.Peek(getNoncePath(state.url, addr)); err != nil {
		panic(err)
	} else {
		v, _, _ := value.(*commutative.Uint64).Get()
		if v.(uint64) != 0 {
			return false
		}
	}

	if value, err := state.url.Peek(getCodePath(state.url, addr)); err != nil {
		panic(err)
	} else {
		return len(value.(*noncommutative.Bytes).Value().(codec.Bytes)) == 0
	}
}

func (state *ImplStateDB) RevertToSnapshot(id int) {
	// Do nothing.
}

func (state *ImplStateDB) Snapshot() int {
	return 0
}

func (state *ImplStateDB) AddLog(log *evmtypes.Log) {
	state.logs[state.txHash] = append(state.logs[state.txHash], log)
}

func (state *ImplStateDB) AddPreimage(hash evmcommon.Hash, preimage []byte) {
	// Do nothing.
}

func (state *ImplStateDB) ForEachStorage(addr evmcommon.Address, f func(evmcommon.Hash, evmcommon.Hash) bool) error {
	return nil
}

func (state *ImplStateDB) PrepareAccessList(sender evmcommon.Address, dest *evmcommon.Address, precompiles []evmcommon.Address, txAccesses evmtypes.AccessList) {
	// Do nothing.
}

func (state *ImplStateDB) AddressInAccessList(addr evmcommon.Address) bool {
	return true
}

func (state *ImplStateDB) SlotInAccessList(addr evmcommon.Address, slot evmcommon.Hash) (addressOk bool, slotOk bool) {
	return true, true
}

func (state *ImplStateDB) AddAddressToAccessList(addr evmcommon.Address) {
	// Do nothing.
}

func (state *ImplStateDB) AddSlotToAccessList(addr evmcommon.Address, slot evmcommon.Hash) {
	// Do nothing.
}

func (state *ImplStateDB) Prepare(txHash, bhash evmcommon.Hash, ti int) {
	state.refund = 0
	state.txHash = txHash
	state.tid = uint32(ti)
	state.logs = make(map[evmcommon.Hash][]*evmtypes.Log)
}

func (state *ImplStateDB) GetLogs(hash evmcommon.Hash) []*evmtypes.Log {
	return state.logs[hash]
}

func (state *ImplStateDB) Copy() *ImplStateDB {
	return &ImplStateDB{
		logs: make(map[evmcommon.Hash][]*evmtypes.Log),
		// kapi: state.kapi,
		// db:   state.db,
		url: state.url,
	}
}

func (state *ImplStateDB) Set(eac EthAccountCache, esc EthStorageCache) {
	// TODO
}

func GetAccountPathSize(url *concurrenturl.ConcurrentUrl) int {
	return len(url.Platform.Eth10()) + 2*evmcommon.AddressLength
}

// func ExportOnFailure(
// 	db urlcommon.DatastoreInterface,
// 	txIndex int,
// 	from evmcommon.Address,
// 	coinbase evmcommon.Address,
// 	gasUsed uint64,
// 	gasPrice *big.Int,
// ) ([]urlcommon.UnivalueInterface, []urlcommon.UnivalueInterface) {
// 	url := concurrenturl.NewConcurrentUrl(db)
// 	state := NewImplStateDB(nil, db, url)
// 	state.Prepare(evmcommon.Hash{}, evmcommon.Hash{}, txIndex)

// 	// Chage the tranasction fees
// 	state.AddBalance(coinbase, new(big.Int).Mul(new(big.Int).SetUint64(gasUsed), gasPrice))
// 	state.SubBalance(from, new(big.Int).Mul(new(big.Int).SetUint64(gasUsed), gasPrice))
// 	//	state.SetNonce(from, 0) don't increase if the transcation failes.
// 	return url.Export(false)
// }

// func ExportOnFailureEx(
// 	db urlcommon.DatastoreInterface,
// 	txIndex int,
// 	from, coinbase evmcommon.Address,
// 	gasUsed uint64, gasPrice *big.Int,
// ) ([][]byte, [][]byte) {
// 	url := concurrenturl.NewConcurrentUrl(db)
// 	state := NewImplStateDB(nil, db, url)
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
// 	state := NewImplStateDB(nil, db, url)
// 	state.Prepare(evmcommon.Hash{}, evmcommon.Hash{}, txIndex)
// 	state.SetNonce(from, 0)
// 	return url.ExportAll()
// }

// func ExportOnConflictionEx(
// 	db urlcommon.DatastoreInterface,
// 	txIndex int,
// 	from evmcommon.Address,
// ) ([][]byte, [][]byte) {
// 	url := concurrenturl.NewConcurrentUrl(db)
// 	state := NewImplStateDB(nil, db, url)
// 	state.Prepare(evmcommon.Hash{}, evmcommon.Hash{}, txIndex)
// 	state.SetNonce(from, 0)
// 	return url.ExportEncoded(func(accesses, transitions []urlcommon.UnivalueInterface) ([]urlcommon.UnivalueInterface, []urlcommon.UnivalueInterface) {
// 		for _, t := range transitions {
// 			t.SetTransitionType(urlcommon.INVARIATE_TRANSITIONS)
// 		}
// 		return accesses, transitions
// 	})
// }
