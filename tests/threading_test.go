package tests

import (
	"os"
	"path/filepath"
	"testing"
)

func TestThreadsWithConflict(t *testing.T) {
	currentPath, _ := os.Getwd()
	targetPath := filepath.Dir(currentPath) + "/api/"

	err, _ := InvokeTestContract(targetPath, "threading/threading_test.sol", "0.8.19", "ThreadingFixedLengthWithConflictTest", "call()", []byte{}, false)
	if err != nil {
		t.Error(err)
	}
}

func TestThreadingBasic(t *testing.T) {
	currentPath, _ := os.Getwd()
	targetPath := filepath.Dir(currentPath) + "/api/"

	err, _ := InvokeTestContract(targetPath, "threading/threading_test.sol", "0.8.19", "ThreadingParaHasherTest", "call()", []byte{}, false)
	if err != nil {
		t.Error(err)
	}
}

func TestParaContainerManipulation(t *testing.T) {
	currentPath, _ := os.Getwd()
	targetPath := filepath.Dir(currentPath) + "/api/"

	err, _ := InvokeTestContract(targetPath, "threading/threading_test.sol", "0.8.19", "ThreadingParaContainerManipulationTest", "call()", []byte{}, false)
	if err != nil {
		t.Error(err)
	}
}

func TestThreadingMultiMPsTest(t *testing.T) {
	currentPath, _ := os.Getwd()
	targetPath := filepath.Dir(currentPath) + "/api/"

	err, _ := InvokeTestContract(targetPath, "threading/threading_test.sol", "0.8.19", "ThreadingMultiMPsTest", "call()", []byte{}, false)
	if err != nil {
		t.Error(err)
	}
}

func TestThreadingMpArray(t *testing.T) {
	currentPath, _ := os.Getwd()
	targetPath := filepath.Dir(currentPath) + "/api/"

	err, _ := InvokeTestContract(targetPath, "threading/threading_test.sol", "0.8.19", "ThreadingMpArrayTest", "call()", []byte{}, false)
	if err != nil {
		t.Error(err)
	}
}

func TestRecursiveThreadingNative(t *testing.T) {
	currentPath, _ := os.Getwd()
	targetPath := filepath.Dir(currentPath) + "/api/"

	err, _ := InvokeTestContract(targetPath, "threading/threading_test.sol", "0.8.19", "RecursiveThreadingTest", "call()", []byte{}, false)
	if err != nil {
		t.Error(err)
	}
}

func TestRecursiveThreadingContainer(t *testing.T) {
	currentPath, _ := os.Getwd()
	targetPath := filepath.Dir(currentPath) + "/api/"

	err, _ := InvokeTestContract(targetPath, "threading/threading_test.sol", "0.8.19", "ThreadingConflictTest", "call()", []byte{}, false)
	if err != nil {
		t.Error(err)
	}
}
