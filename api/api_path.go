package api

import (
	"fmt"

	commonlib "github.com/arcology-network/common-lib/common"
	"github.com/arcology-network/common-lib/types"
	ethCommon "github.com/arcology-network/evm/common"

	"github.com/arcology-network/concurrenturl/v2"
	commutative "github.com/arcology-network/concurrenturl/v2/type/commutative"
)

// Ccurl connectors for Arcology APIs
type CCurlPathBuilder struct {
	txHash  ethCommon.Hash
	txIndex uint32
	ccurl   *concurrenturl.ConcurrentUrl
}

func NewCCurlPathBuilder(txHash ethCommon.Hash, txIndex uint32, ccurl *concurrenturl.ConcurrentUrl) *CCurlPathBuilder {
	return &CCurlPathBuilder{
		txHash:  txHash,
		txIndex: txIndex,
		ccurl:   ccurl,
	}
}

// Make Arcology paths under the current account
func (this *CCurlPathBuilder) New(account types.Address, containerId string, keyType int) bool {
	if !this.makeStorageRootPath(account, this.txIndex) { // Create the root path if has been created yet.
		return false
	}

	if !this.makeContainerRootPath(account, containerId, this.txIndex) { // Create the container path if has been created yet.
		return false
	}

	// Write the container meta data.
	// if err := this.ccurl.Write(this.txIndex, this.buildKeyTypePath(account, containerId), noncommutative.NewInt64(int64(keyType))); err != nil {
	// 	return false
	// }

	// if err := this.ccurl.Write(this.txIndex, this.buildValueTypePath(account, containerId), noncommutative.NewInt64(int64(valueType))); err != nil {
	// 	return false
	// }
	return true
}

func (this *CCurlPathBuilder) makeStorageRootPath(account types.Address, txIndex uint32) bool {
	accountRoot := this.buildAccountRootPath(account)
	if value, err := this.ccurl.Peek(accountRoot); err != nil {
		return false
	} else if value == nil { // The account didn't exist.
		if err := this.ccurl.CreateAccount(txIndex, this.ccurl.Platform.Eth10(), string(account)); err != nil {
			return false
		}
	}

	return true
}

func (this *CCurlPathBuilder) makeContainerRootPath(account types.Address, id string, txIndex uint32) bool {
	containerRoot := this.buildContainerRootPath(account, id)
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

func (this *CCurlPathBuilder) buildAccountRootPath(account types.Address) string {
	return commonlib.StrCat(this.ccurl.Platform.Eth10Account(), string(account), "/")
}

func (this *CCurlPathBuilder) buildDeferCallPath(account types.Address, id string) string {
	return commonlib.StrCat(this.ccurl.Platform.Eth10Account(), string(account), "/defer/", id)
}

func (this *CCurlPathBuilder) buildContainerRootPath(account types.Address, id string) string {
	return commonlib.StrCat(this.ccurl.Platform.Eth10Account(), string(account), "/storage/containers/", id, "/")
}

func (this *CCurlPathBuilder) buildValuePath(account types.Address, id string, key interface{}) string {
	return fmt.Sprintf("%s%v", this.buildContainerRootPath(account, id), key)
}

// func (this *CCurlPathBuilder) buildContainerLength(account types.Address, id string) string {
// 	return commonlib.StrCat(this.ccurl.Platform.Eth10Account(), string(account), "/storage/containers/", id, "/")

// 	if value, err := this.url.Read(this.context.GetIndex(), buildContainerRootPath(this.url, account, id)); err != nil || value == nil {
// 		return ContainerSizeInvalid
// 	} else {
// 		return len(value.(*commutative.Meta).KeyView()) - 2
// 	}
// }

// func (this *CCurlPathBuilder) buildContainerTypePath(account types.Address, id string) string {
// 	return commonlib.StrCat(this.ccurl.Platform.Eth10Account(), string(account), "/storage/containers/!/", id)
// }

// func (this *CCurlPathBuilder) buildSizePath(account types.Address, id string) string {
// 	return commonlib.StrCat(this.ccurl.Platform.Eth10Account(), string(account), "/storage/containers/", id, "/#")
// }

// func (this *CCurlPathBuilder) buildKeyTypePath(account types.Address, id string) string {
// 	return commonlib.StrCat(this.ccurl.Platform.Eth10Account(), string(account), "/storage/containers/", id, "/!")
// }

// func (this *CCurlPathBuilder) buildValueTypePath(account types.Address, id string) string {
// 	return commonlib.StrCat(this.ccurl.Platform.Eth10Account(), string(account), "/storage/containers/", id, "/@")
// }

// func (this *CCurlPathBuilder) buildContainerType(url *concurrenturl.ConcurrentUrl, account types.Address, id string, txIndex uint32) int {
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
