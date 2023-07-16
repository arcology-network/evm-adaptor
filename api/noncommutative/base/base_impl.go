package concurrentcontainer

import (
	"github.com/arcology-network/concurrenturl/noncommutative"
	"github.com/arcology-network/vm-adaptor/execution"
)

// // get the number of elements in the container
func (this *BytesHandlers) Length(path string) (uint64, bool, int64) {
	if len(path) == 0 {
		return 0, false, 0
	}

	if path, _ := this.api.Ccurl().Read(uint32(this.api.GetEU().(*execution.EU).Message().ID), path); path != nil {
		return uint64(len(path.([]string))), true, 0
	}
	return 0, false, 0
}

// // get the number of elements in the container
func (this *BytesHandlers) GetIndex(path string, idx uint64) ([]byte, bool, int64) {
	if value, _, err := this.api.Ccurl().ReadAt(uint32(this.api.GetEU().(*execution.EU).Message().ID), path, idx); err == nil && value != nil {
		return value.([]byte), true, 0
	}
	return []byte{}, false, 0
}

func (this *BytesHandlers) SetIndex(path string, idx uint64, bytes []byte) (bool, int64) {
	if len(path) > 0 {
		value := noncommutative.NewBytes(bytes)
		if _, err := this.api.Ccurl().WriteAt(uint32(this.api.GetEU().(*execution.EU).Message().ID), path, idx, value, true); err == nil {
			return true, 0
		}
	}
	return false, 0
}

func (this *BytesHandlers) GetKey(path string) ([]byte, bool, int64) {
	if value, fees := this.api.Ccurl().Read(uint32(this.api.GetEU().(*execution.EU).Message().ID), path); value != nil {
		return value.([]byte), true, int64(fees)
	}
	return []byte{}, false, 0
}

func (this *BytesHandlers) SetKey(path string, bytes []byte) (bool, int64) {
	if len(path) > 0 {
		value := noncommutative.NewBytes(bytes)
		if _, err := this.api.Ccurl().Write(uint32(this.api.GetEU().(*execution.EU).Message().ID), path, value, true); err == nil {
			return true, 0
		}
	}
	return false, 0
}
