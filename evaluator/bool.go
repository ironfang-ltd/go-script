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
