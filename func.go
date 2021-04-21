package main

import (
	"bufio"
	"fmt"
	"os"
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

type Function struct {
	builtin func(args ...Node) Node
}

func NewBuiltinFunction(f func(args ...Node) Node) *Function {
	return &Function{
		builtin: f,
	}
}

func (f *Function) call(args *ArgList) Node {
	ev := make([]Node, 0)
	for _, v := range args.args {
		ev = append(ev, v.eval())
	}
	return f.builtin(ev...)
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
	scanner := bufio.NewScanner(execContext.in)
	scanner.Scan()
	s := scanner.Text()
	execContext.scope.setDollarSpecialVars(s)
	return NewStringExpression(s)
}

// puts(string)
// Print string and new line to stdout.
func builtinPuts(args ...Node) Node {
	s := args[0].asString()
	fmt.Fprintf(execContext.out, "%s\n", s)
	return nil
}
