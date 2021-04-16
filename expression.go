package main

import (
	"fmt"
	"strconv"
)

// Expression types
const (
	LiteralExpression = iota
	CellReferExpression
	CellAssignExpression
)

type Expression struct {
	exprType int
	left     Node
	right    Node
}

func NewLiteralExpression(prim Primitive) *Expression {
	e := &Expression{exprType: LiteralExpression, left: prim}
	return e
}

func NewCellReferExpression(axis *Expression) *Expression {
	e := &Expression{exprType: CellReferExpression, left: axis}
	return e
}

func NewCellAssignExpression(axis *Expression, expr *Expression) *Expression {
	e := &Expression{exprType: CellAssignExpression, left: axis, right: expr}
	return e
}

func (e *Expression) eval() Node {
	switch e.exprType {
	case LiteralExpression:
		return e.left.eval()
	case CellReferExpression:
		v := execContext.spreadsheet.getCellValue(e.left.asString())

		f, ok := maybeNumber(v)
		if !ok {
			return NewStringValue(v)
		}
		return NewNumberValue(f)
	case CellAssignExpression:
		v := e.right.eval()

		_, isnum := maybeNumber(v.asString())
		if isnum {
			execContext.spreadsheet.setCellValue(e.left.asString(), v.asNumber())
		} else {
			execContext.spreadsheet.setCellValue(e.left.asString(), v.asString())
		}

		return v
	}
	panic("evaluate unknown type.")
}

func maybeNumber(val string) (float64, bool) {
	f, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return 0, false
	}
	return f, true
}

func (e *Expression) asNumber() float64 {
	return e.eval().asNumber()
}

func (e *Expression) asString() string {
	return e.eval().asString()
}

func (e *Expression) nodeType() int {
	return NodeTypeExpression
}

func (e *Expression) String() string {
	return fmt.Sprintf("[Type: Expression]")
}
