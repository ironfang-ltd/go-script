package lexer

import (
	"fmt"
	"strings"
	"testing"
)

// --- Helpers ---

func readAllTokens(l *Lexer) int {
	count := 0
	for {
		tok, err := l.Read()
		if err != nil {
			break
		}
		count++
		if tok.Type == EndOfFile {
			break
		}
	}
	return count
}

// --- Single token type benchmarks ---

func BenchmarkReadIdentifier(b *testing.B) {
	src := "variable"
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		l := NewScript(src)
		l.Read()
	}
}

func BenchmarkReadKeyword(b *testing.B) {
	src := "foreach"
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		l := NewScript(src)
		l.Read()
	}
}

func BenchmarkReadInteger(b *testing.B) {
	src := "1234567890"
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		l := NewScript(src)
		l.Read()
	}
}

func BenchmarkReadFloat(b *testing.B) {
	src := "3.14159265"
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		l := NewScript(src)
		l.Read()
	}
}

func BenchmarkReadString(b *testing.B) {
	src := `"hello world"`
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		l := NewScript(src)
		l.Read()
	}
}

func BenchmarkReadStringWithEscapes(b *testing.B) {
	src := `"hello \"world\" \\path\\to\\file\n"`
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		l := NewScript(src)
		l.Read()
	}
}

func BenchmarkReadOperatorSingle(b *testing.B) {
	src := "+"
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		l := NewScript(src)
		l.Read()
	}
}

func BenchmarkReadOperatorDouble(b *testing.B) {
	src := "=="
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		l := NewScript(src)
		l.Read()
	}
}

// --- Multi-token benchmarks ---

func BenchmarkReadLetStatement(b *testing.B) {
	src := "let x = 42;"
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		l := NewScript(src)
		readAllTokens(l)
	}
}

func BenchmarkReadFunctionDecl(b *testing.B) {
	src := "fn add(a, b) { return a + b; }"
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		l := NewScript(src)
		readAllTokens(l)
	}
}

func BenchmarkReadIfElse(b *testing.B) {
	src := "if (x > 0) { return true; } else { return false; }"
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		l := NewScript(src)
		readAllTokens(l)
	}
}

func BenchmarkReadForeach(b *testing.B) {
	src := `foreach (items as item) { print(item); }`
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		l := NewScript(src)
		readAllTokens(l)
	}
}

// --- Template mode benchmarks ---

func BenchmarkReadTemplateSimple(b *testing.B) {
	src := "Hello {% name %}, welcome!"
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		l := NewTemplate(src)
		readAllTokens(l)
	}
}

func BenchmarkReadTemplateTextOnly(b *testing.B) {
	src := strings.Repeat("Hello world, this is plain text. ", 100)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		l := NewTemplate(src)
		readAllTokens(l)
	}
}

func BenchmarkReadTemplateManyBlocks(b *testing.B) {
	var sb strings.Builder
	for i := 0; i < 100; i++ {
		fmt.Fprintf(&sb, "text %d {%% var%d %%} ", i, i)
	}
	src := sb.String()
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		l := NewTemplate(src)
		readAllTokens(l)
	}
}

// --- Scaling benchmarks ---

func BenchmarkReadIdentifiers10(b *testing.B) {
	src := strings.Repeat("variable ", 10)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		l := NewScript(src)
		readAllTokens(l)
	}
}

func BenchmarkReadIdentifiers100(b *testing.B) {
	src := strings.Repeat("variable ", 100)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		l := NewScript(src)
		readAllTokens(l)
	}
}

func BenchmarkReadIdentifiers1000(b *testing.B) {
	src := strings.Repeat("variable ", 1000)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		l := NewScript(src)
		readAllTokens(l)
	}
}

func BenchmarkReadMixedTokens(b *testing.B) {
	src := `let result = 0; fn add(a, b) { return a + b; }; if (result < 10) { result = add(result, 1); }`
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		l := NewScript(src)
		readAllTokens(l)
	}
}

// --- Constructor benchmarks ---

func BenchmarkNewScript(b *testing.B) {
	src := "let x = 1;"
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		NewScript(src)
	}
}

func BenchmarkNewTemplate(b *testing.B) {
	src := "Hello {% name %}!"
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		NewTemplate(src)
	}
}

// --- Whitespace handling benchmarks ---

func BenchmarkReadWhitespaceHeavy(b *testing.B) {
	src := "a   +   b   *   c   -   d   /   e"
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		l := NewScript(src)
		readAllTokens(l)
	}
}

func BenchmarkReadNewlineHeavy(b *testing.B) {
	var sb strings.Builder
	for i := 0; i < 100; i++ {
		fmt.Fprintf(&sb, "let v%d = %d;\n", i, i)
	}
	src := sb.String()
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		l := NewScript(src)
		readAllTokens(l)
	}
}

// --- Error path benchmarks ---

func BenchmarkReadUnexpectedChar(b *testing.B) {
	src := "#"
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		l := NewScript(src)
		l.Read()
	}
}

func BenchmarkReadUnterminatedString(b *testing.B) {
	src := `"unterminated`
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		l := NewScript(src)
		l.Read()
	}
}

func BenchmarkTokenErrorFormat(b *testing.B) {
	err := NewTokenError("test error", "let x = #;", 1, 9)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = err.Error()
	}
}

func BenchmarkReadLogicalOperators(b *testing.B) {
	src := "true && false || true"
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		l := NewScript(src)
		for {
			tok, _ := l.Read()
			if tok.Type == EndOfFile {
				break
			}
		}
	}
}

func BenchmarkReadNullCoalescing(b *testing.B) {
	src := "x ?? y ?? z"
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		l := NewScript(src)
		for {
			tok, _ := l.Read()
			if tok.Type == EndOfFile {
				break
			}
		}
	}
}

func BenchmarkReadSingleLineComment(b *testing.B) {
	src := "let x = 5; // this is a comment\nlet y = 10;"
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		l := NewScript(src)
		for {
			tok, _ := l.Read()
			if tok.Type == EndOfFile {
				break
			}
		}
	}
}

func BenchmarkReadMultiLineComment(b *testing.B) {
	src := "let x = 5; /* multi\nline\ncomment */ let y = 10;"
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		l := NewScript(src)
		for {
			tok, _ := l.Read()
			if tok.Type == EndOfFile {
				break
			}
		}
	}
}

func BenchmarkReadWhileStatement(b *testing.B) {
	src := "while (x < 10) { break; continue; }"
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		l := NewScript(src)
		for {
			tok, _ := l.Read()
			if tok.Type == EndOfFile {
				break
			}
		}
	}
}
