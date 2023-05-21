package tests

// func TestLocalStructure(t *testing.T) {
// 	eu, config, _, _, _ := NewTestEU()

// 	// ================================== Compile the contract ==================================
// 	currentPath, _ := os.Getwd()
// 	project := filepath.Dir(currentPath)
// 	pyCompiler := project + "/compiler/compiler.py"
// 	targetPath := project + "/api/"
// 	// baseFile := targetPath + "base/Base.sol"

// 	// if err := common.CopyFile(project+"/api/noncommutative/base/Base.sol", targetPath+"bool/Base.sol"); err != nil {
// 	// 	t.Error(err)
// 	// }

// 	// if err := common.CopyFile(project+"/api/threading/Threading.sol", targetPath+"/bool/Threading.sol"); err != nil {
// 	// 	t.Error(err)
// 	// }

// 	code, err := compiler.CompileContracts(pyCompiler, targetPath+"/local/local_test.sol", "LocalTest")
// 	if err != nil || len(code) == 0 {
// 		t.Error(err)
// 	}

// 	// ================================== Deploy the contract ==================================
// 	msg := types.NewMessage(eucommon.Alice, nil, 0, new(big.Int).SetUint64(0), 1e15, new(big.Int).SetUint64(1), evmcommon.Hex2Bytes(code), nil, true) // Build the message
// 	receipt, _, err := eu.Run(evmcommon.BytesToHash([]byte{1, 1, 1}), 1, &msg, ccEu.NewEVMBlockContext(config), ccEu.NewEVMTxContext(msg))            // Execute it
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
