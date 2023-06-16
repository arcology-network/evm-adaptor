package execution

import (
	"strings"

	indexer "github.com/arcology-network/concurrenturl/indexer"
	ccurlinterfaces "github.com/arcology-network/concurrenturl/interfaces"
)

// APIs under the concurrency namespace
type TransitionFilter struct {
	indexer.IPCTransition
	status uint8
}

// Remove nonce
func (this TransitionFilter) From(univ ccurlinterfaces.Univalue) interface{} {
	v := indexer.IPCTransition{Status: this.status}.From(univ)
	if v == nil || v.(ccurlinterfaces.Univalue).Value() == nil {
		return v
	}

	if strings.HasSuffix(*v.(ccurlinterfaces.Univalue).GetPath(), "/nonce") {
		return nil
	}
	return v
}
