package lexer

import (
	"errors"
	"testing"
)

func TestParseStringWithMissingEnd(t *testing.T) {

	script := "\"string"

	l := NewScript(script)

	_, err := l.Read()
	if err == nil {
		t.Fatal("expect error, got nil")
	}

	var lerr *TokenError
	if !errors.As(err, &lerr) {
		t.Fatal("expect TokenError, got nil")
	}

	if lerr.Line != 1 {
		t.Fatalf("expect line 1, got %d", lerr.Line)
	}

	if lerr.Column != 1 {
		t.Fatalf("expect column 1, got %d", lerr.Column)
	}
}

func TestParseStringWithEscaped(t *testing.T) {

	script := "\"string \\\"escaped\\\"\""

	l := NewScript(script)

	tok, err := l.Read()
	if err != nil {
		t.Fatal(err)
	}

	if tok.Type != String {
		t.Fatalf("expect token type String, got %s", tok.Type)
	}

	if tok.Source != "\"string \\\"escaped\\\"\"" {
		t.Fatalf("expect token source \"string \\\"escaped\\\"\", got %s", tok.Source)
	}
}

func TestParseStringWithNewLine(t *testing.T) {

	script := "\"string\n\""

	l := NewScript(script)

	_, err := l.Read()
	if err == nil {
		t.Fatal("expect error, got nil")
	}

	var lerr *TokenError
	if !errors.As(err, &lerr) {
		t.Fatal("expect TokenError, got nil")
	}

	if lerr.Message != "unexpected character '\n'" {
		t.Fatalf("expect message \"unexpected character '\n'\", got \"%s\"", lerr.Message)
	}

	if lerr.Line != 1 {
		t.Fatalf("expect line 1, got %d", lerr.Line)
	}

	if lerr.Column != 8 {
		t.Fatalf("expect column 8, got %d", lerr.Column)
	}
}

func TestParseStringWithMatchingQuote(t *testing.T) {

	script := "\"string\""

	l := NewScript(script)

	tok, err := l.Read()
	if err != nil {
		t.Fatal(err)
	}

	if tok.Type != String {
		t.Fatalf("expect token type String, got %s", tok.Type)
	}

	if tok.Source != "\"string\"" {
		t.Fatalf("expect token source \"string\", got %s", tok.Source)
	}
}

func TestParseScriptTokens(t *testing.T) {

	want := []struct {
		source string
		tok    Token
	}{
		{"(", Token{Type: LeftParen, Source: "("}},
		{")", Token{Type: RightParen, Source: ")"}},
		{"{", Token{Type: LeftBrace, Source: "{"}},
		{"}", Token{Type: RightBrace, Source: "}"}},
		{"[", Token{Type: LeftBracket, Source: "["}},
		{"]", Token{Type: RightBracket, Source: "]"}},
		{",", Token{Type: Comma, Source: ","}},
		{":", Token{Type: Colon, Source: ":"}},
		{";", Token{Type: Semicolon, Source: ";"}},
		{"+", Token{Type: Plus, Source: "+"}},
		{"-", Token{Type: Minus, Source: "-"}},
		{"*", Token{Type: Asterisk, Source: "*"}},
		{".", Token{Type: Dot, Source: "."}},
		{"/", Token{Type: Slash, Source: "/"}},
		{"=", Token{Type: Equal, Source: "="}},
		{"==", Token{Type: Equals, Source: "=="}},
		{"!=", Token{Type: NotEqual, Source: "!="}},
		{"<", Token{Type: LessThan, Source: "<"}},
		{">", Token{Type: GreaterThan, Source: ">"}},
		{"<=", Token{Type: LessOrEqual, Source: "<="}},
		{">=", Token{Type: GreaterOrEqual, Source: ">="}},
		{"\"string\"", Token{Type: String, Source: "\"string\""}},
		{"let", Token{Type: Let, Source: "let"}},
		{"fn", Token{Type: Function, Source: "fn"}},
		{"return", Token{Type: Return, Source: "return"}},
		{"if", Token{Type: If, Source: "if"}},
		{"else", Token{Type: Else, Source: "else"}},
		{"foreach", Token{Type: Foreach, Source: "foreach"}},
		{"true", Token{Type: True, Source: "true"}},
		{"false", Token{Type: False, Source: "false"}},
		{"as", Token{Type: As, Source: "as"}},
		{"1", Token{Type: Integer, Source: "1"}},
		{"1.0", Token{Type: Float, Source: "1.0"}},
		{"!", Token{Type: Bang, Source: "!"}},
	}

	for _, w := range want {

		l := NewScript(w.source)

		tok, err := l.Read()
		if err != nil {
			t.Error(err)
		}

		if tok.Type != w.tok.Type {
			t.Errorf("want %s, got %s", w.tok.Type, tok.Type)
		}

		if tok.Source != w.tok.Source {
			t.Errorf("want %s, got %s", w.tok.Source, tok.Source)
		}
	}
}

