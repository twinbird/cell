package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	spreadsheet *Spreadsheet
	program     string
	topath      string
	frompath    string
)

func main() {
	f := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	f.StringVar(&topath, "to", "", "output spreadsheet filepath")
	f.StringVar(&frompath, "from", "", "input spreadsheet filepath")

	f.Parse(os.Args[2:])

	program = os.Args[1]

	run()
}

func run() {
	sheet, err := NewSpreadsheet(frompath, topath)
	if err != nil {
		fatalError(err)
	}
	spreadsheet = sheet

	execScript()

	if err := spreadsheet.writeSpreadsheet(); err != nil {
		fatalError(err)
	}
}

func execScript() {
}

func fatalError(err error) {
	fmt.Fprintf(os.Stderr, "FATAL ERROR: %v\n", err)
	os.Exit(1)
}
