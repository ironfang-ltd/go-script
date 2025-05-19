package evaluator

import (
	"github.com/ironfang-ltd/go-script/lexer"
	"github.com/ironfang-ltd/go-script/parser"
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

	if ret.Type() != FunctionObject {
		t.Errorf("got=%v, expected=%v", ret.Type(), FunctionObject)
	}
}

func TestEvaluateScopeVariable(t *testing.T) {
	test := "parent.parent2.value"

	l := lexer.NewScript(test)
	p := parser.New(l)

	program, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}

	e := New(os.Stdout)
	scope := NewScope()

	v := NewFileValue("123", "files/123.png", "test.png", "image/png", 100)

	parent2 := NewHashValue()
	parent2.Set(NewStringValue("value"), v)

	parent := NewHashValue()
	parent.Set(NewStringValue("parent2"), parent2)

	scope.Set("parent", parent)

	ret, err := e.Evaluate(program, scope)
	if err != nil {
		t.Fatal(err)
	}

	if ret.Type() != FileObject {
		t.Errorf("got=%v, expected=%v", ret.Type(), FileObject)
	}
}
