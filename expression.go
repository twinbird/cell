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
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return NewNumberValue(0)
		}
		return NewNumberValue(f)
	case CellAssignExpression:
		v := e.right.eval()
		execContext.spreadsheet.setCellValue(e.left.asString(), v.asNumber())
		return v
	}
	panic("evaluate unknown type.")
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
