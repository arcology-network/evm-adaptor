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

	eucommon "github.com/arcology-network/vm-adaptor/common"
)

// Arcology implementation of Eth ImplStateDB interfaces.
type ImplStateDB struct {
	refund uint64
	txHash evmcommon.Hash
	tid    uint32 // tx id
	logs   map[evmcommon.Hash][]*evmtypes.Log

	url *concurrenturl.ConcurrentUrl
	api eucommon.ConcurrentApiRouterInterface
}

func NewImplStateDB(api eucommon.ConcurrentApiRouterInterface) *ImplStateDB {
	return &ImplStateDB{
		logs: make(map[evmcommon.Hash][]*evmtypes.Log),
		url:  api.Ccurl(),
		api:  api,
	}
}

func (this *ImplStateDB) CreateAccount(addr evmcommon.Address) {
	createAccount(this.url, addr, this.tid)
}

func (this *ImplStateDB) SubBalance(addr evmcommon.Address, amount *big.Int) {
	this.AddBalance(addr, new(big.Int).Neg(amount))
}

func (this *ImplStateDB) AddBalance(addr evmcommon.Address, amount *big.Int) {
	if !accountExist(this.url, addr, this.tid) {
		createAccount(this.url, addr, this.tid)
	}

	if delta, ok := commutative.NewU256DeltaFromBigInt(amount); ok {
		if err := this.url.Write(this.tid, getBalancePath(this.url, addr), delta); err != nil {
			panic(err) //should not just panic
		}
		return
	}
	panic("Error: Failed to call AddBalance()")
}

func (this *ImplStateDB) GetBalance(addr evmcommon.Address) *big.Int {
	if !accountExist(this.url, addr, this.tid) {
		return new(big.Int)
	}

	if value, err := this.url.Read(this.tid, getBalancePath(this.url, addr)); err != nil {
		panic(err)
	} else {
		return (*(value.(*uint256.Int))).ToBig() // v.(*commutative.U256).Value().(*big.Int)
	}
}

func (this *ImplStateDB) PeekBalance(addr evmcommon.Address) *big.Int {
	if !accountExist(this.url, addr, this.tid) {
		return new(big.Int)
	}

	if value, err := this.url.Peek(getBalancePath(this.url, addr)); err != nil {
		panic(err)
	} else {
		// return value.(*urltype.Univalue).Value().(*commutative.U256).Value().(*big.Int)
		v, _, _ := value.(*commutative.U256).Get()
		return v.(*uint256.Int).ToBig()
	}
}

func (this *ImplStateDB) SetBalance(addr evmcommon.Address, amount *big.Int) {
	origin := this.PeekBalance(addr)
	this.AddBalance(addr, new(big.Int).Sub(amount, origin))
}

func (this *ImplStateDB) GetNonce(addr evmcommon.Address) uint64 {
	if !accountExist(this.url, addr, this.tid) {
		return 0
	}

	if value, err := this.url.Peek(getNoncePath(this.url, addr)); err != nil {
		panic(err)
	} else {
		v, _, _ := value.(*commutative.Uint64).Get()
		return v.(uint64)
	}
}

func (this *ImplStateDB) SetNonce(addr evmcommon.Address, nonce uint64) {
	if !accountExist(this.url, addr, this.tid) {
		createAccount(this.url, addr, this.tid)
	}

	if err := this.url.Write(this.tid, getNoncePath(this.url, addr), commutative.NewUint64Delta(1)); err != nil {
		panic(err)
	}
}

func (this *ImplStateDB) GetCodeHash(addr evmcommon.Address) evmcommon.Hash {
	code := this.GetCode(addr)
	if len(code) == 0 {
		return evmcommon.Hash{}
	} else {
		return crypto.Keccak256Hash(code)
	}
}

