package main

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
)

func setCellValue(t *testing.T, filepath string, sheet string, axis string, value interface{}) {
	f, err := excelize.OpenFile(filepath)
	if err != nil {
		t.Fatalf("on error occured open '%s'.", filepath)
	}
	err = f.SetCellValue(sheet, axis, value)
	if err != nil {
		t.Fatalf("on error occured set cell value '%s'.(%v)", axis, err)
	}
}

func getCellValue(t *testing.T, filepath string, sheet string, axis string) string {
	f, err := excelize.OpenFile(filepath)
	if err != nil {
		t.Fatalf("on error occured open '%s'.", filepath)
	}
	v, err := f.GetCellValue(sheet, axis)
	if err != nil {
		t.Fatalf("on error occured get cell value '%s'.(%v)", axis, err)
	}
	return v
}

func wrapStdio(t *testing.T, f func()) string {
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	stdout := os.Stdout
	os.Stdout = w

	f()

	os.Stdout = stdout
	w.Close()

	var buf bytes.Buffer
	io.Copy(&buf, r)

	return buf.String()
}

func TestSimpleNumberExpression(t *testing.T) {
	con := NewExecContext()
	con.code = `1`
	run(con)
	if con.exitCode != 0 {
		t.Fatalf("exit code '%s'. want '%d' but got '%d'", con.code, 0, con.exitCode)
	}
}

func TestSimpleStringExpression(t *testing.T) {
	con := NewExecContext()
	con.code = `"str"`
	run(con)
	if con.exitCode != 0 {
		t.Fatalf("exit code '%s'. want '%d' but got '%d'", con.code, 0, con.exitCode)
	}
}

func TestSimpleCellReferExpression(t *testing.T) {
	con := NewExecContext()
	con.frompath = "test/values.xlsx"
	con.code = `exit(["A1"]);`
	run(con)
	if con.exitCode != 2 {
		t.Fatalf("exit code '%s'. want '%d' but got '%d'", con.code, 2, con.exitCode)
	}
}

func TestSimpleCellAssignExpression(t *testing.T) {
	con := NewExecContext()
	con.frompath = "test/values.xlsx"
	con.topath = "TestSimpleCellAssignExpression.xlsx"
	con.code = `["A1"] = 5`
	run(con)
	if con.exitCode != 0 {
		t.Fatalf("exit code '%s'. want '%d' but got '%d'", con.code, 5, con.exitCode)
	}
	v := getCellValue(t, con.topath, "Sheet1", "A1")
	if v != "5" {
		t.Fatalf("want cell value 5, but got %s", v)
	}
}

func TestCellAssignToString(t *testing.T) {
	con := NewExecContext()
	con.frompath = "test/values.xlsx"
	con.topath = "TestCellAssignToString.xlsx"
	con.code = `["A1"] = "abc"`
	run(con)
	if con.exitCode != 0 {
		t.Fatalf("exit code '%s'. want '%d' but got '%d'", con.code, 0, con.exitCode)
	}
	v := getCellValue(t, con.topath, "Sheet1", "A1")
	if v != "abc" {
		t.Fatalf("want cell value 'abc', but got %s", v)
	}
}

func TestCellReferFromString(t *testing.T) {
	con := NewExecContext()
	con.frompath = "test/values.xlsx"
	con.topath = "TestCellReferFromString.xlsx"
	con.code = `["A3"] = ["A2"]`
	run(con)
	if con.exitCode != 0 {
		t.Fatalf("exit code '%s'. want '%d' but got '%d'", con.code, 0, con.exitCode)
	}
	v := getCellValue(t, con.topath, "Sheet1", "A3")
	if v != "test" {
		t.Fatalf("want cell value 'test', but got %s", v)
	}
}

func TestNumberAssignToVar(t *testing.T) {
	con := NewExecContext()
	con.code = `var = 10`
	run(con)
	if con.exitCode != 0 {
		t.Fatalf("exit code '%s'. want '%d' but got '%d'", con.code, 0, con.exitCode)
	}
}

func TestNumberVarRefer(t *testing.T) {
	con := NewExecContext()
	con.code = `var = 10;exit(var)`
	run(con)
	if con.exitCode != 10 {
		t.Fatalf("exit code '%s'. want '%d' but got '%d'", con.code, 10, con.exitCode)
	}
}

