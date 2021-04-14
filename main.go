package main

import (
	"flag"
	"os"
)

var (
	program  string
	topath   string
	frompath string
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
}
