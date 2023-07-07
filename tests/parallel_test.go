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

func TestParaContainerManipulation(t *testing.T) {
	currentPath, _ := os.Getwd()
	targetPath := filepath.Dir(currentPath) + "/api/"

	err, _ := InvokeTestContract(targetPath, "parallel/parallel_test.sol", "0.8.19", "ParaContainerConcurrentPushTest", "call()", []byte{}, false)
	if err != nil {
		t.Error(err)
	}
}

func TestParaMulti(t *testing.T) {
	currentPath, _ := os.Getwd()
	targetPath := filepath.Dir(currentPath) + "/api/"

	err, _ := InvokeTestContract(targetPath, "parallel/parallel_test.sol", "0.8.19", "MultiParaTest", "call()", []byte{}, false)
	if err != nil {
		t.Error(err)
	}
}

func TestParaMultiWithClear(t *testing.T) {
	currentPath, _ := os.Getwd()
	targetPath := filepath.Dir(currentPath) + "/api/"

	err, _ := InvokeTestContract(targetPath, "parallel/parallel_test.sol", "0.8.19", "MultiParaTestWithClear", "call()", []byte{}, false)
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

func TestRecursiveParallelizerOnContainer(t *testing.T) {
	currentPath, _ := os.Getwd()
	targetPath := filepath.Dir(currentPath) + "/api/"

	err, _ := InvokeTestContract(targetPath, "parallel/parallel_test.sol", "0.8.19", "RecursiveParallelizerOnContainerTest", "call()", []byte{}, false)
	if err != nil {
		t.Error(err)
	}
}

// func TestRecursiveParallelizer(t *testing.T) {
// 	currentPath, _ := os.Getwd()
// 	targetPath := filepath.Dir(currentPath) + "/api/"

// 	err, _ := InvokeTestContract(targetPath, "parallel/parallel_test.sol", "0.8.19", "RecursiveThreadingTest", "call()", []byte{}, false)
// 	if err != nil {
// 		t.Error(err)
// 	}
// }
