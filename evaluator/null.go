package evaluator

type NullValue struct{}

func (n *NullValue) Debug() string {
	return "null"
}

func (n *NullValue) Type() ObjectType {
	return NullObject
}
