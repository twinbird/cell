package main

// Expression types
const (
	NumberExpression = iota
	StringExpression
)

type Expression struct {
	exprType int
	number   float64
	str      string
}

func NewNumberExpression(number float64) *Expression {
	e := &Expression{exprType: NumberExpression, number: number}
	return e
}

func NewStringExpression(str string) *Expression {
	e := &Expression{exprType: StringExpression, str: str}
	return e
}

func (e *Expression) eval() float64 {
	return e.number
}
