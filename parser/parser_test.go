package parser

import (
	"fmt"
	"github.com/ironfang-ltd/go-script/lexer"
	"testing"
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
