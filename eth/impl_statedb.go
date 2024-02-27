package eth

import (
	"math/big"

	cache "github.com/arcology-network/eu/cache"
	commutative "github.com/arcology-network/storage-committer/commutative"
	noncommutative "github.com/arcology-network/storage-committer/noncommutative"
	intf "github.com/arcology-network/vm-adaptor/interface"
	"github.com/ethereum/go-ethereum/common"
	evmcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	evmtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/params"
	uint256 "github.com/holiman/uint256"
)

// Arcology implementation of Eth ImplStateDB interfaces.
type ImplStateDB struct {
	refund           uint64
	txHash           evmcommon.Hash
	tid              uint32 // tx id
	logs             map[evmcommon.Hash][]*evmtypes.Log
	transientStorage transientStorage
	api              intf.EthApiRouter
}

func NewImplStateDB(api intf.EthApiRouter) *ImplStateDB {
	return &ImplStateDB{
		logs:             make(map[evmcommon.Hash][]*evmtypes.Log),
		api:              api,
		transientStorage: newTransientStorage(),
	}
}

func (this *ImplStateDB) CreateAccount(addr evmcommon.Address) {
	createAccount(this.api.WriteCache().(*cache.WriteCache), addr, this.tid)
}

func (this *ImplStateDB) SubBalance(addr evmcommon.Address, amount *big.Int) {
	this.AddBalance(addr, new(big.Int).Neg(new(big.Int).Set(amount)))
}

func (this *ImplStateDB) AddBalance(addr evmcommon.Address, amount *big.Int) {
	if !this.Exist(addr) {
		createAccount(this.api.WriteCache().(*cache.WriteCache), addr, this.tid)
	}

	if delta, ok := commutative.NewU256DeltaFromBigInt(new(big.Int).Set(amount)); ok {
		if _, err := this.api.WriteCache().(*cache.WriteCache).Write(this.tid, getBalancePath(this.api.WriteCache().(*cache.WriteCache), addr), delta); err == nil {
			return
		}
	}
	panic("Error: Failed to call AddBalance()")
}

func (this *ImplStateDB) GetBalance(addr evmcommon.Address) *big.Int {
	if !this.Exist(addr) {
		return new(big.Int)
	}

	value, _, _ := this.api.WriteCache().(*cache.WriteCache).Read(this.tid, getBalancePath(this.api.WriteCache().(*cache.WriteCache), addr), new(commutative.U256))
	v := value.(uint256.Int)
	return (&v).ToBig() // v.(*commutative.U256).Value().(*big.Int)
}

func (this *ImplStateDB) PeekBalance(addr evmcommon.Address) *big.Int {
	if !this.Exist(addr) {
		return new(big.Int)
	}

	value, _ := this.api.WriteCache().(*cache.WriteCache).Peek(getBalancePath(this.api.WriteCache().(*cache.WriteCache), addr), new(commutative.U256))
	v := value.(uint256.Int)
	return v.ToBig()

}

func (this *ImplStateDB) SetBalance(addr evmcommon.Address, amount *big.Int) {
	origin := this.PeekBalance(addr)
	this.AddBalance(addr, new(big.Int).Sub(amount, origin))
}

func (this *ImplStateDB) GetNonce(addr evmcommon.Address) uint64 {
	if !this.Exist(addr) {
		return 0
	}

	nonce, _ := this.api.WriteCache().(*cache.WriteCache).Peek(getNoncePath(this.api.WriteCache().(*cache.WriteCache), addr), new(commutative.Uint64))
	return nonce.(uint64) + this.CalculateNonceOffset(addr, nonce.(uint64)) // Add the nonce offset
}

