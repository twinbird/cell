package main

import (
	"flag"
	"fmt"
	"os"
)

type ExecContext struct {
	code        string
	topath      string
	frompath    string
	spreadsheet *Spreadsheet
	exitCode    int
	scope       *Scope
}

var execContext *ExecContext

func NewExecContext() *ExecContext {
	con := &ExecContext{}
	con.scope = NewScope()

	return con
}

func main() {
	con := NewExecContext()

	f := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	f.StringVar(&con.topath, "to", "", "output spreadsheet filepath")
	f.StringVar(&con.frompath, "from", "", "input spreadsheet filepath")

	f.Parse(os.Args[2:])

	con.code = os.Args[1]

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
	ret := execScript()
	afterRun()

	execContext.exitCode = ret
}

func execScript() int {
	yyDebug = 1
	yyErrorVerbose = true

	lexer := NewLexer(execContext.code)
	yyParse(lexer)

	n := lexer.ast.eval()
	return int(n.asNumber())
}

func fatalError(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, "ERROR: "+format+"\n", a)
	os.Exit(1)
}
