package evaluator

import "github.com/ironfang-ltd/go-script/parser"

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

type FunctionValue struct {
	Parameters []*parser.Identifier
	Body       *parser.BlockStatement
	Scope      *Scope
}

func (f *FunctionValue) Debug() string {
	return "Function"
}

func (f *FunctionValue) Type() ObjectType {
	return FunctionObject
}
