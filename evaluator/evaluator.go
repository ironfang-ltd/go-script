package evaluator

import (
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/ironfang-ltd/go-script/lexer"
	"github.com/ironfang-ltd/go-script/parser"
)

var (
	Null = &NullValue{}
)

type Node interface {
	Debug() string
}

type Evaluator struct {
	functions map[string]*BuiltInFunction
}

func New() *Evaluator {
	e := &Evaluator{
		functions: make(map[string]*BuiltInFunction),
	}

	e.RegisterFunction("log", func(ctx *ExecutionContext, scope *Scope, args ...Object) (Object, error) {
		for _, arg := range args {
			_, _ = ctx.Logger.WriteString(arg.Debug())
			_, _ = ctx.Logger.WriteString("\n")
		}
		return Null, nil
	})

	e.RegisterFunction("print", func(ctx *ExecutionContext, scope *Scope, args ...Object) (Object, error) {
		for _, arg := range args {
			ctx.output.WriteString(arg.Debug())
		}
		return Null, nil
	})

	e.RegisterFunction("append", func(ctx *ExecutionContext, scope *Scope, args ...Object) (Object, error) {

		if len(args) != 2 {
			return nil, fmt.Errorf("expected 2 arguments, got %d", len(args))
		}

		arrValue, ok := args[0].(*ArrayValue)
		if !ok {
			return nil, fmt.Errorf("expected array, got %s", args[0].Type())
		}

		if ctx.MaxArraySize > 0 && len(arrValue.Elements) >= ctx.MaxArraySize {
			return nil, fmt.Errorf("maximum array size exceeded: %d", ctx.MaxArraySize)
		}

		arrValue.Elements = append(arrValue.Elements, args[1])

		return arrValue, nil
	})

	e.RegisterFunction("len", func(ctx *ExecutionContext, scope *Scope, args ...Object) (Object, error) {
		if len(args) != 1 {
			return nil, fmt.Errorf("len: expected 1 argument, got %d", len(args))
		}
		switch v := args[0].(type) {
		case *StringValue:
			return &IntegerValue{Value: len(v.Value)}, nil
		case *ArrayValue:
			return &IntegerValue{Value: len(v.Elements)}, nil
		case *HashValue:
			return &IntegerValue{Value: len(v.Pairs)}, nil
		default:
			return nil, fmt.Errorf("len: unsupported type %s", args[0].Type())
		}
	})

	e.RegisterFunction("split", func(ctx *ExecutionContext, scope *Scope, args ...Object) (Object, error) {
		if len(args) != 2 {
			return nil, fmt.Errorf("split: expected 2 arguments, got %d", len(args))
		}
		str, ok := args[0].(*StringValue)
		if !ok {
			return nil, fmt.Errorf("split: first argument must be a string, got %s", args[0].Type())
		}
		delim, ok := args[1].(*StringValue)
		if !ok {
			return nil, fmt.Errorf("split: second argument must be a string, got %s", args[1].Type())
		}
		parts := strings.Split(str.Value, delim.Value)
		elements := make([]Object, len(parts))
		for i, p := range parts {
			elements[i] = &StringValue{Value: p}
		}
		return &ArrayValue{Elements: elements}, nil
	})

	e.RegisterFunction("trim", func(ctx *ExecutionContext, scope *Scope, args ...Object) (Object, error) {
		if len(args) != 1 {
			return nil, fmt.Errorf("trim: expected 1 argument, got %d", len(args))
		}
		str, ok := args[0].(*StringValue)
		if !ok {
			return nil, fmt.Errorf("trim: argument must be a string, got %s", args[0].Type())
		}
		return &StringValue{Value: strings.TrimSpace(str.Value)}, nil
	})

	e.RegisterFunction("toUpper", func(ctx *ExecutionContext, scope *Scope, args ...Object) (Object, error) {
		if len(args) != 1 {
			return nil, fmt.Errorf("toUpper: expected 1 argument, got %d", len(args))
		}
		str, ok := args[0].(*StringValue)
		if !ok {
			return nil, fmt.Errorf("toUpper: argument must be a string, got %s", args[0].Type())
		}
		return &StringValue{Value: strings.ToUpper(str.Value)}, nil
	})

	e.RegisterFunction("toLower", func(ctx *ExecutionContext, scope *Scope, args ...Object) (Object, error) {
		if len(args) != 1 {
			return nil, fmt.Errorf("toLower: expected 1 argument, got %d", len(args))
		}
		str, ok := args[0].(*StringValue)
		if !ok {
			return nil, fmt.Errorf("toLower: argument must be a string, got %s", args[0].Type())
		}
		return &StringValue{Value: strings.ToLower(str.Value)}, nil
	})

	e.RegisterFunction("contains", func(ctx *ExecutionContext, scope *Scope, args ...Object) (Object, error) {
		if len(args) != 2 {
			return nil, fmt.Errorf("contains: expected 2 arguments, got %d", len(args))
		}
		str, ok := args[0].(*StringValue)
		if !ok {
			return nil, fmt.Errorf("contains: first argument must be a string, got %s", args[0].Type())
		}
		substr, ok := args[1].(*StringValue)
		if !ok {
			return nil, fmt.Errorf("contains: second argument must be a string, got %s", args[1].Type())
		}
		return &BooleanValue{Value: strings.Contains(str.Value, substr.Value)}, nil
	})

	e.RegisterFunction("startsWith", func(ctx *ExecutionContext, scope *Scope, args ...Object) (Object, error) {
		if len(args) != 2 {
			return nil, fmt.Errorf("startsWith: expected 2 arguments, got %d", len(args))
		}
		str, ok := args[0].(*StringValue)
		if !ok {
			return nil, fmt.Errorf("startsWith: first argument must be a string, got %s", args[0].Type())
		}
		prefix, ok := args[1].(*StringValue)
		if !ok {
			return nil, fmt.Errorf("startsWith: second argument must be a string, got %s", args[1].Type())
		}
		return &BooleanValue{Value: strings.HasPrefix(str.Value, prefix.Value)}, nil
	})

	e.RegisterFunction("endsWith", func(ctx *ExecutionContext, scope *Scope, args ...Object) (Object, error) {
		if len(args) != 2 {
			return nil, fmt.Errorf("endsWith: expected 2 arguments, got %d", len(args))
		}
		str, ok := args[0].(*StringValue)
		if !ok {
			return nil, fmt.Errorf("endsWith: first argument must be a string, got %s", args[0].Type())
		}
		suffix, ok := args[1].(*StringValue)
		if !ok {
			return nil, fmt.Errorf("endsWith: second argument must be a string, got %s", args[1].Type())
		}
		return &BooleanValue{Value: strings.HasSuffix(str.Value, suffix.Value)}, nil
	})

	e.RegisterFunction("indexOf", func(ctx *ExecutionContext, scope *Scope, args ...Object) (Object, error) {
		if len(args) != 2 {
			return nil, fmt.Errorf("indexOf: expected 2 arguments, got %d", len(args))
		}
		str, ok := args[0].(*StringValue)
		if !ok {
			return nil, fmt.Errorf("indexOf: first argument must be a string, got %s", args[0].Type())
		}
		substr, ok := args[1].(*StringValue)
		if !ok {
			return nil, fmt.Errorf("indexOf: second argument must be a string, got %s", args[1].Type())
		}
		return &IntegerValue{Value: strings.Index(str.Value, substr.Value)}, nil
	})

	e.RegisterFunction("replace", func(ctx *ExecutionContext, scope *Scope, args ...Object) (Object, error) {
		if len(args) != 3 {
			return nil, fmt.Errorf("replace: expected 3 arguments, got %d", len(args))
		}
		str, ok := args[0].(*StringValue)
		if !ok {
			return nil, fmt.Errorf("replace: first argument must be a string, got %s", args[0].Type())
		}
		old, ok := args[1].(*StringValue)
		if !ok {
			return nil, fmt.Errorf("replace: second argument must be a string, got %s", args[1].Type())
		}
		newStr, ok := args[2].(*StringValue)
		if !ok {
			return nil, fmt.Errorf("replace: third argument must be a string, got %s", args[2].Type())
		}
		return &StringValue{Value: strings.ReplaceAll(str.Value, old.Value, newStr.Value)}, nil
	})

	e.RegisterFunction("substring", func(ctx *ExecutionContext, scope *Scope, args ...Object) (Object, error) {
		if len(args) < 2 || len(args) > 3 {
			return nil, fmt.Errorf("substring: expected 2 or 3 arguments, got %d", len(args))
		}
		str, ok := args[0].(*StringValue)
		if !ok {
			return nil, fmt.Errorf("substring: first argument must be a string, got %s", args[0].Type())
		}
		start, ok := args[1].(*IntegerValue)
		if !ok {
			return nil, fmt.Errorf("substring: second argument must be an integer, got %s", args[1].Type())
		}
		s := start.Value
		if s < 0 {
			s = 0
		}
		if s > len(str.Value) {
			s = len(str.Value)
		}
		end := len(str.Value)
		if len(args) == 3 {
			endVal, ok := args[2].(*IntegerValue)
			if !ok {
				return nil, fmt.Errorf("substring: third argument must be an integer, got %s", args[2].Type())
			}
			end = endVal.Value
			if end < s {
				end = s
			}
			if end > len(str.Value) {
				end = len(str.Value)
			}
		}
		return &StringValue{Value: str.Value[s:end]}, nil
	})

	e.RegisterFunction("keys", func(ctx *ExecutionContext, scope *Scope, args ...Object) (Object, error) {
		if len(args) != 1 {
			return nil, fmt.Errorf("keys: expected 1 argument, got %d", len(args))
		}
		hash, ok := args[0].(*HashValue)
		if !ok {
			return nil, fmt.Errorf("keys: argument must be a hash, got %s", args[0].Type())
		}
		ordered := hash.OrderedPairs()
		elements := make([]Object, 0, len(ordered))
		for _, pair := range ordered {
			elements = append(elements, pair.Key)
		}
		return &ArrayValue{Elements: elements}, nil
	})

	e.RegisterFunction("values", func(ctx *ExecutionContext, scope *Scope, args ...Object) (Object, error) {
		if len(args) != 1 {
			return nil, fmt.Errorf("values: expected 1 argument, got %d", len(args))
		}
		hash, ok := args[0].(*HashValue)
		if !ok {
			return nil, fmt.Errorf("values: argument must be a hash, got %s", args[0].Type())
		}
		ordered := hash.OrderedPairs()
		elements := make([]Object, 0, len(ordered))
		for _, pair := range ordered {
			elements = append(elements, pair.Value)
		}
		return &ArrayValue{Elements: elements}, nil
	})

	e.RegisterFunction("type", func(ctx *ExecutionContext, scope *Scope, args ...Object) (Object, error) {
		if len(args) != 1 {
			return nil, fmt.Errorf("type: expected 1 argument, got %d", len(args))
		}
		return &StringValue{Value: string(args[0].Type())}, nil
	})

	e.RegisterFunction("toString", func(ctx *ExecutionContext, scope *Scope, args ...Object) (Object, error) {
		if len(args) != 1 {
			return nil, fmt.Errorf("toString: expected 1 argument, got %d", len(args))
		}
		return &StringValue{Value: args[0].Debug()}, nil
	})

	e.RegisterFunction("parseInt", func(ctx *ExecutionContext, scope *Scope, args ...Object) (Object, error) {
		if len(args) != 1 {
			return nil, fmt.Errorf("parseInt: expected 1 argument, got %d", len(args))
		}
		str, ok := args[0].(*StringValue)
		if !ok {
			return nil, fmt.Errorf("parseInt: argument must be a string, got %s", args[0].Type())
		}
		val, err := strconv.ParseInt(str.Value, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("parseInt: %s", err)
		}
		return &IntegerValue{Value: int(val)}, nil
	})

	e.RegisterFunction("parseFloat", func(ctx *ExecutionContext, scope *Scope, args ...Object) (Object, error) {
		if len(args) != 1 {
			return nil, fmt.Errorf("parseFloat: expected 1 argument, got %d", len(args))
		}
		str, ok := args[0].(*StringValue)
		if !ok {
			return nil, fmt.Errorf("parseFloat: argument must be a string, got %s", args[0].Type())
		}
		val, err := strconv.ParseFloat(str.Value, 64)
		if err != nil {
			return nil, fmt.Errorf("parseFloat: %s", err)
		}
		return &DecimalValue{Value: val}, nil
	})

	e.RegisterFunction("join", func(ctx *ExecutionContext, scope *Scope, args ...Object) (Object, error) {
		if len(args) != 2 {
			return nil, fmt.Errorf("join: expected 2 arguments, got %d", len(args))
		}
		arr, ok := args[0].(*ArrayValue)
		if !ok {
			return nil, fmt.Errorf("join: first argument must be an array, got %s", args[0].Type())
		}
		sep, ok := args[1].(*StringValue)
		if !ok {
			return nil, fmt.Errorf("join: second argument must be a string, got %s", args[1].Type())
		}
		parts := make([]string, len(arr.Elements))
		for i, el := range arr.Elements {
			parts[i] = el.Debug()
		}
		return &StringValue{Value: strings.Join(parts, sep.Value)}, nil
	})

	e.RegisterFunction("map", func(ctx *ExecutionContext, scope *Scope, args ...Object) (Object, error) {
		if len(args) != 2 {
			return nil, fmt.Errorf("map: expected 2 arguments, got %d", len(args))
		}
		arr, ok := args[0].(*ArrayValue)
		if !ok {
			return nil, fmt.Errorf("map: first argument must be an array, got %s", args[0].Type())
		}
		result := make([]Object, len(arr.Elements))
		for i, el := range arr.Elements {
			val, err := e.applyFunction(ctx, scope, args[1], []Object{el})
			if err != nil {
				return nil, err
			}
			result[i] = val
		}
		return &ArrayValue{Elements: result}, nil
	})

	e.RegisterFunction("filter", func(ctx *ExecutionContext, scope *Scope, args ...Object) (Object, error) {
		if len(args) != 2 {
			return nil, fmt.Errorf("filter: expected 2 arguments, got %d", len(args))
		}
		arr, ok := args[0].(*ArrayValue)
		if !ok {
			return nil, fmt.Errorf("filter: first argument must be an array, got %s", args[0].Type())
		}
		var result []Object
		for _, el := range arr.Elements {
			val, err := e.applyFunction(ctx, scope, args[1], []Object{el})
			if err != nil {
				return nil, err
			}
			if isTruthy(val) {
				result = append(result, el)
			}
		}
		return &ArrayValue{Elements: result}, nil
	})

	e.RegisterFunction("floor", func(ctx *ExecutionContext, scope *Scope, args ...Object) (Object, error) {
		if len(args) != 1 {
			return nil, fmt.Errorf("floor: expected 1 argument, got %d", len(args))
		}
		switch v := args[0].(type) {
		case *IntegerValue:
			return v, nil
		case *DecimalValue:
			return &IntegerValue{Value: int(math.Floor(v.Value))}, nil
		default:
			return nil, fmt.Errorf("floor: argument must be a number, got %s", args[0].Type())
		}
	})

	e.RegisterFunction("ceil", func(ctx *ExecutionContext, scope *Scope, args ...Object) (Object, error) {
		if len(args) != 1 {
			return nil, fmt.Errorf("ceil: expected 1 argument, got %d", len(args))
		}
		switch v := args[0].(type) {
		case *IntegerValue:
			return v, nil
		case *DecimalValue:
			return &IntegerValue{Value: int(math.Ceil(v.Value))}, nil
		default:
			return nil, fmt.Errorf("ceil: argument must be a number, got %s", args[0].Type())
		}
	})

	e.RegisterFunction("round", func(ctx *ExecutionContext, scope *Scope, args ...Object) (Object, error) {
		if len(args) != 1 {
			return nil, fmt.Errorf("round: expected 1 argument, got %d", len(args))
		}
		switch v := args[0].(type) {
		case *IntegerValue:
			return v, nil
		case *DecimalValue:
			return &IntegerValue{Value: int(math.Round(v.Value))}, nil
		default:
			return nil, fmt.Errorf("round: argument must be a number, got %s", args[0].Type())
		}
	})

	e.RegisterFunction("abs", func(ctx *ExecutionContext, scope *Scope, args ...Object) (Object, error) {
		if len(args) != 1 {
			return nil, fmt.Errorf("abs: expected 1 argument, got %d", len(args))
		}
		switch v := args[0].(type) {
		case *IntegerValue:
			if v.Value < 0 {
				return &IntegerValue{Value: -v.Value}, nil
			}
			return v, nil
		case *DecimalValue:
			return &DecimalValue{Value: math.Abs(v.Value)}, nil
		default:
			return nil, fmt.Errorf("abs: argument must be a number, got %s", args[0].Type())
		}
	})

	return e
}

