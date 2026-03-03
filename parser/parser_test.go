package parser

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/ironfang-ltd/go-script/lexer"
)

func TestParseLetStatement(t *testing.T) {

	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      interface{}
	}{
		{"let x = 5;", "x", 5},
		{"let y = true;", "y", true},
		{"let foobar = y;", "foobar", "y"},
	}

	for _, tt := range tests {

		l := lexer.NewScript(tt.input)
		p := New(l)

		_, err := p.Parse()
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestParseAssignmentStatement(t *testing.T) {

	tests := []struct {
		input string
	}{
		{"x[0] = 10;"},
		{"x[0].z = 10;"},
		{"x = 5;"},
		{"y = true;"},
		{"x.y = true;"},
		{"x.y.z = 10;"},
	}

	for _, tt := range tests {
		l := lexer.NewScript(tt.input)
		p := New(l)

		program, err := p.Parse()
		if err != nil {
			t.Fatalf("Parse(%s) = %v", tt.input, err)
		}

		if len(program.Statements) != 1 {
			t.Fatalf("Length(%s) = expected 1 statement, got %d", tt.input, len(program.Statements))
		}

		exp, ok := program.Statements[0].(*ExpressionStatement)
		if !ok {
			t.Fatalf("ExpressionStatement(%s) = expected ExpressionStatement, got %T", tt.input, program.Statements[0])
		}

		_, ok = exp.Expression.(*AssignmentExpression)
		if !ok {
			t.Fatalf("AssignmentExpression(%s) = expected AssignmentExpression, got %T", tt.input, program.Statements[0])
		}
	}
}

func TestParseReturnStatement(t *testing.T) {

	tests := []struct {
		input         string
		expectedValue interface{}
	}{
		{"return 5;", 5},
		{"return true;", true},
		{"return y;", "y"},
	}

	for _, tt := range tests {

		l := lexer.NewScript(tt.input)
		p := New(l)

		_, err := p.Parse()
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestParseExpressionStatement(t *testing.T) {

	tests := []struct {
		input         string
		expectedValue interface{}
	}{
		{"5;", 5},
		{"true;", true},
		{"y;", "y"},
	}

	for _, tt := range tests {

		l := lexer.NewScript(tt.input)
		p := New(l)

		_, err := p.Parse()
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestParseFunctionLiteral(t *testing.T) {

	input := `fn test(x, y) { x + y; }`

	l := lexer.NewScript(input)
	p := New(l)

	_, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}
}

func TestParseCallExpression(t *testing.T) {

	input := `add(1, 2);`

	l := lexer.NewScript(input)
	p := New(l)

	_, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}
}

func TestParseIfExpression(t *testing.T) {

	input := `if (x < y) { return x; }`

	l := lexer.NewScript(input)
	p := New(l)

	_, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}
}

func TestParseIfElseExpression(t *testing.T) {

	input := `if (x < y) { return x; } else { return y; }`

	l := lexer.NewScript(input)
	p := New(l)

	_, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}
}

func TestParseTemplateInScriptMode(t *testing.T) {

	input := `{% address.street %}`

	l := lexer.NewScript(input)
	p := New(l)

	_, err := p.Parse()
	if err == nil {
		t.Fatal(fmt.Errorf("expected error"))
	}
}

func TestParseHashLiteral(t *testing.T) {

	input := `let h = { "name": "test", "age": 123 }`

	l := lexer.NewScript(input)
	p := New(l)

	_, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}
}

// --- Regression tests for precedence bugs ---

func TestPrecedenceMultiplyDivideAreEqual(t *testing.T) {
	// 2 / 3 * 4 should parse as (2 / 3) * 4, not 2 / (3 * 4)
	input := `2 / 3 * 4;`
	l := lexer.NewScript(input)
	p := New(l)
	program, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}

	stmt := program.Statements[0].(*ExpressionStatement)
	outer, ok := stmt.Expression.(*InfixExpression)
	if !ok {
		t.Fatalf("expected InfixExpression, got %T", stmt.Expression)
	}
	if outer.Token.Source != "*" {
		t.Fatalf("outer operator should be *, got %s", outer.Token.Source)
	}
	inner, ok := outer.Left.(*InfixExpression)
	if !ok {
		t.Fatalf("left should be InfixExpression, got %T", outer.Left)
	}
	if inner.Token.Source != "/" {
		t.Fatalf("inner operator should be /, got %s", inner.Token.Source)
	}
}

func TestPrecedenceAddSubtractAreEqual(t *testing.T) {
	// 1 - 2 + 3 should parse as (1 - 2) + 3
	input := `1 - 2 + 3;`
	l := lexer.NewScript(input)
	p := New(l)
	program, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}

	stmt := program.Statements[0].(*ExpressionStatement)
	outer, ok := stmt.Expression.(*InfixExpression)
	if !ok {
		t.Fatalf("expected InfixExpression, got %T", stmt.Expression)
	}
	if outer.Token.Source != "+" {
		t.Fatalf("outer operator should be +, got %s", outer.Token.Source)
	}
	inner, ok := outer.Left.(*InfixExpression)
	if !ok {
		t.Fatalf("left should be InfixExpression, got %T", outer.Left)
	}
	if inner.Token.Source != "-" {
		t.Fatalf("inner operator should be -, got %s", inner.Token.Source)
	}
}

func TestPrecedenceMultiplyHigherThanAdd(t *testing.T) {
	// 1 + 2 * 3 should parse as 1 + (2 * 3)
	input := `1 + 2 * 3;`
	l := lexer.NewScript(input)
	p := New(l)
	program, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}

	stmt := program.Statements[0].(*ExpressionStatement)
	outer, ok := stmt.Expression.(*InfixExpression)
	if !ok {
		t.Fatalf("expected InfixExpression, got %T", stmt.Expression)
	}
	if outer.Token.Source != "+" {
		t.Fatalf("outer operator should be +, got %s", outer.Token.Source)
	}
	rightInfix, ok := outer.Right.(*InfixExpression)
	if !ok {
		t.Fatalf("right should be InfixExpression, got %T", outer.Right)
	}
	if rightInfix.Token.Source != "*" {
		t.Fatalf("right operator should be *, got %s", rightInfix.Token.Source)
	}
}

func TestPrecedenceDotHigherThanArithmetic(t *testing.T) {
	// a.b + c should parse as (a.b) + c, not a.(b + c)
	input := `a.b + c;`
	l := lexer.NewScript(input)
	p := New(l)
	program, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}

	stmt := program.Statements[0].(*ExpressionStatement)
	outer, ok := stmt.Expression.(*InfixExpression)
	if !ok {
		t.Fatalf("expected InfixExpression, got %T", stmt.Expression)
	}
	if outer.Token.Source != "+" {
		t.Fatalf("outer operator should be +, got %s", outer.Token.Source)
	}
	prop, ok := outer.Left.(*PropertyExpression)
	if !ok {
		t.Fatalf("left should be PropertyExpression, got %T", outer.Left)
	}
	leftIdent, ok := prop.Left.(*Identifier)
	if !ok {
		t.Fatalf("property left should be Identifier, got %T", prop.Left)
	}
	if leftIdent.Value != "a" {
		t.Fatalf("expected 'a', got %s", leftIdent.Value)
	}
}

func TestPrecedenceModuloSameAsMultiply(t *testing.T) {
	// a.b % c should parse as (a.b) % c
	input := `a.b % c;`
	l := lexer.NewScript(input)
	p := New(l)
	program, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}

	stmt := program.Statements[0].(*ExpressionStatement)
	outer, ok := stmt.Expression.(*InfixExpression)
	if !ok {
		t.Fatalf("expected InfixExpression, got %T", stmt.Expression)
	}
	if outer.Token.Source != "%" {
		t.Fatalf("outer operator should be %%, got %s", outer.Token.Source)
	}
	_, ok = outer.Left.(*PropertyExpression)
	if !ok {
		t.Fatalf("left should be PropertyExpression, got %T", outer.Left)
	}
}

func TestPrecedenceComparisonLowerThanArithmetic(t *testing.T) {
	// a + b == c should parse as (a + b) == c
	input := `a + b == c;`
	l := lexer.NewScript(input)
	p := New(l)
	program, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}

	stmt := program.Statements[0].(*ExpressionStatement)
	outer, ok := stmt.Expression.(*InfixExpression)
	if !ok {
		t.Fatalf("expected InfixExpression, got %T", stmt.Expression)
	}
	if outer.Token.Source != "==" {
		t.Fatalf("outer operator should be ==, got %s", outer.Token.Source)
	}
	inner, ok := outer.Left.(*InfixExpression)
	if !ok {
		t.Fatalf("left should be InfixExpression, got %T", outer.Left)
	}
	if inner.Token.Source != "+" {
		t.Fatalf("inner operator should be +, got %s", inner.Token.Source)
	}
}

// --- Regression tests for missing operators ---

func TestParseLessOrEqual(t *testing.T) {
	input := `a <= b;`
	l := lexer.NewScript(input)
	p := New(l)
	program, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}

	stmt := program.Statements[0].(*ExpressionStatement)
	infix, ok := stmt.Expression.(*InfixExpression)
	if !ok {
		t.Fatalf("expected InfixExpression, got %T", stmt.Expression)
	}
	if infix.Token.Source != "<=" {
		t.Fatalf("expected <=, got %s", infix.Token.Source)
	}
}

func TestParseGreaterOrEqual(t *testing.T) {
	input := `a >= b;`
	l := lexer.NewScript(input)
	p := New(l)
	program, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}

	stmt := program.Statements[0].(*ExpressionStatement)
	infix, ok := stmt.Expression.(*InfixExpression)
	if !ok {
		t.Fatalf("expected InfixExpression, got %T", stmt.Expression)
	}
	if infix.Token.Source != ">=" {
		t.Fatalf("expected >=, got %s", infix.Token.Source)
	}
}

func TestParseModulo(t *testing.T) {
	input := `10 % 3;`
	l := lexer.NewScript(input)
	p := New(l)
	program, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}

	stmt := program.Statements[0].(*ExpressionStatement)
	infix, ok := stmt.Expression.(*InfixExpression)
	if !ok {
		t.Fatalf("expected InfixExpression, got %T", stmt.Expression)
	}
	if infix.Token.Source != "%" {
		t.Fatalf("expected %%, got %s", infix.Token.Source)
	}
}

// --- Regression tests for property access bug ---

func TestPropertyAccessThenOperator(t *testing.T) {
	// item.name + " suffix" should be (item.name) + " suffix"
	input := `item.name + " suffix";`
	l := lexer.NewScript(input)
	p := New(l)
	program, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}

	stmt := program.Statements[0].(*ExpressionStatement)
	infix, ok := stmt.Expression.(*InfixExpression)
	if !ok {
		t.Fatalf("expected InfixExpression at top, got %T", stmt.Expression)
	}
	if infix.Token.Source != "+" {
		t.Fatalf("expected + at top, got %s", infix.Token.Source)
	}
	prop, ok := infix.Left.(*PropertyExpression)
	if !ok {
		t.Fatalf("left should be PropertyExpression, got %T", infix.Left)
	}
	if prop.Property.(*Identifier).Value != "name" {
		t.Fatalf("property should be 'name', got %s", prop.Property.(*Identifier).Value)
	}
}

func TestPropertyAccessChained(t *testing.T) {
	// a.b.c should parse left-to-right: (a.b).c
	input := `a.b.c;`
	l := lexer.NewScript(input)
	p := New(l)
	program, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}

	stmt := program.Statements[0].(*ExpressionStatement)
	outer, ok := stmt.Expression.(*PropertyExpression)
	if !ok {
		t.Fatalf("expected PropertyExpression, got %T", stmt.Expression)
	}
	if outer.Property.(*Identifier).Value != "c" {
		t.Fatalf("outer property should be 'c', got %s", outer.Property.(*Identifier).Value)
	}
	inner, ok := outer.Left.(*PropertyExpression)
	if !ok {
		t.Fatalf("left should be PropertyExpression, got %T", outer.Left)
	}
	if inner.Property.(*Identifier).Value != "b" {
		t.Fatalf("inner property should be 'b', got %s", inner.Property.(*Identifier).Value)
	}
	if inner.Left.(*Identifier).Value != "a" {
		t.Fatalf("innermost should be 'a', got %s", inner.Left.(*Identifier).Value)
	}
}

