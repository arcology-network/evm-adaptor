package tests

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLocalizer(t *testing.T) {
	currentPath, _ := os.Getwd()
	targetPath := filepath.Dir(currentPath) + "/api/"

	err, _ := InvokeTestContract(targetPath, "runtime/runtime_test.sol", "0.8.19", "Deployee", "check()", []byte{}, false)
	if err != nil {
		t.Error(err)
	}
}

func TestDeployer(t *testing.T) {
	currentPath, _ := os.Getwd()
	targetPath := filepath.Dir(currentPath) + "/api/"

	err, _ := InvokeTestContract(targetPath, "runtime/runtime_test.sol", "0.8.19", "Deployer", "", []byte{1}, false)
	if err != nil {
		t.Error(err)
	}
}

// func TestDeployer2(t *testing.T) {
// 	currentPath, _ := os.Getwd()
// 	targetPath := filepath.Dir(currentPath) + "/api/"

// 	err, _ := InvokeTestContract(targetPath, "runtime/runtime_test.sol", "0.8.19", "Deployer2", "", []byte{}, false)
// 	if err != nil {
// 		t.Error(err)
// 	}
// }
