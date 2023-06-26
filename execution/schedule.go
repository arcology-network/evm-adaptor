package execution

import (
	"crypto/sha256"

	commonlibcommon "github.com/arcology-network/common-lib/common"
	"github.com/arcology-network/evm/core"
)

type Branch struct {
	sequences StandardMessages
}

type Generation struct {
	branches []Branch
}

type Schedule struct {
	generations []Generation
	dict        map[[32]byte][]*StandardMessage
}

func (this *Schedule) NewSchedule() *Schedule {
	return &Schedule{
		generations: []Generation{},
		dict:        make(map[[32]byte][]*StandardMessage),
	}
}

func (this Schedule) callSignature(msg *core.Message) [32]byte {
	return sha256.Sum256(commonlibcommon.Flatten([][]byte{
		msg.To[:],
		msg.Data[:4],
	}))
}

func (this *Schedule) IsLast(txhash [32]byte, msg *core.Message) bool { //
	signature := this.callSignature(msg)
	if stdMsgs, ok := this.dict[signature]; ok {
		return txhash == stdMsgs[len(stdMsgs)-1].TxHash
	}
	return false
}