func (e *Evaluator) RegisterFunction(name string, fn Function) {
	e.functions[name] = &BuiltInFunction{Fn: fn}
}

type ExecutionContext struct {
	Program      *parser.Program
	RootScope    *Scope
	Logger       io.StringWriter
	Metadata     map[string]any
	Source       string
	MaxSteps     int
	MaxDepth     int
	MaxArraySize int
	steps        int
	depth        int
	output       *strings.Builder
	templateMode bool
}

func NewExecutionContext(program *parser.Program) *ExecutionContext {
	return &ExecutionContext{
		Program:      program,
		RootScope:    NewScope(),
		Logger:       os.Stdout,
		Metadata:     make(map[string]any),
		MaxSteps:     100_000,
		MaxDepth:     256,
		MaxArraySize: 10_000,
		output:       &strings.Builder{},
	}
}

func NewExecutionContextWithScope(program *parser.Program, rootScope *Scope) *ExecutionContext {
	return &ExecutionContext{
		Program:      program,
		RootScope:    rootScope,
		Logger:       os.Stdout,
		Metadata:     make(map[string]any),
		MaxSteps:     100_000,
		MaxDepth:     256,
		MaxArraySize: 10_000,
		output:       &strings.Builder{},
	}
}

func runtimeError(ctx *ExecutionContext, token lexer.Token, message string) error {
	if ctx.Source != "" {
		return NewRuntimeError(message, ctx.Source, token.Line, token.Column)
	}
	return fmt.Errorf("%s", message)
}

