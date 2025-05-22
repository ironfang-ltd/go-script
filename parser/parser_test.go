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
