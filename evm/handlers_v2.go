package evm

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"

	"github.com/arcology-network/common-lib/types"
	clib "github.com/arcology-network/concurrentlib"
	"github.com/arcology-network/evm/common"
)

func arrayCreateV2(api *APIV2, caller, callee common.Address, input []byte, origin common.Address, nonce uint64, thash, bhash common.Hash) ([]byte, bool) {
	elemType := BytesToInt32(input[64:68])
	size := BytesToInt32(input[96:100])
	idLen := BytesToInt32(input[128:132])
	id := string(input[132 : 132+idLen])

	ok := api.array.Create(types.Address(hex.EncodeToString(caller.Bytes())), id, int(elemType), int(size))

	api.AddLog("arrayCreate", fmt.Sprintf("params: caller=%x elemType=%v size=%v idLen=%v id=%x array.Create=%v", caller.Bytes(), elemType, size, idLen, []byte(id), ok))
	return nil, ok
}

func arraySizeV2(api *APIV2, caller, callee common.Address, input []byte, origin common.Address, nonce uint64, thash, bhash common.Hash) ([]byte, bool) {
	idLen := BytesToInt32(input[64:68])
	id := string(input[68 : 68+idLen])
	value := api.array.GetSize(types.Address(hex.EncodeToString(caller.Bytes())), id)
	if value == -1 {
		api.AddLog("arraySize", fmt.Sprintf("params: caller=%x idLen=%v id=%x array.GetSize=%v", caller.Bytes(), idLen, []byte(id), value))
		return nil, false
	}
	data := padLeftToSize(Int64ToBytes(int64(value)), 32)
	api.AddLog("arraySize", fmt.Sprintf("params: caller=%x idLen=%v id=%x array.GetSize=%v padLeftToSize=%x", caller.Bytes(), idLen, []byte(id), value, data))
	return data, true
}

func arraySetUint256V2(api *APIV2, caller, callee common.Address, input []byte, origin common.Address, nonce uint64, thash, bhash common.Hash) ([]byte, bool) {
	index := BytesToInt32(input[64:68])
	value := input[68:100]
	idLen := BytesToInt32(input[128:132])
	id := string(input[132 : 132+idLen])
	ok := api.array.SetElem(types.Address(hex.EncodeToString(caller.Bytes())), id, int(index), value, clib.DataTypeUint256)
	api.AddLog("arraySetUint256", fmt.Sprintf("params: caller=%x idLen=%v id=%x value=%x index=%v array.SetElem=%v", caller.Bytes(), idLen, []byte(id), value, index, ok))
	return nil, ok
}

func arrayGetUint256V2(api *APIV2, caller, callee common.Address, input []byte, origin common.Address, nonce uint64, thash, bhash common.Hash) ([]byte, bool) {
	index := BytesToInt32(input[64:68])
	idLen := BytesToInt32(input[96:100])
	id := string(input[100 : 100+idLen])
	value, _ := api.array.GetElem(types.Address(hex.EncodeToString(caller.Bytes())), id, int(index), clib.DataTypeUint256)

	data := padLeftToSize(value, 32)
	api.AddLog("arrayGetUint256", fmt.Sprintf("params: caller=%x idLen=%v id=%x array.GetElem=%x index=%v padLeftToSize=%x", caller.Bytes(), idLen, []byte(id), value, index, data))
	return data, true
}

func arraySetAddressV2(api *APIV2, caller, callee common.Address, input []byte, origin common.Address, nonce uint64, thash, bhash common.Hash) ([]byte, bool) {
	index := BytesToInt32(input[64:68])
	value := input[80:100]
	idLen := BytesToInt32(input[128:132])
	id := string(input[132 : 132+idLen])
	ok := api.array.SetElem(types.Address(hex.EncodeToString(caller.Bytes())), id, int(index), value, clib.DataTypeAddress)
	api.AddLog("arraySetAddress", fmt.Sprintf("params: caller=%x idLen=%v id=%x value=%x index=%v array.SetElem=%v", caller.Bytes(), idLen, []byte(id), value, index, ok))
	return nil, ok
}

