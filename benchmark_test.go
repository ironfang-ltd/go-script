package main

import (
	"fmt"
	"strings"
	"testing"

	"github.com/ironfang-ltd/go-script/lexer"
	"github.com/ironfang-ltd/go-script/parser"
)

// --- Script generators ---

// generateArithmeticScript produces N lines of let statements with arithmetic.
//
//	let v0 = 1 + 2 * 3 - 4 / 2;
//	let v1 = v0 + 5 * 3 - 1 / 1;
//	...
func generateArithmeticScript(n int) string {
	var sb strings.Builder
	sb.WriteString("let v0 = 1 + 2 * 3 - 4 / 2;\n")
	for i := 1; i < n; i++ {
		fmt.Fprintf(&sb, "let v%d = v%d + %d * 3 - %d / 1;\n", i, i-1, i%100, i%7+1)
	}
	fmt.Fprintf(&sb, "return v%d;\n", n-1)
	return sb.String()
}

// generateFunctionScript produces N function definitions with bodies that do
// arithmetic and call the previous function.
func generateFunctionScript(n int) string {
	var sb strings.Builder
	sb.WriteString("fn f0(x) { return x + 1; }\n")
	for i := 1; i < n; i++ {
		fmt.Fprintf(&sb, "fn f%d(x) { let y = x + %d; return f%d(y); }\n", i, i%50, i-1)
	}
	fmt.Fprintf(&sb, "return f%d(0);\n", n-1)
	return sb.String()
}

// generateControlFlowScript produces N if/else blocks with comparisons.
func generateControlFlowScript(n int) string {
	var sb strings.Builder
	sb.WriteString("let result = 0;\n")
	for i := 0; i < n; i++ {
		fmt.Fprintf(&sb, "if (result < %d) { result = result + 1; } else { result = result - 1; }\n", i+1)
	}
	sb.WriteString("return result;\n")
	return sb.String()
}

// generateStringHeavyScript produces N let statements assigning and
// concatenating strings.
func generateStringHeavyScript(n int) string {
	var sb strings.Builder
	sb.WriteString("let s = \"hello\";\n")
	for i := 0; i < n; i++ {
		fmt.Fprintf(&sb, "let s%d = s + \" world_%d\";\n", i, i)
	}
	fmt.Fprintf(&sb, "return s%d;\n", n-1)
	return sb.String()
}

// generateArrayHashScript produces code that builds arrays and hashes.
func generateArrayHashScript(n int) string {
	var sb strings.Builder
	sb.WriteString("let arr = [];\n")
	for i := 0; i < n; i++ {
		fmt.Fprintf(&sb, "let h%d = {\"key\": %d, \"val\": \"%d\"};\n", i, i, i)
		fmt.Fprintf(&sb, "append(arr, h%d);\n", i)
	}
	sb.WriteString("return arr;\n")
	return sb.String()
}

// generateMixedScript produces a realistic script mixing all features.
func generateMixedScript(n int) string {
	var sb strings.Builder

	// Functions
	sb.WriteString("fn add(a, b) { return a + b; }\n")
	sb.WriteString("fn mul(a, b) { return a * b; }\n")

	// Variables with arithmetic
	sb.WriteString("let total = 0;\n")
	for i := 0; i < n; i++ {
		fmt.Fprintf(&sb, "let x%d = add(%d, %d);\n", i, i, i*2)
		fmt.Fprintf(&sb, "let y%d = mul(x%d, %d);\n", i, i, i%10+1)
		fmt.Fprintf(&sb, "if (y%d > %d) { total = total + y%d; } else { total = total + 1; }\n", i, i*5, i)
	}
	sb.WriteString("return total;\n")
	return sb.String()
}

// generateForeachScript produces N foreach loops over arrays.
func generateForeachScript(n int) string {
	var sb strings.Builder
	sb.WriteString("let total = 0;\n")
	for i := 0; i < n; i++ {
		fmt.Fprintf(&sb, "let arr%d = [%d, %d, %d];\n", i, i, i+1, i+2)
		fmt.Fprintf(&sb, "foreach (arr%d as item) { total = total + item; }\n", i)
	}
	sb.WriteString("return total;\n")
	return sb.String()
}

// generateTemplate produces a template with N script blocks interleaved with text.
func generateTemplate(n int) string {
	var sb strings.Builder
	for i := 0; i < n; i++ {
		fmt.Fprintf(&sb, "Hello user %d, ", i)
		fmt.Fprintf(&sb, "{%% let x%d = %d + %d; %%}", i, i, i*3)
		fmt.Fprintf(&sb, "your score is {%% x%d %%}.\n", i)
	}
	return sb.String()
}

// --- Helpers ---

