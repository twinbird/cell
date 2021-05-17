package main

import (
	"math"
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

	return s.setVar(name, value)
}

func (s *Scope) setVar(name string, value Node) Node {
	s.vars[name] = value
	return value
}

func (s *Scope) get(name string) Node {
	if s.isSpecialVar(name) {
		return s.getSpecialVar(name)
	}

	return s.getVar(name)
}

func (s *Scope) getVar(name string) Node {
	v, ok := s.vars[name]
	if !ok {
		if s.parent != nil {
			return s.parent.get(name)
		}
		return NewStringExpression("")
	}
	return v
}

func (s *Scope) isSpecialVar(name string) bool {
	switch name {
	case "@":
		return true
	case "LR":
		return true
	case "LC":
		return true
	case "LCC":
		return true
	case "NF":
		return true
	case "FS":
		return true
	case "OFS":
		return true
	case "RS":
		return true
	case "ORS":
		return true
	case "NR":
		return true
	}
	if name[0] == '$' {
		return true
	}
	return false
}

func (s *Scope) getSpecialVar(name string) Node {
	switch name {
	case "@":
		s := execContext.spreadsheet.getActiveSheetName()
		return NewStringExpression(s)
	case "LR":
		// return active sheet Last Row index(start by 1)
		n := execContext.spreadsheet.getRowsCount()
		return NewNumberExpression(float64(n))
	case "LC":
		// return active sheet Last Column index(start by 1)
		n := execContext.spreadsheet.getColsCount()
		return NewNumberExpression(float64(n))
	case "LCC":
		// return active sheet Last Column index char(start by A)
		n := execContext.spreadsheet.getColsCount()
		c, err := columnNumberToName(n)
		if err != nil {
			return NewStringExpression("")
		}
		return NewStringExpression(c)
	default:
		if s.parent != nil {
			return s.parent.get(name)
		}
		return s.getVar(name)
	}
	panic("unknown special var referenced")
}

func (s *Scope) setSpecialVar(name string, value Node) Node {
	switch name {
	case "@":
		s := value.eval().asString()
		if execContext.spreadsheet.existSheetName(s) {
			if execContext.spreadsheet.setActiveSheetByName(s) != nil {
				fatalError("active sheet change error")
			}
		} else {
			if !isValidSheetName(s) {
				fatalError("sheet add error. '%s' is invalid sheet name.", s)
			}
			if err := execContext.spreadsheet.addSheet(s); err != nil {
				fatalError("sheet add error")
			}
		}
		return NewStringExpression(s)
	case "LR":
		fallthrough
	case "LC":
		fallthrough
	case "LCC":
		fatalError("special vars 'LR, LC, LCC' are readonly")
	default:
		if s.parent != nil {
			s.parent.set(name, value)
		} else {
			s.setVar(name, value)
		}
		return value
	}
	panic("assign to unknown special var")
}

func (s *Scope) setDollarSpecialVars(input string) {
	fs := execContext.scope.get("FS").asString()
	reg := s.makeFSSplitReg(fs)
	a := reg.Split(input, -1)
	execContext.scope.set("$0", NewStringExpression(input))

	if len(a) > math.MaxUint16 {
		fatalError("'%s' has too many fields", input)
	}

	s.resetDollarSpecialVars()
	execContext.ndollars = uint16(len(a))
	for i, v := range a {
		idx := strconv.Itoa(i + 1)
		execContext.scope.set("$"+idx, NewStringExpression(v))
	}
	execContext.scope.set("NF", NewNumberExpression(float64(len(a))))
}

func (s *Scope) resetDollarSpecialVars() {
	for i := 0; i < int(execContext.ndollars); i++ {
		idx := strconv.Itoa(i + 1)
		execContext.scope.set("$"+idx, NewStringExpression(""))
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

func (s *Scope) setRegexpSpecialVars(str string, reg string) bool {
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

func (s *Scope) incNR() {
	v := s.getVar("NR")
	newv := NewNumberExpression(v.asNumber() + 1.0)
	s.setVar("NR", newv)
}
