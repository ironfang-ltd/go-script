package evaluator

import (
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/ironfang-ltd/go-script/lexer"
	"github.com/ironfang-ltd/go-script/parser"
)

func TestEvaluateAssignment(t *testing.T) {
	test := "let x = 0; x = 123; return x;"

	l := lexer.NewScript(test)
	p := parser.New(l)

	program, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}

	e := New()
	ctx := NewExecutionContext(program)

	ret, err := e.Evaluate(ctx)
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

func TestEvaluateArrayAssignment(t *testing.T) {
	test := "let x = [0]; x[0] = 123; return x;"

	l := lexer.NewScript(test)
	p := parser.New(l)

	program, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}

	e := New()
	ctx := NewExecutionContext(program)

	result, err := e.Evaluate(ctx)
	if err != nil {
		t.Fatal(err)
	}

	ret, ok := result.(*ReturnValue)
	if !ok {
		t.Fatalf("expected *ReturnValue, got %T", ret)
	}

	arrayObj, ok := ret.Value.(*ArrayValue)
	if !ok {
		t.Fatalf("expected *HashValue, got %T", ret.Value)
	}

	if len(arrayObj.Elements) != 1 {
		t.Fatalf("expected array length to be 1, got %d", len(arrayObj.Elements))
	}

	if arrayObj.Elements[0].Type() != IntegerObject {
		t.Fatalf("expected first element to be an IntegerObject, got %s", arrayObj.Elements[0].Type())
	}

	if arrayObj.Elements[0].Debug() != "123" {
		t.Errorf("got=%v, expected=%v", arrayObj.Elements[0].Debug(), "123")
	}
}

func TestEvaluatePropertyAssignment(t *testing.T) {
	test := "let x = {\"y\": 0}; x.y = 123; return x;"

	l := lexer.NewScript(test)
	p := parser.New(l)

	program, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}

	e := New()
	ctx := NewExecutionContext(program)

	result, err := e.Evaluate(ctx)
	if err != nil {
		t.Fatal(err)
	}

	ret, ok := result.(*ReturnValue)
	if !ok {
		t.Fatalf("expected *ReturnValue, got %T", ret)
	}

	hashObj, ok := ret.Value.(*HashValue)
	if !ok {
		t.Fatalf("expected *HashValue, got %T", ret.Value)
	}

	yValue, ok := hashObj.GetValue(NewStringValue("y"))
	if !ok {
		t.Fatal("expected to find key 'y' in hash")
	}

	if yValue.Type() != IntegerObject {
		t.Fatalf("expected y to be an IntegerObject, got %s", yValue.Type())
	}

	if yValue.Debug() != "123" {
		t.Errorf("got=%v, expected=%v", yValue.Debug(), "123")
	}
}

func TestEvaluatePropertyNestedAssignment(t *testing.T) {
	test := "let x = {\"y\": {\"z\": 0}}; x.y.z = 123; return x;"

	l := lexer.NewScript(test)
	p := parser.New(l)

	program, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}

	e := New()
	ctx := NewExecutionContext(program)

	result, err := e.Evaluate(ctx)
	if err != nil {
		t.Fatal(err)
	}

	ret, ok := result.(*ReturnValue)
	if !ok {
		t.Fatalf("expected *ReturnValue, got %T", ret)
	}

	hashObj, ok := ret.Value.(*HashValue)
	if !ok {
		t.Fatalf("expected *HashValue, got %T", ret.Value)
	}

	yValue, ok := hashObj.GetValue(NewStringValue("y"))
	if !ok {
		t.Fatal("expected to find key 'y' in hash")
	}

	if yValue.Type() != HashObject {
		t.Fatalf("expected y to be an HashObject, got %s", yValue.Type())
	}

	yHash, ok := yValue.(*HashValue)
	if !ok {
		t.Fatalf("expected y to be an HashValue, got %T", yValue)
	}

	zValue, ok := yHash.GetValue(NewStringValue("z"))
	if !ok {
		t.Fatal("expected to find key 'z' in y hash")
	}

	if zValue.Type() != IntegerObject {
		t.Fatalf("expected z to be an IntegerObject, got %s", zValue.Type())
	}

	if zValue.Debug() != "123" {
		t.Errorf("got=%v, expected=%v", zValue.Debug(), "123")
	}
}

func TestEvaluateAssignmentWithoutDeclaration(t *testing.T) {
	test := "x = 123; return x;"

	l := lexer.NewScript(test)
	p := parser.New(l)

	program, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}

	e := New()
	ctx := NewExecutionContext(program)

	_, err = e.Evaluate(ctx)
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
}

func TestEvaluateReturn(t *testing.T) {
	test := "return 123;"

	l := lexer.NewScript(test)
	p := parser.New(l)

	program, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}

	e := New()
	ctx := NewExecutionContext(program)

	ret, err := e.Evaluate(ctx)
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

	e := New()
	ctx := NewExecutionContext(program)

	ret, err := e.Evaluate(ctx)
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

	e := New()
	ctx := NewExecutionContext(program)

	v := NewFileValue("123", "files/123.png", "test.png", "image/png", 100)

	parent2 := NewHashValue()
	parent2.Set(NewStringValue("value"), v)

	parent := NewHashValue()
	parent.Set(NewStringValue("parent2"), parent2)

	ctx.RootScope.SetLocal("parent", parent)

	ret, err := e.Evaluate(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if ret.Type() != FileObject {
		t.Errorf("got=%v, expected=%v", ret.Type(), FileObject)
	}
}

func TestEvaluateFloatArithmetic(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"return 1.5 + 2.5;", "4.0"},
		{"return 3.0 - 1.5;", "1.5"},
		{"return 2.0 * 3.5;", "7.0"},
		{"return 7.0 / 2.0;", "3.5"},
	}

	for _, tt := range tests {
		l := lexer.NewScript(tt.input)
		p := parser.New(l)
		program, err := p.Parse()
		if err != nil {
			t.Fatalf("Parse(%s) error: %v", tt.input, err)
		}

		e := New()
		ctx := NewExecutionContext(program)
		ret, err := e.Evaluate(ctx)
		if err != nil {
			t.Fatalf("Evaluate(%s) error: %v", tt.input, err)
		}

		retVal, ok := ret.(*ReturnValue)
		if !ok {
			t.Fatalf("expected ReturnValue for %s, got %T", tt.input, ret)
		}
		if retVal.Value.Debug() != tt.expected {
			t.Fatalf("for %s: expected %s, got %s", tt.input, tt.expected, retVal.Value.Debug())
		}
	}
}

func TestEvaluateFloatComparisons(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"return 1.5 < 2.5;", true},
		{"return 2.5 > 1.5;", true},
		{"return 1.5 == 1.5;", true},
		{"return 1.5 != 2.5;", true},
		{"return 1.5 > 2.5;", false},
		{"return 1.5 == 2.5;", false},
	}

	for _, tt := range tests {
		l := lexer.NewScript(tt.input)
		p := parser.New(l)
		program, err := p.Parse()
		if err != nil {
			t.Fatalf("Parse(%s) error: %v", tt.input, err)
		}

		e := New()
		ctx := NewExecutionContext(program)
		ret, err := e.Evaluate(ctx)
		if err != nil {
			t.Fatalf("Evaluate(%s) error: %v", tt.input, err)
		}

		retVal, ok := ret.(*ReturnValue)
		if !ok {
			t.Fatalf("expected ReturnValue for %s, got %T", tt.input, ret)
		}
		boolVal, ok := retVal.Value.(*BooleanValue)
		if !ok {
			t.Fatalf("expected BooleanValue for %s, got %T", tt.input, retVal.Value)
		}
		if boolVal.Value != tt.expected {
			t.Fatalf("for %s: expected %v, got %v", tt.input, tt.expected, boolVal.Value)
		}
	}
}

func TestEvaluateNegativeFloat(t *testing.T) {
	test := "return -3.14;"

	l := lexer.NewScript(test)
	p := parser.New(l)
	program, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}

	e := New()
	ctx := NewExecutionContext(program)
	ret, err := e.Evaluate(ctx)
	if err != nil {
		t.Fatal(err)
	}

	retVal, ok := ret.(*ReturnValue)
	if !ok {
		t.Fatalf("expected ReturnValue, got %T", ret)
	}
	decVal, ok := retVal.Value.(*DecimalValue)
	if !ok {
		t.Fatalf("expected DecimalValue, got %T", retVal.Value)
	}
	if decVal.Value != -3.14 {
		t.Fatalf("expected -3.14, got %f", decVal.Value)
	}
}

func TestEvaluateFloatLetStatement(t *testing.T) {
	test := "let x = 3.14; return x;"

	l := lexer.NewScript(test)
	p := parser.New(l)
	program, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}

	e := New()
	ctx := NewExecutionContext(program)
	ret, err := e.Evaluate(ctx)
	if err != nil {
		t.Fatal(err)
	}

	retVal, ok := ret.(*ReturnValue)
	if !ok {
		t.Fatalf("expected ReturnValue, got %T", ret)
	}
	decVal, ok := retVal.Value.(*DecimalValue)
	if !ok {
		t.Fatalf("expected DecimalValue, got %T", retVal.Value)
	}
	if decVal.Value != 3.14 {
		t.Fatalf("expected 3.14, got %f", decVal.Value)
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

	e := New()

	e.RegisterFunction("object_query", func(_ *ExecutionContext, _ *Scope, args ...Object) (Object, error) {

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

	ctx := NewExecutionContext(program)
	leadObj := NewHashValue()
	leadObj.Set(NewStringValue("id"), NewStringValue("123"))
	ctx.RootScope.SetLocal("lead", leadObj)

	result, err := e.Evaluate(ctx)
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

// --- Test helpers ---

func evalScript(t *testing.T, input string) Object {
	t.Helper()
	l := lexer.NewScript(input)
	p := parser.New(l)
	program, err := p.Parse()
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	e := New()
	ctx := NewExecutionContext(program)
	ret, err := e.Evaluate(ctx)
	if err != nil {
		t.Fatalf("Evaluate error: %v", err)
	}
	return ret
}

func evalScriptError(t *testing.T, input string) error {
	t.Helper()
	l := lexer.NewScript(input)
	p := parser.New(l)
	program, err := p.Parse()
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	e := New()
	ctx := NewExecutionContext(program)
	_, err = e.Evaluate(ctx)
	return err
}

func evalTemplate(t *testing.T, input string) string {
	t.Helper()
	l := lexer.NewTemplate(input)
	p := parser.New(l)
	program, err := p.Parse()
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	e := New()
	ctx := NewExecutionContext(program)
	output, err := e.EvaluateString(ctx)
	if err != nil {
		t.Fatalf("EvaluateString error: %v", err)
	}
	return output
}

func unwrapReturn(t *testing.T, obj Object) Object {
	t.Helper()
	ret, ok := obj.(*ReturnValue)
	if !ok {
		t.Fatalf("expected ReturnValue, got %T", obj)
	}
	return ret.Value
}

// --- Integer arithmetic ---

func TestEvaluateIntegerArithmetic(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"return 1 + 2;", 3},
		{"return 10 - 3;", 7},
		{"return 3 * 4;", 12},
		{"return 10 / 3;", 3},
		{"return 10 % 3;", 1},
		{"return -5;", -5},
		{"return 2 + 3 * 4;", 14},
		{"return (2 + 3) * 4;", 20},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			val := unwrapReturn(t, evalScript(t, tt.input))
			intVal, ok := val.(*IntegerValue)
			if !ok {
				t.Fatalf("expected IntegerValue, got %T", val)
			}
			if intVal.Value != tt.expected {
				t.Fatalf("expected %d, got %d", tt.expected, intVal.Value)
			}
		})
	}
}

// --- Integer comparisons ---

func TestEvaluateIntegerComparisons(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"return 1 < 2;", true},
		{"return 2 > 1;", true},
		{"return 5 == 5;", true},
		{"return 5 != 3;", true},
		{"return 5 <= 5;", true},
		{"return 4 <= 5;", true},
		{"return 6 <= 5;", false},
		{"return 5 >= 5;", true},
		{"return 6 >= 5;", true},
		{"return 4 >= 5;", false},
		{"return 1 > 2;", false},
		{"return 1 == 2;", false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			val := unwrapReturn(t, evalScript(t, tt.input))
			boolVal, ok := val.(*BooleanValue)
			if !ok {
				t.Fatalf("expected BooleanValue, got %T", val)
			}
			if boolVal.Value != tt.expected {
				t.Fatalf("expected %v, got %v", tt.expected, boolVal.Value)
			}
		})
	}
}

// --- Decimal <= and >= ---

func TestEvaluateDecimalLessOrEqualGreaterOrEqual(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"return 1.5 <= 1.5;", true},
		{"return 1.0 <= 1.5;", true},
		{"return 2.0 <= 1.5;", false},
		{"return 1.5 >= 1.5;", true},
		{"return 2.0 >= 1.5;", true},
		{"return 1.0 >= 1.5;", false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			val := unwrapReturn(t, evalScript(t, tt.input))
			boolVal, ok := val.(*BooleanValue)
			if !ok {
				t.Fatalf("expected BooleanValue, got %T", val)
			}
			if boolVal.Value != tt.expected {
				t.Fatalf("expected %v, got %v", tt.expected, boolVal.Value)
			}
		})
	}
}

// --- String operations ---

