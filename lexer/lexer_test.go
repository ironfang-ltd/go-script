package lexer

import (
	"errors"
	"testing"
)

func TestLineAndColumnScript(t *testing.T) {

	script := "let a = 1;\nlet b = #;"

	l := NewScript(script)

	for {
		token, err := l.Read()
		if err != nil {
			var lerr *TokenError
			if errors.As(err, &lerr) {
				if lerr.Column != 9 {
					t.Fatalf("expect line 9, got %d\n", lerr.Column)
				}

				if lerr.Line != 2 {
					t.Fatalf("expect line 2, got %d\n", lerr.Line)
				}

				break
			}

			t.Fatal(err)
		}

		if token.Type == EndOfFile {
			break
		}
	}
}

func TestLineAndColumnTemplate(t *testing.T) {

	script := `hello {% # %}`

	want := []Token{
		{Type: Text, Source: "hello ", Line: 1, Column: 1},
		{Type: ScriptStart, Source: "{%", Line: 1, Column: 7},
	}

	l := NewTemplate(script)

	idx := 0
	for {
		token, err := l.Read()
		if err != nil {
			var lerr *TokenError
			if errors.As(err, &lerr) {
				if lerr.Column != 10 {
					t.Fatalf("expect line 10, got %d\n", lerr.Column)
				}

				if lerr.Line != 1 {
					t.Fatalf("expect line 1, got %d\n", lerr.Line)
				}

				break
			}

			t.Fatal(err)
		}

		w := want[idx]
		if w.Type != token.Type {
			t.Fatalf("expect token type %s, got %s\n", w.Type, token.Type)
		}
		if w.Source != token.Source {
			t.Fatalf("expect token source %s, got %s\n", w.Source, token.Source)
		}
		if w.Line != token.Line {
			t.Fatalf("expect token line %d, got %d\n", w.Line, token.Line)
		}
		if w.Column != token.Column {
			t.Fatalf("expect token column %d, got %d\n", w.Column, token.Column)
		}

		if token.Type == EndOfFile {
			break
		}

		idx++
	}
}

func TestTemplateLineAndColumnWithWhitespace(t *testing.T) {

	script := "\n\thello {% # %}"

	want := []Token{
		{Type: Text, Source: "\n\thello ", Line: 1, Column: 1},
		{Type: ScriptStart, Source: "{%", Line: 2, Column: 8},
	}

	l := NewTemplate(script)

	idx := 0
	for {
		token, err := l.Read()
		if err != nil {
			var lerr *TokenError
			if errors.As(err, &lerr) {
				if lerr.Column != 11 {
					t.Fatalf("expect line 11, got %d\n", lerr.Column)
				}

				if lerr.Line != 2 {
					t.Fatalf("expect line 2, got %d\n", lerr.Line)
				}

				break
			}

			t.Fatal(err)
		}

		w := want[idx]
		if w.Type != token.Type {
			t.Fatalf("expect token type %s, got %s\n", w.Type, token.Type)
		}
		if w.Source != token.Source {
			t.Fatalf("expect token source %s, got %s\n", w.Source, token.Source)
		}
		if w.Line != token.Line {
			t.Fatalf("expect token line %d, got %d\n", w.Line, token.Line)
		}
		if w.Column != token.Column {
			t.Fatalf("expect token column %d, got %d\n", w.Column, token.Column)
		}

		if token.Type == EndOfFile {
			break
		}

		idx++
	}
}

func TestTemplateLineAndColumnWithFull(t *testing.T) {

	script := `
	<h1>{% title %}</h1>
	<h2>Items ({% count(items) %}): </h2>
	<ul>
		{% foreach items as item %}
			<li>{% item.name %}</li>
		{% end %}
	</ul>
`

	want := []Token{
		{Type: Text, Source: "\n\t<h1>", Line: 1, Column: 1},
		{Type: ScriptStart, Source: "{%", Line: 2, Column: 6},
		{Type: Identifier, Source: "title", Line: 2, Column: 9},
		{Type: ScriptEnd, Source: "%}", Line: 2, Column: 15},
		{Type: Text, Source: "</h1>\n\t<h2>Items (", Line: 2, Column: 17},
		{Type: ScriptStart, Source: "{%", Line: 3, Column: 13},
		{Type: Identifier, Source: "count", Line: 3, Column: 16},
		{Type: LeftParen, Source: "(", Line: 3, Column: 21},
		{Type: Identifier, Source: "items", Line: 3, Column: 22},
		{Type: RightParen, Source: ")", Line: 3, Column: 27},
		{Type: ScriptEnd, Source: "%}", Line: 3, Column: 29},
	}

	l := NewTemplate(script)

	idx := 0
	for {
		token, err := l.Read()
		if err != nil {
			t.Fatal(err)
		}

		if token.Type == EndOfFile {
			break
		}

		w := want[idx]
		if w.Type != token.Type {
			t.Fatalf("[%d] expect token type %s, got %s\n", idx, w.Type, token.Type)
		}

		if w.Source != token.Source {
			t.Fatalf("[%d] expect token source %s, got %s\n", idx, w.Source, token.Source)
		}

		if w.Line != token.Line {
			t.Fatalf("[%d] expect token line %d, got %d\n", idx, w.Line, token.Line)
		}

		if w.Column != token.Column {
			t.Fatalf("[%d] expect token column %d, got %d\n", idx, w.Column, token.Column)
		}

		idx++
	}
}

func TestLet(t *testing.T) {

	script := `let a = 1;`

	want := []Token{
		{Type: Let, Source: "let"},
		{Type: Identifier, Source: "a"},
		{Type: Equal, Source: "="},
		{Type: Integer, Source: "1"},
		{Type: Semicolon, Source: ";"},
		{Type: EndOfFile, Source: ""},
	}

	l := NewScript(script)

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

	l := NewScript(script)

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

func TestTemplate(t *testing.T) {

	script := `hello {% world %}`

	want := []Token{
		{Type: Text, Source: "hello "},
		{Type: ScriptStart, Source: "{%"},
		{Type: Identifier, Source: "world"},
		{Type: ScriptEnd, Source: "%}"},
		{Type: EndOfFile, Source: ""},
	}

	l := NewTemplate(script)

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