func TestScriptLineAndColumn(t *testing.T) {

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

func TestTemplateLineAndColumn(t *testing.T) {

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
		{% foreach (items as item) { %}
			<li>{% item.name %}</li>
		{% } %}
	</ul>
`

	want := []Token{
		{Type: Text, Source: "\n\t<h1>", Line: 1, Column: 1},                // 0
		{Type: ScriptStart, Source: "{%", Line: 2, Column: 6},               // 1
		{Type: Identifier, Source: "title", Line: 2, Column: 9},             // 2
		{Type: ScriptEnd, Source: "%}", Line: 2, Column: 15},                // 3
		{Type: Text, Source: "</h1>\n\t<h2>Items (", Line: 2, Column: 17},   // 4
		{Type: ScriptStart, Source: "{%", Line: 3, Column: 13},              // 5
		{Type: Identifier, Source: "count", Line: 3, Column: 16},            // 6
		{Type: LeftParen, Source: "(", Line: 3, Column: 21},                 // 7
		{Type: Identifier, Source: "items", Line: 3, Column: 22},            // 8
		{Type: RightParen, Source: ")", Line: 3, Column: 27},                // 9
		{Type: ScriptEnd, Source: "%}", Line: 3, Column: 29},                // 10
		{Type: Text, Source: "): </h2>\n\t<ul>\n\t\t", Line: 3, Column: 31}, // 11
		{Type: ScriptStart, Source: "{%", Line: 5, Column: 3},               // 12
		{Type: Foreach, Source: "foreach", Line: 5, Column: 6},              // 13
		{Type: LeftParen, Source: "(", Line: 5, Column: 14},                 // 14
		{Type: Identifier, Source: "items", Line: 5, Column: 15},            // 15
		{Type: As, Source: "as", Line: 5, Column: 21},                       // 16
		{Type: Identifier, Source: "item", Line: 5, Column: 24},             // 17
		{Type: RightParen, Source: ")", Line: 5, Column: 28},                // 18
		{Type: LeftBrace, Source: "{", Line: 5, Column: 30},                 // 19
		{Type: ScriptEnd, Source: "%}", Line: 5, Column: 32},                // 20
		{Type: Text, Source: "\n\t\t\t<li>", Line: 5, Column: 34},           // 21
		{Type: ScriptStart, Source: "{%", Line: 6, Column: 8},               // 22
		{Type: Identifier, Source: "item", Line: 6, Column: 11},             // 23
		{Type: Dot, Source: ".", Line: 6, Column: 15},                       // 24
		{Type: Identifier, Source: "name", Line: 6, Column: 16},             // 25
		{Type: ScriptEnd, Source: "%}", Line: 6, Column: 21},                // 26
		{Type: Text, Source: "</li>\n\t\t", Line: 6, Column: 23},            // 27
		{Type: ScriptStart, Source: "{%", Line: 7, Column: 3},               // 28
		{Type: RightBrace, Source: "}", Line: 7, Column: 6},                 // 29
		{Type: ScriptEnd, Source: "%}", Line: 7, Column: 8},                 // 30
		{Type: Text, Source: "\n\t</ul>\n", Line: 7, Column: 10},            // 31
		{Type: EndOfFile, Source: "", Line: 8, Column: 1},                   // 32
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

	script := `hello {% world %}!`

	want := []Token{
		{Type: Text, Source: "hello "},
		{Type: ScriptStart, Source: "{%"},
		{Type: Identifier, Source: "world"},
		{Type: ScriptEnd, Source: "%}"},
		{Type: Text, Source: "!"},
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
