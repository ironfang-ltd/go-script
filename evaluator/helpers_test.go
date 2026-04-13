package evaluator

import (
	"testing"
	"time"
)

func TestToObject_Nil(t *testing.T) {
	obj, err := ToObject(nil)
	if err != nil {
		t.Fatal(err)
	}
	if obj.Type() != NullObject {
		t.Fatalf("expected NullObject, got %s", obj.Type())
	}
}

func TestToObject_String(t *testing.T) {
	obj, err := ToObject("hello")
	if err != nil {
		t.Fatal(err)
	}
	sv, ok := obj.(*StringValue)
	if !ok {
		t.Fatalf("expected *StringValue, got %T", obj)
	}
	if sv.Value != "hello" {
		t.Fatalf("expected hello, got %s", sv.Value)
	}
}

func TestToObject_Bool(t *testing.T) {
	obj, err := ToObject(true)
	if err != nil {
		t.Fatal(err)
	}
	bv, ok := obj.(*BooleanValue)
	if !ok {
		t.Fatalf("expected *BooleanValue, got %T", obj)
	}
	if !bv.Value {
		t.Fatal("expected true")
	}
}

func TestToObject_IntVariants(t *testing.T) {
	tests := []struct {
		name string
		val  any
		want int
	}{
		{"int", int(42), 42},
		{"int8", int8(8), 8},
		{"int16", int16(16), 16},
		{"int32", int32(32), 32},
		{"int64", int64(64), 64},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj, err := ToObject(tt.val)
			if err != nil {
				t.Fatal(err)
			}
			iv, ok := obj.(*IntegerValue)
			if !ok {
				t.Fatalf("expected *IntegerValue, got %T", obj)
			}
			if iv.Value != tt.want {
				t.Fatalf("expected %d, got %d", tt.want, iv.Value)
			}
		})
	}
}

func TestToObject_FloatVariants(t *testing.T) {
	tests := []struct {
		name string
		val  any
		want float64
	}{
		{"float32", float32(3.14), float64(float32(3.14))},
		{"float64", float64(2.718), 2.718},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj, err := ToObject(tt.val)
			if err != nil {
				t.Fatal(err)
			}
			dv, ok := obj.(*DecimalValue)
			if !ok {
				t.Fatalf("expected *DecimalValue, got %T", obj)
			}
			if dv.Value != tt.want {
				t.Fatalf("expected %f, got %f", tt.want, dv.Value)
			}
		})
	}
}

func TestToObject_Time(t *testing.T) {
	now := time.Now()
	obj, err := ToObject(now)
	if err != nil {
		t.Fatal(err)
	}
	dt, ok := obj.(*DateTimeValue)
	if !ok {
		t.Fatalf("expected *DateTimeValue, got %T", obj)
	}
	if !dt.Value.Equal(now) {
		t.Fatalf("expected %v, got %v", now, dt.Value)
	}
}

func TestToObject_TimePointerNil(t *testing.T) {
	var tp *time.Time
	obj, err := ToObject(tp)
	if err != nil {
		t.Fatal(err)
	}
	if obj.Type() != NullObject {
		t.Fatalf("expected NullObject, got %s", obj.Type())
	}
}

func TestToObject_TimePointer(t *testing.T) {
	now := time.Now()
	obj, err := ToObject(&now)
	if err != nil {
		t.Fatal(err)
	}
	dt, ok := obj.(*DateTimeValue)
	if !ok {
		t.Fatalf("expected *DateTimeValue, got %T", obj)
	}
	if !dt.Value.Equal(now) {
		t.Fatalf("expected %v, got %v", now, dt.Value)
	}
}

