package api

import (
	"math"

	"github.com/arcology-network/common-lib/common"
	"github.com/arcology-network/common-lib/exp/deltaset"
	"github.com/arcology-network/common-lib/exp/slice"
	"github.com/arcology-network/storage-committer/commutative"
	"github.com/arcology-network/storage-committer/noncommutative"
	cache "github.com/arcology-network/storage-committer/storage/writecache"
)

// Get the number of elements in the container, including the nil elements.
func (this *BaseHandlers) Length(path string) (uint64, bool, int64) {
	if len(path) == 0 {
		return 0, false, 0
	}

	if path, _, _ := this.api.WriteCache().(*cache.WriteCache).Read(this.api.GetEU().(interface{ ID() uint32 }).ID(), path, new(commutative.Path)); path != nil {
		return path.(*deltaset.DeltaSet[string]).NonNilCount(), true, 0
	}
	return 0, false, 0
}

// Export all the elements in the container to a two-dimensional slice.
// This function will read all the elements in the container.
func (this *BaseHandlers) ReadAll(path string) ([][]byte, []bool, []int64) {
	length, _, _ := this.Length(path)
	entries := make([][]byte, length)
	flags := make([]bool, length)
	fees := make([]int64, length)

	slice.NewDo(int(length), func(i int) []byte {
		entries[i], flags[i], fees[i] = this.GetByIndex(path, uint64(i))
		return []byte{}
	})
	return entries, flags, fees
}

// Get the index of the element by its key
func (this *BaseHandlers) GetByIndex(path string, idx uint64) ([]byte, bool, int64) {
	if value, _, err := this.api.WriteCache().(*cache.WriteCache).ReadAt(
		this.api.GetEU().(interface{ ID() uint32 }).ID(),
		path,
		idx,
		new(noncommutative.Bytes),
	); err == nil && value != nil {
		return value.([]byte), true, 0
	}
	return []byte{}, false, 0
}

// Set the element by its index
func (this *BaseHandlers) SetByIndex(path string, idx uint64, bytes []byte) (bool, int64) {
	if len(path) == 0 {
		return false, 0
	}

	value := common.IfThen(bytes == nil, nil, noncommutative.NewBytes(bytes))
	if _, err := this.api.WriteCache().(*cache.WriteCache).WriteAt(this.api.GetEU().(interface{ ID() uint32 }).ID(), path, idx, value); err == nil {
		return true, 0
	}
	return false, 0
}

// Get the element by its key
func (this *BaseHandlers) GetByKey(path string) ([]byte, bool, int64) {
	if value, _, _ := this.api.WriteCache().(*cache.WriteCache).Read(this.api.GetEU().(interface{ ID() uint32 }).ID(), path, new(noncommutative.Bytes)); value != nil {
		return value.([]byte), true, 0
	}
	return []byte{}, false, 0
}

// Set the element by its key
func (this *BaseHandlers) SetByKey(path string, bytes []byte) (bool, int64) {
	if len(path) > 0 {
		value := common.IfThen(bytes == nil, nil, noncommutative.NewBytes(bytes))
		if _, err := this.api.WriteCache().(*cache.WriteCache).Write(this.api.GetEU().(interface{ ID() uint32 }).ID(), path, value); err == nil {
			return true, 0
		}
	}
	return false, 0
}

// Get the index of a key
func (this *BaseHandlers) KeyAt(path string, index uint64) (string, int64) {
	if len(path) > 0 {
		key, _ := this.api.WriteCache().(*cache.WriteCache).KeyAt(this.api.GetEU().(interface{ ID() uint32 }).ID(), path, index, new(noncommutative.Bytes))
		return key, 0
	}
	return "", 0
}

// Get the index of a key
func (this *BaseHandlers) IndexOf(path string, key string) (uint64, int64) {
	if len(path) > 0 {
		index, _ := this.api.WriteCache().(*cache.WriteCache).IndexOf(this.api.GetEU().(interface{ ID() uint32 }).ID(), path, key, new(noncommutative.Bytes))
		return index, 0
	}
	return math.MaxUint64, 0
}
