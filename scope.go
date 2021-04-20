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
	if s.isSpecialVar(name) {
		return s.setSpecialVar(name, value)
	}

	s.vars[name] = value
	return value
}

func (s *Scope) get(name string) Node {
	if s.isSpecialVar(name) {
		return s.getSpecialVar(name)
	}

	v, ok := s.vars[name]
	if !ok {
		return NewStringExpression("")
	}
	return v
}

func (s *Scope) isSpecialVar(name string) bool {
	switch name {
	case "@":
		return true
	}
	return false
}

func (s *Scope) getSpecialVar(name string) Node {
	switch name {
	case "@":
		s := execContext.spreadsheet.getActiveSheetName()
		return NewStringExpression(s)
	}
	panic("unknown special var referenced")
}

func (s *Scope) setSpecialVar(name string, value Node) Node {
	switch name {
	case "@":
		s := value.eval().asString()
		err := execContext.spreadsheet.setActiveSheetByName(s)
		if err != nil {
			if err := execContext.spreadsheet.addSheet(s); err != nil {
				fatalError("active sheet change error")
			}
		}
		return NewStringExpression(s)
	}
	panic("assign to unknown special var")
}
