package evaluator

type Scope struct {
	store  map[string]Object
	parent *Scope
}

func NewScope() *Scope {
	return &Scope{
		store: make(map[string]Object),
	}
}

func NewChildScope(parent *Scope) *Scope {
	s := NewScope()
	s.parent = parent
	return s
}

func (s *Scope) Get(name string) (Object, bool) {
	val, ok := s.store[name]
	if !ok && s.parent != nil {
		val, ok = s.parent.Get(name)
	}
	return val, ok
}

func (s *Scope) GetLocal(name string) (Object, bool) {
	val, ok := s.store[name]
	return val, ok
}

func (s *Scope) Set(name string, val Object) {
	s.store[name] = val
}