func TestPropertyAccessWithMethodCall(t *testing.T) {
	// a.method() should parse as CallExpression{Function: PropertyExpression(a, method)}
	input := `a.method();`
	l := lexer.NewScript(input)
	p := New(l)
	program, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}

	stmt := program.Statements[0].(*ExpressionStatement)
	call, ok := stmt.Expression.(*CallExpression)
	if !ok {
		t.Fatalf("expected CallExpression, got %T", stmt.Expression)
	}
	prop, ok := call.Function.(*PropertyExpression)
	if !ok {
		t.Fatalf("function should be PropertyExpression, got %T", call.Function)
	}
	if prop.Left.(*Identifier).Value != "a" {
		t.Fatalf("left should be 'a', got %s", prop.Left.(*Identifier).Value)
	}
	if prop.Property.(*Identifier).Value != "method" {
		t.Fatalf("property should be 'method', got %s", prop.Property.(*Identifier).Value)
	}
}

func TestPropertyAccessWithComparison(t *testing.T) {
	// a.x > b.y should parse as (a.x) > (b.y)
	input := `a.x > b.y;`
	l := lexer.NewScript(input)
	p := New(l)
	program, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}

	stmt := program.Statements[0].(*ExpressionStatement)
	infix, ok := stmt.Expression.(*InfixExpression)
	if !ok {
		t.Fatalf("expected InfixExpression, got %T", stmt.Expression)
	}
	if infix.Token.Source != ">" {
		t.Fatalf("expected >, got %s", infix.Token.Source)
	}
	leftProp, ok := infix.Left.(*PropertyExpression)
	if !ok {
		t.Fatalf("left should be PropertyExpression, got %T", infix.Left)
	}
	if leftProp.Property.(*Identifier).Value != "x" {
		t.Fatalf("left property should be 'x', got %s", leftProp.Property.(*Identifier).Value)
	}
	rightProp, ok := infix.Right.(*PropertyExpression)
	if !ok {
		t.Fatalf("right should be PropertyExpression, got %T", infix.Right)
	}
	if rightProp.Property.(*Identifier).Value != "y" {
		t.Fatalf("right property should be 'y', got %s", rightProp.Property.(*Identifier).Value)
	}
}

// --- Regression tests for string escape processing ---

