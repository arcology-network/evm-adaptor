package tests

import (
	"os"
	"path"
	"path/filepath"
	"testing"
)

func TestContainerPair(t *testing.T) {
	currentPath, _ := os.Getwd()
	targetPath := path.Join(path.Dir(filepath.Dir(currentPath)), "concurrentlib/lib/")

	err, _ := InvokeTestContract(targetPath, "combo/bytes_bool_test.sol", "0.8.19", "PairTest", "", []byte{}, false)
	if err != nil {
		t.Error(err)
	}
}

func TestBasicSet(t *testing.T) {
	currentPath, _ := os.Getwd()
	targetPath := path.Join(path.Dir(filepath.Dir(currentPath)), "concurrentlib/lib/")

	err, _ := InvokeTestContract(targetPath, "combo/set_test.sol", "0.8.19", "SetTest", "", []byte{}, false)
	if err != nil {
		t.Error(err)
	}
}

func TestBasicMap(t *testing.T) {
	currentPath, _ := os.Getwd()
	targetPath := path.Join(path.Dir(filepath.Dir(currentPath)), "concurrentlib/lib/")

	err, _ := InvokeTestContract(targetPath, "combo/map_test.sol", "0.8.19", "MapTest", "", []byte{}, false)
	if err != nil {
		t.Error(err)
	}
}

func TestConcurrentMap(t *testing.T) {
	currentPath, _ := os.Getwd()
	targetPath := path.Join(path.Dir(filepath.Dir(currentPath)), "concurrentlib/lib/")

	err, _ := InvokeTestContract(targetPath, "combo/map_test.sol", "0.8.19", "ConcurrenctMapTest", "", []byte{}, false)
	if err != nil {
		t.Error(err)
	}
}
