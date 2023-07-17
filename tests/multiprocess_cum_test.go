package tests

import (
	"os"
	"path"
	"path/filepath"
	"testing"
)

func TestRecursiveCumulativeU256(t *testing.T) {
	currentPath, _ := os.Getwd()
	targetPath := path.Join(path.Dir(filepath.Dir(currentPath)), "concurrentlib") + "/lib/"
	err, _ := InvokeTestContract(targetPath, "multiprocess/multiprocess_cum_test.sol", "0.8.19", "MixedRecursiveMultiprocessTest", "call()", []byte{}, false)
	if err != nil {
		t.Error(err)
	}
}

func TestCumulativeU256Case(t *testing.T) {
	currentPath, _ := os.Getwd()
	targetPath := path.Join(path.Dir(filepath.Dir(currentPath)), "concurrentlib") + "/lib/"

	err, _ := InvokeTestContract(targetPath, "multiprocess/multiprocess_cum_test.sol", "0.8.19", "ParallelCumulativeU256", "call()", []byte{}, false)
	if err != nil {
		t.Error(err)
	}
}

func TestCumulativeU256Case1(t *testing.T) {
	currentPath, _ := os.Getwd()
	targetPath := path.Join(path.Dir(filepath.Dir(currentPath)), "concurrentlib") + "/lib/"

	err, _ := InvokeTestContract(targetPath, "multiprocess/multiprocess_cum_test.sol", "0.8.19", "ParallelCumulativeU256", "call1()", []byte{}, false)
	if err != nil {
		t.Error(err)
	}
}

func TestCumulativeU256ThreadingMultiTimes(t *testing.T) {
	currentPath, _ := os.Getwd()
	targetPath := path.Join(path.Dir(filepath.Dir(currentPath)), "concurrentlib") + "/lib/"

	err, _ := InvokeTestContract(targetPath, "multiprocess/multiprocess_cum_test.sol", "0.8.19", "ThreadingCumulativeU256SameMpMulti", "call()", []byte{}, false)
	if err != nil {
		t.Error(err)
	}
}
