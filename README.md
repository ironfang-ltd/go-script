# go-script

A template engine with an embedded scripting language for Go. Templates mix plain text with script blocks delimited by `{% ... %}`, making it ideal for server-side rendering, email templates, and code generation. The interpreter also runs standalone scripts without template syntax.

Pure Go, zero dependencies.

## Installation

```bash
go get github.com/ironfang-ltd/go-script
```

Requires Go 1.25+.

## Quick Start

### Running a Script

```go
result, err := evaluator.RunScript(`
    let greeting = "Hello, " + name + "!";
    return greeting;
`, evaluator.Vars{"name": "world"})

fmt.Println(result.(*evaluator.ReturnValue).Value.Debug()) // Hello, world!
```

### Rendering a Template

```go
output, err := evaluator.RunTemplate(
    `<h1>{% title %}</h1>
<ul>
    {% foreach (items as item) { %}
    <li>{% item.name %}</li>
    {% } %}
</ul>`,
    evaluator.Vars{
        "title": "My List",
        "items": []any{
            map[string]any{"name": "Apple"},
            map[string]any{"name": "Banana"},
        },
    },
)
fmt.Println(output)
```

Output:

```html
<h1>My List</h1>
<ul>
  <li>Apple</li>
  <li>Banana</li>
</ul>
```

### With Custom Functions

When you need to register custom functions, create an evaluator and use its methods:

```go
eval := evaluator.New()
eval.RegisterFunction("shout", func(
    ctx *evaluator.ExecutionContext,
    scope *evaluator.Scope,
    args ...evaluator.Object,
) (evaluator.Object, error) {
    return &evaluator.StringValue{Value: args[0].Debug() + "!"}, nil
})

output, err := eval.RunTemplate(`Say {% shout(name) %}`, evaluator.Vars{"name": "hello"})
// output: "Say hello!"
```

---

## Scripting Language Reference

### Data Types

| Type     | Examples                         | Description                                            |
| -------- | -------------------------------- | ------------------------------------------------------ |
| Integer  | `42`, `0`, `-7`                  | 64-bit signed integer                                  |
| Decimal  | `3.14`, `0.5`                    | 64-bit floating point                                  |
| String   | `"hello"`, `"line\nbreak"`       | Double-quoted, supports `\\`, `\"`, `\n`, `\t` escapes |
| Boolean  | `true`, `false`                  |                                                        |
| Null     | `null`                           | Absence of a value                                     |
| Array    | `[1, 2, 3]`                      | Ordered, mixed-type collection                         |
| Hash     | `{"key": "value"}`               | Ordered key-value map (insertion order preserved)      |
| Function | `fn add(a, b) { return a + b; }` | First-class, supports closures                         |

### Variables

Declare with `let`, reassign with `=`:

```
let name = "Alice";
let age = 30;
age = 31;
```

Variables must be declared with `let` before they can be reassigned. Assignment walks the scope chain — it finds the variable in the nearest enclosing scope that defined it.

### Operators

#### Arithmetic

```
let a = 10 + 3;    // 13
let b = 10 - 3;    // 7
let c = 10 * 3;    // 30
let d = 10 / 3;    // 3 (integer division)
let e = 10 % 3;    // 1
let f = 10.0 / 3;  // 3.3333... (decimal division)
let g = -a;         // -13
```

Integer arithmetic detects overflow and returns an error instead of silently wrapping.

Mixed integer/decimal operations automatically promote the integer to decimal:

```
let x = 5 + 2.5;   // 7.5 (decimal)
```

#### Comparison

```
1 == 1    // true
1 != 2    // true
3 < 5     // true
5 > 3     // true
3 <= 3    // true
3 >= 4    // false
```

#### Logical

```
true && false   // false (short-circuit)
false || true   // true  (short-circuit)
!true           // false
```

#### Null Coalescing

```
let val = null ?? "default";   // "default"
let val2 = "hello" ?? "nope";  // "hello"
```

#### String Concatenation and Auto-Coercion

