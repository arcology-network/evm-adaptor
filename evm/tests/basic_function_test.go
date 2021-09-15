package tests

import (
	"fmt"
	"math/big"
	"reflect"
	"testing"

	cmncommon "github.com/arcology-network/common-lib/common"
	cmntypes "github.com/arcology-network/common-lib/types"
	"github.com/arcology-network/concurrenturl/v2"
	urlcommon "github.com/arcology-network/concurrenturl/v2/common"
	urltype "github.com/arcology-network/concurrenturl/v2/type"
	"github.com/arcology-network/concurrenturl/v2/type/commutative"
	"github.com/arcology-network/concurrenturl/v2/type/noncommutative"
	"github.com/arcology-network/evm/common"
	"github.com/arcology-network/evm/core/types"
	"github.com/arcology-network/evm/crypto"
	adaptor "github.com/arcology-network/vm-adaptor/evm"
)

var (
	basicFunctionCode = "60806040526000805534801561001457600080fd5b50608173ffffffffffffffffffffffffffffffffffffffff1663f02e3aff6001600381111561003f57fe5b6002600381111561004c57fe5b6040518363ffffffff167c010000000000000000000000000000000000000000000000000000000002815260040180806020018460030b60030b81526020018360030b60030b8152602001828103825260078152602001807f62616c616e6365000000000000000000000000000000000000000000000000008152506020019350505050600060405180830381600087803b1580156100ea57600080fd5b505af11580156100fe573d6000803e3d6000fd5b50505050610390806101116000396000f3fe608060405260043610610057576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff1680633825d8281461005c578063569c5f6d146100b7578063853255cc146100ce575b600080fd5b34801561006857600080fd5b506100b56004803603604081101561007f57600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803590602001909291905050506100f9565b005b3480156100c357600080fd5b506100cc610323565b005b3480156100da57600080fd5b506100e361035e565b6040518082815260200191505060405180910390f35b6000608173ffffffffffffffffffffffffffffffffffffffff1663c41eb85a846040518263ffffffff167c010000000000000000000000000000000000000000000000000000000002815260040180806020018373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001828103825260078152602001807f62616c616e6365000000000000000000000000000000000000000000000000008152506020019250505060206040518083038186803b1580156101ce57600080fd5b505afa1580156101e2573d6000803e3d6000fd5b505050506040513d60208110156101f857600080fd5b81019080805190602001909291905050509050608173ffffffffffffffffffffffffffffffffffffffff16634f7c4f4c84846040518363ffffffff167c010000000000000000000000000000000000000000000000000000000002815260040180806020018473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001838152602001828103825260078152602001807f62616c616e6365000000000000000000000000000000000000000000000000008152506020019350505050600060405180830381600087803b1580156102e857600080fd5b505af11580156102fc573d6000803e3d6000fd5b50505050806000808282540392505081905550816000808282540192505081905550505050565b7ff3bb7cb62e792d1ead1cfd4901fdf7fbd9ab1e5db0c69c810ff68318b141e0776000546040518082815260200191505060405180910390a1565b6000548156fea165627a7a7230582050434841bd968433e6f25244b39d803a5fe19d63347f67fe8ba7d7115598bcc50029"
)

func TestEncodeDecode(t *testing.T) {
	config := MainConfig()
	persistentDB := urlcommon.NewDataStore()
	meta, _ := commutative.NewMeta(urlcommon.NewPlatform().Eth10Account())
	persistentDB.Save(urlcommon.NewPlatform().Eth10Account(), meta)
	db := urlcommon.NewTransientDB(persistentDB)

	url := concurrenturl.NewConcurrentUrl(db)
	api := adaptor.NewAPIV2(db, url)
	statedb := adaptor.NewStateDBV2(api, db, url)
	statedb.Prepare(common.Hash{}, common.Hash{}, 0)
	statedb.CreateAccount(coinbase)
	statedb.CreateAccount(user1)
	statedb.AddBalance(user1, new(big.Int).SetUint64(1e18))
	_, transitions := url.Export(true)

	// Deploy.
	url = concurrenturl.NewConcurrentUrl(db)
	url.Commit(transitions, []uint32{0})
	api = adaptor.NewAPIV2(db, url)
	statedb = adaptor.NewStateDBV2(api, db, url)
	eu := adaptor.NewEUV2(config.ChainConfig, *config.VMConfig, config.Chain, statedb, api, db, url)

	config.Coinbase = &coinbase
	config.BlockNumber = new(big.Int).SetUint64(10000000)
	config.Time = new(big.Int).SetUint64(10000000)

	msg := types.NewMessage(user1, nil, 0, new(big.Int).SetUint64(0), 1e15, new(big.Int).SetUint64(1), common.Hex2Bytes(basicFunctionCode), nil, true)
	_, transitions, receipt := eu.Run(common.BytesToHash([]byte{1, 1, 1}), 1, &msg, adaptor.NewEVMBlockContextV2(config), adaptor.NewEVMTxContext(msg))

	t.Log("\n" + formatTransitions(transitions))
	//t.Log(receipt)
	fmt.Printf("receipt.Status=%v,receipt=%v,receipt.ContractAddress=%v\n", receipt.Status, receipt, receipt.ContractAddress)

	// for _, transtion := range transitions {
	// 	transtion.Print()
	// }
	trans := make([]interface{}, len(transitions))
	for i := range trans {
		trans[i] = transitions[i]
	}
	result := cmntypes.EuResult{
		ID:          uint32(1),
		Transitions: trans,
	}
	sendData, err := cmncommon.GobEncode(result)
	if err != nil {
		fmt.Sprintf("encode err=%v\n", err)
		return
	}

	var dresult *cmntypes.EuResult
	err = cmncommon.GobDecode(sendData, &dresult)
	if err != nil {
		fmt.Sprintf("decode err=%v\n", err)
		return
	}

	if len(result.Transitions) != len(dresult.Transitions) ||
		result.H != dresult.H ||
		result.ID != dresult.ID ||
		result.DC != dresult.DC ||
		result.Status != dresult.Status ||
		result.GasUsed != dresult.GasUsed {
		t.Error("Transitions Mismatched !")
	}

	for i := range result.Transitions {
		if !reflect.DeepEqual(result.Transitions[i].(urlcommon.UnivalueInterface), dresult.Transitions[i].(urlcommon.UnivalueInterface)) {
			t.Error("Transitions Mismatched !")
			return
		}
	}

	// fmt.Println("=====================================================")
	// for _, transtion := range dresult.Transitions {
	// 	transtion.Print()
	// }
	detransitions := make([]urlcommon.UnivalueInterface, len(dresult.Transitions))
	for i := range dresult.Transitions {
		detransitions[i] = dresult.Transitions[i].(urlcommon.UnivalueInterface)
	}
	t.Log("\n" + formatTransitions(detransitions))

}

