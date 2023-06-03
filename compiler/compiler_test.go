package compiler

import (
	"fmt"
	"os"
	"testing"
)

func TestPythonContractCompiler(t *testing.T) {
	currentPath, _ := os.Getwd()

	fmt.Println(currentPath)

	path := currentPath
	version := "0.5.0"
	solfilename := "compiler_test.sol"
	contractName := "Example"

	bincode, err := CompileContracts(path, solfilename, version, contractName, false)
	if err != nil {
		fmt.Printf("reading contract err:%v\n", err)
		return
	}
	fmt.Printf("bytes:%v\n", bincode)
}

func TestEnsure(t *testing.T) {
	currentPath, _ := os.Getwd()
	ensureOutpath(currentPath)
}
func TestGetSolMeta(t *testing.T) {
	currentPath, _ := os.Getwd()
	solfilename := "compiler_test.sol"
	contractName, err := GetContractMeta(currentPath + "/" + solfilename)
	if err != nil {
		fmt.Printf("Get contract meta err:%v\n", err)
		return
	}
	fmt.Printf("contractName:%v\n", contractName)
}
