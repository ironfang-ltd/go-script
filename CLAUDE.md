# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build & Test Commands

```bash
go build ./...            # Build all packages
go test ./...             # Run all tests
go test -v ./...          # Run all tests with verbose output
go test ./lexer           # Run tests for a specific package
go test -run TestName ./evaluator  # Run a single test in a package
go run main.go            # Run the example template demo
```

No external dependencies — pure Go stdlib. Module requires Go 1.23+.

Any builds that produce an executable file must be output into the ./bin folder so they can get exluded from the git commit.

## Architecture

This is a **template engine with an embedded scripting language**. Templates mix plain text with script blocks delimited by `{% ... %}`. The interpreter follows a classic three-stage pipeline:

```
Template/Script String → Lexer → Parser → Evaluator → Output
```

### Lexer (`lexer/`)

Tokenizes input in two modes: `ModeTemplate` (plain text until `{%`) and `ModeScript` (code tokens until `%}`). The mode switching is what enables template syntax. Entry points: `lexer.NewTemplate(input)` for templates, `lexer.NewScript(input)` for pure script.

### Parser (`parser/`)

Builds an AST from the token stream using **Pratt parsing** (operator precedence climbing) for expressions and recursive descent for statements. Uses a 3-token lookahead window (`prev`, `current`, `next`). Parse errors are accumulated rather than failing on the first error — check the returned error from `p.Parse()`.

Key AST node categories:

- **Statements** (`statement.go`): `PrintStatement` (template text output), `LetStatement`, `ReturnStatement`, `ExpressionStatement`, `BlockStatement`
- **Expressions** (`expression.go`): `InfixExpression`, `PrefixExpression`, `IfExpression`, `ForeachExpression`, `FunctionLiteral`, `CallExpression`, `PropertyExpression`, `IndexExpression`, `AssignmentExpression`
- **Literals** (`literal.go`): `IntegerLiteral`, `StringLiteral`, `BooleanLiteral`, `ArrayLiteral`, `HashLiteral`

### Evaluator (`evaluator/`)

Tree-walking interpreter that executes AST nodes. Core concepts:

- **Object interface** (`object.go`): All runtime values implement `Object` with `Type()` and `Debug()` methods. Value types: `StringValue`, `IntegerValue`, `BooleanValue`, `DecimalValue`, `DateTimeValue`, `ArrayValue`, `HashValue`, `FunctionValue`, `NullValue`, etc.
- **Scope** (`scope.go`): Lexical scoping with parent chain. `SetLocal()` sets in current scope; `Assign()` walks the chain to find an existing variable; `Get()` reads from current or parent scopes.
- **ExecutionContext** (`evaluator.go`): Holds the parsed program, root scope, logger, metadata map, and output buffer. Created via `NewExecutionContext(program)` or `NewExecutionContextWithScope(program, scope)`.
- **Built-in functions**: `log()`, `print()`, `append()` are registered by default. Custom functions are added via `eval.RegisterFunction(name, fn)` where `fn` has signature `func(ctx *ExecutionContext, scope *Scope, args ...Object) (Object, error)`.
- **Two evaluation modes**: `Evaluate(ctx)` returns the final `Object`; `EvaluateString(ctx)` returns accumulated template output as a string.

### Typical Usage Flow

```go
l := lexer.NewTemplate(templateString)
p := parser.New(l)
program, err := p.Parse()
eval := evaluator.New()
eval.RegisterFunction("myFunc", myFuncImpl)
ctx := evaluator.NewExecutionContext(program)
ctx.RootScope.SetLocal("varName", &evaluator.StringValue{Value: "hello"})
output, err := eval.EvaluateString(ctx)
```

### HashValue

`HashValue` uses the `Hashable` interface for keys. Types implementing `Hashable` (strings, integers, booleans) can be used as hash keys. Use `evaluator.NewHashValue()` to create, then `Set(key, value)` / `Get(key)` to manipulate.