func (e *Evaluator) Evaluate(ctx *ExecutionContext) (Object, error) {

	var result Object = Null

	for _, statement := range ctx.Program.Statements {
		evalResult, err := e.evaluateNode(ctx, statement, ctx.RootScope)
		if err != nil {
			return nil, err
		}

		if _, ok := evalResult.(*ReturnValue); ok {
			return evalResult, nil
		}

		result = evalResult
	}

	return result, nil
}

func (e *Evaluator) EvaluateString(ctx *ExecutionContext) (string, error) {

	ctx.templateMode = true

	for _, statement := range ctx.Program.Statements {
		evalResult, err := e.evaluateNode(ctx, statement, ctx.RootScope)
		if err != nil {
			return "", err
		}

		if _, ok := evalResult.(*ReturnValue); ok {
			return "", nil
		}
	}

	return ctx.output.String(), nil
}

func shouldWriteTemplateOutput(obj Object) bool {
	if obj == nil {
		return false
	}
	switch obj.Type() {
	case NullObject, ReturnValueObject, FunctionObject, BreakSignalObject, ContinueSignalObject:
		return false
	}
	return true
}

func (e *Evaluator) evaluateNode(ctx *ExecutionContext, node Node, scope *Scope) (Object, error) {
	if ctx.MaxSteps > 0 {
		ctx.steps++
		if ctx.steps > ctx.MaxSteps {
			return nil, fmt.Errorf("execution limit exceeded: %d steps", ctx.MaxSteps)
		}
	}

	switch n := node.(type) {
	case *parser.PrintStatement:
		return e.evaluatePrintStatement(ctx, n)
	case *parser.BlockStatement:
		return e.evaluateBlockStatement(ctx, n, scope)
	case *parser.LetStatement:
		return e.evaluateLetStatement(ctx, n, scope)
	case *parser.AssignmentExpression:
		return e.evaluateAssignmentExpression(ctx, n, scope)
	case *parser.ReturnStatement:
		return e.evaluateReturnStatement(ctx, n, scope)
	case *parser.ExpressionStatement:
		result, err := e.evaluateNode(ctx, n.Expression, scope)
		if err != nil {
			return nil, err
		}
		if ctx.templateMode {
			if _, ok := n.Expression.(*parser.AssignmentExpression); !ok {
				if shouldWriteTemplateOutput(result) {
					ctx.output.WriteString(result.Debug())
				}
			}
			return Null, nil
		}
		return result, nil
	case *parser.ForeachExpression:
		return e.evaluateForEach(ctx, n, scope)
	case *parser.IntegerLiteral:
		return &IntegerValue{Value: n.Value}, nil
	case *parser.FloatLiteral:
		return &DecimalValue{Value: n.Value}, nil
	case *parser.BooleanLiteral:
		return &BooleanValue{Value: n.Value}, nil
	case *parser.StringLiteral:
		return &StringValue{Value: n.Value}, nil
	case *parser.PrefixExpression:
		right, err := e.evaluateNode(ctx, n.Right, scope)
		if err != nil {
			return nil, err
		}

		return e.evaluatePrefixExpression(n.Operator, right)
	case *parser.NullLiteral:
		return Null, nil
	case *parser.WhileExpression:
		return e.evaluateWhileExpression(ctx, n, scope)
	case *parser.BreakStatement:
		return &BreakSignal{}, nil
	case *parser.ContinueStatement:
		return &ContinueSignal{}, nil
	case *parser.InfixExpression:
		// Short-circuit for &&, ||, ??
		switch n.Token.Source {
		case "&&":
			left, err := e.evaluateNode(ctx, n.Left, scope)
			if err != nil {
				return nil, err
			}
			if !isTruthy(left) {
				return left, nil
			}
			return e.evaluateNode(ctx, n.Right, scope)
		case "||":
			left, err := e.evaluateNode(ctx, n.Left, scope)
			if err != nil {
				return nil, err
			}
			if isTruthy(left) {
				return left, nil
			}
			return e.evaluateNode(ctx, n.Right, scope)
		case "??":
			left, err := e.evaluateNode(ctx, n.Left, scope)
			if err != nil {
				return nil, err
			}
			if _, isNull := left.(*NullValue); !isNull {
				return left, nil
			}
			return e.evaluateNode(ctx, n.Right, scope)
		}

		left, err := e.evaluateNode(ctx, n.Left, scope)
		if err != nil {
			return nil, err
		}

		right, err := e.evaluateNode(ctx, n.Right, scope)
		if err != nil {
			return nil, err
		}

		return e.evaluateInfixExpression(ctx, n.Token, left, right)
	case *parser.IfExpression:
		return e.evaluateIfExpression(ctx, n, scope)
	case *parser.Identifier:
		return e.evaluateIdentifier(ctx, n, scope)
	case *parser.FunctionLiteral:
		return e.evaluateFunctionLiteral(n, scope)
	case *parser.CallExpression:
		return e.evaluateCallExpression(ctx, n, scope)
	case *parser.ArrayLiteral:
		return e.evaluateArrayLiteral(ctx, n, scope)
	case *parser.IndexExpression:

		left, err := e.evaluateNode(ctx, n.Left, scope)
		if err != nil {
			return nil, err
		}

		index, err := e.evaluateNode(ctx, n.Index, scope)
		if err != nil {
			return nil, err
		}
		return e.evaluateIndexExpression(left, index)
	case *parser.HashLiteral:
		return e.evaluateHashLiteral(ctx, n, scope)
	case *parser.PropertyExpression:
		_, _, v, err := e.evaluatePropertyExpression(ctx, n, scope)
		return v, err

	default:
		return nil, fmt.Errorf("unknown node type: %T", n)
	}
}

