package main

import (
	"fmt"
	"io"
	"math"
	"math/rand"
	"os"
	"strings"
	"time"
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
		execContext.doReturn = false
		f.defineStmt.eval()

		var ret Node
		if execContext.doReturn {
			ret = execContext.funcRet
			execContext.doReturn = false
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
		"exit":   NewBuiltinFunction(builtinExit),
		"abort":  NewBuiltinFunction(builtinAbort),
		"gets":   NewBuiltinFunction(builtinGets),
		"puts":   NewBuiltinFunction(builtinPuts),
		"head":   NewBuiltinFunction(builtinHead),
		"tail":   NewBuiltinFunction(builtinTail),
		"rename": NewBuiltinFunction(builtinRename),
		"exist":  NewBuiltinFunction(builtinExist),
		"count":  NewBuiltinFunction(builtinCount),
		"delete": NewBuiltinFunction(builtinDelete),
		"copy":   NewBuiltinFunction(builtinCopy),
		"srand":  NewBuiltinFunction(builtinSrand),
		"rand":   NewBuiltinFunction(builtinRand),
		"floor":  NewBuiltinFunction(builtinFloor),
		"ceil":   NewBuiltinFunction(builtinCeil),
		"round":  NewBuiltinFunction(builtinRound),
	}

	return f
}

// exit(number) noreturn
// Exit program.If "to" option specified, 'cell' will save editing spreadsheet.
func builtinExit(args ...Node) Node {
	if len(args) != 1 {
		fatalError("invalid as number of arguments for exit()")
	}
	exitCode := args[0]
	execContext.exitCode = int(exitCode.asNumber())
	execContext.doExit = true
	return nil
}

// abort(number) noreturn
// Exit program immediately
func builtinAbort(args ...Node) Node {
	if len(args) != 1 {
		fatalError("invalid as number of arguments for abort()")
	}
	exitCode := args[0]
	os.Exit(int(exitCode.asNumber()))
	return nil
}

// gets(void) string
// Get character line from stdin.
func builtinGets(args ...Node) Node {
	if len(args) != 0 {
		fatalError("invalid as number of arguments for gets()")
	}
	rs := execContext.scope.get("RS").asString()

	s, err := execContext.in.ReadString(rs[0])
	if err != io.EOF && err != nil {
		fatalError("builtin function 'gets' raised error '%v'", err)
	}
	execContext.scope.incNR()

	if rs == "\n" {
		s = strings.TrimRight(s, "\r\n")
	} else {
		s = strings.TrimRight(s, rs)
	}
	execContext.scope.setDollarSpecialVars(s)
	return NewStringExpression(s)
}

// puts(string) string
// Print string and new line to stdout.
// And return puts string(No include ORS).
func builtinPuts(args ...Node) Node {
	ors := execContext.scope.get("ORS").asString()

	if len(args) == 0 {
		v := execContext.scope.get("$0")
		s := v.asString()
		fmt.Fprintf(execContext.out, "%s%s", s, ors)
		return NewStringExpression(s)
	}
	ofs := execContext.scope.get("OFS").asString()
	s := args[0].asString()
	for i := 1; i < len(args); i++ {
		s = ofs + s
		s = args[i].asString() + s
	}
	fmt.Fprintf(execContext.out, "%s%s", s, ors)
	return NewStringExpression(s)
}

// head() string
// Set the active sheet to the first sheet
// And return active sheet name
func builtinHead(args ...Node) Node {
	if len(args) != 0 {
		fatalError("invalid as number of arguments for head()")
	}
	execContext.spreadsheet.setHeadSheet()
	s := execContext.spreadsheet.getActiveSheetName()
	return NewStringExpression(s)
}

// tail() string
// Set the active sheet to the last sheet
// And return active sheet name
func builtinTail(args ...Node) Node {
	if len(args) != 0 {
		fatalError("invalid as number of arguments for tail()")
	}
	execContext.spreadsheet.setTailSheet()
	s := execContext.spreadsheet.getActiveSheetName()
	return NewStringExpression(s)
}

