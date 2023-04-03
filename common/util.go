package common

import (
	"fmt"
	"math/big"
	"math/rand"

	urlcommon "github.com/arcology-network/concurrenturl/v2/common"
	urltype "github.com/arcology-network/concurrenturl/v2/type"
	"github.com/arcology-network/concurrenturl/v2/type/commutative"
	"github.com/arcology-network/concurrenturl/v2/type/noncommutative"
	arbitrator "github.com/arcology-network/urlarbitrator-engine/go-wrapper"
	"github.com/holiman/uint256"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func FormatValue(value interface{}) string {
	switch value.(type) {
	case *commutative.Meta:
		meta := value.(*commutative.Meta)
		var str string
		str += "{"
		for i, k := range meta.Keys() {
			str += k.(string)
			if i != len(meta.Keys())-1 {
				str += ", "
			}
		}
		str += "}"
		if len(meta.Added()) != 0 {
			str += " + {"
			for i, k := range meta.Added() {
				str += k
				if i != len(meta.Added())-1 {
					str += ", "
				}
			}
			str += "}"
		}
		if len(meta.Removed()) != 0 {
			str += " - {"
			for i, k := range meta.Removed() {
				str += k
				if i != len(meta.Removed())-1 {
					str += ", "
				}
			}
			str += "}"
		}
		return str
	case *noncommutative.Int64:
		// uint256.NewInt(0)
		return fmt.Sprintf(" = %v", (*(value.(*uint256.Int))))
	case *noncommutative.Bytes:
		return fmt.Sprintf(" = %v", value.(*noncommutative.Bytes).Data())
	case *commutative.U256:
		v := value.(*commutative.U256).Value()
		d := value.(*commutative.U256).GetDelta()
		return fmt.Sprintf(" = %v + %v", (*(v.(*uint256.Int))), d.(*big.Int).Int64())
	case *commutative.Int64:
		v := value.(*commutative.Int64).Value()
		d := value.(*commutative.Int64).GetDelta()
		return fmt.Sprintf(" = %v + %v", v, d)
	}
	return ""
}

func FormatTransitions(transitions []urlcommon.UnivalueInterface) string {
	var str string
	for _, t := range transitions {
		str += fmt.Sprintf("%v%v%v%v%v%v%v%v%v%v%v%v%v%v",
			"Tx=", t.(*urltype.Univalue).GetTx(),
			" Reads=", t.(*urltype.Univalue).Reads(),
			" Writes=", t.(*urltype.Univalue).Writes(),
			" Preexists=", t.(*urltype.Univalue).Preexist(),
			" Composite=", t.(*urltype.Univalue).Composite(),
			" Path=", *(t.(*urltype.Univalue).GetPath()),
			" Value", FormatValue(t.(*urltype.Univalue).Value())+"\n")
	}
	return str
}

func DetectConflict(transitions []urlcommon.UnivalueInterface) ([]uint32, []uint32, []bool) {
	length := len(transitions)
	txs := make([]uint32, length)
	paths := make([]string, length)
	reads := make([]uint32, length)
	writes := make([]uint32, length)
	composite := make([]bool, length)
	uniqueTxsDict := make(map[uint32]struct{})
	for i, t := range transitions {
		txs[i] = t.(*urltype.Univalue).GetTx()
		paths[i] = *(t.(*urltype.Univalue).GetPath())
		reads[i] = t.(*urltype.Univalue).Reads()
		writes[i] = t.(*urltype.Univalue).Writes()
		composite[i] = t.(*urltype.Univalue).Composite()
		uniqueTxsDict[txs[i]] = struct{}{}
	}

	uniqueTxs := make([]uint32, 0, len(uniqueTxsDict))
	for tx := range uniqueTxsDict {
		uniqueTxs = append(uniqueTxs, tx)
	}
	engine := arbitrator.Start()
	arbitrator.Insert(engine, txs, paths, reads, writes, composite)
	txs, groups, flags := arbitrator.DetectLegacy(engine, uniqueTxs)
	arbitrator.Clear(engine)
	return txs, groups, flags
}