func tokenizeAll(l *lexer.Lexer) int {
	count := 0
	for {
		tok, err := l.Read()
		if err != nil {
			break
		}
		count++
		if tok.Type == lexer.EndOfFile {
			break
		}
	}
	return count
}

// --- Lexer benchmarks ---

func BenchmarkLexerArithmetic(b *testing.B) {
	src := generateArithmeticScript(500)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		l := lexer.NewScript(src)
		tokenizeAll(l)
	}
}

func BenchmarkLexerFunctions(b *testing.B) {
	src := generateFunctionScript(200)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		l := lexer.NewScript(src)
		tokenizeAll(l)
	}
}

func BenchmarkLexerControlFlow(b *testing.B) {
	src := generateControlFlowScript(500)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		l := lexer.NewScript(src)
		tokenizeAll(l)
	}
}

func BenchmarkLexerStrings(b *testing.B) {
	src := generateStringHeavyScript(500)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		l := lexer.NewScript(src)
		tokenizeAll(l)
	}
}

func BenchmarkLexerArrayHash(b *testing.B) {
	src := generateArrayHashScript(250)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		l := lexer.NewScript(src)
		tokenizeAll(l)
	}
}

func BenchmarkLexerMixed(b *testing.B) {
	src := generateMixedScript(200)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		l := lexer.NewScript(src)
		tokenizeAll(l)
	}
}

func BenchmarkLexerForeach(b *testing.B) {
	src := generateForeachScript(250)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		l := lexer.NewScript(src)
		tokenizeAll(l)
	}
}

func BenchmarkLexerTemplate(b *testing.B) {
	src := generateTemplate(250)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		l := lexer.NewTemplate(src)
		tokenizeAll(l)
	}
}

// --- Parser benchmarks ---

func BenchmarkParserArithmetic(b *testing.B) {
	src := generateArithmeticScript(500)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		l := lexer.NewScript(src)
		p := parser.New(l)
		_, _ = p.Parse()
	}
}

func BenchmarkParserFunctions(b *testing.B) {
	src := generateFunctionScript(200)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		l := lexer.NewScript(src)
		p := parser.New(l)
		_, _ = p.Parse()
	}
}

func BenchmarkParserControlFlow(b *testing.B) {
	src := generateControlFlowScript(500)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		l := lexer.NewScript(src)
		p := parser.New(l)
		_, _ = p.Parse()
	}
}

func BenchmarkParserStrings(b *testing.B) {
	src := generateStringHeavyScript(500)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		l := lexer.NewScript(src)
		p := parser.New(l)
		_, _ = p.Parse()
	}
}

func BenchmarkParserArrayHash(b *testing.B) {
	src := generateArrayHashScript(250)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		l := lexer.NewScript(src)
		p := parser.New(l)
		_, _ = p.Parse()
	}
}

func BenchmarkParserMixed(b *testing.B) {
	src := generateMixedScript(200)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		l := lexer.NewScript(src)
		p := parser.New(l)
		_, _ = p.Parse()
	}
}

func BenchmarkParserForeach(b *testing.B) {
	src := generateForeachScript(250)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		l := lexer.NewScript(src)
		p := parser.New(l)
		_, _ = p.Parse()
	}
}

func BenchmarkParserTemplate(b *testing.B) {
	src := generateTemplate(250)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		l := lexer.NewTemplate(src)
		p := parser.New(l)
		_, _ = p.Parse()
	}
}

// --- Scaling benchmarks (same script type, increasing sizes) ---

func BenchmarkLexerArithmeticScale100(b *testing.B) {
	src := generateArithmeticScript(100)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		l := lexer.NewScript(src)
		tokenizeAll(l)
	}
}

func BenchmarkLexerArithmeticScale1000(b *testing.B) {
	src := generateArithmeticScript(1000)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		l := lexer.NewScript(src)
		tokenizeAll(l)
	}
}

func BenchmarkLexerArithmeticScale5000(b *testing.B) {
	src := generateArithmeticScript(5000)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		l := lexer.NewScript(src)
		tokenizeAll(l)
	}
}

func BenchmarkParserArithmeticScale100(b *testing.B) {
	src := generateArithmeticScript(100)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		l := lexer.NewScript(src)
		p := parser.New(l)
		_, _ = p.Parse()
	}
}

func BenchmarkParserArithmeticScale1000(b *testing.B) {
	src := generateArithmeticScript(1000)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		l := lexer.NewScript(src)
		p := parser.New(l)
		_, _ = p.Parse()
	}
}

func BenchmarkParserArithmeticScale5000(b *testing.B) {
	src := generateArithmeticScript(5000)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		l := lexer.NewScript(src)
		p := parser.New(l)
		_, _ = p.Parse()
	}
}
