package main

import "fmt"

const (
	BlankStatement = iota
	ExpressionStatement
)

type Statement struct {
	stmtType int
	expr     *Expression
}

func NewBlankStatement() *Statement {
	s := &Statement{stmtType: BlankStatement}
	return s
}

func NewExpressionStatement(expr *Expression) *Statement {
	s := &Statement{stmtType: ExpressionStatement, expr: expr}
	return s
}

func (s *Statement) eval() Node {
	switch s.stmtType {
	case BlankStatement:
		return s
	case ExpressionStatement:
		return s.expr.eval()
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
