package evaluator

type Function func(args ...Object) (Object, error)

type BuiltInFunction struct {
	Fn Function
}

func (bif *BuiltInFunction) Type() ObjectType {
	return BuiltInFunctionObject
}

func (bif *BuiltInFunction) Debug() string {
	return "builtin function"
}