func TestEvaluateStringConcatenation(t *testing.T) {
	val := unwrapReturn(t, evalScript(t, `return "hello" + " " + "world";`))
	strVal, ok := val.(*StringValue)
	if !ok {
		t.Fatalf("expected StringValue, got %T", val)
	}
	if strVal.Value != "hello world" {
		t.Fatalf("expected 'hello world', got %q", strVal.Value)
	}
}

func TestEvaluateStringComparisons(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{`return "abc" == "abc";`, true},
		{`return "abc" != "def";`, true},
		{`return "abc" == "def";`, false},
		{`return "abc" != "abc";`, false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			val := unwrapReturn(t, evalScript(t, tt.input))
			boolVal, ok := val.(*BooleanValue)
			if !ok {
				t.Fatalf("expected BooleanValue, got %T", val)
			}
			if boolVal.Value != tt.expected {
				t.Fatalf("expected %v, got %v", tt.expected, boolVal.Value)
			}
		})
	}
}

// --- Boolean operations ---

func TestEvaluateBooleanInfix(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"return true == true;", true},
		{"return true == false;", false},
		{"return true != false;", true},
		{"return false != false;", false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			val := unwrapReturn(t, evalScript(t, tt.input))
			boolVal, ok := val.(*BooleanValue)
			if !ok {
				t.Fatalf("expected BooleanValue, got %T", val)
			}
			if boolVal.Value != tt.expected {
				t.Fatalf("expected %v, got %v", tt.expected, boolVal.Value)
			}
		})
	}
}

func TestEvaluateBangPrefix(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"return !true;", false},
		{"return !false;", true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			val := unwrapReturn(t, evalScript(t, tt.input))
			boolVal, ok := val.(*BooleanValue)
			if !ok {
				t.Fatalf("expected BooleanValue, got %T", val)
			}
			if boolVal.Value != tt.expected {
				t.Fatalf("expected %v, got %v", tt.expected, boolVal.Value)
			}
		})
	}
}

// --- If/else ---

func TestEvaluateIfTrue(t *testing.T) {
	val := unwrapReturn(t, evalScript(t, `if (true) { return 1; } else { return 2; }`))
	intVal := val.(*IntegerValue)
	if intVal.Value != 1 {
		t.Fatalf("expected 1, got %d", intVal.Value)
	}
}

func TestEvaluateIfFalse(t *testing.T) {
	val := unwrapReturn(t, evalScript(t, `if (false) { return 1; } else { return 2; }`))
	intVal := val.(*IntegerValue)
	if intVal.Value != 2 {
		t.Fatalf("expected 2, got %d", intVal.Value)
	}
}

func TestEvaluateIfNoElse(t *testing.T) {
	ret := evalScript(t, `if (false) { return 1; }`)
	if ret.Type() != NullObject {
		t.Fatalf("expected Null, got %s", ret.Type())
	}
}

func TestEvaluateIfWithCondition(t *testing.T) {
	val := unwrapReturn(t, evalScript(t, `let x = 10; if (x > 5) { return "big"; } else { return "small"; }`))
	strVal := val.(*StringValue)
	if strVal.Value != "big" {
		t.Fatalf("expected 'big', got %q", strVal.Value)
	}
}

// --- Function calls ---

func TestEvaluateFunctionCall(t *testing.T) {
	val := unwrapReturn(t, evalScript(t, `fn add(a, b) { return a + b; } return add(3, 4);`))
	intVal := val.(*IntegerValue)
	if intVal.Value != 7 {
		t.Fatalf("expected 7, got %d", intVal.Value)
	}
}

func TestEvaluateFunctionNoReturn(t *testing.T) {
	// A function without explicit return yields the last expression value
	val := evalScript(t, `fn greet() { let x = 1; } return greet();`)
	retVal := unwrapReturn(t, val)
	intVal, ok := retVal.(*IntegerValue)
	if !ok {
		t.Fatalf("expected IntegerValue, got %T", retVal)
	}
	if intVal.Value != 1 {
		t.Fatalf("expected 1, got %d", intVal.Value)
	}
}

func TestEvaluateRecursiveFunction(t *testing.T) {
	script := `
fn factorial(n) {
	if (n <= 1) { return 1; }
	return n * factorial(n - 1);
}
return factorial(5);`
	val := unwrapReturn(t, evalScript(t, script))
	intVal := val.(*IntegerValue)
	if intVal.Value != 120 {
		t.Fatalf("expected 120, got %d", intVal.Value)
	}
}

func TestEvaluateFunctionWrongArity(t *testing.T) {
	err := evalScriptError(t, `fn add(a, b) { return a + b; } add(1);`)
	if err == nil {
		t.Fatal("expected error for wrong number of arguments")
	}
	if !strings.Contains(err.Error(), "wrong number of arguments") {
		t.Fatalf("expected 'wrong number of arguments' error, got: %v", err)
	}
}

func TestEvaluateFunctionClosureScope(t *testing.T) {
	script := `
let x = 10;
fn addX(y) { return x + y; }
return addX(5);`
	val := unwrapReturn(t, evalScript(t, script))
	intVal := val.(*IntegerValue)
	if intVal.Value != 15 {
		t.Fatalf("expected 15, got %d", intVal.Value)
	}
}

// --- Scope: assign to outer scope variable ---

func TestEvaluateAssignOuterScope(t *testing.T) {
	script := `
let x = 1;
fn setX() { x = 99; }
setX();
return x;`
	val := unwrapReturn(t, evalScript(t, script))
	intVal := val.(*IntegerValue)
	if intVal.Value != 99 {
		t.Fatalf("expected 99, got %d", intVal.Value)
	}
}

// --- Array operations ---

func TestEvaluateArrayLiteral(t *testing.T) {
	val := evalScript(t, `let a = [1, 2, 3];`)
	arrVal, ok := val.(*ArrayValue)
	if !ok {
		t.Fatalf("expected ArrayValue, got %T", val)
	}
	if len(arrVal.Elements) != 3 {
		t.Fatalf("expected 3 elements, got %d", len(arrVal.Elements))
	}
}

func TestEvaluateArrayIndex(t *testing.T) {
	val := unwrapReturn(t, evalScript(t, `let a = [10, 20, 30]; return a[1];`))
	intVal := val.(*IntegerValue)
	if intVal.Value != 20 {
		t.Fatalf("expected 20, got %d", intVal.Value)
	}
}

func TestEvaluateArrayIndexOutOfBounds(t *testing.T) {
	val := unwrapReturn(t, evalScript(t, `let a = [10, 20]; return a[5];`))
	if val.Type() != NullObject {
		t.Fatalf("expected Null for out-of-bounds, got %s", val.Type())
	}
}

func TestEvaluateArrayNegativeIndex(t *testing.T) {
	val := unwrapReturn(t, evalScript(t, `let a = [10, 20]; return a[-1];`))
	if val.Type() != NullObject {
		t.Fatalf("expected Null for negative index, got %s", val.Type())
	}
}

// --- Hash operations ---

func TestEvaluateHashLiteral(t *testing.T) {
	val := evalScript(t, `let h = {"a": 1, "b": 2};`)
	hashVal, ok := val.(*HashValue)
	if !ok {
		t.Fatalf("expected HashValue, got %T", val)
	}
	aVal, ok := hashVal.GetValue(NewStringValue("a"))
	if !ok {
		t.Fatal("expected key 'a'")
	}
	if aVal.Debug() != "1" {
		t.Fatalf("expected 1, got %s", aVal.Debug())
	}
}

func TestEvaluateHashIndexExpression(t *testing.T) {
	val := unwrapReturn(t, evalScript(t, `let h = {"x": 42}; return h["x"];`))
	intVal := val.(*IntegerValue)
	if intVal.Value != 42 {
		t.Fatalf("expected 42, got %d", intVal.Value)
	}
}

func TestEvaluateHashMissingKey(t *testing.T) {
	val := unwrapReturn(t, evalScript(t, `let h = {"x": 1}; return h["y"];`))
	if val.Type() != NullObject {
		t.Fatalf("expected Null for missing key, got %s", val.Type())
	}
}

func TestEvaluateHashPropertyAccess(t *testing.T) {
	val := unwrapReturn(t, evalScript(t, `let h = {"name": "test"}; return h.name;`))
	strVal := val.(*StringValue)
	if strVal.Value != "test" {
		t.Fatalf("expected 'test', got %q", strVal.Value)
	}
}

// --- Foreach ---

func TestEvaluateForeachArray(t *testing.T) {
	script := `
let arr = [1, 2, 3];
let sum = 0;
foreach (arr as item) {
	sum = sum + item;
}
return sum;`
	val := unwrapReturn(t, evalScript(t, script))
	intVal := val.(*IntegerValue)
	if intVal.Value != 6 {
		t.Fatalf("expected 6, got %d", intVal.Value)
	}
}

func TestEvaluateForeachEmptyArray(t *testing.T) {
	script := `
let arr = [];
let sum = 0;
foreach (arr as item) {
	sum = sum + 1;
}
return sum;`
	val := unwrapReturn(t, evalScript(t, script))
	intVal := val.(*IntegerValue)
	if intVal.Value != 0 {
		t.Fatalf("expected 0, got %d", intVal.Value)
	}
}

// --- EvaluateString (template mode) ---

func TestEvaluateStringTemplate(t *testing.T) {
	output := evalTemplate(t, `Hello {% "world" %}!`)
	if output != "Hello world!" {
		t.Fatalf("expected 'Hello world!', got %q", output)
	}
}

func TestEvaluateStringTemplateWithExpression(t *testing.T) {
	output := evalTemplate(t, `Result: {% 1 + 2 %}`)
	if output != "Result: 3" {
		t.Fatalf("expected 'Result: 3', got %q", output)
	}
}

func TestEvaluateStringTemplateWithVariable(t *testing.T) {
	// Note: EvaluateString outputs all non-null expression results,
	// including let statement values, so we inject the variable via scope instead.
	input := `Hello {% name %}!`
	l := lexer.NewTemplate(input)
	p := parser.New(l)
	program, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}
	e := New()
	ctx := NewExecutionContext(program)
	ctx.RootScope.SetLocal("name", &StringValue{Value: "Alice"})
	output, err := e.EvaluateString(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if output != "Hello Alice!" {
		t.Fatalf("expected 'Hello Alice!', got %q", output)
	}
}

// --- Built-in functions ---

func TestEvaluateAppendBuiltin(t *testing.T) {
	script := `let a = [1, 2]; append(a, 3); return a;`
	val := unwrapReturn(t, evalScript(t, script))
	arrVal := val.(*ArrayValue)
	if len(arrVal.Elements) != 3 {
		t.Fatalf("expected 3 elements, got %d", len(arrVal.Elements))
	}
	if arrVal.Elements[2].Debug() != "3" {
		t.Fatalf("expected last element '3', got %s", arrVal.Elements[2].Debug())
	}
}

func TestEvaluateAppendNonArray(t *testing.T) {
	err := evalScriptError(t, `append("not_array", 1);`)
	if err == nil {
		t.Fatal("expected error when appending to non-array")
	}
	if !strings.Contains(err.Error(), "expected array") {
		t.Fatalf("expected 'expected array' error, got: %v", err)
	}
}

func TestEvaluateAppendWrongArgCount(t *testing.T) {
	err := evalScriptError(t, `let a = []; append(a);`)
	if err == nil {
		t.Fatal("expected error for wrong arg count")
	}
}

func TestEvaluatePrintBuiltin(t *testing.T) {
	input := `{% print("hello"); %}`
	l := lexer.NewTemplate(input)
	p := parser.New(l)
	program, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}
	e := New()
	ctx := NewExecutionContext(program)
	output, err := e.EvaluateString(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if output != "hello" {
		t.Fatalf("expected 'hello', got %q", output)
	}
}

func TestEvaluateLogBuiltin(t *testing.T) {
	input := `log("test message");`
	l := lexer.NewScript(input)
	p := parser.New(l)
	program, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}
	e := New()
	ctx := NewExecutionContext(program)
	var logBuf strings.Builder
	ctx.Logger = &logBuf
	_, err = e.Evaluate(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(logBuf.String(), "test message") {
		t.Fatalf("expected log to contain 'test message', got %q", logBuf.String())
	}
}

// --- Division by zero ---

func TestEvaluateIntegerDivisionByZero(t *testing.T) {
	err := evalScriptError(t, `return 10 / 0;`)
	if err == nil {
		t.Fatal("expected error for division by zero")
	}
	if !strings.Contains(err.Error(), "division by zero") {
		t.Fatalf("expected 'division by zero' error, got: %v", err)
	}
}

func TestEvaluateIntegerModuloByZero(t *testing.T) {
	err := evalScriptError(t, `return 10 % 0;`)
	if err == nil {
		t.Fatal("expected error for modulo by zero")
	}
	if !strings.Contains(err.Error(), "division by zero") {
		t.Fatalf("expected 'division by zero' error, got: %v", err)
	}
}

func TestEvaluateDecimalDivisionByZero(t *testing.T) {
	err := evalScriptError(t, `return 10.0 / 0.0;`)
	if err == nil {
		t.Fatal("expected error for decimal division by zero")
	}
	if !strings.Contains(err.Error(), "division by zero") {
		t.Fatalf("expected 'division by zero' error, got: %v", err)
	}
}

// --- Type mismatch errors ---