func (e *Evaluator) evaluatePrintStatement(ctx *ExecutionContext, print *parser.PrintStatement) (Object, error) {
	ctx.output.WriteString(print.Value)
	return Null, nil
}

func (e *Evaluator) evaluateForEach(ctx *ExecutionContext, foreach *parser.ForeachExpression, scope *Scope) (Object, error) {
	iterable, err := e.evaluateNode(ctx, foreach.Iterable, scope)
	if err != nil {
		return nil, err
	}

	switch i := iterable.(type) {
	case *ArrayValue:
		return e.evaluateArrayForEach(ctx, foreach, i, scope)
	case *HashValue:
		return e.evaluateHashForEach(ctx, foreach, i, scope)
	default:
		return Null, nil
	}
}

func (e *Evaluator) evaluateArrayForEach(ctx *ExecutionContext, foreach *parser.ForeachExpression, array *ArrayValue, scope *Scope) (Object, error) {
	for i, el := range array.Elements {
		extendedScope := NewChildScope(scope)
		extendedScope.SetLocal(foreach.Variable.Value, el)
		if foreach.Index != nil {
			extendedScope.SetLocal(foreach.Index.Value, &IntegerValue{Value: i})
		}

		result, err := e.evaluateBlockStatement(ctx, foreach.Body, extendedScope)
		if err != nil {
			return nil, err
		}

		if _, ok := result.(*BreakSignal); ok {
			break
		}
		if _, ok := result.(*ContinueSignal); ok {
			continue
		}
		if _, ok := result.(*ReturnValue); ok {
			return result, nil
		}
	}

	return Null, nil
}

