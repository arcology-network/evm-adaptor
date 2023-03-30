package basic

import (
	"fmt"
	"math/big"

	urlcommon "github.com/arcology-network/concurrenturl/v2/common"
	urltype "github.com/arcology-network/concurrenturl/v2/type"
	"github.com/arcology-network/concurrenturl/v2/type/commutative"
	"github.com/arcology-network/concurrenturl/v2/type/noncommutative"
	"github.com/holiman/uint256"
)

// var (
// 	basicFunctionCode = "60806040526000805534801561001457600080fd5b50608173ffffffffffffffffffffffffffffffffffffffff1663f02e3aff6001600381111561003f57fe5b6002600381111561004c57fe5b6040518363ffffffff167c010000000000000000000000000000000000000000000000000000000002815260040180806020018460030b60030b81526020018360030b60030b8152602001828103825260078152602001807f62616c616e6365000000000000000000000000000000000000000000000000008152506020019350505050600060405180830381600087803b1580156100ea57600080fd5b505af11580156100fe573d6000803e3d6000fd5b50505050610390806101116000396000f3fe608060405260043610610057576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff1680633825d8281461005c578063569c5f6d146100b7578063853255cc146100ce575b600080fd5b34801561006857600080fd5b506100b56004803603604081101561007f57600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803590602001909291905050506100f9565b005b3480156100c357600080fd5b506100cc610323565b005b3480156100da57600080fd5b506100e361035e565b6040518082815260200191505060405180910390f35b6000608173ffffffffffffffffffffffffffffffffffffffff1663c41eb85a846040518263ffffffff167c010000000000000000000000000000000000000000000000000000000002815260040180806020018373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001828103825260078152602001807f62616c616e6365000000000000000000000000000000000000000000000000008152506020019250505060206040518083038186803b1580156101ce57600080fd5b505afa1580156101e2573d6000803e3d6000fd5b505050506040513d60208110156101f857600080fd5b81019080805190602001909291905050509050608173ffffffffffffffffffffffffffffffffffffffff16634f7c4f4c84846040518363ffffffff167c010000000000000000000000000000000000000000000000000000000002815260040180806020018473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001838152602001828103825260078152602001807f62616c616e6365000000000000000000000000000000000000000000000000008152506020019350505050600060405180830381600087803b1580156102e857600080fd5b505af11580156102fc573d6000803e3d6000fd5b50505050806000808282540392505081905550816000808282540192505081905550505050565b7ff3bb7cb62e792d1ead1cfd4901fdf7fbd9ab1e5db0c69c810ff68318b141e0776000546040518082815260200191505060405180910390a1565b6000548156fea165627a7a7230582050434841bd968433e6f25244b39d803a5fe19d63347f67fe8ba7d7115598bcc50029"
// )

// func TestBasicFunction(t *testing.T) {
// 	config := MainConfig()
// 	persistentDB := cachedstorage.NewDataStore()
// 	meta, _ := commutative.NewMeta(urlcommon.NewPlatform().Eth10Account())
// 	persistentDB.Inject(urlcommon.NewPlatform().Eth10Account(), meta)
// 	db := curstorage.NewTransientDB(persistentDB)

// 	url := concurrenturl.NewConcurrentUrl(db)
// 	api := adaptor.NewAPI(db, url)
// 	statedb := adaptor.NewStateDB(api, db, url)
// 	statedb.Prepare(common.Hash{}, common.Hash{}, 0)
// 	statedb.CreateAccount(coinbase)
// 	statedb.CreateAccount(user1)
// 	statedb.AddBalance(user1, new(big.Int).SetUint64(1e18))
// 	_, transitions := url.Export(true)
// 	t.Log("\n" + FormatTransitions(transitions))

// 	// Deploy.
// 	url = concurrenturl.NewConcurrentUrl(db)
// 	url.Import(transitions)
// 	url.PostImport()
// 	url.Commit([]uint32{0})
// 	api = adaptor.NewAPI(db, url)
// 	statedb = adaptor.NewStateDB(api, db, url)
// 	eu := adaptor.NewEU(config.ChainConfig, *config.VMConfig, config.Chain, statedb, api, db, url)

// 	config.Coinbase = &coinbase
// 	config.BlockNumber = new(big.Int).SetUint64(10000000)
// 	config.Time = new(big.Int).SetUint64(10000000)

// 	// Message to execute
// 	msg := types.NewMessage(user1, nil, 0, new(big.Int).SetUint64(0), 1e15, new(big.Int).SetUint64(1), common.Hex2Bytes(basicFunctionCode), nil, true)

// 	_, transitions, receipt := eu.Run(common.BytesToHash([]byte{1, 1, 1}), 1, &msg, adaptor.NewEVMBlockContextV2(config), adaptor.NewEVMTxContext(msg))
// 	// t.Log("\n" + FormatTransitions(accesses))
// 	t.Log("\n" + FormatTransitions(transitions))
// 	t.Log(receipt)
// 	contractAddress := receipt.ContractAddress

