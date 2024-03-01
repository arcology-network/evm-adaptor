package runtime

import (
	"fmt"

	"github.com/arcology-network/common-lib/codec"
	"github.com/arcology-network/evm-adaptor/abi"
	intf "github.com/arcology-network/evm-adaptor/interface"
	evmcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/holiman/uint256"

	cache "github.com/arcology-network/eu/cache"
	scheduler "github.com/arcology-network/eu/new-scheduler"
	"github.com/arcology-network/evm-adaptor/common"
	adaptorcommon "github.com/arcology-network/evm-adaptor/common"
)

type RuntimeHandlers struct {
	api         intf.EthApiRouter
	pathBuilder *adaptorcommon.PathBuilder
}

func NewRuntimeHandlers(ethApiRouter intf.EthApiRouter) *RuntimeHandlers {
	return &RuntimeHandlers{
		api:         ethApiRouter,
		pathBuilder: adaptorcommon.NewPathBuilder("/storage", ethApiRouter),
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

	case [4]byte{0xa8, 0x7a, 0xe4, 0x81}: // bb 07 e8 5d
		return this.instances(caller, callee, input[4:])

	case [4]byte{0xf5, 0xf0, 0x15, 0xf3}: //
		return this.deferred(caller, callee, input[4:])
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

func (this *RuntimeHandlers) instances(caller evmcommon.Address, callee evmcommon.Address, input []byte) ([]byte, bool, int64) {
	if this.api.GetSchedule() == nil {
		return []byte{}, false, 0
	}

	address, err := abi.DecodeTo(input, 0, [20]byte{}, 1, 4)
	if err != nil {
		return []byte{}, false, 0
	}

	funSign, err := abi.DecodeTo(input, 1, []byte{}, 1, 4)
	if err != nil {
		return []byte{}, false, 0
	}

	dict := this.api.GetSchedule().(*map[string]int)
	key := scheduler.CallToKey(address[:], funSign)

	// Encode the total number of instances and return
	if encoded, err := abi.Encode(uint256.NewInt(uint64((*dict)[key]))); err == nil {
		// encoded, _ := abi.Encode(uint256.NewInt(2))
		// if !bytes.Equal(encoded, encoded2) {
		// 	panic("")
		// }

		return encoded, true, 0
	}
	return []byte{}, false, 0
}

// This function needs to schedule a defer call to the next generation.
func (this *RuntimeHandlers) deferred(caller, callee evmcommon.Address, input []byte) ([]byte, bool, int64) {
	if !this.api.VM().(*vm.EVM).ArcologyNetworkAPIs.IsInConstructor() {
		return []byte{}, false, 0 // Can only be called from a constructor.
	}

	if len(input) >= 4 {
		this.api.AuxDict()["deferrable"] = input[:4]
	}
	return []byte{}, false, 0 // Can only initialize once.
}
