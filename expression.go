package main

import "strconv"

// Expression types
const (
	NumberExpression = iota
	StringExpression
	CellReferExpression
	CellAssignExpression
)

type Expression struct {
	exprType int
	number   float64
	str      string
	axis     string
	right    *Expression
}

func NewNumberExpression(number float64) *Expression {
	e := &Expression{exprType: NumberExpression, number: number}
	return e
}

func NewStringExpression(str string) *Expression {
	e := &Expression{exprType: StringExpression, str: str}
	return e
}

func NewCellReferExpression(axis *Expression) *Expression {
	if axis.exprType != StringExpression {
		fatalError("the axis of the cell reference was specified as a non-string.")
	}
	e := &Expression{exprType: CellReferExpression, axis: axis.str}
	return e
}

func NewCellAssignExpression(axis *Expression, expr *Expression) *Expression {
	if axis.exprType != StringExpression {
		fatalError("the axis of the cell reference was specified as a non-string.")
	}
	e := &Expression{exprType: CellAssignExpression, axis: axis.str, right: expr}
	return e
}

func (e *Expression) eval() float64 {
	switch e.exprType {
	case NumberExpression:
		return e.number
	case StringExpression:
		return 0
	case CellReferExpression:
		v := execContext.spreadsheet.getCellValue(e.axis)
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return 0
		}
		return f
	case CellAssignExpression:
		v := e.right.eval()
		execContext.spreadsheet.setCellValue(e.axis, v)
		return v
	}
	panic("evaluate unknown type.")
}
