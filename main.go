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

const CELL_VERSION = "0.0.1"

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
	con.scope.set("NR", NewNumberExpression(0))

	return con
}

func main() {
	con := NewExecContext()

	var pgpath string
	var showVer bool
	flag.StringVar(&con.topath, "to", "", "output spreadsheet filepath")
	flag.StringVar(&con.frompath, "from", "", "input spreadsheet filepath")
	flag.StringVar(&pgpath, "f", "", "program filepath")
	flag.BoolVar(&showVer, "V", false, "show version")

	flag.Parse()

	if showVer {
		showVersion()
	}

	args := flag.Args()
	if len(args) < 1 && pgpath == "" {
		flag.Usage()
		os.Exit(1)
	}

	if pgpath != "" {
		con.code = readProg(pgpath)
	} else {
		con.code = args[0]
	}

	if 1 < len(args) {
		switchStdin(con, args[1:])
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
		fmt.Fprintf(os.Stderr, "ERROR: "+format+"\n", a...)
	} else {
		fmt.Fprintf(os.Stderr, "ERROR: "+format+"\n")
	}
	os.Exit(1)
}
