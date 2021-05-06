//go:generate goyacc parser.y
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
)

type ExecContext struct {
	code        string
	topath      string
	frompath    string
	spreadsheet *Spreadsheet
	exitCode    int
	scope       *Scope
	ndollars    uint16
	functions   map[string]*Function
	funcRet     Node
	doExit      bool
	doBreak     bool
	doContinue  bool
	doReturn    bool
	in          *bufio.Reader
	out         io.Writer
	errout      io.Writer
}

var execContext *ExecContext

func NewExecContext() *ExecContext {
	con := &ExecContext{}
	con.scope = NewScope()
	con.functions = builtinFunctions()
	con.in = bufio.NewReader(os.Stdin)
	con.out = os.Stdout
	con.errout = os.Stderr
	con.scope.set("FS", NewStringExpression(" "))
	con.scope.set("OFS", NewStringExpression(" "))
	con.scope.set("RS", NewStringExpression("\n"))
	con.scope.set("ORS", NewStringExpression("\n"))

	return con
}

func main() {
	con := NewExecContext()

	flag.StringVar(&con.topath, "to", "", "output spreadsheet filepath")
	flag.StringVar(&con.frompath, "from", "", "input spreadsheet filepath")

	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		flag.Usage()
		os.Exit(1)
	}

	con.code = args[0]

	run(con)
	os.Exit(con.exitCode)
}

func beforeRun() {
	sheet, err := NewSpreadsheet(execContext.frompath, execContext.topath)
	if err != nil {
		fatalError("on error occured creating new spreadsheet")
	}
	execContext.spreadsheet = sheet
}

func afterRun() {
	if execContext.topath != "" {
		if err := execContext.spreadsheet.writeSpreadsheet(); err != nil {
			fatalError("on error occured writting spreadsheet")
		}
	}
}

func run(con *ExecContext) {
	execContext = con

	beforeRun()
	execScript()
	afterRun()
}

func execScript() int {
	yyDebug = 1
	yyErrorVerbose = true

	lexer := NewLexer(execContext.code)
	yyParse(lexer)

	lexer.ast.eval()

	if execContext.doBreak {
		fatalError("'break' is not allowed outside a loop")
	}
	if execContext.doContinue {
		fatalError("'continue' is not allowed outside a loop")
	}

	return 0
}

func fatalError(format string, a ...interface{}) {
	if len(a) > 0 {
		fmt.Fprintf(os.Stderr, "ERROR: "+format+"\n", a)
	} else {
		fmt.Fprintf(os.Stderr, "ERROR: "+format+"\n")
	}
	os.Exit(1)
}
