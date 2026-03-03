package parser

import (
	"fmt"
	"strings"
	"testing"

	"github.com/ironfang-ltd/go-script/lexer"
)

// --- Helpers ---

func parseScript(src string) (*Program, error) {
	l := lexer.NewScript(src)
	p := New(l)
	return p.Parse()
}

func parseTemplate(src string) (*Program, error) {
	l := lexer.NewTemplate(src)
	p := New(l)
	return p.Parse()
}

// --- Single statement benchmarks ---

func BenchmarkParseLetStatement(b *testing.B) {
	src := "let x = 42;"
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		parseScript(src)
	}
}

func BenchmarkParseReturnStatement(b *testing.B) {
	src := "return 42;"
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		parseScript(src)
	}
}

func BenchmarkParseExpressionStatement(b *testing.B) {
	src := "x;"
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		parseScript(src)
	}
}

// --- Expression benchmarks ---

func BenchmarkParseIntegerLiteral(b *testing.B) {
	src := "42;"
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		parseScript(src)
	}
}

func BenchmarkParseFloatLiteral(b *testing.B) {
	src := "3.14;"
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		parseScript(src)
	}
}

func BenchmarkParseStringLiteral(b *testing.B) {
	src := `"hello world";`
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		parseScript(src)
	}
}

func BenchmarkParseStringWithEscapes(b *testing.B) {
	src := `"hello \"world\" \\path\\to\\file\n";`
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		parseScript(src)
	}
}

func BenchmarkParseBooleanLiteral(b *testing.B) {
	src := "true;"
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		parseScript(src)
	}
}

func BenchmarkParsePrefixExpression(b *testing.B) {
	src := "-42;"
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		parseScript(src)
	}
}

func BenchmarkParseInfixExpression(b *testing.B) {
	src := "1 + 2;"
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		parseScript(src)
	}
}

func BenchmarkParseNestedInfix(b *testing.B) {
	src := "1 + 2 * 3 - 4 / 5;"
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		parseScript(src)
	}
}

func BenchmarkParseGroupedExpression(b *testing.B) {
	src := "(1 + 2) * 3;"
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		parseScript(src)
	}
}

func BenchmarkParsePropertyAccess(b *testing.B) {
	src := "obj.name;"
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		parseScript(src)
	}
}

func BenchmarkParseChainedPropertyAccess(b *testing.B) {
	src := "a.b.c.d.e;"
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		parseScript(src)
	}
}

func BenchmarkParseIndexExpression(b *testing.B) {
	src := "arr[0];"
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		parseScript(src)
	}
}

func BenchmarkParseCallExpression(b *testing.B) {
	src := "add(1, 2);"
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		parseScript(src)
	}
}

func BenchmarkParseCallNoArgs(b *testing.B) {
	src := "doSomething();"
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		parseScript(src)
	}
}

func BenchmarkParseAssignment(b *testing.B) {
	src := "x = 42;"
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		parseScript(src)
	}
}

func BenchmarkParsePropertyAssignment(b *testing.B) {
	src := "obj.name = \"test\";"
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		parseScript(src)
	}
}

// --- Compound expression benchmarks ---

func BenchmarkParseIfExpression(b *testing.B) {
	src := "if (x > 0) { return true; }"
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		parseScript(src)
	}
}

func BenchmarkParseIfElseExpression(b *testing.B) {
	src := "if (x > 0) { return true; } else { return false; }"
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		parseScript(src)
	}
}

func BenchmarkParseForeachExpression(b *testing.B) {
	src := "foreach (items as item) { print(item); }"
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		parseScript(src)
	}
}

func BenchmarkParseFunctionLiteral(b *testing.B) {
	src := "fn add(a, b) { return a + b; }"
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		parseScript(src)
	}
}

func BenchmarkParseFunctionNoParams(b *testing.B) {
	src := "fn greet() { return \"hello\"; }"
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		parseScript(src)
	}
}

// --- Collection literal benchmarks ---

func BenchmarkParseArrayLiteral(b *testing.B) {
	src := "[1, 2, 3, 4, 5];"
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		parseScript(src)
	}
}

func BenchmarkParseEmptyArray(b *testing.B) {
	src := "[];"
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		parseScript(src)
	}
}

func BenchmarkParseHashLiteral(b *testing.B) {
	src := `{"name": "Alice", "age": 30};`
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		parseScript(src)
	}
}

