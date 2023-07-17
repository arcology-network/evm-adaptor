package tests

import (
	"os"
	"path"
	"path/filepath"
	"testing"
)

func TestParallelBasic(t *testing.T) {
	currentPath, _ := os.Getwd()
	targetPath := path.Join(path.Dir(filepath.Dir(currentPath)), "concurrentlib/lib/")

	err, _ := InvokeTestContract(targetPath, "multiprocess/multiprocess_container_test.sol", "0.8.19", "ParaNativeAssignmentTest", "call()", []byte{}, false)
	if err != nil {
		t.Error(err)
	}
}

func TestParallelWithConflict(t *testing.T) {
	currentPath, _ := os.Getwd()
	targetPath := path.Join(path.Dir(filepath.Dir(currentPath)), "concurrentlib/lib/")

	err, _ := InvokeTestContract(targetPath, "multiprocess/multiprocess_container_test.sol", "0.8.19", "ParaFixedLengthWithConflictTest", "call()", []byte{}, false)
	if err != nil {
		t.Error(err)
	}
}

func TestParaFixedLengthWithConflictRollback(t *testing.T) {
	currentPath, _ := os.Getwd()
	targetPath := path.Join(path.Dir(filepath.Dir(currentPath)), "concurrentlib/lib/")

	err, _ := InvokeTestContract(targetPath, "multiprocess/multiprocess_container_test.sol", "0.8.19", "ParaFixedLengthWithConflictRollbackTest", "call()", []byte{}, false)
	if err != nil {
		t.Error(err)
	}
}

func TestMultiGlobalParaSingleInUse(t *testing.T) {
	currentPath, _ := os.Getwd()
	targetPath := path.Join(path.Dir(filepath.Dir(currentPath)), "concurrentlib/lib/")

	err, _ := InvokeTestContract(targetPath, "multiprocess/multiprocess_container_test.sol", "0.8.19", "MultiGlobalParaSingleInUse", "call()", []byte{}, false)
	if err != nil {
		t.Error(err)
	}
}

func TestMultiGlobalParaTest(t *testing.T) {
	currentPath, _ := os.Getwd()
	targetPath := path.Join(path.Dir(filepath.Dir(currentPath)), "concurrentlib/lib/")

	err, _ := InvokeTestContract(targetPath, "multiprocess/multiprocess_container_test.sol", "0.8.19", "MultiGlobalPara", "call()", []byte{}, false)
	if err != nil {
		t.Error(err)
	}
}

func TestMultiLocalPara(t *testing.T) {
	currentPath, _ := os.Getwd()
	targetPath := path.Join(path.Dir(filepath.Dir(currentPath)), "concurrentlib/lib/")

	err, _ := InvokeTestContract(targetPath, "multiprocess/multiprocess_container_test.sol", "0.8.19", "MultiTempParaTest", "call()", []byte{}, false)
	if err != nil {
		t.Error(err)
	}
}

func TestParaMultiWithClear(t *testing.T) {
	currentPath, _ := os.Getwd()
	targetPath := path.Join(path.Dir(filepath.Dir(currentPath)), "concurrentlib/lib/")

	err, _ := InvokeTestContract(targetPath, "multiprocess/multiprocess_container_test.sol", "0.8.19", "MultiLocalParaTestWithClear", "call()", []byte{}, false)
	if err != nil {
		t.Error(err)
	}
}

func TestMultiParaCumulativeU256(t *testing.T) {
	currentPath, _ := os.Getwd()
	targetPath := path.Join(path.Dir(filepath.Dir(currentPath)), "concurrentlib/lib/")

	err, _ := InvokeTestContract(targetPath, "multiprocess/multiprocess_container_test.sol", "0.8.19", "MultiParaCumulativeU256", "call()", []byte{}, false)
	if err != nil {
		t.Error(err)
	}
}

func TestParallelizerArray(t *testing.T) {
	currentPath, _ := os.Getwd()
	targetPath := path.Join(path.Dir(filepath.Dir(currentPath)), "concurrentlib/lib/")

	err, _ := InvokeTestContract(targetPath, "multiprocess/multiprocess_container_test.sol", "0.8.19", "ParallelizerArrayTest", "call()", []byte{}, false)
	if err != nil {
		t.Error(err)
	}
}

func TestMultipleParallelArray(t *testing.T) {
	currentPath, _ := os.Getwd()
	targetPath := path.Join(path.Dir(filepath.Dir(currentPath)), "concurrentlib/lib/")

	err, _ := InvokeTestContract(targetPath, "multiprocess/multiprocess_container_test.sol", "0.8.19", "MultiParaCumulativeU256", "call()", []byte{}, false)
	if err != nil {
		t.Error(err)
	}
}

