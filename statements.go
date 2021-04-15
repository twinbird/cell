package main

type Statements struct {
	stmts []*Statement
}

func NewStatements(stmt *Statement) *Statements {
	ary := make([]*Statement, 1)
	ary[0] = stmt
	s := &Statements{stmts: ary}
	return s
}

func (stmts *Statements) appendStatement(stmt *Statement) *Statements {
	stmts.stmts = append(stmts.stmts, stmt)
	return stmts
}

func (stmts *Statements) eval() float64 {
	ret := 0.0
	for _, s := range stmts.stmts {
		ret = s.eval()
	}
	return ret
}
