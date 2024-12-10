package main

import (
	"fmt"
	"github.com/ironfang-ltd/ironscript/evaluator"

	"github.com/ironfang-ltd/ironscript/lexer"
	"github.com/ironfang-ltd/ironscript/parser"
)

func main() {

	t := `
<h1>{% title %}</h1>
<h2>Items ({% print(count(items)) %}): </h2>
{% if (count(items) > 0) { %}
<p>There are items</p>
{% } %}
<ul>
	{% foreach (items as item) { %}
	<li>{% print(item.name) %}</li>
	{% } %}
</ul>
`

	l := lexer.NewTemplate(t)
	p := parser.New(l)

	program, err := p.Parse()
	if err != nil {
		fmt.Println("Parse Errors:")
		fmt.Println(err)
	} else {

		eval := evaluator.New()

		eval.RegisterFunction("count", func(args ...evaluator.Object) (evaluator.Object, error) {
			if len(args) != 1 {
				return evaluator.Null, fmt.Errorf("wrong number of arguments. got=%d, want=1", len(args))
			}

			switch arg := args[0].(type) {
			case *evaluator.ArrayValue:
				return &evaluator.IntegerValue{Value: len(arg.Elements)}, nil
			default:
				return evaluator.Null, fmt.Errorf("argument to `count` not supported, got %s", args[0].Type())
			}
		})

		scope := evaluator.NewScope()
		scope.Set("title", &evaluator.StringValue{Value: "Hello World"})

		hash := evaluator.NewHashValue()
		hash.Set(&evaluator.StringValue{Value: "name"}, &evaluator.StringValue{Value: "Item #1"})
		hash2 := evaluator.NewHashValue()
		hash2.Set(&evaluator.StringValue{Value: "name"}, &evaluator.StringValue{Value: "Item #2"})
		hash3 := evaluator.NewHashValue()
		hash3.Set(&evaluator.StringValue{Value: "name"}, &evaluator.StringValue{Value: "Item #3"})

		scope.Set("items", &evaluator.ArrayValue{Elements: []evaluator.Object{
			hash,
			hash2,
			hash3,
		}})

		result, err := eval.Evaluate(program, scope)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("Result: %s\n", result.Debug())
			fmt.Println("Output:---\n", eval.GetOutput())
		}
	}

	/*script := `
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
		}*/
}