func BenchmarkParseEmptyHash(b *testing.B) {
	src := "{};"
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		parseScript(src)
	}
}

// --- Template mode benchmarks ---

func BenchmarkParseTemplateSimple(b *testing.B) {
	src := "Hello {% name %}, welcome!"
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		parseTemplate(src)
	}
}

func BenchmarkParseTemplateTextOnly(b *testing.B) {
	src := strings.Repeat("Hello world, this is plain text. ", 100)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		parseTemplate(src)
	}
}

func BenchmarkParseTemplateManyBlocks(b *testing.B) {
	var sb strings.Builder
	for i := 0; i < 100; i++ {
		fmt.Fprintf(&sb, "text %d {%% var%d %%} ", i, i)
	}
	src := sb.String()
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		parseTemplate(src)
	}
}

func BenchmarkParseTemplateWithExpressions(b *testing.B) {
	src := "Result: {% 1 + 2 * 3 %} done"
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		parseTemplate(src)
	}
}

// --- Multi-statement benchmarks ---

func BenchmarkParseMultipleStatements(b *testing.B) {
	src := `let result = 0; fn add(a, b) { return a + b; }; if (result < 10) { result = add(result, 1); }`
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		parseScript(src)
	}
}

func BenchmarkParseComparisonOperators(b *testing.B) {
	src := `1 == 1; 2 != 3; 4 > 3; 2 < 5; 3 >= 3; 4 <= 5;`
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		parseScript(src)
	}
}

// --- Scaling benchmarks ---

func BenchmarkParseStatements10(b *testing.B) {
	src := strings.Repeat("let x = 42; ", 10)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		parseScript(src)
	}
}

func BenchmarkParseStatements100(b *testing.B) {
	src := strings.Repeat("let x = 42; ", 100)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		parseScript(src)
	}
}

func BenchmarkParseStatements1000(b *testing.B) {
	src := strings.Repeat("let x = 42; ", 1000)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		parseScript(src)
	}
}

func BenchmarkParseExpressions10(b *testing.B) {
	src := strings.Repeat("1 + 2 * 3; ", 10)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		parseScript(src)
	}
}

func BenchmarkParseExpressions100(b *testing.B) {
	src := strings.Repeat("1 + 2 * 3; ", 100)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		parseScript(src)
	}
}

// --- Constructor benchmarks ---

func BenchmarkNewParser(b *testing.B) {
	src := "let x = 1;"
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		l := lexer.NewScript(src)
		New(l)
	}
}

// --- Error path benchmarks ---

func BenchmarkParseWithError(b *testing.B) {
	src := "let x = ;"
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		parseScript(src)
	}
}

func BenchmarkParseWithLexerError(b *testing.B) {
	src := "let x = #;"
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		parseScript(src)
	}
}

// --- Debug output benchmarks ---

func BenchmarkProgramDebug(b *testing.B) {
	program, _ := parseScript("let x = 1 + 2; fn add(a, b) { return a + b; }")
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = program.Debug()
	}
}

func BenchmarkParseErrorFormat(b *testing.B) {
	err := NewParseError("test error", "let x = #;", lexer.Token{Line: 1, Column: 9})
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = err.Error()
	}
}

func BenchmarkParseLogicalOperators(b *testing.B) {
	src := "true && false || true;"
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		parseScript(src)
	}
}

func BenchmarkParseNullCoalescing(b *testing.B) {
	src := "x ?? y ?? z;"
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		parseScript(src)
	}
}

func BenchmarkParseWhileExpression(b *testing.B) {
	src := "while (x < 10) { x = x + 1; }"
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		parseScript(src)
	}
}

func BenchmarkParseWhileBreakContinue(b *testing.B) {
	src := "while (true) { if (x == 5) { break; } continue; }"
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		parseScript(src)
	}
}

func BenchmarkParseElseIf(b *testing.B) {
	src := "if (x == 1) { 1; } else if (x == 2) { 2; } else { 3; }"
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		parseScript(src)
	}
}

func BenchmarkParseForeachWithIndex(b *testing.B) {
	src := "foreach (arr as i, v) { v; }"
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		parseScript(src)
	}
}

func BenchmarkParseNullLiteral(b *testing.B) {
	src := "null;"
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		parseScript(src)
	}
}

func BenchmarkParseCommentSkipping(b *testing.B) {
	src := "// comment\nlet x = 5; /* block */ let y = 10;"
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		parseScript(src)
	}
}
