package tests

import (
	"bufio"
	"errors"
	"fmt"
	"math"
	"math/big"
	"os"

	"github.com/arcology-network/concurrenturl/v2"
	urlcommon "github.com/arcology-network/concurrenturl/v2/common"
	urltype "github.com/arcology-network/concurrenturl/v2/type"
	"github.com/arcology-network/concurrenturl/v2/type/commutative"
	"github.com/arcology-network/concurrenturl/v2/type/noncommutative"
	"github.com/arcology-network/evm/common"
	evmcommon "github.com/arcology-network/evm/common"
	"github.com/arcology-network/evm/consensus"
	"github.com/arcology-network/evm/core/types"
	evmtypes "github.com/arcology-network/evm/core/types"
	"github.com/arcology-network/evm/core/vm"
	"github.com/arcology-network/evm/crypto"
	"github.com/arcology-network/evm/params"
	cceu "github.com/arcology-network/vm-adaptor"
	eu "github.com/arcology-network/vm-adaptor"

	ccapi "github.com/arcology-network/vm-adaptor/api"
	cceueth "github.com/arcology-network/vm-adaptor/eth"
	"github.com/holiman/uint256"
)

var (
	Coinbase = common.BytesToAddress([]byte("coinbase"))
	Owner    = common.BytesToAddress([]byte("owner"))
	User1    = common.BytesToAddress([]byte("user1"))
	User2    = common.BytesToAddress([]byte("user2"))
)

// fakeChain implements the ChainContext interface.
type fakeChain struct {
}

func (chain *fakeChain) GetHeader(common.Hash, uint64) *types.Header {
	return &types.Header{}
}

func (chain *fakeChain) Engine() consensus.Engine {
	return nil
}

func MainConfig() *eu.Config {
	vmConfig := vm.Config{}
	cfg := &eu.Config{
		ChainConfig: params.MainnetChainConfig,
		VMConfig:    &vmConfig,
		BlockNumber: new(big.Int).SetUint64(10000000),
		ParentHash:  common.Hash{},
		Time:        new(big.Int).SetUint64(10000000),
		Coinbase:    &Coinbase,
		GasLimit:    math.MaxUint64,
		Difficulty:  new(big.Int).SetUint64(10000000),
	}
	cfg.Chain = new(fakeChain)
	return cfg
}

func FormatValue(value interface{}) string {
	switch value.(type) {
	case *commutative.Meta:
		meta := value.(*commutative.Meta)
		var str string
		str += "{"
		for i, k := range meta.Keys() {
			str += k.(string)
			if i != len(meta.Keys())-1 {
				str += ", "
			}
		}
		str += "}"
		if len(meta.Added()) != 0 {
			str += " + {"
			for i, k := range meta.Added() {
				str += k
				if i != len(meta.Added())-1 {
					str += ", "
				}
			}
			str += "}"
		}
		if len(meta.Removed()) != 0 {
			str += " - {"
			for i, k := range meta.Removed() {
				str += k
				if i != len(meta.Removed())-1 {
					str += ", "
				}
			}
			str += "}"
		}
		return str
	case *noncommutative.Int64:
		// uint256.NewInt(0)
		return fmt.Sprintf(" = %v", (*(value.(*noncommutative.Int64))))
	case *noncommutative.Bytes:
		return fmt.Sprintf(" = %v", value.(*noncommutative.Bytes).Data())
	case *commutative.U256:
		v := value.(*commutative.U256).Value()
		d := value.(*commutative.U256).GetDelta()
		return fmt.Sprintf(" = %v + %v", (*(v.(*uint256.Int))), d.(*big.Int).Int64())
	case *commutative.Int64:
		v := value.(*commutative.Int64).Value()
		d := value.(*commutative.Int64).GetDelta()
		return fmt.Sprintf(" = %v + %v", v, d)
	}
	return ""
}

func FormatTransitions(transitions []urlcommon.UnivalueInterface) string {
	var str string
	for _, t := range transitions {
		str += fmt.Sprintf("%v%v%v%v%v%v%v%v%v%v%v%v%v%v",
			"Tx=", t.(*urltype.Univalue).GetTx(),
			" Reads=", t.(*urltype.Univalue).Reads(),
			" Writes=", t.(*urltype.Univalue).Writes(),
			" Preexists=", t.(*urltype.Univalue).Preexist(),
			" Composite=", t.(*urltype.Univalue).Composite(),
			" Path=", *(t.(*urltype.Univalue).GetPath()),
			" Value", FormatValue(t.(*urltype.Univalue).Value())+"\n")
	}
	return str
}