func TestRecursiveParallelizerOnNativeArray(t *testing.T) {
	currentPath, _ := os.Getwd()
	targetPath := path.Join(path.Dir(filepath.Dir(currentPath)), "concurrentlib/lib/")

	err, _ := InvokeTestContract(targetPath, "multiprocess/multiprocess_container_test.sol", "0.8.19", "RecursiveParallelizerOnNativeArrayTest", "call()", []byte{}, false)
	if err != nil {
		t.Error(err)
	}
}

func TestRecursiveParallelizerOnContainer(t *testing.T) {
	currentPath, _ := os.Getwd()
	targetPath := path.Join(path.Dir(filepath.Dir(currentPath)), "concurrentlib/lib/")

	err, _ := InvokeTestContract(targetPath, "multiprocess/multiprocess_container_test.sol", "0.8.19", "RecursiveParallelizerOnContainerTest", "call()", []byte{}, false)
	if err != nil {
		t.Error(err)
	}
}

func TestMaxSelfRecursiveDepth4Test(t *testing.T) {
	currentPath, _ := os.Getwd()
	targetPath := path.Join(path.Dir(filepath.Dir(currentPath)), "concurrentlib/lib/")

	err, _ := InvokeTestContract(targetPath, "multiprocess/multiprocess_container_test.sol", "0.8.19", "MaxSelfRecursiveDepth4Test", "call()", []byte{}, false)
	if err != nil {
		t.Error(err)
	}
}

func TestMaxSelfRecursiveDepth(t *testing.T) {
	currentPath, _ := os.Getwd()
	targetPath := path.Join(path.Dir(filepath.Dir(currentPath)), "concurrentlib/lib/")

	err, _ := InvokeTestContract(targetPath, "multiprocess/multiprocess_container_test.sol", "0.8.19", "MaxRecursiveDepth4Test", "call()", []byte{}, false)
	if err != nil {
		t.Error(err)
	}
}

func TestMaxRecursiveDepthOffLimits(t *testing.T) {
	currentPath, _ := os.Getwd()
	targetPath := path.Join(path.Dir(filepath.Dir(currentPath)), "concurrentlib/lib/")

	err, _ := InvokeTestContract(targetPath, "multiprocess/multiprocess_container_test.sol", "0.8.19", "MaxRecursiveDepthOffLimitTest", "call()", []byte{}, false)
	if err != nil {
		t.Error(err)
	}
}

func TestCumulativeU256Recursive(t *testing.T) {
	currentPath, _ := os.Getwd()
	targetPath := path.Join(path.Dir(filepath.Dir(currentPath)), "concurrentlib/lib/")
	err, _ := InvokeTestContract(targetPath, "multiprocess/multiprocess_cum_test.sol", "0.8.19", "MixedRecursiveMultiprocessTest", "call()", []byte{}, false)
	if err != nil {
		t.Error(err)
	}
}

func TestCumulativeU256Case(t *testing.T) {
	currentPath, _ := os.Getwd()
	targetPath := path.Join(path.Dir(filepath.Dir(currentPath)), "concurrentlib/lib/")

	err, _ := InvokeTestContract(targetPath, "multiprocess/multiprocess_cum_test.sol", "0.8.19", "ParallelCumulativeU256", "call()", []byte{}, false)
	if err != nil {
		t.Error(err)
	}
}

func TestCumulativeU256Case1(t *testing.T) {
	currentPath, _ := os.Getwd()
	targetPath := path.Join(path.Dir(filepath.Dir(currentPath)), "concurrentlib/lib/")

	err, _ := InvokeTestContract(targetPath, "multiprocess/multiprocess_cum_test.sol", "0.8.19", "ParallelCumulativeU256", "call1()", []byte{}, false)
	if err != nil {
		t.Error(err)
	}
}

func TestCumulativeU256ThreadingMultiTimes(t *testing.T) {
	currentPath, _ := os.Getwd()
	targetPath := path.Join(path.Dir(filepath.Dir(currentPath)), "concurrentlib/lib/")

	err, _ := InvokeTestContract(targetPath, "multiprocess/multiprocess_cum_test.sol", "0.8.19", "ThreadingCumulativeU256SameMpMulti", "call()", []byte{}, false)
	if err != nil {
		t.Error(err)
	}
}