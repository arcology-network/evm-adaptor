package tests

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPythonContractCompiler(t *testing.T) {
	currentPath, _ := os.Getwd()
	compiler := filepath.Dir(currentPath) + "/tests/compiler.py"
	if code, err := CompileContracts(compiler, "./compiler_test.sol", "Example"); err != nil || len(code) == 0 {
		t.Error(err)
	}
}
