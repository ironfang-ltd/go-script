package evaluator

import (
	"fmt"
	"hash/fnv"

	"github.com/ironfang-ltd/ironscript/parser"
)

var (
	True  = &BooleanValue{Value: true}
	False = &BooleanValue{Value: false}
	Null  = &NullValue{}
)

type ObjectType string

const (
	NullObject            ObjectType = "NULL"
	ReturnValueObject     ObjectType = "RETURN_VALUE"
	BooleanObject         ObjectType = "BOOLEAN"
	IntegerObject         ObjectType = "INTEGER"
	StringObject          ObjectType = "STRING"
	FunctionObject        ObjectType = "FUNCTION"
	ArrayObject           ObjectType = "ARRAY"
	HashObject            ObjectType = "HASH"
	BuiltInFunctionObject ObjectType = "BUILTIN_FUNCTION"
)

type Object interface {
	Type() ObjectType
	Debug() string
}

type NullValue struct{}

func (n *NullValue) Debug() string {
	return "null"
}

func (n *NullValue) Type() ObjectType {
	return NullObject
}

type ReturnValue struct {
	Value Object
}

func (r *ReturnValue) Debug() string {
	return r.Value.Debug()
}

func (r *ReturnValue) Type() ObjectType {
	return ReturnValueObject
}

type BooleanValue struct {
	Value bool
}

func (b *BooleanValue) Debug() string {
	if b.Value {
		return "true"
	}

	return "false"
}

func (b *BooleanValue) Type() ObjectType {
	return BooleanObject
}

type IntegerValue struct {
	Value int
}

func (i *IntegerValue) Debug() string {
	return fmt.Sprintf("%d", i.Value)
}

func (i *IntegerValue) Type() ObjectType {
	return IntegerObject
}

type StringValue struct {
	Value string
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

type FunctionValue struct {
	Parameters []*parser.Identifier
	Body       *parser.BlockStatement
	Scope      *Scope
}

func (f *FunctionValue) Debug() string {
	return "Function"
}

func (f *FunctionValue) Type() ObjectType {
	return FunctionObject
}

type ArrayValue struct {
	Elements []Object
}

func (a *ArrayValue) Debug() string {
	return "Array"
}

func (a *ArrayValue) Type() ObjectType {
	return ArrayObject
}

type Hashable interface {
	HashKey() HashKey
}

type HashKey struct {
	Type  ObjectType
	Value uint64
}

type HashPair struct {
	Key   Object
	Value Object
}

type HashValue struct {
	Pairs map[HashKey]HashPair
}

func NewHashValue() *HashValue {
	return &HashValue{Pairs: make(map[HashKey]HashPair)}
}

func (h *HashValue) Debug() string {
	return "Hash"
}

func (h *HashValue) Type() ObjectType {
	return HashObject
}

func (h *HashValue) Set(key Object, value Object) {
	hashKey := key.(Hashable).HashKey()
	h.Pairs[hashKey] = HashPair{Key: key, Value: value}
}