func TestToObject_Slice(t *testing.T) {
	obj, err := ToObject([]any{"a", 1, true})
	if err != nil {
		t.Fatal(err)
	}
	arr, ok := obj.(*ArrayValue)
	if !ok {
		t.Fatalf("expected *ArrayValue, got %T", obj)
	}
	if len(arr.Elements) != 3 {
		t.Fatalf("expected 3 elements, got %d", len(arr.Elements))
	}
	if arr.Elements[0].(*StringValue).Value != "a" {
		t.Fatal("expected first element to be 'a'")
	}
	if arr.Elements[1].(*IntegerValue).Value != 1 {
		t.Fatal("expected second element to be 1")
	}
	if !arr.Elements[2].(*BooleanValue).Value {
		t.Fatal("expected third element to be true")
	}
}

func TestToObject_Map(t *testing.T) {
	obj, err := ToObject(map[string]any{"name": "Alice", "age": 30})
	if err != nil {
		t.Fatal(err)
	}
	hash, ok := obj.(*HashValue)
	if !ok {
		t.Fatalf("expected *HashValue, got %T", obj)
	}
	nameVal, ok := hash.GetValue(&StringValue{Value: "name"})
	if !ok {
		t.Fatal("expected 'name' key to exist")
	}
	if nameVal.(*StringValue).Value != "Alice" {
		t.Fatalf("expected Alice, got %s", nameVal.Debug())
	}
	ageVal, ok := hash.GetValue(&StringValue{Value: "age"})
	if !ok {
		t.Fatal("expected 'age' key to exist")
	}
	if ageVal.(*IntegerValue).Value != 30 {
		t.Fatalf("expected 30, got %s", ageVal.Debug())
	}
}

func TestToObject_ObjectPassthrough(t *testing.T) {
	original := &StringValue{Value: "pass"}
	obj, err := ToObject(original)
	if err != nil {
		t.Fatal(err)
	}
	if obj != original {
		t.Fatal("expected same object reference")
	}
}

func TestToObject_TypedSliceOfMaps(t *testing.T) {
	input := []map[string]any{
		{"name": "Alice", "age": 30},
		{"name": "Bob", "age": 25},
	}
	obj, err := ToObject(input)
	if err != nil {
		t.Fatal(err)
	}
	arr, ok := obj.(*ArrayValue)
	if !ok {
		t.Fatalf("expected *ArrayValue, got %T", obj)
	}
	if len(arr.Elements) != 2 {
		t.Fatalf("expected 2 elements, got %d", len(arr.Elements))
	}
	hash, ok := arr.Elements[0].(*HashValue)
	if !ok {
		t.Fatalf("expected *HashValue, got %T", arr.Elements[0])
	}
	name, _ := hash.GetValue(&StringValue{Value: "name"})
	if name.(*StringValue).Value != "Alice" {
		t.Fatalf("expected 'Alice', got %q", name.(*StringValue).Value)
	}
}

func TestToObject_TypedSliceOfStrings(t *testing.T) {
	input := []string{"hello", "world"}
	obj, err := ToObject(input)
	if err != nil {
		t.Fatal(err)
	}
	arr, ok := obj.(*ArrayValue)
	if !ok {
		t.Fatalf("expected *ArrayValue, got %T", obj)
	}
	if len(arr.Elements) != 2 {
		t.Fatalf("expected 2 elements, got %d", len(arr.Elements))
	}
	if arr.Elements[0].(*StringValue).Value != "hello" {
		t.Fatalf("expected 'hello', got %q", arr.Elements[0].(*StringValue).Value)
	}
}

func TestToObject_TypedSliceOfInts(t *testing.T) {
	input := []int{1, 2, 3}
	obj, err := ToObject(input)
	if err != nil {
		t.Fatal(err)
	}
	arr, ok := obj.(*ArrayValue)
	if !ok {
		t.Fatalf("expected *ArrayValue, got %T", obj)
	}
	if len(arr.Elements) != 3 {
		t.Fatalf("expected 3 elements, got %d", len(arr.Elements))
	}
	if arr.Elements[0].(*IntegerValue).Value != 1 {
		t.Fatalf("expected 1, got %d", arr.Elements[0].(*IntegerValue).Value)
	}
}