func arrayGetAddressV2(api *APIV2, caller, callee common.Address, input []byte, origin common.Address, nonce uint64, thash, bhash common.Hash) ([]byte, bool) {
	index := BytesToInt32(input[64:68])
	idLen := BytesToInt32(input[96:100])
	id := string(input[100 : 100+idLen])
	value, _ := api.array.GetElem(types.Address(hex.EncodeToString(caller.Bytes())), id, int(index), clib.DataTypeAddress)
	data := padLeftToSize(value, 32)
	api.AddLog("arrayGetAddress", fmt.Sprintf("params: caller=%x idLen=%v id=%x array.GetElem=%x index=%v padLeftToSize=%x", caller.Bytes(), idLen, []byte(id), value, index, data))
	return data, true
}

func arraySetBytesV2(api *APIV2, caller, callee common.Address, input []byte, origin common.Address, nonce uint64, thash, bhash common.Hash) ([]byte, bool) {
	index := BytesToInt32(input[64:68])
	idLen := BytesToInt32(input[128:132])
	id := string(input[132 : 132+idLen])
	value := input[164:]
	if len(value)%32 != 0 {
		panic("the data was not 32 bytes aligned")
	}
	ok := api.array.SetElem(types.Address(hex.EncodeToString(caller.Bytes())), id, int(index), value, clib.DataTypeBytes)
	api.AddLog("arraySetBytes", fmt.Sprintf("params: caller=%x idLen=%v id=%x value=%x index=%v array.SetElem=%v", caller.Bytes(), idLen, []byte(id), value, index, ok))
	return nil, ok
}

func arrayGetBytesV2(api *APIV2, caller, callee common.Address, input []byte, origin common.Address, nonce uint64, thash, bhash common.Hash) ([]byte, bool) {
	index := BytesToInt32(input[64:68])
	idLen := BytesToInt32(input[96:100])
	id := string(input[100 : 100+idLen])
	value, _ := api.array.GetElem(types.Address(hex.EncodeToString(caller.Bytes())), id, int(index), clib.DataTypeBytes)
	// We cannot supply default value for bytes.
	if value == nil {
		return nil, false
	}
	// This is tricky!!!
	value = append(padLeftToSize([]byte{32}, 32), value...)

	api.AddLog("arrayGetBytes", fmt.Sprintf("params: caller=%x idLen=%v id=%x array.GetElem=%x index=%v ", caller.Bytes(), idLen, []byte(id), value, index))

	return value, true
}

func varCreateV2(api *APIV2, caller, callee common.Address, input []byte, origin common.Address, nonce uint64, thash, bhash common.Hash) ([]byte, bool) {
	valueType := BytesToInt32(input[64:68])
	idLen := BytesToInt32(input[96:100])
	id := string(input[100 : 100+idLen])
	ok := api.array.Create(types.Address(hex.EncodeToString(caller.Bytes())), id, int(valueType), 1)
	api.AddLog("varCreate", fmt.Sprintf("params: caller=%x idLen=%v id=%x valueType=%v array.Create=%v ", caller.Bytes(), idLen, []byte(id), valueType, ok))
	return nil, ok
}

func varSetUint256V2(api *APIV2, caller, callee common.Address, input []byte, origin common.Address, nonce uint64, thash, bhash common.Hash) ([]byte, bool) {
	value := input[36:68]
	idLen := BytesToInt32(input[96:100])
	id := string(input[100 : 100+idLen])
	ok := api.array.SetElem(types.Address(hex.EncodeToString(caller.Bytes())), id, 0, value, clib.DataTypeUint256)
	api.AddLog("varSetUint256", fmt.Sprintf("params: caller=%x idLen=%v id=%x value=%x array.SetElem=%v ", caller.Bytes(), idLen, []byte(id), value, ok))
	return nil, ok
}

func varGetUint256V2(api *APIV2, caller, callee common.Address, input []byte, origin common.Address, nonce uint64, thash, bhash common.Hash) ([]byte, bool) {
	idLen := BytesToInt32(input[64:68])
	id := string(input[68 : 68+idLen])
	value, _ := api.array.GetElem(types.Address(hex.EncodeToString(caller.Bytes())), id, 0, clib.DataTypeUint256)
	api.AddLog("varGetUint256", fmt.Sprintf("params: caller=%x idLen=%v id=%x value=%x ", caller.Bytes(), idLen, []byte(id), value))
	if value == nil {
		return nil, false
	}
	return value, true
}

