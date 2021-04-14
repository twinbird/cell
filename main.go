package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	spreadsheet    *Spreadsheet
	program        string
	optionTopath   string
	optionFrompath string
)

func main() {
	f := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	f.StringVar(&optionTopath, "to", "", "output spreadsheet filepath")
	f.StringVar(&optionFrompath, "from", "", "input spreadsheet filepath")

	f.Parse(os.Args[2:])

	program = os.Args[1]

	run()
}

func beforeRun() {
	sheet, err := NewSpreadsheet(optionFrompath, optionTopath)
	if err != nil {
		fatalError(err)
	}
	spreadsheet = sheet
}

func afterRun() {
	if optionTopath != "" {
		if err := spreadsheet.writeSpreadsheet(); err != nil {
			fatalError(err)
		}
	}
}

func run() {
	beforeRun()
	ret := execScript()
	afterRun()

	os.Exit(ret)
}

func execScript() int {
	yyDebug = 1
	yyErrorVerbose = true

	lexer := NewLexer(program)
	yyParse(lexer)

	return int(lexer.ast.eval())
}

func fatalError(err error) {
	fmt.Fprintf(os.Stderr, "FATAL ERROR: %v\n", err)
	os.Exit(1)
}