// exist(name) number
// return if exists 'name' sheet 1, else 0
func builtinExist(args ...Node) Node {
	if len(args) != 1 {
		fatalError("invalid as number of arguments for exist()")
	}
	b := execContext.spreadsheet.existSheetName(args[0].asString())
	if b {
		return NewNumberExpression(1)
	}
	return NewNumberExpression(0)
}

// rename(old, new)
// rename sheet name
// return the changed name if successful
func builtinRename(args ...Node) Node {
	if len(args) != 2 {
		fatalError("invalid as number of arguments for rename()")
	}
	o := args[1].asString()
	n := args[0].asString()

	if !execContext.spreadsheet.existSheetName(o) {
		fatalError("rename(): sheet '%s' not exist", o)
	}
	if execContext.spreadsheet.existSheetName(n) {
		fatalError("rename(): sheet '%s' already exist", n)
	}
	if !isValidSheetName(n) {
		fatalError("rename(): sheet '%s' is invalid sheet name", n)
	}
	s := execContext.spreadsheet.setSheetName(o, n)

	return NewStringExpression(s)
}

// count() number
// count sheets
func builtinCount(args ...Node) Node {
	if len(args) != 0 {
		fatalError("invalid as number of arguments for count()")
	}
	n := execContext.spreadsheet.countSheet()
	return NewNumberExpression(float64(n))
}

// delete(string)
// delete specify sheet
func builtinDelete(args ...Node) Node {
	if len(args) != 1 {
		fatalError("invalid as number of arguments for delete()")
	}
	s := args[0].asString()

	if !execContext.spreadsheet.existSheetName(s) {
		fatalError("delete(): sheet '%s' not exist", s)
	}
	if execContext.spreadsheet.countSheet() <= 1 {
		fatalError("delete(): could not delete last sheet")
	}

	execContext.spreadsheet.deleteSheet(s)

	return NewStringExpression("")
}

// copy(string[from], string[to]) string[to]
// copy from [from] sheet to [to] sheet
func builtinCopy(args ...Node) Node {
	if len(args) != 2 {
		fatalError("invalid as number of arguments for copy()")
	}
	from := args[1].asString()
	to := args[0].asString()

	if !execContext.spreadsheet.existSheetName(from) {
		fatalError("copy(): sheet '%s' not exist", from)
	}
	if execContext.spreadsheet.existSheetName(to) {
		fatalError("copy(): sheet '%s' already exist", to)
	}
	if !isValidSheetName(to) {
		fatalError("copy(): sheet '%s' is invalid sheet name", to)
	}

	execContext.spreadsheet.copySheet(from, to)

	return NewStringExpression(to)
}

// srand([expr])
// Use expr as the new seed for the random number generator.  If no expr is provided, use the current time.
func builtinSrand(args ...Node) Node {
	if 1 < len(args) {
		fatalError("invalid as number of arguments for srand()")
	}
	n := time.Now().UnixNano()
	if len(args) > 0 {
		n = int64(args[0].asNumber())
	}
	rand.Seed(n)
	return NewStringExpression("")
}

// rand() number
// Return a random number N, between zero and one, such that 0 <= N <= 1.
func builtinRand(args ...Node) Node {
	if len(args) != 0 {
		fatalError("invalid as number of arguments for rand()")
	}
	f := rand.Float64()
	return NewNumberExpression(f)
}

// floor(number) number
func builtinFloor(args ...Node) Node {
	if len(args) != 1 {
		fatalError("invalid as number of arguments for floor()")
	}
	f := args[0].asNumber()
	v := math.Floor(f)
	return NewNumberExpression(v)
}

// ceil(number) number
func builtinCeil(args ...Node) Node {
	if len(args) != 1 {
		fatalError("invalid as number of arguments for ceil()")
	}
	f := args[0].asNumber()
	v := math.Ceil(f)
	return NewNumberExpression(v)
}

// round(number) number
func builtinRound(args ...Node) Node {
	if len(args) != 1 {
		fatalError("invalid as number of arguments for round()")
	}
	f := args[0].asNumber()
	v := math.Round(f)
	return NewNumberExpression(v)
}
