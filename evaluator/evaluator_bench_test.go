package evaluator

import (
	"fmt"
	"strings"
	"testing"

	"github.com/ironfang-ltd/go-script/lexer"
	"github.com/ironfang-ltd/go-script/parser"
)

// --- Helpers ---

func parseScriptProgram(b *testing.B, src string) *parser.Program {
	b.Helper()
	l := lexer.NewScript(src)
	p := parser.New(l)
	program, err := p.Parse()
	if err != nil {
		b.Fatalf("Parse error: %v", err)
	}
	return program
}

func parseTemplateProgram(b *testing.B, src string) *parser.Program {
	b.Helper()
	l := lexer.NewTemplate(src)
	p := parser.New(l)
	program, err := p.Parse()
	if err != nil {
		b.Fatalf("Parse error: %v", err)
	}
	return program
}

// --- Value construction benchmarks ---

func BenchmarkNewIntegerValue(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		NewIntegerValue(42)
	}
}

func BenchmarkNewStringValue(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		NewStringValue("hello world")
	}
}

func BenchmarkNewBooleanValue(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		NewBooleanValue(true)
	}
}

func BenchmarkNewDecimalValue(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		NewDecimalValue(3.14)
	}
}

func BenchmarkNewArrayValue(b *testing.B) {
	elems := []Object{NewIntegerValue(1), NewIntegerValue(2), NewIntegerValue(3)}
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		NewArrayValue(elems)
	}
}

func BenchmarkNewHashValue(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		NewHashValue()
	}
}

func BenchmarkStringHashKey(b *testing.B) {
	s := NewStringValue("hello")
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		s.HashKey()
	}
}

func BenchmarkFileHashKey(b *testing.B) {
	f := NewFileValue("id-1", "/path/to/file", "file.png", "image/png", 1024)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		f.HashKey()
	}
}

// --- Expression evaluation benchmarks ---

func BenchmarkEvalIntegerArithmetic(b *testing.B) {
	prog := parseScriptProgram(b, "return 2 + 3 * 4 - 1;")
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		e := New()
		ctx := NewExecutionContext(prog)
		e.Evaluate(ctx)
	}
}

func BenchmarkEvalIntegerComparison(b *testing.B) {
	prog := parseScriptProgram(b, "return 5 > 3;")
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		e := New()
		ctx := NewExecutionContext(prog)
		e.Evaluate(ctx)
	}
}

func BenchmarkEvalDecimalArithmetic(b *testing.B) {
	prog := parseScriptProgram(b, "return 1.5 + 2.5 * 3.0;")
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		e := New()
		ctx := NewExecutionContext(prog)
		e.Evaluate(ctx)
	}
}

func BenchmarkEvalStringConcat(b *testing.B) {
	prog := parseScriptProgram(b, `return "hello" + " " + "world";`)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		e := New()
		ctx := NewExecutionContext(prog)
		e.Evaluate(ctx)
	}
}

func BenchmarkEvalBooleanOps(b *testing.B) {
	prog := parseScriptProgram(b, "return true == false;")
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		e := New()
		ctx := NewExecutionContext(prog)
		e.Evaluate(ctx)
	}
}

func BenchmarkEvalPrefixBang(b *testing.B) {
	prog := parseScriptProgram(b, "return !true;")
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		e := New()
		ctx := NewExecutionContext(prog)
		e.Evaluate(ctx)
	}
}

func BenchmarkEvalPrefixMinus(b *testing.B) {
	prog := parseScriptProgram(b, "return -42;")
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		e := New()
		ctx := NewExecutionContext(prog)
		e.Evaluate(ctx)
	}
}

func BenchmarkEvalGroupedExpression(b *testing.B) {
	prog := parseScriptProgram(b, "return (2 + 3) * (4 - 1);")
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		e := New()
		ctx := NewExecutionContext(prog)
		e.Evaluate(ctx)
	}
}

// --- Statement evaluation benchmarks ---

func BenchmarkEvalLetStatement(b *testing.B) {
	prog := parseScriptProgram(b, "let x = 42;")
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		e := New()
		ctx := NewExecutionContext(prog)
		e.Evaluate(ctx)
	}
}

func BenchmarkEvalReturnStatement(b *testing.B) {
	prog := parseScriptProgram(b, "return 42;")
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		e := New()
		ctx := NewExecutionContext(prog)
		e.Evaluate(ctx)
	}
}

func BenchmarkEvalIfElse(b *testing.B) {
	prog := parseScriptProgram(b, `if (true) { return 1; } else { return 2; }`)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		e := New()
		ctx := NewExecutionContext(prog)
		e.Evaluate(ctx)
	}
}

