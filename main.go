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

func beforeRun(con *ExecContext) {
	sheet, err := NewSpreadsheet(con.frompath, con.topath)
	if err != nil {
		fatalError(err)
	}
	con.spreadsheet = sheet
}

func afterRun(con *ExecContext) {
	if con.topath != "" {
		if err := con.spreadsheet.writeSpreadsheet(); err != nil {
			fatalError(err)
		}
	}
}

func run(con *ExecContext) {
	beforeRun(con)
	ret := execScript(con)
	afterRun(con)

	con.exitCode = ret
}

func execScript(con *ExecContext) int {
	yyDebug = 1
	yyErrorVerbose = true

	lexer := NewLexer(con.code)
	yyParse(lexer)

	return int(lexer.ast.eval())
}

func fatalError(err error) {
	fmt.Fprintf(os.Stderr, "FATAL ERROR: %v\n", err)
	os.Exit(1)
}
