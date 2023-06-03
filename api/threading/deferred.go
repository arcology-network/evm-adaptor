package threading

type DeferredCall struct {
	Hash [32]byte
	dict map[string][][32]byte
}

func (this *DeferredCall) NewDeferredCall() *DeferredCall {
	return &DeferredCall{
		dict: map[string][][32]byte{},
	}
}

func (this *DeferredCall) Add(funcCall string, calls ...[32]byte) {
	for _, call := range calls {
		this.dict[funcCall] = append(this.dict[funcCall], call)
	}
}
