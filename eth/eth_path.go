package eth

import (
	"encoding/hex"

	codec "github.com/arcology-network/common-lib/codec"
	commonlib "github.com/arcology-network/common-lib/common"
	"github.com/arcology-network/concurrenturl/v2"
	commutative "github.com/arcology-network/concurrenturl/v2/type/commutative"
	evmcommon "github.com/arcology-network/evm/common"
	uint256 "github.com/holiman/uint256"
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
	return commonlib.StrCat(this.ccurl.Platform.Eth10Account(), codec.Bytes20(account).Hex(), "/")
}

func (this *EthCCurlConnector) StorageRootPath(account evmcommon.Address) string {
	return commonlib.StrCat(this.ccurl.Platform.Eth10Account(), codec.Bytes20(account).Hex(), "/storage/native/")
}

func (this *EthCCurlConnector) BalancePath(account evmcommon.Address) string {
	return commonlib.StrCat(this.ccurl.Platform.Eth10Account(), codec.Bytes20(account).Hex(), "/balance")
}

func (this *EthCCurlConnector) NoncePath(account evmcommon.Address) string {
	return commonlib.StrCat(this.ccurl.Platform.Eth10Account(), codec.Bytes20(account).Hex(), "/nonce")
}

func (this *EthCCurlConnector) CodePath(account evmcommon.Address) string {
	return commonlib.StrCat(this.ccurl.Platform.Eth10Account(), codec.Bytes20(account).Hex(), "/code")
}

func getAccountRootPath(url *concurrenturl.ConcurrentUrl, account evmcommon.Address) string {
	var accHex [2 * evmcommon.AddressLength]byte
	hex.Encode(accHex[:], account[:])
	return commonlib.StrCat(url.Platform.Eth10Account(), codec.Bytes20(account).Hex(), "/")
}

func getStorageRootPath(url *concurrenturl.ConcurrentUrl, account evmcommon.Address) string {
	var accHex [2 * evmcommon.AddressLength]byte
	hex.Encode(accHex[:], account[:])
	return commonlib.StrCat(url.Platform.Eth10Account(), string(accHex[:]), "/storage/native/")
}

func getStorageKeyPath(url *concurrenturl.ConcurrentUrl, account evmcommon.Address, key evmcommon.Hash) string {
	var accHex [2 * evmcommon.AddressLength]byte
	hex.Encode(accHex[:], account[:])
	return commonlib.StrCat(url.Platform.Eth10Account(), string(accHex[:]), "/storage/native/", key.Hex())
}

func getBalancePath(url *concurrenturl.ConcurrentUrl, account evmcommon.Address) string {
	var accHex [2 * evmcommon.AddressLength]byte
	hex.Encode(accHex[:], account[:])
	return commonlib.StrCat(url.Platform.Eth10Account(), string(accHex[:]), "/balance")
}

func getNoncePath(url *concurrenturl.ConcurrentUrl, account evmcommon.Address) string {
	var accHex [2 * evmcommon.AddressLength]byte
	hex.Encode(accHex[:], account[:])
	return commonlib.StrCat(url.Platform.Eth10Account(), string(accHex[:]), "/nonce")
}

func getCodePath(url *concurrenturl.ConcurrentUrl, account evmcommon.Address) string {
	var accHex [2 * evmcommon.AddressLength]byte
	hex.Encode(accHex[:], account[:])
	return commonlib.StrCat(url.Platform.Eth10Account(), string(accHex[:]), "/code")
}

func accountExist(url *concurrenturl.ConcurrentUrl, account evmcommon.Address, tid uint32) bool {
	var accHex [2 * evmcommon.AddressLength]byte
	hex.Encode(accHex[:], account[:])
	return url.IfExists(getAccountRootPath(url, account))
}

func createAccount(url *concurrenturl.ConcurrentUrl, account evmcommon.Address, tid uint32) {
	if err := url.CreateAccount(tid, url.Platform.Eth10(), codec.Bytes20(account).Hex()); err != nil {
		panic(err)
	}

	if err := url.Write(tid, getBalancePath(url, account),
		commutative.NewU256(uint256.NewInt(0), commutative.U256MIN, commutative.U256MAX)); err != nil { // Initialize balance
		panic(err)
	}
	if err := url.Write(tid, getNoncePath(url, account), commutative.NewInt64(0, 0)); err != nil {
		panic(err)
	}
	// if err := url.Write(tid, getCodePath(url, account), noncommutative.NewBytes(nil)); err != nil {
	// 	panic(err)
	// }
}