func varSetAddressV2(api *APIV2, caller, callee common.Address, input []byte, origin common.Address, nonce uint64, thash, bhash common.Hash) ([]byte, bool) {
	value := input[48:68]
	idLen := BytesToInt32(input[96:100])
	id := string(input[100 : 100+idLen])
	ok := api.array.SetElem(types.Address(hex.EncodeToString(caller.Bytes())), id, 0, value, clib.DataTypeAddress)
	api.AddLog("varSetAddress", fmt.Sprintf("params: caller=%x idLen=%v id=%x value=%x array.SetElem=%v", caller.Bytes(), idLen, []byte(id), value, ok))
	return nil, ok
}

func varGetAddressV2(api *APIV2, caller, callee common.Address, input []byte, origin common.Address, nonce uint64, thash, bhash common.Hash) ([]byte, bool) {
	idLen := BytesToInt32(input[64:68])
	id := string(input[68 : 68+idLen])
	value, _ := api.array.GetElem(types.Address(hex.EncodeToString(caller.Bytes())), id, 0, clib.DataTypeAddress)
	if value == nil {
		api.AddLog("varGetAddress", fmt.Sprintf("params: caller=%x idLen=%v id=%x array.GetElem=%x", caller.Bytes(), idLen, []byte(id), value))
		return nil, false
	}
	data := padLeftToSize(value, 32)
	api.AddLog("varGetAddress", fmt.Sprintf("params: caller=%x idLen=%v id=%x array.GetElem=%x padLeftToSize=%x", caller.Bytes(), idLen, []byte(id), value, data))
	return data, true
}

func varSetBytesV2(api *APIV2, caller, callee common.Address, input []byte, origin common.Address, nonce uint64, thash, bhash common.Hash) ([]byte, bool) {
	idLen := BytesToInt32(input[96:100])
	id := string(input[100 : 100+idLen])
	value := input[132:]
	if len(value)%32 != 0 {
		panic("the data was not 32 bytes aligned")
	}
	ok := api.array.SetElem(types.Address(hex.EncodeToString(caller.Bytes())), id, 0, value, clib.DataTypeBytes)
	api.AddLog("varSetBytes", fmt.Sprintf("params: caller=%x idLen=%v id=%x array.SetElem=%v value=%x", caller.Bytes(), idLen, []byte(id), ok, value))
	return nil, ok
}

func varGetBytesV2(api *APIV2, caller, callee common.Address, input []byte, origin common.Address, nonce uint64, thash, bhash common.Hash) ([]byte, bool) {
	idLen := BytesToInt32(input[64:68])
	id := string(input[68 : 68+idLen])
	value, _ := api.array.GetElem(types.Address(hex.EncodeToString(caller.Bytes())), id, 0, clib.DataTypeBytes)
	if value == nil {
		api.AddLog("varGetBytes", fmt.Sprintf("params: caller=%x idLen=%v id=%x array.GetElem=%v ", caller.Bytes(), idLen, []byte(id), value))
		return nil, false
	}
	data := append(padLeftToSize([]byte{32}, 32), value...)
	api.AddLog("varGetBytes", fmt.Sprintf("params: caller=%x idLen=%v id=%x array.GetElem=%v padLeftToSize=%x", caller.Bytes(), idLen, []byte(id), value, data))
	return data, true
}

func hashmapCreateV2(api *APIV2, caller, callee common.Address, input []byte, origin common.Address, nonce uint64, thash, bhash common.Hash) ([]byte, bool) {
	keyType := BytesToInt32(input[64:68])
	valueType := BytesToInt32(input[96:100])
	idLen := BytesToInt32(input[128:132])
	id := string(input[132 : 132+idLen])
	ok := api.sortedMap.Create(types.Address(hex.EncodeToString(caller.Bytes())), id, int(keyType), int(valueType))
	api.AddLog("hashmapCreate", fmt.Sprintf("params: caller=%x idLen=%v id=%x sortedMap.Create=%v keyType=%v valueType=%v", caller.Bytes(), idLen, []byte(id), ok, keyType, valueType))
	return nil, ok
}

func hashmapUint256AddressGetV2(api *APIV2, caller, callee common.Address, input []byte, origin common.Address, nonce uint64, thash, bhash common.Hash) ([]byte, bool) {
	key := input[36:68]
	idLen := BytesToInt32(input[96:100])
	id := string(input[100 : 100+idLen])
	value, _ := api.sortedMap.GetValue(types.Address(hex.EncodeToString(caller.Bytes())), id, key, clib.DataTypeUint256, clib.DataTypeAddress)
	data := padLeftToSize(value, 32)
	api.AddLog("hashmapUint256AddressGet", fmt.Sprintf("params: caller=%x idLen=%v id=%x key=%x sortedMap.GetValue=%x padLeftToSize=%x ", caller.Bytes(), idLen, []byte(id), key, value, data))
	return data, true
}

