package evaluator

import (
	"fmt"
	"strings"
)

type DecimalValue struct {
	Value float64
}

func NewDecimalValue(val float64) *DecimalValue {
	return &DecimalValue{val}
}

func (i *DecimalValue) Debug() string {
	s := fmt.Sprintf("%f", i.Value)
	// Trim trailing zeros but keep at least one digit after the decimal point
	if strings.Contains(s, ".") {
		s = strings.TrimRight(s, "0")
		if s[len(s)-1] == '.' {
			s += "0"
		}
	}
	return s
}

func (i *DecimalValue) Type() ObjectType {
	return DecimalObject
}
