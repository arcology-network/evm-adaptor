package execution

import (
	common "github.com/arcology-network/common-lib/common"
	eucommon "github.com/arcology-network/vm-adaptor/common"
	"golang.org/x/crypto/sha3"
)

type DeferredCall struct {
	From         [20]byte
	Signature    [32]byte // Signature == callstack + function name
	Addr         [20]byte
	FuncCallData []byte // function Signature
	gasLimit     uint64
	// gasPrice
}

func NewDeferredCall(gasLimit uint64, from, targetAddr [20]byte, input []byte, api eucommon.EthApiRouter) *DeferredCall {
	hierarchy := api.VM().ArcologyNetworkAPIs.CallHierarchy()
	signature := sha3.Sum256(common.Flatten(common.Reverse[[]byte](&hierarchy)))

	return &DeferredCall{
		from,
		signature,
		targetAddr,
		input,
		gasLimit,
	}
}