func BenchmarkEvalAssignment(b *testing.B) {
	prog := parseScriptProgram(b, "let x = 0; x = 42;")
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		e := New()
		ctx := NewExecutionContext(prog)
		e.Evaluate(ctx)
	}
}

// --- Function benchmarks ---

func BenchmarkEvalFunctionDefinition(b *testing.B) {
	prog := parseScriptProgram(b, "fn add(a, b) { return a + b; }")
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		e := New()
		ctx := NewExecutionContext(prog)
		e.Evaluate(ctx)
	}
}

func BenchmarkEvalFunctionCall(b *testing.B) {
	prog := parseScriptProgram(b, "fn add(a, b) { return a + b; } return add(3, 4);")
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		e := New()
		ctx := NewExecutionContext(prog)
		e.Evaluate(ctx)
	}
}

func BenchmarkEvalRecursiveFunction(b *testing.B) {
	prog := parseScriptProgram(b, `
fn factorial(n) {
	if (n <= 1) { return 1; }
	return n * factorial(n - 1);
}
return factorial(10);`)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		e := New()
		ctx := NewExecutionContext(prog)
		e.Evaluate(ctx)
	}
}

func BenchmarkEvalClosure(b *testing.B) {
	prog := parseScriptProgram(b, `
let x = 10;
fn addX(y) { return x + y; }
return addX(5);`)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		e := New()
		ctx := NewExecutionContext(prog)
		e.Evaluate(ctx)
	}
}

func BenchmarkEvalBuiltinAppend(b *testing.B) {
	prog := parseScriptProgram(b, `let a = []; append(a, 1); append(a, 2); append(a, 3);`)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		e := New()
		ctx := NewExecutionContext(prog)
		e.Evaluate(ctx)
	}
}

// --- Collection benchmarks ---

func BenchmarkEvalArrayLiteral(b *testing.B) {
	prog := parseScriptProgram(b, "let a = [1, 2, 3, 4, 5];")
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		e := New()
		ctx := NewExecutionContext(prog)
		e.Evaluate(ctx)
	}
}

func BenchmarkEvalArrayIndex(b *testing.B) {
	prog := parseScriptProgram(b, "let a = [10, 20, 30]; return a[1];")
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		e := New()
		ctx := NewExecutionContext(prog)
		e.Evaluate(ctx)
	}
}

func BenchmarkEvalHashLiteral(b *testing.B) {
	prog := parseScriptProgram(b, `let h = {"a": 1, "b": 2, "c": 3};`)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		e := New()
		ctx := NewExecutionContext(prog)
		e.Evaluate(ctx)
	}
}

func BenchmarkEvalHashIndex(b *testing.B) {
	prog := parseScriptProgram(b, `let h = {"x": 42}; return h["x"];`)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		e := New()
		ctx := NewExecutionContext(prog)
		e.Evaluate(ctx)
	}
}

func BenchmarkEvalPropertyAccess(b *testing.B) {
	prog := parseScriptProgram(b, `let h = {"name": "test"}; return h.name;`)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		e := New()
		ctx := NewExecutionContext(prog)
		e.Evaluate(ctx)
	}
}

func BenchmarkEvalNestedProperty(b *testing.B) {
	prog := parseScriptProgram(b, `let h = {"a": {"b": {"c": 42}}}; return h.a.b.c;`)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		e := New()
		ctx := NewExecutionContext(prog)
		e.Evaluate(ctx)
	}
}

func BenchmarkEvalForeachArray(b *testing.B) {
	prog := parseScriptProgram(b, `
let arr = [1, 2, 3, 4, 5];
let sum = 0;
foreach (arr as item) { sum = sum + item; }
return sum;`)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		e := New()
		ctx := NewExecutionContext(prog)
		e.Evaluate(ctx)
	}
}

func BenchmarkEvalForeachHash(b *testing.B) {
	prog := parseScriptProgram(b, `
let h = {"a": 1, "b": 2, "c": 3};
let sum = 0;
foreach (h as val) { sum = sum + val; }
return sum;`)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		e := New()
		ctx := NewExecutionContext(prog)
		e.Evaluate(ctx)
	}
}

// --- Template mode benchmarks ---

func BenchmarkEvalTemplateSimple(b *testing.B) {
	prog := parseTemplateProgram(b, `Hello {% "world" %}!`)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		e := New()
		ctx := NewExecutionContext(prog)
		e.EvaluateString(ctx)
	}
}

func BenchmarkEvalTemplateVariable(b *testing.B) {
	prog := parseTemplateProgram(b, `Hello {% name %}!`)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		e := New()
		ctx := NewExecutionContext(prog)
		ctx.RootScope.SetLocal("name", &StringValue{Value: "Alice"})
		e.EvaluateString(ctx)
	}
}