```
"hello" + " " + "world"    // "hello world"

// When one side of + is a string, the other is auto-coerced:
"count: " + 42              // "count: 42"
3.14 + " is pi"             // "3.14 is pi"
true + " story"             // "true story"
```

#### Compound Assignment

```
let x = 10;
x += 5;    // x is now 15
x -= 3;    // x is now 12
x *= 2;    // x is now 24
x /= 4;    // x is now 6
x %= 4;    // x is now 2
```

Works with property and index access too:

```
let obj = {"count": 0};
obj.count += 1;

let arr = [10, 20, 30];
arr[0] += 5;
```

### Control Flow

#### If / Else

```
if (x > 10) {
    return "big";
} else if (x > 5) {
    return "medium";
} else {
    return "small";
}
```

`if` is an expression — it returns the value of the last statement in the taken branch:

```
let label = if (score >= 90) { "A"; } else { "B"; };
```

#### While Loop

```
let i = 0;
while (i < 10) {
    log(i);
    i += 1;
}
```

#### Foreach Loop

Iterate over arrays:

```
let fruits = ["apple", "banana", "cherry"];

foreach (fruits as fruit) {
    log(fruit);
}
```

With index:

```
foreach (fruits as fruit, i) {
    log(i + ": " + fruit);
}
// 0: apple
// 1: banana
// 2: cherry
```

Iterate over hashes (insertion order is preserved):

```
let person = {"name": "Alice", "age": 30, "city": "London"};

foreach (person as value, key) {
    log(key + " = " + toString(value));
}
// name = Alice
// age = 30
// city = London
```

#### Break and Continue

```
let i = 0;
while (true) {
    if (i >= 5) { break; }
    i += 1;
}

foreach ([1, 2, 3, 4, 5] as n) {
    if (n % 2 == 0) { continue; }
    log(n);  // 1, 3, 5
}
```

### Functions

#### Named Functions

```
fn greet(name) {
    return "Hello, " + name + "!";
}

let msg = greet("Alice");  // "Hello, Alice!"
```

Named functions are statements — no semicolon needed after the closing brace.

#### Anonymous Functions

```
let double = fn(x) { return x * 2; };

let result = double(5);  // 10
```

#### Closures

Functions capture their enclosing scope:

```
fn makeCounter() {
    let count = 0;
    return fn() {
        count = count + 1;
        return count;
    };
}

let counter = makeCounter();
counter();  // 1
counter();  // 2
counter();  // 3
```

#### Recursion

```
fn fibonacci(n) {
    if (n <= 1) { return n; }
    return fibonacci(n - 1) + fibonacci(n - 2);
}

fibonacci(10);  // 55
```

#### Higher-Order Functions

Pass functions as arguments:

```
let nums = [1, 2, 3, 4, 5];

let evens = filter(nums, fn(x) { return x % 2 == 0; });
// [2, 4]

let doubled = map(nums, fn(x) { return x * 2; });
// [2, 4, 6, 8, 10]

let csv = join(doubled, ", ");
// "2, 4, 6, 8, 10"
```

### Arrays

```
let arr = [1, "two", true, null];

// Access by index (0-based)
arr[0];     // 1
arr[1];     // "two"

// Modify
arr[0] = 99;

// Append
arr = append(arr, "new");

// Length
len(arr);   // 5

// Nested arrays
let matrix = [[1, 2], [3, 4]];
matrix[1][0];  // 3
```

### Hashes

Hashes are ordered maps — keys maintain insertion order during iteration.

```
let user = {
    "name": "Alice",
    "age": 30,
    "active": true
};

// Access via index
user["name"];    // "Alice"

// Access via dot notation
user.name;       // "Alice"

// Set/modify
user["email"] = "alice@example.com";
user.email = "alice@example.com";

// Get keys and values (in insertion order)
keys(user);      // ["name", "age", "active", "email"]
values(user);    // ["Alice", 30, true, "alice@example.com"]

// Nested access
let data = {
    "user": {
        "address": {
            "city": "London"
        }
    }
};
data.user.address.city;  // "London"
```

Valid key types: strings, integers, booleans.

### Comments

