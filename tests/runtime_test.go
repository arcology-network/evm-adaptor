package tests

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResettable(t *testing.T) {
	currentPath, _ := os.Getwd()
	targetPath := filepath.Dir(currentPath) + "/api/"

	err, _ := InvokeTestContract(targetPath, "runtime/runtime_test.sol", "0.8.19", "ResettableDeployer", "afterCheck()", []byte{}, false)
	if err != nil {
		t.Error(err)
	}
}
