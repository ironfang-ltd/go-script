package evaluator

import (
	"github.com/ironfang-ltd/ironscript/lexer"
	"github.com/ironfang-ltd/ironscript/parser"
	"os"
	"testing"
)

func TestEvaluateReturn(t *testing.T) {
	test := "return 123;"

	l := lexer.NewScript(test)
	p := parser.New(l)

	program, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}

	e := New(os.Stdout)
	scope := NewScope()

	ret, err := e.Evaluate(program, scope)
	if err != nil {
		t.Fatal(err)
	}

	if ret.Type() != ReturnValueObject {
		t.Errorf("got=%v, expected=%v", ret.Type(), ReturnValueObject)
	}

	if ret.Debug() != "123" {
		t.Errorf("got=%v, expected=%v", ret.Debug(), "123")
	}
}

func TestEvaluateFnLiteral(t *testing.T) {
	test := "fn add(x, y) { return x + y; }"

	l := lexer.NewScript(test)
	p := parser.New(l)

	program, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}

	e := New(os.Stdout)
	scope := NewScope()

	ret, err := e.Evaluate(program, scope)
	if err != nil {
		t.Fatal(err)
	}

	if ret.Type() != NullObject {
		t.Errorf("got=%v, expected=%v", ret.Type(), NullObject)
	}
}
