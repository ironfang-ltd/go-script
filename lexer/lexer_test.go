package lexer

import (
	"testing"
)

func TestFn(t *testing.T) {

	script := `fn test(a,b) { return a + b; }`

	want := []Token{
		{Type: Function, Source: "fn"},
		{Type: Identifier, Source: "test"},
		{Type: LeftParen, Source: "("},
		{Type: Identifier, Source: "a"},
		{Type: Comma, Source: ","},
		{Type: Identifier, Source: "b"},
		{Type: RightParen, Source: ")"},
		{Type: LeftBrace, Source: "{"},
		{Type: Return, Source: "return"},
		{Type: Identifier, Source: "a"},
		{Type: Plus, Source: "+"},
		{Type: Identifier, Source: "b"},
		{Type: Semicolon, Source: ";"},
		{Type: RightBrace, Source: "}"},
		{Type: EndOfFile, Source: ""},
	}

	l := New(script)

	for _, w := range want {
		tok, err := l.Read()
		if err != nil {
			t.Error(err)
		}

		if tok.Type != w.Type {
			t.Errorf("want %s, got %s", w.Type, tok.Type)
		}
	}
}