func Prepare(db urlcommon.DatastoreInterface, height uint64, transitions []urlcommon.UnivalueInterface, txs []uint32) (*cceu.EU, *cceu.Config) {
	url := concurrenturl.NewConcurrentUrl(db)
	if transitions != nil && len(transitions) != 0 {
		url.Import(transitions)
		url.PostImport()
		url.Commit(txs)
	}

	api := ccapi.NewAPI(url)
	statedb := cceueth.NewImplStateDB(url)

	config := MainConfig()
	config.Coinbase = &Coinbase
	config.BlockNumber = new(big.Int).SetUint64(height)
	config.Time = new(big.Int).SetUint64(height)

	return cceu.NewEU(config.ChainConfig, *config.VMConfig, config.Chain, statedb, api, url), config
}

func Deploy(eu *cceu.EU, config *cceu.Config, owner evmcommon.Address, nonce uint64, code string, args ...[]byte) ([]urlcommon.UnivalueInterface, *evmtypes.Receipt, error) {
	data := evmcommon.Hex2Bytes(code)
	for _, arg := range args {
		data = append(data, evmcommon.BytesToHash(arg).Bytes()...)
	}
	msg := evmtypes.NewMessage(owner, nil, nonce, new(big.Int).SetUint64(0), 1e15, new(big.Int).SetUint64(1), data, nil, false)
	_, transitions, receipt, err := eu.Run(evmcommon.BytesToHash([]byte{byte(nonce + 1), byte(nonce + 1), byte(nonce + 1)}), int(nonce+1), &msg, cceu.NewEVMBlockContextV2(config), cceu.NewEVMTxContext(msg))
	return transitions, receipt, err
}

// func Run(eu *cceu.EU, config *cceu.Config, from, to *evmcommon.Address, nonce uint64, checkNonce bool, function string, args ...[]byte) ([]urlcommon.UnivalueInterface, *evmtypes.Receipt, error) {
// 	data := crypto.Keccak256([]byte(function))[:4]
// 	for _, arg := range args {
// 		data = append(data, evmcommon.BytesToHash(arg).Bytes()...)
// 	}
// 	msg := evmtypes.NewMessage(*from, to, nonce, new(big.Int).SetUint64(0), 1e15, new(big.Int).SetUint64(1), data, nil, checkNonce)
// 	_, transitions, receipt, err := eu.Run(evmcommon.BytesToHash([]byte{byte((nonce + 1) / 65536), byte((nonce + 1) / 256), byte((nonce + 1) % 256)}), int(nonce+1), &msg, cceu.NewEVMBlockContextV2(config), cceu.NewEVMTxContext(msg))
// 	return transitions, receipt, err
// }

func CallFunc(eu *cceu.EU, config *cceu.Config, from, to *evmcommon.Address, nonce uint64, checkNonce bool, function string, encodedArgs ...[]byte) ([]urlcommon.UnivalueInterface, []urlcommon.UnivalueInterface, *evmtypes.Receipt, error) {
	data := crypto.Keccak256([]byte(function))[:4]
	for _, arg := range encodedArgs {
		data = append(data, arg...)
	}
	msg := evmtypes.NewMessage(*from, to, nonce, new(big.Int).SetUint64(0), 1e15, new(big.Int).SetUint64(1), data, nil, checkNonce)
	accesses, transitions, receipt, err := eu.Run(evmcommon.BytesToHash([]byte{byte((nonce + 1) / 65536), byte((nonce + 1) / 256), byte((nonce + 1) % 256)}), int(nonce+1), &msg, cceu.NewEVMBlockContextV2(config), cceu.NewEVMTxContext(msg))
	return accesses, transitions, receipt, err
}

func BytecodeReader(fileName string) (string, error) {
	// Open the file
	file, err := os.Open(fileName)
	if err != nil {
		return "", errors.New("Failed to read " + fileName)
	}
	defer file.Close()

	// Create a scanner to read the file with the specified encoding
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	scanner.Scan() // read the first line of the file

	return scanner.Text(), nil // Convert the scanned text to a string
}
