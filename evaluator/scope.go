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

func (s *Scope) Assign(name string, val Object) bool {

	_, ok := s.store[name]
	if !ok {
		// if the variable is not found in the current scope,
		// check the parent scope if it exists
		if s.parent != nil {
			s.parent.Assign(name, val)
		}
	} else {
		// if the variable is found in the current scope,
		// update its value
		s.store[name] = val
		return true
	}

	return false
}

func (s *Scope) SetLocal(name string, val Object) {
	s.store[name] = val
}

func (s *Scope) DeleteLocal(name string) {
	delete(s.store, name)
}
