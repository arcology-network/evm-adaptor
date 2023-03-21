package evm

import (
	"encoding/hex"
	"math/big"

	commonlib "github.com/arcology-network/common-lib/common"
	"github.com/arcology-network/concurrenturl/v2"
	commutative "github.com/arcology-network/concurrenturl/v2/type/commutative"
	evmcommon "github.com/arcology-network/evm/common"
	uint256 "github.com/holiman/uint256"
)

func addressToHex(addr evmcommon.Address) string {
	var accHex [2 * evmcommon.AddressLength]byte
	hex.Encode(accHex[:], addr[:])
	return string(accHex[:])
}

func getAccountRootPath(url *concurrenturl.ConcurrentUrl, account evmcommon.Address) string {
	var accHex [2 * evmcommon.AddressLength]byte
	hex.Encode(accHex[:], account[:])
	return commonlib.StrCat(url.Platform.Eth10Account(), string(accHex[:]), "/")
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
	return url.IfExists(getAccountRootPath(url, account))
}

func createAccount(url *concurrenturl.ConcurrentUrl, account evmcommon.Address, tid uint32) {
	if err := url.CreateAccount(tid, url.Platform.Eth10(), addressToHex(account)); err != nil {
		panic(err)
	}

	if err := url.Write(tid, getBalancePath(url, account), commutative.NewBalance(uint256.NewInt(0), new(big.Int))); err != nil {
		panic(err)
	}
	if err := url.Write(tid, getNoncePath(url, account), commutative.NewInt64(0, 0)); err != nil {
		panic(err)
	}
	// if err := url.Write(tid, getCodePath(url, account), noncommutative.NewBytes(nil)); err != nil {
	// 	panic(err)
	// }
}
