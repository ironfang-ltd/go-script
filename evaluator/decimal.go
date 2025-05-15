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
	return strings.TrimRight(fmt.Sprintf("%f", i.Value), "0")
}

func (i *DecimalValue) Type() ObjectType {
	return DecimalObject
}
