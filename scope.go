package main

import (
	"regexp"
	"strconv"
)

type Scope struct {
	vars   map[string]Node
	parent *Scope
}

func NewScope() *Scope {
	s := &Scope{}
	s.vars = make(map[string]Node)
	return s
}

func AppendScope(s *Scope) *Scope {
	ns := &Scope{parent: s}
	ns.vars = make(map[string]Node)
	return ns
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

func (s *Scope) setDollarSpecialVars(input string) {
	fs := execContext.scope.get("FS").asString()
	reg := s.makeFSSplitReg(fs)
	a := reg.Split(input, -1)
	execContext.scope.set("$0", NewStringExpression(input))

	for i, v := range a {
		idx := strconv.Itoa(i + 1)
		execContext.scope.set("$"+idx, NewStringExpression(v))
	}
}

func (s *Scope) makeFSSplitReg(fs string) *regexp.Regexp {
	// Note:
	//   FS rule imitated gawk style
	//   1. just one char space => space or tab or new line
	//   2. just one char       => as it is
	//   3. other               => consider as a regexp

	r := fs
	if fs == " " {
		r = "( |\t|\n)+"
	} else if len(fs) == 1 {
		r = regexp.QuoteMeta(fs)
	}
	reg, err := regexp.Compile(r)
	if err != nil {
		fatalError("FS '%s' is invalid format", fs)
	}
	return reg
}

func (s *Scope) setAmpersandSpecialVars(str string, reg string) bool {
	r := regexp.MustCompile(reg)
	if !r.MatchString(str) {
		return false
	}
	g := r.FindStringSubmatch(str)

	for i, v := range g {
		idx := strconv.Itoa(i)
		execContext.scope.set("$_"+idx, NewStringExpression(v))
	}

	return true
}
