package lexer

import (
	"errors"
	"strings"
	"testing"
)

func TestParseStringWithMissingEnd(t *testing.T) {

	script := "\"string"

	l := NewScript(script)

	_, err := l.Read()
	if err == nil {
		t.Fatal("expect error, got nil")
	}

	var tokenErr *TokenError
	if !errors.As(err, &tokenErr) {
		t.Fatal("expect TokenError, got nil")
	}

	if tokenErr.Line != 1 {
		t.Fatalf("expect line 1, got %d", tokenErr.Line)
	}

	if tokenErr.Column != 1 {
		t.Fatalf("expect column 1, got %d", tokenErr.Column)
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

	var tokenErr *TokenError
	if !errors.As(err, &tokenErr) {
		t.Fatal("expect TokenError, got nil")
	}

	if tokenErr.Message != "unterminated string literal" {
		t.Fatalf("expect message \"unterminated string literal\", got \"%s\"", tokenErr.Message)
	}

	if tokenErr.Line != 1 {
		t.Fatalf("expect line 1, got %d", tokenErr.Line)
	}

	if tokenErr.Column != 1 {
		t.Fatalf("expect column 1, got %d", tokenErr.Column)
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
			var tokenErr *TokenError
			if errors.As(err, &tokenErr) {
				if tokenErr.Column != 9 {
					t.Fatalf("expect line 9, got %d\n", tokenErr.Column)
				}

				if tokenErr.Line != 2 {
					t.Fatalf("expect line 2, got %d\n", tokenErr.Line)
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
			var tokenErr *TokenError
			if errors.As(err, &tokenErr) {
				if tokenErr.Column != 10 {
					t.Fatalf("expect line 10, got %d\n", tokenErr.Column)
				}

				if tokenErr.Line != 1 {
					t.Fatalf("expect line 1, got %d\n", tokenErr.Line)
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
			var tokenErr *TokenError
			if errors.As(err, &tokenErr) {
				if tokenErr.Column != 11 {
					t.Fatalf("expect line 11, got %d\n", tokenErr.Column)
				}

				if tokenErr.Line != 2 {
					t.Fatalf("expect line 2, got %d\n", tokenErr.Line)
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

func TestFnCall(t *testing.T) {

	script := `test64(a,b)`

	want := []Token{
		{Type: Identifier, Source: "test64"},
		{Type: LeftParen, Source: "("},
		{Type: Identifier, Source: "a"},
		{Type: Comma, Source: ","},
		{Type: Identifier, Source: "b"},
		{Type: RightParen, Source: ")"},
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

// --- Regression tests for keyword boundary bug ---

func TestKeywordPrefixInIdentifier(t *testing.T) {
	tests := []struct {
		input string
		want  TokenType
	}{
		{"letter", Identifier},
		{"letters", Identifier},
		{"letting", Identifier},
		{"asset", Identifier},
		{"assign", Identifier},
		{"assume", Identifier},
		{"ifBlock", Identifier},
		{"iffy", Identifier},
		{"trueValue", Identifier},
		{"truer", Identifier},
		{"falsehood", Identifier},
		{"falsify", Identifier},
		{"fns", Identifier},
		{"fnCall", Identifier},
		{"returnValue", Identifier},
		{"returning", Identifier},
		{"elsewhere", Identifier},
		{"elseif", Identifier},
		{"foreach2", Identifier},
		{"foreachItem", Identifier},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			l := NewScript(tt.input)
			tok, err := l.Read()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tok.Type != tt.want {
				t.Fatalf("input %q: want token type %s, got %s (source: %q)", tt.input, tt.want, tok.Type, tok.Source)
			}
			if tok.Source != tt.input {
				t.Fatalf("input %q: want source %q, got %q", tt.input, tt.input, tok.Source)
			}
		})
	}
}

func TestKeywordsStillWork(t *testing.T) {
	tests := []struct {
		input string
		want  TokenType
	}{
		{"let", Let},
		{"fn", Function},
		{"return", Return},
		{"true", True},
		{"false", False},
		{"if", If},
		{"else", Else},
		{"foreach", Foreach},
		{"as", As},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			l := NewScript(tt.input)
			tok, err := l.Read()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tok.Type != tt.want {
				t.Fatalf("input %q: want token type %s, got %s", tt.input, tt.want, tok.Type)
			}
		})
	}
}

func TestKeywordFollowedBySpace(t *testing.T) {
	l := NewScript("let x = 1;")
	tok, _ := l.Read()
	if tok.Type != Let {
		t.Fatalf("want Let, got %s", tok.Type)
	}
	tok, _ = l.Read()
	if tok.Type != Identifier || tok.Source != "x" {
		t.Fatalf("want Identifier 'x', got %s %q", tok.Type, tok.Source)
	}
}

func TestKeywordFollowedByParen(t *testing.T) {
	l := NewScript("if(true)")

	want := []Token{
		{Type: If, Source: "if"},
		{Type: LeftParen, Source: "("},
		{Type: True, Source: "true"},
		{Type: RightParen, Source: ")"},
	}

	for _, w := range want {
		tok, err := l.Read()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if tok.Type != w.Type || tok.Source != w.Source {
			t.Fatalf("want %s %q, got %s %q", w.Type, w.Source, tok.Type, tok.Source)
		}
	}
}

// --- Regression tests for string escape bug ---

func TestStringDoubleBackslash(t *testing.T) {
	// Source: "he\\" — the string value is he\ and the closing quote terminates it
	script := `"he\\"`
	l := NewScript(script)
	tok, err := l.Read()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tok.Type != String {
		t.Fatalf("want String, got %s", tok.Type)
	}
	if tok.Source != `"he\\"` {
		t.Fatalf("want source %q, got %q", `"he\\"`, tok.Source)
	}
}

func TestStringDoubleBackslashFollowedByMore(t *testing.T) {
	// "hello\\world" — backslash-backslash then 'w'
	script := `"hello\\world"`
	l := NewScript(script)
	tok, err := l.Read()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tok.Type != String {
		t.Fatalf("want String, got %s", tok.Type)
	}
	if tok.Source != `"hello\\world"` {
		t.Fatalf("want source %q, got %q", `"hello\\world"`, tok.Source)
	}
}

func TestStringEscapedQuoteFollowedByBackslash(t *testing.T) {
	// "test\\\"end" — backslash-backslash, backslash-quote, then end
	script := `"test\\\"end"`
	l := NewScript(script)
	tok, err := l.Read()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tok.Type != String {
		t.Fatalf("want String, got %s", tok.Type)
	}
	if tok.Source != `"test\\\"end"` {
		t.Fatalf("want source %q, got %q", `"test\\\"end"`, tok.Source)
	}
}

func TestStringEscapeSequences(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		source string
	}{
		{"escaped_n", `"hello\nworld"`, `"hello\nworld"`},
		{"escaped_t", `"hello\tworld"`, `"hello\tworld"`},
		{"escaped_r", `"hello\rworld"`, `"hello\rworld"`},
		{"escaped_quote", `"say \"hi\""`, `"say \"hi\""`},
		{"escaped_backslash", `"path\\to"`, `"path\\to"`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewScript(tt.input)
			tok, err := l.Read()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tok.Type != String {
				t.Fatalf("want String, got %s", tok.Type)
			}
			if tok.Source != tt.source {
				t.Fatalf("want source %q, got %q", tt.source, tok.Source)
			}
		})
	}
}

func TestStringBackslashAtEndUnterminated(t *testing.T) {
	// "test\ — backslash at EOF without a following character
	script := `"test\`
	l := NewScript(script)
	_, err := l.Read()
	if err == nil {
		t.Fatal("expected error for unterminated string, got nil")
	}
	var tokenErr *TokenError
	if !errors.As(err, &tokenErr) {
		t.Fatal("expected TokenError")
	}
}

// --- Regression tests for number parsing bug ---

func TestNumberMultipleDotsProducesMultipleTokens(t *testing.T) {
	// 1.2.3 should be Float(1.2), Dot, Integer(3)
	l := NewScript("1.2.3")

	tok, err := l.Read()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tok.Type != Float || tok.Source != "1.2" {
		t.Fatalf("want Float '1.2', got %s %q", tok.Type, tok.Source)
	}

	tok, err = l.Read()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tok.Type != Dot {
		t.Fatalf("want Dot, got %s %q", tok.Type, tok.Source)
	}

	tok, err = l.Read()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tok.Type != Integer || tok.Source != "3" {
		t.Fatalf("want Integer '3', got %s %q", tok.Type, tok.Source)
	}
}

func TestNumberTrailingDotIsIntegerThenDot(t *testing.T) {
	// 1. should be Integer(1), Dot
	l := NewScript("1.")

	tok, err := l.Read()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tok.Type != Integer || tok.Source != "1" {
		t.Fatalf("want Integer '1', got %s %q", tok.Type, tok.Source)
	}

	tok, err = l.Read()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tok.Type != Dot {
		t.Fatalf("want Dot, got %s %q", tok.Type, tok.Source)
	}
}

func TestNumberPropertyAccess(t *testing.T) {
	// 1.toString should be Integer(1), Dot, Identifier(toString)
	l := NewScript("1.toString")

	want := []Token{
		{Type: Integer, Source: "1"},
		{Type: Dot, Source: "."},
		{Type: Identifier, Source: "toString"},
	}

	for _, w := range want {
		tok, err := l.Read()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if tok.Type != w.Type || tok.Source != w.Source {
			t.Fatalf("want %s %q, got %s %q", w.Type, w.Source, tok.Type, tok.Source)
		}
	}
}

func TestNumberVariants(t *testing.T) {
	tests := []struct {
		input     string
		wantType  TokenType
		wantValue string
	}{
		{"0", Integer, "0"},
		{"42", Integer, "42"},
		{"123456", Integer, "123456"},
		{"3.14", Float, "3.14"},
		{"0.5", Float, "0.5"},
		{"100.001", Float, "100.001"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			l := NewScript(tt.input)
			tok, err := l.Read()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tok.Type != tt.wantType {
				t.Fatalf("want type %s, got %s", tt.wantType, tok.Type)
			}
			if tok.Source != tt.wantValue {
				t.Fatalf("want source %q, got %q", tt.wantValue, tok.Source)
			}
		})
	}
}

// --- Identifier tests ---

func TestIdentifierVariants(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"x", "x"},
		{"myVar", "myVar"},
		{"_private", "_private"},
		{"__double", "__double"},
		{"camelCase", "camelCase"},
		{"PascalCase", "PascalCase"},
		{"snake_case", "snake_case"},
		{"name123", "name123"},
		{"a1b2c3", "a1b2c3"},
		{"_", "_"},
		{"_1", "_1"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			l := NewScript(tt.input)
			tok, err := l.Read()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tok.Type != Identifier {
				t.Fatalf("want Identifier, got %s", tok.Type)
			}
			if tok.Source != tt.want {
				t.Fatalf("want source %q, got %q", tt.want, tok.Source)
			}
		})
	}
}

// --- Operator tests ---

func TestMultiCharOperators(t *testing.T) {
	tests := []struct {
		input string
		want  TokenType
	}{
		{"==", Equals},
		{"!=", NotEqual},
		{"<=", LessOrEqual},
		{">=", GreaterOrEqual},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			l := NewScript(tt.input)
			tok, err := l.Read()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tok.Type != tt.want {
				t.Fatalf("want %s, got %s", tt.want, tok.Type)
			}
			if tok.Source != tt.input {
				t.Fatalf("want source %q, got %q", tt.input, tok.Source)
			}
		})
	}
}

func TestOperatorFollowedByOperand(t *testing.T) {
	// a == b
	l := NewScript("a == b")
	want := []Token{
		{Type: Identifier, Source: "a"},
		{Type: Equals, Source: "=="},
		{Type: Identifier, Source: "b"},
	}
	for _, w := range want {
		tok, _ := l.Read()
		if tok.Type != w.Type || tok.Source != w.Source {
			t.Fatalf("want %s %q, got %s %q", w.Type, w.Source, tok.Type, tok.Source)
		}
	}
}

// --- Empty/edge input tests ---

func TestEmptyScript(t *testing.T) {
	l := NewScript("")
	tok, err := l.Read()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tok.Type != EndOfFile {
		t.Fatalf("want EOF, got %s", tok.Type)
	}
}

func TestEmptyTemplate(t *testing.T) {
	l := NewTemplate("")
	tok, err := l.Read()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tok.Type != EndOfFile {
		t.Fatalf("want EOF, got %s", tok.Type)
	}
}

func TestTemplateTextOnly(t *testing.T) {
	l := NewTemplate("just plain text")
	tok, err := l.Read()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tok.Type != Text || tok.Source != "just plain text" {
		t.Fatalf("want Text 'just plain text', got %s %q", tok.Type, tok.Source)
	}
	tok, err = l.Read()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tok.Type != EndOfFile {
		t.Fatalf("want EOF, got %s", tok.Type)
	}
}

func TestWhitespaceOnlyScript(t *testing.T) {
	l := NewScript("   \t  ")
	tok, err := l.Read()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tok.Type != EndOfFile {
		t.Fatalf("want EOF, got %s", tok.Type)
	}
}

func TestNewlinesOnlyScript(t *testing.T) {
	l := NewScript("\n\n\n")
	tok, err := l.Read()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tok.Type != EndOfFile {
		t.Fatalf("want EOF, got %s", tok.Type)
	}
}

// --- Error tests ---

func TestUnexpectedCharacter(t *testing.T) {
	l := NewScript("#")
	_, err := l.Read()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var tokenErr *TokenError
	if !errors.As(err, &tokenErr) {
		t.Fatal("expected TokenError")
	}
	if tokenErr.Line != 1 || tokenErr.Column != 1 {
		t.Fatalf("want error at 1:1, got %d:%d", tokenErr.Line, tokenErr.Column)
	}
}

func TestUnterminatedStringEOF(t *testing.T) {
	l := NewScript(`"unterminated`)
	_, err := l.Read()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var tokenErr *TokenError
	if !errors.As(err, &tokenErr) {
		t.Fatal("expected TokenError")
	}
	if tokenErr.Message != "unterminated string literal" {
		t.Fatalf("want 'unterminated string literal', got %q", tokenErr.Message)
	}
}

func TestEmptyString(t *testing.T) {
	l := NewScript(`""`)
	tok, err := l.Read()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tok.Type != String || tok.Source != `""` {
		t.Fatalf("want String '\"\"', got %s %q", tok.Type, tok.Source)
	}
}

// --- Line/column tracking tests ---

func TestScriptMultilineTracking(t *testing.T) {
	script := "let a = 1;\nlet b = 2;\nlet c = 3;"

	l := NewScript(script)

	// Read all tokens, track the 'let' on each line
	letPositions := []struct {
		line int
		col  int
	}{
		{1, 1},
		{2, 1},
		{3, 1},
	}

	idx := 0
	for {
		tok, err := l.Read()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if tok.Type == EndOfFile {
			break
		}
		if tok.Type == Let {
			if idx >= len(letPositions) {
				t.Fatalf("more Let tokens than expected")
			}
			if tok.Line != letPositions[idx].line {
				t.Fatalf("let #%d: want line %d, got %d", idx, letPositions[idx].line, tok.Line)
			}
			if tok.Column != letPositions[idx].col {
				t.Fatalf("let #%d: want col %d, got %d", idx, letPositions[idx].col, tok.Column)
			}
			idx++
		}
	}
	if idx != len(letPositions) {
		t.Fatalf("expected %d Let tokens, got %d", len(letPositions), idx)
	}
}

// --- Template mode specific tests ---

func TestTemplateMultipleScriptBlocks(t *testing.T) {
	script := `Hello {% name %}, you have {% count %} items.`

	want := []Token{
		{Type: Text, Source: "Hello "},
		{Type: ScriptStart, Source: "{%"},
		{Type: Identifier, Source: "name"},
		{Type: ScriptEnd, Source: "%}"},
		{Type: Text, Source: ", you have "},
		{Type: ScriptStart, Source: "{%"},
		{Type: Identifier, Source: "count"},
		{Type: ScriptEnd, Source: "%}"},
		{Type: Text, Source: " items."},
		{Type: EndOfFile, Source: ""},
	}

	l := NewTemplate(script)
	for i, w := range want {
		tok, err := l.Read()
		if err != nil {
			t.Fatalf("[%d] unexpected error: %v", i, err)
		}
		if tok.Type != w.Type {
			t.Fatalf("[%d] want type %s, got %s", i, w.Type, tok.Type)
		}
		if tok.Source != w.Source {
			t.Fatalf("[%d] want source %q, got %q", i, w.Source, tok.Source)
		}
	}
}

func TestTemplateAdjacentScriptBlocks(t *testing.T) {
	script := `{% a %}{% b %}`

	want := []Token{
		{Type: ScriptStart, Source: "{%"},
		{Type: Identifier, Source: "a"},
		{Type: ScriptEnd, Source: "%}"},
		{Type: ScriptStart, Source: "{%"},
		{Type: Identifier, Source: "b"},
		{Type: ScriptEnd, Source: "%}"},
		{Type: EndOfFile, Source: ""},
	}

	l := NewTemplate(script)
	for i, w := range want {
		tok, err := l.Read()
		if err != nil {
			t.Fatalf("[%d] unexpected error: %v", i, err)
		}
		if tok.Type != w.Type {
			t.Fatalf("[%d] want type %s, got %s", i, w.Type, tok.Type)
		}
		if tok.Source != w.Source {
			t.Fatalf("[%d] want source %q, got %q", i, w.Source, tok.Source)
		}
	}
}

func TestTemplateScriptWithExpression(t *testing.T) {
	script := `Result: {% 1 + 2 %}`

	want := []Token{
		{Type: Text, Source: "Result: "},
		{Type: ScriptStart, Source: "{%"},
		{Type: Integer, Source: "1"},
		{Type: Plus, Source: "+"},
		{Type: Integer, Source: "2"},
		{Type: ScriptEnd, Source: "%}"},
		{Type: EndOfFile, Source: ""},
	}

	l := NewTemplate(script)
	for i, w := range want {
		tok, err := l.Read()
		if err != nil {
			t.Fatalf("[%d] unexpected error: %v", i, err)
		}
		if tok.Type != w.Type || tok.Source != w.Source {
			t.Fatalf("[%d] want %s %q, got %s %q", i, w.Type, w.Source, tok.Type, tok.Source)
		}
	}
}

// --- Modulo vs ScriptEnd in template mode ---

func TestModuloInTemplateScript(t *testing.T) {
	// x % y inside a template script block
	script := `{% x % y %}`

	want := []Token{
		{Type: ScriptStart, Source: "{%"},
		{Type: Identifier, Source: "x"},
		{Type: Modulo, Source: "%"},
		{Type: Identifier, Source: "y"},
		{Type: ScriptEnd, Source: "%}"},
		{Type: EndOfFile, Source: ""},
	}

	l := NewTemplate(script)
	for i, w := range want {
		tok, err := l.Read()
		if err != nil {
			t.Fatalf("[%d] unexpected error: %v", i, err)
		}
		if tok.Type != w.Type || tok.Source != w.Source {
			t.Fatalf("[%d] want %s %q, got %s %q", i, w.Type, w.Source, tok.Type, tok.Source)
		}
	}
}

// --- Complex expression tests ---

func TestComplexExpression(t *testing.T) {
	script := `fn add(a, b) { return a + b; }`

	want := []Token{
		{Type: Function, Source: "fn"},
		{Type: Identifier, Source: "add"},
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
	for i, w := range want {
		tok, err := l.Read()
		if err != nil {
			t.Fatalf("[%d] unexpected error: %v", i, err)
		}
		if tok.Type != w.Type || tok.Source != w.Source {
			t.Fatalf("[%d] want %s %q, got %s %q", i, w.Type, w.Source, tok.Type, tok.Source)
		}
	}
}

func TestHashLiteral(t *testing.T) {
	script := `{"key": "value", "num": 42}`

	want := []Token{
		{Type: LeftBrace, Source: "{"},
		{Type: String, Source: `"key"`},
		{Type: Colon, Source: ":"},
		{Type: String, Source: `"value"`},
		{Type: Comma, Source: ","},
		{Type: String, Source: `"num"`},
		{Type: Colon, Source: ":"},
		{Type: Integer, Source: "42"},
		{Type: RightBrace, Source: "}"},
		{Type: EndOfFile, Source: ""},
	}

	l := NewScript(script)
	for i, w := range want {
		tok, err := l.Read()
		if err != nil {
			t.Fatalf("[%d] unexpected error: %v", i, err)
		}
		if tok.Type != w.Type || tok.Source != w.Source {
			t.Fatalf("[%d] want %s %q, got %s %q", i, w.Type, w.Source, tok.Type, tok.Source)
		}
	}
}

func TestArrayLiteral(t *testing.T) {
	script := `[1, 2, 3]`

	want := []Token{
		{Type: LeftBracket, Source: "["},
		{Type: Integer, Source: "1"},
		{Type: Comma, Source: ","},
		{Type: Integer, Source: "2"},
		{Type: Comma, Source: ","},
		{Type: Integer, Source: "3"},
		{Type: RightBracket, Source: "]"},
		{Type: EndOfFile, Source: ""},
	}

	l := NewScript(script)
	for i, w := range want {
		tok, err := l.Read()
		if err != nil {
			t.Fatalf("[%d] unexpected error: %v", i, err)
		}
		if tok.Type != w.Type || tok.Source != w.Source {
			t.Fatalf("[%d] want %s %q, got %s %q", i, w.Type, w.Source, tok.Type, tok.Source)
		}
	}
}

func TestPropertyAccessChain(t *testing.T) {
	script := `a.b.c`

	want := []Token{
		{Type: Identifier, Source: "a"},
		{Type: Dot, Source: "."},
		{Type: Identifier, Source: "b"},
		{Type: Dot, Source: "."},
		{Type: Identifier, Source: "c"},
		{Type: EndOfFile, Source: ""},
	}

	l := NewScript(script)
	for i, w := range want {
		tok, err := l.Read()
		if err != nil {
			t.Fatalf("[%d] unexpected error: %v", i, err)
		}
		if tok.Type != w.Type || tok.Source != w.Source {
			t.Fatalf("[%d] want %s %q, got %s %q", i, w.Type, w.Source, tok.Type, tok.Source)
		}
	}
}

func TestIndexAccess(t *testing.T) {
	script := `arr[0]`

	want := []Token{
		{Type: Identifier, Source: "arr"},
		{Type: LeftBracket, Source: "["},
		{Type: Integer, Source: "0"},
		{Type: RightBracket, Source: "]"},
		{Type: EndOfFile, Source: ""},
	}

	l := NewScript(script)
	for i, w := range want {
		tok, err := l.Read()
		if err != nil {
			t.Fatalf("[%d] unexpected error: %v", i, err)
		}
		if tok.Type != w.Type || tok.Source != w.Source {
			t.Fatalf("[%d] want %s %q, got %s %q", i, w.Type, w.Source, tok.Type, tok.Source)
		}
	}
}

func TestForeachStatement(t *testing.T) {
	script := `foreach (items as item) { print(item); }`

	want := []Token{
		{Type: Foreach, Source: "foreach"},
		{Type: LeftParen, Source: "("},
		{Type: Identifier, Source: "items"},
		{Type: As, Source: "as"},
		{Type: Identifier, Source: "item"},
		{Type: RightParen, Source: ")"},
		{Type: LeftBrace, Source: "{"},
		{Type: Identifier, Source: "print"},
		{Type: LeftParen, Source: "("},
		{Type: Identifier, Source: "item"},
		{Type: RightParen, Source: ")"},
		{Type: Semicolon, Source: ";"},
		{Type: RightBrace, Source: "}"},
		{Type: EndOfFile, Source: ""},
	}

	l := NewScript(script)
	for i, w := range want {
		tok, err := l.Read()
		if err != nil {
			t.Fatalf("[%d] unexpected error: %v", i, err)
		}
		if tok.Type != w.Type || tok.Source != w.Source {
			t.Fatalf("[%d] want %s %q, got %s %q", i, w.Type, w.Source, tok.Type, tok.Source)
		}
	}
}

func TestIfElseStatement(t *testing.T) {
	script := `if (x > 0) { return true; } else { return false; }`

	want := []Token{
		{Type: If, Source: "if"},
		{Type: LeftParen, Source: "("},
		{Type: Identifier, Source: "x"},
		{Type: GreaterThan, Source: ">"},
		{Type: Integer, Source: "0"},
		{Type: RightParen, Source: ")"},
		{Type: LeftBrace, Source: "{"},
		{Type: Return, Source: "return"},
		{Type: True, Source: "true"},
		{Type: Semicolon, Source: ";"},
		{Type: RightBrace, Source: "}"},
		{Type: Else, Source: "else"},
		{Type: LeftBrace, Source: "{"},
		{Type: Return, Source: "return"},
		{Type: False, Source: "false"},
		{Type: Semicolon, Source: ";"},
		{Type: RightBrace, Source: "}"},
		{Type: EndOfFile, Source: ""},
	}

	l := NewScript(script)
	for i, w := range want {
		tok, err := l.Read()
		if err != nil {
			t.Fatalf("[%d] unexpected error: %v", i, err)
		}
		if tok.Type != w.Type || tok.Source != w.Source {
			t.Fatalf("[%d] want %s %q, got %s %q", i, w.Type, w.Source, tok.Type, tok.Source)
		}
	}
}

func TestAssignment(t *testing.T) {
	script := `x = 10;`

	want := []Token{
		{Type: Identifier, Source: "x"},
		{Type: Equal, Source: "="},
		{Type: Integer, Source: "10"},
		{Type: Semicolon, Source: ";"},
		{Type: EndOfFile, Source: ""},
	}

	l := NewScript(script)
	for i, w := range want {
		tok, err := l.Read()
		if err != nil {
			t.Fatalf("[%d] unexpected error: %v", i, err)
		}
		if tok.Type != w.Type || tok.Source != w.Source {
			t.Fatalf("[%d] want %s %q, got %s %q", i, w.Type, w.Source, tok.Type, tok.Source)
		}
	}
}

func TestNegativeNumber(t *testing.T) {
	// Negative numbers are prefix expression: Minus then Integer
	script := `-5`

	want := []Token{
		{Type: Minus, Source: "-"},
		{Type: Integer, Source: "5"},
		{Type: EndOfFile, Source: ""},
	}

	l := NewScript(script)
	for i, w := range want {
		tok, err := l.Read()
		if err != nil {
			t.Fatalf("[%d] unexpected error: %v", i, err)
		}
		if tok.Type != w.Type || tok.Source != w.Source {
			t.Fatalf("[%d] want %s %q, got %s %q", i, w.Type, w.Source, tok.Type, tok.Source)
		}
	}
}

func TestBangPrefix(t *testing.T) {
	script := `!true`

	want := []Token{
		{Type: Bang, Source: "!"},
		{Type: True, Source: "true"},
		{Type: EndOfFile, Source: ""},
	}

	l := NewScript(script)
	for i, w := range want {
		tok, err := l.Read()
		if err != nil {
			t.Fatalf("[%d] unexpected error: %v", i, err)
		}
		if tok.Type != w.Type || tok.Source != w.Source {
			t.Fatalf("[%d] want %s %q, got %s %q", i, w.Type, w.Source, tok.Type, tok.Source)
		}
	}
}

// --- ScriptEnd not recognized in pure script mode ---

func TestScriptModeIgnoresScriptEnd(t *testing.T) {
	// In pure script mode, %} should be Modulo then RightBrace, not ScriptEnd
	l := NewScript("%}")

	want := []Token{
		{Type: Modulo, Source: "%"},
		{Type: RightBrace, Source: "}"},
		{Type: EndOfFile, Source: ""},
	}

	for i, w := range want {
		tok, err := l.Read()
		if err != nil {
			t.Fatalf("[%d] unexpected error: %v", i, err)
		}
		if tok.Type != w.Type || tok.Source != w.Source {
			t.Fatalf("[%d] want %s %q, got %s %q", i, w.Type, w.Source, tok.Type, tok.Source)
		}
	}
}

func TestGetSource(t *testing.T) {
	input := "let x = 1;"
	l := NewScript(input)
	if l.GetSource() != input {
		t.Fatalf("want %q, got %q", input, l.GetSource())
	}
}

func TestConsecutiveReadsAfterEOF(t *testing.T) {
	l := NewScript("1")
	tok, _ := l.Read()
	if tok.Type != Integer {
		t.Fatalf("want Integer, got %s", tok.Type)
	}
	// Multiple EOF reads should be safe
	for i := 0; i < 3; i++ {
		tok, err := l.Read()
		if err != nil {
			t.Fatalf("unexpected error on EOF read %d: %v", i, err)
		}
		if tok.Type != EndOfFile {
			t.Fatalf("read %d: want EOF, got %s", i, tok.Type)
		}
	}
}

// --- TokenError.Error() formatting tests ---

func TestTokenErrorError(t *testing.T) {
	err := NewTokenError("unexpected character '#'", "let x = #;", 1, 9)
	msg := err.Error()

	if !strings.Contains(msg, "unexpected character '#'") {
		t.Fatalf("error message should contain the message text, got: %s", msg)
	}
	if !strings.Contains(msg, "line 1") {
		t.Fatalf("error message should contain line number, got: %s", msg)
	}
	if !strings.Contains(msg, "column 9") {
		t.Fatalf("error message should contain column number, got: %s", msg)
	}
	if !strings.Contains(msg, "let x = #;") {
		t.Fatalf("error message should contain source line, got: %s", msg)
	}
	if !strings.Contains(msg, "^") {
		t.Fatalf("error message should contain caret pointer, got: %s", msg)
	}
}

func TestTokenErrorErrorWithTabs(t *testing.T) {
	// Tab should be expanded to 4 spaces in the output
	err := NewTokenError("bad char", "\t#", 1, 2)
	msg := err.Error()

	// The tab in the source should be replaced with 4 spaces
	if !strings.Contains(msg, "    #") {
		t.Fatalf("tabs should be expanded to 4 spaces, got: %s", msg)
	}
}

func TestTokenErrorErrorMultiline(t *testing.T) {
	source := "let a = 1;\nlet b = #;"
	err := NewTokenError("unexpected character", source, 2, 9)
	msg := err.Error()

	if !strings.Contains(msg, "line 2") {
		t.Fatalf("should reference line 2, got: %s", msg)
	}
	// Should show the second line of source
	if !strings.Contains(msg, "let b = #;") {
		t.Fatalf("should contain the second source line, got: %s", msg)
	}
}

// --- Internal method guard branch tests (for 100% coverage) ---

func TestTryStringNotAtQuote(t *testing.T) {
	l := NewScript("abc")
	tok, ok, err := l.tryString('"', String)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Fatal("expected ok=false when not at a quote character")
	}
	if tok.Type != None {
		t.Fatalf("expected None token, got %s", tok.Type)
	}
}

func TestTryNumberNotAtDigit(t *testing.T) {
	l := NewScript("abc")
	tok, ok := l.tryNumber()
	if ok {
		t.Fatal("expected ok=false when not at a digit")
	}
	if tok != (Token{}) {
		t.Fatalf("expected zero Token, got %+v", tok)
	}
}

func TestTryIdentifierOrKeywordNotAtIdentStart(t *testing.T) {
	l := NewScript("123")
	tok, ok := l.tryIdentifierOrKeyword()
	if ok {
		t.Fatal("expected ok=false when not at an identifier start character")
	}
	if tok.Type != None {
		t.Fatalf("expected None token, got %s", tok.Type)
	}
}

// --- Additional edge case tests ---

func TestModuloInPureScript(t *testing.T) {
	l := NewScript("5 % 3")
	want := []Token{
		{Type: Integer, Source: "5"},
		{Type: Modulo, Source: "%"},
		{Type: Integer, Source: "3"},
		{Type: EndOfFile, Source: ""},
	}
	for i, w := range want {
		tok, err := l.Read()
		if err != nil {
			t.Fatalf("[%d] unexpected error: %v", i, err)
		}
		if tok.Type != w.Type || tok.Source != w.Source {
			t.Fatalf("[%d] want %s %q, got %s %q", i, w.Type, w.Source, tok.Type, tok.Source)
		}
	}
}

func TestScriptStartInPureScript(t *testing.T) {
	// In pure script mode, {% should be LeftBrace then Modulo
	l := NewScript("{%")
	want := []Token{
		{Type: LeftBrace, Source: "{"},
		{Type: Modulo, Source: "%"},
		{Type: EndOfFile, Source: ""},
	}
	for i, w := range want {
		tok, err := l.Read()
		if err != nil {
			t.Fatalf("[%d] unexpected error: %v", i, err)
		}
		if tok.Type != w.Type || tok.Source != w.Source {
			t.Fatalf("[%d] want %s %q, got %s %q", i, w.Type, w.Source, tok.Type, tok.Source)
		}
	}
}

func TestTemplateStartsWithScriptBlock(t *testing.T) {
	l := NewTemplate("{% x %}tail")
	want := []Token{
		{Type: ScriptStart, Source: "{%"},
		{Type: Identifier, Source: "x"},
		{Type: ScriptEnd, Source: "%}"},
		{Type: Text, Source: "tail"},
		{Type: EndOfFile, Source: ""},
	}
	for i, w := range want {
		tok, err := l.Read()
		if err != nil {
			t.Fatalf("[%d] unexpected error: %v", i, err)
		}
		if tok.Type != w.Type || tok.Source != w.Source {
			t.Fatalf("[%d] want %s %q, got %s %q", i, w.Type, w.Source, tok.Type, tok.Source)
		}
	}
}

func TestTemplateEndsWithScriptBlock(t *testing.T) {
	l := NewTemplate("head{% x %}")
	want := []Token{
		{Type: Text, Source: "head"},
		{Type: ScriptStart, Source: "{%"},
		{Type: Identifier, Source: "x"},
		{Type: ScriptEnd, Source: "%}"},
		{Type: EndOfFile, Source: ""},
	}
	for i, w := range want {
		tok, err := l.Read()
		if err != nil {
			t.Fatalf("[%d] unexpected error: %v", i, err)
		}
		if tok.Type != w.Type || tok.Source != w.Source {
			t.Fatalf("[%d] want %s %q, got %s %q", i, w.Type, w.Source, tok.Type, tok.Source)
		}
	}
}

func TestTemplateEmptyScriptBlock(t *testing.T) {
	l := NewTemplate("{% %}")
	want := []Token{
		{Type: ScriptStart, Source: "{%"},
		{Type: ScriptEnd, Source: "%}"},
		{Type: EndOfFile, Source: ""},
	}
	for i, w := range want {
		tok, err := l.Read()
		if err != nil {
			t.Fatalf("[%d] unexpected error: %v", i, err)
		}
		if tok.Type != w.Type || tok.Source != w.Source {
			t.Fatalf("[%d] want %s %q, got %s %q", i, w.Type, w.Source, tok.Type, tok.Source)
		}
	}
}

func TestTemplateOnlyScriptBlock(t *testing.T) {
	l := NewTemplate("{% 42 %}")
	want := []Token{
		{Type: ScriptStart, Source: "{%"},
		{Type: Integer, Source: "42"},
		{Type: ScriptEnd, Source: "%}"},
		{Type: EndOfFile, Source: ""},
	}
	for i, w := range want {
		tok, err := l.Read()
		if err != nil {
			t.Fatalf("[%d] unexpected error: %v", i, err)
		}
		if tok.Type != w.Type || tok.Source != w.Source {
			t.Fatalf("[%d] want %s %q, got %s %q", i, w.Type, w.Source, tok.Type, tok.Source)
		}
	}
}

func TestScriptNewlinesBetweenTokens(t *testing.T) {
	l := NewScript("a\n+\nb")
	want := []Token{
		{Type: Identifier, Source: "a", Line: 1, Column: 1},
		{Type: Plus, Source: "+", Line: 2, Column: 1},
		{Type: Identifier, Source: "b", Line: 3, Column: 1},
		{Type: EndOfFile, Source: ""},
	}
	for i, w := range want {
		tok, err := l.Read()
		if err != nil {
			t.Fatalf("[%d] unexpected error: %v", i, err)
		}
		if tok.Type != w.Type || tok.Source != w.Source {
			t.Fatalf("[%d] want %s %q, got %s %q", i, w.Type, w.Source, tok.Type, tok.Source)
		}
		if w.Line != 0 && tok.Line != w.Line {
			t.Fatalf("[%d] want line %d, got %d", i, w.Line, tok.Line)
		}
		if w.Column != 0 && tok.Column != w.Column {
			t.Fatalf("[%d] want column %d, got %d", i, w.Column, tok.Column)
		}
	}
}

func TestCarriageReturnIsWhitespace(t *testing.T) {
	l := NewScript("a\r+")
	want := []Token{
		{Type: Identifier, Source: "a"},
		{Type: Plus, Source: "+"},
		{Type: EndOfFile, Source: ""},
	}
	for i, w := range want {
		tok, err := l.Read()
		if err != nil {
			t.Fatalf("[%d] unexpected error: %v", i, err)
		}
		if tok.Type != w.Type || tok.Source != w.Source {
			t.Fatalf("[%d] want %s %q, got %s %q", i, w.Type, w.Source, tok.Type, tok.Source)
		}
	}
}

func TestTabIsWhitespace(t *testing.T) {
	l := NewScript("a\t+")
	want := []Token{
		{Type: Identifier, Source: "a"},
		{Type: Plus, Source: "+"},
		{Type: EndOfFile, Source: ""},
	}
	for i, w := range want {
		tok, err := l.Read()
		if err != nil {
			t.Fatalf("[%d] unexpected error: %v", i, err)
		}
		if tok.Type != w.Type || tok.Source != w.Source {
			t.Fatalf("[%d] want %s %q, got %s %q", i, w.Type, w.Source, tok.Type, tok.Source)
		}
	}
}

func TestMultipleUnexpectedCharacters(t *testing.T) {
	tests := []struct {
		input string
		char  string
	}{
		{"#", "#"},
		{"@", "@"},
		{"$", "$"},
		{"~", "~"},
		{"&", "&"},
		{"|", "|"},
		{"^", "^"},
		{"?", "?"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			l := NewScript(tt.input)
			_, err := l.Read()
			if err == nil {
				t.Fatalf("expected error for character %q, got nil", tt.input)
			}
			var tokenErr *TokenError
			if !errors.As(err, &tokenErr) {
				t.Fatal("expected TokenError")
			}
			if !strings.Contains(tokenErr.Message, tt.char) {
				t.Fatalf("error message should mention character %q, got: %s", tt.char, tokenErr.Message)
			}
		})
	}
}

func TestLargeInteger(t *testing.T) {
	l := NewScript("99999999999999")
	tok, err := l.Read()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tok.Type != Integer || tok.Source != "99999999999999" {
		t.Fatalf("want Integer '99999999999999', got %s %q", tok.Type, tok.Source)
	}
}

func TestStringWithOnlyEscapes(t *testing.T) {
	l := NewScript(`"\\\\"`)
	tok, err := l.Read()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tok.Type != String {
		t.Fatalf("want String, got %s", tok.Type)
	}
	if tok.Source != `"\\\\"` {
		t.Fatalf("want source %q, got %q", `"\\\\"`, tok.Source)
	}
}

func TestNestedBraces(t *testing.T) {
	l := NewScript("{{}}")
	want := []Token{
		{Type: LeftBrace, Source: "{"},
		{Type: LeftBrace, Source: "{"},
		{Type: RightBrace, Source: "}"},
		{Type: RightBrace, Source: "}"},
		{Type: EndOfFile, Source: ""},
	}
	for i, w := range want {
		tok, err := l.Read()
		if err != nil {
			t.Fatalf("[%d] unexpected error: %v", i, err)
		}
		if tok.Type != w.Type || tok.Source != w.Source {
			t.Fatalf("[%d] want %s %q, got %s %q", i, w.Type, w.Source, tok.Type, tok.Source)
		}
	}
}

func TestNestedBrackets(t *testing.T) {
	l := NewScript("a[b[0]]")
	want := []Token{
		{Type: Identifier, Source: "a"},
		{Type: LeftBracket, Source: "["},
		{Type: Identifier, Source: "b"},
		{Type: LeftBracket, Source: "["},
		{Type: Integer, Source: "0"},
		{Type: RightBracket, Source: "]"},
		{Type: RightBracket, Source: "]"},
		{Type: EndOfFile, Source: ""},
	}
	for i, w := range want {
		tok, err := l.Read()
		if err != nil {
			t.Fatalf("[%d] unexpected error: %v", i, err)
		}
		if tok.Type != w.Type || tok.Source != w.Source {
			t.Fatalf("[%d] want %s %q, got %s %q", i, w.Type, w.Source, tok.Type, tok.Source)
		}
	}
}

func TestAllComparisonOperators(t *testing.T) {
	l := NewScript("a < b > c <= d >= e == f != g")
	want := []Token{
		{Type: Identifier, Source: "a"},
		{Type: LessThan, Source: "<"},
		{Type: Identifier, Source: "b"},
		{Type: GreaterThan, Source: ">"},
		{Type: Identifier, Source: "c"},
		{Type: LessOrEqual, Source: "<="},
		{Type: Identifier, Source: "d"},
		{Type: GreaterOrEqual, Source: ">="},
		{Type: Identifier, Source: "e"},
		{Type: Equals, Source: "=="},
		{Type: Identifier, Source: "f"},
		{Type: NotEqual, Source: "!="},
		{Type: Identifier, Source: "g"},
		{Type: EndOfFile, Source: ""},
	}
	for i, w := range want {
		tok, err := l.Read()
		if err != nil {
			t.Fatalf("[%d] unexpected error: %v", i, err)
		}
		if tok.Type != w.Type || tok.Source != w.Source {
			t.Fatalf("[%d] want %s %q, got %s %q", i, w.Type, w.Source, tok.Type, tok.Source)
		}
	}
}

func TestAllArithmeticOperators(t *testing.T) {
	l := NewScript("a + b - c * d / e % f")
	want := []Token{
		{Type: Identifier, Source: "a"},
		{Type: Plus, Source: "+"},
		{Type: Identifier, Source: "b"},
		{Type: Minus, Source: "-"},
		{Type: Identifier, Source: "c"},
		{Type: Asterisk, Source: "*"},
		{Type: Identifier, Source: "d"},
		{Type: Slash, Source: "/"},
		{Type: Identifier, Source: "e"},
		{Type: Modulo, Source: "%"},
		{Type: Identifier, Source: "f"},
		{Type: EndOfFile, Source: ""},
	}
	for i, w := range want {
		tok, err := l.Read()
		if err != nil {
			t.Fatalf("[%d] unexpected error: %v", i, err)
		}
		if tok.Type != w.Type || tok.Source != w.Source {
			t.Fatalf("[%d] want %s %q, got %s %q", i, w.Type, w.Source, tok.Type, tok.Source)
		}
	}
}

func TestTemplateWithNewlinesInText(t *testing.T) {
	l := NewTemplate("line1\nline2\n{% x %}\nline4")
	want := []Token{
		{Type: Text, Source: "line1\nline2\n"},
		{Type: ScriptStart, Source: "{%"},
		{Type: Identifier, Source: "x"},
		{Type: ScriptEnd, Source: "%}"},
		{Type: Text, Source: "\nline4"},
		{Type: EndOfFile, Source: ""},
	}
	for i, w := range want {
		tok, err := l.Read()
		if err != nil {
			t.Fatalf("[%d] unexpected error: %v", i, err)
		}
		if tok.Type != w.Type || tok.Source != w.Source {
			t.Fatalf("[%d] want %s %q, got %s %q", i, w.Type, w.Source, tok.Type, tok.Source)
		}
	}
}

func TestTemplatePartialScriptStart(t *testing.T) {
	// A lone { not followed by % should be text
	l := NewTemplate("hello { world")
	tok, err := l.Read()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tok.Type != Text || tok.Source != "hello { world" {
		t.Fatalf("want Text 'hello { world', got %s %q", tok.Type, tok.Source)
	}
}

func TestTemplateTrailingOpenBrace(t *testing.T) {
	// A { at end of input (not followed by %) should be text
	l := NewTemplate("hello {")
	tok, err := l.Read()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tok.Type != Text || tok.Source != "hello {" {
		t.Fatalf("want Text 'hello {', got %s %q", tok.Type, tok.Source)
	}
}

func TestStringInScriptBlock(t *testing.T) {
	l := NewTemplate(`{% "hello world" %}`)
	want := []Token{
		{Type: ScriptStart, Source: "{%"},
		{Type: String, Source: `"hello world"`},
		{Type: ScriptEnd, Source: "%}"},
		{Type: EndOfFile, Source: ""},
	}
	for i, w := range want {
		tok, err := l.Read()
		if err != nil {
			t.Fatalf("[%d] unexpected error: %v", i, err)
		}
		if tok.Type != w.Type || tok.Source != w.Source {
			t.Fatalf("[%d] want %s %q, got %s %q", i, w.Type, w.Source, tok.Type, tok.Source)
		}
	}
}

func TestEqualNotFollowedByEqual(t *testing.T) {
	// = at end of input
	l := NewScript("=")
	tok, err := l.Read()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tok.Type != Equal || tok.Source != "=" {
		t.Fatalf("want Equal '=', got %s %q", tok.Type, tok.Source)
	}
}

func TestBangNotFollowedByEqual(t *testing.T) {
	// ! at end of input
	l := NewScript("!")
	tok, err := l.Read()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tok.Type != Bang || tok.Source != "!" {
		t.Fatalf("want Bang '!', got %s %q", tok.Type, tok.Source)
	}
}

func TestLessThanNotFollowedByEqual(t *testing.T) {
	// < at end of input
	l := NewScript("<")
	tok, err := l.Read()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tok.Type != LessThan || tok.Source != "<" {
		t.Fatalf("want LessThan '<', got %s %q", tok.Type, tok.Source)
	}
}

func TestGreaterThanNotFollowedByEqual(t *testing.T) {
	// > at end of input
	l := NewScript(">")
	tok, err := l.Read()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tok.Type != GreaterThan || tok.Source != ">" {
		t.Fatalf("want GreaterThan '>', got %s %q", tok.Type, tok.Source)
	}
}

func TestModuloAtEndOfInput(t *testing.T) {
	// % at end of input in pure script mode
	l := NewScript("%")
	tok, err := l.Read()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tok.Type != Modulo || tok.Source != "%" {
		t.Fatalf("want Modulo '%%', got %s %q", tok.Type, tok.Source)
	}
}

func TestModuloNotFollowedByBrace(t *testing.T) {
	// In template mode, % not followed by } should be Modulo
	l := NewTemplate("{% % x %}")
	want := []Token{
		{Type: ScriptStart, Source: "{%"},
		{Type: Modulo, Source: "%"},
		{Type: Identifier, Source: "x"},
		{Type: ScriptEnd, Source: "%}"},
		{Type: EndOfFile, Source: ""},
	}
	for i, w := range want {
		tok, err := l.Read()
		if err != nil {
			t.Fatalf("[%d] unexpected error: %v", i, err)
		}
		if tok.Type != w.Type || tok.Source != w.Source {
			t.Fatalf("[%d] want %s %q, got %s %q", i, w.Type, w.Source, tok.Type, tok.Source)
		}
	}
}

func TestSingleLineComment(t *testing.T) {
	script := "let x = 5; // this is a comment\nlet y = 10;"

	l := NewScript(script)

	expected := []Token{
		{Type: Let, Source: "let"},
		{Type: Identifier, Source: "x"},
		{Type: Equal, Source: "="},
		{Type: Integer, Source: "5"},
		{Type: Semicolon, Source: ";"},
		{Type: Let, Source: "let"},
		{Type: Identifier, Source: "y"},
		{Type: Equal, Source: "="},
		{Type: Integer, Source: "10"},
		{Type: Semicolon, Source: ";"},
		{Type: EndOfFile, Source: ""},
	}

	for i, exp := range expected {
		tok, err := l.Read()
		if err != nil {
			t.Fatalf("[%d] unexpected error: %v", i, err)
		}
		if tok.Type != exp.Type || tok.Source != exp.Source {
			t.Fatalf("[%d] want %s %q, got %s %q", i, exp.Type, exp.Source, tok.Type, tok.Source)
		}
	}
}

func TestMultiLineComment(t *testing.T) {
	script := "let x = 5; /* this is\na multi-line\ncomment */ let y = 10;"

	l := NewScript(script)

	expected := []Token{
		{Type: Let, Source: "let"},
		{Type: Identifier, Source: "x"},
		{Type: Equal, Source: "="},
		{Type: Integer, Source: "5"},
		{Type: Semicolon, Source: ";"},
		{Type: Let, Source: "let"},
		{Type: Identifier, Source: "y"},
		{Type: Equal, Source: "="},
		{Type: Integer, Source: "10"},
		{Type: Semicolon, Source: ";"},
		{Type: EndOfFile, Source: ""},
	}

	for i, exp := range expected {
		tok, err := l.Read()
		if err != nil {
			t.Fatalf("[%d] unexpected error: %v", i, err)
		}
		if tok.Type != exp.Type || tok.Source != exp.Source {
			t.Fatalf("[%d] want %s %q, got %s %q", i, exp.Type, exp.Source, tok.Type, tok.Source)
		}
	}
}

func TestCommentOnlyLine(t *testing.T) {
	script := "// entire line is a comment"

	l := NewScript(script)

	tok, err := l.Read()
	if err != nil {
		t.Fatal(err)
	}
	if tok.Type != EndOfFile {
		t.Fatalf("expected EOF, got %s %q", tok.Type, tok.Source)
	}
}

func TestLogicalAndToken(t *testing.T) {
	script := "true && false"

	l := NewScript(script)

	expected := []Token{
		{Type: True, Source: "true"},
		{Type: And, Source: "&&"},
		{Type: False, Source: "false"},
		{Type: EndOfFile, Source: ""},
	}

	for i, exp := range expected {
		tok, err := l.Read()
		if err != nil {
			t.Fatalf("[%d] unexpected error: %v", i, err)
		}
		if tok.Type != exp.Type || tok.Source != exp.Source {
			t.Fatalf("[%d] want %s %q, got %s %q", i, exp.Type, exp.Source, tok.Type, tok.Source)
		}
	}
}

func TestLogicalOrToken(t *testing.T) {
	script := "true || false"

	l := NewScript(script)

	expected := []Token{
		{Type: True, Source: "true"},
		{Type: Or, Source: "||"},
		{Type: False, Source: "false"},
		{Type: EndOfFile, Source: ""},
	}

	for i, exp := range expected {
		tok, err := l.Read()
		if err != nil {
			t.Fatalf("[%d] unexpected error: %v", i, err)
		}
		if tok.Type != exp.Type || tok.Source != exp.Source {
			t.Fatalf("[%d] want %s %q, got %s %q", i, exp.Type, exp.Source, tok.Type, tok.Source)
		}
	}
}

func TestNullCoalescingToken(t *testing.T) {
	script := "x ?? 5"

	l := NewScript(script)

	expected := []Token{
		{Type: Identifier, Source: "x"},
		{Type: NullCoalescing, Source: "??"},
		{Type: Integer, Source: "5"},
		{Type: EndOfFile, Source: ""},
	}

	for i, exp := range expected {
		tok, err := l.Read()
		if err != nil {
			t.Fatalf("[%d] unexpected error: %v", i, err)
		}
		if tok.Type != exp.Type || tok.Source != exp.Source {
			t.Fatalf("[%d] want %s %q, got %s %q", i, exp.Type, exp.Source, tok.Type, tok.Source)
		}
	}
}

func TestWhileBreakContinueTokens(t *testing.T) {
	script := "while (true) { break; continue; }"

	l := NewScript(script)

	expected := []Token{
		{Type: While, Source: "while"},
		{Type: LeftParen, Source: "("},
		{Type: True, Source: "true"},
		{Type: RightParen, Source: ")"},
		{Type: LeftBrace, Source: "{"},
		{Type: Break, Source: "break"},
		{Type: Semicolon, Source: ";"},
		{Type: Continue, Source: "continue"},
		{Type: Semicolon, Source: ";"},
		{Type: RightBrace, Source: "}"},
		{Type: EndOfFile, Source: ""},
	}

	for i, exp := range expected {
		tok, err := l.Read()
		if err != nil {
			t.Fatalf("[%d] unexpected error: %v", i, err)
		}
		if tok.Type != exp.Type || tok.Source != exp.Source {
			t.Fatalf("[%d] want %s %q, got %s %q", i, exp.Type, exp.Source, tok.Type, tok.Source)
		}
	}
}

func TestNullToken(t *testing.T) {
	script := "null"

	l := NewScript(script)

	tok, err := l.Read()
	if err != nil {
		t.Fatal(err)
	}
	if tok.Type != Null {
		t.Fatalf("expected NULL, got %s %q", tok.Type, tok.Source)
	}
}

func TestSingleAmpersandError(t *testing.T) {
	script := "x & y"

	l := NewScript(script)

	_, _ = l.Read() // skip 'x'

	_, err := l.Read()
	if err == nil {
		t.Fatal("expected error for single '&', got nil")
	}
}

func TestSinglePipeError(t *testing.T) {
	script := "x | y"

	l := NewScript(script)

	_, _ = l.Read() // skip 'x'

	_, err := l.Read()
	if err == nil {
		t.Fatal("expected error for single '|', got nil")
	}
}

func TestSingleQuestionMarkError(t *testing.T) {
	script := "x ? y"

	l := NewScript(script)

	_, _ = l.Read() // skip 'x'

	_, err := l.Read()
	if err == nil {
		t.Fatal("expected error for single '?', got nil")
	}
}

func TestSlashStillWorksDivision(t *testing.T) {
	script := "10 / 2"

	l := NewScript(script)

	expected := []Token{
		{Type: Integer, Source: "10"},
		{Type: Slash, Source: "/"},
		{Type: Integer, Source: "2"},
		{Type: EndOfFile, Source: ""},
	}

	for i, exp := range expected {
		tok, err := l.Read()
		if err != nil {
			t.Fatalf("[%d] unexpected error: %v", i, err)
		}
		if tok.Type != exp.Type || tok.Source != exp.Source {
			t.Fatalf("[%d] want %s %q, got %s %q", i, exp.Type, exp.Source, tok.Type, tok.Source)
		}
	}
}

// --- Bug fix verification tests ---

func TestUnterminatedMultiLineCommentAtEOF(t *testing.T) {
	l := NewScript("let x = 5; /* unterminated")
	for {
		tok, err := l.Read()
		if err != nil {
			var tokenErr *TokenError
			if !errors.As(err, &tokenErr) {
				t.Fatalf("expected TokenError, got %T: %v", err, err)
			}
			if !strings.Contains(tokenErr.Message, "unterminated comment") {
				t.Fatalf("expected 'unterminated comment' error, got: %s", tokenErr.Message)
			}
			return
		}
		if tok.Type == EndOfFile {
			t.Fatal("expected error for unterminated comment, but got EOF")
		}
	}
}

func TestTokenErrorRustStyleFormat(t *testing.T) {
	err := NewTokenError("unexpected character '#'", "let x = #;", 1, 9)
	msg := err.Error()

	// Verify Rust-style format elements
	if !strings.Contains(msg, "error: unexpected character '#'") {
		t.Fatalf("should start with 'error:' prefix, got: %s", msg)
	}
	if !strings.Contains(msg, "-->") {
		t.Fatalf("should contain '-->' location indicator, got: %s", msg)
	}
	if !strings.Contains(msg, "|") {
		t.Fatalf("should contain pipe separators, got: %s", msg)
	}
	if !strings.Contains(msg, "1 | let x = #;") {
		t.Fatalf("should contain line number with pipe, got: %s", msg)
	}
}

func TestTokenErrorBoundsCheck(t *testing.T) {
	// Line number out of range should not panic
	err := NewTokenError("test error", "single line", 99, 1)
	msg := err.Error()
	if !strings.Contains(msg, "error: test error") {
		t.Fatalf("should contain error message, got: %s", msg)
	}
}

func TestTokenErrorLineZero(t *testing.T) {
	// Line 0 should not panic
	err := NewTokenError("test error", "single line", 0, 1)
	msg := err.Error()
	if !strings.Contains(msg, "error: test error") {
		t.Fatalf("should contain error message, got: %s", msg)
	}
}

// --- Compound Assignment Operators ---

func TestLexCompoundAssignmentOperators(t *testing.T) {
	tests := []struct {
		input    string
		expected []TokenType
	}{
		{"x += 1", []TokenType{Identifier, PlusEqual, Integer}},
		{"x -= 1", []TokenType{Identifier, MinusEqual, Integer}},
		{"x *= 1", []TokenType{Identifier, AsteriskEqual, Integer}},
		{"x /= 1", []TokenType{Identifier, SlashEqual, Integer}},
		{"x %= 1", []TokenType{Identifier, ModuloEqual, Integer}},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			l := NewScript(tt.input)
			for i, expectedType := range tt.expected {
				tok, err := l.Read()
				if err != nil {
					t.Fatalf("token %d: unexpected error: %v", i, err)
				}
				if tok.Type != expectedType {
					t.Fatalf("token %d: expected %s, got %s (%q)", i, expectedType, tok.Type, tok.Source)
				}
			}
		})
	}
}

func TestLexCompoundAssignmentDoesNotConflictWithComments(t *testing.T) {
	// /= should not conflict with // or /*
	input := `// comment
x /= 2`
	l := NewScript(input)
	tok, err := l.Read()
	if err != nil {
		t.Fatal(err)
	}
	if tok.Type != Identifier {
		t.Fatalf("expected Identifier, got %s", tok.Type)
	}
	tok, err = l.Read()
	if err != nil {
		t.Fatal(err)
	}
	if tok.Type != SlashEqual {
		t.Fatalf("expected SlashEqual, got %s", tok.Type)
	}
}

func TestLexCompoundAssignmentModuloInTemplate(t *testing.T) {
	// %= should not conflict with %} in template mode
	input := `{% x %= 2; %}`
	l := NewTemplate(input)
	tok, err := l.Read() // {%
	if err != nil {
		t.Fatal(err)
	}
	if tok.Type != ScriptStart {
		t.Fatalf("expected ScriptStart, got %s", tok.Type)
	}
	tok, err = l.Read() // x
	if err != nil {
		t.Fatal(err)
	}
	if tok.Type != Identifier {
		t.Fatalf("expected Identifier, got %s", tok.Type)
	}
	tok, err = l.Read() // %=
	if err != nil {
		t.Fatal(err)
	}
	if tok.Type != ModuloEqual {
		t.Fatalf("expected ModuloEqual, got %s (%q)", tok.Type, tok.Source)
	}
}

func TestLexCompoundAssignmentSources(t *testing.T) {
	tests := []struct {
		input    string
		expected string
		typ      TokenType
	}{
		{"+=", "+=", PlusEqual},
		{"-=", "-=", MinusEqual},
		{"*=", "*=", AsteriskEqual},
		{"/=", "/=", SlashEqual},
		{"%=", "%=", ModuloEqual},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			l := NewScript(tt.input)
			tok, err := l.Read()
			if err != nil {
				t.Fatal(err)
			}
			if tok.Source != tt.expected {
				t.Fatalf("expected source %q, got %q", tt.expected, tok.Source)
			}
			if tok.Type != tt.typ {
				t.Fatalf("expected type %s, got %s", tt.typ, tok.Type)
			}
		})
	}
}
