package tests

import (
	"os"
	"path/filepath"
	"testing"
)

func TestAtomicDeferredInThreading(t *testing.T) {
	currentPath, _ := os.Getwd()
	targetPath := filepath.Dir(currentPath) + "/api/"

	err, _ := InvokeTestContract(targetPath, "atomic/atomic_test.sol", "0.8.19", "AtomicDeferredInThreadingTest", "call()", []byte{}, false)
	if err != nil {
		t.Error(err)
	}
}

func TestConflictInThreads(t *testing.T) {
	currentPath, _ := os.Getwd()
	targetPath := filepath.Dir(currentPath) + "/api/"

	err, _ := InvokeTestContract(targetPath, "atomic/atomic_test.sol", "0.8.19", "ConflictInThreadsFixedLengthTest", "call()", []byte{}, false)
	if err != nil {
		t.Error(err)
	}
}
