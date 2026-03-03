# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build & Test Commands

```bash
go build ./...            # Build all packages
go test ./...             # Run all tests
go test -v ./...          # Run all tests with verbose output
go test ./lexer           # Run tests for a specific package
go test -run TestName ./evaluator  # Run a single test in a package
go test -bench=. -benchmem ./...  # Run benchmarks
go test -cover ./...      # Run tests with coverage
go run main.go            # Run the example template demo
```

No external dependencies — pure Go stdlib. Module requires Go 1.25+.

Any builds that produce an executable file must be output into the ./bin folder so they can get excluded from the git commit.

## Architecture

This is a **template engine with an embedded scripting language**. Templates mix plain text with script blocks delimited by `{% ... %}`. The interpreter follows a classic three-stage pipeline:

```
Template/Script String → Lexer → Parser → Evaluator → Output
```

### Lexer (`lexer/`)

Tokenizes input in two modes: `ModeTemplate` (plain text until `{%`) and `ModeScript` (code tokens until `%}`). The mode switching is what enables template syntax. Entry points: `lexer.NewTemplate(input)` for templates, `lexer.NewScript(input)` for pure script.

Keywords are matched via `tryIdentifierOrKeyword()` which reads the full identifier then checks a keyword map. Numbers only consume `.` if followed by a digit (prevents `1.method` being parsed as Float). String escape uses forward-looking: when `\` is seen, skip next char unconditionally.

### Parser (`parser/`)

Builds an AST from the token stream using **Pratt parsing** (operator precedence climbing) for expressions and recursive descent for statements. Uses a 3-token lookahead window (`prev`, `current`, `next`). Parse errors are accumulated (via `errors.Join()`) rather than failing on the first error — check the returned error from `p.Parse()`.

**Operator Precedence** (lowest to highest):
1. `??` (null coalescing)
2. `||`
3. `&&`
4. `==` `!=`
5. `<` `>` `<=` `>=`
6. `+` `-`
7. `/` `*` `%`
8. `.` (property access)
9. `()` `[]` (call and index)

Key AST node categories:

- **Statements** (`statement.go`): `PrintStatement` (template text output), `LetStatement`, `ReturnStatement`, `ExpressionStatement`, `BlockStatement`, `BreakStatement`, `ContinueStatement`, `AssignmentExpression`
- **Expressions** (`expression.go`): `InfixExpression`, `PrefixExpression`, `IfExpression`, `ForeachExpression`, `WhileExpression`, `FunctionLiteral`, `CallExpression`, `PropertyExpression`, `IndexExpression`
- **Literals** (`literal.go`): `IntegerLiteral`, `FloatLiteral`, `StringLiteral`, `BooleanLiteral`, `NullLiteral`, `ArrayLiteral`, `HashLiteral`

**Compound assignment** (`+=`, `-=`, `*=`, `/=`, `%=`) is desugared in the parser: `x += expr` becomes `AssignmentExpression{Left: x, Right: InfixExpression{Left: x, Op: "+", Right: expr}}`. No evaluator changes needed.

**Anonymous functions** are supported: `fn(x) { return x * 2; }`. The parser makes `FunctionLiteral.Identifier` optional — if `next` is `LeftParen`, it skips the name. Only named functions skip semicolons in `parseExpressionStatement`.

**Hash literals** use ordered `[]HashPair` slices (not maps) to preserve insertion order.

### Evaluator (`evaluator/`)

Tree-walking interpreter that executes AST nodes. Core concepts:

- **Object interface** (`object.go`): All runtime values implement `Object` with `Type()` and `Debug()` methods. Value types: `StringValue`, `IntegerValue`, `BooleanValue`, `DecimalValue`, `DateTimeValue`, `ArrayValue`, `HashValue`, `FileValue`, `FunctionValue`, `NullValue`, `BuiltInFunction`.
- **Scope** (`scope.go`): Lexical scoping with parent chain. `SetLocal()` sets in current scope; `Assign()` walks the chain to find an existing variable; `Get()` reads from current or parent scopes.
- **ExecutionContext** (`evaluator.go`): Holds the parsed program, root scope, logger, metadata map, output buffer, source string, and security limits. Created via `NewExecutionContext(program)` or `NewExecutionContextWithScope(program, scope)`.
- **Two evaluation modes**: `Evaluate(ctx)` returns the final `Object`; `EvaluateString(ctx)` returns accumulated template output as a string. `EvaluateString` sets `ctx.templateMode = true`, which causes `ExpressionStatement` nodes to auto-write their results to the output buffer at any nesting depth. `LetStatement`, `AssignmentExpression`, function bodies, and return values suppress `templateMode` to prevent value leakage. Template text is always written via `PrintStatement`.
- **Property access**: Dot notation (`obj.name`) only works on `HashValue` — it converts to an index lookup (`obj["name"]`). `PropertyExpression` evaluation returns `(parent, index, value, error)` to support both reads and writes.
- **Assignment**: Three forms — variable (`x = 5`), index (`arr[0] = val`), property (`obj.prop = val`). Variables must already exist via `let`; `Assign()` walks the scope chain to find them.
- **String auto-coercion**: `"str" + 42` → `"str42"`. If `+` operator has one `*StringValue` side, the other is coerced via `.Debug()`.

### Execution Security Limits

`ExecutionContext` enforces three limits (all configurable, override via `ctx.MaxSteps = N`):

- **`MaxSteps`** (default 100,000): Total node evaluations. Incremented at top of `evaluateNode`. Prevents infinite loops.
- **`MaxDepth`** (default 256): Function call nesting. Checked in `applyFunction`. Prevents stack overflow from recursion.
- **`MaxArraySize`** (default 10,000): Maximum array length. Checked in `append()` built-in and `evaluateArrayLiteral`.

Set any limit to 0 to disable it.

### Built-in Functions

**String**: `len`, `split`, `trim`, `toUpper`, `toLower`, `contains`, `startsWith`, `endsWith`, `indexOf`, `replace`, `substring`, `join`
**Type conversion**: `toString`, `parseInt`, `parseFloat`, `type`
**Math**: `floor`, `ceil`, `round`, `abs`
**Functional**: `map(arr, fn)`, `filter(arr, fn)` — these call `applyFunction` internally
**I/O**: `log` (writes to Logger), `print` (writes to output buffer), `append`
**Hash**: `keys`, `values` — both return arrays in insertion order

Custom functions: `eval.RegisterFunction(name, fn)` where `fn` has signature `func(ctx *ExecutionContext, scope *Scope, args ...Object) (Object, error)`.

### HashValue

`HashValue` maintains insertion order via an `order []HashKey` slice alongside the `Pairs` map. Types implementing `Hashable` (strings, integers, booleans, files) can be used as hash keys via FNV-1a hashing. Use `NewHashValue()` to create, `Set(key, value)` / `GetValue(key)` / `Delete(key)` to manipulate, `OrderedPairs()` for ordered iteration.

### Integer Overflow Protection

Integer arithmetic (`+`, `-`, `*`) checks for overflow and returns errors. Negating `math.MinInt` is also caught. Decimal modulo uses `math.Mod`.

### Error Handling

Three error types provide rich, formatted output with source location:
- **`TokenError`** (`lexer/error.go`): Lexer errors with line, column, and source context with caret pointer.
- **`ParseError`** (`parser/error.go`): Parser errors with similar formatting plus token info.
- **`RuntimeError`** (`evaluator/error.go`): Evaluator errors with same formatting. Only produced when `ctx.Source` is set; otherwise falls back to plain `fmt.Errorf`.

All three expand tabs to 4 spaces for proper alignment in error displays. The parser accumulates all errors and returns them joined at the end.

To enable `RuntimeError` formatting, set `ctx.Source` to the original source string before evaluation.

### Convenience Helpers (`helpers.go`)

- **`Vars`** (`map[string]any`): Enables `evaluator.Vars{"name": "Alice", "age": 30}` for variable injection.
- **`ToObject(v any) (Object, error)`**: Converts Go values to script Objects (nil→Null, string, bool, int variants, float variants, time.Time, `*time.Time`, `[]any`, `map[string]any`, Object passthrough).
- **`(*Evaluator).RunScript(source, vars...)`** / **`(*Evaluator).RunTemplate(source, vars...)`**: One-call lex→parse→evaluate with vars and custom functions.
- **`RunScript(source, vars...)`** / **`RunTemplate(source, vars...)`**: Package-level convenience (fresh evaluator, no custom functions).

### Typical Usage Flow

```go
// Simple (using helpers):
output, err := evaluator.RunTemplate(tmpl, evaluator.Vars{"name": "Alice"})

// With custom functions:
eval := evaluator.New()
eval.RegisterFunction("myFunc", myFuncImpl)
output, err := eval.RunTemplate(tmpl, evaluator.Vars{"name": "Alice"})

// Low-level (full control):
l := lexer.NewTemplate(templateString)
p := parser.New(l)
program, err := p.Parse()
eval := evaluator.New()
eval.RegisterFunction("myFunc", myFuncImpl)
ctx := evaluator.NewExecutionContext(program)
ctx.Source = templateString // enables RuntimeError source formatting
ctx.RootScope.SetLocal("varName", &evaluator.StringValue{Value: "hello"})
output, err := eval.EvaluateString(ctx)
```
