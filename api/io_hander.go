package api

import (
	"fmt"

	"github.com/arcology-network/vm-adaptor/common"
	adaptorcommon "github.com/arcology-network/vm-adaptor/common"
)

// APIs under the concurrency namespace
type IoHandlers struct {
	api adaptorcommon.EthApiRouter
}

func NewIoHandlers(api adaptorcommon.EthApiRouter) *IoHandlers {
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

	// switch signature {
	// case [4]byte{0x13, 0xbd, 0xfa, 0xcd}: // 13 bd fa cd
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