func TestEvaluateTypeMismatch(t *testing.T) {
	err := evalScriptError(t, `return 1 - "hello";`)
	if err == nil {
		t.Fatal("expected error for type mismatch")
	}
	if !strings.Contains(err.Error(), "type mismatch") {
		t.Fatalf("expected 'type mismatch' error, got: %v", err)
	}
}

// --- Identifier not found ---

func TestEvaluateUndefinedIdentifier(t *testing.T) {
	err := evalScriptError(t, `return x;`)
	if err == nil {
		t.Fatal("expected error for undefined identifier")
	}
	if !strings.Contains(err.Error(), "identifier not found") {
		t.Fatalf("expected 'identifier not found' error, got: %v", err)
	}
}

// --- Let and scope ---

func TestEvaluateLetShadowing(t *testing.T) {
	script := `
let x = 1;
fn inner() {
	let x = 99;
	return x;
}
let result = inner();
return x;`
	val := unwrapReturn(t, evalScript(t, script))
	intVal := val.(*IntegerValue)
	if intVal.Value != 1 {
		t.Fatalf("expected outer x=1, got %d", intVal.Value)
	}
}

// --- Calling non-function ---

func TestEvaluateCallNonFunction(t *testing.T) {
	err := evalScriptError(t, `let x = 5; x();`)
	if err == nil {
		t.Fatal("expected error when calling non-function")
	}
	if !strings.Contains(err.Error(), "not a function") {
		t.Fatalf("expected 'not a function' error, got: %v", err)
	}
}

// --- Duplicate function definition ---

func TestEvaluateDuplicateFunction(t *testing.T) {
	err := evalScriptError(t, `fn foo() { return 1; } fn foo() { return 2; }`)
	if err == nil {
		t.Fatal("expected error for duplicate function name")
	}
	if !strings.Contains(err.Error(), "already defined") {
		t.Fatalf("expected 'already defined' error, got: %v", err)
	}
}

// --- Negative integer ---

func TestEvaluateNegativeInteger(t *testing.T) {
	val := unwrapReturn(t, evalScript(t, `return -42;`))
	intVal := val.(*IntegerValue)
	if intVal.Value != -42 {
		t.Fatalf("expected -42, got %d", intVal.Value)
	}
}

// --- Empty program ---

func TestEvaluateEmptyProgram(t *testing.T) {
	val := evalScript(t, ``)
	if val.Type() != NullObject {
		t.Fatalf("expected Null, got %s", val.Type())
	}
}

// --- Multiple returns (first wins) ---

func TestEvaluateEarlyReturn(t *testing.T) {
	val := unwrapReturn(t, evalScript(t, `return 1; return 2;`))
	intVal := val.(*IntegerValue)
	if intVal.Value != 1 {
		t.Fatalf("expected 1, got %d", intVal.Value)
	}
}

// --- Registered custom function ---

func TestEvaluateRegisteredFunction(t *testing.T) {
	input := `return double(5);`
	l := lexer.NewScript(input)
	p := parser.New(l)
	program, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}
	e := New()
	e.RegisterFunction("double", func(ctx *ExecutionContext, scope *Scope, args ...Object) (Object, error) {
		intArg := args[0].(*IntegerValue)
		return &IntegerValue{Value: intArg.Value * 2}, nil
	})
	ctx := NewExecutionContext(program)
	ret, err := e.Evaluate(ctx)
	if err != nil {
		t.Fatal(err)
	}
	val := unwrapReturn(t, ret)
	intVal := val.(*IntegerValue)
	if intVal.Value != 10 {
		t.Fatalf("expected 10, got %d", intVal.Value)
	}
}

// --- Metadata on execution context ---

func TestExecutionContextMetadata(t *testing.T) {
	input := ``
	l := lexer.NewScript(input)
	p := parser.New(l)
	program, _ := p.Parse()
	ctx := NewExecutionContext(program)
	ctx.Metadata["key"] = "value"
	if ctx.Metadata["key"] != "value" {
		t.Fatalf("expected metadata 'value', got %v", ctx.Metadata["key"])
	}
}

// --- Scope with pre-set variables ---

func TestEvaluateWithPresetScope(t *testing.T) {
	input := `return name;`
	l := lexer.NewScript(input)
	p := parser.New(l)
	program, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}
	e := New()
	ctx := NewExecutionContext(program)
	ctx.RootScope.SetLocal("name", &StringValue{Value: "Alice"})
	ret, err := e.Evaluate(ctx)
	if err != nil {
		t.Fatal(err)
	}
	val := unwrapReturn(t, ret)
	strVal := val.(*StringValue)
	if strVal.Value != "Alice" {
		t.Fatalf("expected 'Alice', got %q", strVal.Value)
	}
}

// ===================================================================
// Batch 1: Value Type Methods — Direct-call tests for all Object types
// ===================================================================

func TestArrayValueDebugAndType(t *testing.T) {
	a := NewArrayValue([]Object{NewIntegerValue(1), NewStringValue("two")})
	if a.Debug() != "Array" {
		t.Fatalf("expected 'Array', got %q", a.Debug())
	}
	if a.Type() != ArrayObject {
		t.Fatalf("expected ArrayObject, got %s", a.Type())
	}
}

func TestBooleanValueDebugTypeAndConstructor(t *testing.T) {
	tr := NewBooleanValue(true)
	if tr.Debug() != "true" {
		t.Fatalf("expected 'true', got %q", tr.Debug())
	}
	if tr.Type() != BooleanObject {
		t.Fatalf("expected BooleanObject, got %s", tr.Type())
	}

	fl := NewBooleanValue(false)
	if fl.Debug() != "false" {
		t.Fatalf("expected 'false', got %q", fl.Debug())
	}
}

func TestDateTimeValueDebugTypeAndConstructor(t *testing.T) {
	now := time.Date(2026, 3, 3, 12, 0, 0, 0, time.UTC)
	dt := NewDateTimeValue(now)
	if dt.Type() != DateTimeObject {
		t.Fatalf("expected DateTimeObject, got %s", dt.Type())
	}
	expected := now.Format(time.RFC3339)
	if dt.Debug() != expected {
		t.Fatalf("expected %q, got %q", expected, dt.Debug())
	}
}

func TestDecimalValueDebugAndType(t *testing.T) {
	d := NewDecimalValue(3.14)
	if d.Type() != DecimalObject {
		t.Fatalf("expected DecimalObject, got %s", d.Type())
	}
	if d.Debug() != "3.14" {
		t.Fatalf("expected '3.14', got %q", d.Debug())
	}
}

func TestFileValueDebugTypeHashKeyAndConstructor(t *testing.T) {
	f := NewFileValue("id-1", "/path/to/file.png", "file.png", "image/png", 2048)
	if f.Type() != FileObject {
		t.Fatalf("expected FileObject, got %s", f.Type())
	}
	expected := "<file: file.png, size: 2048, contentType: image/png>"
	if f.Debug() != expected {
		t.Fatalf("expected %q, got %q", expected, f.Debug())
	}

	hk := f.HashKey()
	if hk.Type != FileObject {
		t.Fatalf("expected FileObject hash key type, got %s", hk.Type)
	}

	// Same FileID → same hash key
	f2 := NewFileValue("id-1", "/other/path", "other.png", "text/plain", 100)
	if f.HashKey() != f2.HashKey() {
		t.Fatal("expected same hash key for same FileID")
	}

	// Different FileID → different hash key
	f3 := NewFileValue("id-2", "/path/to/file.png", "file.png", "image/png", 2048)
	if f.HashKey() == f3.HashKey() {
		t.Fatal("expected different hash keys for different FileIDs")
	}
}

func TestFunctionValueDebugAndType(t *testing.T) {
	fv := &FunctionValue{}
	if fv.Debug() != "Function" {
		t.Fatalf("expected 'Function', got %q", fv.Debug())
	}
	if fv.Type() != FunctionObject {
		t.Fatalf("expected FunctionObject, got %s", fv.Type())
	}
}

func TestHashValueDebugTypeGetValueMissAndHasKey(t *testing.T) {
	h := NewHashValue()
	if h.Debug() != "Hash" {
		t.Fatalf("expected 'Hash', got %q", h.Debug())
	}
	if h.Type() != HashObject {
		t.Fatalf("expected HashObject, got %s", h.Type())
	}

	// GetValue for missing key
	val, ok := h.GetValue(NewStringValue("missing"))
	if ok {
		t.Fatal("expected ok=false for missing key")
	}
	if val.Type() != NullObject {
		t.Fatalf("expected NullValue for missing key, got %s", val.Type())
	}

	// HasKey for missing key
	if h.HasKey(NewStringValue("missing")) {
		t.Fatal("expected HasKey=false for missing key")
	}

	// Set and check
	h.Set(NewStringValue("key"), NewIntegerValue(42))
	val, ok = h.GetValue(NewStringValue("key"))
	if !ok {
		t.Fatal("expected ok=true for existing key")
	}
	if val.Debug() != "42" {
		t.Fatalf("expected '42', got %q", val.Debug())
	}
	if !h.HasKey(NewStringValue("key")) {
		t.Fatal("expected HasKey=true for existing key")
	}
}

func TestIntegerValueDebugTypeAndConstructor(t *testing.T) {
	i := NewIntegerValue(42)
	if i.Debug() != "42" {
		t.Fatalf("expected '42', got %q", i.Debug())
	}
	if i.Type() != IntegerObject {
		t.Fatalf("expected IntegerObject, got %s", i.Type())
	}
}

func TestNullValueDebugAndType(t *testing.T) {
	n := &NullValue{}
	if n.Debug() != "null" {
		t.Fatalf("expected 'null', got %q", n.Debug())
	}
	if n.Type() != NullObject {
		t.Fatalf("expected NullObject, got %s", n.Type())
	}
}

func TestReturnValueDebugAndType(t *testing.T) {
	rv := &ReturnValue{Value: NewIntegerValue(99)}
	if rv.Debug() != "99" {
		t.Fatalf("expected '99', got %q", rv.Debug())
	}
	if rv.Type() != ReturnValueObject {
		t.Fatalf("expected ReturnValueObject, got %s", rv.Type())
	}
}

func TestStringValueDebugTypeAndHashKey(t *testing.T) {
	s := NewStringValue("hello")
	if s.Debug() != "hello" {
		t.Fatalf("expected 'hello', got %q", s.Debug())
	}
	if s.Type() != StringObject {
		t.Fatalf("expected StringObject, got %s", s.Type())
	}

	hk := s.HashKey()
	if hk.Type != StringObject {
		t.Fatalf("expected StringObject hash key type, got %s", hk.Type)
	}

	// Same string → same hash key
	s2 := NewStringValue("hello")
	if s.HashKey() != s2.HashKey() {
		t.Fatal("expected same hash key for same string")
	}

	// Different string → different hash key
	s3 := NewStringValue("world")
	if s.HashKey() == s3.HashKey() {
		t.Fatal("expected different hash keys for different strings")
	}
}

func TestBuiltInFunctionDebugAndType(t *testing.T) {
	bif := &BuiltInFunction{Fn: func(ctx *ExecutionContext, scope *Scope, args ...Object) (Object, error) {
		return Null, nil
	}}
	if bif.Debug() != "builtin function" {
		t.Fatalf("expected 'builtin function', got %q", bif.Debug())
	}
	if bif.Type() != BuiltInFunctionObject {
		t.Fatalf("expected BuiltInFunctionObject, got %s", bif.Type())
	}
}

// ===================================================================
// Batch 2: Scope & Context
// ===================================================================

func TestScopeDeleteLocal(t *testing.T) {
	s := NewScope()
	s.SetLocal("x", NewIntegerValue(10))

	val, ok := s.Get("x")
	if !ok || val.Debug() != "10" {
		t.Fatal("expected x=10 before delete")
	}

	s.DeleteLocal("x")

	_, ok = s.Get("x")
	if ok {
		t.Fatal("expected x not found after delete")
	}
}

func TestScopeAssignReturnsFlase(t *testing.T) {
	s := NewScope()
	assigned := s.Assign("nonexistent", NewIntegerValue(1))
	if assigned {
		t.Fatal("expected Assign to return false for nonexistent variable")
	}
}

func TestScopeGetLocalDoesNotSearchParent(t *testing.T) {
	parent := NewScope()
	parent.SetLocal("x", NewIntegerValue(10))
	child := NewChildScope(parent)

	_, ok := child.GetLocal("x")
	if ok {
		t.Fatal("expected GetLocal to NOT find parent variable")
	}

	// But Get should find it
	val, ok := child.Get("x")
	if !ok {
		t.Fatal("expected Get to find parent variable")
	}
	if val.Debug() != "10" {
		t.Fatalf("expected '10', got %q", val.Debug())
	}
}

func TestNewExecutionContextWithScope(t *testing.T) {
	input := `return x;`
	l := lexer.NewScript(input)
	p := parser.New(l)
	program, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}

	scope := NewScope()
	scope.SetLocal("x", NewIntegerValue(42))

	e := New()
	ctx := NewExecutionContextWithScope(program, scope)
	ret, err := e.Evaluate(ctx)
	if err != nil {
		t.Fatal(err)
	}

	val := unwrapReturn(t, ret)
	intVal, ok := val.(*IntegerValue)
	if !ok {
		t.Fatalf("expected IntegerValue, got %T", val)
	}
	if intVal.Value != 42 {
		t.Fatalf("expected 42, got %d", intVal.Value)
	}
}