func TestStringAssignToVar(t *testing.T) {
	con := NewExecContext()
	con.topath = "TestStringAssignToVar.xlsx"
	con.code = `var = "test string";["A1"] = var;0`
	run(con)
	if con.exitCode != 0 {
		t.Fatalf("exit code '%s'. want '%d' but got '%d'", con.code, 0, con.exitCode)
	}
	v := getCellValue(t, con.topath, "Sheet1", "A1")
	if v != "test string" {
		t.Fatalf("want cell value 'test string', but got %s", v)
	}
}

func TestBuiltinPuts(t *testing.T) {
	out := new(bytes.Buffer)

	con := NewExecContext()
	con.out = out

	con.code = `puts("test string")`
	run(con)

	if con.exitCode != 0 {
		t.Fatalf("exit code '%s'. want '%d' but got '%d'", con.code, 0, con.exitCode)
	}
	if out.String() != "test string\n" {
		t.Fatalf("want stdout 'test string\n', but got %s", out)
	}
}

func TestBuiltinGets(t *testing.T) {
	in := bytes.NewBufferString("test string")
	out := new(bytes.Buffer)

	con := NewExecContext()
	con.in = in
	con.out = out

	con.code = `puts(gets())`
	run(con)

	if con.exitCode != 0 {
		t.Fatalf("exit code '%s'. want '%d' but got '%d'", con.code, 0, con.exitCode)
	}

	if out.String() != "test string\n" {
		t.Fatalf("want stdout 'test string\n', but got '%s'", out)
	}
}

func TestSpecialVarAtMarkRefer(t *testing.T) {
	out := new(bytes.Buffer)

	con := NewExecContext()
	con.out = out

	con.code = `puts(@)`
	run(con)

	if con.exitCode != 0 {
		t.Fatalf("exit code '%s'. want '%d', but got '%d'", con.code, 0, con.exitCode)
	}

	if out.String() != "Sheet1\n" {
		t.Fatalf("want stdout 'Sheet1\n', but got '%s'", out)
	}
}

func TestSpecialVarAtMarkAssign(t *testing.T) {
	out := new(bytes.Buffer)

	con := NewExecContext()
	con.out = out
	con.frompath = "test/values.xlsx"

	con.code = `@="Sheet2";puts(["A1"])`
	run(con)

	if con.exitCode != 0 {
		t.Fatalf("exit code '%s'. want '%d', but got '%d'", con.code, 0, con.exitCode)
	}

	if out.String() != "sheet2\n" {
		t.Fatalf("want stdout 'sheet2\n', but got '%s'", out)
	}
}

func TestSpecialVarAtMarkAssignUndefinedSheet(t *testing.T) {
	con := NewExecContext()
	con.topath = "TestSpecialVarAStMarkAssignUndefinedSheet.xlsx"

	con.code = `@="foo";["A1"] = "new sheet"`
	run(con)

	if con.exitCode != 0 {
		t.Fatalf("exit code '%s'. want '%d', but got '%d'", con.code, 0, con.exitCode)
	}

	actual := getCellValue(t, con.topath, "foo", "A1")
	if actual != "new sheet" {
		t.Fatalf("no new sheet has been added or no value has been set.")
	}
}

func TestSpecialVarDollarDefault(t *testing.T) {
	expect := "aa bb  cc\t \tstring"
	out := new(bytes.Buffer)
	in := bytes.NewBufferString(expect)

	con := NewExecContext()
	con.topath = "TestSpecialVarDollarDefault.xlsx"
	con.out = out
	con.in = in
	con.code = `gets();["A1"]=$0;["A2"]=$1;["A3"]=$2;["A4"]=$3;["A5"]=$4;`
	run(con)

	if con.exitCode != 0 {
		t.Fatalf("exit code '%s'. want '%d', but got '%d'", con.code, 0, con.exitCode)
	}

	actual := getCellValue(t, con.topath, "Sheet1", "A1")
	if actual != expect {
		t.Fatalf("$0 value wrong. want '%s', but got '%s'", expect, actual)
	}

	actual = getCellValue(t, con.topath, "Sheet1", "A2")
	if actual != "aa" {
		t.Fatalf("$1 value wrong. want '%s', but got '%s'", "aa", actual)
	}

	actual = getCellValue(t, con.topath, "Sheet1", "A3")
	if actual != "bb" {
		t.Fatalf("$2 value wrong. want '%s', but got '%s'", "bb", actual)
	}

	actual = getCellValue(t, con.topath, "Sheet1", "A4")
	if actual != "cc" {
		t.Fatalf("$3 value wrong. want '%s', but got '%s'", "cc", actual)
	}

	actual = getCellValue(t, con.topath, "Sheet1", "A5")
	if actual != "string" {
		t.Fatalf("$4 value wrong. want '%s', but got '%s'", "string", actual)
	}
}

