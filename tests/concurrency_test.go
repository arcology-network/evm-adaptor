package tests

// func TestConcurrencyWithThreading(t *testing.T) {
// 	eu, config, db, url, _ := NewTestEU()

// 	// ================================== Compile the contract ==================================
// 	currentPath, _ := os.Getwd()
// 	project := filepath.Dir(currentPath)
// 	pyCompiler := project + "/compiler/compiler.py"
// 	targetPath := project + "/api/concurrency/"

// 	if err := common.CopyFile(project+"/api/threading/Threading.sol", targetPath+"/Threading.sol"); err != nil {
// 		t.Error(err)
// 	}

// 	code, err := compiler.CompileContracts(pyCompiler, project+"/api/concurrency/concurrency_test.sol", "ConcurrencyDeferredInThreadingTest")
// 	if err != nil || len(code) == 0 {
// 		t.Error("Error: Failed to generate the byte code")
// 	}
// 	// ================================== Deploy the contract ==================================
// 	msg := types.NewMessage(eucommon.Alice, nil, 0, new(big.Int).SetUint64(0), 1e15, new(big.Int).SetUint64(1), evmcommon.Hex2Bytes(code), nil, false) // Build the message
// 	receipt, _, err := eu.Run(evmcommon.BytesToHash([]byte{1, 1, 1}), 1, &msg, cceu.NewEVMBlockContext(config), cceu.NewEVMTxContext(msg))             // Execute it
// 	_, transitions := eu.Api().Ccurl().ExportAll()

// 	if receipt.Status != 1 || err != nil {
// 		t.Error("Error: Deployment failed!!!", err)
// 	}
// 	fmt.Println(receipt.ContractAddress)

// 	url = concurrenturl.NewConcurrentUrl(db)
// 	url.Import(transitions)
// 	url.Sort()
// 	url.Commit([]uint32{1})

// 	receipt, _, err = eu.Run(evmcommon.BytesToHash([]byte{1, 1, 1}), 1, &msg, cceu.NewEVMBlockContext(config), cceu.NewEVMTxContext(msg))
// 	// _, transitions = eu.Api().Ccurl().ExportAll()

// 	if err != nil {
// 		fmt.Print(err)
// 	}

// 	contractAddress := receipt.ContractAddress
// 	data := crypto.Keccak256([]byte("call()"))[:4]
// 	msg = types.NewMessage(eucommon.Alice, &contractAddress, 1, new(big.Int).SetUint64(0), 1e15, new(big.Int).SetUint64(1), data, nil, false)
// 	receipt, execResult, err := eu.Run(evmcommon.BytesToHash([]byte{1, 1, 1}), 1, &msg, cceu.NewEVMBlockContext(config), cceu.NewEVMTxContext(msg))
// 	_, transitions = eu.Api().Ccurl().ExportAll()

// 	if err != nil {
// 		t.Error(err)
// 	}

// 	if execResult != nil && execResult.Err != nil {
// 		t.Error(execResult.Err)
// 	}

// 	if receipt.Status != 1 || err != nil {
// 		t.Error("Error: Failed to call!!!", err)
// 	}
// }

// func TestConcurrencyDeferredTest(t *testing.T) {
// 	eu, config, db, url, _ := NewTestEU()

// 	// ================================== Compile the contract ==================================
// 	currentPath, _ := os.Getwd()
// 	project := filepath.Dir(currentPath)
// 	pyCompiler := project + "/compiler/compiler.py"
// 	targetPath := project + "/api/concurrency/"

// 	if err := common.CopyFile(project+"/api/threading/Threading.sol", targetPath+"/Threading.sol"); err != nil {
// 		t.Error(err)
// 	}

// 	code, err := compiler.CompileContracts(pyCompiler, project+"/api/concurrency/concurrency_test.sol", "ConcurrencyDeferredTest")
// 	if err != nil || len(code) == 0 {
// 		t.Error("Error: Failed to generate the byte code")
// 	}
// 	// ================================== Deploy the contract ==================================
// 	msg := types.NewMessage(eucommon.Alice, nil, 0, new(big.Int).SetUint64(0), 1e15, new(big.Int).SetUint64(1), evmcommon.Hex2Bytes(code), nil, false) // Build the message
// 	receipt, _, err := eu.Run(evmcommon.BytesToHash([]byte{1, 1, 1}), 1, &msg, cceu.NewEVMBlockContext(config), cceu.NewEVMTxContext(msg))             // Execute it
// 	_, transitions := eu.Api().Ccurl().ExportAll()

// 	if receipt.Status != 1 || err != nil {
// 		t.Error("Error: Deployment failed!!!", err)
// 	}
// 	fmt.Println(receipt.ContractAddress)

// 	url = concurrenturl.NewConcurrentUrl(db)
// 	url.Import(transitions)
// 	url.Sort()
// 	url.Commit([]uint32{1})

// 	receipt, _, err = eu.Run(evmcommon.BytesToHash([]byte{1, 1, 1}), 1, &msg, cceu.NewEVMBlockContext(config), cceu.NewEVMTxContext(msg))
// 	// _, transitions = eu.Api().Ccurl().ExportAll()

// 	if err != nil {
// 		fmt.Print(err)
// 	}

// 	contractAddress := receipt.ContractAddress
// 	data := crypto.Keccak256([]byte("call()"))[:4]
// 	msg = types.NewMessage(eucommon.Alice, &contractAddress, 1, new(big.Int).SetUint64(0), 1e15, new(big.Int).SetUint64(1), data, nil, false)
// 	receipt, execResult, err := eu.Run(evmcommon.BytesToHash([]byte{1, 1, 1}), 1, &msg, cceu.NewEVMBlockContext(config), cceu.NewEVMTxContext(msg))
// 	_, transitions = eu.Api().Ccurl().ExportAll()

// 	if err != nil {
// 		t.Error(err)
// 	}

// 	if execResult != nil && execResult.Err != nil {
// 		t.Error(execResult.Err)
// 	}

// 	if receipt.Status != 1 || err != nil {
// 		t.Error("Error: Failed to call!!!", err)
// 	}
// }
