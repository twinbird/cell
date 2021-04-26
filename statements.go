package main

import "fmt"

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

func (stmts *Statements) eval() Node {
	var ret Node
	ret = nil

	for _, s := range stmts.stmts {
		ret = s.eval()
		if execContext.doExit {
			break
		}
		if execContext.doBreak {
			break
		}
	}
	return ret
}

func (stmts *Statements) asNumber() float64 {
	panic("statements can not evaluate as a number")
}

func (stmts *Statements) asString() string {
	panic("statements can not evaluate as a string")
}

func (stmts *Statements) isTruthy() bool {
	panic("statements can not evaluate as a truthy")
}

func (stmts *Statements) nodeType() int {
	return NodeTypeStatements
}

func (stmts *Statements) String() string {
	return fmt.Sprintf("[Type: Statements] (has %d statements)", len(stmts.stmts))
}
