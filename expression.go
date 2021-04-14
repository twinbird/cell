package main

// Expression types
const (
	NumberExpression = iota
)

type Expression struct {
	exprType int
	number   float64
}

func NewNumberExpression(number float64) *Expression {
	e := &Expression{exprType: NumberExpression, number: number}
	return e
}

func (e *Expression) eval() float64 {
	return e.number
}
