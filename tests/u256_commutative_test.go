package tests

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCumulativeU256(t *testing.T) {
	currentPath, _ := os.Getwd()
	targetPath := filepath.Dir(currentPath) + "/api/"

	err, _ := InvokeTestContract(targetPath, "commutative/u256/u256Cumulative_test.sol", "0.8.19", "CumulativeU256Test", "", []byte{}, false)
	if err != nil {
		t.Error(err)
	}
}

func TestCumulativeU256Case(t *testing.T) {
	currentPath, _ := os.Getwd()
	targetPath := filepath.Dir(currentPath) + "/api"
	err, _ := InvokeTestContract(targetPath, "commutative/u256/u256Cumulative_test.sol", "0.8.19", "ThreadingCumulativeU256", "call()", []byte{}, false)
	if err != nil {
		t.Error(err)
	}
}

func TestCumulativeU256Case1(t *testing.T) {
	currentPath, _ := os.Getwd()
	targetPath := filepath.Dir(currentPath) + "/api/"

	err, _ := InvokeTestContract(targetPath, "commutative/u256/u256Cumulative_test.sol", "0.8.19", "ThreadingCumulativeU256", "call1()", []byte{}, false)
	if err != nil {
		t.Error(err)
	}
}

func TestCumulativeU256ThreadingMultiTimes(t *testing.T) {
	currentPath, _ := os.Getwd()
	targetPath := filepath.Dir(currentPath) + "/api/"

	err, _ := InvokeTestContract(targetPath, "commutative/u256/u256Cumulative_test.sol", "0.8.19", "ThreadingCumulativeU256SameMpMulti", "call()", []byte{}, false)
	if err != nil {
		t.Error(err)
	}
}

func TestCumulativeU256ThreadingMulti(t *testing.T) {
	currentPath, _ := os.Getwd()
	targetPath := filepath.Dir(currentPath) + "/api/"

	err, _ := InvokeTestContract(targetPath, "commutative/u256/u256Cumulative_test.sol", "0.8.19", "ThreadingCumulativeU256DifferentMPMulti", "call()", []byte{}, false)
	if err != nil {
		t.Error(err)
	}
}
