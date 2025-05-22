package evaluator

import (
	"fmt"
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

func TestEvaluateComplexExample(t *testing.T) {
	test := `
let lenders = object_query( "lender", lead.id );
let result = [];

foreach (lenders as lender) {
	
	let lender = { "id": lender.id, "name": lender.name, "documents": [] };	
	let documents = object_query( "document", lender.id );
	
	foreach (documents as document) {
		append(lender.documents, document);
	}
	
	append(result, lender);
}

return result;`

	l := lexer.NewScript(test)
	p := parser.New(l)

	program, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}

	e := New(os.Stdout)

	e.RegisterFunction("object_query", func(args ...Object) (Object, error) {

		if len(args) != 2 {
			return nil, fmt.Errorf("object_query() expects 2 arguments")
		}

		typeId, ok := args[0].(*StringValue)
		if !ok {
			return nil, fmt.Errorf("object_query() expects a string as the first argument")
		}

		id, ok := args[1].(*StringValue)
		if !ok {
			return nil, fmt.Errorf("object_query() expects a string as the second argument")
		}

		if typeId.Value == "lender" {
			objArr := make([]Object, 0)
			lender1 := NewHashValue()
			lender1.Set(NewStringValue("id"), NewStringValue("1"))
			lender1.Set(NewStringValue("name"), NewStringValue("lender-1"))
			lender2 := NewHashValue()
			lender2.Set(NewStringValue("id"), NewStringValue("2"))
			lender2.Set(NewStringValue("name"), NewStringValue("lender-2"))
			objArr = append(objArr, lender1, lender2)
			return NewArrayValue(objArr), nil
		}

		if typeId.Value == "document" && id.Value == "1" {
			objArr := make([]Object, 0)
			doc := NewHashValue()
			doc.Set(NewStringValue("id"), NewStringValue("doc-1"))
			objArr = append(objArr, doc)
			return NewArrayValue(objArr), nil
		}

		if typeId.Value == "document" && id.Value == "2" {
			objArr := make([]Object, 0)
			doc := NewHashValue()
			doc.Set(NewStringValue("id"), NewStringValue("doc-2"))
			objArr = append(objArr, doc)
			return NewArrayValue(objArr), nil
		}

		return nil, nil
	})

	scope := NewScope()
	leadObj := NewHashValue()
	leadObj.Set(NewStringValue("id"), NewStringValue("123"))
	scope.Set("lead", leadObj)

	result, err := e.Evaluate(program, scope)
	if err != nil {
		t.Fatal(err)
	}

	ret, ok := result.(*ReturnValue)
	if !ok {
		t.Fatalf("expected *ReturnValue, got %s", ret.Type())
	}

	retVal, ok := ret.Value.(*ArrayValue)
	if !ok {
		t.Fatalf("expected *ArrayValue, got %s", ret.Type())
	}

	if len(retVal.Elements) != 2 {
		t.Fatalf("expected: 2, got %d", len(retVal.Elements))
	}

	for _, elem := range retVal.Elements {

		objVal, ok := elem.(*HashValue)
		if !ok {
			t.Fatalf("expected *HashValue, got %s", ret.Type())
		}

		if !objVal.HasKey(NewStringValue("id")) {
			t.Fatalf("expected: id, got %s", ret.Type())
		}

		if !objVal.HasKey(NewStringValue("name")) {
			t.Fatalf("expected: name, got %s", ret.Type())
		}

		if !objVal.HasKey(NewStringValue("documents")) {
			t.Fatalf("expected: document, got %s", ret.Type())
		}
	}
}
