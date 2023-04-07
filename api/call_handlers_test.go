package api

import (
	"bytes"
	"testing"

	"github.com/arcology-network/common-lib/codec"
)

func TestUUID(t *testing.T) {
	api := NewAPI(nil)
	id0 := codec.Hash32([32]byte{}).UUID(999999999999)
	id1 := codec.Hash32([32]byte{}).UUID(999999999999)
	id2 := codec.Hash32([32]byte{}).UUID(999999999999)

	if id2 == [32]byte{} || id0 != id1 || id0 != id2 {
		t.Error("Error!") // same with the same seed
	}

	id3 := api.GenUUID()
	id4 := api.GenUUID()
	id5 := api.GenUUID()

	if bytes.Equal(id3, id4) || bytes.Equal(id3, id5) || bytes.Equal(id4, id5) {
		t.Error("Error!") // Should be different with different seeds
	}
}
