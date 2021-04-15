package main

const (
	ExpressionStatement = iota
)

type Statement struct {
	stmtType int
	expr     *Expression
}

func NewExpressionStatement(expr *Expression) *Statement {
	s := &Statement{stmtType: ExpressionStatement, expr: expr}
	return s
}

func (s *Statement) eval() float64 {
	switch s.stmtType {
	case ExpressionStatement:
		return s.expr.eval()
	}
	panic("evaluate unknown type.")
}