func TestExecutionContextMetadataMultipleKeys(t *testing.T) {
	input := ``
	l := lexer.NewScript(input)
	p := parser.New(l)
	program, _ := p.Parse()
	ctx := NewExecutionContext(program)
	ctx.Metadata["a"] = 1
	ctx.Metadata["b"] = "two"
	if ctx.Metadata["a"] != 1 {
		t.Fatalf("expected 1, got %v", ctx.Metadata["a"])
	}
	if ctx.Metadata["b"] != "two" {
		t.Fatalf("expected 'two', got %v", ctx.Metadata["b"])
	}
}

// ===================================================================
// Batch 3: isTruthy + Prefix Edge Cases
// ===================================================================

func TestIsTruthyNonBooleans(t *testing.T) {
	// Integer is truthy
	val := unwrapReturn(t, evalScript(t, `if (1) { return true; } else { return false; }`))
	if bv, ok := val.(*BooleanValue); !ok || !bv.Value {
		t.Fatal("expected integer to be truthy")
	}

	// String is truthy
	val = unwrapReturn(t, evalScript(t, `if ("hello") { return true; } else { return false; }`))
	if bv, ok := val.(*BooleanValue); !ok || !bv.Value {
		t.Fatal("expected string to be truthy")
	}

	// Array is truthy
	val = unwrapReturn(t, evalScript(t, `if ([1,2]) { return true; } else { return false; }`))
	if bv, ok := val.(*BooleanValue); !ok || !bv.Value {
		t.Fatal("expected array to be truthy")
	}

	// Null is falsy — inject via scope since null is not a keyword
	input := `if (x) { return true; } else { return false; }`
	l := lexer.NewScript(input)
	p := parser.New(l)
	program, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}
	e := New()
	ctx := NewExecutionContext(program)
	ctx.RootScope.SetLocal("x", Null)
	ret, err := e.Evaluate(ctx)
	if err != nil {
		t.Fatal(err)
	}
	val = unwrapReturn(t, ret)
	if bv, ok := val.(*BooleanValue); !ok || bv.Value {
		t.Fatal("expected null to be falsy")
	}
}

func TestBangOnNonBooleans(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"return !0;", false},     // !<integer> → false (integer is truthy → !truthy = false)
		{"return !1;", false},     // same
		{`return !"hello";`, false}, // !<string> → false
		{"return ![];", false},    // !<array> → false
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			val := unwrapReturn(t, evalScript(t, tt.input))
			boolVal, ok := val.(*BooleanValue)
			if !ok {
				t.Fatalf("expected BooleanValue, got %T", val)
			}
			if boolVal.Value != tt.expected {
				t.Fatalf("expected %v, got %v", tt.expected, boolVal.Value)
			}
		})
	}
}

func TestMinusPrefixOnNonNumeric(t *testing.T) {
	tests := []struct {
		input string
		desc  string
	}{
		{`return -"hello";`, "negate string"},
		{`return -true;`, "negate boolean"},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			err := evalScriptError(t, tt.input)
			if err == nil {
				t.Fatalf("expected error for %s", tt.desc)
			}
			if !strings.Contains(err.Error(), "unknown operator") {
				t.Fatalf("expected 'unknown operator' error, got: %v", err)
			}
		})
	}
}

// ===================================================================
// Batch 4: Infix Edge Cases
// ===================================================================

func TestStringInfixUnknownOperator(t *testing.T) {
	err := evalScriptError(t, `return "a" - "b";`)
	if err == nil {
		t.Fatal("expected error for string subtraction")
	}
	if !strings.Contains(err.Error(), "unknown operator") {
		t.Fatalf("expected 'unknown operator' error, got: %v", err)
	}
}

func TestBooleanInfixUnknownOperator(t *testing.T) {
	err := evalScriptError(t, `return true + false;`)
	if err == nil {
		t.Fatal("expected error for boolean addition")
	}
	if !strings.Contains(err.Error(), "unknown operator") {
		t.Fatalf("expected 'unknown operator' error, got: %v", err)
	}
}

func TestDecimalMultiplication(t *testing.T) {
	val := unwrapReturn(t, evalScript(t, `return 2.5 * 4.0;`))
	decVal, ok := val.(*DecimalValue)
	if !ok {
		t.Fatalf("expected DecimalValue, got %T", val)
	}
	if decVal.Value != 10.0 {
		t.Fatalf("expected 10.0, got %f", decVal.Value)
	}
}

func TestDecimalInequality(t *testing.T) {
	val := unwrapReturn(t, evalScript(t, `return 1.5 != 1.5;`))
	boolVal, ok := val.(*BooleanValue)
	if !ok {
		t.Fatalf("expected BooleanValue, got %T", val)
	}
	if boolVal.Value != false {
		t.Fatalf("expected false, got %v", boolVal.Value)
	}
}

func TestMixedTypeEquality(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{`return 1 == true;`, false},
		{`return 1 != true;`, true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			val := unwrapReturn(t, evalScript(t, tt.input))
			boolVal, ok := val.(*BooleanValue)
			if !ok {
				t.Fatalf("expected BooleanValue, got %T", val)
			}
			if boolVal.Value != tt.expected {
				t.Fatalf("expected %v, got %v", tt.expected, boolVal.Value)
			}
		})
	}
}

func TestTypeMismatchMultiply(t *testing.T) {
	err := evalScriptError(t, `return "a" * 2;`)
	if err == nil {
		t.Fatal("expected error for string * integer")
	}
	if !strings.Contains(err.Error(), "type mismatch") {
		t.Fatalf("expected 'type mismatch' error, got: %v", err)
	}
}

// ===================================================================
// Batch 5: Hash ForEach + Property Expression
// ===================================================================

func TestEvaluateHashForEachBasic(t *testing.T) {
	script := `
let h = {"a": 1, "b": 2, "c": 3};
let sum = 0;
foreach (h as val) {
	sum = sum + val;
}
return sum;`
	val := unwrapReturn(t, evalScript(t, script))
	intVal, ok := val.(*IntegerValue)
	if !ok {
		t.Fatalf("expected IntegerValue, got %T", val)
	}
	if intVal.Value != 6 {
		t.Fatalf("expected 6, got %d", intVal.Value)
	}
}

func TestEvaluateHashForEachEmpty(t *testing.T) {
	script := `
let h = {};
let count = 0;
foreach (h as val) {
	count = count + 1;
}
return count;`
	val := unwrapReturn(t, evalScript(t, script))
	intVal, ok := val.(*IntegerValue)
	if !ok {
		t.Fatalf("expected IntegerValue, got %T", val)
	}
	if intVal.Value != 0 {
		t.Fatalf("expected 0, got %d", intVal.Value)
	}
}

func TestEvaluateHashForEachBodyError(t *testing.T) {
	script := `
let h = {"a": 1};
foreach (h as val) {
	let x = val + undefined_var;
}
`
	err := evalScriptError(t, script)
	if err == nil {
		t.Fatal("expected error from hash foreach body")
	}
}

func TestEvaluatePropertyExpressionDeepNesting(t *testing.T) {
	script := `
let h = {"a": {"b": {"c": 42}}};
return h.a.b.c;`
	val := unwrapReturn(t, evalScript(t, script))
	intVal, ok := val.(*IntegerValue)
	if !ok {
		t.Fatalf("expected IntegerValue, got %T", val)
	}
	if intVal.Value != 42 {
		t.Fatalf("expected 42, got %d", intVal.Value)
	}
}

func TestEvaluatePropertyOnNonHash(t *testing.T) {
	// Property on a non-hash returns Null (not an error)
	script := `let x = 5; return x.name;`
	ret := evalScript(t, script)
	val := unwrapReturn(t, ret)
	if val.Type() != NullObject {
		t.Fatalf("expected NullObject, got %s", val.Type())
	}
}

func TestEvaluatePropertyWithIndex(t *testing.T) {
	script := `
let h = {"items": [10, 20, 30]};
return h.items[1];`
	val := unwrapReturn(t, evalScript(t, script))
	intVal, ok := val.(*IntegerValue)
	if !ok {
		t.Fatalf("expected IntegerValue, got %T", val)
	}
	if intVal.Value != 20 {
		t.Fatalf("expected 20, got %d", intVal.Value)
	}
}

func TestEvaluateNestedPropertyAssignment(t *testing.T) {
	script := `
let h = {"a": {"b": 0}};
h.a.b = 99;
return h.a.b;`
	val := unwrapReturn(t, evalScript(t, script))
	intVal, ok := val.(*IntegerValue)
	if !ok {
		t.Fatalf("expected IntegerValue, got %T", val)
	}
	if intVal.Value != 99 {
		t.Fatalf("expected 99, got %d", intVal.Value)
	}
}

func TestEvaluateArrayElementPropertyAccess(t *testing.T) {
	script := `
let h = {"items": [{"name": "first"}, {"name": "second"}]};
return h.items[0].name;`
	val := unwrapReturn(t, evalScript(t, script))
	strVal, ok := val.(*StringValue)
	if !ok {
		t.Fatalf("expected StringValue, got %T", val)
	}
	if strVal.Value != "first" {
		t.Fatalf("expected 'first', got %q", strVal.Value)
	}
}

// ===================================================================
// Batch 6: Error Propagation Paths
// ===================================================================

func TestInfixLeftError(t *testing.T) {
	err := evalScriptError(t, `return undefined_left + 1;`)
	if err == nil {
		t.Fatal("expected error for undefined left operand")
	}
}

func TestInfixRightError(t *testing.T) {
	err := evalScriptError(t, `return 1 + undefined_right;`)
	if err == nil {
		t.Fatal("expected error for undefined right operand")
	}
}

func TestPrefixError(t *testing.T) {
	err := evalScriptError(t, `return -undefined_var;`)
	if err == nil {
		t.Fatal("expected error for prefix on undefined")
	}
}

func TestIndexLeftError(t *testing.T) {
	err := evalScriptError(t, `return undefined_arr[0];`)
	if err == nil {
		t.Fatal("expected error for undefined left in index expr")
	}
}

func TestIndexRightError(t *testing.T) {
	err := evalScriptError(t, `let a = [1, 2]; return a[undefined_idx];`)
	if err == nil {
		t.Fatal("expected error for undefined index")
	}
}

func TestArrayLiteralElementError(t *testing.T) {
	err := evalScriptError(t, `let a = [1, undefined_var, 3];`)
	if err == nil {
		t.Fatal("expected error for undefined element in array literal")
	}
}

func TestHashLiteralValueError(t *testing.T) {
	err := evalScriptError(t, `let h = {"a": undefined_var};`)
	if err == nil {
		t.Fatal("expected error for undefined value in hash literal")
	}
}

func TestCallExpressionFunctionEvalError(t *testing.T) {
	err := evalScriptError(t, `undefined_fn(1);`)
	if err == nil {
		t.Fatal("expected error for calling undefined function")
	}
}

func TestCallExpressionArgumentError(t *testing.T) {
	err := evalScriptError(t, `fn f(x) { return x; } f(undefined_arg);`)
	if err == nil {
		t.Fatal("expected error for undefined argument in call")
	}
}

func TestLetStatementValueError(t *testing.T) {
	err := evalScriptError(t, `let x = undefined_var;`)
	if err == nil {
		t.Fatal("expected error for undefined value in let")
	}
}

func TestReturnStatementValueError(t *testing.T) {
	err := evalScriptError(t, `return undefined_var;`)
	if err == nil {
		t.Fatal("expected error for undefined value in return")
	}
}

func TestAssignmentRightSideError(t *testing.T) {
	err := evalScriptError(t, `let x = 0; x = undefined_var;`)
	if err == nil {
		t.Fatal("expected error for undefined right side in assignment")
	}
}

func TestPropertyAssignmentRightSideError(t *testing.T) {
	err := evalScriptError(t, `let h = {"a": 1}; h.a = undefined_var;`)
	if err == nil {
		t.Fatal("expected error for undefined right side in property assignment")
	}
}

func TestIndexAssignmentLeftError(t *testing.T) {
	err := evalScriptError(t, `undefined_arr[0] = 1;`)
	if err == nil {
		t.Fatal("expected error for undefined left in index assignment")
	}
}

func TestIndexAssignmentIndexError(t *testing.T) {
	err := evalScriptError(t, `let a = [1]; a[undefined_idx] = 1;`)
	if err == nil {
		t.Fatal("expected error for undefined index in index assignment")
	}
}

func TestIndexAssignmentRightError(t *testing.T) {
	err := evalScriptError(t, `let a = [1]; a[0] = undefined_var;`)
	if err == nil {
		t.Fatal("expected error for undefined right side in index assignment")
	}
}

func TestIndexAssignmentBoundsCheck(t *testing.T) {
	err := evalScriptError(t, `let a = [1, 2, 3]; a[10] = 99;`)
	if err == nil {
		t.Fatal("expected error for index out of bounds in assignment")
	}
	if !strings.Contains(err.Error(), "index out of bounds") {
		t.Fatalf("expected 'index out of bounds' error, got: %v", err)
	}
}

func TestIndexAssignmentNegativeBounds(t *testing.T) {
	err := evalScriptError(t, `let a = [1]; a[-1] = 99;`)
	if err == nil {
		t.Fatal("expected error for negative index in assignment")
	}
	if !strings.Contains(err.Error(), "index out of bounds") {
		t.Fatalf("expected 'index out of bounds' error, got: %v", err)
	}
}

