package evaluator

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
