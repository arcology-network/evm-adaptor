package eth

import (
	"encoding/hex"
	"fmt"

	codec "github.com/arcology-network/common-lib/codec"
	commonlib "github.com/arcology-network/common-lib/common"
	"github.com/arcology-network/concurrenturl"
	ccurlcommon "github.com/arcology-network/concurrenturl/common"
	commutative "github.com/arcology-network/concurrenturl/commutative"
	evmcommon "github.com/arcology-network/evm/common"
	vmCommon "github.com/arcology-network/vm-adaptor/common"
)

type EthCCurlConnector struct {
	ccurl *concurrenturl.ConcurrentUrl
}

func (this *EthCCurlConnector) AccountExist(url *concurrenturl.ConcurrentUrl, account evmcommon.Address, tid uint32) bool {
	var accHex [2 * evmcommon.AddressLength]byte
	hex.Encode(accHex[:], account[:])
	return url.IfExists(this.AccountRootPath(account))
}

func (this *EthCCurlConnector) AccountRootPath(account evmcommon.Address) string {
	return commonlib.StrCat(ccurlcommon.ETH10_ACCOUNT_PREFIX, codec.Bytes20(account).Hex(), "/")
}

func (this *EthCCurlConnector) StorageRootPath(account evmcommon.Address) string {
	return commonlib.StrCat(ccurlcommon.ETH10_ACCOUNT_PREFIX, codec.Bytes20(account).Hex(), "/storage/native/")
}

func (this *EthCCurlConnector) BalancePath(account evmcommon.Address) string {
	return commonlib.StrCat(ccurlcommon.ETH10_ACCOUNT_PREFIX, codec.Bytes20(account).Hex(), "/balance")
}

func (this *EthCCurlConnector) NoncePath(account evmcommon.Address) string {
	return commonlib.StrCat(ccurlcommon.ETH10_ACCOUNT_PREFIX, codec.Bytes20(account).Hex(), "/nonce")
}

func (this *EthCCurlConnector) CodePath(account evmcommon.Address) string {
	return commonlib.StrCat(ccurlcommon.ETH10_ACCOUNT_PREFIX, codec.Bytes20(account).Hex(), "/code")
}

func getAccountRootPath(url *concurrenturl.ConcurrentUrl, account evmcommon.Address) string {
	var accHex [2 * evmcommon.AddressLength]byte
	hex.Encode(accHex[:], account[:])
	return commonlib.StrCat(ccurlcommon.ETH10_ACCOUNT_PREFIX, codec.Bytes20(account).Hex(), "/")
}

func getStorageRootPath(url *concurrenturl.ConcurrentUrl, account evmcommon.Address) string {
	var accHex [2 * evmcommon.AddressLength]byte
	hex.Encode(accHex[:], account[:])
	return commonlib.StrCat(ccurlcommon.ETH10_ACCOUNT_PREFIX, string(accHex[:]), "/storage/native/")
}

func getLocalStorageKeyPath(api vmCommon.EthApiRouter, account evmcommon.Address, key evmcommon.Hash) string {
	var accHex [2 * evmcommon.AddressLength]byte
	hex.Encode(accHex[:], account[:])

	mem := api.VM().ArcologyNetworkAPIs.CallContext.Memory.Data()
	fmt.Print(len(mem))
	return ccurlcommon.ETH10_ACCOUNT_PREFIX + string(accHex[:]) + "/storage/native/local/" + "0"
}

func getStorageKeyPath(api vmCommon.EthApiRouter, account evmcommon.Address, key evmcommon.Hash) string {
	var accHex [2 * evmcommon.AddressLength]byte
	hex.Encode(accHex[:], account[:])
	return commonlib.StrCat(ccurlcommon.ETH10_ACCOUNT_PREFIX, string(accHex[:]), "/storage/native/", key.Hex())
}

func getBalancePath(url *concurrenturl.ConcurrentUrl, account evmcommon.Address) string {
	var accHex [2 * evmcommon.AddressLength]byte
	hex.Encode(accHex[:], account[:])
	return commonlib.StrCat(ccurlcommon.ETH10_ACCOUNT_PREFIX, string(accHex[:]), "/balance")
}

func getNoncePath(url *concurrenturl.ConcurrentUrl, account evmcommon.Address) string {
	var accHex [2 * evmcommon.AddressLength]byte
	hex.Encode(accHex[:], account[:])
	return commonlib.StrCat(ccurlcommon.ETH10_ACCOUNT_PREFIX, string(accHex[:]), "/nonce")
}

func getCodePath(url *concurrenturl.ConcurrentUrl, account evmcommon.Address) string {
	var accHex [2 * evmcommon.AddressLength]byte
	hex.Encode(accHex[:], account[:])
	return commonlib.StrCat(ccurlcommon.ETH10_ACCOUNT_PREFIX, string(accHex[:]), "/code")
}

func accountExist(url *concurrenturl.ConcurrentUrl, account evmcommon.Address, tid uint32) bool {
	var accHex [2 * evmcommon.AddressLength]byte
	hex.Encode(accHex[:], account[:])
	return url.IfExists(getAccountRootPath(url, account))
}

func createAccount(url *concurrenturl.ConcurrentUrl, account evmcommon.Address, tid uint32) {
	if err := url.NewAccount(tid, codec.Bytes20(account).Hex()); err != nil {
		panic(err)
	}

	if _, err := url.Write(tid, getBalancePath(url, account), commutative.NewU256(commutative.U256_MIN, commutative.U256_MAX)); err != nil { // Initialize balance
		panic(err)
	}

	if _, err := url.Write(tid, getNoncePath(url, account), commutative.NewUint64()); err != nil {
		panic(err)
	}
	// if err := url.Write(tid, getCodePath(url, account), noncommutative.NewBytes(nil)); err != nil {
	// 	panic(err)
	// }
}