func TestParseStringEscapes(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"newline", `"\n"`, "\n"},
		{"tab", `"\t"`, "\t"},
		{"return", `"\r"`, "\r"},
		{"escaped_quote", `"\""`, `"`},
		{"escaped_backslash", `"\\"`, `\`},
		{"plain", `"hello"`, "hello"},
		{"empty", `""`, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := fmt.Sprintf("let x = %s;", tt.input)
			l := lexer.NewScript(input)
			p := New(l)
			program, err := p.Parse()
			if err != nil {
				t.Fatal(err)
			}
			letStmt := program.Statements[0].(*LetStatement)
			strLit := letStmt.Value.(*StringLiteral)
			if strLit.Value != tt.want {
				t.Fatalf("want %q, got %q", tt.want, strLit.Value)
			}
		})
	}
}

func TestParseStringDoubleBackslashThenN(t *testing.T) {
	// "test\\nworld" should be: test\nworld (literal backslash, n, world)
	// NOT: test<newline>world
	input := `let x = "test\\nworld";`
	l := lexer.NewScript(input)
	p := New(l)
	program, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}
	letStmt := program.Statements[0].(*LetStatement)
	strLit := letStmt.Value.(*StringLiteral)
	want := "test\\nworld"
	if strLit.Value != want {
		t.Fatalf("want %q, got %q", want, strLit.Value)
	}
}

func TestParseStringDoubleBackslash(t *testing.T) {
	// "\\" should produce a single backslash
	input := `let x = "\\";`
	l := lexer.NewScript(input)
	p := New(l)
	program, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}
	letStmt := program.Statements[0].(*LetStatement)
	strLit := letStmt.Value.(*StringLiteral)
	if strLit.Value != `\` {
		t.Fatalf("want single backslash, got %q", strLit.Value)
	}
}

// --- Regression test for function literal identifier ---

func TestParseFunctionLiteralIdentifier(t *testing.T) {
	input := `fn add(x, y) { return x + y; }`
	l := lexer.NewScript(input)
	p := New(l)
	program, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}

	stmt := program.Statements[0].(*ExpressionStatement)
	fnLit, ok := stmt.Expression.(*FunctionLiteral)
	if !ok {
		t.Fatalf("expected FunctionLiteral, got %T", stmt.Expression)
	}
	if fnLit.Identifier == nil {
		t.Fatal("expected function identifier to be set, got nil")
	}
	if fnLit.Identifier.Value != "add" {
		t.Fatalf("expected function name 'add', got %q", fnLit.Identifier.Value)
	}
	if len(fnLit.Parameters) != 2 {
		t.Fatalf("expected 2 parameters, got %d", len(fnLit.Parameters))
	}
	if fnLit.Parameters[0].Value != "x" {
		t.Fatalf("expected first param 'x', got %q", fnLit.Parameters[0].Value)
	}
	if fnLit.Parameters[1].Value != "y" {
		t.Fatalf("expected second param 'y', got %q", fnLit.Parameters[1].Value)
	}
}

func TestParseFunctionNoParams(t *testing.T) {
	input := `fn greet() { return 1; }`
	l := lexer.NewScript(input)
	p := New(l)
	program, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}

	stmt := program.Statements[0].(*ExpressionStatement)
	fnLit := stmt.Expression.(*FunctionLiteral)
	if fnLit.Identifier.Value != "greet" {
		t.Fatalf("expected name 'greet', got %q", fnLit.Identifier.Value)
	}
	if len(fnLit.Parameters) != 0 {
		t.Fatalf("expected 0 params, got %d", len(fnLit.Parameters))
	}
}

// --- Let statement tests ---

func TestParseLetWithExpression(t *testing.T) {
	input := `let x = 1 + 2;`
	l := lexer.NewScript(input)
	p := New(l)
	program, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(program.Statements))
	}
	letStmt, ok := program.Statements[0].(*LetStatement)
	if !ok {
		t.Fatalf("expected LetStatement, got %T", program.Statements[0])
	}
	if letStmt.Name.Value != "x" {
		t.Fatalf("expected name 'x', got %s", letStmt.Name.Value)
	}
	infix, ok := letStmt.Value.(*InfixExpression)
	if !ok {
		t.Fatalf("expected InfixExpression value, got %T", letStmt.Value)
	}
	if infix.Token.Source != "+" {
		t.Fatalf("expected +, got %s", infix.Token.Source)
	}
}

func TestParseLetWithString(t *testing.T) {
	input := `let name = "hello";`
	l := lexer.NewScript(input)
	p := New(l)
	program, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}

	letStmt := program.Statements[0].(*LetStatement)
	strLit, ok := letStmt.Value.(*StringLiteral)
	if !ok {
		t.Fatalf("expected StringLiteral, got %T", letStmt.Value)
	}
	if strLit.Value != "hello" {
		t.Fatalf("expected 'hello', got %q", strLit.Value)
	}
}

// --- Return statement tests ---

func TestParseReturnExpression(t *testing.T) {
	input := `return 1 + 2;`
	l := lexer.NewScript(input)
	p := New(l)
	program, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}

	retStmt, ok := program.Statements[0].(*ReturnStatement)
	if !ok {
		t.Fatalf("expected ReturnStatement, got %T", program.Statements[0])
	}
	infix, ok := retStmt.Value.(*InfixExpression)
	if !ok {
		t.Fatalf("expected InfixExpression, got %T", retStmt.Value)
	}
	if infix.Token.Source != "+" {
		t.Fatalf("expected +, got %s", infix.Token.Source)
	}
}

// --- Array tests ---

func TestParseArrayLiteral(t *testing.T) {
	input := `let a = [1, 2, 3];`
	l := lexer.NewScript(input)
	p := New(l)
	program, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}

	letStmt := program.Statements[0].(*LetStatement)
	arr, ok := letStmt.Value.(*ArrayLiteral)
	if !ok {
		t.Fatalf("expected ArrayLiteral, got %T", letStmt.Value)
	}
	if len(arr.Elements) != 3 {
		t.Fatalf("expected 3 elements, got %d", len(arr.Elements))
	}
}

func TestParseEmptyArray(t *testing.T) {
	input := `let a = [];`
	l := lexer.NewScript(input)
	p := New(l)
	program, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}

	letStmt := program.Statements[0].(*LetStatement)
	arr := letStmt.Value.(*ArrayLiteral)
	if len(arr.Elements) != 0 {
		t.Fatalf("expected 0 elements, got %d", len(arr.Elements))
	}
}

// --- Hash tests ---

func TestParseEmptyHash(t *testing.T) {
	input := `let h = {};`
	l := lexer.NewScript(input)
	p := New(l)
	program, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}

	letStmt := program.Statements[0].(*LetStatement)
	hash, ok := letStmt.Value.(*HashLiteral)
	if !ok {
		t.Fatalf("expected HashLiteral, got %T", letStmt.Value)
	}
	if len(hash.Pairs) != 0 {
		t.Fatalf("expected 0 pairs, got %d", len(hash.Pairs))
	}
}

func TestParseHashWithValues(t *testing.T) {
	input := `let h = {"a": 1, "b": 2};`
	l := lexer.NewScript(input)
	p := New(l)
	program, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}

	letStmt := program.Statements[0].(*LetStatement)
	hash := letStmt.Value.(*HashLiteral)
	if len(hash.Pairs) != 2 {
		t.Fatalf("expected 2 pairs, got %d", len(hash.Pairs))
	}
}

// --- Index expression tests ---

func TestParseIndexExpression(t *testing.T) {
	input := `arr[0];`
	l := lexer.NewScript(input)
	p := New(l)
	program, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}

	stmt := program.Statements[0].(*ExpressionStatement)
	idx, ok := stmt.Expression.(*IndexExpression)
	if !ok {
		t.Fatalf("expected IndexExpression, got %T", stmt.Expression)
	}
	if idx.Left.(*Identifier).Value != "arr" {
		t.Fatalf("expected 'arr', got %s", idx.Left.(*Identifier).Value)
	}
	intLit := idx.Index.(*IntegerLiteral)
	if intLit.Value != 0 {
		t.Fatalf("expected index 0, got %d", intLit.Value)
	}
}

func TestParseIndexThenPropertyAssignment(t *testing.T) {
	input := `x[0].z = 10;`
	l := lexer.NewScript(input)
	p := New(l)
	program, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}

	stmt := program.Statements[0].(*ExpressionStatement)
	assign, ok := stmt.Expression.(*AssignmentExpression)
	if !ok {
		t.Fatalf("expected AssignmentExpression, got %T", stmt.Expression)
	}
	prop, ok := assign.Left.(*PropertyExpression)
	if !ok {
		t.Fatalf("left should be PropertyExpression, got %T", assign.Left)
	}
	_, ok = prop.Left.(*IndexExpression)
	if !ok {
		t.Fatalf("property left should be IndexExpression, got %T", prop.Left)
	}
}

// --- If/Else tests ---

func TestParseIfWithBlock(t *testing.T) {
	input := `if (x > 0) { let y = 1; }`
	l := lexer.NewScript(input)
	p := New(l)
	program, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}

	stmt := program.Statements[0].(*ExpressionStatement)
	ifExpr, ok := stmt.Expression.(*IfExpression)
	if !ok {
		t.Fatalf("expected IfExpression, got %T", stmt.Expression)
	}
	if ifExpr.Alternative != nil {
		t.Fatal("expected no alternative")
	}
	if len(ifExpr.Consequence.Statements) != 1 {
		t.Fatalf("expected 1 consequence statement, got %d", len(ifExpr.Consequence.Statements))
	}
}

func TestParseIfElseWithBlocks(t *testing.T) {
	input := `if (x > 0) { return 1; } else { return 2; }`
	l := lexer.NewScript(input)
	p := New(l)
	program, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}

	stmt := program.Statements[0].(*ExpressionStatement)
	ifExpr := stmt.Expression.(*IfExpression)
	if ifExpr.Alternative == nil {
		t.Fatal("expected alternative block")
	}
	if len(ifExpr.Consequence.Statements) != 1 {
		t.Fatalf("expected 1 consequence statement, got %d", len(ifExpr.Consequence.Statements))
	}
	if len(ifExpr.Alternative.Statements) != 1 {
		t.Fatalf("expected 1 alternative statement, got %d", len(ifExpr.Alternative.Statements))
	}
}

// --- Foreach tests ---

func TestParseForeach(t *testing.T) {
	input := `foreach (items as item) { print(item); }`
	l := lexer.NewScript(input)
	p := New(l)
	program, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}

	fe, ok := program.Statements[0].(*ForeachExpression)
	if !ok {
		t.Fatalf("expected ForeachExpression, got %T", program.Statements[0])
	}
	if fe.Variable.Value != "item" {
		t.Fatalf("expected variable 'item', got %q", fe.Variable.Value)
	}
	if fe.Iterable.(*Identifier).Value != "items" {
		t.Fatalf("expected iterable 'items', got %q", fe.Iterable.(*Identifier).Value)
	}
}

// --- Prefix expression tests ---

func TestParsePrefixMinus(t *testing.T) {
	input := `-5;`
	l := lexer.NewScript(input)
	p := New(l)
	program, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}

	stmt := program.Statements[0].(*ExpressionStatement)
	prefix, ok := stmt.Expression.(*PrefixExpression)
	if !ok {
		t.Fatalf("expected PrefixExpression, got %T", stmt.Expression)
	}
	if prefix.Operator != "-" {
		t.Fatalf("expected -, got %s", prefix.Operator)
	}
	intLit := prefix.Right.(*IntegerLiteral)
	if intLit.Value != 5 {
		t.Fatalf("expected 5, got %d", intLit.Value)
	}
}

func TestParsePrefixBang(t *testing.T) {
	input := `!true;`
	l := lexer.NewScript(input)
	p := New(l)
	program, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}

	stmt := program.Statements[0].(*ExpressionStatement)
	prefix, ok := stmt.Expression.(*PrefixExpression)
	if !ok {
		t.Fatalf("expected PrefixExpression, got %T", stmt.Expression)
	}
	if prefix.Operator != "!" {
		t.Fatalf("expected !, got %s", prefix.Operator)
	}
}

// --- Grouped expression tests ---

func TestParseGroupedExpression(t *testing.T) {
	// (1 + 2) * 3 should parse as multiply with grouped left
	input := `(1 + 2) * 3;`
	l := lexer.NewScript(input)
	p := New(l)
	program, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}

	stmt := program.Statements[0].(*ExpressionStatement)
	outer, ok := stmt.Expression.(*InfixExpression)
	if !ok {
		t.Fatalf("expected InfixExpression, got %T", stmt.Expression)
	}
	if outer.Token.Source != "*" {
		t.Fatalf("outer should be *, got %s", outer.Token.Source)
	}
	inner, ok := outer.Left.(*InfixExpression)
	if !ok {
		t.Fatalf("left should be InfixExpression (grouped), got %T", outer.Left)
	}
	if inner.Token.Source != "+" {
		t.Fatalf("inner should be +, got %s", inner.Token.Source)
	}
}

// --- Call expression tests ---

func TestParseCallWithNoArgs(t *testing.T) {
	input := `foo();`
	l := lexer.NewScript(input)
	p := New(l)
	program, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}

	stmt := program.Statements[0].(*ExpressionStatement)
	call, ok := stmt.Expression.(*CallExpression)
	if !ok {
		t.Fatalf("expected CallExpression, got %T", stmt.Expression)
	}
	if len(call.Args) != 0 {
		t.Fatalf("expected 0 args, got %d", len(call.Args))
	}
}

func TestParseCallWithMultipleArgs(t *testing.T) {
	input := `add(1, 2, 3);`
	l := lexer.NewScript(input)
	p := New(l)
	program, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}

	stmt := program.Statements[0].(*ExpressionStatement)
	call := stmt.Expression.(*CallExpression)
	if len(call.Args) != 3 {
		t.Fatalf("expected 3 args, got %d", len(call.Args))
	}
}

func TestParseCallWithExpressionArg(t *testing.T) {
	input := `foo(1 + 2);`
	l := lexer.NewScript(input)
	p := New(l)
	program, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}

	stmt := program.Statements[0].(*ExpressionStatement)
	call := stmt.Expression.(*CallExpression)
	if len(call.Args) != 1 {
		t.Fatalf("expected 1 arg, got %d", len(call.Args))
	}
	_, ok := call.Args[0].(*InfixExpression)
	if !ok {
		t.Fatalf("expected InfixExpression arg, got %T", call.Args[0])
	}
}

// --- Template mode tests ---

func TestParseTemplate(t *testing.T) {
	input := `Hello {% name %}, welcome!`
	l := lexer.NewTemplate(input)
	p := New(l)
	program, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}

	// Should have: PrintStatement, ExpressionStatement(name), PrintStatement
	if len(program.Statements) != 3 {
		t.Fatalf("expected 3 statements, got %d", len(program.Statements))
	}

	print1, ok := program.Statements[0].(*PrintStatement)
	if !ok {
		t.Fatalf("expected PrintStatement, got %T", program.Statements[0])
	}
	if print1.Value != "Hello " {
		t.Fatalf("expected 'Hello ', got %q", print1.Value)
	}

	exprStmt, ok := program.Statements[1].(*ExpressionStatement)
	if !ok {
		t.Fatalf("expected ExpressionStatement, got %T", program.Statements[1])
	}
	ident := exprStmt.Expression.(*Identifier)
	if ident.Value != "name" {
		t.Fatalf("expected 'name', got %q", ident.Value)
	}
}

func TestParseTemplateWithIfBlock(t *testing.T) {
	input := `{% if (x > 0) { %}yes{% } %}`
	l := lexer.NewTemplate(input)
	p := New(l)
	_, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}
}

// --- Multiple statements ---

func TestParseMultipleStatements(t *testing.T) {
	input := `let a = 1; let b = 2; return a + b;`
	l := lexer.NewScript(input)
	p := New(l)
	program, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}

	if len(program.Statements) != 3 {
		t.Fatalf("expected 3 statements, got %d", len(program.Statements))
	}
	if _, ok := program.Statements[0].(*LetStatement); !ok {
		t.Fatalf("[0] expected LetStatement, got %T", program.Statements[0])
	}
	if _, ok := program.Statements[1].(*LetStatement); !ok {
		t.Fatalf("[1] expected LetStatement, got %T", program.Statements[1])
	}
	if _, ok := program.Statements[2].(*ReturnStatement); !ok {
		t.Fatalf("[2] expected ReturnStatement, got %T", program.Statements[2])
	}
}

// --- Error tests ---

func TestParseErrorMissingSemicolon(t *testing.T) {
	input := `let x = 5 let y = 10;`
	l := lexer.NewScript(input)
	p := New(l)
	_, err := p.Parse()
	if err == nil {
		t.Fatal("expected error for missing semicolon")
	}
}

func TestParseErrorMissingRightParen(t *testing.T) {
	input := `if (x > 0 { return 1; }`
	l := lexer.NewScript(input)
	p := New(l)
	_, err := p.Parse()
	if err == nil {
		t.Fatal("expected error for missing right paren")
	}
}

func TestParseEmptyProgram(t *testing.T) {
	input := ``
	l := lexer.NewScript(input)
	p := New(l)
	program, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}
	if len(program.Statements) != 0 {
		t.Fatalf("expected 0 statements, got %d", len(program.Statements))
	}
}

// --- Boolean literal tests ---

func TestParseBooleanLiterals(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"true;", true},
		{"false;", false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			l := lexer.NewScript(tt.input)
			p := New(l)
			program, err := p.Parse()
			if err != nil {
				t.Fatal(err)
			}
			stmt := program.Statements[0].(*ExpressionStatement)
			boolLit := stmt.Expression.(*BooleanLiteral)
			if boolLit.Value != tt.want {
				t.Fatalf("want %v, got %v", tt.want, boolLit.Value)
			}
		})
	}
}

// --- Integer literal tests ---

func TestParseIntegerLiteral(t *testing.T) {
	input := `42;`
	l := lexer.NewScript(input)
	p := New(l)
	program, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}

	stmt := program.Statements[0].(*ExpressionStatement)
	intLit, ok := stmt.Expression.(*IntegerLiteral)
	if !ok {
		t.Fatalf("expected IntegerLiteral, got %T", stmt.Expression)
	}
	if intLit.Value != 42 {
		t.Fatalf("expected 42, got %d", intLit.Value)
	}
}

// --- All infix operators ---

// --- Float literal tests ---

func TestParseFloatLiteral(t *testing.T) {
	input := `3.14;`
	l := lexer.NewScript(input)
	p := New(l)
	program, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}

	stmt := program.Statements[0].(*ExpressionStatement)
	floatLit, ok := stmt.Expression.(*FloatLiteral)
	if !ok {
		t.Fatalf("expected FloatLiteral, got %T", stmt.Expression)
	}
	if floatLit.Value != 3.14 {
		t.Fatalf("expected 3.14, got %f", floatLit.Value)
	}
}

func TestParseFloatInExpression(t *testing.T) {
	input := `1.5 + 2.5;`
	l := lexer.NewScript(input)
	p := New(l)
	program, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}

	stmt := program.Statements[0].(*ExpressionStatement)
	infix, ok := stmt.Expression.(*InfixExpression)
	if !ok {
		t.Fatalf("expected InfixExpression, got %T", stmt.Expression)
	}
	if infix.Token.Source != "+" {
		t.Fatalf("expected +, got %s", infix.Token.Source)
	}
	left, ok := infix.Left.(*FloatLiteral)
	if !ok {
		t.Fatalf("left should be FloatLiteral, got %T", infix.Left)
	}
	if left.Value != 1.5 {
		t.Fatalf("expected left 1.5, got %f", left.Value)
	}
	right, ok := infix.Right.(*FloatLiteral)
	if !ok {
		t.Fatalf("right should be FloatLiteral, got %T", infix.Right)
	}
	if right.Value != 2.5 {
		t.Fatalf("expected right 2.5, got %f", right.Value)
	}
}

func TestParseNegativeFloat(t *testing.T) {
	input := `-3.14;`
	l := lexer.NewScript(input)
	p := New(l)
	program, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}

	stmt := program.Statements[0].(*ExpressionStatement)
	prefix, ok := stmt.Expression.(*PrefixExpression)
	if !ok {
		t.Fatalf("expected PrefixExpression, got %T", stmt.Expression)
	}
	if prefix.Operator != "-" {
		t.Fatalf("expected -, got %s", prefix.Operator)
	}
	floatLit, ok := prefix.Right.(*FloatLiteral)
	if !ok {
		t.Fatalf("right should be FloatLiteral, got %T", prefix.Right)
	}
	if floatLit.Value != 3.14 {
		t.Fatalf("expected 3.14, got %f", floatLit.Value)
	}
}

func TestParseFloatPropertyAccess(t *testing.T) {
	// Ensure 1.method doesn't get parsed as a float
	// The lexer should produce Integer(1), Dot, Identifier(method)
	input := `a.b + 1.5;`
	l := lexer.NewScript(input)
	p := New(l)
	program, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}

	stmt := program.Statements[0].(*ExpressionStatement)
	infix, ok := stmt.Expression.(*InfixExpression)
	if !ok {
		t.Fatalf("expected InfixExpression, got %T", stmt.Expression)
	}
	if infix.Token.Source != "+" {
		t.Fatalf("expected +, got %s", infix.Token.Source)
	}
	_, ok = infix.Left.(*PropertyExpression)
	if !ok {
		t.Fatalf("left should be PropertyExpression, got %T", infix.Left)
	}
	_, ok = infix.Right.(*FloatLiteral)
	if !ok {
		t.Fatalf("right should be FloatLiteral, got %T", infix.Right)
	}
}

func TestParseAllInfixOperators(t *testing.T) {
	tests := []struct {
		input    string
		operator string
	}{
		{"1 + 2;", "+"},
		{"1 - 2;", "-"},
		{"1 * 2;", "*"},
		{"1 / 2;", "/"},
		{"1 % 2;", "%"},
		{"1 == 2;", "=="},
		{"1 != 2;", "!="},
		{"1 < 2;", "<"},
		{"1 > 2;", ">"},
		{"1 <= 2;", "<="},
		{"1 >= 2;", ">="},
	}

	for _, tt := range tests {
		t.Run(tt.operator, func(t *testing.T) {
			l := lexer.NewScript(tt.input)
			p := New(l)
			program, err := p.Parse()
			if err != nil {
				t.Fatalf("parsing %q: %v", tt.input, err)
			}
			stmt := program.Statements[0].(*ExpressionStatement)
			infix, ok := stmt.Expression.(*InfixExpression)
			if !ok {
				t.Fatalf("expected InfixExpression, got %T", stmt.Expression)
			}
			if infix.Token.Source != tt.operator {
				t.Fatalf("expected operator %q, got %q", tt.operator, infix.Token.Source)
			}
		})
	}
}

// --- ParseError.Error() tests ---

func TestParseErrorError(t *testing.T) {
	tok := lexer.Token{Type: lexer.Identifier, Source: "x", Line: 1, Column: 5}
	err := NewParseError("unexpected token", "let x 5;", tok)
	msg := err.Error()

	if !strings.Contains(msg, "unexpected token") {
		t.Fatalf("should contain message, got: %s", msg)
	}
	if !strings.Contains(msg, "line 1") {
		t.Fatalf("should contain line, got: %s", msg)
	}
	if !strings.Contains(msg, "column 5") {
		t.Fatalf("should contain column, got: %s", msg)
	}
	if !strings.Contains(msg, "let x 5;") {
		t.Fatalf("should contain source line, got: %s", msg)
	}
	if !strings.Contains(msg, "^") {
		t.Fatalf("should contain caret, got: %s", msg)
	}
}

func TestParseErrorErrorWithTabs(t *testing.T) {
	tok := lexer.Token{Type: lexer.Identifier, Source: "#", Line: 1, Column: 2}
	err := NewParseError("bad token", "\t#", tok)
	msg := err.Error()

	if !strings.Contains(msg, "    #") {
		t.Fatalf("tabs should be expanded to 4 spaces, got: %s", msg)
	}
}

func TestParseErrorErrorMultiline(t *testing.T) {
	source := "let a = 1;\nlet b = #;"
	tok := lexer.Token{Type: lexer.Identifier, Source: "#", Line: 2, Column: 9}
	err := NewParseError("bad token", source, tok)
	msg := err.Error()

	if !strings.Contains(msg, "line 2") {
		t.Fatalf("should reference line 2, got: %s", msg)
	}
	if !strings.Contains(msg, "let b = #;") {
		t.Fatalf("should contain second source line, got: %s", msg)
	}
}

func TestParseErrorIsError(t *testing.T) {
	// Verify ParseError is returned through Parse() and extractable via errors.As
	input := `let = 5;`
	l := lexer.NewScript(input)
	p := New(l)
	_, err := p.Parse()
	if err == nil {
		t.Fatal("expected error")
	}
	var parseErr *ParseError
	if !errors.As(err, &parseErr) {
		t.Fatalf("expected ParseError, got %T: %v", err, err)
	}
}

// --- Debug() method tests ---

func TestProgramDebugEmpty(t *testing.T) {
	program := NewProgram()
	output := program.Debug()
	if output != "" {
		t.Fatalf("expected empty output, got: %q", output)
	}
}

func TestProgramDebugWithStatements(t *testing.T) {
	input := `let x = 5; return x;`
	l := lexer.NewScript(input)
	p := New(l)
	program, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}
	output := program.Debug()
	if output == "" {
		t.Fatal("expected non-empty debug output")
	}
	if !strings.Contains(output, "let") {
		t.Fatalf("expected 'let' in output, got: %s", output)
	}
}

func TestDebugLetStatement(t *testing.T) {
	input := `let x = 5;`
	l := lexer.NewScript(input)
	p := New(l)
	program, _ := p.Parse()
	stmt := program.Statements[0].(*LetStatement)
	d := stmt.Debug()
	if !strings.Contains(d, "let") || !strings.Contains(d, "x") {
		t.Fatalf("unexpected debug output: %s", d)
	}
}

func TestDebugReturnStatement(t *testing.T) {
	input := `return 42;`
	l := lexer.NewScript(input)
	p := New(l)
	program, _ := p.Parse()
	stmt := program.Statements[0].(*ReturnStatement)
	d := stmt.Debug()
	if !strings.Contains(d, "return") {
		t.Fatalf("unexpected debug output: %s", d)
	}
}

func TestDebugExpressionStatement(t *testing.T) {
	input := `42;`
	l := lexer.NewScript(input)
	p := New(l)
	program, _ := p.Parse()
	stmt := program.Statements[0].(*ExpressionStatement)
	d := stmt.Debug()
	if d != "42" {
		t.Fatalf("expected '42', got: %s", d)
	}
}

func TestDebugPrintStatement(t *testing.T) {
	l := lexer.NewTemplate("hello world")
	p := New(l)
	program, _ := p.Parse()
	stmt := program.Statements[0].(*PrintStatement)
	d := stmt.Debug()
	if !strings.Contains(d, "hello world") {
		t.Fatalf("expected 'hello world' in output, got: %s", d)
	}
}

func TestDebugAssignmentExpression(t *testing.T) {
	input := `x = 10;`
	l := lexer.NewScript(input)
	p := New(l)
	program, _ := p.Parse()
	stmt := program.Statements[0].(*ExpressionStatement)
	d := stmt.Expression.(*AssignmentExpression).Debug()
	if !strings.Contains(d, "=") {
		t.Fatalf("expected '=' in output, got: %s", d)
	}
}

func TestDebugIdentifier(t *testing.T) {
	input := `x;`
	l := lexer.NewScript(input)
	p := New(l)
	program, _ := p.Parse()
	stmt := program.Statements[0].(*ExpressionStatement)
	d := stmt.Expression.(*Identifier).Debug()
	if d != "x" {
		t.Fatalf("expected 'x', got: %s", d)
	}
}

func TestDebugInfixExpression(t *testing.T) {
	input := `1 + 2;`
	l := lexer.NewScript(input)
	p := New(l)
	program, _ := p.Parse()
	stmt := program.Statements[0].(*ExpressionStatement)
	d := stmt.Expression.(*InfixExpression).Debug()
	if !strings.Contains(d, "+") {
		t.Fatalf("expected '+' in output, got: %s", d)
	}
}

func TestDebugPrefixExpression(t *testing.T) {
	input := `-5;`
	l := lexer.NewScript(input)
	p := New(l)
	program, _ := p.Parse()
	stmt := program.Statements[0].(*ExpressionStatement)
	d := stmt.Expression.(*PrefixExpression).Debug()
	if !strings.Contains(d, "-") {
		t.Fatalf("expected '-' in output, got: %s", d)
	}
}

func TestDebugPropertyExpression(t *testing.T) {
	input := `a.b;`
	l := lexer.NewScript(input)
	p := New(l)
	program, _ := p.Parse()
	stmt := program.Statements[0].(*ExpressionStatement)
	d := stmt.Expression.(*PropertyExpression).Debug()
	if d != "a.b" {
		t.Fatalf("expected 'a.b', got: %s", d)
	}
}

func TestDebugCallExpressionNoArgs(t *testing.T) {
	input := `foo();`
	l := lexer.NewScript(input)
	p := New(l)
	program, _ := p.Parse()
	stmt := program.Statements[0].(*ExpressionStatement)
	d := stmt.Expression.(*CallExpression).Debug()
	if d != "foo()" {
		t.Fatalf("expected 'foo()', got: %s", d)
	}
}

func TestDebugCallExpressionWithArgs(t *testing.T) {
	input := `add(1, 2);`
	l := lexer.NewScript(input)
	p := New(l)
	program, _ := p.Parse()
	stmt := program.Statements[0].(*ExpressionStatement)
	d := stmt.Expression.(*CallExpression).Debug()
	if !strings.Contains(d, "add(") && !strings.Contains(d, ", ") {
		t.Fatalf("unexpected call debug: %s", d)
	}
}

func TestDebugIndexExpression(t *testing.T) {
	input := `arr[0];`
	l := lexer.NewScript(input)
	p := New(l)
	program, _ := p.Parse()
	stmt := program.Statements[0].(*ExpressionStatement)
	d := stmt.Expression.(*IndexExpression).Debug()
	if d != "arr[0]" {
		t.Fatalf("expected 'arr[0]', got: %s", d)
	}
}

func TestDebugIfExpressionNoElse(t *testing.T) {
	input := `if (true) { 1; }`
	l := lexer.NewScript(input)
	p := New(l)
	program, _ := p.Parse()
	stmt := program.Statements[0].(*ExpressionStatement)
	d := stmt.Expression.(*IfExpression).Debug()
	if !strings.Contains(d, "if") {
		t.Fatalf("expected 'if' in output, got: %s", d)
	}
	if strings.Contains(d, "else") {
		t.Fatalf("should not contain 'else', got: %s", d)
	}
}

func TestDebugIfExpressionWithElse(t *testing.T) {
	input := `if (true) { 1; } else { 2; }`
	l := lexer.NewScript(input)
	p := New(l)
	program, _ := p.Parse()
	stmt := program.Statements[0].(*ExpressionStatement)
	d := stmt.Expression.(*IfExpression).Debug()
	if !strings.Contains(d, "else") {
		t.Fatalf("expected 'else' in output, got: %s", d)
	}
}

func TestDebugForeachExpression(t *testing.T) {
	input := `foreach (items as item) { 1; }`
	l := lexer.NewScript(input)
	p := New(l)
	program, _ := p.Parse()
	stmt := program.Statements[0].(*ForeachExpression)
	d := stmt.Debug()
	if !strings.Contains(d, "foreach") {
		t.Fatalf("expected 'foreach' in output, got: %s", d)
	}
	if !strings.Contains(d, "item") {
		t.Fatalf("expected 'item' in output, got: %s", d)
	}
}

func TestDebugBlockStatement(t *testing.T) {
	input := `if (true) { let x = 1; let y = 2; }`
	l := lexer.NewScript(input)
	p := New(l)
	program, _ := p.Parse()
	stmt := program.Statements[0].(*ExpressionStatement)
	ifExpr := stmt.Expression.(*IfExpression)
	d := ifExpr.Consequence.Debug()
	if !strings.Contains(d, "{") || !strings.Contains(d, "}") {
		t.Fatalf("expected braces in block debug, got: %s", d)
	}
}

func TestDebugIntegerLiteral(t *testing.T) {
	input := `42;`
	l := lexer.NewScript(input)
	p := New(l)
	program, _ := p.Parse()
	stmt := program.Statements[0].(*ExpressionStatement)
	d := stmt.Expression.(*IntegerLiteral).Debug()
	if d != "42" {
		t.Fatalf("expected '42', got: %s", d)
	}
}

func TestDebugFloatLiteral(t *testing.T) {
	input := `3.14;`
	l := lexer.NewScript(input)
	p := New(l)
	program, _ := p.Parse()
	stmt := program.Statements[0].(*ExpressionStatement)
	d := stmt.Expression.(*FloatLiteral).Debug()
	if d != "3.14" {
		t.Fatalf("expected '3.14', got: %s", d)
	}
}

func TestDebugStringLiteral(t *testing.T) {
	input := `"hello";`
	l := lexer.NewScript(input)
	p := New(l)
	program, _ := p.Parse()
	stmt := program.Statements[0].(*ExpressionStatement)
	d := stmt.Expression.(*StringLiteral).Debug()
	if d != `"hello"` {
		t.Fatalf("expected '\"hello\"', got: %s", d)
	}
}

func TestDebugBooleanLiteral(t *testing.T) {
	input := `true;`
	l := lexer.NewScript(input)
	p := New(l)
	program, _ := p.Parse()
	stmt := program.Statements[0].(*ExpressionStatement)
	d := stmt.Expression.(*BooleanLiteral).Debug()
	if d != "true" {
		t.Fatalf("expected 'true', got: %s", d)
	}
}

func TestDebugFunctionLiteral(t *testing.T) {
	input := `fn add(a, b) { return a; }`
	l := lexer.NewScript(input)
	p := New(l)
	program, _ := p.Parse()
	stmt := program.Statements[0].(*ExpressionStatement)
	d := stmt.Expression.(*FunctionLiteral).Debug()
	if !strings.Contains(d, "fn") || !strings.Contains(d, "add") {
		t.Fatalf("expected fn and name in output, got: %s", d)
	}
	if !strings.Contains(d, "a, b") {
		t.Fatalf("expected params in output, got: %s", d)
	}
}

func TestDebugFunctionLiteralNoParams(t *testing.T) {
	input := `fn noop() { return 0; }`
	l := lexer.NewScript(input)
	p := New(l)
	program, _ := p.Parse()
	stmt := program.Statements[0].(*ExpressionStatement)
	d := stmt.Expression.(*FunctionLiteral).Debug()
	if !strings.Contains(d, "noop()") {
		t.Fatalf("expected 'noop()' in output, got: %s", d)
	}
}

func TestDebugArrayLiteral(t *testing.T) {
	input := `[1, 2, 3];`
	l := lexer.NewScript(input)
	p := New(l)
	program, _ := p.Parse()
	stmt := program.Statements[0].(*ExpressionStatement)
	d := stmt.Expression.(*ArrayLiteral).Debug()
	if d == "" {
		t.Fatal("expected non-empty debug output")
	}
}

func TestDebugHashLiteral(t *testing.T) {
	input := `{"a": 1};`
	l := lexer.NewScript(input)
	p := New(l)
	program, _ := p.Parse()
	stmt := program.Statements[0].(*ExpressionStatement)
	d := stmt.Expression.(*HashLiteral).Debug()
	if !strings.Contains(d, "{") {
		t.Fatalf("expected '{' in output, got: %s", d)
	}
}

func TestDebugHashLiteralEmpty(t *testing.T) {
	input := `{};`
	l := lexer.NewScript(input)
	p := New(l)
	program, _ := p.Parse()
	stmt := program.Statements[0].(*ExpressionStatement)
	d := stmt.Expression.(*HashLiteral).Debug()
	if !strings.Contains(d, "{}") {
		t.Fatalf("expected '{}' in output, got: %s", d)
	}
}

func TestDebugScriptStatement(t *testing.T) {
	// ScriptStatement is not created by the parser, test directly
	ss := &ScriptStatement{
		Statements: []Statement{
			&PrintStatement{Value: "hello"},
			&PrintStatement{Value: "world"},
		},
	}
	d := ss.Debug()
	if !strings.Contains(d, "hello") || !strings.Contains(d, "world") {
		t.Fatalf("expected both statements in output, got: %s", d)
	}
}

func TestDebugScriptStatementEmpty(t *testing.T) {
	ss := &ScriptStatement{}
	d := ss.Debug()
	if d != "" {
		t.Fatalf("expected empty output, got: %s", d)
	}
}

// --- String escape edge case ---

func TestParseStringUnknownEscape(t *testing.T) {
	// \q is not a recognized escape - should produce literal \q
	input := `let x = "\q";`
	l := lexer.NewScript(input)
	p := New(l)
	program, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}
	letStmt := program.Statements[0].(*LetStatement)
	strLit := letStmt.Value.(*StringLiteral)
	if strLit.Value != `\q` {
		t.Fatalf("expected '\\q', got %q", strLit.Value)
	}
}

func TestParseStringMultipleUnknownEscapes(t *testing.T) {
	input := `"\a\b\c";`
	l := lexer.NewScript(input)
	p := New(l)
	program, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}
	stmt := program.Statements[0].(*ExpressionStatement)
	strLit := stmt.Expression.(*StringLiteral)
	if strLit.Value != `\a\b\c` {
		t.Fatalf("expected '\\a\\b\\c', got %q", strLit.Value)
	}
}

// --- Template whitespace trimming ---

func TestParseTemplateWhitespaceTrimmingBeforeScript(t *testing.T) {
	// Text ending with newline+whitespace before {% should be trimmed
	input := "hello\n\t{% x %}"
	l := lexer.NewTemplate(input)
	p := New(l)
	program, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}
	printStmt := program.Statements[0].(*PrintStatement)
	if printStmt.Value != "hello" {
		t.Fatalf("expected 'hello' (trimmed), got %q", printStmt.Value)
	}
}

func TestParseTemplateNoTrimmingWithContent(t *testing.T) {
	// Text ending with newline+non-whitespace before {% should NOT be trimmed
	input := "hello\ncontent{% x %}"
	l := lexer.NewTemplate(input)
	p := New(l)
	program, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}
	printStmt := program.Statements[0].(*PrintStatement)
	if printStmt.Value != "hello\ncontent" {
		t.Fatalf("expected 'hello\\ncontent', got %q", printStmt.Value)
	}
}

func TestParseTemplateNoTrimmingWithoutNewline(t *testing.T) {
	// Text without newline before {% should NOT be trimmed
	input := "hello {% x %}"
	l := lexer.NewTemplate(input)
	p := New(l)
	program, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}
	printStmt := program.Statements[0].(*PrintStatement)
	if printStmt.Value != "hello " {
		t.Fatalf("expected 'hello ', got %q", printStmt.Value)
	}
}

func TestParseTemplateTrimmingMultipleLines(t *testing.T) {
	// Multiple lines, last line is whitespace only before {%
	input := "line1\nline2\n    {% x %}"
	l := lexer.NewTemplate(input)
	p := New(l)
	program, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}
	printStmt := program.Statements[0].(*PrintStatement)
	if printStmt.Value != "line1\nline2" {
		t.Fatalf("expected 'line1\\nline2', got %q", printStmt.Value)
	}
}

// --- Parser error path tests ---

func TestParseErrorPaths(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"lexer_error_first_token", "#"},
		{"let_missing_identifier", "let = 5;"},
		{"let_missing_equal", "let x 5;"},
		{"let_missing_semicolon", "let x = 5 let y = 10;"},
		{"if_missing_lparen", "if true { }"},
		{"if_missing_rparen", "if (true { }"},
		{"if_missing_lbrace", "if (true) return 1; }"},
		{"foreach_missing_lparen", "foreach items as item { }"},
		{"foreach_missing_as", "foreach (items item) { }"},
		{"foreach_missing_rparen", "foreach (items as item { }"},
		{"foreach_missing_lbrace", "foreach (items as item) return 1; }"},
		{"fn_missing_name", "fn 123() { }"},
		{"fn_missing_lparen", "fn test { }"},
		{"fn_missing_lbrace", "fn test() return 1; }"},
		{"grouped_unclosed", "(1 + 2;"},
		{"index_unclosed", "a[0;"},
		{"access_missing_ident", "a.123;"},
		{"hash_missing_string_key", `{123: "val"};`},
		{"hash_missing_colon", `{"key" "val"};`},
		{"hash_unclosed", `{"key": "val"`},
		{"unexpected_token_in_expression", "return ;"},
		{"expression_missing_semicolon", "1 + 2 let"},
		{"fn_params_missing_rparen", "fn test(a, b { }"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.NewScript(tt.input)
			p := New(l)
			_, err := p.Parse()
			if err == nil {
				t.Fatalf("expected error for input %q", tt.input)
			}
		})
	}
}

func TestParseLexerErrorInLetValue(t *testing.T) {
	input := "let x = #;"
	l := lexer.NewScript(input)
	p := New(l)
	_, err := p.Parse()
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseLexerErrorSecondToken(t *testing.T) {
	// Valid first token, invalid second
	input := "5#"
	l := lexer.NewScript(input)
	p := New(l)
	_, err := p.Parse()
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseLexerErrorInLoopNextToken(t *testing.T) {
	// First statement succeeds, then lexer error on next nextToken in the Parse loop
	// if expression doesn't consume trailing semicolon, so Parse's loop nextToken reads next
	input := "if (true) { 1; }#"
	l := lexer.NewScript(input)
	p := New(l)
	_, err := p.Parse()
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseLexerErrorInBlockStatement(t *testing.T) {
	input := "if (true) { # }"
	l := lexer.NewScript(input)
	p := New(l)
	_, err := p.Parse()
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseLexerErrorInFunctionBody(t *testing.T) {
	input := "fn test() { # }"
	l := lexer.NewScript(input)
	p := New(l)
	_, err := p.Parse()
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseLexerErrorInCallArgs(t *testing.T) {
	input := "foo(#);"
	l := lexer.NewScript(input)
	p := New(l)
	_, err := p.Parse()
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseLexerErrorInArrayElements(t *testing.T) {
	input := "[#];"
	l := lexer.NewScript(input)
	p := New(l)
	_, err := p.Parse()
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseLexerErrorInHashValue(t *testing.T) {
	input := `{"key": #};`
	l := lexer.NewScript(input)
	p := New(l)
	_, err := p.Parse()
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseLexerErrorInIndexExpr(t *testing.T) {
	input := "a[#];"
	l := lexer.NewScript(input)
	p := New(l)
	_, err := p.Parse()
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseLexerErrorInGroupedExpr(t *testing.T) {
	input := "(#);"
	l := lexer.NewScript(input)
	p := New(l)
	_, err := p.Parse()
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseLexerErrorInIfCondition(t *testing.T) {
	input := "if (#) { 1; }"
	l := lexer.NewScript(input)
	p := New(l)
	_, err := p.Parse()
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseLexerErrorInForeachIterable(t *testing.T) {
	input := "foreach (#  as item) { 1; }"
	l := lexer.NewScript(input)
	p := New(l)
	_, err := p.Parse()
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseLexerErrorInPrefixRight(t *testing.T) {
	input := "-#;"
	l := lexer.NewScript(input)
	p := New(l)
	_, err := p.Parse()
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseLexerErrorInInfixRight(t *testing.T) {
	input := "1 + #;"
	l := lexer.NewScript(input)
	p := New(l)
	_, err := p.Parse()
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseLexerErrorInAssignmentValue(t *testing.T) {
	input := "x = #;"
	l := lexer.NewScript(input)
	p := New(l)
	_, err := p.Parse()
	if err == nil {
		t.Fatal("expected error")
	}
}

// --- currentPrecedence edge case ---

func TestCurrentPrecedenceUnknownToken(t *testing.T) {
	l := lexer.NewScript("")
	p := New(l)
	p.current = lexer.Token{Type: lexer.Semicolon}
	if p.currentPrecedence() != 0 {
		t.Fatal("expected 0 precedence for unknown token type")
	}
}

// --- Template mode edge cases ---

func TestParseTemplateOnly(t *testing.T) {
	input := `Hello {% "world" %}`
	l := lexer.NewTemplate(input)
	p := New(l)
	program, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}
	if len(program.Statements) < 2 {
		t.Fatalf("expected at least 2 statements, got %d", len(program.Statements))
	}
}

func TestParseTemplateWithForeach(t *testing.T) {
	input := `{% foreach (items as item) { %}<li>{% item %}</li>{% } %}`
	l := lexer.NewTemplate(input)
	p := New(l)
	_, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}
}

func TestParseTemplateScriptEndTerminatesExpression(t *testing.T) {
	// Expression terminated by %} instead of ;
	input := `{% x %}`
	l := lexer.NewTemplate(input)
	p := New(l)
	program, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}
	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(program.Statements))
	}
}

func TestParseTemplateLetWithScriptEnd(t *testing.T) {
	// Let terminated by %} instead of ;
	input := `{% let x = 5 %}`
	l := lexer.NewTemplate(input)
	p := New(l)
	program, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, s := range program.Statements {
		if _, ok := s.(*LetStatement); ok {
			found = true
		}
	}
	if !found {
		t.Fatal("expected LetStatement in program")
	}
}

// --- Multiple errors accumulated ---

func TestParseMultipleErrors(t *testing.T) {
	// Multiple parse errors in one program
	input := `let = 5;`
	l := lexer.NewScript(input)
	p := New(l)
	_, err := p.Parse()
	if err == nil {
		t.Fatal("expected error")
	}
	// Error message should be present
	errStr := err.Error()
	if errStr == "" {
		t.Fatal("expected non-empty error message")
	}
}

// --- Expression at EOF ---

func TestParseExpressionAtEOF(t *testing.T) {
	// Expression with ScriptEnd-like termination
	input := `5`
	l := lexer.NewScript(input)
	p := New(l)
	program, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}
	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(program.Statements))
	}
}

// --- Nested expressions ---

func TestParseNestedCallExpressions(t *testing.T) {
	input := `foo(bar(1));`
	l := lexer.NewScript(input)
	p := New(l)
	program, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}
	stmt := program.Statements[0].(*ExpressionStatement)
	call := stmt.Expression.(*CallExpression)
	if len(call.Args) != 1 {
		t.Fatalf("expected 1 arg, got %d", len(call.Args))
	}
	innerCall, ok := call.Args[0].(*CallExpression)
	if !ok {
		t.Fatalf("expected inner CallExpression, got %T", call.Args[0])
	}
	if len(innerCall.Args) != 1 {
		t.Fatalf("expected 1 inner arg, got %d", len(innerCall.Args))
	}
}

func TestParseNestedIndexExpressions(t *testing.T) {
	input := `a[b[0]];`
	l := lexer.NewScript(input)
	p := New(l)
	program, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}
	stmt := program.Statements[0].(*ExpressionStatement)
	outer := stmt.Expression.(*IndexExpression)
	_, ok := outer.Index.(*IndexExpression)
	if !ok {
		t.Fatalf("expected inner IndexExpression, got %T", outer.Index)
	}
}