```
// Single-line comment

/*
   Multi-line
   comment
*/

let x = 42; // inline comment
```

---

## Template Mode

Templates are plain text with embedded script blocks between `{% %}` delimiters. Text outside script blocks is emitted as-is. Expressions evaluated inside script blocks are output directly; use `print()` for explicit output control.

### Basic Output

Variable references inside `{% %}` output their value:

```
Hello, {% name %}!
```

With `name` set to `"World"`, produces:

```
Hello, World!
```

### Script Blocks

Use script blocks for logic — variables, conditions, loops:

```
{% let x = 10; %}
The value is {% x %}.
```

Output:

```
The value is 10.
```

### Conditionals in Templates

```
{% if (user.role == "admin") { %}
    <div class="admin-panel">Welcome, admin!</div>
{% } else { %}
    <div>Welcome, user!</div>
{% } %}
```

### Loops in Templates

```
<table>
{% foreach (rows as row, i) { %}
    <tr class="{% if (i % 2 == 0) { print("even"); } else { print("odd"); } %}">
        <td>{% row.name %}</td>
        <td>{% toString(row.score) %}</td>
    </tr>
{% } %}
</table>
```

Expressions inside `{% %}` blocks auto-output their value, even inside loops and conditionals. Use `print()` when you need explicit control (e.g., conditional output within a single script block).

### Using `print()`

`print()` writes directly to the template output. It's useful inside multi-statement script blocks where you need fine-grained control:

```
{%
    let items = ["one", "two", "three"];
    foreach (items as item, i) {
        if (i > 0) { print(", "); }
        print(item);
    }
%}
```

Output: `one, two, three`

### Full Template Example

```
<html>
<head><title>{% title %}</title></head>
<body>
    <h1>{% title %}</h1>

    {% if (len(notifications) > 0) { %}
    <div class="alerts">
        {% foreach (notifications as note) { %}
        <div class="alert">{% note %}</div>
        {% } %}
    </div>
    {% } %}

    <ul>
    {% foreach (items as item) { %}
        <li>
            {% item.name %} - ${% toString(item.price) %}
            {% if (item.onSale) { %}
                <span class="sale">SALE!</span>
            {% } %}
        </li>
    {% } %}
    </ul>

    <footer>{% toString(len(items)) %} items listed</footer>
</body>
</html>
```

---

## Built-in Functions

### Output

| Function          | Description                                | Example            |
| ----------------- | ------------------------------------------ | ------------------ |
| `print(val, ...)` | Write values to template output            | `print("hello")`   |
| `log(val, ...)`   | Write values to logger (stdout by default) | `log("debug:", x)` |

### Type Conversion

| Function          | Description                 | Example                       |
| ----------------- | --------------------------- | ----------------------------- |
| `toString(val)`   | Convert any value to string | `toString(42)` → `"42"`       |
| `parseInt(str)`   | Parse string to integer     | `parseInt("42")` → `42`       |
| `parseFloat(str)` | Parse string to decimal     | `parseFloat("3.14")` → `3.14` |
| `type(val)`       | Get type name as string     | `type(42)` → `"INTEGER"`      |

### String Functions

| Function                     | Description                               | Example                                   |
| ---------------------------- | ----------------------------------------- | ----------------------------------------- |
| `len(str)`                   | String length                             | `len("hello")` → `5`                      |
| `toUpper(str)`               | Convert to uppercase                      | `toUpper("hello")` → `"HELLO"`            |
| `toLower(str)`               | Convert to lowercase                      | `toLower("HELLO")` → `"hello"`            |
| `trim(str)`                  | Remove leading/trailing whitespace        | `trim("  hi  ")` → `"hi"`                 |
| `contains(str, sub)`         | Check if string contains substring        | `contains("hello", "ell")` → `true`       |
| `startsWith(str, prefix)`    | Check string prefix                       | `startsWith("hello", "he")` → `true`      |
| `endsWith(str, suffix)`      | Check string suffix                       | `endsWith("hello", "lo")` → `true`        |
| `indexOf(str, sub)`          | Find substring position (-1 if not found) | `indexOf("hello", "ll")` → `2`            |
| `replace(str, old, new)`     | Replace all occurrences                   | `replace("aabb", "a", "x")` → `"xxbb"`    |
| `substring(str, start)`      | Extract from start to end                 | `substring("hello", 2)` → `"llo"`         |
| `substring(str, start, end)` | Extract from start to end (exclusive)     | `substring("hello", 1, 4)` → `"ell"`      |
| `split(str, delim)`          | Split string into array                   | `split("a,b,c", ",")` → `["a", "b", "c"]` |
| `join(arr, sep)`             | Join array elements into string           | `join([1, 2, 3], "-")` → `"1-2-3"`        |

