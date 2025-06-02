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