// 	// Set.
// 	url = concurrenturl.NewConcurrentUrl(db)
// 	url.Import(transitions)
// 	errs := url.Commit([]uint32{1})
// 	if len(errs) != 0 {
// 		t.Error(errs)
// 		return
// 	}
// 	api = adaptor.NewAPI(db, url)
// 	statedb = adaptor.NewStateDB(api, db, url)
// 	eu = adaptor.NewEU(config.ChainConfig, *config.VMConfig, config.Chain, statedb, api, db, url)

// 	config.BlockNumber = new(big.Int).SetUint64(10000001)
// 	config.Time = new(big.Int).SetUint64(10000001)

// 	data := crypto.Keccak256([]byte("set(address,uint256)"))[:4]
// 	data = append(data, common.BytesToHash(user1.Bytes()).Bytes()...)
// 	data = append(data, common.BytesToHash([]byte{0xcc}).Bytes()...)
// 	msg = types.NewMessage(user1, &contractAddress, 1, new(big.Int).SetUint64(0), 1e15, new(big.Int).SetUint64(1), data, nil, true)
// 	_, transitions, receipt = eu.Run(common.BytesToHash([]byte{2, 2, 2}), 2, &msg, adaptor.NewEVMBlockContextV2(config), adaptor.NewEVMTxContext(msg))
// 	t.Log("\n" + FormatTransitions(transitions))
// 	t.Log(receipt)

// 	// Get.
// 	url = concurrenturl.NewConcurrentUrl(db)
// 	url.Import(transitions)
// 	url.PostImport()
// 	errs = url.Commit([]uint32{2})
// 	if len(errs) != 0 {
// 		t.Error(errs)
// 		return
// 	}
// 	api = adaptor.NewAPI(db, url)
// 	statedb = adaptor.NewStateDB(api, db, url)
// 	eu = adaptor.NewEU(config.ChainConfig, *config.VMConfig, config.Chain, statedb, api, db, url)

// 	config.BlockNumber = new(big.Int).SetUint64(10000002)
// 	config.Time = new(big.Int).SetUint64(10000002)

// 	data = crypto.Keccak256([]byte("getSum()"))[:4]
// 	msg = types.NewMessage(user1, &contractAddress, 2, new(big.Int).SetUint64(0), 1e15, new(big.Int).SetUint64(1), data, nil, true)
// 	accesses, transitions, receipt := eu.Run(common.BytesToHash([]byte{3, 3, 3}), 3, &msg, adaptor.NewEVMBlockContextV2(config), adaptor.NewEVMTxContext(msg))
// 	t.Log("\n" + FormatTransitions(accesses))
// 	t.Log("\n" + FormatTransitions(transitions))

// 	t.Log(receipt)
// }

func FormatValue(value interface{}) string {
	switch value.(type) {
	case *commutative.Meta:
		meta := value.(*commutative.Meta)
		var str string
		str += "{"
		for i, k := range meta.PeekKeys() {
			str += k
			if i != len(meta.PeekKeys())-1 {
				str += ", "
			}
		}
		str += "}"
		if len(meta.PeekAdded()) != 0 {
			str += " + {"
			for i, k := range meta.PeekAdded() {
				str += k
				if i != len(meta.PeekAdded())-1 {
					str += ", "
				}
			}
			str += "}"
		}
		if len(meta.PeekRemoved()) != 0 {
			str += " - {"
			for i, k := range meta.PeekRemoved() {
				str += k
				if i != len(meta.PeekRemoved())-1 {
					str += ", "
				}
			}
			str += "}"
		}
		return str
	case *noncommutative.Int64:
		// uint256.NewInt(0)
		return fmt.Sprintf(" = %v", (*(value.(*noncommutative.Int64))))
	case *noncommutative.Bytes:
		return fmt.Sprintf(" = %v", value.(*noncommutative.Bytes).Data())
	case *commutative.U256:
		v := value.(*commutative.U256).Value()
		d := value.(*commutative.U256).GetDelta()
		return fmt.Sprintf(" = %v + %v", (*(v.(*uint256.Int))), d.(*big.Int).Int64())
	case *commutative.Int64:
		v := value.(*commutative.Int64).Value()
		d := value.(*commutative.Int64).GetDelta()
		return fmt.Sprintf(" = %v + %v", v, d)
	}
	return ""
}

func FormatTransitions(transitions []urlcommon.UnivalueInterface) string {
	var str string
	for _, t := range transitions {
		str += fmt.Sprintf("%v%v%v%v%v%v%v%v%v%v%v%v%v%v",
			"Tx=", t.(*urltype.Univalue).GetTx(),
			" Reads=", t.(*urltype.Univalue).Reads(),
			" Writes=", t.(*urltype.Univalue).Writes(),
			" Preexists=", t.(*urltype.Univalue).Preexist(),
			" Composite=", t.(*urltype.Univalue).Composite(),
			" Path=", *(t.(*urltype.Univalue).GetPath()),
			" Value", FormatValue(t.(*urltype.Univalue).Value())+"\n")
	}
	return str
}