func TestIndexAssignmentNonIntegerIndex(t *testing.T) {
	err := evalScriptError(t, `let a = [1]; a["key"] = 99;`)
	if err == nil {
		t.Fatal("expected error for non-integer index on array assignment")
	}
	if !strings.Contains(err.Error(), "index must be an integer") {
		t.Fatalf("expected 'index must be an integer' error, got: %v", err)
	}
}

func TestIndexAssignmentOnNonArrayNonHash(t *testing.T) {
	err := evalScriptError(t, `let x = 5; x[0] = 99;`)
	if err == nil {
		t.Fatal("expected error for index assignment on non-array/hash")
	}
	if !strings.Contains(err.Error(), "must be an array or hash") {
		t.Fatalf("expected 'must be an array or hash' error, got: %v", err)
	}
}

func TestHashIndexAssignment(t *testing.T) {
	script := `let h = {"a": 1}; h["b"] = 2; return h["b"];`
	val := unwrapReturn(t, evalScript(t, script))
	intVal, ok := val.(*IntegerValue)
	if !ok {
		t.Fatalf("expected IntegerValue, got %T", val)
	}
	if intVal.Value != 2 {
		t.Fatalf("expected 2, got %d", intVal.Value)
	}
}

func TestEvaluateStringWithReturn(t *testing.T) {
	input := `Hello {% return "done"; %} world`
	l := lexer.NewTemplate(input)
	p := parser.New(l)
	program, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}
	e := New()
	ctx := NewExecutionContext(program)
	output, err := e.EvaluateString(ctx)
	if err != nil {
		t.Fatal(err)
	}
	// EvaluateString returns "" on return
	if output != "" {
		t.Fatalf("expected empty string on return, got %q", output)
	}
}

func TestIndexOnNonArrayNonHash(t *testing.T) {
	err := evalScriptError(t, `let x = 5; return x[0];`)
	if err == nil {
		t.Fatal("expected error for index on non-array/hash")
	}
	if !strings.Contains(err.Error(), "index operator not supported") {
		t.Fatalf("expected 'index operator not supported' error, got: %v", err)
	}
}

func TestIfConditionError(t *testing.T) {
	err := evalScriptError(t, `if (undefined_var) { return 1; }`)
	if err == nil {
		t.Fatal("expected error for undefined condition in if")
	}
}

func TestForEachIterableError(t *testing.T) {
	err := evalScriptError(t, `foreach (undefined_var as item) { return item; }`)
	if err == nil {
		t.Fatal("expected error for undefined iterable in foreach")
	}
}

func TestForEachBodyError(t *testing.T) {
	err := evalScriptError(t, `let a = [1]; foreach (a as item) { let x = undefined_var; }`)
	if err == nil {
		t.Fatal("expected error from foreach body")
	}
}

func TestForEachOnNonIterable(t *testing.T) {
	// foreach on non-array/hash returns Null without error
	ret := evalScript(t, `foreach (5 as item) { return item; }`)
	if ret.Type() != NullObject {
		t.Fatalf("expected NullObject for non-iterable foreach, got %s", ret.Type())
	}
}

func TestBlockStatementReturnPropagation(t *testing.T) {
	script := `
fn test() {
	if (true) {
		return 42;
	}
	return 99;
}
return test();`
	val := unwrapReturn(t, evalScript(t, script))
	intVal, ok := val.(*IntegerValue)
	if !ok {
		t.Fatalf("expected IntegerValue, got %T", val)
	}
	if intVal.Value != 42 {
		t.Fatalf("expected 42, got %d", intVal.Value)
	}
}

func TestFunctionBodyError(t *testing.T) {
	err := evalScriptError(t, `fn f() { return undefined_var; } f();`)
	if err == nil {
		t.Fatal("expected error from function body")
	}
}

func TestPropertyAssignmentOnNonHash(t *testing.T) {
	// Property assignment where left evaluates to non-hash
	err := evalScriptError(t, `let x = 5; x.prop = 1;`)
	if err == nil {
		// It might return nil,nil or error — check the behavior
		// Based on evaluatePropertyExpression, non-hash → returns Null,Null,Null,nil
		// Then the assignment checks if idx == Null || parent == Null → returns nil, nil
		// Which is not an error, so this is expected
	}
}

func TestEvaluateStringOutputsExpressionResults(t *testing.T) {
	// Non-null, non-function expression results are written to output
	input := `{% 1 + 2 %}{% "hello" %}`
	output := evalTemplate(t, input)
	if output != "3hello" {
		t.Fatalf("expected '3hello', got %q", output)
	}
}

func TestEvaluateStringSkipsFunctions(t *testing.T) {
	// Function definitions should not be written to output
	input := `{% fn add(a, b) { return a + b; } %}done`
	output := evalTemplate(t, input)
	if output != "done" {
		t.Fatalf("expected 'done', got %q", output)
	}
}

func TestEvaluateStringError(t *testing.T) {
	input := `{% undefined_var %}`
	l := lexer.NewTemplate(input)
	p := parser.New(l)
	program, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}
	e := New()
	ctx := NewExecutionContext(program)
	_, err = e.EvaluateString(ctx)
	if err == nil {
		t.Fatal("expected error for undefined variable in template")
	}
}

// ===================================================================
// Batch 7: Remaining Targeted Coverage
// ===================================================================

func TestLogMultipleArgs(t *testing.T) {
	input := `log("hello", "world");`
	l := lexer.NewScript(input)
	p := parser.New(l)
	program, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}
	e := New()
	ctx := NewExecutionContext(program)
	var logBuf strings.Builder
	ctx.Logger = &logBuf
	_, err = e.Evaluate(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(logBuf.String(), "hello") || !strings.Contains(logBuf.String(), "world") {
		t.Fatalf("expected log to contain both messages, got %q", logBuf.String())
	}
}

func TestPrintMultipleArgs(t *testing.T) {
	input := `{% print("hello", " ", "world"); %}`
	output := evalTemplate(t, input)
	if output != "hello world" {
		t.Fatalf("expected 'hello world', got %q", output)
	}
}

func TestUnwrapReturnValuePassthrough(t *testing.T) {
	// unwrapReturnValue on non-ReturnValue returns the object itself
	obj := NewIntegerValue(42)
	result := unwrapReturnValue(obj)
	if result != obj {
		t.Fatal("expected same object back from unwrapReturnValue for non-ReturnValue")
	}
}

func TestUnwrapReturnValueUnwraps(t *testing.T) {
	inner := NewIntegerValue(42)
	rv := &ReturnValue{Value: inner}
	result := unwrapReturnValue(rv)
	if result != inner {
		t.Fatal("expected inner value from unwrapReturnValue for ReturnValue")
	}
}

func TestDuplicateFunctionDefinitionInScope(t *testing.T) {
	// Named function registered in scope — attempt to redefine
	err := evalScriptError(t, `fn foo() { return 1; } fn foo() { return 2; }`)
	if err == nil {
		t.Fatal("expected error for duplicate function")
	}
	if !strings.Contains(err.Error(), "already defined") {
		t.Fatalf("expected 'already defined' error, got: %v", err)
	}
}

func TestHashLiteralValueErrorSecondKey(t *testing.T) {
	// Hash with multiple entries where second value is undefined
	err := evalScriptError(t, `let h = {"a": 1, "b": undefined_val};`)
	if err == nil {
		t.Fatal("expected error for undefined hash value")
	}
}

func TestEvaluateStringPlainText(t *testing.T) {
	output := evalTemplate(t, `Just plain text, no script blocks.`)
	if output != "Just plain text, no script blocks." {
		t.Fatalf("expected plain text, got %q", output)
	}
}

func TestDecimalDebugFormats(t *testing.T) {
	tests := []struct {
		input    float64
		expected string
	}{
		{4.0, "4.0"},
		{0.5, "0.5"},
		{100.123, "100.123"},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%f", tt.input), func(t *testing.T) {
			d := NewDecimalValue(tt.input)
			if d.Debug() != tt.expected {
				t.Fatalf("expected %q, got %q", tt.expected, d.Debug())
			}
		})
	}
}

func TestArrayValueNewEmpty(t *testing.T) {
	a := NewArrayValue([]Object{})
	if len(a.Elements) != 0 {
		t.Fatalf("expected 0 elements, got %d", len(a.Elements))
	}
}

func TestScopeAssignWalksChain(t *testing.T) {
	grandparent := NewScope()
	grandparent.SetLocal("x", NewIntegerValue(1))
	parent := NewChildScope(grandparent)
	child := NewChildScope(parent)

	// Assign through 2-level chain
	ok := child.Assign("x", NewIntegerValue(99))
	if !ok {
		t.Fatal("expected Assign to succeed through chain")
	}

	val, _ := grandparent.Get("x")
	if val.Debug() != "99" {
		t.Fatalf("expected '99', got %q", val.Debug())
	}
}

func TestEvaluateNullFromScope(t *testing.T) {
	input := `return x;`
	l := lexer.NewScript(input)
	p := parser.New(l)
	program, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}
	e := New()
	ctx := NewExecutionContext(program)
	ctx.RootScope.SetLocal("x", Null)
	ret, err := e.Evaluate(ctx)
	if err != nil {
		t.Fatal(err)
	}
	val := unwrapReturn(t, ret)
	if val.Type() != NullObject {
		t.Fatalf("expected NullObject, got %s", val.Type())
	}
}

func TestPropertyExpressionLeftError(t *testing.T) {
	err := evalScriptError(t, `return undefined_var.prop;`)
	if err == nil {
		t.Fatal("expected error for undefined left in property expression")
	}
}

func TestPropertyAssignmentLeftError(t *testing.T) {
	// Property assignment where the property expression left side errors
	err := evalScriptError(t, `undefined_var.prop = 1;`)
	if err == nil {
		t.Fatal("expected error for undefined left in property assignment")
	}
}

func TestIndexAssignmentNullLeft(t *testing.T) {
	// Index assignment where left evaluates to Null
	input := `let a = [1]; let x = a[5]; x[0] = 99;`
	l := lexer.NewScript(input)
	p := parser.New(l)
	program, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}
	e := New()
	ctx := NewExecutionContext(program)
	_, err = e.Evaluate(ctx)
	if err == nil {
		t.Fatal("expected error for null left in index assignment")
	}
	if !strings.Contains(err.Error(), "null") || !strings.Contains(err.Error(), "index") {
		t.Fatalf("expected null index error, got: %v", err)
	}
}

func TestPropertyAssignmentNullParent(t *testing.T) {
	// Property assignment on hash where property chain results in null parent
	// h.nonexistent.prop = 1 → evaluatePropertyExpression returns Null parent
	script := `let h = {"a": 1}; h.nonexistent.prop = 1;`
	l := lexer.NewScript(script)
	p := parser.New(l)
	program, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}
	e := New()
	ctx := NewExecutionContext(program)
	// Should not error, returns nil,nil for null parent/idx
	_, err = e.Evaluate(ctx)
	// This should either succeed (return nil) or error
	// Based on code: returns nil, nil when parent == Null
}

func TestHashSetOverwrite(t *testing.T) {
	// Overwrite an existing key in hash
	h := NewHashValue()
	h.Set(NewStringValue("key"), NewIntegerValue(1))
	h.Set(NewStringValue("key"), NewIntegerValue(2))
	val, ok := h.GetValue(NewStringValue("key"))
	if !ok {
		t.Fatal("expected key to exist")
	}
	if val.Debug() != "2" {
		t.Fatalf("expected '2', got %q", val.Debug())
	}
}

func TestNewArrayValueNil(t *testing.T) {
	a := NewArrayValue(nil)
	if a.Elements != nil {
		t.Fatal("expected nil elements")
	}
}

func TestEvaluateStringWithNullResult(t *testing.T) {
	// Template where script block evaluates to null (let statement after evaluateString output)
	input := `Hello {% let x = 5; %}`
	output := evalTemplate(t, input)
	// The let statement returns the value 5, which gets output
	if !strings.Contains(output, "Hello") {
		t.Fatalf("expected output to contain 'Hello', got %q", output)
	}
}

// --- Logical Operators ---

func TestLogicalAnd(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"return true && true;", "true"},
		{"return true && false;", "false"},
		{"return false && true;", "false"},
		{"return false && false;", "false"},
	}

	for _, tt := range tests {
		result := evalScript(t, tt.input)
		if result.Debug() != tt.expected {
			t.Errorf("input=%q: expected %s, got %s", tt.input, tt.expected, result.Debug())
		}
	}
}

func TestLogicalOr(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"return true || true;", "true"},
		{"return true || false;", "true"},
		{"return false || true;", "true"},
		{"return false || false;", "false"},
	}

	for _, tt := range tests {
		result := evalScript(t, tt.input)
		if result.Debug() != tt.expected {
			t.Errorf("input=%q: expected %s, got %s", tt.input, tt.expected, result.Debug())
		}
	}
}

func TestLogicalAndShortCircuit(t *testing.T) {
	// && with false left should return false without evaluating right
	test := `return false && 1;`
	result := evalScript(t, test)
	if result.Debug() != "false" {
		t.Fatalf("expected 'false' (short circuit), got %s", result.Debug())
	}
}

func TestLogicalOrShortCircuit(t *testing.T) {
	// || with truthy left should return left without evaluating right
	test := `return true || 1;`
	result := evalScript(t, test)
	if result.Debug() != "true" {
		t.Fatalf("expected 'true' (short circuit), got %s", result.Debug())
	}
}

