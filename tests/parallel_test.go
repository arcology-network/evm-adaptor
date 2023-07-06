package tests

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParallelBasic(t *testing.T) {
	currentPath, _ := os.Getwd()
	targetPath := filepath.Dir(currentPath) + "/api/"

	err, _ := InvokeTestContract(targetPath, "parallel/parallel_test.sol", "0.8.19", "ParaHasherTest", "call()", []byte{}, false)
	if err != nil {
		t.Error(err)
	}
}

func TestParallelWithConflict(t *testing.T) {
	currentPath, _ := os.Getwd()
	targetPath := filepath.Dir(currentPath) + "/api/"

	err, _ := InvokeTestContract(targetPath, "parallel/parallel_test.sol", "0.8.19", "ParaFixedLengthWithConflictTest", "call()", []byte{}, false)
	if err != nil {
		t.Error(err)
	}
}

// func TestParaContainerManipulation(t *testing.T) {
// 	currentPath, _ := os.Getwd()
// 	targetPath := filepath.Dir(currentPath) + "/api/"

// 	err, _ := InvokeTestContract(targetPath, "parallel/parallel_test.sol", "0.8.19", "ParaContainerManipulationTest", "call()", []byte{}, false)
// 	if err != nil {
// 		t.Error(err)
// 	}
// }

// func TestparallelMultiMPsTest(t *testing.T) {
// 	currentPath, _ := os.Getwd()
// 	targetPath := filepath.Dir(currentPath) + "/api/"

// 	err, _ := InvokeTestContract(targetPath, "parallel/parallel_test.sol", "0.8.19", "parallelMultiMPsTest", "call()", []byte{}, false)
// 	if err != nil {
// 		t.Error(err)
// 	}
// }

// func TestparallelMpArray(t *testing.T) {
// 	currentPath, _ := os.Getwd()
// 	targetPath := filepath.Dir(currentPath) + "/api/"

// 	err, _ := InvokeTestContract(targetPath, "parallel/parallel_test.sol", "0.8.19", "parallelMpArrayTest", "call()", []byte{}, false)
// 	if err != nil {
// 		t.Error(err)
// 	}
// }

// func TestRecursiveparallelNative(t *testing.T) {
// 	currentPath, _ := os.Getwd()
// 	targetPath := filepath.Dir(currentPath) + "/api/"

// 	err, _ := InvokeTestContract(targetPath, "parallel/parallel_test.sol", "0.8.19", "RecursiveparallelTest", "call()", []byte{}, false)
// 	if err != nil {
// 		t.Error(err)
// 	}
// }

// func TestRecursiveparallelContainer(t *testing.T) {
// 	currentPath, _ := os.Getwd()
// 	targetPath := filepath.Dir(currentPath) + "/api/"

// 	err, _ := InvokeTestContract(targetPath, "parallel/parallel_test.sol", "0.8.19", "parallelMpArraySubprocessTest", "call()", []byte{}, false)
// 	if err != nil {
// 		t.Error(err)
// 	}
// }

// func TestparallelDeployment(t *testing.T) {
// 	currentPath, _ := os.Getwd()
// 	targetPath := filepath.Dir(currentPath) + "/api/"

// 	err, _ := InvokeTestContract(targetPath, "parallel/parallel_test.sol", "0.8.19", "parallelDeploymentAddressTest", "call()", []byte{}, false)
// 	if err != nil {
// 		t.Error(err)
// 	}
// }
