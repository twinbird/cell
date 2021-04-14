package main

import (
	"flag"
	"log"
	"os"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
)

var (
	spreadsheet *excelize.File
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
	openSpreadsheet()
	defer closeSpreadsheet()

	execScript()
}

func openSpreadsheet() {
	if frompath != "" {
		err := readSpreadsheet()
		if err != nil {
			log.Fatalf("on error spreadsheet reading '%v'", err)
		}
	} else {
		err := createSpreadsheet()
		if err != nil {
			log.Fatalf("on error spreadsheet creating '%v'", err)
		}
	}
}

func createSpreadsheet() error {
	spreadsheet = excelize.NewFile()
	return nil
}

func readSpreadsheet() error {
	f, err := excelize.OpenFile(frompath)
	if err != nil {
		return err
	}
	spreadsheet = f

	return nil
}

func closeSpreadsheet() {
	if topath != "" {
		err := spreadsheet.SaveAs(topath)
		if err != nil {
			log.Fatalf("on error spreadsheet closing. '%v'", err)
		}
	}
}

func execScript() {
}
