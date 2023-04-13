package compiler

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/arcology-network/common-lib/common"
)

func CompileContracts(compiler, file, contract string) (string, error) {
	if !common.FileExists(compiler) {
		return "", errors.New("Error: Compiler doesn't exist in " + compiler)
	}

	if !common.FileExists(file) {
		return "", errors.New("Error: The contract file doesn't exist in " + file)
	}

	currentPath, _ := os.Getwd()
	fmt.Println(currentPath)
	if code, err := exec.Command("python", compiler, file, contract).Output(); err == nil && len(code) > 0 {
		return string(code), nil
	} else {
		return "", (err)
	}
}
