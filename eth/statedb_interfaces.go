package adaptorcommon

// type StateDB interface {
// 	CreateAccount(common.Address)

// 	SubBalance(common.Address, *big.Int)
// 	AddBalance(common.Address, *big.Int)
// 	GetBalance(common.Address) *big.Int
// 	SetBalance(common.Address, *big.Int)

// 	PeekBalance(common.Address) *big.Int // For Arcology

// 	GetNonce(common.Address) uint64
// 	SetNonce(common.Address, uint64)

// 	GetCodeHash(common.Address) common.Hash
// 	GetCode(common.Address) []byte
// 	SetCode(common.Address, []byte)
// 	GetCodeSize(common.Address) int

// 	AddRefund(uint64)
// 	SubRefund(uint64)
// 	GetRefund() uint64

// 	GetCommittedState(common.Address, common.Hash) common.Hash
// 	GetState(common.Address, common.Hash) common.Hash
// 	SetState(common.Address, common.Hash, common.Hash)

// 	GetTransientState(addr common.Address, key common.Hash) common.Hash
// 	SetTransientState(addr common.Address, key, value common.Hash)

// 	SelfDestruct(common.Address)
// 	HasSelfDestructed(common.Address) bool

// 	Selfdestruct6780(common.Address)

// 	// Exist reports whether the given account exists in state.
// 	// Notably this should also return true for suicided accounts.
// 	Exist(common.Address) bool
// 	// Empty returns whether the given account is empty. Empty
// 	// is defined according to EIP161 (balance = nonce = code = 0).
// 	Empty(common.Address) bool

// 	RevertToSnapshot(int)
// 	Snapshot() int

// 	AddLog(*types.Log)
// 	AddPreimage(common.Hash, []byte)

// 	ForEachStorage(common.Address, func(common.Hash, common.Hash) bool) error

// 	PrepareAccessList(sender common.Address, dest *common.Address, precompiles []common.Address, txAccesses types.AccessList)
// 	AddressInAccessList(addr common.Address) bool
// 	SlotInAccessList(addr common.Address, slot common.Hash) (addressOk bool, slotOk bool)
// 	// AddAddressToAccessList adds the given address to the access list. This operation is safe to perform
// 	// even if the feature/fork is not active yet
// 	AddAddressToAccessList(addr common.Address)
// 	// AddSlotToAccessList adds the given (address,slot) to the access list. This operation is safe to perform
// 	// even if the feature/fork is not active yet
// 	AddSlotToAccessList(addr common.Address, slot common.Hash)

// 	PrepareFormer(txHash, bhash common.Hash, ti int)
// 	GetLogs(hash common.Hash) []*types.Log
// 	Copy() StateDB

// 	Set(eac EthAccountCache, esc EthStorageCache)

// 	Prepare(rules params.Rules, sender, coinbase common.Address, dest *common.Address, precompiles []common.Address, txAccesses types.AccessList)
// }

// type Account interface {
// 	GetBalance() *big.Int
// 	GetNonce() uint64
// 	GetCodeHash() []byte
// }

// type EthAccountCache interface {
// 	GetAccount(string) (Account, error)
// 	GetCode(string) ([]byte, error)
// }

// type EthStorageCache interface {
// 	GetState(string, []byte) []byte
// }

// type KernelAPI interface {
// 	AddLog(key, value string)
// 	GetLogs() []intf.ILog
// 	ClearLogs()
// 	Prepare(txHash common.Hash)
// 	Call(caller, callee common.Address, input []byte, origin common.Address, nonce uint64, blockhash common.Hash) ([]byte, bool)
// }