func TestBasicFunction(t *testing.T) {
	config := MainConfig()
	persistentDB := urlcommon.NewDataStore()
	meta, _ := commutative.NewMeta(urlcommon.NewPlatform().Eth10Account())
	persistentDB.Save(urlcommon.NewPlatform().Eth10Account(), meta)
	db := urlcommon.NewTransientDB(persistentDB)

	url := concurrenturl.NewConcurrentUrl(db)
	api := adaptor.NewAPIV2(db, url)
	statedb := adaptor.NewStateDBV2(api, db, url)
	statedb.Prepare(common.Hash{}, common.Hash{}, 0)
	statedb.CreateAccount(coinbase)
	statedb.CreateAccount(user1)
	statedb.AddBalance(user1, new(big.Int).SetUint64(1e18))
	_, transitions := url.Export(true)
	t.Log("\n" + formatTransitions(transitions))

	// Deploy.
	url = concurrenturl.NewConcurrentUrl(db)
	url.Commit(transitions, []uint32{0})
	api = adaptor.NewAPIV2(db, url)
	statedb = adaptor.NewStateDBV2(api, db, url)
	eu := adaptor.NewEUV2(config.ChainConfig, *config.VMConfig, config.Chain, statedb, api, db, url)

	config.Coinbase = &coinbase
	config.BlockNumber = new(big.Int).SetUint64(10000000)
	config.Time = new(big.Int).SetUint64(10000000)

	msg := types.NewMessage(user1, nil, 0, new(big.Int).SetUint64(0), 1e15, new(big.Int).SetUint64(1), common.Hex2Bytes(basicFunctionCode), nil, true)
	_, transitions, receipt := eu.Run(common.BytesToHash([]byte{1, 1, 1}), 1, &msg, adaptor.NewEVMBlockContextV2(config), adaptor.NewEVMTxContext(msg))
	// t.Log("\n" + formatTransitions(accesses))
	t.Log("\n" + formatTransitions(transitions))
	t.Log(receipt)
	contractAddress := receipt.ContractAddress

	// Set.
	url = concurrenturl.NewConcurrentUrl(db)
	errs := url.Commit(transitions, []uint32{1})
	if len(errs) != 0 {
		t.Error(errs)
		return
	}
	api = adaptor.NewAPIV2(db, url)
	statedb = adaptor.NewStateDBV2(api, db, url)
	eu = adaptor.NewEUV2(config.ChainConfig, *config.VMConfig, config.Chain, statedb, api, db, url)

	config.BlockNumber = new(big.Int).SetUint64(10000001)
	config.Time = new(big.Int).SetUint64(10000001)

	data := crypto.Keccak256([]byte("set(address,uint256)"))[:4]
	data = append(data, common.BytesToHash(user1.Bytes()).Bytes()...)
	data = append(data, common.BytesToHash([]byte{0xcc}).Bytes()...)
	msg = types.NewMessage(user1, &contractAddress, 1, new(big.Int).SetUint64(0), 1e15, new(big.Int).SetUint64(1), data, nil, true)
	_, transitions, receipt = eu.Run(common.BytesToHash([]byte{2, 2, 2}), 2, &msg, adaptor.NewEVMBlockContextV2(config), adaptor.NewEVMTxContext(msg))
	t.Log("\n" + formatTransitions(transitions))
	t.Log(receipt)

	// Get.
	url = concurrenturl.NewConcurrentUrl(db)
	errs = url.Commit(transitions, []uint32{2})
	if len(errs) != 0 {
		t.Error(errs)
		return
	}
	api = adaptor.NewAPIV2(db, url)
	statedb = adaptor.NewStateDBV2(api, db, url)
	eu = adaptor.NewEUV2(config.ChainConfig, *config.VMConfig, config.Chain, statedb, api, db, url)

	config.BlockNumber = new(big.Int).SetUint64(10000002)
	config.Time = new(big.Int).SetUint64(10000002)

	data = crypto.Keccak256([]byte("getSum()"))[:4]
	msg = types.NewMessage(user1, &contractAddress, 2, new(big.Int).SetUint64(0), 1e15, new(big.Int).SetUint64(1), data, nil, true)
	accesses, transitions, receipt := eu.Run(common.BytesToHash([]byte{3, 3, 3}), 3, &msg, adaptor.NewEVMBlockContextV2(config), adaptor.NewEVMTxContext(msg))
	t.Log("\n" + formatTransitions(accesses))
	t.Log("\n" + formatTransitions(transitions))

	t.Log(receipt)
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
