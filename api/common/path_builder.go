package common

import (
	commonlib "github.com/arcology-network/common-lib/common"
	"github.com/arcology-network/common-lib/types"
	ethCommon "github.com/arcology-network/evm/common"

	"github.com/arcology-network/concurrenturl"
	commutative "github.com/arcology-network/concurrenturl/commutative"
)

// Ccurl connectors for Arcology APIs
type CcurlConnector struct {
	txHash  ethCommon.Hash
	txIndex uint32
	ccurl   *concurrenturl.ConcurrentUrl
	subDir  string
}

func NewCCurlConnector(subDir string, txHash ethCommon.Hash, txIndex uint32, ccurl *concurrenturl.ConcurrentUrl) *CcurlConnector {
	return &CcurlConnector{
		subDir:  subDir,
		txHash:  txHash,
		txIndex: txIndex,
		ccurl:   ccurl,
	}
}

func (this *CcurlConnector) TxHash() ethCommon.Hash              { return this.txHash }
func (this *CcurlConnector) TxIndex() uint32                     { return this.txIndex }
func (this *CcurlConnector) Ccurl() *concurrenturl.ConcurrentUrl { return this.ccurl }

// Make Arcology paths under the current account
func (this *CcurlConnector) New(account types.Address, containerId string, keyType int) bool {
	if !this.newStorageRoot(account, this.txIndex) { // Create the root path if has been created yet.
		return false
	}

	// Create the container path if has been created yet.
	if !this.newContainerRoot(account, containerId[:], this.txIndex) {
		return false
	}

	return true
}

func (this *CcurlConnector) newStorageRoot(account types.Address, txIndex uint32) bool {
	accountRoot := commonlib.StrCat(this.ccurl.Platform.Eth10Account(), string(account), "/")
	if value, err := this.ccurl.Peek(accountRoot); err != nil {
		return false
	} else if value == nil { // The account didn't exist.
		if err := this.ccurl.CreateAccount(txIndex, this.ccurl.Platform.Eth10(), string(account)); err != nil {
			return false
		}
	}

	return true
}

func (this *CcurlConnector) newContainerRoot(account types.Address, id string, txIndex uint32) bool {
	containerRoot := this.Key(account, id)
	if value, err := this.ccurl.Peek(containerRoot); err != nil || value != nil {
		return false
	}

	if err := this.ccurl.Write(txIndex, containerRoot, commutative.NewPath()); err != nil {
		return false
	}

	return true
}

func (this *CcurlConnector) Key(account types.Address, id string) string { // container ID
	return commonlib.StrCat(this.ccurl.Platform.Eth10Account(), string(account), "/storage", this.subDir, id, "/")
}