func (e *Evaluator) evaluateHashForEach(ctx *ExecutionContext, foreach *parser.ForeachExpression, hash *HashValue, scope *Scope) (Object, error) {
	for _, pair := range hash.OrderedPairs() {
		extendedScope := NewChildScope(scope)
		extendedScope.SetLocal(foreach.Variable.Value, pair.Value)
		if foreach.Index != nil {
			extendedScope.SetLocal(foreach.Index.Value, pair.Key)
		}

		result, err := e.evaluateBlockStatement(ctx, foreach.Body, extendedScope)
		if err != nil {
			return nil, err
		}

		if _, ok := result.(*BreakSignal); ok {
			break
		}
		if _, ok := result.(*ContinueSignal); ok {
			continue
		}
		if _, ok := result.(*ReturnValue); ok {
			return result, nil
		}
	}

	return Null, nil
}

func (e *Evaluator) evaluateWhileExpression(ctx *ExecutionContext, we *parser.WhileExpression, scope *Scope) (Object, error) {
	for {
		condition, err := e.evaluateNode(ctx, we.Condition, scope)
		if err != nil {
			return nil, err
		}

		if !isTruthy(condition) {
			break
		}

		result, err := e.evaluateBlockStatement(ctx, we.Body, scope)
		if err != nil {
			return nil, err
		}

		if _, ok := result.(*BreakSignal); ok {
			break
		}
		if _, ok := result.(*ContinueSignal); ok {
			continue
		}
		if _, ok := result.(*ReturnValue); ok {
			return result, nil
		}
	}

	return Null, nil
}

func (e *Evaluator) evaluateBlockStatement(ctx *ExecutionContext, block *parser.BlockStatement, scope *Scope) (Object, error) {
	var result Object

	for _, statement := range block.Statements {
		evalResult, err := e.evaluateNode(ctx, statement, scope)
		if err != nil {
			return nil, err
		}

		// Propagate signals immediately
		switch evalResult.(type) {
		case *ReturnValue, *BreakSignal, *ContinueSignal:
			return evalResult, nil
		}

		result = evalResult
	}

	return result, nil
}

func (e *Evaluator) evaluateLetStatement(ctx *ExecutionContext, let *parser.LetStatement, scope *Scope) (Object, error) {
	prevTM := ctx.templateMode
	ctx.templateMode = false
	val, err := e.evaluateNode(ctx, let.Value, scope)
	ctx.templateMode = prevTM
	if err != nil {
		return nil, err
	}

	scope.SetLocal(let.Name.Value, val)

	return val, nil
}

func (e *Evaluator) evaluateAssignmentExpression(ctx *ExecutionContext, assign *parser.AssignmentExpression, scope *Scope) (Object, error) {

	prevTM := ctx.templateMode
	ctx.templateMode = false
	defer func() { ctx.templateMode = prevTM }()

	if ident, ok := assign.Left.(*parser.Identifier); ok {

		right, err := e.evaluateNode(ctx, assign.Right, scope)
		if err != nil {
			return nil, err
		}

		assigned := scope.Assign(ident.Value, right)
		if !assigned {
			return nil, runtimeError(ctx, assign.Token, fmt.Sprintf("identifier not found in scope: %s", ident.Value))
		}
		return right, nil
	}

	if propExpr, ok := assign.Left.(*parser.PropertyExpression); ok {
		parent, idx, _, err := e.evaluatePropertyExpression(ctx, propExpr, scope)
		if err != nil {
			return nil, err
		}

		if idx == Null || parent == Null {
			return nil, fmt.Errorf("cannot assign to property: left side evaluated to null")
		}

		hashValue, ok := parent.(*HashValue)
		if !ok {
			return nil, fmt.Errorf("left side of property expression must be a hash, got %T", parent)
		}

		right, err := e.evaluateNode(ctx, assign.Right, scope)
		if err != nil {
			return nil, err
		}

		if err := hashValue.Set(idx, right); err != nil {
			return nil, err
		}

		return right, nil
	}

	if indexExpr, ok := assign.Left.(*parser.IndexExpression); ok {

		left, err := e.evaluateNode(ctx, indexExpr.Left, scope)
		if err != nil {
			return nil, err
		}

		index, err := e.evaluateNode(ctx, indexExpr.Index, scope)
		if err != nil {
			return nil, err
		}

		right, err := e.evaluateNode(ctx, assign.Right, scope)
		if err != nil {
			return nil, err
		}

		if left == Null || index == Null {
			return nil, fmt.Errorf("index expression left side or index evaluated to null")
		}

		if arrayValue, ok := left.(*ArrayValue); ok {

			if i, ok := index.(*IntegerValue); ok {
				if i.Value < 0 || i.Value >= len(arrayValue.Elements) {
					return nil, fmt.Errorf("index out of bounds: %d", i.Value)
				}
				arrayValue.Elements[i.Value] = right
			} else {
				return nil, fmt.Errorf("index must be an integer, got %T", index)
			}

		} else if hashValue, ok := left.(*HashValue); ok {
			if err := hashValue.Set(index, right); err != nil {
				return nil, err
			}
		} else {
			return nil, fmt.Errorf("left side of index expression must be an array or hash, got %T", left)
		}

		return right, nil
	}

	return nil, fmt.Errorf("unknown expression type in assignment: %T", assign.Left)
}