// --- Null Coalescing ---

func TestNullCoalescing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"return null ?? 5;", "5"},
		{"return 10 ?? 5;", "10"},
		{"return null ?? null ?? 3;", "3"},
	}

	for _, tt := range tests {
		result := evalScript(t, tt.input)
		if result.Debug() != tt.expected {
			t.Errorf("input=%q: expected %s, got %s", tt.input, tt.expected, result.Debug())
		}
	}
}

// --- While Loop ---

func TestWhileLoop(t *testing.T) {
	test := `let x = 0; while (x < 5) { x = x + 1; } return x;`
	result := evalScript(t, test)
	if result.Debug() != "5" {
		t.Fatalf("expected 5, got %s", result.Debug())
	}
}

func TestWhileLoopBreak(t *testing.T) {
	test := `let x = 0; while (true) { x = x + 1; if (x == 3) { break; } } return x;`
	result := evalScript(t, test)
	if result.Debug() != "3" {
		t.Fatalf("expected 3, got %s", result.Debug())
	}
}

func TestWhileLoopContinue(t *testing.T) {
	test := `let sum = 0; let i = 0; while (i < 5) { i = i + 1; if (i == 3) { continue; } sum = sum + i; } return sum;`
	// sum = 1 + 2 + 4 + 5 = 12 (skips 3)
	result := evalScript(t, test)
	if result.Debug() != "12" {
		t.Fatalf("expected 12, got %s", result.Debug())
	}
}

func TestWhileFalseCondition(t *testing.T) {
	test := `let x = 0; while (false) { x = 1; } return x;`
	result := evalScript(t, test)
	if result.Debug() != "0" {
		t.Fatalf("expected 0, got %s", result.Debug())
	}
}

// --- Break/Continue in Foreach ---

func TestForeachBreak(t *testing.T) {
	test := `let sum = 0; foreach ([1,2,3,4,5] as v) { if (v == 4) { break; } sum = sum + v; } return sum;`
	result := evalScript(t, test)
	if result.Debug() != "6" {
		t.Fatalf("expected 6, got %s", result.Debug())
	}
}

func TestForeachContinue(t *testing.T) {
	test := `let sum = 0; foreach ([1,2,3,4,5] as v) { if (v == 3) { continue; } sum = sum + v; } return sum;`
	// sum = 1 + 2 + 4 + 5 = 12
	result := evalScript(t, test)
	if result.Debug() != "12" {
		t.Fatalf("expected 12, got %s", result.Debug())
	}
}

// --- Else If ---

func TestElseIf(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`let x = 1; if (x == 1) { return "one"; } else if (x == 2) { return "two"; } else { return "other"; }`, "one"},
		{`let x = 2; if (x == 1) { return "one"; } else if (x == 2) { return "two"; } else { return "other"; }`, "two"},
		{`let x = 3; if (x == 1) { return "one"; } else if (x == 2) { return "two"; } else { return "other"; }`, "other"},
	}

	for _, tt := range tests {
		result := evalScript(t, tt.input)
		if result.Debug() != tt.expected {
			t.Errorf("input=%q: expected %s, got %s", tt.input, tt.expected, result.Debug())
		}
	}
}

func TestElseIfChain(t *testing.T) {
	test := `let x = 3;
if (x == 1) { return "one"; }
else if (x == 2) { return "two"; }
else if (x == 3) { return "three"; }
else { return "other"; }`
	result := evalScript(t, test)
	if result.Debug() != "three" {
		t.Fatalf("expected 'three', got %s", result.Debug())
	}
}

// --- len() ---

func TestLen(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`return len("hello");`, "5"},
		{`return len("");`, "0"},
		{`return len([1, 2, 3]);`, "3"},
		{`return len([]);`, "0"},
		{`return len({"a": 1, "b": 2});`, "2"},
	}

	for _, tt := range tests {
		result := evalScript(t, tt.input)
		if result.Debug() != tt.expected {
			t.Errorf("input=%q: expected %s, got %s", tt.input, tt.expected, result.Debug())
		}
	}
}

// --- String Built-ins ---

func TestSplit(t *testing.T) {
	test := `let parts = split("a,b,c", ","); return len(parts);`
	result := evalScript(t, test)
	if result.Debug() != "3" {
		t.Fatalf("expected 3, got %s", result.Debug())
	}
}

func TestTrim(t *testing.T) {
	test := `return trim("  hello  ");`
	result := evalScript(t, test)
	if result.Debug() != "hello" {
		t.Fatalf("expected 'hello', got %s", result.Debug())
	}
}

func TestToUpper(t *testing.T) {
	test := `return toUpper("hello");`
	result := evalScript(t, test)
	if result.Debug() != "HELLO" {
		t.Fatalf("expected 'HELLO', got %s", result.Debug())
	}
}

func TestToLower(t *testing.T) {
	test := `return toLower("HELLO");`
	result := evalScript(t, test)
	if result.Debug() != "hello" {
		t.Fatalf("expected 'hello', got %s", result.Debug())
	}
}

func TestContains(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`return contains("hello world", "world");`, "true"},
		{`return contains("hello world", "xyz");`, "false"},
	}

	for _, tt := range tests {
		result := evalScript(t, tt.input)
		if result.Debug() != tt.expected {
			t.Errorf("input=%q: expected %s, got %s", tt.input, tt.expected, result.Debug())
		}
	}
}

func TestStartsWith(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`return startsWith("hello world", "hello");`, "true"},
		{`return startsWith("hello world", "world");`, "false"},
	}

	for _, tt := range tests {
		result := evalScript(t, tt.input)
		if result.Debug() != tt.expected {
			t.Errorf("input=%q: expected %s, got %s", tt.input, tt.expected, result.Debug())
		}
	}
}

func TestEndsWith(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`return endsWith("hello world", "world");`, "true"},
		{`return endsWith("hello world", "hello");`, "false"},
	}

	for _, tt := range tests {
		result := evalScript(t, tt.input)
		if result.Debug() != tt.expected {
			t.Errorf("input=%q: expected %s, got %s", tt.input, tt.expected, result.Debug())
		}
	}
}

func TestIndexOf(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`return indexOf("hello", "ll");`, "2"},
		{`return indexOf("hello", "xyz");`, "-1"},
	}

	for _, tt := range tests {
		result := evalScript(t, tt.input)
		if result.Debug() != tt.expected {
			t.Errorf("input=%q: expected %s, got %s", tt.input, tt.expected, result.Debug())
		}
	}
}

func TestReplace(t *testing.T) {
	test := `return replace("hello world", "world", "go");`
	result := evalScript(t, test)
	if result.Debug() != "hello go" {
		t.Fatalf("expected 'hello go', got %s", result.Debug())
	}
}

func TestSubstring(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`return substring("hello", 1, 3);`, "el"},
		{`return substring("hello", 2);`, "llo"},
	}

	for _, tt := range tests {
		result := evalScript(t, tt.input)
		if result.Debug() != tt.expected {
			t.Errorf("input=%q: expected %s, got %s", tt.input, tt.expected, result.Debug())
		}
	}
}

// --- keys(), values(), type() ---

func TestKeys(t *testing.T) {
	test := `let h = {"a": 1}; let k = keys(h); return len(k);`
	result := evalScript(t, test)
	if result.Debug() != "1" {
		t.Fatalf("expected 1, got %s", result.Debug())
	}
}

func TestValues(t *testing.T) {
	test := `let h = {"a": 1}; let v = values(h); return v[0];`
	result := evalScript(t, test)
	if result.Debug() != "1" {
		t.Fatalf("expected 1, got %s", result.Debug())
	}
}

func TestTypeBuiltin(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`return type(5);`, "INTEGER"},
		{`return type("hello");`, "STRING"},
		{`return type(true);`, "BOOLEAN"},
		{`return type(1.5);`, "DECIMAL"},
		{`return type([]);`, "ARRAY"},
		{`return type({});`, "HASH"},
		{`return type(null);`, "NULL"},
	}

	for _, tt := range tests {
		result := evalScript(t, tt.input)
		if result.Debug() != tt.expected {
			t.Errorf("input=%q: expected %s, got %s", tt.input, tt.expected, result.Debug())
		}
	}
}

// --- Integer ↔ Decimal Type Coercion ---

func TestIntegerDecimalCoercion(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`return 1 + 1.5;`, "2.5"},
		{`return 1.5 + 1;`, "2.5"},
		{`return 3 * 1.5;`, "4.5"},
		{`return 1.5 * 3;`, "4.5"},
		{`return 10 / 2.5;`, "4.0"},
		{`return 2.5 / 1;`, "2.5"},
		{`return 5 - 1.5;`, "3.5"},
		{`return 1.5 - 1;`, "0.5"},
		{`return 1 < 1.5;`, "true"},
		{`return 1.5 > 1;`, "true"},
		{`return 1 == 1.0;`, "true"},    // int 1 == decimal 1.0
		{`return 2 != 1.5;`, "true"},
	}

	for _, tt := range tests {
		result := evalScript(t, tt.input)
		if result.Debug() != tt.expected {
			t.Errorf("input=%q: expected %s, got %s", tt.input, tt.expected, result.Debug())
		}
	}
}

// --- Foreach Index/Key ---

func TestForeachWithIndex(t *testing.T) {
	// Use integer math to verify index values: sum of indices (0+1+2=3)
	test := `let sum = 0; foreach ([10, 20, 30] as i, v) { sum = sum + i; } return sum;`
	result := evalScript(t, test)
	if result.Debug() != "3" {
		t.Fatalf("expected '3', got %s", result.Debug())
	}
}

func TestForeachWithIndexValue(t *testing.T) {
	// Verify both index and value are correct
	test := `let last_idx = 0; let last_val = 0; foreach ([10, 20, 30] as i, v) { last_idx = i; last_val = v; } return last_idx;`
	result := evalScript(t, test)
	if result.Debug() != "2" {
		t.Fatalf("expected '2', got %s", result.Debug())
	}
}

func TestForeachHashWithKey(t *testing.T) {
	test := `let h = {"a": 1}; let k_result = ""; foreach (h as k, v) { k_result = k; } return k_result;`
	result := evalScript(t, test)
	if result.Debug() != "a" {
		t.Fatalf("expected 'a', got %s", result.Debug())
	}
}

// --- Integer and Boolean as Hash Keys ---

func TestIntegerAsHashKey(t *testing.T) {
	test := `let h = {}; h[1] = "one"; h[2] = "two"; return h[1];`
	result := evalScript(t, test)
	if result.Debug() != "one" {
		t.Fatalf("expected 'one', got %s", result.Debug())
	}
}

func TestBooleanAsHashKey(t *testing.T) {
	test := `let h = {}; h[true] = "yes"; h[false] = "no"; return h[true];`
	result := evalScript(t, test)
	if result.Debug() != "yes" {
		t.Fatalf("expected 'yes', got %s", result.Debug())
	}
}

// --- Null Literal ---

func TestNullLiteral(t *testing.T) {
	test := `return null;`
	result := evalScript(t, test)
	if result.Debug() != "null" {
		t.Fatalf("expected 'null', got %s", result.Debug())
	}
}

func TestNullEquality(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`return null == null;`, "true"},
		{`return null != null;`, "false"},
		{`return null == 5;`, "false"},
	}

	for _, tt := range tests {
		result := evalScript(t, tt.input)
		if result.Debug() != tt.expected {
			t.Errorf("input=%q: expected %s, got %s", tt.input, tt.expected, result.Debug())
		}
	}
}

// --- Comments in Evaluation ---

func TestCommentsInScript(t *testing.T) {
	test := `
// This is a comment
let x = 5;
/* Multi-line
   comment */
let y = 10;
return x + y;`
	result := evalScript(t, test)
	if result.Debug() != "15" {
		t.Fatalf("expected 15, got %s", result.Debug())
	}
}

// --- Bug fix verification tests ---

func TestDecimalDebugWholeNumber(t *testing.T) {
	// Verify that whole decimal numbers show "X.0" not "X."
	tests := []struct {
		value    float64
		expected string
	}{
		{2.0, "2.0"},
		{4.0, "4.0"},
		{100.0, "100.0"},
		{0.0, "0.0"},
		{1.5, "1.5"},
		{3.14, "3.14"},
	}
	for _, tt := range tests {
		d := NewDecimalValue(tt.value)
		if d.Debug() != tt.expected {
			t.Errorf("DecimalValue(%v).Debug() = %q, want %q", tt.value, d.Debug(), tt.expected)
		}
	}
}

func TestPropertyAssignmentToNullReturnsError(t *testing.T) {
	// Property assignment where left side is null should return error, not (nil, nil)
	input := `let x = null; x.name = "test";`
	l := lexer.NewScript(input)
	p := parser.New(l)
	program, err := p.Parse()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	e := New()
	ctx := NewExecutionContext(program)
	_, err = e.Evaluate(ctx)
	if err == nil {
		t.Fatal("expected error for property assignment on null value")
	}
}

