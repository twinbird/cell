package main

import (
	"fmt"
	"math"
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
	AddAssignExpression
	SubAssignExpression
	MulAssignExpression
	DivAssignExpression
	ModAssignExpression
	PowAssignExpression
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
	StringMatchExpression
	NumberPowerExpression
	LogicalAndExpression
	LogicalOrExpression
	LogicalNotExpression
	MinusExpression
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

func NewAddAssignExpression(ident string, expr *Expression) *Expression {
	e := &Expression{exprType: AddAssignExpression, ident: ident, right: expr}
	return e
}

func NewSubAssignExpression(ident string, expr *Expression) *Expression {
	e := &Expression{exprType: SubAssignExpression, ident: ident, right: expr}
	return e
}

func NewMulAssignExpression(ident string, expr *Expression) *Expression {
	e := &Expression{exprType: MulAssignExpression, ident: ident, right: expr}
	return e
}

func NewDivAssignExpression(ident string, expr *Expression) *Expression {
	e := &Expression{exprType: DivAssignExpression, ident: ident, right: expr}
	return e
}

func NewModAssignExpression(ident string, expr *Expression) *Expression {
	e := &Expression{exprType: ModAssignExpression, ident: ident, right: expr}
	return e
}

func NewPowAssignExpression(ident string, expr *Expression) *Expression {
	e := &Expression{exprType: PowAssignExpression, ident: ident, right: expr}
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

func NewStringMatchExpression(left *Expression, right *Expression) *Expression {
	e := &Expression{exprType: StringMatchExpression, left: left, right: right}
	return e
}

func NewNumberPowerExpression(left *Expression, right *Expression) *Expression {
	e := &Expression{exprType: NumberPowerExpression, left: left, right: right}
	return e
}

func NewLogicalAndExpression(left *Expression, right *Expression) *Expression {
	e := &Expression{exprType: LogicalAndExpression, left: left, right: right}
	return e
}

func NewLogicalOrExpression(left *Expression, right *Expression) *Expression {
	e := &Expression{exprType: LogicalOrExpression, left: left, right: right}
	return e
}

func NewLogicalNotExpression(left *Expression) *Expression {
	e := &Expression{exprType: LogicalNotExpression, left: left}
	return e
}

func NewMinusExpression(left *Expression) *Expression {
	e := &Expression{exprType: MinusExpression, left: left}
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
	case AddAssignExpression:
		r := e.right.eval()
		l := execContext.scope.get(e.ident)
		v := NewNumberExpression(l.asNumber() + r.asNumber())
		execContext.scope.set(e.ident, v)
		return v
	case SubAssignExpression:
		r := e.right.eval()
		l := execContext.scope.get(e.ident)
		v := NewNumberExpression(l.asNumber() - r.asNumber())
		execContext.scope.set(e.ident, v)
		return v
	case MulAssignExpression:
		r := e.right.eval()
		l := execContext.scope.get(e.ident)
		v := NewNumberExpression(l.asNumber() * r.asNumber())
		execContext.scope.set(e.ident, v)
		return v
	case DivAssignExpression:
		r := e.right.eval()
		l := execContext.scope.get(e.ident)
		v := NewNumberExpression(l.asNumber() / r.asNumber())
		execContext.scope.set(e.ident, v)
		return v
	case ModAssignExpression:
		r := e.right.eval()
		l := execContext.scope.get(e.ident)
		v := NewNumberExpression(float64(int(l.asNumber()) % int(r.asNumber())))
		execContext.scope.set(e.ident, v)
		return v
	case PowAssignExpression:
		r := e.right.eval()
		l := execContext.scope.get(e.ident)
		v := NewNumberExpression(math.Pow(l.asNumber(), r.asNumber()))
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
	case StringMatchExpression:
		left := e.left.eval().asString()
		right := e.right.eval().asString()
		b := execContext.scope.setAmpersandSpecialVars(left, right)

		if b {
			return NewNumberExpression(1)
		} else {
			return NewNumberExpression(0)
		}
	case NumberPowerExpression:
		left := e.left.eval().asNumber()
		right := e.right.eval().asNumber()

		return NewNumberExpression(math.Pow(left, right))
	case LogicalAndExpression:
		left := e.left.eval().isTruthy()
		if !left {
			return NewNumberExpression(0)
		}
		right := e.right.eval().isTruthy()
		if !right {
			return NewNumberExpression(0)
		}
		return NewNumberExpression(1)
	case LogicalOrExpression:
		left := e.left.eval().isTruthy()
		if left {
			return NewNumberExpression(1)
		}
		right := e.right.eval().isTruthy()
		if right {
			return NewNumberExpression(1)
		}
		return NewNumberExpression(0)
	case LogicalNotExpression:
		left := e.left.eval().isTruthy()
		if left {
			return NewNumberExpression(0)
		}
		return NewNumberExpression(1)
	case MinusExpression:
		left := e.left.eval().asNumber()
		return NewNumberExpression(-left)
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

func (e *Expression) isTruthy() bool {
	if e.exprType == StringExpression {
		if e.str == "" {
			return false
		} else {
			return true
		}
	}
	if e.exprType == NumberExpression {
		if e.number == 0 {
			return false
		} else {
			return true
		}
	}
	return e.eval().isTruthy()
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
	case AddAssignExpression:
		et = "AddAssignExpression"
	case SubAssignExpression:
		et = "SubAssignExpression"
	case MulAssignExpression:
		et = "MulAssignExpression"
	case DivAssignExpression:
		et = "DivAssignExpression"
	case ModAssignExpression:
		et = "ModAssignExpression"
	case PowAssignExpression:
		et = "PowAssignExpression"
	case FuncCallExpression:
		et = "FuncCallExpression"
	case NumberEQExpression:
		et = "NumberEQExpression"
	case NumberNEExpression:
		et = "NumberNEExpression"
	case NumberLTExpression:
		et = "NumberLTExpression"
	case NumberLEExpression:
		et = "NumberLEExpression"
	case NumberGTExpression:
		et = "NumberGTExpression"
	case NumberGEExpression:
		et = "NumberGEExpression"
	case StringEQExpression:
		et = "StringEQExpression"
	case StringNEExpression:
		et = "StringNEExpression"
	case StringConcatExpression:
		et = "StringConcatExpression"
	case NumberAddExpression:
		et = "NumberAddExpression"
	case NumberSubExpression:
		et = "NumberSubExpression"
	case NumberMulExpression:
		et = "NumberMulExpression"
	case NumberDivExpression:
		et = "NumberDivExpression"
	case NumberModuloExpression:
		et = "NumberModuloExpression"
	case NumberPowerExpression:
		et = "NumberPowerExpression"
	case LogicalAndExpression:
		et = "LogicalAndExpression"
	case LogicalOrExpression:
		et = "LogicalOrExpression"
	case LogicalNotExpression:
		et = "LogicalNotExpression"
	case MinusExpression:
		et = "MinusExpression"
	}
	return fmt.Sprintf("[Type: Expression] expr type: %s\n", et)
}
