package evaluator

import (
	"fmt"
	"reflect"
	"time"

	"github.com/ironfang-ltd/go-script/lexer"
	"github.com/ironfang-ltd/go-script/parser"
)

// Vars is a convenience type for passing variables to scripts and templates.
type Vars map[string]any

// ToObject converts a Go value to a script Object.
func ToObject(v any) (Object, error) {
	if v == nil {
		return Null, nil
	}

	if obj, ok := v.(Object); ok {
		return obj, nil
	}

	switch val := v.(type) {
	case string:
		return &StringValue{Value: val}, nil
	case bool:
		return &BooleanValue{Value: val}, nil
	case int:
		return &IntegerValue{Value: val}, nil
	case int8:
		return &IntegerValue{Value: int(val)}, nil
	case int16:
		return &IntegerValue{Value: int(val)}, nil
	case int32:
		return &IntegerValue{Value: int(val)}, nil
	case int64:
		return &IntegerValue{Value: int(val)}, nil
	case float32:
		return &DecimalValue{Value: float64(val)}, nil
	case float64:
		return &DecimalValue{Value: val}, nil
	case time.Time:
		return &DateTimeValue{Value: val}, nil
	case *time.Time:
		if val == nil {
			return Null, nil
		}
		return &DateTimeValue{Value: *val}, nil
	case []any:
		elements := make([]Object, len(val))
		for i, elem := range val {
			obj, err := ToObject(elem)
			if err != nil {
				return nil, fmt.Errorf("element [%d]: %w", i, err)
			}
			elements[i] = obj
		}
		return &ArrayValue{Elements: elements}, nil
	case map[string]any:
		hash := NewHashValue()
		for k, v := range val {
			obj, err := ToObject(v)
			if err != nil {
				return nil, fmt.Errorf("key %q: %w", k, err)
			}
			if err := hash.Set(&StringValue{Value: k}, obj); err != nil {
				return nil, err
			}
		}
		return hash, nil
	default:
		rv := reflect.ValueOf(v)
		if rv.Kind() == reflect.Slice || rv.Kind() == reflect.Array {
			elements := make([]Object, rv.Len())
			for i := range rv.Len() {
				obj, err := ToObject(rv.Index(i).Interface())
				if err != nil {
					return nil, fmt.Errorf("element [%d]: %w", i, err)
				}
				elements[i] = obj
			}
			return &ArrayValue{Elements: elements}, nil
		}
		if rv.Kind() == reflect.Map && rv.Type().Key().Kind() == reflect.String {
			hash := NewHashValue()
			iter := rv.MapRange()
			for iter.Next() {
				k := iter.Key().String()
				obj, err := ToObject(iter.Value().Interface())
				if err != nil {
					return nil, fmt.Errorf("key %q: %w", k, err)
				}
				if err := hash.Set(&StringValue{Value: k}, obj); err != nil {
					return nil, err
				}
			}
			return hash, nil
		}
		return nil, fmt.Errorf("unsupported type: %T", v)
	}
}

func applyVars(scope *Scope, vars []Vars) error {
	for _, m := range vars {
		for k, v := range m {
			obj, err := ToObject(v)
			if err != nil {
				return fmt.Errorf("variable %q: %w", k, err)
			}
			scope.SetLocal(k, obj)
		}
	}
	return nil
}

// RunScript parses and evaluates source as a script, returning the result Object.
func (e *Evaluator) RunScript(source string, vars ...Vars) (Object, error) {
	l := lexer.NewScript(source)
	p := parser.New(l)
	program, err := p.Parse()
	if err != nil {
		return nil, err
	}

	ctx := NewExecutionContext(program)
	ctx.Source = source

	if err := applyVars(ctx.RootScope, vars); err != nil {
		return nil, err
	}

	return e.Evaluate(ctx)
}

// RunTemplate parses and evaluates source as a template, returning the output string.
func (e *Evaluator) RunTemplate(source string, vars ...Vars) (string, error) {
	l := lexer.NewTemplate(source)
	p := parser.New(l)
	program, err := p.Parse()
	if err != nil {
		return "", err
	}

	ctx := NewExecutionContext(program)
	ctx.Source = source

	if err := applyVars(ctx.RootScope, vars); err != nil {
		return "", err
	}

	return e.EvaluateString(ctx)
}

// RunScript creates a fresh evaluator and evaluates source as a script.
func RunScript(source string, vars ...Vars) (Object, error) {
	return New().RunScript(source, vars...)
}

// RunTemplate creates a fresh evaluator and evaluates source as a template.
func RunTemplate(source string, vars ...Vars) (string, error) {
	return New().RunTemplate(source, vars...)
}