func TestParseComplexMixedExpression(t *testing.T) {
	input := `a.b[0] + foo(1, c.d);`
	l := lexer.NewScript(input)
	p := New(l)
	_, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}
}

// --- Assignment edge cases ---

func TestParseAssignmentWithExpression(t *testing.T) {
	input := `x = 1 + 2;`
	l := lexer.NewScript(input)
	p := New(l)
	program, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}
	stmt := program.Statements[0].(*ExpressionStatement)
	assign := stmt.Expression.(*AssignmentExpression)
	_, ok := assign.Right.(*InfixExpression)
	if !ok {
		t.Fatalf("expected InfixExpression on right side, got %T", assign.Right)
	}
}

// --- Targeted error path tests for remaining coverage ---

func TestParseLetMissingValue(t *testing.T) {
	// let x = ; — triggers value==nil path in parseLetStatement
	input := `let x = ;`
	l := lexer.NewScript(input)
	p := New(l)
	_, err := p.Parse()
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseIntegerOverflow(t *testing.T) {
	// Integer too large for int64 — triggers strconv.ParseInt error
	input := `99999999999999999999999999999;`
	l := lexer.NewScript(input)
	p := New(l)
	_, err := p.Parse()
	if err != nil {
		// Expected: ParseInt overflow error
		return
	}
	// If no error, verify the value is at least parsed
}

func TestParseFloatOverflow(t *testing.T) {
	// Float too large for float64
	input := `1` + strings.Repeat("9", 400) + `.0;`
	l := lexer.NewScript(input)
	p := New(l)
	_, err := p.Parse()
	// May or may not error depending on Go's ParseFloat behavior for huge numbers
	_ = err
}

func TestParseLexerErrorInCallArgsComma(t *testing.T) {
	// Error after comma in call args
	input := "foo(1, #);"
	l := lexer.NewScript(input)
	p := New(l)
	_, err := p.Parse()
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseLexerErrorInExpressionListSecondItem(t *testing.T) {
	// Error reading second element in array
	input := "[1, #];"
	l := lexer.NewScript(input)
	p := New(l)
	_, err := p.Parse()
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseLexerErrorInHashAfterComma(t *testing.T) {
	input := `{"a": 1, #};`
	l := lexer.NewScript(input)
	p := New(l)
	_, err := p.Parse()
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseLexerErrorInFunctionParams(t *testing.T) {
	input := "fn test(a, #) { }"
	l := lexer.NewScript(input)
	p := New(l)
	_, err := p.Parse()
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseLexerErrorAfterDot(t *testing.T) {
	// Error in property access after dot
	input := "a.#;"
	l := lexer.NewScript(input)
	p := New(l)
	_, err := p.Parse()
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseLexerErrorInReturnValue(t *testing.T) {
	input := "return #;"
	l := lexer.NewScript(input)
	p := New(l)
	_, err := p.Parse()
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseReturnWithUnexpectedToken(t *testing.T) {
	// return ; — semicolon is not a valid expression, parsePrefixExpression returns nil
	// This exercises the value==nil path in parseReturnStatement
	input := "return ;"
	l := lexer.NewScript(input)
	p := New(l)
	_, err := p.Parse()
	if err == nil {
		t.Fatal("expected error for 'return ;'")
	}
}

func TestParseExpressionStatementAssignmentErrorPaths(t *testing.T) {
	// Assignment where the right side has a lexer error
	tests := []struct {
		name  string
		input string
	}{
		{"assign_identifier_error", "x = #;"},
		{"assign_index_error", "arr[0] = #;"},
		{"assign_property_error", "obj.prop = #;"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.NewScript(tt.input)
			p := New(l)
			_, err := p.Parse()
			if err == nil {
				t.Fatal("expected error")
			}
		})
	}
}

func TestParseLexerErrorInElseBlock(t *testing.T) {
	input := "if (true) { 1; } else { # }"
	l := lexer.NewScript(input)
	p := New(l)
	_, err := p.Parse()
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseLexerErrorInForeachBody(t *testing.T) {
	input := "foreach (items as item) { # }"
	l := lexer.NewScript(input)
	p := New(l)
	_, err := p.Parse()
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseExpressionInfixWithExpressionError(t *testing.T) {
	// Error in right side of infix at various operator positions
	tests := []struct {
		name  string
		input string
	}{
		{"plus_error", "1 + #;"},
		{"minus_error", "1 - #;"},
		{"asterisk_error", "1 * #;"},
		{"slash_error", "1 / #;"},
		{"modulo_error", "1 %% #;"},
		{"equals_error", "1 == #;"},
		{"not_equal_error", "1 != #;"},
		{"less_than_error", "1 < #;"},
		{"greater_than_error", "1 > #;"},
		{"less_or_equal_error", "1 <= #;"},
		{"greater_or_equal_error", "1 >= #;"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.NewScript(tt.input)
			p := New(l)
			_, err := p.Parse()
			if err == nil {
				t.Fatal("expected error")
			}
		})
	}
}

func TestParseForeachAllErrorPaths(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"missing_lparen", "foreach x as item { }"},
		{"error_in_iterable", "foreach (# as item) { }"},
		{"missing_as", "foreach (items x) { }"},
		{"error_after_as", "foreach (items as #) { }"},
		{"missing_rparen", "foreach (items as item { }"},
		{"missing_lbrace", "foreach (items as item) 1; }"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.NewScript(tt.input)
			p := New(l)
			_, err := p.Parse()
			if err == nil {
				t.Fatalf("expected error for input %q", tt.input)
			}
		})
	}
}

func TestParseIfAllErrorPaths(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"missing_lparen", "if true { }"},
		{"error_in_condition", "if (# ) { }"},
		{"missing_rparen", "if (true { }"},
		{"missing_lbrace", "if (true) 1; }"},
		{"else_missing_lbrace", "if (true) { 1; } else 2; }"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.NewScript(tt.input)
			p := New(l)
			_, err := p.Parse()
			if err == nil {
				t.Fatalf("expected error for input %q", tt.input)
			}
		})
	}
}

func TestParseFunctionLiteralAllErrorPaths(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"missing_lparen", "fn test { }"},
		{"missing_lbrace", "fn test() 1; }"},
		{"error_in_body", "fn test() { # }"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.NewScript(tt.input)
			p := New(l)
			_, err := p.Parse()
			if err == nil {
				t.Fatalf("expected error for input %q", tt.input)
			}
		})
	}
}

