package evaluator

type BooleanValue struct {
	Value bool
}

func NewBooleanValue(val bool) *BooleanValue {
	return &BooleanValue{val}
}

func (b *BooleanValue) Debug() string {
	if b.Value {
		return "true"
	}

	return "false"
}

func (b *BooleanValue) Type() ObjectType {
	return BooleanObject
}

func (b *BooleanValue) HashKey() HashKey {
	var val uint64
	if b.Value {
		val = 1
	}
	return HashKey{Type: b.Type(), Value: val}
}
