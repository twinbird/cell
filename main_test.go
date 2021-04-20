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
		t.Fatalf("on error occured set cell value '%s'.", axis)
	}
}

func getCellValue(t *testing.T, filepath string, sheet string, axis string) string {
	f, err := excelize.OpenFile(filepath)
	if err != nil {
		t.Fatalf("on error occured open '%s'.", filepath)
	}
	v, err := f.GetCellValue(sheet, axis)
	if err != nil {
		t.Fatalf("on error occured get cell value '%s'.", axis)
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

func TestSpecialVarAStMarkAssignUndefinedSheet(t *testing.T) {
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
