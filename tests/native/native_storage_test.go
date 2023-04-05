package native

import (
	"math/big"
	"os/exec"
	"testing"

	cachedstorage "github.com/arcology-network/common-lib/cachedstorage"
	"github.com/arcology-network/concurrenturl/v2"
	urlcommon "github.com/arcology-network/concurrenturl/v2/common"
	curstorage "github.com/arcology-network/concurrenturl/v2/storage"
	"github.com/arcology-network/concurrenturl/v2/type/commutative"
	evmcommon "github.com/arcology-network/evm/common"
	cceueth "github.com/arcology-network/vm-adaptor/eth"
	"github.com/arcology-network/vm-adaptor/tests"
)

func TestNativeStorage(t *testing.T) {
	persistentDB := cachedstorage.NewDataStore()
	meta, _ := commutative.NewMeta(urlcommon.NewPlatform().Eth10Account())
	persistentDB.Inject(urlcommon.NewPlatform().Eth10Account(), meta)
	db := curstorage.NewTransientDB(persistentDB)

	url := concurrenturl.NewConcurrentUrl(db)
	statedb := cceueth.NewImplStateDB(url)
	statedb.Prepare(evmcommon.Hash{}, evmcommon.Hash{}, 0)
	statedb.CreateAccount(tests.Coinbase)
	statedb.CreateAccount(tests.User1)
	statedb.AddBalance(tests.User1, new(big.Int).SetUint64(1e18))

	statedb.CreateAccount(tests.Owner)
	statedb.AddBalance(tests.Owner, new(big.Int).SetUint64(1e18))

	_, transitions := url.Export(true)

	// ================================== Compile ==================================
	_, err := exec.Command("python", "./compiler.py").Output() // capture the output of the Python script
	if err != nil {
		panic(err)
	}

	bytecode, err := tests.BytecodeReader("./bytecode.txt") // Read the byte code
	if err != nil {
		t.Error("Error: ", err)
	}

	// Compile
	// ================= Deploy the contract==================================
	eu, config := tests.Prepare(db, 10000000, transitions, []uint32{0})
	transitions, receipt, err := tests.Deploy(eu, config, tests.Owner, 0, bytecode)
	t.Log("\n" + tests.FormatTransitions(transitions))
	t.Log(receipt)
	address := receipt.ContractAddress
	t.Log(address)
	if receipt.Status != 1 {
		t.Error("Error: Failed to deploy!!!", err)
	}

	// ================================== Call the contract ==================================
	// Call accessX.
	eu, config = tests.Prepare(db, 10000001, transitions, []uint32{1})
	acc, transitions, receipt, err := tests.RunEx(eu, config, &tests.User1, &address, 1, true, "accessX()")
	t.Log("\n" + tests.FormatTransitions(acc))
	t.Log(receipt)

	if receipt.Status != 1 {
		t.Error("Error: Failed to call accessX()!!!", err)
	}

	// Call accessY.
	eu, config = tests.Prepare(db, 10000002, transitions, []uint32{2})
	acc, _, receipt, err = tests.RunEx(eu, config, &tests.User1, &address, 2, true, "accessY()")
	t.Log("\n" + tests.FormatTransitions(acc))
	t.Log(receipt)
	if receipt.Status != 1 {
		t.Error("Error: Failed to call accessY()!!!", err)
	}
}
