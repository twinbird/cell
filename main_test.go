package main

import (
	"bufio"
	"bytes"
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

func TestCellAssignNumberString(t *testing.T) {
	con := NewExecContext()
	con.topath = "TestCellAssignNumberString.xlsx"
	con.code = `["A1"] = "5"`
	run(con)
	if con.exitCode != 0 {
		t.Fatalf("exit code '%s'. want '%d' but got '%d'", con.code, 0, con.exitCode)
	}
	v := getCellValue(t, con.topath, "Sheet1", "A1")
	if v != "5" {
		t.Fatalf("want cell value '5', but got %s", v)
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
	in := bufio.NewReader(bytes.NewBufferString("test string"))
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
	in := bufio.NewReader(bytes.NewBufferString(expect))

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
	in := bufio.NewReader(bytes.NewBufferString(expect))

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
	in := bufio.NewReader(bytes.NewBufferString(expect))

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
	in := bufio.NewReader(bytes.NewBufferString(expect))

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
	in := bufio.NewReader(bytes.NewBufferString(src))

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
	in := bufio.NewReader(bytes.NewBufferString(src))

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

func TestNumberSubExpression(t *testing.T) {
	con := NewExecContext()
	con.topath = "TestNumberSubExpression.xlsx"
	con.code = `["A1"] = 1-3`
	run(con)

	if con.exitCode != 0 {
		t.Fatalf("exit code '%s'. want '%d', but got '%d'", con.code, 0, con.exitCode)
	}

	v := getCellValue(t, con.topath, "Sheet1", "A1")
	if v != "-2" {
		t.Fatalf("want cell value '-2', but got %s", v)
	}
}

func TestNumberDivExpression(t *testing.T) {
	con := NewExecContext()
	con.topath = "TestNumberDivExpression.xlsx"
	con.code = `["A1"] = 9/3`
	run(con)

	if con.exitCode != 0 {
		t.Fatalf("exit code '%s'. want '%d', but got '%d'", con.code, 0, con.exitCode)
	}

	v := getCellValue(t, con.topath, "Sheet1", "A1")
	if v != "3" {
		t.Fatalf("want cell value '3', but got %s", v)
	}
}

func TestNumberModuloExpression(t *testing.T) {
	con := NewExecContext()
	con.topath = "TestNumberModuloExpression.xlsx"
	con.code = `["A1"] = 10%3;["A2"]=9.99%3.33`
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

func TestNumberPowerExpression(t *testing.T) {
	con := NewExecContext()
	con.topath = "TestNumberPowerExpression.xlsx"
	con.code = `["A1"] = 2**3;["A2"]=3**0`
	run(con)

	if con.exitCode != 0 {
		t.Fatalf("exit code '%s'. want '%d', but got '%d'", con.code, 0, con.exitCode)
	}

	v := getCellValue(t, con.topath, "Sheet1", "A1")
	if v != "8" {
		t.Fatalf("want cell value '8', but got %s", v)
	}
	v = getCellValue(t, con.topath, "Sheet1", "A2")
	if v != "1" {
		t.Fatalf("want cell value '1', but got %s", v)
	}
}

func TestLogicalAndExpression(t *testing.T) {
	con := NewExecContext()
	con.topath = "TestLogicalAndExpression.xlsx"
	con.code = `["A1"] = 1 && 0;["A2"]="" && 1;["A3"]="a"&&1`
	run(con)

	if con.exitCode != 0 {
		t.Fatalf("exit code '%s'. want '%d', but got '%d'", con.code, 0, con.exitCode)
	}

	v := getCellValue(t, con.topath, "Sheet1", "A1")
	if v != "0" {
		t.Fatalf("want cell value '0', but got %s", v)
	}
	v = getCellValue(t, con.topath, "Sheet1", "A2")
	if v != "0" {
		t.Fatalf("want cell value '0', but got %s", v)
	}
	v = getCellValue(t, con.topath, "Sheet1", "A3")
	if v != "1" {
		t.Fatalf("want cell value '1', but got %s", v)
	}
}

func TestLogicalOrExpression(t *testing.T) {
	con := NewExecContext()
	con.topath = "TestLogicalOrExpression.xlsx"
	con.code = `["A1"] = 1 || 0;["A2"]="" || 1;["A3"]="a"||1;["A4"]="" || 0`
	run(con)

	if con.exitCode != 0 {
		t.Fatalf("exit code '%s'. want '%d', but got '%d'", con.code, 0, con.exitCode)
	}

	v := getCellValue(t, con.topath, "Sheet1", "A1")
	if v != "1" {
		t.Fatalf("want cell value '1', but got %s", v)
	}
	v = getCellValue(t, con.topath, "Sheet1", "A2")
	if v != "1" {
		t.Fatalf("want cell value '1', but got %s", v)
	}
	v = getCellValue(t, con.topath, "Sheet1", "A3")
	if v != "1" {
		t.Fatalf("want cell value '1', but got %s", v)
	}
	v = getCellValue(t, con.topath, "Sheet1", "A4")
	if v != "0" {
		t.Fatalf("want cell value '0', but got %s", v)
	}
}

func TestLogicalNotExpression(t *testing.T) {
	con := NewExecContext()
	con.topath = "TestLogicalNotExpression.xlsx"
	con.code = `["A1"] = !1;["A2"]=!0;["A3"]=!"a";["A4"]=!""`
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
	v = getCellValue(t, con.topath, "Sheet1", "A3")
	if v != "0" {
		t.Fatalf("want cell value '0', but got %s", v)
	}
	v = getCellValue(t, con.topath, "Sheet1", "A4")
	if v != "1" {
		t.Fatalf("want cell value '1', but got %s", v)
	}
}

func TestParenthesesOperatorExpression(t *testing.T) {
	con := NewExecContext()
	con.topath = "TestParenthesesOperatorExpression.xlsx"
	con.code = `["A1"] = (0+1)&&1;["A2"]=(1+3)*2;`
	run(con)

	if con.exitCode != 0 {
		t.Fatalf("exit code '%s'. want '%d', but got '%d'", con.code, 0, con.exitCode)
	}

	v := getCellValue(t, con.topath, "Sheet1", "A1")
	if v != "1" {
		t.Fatalf("want cell value '1', but got %s", v)
	}
	v = getCellValue(t, con.topath, "Sheet1", "A2")
	if v != "8" {
		t.Fatalf("want cell value '8', but got %s", v)
	}
}

func TestMinusExpression(t *testing.T) {
	con := NewExecContext()
	con.topath = "TestMinusExpression.xlsx"
	con.code = `["A1"] = -1; ["A2"]=-(-1)`
	run(con)

	if con.exitCode != 0 {
		t.Fatalf("exit code '%s'. want '%d', but got '%d'", con.code, 0, con.exitCode)
	}

	v := getCellValue(t, con.topath, "Sheet1", "A1")
	if v != "-1" {
		t.Fatalf("want cell value '-1', but got %s", v)
	}
	v = getCellValue(t, con.topath, "Sheet1", "A2")
	if v != "1" {
		t.Fatalf("want cell value '1', but got %s", v)
	}
}

func TestPlusExpression(t *testing.T) {
	con := NewExecContext()
	con.topath = "TestPlusExpression.xlsx"
	con.code = `["A1"] = +1; ["A2"]=+"a"`
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

func TestAddAssignExpression(t *testing.T) {
	con := NewExecContext()
	con.topath = "TestAddAssignExpression.xlsx"
	con.code = `a = 10; a += 5; ["A1"]=a; b+=10;["A2"]=b;`
	run(con)

	if con.exitCode != 0 {
		t.Fatalf("exit code '%s'. want '%d', but got '%d'", con.code, 0, con.exitCode)
	}

	v := getCellValue(t, con.topath, "Sheet1", "A1")
	if v != "15" {
		t.Fatalf("want cell value '15', but got %s", v)
	}
	v = getCellValue(t, con.topath, "Sheet1", "A2")
	if v != "10" {
		t.Fatalf("want cell value '10', but got %s", v)
	}
}

func TestSubAssignExpression(t *testing.T) {
	con := NewExecContext()
	con.topath = "TestSubAssignExpression.xlsx"
	con.code = `a = 10; a -= 5; ["A1"]=a; b-=10;["A2"]=b;`
	run(con)

	if con.exitCode != 0 {
		t.Fatalf("exit code '%s'. want '%d', but got '%d'", con.code, 0, con.exitCode)
	}

	v := getCellValue(t, con.topath, "Sheet1", "A1")
	if v != "5" {
		t.Fatalf("want cell value '5', but got %s", v)
	}
	v = getCellValue(t, con.topath, "Sheet1", "A2")
	if v != "-10" {
		t.Fatalf("want cell value '-10', but got %s", v)
	}
}

func TestMulAssignExpression(t *testing.T) {
	con := NewExecContext()
	con.topath = "TestMulAssignExpression.xlsx"
	con.code = `a = 10; a *= 5; ["A1"]=a; b*=10;["A2"]=b;`
	run(con)

	if con.exitCode != 0 {
		t.Fatalf("exit code '%s'. want '%d', but got '%d'", con.code, 0, con.exitCode)
	}

	v := getCellValue(t, con.topath, "Sheet1", "A1")
	if v != "50" {
		t.Fatalf("want cell value '50', but got %s", v)
	}
	v = getCellValue(t, con.topath, "Sheet1", "A2")
	if v != "0" {
		t.Fatalf("want cell value '0', but got %s", v)
	}
}

func TestDivAssignExpression(t *testing.T) {
	con := NewExecContext()
	con.topath = "TestDivAssignExpression.xlsx"
	con.code = `a = 10; a /= 5; ["A1"]=a; b/=10;["A2"]=b;`
	run(con)

	if con.exitCode != 0 {
		t.Fatalf("exit code '%s'. want '%d', but got '%d'", con.code, 0, con.exitCode)
	}

	v := getCellValue(t, con.topath, "Sheet1", "A1")
	if v != "2" {
		t.Fatalf("want cell value '2', but got %s", v)
	}
	v = getCellValue(t, con.topath, "Sheet1", "A2")
	if v != "0" {
		t.Fatalf("want cell value '0', but got %s", v)
	}
}

func TestModAssignExpression(t *testing.T) {
	con := NewExecContext()
	con.topath = "TestModAssignExpression.xlsx"
	con.code = `a = 11; a %= 5; ["A1"]=a; b%=10;["A2"]=b;`
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

func TestPowAssignExpression(t *testing.T) {
	con := NewExecContext()
	con.topath = "TestPowAssignExpression.xlsx"
	con.code = `a = 2; a **= 3; ["A1"]=a; b**=10;["A2"]=b;`
	run(con)

	if con.exitCode != 0 {
		t.Fatalf("exit code '%s'. want '%d', but got '%d'", con.code, 0, con.exitCode)
	}

	v := getCellValue(t, con.topath, "Sheet1", "A1")
	if v != "8" {
		t.Fatalf("want cell value '8', but got %s", v)
	}
	v = getCellValue(t, con.topath, "Sheet1", "A2")
	if v != "0" {
		t.Fatalf("want cell value '0', but got %s", v)
	}
}

func TestConcatAssignExpression(t *testing.T) {
	con := NewExecContext()
	con.topath = "TestConcatAssignExpression.xlsx"
	con.code = `a = "Hello, "; a .= "world"; ["A1"]=a;`
	run(con)

	if con.exitCode != 0 {
		t.Fatalf("exit code '%s'. want '%d', but got '%d'", con.code, 0, con.exitCode)
	}

	v := getCellValue(t, con.topath, "Sheet1", "A1")
	if v != "Hello, world" {
		t.Fatalf("concat assign could not working. want cell value 'Hello, world', but got %s", v)
	}
}

func TestComment(t *testing.T) {
	out := new(bytes.Buffer)
	con := NewExecContext()
	con.out = out
	con.code = `
	# this is a comment
	a=10 # and also 
	# puts(a) <-no output
	puts(a)
	`
	run(con)

	if out.String() != "10\n" {
		t.Fatalf("comments are not working")
	}
}

func TestMatchExpression(t *testing.T) {
	con := NewExecContext()
	con.topath = "TestMatchExpression.xlsx"
	con.code = `["A1"] = "Hello, world" ~ "Hell(o)?";["A2"]="Hello, world" ~ "foo"`
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

func TestMatchSpecialVar(t *testing.T) {
	con := NewExecContext()
	con.topath = "TestMatchSpecialVar.xlsx"
	con.code = `"Hello, world" ~ "Hell(o)?";["A1"] = $_0;["A2"]=$_1;`
	run(con)

	if con.exitCode != 0 {
		t.Fatalf("exit code '%s'. want '%d', but got '%d'", con.code, 0, con.exitCode)
	}

	v := getCellValue(t, con.topath, "Sheet1", "A1")
	if v != "Hello" {
		t.Fatalf("special variable $_0 dont working. want '%s', but got '%s'", "Hello", v)
	}
	v = getCellValue(t, con.topath, "Sheet1", "A2")
	if v != "o" {
		t.Fatalf("special variable $_0 dont working. want '%s', but got '%s'", "o", v)
	}
}

func TestNotMatchExpression(t *testing.T) {
	con := NewExecContext()
	con.topath = "TestNotMatchExpression.xlsx"
	con.code = `["A1"] = "Hello, world" !~ "Hell(o)?";["A2"]="Hello, world" !~ "foo"`
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

func TestMatchSpecialVarWhenNotMatch(t *testing.T) {
	con := NewExecContext()
	con.topath = "TestMatchSpecialVarWhenNotMatch.xlsx"
	con.code = `"Hello, world" !~ "Hell(o)?";["A1"] = $_0;["A2"]=$_1;`
	run(con)

	if con.exitCode != 0 {
		t.Fatalf("exit code '%s'. want '%d', but got '%d'", con.code, 0, con.exitCode)
	}

	v := getCellValue(t, con.topath, "Sheet1", "A1")
	if v != "Hello" {
		t.Fatalf("special variable $_0 could not working. want '%s', but got '%s'", "Hello", v)
	}
	v = getCellValue(t, con.topath, "Sheet1", "A2")
	if v != "o" {
		t.Fatalf("special variable $_0 could not working. want '%s', but got '%s'", "o", v)
	}
}

func TestIfStatement(t *testing.T) {
	con := NewExecContext()
	con.topath = "TestIfStatement.xlsx"
	con.code = `["A1"]=1;if(0)["A1"]=2;["A2"]=1;if(1)["A2"]=2;`
	run(con)

	if con.exitCode != 0 {
		t.Fatalf("exit code '%s'. want '%d', but got '%d'", con.code, 0, con.exitCode)
	}

	v := getCellValue(t, con.topath, "Sheet1", "A1")
	if v != "1" {
		t.Fatalf("if statement could not working")
	}
	v = getCellValue(t, con.topath, "Sheet1", "A2")
	if v != "2" {
		t.Fatalf("if then statement could not working")
	}
}

func TestIfElseStatement(t *testing.T) {
	con := NewExecContext()
	con.topath = "TestIfElseStatement.xlsx"
	con.code = `["A1"]=0;if(0)["A1"]=2;else["A1"]=1;if(1)["A2"]=2;else["A2"]=3;`
	run(con)

	if con.exitCode != 0 {
		t.Fatalf("exit code '%s'. want '%d', but got '%d'", con.code, 0, con.exitCode)
	}

	v := getCellValue(t, con.topath, "Sheet1", "A1")
	if v != "1" {
		t.Fatalf("else statement could not working")
	}
	v = getCellValue(t, con.topath, "Sheet1", "A2")
	if v != "2" {
		t.Fatalf("else then statement could not working")
	}
}

func TestBlockStatement(t *testing.T) {
	con := NewExecContext()
	con.topath = "TestBlockStatement.xlsx"
	con.code = `if(0){["A1"]="hello";["A2"]="world";}else{["A1"]="Bye";["A2"]="bye";}`
	run(con)

	if con.exitCode != 0 {
		t.Fatalf("exit code '%s'. want '%d', but got '%d'", con.code, 0, con.exitCode)
	}

	v := getCellValue(t, con.topath, "Sheet1", "A1")
	if v != "Bye" {
		t.Fatalf("else block statement could not working")
	}
	v = getCellValue(t, con.topath, "Sheet1", "A2")
	if v != "bye" {
		t.Fatalf("else block statement could not working")
	}
}

func TestWhileStatement(t *testing.T) {
	con := NewExecContext()
	con.topath = "TestWhileStatement.xlsx"
	con.code = `while(i<10){sum+=i;i+=1;}["A1"]=sum`
	run(con)

	if con.exitCode != 0 {
		t.Fatalf("exit code '%s'. want '%d', but got '%d'", con.code, 0, con.exitCode)
	}

	v := getCellValue(t, con.topath, "Sheet1", "A1")
	if v != "45" {
		t.Fatalf("while statement could not working. want '45', but got %v", v)
	}
}

func TestCalculatedCellAssignToString(t *testing.T) {
	con := NewExecContext()
	con.topath = "TestCalculatedCellAssignToString.xlsx"
	con.code = `["A"."1"] = "abc"`
	run(con)
	if con.exitCode != 0 {
		t.Fatalf("exit code '%s'. want '%d' but got '%d'", con.code, 0, con.exitCode)
	}
	v := getCellValue(t, con.topath, "Sheet1", "A1")
	if v != "abc" {
		t.Fatalf("want cell value 'abc', but got %s", v)
	}
}

func TestAddAndCellAssign(t *testing.T) {
	con := NewExecContext()
	con.topath = "TestAddAndCellAssign.xlsx"
	con.code = `["A1"] = 1;["A1"]+=3;`
	run(con)
	if con.exitCode != 0 {
		t.Fatalf("exit code '%s'. want '%d' but got '%d'", con.code, 0, con.exitCode)
	}
	v := getCellValue(t, con.topath, "Sheet1", "A1")
	if v != "4" {
		t.Fatalf("want cell value '4', but got %s", v)
	}
}

func TestSubAndCellAssign(t *testing.T) {
	con := NewExecContext()
	con.topath = "TestSubAndCellAssign.xlsx"
	con.code = `["A1"] = 1;["A1"]-=3;`
	run(con)
	if con.exitCode != 0 {
		t.Fatalf("exit code '%s'. want '%d' but got '%d'", con.code, 0, con.exitCode)
	}
	v := getCellValue(t, con.topath, "Sheet1", "A1")
	if v != "-2" {
		t.Fatalf("want cell value '-2', but got %s", v)
	}
}

func TestMulAndCellAssign(t *testing.T) {
	con := NewExecContext()
	con.topath = "TestMulAndCellAssign.xlsx"
	con.code = `["A1"] = 2;["A1"]*=3;`
	run(con)
	if con.exitCode != 0 {
		t.Fatalf("exit code '%s'. want '%d' but got '%d'", con.code, 0, con.exitCode)
	}
	v := getCellValue(t, con.topath, "Sheet1", "A1")
	if v != "6" {
		t.Fatalf("want cell value '6', but got %s", v)
	}
}

func TestDivAndCellAssign(t *testing.T) {
	con := NewExecContext()
	con.topath = "TestDivAndCellAssign.xlsx"
	con.code = `["A1"] = 9;["A1"]/=3;`
	run(con)
	if con.exitCode != 0 {
		t.Fatalf("exit code '%s'. want '%d' but got '%d'", con.code, 0, con.exitCode)
	}
	v := getCellValue(t, con.topath, "Sheet1", "A1")
	if v != "3" {
		t.Fatalf("want cell value '3', but got %s", v)
	}
}

func TestModAndCellAssign(t *testing.T) {
	con := NewExecContext()
	con.topath = "TestModAndCellAssign.xlsx"
	con.code = `["A1"] = 10;["A1"]%=3;`
	run(con)
	if con.exitCode != 0 {
		t.Fatalf("exit code '%s'. want '%d' but got '%d'", con.code, 0, con.exitCode)
	}
	v := getCellValue(t, con.topath, "Sheet1", "A1")
	if v != "1" {
		t.Fatalf("want cell value '1', but got %s", v)
	}
}

func TestPowAndCellAssign(t *testing.T) {
	con := NewExecContext()
	con.topath = "TestPowAndCellAssign.xlsx"
	con.code = `["A1"] = 10;["A1"]**=2;`
	run(con)
	if con.exitCode != 0 {
		t.Fatalf("exit code '%s'. want '%d' but got '%d'", con.code, 0, con.exitCode)
	}
	v := getCellValue(t, con.topath, "Sheet1", "A1")
	if v != "100" {
		t.Fatalf("want cell value '100', but got %s", v)
	}
}

func TestConcatAndCellAssign(t *testing.T) {
	con := NewExecContext()
	con.topath = "TestConcatAndCellAssign.xlsx"
	con.code = `["A1"] = "Hello";["A1"].=" world";`
	run(con)
	if con.exitCode != 0 {
		t.Fatalf("exit code '%s'. want '%d' but got '%d'", con.code, 0, con.exitCode)
	}
	v := getCellValue(t, con.topath, "Sheet1", "A1")
	if v != "Hello world" {
		t.Fatalf("want cell value 'Hello world', but got %s", v)
	}
}

func TestBreakStatement(t *testing.T) {
	con := NewExecContext()
	con.topath = "TestBreakStatement.xlsx"
	con.code = `while(i<10){sum+=i;if(i==3)break;i+=1;}["A1"]=sum`
	run(con)
	if con.exitCode != 0 {
		t.Fatalf("exit code '%s'. want '%d' but got '%d'", con.code, 0, con.exitCode)
	}
	v := getCellValue(t, con.topath, "Sheet1", "A1")
	if v != "6" {
		t.Fatalf("want cell value '6', but got %s", v)
	}
}

func TestContinueStatement(t *testing.T) {
	con := NewExecContext()
	con.topath = "TestContinueStatement.xlsx"
	con.code = `while(i<10){i+=1;if(i==3)continue;sum+=i;}["A1"]=sum`
	run(con)
	if con.exitCode != 0 {
		t.Fatalf("exit code '%s'. want '%d' but got '%d'", con.code, 0, con.exitCode)
	}
	v := getCellValue(t, con.topath, "Sheet1", "A1")
	if v != "52" {
		t.Fatalf("want cell value '52', but got %s", v)
	}
}

func TestIncrementExpression(t *testing.T) {
	con := NewExecContext()
	con.topath = "TestIncrementExpression.xlsx"
	con.code = `a=0;if(a++)a=100;["A1"]=a;`
	run(con)
	if con.exitCode != 0 {
		t.Fatalf("exit code '%s'. want '%d' but got '%d'", con.code, 0, con.exitCode)
	}
	v := getCellValue(t, con.topath, "Sheet1", "A1")
	if v != "1" {
		t.Fatalf("want cell value '1', but got %s", v)
	}
}

func TestIncrementCellExpression(t *testing.T) {
	con := NewExecContext()
	con.topath = "TestIncrementCellExpression.xlsx"
	con.code = `["A1"]=0;if(["A1"]++)["A1"]=100;`
	run(con)
	if con.exitCode != 0 {
		t.Fatalf("exit code '%s'. want '%d' but got '%d'", con.code, 0, con.exitCode)
	}
	v := getCellValue(t, con.topath, "Sheet1", "A1")
	if v != "1" {
		t.Fatalf("want cell value '1', but got %s", v)
	}
}

func TestDecrementExpression(t *testing.T) {
	con := NewExecContext()
	con.topath = "TestDecrementExpression.xlsx"
	con.code = `a=0;if(a--)a=100;["A1"]=a;`
	run(con)
	if con.exitCode != 0 {
		t.Fatalf("exit code '%s'. want '%d' but got '%d'", con.code, 0, con.exitCode)
	}
	v := getCellValue(t, con.topath, "Sheet1", "A1")
	if v != "-1" {
		t.Fatalf("want cell value '-1', but got %s", v)
	}
}

func TestDecrementCellExpression(t *testing.T) {
	con := NewExecContext()
	con.topath = "TestDecrementCellExpression.xlsx"
	con.code = `["A1"]=0;if(["A1"]--)["A1"]=100;`
	run(con)
	if con.exitCode != 0 {
		t.Fatalf("exit code '%s'. want '%d' but got '%d'", con.code, 0, con.exitCode)
	}
	v := getCellValue(t, con.topath, "Sheet1", "A1")
	if v != "-1" {
		t.Fatalf("want cell value '-1', but got %s", v)
	}
}

func TestPreIncrementExpression(t *testing.T) {
	con := NewExecContext()
	con.topath = "TestPreIncrementExpression.xlsx"
	con.code = `a=0;if(++a)a=100;["A1"]=++a;`
	run(con)
	if con.exitCode != 0 {
		t.Fatalf("exit code '%s'. want '%d' but got '%d'", con.code, 0, con.exitCode)
	}
	v := getCellValue(t, con.topath, "Sheet1", "A1")
	if v != "101" {
		t.Fatalf("want cell value '101', but got %s", v)
	}
}

func TestPreIncrementCellExpression(t *testing.T) {
	con := NewExecContext()
	con.topath = "TestPreIncrementCellExpression.xlsx"
	con.code = `["A1"]=0;if(++["A1"])["A1"]=100;++["A1"]`
	run(con)
	if con.exitCode != 0 {
		t.Fatalf("exit code '%s'. want '%d' but got '%d'", con.code, 0, con.exitCode)
	}
	v := getCellValue(t, con.topath, "Sheet1", "A1")
	if v != "101" {
		t.Fatalf("want cell value '101', but got %s", v)
	}
}

func TestPreDecrementExpression(t *testing.T) {
	con := NewExecContext()
	con.topath = "TestPreDecrementExpression.xlsx"
	con.code = `a=0;if(--a)a=100;["A1"]=--a;`
	run(con)
	if con.exitCode != 0 {
		t.Fatalf("exit code '%s'. want '%d' but got '%d'", con.code, 0, con.exitCode)
	}
	v := getCellValue(t, con.topath, "Sheet1", "A1")
	if v != "99" {
		t.Fatalf("want cell value '99', but got %s", v)
	}
}

func TestPreDecrementCellExpression(t *testing.T) {
	con := NewExecContext()
	con.topath = "TestPreDecrementCellExpression.xlsx"
	con.code = `["A1"]=0;if(--["A1"])["A1"]=100;--["A1"]`
	run(con)
	if con.exitCode != 0 {
		t.Fatalf("exit code '%s'. want '%d' but got '%d'", con.code, 0, con.exitCode)
	}
	v := getCellValue(t, con.topath, "Sheet1", "A1")
	if v != "99" {
		t.Fatalf("want cell value '99', but got %s", v)
	}
}

func TestDoWhileStatement(t *testing.T) {
	con := NewExecContext()
	con.topath = "TestDoWhileStatement.xlsx"
	con.code = `do{sum+=i;i+=1;}while(i<10);["A1"]=sum;b=0;do["A2"]=b++;while(0);`
	run(con)

	if con.exitCode != 0 {
		t.Fatalf("exit code '%s'. want '%d', but got '%d'", con.code, 0, con.exitCode)
	}

	v := getCellValue(t, con.topath, "Sheet1", "A1")
	if v != "45" {
		t.Fatalf("do while statement could not working. want '45', but got %v", v)
	}
	v = getCellValue(t, con.topath, "Sheet1", "A2")
	if v != "0" {
		t.Fatalf("do while statement could not working. want '0', but got %v", v)
	}
}

func TestBreakStatementInDoWhile(t *testing.T) {
	con := NewExecContext()
	con.topath = "TestBreakStatementInDoWhile.xlsx"
	con.code = `do{sum+=i;if(i==3)break;i+=1;}while(i<10);["A1"]=sum`
	run(con)
	if con.exitCode != 0 {
		t.Fatalf("exit code '%s'. want '%d' but got '%d'", con.code, 0, con.exitCode)
	}
	v := getCellValue(t, con.topath, "Sheet1", "A1")
	if v != "6" {
		t.Fatalf("want cell value '6', but got %s", v)
	}
}

func TestContinueStatementInDoWhile(t *testing.T) {
	con := NewExecContext()
	con.topath = "TestContinueStatementInDoWhile.xlsx"
	con.code = `do{i+=1;if(i==3)continue;sum+=i;}while(i<10);["A1"]=sum`
	run(con)
	if con.exitCode != 0 {
		t.Fatalf("exit code '%s'. want '%d' but got '%d'", con.code, 0, con.exitCode)
	}
	v := getCellValue(t, con.topath, "Sheet1", "A1")
	if v != "52" {
		t.Fatalf("want cell value '52', but got %s", v)
	}
}

func TestForStatement(t *testing.T) {
	con := NewExecContext()
	con.topath = "TestForStatement.xlsx"
	con.code = `sum=0;for(i=0; i<10; i++) sum+=i; ["A1"]=sum`
	run(con)

	if con.exitCode != 0 {
		t.Fatalf("exit code '%s'. want '%d', but got '%d'", con.code, 0, con.exitCode)
	}

	v := getCellValue(t, con.topath, "Sheet1", "A1")
	if v != "45" {
		t.Fatalf("do while statement could not working. want '45', but got %v", v)
	}
}

func TestBreakStatementInFor(t *testing.T) {
	con := NewExecContext()
	con.topath = "TestBreakStatementInFor.xlsx"
	con.code = `for(i=0;i<10;i++){if(i==3)break;sum+=i;}["A1"]=sum;`
	run(con)
	if con.exitCode != 0 {
		t.Fatalf("exit code '%s'. want '%d' but got '%d'", con.code, 0, con.exitCode)
	}
	v := getCellValue(t, con.topath, "Sheet1", "A1")
	if v != "3" {
		t.Fatalf("want cell value '3', but got %s", v)
	}
}

func TestContinueStatementInFor(t *testing.T) {
	con := NewExecContext()
	con.topath = "TestContinueStatementInFor.xlsx"
	con.code = `for(i=0;i<10;i++){if(i==3)continue;sum+=i;}["A1"]=sum;`
	run(con)
	if con.exitCode != 0 {
		t.Fatalf("exit code '%s'. want '%d' but got '%d'", con.code, 0, con.exitCode)
	}
	v := getCellValue(t, con.topath, "Sheet1", "A1")
	if v != "42" {
		t.Fatalf("want cell value '42', but got %s", v)
	}
}

func TestDefineNoArgFunction(t *testing.T) {
	con := NewExecContext()
	con.topath = "TestDefineNoArgFunction.xlsx"
	con.code = `function answer() {["A1"] = 42;} answer();`
	run(con)
	if con.exitCode != 0 {
		t.Fatalf("exit code '%s'. want '%d' but got '%d'", con.code, 0, con.exitCode)
	}
	v := getCellValue(t, con.topath, "Sheet1", "A1")
	if v != "42" {
		t.Fatalf("want cell value '42', but got %s", v)
	}
}

func TestDefineWithArgFunction(t *testing.T) {
	con := NewExecContext()
	con.topath = "TestDefineWithArgFunction.xlsx"
	con.code = `function answer(one, two) {["A1"] = one;["A2"]=two;} answer(1,2);`
	run(con)
	if con.exitCode != 0 {
		t.Fatalf("exit code '%s'. want '%d' but got '%d'", con.code, 0, con.exitCode)
	}
	v := getCellValue(t, con.topath, "Sheet1", "A1")
	if v != "1" {
		t.Fatalf("want cell value '1', but got %s", v)
	}
	v = getCellValue(t, con.topath, "Sheet1", "A2")
	if v != "2" {
		t.Fatalf("want cell value '2', but got %s", v)
	}
}

func TestFunctionLocalScope(t *testing.T) {
	con := NewExecContext()
	con.topath = "TestFunctionLocalScope.xlsx"
	con.code = `a=10;b=20;function f(){a=100;["A1"]=a;["A2"]=b;} f();`
	run(con)
	if con.exitCode != 0 {
		t.Fatalf("exit code '%s'. want '%d' but got '%d'", con.code, 0, con.exitCode)
	}
	v := getCellValue(t, con.topath, "Sheet1", "A1")
	if v != "100" {
		t.Fatalf("want cell value '100', but got %s", v)
	}
	v = getCellValue(t, con.topath, "Sheet1", "A2")
	if v != "20" {
		t.Fatalf("want cell value '20', but got %s", v)
	}
}

func TestReturnStatement(t *testing.T) {
	con := NewExecContext()
	con.topath = "TestReturnStatement.xlsx"
	con.code = `function f(){return 100;["A2"] = 20;} ["A1"] = f();`
	run(con)
	if con.exitCode != 0 {
		t.Fatalf("exit code '%s'. want '%d' but got '%d'", con.code, 0, con.exitCode)
	}
	v := getCellValue(t, con.topath, "Sheet1", "A1")
	if v != "100" {
		t.Fatalf("want cell value '100', but got %s", v)
	}
	v = getCellValue(t, con.topath, "Sheet1", "A2")
	if v != "" {
		t.Fatalf("want cell value '', but got %s", v)
	}
}

func TestNestedReturnStatement(t *testing.T) {
	con := NewExecContext()
	con.topath = "TestNestedReturnStatement.xlsx"
	con.code = `function f1(){return 100;["A2"] = 20;} function f2(){ return f1() + 200; } ["A1"]=f2();`
	run(con)
	if con.exitCode != 0 {
		t.Fatalf("exit code '%s'. want '%d' but got '%d'", con.code, 0, con.exitCode)
	}
	v := getCellValue(t, con.topath, "Sheet1", "A1")
	if v != "300" {
		t.Fatalf("want cell value '300', but got %s", v)
	}
	v = getCellValue(t, con.topath, "Sheet1", "A2")
	if v != "" {
		t.Fatalf("want cell value '', but got %s", v)
	}
}

func TestRecursiveFunction(t *testing.T) {
	con := NewExecContext()
	con.topath = "TestRecursiveFunction.xlsx"
	con.code = `function fib(n) {if(n == 0 || n == 1) { return 1;} else { return fib(n-1)+fib(n-2);}} ["A1"]=fib(7);`
	run(con)
	if con.exitCode != 0 {
		t.Fatalf("exit code '%s'. want '%d' but got '%d'", con.code, 0, con.exitCode)
	}
	v := getCellValue(t, con.topath, "Sheet1", "A1")
	if v != "21" {
		t.Fatalf("want cell value '21', but got %s", v)
	}
}

func TestEvalStringAsNumber(t *testing.T) {
	con := NewExecContext()
	con.topath = "TestEvalStringAsNumber.xlsx"
	con.code = `["A1"] = "1" == 1;`
	run(con)
	if con.exitCode != 0 {
		t.Fatalf("exit code '%s'. want '%d' but got '%d'", con.code, 0, con.exitCode)
	}
	v := getCellValue(t, con.topath, "Sheet1", "A1")
	if v != "1" {
		t.Fatalf("want cell value '1', but got %s", v)
	}
}

func TestEvalNumberAsString(t *testing.T) {
	con := NewExecContext()
	con.topath = "TestEvalNumberAsString.xlsx"
	con.code = `["A1"] = "1" eq 1;`
	run(con)
	if con.exitCode != 0 {
		t.Fatalf("exit code '%s'. want '%d' but got '%d'", con.code, 0, con.exitCode)
	}
	v := getCellValue(t, con.topath, "Sheet1", "A1")
	if v != "1" {
		t.Fatalf("want cell value '1', but got %s", v)
	}
}

func TestIncrementExpressionForAtmark(t *testing.T) {
	con := NewExecContext()
	con.topath = "TestIncrementExpressionForAtmark.xlsx"
	con.code = `@="add1";@="add2";@="Sheet1";["A1"]=@++;`
	run(con)
	if con.exitCode != 0 {
		t.Fatalf("exit code '%s'. want '%d' but got '%d'", con.code, 0, con.exitCode)
	}
	v := getCellValue(t, con.topath, "add1", "A1")
	if v != "Sheet1" {
		t.Fatalf("want cell value 'Sheet1', but got %s", v)
	}
}

func TestPreIncrementExpressionForAtmark(t *testing.T) {
	con := NewExecContext()
	con.topath = "TestPreIncrementExpressionForAtmark.xlsx"
	con.code = `@="add1";@="add2";@="Sheet1";["A1"]=++@;`
	run(con)
	if con.exitCode != 0 {
		t.Fatalf("exit code '%s'. want '%d' but got '%d'", con.code, 0, con.exitCode)
	}
	v := getCellValue(t, con.topath, "add1", "A1")
	if v != "add1" {
		t.Fatalf("want cell value 'add1', but got %s", v)
	}
}

func TestDecrementExpressionForAtmark(t *testing.T) {
	con := NewExecContext()
	con.topath = "TestDecrementExpressionForAtmark.xlsx"
	con.code = `@="add1";@="add2";["A1"]=@--;`
	run(con)
	if con.exitCode != 0 {
		t.Fatalf("exit code '%s'. want '%d' but got '%d'", con.code, 0, con.exitCode)
	}
	v := getCellValue(t, con.topath, "add1", "A1")
	if v != "add2" {
		t.Fatalf("want cell value 'add2', but got %s", v)
	}
}

func TestPreDecrementExpressionForAtmark(t *testing.T) {
	con := NewExecContext()
	con.topath = "TestPreDecrementExpressionForAtmark.xlsx"
	con.code = `@="add1";@="add2";["A1"]=--@;`
	run(con)
	if con.exitCode != 0 {
		t.Fatalf("exit code '%s'. want '%d' but got '%d'", con.code, 0, con.exitCode)
	}
	v := getCellValue(t, con.topath, "add1", "A1")
	if v != "add1" {
		t.Fatalf("want cell value 'add1', but got %s", v)
	}
}

func TestResetInputSpecialVars(t *testing.T) {
	expect := "aa bb cc\ndd ee"
	in := bufio.NewReader(bytes.NewBufferString(expect))

	con := NewExecContext()
	con.topath = "TestResetInputSpecialVars.xlsx"
	con.in = in
	con.code = `gets();["A1"]=$1;["A2"]=$2;["A3"]=$3;gets();["A4"]=$1;["A5"]=$2;["A6"]=$3;`
	run(con)

	if con.exitCode != 0 {
		t.Fatalf("exit code '%s'. want '%d', but got '%d'", con.code, 0, con.exitCode)
	}

	actual := getCellValue(t, con.topath, "Sheet1", "A1")
	if actual != "aa" {
		t.Fatalf("A1 value wrong. want '%s', but got '%s'", "aa", actual)
	}

	actual = getCellValue(t, con.topath, "Sheet1", "A2")
	if actual != "bb" {
		t.Fatalf("A2 value wrong. want '%s', but got '%s'", "bb", actual)
	}

	actual = getCellValue(t, con.topath, "Sheet1", "A3")
	if actual != "cc" {
		t.Fatalf("A3 value wrong. want '%s', but got '%s'", "cc", actual)
	}

	actual = getCellValue(t, con.topath, "Sheet1", "A4")
	if actual != "dd" {
		t.Fatalf("A4 value wrong. want '%s', but got '%s'", "dd", actual)
	}

	actual = getCellValue(t, con.topath, "Sheet1", "A5")
	if actual != "ee" {
		t.Fatalf("A5 value wrong. want '%s', but got '%s'", "ee", actual)
	}

	actual = getCellValue(t, con.topath, "Sheet1", "A6")
	if actual != "" {
		t.Fatalf("A6 is '%s'. $3 should reset.", actual)
	}
}

func TestLastRowSpecialVar(t *testing.T) {
	con := NewExecContext()
	con.topath = "TestLastRowSpecialVar.xlsx"
	con.code = `["A1"]=LR;["A20"] = "a";["A2"]=LR;`
	run(con)
	if con.exitCode != 0 {
		t.Fatalf("exit code '%s'. want '%d' but got '%d'", con.code, 0, con.exitCode)
	}
	v := getCellValue(t, con.topath, "Sheet1", "A1")
	if v != "0" {
		t.Fatalf("LR want '%s', but got '%s'", "0", v)
	}
	v = getCellValue(t, con.topath, "Sheet1", "A2")
	if v != "20" {
		t.Fatalf("LR want '%s', but got '%s'", "20", v)
	}
}

func TestLastColSpecialVar(t *testing.T) {
	con := NewExecContext()
	con.topath = "TestLastColSpecialVar.xlsx"
	con.code = `["A1"]=LR;["E20"] = "a";["A2"]=LC;`
	run(con)
	if con.exitCode != 0 {
		t.Fatalf("exit code '%s'. want '%d' but got '%d'", con.code, 0, con.exitCode)
	}
	v := getCellValue(t, con.topath, "Sheet1", "A1")
	if v != "0" {
		t.Fatalf("LR want '%s', but got '%s'", "0", v)
	}
	v = getCellValue(t, con.topath, "Sheet1", "A2")
	if v != "5" {
		t.Fatalf("LR want '%s', but got '%s'", "5", v)
	}
}

func TestStringIncrement(t *testing.T) {
	con := NewExecContext()
	con.topath = "TestStringIncrement.xlsx"
	con.code = `col="z";b=col++;["A1"]=col;["A2"]=b;col="a1";col++;["A3"]=col;`
	run(con)
	if con.exitCode != 0 {
		t.Fatalf("exit code '%s'. want '%d' but got '%d'", con.code, 0, con.exitCode)
	}
	v := getCellValue(t, con.topath, "Sheet1", "A1")
	if v != "AA" {
		t.Fatalf("LR want '%s', but got '%s'", "AA", v)
	}
	v = getCellValue(t, con.topath, "Sheet1", "A2")
	if v != "z" {
		t.Fatalf("LR want '%s', but got '%s'", "z", v)
	}
	v = getCellValue(t, con.topath, "Sheet1", "A3")
	if v != "1" {
		t.Fatalf("LR want '%s', but got '%s'", "1", v)
	}
}

func TestStringPreIncrement(t *testing.T) {
	con := NewExecContext()
	con.topath = "TestStringPreIncrement.xlsx"
	con.code = `col="z";b=++col;["A1"]=col;["A2"]=b;col="a1";++col;["A3"]=col;`
	run(con)
	if con.exitCode != 0 {
		t.Fatalf("exit code '%s'. want '%d' but got '%d'", con.code, 0, con.exitCode)
	}
	v := getCellValue(t, con.topath, "Sheet1", "A1")
	if v != "AA" {
		t.Fatalf("LR want '%s', but got '%s'", "AA", v)
	}
	v = getCellValue(t, con.topath, "Sheet1", "A2")
	if v != "AA" {
		t.Fatalf("LR want '%s', but got '%s'", "AA", v)
	}
	v = getCellValue(t, con.topath, "Sheet1", "A3")
	if v != "1" {
		t.Fatalf("LR want '%s', but got '%s'", "1", v)
	}
}

func TestStringDecrement(t *testing.T) {
	con := NewExecContext()
	con.topath = "TestStringDecrement.xlsx"
	con.code = `col="aa";b=col--;["A1"]=col;["A2"]=b;col="a1";col--;["A3"]=col;`
	run(con)
	if con.exitCode != 0 {
		t.Fatalf("exit code '%s'. want '%d' but got '%d'", con.code, 0, con.exitCode)
	}
	v := getCellValue(t, con.topath, "Sheet1", "A1")
	if v != "Z" {
		t.Fatalf("LR want '%s', but got '%s'", "Z", v)
	}
	v = getCellValue(t, con.topath, "Sheet1", "A2")
	if v != "aa" {
		t.Fatalf("LR want '%s', but got '%s'", "aa", v)
	}
	v = getCellValue(t, con.topath, "Sheet1", "A3")
	if v != "-1" {
		t.Fatalf("LR want '%s', but got '%s'", "-1", v)
	}
}

func TestStringPreDecrement(t *testing.T) {
	con := NewExecContext()
	con.topath = "TestStringPreDecrement.xlsx"
	con.code = `col="aa";b=--col;["A1"]=col;["A2"]=b;col="a1";--col;["A3"]=col;`
	run(con)
	if con.exitCode != 0 {
		t.Fatalf("exit code '%s'. want '%d' but got '%d'", con.code, 0, con.exitCode)
	}
	v := getCellValue(t, con.topath, "Sheet1", "A1")
	if v != "Z" {
		t.Fatalf("LR want '%s', but got '%s'", "Z", v)
	}
	v = getCellValue(t, con.topath, "Sheet1", "A2")
	if v != "Z" {
		t.Fatalf("LR want '%s', but got '%s'", "Z", v)
	}
	v = getCellValue(t, con.topath, "Sheet1", "A3")
	if v != "-1" {
		t.Fatalf("LR want '%s', but got '%s'", "-1", v)
	}
}

func TestStringCellDecrement(t *testing.T) {
	con := NewExecContext()
	con.topath = "TestStringCellDecrement.xlsx"
	con.code = `["A1"]=2;b=["A1"]--;["A2"]=b;["A3"]="a";["A3"]--`
	run(con)
	if con.exitCode != 0 {
		t.Fatalf("exit code '%s'. want '%d' but got '%d'", con.code, 0, con.exitCode)
	}
	v := getCellValue(t, con.topath, "Sheet1", "A1")
	if v != "1" {
		t.Fatalf("LR want '%s', but got '%s'", "1", v)
	}

	v = getCellValue(t, con.topath, "Sheet1", "A2")
	if v != "2" {
		t.Fatalf("LR want '%s', but got '%s'", "2", v)
	}

	v = getCellValue(t, con.topath, "Sheet1", "A3")
	if v != "-1" {
		t.Fatalf("LR want '%s', but got '%s'", "-1", v)
	}
}

func TestStringCellPreDecrement(t *testing.T) {
	con := NewExecContext()
	con.topath = "TestStringCellPreDecrement.xlsx"
	con.code = `["A1"]=2;b=--["A1"];["A2"]=b;["A3"]="a";--["A3"]`
	run(con)
	if con.exitCode != 0 {
		t.Fatalf("exit code '%s'. want '%d' but got '%d'", con.code, 0, con.exitCode)
	}
	v := getCellValue(t, con.topath, "Sheet1", "A1")
	if v != "1" {
		t.Fatalf("LR want '%s', but got '%s'", "1", v)
	}
	v = getCellValue(t, con.topath, "Sheet1", "A2")
	if v != "1" {
		t.Fatalf("LR want '%s', but got '%s'", "1", v)
	}
	v = getCellValue(t, con.topath, "Sheet1", "A3")
	if v != "-1" {
		t.Fatalf("LR want '%s', but got '%s'", "-1", v)
	}
}

func TestRSSpecialVar(t *testing.T) {
	in := bufio.NewReader(bytes.NewBufferString("1 2 3\t4 5 6"))
	out := new(bytes.Buffer)

	con := NewExecContext()
	con.in = in
	con.out = out

	con.code = `RS="\t";while(gets())puts();`
	run(con)

	if con.exitCode != 0 {
		t.Fatalf("exit code '%s'. want '%d' but got '%d'", con.code, 0, con.exitCode)
	}

	if out.String() != "1 2 3\n4 5 6\n" {
		t.Fatalf("want stdout '1 2 3\n4 5 6\n', but got '%s'", out)
	}
}
