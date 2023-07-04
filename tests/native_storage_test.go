package tests

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/holiman/uint256"
	sha3 "golang.org/x/crypto/sha3"
)

func TestSlotHash(t *testing.T) {
	_ctrn := uint256.NewInt(2).Bytes32()
	_elem := uint256.NewInt(2).Bytes32()

	hash := sha3.NewLegacyKeccak256()
	hash.Write(append(_elem[:], _ctrn[:]...))
	v := hash.Sum(nil)
	fmt.Println(v)
	fmt.Println("0x" + hex.EncodeToString(v))

}

// func TestStorageSlot(t *testing.T) {
// 	eu, config, _, _, _ := NewTestEU()

// 	// ================================== Compile the contract ==================================
// 	currentPath, _ := os.Getwd()
// 	project := filepath.Dir(currentPath)
// 	code, err := compiler.CompileContracts(project+"/apps", "/storagenative/local_test.sol", "0.8.19", "LocalTest", false)
// 	if err != nil || len(code) == 0 {
// 		t.Error(err)
// 	}

// 	// ================================== Deploy the contract ==================================
// 	msg := core.NewMessage(eucommon.Alice, nil, 0, new(big.Int).SetUint64(0), 1e15, new(big.Int).SetUint64(1), evmcommon.Hex2Bytes(code), nil, true) // Build the message
// 	receipt, _, err := eu.Run(evmcommon.BytesToHash([]byte{1, 1, 1}), 1, &msg, ccEu.NewEVMBlockContext(config), ccEu.NewEVMTxContext(msg))           // Execute it
// 	_, transitions := eu.Api().Ccurl().ExportAll()

// 	// ---------------

// 	// t.Log("\n" + FormatTransitions(accesses))
// 	t.Log("\n" + eucommon.FormatTransitions(transitions))
// 	t.Log(receipt)
// 	// contractAddress := receipt.ContractAddress
// 	if receipt.Status != 1 || err != nil {
// 		t.Error("Error: Deployment failed!!!", err)
// 	}
// }
