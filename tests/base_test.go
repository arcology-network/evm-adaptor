/*
 *   Copyright (c) 2024 Arcology Network

 *   This program is free software: you can redistribute it and/or modify
 *   it under the terms of the GNU General Public License as published by
 *   the Free Software Foundation, either version 3 of the License, or
 *   (at your option) any later version.

 *   This program is distributed in the hope that it will be useful,
 *   but WITHOUT ANY WARRANTY; without even the implied warranty of
 *   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *   GNU General Public License for more details.

 *   You should have received a copy of the GNU General Public License
 *   along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package api

import (
	"bytes"
	"fmt"
	"math"
	"testing"

	"github.com/arcology-network/evm-adaptor/abi"
	apihandler "github.com/arcology-network/evm-adaptor/apihandler"
	base "github.com/arcology-network/evm-adaptor/apihandler/container"
	cache "github.com/arcology-network/storage-committer/storage/writecache"
	"github.com/holiman/uint256"
)

type MockID struct{}

func (MockID) ID() uint32 { return 0 }

func TestBaseHandlers(t *testing.T) {
	api := apihandler.NewAPIHandler(cache.NewWriteCache(nil))
	api.SetEU(MockID{})
	baseContainer := base.NewBaseHandlers(api)

	// Create a new container
	baseContainer.Call(
		[20]byte{0xc2},
		[20]byte{0xc2},
		[]byte{0xcd, 0xbf, 0x60, 0x8d},
		[20]byte{},
		0) // Nonce

	// Push a new element by calling setByKey()
	data, _ := abi.Encode(uint256.NewInt(11))
	d, _ := abi.Decode(data, 0, new(uint256.Int), 1, math.MaxInt)
	// if d, _ := abi.Decode(data, 0, new(uint256.Int), 1, math.MaxInt); !u256.Eq(d.(*uint256.Int)) {
	// 	t.Error("Error: Should be equal!")
	// }

	fmt.Println(d)

	baseContainer.Call(
		[20]byte{0xc2},
		[20]byte{0xc2},
		append([]byte{0xc2, 0x78, 0xb7, 0x99}, data...),
		[20]byte{},
		0)

	// Get length
	data, _, _ = baseContainer.Call(
		[20]byte{0xc2},
		[20]byte{0xc2},
		[]byte{0x1f, 0x7b, 0x6d, 0x32},
		[20]byte{},
		0) // Nonce

	buffer := [32]byte{}
	if bytes.Equal(data, buffer[:]) {
		t.Log("Success")
	} else {
		t.Error("Failed")
	}
}
