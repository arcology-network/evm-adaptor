package types

import (
	"bytes"

	common "github.com/arcology-network/common-lib/common"
	"github.com/arcology-network/vm-adaptor/interfaces"
	"golang.org/x/crypto/sha3"
)

type DeferredCall struct {
	From         [20]byte
	Signature    [32]byte // Signature
	FuncCallData []byte
}

func NewDeferredCall(from [20]byte, api interfaces.EthApiRouter) *DeferredCall {
	hierarchy := api.VM().ArcologyNetworkAPIs.CallHierarchy()
	signature := sha3.Sum256(common.Flatten(common.Reverse[[]byte](&hierarchy)))

	return &DeferredCall{
		from,
		signature,
		bytes.Clone(api.VM().ArcologyNetworkAPIs.CallContext.Contract.Input),
	}
}
