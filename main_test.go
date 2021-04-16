package main

import (
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

func TestSimpleNumberExpression(t *testing.T) {
	con := &ExecContext{}
	con.code = `1`
	run(con)
	if con.exitCode != 1 {
		t.Fatalf("exec code '%s'. want '%d' but got '%d'", con.code, 1, con.exitCode)
	}
}

func TestSimpleStringExpression(t *testing.T) {
	con := &ExecContext{}
	con.code = `"str"`
	run(con)
	if con.exitCode != 0 {
		t.Fatalf("exec code '%s'. want '%d' but got '%d'", con.code, 0, con.exitCode)
	}
}

func TestSimpleCellReferExpression(t *testing.T) {
	con := &ExecContext{}
	con.frompath = "test/values.xlsx"
	con.code = `["A1"]`
	run(con)
	if con.exitCode != 2 {
		t.Fatalf("exec code '%s'. want '%d' but got '%d'", con.code, 2, con.exitCode)
	}
}

func TestSimpleCellAssignExpression(t *testing.T) {
	con := &ExecContext{}
	con.frompath = "test/values.xlsx"
	con.topath = "TestSimpleCellAssignExpression.xlsx"
	con.code = `["A1"] = 5`
	run(con)
	if con.exitCode != 5 {
		t.Fatalf("exec code '%s'. want '%d' but got '%d'", con.code, 5, con.exitCode)
	}
	v := getCellValue(t, con.topath, "Sheet1", "A1")
	if v != "5" {
		t.Fatalf("want cell value 5, but got %s", v)
	}
}

func TestCellAssignToString(t *testing.T) {
	con := &ExecContext{}
	con.frompath = "test/values.xlsx"
	con.topath = "TestCellAssignToString.xlsx"
	con.code = `["A1"] = "abc"`
	run(con)
	if con.exitCode != 0 {
		t.Fatalf("exec code '%s'. want '%d' but got '%d'", con.code, 0, con.exitCode)
	}
	v := getCellValue(t, con.topath, "Sheet1", "A1")
	if v != "abc" {
		t.Fatalf("want cell value 'abc', but got %s", v)
	}
}

func TestCellReferFromString(t *testing.T) {
	con := &ExecContext{}
	con.frompath = "test/values.xlsx"
	con.topath = "TestCellReferFromString.xlsx"
	con.code = `["A3"] = ["A2"]`
	run(con)
	if con.exitCode != 0 {
		t.Fatalf("exec code '%s'. want '%d' but got '%d'", con.code, 0, con.exitCode)
	}
	v := getCellValue(t, con.topath, "Sheet1", "A3")
	if v != "test" {
		t.Fatalf("want cell value 'test', but got %s", v)
	}
}