func hashmapUint256AddressSetV2(api *APIV2, caller, callee common.Address, input []byte, origin common.Address, nonce uint64, thash, bhash common.Hash) ([]byte, bool) {
	key := input[36:68]
	value := input[80:100]
	idLen := BytesToInt32(input[128:132])
	id := string(input[132 : 132+idLen])
	ok := api.sortedMap.SetValue(types.Address(hex.EncodeToString(caller.Bytes())), id, key, value, clib.DataTypeUint256, clib.DataTypeAddress)
	api.AddLog("hashmapUint256AddressSet", fmt.Sprintf("params: caller=%x idLen=%v id=%x key=%x value=%x sortedMap.SetValue=%v ", caller.Bytes(), idLen, []byte(id), key, value, ok))
	return nil, ok
}

func hashmapUint256Uint256GetV2(api *APIV2, caller, callee common.Address, input []byte, origin common.Address, nonce uint64, thash, bhash common.Hash) ([]byte, bool) {
	key := input[36:68]
	idLen := BytesToInt32(input[96:100])
	id := string(input[100 : 100+idLen])
	value, _ := api.sortedMap.GetValue(types.Address(hex.EncodeToString(caller.Bytes())), id, key, clib.DataTypeUint256, clib.DataTypeUint256)
	data := padLeftToSize(value, 32)
	api.AddLog("hashmapUint256Uint256Get", fmt.Sprintf("params: caller=%x idLen=%v id=%x key=%x sortedMap.GetValue=%x padLeftToSize=%x", caller.Bytes(), idLen, []byte(id), key, value, data))
	return data, true
}

func hashmapUint256Uint256SetV2(api *APIV2, caller, callee common.Address, input []byte, origin common.Address, nonce uint64, thash, bhash common.Hash) ([]byte, bool) {
	key := input[36:68]
	value := input[68:100]
	idLen := BytesToInt32(input[128:132])
	id := string(input[132 : 132+idLen])
	ok := api.sortedMap.SetValue(types.Address(hex.EncodeToString(caller.Bytes())), id, key, value, clib.DataTypeUint256, clib.DataTypeUint256)
	api.AddLog("hashmapUint256Uint256Set", fmt.Sprintf("params: caller=%x idLen=%v id=%x key=%x value=%x sortedMap.SetValue=%v ", caller.Bytes(), idLen, []byte(id), key, value, ok))
	return nil, ok
}

func hashmapUint256BytesGetV2(api *APIV2, caller, callee common.Address, input []byte, origin common.Address, nonce uint64, thash, bhash common.Hash) ([]byte, bool) {
	key := input[36:68]
	idLen := BytesToInt32(input[96:100])
	id := string(input[100 : 100+idLen])
	value, _ := api.sortedMap.GetValue(types.Address(hex.EncodeToString(caller.Bytes())), id, key, clib.DataTypeUint256, clib.DataTypeBytes)
	// We cannot supply default value for bytes.
	if value == nil {
		api.AddLog("hashmapUint256BytesGet", fmt.Sprintf("params: caller=%x idLen=%v id=%x key=%x sortedMap.GetValue=%x ", caller.Bytes(), idLen, []byte(id), key, value))

		return nil, false
	}
	// This is tricky!!!
	data := append(padLeftToSize([]byte{32}, 32), value...)
	api.AddLog("hashmapUint256BytesGet", fmt.Sprintf("params: caller=%x idLen=%v id=%x key=%x sortedMap.GetValue=%x padLeftToSize=%x", caller.Bytes(), idLen, []byte(id), key, value, data))

	return data, true
}

func hashmapUint256BytesSetV2(api *APIV2, caller, callee common.Address, input []byte, origin common.Address, nonce uint64, thash, bhash common.Hash) ([]byte, bool) {
	key := input[36:68]
	idLen := BytesToInt32(input[128:132])
	id := string(input[132 : 132+idLen])
	value := input[164:]
	if len(value)%32 != 0 {
		panic("the data was not 32 bytes aligned")
	}
	ok := api.sortedMap.SetValue(types.Address(hex.EncodeToString(caller.Bytes())), id, key, value, clib.DataTypeUint256, clib.DataTypeBytes)
	api.AddLog("hashmapUint256BytesSet", fmt.Sprintf("params: caller=%x idLen=%v id=%x key=%x sortedMap.SetValue=%v value=%x", caller.Bytes(), idLen, []byte(id), key, ok, value))

	return nil, ok
}

