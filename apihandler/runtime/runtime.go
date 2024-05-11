package runtime

import (
	"encoding/hex"
	"fmt"
	"math"

	"github.com/arcology-network/common-lib/codec"
	"github.com/arcology-network/common-lib/exp/slice"
	"github.com/arcology-network/evm-adaptor/abi"
	intf "github.com/arcology-network/evm-adaptor/interface"
	evmcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/holiman/uint256"

	scheduler "github.com/arcology-network/eu/new-scheduler"
	"github.com/arcology-network/evm-adaptor/common"
	adaptorcommon "github.com/arcology-network/evm-adaptor/common"
	pathbuilder "github.com/arcology-network/evm-adaptor/pathbuilder"
	"github.com/arcology-network/storage-committer/commutative"
	cache "github.com/arcology-network/storage-committer/storage/writecache"
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

	case [4]byte{0x83, 0x2e, 0x49, 0xcf}: //
		return this.sequentializeAll(caller, callee, input[4:])

	case [4]byte{0xc4, 0xdf, 0xfe, 0x6e}: //
		return this.sequentialize(caller, callee, input[4:])

	case [4]byte{0x68, 0x7b, 0x09, 0xb7}: //
		return this.parallelize(caller, callee, input[4:])

	case [4]byte{0x68, 0x7b, 0x09, 0xb7}: //
		return this.parallelize(caller, callee, input[4:])

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

// This function rolls back the storage to the previous generation. It should be used with extreme caution.
func (this *RuntimeHandlers) rollback(caller evmcommon.Address, input []byte) ([]byte, bool, int64) {
	cache.NewWriteCacheFilter(this.api.WriteCache()).RemoveByAddress(codec.Bytes20(caller).Hex())
	return []byte{}, true, 0
}

func (this *RuntimeHandlers) uuid(caller, callee evmcommon.Address, input []byte) ([]byte, bool, int64) {
	return this.api.ElementUID(), true, 0
}

// Get the number of running instances of a function.
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

func (this *RuntimeHandlers) sequentializeAll(caller, _ evmcommon.Address, input []byte) ([]byte, bool, int64) {
	if !this.api.VM().(*vm.EVM).ArcologyNetworkAPIs.IsInConstructor() {
		return []byte{}, false, 0 // Can only be called from a constructor.
	}

	cache := this.api.WriteCache().(*cache.WriteCache)
	txID := this.api.GetEU().(interface{ ID() uint32 }).ID()

	signatures, err := abi.DecodeTo(input[32:], 1, [][4]byte{}, 2, math.MaxInt)
	for _, signature := range signatures {
		path := pathbuilder.FuncPropertyPath(caller, signature)
		_, err = cache.Write(txID, path, commutative.NewPath())
		if err != nil {
			return []byte{}, false, 0 // Create a new sequentializer path.
		}
	}
	return []byte{}, true, 0 // Create a new sequentializer path.
}

func (this *RuntimeHandlers) sequentialize(caller, _ evmcommon.Address, input []byte) ([]byte, bool, int64) {
	if !this.api.VM().(*vm.EVM).ArcologyNetworkAPIs.IsInConstructor() {
		return []byte{}, false, 0 // Can only be called from a constructor.
	}

	// Get the target contract address.
	sourceFunc, err := abi.DecodeTo(input, 0, [4]byte{}, 1, 4)
	if err != nil {
		return []byte{}, false, 0
	}

	targetAddr, err := abi.DecodeTo(input, 1, [20]byte{}, 1, math.MaxInt)
	if err != nil {
		return []byte{}, false, 0
	}

	// Get the target function signatures
	signBytes, err := abi.DecodeTo(input, 2, []byte{}, 1, math.MaxInt)
	if err != nil || len(signBytes) <= 32 {
		return []byte{}, false, 0
	}

	// Parse the function signatures.
	signatures, err := abi.DecodeTo(signBytes[32:], 2, [][4]byte{}, 2, math.MaxInt)

	cache := this.api.WriteCache().(*cache.WriteCache)
	txID := this.api.GetEU().(interface{ ID() uint32 }).ID()

	parentPath := pathbuilder.FuncPropertyPath(caller, sourceFunc)
	if _, err = cache.Write(txID, parentPath, commutative.NewPath()); err != nil { // Create a new sequentializer path.
		return []byte{}, err == nil, 0
	}

	callees := slice.Transform(signatures, func(i int, signature [4]byte) string {
		return hex.EncodeToString(new(scheduler.Callee).Compact(targetAddr[:], signature[:]))
	})

	path := pathbuilder.SequentializerPath(caller, sourceFunc)
	_, err = cache.Write(txID, path, commutative.NewPath(callees...)) // Write the sequentializer path regardless of its existence.
	return []byte{}, err == nil, 0
}

func (this *RuntimeHandlers) parallelize(caller, callee evmcommon.Address, input []byte) ([]byte, bool, int64) {
	fmt.Println(input)

	if !this.api.VM().(*vm.EVM).ArcologyNetworkAPIs.IsInConstructor() {
		return []byte{}, false, 0 // Can only be called from a constructor.
	}
	return []byte{}, true, 0 // Can only initialize once.
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
