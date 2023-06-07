package common

import (
	commonlib "github.com/arcology-network/common-lib/common"
	"github.com/arcology-network/common-lib/types"
	interfaces "github.com/arcology-network/vm-adaptor/interfaces"

	"github.com/arcology-network/concurrenturl"
	commutative "github.com/arcology-network/concurrenturl/commutative"
)

// Ccurl connectors for Arcology APIs
type CcurlConnector struct {
	apiRouter interfaces.EthApiRouter
	ccurl     *concurrenturl.ConcurrentUrl
	subDir    string
}

func NewCCurlConnector(subDir string, api interfaces.EthApiRouter, ccurl *concurrenturl.ConcurrentUrl) *CcurlConnector {
	return &CcurlConnector{
		subDir:    subDir,
		apiRouter: api,
		ccurl:     ccurl,
	}
}

// Make Arcology paths under the current account
func (this *CcurlConnector) New(account types.Address, containerId string) bool {
	if !this.newStorageRoot(account, this.apiRouter.TxIndex()) { // Create the root path if has been created yet.
		return false
	}
	return this.newContainerRoot(account, containerId[:], this.apiRouter.TxIndex()) //
}

func (this *CcurlConnector) newStorageRoot(account types.Address, txIndex uint32) bool {
	accountRoot := commonlib.StrCat(this.ccurl.Platform.Eth10Account(), string(account), "/")
	if value, _ := this.ccurl.Peek(accountRoot); value == nil {
		return this.ccurl.CreateAccount(txIndex, this.ccurl.Platform.Eth10(), string(account)) != nil // Create a new account
	}
	return true // ALready exists
}

func (this *CcurlConnector) newContainerRoot(account types.Address, id string, txIndex uint32) bool {
	containerRoot := this.Key(account, id)
	if value, _ := this.ccurl.Peek(containerRoot); value == nil {
		_, err := this.ccurl.Write(txIndex, containerRoot, commutative.NewPath()) // Create a new container
		return err == nil
	}
	return true // Already exists

}

func (this *CcurlConnector) Key(account types.Address, id string) string { // container ID
	return commonlib.StrCat(this.ccurl.Platform.Eth10Account(), string(account), "/storage", this.subDir, id, "/")
}
