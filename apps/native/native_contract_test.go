package native

import (
	"math/big"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	cachedstorage "github.com/arcology-network/common-lib/cachedstorage"
	"github.com/arcology-network/concurrenturl/v2"
	urlcommon "github.com/arcology-network/concurrenturl/v2/common"
	curstorage "github.com/arcology-network/concurrenturl/v2/storage"
	"github.com/arcology-network/concurrenturl/v2/type/commutative"
	evmcommon "github.com/arcology-network/evm/common"
	"github.com/arcology-network/vm-adaptor/abi"
	eucommon "github.com/arcology-network/vm-adaptor/common"
	cceueth "github.com/arcology-network/vm-adaptor/eth"
	tests "github.com/arcology-network/vm-adaptor/tests"
)

func TestNativeContractSameBlock(t *testing.T) {
	persistentDB := cachedstorage.NewDataStore()
	meta, _ := commutative.NewMeta(urlcommon.NewPlatform().Eth10Account())
	persistentDB.Inject(urlcommon.NewPlatform().Eth10Account(), meta)
	db := curstorage.NewTransientDB(persistentDB)

	url := concurrenturl.NewConcurrentUrl(db)
	statedb := cceueth.NewImplStateDB(url)
	statedb.Prepare(evmcommon.Hash{}, evmcommon.Hash{}, 0)
	statedb.CreateAccount(eucommon.Coinbase)
	// User 1
	statedb.CreateAccount(eucommon.User1)
	statedb.AddBalance(eucommon.User1, new(big.Int).SetUint64(1e18))
	// user2
	statedb.CreateAccount(eucommon.User2)
	statedb.AddBalance(eucommon.User2, new(big.Int).SetUint64(1e18))
	// Contract owner
	statedb.CreateAccount(eucommon.Owner)
	statedb.AddBalance(eucommon.Owner, new(big.Int).SetUint64(1e18))

	_, transitions := url.Export(true)

	// ================================== Compile ==================================
	_, err := exec.Command("python", "./compiler.py").Output() // capture the output of the Python script
	if err != nil {
		panic(err)
	}

	bytecode, err := eucommon.BytecodeReader("./bytecode.txt") // Read the byte code
	if err != nil {
		t.Error("Error: ", err)
	}

	// Compile
	// ================================ Deploy the contract==================================
	eu, config := tests.Prepare(db, 10000000, transitions, []uint32{0})
	transitions, receipt, _, err := tests.Deploy(eu, config, eucommon.Owner, 0, bytecode)
	t.Log("\n" + eucommon.FormatTransitions(transitions))
	t.Log(receipt)
	address := receipt.ContractAddress
	t.Log(address)
	if receipt.Status != 1 {
		t.Error("Error: Failed to deploy!!!", err)
	}

	// Increment x by one
	if _, _, receipt, err = tests.CallFunc(eu, config, &eucommon.User1, &address, 0, true, "incrementX()"); receipt.Status != 1 || err != nil {
		t.Error("Error: Failed to call incrementX() 1!!!", err)
	}

	if _, _, receipt, err = tests.CallFunc(eu, config, &eucommon.User1, &address, 0, true, "incrementX()"); receipt.Status != 1 || err != nil {
		t.Error("Error: Failed to call incrementX() 2!!!", err)
	}

	if _, _, receipt, err = tests.CallFunc(eu, config, &eucommon.User1, &address, 0, true, "incrementX()"); receipt.Status != 1 || err != nil {
		t.Error("Error: Failed to call incrementX() 3!!!", err)
	}

	encoded, _ := abi.Encode(uint64(102))
	if _, _, receipt, err := tests.CallFunc(eu, config, &eucommon.User1, &address, 0, true, "checkY(uint256)", encoded); receipt.Status != 1 || err != nil {
		t.Error("Error: Failed to check checkY() 1!!!", err)
	}

	if _, _, receipt, err = tests.CallFunc(eu, config, &eucommon.User1, &address, 0, true, "incrementY()"); receipt.Status != 1 {
		t.Error("Error: Failed to call incrementY() 1!!!", err)
	}

	encoded, _ = abi.Encode(uint64(104))
	if _, _, receipt, err := tests.CallFunc(eu, config, &eucommon.User1, &address, 0, true, "checkY(uint256)", encoded); receipt.Status != 1 || err != nil {
		t.Error("Error: Failed to check checkY() 2!!!", err)
	}
}

