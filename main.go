package main

import (
	"fmt"

	"github.com/ironfang-ltd/ironscript/evaluator"
	"github.com/ironfang-ltd/ironscript/lexer"
	"github.com/ironfang-ltd/ironscript/parser"
)

func main() {

	script := `
fn add(a, b) {
	return a + b;
}

let result = add(5, 10);

print (result);
`

	l := lexer.New(script)
	p := parser.New(l)

	program, err := p.Parse()
	if err != nil {
		fmt.Println(err)
	}

	eval := evaluator.New()

	eval.RegisterFunction("print", func(args ...evaluator.Object) (evaluator.Object, error) {
		for _, arg := range args {
			fmt.Println(arg.Debug())
		}
		return evaluator.Null, nil
	})

	eval.RegisterFunction("add", func(args ...evaluator.Object) (evaluator.Object, error) {
		if len(args) != 2 {
			return evaluator.Null, fmt.Errorf("wrong number of arguments. got=%d, want=2", len(args))
		}

		left, ok := args[0].(*evaluator.IntegerValue)
		if !ok {
			return evaluator.Null, fmt.Errorf("arguments to `add` must be INTEGER, got %s and %s", args[0].Type(), args[1].Type())
		}

		right, ok := args[1].(*evaluator.IntegerValue)
		if !ok {
			return evaluator.Null, fmt.Errorf("arguments to `add` must be INTEGER, got %s and %s", args[0].Type(), args[1].Type())
		}

		return &evaluator.IntegerValue{Value: left.Value + right.Value}, nil
	})

	scope := evaluator.NewScope()
	result, err := eval.Evaluate(program, scope)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(result.Debug())
	}
}