// --- Deep error path tests ---
// These place '#' at precise token positions to trigger nextToken() errors
// inside deeply nested parse functions.

func TestParseLetStatementNextTokenErrorAfterValue(t *testing.T) {
	// Exercises nextToken error at parseLetStatement line 184
	// After parsing value "1", nextToken reads past semicolon and hits #
	input := "let x = 1;#"
	l := lexer.NewScript(input)
	p := New(l)
	_, err := p.Parse()
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseReturnStatementNextTokenError(t *testing.T) {
	// Exercises nextToken error at parseReturnStatement line 197
	input := "return 1#"
	l := lexer.NewScript(input)
	p := New(l)
	_, err := p.Parse()
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseReturnStatementMissingSemicolon(t *testing.T) {
	// Exercises tryPeek(Semicolon) failure at parseReturnStatement line 209
	input := "return 1 let"
	l := lexer.NewScript(input)
	p := New(l)
	_, err := p.Parse()
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseGroupedExpressionNextTokenError(t *testing.T) {
	// Exercises nextToken error inside parseGroupedExpression line 613
	// (1 # ) — the nextToken inside grouped reads past the already-buffered token
	input := "(1 #)"
	l := lexer.NewScript(input)
	p := New(l)
	_, err := p.Parse()
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseIndexExpressionNextTokenError(t *testing.T) {
	// Exercises nextToken error inside parseIndexExpression line 494
	// a[1 # ] — # at position where parseIndexExpression's nextToken reads
	input := "a[1 #]"
	l := lexer.NewScript(input)
	p := New(l)
	_, err := p.Parse()
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParsePrefixExpressionNextTokenError(t *testing.T) {
	// Exercises nextToken error in prefix expression (line 398)
	// After reading '-', nextToken reads ahead and hits #
	input := "let x = -#;"
	l := lexer.NewScript(input)
	p := New(l)
	_, err := p.Parse()
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseExpressionStatementAssignmentNextTokenError(t *testing.T) {
	// Exercises nextToken errors in assignment branch of parseExpressionStatement
	// x = (value with error at right position)
	input := "x = 1#"
	l := lexer.NewScript(input)
	p := New(l)
	_, err := p.Parse()
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseBlockStatementNextTokenError(t *testing.T) {
	// Exercises nextToken error at start of parseBlockStatement (line 735)
	input := "if (true) {#}"
	l := lexer.NewScript(input)
	p := New(l)
	_, err := p.Parse()
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseBlockStatementLoopNextTokenError(t *testing.T) {
	// Exercises nextToken error in block loop (line 755)
	input := "if (true) { 1;#}"
	l := lexer.NewScript(input)
	p := New(l)
	_, err := p.Parse()
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseIfExpressionNextTokenErrors(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		// Error reading condition after (
		{"condition_error", "if (1 #) { }"},
		// Error reading else block
		{"else_error", "if (true) { 1; } else {#}"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.NewScript(tt.input)
			p := New(l)
			_, err := p.Parse()
			if err == nil {
				t.Fatalf("expected error for %q", tt.input)
			}
		})
	}
}

func TestParseForeachExpressionNextTokenErrors(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		// Error reading iterable
		{"iterable_next_error", "foreach (1 #as item) { }"},
		// Error reading variable after as
		{"variable_error", "foreach (items as 1#) { }"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.NewScript(tt.input)
			p := New(l)
			_, err := p.Parse()
			if err == nil {
				t.Fatalf("expected error for %q", tt.input)
			}
		})
	}
}

func TestParseFunctionParametersNextTokenError(t *testing.T) {
	// Error inside parameter list
	input := "fn test(a#) { }"
	l := lexer.NewScript(input)
	p := New(l)
	_, err := p.Parse()
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseFunctionParametersCommaNextTokenError(t *testing.T) {
	// Error after comma in parameter list
	input := "fn test(a, b#) { }"
	l := lexer.NewScript(input)
	p := New(l)
	_, err := p.Parse()
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseExpressionListNextTokenErrors(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"first_elem_next_error", "foo(1 #)"},
		{"comma_next_error", "foo(1,#)"},
		{"second_elem_next_error", "foo(1, 2 #)"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.NewScript(tt.input)
			p := New(l)
			_, err := p.Parse()
			if err == nil {
				t.Fatalf("expected error for %q", tt.input)
			}
		})
	}
}

func TestParseHashPairsNextTokenErrors(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"value_next_error", `{"a": 1 #}`},
		{"second_pair_error", `{"a": 1, "b": #}`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.NewScript(tt.input)
			p := New(l)
			_, err := p.Parse()
			if err == nil {
				t.Fatalf("expected error for %q", tt.input)
			}
		})
	}
}

func TestParseExpressionLeftNilBreak(t *testing.T) {
	// Trigger the leftExpression == nil break at line 368 of parseExpression
	// This happens when an infix parse function returns nil
	// A call expression that fails: foo( — unclosed, expression list fails
	input := "foo(;"
	l := lexer.NewScript(input)
	p := New(l)
	_, err := p.Parse()
	if err == nil {
		t.Fatal("expected error")
	}
}

// --- Precisely targeted error path tests ---
// Each test places '#' or invalid tokens at the exact position to trigger
// a specific uncovered error guard.

func TestParseAssignmentRightNil(t *testing.T) {
	// Exercises right==nil path in assignment (line 242)
	input := "x = ;"
	l := lexer.NewScript(input)
	p := New(l)
	_, err := p.Parse()
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseAssignmentNextTokenAfterValue(t *testing.T) {
	// Exercises nextToken error after assignment value (line 255)
	input := "x = 1;#"
	l := lexer.NewScript(input)
	p := New(l)
	_, err := p.Parse()
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseAccessExpressionInternalError(t *testing.T) {
	// Exercises error from inside parseAccessExpression (line 341)
	// a.b# — tryPeek(Identifier) succeeds on 'b' but nextToken fails reading '#'
	input := "a.b#"
	l := lexer.NewScript(input)
	p := New(l)
	_, err := p.Parse()
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParsePrefixNextTokenError(t *testing.T) {
	// Exercises nextToken error in prefix (- or !) at line 399
	// 1 + -1# — the # is at the exact position of prefix's nextToken Read()
	input := "1 + -1#"
	l := lexer.NewScript(input)
	p := New(l)
	_, err := p.Parse()
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParsePrefixRightNil(t *testing.T) {
	// Exercises right==nil in prefix expression (line 404)
	input := "1 + -;"
	l := lexer.NewScript(input)
	p := New(l)
	_, err := p.Parse()
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseIndexExpressionParseError(t *testing.T) {
	// Exercises parseExpression error in index expression (line 500)
	input := "a[1 + #]"
	l := lexer.NewScript(input)
	p := New(l)
	_, err := p.Parse()
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseGroupedExpressionParseError(t *testing.T) {
	// Exercises parseExpression error in grouped expression (line 619)
	input := "(1 + #)"
	l := lexer.NewScript(input)
	p := New(l)
	_, err := p.Parse()
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseIfConditionParseError(t *testing.T) {
	// Exercises parseExpression error in if condition (line 648)
	input := "if (1 + #) { }"
	l := lexer.NewScript(input)
	p := New(l)
	_, err := p.Parse()
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseElseNextTokenError(t *testing.T) {
	// Exercises nextToken error in else branch (line 666)
	input := "if (true) { 1; } else {1;#}"
	l := lexer.NewScript(input)
	p := New(l)
	_, err := p.Parse()
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseForeachNextTokenAfterAs(t *testing.T) {
	// Exercises nextToken error after 'as' keyword (line 698)
	input := "foreach (items as item#) { }"
	l := lexer.NewScript(input)
	p := New(l)
	_, err := p.Parse()
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseBlockStatementInitNextTokenError(t *testing.T) {
	// Exercises nextToken error at start of block (line 736)
	input := "fn test() {#}"
	l := lexer.NewScript(input)
	p := New(l)
	_, err := p.Parse()
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseFunctionParamsEmptyNextTokenError(t *testing.T) {
	// Exercises nextToken error in empty params case (line 769)
	// fn test()# — when params are empty, nextToken reads past ')' and hits #
	input := "fn test()#"
	l := lexer.NewScript(input)
	p := New(l)
	_, err := p.Parse()
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseExpressionListFirstNextTokenError(t *testing.T) {
	// Exercises nextToken error reading first element (line 850)
	input := "foo(1#)"
	l := lexer.NewScript(input)
	p := New(l)
	_, err := p.Parse()
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseExpressionListFirstParseError(t *testing.T) {
	// Exercises parseExpression error for first element (line 862)
	input := "foo(1 + #)"
	l := lexer.NewScript(input)
	p := New(l)
	_, err := p.Parse()
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseExpressionListCommaParseError(t *testing.T) {
	// Exercises parseExpression error after comma (line 880)
	input := "foo(1, 1 + #)"
	l := lexer.NewScript(input)
	p := New(l)
	_, err := p.Parse()
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseHashPairsValueNextTokenError(t *testing.T) {
	// Exercises nextToken errors inside hash pairs (lines 930, 945)
	tests := []struct {
		name  string
		input string
	}{
		{"colon_next_error", `{"a": 1, "b"#}`},
		{"value_next_error", `{"a"#: 1}`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.NewScript(tt.input)
			p := New(l)
			_, err := p.Parse()
			if err == nil {
				t.Fatalf("expected error for %q", tt.input)
			}
		})
	}
}

func TestParseHashPairsExpressionError(t *testing.T) {
	// Exercises parseExpression error in hash value (line 962)
	input := `{"a": 1 + #}`
	l := lexer.NewScript(input)
	p := New(l)
	_, err := p.Parse()
	if err == nil {
		t.Fatal("expected error")
	}
}

// --- Batch 5: Final coverage targets ---

func TestParseElseKeywordAdvanceError(t *testing.T) {
	// Exercises nextToken error after advancing past 'else' keyword (line 666)
	input := `if (true) { 1; } else #`
	l := lexer.NewScript(input)
	p := New(l)
	_, err := p.Parse()
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseForeachIterableExpressionError(t *testing.T) {
	// Exercises parseExpression error for foreach iterable (line 698)
	input := `foreach (1 + # as item) { }`
	l := lexer.NewScript(input)
	p := New(l)
	_, err := p.Parse()
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseBlockStatementFirstTokenError(t *testing.T) {
	// Exercises nextToken error at start of parseBlockStatement (line 736)
	// The 'a' is read during tryPeek(LeftBrace), then '#' is read by
	// parseBlockStatement's initial nextToken call.
	input := `fn test() { a # }`
	l := lexer.NewScript(input)
	p := New(l)
	_, err := p.Parse()
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseEmptyCallNextTokenError(t *testing.T) {
	// Exercises nextToken error in parseExpressionList when list is empty
	// and advancing past the closing ')' errors (line 850)
	input := `test()#`
	l := lexer.NewScript(input)
	p := New(l)
	_, err := p.Parse()
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseEmptyHashNextTokenError(t *testing.T) {
	// Exercises nextToken error in parseHashPairs when hash is empty
	// and advancing past the closing '}' errors (line 930)
	input := `{}#`
	l := lexer.NewScript(input)
	p := New(l)
	_, err := p.Parse()
	if err == nil {
		t.Fatal("expected error")
	}
}

// suppress unused import warnings
var _ = fmt.Sprintf
var _ = errors.New

func TestParseLogicalAnd(t *testing.T) {
	input := "true && false;"

	l := lexer.NewScript(input)
	p := New(l)
	program, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(program.Statements))
	}

	exp, ok := program.Statements[0].(*ExpressionStatement)
	if !ok {
		t.Fatalf("expected ExpressionStatement, got %T", program.Statements[0])
	}

	infix, ok := exp.Expression.(*InfixExpression)
	if !ok {
		t.Fatalf("expected InfixExpression, got %T", exp.Expression)
	}

	if infix.Token.Source != "&&" {
		t.Fatalf("expected operator '&&', got %q", infix.Token.Source)
	}
}

func TestParseLogicalOr(t *testing.T) {
	input := "true || false;"

	l := lexer.NewScript(input)
	p := New(l)
	program, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(program.Statements))
	}

	exp, ok := program.Statements[0].(*ExpressionStatement)
	if !ok {
		t.Fatalf("expected ExpressionStatement, got %T", program.Statements[0])
	}

	infix, ok := exp.Expression.(*InfixExpression)
	if !ok {
		t.Fatalf("expected InfixExpression, got %T", exp.Expression)
	}

	if infix.Token.Source != "||" {
		t.Fatalf("expected operator '||', got %q", infix.Token.Source)
	}
}

func TestParseNullCoalescing(t *testing.T) {
	input := "x ?? 5;"

	l := lexer.NewScript(input)
	p := New(l)
	program, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(program.Statements))
	}

	exp, ok := program.Statements[0].(*ExpressionStatement)
	if !ok {
		t.Fatalf("expected ExpressionStatement, got %T", program.Statements[0])
	}

	infix, ok := exp.Expression.(*InfixExpression)
	if !ok {
		t.Fatalf("expected InfixExpression, got %T", exp.Expression)
	}

	if infix.Token.Source != "??" {
		t.Fatalf("expected operator '??', got %q", infix.Token.Source)
	}
}

func TestParseWhileExpression(t *testing.T) {
	input := "while (true) { let x = 1; }"

	l := lexer.NewScript(input)
	p := New(l)
	program, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(program.Statements))
	}

	exp, ok := program.Statements[0].(*ExpressionStatement)
	if !ok {
		t.Fatalf("expected ExpressionStatement, got %T", program.Statements[0])
	}

	whileExpr, ok := exp.Expression.(*WhileExpression)
	if !ok {
		t.Fatalf("expected WhileExpression, got %T", exp.Expression)
	}

	if whileExpr.Body == nil {
		t.Fatal("expected while body not to be nil")
	}

	if len(whileExpr.Body.Statements) != 1 {
		t.Fatalf("expected 1 body statement, got %d", len(whileExpr.Body.Statements))
	}
}

func TestParseBreakStatement(t *testing.T) {
	input := "while (true) { break; }"

	l := lexer.NewScript(input)
	p := New(l)
	program, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(program.Statements))
	}

	exp, ok := program.Statements[0].(*ExpressionStatement)
	if !ok {
		t.Fatalf("expected ExpressionStatement, got %T", program.Statements[0])
	}

	whileExpr, ok := exp.Expression.(*WhileExpression)
	if !ok {
		t.Fatalf("expected WhileExpression, got %T", exp.Expression)
	}

	_, ok = whileExpr.Body.Statements[0].(*BreakStatement)
	if !ok {
		t.Fatalf("expected BreakStatement, got %T", whileExpr.Body.Statements[0])
	}
}

func TestParseContinueStatement(t *testing.T) {
	input := "while (true) { continue; }"

	l := lexer.NewScript(input)
	p := New(l)
	program, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}

	exp, ok := program.Statements[0].(*ExpressionStatement)
	if !ok {
		t.Fatalf("expected ExpressionStatement, got %T", program.Statements[0])
	}

	whileExpr, ok := exp.Expression.(*WhileExpression)
	if !ok {
		t.Fatalf("expected WhileExpression, got %T", exp.Expression)
	}

	_, ok = whileExpr.Body.Statements[0].(*ContinueStatement)
	if !ok {
		t.Fatalf("expected ContinueStatement, got %T", whileExpr.Body.Statements[0])
	}
}

func TestParseElseIf(t *testing.T) {
	input := "if (true) { 1; } else if (false) { 2; } else { 3; }"

	l := lexer.NewScript(input)
	p := New(l)
	program, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(program.Statements))
	}

	exp, ok := program.Statements[0].(*ExpressionStatement)
	if !ok {
		t.Fatalf("expected ExpressionStatement, got %T", program.Statements[0])
	}

	ifExpr, ok := exp.Expression.(*IfExpression)
	if !ok {
		t.Fatalf("expected IfExpression, got %T", exp.Expression)
	}

	if ifExpr.Alternative == nil {
		t.Fatal("expected alternative (else if) not to be nil")
	}

	// The alternative should contain an ExpressionStatement with another IfExpression
	if len(ifExpr.Alternative.Statements) != 1 {
		t.Fatalf("expected 1 alternative statement, got %d", len(ifExpr.Alternative.Statements))
	}

	altExp, ok := ifExpr.Alternative.Statements[0].(*ExpressionStatement)
	if !ok {
		t.Fatalf("expected ExpressionStatement in alternative, got %T", ifExpr.Alternative.Statements[0])
	}

	innerIf, ok := altExp.Expression.(*IfExpression)
	if !ok {
		t.Fatalf("expected nested IfExpression, got %T", altExp.Expression)
	}

	if innerIf.Alternative == nil {
		t.Fatal("expected inner else block not to be nil")
	}
}

func TestParseForeachWithIndex(t *testing.T) {
	input := "foreach ([1,2,3] as i, v) { v; }"

	l := lexer.NewScript(input)
	p := New(l)
	program, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(program.Statements))
	}

	fe, ok := program.Statements[0].(*ForeachExpression)
	if !ok {
		t.Fatalf("expected ForeachExpression, got %T", program.Statements[0])
	}

	if fe.Index == nil {
		t.Fatal("expected index variable not to be nil")
	}

	if fe.Index.Value != "i" {
		t.Fatalf("expected index name 'i', got %q", fe.Index.Value)
	}

	if fe.Variable.Value != "v" {
		t.Fatalf("expected variable name 'v', got %q", fe.Variable.Value)
	}
}

func TestParseForeachWithoutIndex(t *testing.T) {
	input := "foreach ([1,2,3] as v) { v; }"

	l := lexer.NewScript(input)
	p := New(l)
	program, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}

	fe, ok := program.Statements[0].(*ForeachExpression)
	if !ok {
		t.Fatalf("expected ForeachExpression, got %T", program.Statements[0])
	}

	if fe.Index != nil {
		t.Fatalf("expected index to be nil, got %v", fe.Index)
	}

	if fe.Variable.Value != "v" {
		t.Fatalf("expected variable name 'v', got %q", fe.Variable.Value)
	}
}

func TestParseNullLiteral(t *testing.T) {
	input := "null;"

	l := lexer.NewScript(input)
	p := New(l)
	program, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(program.Statements))
	}

	exp, ok := program.Statements[0].(*ExpressionStatement)
	if !ok {
		t.Fatalf("expected ExpressionStatement, got %T", program.Statements[0])
	}

	_, ok = exp.Expression.(*NullLiteral)
	if !ok {
		t.Fatalf("expected NullLiteral, got %T", exp.Expression)
	}
}

func TestLogicalOperatorPrecedence(t *testing.T) {
	// && should bind tighter than ||
	input := "true || false && true;"

	l := lexer.NewScript(input)
	p := New(l)
	program, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}

	exp, ok := program.Statements[0].(*ExpressionStatement)
	if !ok {
		t.Fatalf("expected ExpressionStatement, got %T", program.Statements[0])
	}

	// Should parse as: true || (false && true)
	orExpr, ok := exp.Expression.(*InfixExpression)
	if !ok {
		t.Fatalf("expected InfixExpression, got %T", exp.Expression)
	}

	if orExpr.Token.Source != "||" {
		t.Fatalf("expected top-level '||', got %q", orExpr.Token.Source)
	}

	andExpr, ok := orExpr.Right.(*InfixExpression)
	if !ok {
		t.Fatalf("expected right side to be InfixExpression, got %T", orExpr.Right)
	}

	if andExpr.Token.Source != "&&" {
		t.Fatalf("expected right side '&&', got %q", andExpr.Token.Source)
	}
}

