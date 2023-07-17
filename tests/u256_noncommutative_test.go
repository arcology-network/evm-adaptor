package tests

import (
	"os"
	"path"
	"path/filepath"
	"testing"
)

func TestU256Dynamic(t *testing.T) {
	currentPath, _ := os.Getwd()
	targetPath := path.Join(path.Dir(filepath.Dir(currentPath)), "concurrentlib") + "/lib"
	err, _ := InvokeTestContract(targetPath, "/u256/u256_test.sol", "0.8.19", "U256DynamicTest", "", []byte{}, false)
	if err != nil {
		t.Error(err)
	}
}

func TestU256Multiprocess(t *testing.T) {
	currentPath, _ := os.Getwd()
	targetPath := path.Join(path.Dir(filepath.Dir(currentPath)), "concurrentlib") + "/lib"
	err, _ := InvokeTestContract(targetPath, "/u256/u256_test.sol", "0.8.19", "U256ParallelTest", "call()", []byte{}, false)
	if err != nil {
		t.Error(err)
	}
}

func TestArrayMultiprocess(t *testing.T) {
	currentPath, _ := os.Getwd()
	targetPath := path.Join(path.Dir(filepath.Dir(currentPath)), "concurrentlib") + "/lib"
	err, _ := InvokeTestContract(targetPath, "/u256/u256_test.sol", "0.8.19", "ArrayParallelTest", "call()", []byte{}, false)
	if err != nil {
		t.Error(err)
	}
}

func Test256N(t *testing.T) {
	currentPath, _ := os.Getwd()
	targetPath := path.Join(path.Dir(filepath.Dir(currentPath)), "concurrentlib") + "/lib"
	err, _ := InvokeTestContract(targetPath, "/u256/u256N_test.sol", "0.8.19", "U256NTest", "", []byte{}, false)
	if err != nil {
		t.Error(err)
	}
}
