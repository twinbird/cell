package main

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

func (s *Statement) eval() float64 {
	switch s.stmtType {
	case BlankStatement:
		return 0
	case ExpressionStatement:
		return s.expr.eval()
	}
	panic("evaluate unknown type.")
}