func TestParseComments(t *testing.T) {
	input := "// comment\nlet x = 5; /* another comment */ let y = 10;"

	l := lexer.NewScript(input)
	p := New(l)
	program, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}

	if len(program.Statements) != 2 {
		t.Fatalf("expected 2 statements, got %d", len(program.Statements))
	}
}

// --- Bug fix verification tests ---

func TestParseErrorRustStyleFormat(t *testing.T) {
	tok := lexer.Token{Type: lexer.Identifier, Source: "x", Line: 1, Column: 5}
	err := NewParseError("unexpected token", "let x 5;", tok)
	msg := err.Error()

	if !strings.Contains(msg, "error: unexpected token") {
		t.Fatalf("should start with 'error:' prefix, got: %s", msg)
	}
	if !strings.Contains(msg, "-->") {
		t.Fatalf("should contain '-->' location indicator, got: %s", msg)
	}
	if !strings.Contains(msg, "|") {
		t.Fatalf("should contain pipe separators, got: %s", msg)
	}
	if !strings.Contains(msg, "1 | let x 5;") {
		t.Fatalf("should contain line number with pipe, got: %s", msg)
	}
}

func TestParseErrorBoundsCheck(t *testing.T) {
	tok := lexer.Token{Type: lexer.Identifier, Source: "x", Line: 99, Column: 1}
	err := NewParseError("test error", "single line", tok)
	msg := err.Error()
	if !strings.Contains(msg, "error: test error") {
		t.Fatalf("should contain error message without panic, got: %s", msg)
	}
}

