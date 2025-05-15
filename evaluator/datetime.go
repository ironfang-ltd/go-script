package evaluator

import "time"

type DateTimeValue struct {
	Value time.Time
}

func NewDateTimeValue(t time.Time) *DateTimeValue {
	return &DateTimeValue{Value: t}
}

func (v *DateTimeValue) Debug() string {
	return v.Value.Format(time.RFC3339)
}

func (v *DateTimeValue) Type() ObjectType {
	return DateTimeObject
}
