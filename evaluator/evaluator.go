package evaluator

import (
	"fmt"
	"io"
	"strings"

	"github.com/ironfang-ltd/go-script/parser"
)

type Node interface {
	Debug() string
}

type Evaluator struct {
	funcs  map[string]*BuiltInFunction
	output strings.Builder
	log    io.StringWriter
}

func New(log io.StringWriter) *Evaluator {
	e := &Evaluator{
		funcs:  make(map[string]*BuiltInFunction),
		output: strings.Builder{},
		log:    log,
	}

	e.RegisterFunction("log", func(args ...Object) (Object, error) {
		for _, arg := range args {
			_, _ = e.log.WriteString(arg.Debug())
			_, _ = e.log.WriteString("\n")
		}
		return Null, nil
	})

	e.RegisterFunction("print", func(args ...Object) (Object, error) {
		for _, arg := range args {
			e.output.WriteString(arg.Debug())
		}
		return Null, nil
	})

	return e
}

func (e *Evaluator) RegisterFunction(name string, fn Function) {
	e.funcs[name] = &BuiltInFunction{Fn: fn}
}

func (e *Evaluator) Evaluate(program *parser.Program, scope *Scope) (Object, error) {

	e.output.Reset()

	var result Object = Null

	for _, statement := range program.Statements {
		evalResult, err := e.evaluateNode(statement, scope)
		if err != nil {
			return nil, err
		}

		if _, ok := evalResult.(*ReturnValue); ok {
			return evalResult, nil
		}

		if evalResult != nil && evalResult.Type() != NullObject && evalResult.Type() != ReturnValueObject && evalResult.Type() != FunctionObject {
			e.output.WriteString(evalResult.Debug())
		}

		result = evalResult
	}

	return result, nil
}

func (e *Evaluator) GetOutput() string {
	return e.output.String()
}

func (e *Evaluator) evaluateNode(node Node, scope *Scope) (Object, error) {
	switch n := node.(type) {
	case *parser.PrintStatement:
		return e.evaluatePrintStatement(n)
	case *parser.BlockStatement:
		return e.evaluateBlockStatement(n, scope)
	case *parser.LetStatement:
		return e.evaluateLetStatement(n, scope)
	case *parser.ReturnStatement:
		return e.evaluateReturnStatement(n, scope)
	case *parser.ExpressionStatement:
		return e.evaluateNode(n.Expression, scope)
	case *parser.ForeachExpression:
		return e.evaluateForEach(n, scope)
	case *parser.IntegerLiteral:
		return &IntegerValue{Value: n.Value}, nil
	case *parser.BooleanLiteral:
		return &BooleanValue{Value: n.Value}, nil
	case *parser.StringLiteral:
		return &StringValue{Value: n.Value}, nil
	case *parser.PrefixExpression:
		right, err := e.evaluateNode(n.Right, scope)
		if err != nil {
			return nil, err
		}

		return e.evaluatePrefixExpression(n.Operator, right)
	case *parser.InfixExpression:
		left, err := e.evaluateNode(n.Left, scope)
		if err != nil {
			return nil, err
		}

		right, err := e.evaluateNode(n.Right, scope)
		if err != nil {
			return nil, err
		}

		return e.evaluateInfixExpression(n.Token.Source, left, right)
	case *parser.IfExpression:
		return e.evaluateIfExpression(n, scope)
	case *parser.Identifier:
		return e.evaluateIdentifier(n, scope)
	case *parser.FunctionLiteral:
		return e.evaluateFunctionLiteral(n, scope)
	case *parser.CallExpression:
		return e.evaluateCallExpression(n, scope)
	case *parser.ArrayLiteral:
		return e.evaluateArrayLiteral(n, scope)
	case *parser.IndexExpression:

		left, err := e.evaluateNode(n.Left, scope)
		if err != nil {
			return nil, err
		}

		index, err := e.evaluateNode(n.Index, scope)
		if err != nil {
			return nil, err
		}
		return e.evaluateIndexExpression(left, index)
	case *parser.HashLiteral:
		return e.evaluateHashLiteral(n, scope)
	case *parser.PropertyExpression:
		return e.evaluatePropertyExpression(n, scope)

	default:
		return nil, fmt.Errorf("unknown node type: %T", n)
	}
}

func (e *Evaluator) evaluatePrintStatement(print *parser.PrintStatement) (Object, error) {
	e.output.WriteString(print.Value)
	return Null, nil
}

func (e *Evaluator) evaluateForEach(foreach *parser.ForeachExpression, scope *Scope) (Object, error) {
	iterable, err := e.evaluateNode(foreach.Iterable, scope)
	if err != nil {
		return nil, err
	}

	switch i := iterable.(type) {
	case *ArrayValue:
		return e.evaluateArrayForEach(foreach, i, scope)
	case *HashValue:
		return e.evaluateHashForEach(foreach, i, scope)
	default:
		return Null, nil
	}
}