func BenchmarkEvalTemplateExpression(b *testing.B) {
	prog := parseTemplateProgram(b, `Result: {% 1 + 2 * 3 %}`)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		e := New()
		ctx := NewExecutionContext(prog)
		e.EvaluateString(ctx)
	}
}

func BenchmarkEvalTemplateManyBlocks(b *testing.B) {
	var sb strings.Builder
	for j := 0; j < 50; j++ {
		fmt.Fprintf(&sb, "text %d {%% %d + %d %%} ", j, j, j+1)
	}
	src := sb.String()
	prog := parseTemplateProgram(b, src)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		e := New()
		ctx := NewExecutionContext(prog)
		e.EvaluateString(ctx)
	}
}

func BenchmarkEvalTemplateTextOnly(b *testing.B) {
	src := strings.Repeat("Hello world, this is plain text. ", 100)
	prog := parseTemplateProgram(b, src)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		e := New()
		ctx := NewExecutionContext(prog)
		e.EvaluateString(ctx)
	}
}

// --- Scaling benchmarks ---

func BenchmarkEvalStatements10(b *testing.B) {
	var sb strings.Builder
	for j := 0; j < 10; j++ {
		fmt.Fprintf(&sb, "let v%d = %d; ", j, j)
	}
	prog := parseScriptProgram(b, sb.String())
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		e := New()
		ctx := NewExecutionContext(prog)
		e.Evaluate(ctx)
	}
}

func BenchmarkEvalStatements100(b *testing.B) {
	var sb strings.Builder
	for j := 0; j < 100; j++ {
		fmt.Fprintf(&sb, "let v%d = %d; ", j, j)
	}
	prog := parseScriptProgram(b, sb.String())
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		e := New()
		ctx := NewExecutionContext(prog)
		e.Evaluate(ctx)
	}
}

func BenchmarkEvalStatements1000(b *testing.B) {
	var sb strings.Builder
	for j := 0; j < 1000; j++ {
		fmt.Fprintf(&sb, "let v%d = %d; ", j, j)
	}
	prog := parseScriptProgram(b, sb.String())
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		e := New()
		ctx := NewExecutionContext(prog)
		e.Evaluate(ctx)
	}
}

func BenchmarkEvalExpressions10(b *testing.B) {
	parts := make([]string, 10)
	for j := 0; j < 10; j++ {
		parts[j] = fmt.Sprintf("%d", j+1)
	}
	src := "return " + strings.Join(parts, " + ") + ";"
	prog := parseScriptProgram(b, src)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		e := New()
		ctx := NewExecutionContext(prog)
		e.Evaluate(ctx)
	}
}

func BenchmarkEvalExpressions100(b *testing.B) {
	parts := make([]string, 100)
	for j := 0; j < 100; j++ {
		parts[j] = fmt.Sprintf("%d", j+1)
	}
	src := "return " + strings.Join(parts, " + ") + ";"
	prog := parseScriptProgram(b, src)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		e := New()
		ctx := NewExecutionContext(prog)
		e.Evaluate(ctx)
	}
}

// --- Scope benchmarks ---

func BenchmarkNewScope(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		NewScope()
	}
}

func BenchmarkNewChildScope(b *testing.B) {
	parent := NewScope()
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		NewChildScope(parent)
	}
}

func BenchmarkScopeSetLocalAndGet(b *testing.B) {
	s := NewScope()
	val := NewIntegerValue(42)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		s.SetLocal("x", val)
		s.Get("x")
	}
}

func BenchmarkScopeDeepChainLookup(b *testing.B) {
	root := NewScope()
	root.SetLocal("x", NewIntegerValue(42))
	current := root
	for j := 0; j < 10; j++ {
		current = NewChildScope(current)
	}
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		current.Get("x")
	}
}

// --- Error path benchmarks ---

func BenchmarkEvalUndefinedIdentifier(b *testing.B) {
	prog := parseScriptProgram(b, "return x;")
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		e := New()
		ctx := NewExecutionContext(prog)
		e.Evaluate(ctx)
	}
}

func BenchmarkEvalDivisionByZero(b *testing.B) {
	prog := parseScriptProgram(b, "return 10 / 0;")
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		e := New()
		ctx := NewExecutionContext(prog)
		e.Evaluate(ctx)
	}
}

func BenchmarkEvalTypeMismatch(b *testing.B) {
	prog := parseScriptProgram(b, `return 1 + "hello";`)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		e := New()
		ctx := NewExecutionContext(prog)
		e.Evaluate(ctx)
	}
}

// --- HashValue operation benchmarks ---

