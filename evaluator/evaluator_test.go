package evaluator

import (
	"github.com/ironfang-ltd/ironscript/lexer"
	"github.com/ironfang-ltd/ironscript/parser"
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

	e := New()
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
