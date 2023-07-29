package api

import (
	"github.com/arcology-network/common-lib/codec"
	commonlib "github.com/arcology-network/common-lib/common"
	"github.com/arcology-network/common-lib/types"
	ccurlcommon "github.com/arcology-network/concurrenturl/common"
	eucommon "github.com/arcology-network/vm-adaptor/common"

	"github.com/arcology-network/concurrenturl"
	commutative "github.com/arcology-network/concurrenturl/commutative"
)

// Ccurl connectors for Arcology APIs
type CcurlConnector struct {
	apiRouter eucommon.EthApiRouter
	ccurl     *concurrenturl.ConcurrentUrl
	subDir    string
}

func NewCCurlConnector(subDir string, api eucommon.EthApiRouter, ccurl *concurrenturl.ConcurrentUrl) *CcurlConnector {
	return &CcurlConnector{
		subDir:    subDir,
		apiRouter: api,
		ccurl:     ccurl,
	}
}

// Make Arcology paths under the current account
func (this *CcurlConnector) New(txIndex uint32, deploymentAddr types.Address) bool {
	if !this.newStorageRoot(deploymentAddr, txIndex) { // Create the root path if has been created yet.
		return false
	}
	return this.newContainerRoot(deploymentAddr, txIndex) //
}

func (this *CcurlConnector) newStorageRoot(account types.Address, txIndex uint32) bool {
	accountRoot := commonlib.StrCat(ccurlcommon.ETH10_ACCOUNT_PREFIX, string(account), "/")
	if value, _ := this.ccurl.Peek(accountRoot); value == nil {
		return this.ccurl.NewAccount(txIndex, string(account)) != nil // Create a new account
	}
	return true // ALready exists
}

func (this *CcurlConnector) newContainerRoot(account types.Address, txIndex uint32) bool {
	containerRoot := this.key(account)

	if value, _ := this.ccurl.Peek(containerRoot); value == nil {
		_, err := this.ccurl.Write(txIndex, containerRoot, commutative.NewPath()) // Create a new container
		return err == nil
	}
	return true // Already exists
}

func (this *CcurlConnector) Key(caller [20]byte) string { // container ID
	return this.key(types.Address(codec.Bytes20(caller).Hex()))
}

func (this *CcurlConnector) key(account types.Address) string { // container ID
	return commonlib.StrCat(ccurlcommon.ETH10_ACCOUNT_PREFIX, string(account), "/storage", this.subDir, "/")
}
