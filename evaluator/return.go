package evaluator

type ReturnValue struct {
	Value Object
}

func (r *ReturnValue) Debug() string {
	return r.Value.Debug()
}

func (r *ReturnValue) Type() ObjectType {
	return ReturnValueObject
}
