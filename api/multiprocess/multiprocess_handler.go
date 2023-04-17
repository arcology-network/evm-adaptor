package multiprocess

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math"

	evmcommon "github.com/arcology-network/evm/common"
	"github.com/arcology-network/vm-adaptor/abi"
	eucommon "github.com/arcology-network/vm-adaptor/common"
)

// APIs under the concurrency namespace
type MultiprocessHandler struct {
	api        eucommon.ConcurrentApiRouterInterface
	jobManager *JobManager
}

func NewParallelHandler(apiRounter eucommon.ConcurrentApiRouterInterface) *MultiprocessHandler {
	return &MultiprocessHandler{
		api:        apiRounter,
		jobManager: NewJobManager(apiRounter),
	}
}

func (this *MultiprocessHandler) Address() [20]byte {
	return [20]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0x90}
}

func (this *MultiprocessHandler) Call(caller, callee evmcommon.Address, input []byte, origin evmcommon.Address, nonce uint64) ([]byte, bool) {
	signature := [4]byte{}
	copy(signature[:], input)

	fmt.Println(input)
	fmt.Println("==============================")

	switch signature { // bf 22 6c 78
	case [4]byte{0xa4, 0x62, 0x12, 0x2d}: // a4 62 12 2d
		return this.addJob(caller, callee, input[4:])

	case [4]byte{0xb6, 0xff, 0x8b, 0xd9}:
		return this.delJob(caller, callee, input[4:])

	case [4]byte{0xc0, 0x40, 0x62, 0x26}:
		return this.run(caller, callee, input[4:])

	case [4]byte{0x64, 0x17, 0x43, 0x08}:
		return this.clear(caller, callee, input[4:])

	case [4]byte{0x1f, 0x7b, 0x6d, 0x32}:
		return this.length()
	}
	return this.unknow(caller, callee, input)
}

func (this *MultiprocessHandler) length() ([]byte, bool) {
	if v, err := abi.Encode(uint64(len(this.jobManager.jobs))); err == nil {
		return v, true
	}
	return []byte{}, false
}

func (this *MultiprocessHandler) unknow(caller, callee evmcommon.Address, input []byte) ([]byte, bool) {
	this.api.AddLog("Unhandled function call in cumulative handler router", hex.EncodeToString(input))
	return []byte{}, false
}

func (this *MultiprocessHandler) run(caller, callee evmcommon.Address, input []byte) ([]byte, bool) {
	this.jobManager.Start()
	return []byte{}, true
}

func (this *MultiprocessHandler) addJob(caller, callee evmcommon.Address, input []byte) ([]byte, bool) {
	if len(input) < 4 {
		return []byte(errors.New("Error: Invalid input").Error()), false
	}

	fmt.Println(input)
	rawAddr, err := abi.DecodeTo(input, 0, [20]byte{}, 1, 32)
	if err != nil {
		return []byte(err.Error()), false
	}
	calleeAddr := evmcommon.BytesToAddress(rawAddr[:]) // Callee contract

	funCall, err := abi.DecodeTo(input, 1, []byte{}, 2, math.MaxUint32)
	if err != nil {
		return []byte(err.Error()), false
	}

	jobID := this.jobManager.Add(calleeAddr, funCall)

	if buffer, err := abi.Encode(uint64(jobID)); err != nil {
		return []byte(err.Error()), false
	} else {
		return buffer, true
	}

	// msg := types.NewMessage( // Build the message
	// 	eucommon.User1,
	// 	&calleeAddr,
	// 	0,
	// 	new(big.Int).SetUint64(0), // Amount to transfer
	// 	1e15,
	// 	new(big.Int).SetUint64(1),
	// 	funCall, //need to remove the wrapper first
	// 	nil,
	// 	false, // Stop checking nonce
	// )

	// ccurl := concurrenturl.NewConcurrentUrl(ccurlstorage.NewTransientDB(*(this.api.Ccurl().Store())))
	// _, transitions := this.api.Ccurl().Export(false)
	// this.api.Ccurl().Import(transitions)
	// this.api.Ccurl().PostImport()
	// if errs := this.api.Ccurl().Commit([]uint32{1}); errs != nil && len(errs) != 0 {
	// 	return []byte("Error: Failed to import transitions"), false
	// }

	// // this.api.Ccurl().
	// // ccurl := this.api.Ccurl()
	// statedb := eth.NewImplStateDB(ccurl) // Eth state DB
	// statedb.Prepare([32]byte{}, [32]byte{}, 0)

	// eu := cceu.NewEU(
	// 	params.MainnetChainConfig,
	// 	vm.Config{},
	// 	statedb,
	// 	this.api.New(common.Hash{}, 0, ccurl), // Call function
	// )

	// config := cceu.NewConfig()
	// _, _, receipt, exeResult, err := eu.Run(evmcommon.BytesToHash([]byte{1, 1, 1}), 1, &msg, cceu.NewEVMBlockContext(config), cceu.NewEVMTxContext(msg))
	// if err != nil {
	// 	return []byte(err.Error()), false
	// }

	// if exeResult.Err != nil || receipt.Status != 1 {
	// 	return []byte(exeResult.Err.Error()), false
	// }

	// v, err := abi.Encode(uint32(99))
	// return v, err == nil
}

func (this *MultiprocessHandler) delJob(caller, callee evmcommon.Address, input []byte) ([]byte, bool) {
	id := this.api.GenUUID()
	return id, true
}

// func (this *MultiprocessHandler) run(caller, callee evmcommon.Address, input []byte) ([]byte, bool) {
// 	// id := this.api.GenUUID()
// 	// delta, err := abi.Decode(input, 1, &uint256.Int{}, 1, 32)
// 	return []byte{}, true
// }

func (this *MultiprocessHandler) clear(caller, callee evmcommon.Address, input []byte) ([]byte, bool) {
	// id := this.api.GenUUID()
	// delta, err := abi.Decode(input, 1, &uint256.Int{}, 1, 32)
	return []byte{}, true
}
