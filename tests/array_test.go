package tests

import (
	"os"
	"path"
	"path/filepath"
	"testing"
)

func TestAddressContainer(t *testing.T) {
	currentPath, _ := os.Getwd()
	targetPath := path.Join(path.Dir(filepath.Dir(currentPath)), "concurrentlib/lib/")
	err, _ := InvokeTestContract(targetPath, "array/address_test.sol", "0.8.19", "AddressTest", "", []byte{}, false)
	if err != nil {
		t.Error(err)
	}
}

func TestBoolContainer(t *testing.T) {
	currentPath, _ := os.Getwd()
	targetPath := path.Join(path.Dir(filepath.Dir(currentPath)), "concurrentlib/lib/")
	err, _ := InvokeTestContract(targetPath, "array/bool_test.sol", "0.8.19", "BoolTest", "", []byte{}, false)
	if err != nil {
		t.Error(err)
	}
}

func TestBytesContainer(t *testing.T) {
	currentPath, _ := os.Getwd()
	targetPath := path.Join(path.Dir(filepath.Dir(currentPath)), "concurrentlib/lib/")
	err, _ := InvokeTestContract(targetPath, "array/bytes_test.sol", "0.8.19", "ByteTest", "", []byte{}, false)
	if err != nil {
		t.Error(err)
	}
}

func TestContractBytes32(t *testing.T) {
	currentPath, _ := os.Getwd()
	targetPath := path.Join(path.Dir(filepath.Dir(currentPath)), "concurrentlib/lib/")
	err, _ := InvokeTestContract(targetPath, "array/bytes32_test.sol", "0.8.19", "Bytes32Test", "", []byte{}, false)
	if err != nil {
		t.Error(err)
	}
}

func TestContractNoncommutativeInt256(t *testing.T) {
	currentPath, _ := os.Getwd()
	targetPath := path.Join(path.Dir(filepath.Dir(currentPath)), "concurrentlib/lib/")
	err, _ := InvokeTestContract(targetPath, "array/int256_test.sol", "0.8.19", "Int256Test", "", []byte{}, false)
	if err != nil {
		t.Error(err)
	}
}

func TestContractString(t *testing.T) {
	currentPath, _ := os.Getwd()
	targetPath := path.Join(path.Dir(filepath.Dir(currentPath)), "concurrentlib/lib/")
	err, _ := InvokeTestContract(targetPath, "array/string_test.sol", "0.8.19", "StringTest", "", []byte{}, false)
	if err != nil {
		t.Error(err)
	}
}

func TestCumulativeU256(t *testing.T) {
	currentPath, _ := os.Getwd()
	targetPath := path.Join(path.Dir(filepath.Dir(currentPath)), "concurrentlib/lib/")

	err, _ := InvokeTestContract(targetPath, "array/u256Cum_test.sol", "0.8.19", "CumulativeU256Test", "", []byte{}, false)
	if err != nil {
		t.Error(err)
	}
}

func TestU256Dynamic(t *testing.T) {
	currentPath, _ := os.Getwd()
	targetPath := path.Join(path.Dir(filepath.Dir(currentPath)), "concurrentlib") + "/lib"
	err, _ := InvokeTestContract(targetPath, "/multiprocess/u256_mp_test.sol", "0.8.19", "U256ParallelTest", "", []byte{}, false)
	if err != nil {
		t.Error(err)
	}
}

func TestU256Multiprocess(t *testing.T) {
	currentPath, _ := os.Getwd()
	targetPath := path.Join(path.Dir(filepath.Dir(currentPath)), "concurrentlib") + "/lib"
	err, _ := InvokeTestContract(targetPath, "/multiprocess/u256_mp_test.sol", "0.8.19", "U256ParallelTest", "call()", []byte{}, false)
	if err != nil {
		t.Error(err)
	}
}

func TestArrayMultiprocess(t *testing.T) {
	currentPath, _ := os.Getwd()
	targetPath := path.Join(path.Dir(filepath.Dir(currentPath)), "concurrentlib") + "/lib"
	err, _ := InvokeTestContract(targetPath, "/multiprocess/u256_mp_test.sol", "0.8.19", "ArrayParallelTest", "call()", []byte{}, false)
	if err != nil {
		t.Error(err)
	}
}
