package common

import "encoding/hex"

const (
	MAX_RECURSIION_DEPTH = uint8(4)
	MAX_SUB_PROCESSES    = uint64(2048)
)

const (
	SUB_PROCESS = iota
	CONTAINER_ID
	ELEMENT_ID
	UUID
)

var IO_HANDLER = [20]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0x60}
var BYTES_HANDLER = [20]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0x84}
var CUMULATIVE_U256_HANDLER = [20]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0x85}
var CUMULATIVE_I256_HANDLER = [20]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0x86}
var MULTIPROCESS_HANDLER = [20]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0xb0}
var RUNTIME_HANDLER = [20]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0xa0}

var TotalSubProcesses uint64

func ToValidName(bytes []byte) string {
	for _, c := range bytes {
		if !(('0' <= c && c <= '9') || ('a' <= c && c <= 'f') || ('A' <= c && c <= 'F')) && (c != '-') {
			return hex.EncodeToString(bytes)
		}
	}
	return string(bytes)
}