func BenchmarkHashValueSet(b *testing.B) {
	h := NewHashValue()
	key := NewStringValue("key")
	val := NewIntegerValue(42)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		h.Set(key, val)
	}
}

func BenchmarkHashValueGet(b *testing.B) {
	h := NewHashValue()
	key := NewStringValue("key")
	h.Set(key, NewIntegerValue(42))
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		h.GetValue(key)
	}
}

func BenchmarkHashValueHasKey(b *testing.B) {
	h := NewHashValue()
	key := NewStringValue("key")
	h.Set(key, NewIntegerValue(42))
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		h.HasKey(key)
	}
}

// --- Evaluator constructor benchmark ---

func BenchmarkNewEvaluator(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		New()
	}
}

func BenchmarkNewExecutionContext(b *testing.B) {
	prog := parseScriptProgram(b, "return 1;")
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		NewExecutionContext(prog)
	}
}

func BenchmarkEvalLogicalAnd(b *testing.B) {
	prog := parseScriptProgram(b, "true && false;")
	e := New()
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ctx := NewExecutionContext(prog)
		e.Evaluate(ctx)
	}
}

func BenchmarkEvalLogicalOr(b *testing.B) {
	prog := parseScriptProgram(b, "false || true;")
	e := New()
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ctx := NewExecutionContext(prog)
		e.Evaluate(ctx)
	}
}

func BenchmarkEvalNullCoalescing(b *testing.B) {
	prog := parseScriptProgram(b, "null ?? 5;")
	e := New()
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ctx := NewExecutionContext(prog)
		e.Evaluate(ctx)
	}
}

func BenchmarkEvalWhileLoop(b *testing.B) {
	prog := parseScriptProgram(b, "let x = 0; while (x < 100) { x = x + 1; }")
	e := New()
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ctx := NewExecutionContext(prog)
		e.Evaluate(ctx)
	}
}

func BenchmarkEvalWhileBreak(b *testing.B) {
	prog := parseScriptProgram(b, "let x = 0; while (true) { x = x + 1; if (x == 50) { break; } }")
	e := New()
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ctx := NewExecutionContext(prog)
		e.Evaluate(ctx)
	}
}

func BenchmarkEvalElseIf(b *testing.B) {
	prog := parseScriptProgram(b, `let x = 3; if (x == 1) { 1; } else if (x == 2) { 2; } else if (x == 3) { 3; } else { 4; }`)
	e := New()
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ctx := NewExecutionContext(prog)
		e.Evaluate(ctx)
	}
}

func BenchmarkEvalLen(b *testing.B) {
	prog := parseScriptProgram(b, `len("hello world");`)
	e := New()
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ctx := NewExecutionContext(prog)
		e.Evaluate(ctx)
	}
}

func BenchmarkEvalStringSplit(b *testing.B) {
	prog := parseScriptProgram(b, `split("a,b,c,d,e", ",");`)
	e := New()
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ctx := NewExecutionContext(prog)
		e.Evaluate(ctx)
	}
}

func BenchmarkEvalStringBuiltins(b *testing.B) {
	prog := parseScriptProgram(b, `toUpper("hello"); toLower("HELLO"); trim("  x  "); contains("hello", "ell"); replace("hello", "l", "r");`)
	e := New()
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ctx := NewExecutionContext(prog)
		e.Evaluate(ctx)
	}
}

func BenchmarkEvalTypeCoercion(b *testing.B) {
	prog := parseScriptProgram(b, "1 + 2.5;")
	e := New()
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ctx := NewExecutionContext(prog)
		e.Evaluate(ctx)
	}
}

func BenchmarkEvalForeachWithIndex(b *testing.B) {
	prog := parseScriptProgram(b, `let sum = 0; foreach ([1,2,3,4,5] as i, v) { sum = sum + i; }`)
	e := New()
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ctx := NewExecutionContext(prog)
		e.Evaluate(ctx)
	}
}

func BenchmarkEvalIntegerHashKey(b *testing.B) {
	v := &IntegerValue{Value: 42}
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		v.HashKey()
	}
}

func BenchmarkEvalBooleanHashKey(b *testing.B) {
	v := &BooleanValue{Value: true}
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		v.HashKey()
	}
}

func BenchmarkEvalKeysValues(b *testing.B) {
	prog := parseScriptProgram(b, `let h = {"a": 1, "b": 2, "c": 3}; keys(h); values(h);`)
	e := New()
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ctx := NewExecutionContext(prog)
		e.Evaluate(ctx)
	}
}

func BenchmarkEvalTypeBuiltin(b *testing.B) {
	prog := parseScriptProgram(b, `type(42);`)
	e := New()
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ctx := NewExecutionContext(prog)
		e.Evaluate(ctx)
	}
}
