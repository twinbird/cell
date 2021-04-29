package main

import "fmt"

const (
	BlankStatement = iota
	ExpressionStatement
	IfStatement
	IfElseStatement
	BlockStatement
	WhileStatement
	DoWhileStatement
	ForStatement
	BreakStatement
	ContinueStatement
	FunctionStatement
	ReturnStatement
)

type Statement struct {
	stmtType int
	expr     *Expression
	init     *Expression
	inc      *Expression
	thenStmt *Statement
	elseStmt *Statement
	block    *Statements
	params   *ParamList
	funcName string
}

func NewBlankStatement() *Statement {
	s := &Statement{stmtType: BlankStatement}
	return s
}

func NewExpressionStatement(expr *Expression) *Statement {
	s := &Statement{stmtType: ExpressionStatement, expr: expr}
	return s
}

func NewIfStatement(expr *Expression, then *Statement) *Statement {
	s := &Statement{stmtType: IfStatement, expr: expr, thenStmt: then}
	return s
}

func NewIfElseStatement(expr *Expression, then *Statement, els *Statement) *Statement {
	s := &Statement{stmtType: IfElseStatement, expr: expr, thenStmt: then, elseStmt: els}
	return s
}

func NewBlockStatement(block *Statements) *Statement {
	s := &Statement{stmtType: BlockStatement, block: block}
	return s
}

func NewWhileStatement(expr *Expression, then *Statement) *Statement {
	s := &Statement{stmtType: WhileStatement, expr: expr, thenStmt: then}
	return s
}

func NewDoWhileStatement(then *Statement, expr *Expression) *Statement {
	s := &Statement{stmtType: DoWhileStatement, expr: expr, thenStmt: then}
	return s
}

func NewForStatement(init *Expression, cond *Expression, inc *Expression, then *Statement) *Statement {
	s := &Statement{stmtType: ForStatement, init: init, expr: cond, inc: inc, thenStmt: then}
	return s
}

func NewBreakStatement() *Statement {
	s := &Statement{stmtType: BreakStatement}
	return s
}

func NewContinueStatement() *Statement {
	s := &Statement{stmtType: ContinueStatement}
	return s
}

func NewFunctionDefineStatement(name string, params *ParamList, then *Statement) *Statement {
	s := &Statement{stmtType: FunctionStatement, funcName: name, params: params, thenStmt: then}
	return s
}

func NewReturnStatement(expr *Expression) *Statement {
	s := &Statement{stmtType: ReturnStatement, expr: expr}
	return s
}

func (s *Statement) eval() Node {
	switch s.stmtType {
	case BlankStatement:
		return s
	case ExpressionStatement:
		return s.expr.eval()
	case IfStatement:
		if s.expr.eval().isTruthy() {
			s.thenStmt.eval()
		}
		return NewBlankStatement()
	case IfElseStatement:
		if s.expr.eval().isTruthy() {
			s.thenStmt.eval()
		} else {
			s.elseStmt.eval()
		}
		return NewBlankStatement()
	case BlockStatement:
		return s.block.eval()
	case WhileStatement:
		for s.expr.eval().isTruthy() {
			s.thenStmt.eval()
			if execContext.doExit {
				break
			}
			if execContext.doBreak {
				execContext.doBreak = false
				break
			}
			if execContext.doContinue {
				execContext.doContinue = false
				continue
			}
		}
		return NewBlankStatement()
	case DoWhileStatement:
		s.thenStmt.eval()
		for s.expr.eval().isTruthy() {
			s.thenStmt.eval()
			if execContext.doExit {
				break
			}
			if execContext.doBreak {
				execContext.doBreak = false
				break
			}
			if execContext.doContinue {
				execContext.doContinue = false
				continue
			}
		}
		return NewBlankStatement()
	case ForStatement:
		s.init.eval()
		for s.expr.eval().isTruthy() {
			s.thenStmt.eval()
			if execContext.doExit {
				break
			}
			if execContext.doBreak {
				execContext.doBreak = false
				break
			}
			s.inc.eval()
			if execContext.doContinue {
				execContext.doContinue = false
				continue
			}
		}
		return NewBlankStatement()
	case BreakStatement:
		execContext.doBreak = true
		return NewBlankStatement()
	case ContinueStatement:
		execContext.doContinue = true
		return NewBlankStatement()
	case FunctionStatement:
		defineFunction(s.funcName, s.params, s.thenStmt)
		return NewBlankStatement()
	case ReturnStatement:
		execContext.funcRet = s.expr.eval()
		execContext.doReturn = true
		return NewBlankStatement()
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
