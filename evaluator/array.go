package evaluator

type ArrayValue struct {
	Elements []Object
}

func NewArrayValue(elements []Object) *ArrayValue {
	return &ArrayValue{Elements: elements}
}

func (a *ArrayValue) Debug() string {
	return "Array"
}

func (a *ArrayValue) Type() ObjectType {
	return ArrayObject
}
