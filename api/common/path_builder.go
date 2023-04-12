package common

import (
	commonlib "github.com/arcology-network/common-lib/common"
	"github.com/arcology-network/common-lib/types"
	ethCommon "github.com/arcology-network/evm/common"

	"github.com/arcology-network/concurrenturl/v2"
	commutative "github.com/arcology-network/concurrenturl/v2/type/commutative"
)

// Ccurl connectors for Arcology APIs
type CcurlConnector struct {
	txHash  ethCommon.Hash
	txIndex uint32
	ccurl   *concurrenturl.ConcurrentUrl
	prefix  string
}

func NewCCurlConnector(prefix string, txHash ethCommon.Hash, txIndex uint32, ccurl *concurrenturl.ConcurrentUrl) *CcurlConnector {
	return &CcurlConnector{
		prefix:  prefix,
		txHash:  txHash,
		txIndex: txIndex,
		ccurl:   ccurl,
	}
}

func (this *CcurlConnector) Ccurl() *concurrenturl.ConcurrentUrl { return this.ccurl }

// Make Arcology paths under the current account
func (this *CcurlConnector) New(account types.Address, containerId string, keyType int) bool {
	if !this.newStorageRoot(account, this.txIndex) { // Create the root path if has been created yet.
		return false
	}

	if !this.newContainerRoot(account, containerId, this.txIndex) { // Create the container path if has been created yet.
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

	if path, err := commutative.NewMeta(containerRoot); err != nil {
		return false
	} else if err := this.ccurl.Write(txIndex, containerRoot, path); err != nil {
		return false
	}
	return true
}

func (this *CcurlConnector) Key(account types.Address, id string) string { // container ID
	return commonlib.StrCat(this.ccurl.Platform.Eth10Account(), string(account), this.prefix, id, "/")
}

// func (this *CcurlConnector) buildContainerLength(account types.Address, id string) string {
// 	return commonlib.StrCat(this.ccurl.Platform.Eth10Account(), string(account), "/storage/containers/", id, "/")

// 	if value, err := this.url.Read(this.context.GetIndex(), BuildContainerRootPath(this.url, account, id)); err != nil || value == nil {
// 		return ContainerSizeInvalid
// 	} else {
// 		return len(value.(*commutative.Meta).KeyView()) - 2
// 	}
// }

// func (this *CcurlConnector) buildContainerTypePath(account types.Address, id string) string {
// 	return commonlib.StrCat(this.ccurl.Platform.Eth10Account(), string(account), "/storage/containers/!/", id)
// }

// func (this *CcurlConnector) buildSizePath(account types.Address, id string) string {
// 	return commonlib.StrCat(this.ccurl.Platform.Eth10Account(), string(account), "/storage/containers/", id, "/#")
// }

// func (this *CcurlConnector) buildKeyTypePath(account types.Address, id string) string {
// 	return commonlib.StrCat(this.ccurl.Platform.Eth10Account(), string(account), "/storage/containers/", id, "/!")
// }

// func (this *CcurlConnector) buildValueTypePath(account types.Address, id string) string {
// 	return commonlib.StrCat(this.ccurl.Platform.Eth10Account(), string(account), "/storage/containers/", id, "/@")
// }

// func (this *CcurlConnector) buildContainerType(url *concurrenturl.ConcurrentUrl, account types.Address, id string, txIndex uint32) int {
// 	if value, err := url.Read(txIndex, this.buildContainerTypePath(account, id)); err != nil || value == nil {
// 		return -1
// 	} else {
// 		return int(*value.(*noncommutative.Int64))
// 	}
// }

// func buildDefaultValue(dataType int) ([]byte, bool) {
// 	switch dataType {
// 	case DataTypeAddress:
// 		return common.Address{}.Bytes(), true
// 	case DataTypeUint256:
// 		return common.Hash{}.Bytes(), true
// 	default:
// 		return nil, true
// 	}
// }
