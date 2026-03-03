package evaluator

import "fmt"

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
	order []HashKey
}

func NewHashValue() *HashValue {
	return &HashValue{
		Pairs: make(map[HashKey]HashPair),
		order: nil,
	}
}

func (h *HashValue) Debug() string {
	return "Hash"
}

func (h *HashValue) Type() ObjectType {
	return HashObject
}

func (h *HashValue) GetValue(key Hashable) (Object, bool) {
	v, ok := h.Pairs[key.HashKey()]
	if !ok {
		return &NullValue{}, ok
	}
	return v.Value, ok
}

func (h *HashValue) HasKey(key Hashable) bool {
	_, ok := h.Pairs[key.HashKey()]
	return ok
}

func (h *HashValue) Set(key Object, value Object) error {
	hashable, ok := key.(Hashable)
	if !ok {
		return fmt.Errorf("unusable as hash key: %s", key.Type())
	}
	hk := hashable.HashKey()
	if _, exists := h.Pairs[hk]; !exists {
		h.order = append(h.order, hk)
	}
	h.Pairs[hk] = HashPair{Key: key, Value: value}
	return nil
}

func (h *HashValue) Delete(key Object) error {
	hashable, ok := key.(Hashable)
	if !ok {
		return fmt.Errorf("unusable as hash key: %s", key.Type())
	}
	hk := hashable.HashKey()
	if _, exists := h.Pairs[hk]; !exists {
		return nil
	}
	delete(h.Pairs, hk)
	for i, k := range h.order {
		if k == hk {
			h.order = append(h.order[:i], h.order[i+1:]...)
			break
		}
	}
	return nil
}

func (h *HashValue) OrderedPairs() []HashPair {
	pairs := make([]HashPair, 0, len(h.order))
	for _, hk := range h.order {
		if pair, ok := h.Pairs[hk]; ok {
			pairs = append(pairs, pair)
		}
	}
	return pairs
}
