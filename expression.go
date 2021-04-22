package main

import (
	"fmt"
	"strconv"
)

// Expression types
const (
	NumberExpression = iota
	StringExpression
	CellReferExpression
	CellAssignExpression
	VarReferExpression
	VarAssignExpression
	FuncCallExpression
	NumberEQExpression
	NumberNEExpression
	NumberLTExpression
	NumberLEExpression
	NumberGTExpression
	NumberGEExpression
	StringEQExpression
	StringNEExpression
	StringConcatExpression
	NumberAddExpression
	NumberSubExpression
	NumberMulExpression
	NumberDivExpression
	NumberModuloExpression
)

type Expression struct {
	exprType int
	left     Node
	right    Node
	ident    string
	number   float64
	str      string
	args     *ArgList
}

func NewNumberExpression(f float64) *Expression {
	n := &Expression{exprType: NumberExpression, number: f}
	return n
}

func NewStringExpression(str string) *Expression {
	s := &Expression{exprType: StringExpression, str: str}
	return s
}

func NewCellReferExpression(axis *Expression) *Expression {
	e := &Expression{exprType: CellReferExpression, left: axis}
	return e
}

func NewCellAssignExpression(axis *Expression, expr *Expression) *Expression {
	e := &Expression{exprType: CellAssignExpression, left: axis, right: expr}
	return e
}

func NewVarReferExpression(ident string) *Expression {
	e := &Expression{exprType: VarReferExpression, ident: ident}
	return e
}

func NewVarAssignExpression(ident string, expr *Expression) *Expression {
	e := &Expression{exprType: VarAssignExpression, ident: ident, right: expr}
	return e
}

func NewFuncCallExpression(ident string, args *ArgList) *Expression {
	e := &Expression{exprType: FuncCallExpression, ident: ident, args: args}
	return e
}

func NewNumberEQExpression(left *Expression, right *Expression) *Expression {
	e := &Expression{exprType: NumberEQExpression, left: left, right: right}
	return e
}

func NewNumberNEExpression(left *Expression, right *Expression) *Expression {
	e := &Expression{exprType: NumberNEExpression, left: left, right: right}
	return e
}

func NewNumberLTExpression(left *Expression, right *Expression) *Expression {
	e := &Expression{exprType: NumberLTExpression, left: left, right: right}
	return e
}

func NewNumberLEExpression(left *Expression, right *Expression) *Expression {
	e := &Expression{exprType: NumberLEExpression, left: left, right: right}
	return e
}

func NewNumberGTExpression(left *Expression, right *Expression) *Expression {
	e := &Expression{exprType: NumberGTExpression, left: left, right: right}
	return e
}

func NewNumberGEExpression(left *Expression, right *Expression) *Expression {
	e := &Expression{exprType: NumberGEExpression, left: left, right: right}
	return e
}

func NewStringEQExpression(left *Expression, right *Expression) *Expression {
	e := &Expression{exprType: StringEQExpression, left: left, right: right}
	return e
}

func NewStringNEExpression(left *Expression, right *Expression) *Expression {
	e := &Expression{exprType: StringNEExpression, left: left, right: right}
	return e
}

func NewStringConcatExpression(left *Expression, right *Expression) *Expression {
	e := &Expression{exprType: StringConcatExpression, left: left, right: right}
	return e
}

func NewNumberAddExpression(left *Expression, right *Expression) *Expression {
	e := &Expression{exprType: NumberAddExpression, left: left, right: right}
	return e
}

func NewNumberSubExpression(left *Expression, right *Expression) *Expression {
	e := &Expression{exprType: NumberSubExpression, left: left, right: right}
	return e
}

func NewNumberMulExpression(left *Expression, right *Expression) *Expression {
	e := &Expression{exprType: NumberMulExpression, left: left, right: right}
	return e
}

func NewNumberDivExpression(left *Expression, right *Expression) *Expression {
	e := &Expression{exprType: NumberDivExpression, left: left, right: right}
	return e
}

func NewNumberModuloExpression(left *Expression, right *Expression) *Expression {
	e := &Expression{exprType: NumberModuloExpression, left: left, right: right}
	return e
}

