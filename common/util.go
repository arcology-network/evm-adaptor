package common

import (
	"reflect"

	"github.com/arcology-network/storage-committer/commutative"
	"github.com/arcology-network/storage-committer/noncommutative"
	"github.com/arcology-network/storage-committer/platform"
	"github.com/arcology-network/storage-committer/univalue"
)

// CreateNewAccount creates a new account in the write cache.
// It returns the transitions and an error, if any.
func CreateNewAccount(tx uint32, acct string, store interface {
	IfExists(string) bool
	Write(uint32, string, interface{}) (int64, error)
}) ([]*univalue.Univalue, error) {
	paths, typeids := platform.NewPlatform().GetBuiltins(acct)

	transitions := []*univalue.Univalue{}
	for i, path := range paths {
		var v interface{}
		switch typeids[i] {
		case commutative.PATH: // Path
			v = commutative.NewPath()

		case uint8(reflect.Kind(noncommutative.STRING)): // delta big int
			v = noncommutative.NewString("")

		case uint8(reflect.Kind(commutative.UINT256)): // delta big int
			v = commutative.NewUnboundedU256()

		case uint8(reflect.Kind(commutative.UINT64)):
			v = commutative.NewUnboundedUint64()

		case uint8(reflect.Kind(noncommutative.INT64)):
			v = new(noncommutative.Int64)

		case uint8(reflect.Kind(noncommutative.BYTES)):
			v = noncommutative.NewBytes([]byte{})
		}

		// fmt.Println(path)
		if !store.IfExists(path) {
			transitions = append(transitions, univalue.NewUnivalue(tx, path, 0, 1, 0, v, nil))

			if _, err := store.Write(tx, path, v); err != nil { // root path
				return nil, err
			}

			if !store.IfExists(path) {
				_, err := store.Write(tx, path, v)
				return transitions, err // root path
			}
		}
	}
	return transitions, nil
}
