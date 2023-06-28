package tests

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNoncommutativeU256Dynamic(t *testing.T) {
	currentPath, _ := os.Getwd()
	targetPath := filepath.Dir(currentPath) + "/api"
	err, _ := InvokeTestContract(targetPath, "noncommutative/u256/u256_test.sol", "0.8.19", "U256DynamicTest", "", []byte{}, false)
	if err != nil {
		t.Error(err)
	}
}

func TestNonCommutativeU256Threading(t *testing.T) {
	currentPath, _ := os.Getwd()
	targetPath := filepath.Dir(currentPath) + "/api"
	err, _ := InvokeTestContract(targetPath, "noncommutative/u256/u256_test.sol", "0.8.19", "U256ThreadingTest", "call()", []byte{}, false)
	if err != nil {
		t.Error(err)
	}
}

func TestArrayThreading(t *testing.T) {
	currentPath, _ := os.Getwd()
	targetPath := filepath.Dir(currentPath) + "/api"
	err, _ := InvokeTestContract(targetPath, "noncommutative/u256/u256_test.sol", "0.8.19", "ArrayThreadingTest", "call()", []byte{}, false)
	if err != nil {
		t.Error(err)
	}
}

func TestNoncommutative256N(t *testing.T) {
	currentPath, _ := os.Getwd()
	targetPath := filepath.Dir(currentPath) + "/api"
	err, _ := InvokeTestContract(targetPath, "noncommutative/u256/u256N_test.sol", "0.8.19", "U256NTest", "", []byte{}, false)
	if err != nil {
		t.Error(err)
	}
}
