package common

import (
	"fmt"
	"math/rand"

	"github.com/arcology-network/storage-committer/commutative"
	"github.com/arcology-network/storage-committer/univalue"
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
	case *commutative.Path:
		meta := value.(*commutative.Path)
		var str string
		str += "{"
		for i, k := range meta.Committed().Elements() {
			str += k
			if i != meta.Committed().Length()-1 {
				str += ", "
			}
		}
		str += "}"
		if len(meta.Committed().Elements()) != 0 {
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
		// case *noncommutative.Int64:
		// 	// uint256.NewInt(0)
		// 	return fmt.Sprintf(" = %v", (*(value.(*codec.Int64))))
		// case *noncommutative.Bytes:
		// 	return fmt.Sprintf(" = %v", value.(*noncommutative.Bytes).Value())
		// case *commutative.U256:
		// 	v := value.(*commutative.U256).Value()
		// 	d := value.(*commutative.U256).Delta()
		// 	return fmt.Sprintf(" = %v + %v", (*(v.(*codec.Uint256))), d.(*codec.Uint256).Uint64())
		// case *commutative.Uint64:
		// 	v := value.(*commutative.Uint64).Value()
		// 	d := value.(*commutative.Uint64).Delta()
		// 	return fmt.Sprintf(" = %v + %v", v, d)
	}
	return ""
}

func FormatTransitions(transitions []*univalue.Univalue) string {
	var str string
	for _, t := range transitions {
		str += fmt.Sprintf("%v%v%v%v%v%v%v%v%v%v%v%v%v%v",
			"Tx=", t.GetTx(),
			" Reads=", t.Reads(),
			" Writes=", t.Writes(),
			" Delta Writes=", t.DeltaWrites(),
			" Preexists=", t.Preexist(),
			" Path=", *(t.GetPath()),
			" Value", FormatValue(t.Value())+"\n")
	}
	return str
}

// func DetectConflict(transitions [] *univalue.Univalue) ([]uint32, []uint32, []bool) {
// 	length := len(transitions)
// 	txs := make([]uint32, length)
// 	paths := make([]string, length)
// 	reads := make([]uint32, length)
// 	writes := make([]uint32, length)
// 	composite := make([]bool, length)
// 	uniqueTxsDict := make(map[uint32]struct{})
// 	for i, t := range transitions {
// 		txs[i] = t.GetTx()
// 		paths[i] = *(t.GetPath())
// 		reads[i] = t.Reads()
// 		writes[i] = t.Writes()
// 		composite[i] = t.Reads() == 0 && t.Writes() == 0 && t.DeltaWrites() >= 0
// 		uniqueTxsDict[txs[i]] = struct{}{}
// 	}

// 	uniqueTxs := make([]uint32, 0, len(uniqueTxsDict))
// 	for tx := range uniqueTxsDict {
// 		uniqueTxs = append(uniqueTxs, tx)
// 	}
// 	engine := arbitrator.Start()
// 	arbitrator.Insert(engine, txs, paths, reads, writes, composite)
// 	txs, groups, flags := arbitrator.DetectLegacy(engine, uniqueTxs)
// 	arbitrator.Clear(engine)
// 	return txs, groups, flags
// }

// func prepare(db interfaces.Datastore, height uint64, transitions [] *univalue.Univalue, txs []uint32) (*vmadaptor.EU, *vmadaptor.Config) {
// 	url := concurrenturl.NewConcurrentUrl(db)
// 	url.Import(transitions)
// 	url.Sort()
// 	url.Commit(txs)
// 	api := eu.NewAPI(url)
// 	statedb := eth.NewImplStateDB(api)

// 	config := vmadaptor.NewConfig()
// 	config.Coinbase = &Coinbase
// 	config.BlockNumber = new(big.Int).SetUint64(height)
// 	config.Time = new(big.Int).SetUint64(height)
// 	return vmadaptor.NewEU(config.ChainConfig, *config.VMConfig, statedb, api), config
// }