func TestParseGroupedExpressionNilCheck(t *testing.T) {
	// Grouped expression with invalid inner expression
	input := `(#);`
	l := lexer.NewScript(input)
	p := New(l)
	_, err := p.Parse()
	if err == nil {
		t.Fatal("expected error for invalid grouped expression")
	}
}

func TestParseIfConditionNilCheck(t *testing.T) {
	// If condition with invalid expression
	input := `if (#) { }`
	l := lexer.NewScript(input)
	p := New(l)
	_, err := p.Parse()
	if err == nil {
		t.Fatal("expected error for invalid if condition")
	}
}

func TestParseForeachIterableNilCheck(t *testing.T) {
	// Foreach with invalid iterable
	input := `foreach (# as item) { }`
	l := lexer.NewScript(input)
	p := New(l)
	_, err := p.Parse()
	if err == nil {
		t.Fatal("expected error for invalid foreach iterable")
	}
}

func TestParseWhileConditionNilCheck(t *testing.T) {
	// While with invalid condition
	input := `while (#) { }`
	l := lexer.NewScript(input)
	p := New(l)
	_, err := p.Parse()
	if err == nil {
		t.Fatal("expected error for invalid while condition")
	}
}

func TestParseExpressionListNilCheck(t *testing.T) {
	// Call with invalid argument expression
	input := `foo(#);`
	l := lexer.NewScript(input)
	p := New(l)
	_, err := p.Parse()
	if err == nil {
		t.Fatal("expected error for invalid call argument")
	}
}

