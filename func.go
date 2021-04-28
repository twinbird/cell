package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

type ArgList struct {
	args []*Expression
}

func NewArgList(expr *Expression) *ArgList {
	a := &ArgList{}
	a.args = make([]*Expression, 1)
	a.args[0] = expr

	return a
}

func (args *ArgList) appendArg(expr *Expression) *ArgList {
	args.args = append(args.args, expr)
	return args
}

func NewEmptyArgList() *ArgList {
	a := &ArgList{}
	a.args = make([]*Expression, 0)
	return a
}

type ParamList struct {
	params []string
}

func NewParamList(ident string) *ParamList {
	p := &ParamList{}
	p.params = make([]string, 1)
	p.params[0] = ident

	return p
}

func (params *ParamList) appendParam(ident string) *ParamList {
	params.params = append(params.params, ident)
	return params
}

func NewEmptyParamList() *ParamList {
	p := &ParamList{}
	p.params = make([]string, 0)

	return p
}

const (
	FunctionTypeBuiltin = iota
	FunctionTypeDefine
)

type Function struct {
	funcType       int
	builtin        func(args ...Node) Node
	defineParams   *ParamList
	defineStmt     *Statement
	defineFuncName string
}

func NewBuiltinFunction(f func(args ...Node) Node) *Function {
	return &Function{
		funcType: FunctionTypeBuiltin,
		builtin:  f,
	}
}

func defineFunction(name string, params *ParamList, stmt *Statement) {
	f := &Function{
		funcType:       FunctionTypeDefine,
		defineParams:   params,
		defineStmt:     stmt,
		defineFuncName: name,
	}
	_, exist := execContext.functions[name]
	if exist {
		fatalError("function '%s' is already defined", name)
	}
	execContext.functions[name] = f
}

func (f *Function) call(args *ArgList) Node {
	if f.funcType == FunctionTypeBuiltin {
		ev := make([]Node, 0)
		for _, v := range args.args {
			ev = append(ev, v.eval())
		}
		return f.builtin(ev...)
	} else {
		if len(args.args) != len(f.defineParams.params) {
			fatalError("invalid as number of arguments for %s", f.defineFuncName)
		}

		execContext.scope = AppendScope(execContext.scope)

		for i, p := range f.defineParams.params {
			if i < len(args.args) {
				execContext.scope.set(p, args.args[i].eval())
			} else {
				execContext.scope.set(p, NewStringExpression(""))
			}
		}
		execContext.funcRet = nil
		f.defineStmt.eval()

		var ret Node
		if execContext.funcRet != nil {
			ret = execContext.funcRet
			execContext.funcRet = nil
		} else {
			ret = NewStringExpression("")
		}

		execContext.scope = execContext.scope.parent

		return ret
	}
}

func builtinFunctions() map[string]*Function {
	f := map[string]*Function{
		"exit":  NewBuiltinFunction(builtinExit),
		"abort": NewBuiltinFunction(builtinAbort),
		"gets":  NewBuiltinFunction(builtinGets),
		"puts":  NewBuiltinFunction(builtinPuts),
	}

	return f
}

// exit(number) void
// Exit program.If "to" option specified, 'cell' will save editing spreadsheet.
func builtinExit(args ...Node) Node {
	exitCode := args[0]
	execContext.exitCode = int(exitCode.asNumber())
	execContext.doExit = true
	return nil
}

// abort(number)
// Exit program immediately
func builtinAbort(args ...Node) Node {
	exitCode := args[0]
	os.Exit(int(exitCode.asNumber()))
	return nil
}

// gets(void) string
// Get character line from stdin.
func builtinGets(args ...Node) Node {
	s, err := execContext.in.ReadString('\n')
	if err != io.EOF && err != nil {
		fatalError("builtin function 'gets' raised error '%v'", err)
	}
	s = strings.TrimRight(s, "\r\n")
	execContext.scope.setDollarSpecialVars(s)
	return NewStringExpression(s)
}

// puts(string)
// Print string and new line to stdout.
func builtinPuts(args ...Node) Node {
	if len(args) == 0 {
		v := execContext.scope.get("$0")
		s := v.asString()
		fmt.Fprintf(execContext.out, "%s\n", s)
		return nil
	}
	ofs := execContext.scope.get("OFS").asString()
	s := args[0].asString()
	for i := 1; i < len(args); i++ {
		s = ofs + s
		s = args[i].asString() + s
	}
	fmt.Fprintf(execContext.out, "%s\n", s)
	return nil
}