func (e *Expression) eval() Node {
	switch e.exprType {
	case NumberExpression:
		return e
	case StringExpression:
		return e
	case CellReferExpression:
		v := execContext.spreadsheet.getCellValue(e.left.asString())

		f, ok := maybeNumber(v)
		if !ok {
			return NewStringExpression(v)
		}
		return NewNumberExpression(f)
	case CellAssignExpression:
		v := e.right.eval()

		_, isnum := maybeNumber(v.asString())
		if isnum {
			execContext.spreadsheet.setCellValue(e.left.asString(), v.asNumber())
		} else {
			execContext.spreadsheet.setCellValue(e.left.asString(), v.asString())
		}

		return v
	case VarReferExpression:
		return execContext.scope.get(e.ident)
	case VarAssignExpression:
		v := e.right.eval()
		execContext.scope.set(e.ident, v)
		return v
	case FuncCallExpression:
		f, found := execContext.functions[e.ident]
		if !found {
			fatalError("function '%s' is not found.", e.ident)
		}
		return f.call(e.args)
	case NumberEQExpression:
		left := e.left.eval().asNumber()
		right := e.right.eval().asNumber()

		if left == right {
			return NewNumberExpression(1)
		} else {
			return NewNumberExpression(0)
		}
	case NumberNEExpression:
		left := e.left.eval().asNumber()
		right := e.right.eval().asNumber()

		if left != right {
			return NewNumberExpression(1)
		} else {
			return NewNumberExpression(0)
		}
	case NumberLTExpression:
		left := e.left.eval().asNumber()
		right := e.right.eval().asNumber()

		if left < right {
			return NewNumberExpression(1)
		} else {
			return NewNumberExpression(0)
		}
	case NumberLEExpression:
		left := e.left.eval().asNumber()
		right := e.right.eval().asNumber()

		if left <= right {
			return NewNumberExpression(1)
		} else {
			return NewNumberExpression(0)
		}
	case NumberGTExpression:
		left := e.left.eval().asNumber()
		right := e.right.eval().asNumber()

		if left > right {
			return NewNumberExpression(1)
		} else {
			return NewNumberExpression(0)
		}
	case NumberGEExpression:
		left := e.left.eval().asNumber()
		right := e.right.eval().asNumber()

		if left >= right {
			return NewNumberExpression(1)
		} else {
			return NewNumberExpression(0)
		}
	case StringEQExpression:
		left := e.left.eval().asString()
		right := e.right.eval().asString()

		if left == right {
			return NewNumberExpression(1)
		} else {
			return NewNumberExpression(0)
		}
	case StringNEExpression:
		left := e.left.eval().asString()
		right := e.right.eval().asString()

		if left != right {
			return NewNumberExpression(1)
		} else {
			return NewNumberExpression(0)
		}
	case StringConcatExpression:
		left := e.left.eval().asString()
		right := e.right.eval().asString()

		return NewStringExpression(left + right)
	case NumberAddExpression:
		left := e.left.eval().asNumber()
		right := e.right.eval().asNumber()

		return NewNumberExpression(left + right)
	case NumberSubExpression:
		left := e.left.eval().asNumber()
		right := e.right.eval().asNumber()

		return NewNumberExpression(left - right)
	case NumberMulExpression:
		left := e.left.eval().asNumber()
		right := e.right.eval().asNumber()

		return NewNumberExpression(left * right)
	case NumberDivExpression:
		left := e.left.eval().asNumber()
		right := e.right.eval().asNumber()

		return NewNumberExpression(left / right)
	case NumberModuloExpression:
		left := e.left.eval().asNumber()
		right := e.right.eval().asNumber()

		return NewNumberExpression(float64(int(left) % int(right)))
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
	if e.exprType == NumberExpression {
		return e.number
	}
	if e.exprType == StringExpression {
		return 0
	}
	return e.asNumber()
}

func (e *Expression) asString() string {
	if e.exprType == StringExpression {
		return e.str
	}
	if e.exprType == NumberExpression {
		return fmt.Sprintf("%g", e.number)
	}
	return e.asString()
}

func (e *Expression) nodeType() int {
	return NodeTypeExpression
}

func (e *Expression) String() string {
	et := "unknown"
	switch e.exprType {
	case NumberExpression:
		et = "NumberExpression"
	case StringExpression:
		et = "StringExpression"
	case CellReferExpression:
		et = "CellReferExpression"
	case CellAssignExpression:
		et = "CellAssignExpression"
	case VarReferExpression:
		et = "VarReferExpression"
	case VarAssignExpression:
		et = "VarAssignExpression"
	case FuncCallExpression:
		et = "FuncCallExpression"
	}
	return fmt.Sprintf("[Type: Expression] expr type: %s\n", et)
}