### Array Functions

| Function           | Description                           | Example                                                |
| ------------------ | ------------------------------------- | ------------------------------------------------------ |
| `len(arr)`         | Array length                          | `len([1,2,3])` → `3`                                   |
| `append(arr, val)` | Append element (mutates array)        | `append(arr, 4)`                                       |
| `map(arr, fn)`     | Transform each element                | `map([1,2,3], fn(x) { return x * 2; })` → `[2,4,6]`    |
| `filter(arr, fn)`  | Keep elements where fn returns truthy | `filter([1,2,3,4], fn(x) { return x > 2; })` → `[3,4]` |

### Hash Functions

| Function       | Description                        | Example                                 |
| -------------- | ---------------------------------- | --------------------------------------- |
| `len(hash)`    | Number of key-value pairs          | `len({"a": 1})` → `1`                   |
| `keys(hash)`   | Get keys array (insertion order)   | `keys({"b": 2, "a": 1})` → `["b", "a"]` |
| `values(hash)` | Get values array (insertion order) | `values({"b": 2, "a": 1})` → `[2, 1]`   |

### Math Functions

| Function     | Description              | Example            |
| ------------ | ------------------------ | ------------------ |
| `floor(num)` | Round down to integer    | `floor(3.7)` → `3` |
| `ceil(num)`  | Round up to integer      | `ceil(3.2)` → `4`  |
| `round(num)` | Round to nearest integer | `round(3.5)` → `4` |
| `abs(num)`   | Absolute value           | `abs(-5)` → `5`    |

---

## Go API Reference

### Convenience Helpers

The simplest way to evaluate scripts and templates. These handle lexing, parsing, context creation, and variable injection in a single call.

#### Package-Level Functions

For the common case where no custom functions are needed:

```go
// Run a script, get the result Object
result, err := evaluator.RunScript(`return a + b;`, evaluator.Vars{"a": 1, "b": 2})

// Render a template, get the output string
output, err := evaluator.RunTemplate(`Hello, {% name %}!`, evaluator.Vars{"name": "World"})
```

#### Evaluator Methods

When you need custom functions, use the evaluator's `RunScript` and `RunTemplate` methods:

```go
eval := evaluator.New()
eval.RegisterFunction("upper", myUpperFunc)

result, err := eval.RunScript(`return upper(name);`, evaluator.Vars{"name": "alice"})
output, err := eval.RunTemplate(`Hello {% upper(name) %}!`, evaluator.Vars{"name": "world"})
```

#### `Vars` Type

`evaluator.Vars` is a `map[string]any` that supports automatic Go-to-script type conversion:

| Go type                                  | Script Object                           |
| ---------------------------------------- | --------------------------------------- |
| `nil`                                    | `Null`                                  |
| `string`                                 | `*StringValue`                          |
| `bool`                                   | `*BooleanValue`                         |
| `int`, `int8`, `int16`, `int32`, `int64` | `*IntegerValue`                         |
| `float32`, `float64`                     | `*DecimalValue`                         |
| `time.Time`, `*time.Time`                | `*DateTimeValue` (nil pointer → `Null`) |
| `[]any`                                  | `*ArrayValue` (recursive)               |
| `map[string]any`                         | `*HashValue` (recursive)                |
| any `evaluator.Object`                   | pass-through                            |

Multiple `Vars` maps can be passed — later maps overwrite earlier ones:

