package evm

import (
	"bytes"
	"testing"

	"github.com/arcology-network/evm/crypto"
)

func TestBytesToInt32(t *testing.T) {
	input := []byte{0, 0, 0, 0xcc}
	output := BytesToInt32(input)
	if output != 0xcc {
		t.Error("failed")
	}
}

func TestInt64ToBytes(t *testing.T) {
	value := 0x1122334455667788
	output := Int64ToBytes(int64(value))
	if bytes.Compare(output, []byte{0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88}) != 0 {
		t.Error("failed")
	}

	value = -0x7feeddccbbaa9988
	output = Int64ToBytes(int64(value))
	if bytes.Compare(output, []byte{0x80, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x78}) != 0 {
		t.Error("failed")
	}
}

func TestPrintMethodID(t *testing.T) {
	// ConcurrentArray
	printMethodID("ConcurrentArray", "create(string,int32,int32)", t)
	printMethodID("ConcurrentArray", "size(string)", t)
	printMethodID("ConcurrentArray", "set(string,int32,uint256)", t)
	printMethodID("ConcurrentArray", "getUint256(string,int32)", t)
	printMethodID("ConcurrentArray", "set(string,int32,address)", t)
	printMethodID("ConcurrentArray", "getAddress(string,int32)", t)
	printMethodID("ConcurrentArray", "set(string,int32,bytes)", t)
	printMethodID("ConcurrentArray", "getBytes(string,int32)", t)
	// ConcurrentVariable
	printMethodID("ConcurrentVariable", "create(string,int32)", t)
	printMethodID("ConcurrentVariable", "set(string,uint256)", t)
	printMethodID("ConcurrentVariable", "getUint256(string)", t)
	printMethodID("ConcurrentVariable", "set(string,address)", t)
	printMethodID("ConcurrentVariable", "getAddress(string)", t)
	printMethodID("ConcurrentVariable", "set(string,bytes)", t)
	printMethodID("ConcurrentVariable", "getBytes(string)", t)
	// ConcurrentHashMap
	printMethodID("ConcurrentHashMap", "create(string,int32,int32)", t)
	printMethodID("ConcurrentHashMap", "getAddress(string,uint256)", t)
	printMethodID("ConcurrentHashMap", "set(string,uint256,address)", t)
	printMethodID("ConcurrentHashMap", "getUint256(string,uint256)", t)
	printMethodID("ConcurrentHashMap", "set(string,uint256,uint256)", t)
	printMethodID("ConcurrentHashMap", "getBytes(string,uint256)", t)
	printMethodID("ConcurrentHashMap", "set(string,uint256,bytes)", t)
	printMethodID("ConcurrentHashMap", "getAddress(string,address)", t)
	printMethodID("ConcurrentHashMap", "set(string,address,address)", t)
	printMethodID("ConcurrentHashMap", "getUint256(string,address)", t)
	printMethodID("ConcurrentHashMap", "set(string,address,uint256)", t)
	printMethodID("ConcurrentHashMap", "getBytes(string,address)", t)
	printMethodID("ConcurrentHashMap", "set(string,address,bytes)", t)
	printMethodID("ConcurrentHashMap", "deleteKey(string,uint256)", t)
	printMethodID("ConcurrentHashMap", "deleteKey(string,address)", t)
	// UUID
	printMethodID("UUID", "gen(string)", t)
	// System
	printMethodID("System", "getPid()", t)
	printMethodID("System", "revertPid(uint256)", t)
	printMethodID("System", "create(string,string)", t)
	printMethodID("System", "call(string)", t)
	// Queue
	printMethodID("Queue", "create(string,uint256)", t)
	printMethodID("Queue", "size(string)", t)
	printMethodID("Queue", "pushUint256(string,uint256)", t)
	printMethodID("Queue", "popUint256(string)", t)
	// Others
	printMethodID("Others", "defer(string)", t)
}

func printMethodID(className, signature string, t *testing.T) {
	t.Logf("%v.%v: %x", className, signature, crypto.Keccak256([]byte(signature))[:4])
}
