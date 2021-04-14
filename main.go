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
}

var execContext *ExecContext

func main() {
	con := &ExecContext{}

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
		fatalError(err)
	}
	execContext.spreadsheet = sheet
}

func afterRun() {
	if execContext.topath != "" {
		if err := execContext.spreadsheet.writeSpreadsheet(); err != nil {
			fatalError(err)
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

	return int(lexer.ast.eval())
}

func fatalError(err error) {
	fmt.Fprintf(os.Stderr, "FATAL ERROR: %v\n", err)
	os.Exit(1)
}