func hashmapAddressAddressGetV2(api *APIV2, caller, callee common.Address, input []byte, origin common.Address, nonce uint64, thash, bhash common.Hash) ([]byte, bool) {
	key := input[48:68]
	idLen := BytesToInt32(input[96:100])
	id := string(input[100 : 100+idLen])
	value, _ := api.sortedMap.GetValue(types.Address(hex.EncodeToString(caller.Bytes())), id, key, clib.DataTypeAddress, clib.DataTypeAddress)
	data := padLeftToSize(value, 32)
	api.AddLog("hashmapAddressAddressGet", fmt.Sprintf("params: caller=%x idLen=%v id=%x key=%x sortedMap.GetValue=%x padLeftToSize=%x", caller.Bytes(), idLen, []byte(id), key, value, data))

	return data, true
}

func hashmapAddressAddressSetV2(api *APIV2, caller, callee common.Address, input []byte, origin common.Address, nonce uint64, thash, bhash common.Hash) ([]byte, bool) {
	key := input[48:68]
	value := input[80:100]
	idLen := BytesToInt32(input[128:132])
	id := string(input[132 : 132+idLen])
	ok := api.sortedMap.SetValue(types.Address(hex.EncodeToString(caller.Bytes())), id, key, value, clib.DataTypeAddress, clib.DataTypeAddress)
	api.AddLog("hashmapAddressAddressSet", fmt.Sprintf("params: caller=%x idLen=%v id=%x key=%x value=%x sortedMap.SetValue=%v", caller.Bytes(), idLen, []byte(id), key, value, ok))

	return nil, ok
}

func hashmapAddressUint256GetV2(api *APIV2, caller, callee common.Address, input []byte, origin common.Address, nonce uint64, thash, bhash common.Hash) ([]byte, bool) {
	key := input[48:68]
	idLen := BytesToInt32(input[96:100])
	id := string(input[100 : 100+idLen])
	value, _ := api.sortedMap.GetValue(types.Address(hex.EncodeToString(caller.Bytes())), id, key, clib.DataTypeAddress, clib.DataTypeUint256)

	data := padLeftToSize(value, 32)
	api.AddLog("hashmapAddressUint256Get", fmt.Sprintf("params: caller=%x idLen=%v id=%x key=%x  sortedMap.GetValue=%x padLeftToSize=%x", caller.Bytes(), idLen, []byte(id), key, value, data))

	return data, true
}

func hashmapAddressUint256SetV2(api *APIV2, caller, callee common.Address, input []byte, origin common.Address, nonce uint64, thash, bhash common.Hash) ([]byte, bool) {
	key := input[48:68]
	value := input[68:100]
	idLen := BytesToInt32(input[128:132])
	id := string(input[132 : 132+idLen])
	ok := api.sortedMap.SetValue(types.Address(hex.EncodeToString(caller.Bytes())), id, key, value, clib.DataTypeAddress, clib.DataTypeUint256)
	api.AddLog("hashmapAddressUint256Set", fmt.Sprintf("params: caller=%x idLen=%v id=%x key=%x value=%x sortedMap.SetValue=%v ", caller.Bytes(), idLen, []byte(id), key, value, ok))

	return nil, ok
}

func hashmapAddressBytesGetV2(api *APIV2, caller, callee common.Address, input []byte, origin common.Address, nonce uint64, thash, bhash common.Hash) ([]byte, bool) {
	key := input[48:68]
	idLen := BytesToInt32(input[96:100])
	id := string(input[100 : 100+idLen])
	value, _ := api.sortedMap.GetValue(types.Address(hex.EncodeToString(caller.Bytes())), id, key, clib.DataTypeAddress, clib.DataTypeBytes)
	if value == nil {
		api.AddLog("hashmapAddressBytesGet", fmt.Sprintf("params: caller=%x idLen=%v id=%x key=%x sortedMap.GetValue=%x ", caller.Bytes(), idLen, []byte(id), key, value))

		return nil, false
	}
	data := append(padLeftToSize([]byte{32}, 32), value...)
	api.AddLog("hashmapAddressBytesGet", fmt.Sprintf("params: caller=%x idLen=%v id=%x key=%x sortedMap.GetValue=%x padLeftToSize=%x", caller.Bytes(), idLen, []byte(id), key, value, data))

	return data, true
}