```go
result, err := evaluator.RunScript(`return x;`,
    evaluator.Vars{"x": 1},
    evaluator.Vars{"x": 2}, // overwrites
)
// result is 2
```

#### `ToObject`

Convert individual Go values when working with the lower-level API:

```go
obj, err := evaluator.ToObject(42)           // *IntegerValue
obj, err := evaluator.ToObject("hello")      // *StringValue
obj, err := evaluator.ToObject([]any{1, 2})  // *ArrayValue
```

#### JSON Interop

Since `encoding/json.Unmarshal` produces `map[string]any` and `[]any`, JSON data works directly:

```go
var data map[string]any
json.Unmarshal(jsonBytes, &data)
output, err := evaluator.RunTemplate(tmpl, evaluator.Vars{"data": data})
```

### Pipeline (Low-Level)

Every evaluation follows the same three-step pipeline:

```go
// 1. Lex
l := lexer.NewScript(source)   // or lexer.NewTemplate(source)

// 2. Parse
p := parser.New(l)
program, err := p.Parse()

// 3. Evaluate
eval := evaluator.New()
ctx := evaluator.NewExecutionContext(program)
result, err := eval.Evaluate(ctx)       // returns Object
// or
output, err := eval.EvaluateString(ctx) // returns template string
```

Use `lexer.NewScript()` for pure scripts and `eval.Evaluate()` to get the result as an `Object`.

Use `lexer.NewTemplate()` for templates and `eval.EvaluateString()` to get the rendered output as a string.

### Injecting Variables

Set variables on the root scope before evaluation:

```go
ctx := evaluator.NewExecutionContext(program)

// Primitives
ctx.RootScope.SetLocal("name", &evaluator.StringValue{Value: "Alice"})
ctx.RootScope.SetLocal("age", &evaluator.IntegerValue{Value: 30})
ctx.RootScope.SetLocal("score", &evaluator.DecimalValue{Value: 95.5})
ctx.RootScope.SetLocal("active", &evaluator.BooleanValue{Value: true})

// Array
ctx.RootScope.SetLocal("tags", &evaluator.ArrayValue{
    Elements: []evaluator.Object{
        &evaluator.StringValue{Value: "go"},
        &evaluator.StringValue{Value: "template"},
    },
})

// Hash (ordered map)
user := evaluator.NewHashValue()
user.Set(&evaluator.StringValue{Value: "name"}, &evaluator.StringValue{Value: "Alice"})
user.Set(&evaluator.StringValue{Value: "role"}, &evaluator.StringValue{Value: "admin"})
ctx.RootScope.SetLocal("user", user)
```

### Registering Custom Functions

```go
eval := evaluator.New()

eval.RegisterFunction("count", func(
    ctx *evaluator.ExecutionContext,
    scope *evaluator.Scope,
    args ...evaluator.Object,
) (evaluator.Object, error) {
    if len(args) != 1 {
        return nil, fmt.Errorf("count: expected 1 argument, got %d", len(args))
    }
    arr, ok := args[0].(*evaluator.ArrayValue)
    if !ok {
        return nil, fmt.Errorf("count: expected array, got %s", args[0].Type())
    }
    return &evaluator.IntegerValue{Value: len(arr.Elements)}, nil
})
```

### Sharing Scope Between Evaluations

Use `NewExecutionContextWithScope` to share a scope across multiple evaluation runs:

```go
scope := evaluator.NewScope()
scope.SetLocal("counter", &evaluator.IntegerValue{Value: 0})

ctx1 := evaluator.NewExecutionContextWithScope(program1, scope)
eval.Evaluate(ctx1)

ctx2 := evaluator.NewExecutionContextWithScope(program2, scope)
eval.Evaluate(ctx2) // sees counter modifications from program1
```

### Metadata

Store arbitrary Go values on the execution context, accessible from custom functions:

```go
ctx.Metadata["requestID"] = "abc-123"
ctx.Metadata["db"] = myDBConnection

eval.RegisterFunction("getRequestID", func(
    ctx *evaluator.ExecutionContext,
    scope *evaluator.Scope,
    args ...evaluator.Object,
) (evaluator.Object, error) {
    id := ctx.Metadata["requestID"].(string)
    return &evaluator.StringValue{Value: id}, nil
})
```

