package tests

import (
	"os"
	"path"
	"path/filepath"
	"testing"
)

func TestCumulativeU256(t *testing.T) {
	currentPath, _ := os.Getwd()
	targetPath := path.Join(path.Dir(filepath.Dir(currentPath)), "concurrentlib") + "/lib/"

	err, _ := InvokeTestContract(targetPath, "u256cum/u256Cum_test.sol", "0.8.19", "CumulativeU256Test", "", []byte{}, false)
	if err != nil {
		t.Error(err)
	}
}