func (e *Evaluator) evaluateArrayForEach(foreach *parser.ForeachExpression, array *ArrayValue, scope *Scope) (Object, error) {
	for _, el := range array.Elements {
		extendedScope := NewChildScope(scope)
		extendedScope.Set(foreach.Variable.Value, el)

		_, err := e.evaluateBlockStatement(foreach.Body, extendedScope)
		if err != nil {
			return nil, err
		}
	}

	return Null, nil
}

func (e *Evaluator) evaluateHashForEach(foreach *parser.ForeachExpression, hash *HashValue, scope *Scope) (Object, error) {
	for _, pair := range hash.Pairs {
		extendedScope := NewChildScope(scope)
		extendedScope.Set(foreach.Variable.Value, pair.Value)

		_, err := e.evaluateBlockStatement(foreach.Body, extendedScope)
		if err != nil {
			return nil, err
		}
	}

	return Null, nil
}

func (e *Evaluator) evaluateBlockStatement(block *parser.BlockStatement, scope *Scope) (Object, error) {
	var result Object

	for _, statement := range block.Statements {
		evalResult, err := e.evaluateNode(statement, scope)
		if err != nil {
			return nil, err
		}

		// If the statement is a return statement, we should return the value
		// immediately.
		if _, ok := evalResult.(*ReturnValue); ok {
			return evalResult, nil
		}

		result = evalResult
	}

	return result, nil
}

func (e *Evaluator) evaluateLetStatement(let *parser.LetStatement, scope *Scope) (Object, error) {
	val, err := e.evaluateNode(let.Value, scope)
	if err != nil {
		return nil, err
	}

	scope.Set(let.Name.Value, val)

	return val, nil
}

func (e *Evaluator) evaluateReturnStatement(ret *parser.ReturnStatement, scope *Scope) (Object, error) {
	val, err := e.evaluateNode(ret.Value, scope)
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
		return &IntegerValue{Value: -r.Value}, nil
	default:
		return nil, fmt.Errorf("unknown operator: -%T", r)
	}
}

func (e *Evaluator) evaluateInfixExpression(operator string, left, right Object) (Object, error) {

	if i1, ok := left.(*IntegerValue); ok {
		if i2, ok := right.(*IntegerValue); ok {
			return e.evaluateIntegerInfixExpression(operator, i1, i2)
		}
	}

	if b1, ok := left.(*BooleanValue); ok {
		if b2, ok := right.(*BooleanValue); ok {
			return e.evaluateBooleanInfixExpression(operator, b1, b2)
		}
	}

	if s1, ok := left.(*StringValue); ok {
		if s2, ok := right.(*StringValue); ok {
			return e.evaluateStringInfixExpression(operator, s1, s2)
		}
	}

	if operator == "==" {
		return &BooleanValue{Value: left == right}, nil
	}

	if operator == "!=" {
		return &BooleanValue{Value: left != right}, nil
	}

	if left.Type() != right.Type() {
		return nil, fmt.Errorf("type mismatch: %T %s %T", left, operator, right)
	}

	return nil, fmt.Errorf("unknown operator: %T %s %T", left, operator, right)
}