func TestToObject_TypedMapStringString(t *testing.T) {
	input := map[string]string{"foo": "bar", "baz": "qux"}
	obj, err := ToObject(input)
	if err != nil {
		t.Fatal(err)
	}
	hash, ok := obj.(*HashValue)
	if !ok {
		t.Fatalf("expected *HashValue, got %T", obj)
	}
	val, _ := hash.GetValue(&StringValue{Value: "foo"})
	if val.(*StringValue).Value != "bar" {
		t.Fatalf("expected 'bar', got %q", val.(*StringValue).Value)
	}
}

func TestToObject_TypedSliceNestedError(t *testing.T) {
	input := []struct{ X int }{{1}}
	_, err := ToObject(input)
	if err == nil {
		t.Fatal("expected error for unsupported nested type in typed slice")
	}
}

func TestToObject_UnsupportedType(t *testing.T) {
	_, err := ToObject(struct{}{})
	if err == nil {
		t.Fatal("expected error for unsupported type")
	}
}

func TestToObject_SliceNestedError(t *testing.T) {
	_, err := ToObject([]any{struct{}{}})
	if err == nil {
		t.Fatal("expected error for unsupported nested type")
	}
}

func TestToObject_MapNestedError(t *testing.T) {
	_, err := ToObject(map[string]any{"bad": struct{}{}})
	if err == nil {
		t.Fatal("expected error for unsupported nested type")
	}
}

func TestRunScript_Simple(t *testing.T) {
	result, err := RunScript(`return 1 + 2;`)
	if err != nil {
		t.Fatal(err)
	}
	rv, ok := result.(*ReturnValue)
	if !ok {
		t.Fatalf("expected *ReturnValue, got %T", result)
	}
	iv, ok := rv.Value.(*IntegerValue)
	if !ok {
		t.Fatalf("expected *IntegerValue, got %T", rv.Value)
	}
	if iv.Value != 3 {
		t.Fatalf("expected 3, got %d", iv.Value)
	}
}

func TestRunScript_WithVars(t *testing.T) {
	result, err := RunScript(`return name + " is " + toString(age);`, Vars{
		"name": "Alice",
		"age":  30,
	})
	if err != nil {
		t.Fatal(err)
	}
	rv := result.(*ReturnValue)
	sv := rv.Value.(*StringValue)
	if sv.Value != "Alice is 30" {
		t.Fatalf("expected 'Alice is 30', got %q", sv.Value)
	}
}

func TestRunScript_MultipleVarsMerge(t *testing.T) {
	result, err := RunScript(`return a + b;`,
		Vars{"a": 10, "b": 20},
		Vars{"b": 30},
	)
	if err != nil {
		t.Fatal(err)
	}
	rv := result.(*ReturnValue)
	iv := rv.Value.(*IntegerValue)
	if iv.Value != 40 {
		t.Fatalf("expected 40, got %d", iv.Value)
	}
}

func TestRunScript_ParseError(t *testing.T) {
	_, err := RunScript(`return +;`)
	if err == nil {
		t.Fatal("expected parse error")
	}
}

func TestRunScript_RuntimeError(t *testing.T) {
	_, err := RunScript(`return x;`)
	if err == nil {
		t.Fatal("expected runtime error")
	}
}

func TestRunScript_NestedVars(t *testing.T) {
	result, err := RunScript(`return data.name;`, Vars{
		"data": map[string]any{"name": "Bob"},
	})
	if err != nil {
		t.Fatal(err)
	}
	rv := result.(*ReturnValue)
	sv := rv.Value.(*StringValue)
	if sv.Value != "Bob" {
		t.Fatalf("expected 'Bob', got %q", sv.Value)
	}
}