func hashmapAddressBytesSetV2(api *APIV2, caller, callee common.Address, input []byte, origin common.Address, nonce uint64, thash, bhash common.Hash) ([]byte, bool) {
	key := input[48:68]
	idLen := BytesToInt32(input[128:132])
	id := string(input[132 : 132+idLen])
	value := input[164:]
	if len(value)%32 != 0 {
		panic("the data was not 32 bytes aligned")
	}
	ok := api.sortedMap.SetValue(types.Address(hex.EncodeToString(caller.Bytes())), id, key, value, clib.DataTypeAddress, clib.DataTypeBytes)
	api.AddLog("hashmapAddressBytesSet", fmt.Sprintf("params: caller=%x idLen=%v id=%x key=%x value=%x sortedMap.SetValue=%v", caller.Bytes(), idLen, []byte(id), key, value, ok))

	return nil, ok
}

func hashmapDeleteKeyUint256V2(api *APIV2, caller, callee common.Address, input []byte, origin common.Address, nonce uint64, thash, bhash common.Hash) ([]byte, bool) {
	key := input[36:68]
	idLen := BytesToInt32(input[96:100])
	id := string(input[100 : 100+idLen])
	ok := api.sortedMap.DeleteKey(types.Address(hex.EncodeToString(caller.Bytes())), id, key, clib.DataTypeUint256)
	api.AddLog("hashmapDeleteKeyUint256", fmt.Sprintf("params: caller=%x idLen=%v id=%x key=%x sortedMap.DeleteKey=%v ", caller.Bytes(), idLen, []byte(id), key, ok))

	return nil, ok
}

func hashmapDeleteKeyAddressV2(api *APIV2, caller, callee common.Address, input []byte, origin common.Address, nonce uint64, thash, bhash common.Hash) ([]byte, bool) {
	key := input[48:68]
	idLen := BytesToInt32(input[96:100])
	id := string(input[100 : 100+idLen])
	ok := api.sortedMap.DeleteKey(types.Address(hex.EncodeToString(caller.Bytes())), id, key, clib.DataTypeAddress)
	api.AddLog("hashmapDeleteKeyAddress", fmt.Sprintf("params: caller=%x idLen=%v id=%x key=%x sortedMap.DeleteKey=%v ", caller.Bytes(), idLen, []byte(id), key, ok))

	return nil, ok
}

func uuidGenV2(api *APIV2, caller, callee common.Address, input []byte, origin common.Address, nonce uint64, thash, bhash common.Hash) ([]byte, bool) {
	uuid := clib.UUIDGen(origin.Bytes(), nonce, caller.Bytes(), input, thash.Bytes(), bhash.Bytes())
	api.AddLog("uuidGen", fmt.Sprintf("params: caller=%x origin=%x nonce=%v input=%x bhash=%x uuid=%v ", caller.Bytes(), origin.Bytes(), nonce, input, bhash.Bytes(), uuid))

	return uuid, true
}

func systemCreateDeferV2(api *APIV2, caller, callee common.Address, input []byte, origin common.Address, nonce uint64, thash, bhash common.Hash) ([]byte, bool) {
	idLen := BytesToInt32(input[96:100])
	id := string(input[100 : 100+idLen])
	sigLen := BytesToInt32(input[160:164])
	sig := string(input[164 : 164+sigLen])
	ok := api.deferCall.Create(types.Address(hex.EncodeToString(caller.Bytes())), id, sig)
	api.AddLog("systemCreateDefer", fmt.Sprintf("params: caller=%x idLen=%v id=%x sigLen=%v sig=%x deferCall.Create=%v", caller.Bytes(), idLen, []byte(id), sigLen, sig, ok))
	return nil, ok
}

