package evaluator

import "testing"

func TestDecimalValue(t *testing.T) {

	v := NewDecimalValue(2.56)

	if v.Type() != DecimalObject {
		t.Fatalf("expected DecimalObject but got %s", v.Type())
	}

	if v.Value != 2.56 {
		t.Fatalf("expected 2.56 but got %f", v.Value)
	}

	if v.Debug() != "2.56" {
		t.Fatalf("expected \"2.56\" but got %s", v.Debug())
	}
}
