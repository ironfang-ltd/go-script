package evaluator

import (
	"fmt"
	"strings"
	"testing"

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
		{"return 1.5 + 2.5;", "4."},
		{"return 3.0 - 1.5;", "1.5"},
		{"return 2.0 * 3.5;", "7."},
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
	err := evalScriptError(t, `return 1 + "hello";`)
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