func systemCallDeferV2(api *APIV2, caller, callee common.Address, input []byte, origin common.Address, nonce uint64, thash, bhash common.Hash) ([]byte, bool) {
	if api.GetDeferCall() != nil || api.IsInDeferCall() {
		api.AddLog("systemCallDefer", fmt.Sprintf("params: caller=%x api.GetDeferCall()=%v api.IsInDeferCall()=%v", caller.Bytes(), api.GetDeferCall(), api.IsInDeferCall()))
		return nil, false
	}

	idLen := BytesToInt32(input[64:68])
	id := string(input[68 : 68+idLen])
	if !api.deferCall.IsExist(types.Address(hex.EncodeToString(caller.Bytes())), id) {
		api.AddLog("systemCallDefer", fmt.Sprintf("params: caller=%x idLen=%v id=%x deferCall.IsExist=false", caller.Bytes(), idLen, []byte(id)))
		return nil, false
	}
	api.SetDeferCall(types.Address(hex.EncodeToString(caller.Bytes())), id)
	api.AddLog("systemCallDefer", fmt.Sprintf("params: caller=%x idLen=%v  id=%x deferCall.IsExist=true SetDeferCall", caller.Bytes(), idLen, []byte(id)))
	return nil, true
}

func queueCreateV2(api *APIV2, caller, callee common.Address, input []byte, origin common.Address, nonce uint64, thash, bhash common.Hash) ([]byte, bool) {
	elemType := BytesToInt32(input[64:68])
	idLen := BytesToInt32(input[96:100])
	id := string(input[100 : 100+idLen])
	ok := api.queue.Create(types.Address(hex.EncodeToString(caller.Bytes())), id, int(elemType))
	api.AddLog("queueCreate", fmt.Sprintf("params: caller=%x idLen=%v  id=%x elemType=%v queue.Create=%v", caller.Bytes(), idLen, []byte(id), elemType, ok))
	return nil, ok
}

func queueSizeV2(api *APIV2, caller, callee common.Address, input []byte, origin common.Address, nonce uint64, thash, bhash common.Hash) ([]byte, bool) {
	idLen := BytesToInt32(input[64:68])
	id := string(input[68 : 68+idLen])
	value := api.queue.GetSize(types.Address(hex.EncodeToString(caller.Bytes())), id)
	if value == -1 {
		api.AddLog("queueSize", fmt.Sprintf("params: caller=%x idLen=%v  id=%x queue.GetSize=%x ", caller.Bytes(), idLen, []byte(id), value))
		return nil, false
	}
	data := padLeftToSize(Int64ToBytes(int64(value)), 32)
	api.AddLog("queueSize", fmt.Sprintf("params: caller=%x idLen=%v  id=%x queue.GetSize=%x  padLeftToSize=%x", caller.Bytes(), idLen, []byte(id), value, data))
	return data, true
}

func queuePushUint256V2(api *APIV2, caller, callee common.Address, input []byte, origin common.Address, nonce uint64, thash, bhash common.Hash) ([]byte, bool) {
	value := input[36:68]
	idLen := BytesToInt32(input[96:100])
	id := string(input[100 : 100+idLen])
	ok := api.queue.Push(types.Address(hex.EncodeToString(caller.Bytes())), id, value, clib.DataTypeUint256)
	api.AddLog("queuePushUint256", fmt.Sprintf("params: caller=%x idLen=%v  id=%x value=%x queue.Push=%v", caller.Bytes(), idLen, []byte(id), value, ok))
	return nil, ok
}

func queuePopUint256V2(api *APIV2, caller, callee common.Address, input []byte, origin common.Address, nonce uint64, thash, bhash common.Hash) ([]byte, bool) {
	idLen := BytesToInt32(input[64:68])
	id := string(input[68 : 68+idLen])
	value, _ := api.queue.Pop(types.Address(hex.EncodeToString(caller.Bytes())), id, clib.DataTypeUint256)
	api.AddLog("queuePopUint256", fmt.Sprintf("params: caller=%x idLen=%v  id=%x queue.Pop=%v", caller.Bytes(), idLen, []byte(id), value))
	if value == nil {
		return nil, false
	}
	return value, true
}

func BytesToInt32(input []byte) int32 {
	var i32 int32
	_ = binary.Read(bytes.NewReader(input), binary.BigEndian, &i32)
	return i32
}

func BytesToInt64(input []byte) int64 {
	var i64 int64
	_ = binary.Read(bytes.NewReader(input), binary.BigEndian, &i64)
	return i64
}

func Int64ToBytes(value int64) []byte {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(value))
	return buf
}

func padLeftToSize(b []byte, size int) []byte {
	if len(b) > size {
		panic("unexpected")
	}

	if len(b) == size {
		return b
	}

	tmp := make([]byte, size)
	copy(tmp[size-len(b):], b)
	return tmp
}