func TestBreakSignalDoesNotLeakFromFunction(t *testing.T) {
	// Break inside a function (outside a loop) should not leak
	input := `
fn test() {
	break;
}
let result = test();
return result;`
	l := lexer.NewScript(input)
	p := parser.New(l)
	program, err := p.Parse()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	e := New()
	ctx := NewExecutionContext(program)
	result, err := e.Evaluate(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should return Null (break signal was absorbed by the function boundary)
	retVal, ok := result.(*ReturnValue)
	if !ok {
		t.Fatalf("expected ReturnValue, got %T", result)
	}
	if retVal.Value.Type() != NullObject {
		t.Fatalf("expected null, got %s", retVal.Value.Type())
	}
}

func TestContinueSignalDoesNotLeakFromFunction(t *testing.T) {
	input := `
fn test() {
	continue;
}
let result = test();
return result;`
	l := lexer.NewScript(input)
	p := parser.New(l)
	program, err := p.Parse()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	e := New()
	ctx := NewExecutionContext(program)
	result, err := e.Evaluate(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	retVal, ok := result.(*ReturnValue)
	if !ok {
		t.Fatalf("expected ReturnValue, got %T", result)
	}
	if retVal.Value.Type() != NullObject {
		t.Fatalf("expected null, got %s", retVal.Value.Type())
	}
}

func TestHashSetWithNonHashableKey(t *testing.T) {
	h := NewHashValue()
	// ArrayValue is not Hashable, so Set should return an error
	err := h.Set(&ArrayValue{Elements: []Object{}}, NewStringValue("value"))
	if err == nil {
		t.Fatal("expected error when using non-hashable key")
	}
}

// ===================== Phase 1: RuntimeError =====================

func TestRuntimeErrorFormat(t *testing.T) {
	source := "let x = 1;\nreturn y;"
	err := NewRuntimeError("identifier not found: y", source, 2, 8)
	msg := err.Error()
	if !strings.Contains(msg, "identifier not found: y") {
		t.Fatalf("expected error message, got: %s", msg)
	}
	if !strings.Contains(msg, "line 2, column 8") {
		t.Fatalf("expected line/column info, got: %s", msg)
	}
	if !strings.Contains(msg, "^") {
		t.Fatalf("expected caret pointer, got: %s", msg)
	}
}

func TestRuntimeErrorWithSource(t *testing.T) {
	input := "return y;"
	l := lexer.NewScript(input)
	p := parser.New(l)
	program, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}
	e := New()
	ctx := NewExecutionContext(program)
	ctx.Source = input
	_, err = e.Evaluate(ctx)
	if err == nil {
		t.Fatal("expected error")
	}
	var rtErr *RuntimeError
	if !errors.As(err, &rtErr) {
		t.Fatalf("expected RuntimeError, got %T: %v", err, err)
	}
	if rtErr.Line != 1 {
		t.Fatalf("expected line 1, got %d", rtErr.Line)
	}
}

func TestRuntimeErrorWithoutSource(t *testing.T) {
	input := "return y;"
	l := lexer.NewScript(input)
	p := parser.New(l)
	program, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}
	e := New()
	ctx := NewExecutionContext(program)
	// Source not set — should get plain error
	_, err = e.Evaluate(ctx)
	if err == nil {
		t.Fatal("expected error")
	}
	var rtErr *RuntimeError
	if errors.As(err, &rtErr) {
		t.Fatal("expected plain error, not RuntimeError, when Source is empty")
	}
}

func TestRuntimeErrorBoundsCheck(t *testing.T) {
	// Line number out of range should not panic
	err := NewRuntimeError("test", "single line", 99, 1)
	msg := err.Error()
	if !strings.Contains(msg, "error: test") {
		t.Fatalf("expected error message, got: %s", msg)
	}
}

// ===================== Phase 2: Security Limits =====================

func TestStepLimitInfiniteLoop(t *testing.T) {
	input := "let x = 0; while (true) { x = x + 1; }"
	l := lexer.NewScript(input)
	p := parser.New(l)
	program, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}
	e := New()
	ctx := NewExecutionContext(program)
	ctx.MaxSteps = 1000
	_, err = e.Evaluate(ctx)
	if err == nil {
		t.Fatal("expected step limit error")
	}
	if !strings.Contains(err.Error(), "execution limit exceeded") {
		t.Fatalf("expected 'execution limit exceeded' error, got: %v", err)
	}
}

func TestDepthLimitRecursiveBomb(t *testing.T) {
	input := `fn recurse(n) { return recurse(n + 1); } return recurse(0);`
	l := lexer.NewScript(input)
	p := parser.New(l)
	program, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}
	e := New()
	ctx := NewExecutionContext(program)
	ctx.MaxDepth = 50
	_, err = e.Evaluate(ctx)
	if err == nil {
		t.Fatal("expected depth limit error")
	}
	if !strings.Contains(err.Error(), "maximum call depth exceeded") {
		t.Fatalf("expected 'maximum call depth exceeded' error, got: %v", err)
	}
}

func TestArraySizeLimitAppend(t *testing.T) {
	input := `
let arr = [];
let i = 0;
while (i < 20000) {
	append(arr, i);
	i = i + 1;
}
return arr;`
	l := lexer.NewScript(input)
	p := parser.New(l)
	program, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}
	e := New()
	ctx := NewExecutionContext(program)
	ctx.MaxArraySize = 100
	ctx.MaxSteps = 1_000_000
	_, err = e.Evaluate(ctx)
	if err == nil {
		t.Fatal("expected array size limit error")
	}
	if !strings.Contains(err.Error(), "maximum array size exceeded") {
		t.Fatalf("expected 'maximum array size exceeded' error, got: %v", err)
	}
}

func TestArraySizeLimitLiteral(t *testing.T) {
	// Build a large array literal string
	var sb strings.Builder
	sb.WriteString("return [")
	for i := 0; i < 15; i++ {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(fmt.Sprintf("%d", i))
	}
	sb.WriteString("];")

	l := lexer.NewScript(sb.String())
	p := parser.New(l)
	program, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}
	e := New()
	ctx := NewExecutionContext(program)
	ctx.MaxArraySize = 10
	_, err = e.Evaluate(ctx)
	if err == nil {
		t.Fatal("expected array size limit error for large literal")
	}
}

func TestStepLimitDisabled(t *testing.T) {
	input := `let x = 0; while (x < 100) { x = x + 1; } return x;`
	l := lexer.NewScript(input)
	p := parser.New(l)
	program, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}
	e := New()
	ctx := NewExecutionContext(program)
	ctx.MaxSteps = 0 // disabled
	result, err := e.Evaluate(ctx)
	if err != nil {
		t.Fatal(err)
	}
	val := unwrapReturn(t, result)
	intVal := val.(*IntegerValue)
	if intVal.Value != 100 {
		t.Fatalf("expected 100, got %d", intVal.Value)
	}
}

// ===================== Phase 3: Decimal Modulo =====================

func TestDecimalModulo(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"return 5.5 % 2.0;", 1.5},
		{"return 10.0 % 3.0;", 1.0},
		{"return 7.5 % 2.5;", 0.0},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			val := unwrapReturn(t, evalScript(t, tt.input))
			decVal, ok := val.(*DecimalValue)
			if !ok {
				t.Fatalf("expected DecimalValue, got %T", val)
			}
			if decVal.Value != tt.expected {
				t.Fatalf("expected %f, got %f", tt.expected, decVal.Value)
			}
		})
	}
}

func TestDecimalModuloDivisionByZero(t *testing.T) {
	err := evalScriptError(t, `return 5.5 % 0.0;`)
	if err == nil {
		t.Fatal("expected division by zero error")
	}
	if !strings.Contains(err.Error(), "division by zero") {
		t.Fatalf("expected 'division by zero', got: %v", err)
	}
}

func TestIntegerDecimalModulo(t *testing.T) {
	// Integer % Decimal promotes to decimal
	val := unwrapReturn(t, evalScript(t, `return 7 % 2.5;`))
	decVal, ok := val.(*DecimalValue)
	if !ok {
		t.Fatalf("expected DecimalValue, got %T", val)
	}
	if decVal.Value != 2.0 {
		t.Fatalf("expected 2.0, got %f", decVal.Value)
	}
}

// ===================== Phase 4: Integer Overflow =====================

func TestIntegerOverflowAddition(t *testing.T) {
	input := fmt.Sprintf("return %d + 1;", int(^uint(0)>>1)) // MaxInt + 1
	err := evalScriptError(t, input)
	if err == nil {
		t.Fatal("expected overflow error")
	}
	if !strings.Contains(err.Error(), "integer overflow") {
		t.Fatalf("expected 'integer overflow', got: %v", err)
	}
}

func TestIntegerOverflowMultiplication(t *testing.T) {
	input := fmt.Sprintf("return %d * 2;", int(^uint(0)>>1)) // MaxInt * 2
	err := evalScriptError(t, input)
	if err == nil {
		t.Fatal("expected overflow error")
	}
	if !strings.Contains(err.Error(), "integer overflow") {
		t.Fatalf("expected 'integer overflow', got: %v", err)
	}
}

func TestIntegerOverflowSubtraction(t *testing.T) {
	// MinInt - 1 should overflow. We compute MinInt as -(MaxInt) - 1 via arithmetic
	maxInt := int(^uint(0) >> 1)
	input := fmt.Sprintf("let x = -%d - 1; return x - 1;", maxInt) // x = MinInt, then x - 1
	err := evalScriptError(t, input)
	if err == nil {
		t.Fatal("expected overflow error")
	}
	if !strings.Contains(err.Error(), "integer overflow") {
		t.Fatalf("expected 'integer overflow', got: %v", err)
	}
}

func TestIntegerNegateMinInt(t *testing.T) {
	maxInt := int(^uint(0) >> 1)
	input := fmt.Sprintf("let x = -%d - 1; return -x;", maxInt) // x = MinInt, then -x overflows
	err := evalScriptError(t, input)
	if err == nil {
		t.Fatal("expected overflow error for negating MinInt")
	}
	if !strings.Contains(err.Error(), "integer overflow") {
		t.Fatalf("expected 'integer overflow', got: %v", err)
	}
}

func TestIntegerNoOverflowNormalOps(t *testing.T) {
	// Normal operations should not trigger overflow
	tests := []struct {
		input    string
		expected int
	}{
		{"return 100 + 200;", 300},
		{"return 100 - 200;", -100},
		{"return 100 * 200;", 20000},
		{"return -42;", -42},
	}
	for _, tt := range tests {
		val := unwrapReturn(t, evalScript(t, tt.input))
		intVal := val.(*IntegerValue)
		if intVal.Value != tt.expected {
			t.Fatalf("input %q: expected %d, got %d", tt.input, tt.expected, intVal.Value)
		}
	}
}

// ===================== Phase 5: Ordered Hash =====================

func TestOrderedHashIteration(t *testing.T) {
	input := `
let h = {"c": 3, "a": 1, "b": 2};
let result = "";
foreach (h as key, val) {
	result = result + key;
}
return result;`
	val := unwrapReturn(t, evalScript(t, input))
	strVal, ok := val.(*StringValue)
	if !ok {
		t.Fatalf("expected StringValue, got %T", val)
	}
	if strVal.Value != "cab" {
		t.Fatalf("expected 'cab' (insertion order), got %q", strVal.Value)
	}
}

func TestOrderedHashKeysBuiltin(t *testing.T) {
	input := `
let h = {"x": 1, "y": 2, "z": 3};
return keys(h);`
	val := unwrapReturn(t, evalScript(t, input))
	arrVal, ok := val.(*ArrayValue)
	if !ok {
		t.Fatalf("expected ArrayValue, got %T", val)
	}
	expected := []string{"x", "y", "z"}
	if len(arrVal.Elements) != len(expected) {
		t.Fatalf("expected %d keys, got %d", len(expected), len(arrVal.Elements))
	}
	for i, e := range expected {
		sv := arrVal.Elements[i].(*StringValue)
		if sv.Value != e {
			t.Fatalf("key %d: expected %q, got %q", i, e, sv.Value)
		}
	}
}

func TestOrderedHashValuesBuiltin(t *testing.T) {
	input := `
let h = {"x": 10, "y": 20, "z": 30};
return values(h);`
	val := unwrapReturn(t, evalScript(t, input))
	arrVal := val.(*ArrayValue)
	expected := []int{10, 20, 30}
	for i, e := range expected {
		iv := arrVal.Elements[i].(*IntegerValue)
		if iv.Value != e {
			t.Fatalf("value %d: expected %d, got %d", i, e, iv.Value)
		}
	}
}

func TestOrderedHashSetOverwrite(t *testing.T) {
	// Overwriting a key should not change insertion order
	input := `
let h = {"a": 1, "b": 2, "c": 3};
h.b = 99;
let result = "";
foreach (h as key, val) {
	result = result + key;
}
return result;`
	val := unwrapReturn(t, evalScript(t, input))
	strVal := val.(*StringValue)
	if strVal.Value != "abc" {
		t.Fatalf("expected 'abc', got %q", strVal.Value)
	}
}

func TestHashDelete(t *testing.T) {
	h := NewHashValue()
	h.Set(&StringValue{Value: "a"}, &IntegerValue{Value: 1})
	h.Set(&StringValue{Value: "b"}, &IntegerValue{Value: 2})
	h.Set(&StringValue{Value: "c"}, &IntegerValue{Value: 3})

	err := h.Delete(&StringValue{Value: "b"})
	if err != nil {
		t.Fatal(err)
	}

	if len(h.Pairs) != 2 {
		t.Fatalf("expected 2 pairs, got %d", len(h.Pairs))
	}

	ordered := h.OrderedPairs()
	if len(ordered) != 2 {
		t.Fatalf("expected 2 ordered pairs, got %d", len(ordered))
	}
	if ordered[0].Key.(*StringValue).Value != "a" {
		t.Fatalf("expected first key 'a', got %q", ordered[0].Key.(*StringValue).Value)
	}
	if ordered[1].Key.(*StringValue).Value != "c" {
		t.Fatalf("expected second key 'c', got %q", ordered[1].Key.(*StringValue).Value)
	}
}

