package evaluator

import "fmt"

type IntegerValue struct {
	Value int
}

func NewIntegerValue(val int) *IntegerValue {
	return &IntegerValue{val}
}

func (i *IntegerValue) Debug() string {
	return fmt.Sprintf("%d", i.Value)
}

func (i *IntegerValue) Type() ObjectType {
	return IntegerObject
}

func (i *IntegerValue) HashKey() HashKey {
	return HashKey{Type: i.Type(), Value: uint64(i.Value)}
}