func (this *ImplStateDB) SetNonce(addr evmcommon.Address, nonce uint64) {
	if !this.Exist(addr) {
		createAccount(this.api.WriteCache().(*cache.WriteCache), addr, this.tid)
	}

	// This original implementation will set the nonce to the given value, but here we just write the nonce delta, which is 1 to the cache, becuase the nonce increment is always 1
	// This is Arcology's way to handle the nonce, and the actual nonce will be calculated when it is read or at commit time.
	if _, err := this.api.WriteCache().(*cache.WriteCache).Write(this.tid, getNoncePath(this.api.WriteCache().(*cache.WriteCache), addr), commutative.NewUint64Delta(1)); err != nil {
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
	if !this.Exist(addr) {
		return nil
	}

	value, _, _ := this.api.WriteCache().(*cache.WriteCache).Read(this.tid, getCodePath(this.api.WriteCache().(*cache.WriteCache), addr), new(noncommutative.Bytes))
	return value.([]byte)

}

func (this *ImplStateDB) SetCode(addr evmcommon.Address, code []byte) {
	if !this.Exist(addr) {
		createAccount(this.api.WriteCache().(*cache.WriteCache), addr, this.tid)
	}

	if _, err := this.api.WriteCache().(*cache.WriteCache).Write(this.tid, getCodePath(this.api.WriteCache().(*cache.WriteCache), addr), noncommutative.NewBytes(code)); err != nil {
		panic(err)
	}
}

func (this *ImplStateDB) SelfDestruct(addr evmcommon.Address)           { return }
func (this *ImplStateDB) HasSelfDestructed(addr evmcommon.Address) bool { return false }
func (this *ImplStateDB) Selfdestruct6780(common.Address)               {}

func (this *ImplStateDB) GetCodeSize(addr evmcommon.Address) int                          { return len(this.GetCode(addr)) }
func (this *ImplStateDB) AddRefund(amount uint64)                                         { this.refund += amount }
func (this *ImplStateDB) SubRefund(amount uint64)                                         { this.refund -= amount }
func (this *ImplStateDB) GetRefund() uint64                                               { return this.refund }
func (this *ImplStateDB) RevertToSnapshot(id int)                                         {}
func (this *ImplStateDB) Snapshot() int                                                   { return 0 }
func (this *ImplStateDB) AddPreimage(hash evmcommon.Hash, preimage []byte)                {}
func (this *ImplStateDB) AddAddressToAccessList(addr evmcommon.Address)                   {} // Do nothing.
func (this *ImplStateDB) AddSlotToAccessList(addr evmcommon.Address, slot evmcommon.Hash) {}

// func (this *ImplStateDB) Set(eac EthAccountCache, esc EthStorageCache)                    {} // TODO

// Get from DB directly, bypassing ccurl since it make have some temporary states
func (this *ImplStateDB) GetCommittedState(addr evmcommon.Address, key evmcommon.Hash) evmcommon.Hash {
	if value, _ := this.api.WriteCache().(*cache.WriteCache).ReadCommitted(this.tid, getStorageKeyPath(this.api, addr, key), new(noncommutative.Bytes)); value != nil {
		// v, _, _ := value.(interfaces.Type).Get()
		return evmcommon.BytesToHash(value.([]byte))
	}
	return evmcommon.Hash{}
}

func (this *ImplStateDB) GetState(addr evmcommon.Address, key evmcommon.Hash) evmcommon.Hash {
	if value, _, _ := this.api.WriteCache().(*cache.WriteCache).Read(this.tid, getStorageKeyPath(this.api, addr, key), new(noncommutative.Bytes)); value != nil {
		return evmcommon.BytesToHash(value.([]byte))
	}
	return evmcommon.Hash{}
}

func (this *ImplStateDB) SetState(addr evmcommon.Address, key, value evmcommon.Hash) {
	if !this.Exist(addr) {
		createAccount(this.api.WriteCache().(*cache.WriteCache), addr, this.tid)
	}

	path := getStorageKeyPath(this.api, addr, key)
	if _, err := this.api.WriteCache().(*cache.WriteCache).Write(this.tid, path, noncommutative.NewBytes(value.Bytes())); err != nil {
		panic(err)
	}
}

func (this *ImplStateDB) Exist(addr evmcommon.Address) bool {
	return accountExist(this.api.WriteCache().(*cache.WriteCache), addr, this.tid)
}

func (this *ImplStateDB) Empty(addr evmcommon.Address) bool {
	return (!this.Exist(addr)) || (this.PeekBalance(addr).BitLen() == 0 && this.GetNonce(addr) == 0 && this.GetCodeSize(addr) == 0)
}

// SetTransientState sets transient storage for a given account. It
// adds the change to the journal so that it can be rolled back
// to its previous value if there is a revert.
func (s *ImplStateDB) SetTransientState(addr evmcommon.Address, key, value evmcommon.Hash) {
	prev := s.GetTransientState(addr, key)
	if prev == value {
		return
	}
	s.setTransientState(addr, key, value)
}

// setTransientState is a lower level setter for transient storage. It
// is called during a revert to prevent modifications to the journal.
func (s *ImplStateDB) setTransientState(addr evmcommon.Address, key, value evmcommon.Hash) {
	s.transientStorage.Set(addr, key, value)
}

// GetTransientState gets transient storage for a given account.
func (s *ImplStateDB) GetTransientState(addr evmcommon.Address, key evmcommon.Hash) evmcommon.Hash {
	return s.transientStorage.Get(addr, key)
}
func (this *ImplStateDB) Prepare(rules params.Rules, sender, coinbase common.Address, dest *common.Address, precompiles []common.Address, txAccesses types.AccessList) {
	//Do nothing
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

func (this *ImplStateDB) PrepareFormer(txHash, bhash evmcommon.Hash, ti uint32) {
	this.refund = 0
	this.txHash = txHash
	this.tid = ti
	this.logs = make(map[evmcommon.Hash][]*evmtypes.Log)
}

func (this *ImplStateDB) GetLogs(hash evmcommon.Hash) []*evmtypes.Log {
	return this.logs[hash]
}
