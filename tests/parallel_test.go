package tests

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParallelBasic(t *testing.T) {
	currentPath, _ := os.Getwd()
	targetPath := filepath.Dir(currentPath) + "/api/"

	err, _ := InvokeTestContract(targetPath, "parallel/parallel_test.sol", "0.8.19", "ParaNativeAssignmentTest", "call()", []byte{}, false)
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

func TestMultiLocalPara(t *testing.T) {
	currentPath, _ := os.Getwd()
	targetPath := filepath.Dir(currentPath) + "/api/"

	err, _ := InvokeTestContract(targetPath, "parallel/parallel_test.sol", "0.8.19", "MultiTempParaTest", "call()", []byte{}, false)
	if err != nil {
		t.Error(err)
	}
}

func TestParaMultiWithClear(t *testing.T) {
	currentPath, _ := os.Getwd()
	targetPath := filepath.Dir(currentPath) + "/api/"

	err, _ := InvokeTestContract(targetPath, "parallel/parallel_test.sol", "0.8.19", "MultiLocalParaTestWithClear", "call()", []byte{}, false)
	if err != nil {
		t.Error(err)
	}
}

func TestMultiParaCumulativeU256(t *testing.T) {
	currentPath, _ := os.Getwd()
	targetPath := filepath.Dir(currentPath) + "/api/"

	err, _ := InvokeTestContract(targetPath, "parallel/parallel_test.sol", "0.8.19", "MultiParaCumulativeU256", "call()", []byte{}, false)
	if err != nil {
		t.Error(err)
	}
}

func TestParallelizerArray(t *testing.T) {
	currentPath, _ := os.Getwd()
	targetPath := filepath.Dir(currentPath) + "/api/"

	err, _ := InvokeTestContract(targetPath, "parallel/parallel_test.sol", "0.8.19", "ParallelizerArrayTest", "call()", []byte{}, false)
	if err != nil {
		t.Error(err)
	}
}

func TestMultipleParallelArray(t *testing.T) {
	currentPath, _ := os.Getwd()
	targetPath := filepath.Dir(currentPath) + "/api/"

	err, _ := InvokeTestContract(targetPath, "parallel/parallel_test.sol", "0.8.19", "MultiParaCumulativeU256", "call()", []byte{}, false)
	if err != nil {
		t.Error(err)
	}
}

func TestRecursiveParallelizerOnNativeArray(t *testing.T) {
	currentPath, _ := os.Getwd()
	targetPath := filepath.Dir(currentPath) + "/api/"

	err, _ := InvokeTestContract(targetPath, "parallel/parallel_test.sol", "0.8.19", "RecursiveParallelizerOnNativeArrayTest", "call()", []byte{}, false)
	if err != nil {
		t.Error(err)
	}
}

func TestRecursiveParallelizerOnContainer(t *testing.T) {
	currentPath, _ := os.Getwd()
	targetPath := filepath.Dir(currentPath) + "/api/"

	err, _ := InvokeTestContract(targetPath, "parallel/parallel_test.sol", "0.8.19", "RecursiveParallelizerOnContainerTest", "call()", []byte{}, false)
	if err != nil {
		t.Error(err)
	}
}

// func TestMixedRecursiveParallelizer(t *testing.T) {
// 	currentPath, _ := os.Getwd()
// 	targetPath := filepath.Dir(currentPath) + "/api/"

// 	err, _ := InvokeTestContract(targetPath, "parallel/parallel_test.sol", "0.8.19", "MixedRecursiveParallelizerTest", "call()", []byte{}, false)
// 	if err != nil {
// 		t.Error(err)
// 	}
// }

// func TestParaContainerManipulation(t *testing.T) {
// 	currentPath, _ := os.Getwd()
// 	targetPath := filepath.Dir(currentPath) + "/api/"

// 	err, _ := InvokeTestContract(targetPath, "parallel/parallel_test.sol", "0.8.19", "ParaContainerConcurrentPushTest", "call()", []byte{}, false)
// 	if err != nil {
// 		t.Error(err)
// 	}
// }

// func TestMultiGlobalParaTest(t *testing.T) {
// 	currentPath, _ := os.Getwd()
// 	targetPath := filepath.Dir(currentPath) + "/api/"

// 	err, _ := InvokeTestContract(targetPath, "parallel/parallel_test.sol", "0.8.19", "MultiGlobalPara", "call()", []byte{}, false)
// 	if err != nil {
// 		t.Error(err)
// 	}
// }

func TestForeachRun(t *testing.T) {
	currentPath, _ := os.Getwd()
	targetPath := filepath.Dir(currentPath) + "/api/"

	err, _ := InvokeTestContract(targetPath, "parallel/parallel_test.sol", "0.8.19", "ForeachTest", "call()", []byte{}, false)
	if err != nil {
		t.Error(err)
	}
}
