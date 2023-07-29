package tests

import (
	"os"
	"path"
	"path/filepath"
	"testing"

	evmcommon "github.com/arcology-network/evm/common"
	eucommon "github.com/arcology-network/vm-adaptor/common"
	"github.com/arcology-network/vm-adaptor/tests"
)

func TestVote(t *testing.T) {
	currentPath, _ := os.Getwd()
	targetPath := path.Join((path.Dir(filepath.Dir(currentPath))), "concurrentlib/")

	contract, err := tests.NewContract(
		[]evmcommon.Address{
			eucommon.Coinbase,
			eucommon.Alice,
			eucommon.Bob,
			eucommon.Abby,
			eucommon.Abu,
			eucommon.Andy,
			eucommon.Anna,
			eucommon.Antonio,
			eucommon.Bailey,
			eucommon.Baloo,
			eucommon.Bambi,
			eucommon.Banza,
			eucommon.Beast,
		},
		eucommon.Alice,
		targetPath, "examples/vote/vote_test.sol", "0.8.19", "SerialBallotCaller")

	if err != nil {
		return
	}

	if _, err = contract.Deploy(eucommon.Alice, 0); err != nil {
		t.Error(err)
	}

	if _, err = contract.Invoke(eucommon.Alice, []byte{}); err != nil {
		t.Error(err)
	}
}