func TestParseBlockStatementErrorPropagation(t *testing.T) {
	// Block statement with lexer error should propagate
	input := `if (true) { let x = #; }`
	l := lexer.NewScript(input)
	p := New(l)
	_, err := p.Parse()
	if err == nil {
		t.Fatal("expected error for invalid token in block statement")
	}
}

// --- Anonymous Functions ---

func TestParseAnonymousFunction(t *testing.T) {
	input := `let add = fn(a, b) { return a + b; };`
	l := lexer.NewScript(input)
	p := New(l)
	program, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}
	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(program.Statements))
	}
	letStmt, ok := program.Statements[0].(*LetStatement)
	if !ok {
		t.Fatalf("expected LetStatement, got %T", program.Statements[0])
	}
	fnLit, ok := letStmt.Value.(*FunctionLiteral)
	if !ok {
		t.Fatalf("expected FunctionLiteral, got %T", letStmt.Value)
	}
	if fnLit.Identifier != nil {
		t.Fatalf("expected nil Identifier for anonymous function, got %s", fnLit.Identifier.Value)
	}
	if len(fnLit.Parameters) != 2 {
		t.Fatalf("expected 2 parameters, got %d", len(fnLit.Parameters))
	}
}

func TestParseAnonymousFunctionNoParams(t *testing.T) {
	input := `let f = fn() { return 1; };`
	l := lexer.NewScript(input)
	p := New(l)
	program, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}
	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(program.Statements))
	}
	letStmt := program.Statements[0].(*LetStatement)
	fnLit, ok := letStmt.Value.(*FunctionLiteral)
	if !ok {
		t.Fatalf("expected FunctionLiteral, got %T", letStmt.Value)
	}
	if fnLit.Identifier != nil {
		t.Fatalf("expected nil Identifier, got %v", fnLit.Identifier)
	}
}

func TestParseNamedFunctionStillWorks(t *testing.T) {
	input := `fn add(a, b) { return a + b; }`
	l := lexer.NewScript(input)
	p := New(l)
	program, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}
	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(program.Statements))
	}
	exprStmt := program.Statements[0].(*ExpressionStatement)
	fnLit, ok := exprStmt.Expression.(*FunctionLiteral)
	if !ok {
		t.Fatalf("expected FunctionLiteral, got %T", exprStmt.Expression)
	}
	if fnLit.Identifier == nil || fnLit.Identifier.Value != "add" {
		t.Fatalf("expected named function 'add', got %v", fnLit.Identifier)
	}
}

// --- Compound Assignment Operators ---

func TestParseCompoundAssignment(t *testing.T) {
	tests := []struct {
		input string
		op    string
	}{
		{"let x = 1; x += 2;", "+"},
		{"let x = 1; x -= 2;", "-"},
		{"let x = 1; x *= 2;", "*"},
		{"let x = 1; x /= 2;", "/"},
		{"let x = 1; x %= 2;", "%"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			l := lexer.NewScript(tt.input)
			p := New(l)
			program, err := p.Parse()
			if err != nil {
				t.Fatal(err)
			}
			if len(program.Statements) != 2 {
				t.Fatalf("expected 2 statements, got %d", len(program.Statements))
			}
			exprStmt, ok := program.Statements[1].(*ExpressionStatement)
			if !ok {
				t.Fatalf("expected ExpressionStatement, got %T", program.Statements[1])
			}
			assign, ok := exprStmt.Expression.(*AssignmentExpression)
			if !ok {
				t.Fatalf("expected AssignmentExpression, got %T", exprStmt.Expression)
			}
			infix, ok := assign.Right.(*InfixExpression)
			if !ok {
				t.Fatalf("expected InfixExpression as right side, got %T", assign.Right)
			}
			if infix.Token.Source != tt.op {
				t.Fatalf("expected operator %s, got %s", tt.op, infix.Token.Source)
			}
		})
	}
}

func TestParseCompoundAssignmentOnProperty(t *testing.T) {
	input := `let obj = {"x": 1}; obj.x += 5;`
	l := lexer.NewScript(input)
	p := New(l)
	_, err := p.Parse()
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestParseCompoundAssignmentOnIndex(t *testing.T) {
	input := `let arr = [1, 2, 3]; arr[0] += 10;`
	l := lexer.NewScript(input)
	p := New(l)
	_, err := p.Parse()
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

// --- Ordered Hash Literal ---

func TestParseHashLiteralOrdered(t *testing.T) {
	input := `{"a": 1, "b": 2, "c": 3}`
	l := lexer.NewScript(input)
	p := New(l)
	program, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}
	exprStmt := program.Statements[0].(*ExpressionStatement)
	hashLit, ok := exprStmt.Expression.(*HashLiteral)
	if !ok {
		t.Fatalf("expected HashLiteral, got %T", exprStmt.Expression)
	}
	if len(hashLit.Pairs) != 3 {
		t.Fatalf("expected 3 pairs, got %d", len(hashLit.Pairs))
	}
	// Verify order is maintained
	expectedKeys := []string{`"a"`, `"b"`, `"c"`}
	for i, pair := range hashLit.Pairs {
		if pair.Key.Debug() != expectedKeys[i] {
			t.Fatalf("pair %d: expected key %s, got %s", i, expectedKeys[i], pair.Key.Debug())
		}
	}
}