func TestRunTemplate_PlainText(t *testing.T) {
	output, err := RunTemplate(`Hello, World!`)
	if err != nil {
		t.Fatal(err)
	}
	if output != "Hello, World!" {
		t.Fatalf("expected 'Hello, World!', got %q", output)
	}
}

func TestRunTemplate_WithVars(t *testing.T) {
	output, err := RunTemplate(`Hello, {% name %}!`, Vars{"name": "Alice"})
	if err != nil {
		t.Fatal(err)
	}
	if output != "Hello, Alice!" {
		t.Fatalf("expected 'Hello, Alice!', got %q", output)
	}
}

func TestRunTemplate_ParseError(t *testing.T) {
	_, err := RunTemplate(`Hello {% +; %}`)
	if err == nil {
		t.Fatal("expected parse error")
	}
}

func TestEvaluator_RunScript_CustomFunction(t *testing.T) {
	eval := New()
	eval.RegisterFunction("double", func(ctx *ExecutionContext, scope *Scope, args ...Object) (Object, error) {
		iv := args[0].(*IntegerValue)
		return &IntegerValue{Value: iv.Value * 2}, nil
	})

	result, err := eval.RunScript(`return double(21);`)
	if err != nil {
		t.Fatal(err)
	}
	rv := result.(*ReturnValue)
	iv := rv.Value.(*IntegerValue)
	if iv.Value != 42 {
		t.Fatalf("expected 42, got %d", iv.Value)
	}
}

func TestEvaluator_RunTemplate_CustomFunction(t *testing.T) {
	eval := New()
	eval.RegisterFunction("shout", func(ctx *ExecutionContext, scope *Scope, args ...Object) (Object, error) {
		sv := args[0].(*StringValue)
		return &StringValue{Value: sv.Value + "!"}, nil
	})

	output, err := eval.RunTemplate(`Say {% shout("hello") %}`, Vars{"x": 1})
	if err != nil {
		t.Fatal(err)
	}
	if output != "Say hello!" {
		t.Fatalf("expected 'Say hello!', got %q", output)
	}
}

func TestPackageRunScript(t *testing.T) {
	result, err := RunScript(`return 42;`)
	if err != nil {
		t.Fatal(err)
	}
	rv := result.(*ReturnValue)
	if rv.Value.(*IntegerValue).Value != 42 {
		t.Fatal("expected 42")
	}
}

func TestPackageRunTemplate(t *testing.T) {
	output, err := RunTemplate(`Count: {% 1 + 2 %}`)
	if err != nil {
		t.Fatal(err)
	}
	if output != "Count: 3" {
		t.Fatalf("expected 'Count: 3', got %q", output)
	}
}

// --- Template output cleanliness tests ---

func TestRunTemplate_LetNoLeak(t *testing.T) {
	output, err := RunTemplate("start\n{% let x = 42; %}\nend")
	if err != nil {
		t.Fatal(err)
	}
	if output != "start\nend" {
		t.Fatalf("expected 'start\\nend', got %q", output)
	}
}

func TestRunTemplate_AssignmentNoLeak(t *testing.T) {
	output, err := RunTemplate("{% let x = 0; %}\n{% x = 5; %}\nval={% x %}")
	if err != nil {
		t.Fatal(err)
	}
	if output != "\nval=5" {
		t.Fatalf("expected '\\nval=5', got %q", output)
	}
}

func TestRunTemplate_CompoundAssignmentNoLeak(t *testing.T) {
	output, err := RunTemplate("{% let x = 10; %}\n{% x += 5; %}\nval={% x %}")
	if err != nil {
		t.Fatal(err)
	}
	if output != "\nval=15" {
		t.Fatalf("expected '\\nval=15', got %q", output)
	}
}

func TestRunTemplate_MultipleLetNoLeak(t *testing.T) {
	output, err := RunTemplate("A\n{% let x = 1; %}\n{% let y = 2; %}\nB")
	if err != nil {
		t.Fatal(err)
	}
	if output != "A\nB" {
		t.Fatalf("expected 'A\\nB', got %q", output)
	}
}