func (e *Evaluator) evaluateReturnStatement(ctx *ExecutionContext, ret *parser.ReturnStatement, scope *Scope) (Object, error) {
	prevTM := ctx.templateMode
	ctx.templateMode = false
	val, err := e.evaluateNode(ctx, ret.Value, scope)
	ctx.templateMode = prevTM
	if err != nil {
		return nil, err
	}

	return &ReturnValue{Value: val}, nil
}

func (e *Evaluator) evaluatePrefixExpression(operator string, right Object) (Object, error) {
	switch operator {
	case "!":
		return e.evaluateBangOperatorExpression(right)
	case "-":
		return e.evaluateMinusPrefixOperatorExpression(right)
	default:
		return nil, fmt.Errorf("unknown operator: %s", operator)
	}
}

func (e *Evaluator) evaluateBangOperatorExpression(right Object) (Object, error) {
	switch r := right.(type) {
	case *BooleanValue:
		return &BooleanValue{Value: !r.Value}, nil
	default:
		return &BooleanValue{Value: false}, nil
	}
}

func (e *Evaluator) evaluateMinusPrefixOperatorExpression(right Object) (Object, error) {
	switch r := right.(type) {
	case *IntegerValue:
		if r.Value == math.MinInt {
			return nil, fmt.Errorf("integer overflow")
		}
		return &IntegerValue{Value: -r.Value}, nil
	case *DecimalValue:
		return &DecimalValue{Value: -r.Value}, nil
	default:
		return nil, fmt.Errorf("unknown operator: -%T", r)
	}
}

func (e *Evaluator) evaluateInfixExpression(ctx *ExecutionContext, token lexer.Token, left, right Object) (Object, error) {
	operator := token.Source

	if i1, ok := left.(*IntegerValue); ok {
		if i2, ok := right.(*IntegerValue); ok {
			return e.evaluateIntegerInfixExpression(ctx, token, i1, i2)
		}
	}

	if b1, ok := left.(*BooleanValue); ok {
		if b2, ok := right.(*BooleanValue); ok {
			return e.evaluateBooleanInfixExpression(operator, b1, b2)
		}
	}

	if d1, ok := left.(*DecimalValue); ok {
		if d2, ok := right.(*DecimalValue); ok {
			return e.evaluateDecimalInfixExpression(ctx, token, d1, d2)
		}
	}

	// Integer op Decimal → promote integer to decimal
	if i, ok := left.(*IntegerValue); ok {
		if d, ok := right.(*DecimalValue); ok {
			return e.evaluateDecimalInfixExpression(ctx, token, &DecimalValue{Value: float64(i.Value)}, d)
		}
	}

	// Decimal op Integer → promote integer to decimal
	if d, ok := left.(*DecimalValue); ok {
		if i, ok := right.(*IntegerValue); ok {
			return e.evaluateDecimalInfixExpression(ctx, token, d, &DecimalValue{Value: float64(i.Value)})
		}
	}

	if s1, ok := left.(*StringValue); ok {
		if s2, ok := right.(*StringValue); ok {
			return e.evaluateStringInfixExpression(operator, s1, s2)
		}
	}

	// String auto-coercion: "str" + other → "str" + other.Debug()
	if operator == "+" {
		if s, ok := left.(*StringValue); ok {
			return &StringValue{Value: s.Value + right.Debug()}, nil
		}
		if s, ok := right.(*StringValue); ok {
			return &StringValue{Value: left.Debug() + s.Value}, nil
		}
	}

	if operator == "==" {
		return &BooleanValue{Value: left == right}, nil
	}

	if operator == "!=" {
		return &BooleanValue{Value: left != right}, nil
	}

	if left.Type() != right.Type() {
		return nil, runtimeError(ctx, token, fmt.Sprintf("type mismatch: %s %s %s", left.Type(), operator, right.Type()))
	}

	return nil, runtimeError(ctx, token, fmt.Sprintf("unknown operator: %s %s %s", left.Type(), operator, right.Type()))
}

