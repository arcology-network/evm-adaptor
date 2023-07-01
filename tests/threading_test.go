package tests

import (
	"os"
	"path/filepath"
	"testing"
)

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

func TestRecursiveThreading(t *testing.T) {
	currentPath, _ := os.Getwd()
	targetPath := filepath.Dir(currentPath) + "/api/"

	err, _ := InvokeTestContract(targetPath, "threading/threading_test.sol", "0.8.19", "RecursiveThreadingTest", "call()", []byte{}, false)
	if err != nil {
		t.Error(err)
	}
}