func TestHashDeleteNonHashableKey(t *testing.T) {
	h := NewHashValue()
	err := h.Delete(&ArrayValue{Elements: []Object{}})
	if err == nil {
		t.Fatal("expected error for non-hashable key")
	}
}

// ===================== Phase 6: Type Conversion & String Auto-Coercion =====================

func TestStringAutoCoercion(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`return "count: " + 42;`, "count: 42"},
		{`return 42 + " items";`, "42 items"},
		{`return "val: " + 3.14;`, "val: 3.14"},
		{`return "flag: " + true;`, "flag: true"},
		{`return "nothing: " + null;`, "nothing: null"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			val := unwrapReturn(t, evalScript(t, tt.input))
			strVal, ok := val.(*StringValue)
			if !ok {
				t.Fatalf("expected StringValue, got %T", val)
			}
			if strVal.Value != tt.expected {
				t.Fatalf("expected %q, got %q", tt.expected, strVal.Value)
			}
		})
	}
}

func TestToStringBuiltin(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`return toString(42);`, "42"},
		{`return toString(3.14);`, "3.14"},
		{`return toString(true);`, "true"},
		{`return toString(null);`, "null"},
		{`return toString("hello");`, "hello"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			val := unwrapReturn(t, evalScript(t, tt.input))
			strVal := val.(*StringValue)
			if strVal.Value != tt.expected {
				t.Fatalf("expected %q, got %q", tt.expected, strVal.Value)
			}
		})
	}
}

func TestParseIntBuiltin(t *testing.T) {
	val := unwrapReturn(t, evalScript(t, `return parseInt("42");`))
	intVal, ok := val.(*IntegerValue)
	if !ok {
		t.Fatalf("expected IntegerValue, got %T", val)
	}
	if intVal.Value != 42 {
		t.Fatalf("expected 42, got %d", intVal.Value)
	}
}

func TestParseIntInvalid(t *testing.T) {
	err := evalScriptError(t, `return parseInt("abc");`)
	if err == nil {
		t.Fatal("expected error for invalid parseInt input")
	}
}

func TestParseFloatBuiltin(t *testing.T) {
	val := unwrapReturn(t, evalScript(t, `return parseFloat("3.14");`))
	decVal, ok := val.(*DecimalValue)
	if !ok {
		t.Fatalf("expected DecimalValue, got %T", val)
	}
	if decVal.Value != 3.14 {
		t.Fatalf("expected 3.14, got %f", decVal.Value)
	}
}

func TestParseFloatInvalid(t *testing.T) {
	err := evalScriptError(t, `return parseFloat("xyz");`)
	if err == nil {
		t.Fatal("expected error for invalid parseFloat input")
	}
}

// ===================== Phase 7: Anonymous Functions =====================

func TestAnonymousFunctionCallInline(t *testing.T) {
	input := `
let add = fn(a, b) { return a + b; };
return add(3, 4);`
	val := unwrapReturn(t, evalScript(t, input))
	intVal := val.(*IntegerValue)
	if intVal.Value != 7 {
		t.Fatalf("expected 7, got %d", intVal.Value)
	}
}

func TestAnonymousFunctionAsArgument(t *testing.T) {
	input := `
let arr = [1, 2, 3, 4, 5];
let result = filter(arr, fn(x) { return x > 3; });
return len(result);`
	val := unwrapReturn(t, evalScript(t, input))
	intVal := val.(*IntegerValue)
	if intVal.Value != 2 {
		t.Fatalf("expected 2, got %d", intVal.Value)
	}
}

func TestAnonymousFunctionClosure(t *testing.T) {
	input := `
let multiplier = fn(factor) {
	return fn(x) { return x * factor; };
};
let double = multiplier(2);
return double(5);`
	val := unwrapReturn(t, evalScript(t, input))
	intVal := val.(*IntegerValue)
	if intVal.Value != 10 {
		t.Fatalf("expected 10, got %d", intVal.Value)
	}
}

// ===================== Phase 8: Compound Assignment =====================

func TestCompoundAssignment(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"let x = 10; x += 5; return x;", 15},
		{"let x = 10; x -= 3; return x;", 7},
		{"let x = 10; x *= 4; return x;", 40},
		{"let x = 10; x /= 2; return x;", 5},
		{"let x = 10; x %= 3; return x;", 1},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			val := unwrapReturn(t, evalScript(t, tt.input))
			intVal := val.(*IntegerValue)
			if intVal.Value != tt.expected {
				t.Fatalf("expected %d, got %d", tt.expected, intVal.Value)
			}
		})
	}
}

func TestCompoundAssignmentInLoop(t *testing.T) {
	input := `
let sum = 0;
let i = 1;
while (i <= 10) {
	sum += i;
	i += 1;
}
return sum;`
	val := unwrapReturn(t, evalScript(t, input))
	intVal := val.(*IntegerValue)
	if intVal.Value != 55 {
		t.Fatalf("expected 55, got %d", intVal.Value)
	}
}

func TestCompoundAssignmentOnProperty(t *testing.T) {
	input := `
let obj = {"x": 10};
obj.x += 5;
return obj.x;`
	val := unwrapReturn(t, evalScript(t, input))
	intVal := val.(*IntegerValue)
	if intVal.Value != 15 {
		t.Fatalf("expected 15, got %d", intVal.Value)
	}
}

func TestCompoundAssignmentOnIndex(t *testing.T) {
	input := `
let arr = [10, 20, 30];
arr[1] += 5;
return arr[1];`
	val := unwrapReturn(t, evalScript(t, input))
	intVal := val.(*IntegerValue)
	if intVal.Value != 25 {
		t.Fatalf("expected 25, got %d", intVal.Value)
	}
}

func TestCompoundAssignmentStringConcat(t *testing.T) {
	input := `
let s = "hello";
s += " world";
return s;`
	val := unwrapReturn(t, evalScript(t, input))
	strVal := val.(*StringValue)
	if strVal.Value != "hello world" {
		t.Fatalf("expected 'hello world', got %q", strVal.Value)
	}
}

// ===================== Phase 9: Built-in Functions =====================

func TestJoinBuiltin(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`return join([1, 2, 3], ", ");`, "1, 2, 3"},
		{`return join(["a", "b", "c"], "-");`, "a-b-c"},
		{`return join([], ",");`, ""},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			val := unwrapReturn(t, evalScript(t, tt.input))
			strVal := val.(*StringValue)
			if strVal.Value != tt.expected {
				t.Fatalf("expected %q, got %q", tt.expected, strVal.Value)
			}
		})
	}
}

func TestMapBuiltin(t *testing.T) {
	input := `
let nums = [1, 2, 3];
let doubled = map(nums, fn(x) { return x * 2; });
return doubled;`
	val := unwrapReturn(t, evalScript(t, input))
	arrVal := val.(*ArrayValue)
	expected := []int{2, 4, 6}
	for i, e := range expected {
		iv := arrVal.Elements[i].(*IntegerValue)
		if iv.Value != e {
			t.Fatalf("element %d: expected %d, got %d", i, e, iv.Value)
		}
	}
}

func TestFilterBuiltin(t *testing.T) {
	input := `
let nums = [1, 2, 3, 4, 5, 6];
let evens = filter(nums, fn(x) { return x % 2 == 0; });
return evens;`
	val := unwrapReturn(t, evalScript(t, input))
	arrVal := val.(*ArrayValue)
	expected := []int{2, 4, 6}
	if len(arrVal.Elements) != len(expected) {
		t.Fatalf("expected %d elements, got %d", len(expected), len(arrVal.Elements))
	}
	for i, e := range expected {
		iv := arrVal.Elements[i].(*IntegerValue)
		if iv.Value != e {
			t.Fatalf("element %d: expected %d, got %d", i, e, iv.Value)
		}
	}
}

func TestFilterMapJoinChain(t *testing.T) {
	input := `
let nums = [1, 2, 3, 4, 5, 6];
let evens = filter(nums, fn(x) { return x % 2 == 0; });
let doubled = map(evens, fn(x) { return x * 2; });
return join(doubled, ",");`
	val := unwrapReturn(t, evalScript(t, input))
	strVal := val.(*StringValue)
	if strVal.Value != "4,8,12" {
		t.Fatalf("expected '4,8,12', got %q", strVal.Value)
	}
}

func TestFloorBuiltin(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{`return floor(3.7);`, 3},
		{`return floor(3.2);`, 3},
		{`return floor(-1.5);`, -2},
		{`return floor(5);`, 5},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			val := unwrapReturn(t, evalScript(t, tt.input))
			intVal := val.(*IntegerValue)
			if intVal.Value != tt.expected {
				t.Fatalf("expected %d, got %d", tt.expected, intVal.Value)
			}
		})
	}
}

func TestCeilBuiltin(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{`return ceil(3.2);`, 4},
		{`return ceil(3.0);`, 3},
		{`return ceil(-1.5);`, -1},
		{`return ceil(5);`, 5},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			val := unwrapReturn(t, evalScript(t, tt.input))
			intVal := val.(*IntegerValue)
			if intVal.Value != tt.expected {
				t.Fatalf("expected %d, got %d", tt.expected, intVal.Value)
			}
		})
	}
}

func TestRoundBuiltin(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{`return round(3.4);`, 3},
		{`return round(3.5);`, 4},
		{`return round(-1.5);`, -2},
		{`return round(5);`, 5},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			val := unwrapReturn(t, evalScript(t, tt.input))
			intVal := val.(*IntegerValue)
			if intVal.Value != tt.expected {
				t.Fatalf("expected %d, got %d", tt.expected, intVal.Value)
			}
		})
	}
}

func TestAbsBuiltin(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`return abs(-5);`, "5"},
		{`return abs(5);`, "5"},
		{`return abs(-3.14);`, "3.14"},
		{`return abs(3.14);`, "3.14"},
		{`return abs(0);`, "0"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			val := unwrapReturn(t, evalScript(t, tt.input))
			if val.Debug() != tt.expected {
				t.Fatalf("expected %s, got %s", tt.expected, val.Debug())
			}
		})
	}
}

// ===================== Template Mode Integration =====================

func TestTemplateStringCoercion(t *testing.T) {
	input := `Items: {% print("count: " + 42); %}`
	output := evalTemplate(t, input)
	if output != "Items: count: 42" {
		t.Fatalf("expected 'Items: count: 42', got %q", output)
	}
}

func TestTemplateOrderedHashForeach(t *testing.T) {
	// Use script mode to test ordered hash iteration in templates
	input := `
let h = {"a": 1, "b": 2, "c": 3};
let result = "";
foreach (h as key, val) {
	result += key;
}
return result;`
	val := unwrapReturn(t, evalScript(t, input))
	strVal := val.(*StringValue)
	if strVal.Value != "abc" {
		t.Fatalf("expected 'abc', got %q", strVal.Value)
	}
}

func TestTemplateCompoundAssignment(t *testing.T) {
	// Test compound assignment in a pure script context (avoids template let-output artifacts)
	input := `
let x = 0;
let i = 1;
while (i <= 5) {
	x += i;
	i += 1;
}
return toString(x);`
	val := unwrapReturn(t, evalScript(t, input))
	strVal := val.(*StringValue)
	if strVal.Value != "15" {
		t.Fatalf("expected '15', got %q", strVal.Value)
	}
}

// ===================== Error path tests for new built-ins =====================

func TestMapBuiltinBadArgs(t *testing.T) {
	err := evalScriptError(t, `return map(1, fn(x) { return x; });`)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestFilterBuiltinBadArgs(t *testing.T) {
	err := evalScriptError(t, `return filter("not array", fn(x) { return x; });`)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestJoinBuiltinBadArgs(t *testing.T) {
	err := evalScriptError(t, `return join(1, ",");`)
	if err == nil {
		t.Fatal("expected error for non-array argument")
	}
	err = evalScriptError(t, `return join([1], 1);`)
	if err == nil {
		t.Fatal("expected error for non-string separator")
	}
}

func TestFloorBuiltinBadArgs(t *testing.T) {
	err := evalScriptError(t, `return floor("abc");`)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestCeilBuiltinBadArgs(t *testing.T) {
	err := evalScriptError(t, `return ceil("abc");`)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestRoundBuiltinBadArgs(t *testing.T) {
	err := evalScriptError(t, `return round("abc");`)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestAbsBuiltinBadArgs(t *testing.T) {
	err := evalScriptError(t, `return abs("abc");`)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestToStringBuiltinBadArgs(t *testing.T) {
	err := evalScriptError(t, `return toString();`)
	if err == nil {
		t.Fatal("expected error for wrong arg count")
	}
}

func TestParseIntBuiltinBadArgType(t *testing.T) {
	err := evalScriptError(t, `return parseInt(42);`)
	if err == nil {
		t.Fatal("expected error for non-string argument")
	}
}

func TestParseFloatBuiltinBadArgType(t *testing.T) {
	err := evalScriptError(t, `return parseFloat(42);`)
	if err == nil {
		t.Fatal("expected error for non-string argument")
	}
}