func (e *Evaluator) evaluateIntegerInfixExpression(ctx *ExecutionContext, token lexer.Token, l, r *IntegerValue) (Object, error) {
	operator := token.Source

	switch operator {
	case "+":
		result := l.Value + r.Value
		if (r.Value > 0 && result < l.Value) || (r.Value < 0 && result > l.Value) {
			return nil, runtimeError(ctx, token, "integer overflow")
		}
		return &IntegerValue{Value: result}, nil
	case "-":
		result := l.Value - r.Value
		if (r.Value > 0 && result > l.Value) || (r.Value < 0 && result < l.Value) {
			return nil, runtimeError(ctx, token, "integer overflow")
		}
		return &IntegerValue{Value: result}, nil
	case "*":
		result := l.Value * r.Value
		if l.Value != 0 && r.Value != 0 && result/l.Value != r.Value {
			return nil, runtimeError(ctx, token, "integer overflow")
		}
		return &IntegerValue{Value: result}, nil
	case "/":
		if r.Value == 0 {
			return nil, runtimeError(ctx, token, "division by zero")
		}
		return &IntegerValue{Value: l.Value / r.Value}, nil
	case "%":
		if r.Value == 0 {
			return nil, runtimeError(ctx, token, "division by zero")
		}
		return &IntegerValue{Value: l.Value % r.Value}, nil
	case "<":
		return &BooleanValue{Value: l.Value < r.Value}, nil
	case ">":
		return &BooleanValue{Value: l.Value > r.Value}, nil
	case "<=":
		return &BooleanValue{Value: l.Value <= r.Value}, nil
	case ">=":
		return &BooleanValue{Value: l.Value >= r.Value}, nil
	case "==":
		return &BooleanValue{Value: l.Value == r.Value}, nil
	case "!=":
		return &BooleanValue{Value: l.Value != r.Value}, nil
	default:
		return nil, runtimeError(ctx, token, fmt.Sprintf("unknown operator: %s", operator))
	}
}

func (e *Evaluator) evaluateDecimalInfixExpression(ctx *ExecutionContext, token lexer.Token, l, r *DecimalValue) (Object, error) {
	operator := token.Source

	switch operator {
	case "+":
		return &DecimalValue{Value: l.Value + r.Value}, nil
	case "-":
		return &DecimalValue{Value: l.Value - r.Value}, nil
	case "*":
		return &DecimalValue{Value: l.Value * r.Value}, nil
	case "/":
		if r.Value == 0 {
			return nil, runtimeError(ctx, token, "division by zero")
		}
		return &DecimalValue{Value: l.Value / r.Value}, nil
	case "%":
		if r.Value == 0 {
			return nil, runtimeError(ctx, token, "division by zero")
		}
		return &DecimalValue{Value: math.Mod(l.Value, r.Value)}, nil
	case "<":
		return &BooleanValue{Value: l.Value < r.Value}, nil
	case ">":
		return &BooleanValue{Value: l.Value > r.Value}, nil
	case "<=":
		return &BooleanValue{Value: l.Value <= r.Value}, nil
	case ">=":
		return &BooleanValue{Value: l.Value >= r.Value}, nil
	case "==":
		return &BooleanValue{Value: l.Value == r.Value}, nil
	case "!=":
		return &BooleanValue{Value: l.Value != r.Value}, nil
	default:
		return nil, runtimeError(ctx, token, fmt.Sprintf("unknown operator: %s", operator))
	}
}

func (e *Evaluator) evaluateBooleanInfixExpression(operator string, l, r *BooleanValue) (Object, error) {
	switch operator {
	case "==":
		return &BooleanValue{Value: l.Value == r.Value}, nil
	case "!=":
		return &BooleanValue{Value: l.Value != r.Value}, nil
	default:
		return nil, fmt.Errorf("unknown operator: %s", operator)
	}
}

func (e *Evaluator) evaluateStringInfixExpression(operator string, l, r *StringValue) (Object, error) {

	switch operator {
	case "+":
		return &StringValue{Value: l.Value + r.Value}, nil
	case "==":
		return &BooleanValue{Value: l.Value == r.Value}, nil
	case "!=":
		return &BooleanValue{Value: l.Value != r.Value}, nil
	default:
		return nil, fmt.Errorf("unknown operator: %s", operator)
	}
}

func (e *Evaluator) evaluateIfExpression(ctx *ExecutionContext, ie *parser.IfExpression, scope *Scope) (Object, error) {
	condition, err := e.evaluateNode(ctx, ie.Condition, scope)
	if err != nil {
		return nil, err
	}

	if isTruthy(condition) {
		return e.evaluateNode(ctx, ie.Consequence, scope)
	}

	if ie.Alternative != nil {
		return e.evaluateNode(ctx, ie.Alternative, scope)
	}

	return Null, nil
}

func (e *Evaluator) evaluateIdentifier(ctx *ExecutionContext, ident *parser.Identifier, scope *Scope) (Object, error) {
	if val, ok := scope.Get(ident.Value); ok {
		return val, nil
	}

	if builtin, ok := e.functions[ident.Value]; ok {
		return builtin, nil
	}

	return nil, runtimeError(ctx, ident.Token, fmt.Sprintf("identifier not found: %s", ident.Value))
}

func (e *Evaluator) evaluateFunctionLiteral(fl *parser.FunctionLiteral, scope *Scope) (Object, error) {

	fv := &FunctionValue{Parameters: fl.Parameters, Body: fl.Body, Scope: scope}

	if fl.Identifier != nil {

		if _, ok := scope.GetLocal(fl.Identifier.Value); ok {
			return nil, fmt.Errorf("identifier already defined in local scope: %s", fl.Identifier.Value)
		}

		scope.SetLocal(fl.Identifier.Value, fv)
	}

	return fv, nil
}

func (e *Evaluator) evaluateCallExpression(ctx *ExecutionContext, ce *parser.CallExpression, scope *Scope) (Object, error) {
	function, err := e.evaluateNode(ctx, ce.Function, scope)
	if err != nil {
		return nil, err
	}

	args, err := e.evaluateExpressions(ctx, ce.Args, scope)
	if err != nil {
		return nil, err
	}

	return e.applyFunction(ctx, scope, function, args)
}

func (e *Evaluator) evaluateExpressions(ctx *ExecutionContext, exps []parser.Expression, scope *Scope) ([]Object, error) {
	var result []Object

	for _, exp := range exps {
		evaluated, err := e.evaluateNode(ctx, exp, scope)
		if err != nil {
			return nil, err
		}

		result = append(result, evaluated)
	}

	return result, nil
}

