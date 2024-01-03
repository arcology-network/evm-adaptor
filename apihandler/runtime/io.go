package api

import (
	"fmt"

	"github.com/arcology-network/vm-adaptor/common"
	intf "github.com/arcology-network/vm-adaptor/interface"
)

// APIs under the concurrency namespace
type IoHandlers struct {
	api intf.EthApiRouter
}

func NewIoHandlers(api intf.EthApiRouter) *IoHandlers {
	return &IoHandlers{
		api: api,
	}
}
func (this *IoHandlers) Address() [20]byte {
	return common.IO_HANDLER
}

func (this *IoHandlers) Call(caller, callee [20]byte, input []byte, origin [20]byte, nonce uint64) ([]byte, bool, int64) {
	signature := [4]byte{}
	copy(signature[:], input)

	return this.print(caller, callee, input, origin, nonce)
	// }
	// return []byte{}, false, 0
}

func (this *IoHandlers) print(caller, callee [20]byte, input []byte, origin [20]byte, nonce uint64) ([]byte, bool, int64) {
	fmt.Println("caller:", caller)
	fmt.Println("input:", input)
	fmt.Println("origin:", origin)
	fmt.Println("nonce:", nonce)
	return []byte{}, true, 0
}