func (this *ImplStateDB) GetCode(addr evmcommon.Address) []byte {
	if !accountExist(this.url, addr, this.tid) {
		return nil
	}
	if value, err := this.url.Read(this.tid, getCodePath(this.url, addr)); err != nil {
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

func (this *ImplStateDB) GetCodeSize(addr evmcommon.Address) int                          { return len(this.GetCode(addr)) }
func (this *ImplStateDB) AddRefund(amount uint64)                                         { this.refund += amount }
func (this *ImplStateDB) SubRefund(amount uint64)                                         { this.refund -= amount }
func (this *ImplStateDB) GetRefund() uint64                                               { return this.refund }
func (this *ImplStateDB) Suicide(addr evmcommon.Address) bool                             { return true }
func (this *ImplStateDB) HasSuicided(addr evmcommon.Address) bool                         { return false }
func (this *ImplStateDB) RevertToSnapshot(id int)                                         {}
func (this *ImplStateDB) Snapshot() int                                                   { return 0 }
func (this *ImplStateDB) AddPreimage(hash evmcommon.Hash, preimage []byte)                {}
func (this *ImplStateDB) AddAddressToAccessList(addr evmcommon.Address)                   {} // Do nothing.
func (this *ImplStateDB) AddSlotToAccessList(addr evmcommon.Address, slot evmcommon.Hash) {}

// Get from DB directly, bypassing ccurl since it make have some temporary states
func (this *ImplStateDB) GetCommittedState(addr evmcommon.Address, key evmcommon.Hash) evmcommon.Hash {
	if value, _ := this.url.ReadCommitted(this.tid, getStorageKeyPath(this.url, addr, key)); value == nil {
		return evmcommon.Hash{}
	} else {
		v, _, _ := value.(*noncommutative.Bytes).Get()
		return evmcommon.BytesToHash(v.([]byte))
	}
}

func (this *ImplStateDB) GetState(addr evmcommon.Address, key evmcommon.Hash) evmcommon.Hash {
	if value, err := this.url.Read(this.tid, getStorageKeyPath(this.url, addr, key)); err != nil {
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

	path := getStorageKeyPath(this.url, addr, key)
	if err := this.url.Write(this.tid, path, noncommutative.NewBytes(value.Bytes())); err != nil {
		panic(err)
	}
	this.url.LatestKey(path)
}

func (this *ImplStateDB) Exist(addr evmcommon.Address) bool {
	return accountExist(this.url, addr, this.tid)
}

func (this *ImplStateDB) Empty(addr evmcommon.Address) bool {
	if !accountExist(this.url, addr, this.tid) {
		return true
	}

	if value, err := this.url.Peek(getBalancePath(this.url, addr)); err != nil {
		panic(err)
	} else {
		v, _, _ := value.(*commutative.U256).Get()
		if v.(*uint256.Int).Cmp(commutative.U256_ZERO) != 0 {
			return false
		}
	}

	if value, err := this.url.Peek(getNoncePath(this.url, addr)); err != nil {
		panic(err)
	} else {
		v, _, _ := value.(*commutative.Uint64).Get()
		if v.(uint64) != 0 {
			return false
		}
	}

	if value, err := this.url.Peek(getCodePath(this.url, addr)); err != nil {
		panic(err)
	} else {
		return len(value.(*noncommutative.Bytes).Value().(codec.Bytes)) == 0
	}
}

func (this *ImplStateDB) AddLog(log *evmtypes.Log) {
	this.logs[this.txHash] = append(this.logs[this.txHash], log)
}

func (this *ImplStateDB) ForEachStorage(addr evmcommon.Address, f func(evmcommon.Hash, evmcommon.Hash) bool) error {
	return nil
}

func (this *ImplStateDB) PrepareAccessList(sender evmcommon.Address, dest *evmcommon.Address, precompiles []evmcommon.Address, txAccesses evmtypes.AccessList) {
	// Do nothing.
}

func (this *ImplStateDB) AddressInAccessList(addr evmcommon.Address) bool { return true }

func (this *ImplStateDB) SlotInAccessList(addr evmcommon.Address, slot evmcommon.Hash) (addressOk bool, slotOk bool) {
	return true, true
}

func (this *ImplStateDB) Prepare(txHash, bhash evmcommon.Hash, ti int) {
	this.refund = 0
	this.txHash = txHash
	this.tid = uint32(ti)
	this.logs = make(map[evmcommon.Hash][]*evmtypes.Log)
}

func (this *ImplStateDB) GetLogs(hash evmcommon.Hash) []*evmtypes.Log {
	return this.logs[hash]
}

func (this *ImplStateDB) Copy() *ImplStateDB {
	return &ImplStateDB{
		logs: make(map[evmcommon.Hash][]*evmtypes.Log),
		// kapi: this.kapi,
		// db:   this.db,
		url: this.url,
	}
}

func (this *ImplStateDB) Set(eac EthAccountCache, esc EthStorageCache) {
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
// 	this := NewImplStateDB(nil, db, url)
// 	this.Prepare(evmcommon.Hash{}, evmcommon.Hash{}, txIndex)

// 	// Chage the tranasction fees
// 	this.AddBalance(coinbase, new(big.Int).Mul(new(big.Int).SetUint64(gasUsed), gasPrice))
// 	this.SubBalance(from, new(big.Int).Mul(new(big.Int).SetUint64(gasUsed), gasPrice))
// 	//	this.SetNonce(from, 0) don't increase if the transcation failes.
// 	return url.Export(false)
// }

// func ExportOnFailureEx(
// 	db urlcommon.DatastoreInterface,
// 	txIndex int,
// 	from, coinbase evmcommon.Address,
// 	gasUsed uint64, gasPrice *big.Int,
// ) ([][]byte, [][]byte) {
// 	url := concurrenturl.NewConcurrentUrl(db)
// 	this := NewImplStateDB(nil, db, url)
// 	this.Prepare(evmcommon.Hash{}, evmcommon.Hash{}, txIndex)
// 	this.AddBalance(coinbase, new(big.Int).Mul(new(big.Int).SetUint64(gasUsed), gasPrice))
// 	this.SubBalance(from, new(big.Int).Mul(new(big.Int).SetUint64(gasUsed), gasPrice))
// 	// this.SetNonce(from, 0)
// 	return url.ExportEncoded(nil)
// }

// func ExportOnConfliction(
// 	db urlcommon.DatastoreInterface,
// 	txIndex int,
// 	from evmcommon.Address,
// ) ([]urlcommon.UnivalueInterface, []urlcommon.UnivalueInterface) {
// 	url := concurrenturl.NewConcurrentUrl(db)
// 	this := NewImplStateDB(nil, db, url)
// 	this.Prepare(evmcommon.Hash{}, evmcommon.Hash{}, txIndex)
// 	this.SetNonce(from, 0)
// 	return url.ExportAll()
// }

// func ExportOnConflictionEx(
// 	db urlcommon.DatastoreInterface,
// 	txIndex int,
// 	from evmcommon.Address,
// ) ([][]byte, [][]byte) {
// 	url := concurrenturl.NewConcurrentUrl(db)
// 	this := NewImplStateDB(nil, db, url)
// 	this.Prepare(evmcommon.Hash{}, evmcommon.Hash{}, txIndex)
// 	this.SetNonce(from, 0)
// 	return url.ExportEncoded(func(accesses, transitions []urlcommon.UnivalueInterface) ([]urlcommon.UnivalueInterface, []urlcommon.UnivalueInterface) {
// 		for _, t := range transitions {
// 			t.SetTransitionType(urlcommon.INVARIATE_TRANSITIONS)
// 		}
// 		return accesses, transitions
// 	})
// }
