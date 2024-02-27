package eth

import (
	commonlib "github.com/arcology-network/common-lib/common"
	cache "github.com/arcology-network/eu/cache"
	ccurlcommon "github.com/arcology-network/storage-committer/common"
	commutative "github.com/arcology-network/storage-committer/commutative"
	intf "github.com/arcology-network/vm-adaptor/interface"
	evmcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

type EthPathBuilder struct {
	// ccurl cache.WriteCache
}

func (this *EthPathBuilder) AccountExist(writeCache *cache.WriteCache, account evmcommon.Address, tid uint32) bool {
	return writeCache.IfExists(this.AccountRootPath(account))
}

func (this *EthPathBuilder) AccountRootPath(account evmcommon.Address) string {
	return commonlib.StrCat(ccurlcommon.ETH10_ACCOUNT_PREFIX, hexutil.Encode(account[:]), "/")
}

func (this *EthPathBuilder) StorageRootPath(account evmcommon.Address) string {
	return commonlib.StrCat(ccurlcommon.ETH10_ACCOUNT_PREFIX, hexutil.Encode(account[:]), "/storage/native/")
}

func (this *EthPathBuilder) BalancePath(account evmcommon.Address) string {
	return commonlib.StrCat(ccurlcommon.ETH10_ACCOUNT_PREFIX, hexutil.Encode(account[:]), "/balance")
}

func (this *EthPathBuilder) NoncePath(account evmcommon.Address) string {
	return commonlib.StrCat(ccurlcommon.ETH10_ACCOUNT_PREFIX, hexutil.Encode(account[:]), "/nonce")
}

func (this *EthPathBuilder) CodePath(account evmcommon.Address) string {
	return commonlib.StrCat(ccurlcommon.ETH10_ACCOUNT_PREFIX, hexutil.Encode(account[:]), "/code")
}

func getAccountRootPath(writeCache *cache.WriteCache, account evmcommon.Address) string {
	return commonlib.StrCat(ccurlcommon.ETH10_ACCOUNT_PREFIX, hexutil.Encode(account[:]), "/")
}

func getStorageRootPath(writeCache *cache.WriteCache, account evmcommon.Address) string {
	return commonlib.StrCat(ccurlcommon.ETH10_ACCOUNT_PREFIX, hexutil.Encode(account[:]), "/storage/native/")
}

func getLocalStorageKeyPath(api intf.EthApiRouter, account evmcommon.Address, key evmcommon.Hash) string {
	return ccurlcommon.ETH10_ACCOUNT_PREFIX + hexutil.Encode(account[:]) + "/storage/native/local/" + "0"
}

func getStorageKeyPath(api intf.EthApiRouter, account evmcommon.Address, key evmcommon.Hash) string {
	return commonlib.StrCat(ccurlcommon.ETH10_ACCOUNT_PREFIX, hexutil.Encode(account[:]), "/storage/native/", key.Hex())
}

func getBalancePath(writeCache *cache.WriteCache, account evmcommon.Address) string {
	return commonlib.StrCat(ccurlcommon.ETH10_ACCOUNT_PREFIX, hexutil.Encode(account[:]), "/balance")
}

func getNoncePath(writeCache *cache.WriteCache, account evmcommon.Address) string {
	return commonlib.StrCat(ccurlcommon.ETH10_ACCOUNT_PREFIX, hexutil.Encode(account[:]), "/nonce")
}

func getCodePath(writeCache *cache.WriteCache, account evmcommon.Address) string {
	return commonlib.StrCat(ccurlcommon.ETH10_ACCOUNT_PREFIX, hexutil.Encode(account[:]), "/code")
}

func accountExist(writeCache *cache.WriteCache, account evmcommon.Address, tid uint32) bool {
	return writeCache.IfExists(getAccountRootPath(writeCache, account))
}

func createAccount(writeCache *cache.WriteCache, account evmcommon.Address, tid uint32) {
	if _, err := writeCache.CreateNewAccount(tid, hexutil.Encode(account[:])); err != nil {
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
