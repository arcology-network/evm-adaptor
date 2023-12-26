package eth

import (
	"encoding/hex"
	"fmt"

	codec "github.com/arcology-network/common-lib/codec"
	commonlib "github.com/arcology-network/common-lib/common"
	ccurlcommon "github.com/arcology-network/concurrenturl/common"
	commutative "github.com/arcology-network/concurrenturl/commutative"
	cache "github.com/arcology-network/eu/cache"
	intf "github.com/arcology-network/vm-adaptor/interface"
	evmcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
)

type EthPathBuilder struct {
	// ccurl cache.WriteCache
}

func (this *EthPathBuilder) AccountExist(writeCache *cache.WriteCache, account evmcommon.Address, tid uint32) bool {
	var accHex [2 * evmcommon.AddressLength]byte
	hex.Encode(accHex[:], account[:])
	return writeCache.IfExists(this.AccountRootPath(account))
	// return writeCache.IfExists(this.AccountRootPath(account))
}

func (this *EthPathBuilder) AccountRootPath(account evmcommon.Address) string {
	return commonlib.StrCat(ccurlcommon.ETH10_ACCOUNT_PREFIX, codec.Bytes20(account).Hex(), "/")
}

func (this *EthPathBuilder) StorageRootPath(account evmcommon.Address) string {
	return commonlib.StrCat(ccurlcommon.ETH10_ACCOUNT_PREFIX, codec.Bytes20(account).Hex(), "/storage/native/")
}

func (this *EthPathBuilder) BalancePath(account evmcommon.Address) string {
	return commonlib.StrCat(ccurlcommon.ETH10_ACCOUNT_PREFIX, codec.Bytes20(account).Hex(), "/balance")
}

func (this *EthPathBuilder) NoncePath(account evmcommon.Address) string {
	return commonlib.StrCat(ccurlcommon.ETH10_ACCOUNT_PREFIX, codec.Bytes20(account).Hex(), "/nonce")
}

func (this *EthPathBuilder) CodePath(account evmcommon.Address) string {
	return commonlib.StrCat(ccurlcommon.ETH10_ACCOUNT_PREFIX, codec.Bytes20(account).Hex(), "/code")
}

func getAccountRootPath(writeCache *cache.WriteCache, account evmcommon.Address) string {
	var accHex [2 * evmcommon.AddressLength]byte
	hex.Encode(accHex[:], account[:])
	return commonlib.StrCat(ccurlcommon.ETH10_ACCOUNT_PREFIX, codec.Bytes20(account).Hex(), "/")
}

func getStorageRootPath(writeCache *cache.WriteCache, account evmcommon.Address) string {
	var accHex [2 * evmcommon.AddressLength]byte
	hex.Encode(accHex[:], account[:])
	return commonlib.StrCat(ccurlcommon.ETH10_ACCOUNT_PREFIX, string(accHex[:]), "/storage/native/")
}

func getLocalStorageKeyPath(api intf.EthApiRouter, account evmcommon.Address, key evmcommon.Hash) string {
	var accHex [2 * evmcommon.AddressLength]byte
	hex.Encode(accHex[:], account[:])

	mem := api.VM().(*vm.EVM).ArcologyNetworkAPIs.CallContext.Memory.Data()
	fmt.Print(len(mem))
	return ccurlcommon.ETH10_ACCOUNT_PREFIX + string(accHex[:]) + "/storage/native/local/" + "0"
}

func getStorageKeyPath(api intf.EthApiRouter, account evmcommon.Address, key evmcommon.Hash) string {
	var accHex [2 * evmcommon.AddressLength]byte
	hex.Encode(accHex[:], account[:])
	return commonlib.StrCat(ccurlcommon.ETH10_ACCOUNT_PREFIX, string(accHex[:]), "/storage/native/", key.Hex())
}

func getBalancePath(writeCache *cache.WriteCache, account evmcommon.Address) string {
	var accHex [2 * evmcommon.AddressLength]byte
	hex.Encode(accHex[:], account[:])
	return commonlib.StrCat(ccurlcommon.ETH10_ACCOUNT_PREFIX, string(accHex[:]), "/balance")
}

func getNoncePath(writeCache *cache.WriteCache, account evmcommon.Address) string {
	var accHex [2 * evmcommon.AddressLength]byte
	hex.Encode(accHex[:], account[:])
	return commonlib.StrCat(ccurlcommon.ETH10_ACCOUNT_PREFIX, string(accHex[:]), "/nonce")
}

func getCodePath(writeCache *cache.WriteCache, account evmcommon.Address) string {
	var accHex [2 * evmcommon.AddressLength]byte
	hex.Encode(accHex[:], account[:])
	return commonlib.StrCat(ccurlcommon.ETH10_ACCOUNT_PREFIX, string(accHex[:]), "/code")
}

func accountExist(writeCache *cache.WriteCache, account evmcommon.Address, tid uint32) bool {
	var accHex [2 * evmcommon.AddressLength]byte
	hex.Encode(accHex[:], account[:])
	return writeCache.IfExists(getAccountRootPath(writeCache, account))
}

func createAccount(writeCache *cache.WriteCache, account evmcommon.Address, tid uint32) {
	if _, err := writeCache.CreateNewAccount(tid, codec.Bytes20(account).Hex()); err != nil {
		panic(err)
	}

	if _, err := writeCache.Write(tid, getBalancePath(writeCache, account), commutative.NewUnboundedU256()); err != nil { // Initialize balance
		panic(err)
	}

	if _, err := writeCache.Write(tid, getNoncePath(writeCache, account), commutative.NewUnboundedUint64()); err != nil {
		panic(err)
	}
	// if err := writeCache.Write(tid, getCodePath(writeCache, account), noncommutative.NewBytes(nil)); err != nil {
	// 	panic(err)
	// }
}
