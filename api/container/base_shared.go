package api

import (
	"math"

	"github.com/arcology-network/common-lib/common"
	orderedset "github.com/arcology-network/common-lib/container/set"
	"github.com/arcology-network/concurrenturl/commutative"
	"github.com/arcology-network/concurrenturl/noncommutative"
	"github.com/arcology-network/eu/cache"
	intf "github.com/arcology-network/vm-adaptor/interface"
)

// Get the number of elements in the container
func (this *BaseHandlers) Length(path string) (uint64, bool, int64) {
	if len(path) == 0 {
		return 0, false, 0
	}

	if path, _, _ := this.api.WriteCache().(*cache.WriteCache).Read(this.api.GetEU().(intf.EU).ID(), path, new(commutative.Path)); path != nil {
		keys := path.(*orderedset.OrderedSet).Keys()
		return uint64(len(keys)), true, 0
	}
	return 0, false, 0
}

// Get the index of the element by its key
func (this *BaseHandlers) GetByIndex(path string, idx uint64) ([]byte, bool, int64) {
	if value, _, err := this.api.WriteCache().(*cache.WriteCache).ReadAt(this.api.GetEU().(intf.EU).ID(), path, idx, new(noncommutative.Bytes)); err == nil && value != nil {
		return value.([]byte), true, 0
	}
	return []byte{}, false, 0
}

// Set the element by its index
func (this *BaseHandlers) SetByIndex(path string, idx uint64, bytes []byte) (bool, int64) {
	if len(path) > 0 {
		value := common.IfThen(bytes == nil, nil, noncommutative.NewBytes(bytes))
		if _, err := this.api.WriteCache().(*cache.WriteCache).WriteAt(this.api.GetEU().(intf.EU).ID(), path, idx, value); err == nil {
			return true, 0
		}
	}
	return false, 0
}

// Get the element by its key
func (this *BaseHandlers) GetByKey(path string) ([]byte, bool, int64) {
	if value, _, _ := this.api.WriteCache().(*cache.WriteCache).Read(this.api.GetEU().(intf.EU).ID(), path, new(noncommutative.Bytes)); value != nil {
		return value.([]byte), true, 0
	}
	return []byte{}, false, 0
}

// Set the element by its key
func (this *BaseHandlers) SetByKey(path string, bytes []byte) (bool, int64) {
	if len(path) > 0 {
		value := common.IfThen(bytes == nil, nil, noncommutative.NewBytes(bytes))
		if _, err := this.api.WriteCache().(*cache.WriteCache).Write(this.api.GetEU().(intf.EU).ID(), path, value); err == nil {
			return true, 0
		}
	}
	return false, 0
}

// Get the index of a key
func (this *BaseHandlers) KeyAt(path string, index uint64) (string, int64) {
	if len(path) > 0 {
		key, _ := this.api.WriteCache().(*cache.WriteCache).KeyAt(this.api.GetEU().(intf.EU).ID(), path, index, new(noncommutative.Bytes))
		return key, 0
	}
	return "", 0
}

// Get the index of a key
func (this *BaseHandlers) IndexOf(path string, key string) (uint64, int64) {
	if len(path) > 0 {
		index, _ := this.api.WriteCache().(*cache.WriteCache).IndexOf(this.api.GetEU().(intf.EU).ID(), path, key, new(noncommutative.Bytes))
		return index, 0
	}
	return math.MaxUint64, 0
}
