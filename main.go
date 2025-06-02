package main

import (
	"fmt"
	"github.com/ironfang-ltd/go-script/evaluator"
	"github.com/ironfang-ltd/go-script/lexer"
	"github.com/ironfang-ltd/go-script/parser"
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

		eval.RegisterFunction("count", func(ctx *evaluator.ExecutionContext, scope *evaluator.Scope, args ...evaluator.Object) (evaluator.Object, error) {
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

		ctx := evaluator.NewExecutionContext(program)
		ctx.RootScope.SetLocal("title", &evaluator.StringValue{Value: "Hello World"})

		hash := evaluator.NewHashValue()
		hash.Set(&evaluator.StringValue{Value: "name"}, &evaluator.StringValue{Value: "Item #1"})
		hash2 := evaluator.NewHashValue()
		hash2.Set(&evaluator.StringValue{Value: "name"}, &evaluator.StringValue{Value: "Item #2"})
		hash3 := evaluator.NewHashValue()
		hash3.Set(&evaluator.StringValue{Value: "name"}, &evaluator.StringValue{Value: "Item #3"})

		ctx.RootScope.SetLocal("items", &evaluator.ArrayValue{Elements: []evaluator.Object{
			hash,
			hash2,
			hash3,
		}})

		result, err := eval.EvaluateString(ctx)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("Output:---\n", result)
		}
	}
}
