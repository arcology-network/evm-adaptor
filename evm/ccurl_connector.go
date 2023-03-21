package evm

import (
	"encoding/hex"

	"github.com/arcology-network/common-lib/types"
	"github.com/arcology-network/evm/common"
)

// func decoder() {
// 	// ABI-encoded data
// 	data := "0x0123456789abcdef"

// 	// ABI definition of the data type
// 	abiDef := "[{\"constant\":false,\"inputs\":[{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"setValue\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// 	// Parse the ABI definition
// 	parsedAbi, err := abi.JSON(strings.NewReader(abiDef))
// 	if err != nil {
// 		panic(err)
// 	}

// 	// Decode the ABI-encoded data
// 	decodedData, err := parsedAbi.Unpack([]byte(data))
// 	if err != nil {
// 		panic(err)
// 	}

// 	// Print the decoded data
// 	fmt.Println(hexutil.Encode(decodedData[0].([]byte)))
// }

func New(api *API, caller, callee common.Address, input []byte, origin common.Address, nonce uint64, thash, bhash common.Hash) ([]byte, bool) {
	elemType := BytesToInt32(input[64:68])
	idLen := BytesToInt32(input[96:100])
	id := string(input[100 : 100+idLen])
	ok := api.dynarray.Create(types.Address(hex.EncodeToString(caller.Bytes())), id, int(elemType))
	return nil, ok
}