func TestSpecialVarDollarOneChar(t *testing.T) {
	expect := "aa.bb..cc.test string"
	out := new(bytes.Buffer)
	in := bytes.NewBufferString(expect)

	con := NewExecContext()
	con.topath = "TestSpecialVarDollarOneChar.xlsx"
	con.out = out
	con.in = in
	con.code = `FS=".";gets();["A1"]=$0;["A2"]=$1;["A3"]=$2;["A4"]=$3;["A5"]=$4;["A6"]=$5`
	run(con)

	if con.exitCode != 0 {
		t.Fatalf("exit code '%s'. want '%d', but got '%d'", con.code, 0, con.exitCode)
	}

	actual := getCellValue(t, con.topath, "Sheet1", "A1")
	if actual != expect {
		t.Fatalf("$0 value wrong. want '%s', but got '%s'", expect, actual)
	}

	actual = getCellValue(t, con.topath, "Sheet1", "A2")
	if actual != "aa" {
		t.Fatalf("$1 value wrong. want '%s', but got '%s'", "aa", actual)
	}

	actual = getCellValue(t, con.topath, "Sheet1", "A3")
	if actual != "bb" {
		t.Fatalf("$2 value wrong. want '%s', but got '%s'", "bb", actual)
	}

	actual = getCellValue(t, con.topath, "Sheet1", "A4")
	if actual != "" {
		t.Fatalf("$3 value wrong. want '%s', but got '%s'", "", actual)
	}

	actual = getCellValue(t, con.topath, "Sheet1", "A5")
	if actual != "cc" {
		t.Fatalf("$4 value wrong. want '%s', but got '%s'", "cc", actual)
	}

	actual = getCellValue(t, con.topath, "Sheet1", "A6")
	if actual != "test string" {
		t.Fatalf("$5 value wrong. want '%s', but got '%s'", "test string", actual)
	}
}

func TestSpecialVarDollarRegexp(t *testing.T) {
	expect := "aa,bb,  cc|string"
	out := new(bytes.Buffer)
	in := bytes.NewBufferString(expect)

	con := NewExecContext()
	con.topath = "TestSpecialVarDollarRegexp.xlsx"
	con.out = out
	con.in = in
	con.code = `FS="[,|] *";gets();["A1"]=$0;["A2"]=$1;["A3"]=$2;["A4"]=$3;["A5"]=$4;`
	run(con)

	if con.exitCode != 0 {
		t.Fatalf("exit code '%s'. want '%d', but got '%d'", con.code, 0, con.exitCode)
	}

	actual := getCellValue(t, con.topath, "Sheet1", "A1")
	if actual != expect {
		t.Fatalf("$0 value wrong. want '%s', but got '%s'", expect, actual)
	}

	actual = getCellValue(t, con.topath, "Sheet1", "A2")
	if actual != "aa" {
		t.Fatalf("$1 value wrong. want '%s', but got '%s'", "aa", actual)
	}

	actual = getCellValue(t, con.topath, "Sheet1", "A3")
	if actual != "bb" {
		t.Fatalf("$2 value wrong. want '%s', but got '%s'", "bb", actual)
	}

	actual = getCellValue(t, con.topath, "Sheet1", "A4")
	if actual != "cc" {
		t.Fatalf("$3 value wrong. want '%s', but got '%s'", "cc", actual)
	}

	actual = getCellValue(t, con.topath, "Sheet1", "A5")
	if actual != "string" {
		t.Fatalf("$4 value wrong. want '%s', but got '%s'", "string", actual)
	}
}