func (e *Evaluator) applyFunction(ctx *ExecutionContext, scope *Scope, fn Object, args []Object) (Object, error) {
	switch f := fn.(type) {
	case *FunctionValue:
		if len(args) != len(f.Parameters) {
			return nil, fmt.Errorf("wrong number of arguments: expected %d, got %d", len(f.Parameters), len(args))
		}

		if ctx.MaxDepth > 0 {
			ctx.depth++
			if ctx.depth > ctx.MaxDepth {
				return nil, fmt.Errorf("maximum call depth exceeded: %d", ctx.MaxDepth)
			}
			defer func() { ctx.depth-- }()
		}

		prevTM := ctx.templateMode
		ctx.templateMode = false
		extendedScope := e.extendFunctionScope(f, args)
		evaluated, err := e.evaluateNode(ctx, f.Body, extendedScope)
		ctx.templateMode = prevTM
		if err != nil {
			return nil, err
		}

		// Unwrap return values and discard break/continue signals that leaked
		// out of loops within the function body
		switch evaluated.(type) {
		case *BreakSignal, *ContinueSignal:
			return Null, nil
		}

		return unwrapReturnValue(evaluated), nil
	case *BuiltInFunction:
		return f.Fn(ctx, scope, args...)
	default:
		return nil, fmt.Errorf("not a function: %T", fn)
	}
}

func (e *Evaluator) extendFunctionScope(f *FunctionValue, args []Object) *Scope {
	extended := NewChildScope(f.Scope)

	for i, param := range f.Parameters {
		extended.SetLocal(param.Value, args[i])
	}

	return extended
}

func (e *Evaluator) evaluateArrayLiteral(ctx *ExecutionContext, al *parser.ArrayLiteral, scope *Scope) (Object, error) {
	if ctx.MaxArraySize > 0 && len(al.Elements) > ctx.MaxArraySize {
		return nil, fmt.Errorf("maximum array size exceeded: %d", ctx.MaxArraySize)
	}

	elements, err := e.evaluateExpressions(ctx, al.Elements, scope)
	if err != nil {
		return nil, err
	}

	return &ArrayValue{Elements: elements}, nil
}

func (e *Evaluator) evaluateIndexExpression(left, index Object) (Object, error) {

	if a, ok := left.(*ArrayValue); ok {
		if i, ok := index.(*IntegerValue); ok {
			return e.evaluateArrayIndexExpression(a, i)
		}
	}

	if h, ok := left.(*HashValue); ok {
		return e.evaluateHashIndexExpression(h, index)
	}

	return nil, fmt.Errorf("index operator not supported: %T", left)
}

func (e *Evaluator) evaluateArrayIndexExpression(array *ArrayValue, index *IntegerValue) (Object, error) {
	if index.Value < 0 || index.Value >= len(array.Elements) {
		return Null, nil
	}

	return array.Elements[index.Value], nil
}

func (e *Evaluator) evaluateHashIndexExpression(hash *HashValue, index Object) (Object, error) {
	if key, ok := index.(Hashable); ok {
		if pair, ok := hash.Pairs[key.HashKey()]; ok {
			return pair.Value, nil
		}
	}

	return Null, nil
}

func (e *Evaluator) evaluateHashLiteral(ctx *ExecutionContext, hl *parser.HashLiteral, scope *Scope) (Object, error) {
	hash := NewHashValue()

	for _, pair := range hl.Pairs {
		key, err := e.evaluateNode(ctx, pair.Key, scope)
		if err != nil {
			return nil, err
		}

		value, err := e.evaluateNode(ctx, pair.Value, scope)
		if err != nil {
			return nil, err
		}

		if err := hash.Set(key, value); err != nil {
			return nil, err
		}
	}

	return hash, nil
}

func (e *Evaluator) evaluatePropertyExpression(ctx *ExecutionContext, pe *parser.PropertyExpression, scope *Scope) (Object, Object, Object, error) {

	left, err := e.evaluateNode(ctx, pe.Left, scope)
	if err != nil {
		return nil, nil, nil, err
	}

	if _, ok := left.(*HashValue); !ok {
		return Null, Null, Null, nil
	}

	right := pe.Property

	for {
		switch r := right.(type) {
		case *parser.PropertyExpression:

			if ident, ok := r.Left.(*parser.Identifier); ok {
				leftValue, err := e.evaluateIndexExpression(left, &StringValue{Value: ident.Value})
				if err != nil {
					return nil, nil, nil, err
				}

				leftHash, ok := leftValue.(*HashValue)
				if !ok {
					return Null, Null, Null, nil
				}

				left = leftHash
				right = r.Property
				break
			}

			if indexIdent, ok := r.Left.(*parser.IndexExpression); ok {
				indexValue, err := e.evaluateNode(ctx, indexIdent.Index, scope)
				if err != nil {
					return nil, nil, nil, err
				}

				leftIdentifier, ok := indexIdent.Left.(*parser.Identifier)
				if !ok {
					return Null, Null, Null, nil
				}

				leftValue, err := e.evaluateIndexExpression(left, &StringValue{Value: leftIdentifier.Value})
				if err != nil {
					return nil, nil, nil, err
				}

				arrayObject, ok := leftValue.(*ArrayValue)
				if !ok {
					return Null, Null, Null, nil
				}

				left, err = e.evaluateIndexExpression(arrayObject, indexValue)
				if err != nil {
					return nil, nil, nil, err
				}

				right = r.Property
				break
			}

			return Null, Null, Null, nil
		case *parser.Identifier:
			idx := &StringValue{Value: r.Value}
			v, err := e.evaluateIndexExpression(left, &StringValue{Value: r.Value})
			if err != nil {
				return nil, nil, nil, err
			}

			return left, idx, v, nil
		default:
			return Null, Null, Null, nil
		}
	}
}

func unwrapReturnValue(obj Object) Object {
	if returnValue, ok := obj.(*ReturnValue); ok {
		return returnValue.Value
	}

	return obj
}

func isTruthy(obj Object) bool {
	switch o := obj.(type) {
	case *BooleanValue:
		return o.Value
	case *NullValue:
		return false
	default:
		return true
	}
}
