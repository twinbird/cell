package main

type Scope struct {
	vars map[string]Node
}

func NewScope() *Scope {
	s := &Scope{}
	s.vars = make(map[string]Node)
	return s
}

func (s *Scope) set(name string, value Node) Node {
	s.vars[name] = value
	return value
}

func (s *Scope) get(name string) Node {
	v, ok := s.vars[name]
	if !ok {
		return NewStringExpression("")
	}
	return v
}
