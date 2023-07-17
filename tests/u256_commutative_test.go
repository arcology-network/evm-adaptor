package tests

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCumulativeU256(t *testing.T) {
	currentPath, _ := os.Getwd()
	targetPath := filepath.Dir(currentPath) + "/lib/"

	err, _ := InvokeTestContract(targetPath, "u256cum/u256Cum_test.sol", "0.8.19", "CumulativeU256Test", "", []byte{}, false)
	if err != nil {
		t.Error(err)
	}
}
