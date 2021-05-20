//go:generate goyacc parser.y
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

const CELL_VERSION = "0.1.0"

type ExecContext struct {
	code           string
	topath         string
	frompath       string
	spreadsheet    *Spreadsheet
	exitCode       int
	scope          *Scope
	ndollars       uint16
	functions      map[string]*Function
	funcRet        Node
	doExit         bool
	doBreak        bool
	doContinue     bool
	doReturn       bool
	in             *bufio.Reader
	out            io.Writer
	errout         io.Writer
	doTextRowLoop  bool
	doExcelRowLoop bool
	initSheet      string
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
	con.scope.set("NR", NewNumberExpression(0))
	con.scope.set("SER", NewNumberExpression(1))

	return con
}

func main() {
	flag.Usage = usage

	con := NewExecContext()

	var pgpath string
	var showVer bool
	var fs string
	var ser int
	flag.StringVar(&con.topath, "to", "", "output xlsx filepath")
	flag.StringVar(&con.frompath, "from", "", "input xlsx filepath")
	flag.StringVar(&pgpath, "f", "", "program filepath")
	flag.StringVar(&fs, "F", "", "specify field separator")
	flag.BoolVar(&showVer, "V", false, "show version")
	flag.BoolVar(&con.doTextRowLoop, "n", false, "wrap your script inside while(gets()){... ;} loop")
	flag.BoolVar(&con.doExcelRowLoop, "N", false, "wrap your script inside for(NER = SER; NER <= LR; NER++){... ;} loop")
	flag.IntVar(&ser, "s", 1, "specify special var SER(start excel row)")
	flag.StringVar(&con.initSheet, "S", "", "specify active sheet by name")

	flag.Parse()

	// -V option
	if showVer {
		showVersion()
	}

	args := flag.Args()
	if len(args) < 1 && pgpath == "" {
		flag.Usage()
		os.Exit(1)
	}

	// -f option
	if pgpath != "" {
		con.code = readProg(pgpath)
	} else {
		con.code = args[0]
	}

	// -F option
	if fs != "" {
		con.scope.set("FS", NewStringExpression(fs))
	}

	// -s option
	con.scope.set("SER", NewNumberExpression(float64(ser)))

	// text file specify
	if pgpath == "" {
		if 1 < len(args) {
			switchStdin(con, args[1:])
		}
	} else {
		if 0 < len(args) {
			switchStdin(con, args)
		}
	}

	run(con)
	os.Exit(con.exitCode)
}

func switchStdin(con *ExecContext, files []string) {
	rary := make([]io.Reader, len(files))

	for i := 0; i < len(files); i++ {
		f, err := os.Open(files[i])
		if err != nil {
			fatalError("could not open file '%s'", f)
		}

		rary[i] = f
	}
	con.in = bufio.NewReader(io.MultiReader(rary...))
}

func showVersion() {
	fmt.Printf("Cell %s\n", CELL_VERSION)
	os.Exit(0)
}

func usage() {
	msg := `Cell is an Excel file(xlsx) processing language for command line.

Usage: cell [options] 'program' [file...]
Usage: cell [options] -f programfile [file...]

Options:
  -to output-xlsx-file-path
      Specify the path of the processed Excel file that will be saved
  -from input-xlsx-file-path
      Specify the Excel file to be processed. No overwriting will be done. The default is an empty book containing only Sheet1.
  -f program-file
      Read the Cell program source from the file program-file, instead of from the first command line argument.
  -F fs
      Use fs for the input field separator (the value of the FS predefined variable).
  -n
      Wrap your script inside while(gets()){... ;} loop
  -N
      Wrap your script inside for(NER = SER; NER <= LR; NER++){... ;} loop (NER and SER, LR are predefined variables)
  -s row-no
      Specify the special variable SER(Start Excel Row) (default 1)
  -S
      Specify default active sheet by name
  -V
      Print version information.
  -h
      Show this help

Examples:
        cell -to greeting.xlsx '["A1"] = "Hello, world"'
        cell -F ":" -to users.xlsx -n '["A".NR] = $1' /etc/passwd`

	fmt.Fprintf(os.Stderr, "%s\n", msg)
}

func readProg(filename string) string {
	if !fileExist(filename) {
		fatalError("program file '%s' is not found", filename)
	}

	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		fatalError("on error occured reading file '%s'", filename)
	}

	return string(bytes)
}

func fileExist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func beforeRun() {
	// setup spreadsheet
	sheet, err := NewSpreadsheet(execContext.frompath, execContext.topath)
	if err != nil {
		fatalError("on error occured loading xlsx file")
	}
	execContext.spreadsheet = sheet

	// -S option
	if execContext.initSheet != "" {
		execContext.scope.set("@", NewStringExpression(execContext.initSheet))
	}

	// -N option
	if execContext.doExcelRowLoop {
		execContext.code = "for(NER = SER; NER <= LR; NER++){ " + execContext.code + "; }"
	}

	// -n option
	if execContext.doTextRowLoop {
		execContext.code = "while(gets()){ " + execContext.code + "; }"
	}
}

func afterRun() {
	// request full calculate to excel
	execContext.spreadsheet.file.WorkBook.CalcPr.FullCalcOnLoad = true

	if execContext.topath != "" {
		if err := execContext.spreadsheet.writeSpreadsheet(); err != nil {
			fatalError("on error occured writting xlsx file")
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
		fmt.Fprintf(os.Stderr, "ERROR: "+format+"\n", a...)
	} else {
		fmt.Fprintf(os.Stderr, "ERROR: "+format+"\n")
	}
	os.Exit(1)
}
