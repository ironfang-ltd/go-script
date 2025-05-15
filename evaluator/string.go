package evaluator

import "hash/fnv"

type StringValue struct {
	Value string
}

func NewStringValue(s string) *StringValue {
	return &StringValue{s}
}

func (s *StringValue) Debug() string {
	return s.Value
}

func (s *StringValue) Type() ObjectType {
	return StringObject
}

func (s *StringValue) HashKey() HashKey {

	h := fnv.New64a()
	_, _ = h.Write([]byte(s.Value))

	return HashKey{Type: s.Type(), Value: h.Sum64()}
}