func TestBuiltinFuncPutsNoArg(t *testing.T) {
	expect := "aa bb cc"
	out := new(bytes.Buffer)
	in := bytes.NewBufferString(expect)

	con := NewExecContext()
	con.out = out
	con.in = in
	con.code = `gets();puts();`
	run(con)

	if con.exitCode != 0 {
		t.Fatalf("exit code '%s'. want '%d', but got '%d'", con.code, 0, con.exitCode)
	}

	actual := out.String()
	if actual != expect+"\n" {
		t.Fatalf("puts() want '%s', but got '%s'", expect+"\n", actual)
	}
}

func TestSpecialVarOFSAndBuiltinPutsFuncMultiArgs(t *testing.T) {
	src := "aa bb cc"
	out := new(bytes.Buffer)
	in := bytes.NewBufferString(src)

	con := NewExecContext()
	con.out = out
	con.in = in
	con.code = `OFS="  ";gets();puts($1, $3);`
	run(con)

	if con.exitCode != 0 {
		t.Fatalf("exit code '%s'. want '%d', but got '%d'", con.code, 0, con.exitCode)
	}

	expect := "aa  cc\n"
	actual := out.String()
	if actual != expect {
		t.Fatalf("puts($1, $3) want '%s', but got '%s'", expect, actual)
	}
}

func TestEscapeString(t *testing.T) {
	src := "aa bb cc"
	out := new(bytes.Buffer)
	in := bytes.NewBufferString(src)

	con := NewExecContext()
	con.out = out
	con.in = in
	con.code = `OFS="\t";gets();puts($1, $3);`
	run(con)

	if con.exitCode != 0 {
		t.Fatalf("exit code '%s'. want '%d', but got '%d'", con.code, 0, con.exitCode)
	}

	expect := "aa\tcc\n"
	actual := out.String()
	if actual != expect {
		t.Fatalf("puts($1, $3) want '%s', but got '%s'", expect, actual)
	}
}

func TestNumberEQExpression(t *testing.T) {
	con := NewExecContext()
	con.topath = "TestNumberEQExpression.xlsx"
	con.code = `["A1"] = 1==1;["A2"] = 2==1;`
	run(con)

	if con.exitCode != 0 {
		t.Fatalf("exit code '%s'. want '%d', but got '%d'", con.code, 0, con.exitCode)
	}

	v := getCellValue(t, con.topath, "Sheet1", "A1")
	if v != "1" {
		t.Fatalf("want cell value '1', but got %s", v)
	}
	v = getCellValue(t, con.topath, "Sheet1", "A2")
	if v != "0" {
		t.Fatalf("want cell value '0', but got %s", v)
	}
}

func TestNumberNEExpression(t *testing.T) {
	con := NewExecContext()
	con.topath = "TestNumberNEExpression.xlsx"
	con.code = `["A1"] = 1!=1;["A2"] = 2!=1;`
	run(con)

	if con.exitCode != 0 {
		t.Fatalf("exit code '%s'. want '%d', but got '%d'", con.code, 0, con.exitCode)
	}

	v := getCellValue(t, con.topath, "Sheet1", "A1")
	if v != "0" {
		t.Fatalf("want cell value '1', but got %s", v)
	}
	v = getCellValue(t, con.topath, "Sheet1", "A2")
	if v != "1" {
		t.Fatalf("want cell value '0', but got %s", v)
	}
}

func TestNumberLTExpression(t *testing.T) {
	con := NewExecContext()
	con.topath = "TestNumberLTExpression.xlsx"
	con.code = `["A1"] = 1<1;["A2"] = 0<1;`
	run(con)

	if con.exitCode != 0 {
		t.Fatalf("exit code '%s'. want '%d', but got '%d'", con.code, 0, con.exitCode)
	}

	v := getCellValue(t, con.topath, "Sheet1", "A1")
	if v != "0" {
		t.Fatalf("want cell value '1', but got %s", v)
	}
	v = getCellValue(t, con.topath, "Sheet1", "A2")
	if v != "1" {
		t.Fatalf("want cell value '0', but got %s", v)
	}
}

func TestNumberLEExpression(t *testing.T) {
	con := NewExecContext()
	con.topath = "TestNumberLEExpression.xlsx"
	con.code = `["A1"] = 1<=1;["A2"] = 2<=1;`
	run(con)

	if con.exitCode != 0 {
		t.Fatalf("exit code '%s'. want '%d', but got '%d'", con.code, 0, con.exitCode)
	}

	v := getCellValue(t, con.topath, "Sheet1", "A1")
	if v != "1" {
		t.Fatalf("want cell value '1', but got %s", v)
	}
	v = getCellValue(t, con.topath, "Sheet1", "A2")
	if v != "0" {
		t.Fatalf("want cell value '0', but got %s", v)
	}
}