func TestNativeContractAcrossBlocks(t *testing.T) {
	persistentDB := cachedstorage.NewDataStore()
	meta, _ := commutative.NewMeta(urlcommon.NewPlatform().Eth10Account())
	persistentDB.Inject(urlcommon.NewPlatform().Eth10Account(), meta)
	db := curstorage.NewTransientDB(persistentDB)

	url := concurrenturl.NewConcurrentUrl(db)
	statedb := cceueth.NewImplStateDB(url)
	statedb.Prepare(evmcommon.Hash{}, evmcommon.Hash{}, 0)
	statedb.CreateAccount(eucommon.Coinbase)
	// User 1
	statedb.CreateAccount(eucommon.User1)
	statedb.AddBalance(eucommon.User1, new(big.Int).SetUint64(1e18))
	// user2
	statedb.CreateAccount(eucommon.User2)
	statedb.AddBalance(eucommon.User2, new(big.Int).SetUint64(1e18))
	// Contract owner
	statedb.CreateAccount(eucommon.Owner)
	statedb.AddBalance(eucommon.Owner, new(big.Int).SetUint64(1e18))

	_, transitions := url.Export(true)

	// ================================== Compile ==================================
	currentPath, _ := os.Getwd()
	compiler := filepath.Dir(currentPath) + "/compiler.py"
	code, err := eucommon.CompileContracts(compiler, "./NativeStorage.sol", "NativeStorage")
	if err != nil || len(code) == 0 {
		t.Error("Error: Failed to generate the byte code")
	}

	// ================================ Deploy the contract==================================
	eu, config := tests.Prepare(db, 10000000, transitions, []uint32{0})
	transitions, receipt, _, err := tests.Deploy(eu, config, eucommon.Owner, 0, code)
	t.Log("\n" + eucommon.FormatTransitions(transitions))
	t.Log(receipt)
	address := receipt.ContractAddress
	t.Log(address)
	if receipt.Status != 1 {
		t.Error("Error: Failed to deploy!!!", err)
	}

	eu, config = tests.Prepare(db, 10000001, transitions, []uint32{1})
	// encoded, _ := abi.Encode(uint64(2))
	_, transitions, receipt, err = tests.CallFunc(eu, config, &eucommon.User1, &address, 0, true, "incrementX()")
	t.Log("\n" + eucommon.FormatTransitions(transitions))
	t.Log(receipt)
	if receipt.Status != 1 || err != nil {
		t.Error("Error: Failed to call incrementX() 1!!!", err)
	}

	eu, config = tests.Prepare(db, 10000001, transitions, []uint32{0})
	encoded, _ := abi.Encode(uint64(3))
	acc, transitions, receipt, _, err := tests.CallFunc(eu, config, &eucommon.User1, &address, 0, true, "checkX(uint256)", encoded)
	t.Log("\n" + eucommon.FormatTransitions(acc))
	t.Log(receipt)
	if receipt.Status != 1 {
		t.Error("Error: Failed to call checkX() 1!!!", err)
	}

	eu, config = tests.Prepare(db, 10000001, transitions, []uint32{0})
	encoded, _ = abi.Encode(uint64(102))
	acc, transitions, receipt, err = tests.CallFunc(eu, config, &eucommon.User1, &address, 0, true, "checkY(uint256)", encoded)
	t.Log("\n" + eucommon.FormatTransitions(acc))
	t.Log(receipt)
	if receipt.Status != 1 {
		t.Error("Error: Failed to call checkY() 1!!!", err)
	}

	// eu, config = eucommon.Prepare(db, 10000001, transitions, []uint32{0})
	// encoded, _ = abi.Encode(uint64(2))
	// acc, transitions, receipt, err = tests.CallFunc(eu, config, &eucommon.User1, &address, 0, true, "checkX(uint256)", encoded)
	// t.Log("\n" + eucommon.FormatTransitions(acc))
	// t.Log(receipt)
	// if receipt.Status != 1 {
	// 	t.Error("Error: Failed to call checkX() 1!!!", err)
	// }

}
