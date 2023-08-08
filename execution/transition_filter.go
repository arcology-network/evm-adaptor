package execution

import (
	"encoding/hex"
	"strings"

	indexer "github.com/arcology-network/concurrenturl/indexer"
	ccurlinterfaces "github.com/arcology-network/concurrenturl/interfaces"
	"github.com/arcology-network/concurrenturl/univalue"
)

// APIs under the concurrency namespace
type TransitionFilter struct {
	indexer.ITCTransition
	Sender   [20]byte
	Coinbase [20]byte
	Err      error
}

// Remove nonce
func (this TransitionFilter) From(univ ccurlinterfaces.Univalue) interface{} {
	path := *univ.GetPath()
	if (strings.HasSuffix(path, "/balance") || strings.HasSuffix(path, "/nonce")) &&
		(strings.Contains(path, hex.EncodeToString(this.Sender[:])) || strings.Contains(path, hex.EncodeToString(this.Coinbase[:]))) {
		univ.GetUnimeta().(*univalue.Unimeta).SetPersistent(true) // Keep balance transitions regardless the execution status
	}

	v := indexer.ITCTransition{Err: this.Err}.From(univ)
	if v == nil || v.(ccurlinterfaces.Univalue).Value() == nil {
		return v
	}
	return v
}