func TestNumberGTExpression(t *testing.T) {
	con := NewExecContext()
	con.topath = "TestNumberGTExpression.xlsx"
	con.code = `["A1"] = 1>1;["A2"] = 1>0;`
	run(con)

	if con.exitCode != 0 {
		t.Fatalf("exit code '%s'. want '%d', but got '%d'", con.code, 0, con.exitCode)
	}

	v := getCellValue(t, con.topath, "Sheet1", "A1")
	if v != "0" {
		t.Fatalf("want cell value '1', but got %s", v)
	}
	v = getCellValue(t, con.topath, "Sheet1", "A2")
	if v != "1" {
		t.Fatalf("want cell value '0', but got %s", v)
	}
}

func TestNumberGEExpression(t *testing.T) {
	con := NewExecContext()
	con.topath = "TestNumberGEExpression.xlsx"
	con.code = `["A1"] = 1>=1;["A2"] = 1>=2;`
	run(con)

	if con.exitCode != 0 {
		t.Fatalf("exit code '%s'. want '%d', but got '%d'", con.code, 0, con.exitCode)
	}

	v := getCellValue(t, con.topath, "Sheet1", "A1")
	if v != "1" {
		t.Fatalf("want cell value '1', but got %s", v)
	}
	v = getCellValue(t, con.topath, "Sheet1", "A2")
	if v != "0" {
		t.Fatalf("want cell value '0', but got %s", v)
	}
}

func TestStringEQExpression(t *testing.T) {
	con := NewExecContext()
	con.topath = "TestStringEQExpression.xlsx"
	con.code = `["A1"] = "hello" eq "hello";["A2"] = "hello" eq "bye";`
	run(con)

	if con.exitCode != 0 {
		t.Fatalf("exit code '%s'. want '%d', but got '%d'", con.code, 0, con.exitCode)
	}

	v := getCellValue(t, con.topath, "Sheet1", "A1")
	if v != "1" {
		t.Fatalf("want cell value '1', but got %s", v)
	}
	v = getCellValue(t, con.topath, "Sheet1", "A2")
	if v != "0" {
		t.Fatalf("want cell value '0', but got %s", v)
	}
}

func TestStringNEExpression(t *testing.T) {
	con := NewExecContext()
	con.topath = "TestStringNEExpression.xlsx"
	con.code = `["A1"] = "hello" ne "hello";["A2"] = "hello" ne "bye";`
	run(con)

	if con.exitCode != 0 {
		t.Fatalf("exit code '%s'. want '%d', but got '%d'", con.code, 0, con.exitCode)
	}

	v := getCellValue(t, con.topath, "Sheet1", "A1")
	if v != "0" {
		t.Fatalf("want cell value '0', but got %s", v)
	}
	v = getCellValue(t, con.topath, "Sheet1", "A2")
	if v != "1" {
		t.Fatalf("want cell value '1', but got %s", v)
	}
}

func TestStringConcatExpression(t *testing.T) {
	con := NewExecContext()
	con.topath = "TestStringConcatExpression.xlsx"
	con.code = `["A1"] = "hello"." world"`
	run(con)

	if con.exitCode != 0 {
		t.Fatalf("exit code '%s'. want '%d', but got '%d'", con.code, 0, con.exitCode)
	}

	v := getCellValue(t, con.topath, "Sheet1", "A1")
	if v != "hello world" {
		t.Fatalf("want cell value 'hello world', but got %s", v)
	}
}

func TestNumberAddExpression(t *testing.T) {
	con := NewExecContext()
	con.topath = "TestNumberAddExpression.xlsx"
	con.code = `["A1"] = 1+3`
	run(con)

	if con.exitCode != 0 {
		t.Fatalf("exit code '%s'. want '%d', but got '%d'", con.code, 0, con.exitCode)
	}

	v := getCellValue(t, con.topath, "Sheet1", "A1")
	if v != "4" {
		t.Fatalf("want cell value '4', but got %s", v)
	}
}