func (e *Evaluator) evaluateIntegerInfixExpression(operator string, l, r *IntegerValue) (Object, error) {

	switch operator {
	case "+":
		return &IntegerValue{Value: l.Value + r.Value}, nil
	case "-":
		return &IntegerValue{Value: l.Value - r.Value}, nil
	case "*":
		return &IntegerValue{Value: l.Value * r.Value}, nil
	case "/":
		return &IntegerValue{Value: l.Value / r.Value}, nil
	case "<":
		return &BooleanValue{Value: l.Value < r.Value}, nil
	case ">":
		return &BooleanValue{Value: l.Value > r.Value}, nil
	case "==":
		return &BooleanValue{Value: l.Value == r.Value}, nil
	case "!=":
		return &BooleanValue{Value: l.Value != r.Value}, nil
	default:
		return nil, fmt.Errorf("unknown operator: %s", operator)
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

func (e *Evaluator) evaluateIfExpression(ie *parser.IfExpression, scope *Scope) (Object, error) {
	condition, err := e.evaluateNode(ie.Condition, scope)
	if err != nil {
		return nil, err
	}

	if isTruthy(condition) {
		return e.evaluateNode(ie.Consequence, scope)
	}

	if ie.Alternative != nil {
		return e.evaluateNode(ie.Alternative, scope)
	}

	return Null, nil
}

func (e *Evaluator) evaluateIdentifier(ident *parser.Identifier, scope *Scope) (Object, error) {
	if val, ok := scope.Get(ident.Value); ok {
		return val, nil
	}

	if builtin, ok := e.funcs[ident.Value]; ok {
		return builtin, nil
	}

	return nil, fmt.Errorf("identifier not found: %s", ident.Value)
}

func (e *Evaluator) evaluateFunctionLiteral(fl *parser.FunctionLiteral, scope *Scope) (Object, error) {

	fv := &FunctionValue{Parameters: fl.Parameters, Body: fl.Body, Scope: scope}

	if fl.Identifier != nil {

		if _, ok := scope.GetLocal(fl.Identifier.Value); ok {
			return nil, fmt.Errorf("identifier already defined in local scope: %s", fl.Identifier.Value)
		}

		scope.Set(fl.Identifier.Value, fv)
	}

	return fv, nil
}

func (e *Evaluator) evaluateCallExpression(ce *parser.CallExpression, scope *Scope) (Object, error) {
	function, err := e.evaluateNode(ce.Function, scope)
	if err != nil {
		return nil, err
	}

	args, err := e.evaluateExpressions(ce.Args, scope)
	if err != nil {
		return nil, err
	}

	return e.applyFunction(function, args)
}

func (e *Evaluator) evaluateExpressions(exps []parser.Expression, scope *Scope) ([]Object, error) {
	var result []Object

	for _, exp := range exps {
		evaluated, err := e.evaluateNode(exp, scope)
		if err != nil {
			return nil, err
		}

		result = append(result, evaluated)
	}

	return result, nil
}

func (e *Evaluator) applyFunction(fn Object, args []Object) (Object, error) {
	switch f := fn.(type) {
	case *FunctionValue:
		extendedScope := e.extendFunctionScope(f, args)
		evaluated, err := e.evaluateNode(f.Body, extendedScope)
		if err != nil {
			return nil, err
		}

		return unwrapReturnValue(evaluated), nil
	case *BuiltInFunction:
		return f.Fn(args...)
	default:
		return nil, fmt.Errorf("not a function: %T", fn)
	}
}

func (e *Evaluator) extendFunctionScope(f *FunctionValue, args []Object) *Scope {
	extended := NewChildScope(f.Scope)

	for i, param := range f.Parameters {
		extended.Set(param.Value, args[i])
	}

	return extended
}

func (e *Evaluator) evaluateArrayLiteral(al *parser.ArrayLiteral, scope *Scope) (Object, error) {
	elements, err := e.evaluateExpressions(al.Elements, scope)
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

func (e *Evaluator) evaluateHashLiteral(hl *parser.HashLiteral, scope *Scope) (Object, error) {
	pairs := make(map[HashKey]HashPair)

	for keyNode, valueNode := range hl.Pairs {
		key, err := e.evaluateNode(keyNode, scope)
		if err != nil {
			return nil, err
		}

		hashable, ok := key.(Hashable)
		if !ok {
			return nil, fmt.Errorf("unusable as hash key: %T", key)
		}

		value, err := e.evaluateNode(valueNode, scope)
		if err != nil {
			return nil, err
		}

		pairs[hashable.HashKey()] = HashPair{Key: key, Value: value}
	}

	return &HashValue{Pairs: pairs}, nil
}

func (e *Evaluator) evaluatePropertyExpression(pe *parser.PropertyExpression, scope *Scope) (Object, error) {

	left, err := e.evaluateNode(pe.Left, scope)
	if err != nil {
		return nil, err
	}

	if _, ok := left.(*HashValue); !ok {
		return Null, nil
	}

	right := pe.Property

	for {
		switch r := right.(type) {
		case *parser.PropertyExpression:

			if ident, ok := r.Left.(*parser.Identifier); ok {
				leftValue, err := e.evaluateIndexExpression(left, &StringValue{Value: ident.Value})
				if err != nil {
					return nil, err
				}

				leftHash, ok := leftValue.(*HashValue)
				if !ok {
					return Null, nil
				}

				left = leftHash
				right = r.Property
				break
			}

			if indexIdent, ok := r.Left.(*parser.IndexExpression); ok {
				indexValue, err := e.evaluateNode(indexIdent.Index, scope)
				if err != nil {
					return nil, err
				}

				leftIdentifier, ok := indexIdent.Left.(*parser.Identifier)
				if !ok {
					return Null, nil
				}

				leftValue, err := e.evaluateIndexExpression(left, &StringValue{Value: leftIdentifier.Value})
				if err != nil {
					return nil, err
				}

				arrayObject, ok := leftValue.(*ArrayValue)
				if !ok {
					return Null, nil
				}

				left, err = e.evaluateIndexExpression(arrayObject, indexValue)
				if err != nil {
					return nil, err
				}

				right = r.Property
				break
			}

			return Null, nil
		case *parser.Identifier:
			return e.evaluateIndexExpression(left, &StringValue{Value: r.Value})
		default:
			return Null, nil
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
