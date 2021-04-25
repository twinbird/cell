package main

import "fmt"

const (
	BlankStatement = iota
	ExpressionStatement
	IfStatement
	IfElseStatement
	BlockStatement
)

type Statement struct {
	stmtType int
	expr     *Expression
	thenStmt *Statement
	elseStmt *Statement
	block    *Statements
}

func NewBlankStatement() *Statement {
	s := &Statement{stmtType: BlankStatement}
	return s
}

func NewExpressionStatement(expr *Expression) *Statement {
	s := &Statement{stmtType: ExpressionStatement, expr: expr}
	return s
}

func NewIfStatement(expr *Expression, then *Statement) *Statement {
	s := &Statement{stmtType: IfStatement, expr: expr, thenStmt: then}
	return s
}

func NewIfElseStatement(expr *Expression, then *Statement, els *Statement) *Statement {
	s := &Statement{stmtType: IfElseStatement, expr: expr, thenStmt: then, elseStmt: els}
	return s
}

func NewBlockStatement(block *Statements) *Statement {
	s := &Statement{stmtType: BlockStatement, block: block}
	return s
}

func (s *Statement) eval() Node {
	switch s.stmtType {
	case BlankStatement:
		return s
	case ExpressionStatement:
		return s.expr.eval()
	case IfStatement:
		if s.expr.eval().isTruthy() {
			s.thenStmt.eval()
		}
		return NewBlankStatement()
	case IfElseStatement:
		if s.expr.eval().isTruthy() {
			s.thenStmt.eval()
		} else {
			s.elseStmt.eval()
		}
		return NewBlankStatement()
	case BlockStatement:
		return s.block.eval()
	}
	panic("evaluate unknown type.")
}

func (stmts *Statement) asNumber() float64 {
	panic("statement can not evaluate as a number")
}

func (stmts *Statement) asString() string {
	panic("statement can not evaluate as a string")
}

func (stmts *Statement) isTruthy() bool {
	panic("statement can not evaluate as a truthy")
}

func (stmts *Statement) nodeType() int {
	return NodeTypeStatement
}

func (stmts *Statement) String() string {
	return fmt.Sprintf("[Type: Statement]")
}
