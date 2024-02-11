package api

import (
	"fmt"

	"github.com/arcology-network/common-lib/codec"
	"github.com/arcology-network/vm-adaptor/abi"
	intf "github.com/arcology-network/vm-adaptor/interface"
	evmcommon "github.com/ethereum/go-ethereum/common"

	cache "github.com/arcology-network/eu/cache"
	"github.com/arcology-network/vm-adaptor/common"
	adaptorcommon "github.com/arcology-network/vm-adaptor/common"
)

type RuntimeHandlers struct {
	api       intf.EthApiRouter
	connector *adaptorcommon.PathBuilder
}

func NewRuntimeHandlers(ethApiRouter intf.EthApiRouter) *RuntimeHandlers {
	return &RuntimeHandlers{
		api:       ethApiRouter,
		connector: adaptorcommon.NewPathBuilder("/native/local/", ethApiRouter),
	}
}

func (this *RuntimeHandlers) Address() [20]byte {
	return common.RUNTIME_HANDLER
}

func (this *RuntimeHandlers) Call(caller, callee [20]byte, input []byte, origin [20]byte, nonce uint64) ([]byte, bool, int64) {
	signature := [4]byte{}
	copy(signature[:], input)

	switch signature {

	case [4]byte{0xf1, 0x06, 0x84, 0x54}: // 79 fc 09 a2
		return this.pid(caller, input[4:])

	case [4]byte{0x64, 0x23, 0xdb, 0x34}: // d3 01 e8 fe
		return this.rollback(caller, input[4:])

	case [4]byte{0xbb, 0x07, 0xe8, 0x5d}: // bb 07 e8 5d
		return this.uuid(caller, callee, input[4:])
	}

	fmt.Println(input)
	return []byte{}, false, 0
}

func (this *RuntimeHandlers) pid(caller evmcommon.Address, input []byte) ([]byte, bool, int64) {
	if encoded, err := abi.Encode(this.api.Pid()); err == nil {
		return encoded, true, 0
	}
	return []byte{}, false, 0
}

func (this *RuntimeHandlers) rollback(caller evmcommon.Address, input []byte) ([]byte, bool, int64) {
	cache.NewWriteCacheFilter(this.api.WriteCache()).RemoveByAddress(codec.Bytes20(caller).Hex())
	return []byte{}, true, 0
}

func (this *RuntimeHandlers) uuid(caller, callee evmcommon.Address, input []byte) ([]byte, bool, int64) {
	return this.api.ElementUID(), true, 0
}
