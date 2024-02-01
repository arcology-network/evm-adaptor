package common

import (
	common "github.com/arcology-network/common-lib/common"
	"github.com/arcology-network/common-lib/types"
	ccurlcommon "github.com/arcology-network/concurrenturl/common"
	"github.com/ethereum/go-ethereum/common/hexutil"

	commutative "github.com/arcology-network/concurrenturl/commutative"
	cache "github.com/arcology-network/eu/cache"

	intf "github.com/arcology-network/vm-adaptor/interface"
)

// Ccurl connectors for Arcology APIs
type PathBuilder struct {
	apiRouter intf.EthApiRouter
	// apiRouter.WriteCache()     *concurrenturl.ConcurrentUrl
	subDir string
}

func NewPathBuilder(subDir string, api intf.EthApiRouter) *PathBuilder {
	return &PathBuilder{
		subDir:    subDir,
		apiRouter: api,
		// apiRouter.WriteCache():     apiRouter.WriteCache(),
	}
}

// Make Arcology paths under the current account
func (this *PathBuilder) New(txIndex uint32, deploymentAddr types.Address) bool {
	if !this.newStorageRoot(deploymentAddr, txIndex) { // Create the root path if has been created yet.
		return false
	}
	return this.newContainerRoot(deploymentAddr, txIndex) //
}

func (this *PathBuilder) newStorageRoot(account types.Address, txIndex uint32) bool {
	accountRoot := common.StrCat(ccurlcommon.ETH10_ACCOUNT_PREFIX, string(account), "/")
	if !this.apiRouter.WriteCache().(*cache.WriteCache).IfExists(accountRoot) {
		return common.FilterFirst(this.apiRouter.WriteCache().(*cache.WriteCache).CreateNewAccount(txIndex, string(account))) != nil // Create a new account
	}
	return true // ALready exists
}

func (this *PathBuilder) newContainerRoot(account types.Address, txIndex uint32) bool {
	containerRoot := this.key(account)

	if !this.apiRouter.WriteCache().(*cache.WriteCache).IfExists(containerRoot) {
		_, err := this.apiRouter.WriteCache().(*cache.WriteCache).Write(txIndex, containerRoot, commutative.NewPath()) // Create a new container
		return err == nil
	}
	return true // Already exists
}

func (this *PathBuilder) Key(caller [20]byte) string { // container ID
	return this.key(types.Address(hexutil.Encode(caller[:])))
}

func (this *PathBuilder) key(account types.Address) string { // container ID
	return common.StrCat(ccurlcommon.ETH10_ACCOUNT_PREFIX, string(account), "/storage", this.subDir, "/")
}