func TestRunTemplate_LetIfValueNoLeak(t *testing.T) {
	output, err := RunTemplate("{% let x = if (true) { 5; }; %}\nval={% x %}")
	if err != nil {
		t.Fatal(err)
	}
	if output != "\nval=5" {
		t.Fatalf("expected '\\nval=5', got %q", output)
	}
}

func TestRunTemplate_ExprInsideForeach(t *testing.T) {
	output, err := RunTemplate(
		`{% foreach (items as item) { %}[{% item %}]{% } %}`,
		Vars{"items": []any{"a", "b", "c"}})
	if err != nil {
		t.Fatal(err)
	}
	if output != "[a][b][c]" {
		t.Fatalf("expected '[a][b][c]', got %q", output)
	}
}

func TestRunTemplate_ExprInsideIf(t *testing.T) {
	output, err := RunTemplate(`{% if (true) { %}{% "yes" %}{% } %}`)
	if err != nil {
		t.Fatal(err)
	}
	if output != "yes" {
		t.Fatalf("expected 'yes', got %q", output)
	}
}

func TestRunTemplate_ForeachWithTemplateText(t *testing.T) {
	output, err := RunTemplate(
		"<ul>\n{% foreach (items as item) { %}\n<li>{% item %}</li>\n{% } %}\n</ul>",
		Vars{"items": []any{"Apple", "Banana"}})
	if err != nil {
		t.Fatal(err)
	}
	expected := "<ul>\n<li>Apple</li>\n<li>Banana</li>\n</ul>"
	if output != expected {
		t.Fatalf("expected %q, got %q", expected, output)
	}
}

func TestRunTemplate_IfBlockClean(t *testing.T) {
	output, err := RunTemplate(
		"<div>\n{% if (show) { %}\n<p>hello</p>\n{% } %}\n</div>",
		Vars{"show": true})
	if err != nil {
		t.Fatal(err)
	}
	expected := "<div>\n<p>hello</p>\n</div>"
	if output != expected {
		t.Fatalf("expected %q, got %q", expected, output)
	}
}

func TestRunTemplate_FunctionDefNoLeak(t *testing.T) {
	output, err := RunTemplate("{% fn double(n) { return n * 2; } %}\nval={% double(21) %}")
	if err != nil {
		t.Fatal(err)
	}
	if output != "\nval=42" {
		t.Fatalf("expected '\\nval=42', got %q", output)
	}
}

func TestRunTemplate_DotAccessInsideForeach(t *testing.T) {
	output, err := RunTemplate(
		`{% foreach (items as item) { %}{% item.name %}:{% item.price %} {% } %}`,
		Vars{"items": []any{
			map[string]any{"name": "Apple", "price": 1.5},
			map[string]any{"name": "Banana", "price": 0.75},
		}})
	if err != nil {
		t.Fatal(err)
	}
	expected := "Apple:1.5 Banana:0.75 "
	if output != expected {
		t.Fatalf("expected %q, got %q", expected, output)
	}
}

func TestRunTemplate_PrintStillWorks(t *testing.T) {
	output, err := RunTemplate(
		`{% foreach (items as item) { %}[{% print(item) %}]{% } %}`,
		Vars{"items": []any{"x", "y"}})
	if err != nil {
		t.Fatal(err)
	}
	if output != "[x][y]" {
		t.Fatalf("expected '[x][y]', got %q", output)
	}
}

func TestRunTemplate_LetInsideIfNoLeak(t *testing.T) {
	output, err := RunTemplate(
		"{% if (true) { %}\n{% let x = 99; %}\nval={% x %}\n{% } %}")
	if err != nil {
		t.Fatal(err)
	}
	if output != "\nval=99" {
		t.Fatalf("expected '\\nval=99', got %q", output)
	}
}