### Execution Security Limits

When evaluating untrusted templates (e.g. user-provided input in a SaaS application), configure execution limits to prevent denial-of-service:

```go
ctx := evaluator.NewExecutionContext(program)

// Defaults shown — override as needed:
ctx.MaxSteps = 100_000     // Total AST node evaluations
ctx.MaxDepth = 256         // Maximum function call nesting
ctx.MaxArraySize = 10_000  // Maximum array length
```

Exceeding any limit returns an error:

| Limit          | Error Message                            | What It Prevents                              |
| -------------- | ---------------------------------------- | --------------------------------------------- |
| `MaxSteps`     | `execution limit exceeded: 100000 steps` | Infinite loops, runaway computation           |
| `MaxDepth`     | `maximum call depth exceeded: 256`       | Stack overflow from deep/infinite recursion   |
| `MaxArraySize` | `maximum array size exceeded: 10000`     | Memory exhaustion from unbounded array growth |

Set any limit to `0` to disable it.

### Runtime Error Locations

Set `ctx.Source` to the original source string to get rich error messages with source location:

```go
source := `let x = 1 / 0;`

l := lexer.NewScript(source)
p := parser.New(l)
program, _ := p.Parse()

ctx := evaluator.NewExecutionContext(program)
ctx.Source = source  // enables source-located errors

_, err := evaluator.New().Evaluate(ctx)
// Error output:
//   division by zero
//     at line 1, column 13
//    |
//  1 | let x = 1 / 0;
//    |             ^
```

Without `ctx.Source`, errors still report the message but without the source context.

### Error Handling

The library provides three error types with rich formatting:

- **`lexer.TokenError`** — Malformed tokens (unterminated strings, unexpected characters)
- **`parser.ParseError`** — Syntax errors (missing braces, invalid expressions)
- **`evaluator.RuntimeError`** — Execution errors (type mismatches, undefined variables, division by zero)

All three include line/column information and formatted source context when available. Use `errors.As` to inspect them:

```go
import "errors"

var runtimeErr *evaluator.RuntimeError
if errors.As(err, &runtimeErr) {
    fmt.Println("Line:", runtimeErr.Line)
    fmt.Println("Column:", runtimeErr.Column)
}
```

The parser accumulates all errors rather than failing on the first one, so a single `Parse()` call can report multiple issues.

---

## Complete Examples

### Fizzbuzz

```
fn fizzbuzz(n) {
    let i = 1;
    while (i <= n) {
        if (i % 15 == 0) {
            log("FizzBuzz");
        } else if (i % 3 == 0) {
            log("Fizz");
        } else if (i % 5 == 0) {
            log("Buzz");
        } else {
            log(i);
        }
        i += 1;
    }
}

fizzbuzz(30);
```

### Functional Pipeline

```
let numbers = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10];

let result = join(
    map(
        filter(numbers, fn(x) { return x % 2 == 0; }),
        fn(x) { return x * x; }
    ),
    ", "
);

return result;  // "4, 16, 36, 64, 100"
```

### Email Template

```
Dear {% user.name %},

{% if (len(orders) > 0) { %}
You have {% toString(len(orders)) %} recent order(s):

{% foreach (orders as order, i) { %}
  {% toString(i + 1) %}. {% order.product %} — ${% toString(order.total) %}
{% } %}

Total: ${%
    let sum = 0;
    foreach (orders as order) {
        sum += order.total;
    }
    print(toString(sum));
%}
{% } else { %}
You have no recent orders.
{% } %}

Best regards,
{% company %}
```

### Accumulator Pattern

```
fn makeAccumulator(initial) {
    let total = initial;
    return {
        "add": fn(n) { total = total + n; return total; },
        "get": fn() { return total; }
    };
}

let acc = makeAccumulator(0);
acc.add(10);   // 10
acc.add(20);   // 30
acc.add(5);    // 35
return acc.get();  // 35
```

## License

See [LICENSE](LICENSE) for details.
