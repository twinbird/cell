package main

import "testing"

func TestOpenFileSpecifiedFromPath(t *testing.T) {
	sheet, err := NewSpreadsheet("test/empty.xlsx", "")
	if err != nil {
		t.Fatalf("An error occurred when opening a file that exists: %v", err)
	}
	if sheet.file == nil {
		t.Fatal("File object was not opened when creating the Spreadsheet object.")
	}
}

func TestOpenFileSpecifiedNotExistFile(t *testing.T) {
	sheet, err := NewSpreadsheet("test/notexist.xlsx", "")
	if err == nil {
		t.Fatal("Opened a non-existent file, but no error occurred.")
	}
	if sheet != nil {
		t.Fatal("A non-existent file was opened but sheet was not returned as nil.")
	}
}

func TestOpenFileUnspecifiedFromPath(t *testing.T) {
	sheet, err := NewSpreadsheet("", "")
	if err != nil {
		t.Fatal("An error occurred when creating Spreadsheet without specifying Frompath.")
	}
	if sheet == nil {
		t.Fatal("Could not create Spreadsheet without specifying Frompath.")
	}
}

func TestWriteSpreadsheetUnspecifiedToPath(t *testing.T) {
	sheet, _ := NewSpreadsheet("", "")
	err := sheet.writeSpreadsheet()
	if err == nil {
		t.Fatal("Calling WriteSpreadsheet without specifying topath did not generate an error.")
	}
}

func TestWriteSpreadsheetSpecifiedToPath(t *testing.T) {
	sheet, _ := NewSpreadsheet("", "test.xlsx")
	err := sheet.writeSpreadsheet()
	if err != nil {
		t.Fatalf("Calling WriteSpreadsheet with topath specified caused an error: %v", err)
	}
}

func TestGetActiveSheetName(t *testing.T) {
	sheet, _ := NewSpreadsheet("test/sheet.xlsx", "")
	name := sheet.getActiveSheetName()
	if name != "Sheet3" {
		t.Fatalf("active sheet want %s, but got %s", "Sheet3", name)
	}
}

func TestSetActiveSheetByName(t *testing.T) {
	sheet, _ := NewSpreadsheet("test/sheet.xlsx", "")
	if err := sheet.setActiveSheetByName("Sheet1"); err != nil {
		t.Fatalf("error '%v' on active sheet set.", err)
	}
	name := sheet.getActiveSheetName()
	if name != "Sheet1" {
		t.Fatalf("active sheet want %s, but got %s", "Sheet1", name)
	}
}

func TestSetActiveSheetByNotExistName(t *testing.T) {
	sheet, _ := NewSpreadsheet("test/sheet.xlsx", "")
	if err := sheet.setActiveSheetByName("foobar"); err == nil {
		t.Fatalf("No error occured even through 'setActiveSheetByName' with not exist sheet name")
	}
	name := sheet.getActiveSheetName()
	if name != "Sheet3" || sheet.activeSheet != name {
		t.Fatalf("active sheet changed by 'setActiveSheetByName' invalided call.")
	}
}

func TestGetCellValue(t *testing.T) {
	sheet, _ := NewSpreadsheet("test/values.xlsx", "")

	sheet.setActiveSheetByName("Sheet1")
	v := sheet.getCellValue("A1")
	if v != "2" {
		t.Fatalf("Sheet1[A1] want %s, but got %s", "2", v)
	}

	err := sheet.setActiveSheetByName("Sheet2")
	if err != nil {
		t.Fatalf("%v", err)
	}
	v = sheet.getCellValue("A1")
	if v != "sheet2" {
		t.Fatalf("Sheet2[A1] want %s, but got %s", "sheet2", v)
	}
}

func TestSetCellValue(t *testing.T) {
	sheet, _ := NewSpreadsheet("test/values.xlsx", "")

	sheet.setActiveSheetByName("Sheet1")
	v := sheet.getCellValue("A1")
	if v != "2" {
		t.Fatalf("Sheet1[A1] want %s, but got %s", "2", v)
	}

	sheet.setCellValue("A1", "20")
	v = sheet.getCellValue("A1")
	if v != "20" {
		t.Fatalf("Sheet1[A1] want %s, but got %s", "20", v)
	}
}

func TestGetSheetList(t *testing.T) {
	sheet, _ := NewSpreadsheet("test/values.xlsx", "")
	names := sheet.getSheetList()

	if len(names) != 3 {
		t.Fatalf("test/values.xlsx has sheets want %d, but got %d", 3, len(names))
	}

	expects := []string{"Sheet1", "Sheet2", "Sheet3"}
	if names[0] != expects[0] || names[1] != expects[1] || names[2] != expects[2] {
		t.Fatalf("test/values.xlsx has sheets want '%v', but got '%v'", expects, names)
	}
}

func TestAddSheet(t *testing.T) {
	sheet, _ := NewSpreadsheet("", "")
	names := sheet.getSheetList()

	if len(names) != 1 {
		t.Fatalf("sheet counts want %d, but got %d", 1, len(names))
	}

	sheet.addSheet("foo")
	n := sheet.getActiveSheetName()
	if n != "foo" {
		t.Fatalf("active sheet has not changed even through calling addsheet()")
	}

	names = sheet.getSheetList()
	if len(names) != 2 {
		t.Fatalf("sheet counts want %d, but got %d", 2, len(names))
	}
	if names[0] != "Sheet1" || names[1] != "foo" {
		t.Fatalf("addsheet() called, but no sheet has been added")
	}
}

func TestSetNextSheet(t *testing.T) {
	sheet, _ := NewSpreadsheet("", "")
	sheet.addSheet("foo")
	sheet.addSheet("bar")
	sheet.addSheet("baz")

	sheet.setActiveSheetByName("Sheet1")
	n := sheet.setNextSheet()
	if n != "foo" {
		t.Fatalf("setNextSheet() return '%s', but want '%s'", n, "foo")
	}
	s := sheet.getActiveSheetName()
	if s != "foo" {
		t.Fatalf("setNextSheet() could not change sheet. want '%s', but '%s'", "foo", s)
	}

	n = sheet.setNextSheet()
	if n != "bar" {
		t.Fatalf("setNextSheet() return '%s', but want '%s'", n, "bar")
	}
	s = sheet.getActiveSheetName()
	if s != "bar" {
		t.Fatalf("setNextSheet() could not change sheet. want '%s', but '%s'", "bar", s)
	}

	n = sheet.setNextSheet()
	if n != "baz" {
		t.Fatalf("setNextSheet() return '%s', but want '%s'", n, "baz")
	}
	s = sheet.getActiveSheetName()
	if s != "baz" {
		t.Fatalf("setNextSheet() could not change sheet. want '%s', but '%s'", "baz", s)
	}

	n = sheet.setNextSheet()
	if n != "" {
		t.Fatalf("setNextSheet() return '%s', but want '%s'", n, "")
	}
	s = sheet.getActiveSheetName()
	if s != "baz" {
		t.Fatalf("setNextSheet() changed sheet on last sheet. want '%s', but '%s'", "baz", s)
	}
}

func TestSetPrevSheet(t *testing.T) {
	sheet, _ := NewSpreadsheet("", "")
	sheet.addSheet("foo")
	sheet.addSheet("bar")
	sheet.addSheet("baz")

	sheet.setActiveSheetByName("baz")
	n := sheet.setPrevSheet()
	if n != "bar" {
		t.Fatalf("setPrevSheet() return '%s', but want '%s'", n, "bar")
	}
	s := sheet.getActiveSheetName()
	if s != "bar" {
		t.Fatalf("setPrevSheet() could not change sheet. want '%s', but '%s'", "bar", s)
	}

	n = sheet.setPrevSheet()
	if n != "foo" {
		t.Fatalf("setPrevSheet() return '%s', but want '%s'", n, "foo")
	}
	s = sheet.getActiveSheetName()
	if s != "foo" {
		t.Fatalf("setPrevSheet() could not change sheet. want '%s', but '%s'", "foo", s)
	}

	n = sheet.setPrevSheet()
	if n != "Sheet1" {
		t.Fatalf("setPrevSheet() return '%s', but want '%s'", n, "Sheet1")
	}
	s = sheet.getActiveSheetName()
	if s != "Sheet1" {
		t.Fatalf("setPrevSheet() could not change sheet. want '%s', but '%s'", "Sheet1", s)
	}

	n = sheet.setPrevSheet()
	if n != "" {
		t.Fatalf("setPrevSheet() return '%s', but want '%s'", n, "")
	}
	s = sheet.getActiveSheetName()
	if s != "Sheet1" {
		t.Fatalf("setPrevSheet() changed sheet on last sheet. want '%s', but '%s'", "Sheet1", s)
	}
}

func TestSetHeadSheet(t *testing.T) {
	sheet, _ := NewSpreadsheet("", "")
	sheet.addSheet("foo")
	sheet.addSheet("bar")
	sheet.addSheet("baz")

	sheet.setActiveSheetByName("bar")

	n := sheet.setHeadSheet()
	if n != "Sheet1" {
		t.Fatalf("setHeadSheet() could not changed active sheet. want '%s', but '%s'", "Sheet1", n)
	}
	s := sheet.getActiveSheetName()
	if s != "Sheet1" {
		t.Fatalf("setHeadSheet() could not changed active sheet name. want '%s', but '%s'", "Sheet1", s)
	}
}

func TestSetTailSheet(t *testing.T) {
	sheet, _ := NewSpreadsheet("", "")
	sheet.addSheet("foo")
	sheet.addSheet("bar")
	sheet.addSheet("baz")

	sheet.setActiveSheetByName("bar")

	n := sheet.setTailSheet()
	if n != "baz" {
		t.Fatalf("setTailSheet() could not changed active sheet. want '%s', but '%s'", "baz", n)
	}
	s := sheet.getActiveSheetName()
	if s != "baz" {
		t.Fatalf("setTailSheet() could not changed active sheet name. want '%s', but '%s'", "baz", s)
	}
}

func TestGetColsCount(t *testing.T) {
	sheet, _ := NewSpreadsheet("", "")
	sheet.setCellValue("D10", "20")
	sheet.setCellValue("E11", "21")

	v := sheet.getColsCount()
	if v != 5 {
		t.Fatalf("cols count want %d, but got %d", 5, v)
	}
}

func TestGetRowsCount(t *testing.T) {
	sheet, _ := NewSpreadsheet("", "")
	sheet.setCellValue("D10", "20")
	sheet.setCellValue("E11", "21")

	v := sheet.getRowsCount()
	if v != 11 {
		t.Fatalf("cols count want %d, but got %d", 11, v)
	}
}

func TestGetAlphaColsCount(t *testing.T) {
	sheet, _ := NewSpreadsheet("", "")
	sheet.setCellValue("D10", "20")
	sheet.setCellValue("E11", "21")

	v := sheet.getAlphaColsCount()
	if v != "E" {
		t.Fatalf("cols count want %s, but got %s", "E", v)
	}
}

func TestIncrementColumnNumber(t *testing.T) {
	c := "A"
	v, err := incrementColumnNumber(c)
	if err != nil {
		t.Fatalf("error '%v' occured at 'A' increment", err)
	}
	if v != "B" {
		t.Fatalf("cols count want %s, but got %s", "B", v)
	}

	c = "Z"
	v, err = incrementColumnNumber(c)
	if err != nil {
		t.Fatalf("error '%v' occured at 'Z' increment", err)
	}
	if v != "AA" {
		t.Fatalf("cols count want %s, but got %s", "AA", v)
	}

	c = "A1"
	v, err = incrementColumnNumber(c)
	if err == nil {
		t.Fatalf("converted 'A1' but not error occured.")
	}
}

func TestDecrementColumnNumber(t *testing.T) {
	c := "B"
	v, err := decrementColumnNumber(c)
	if err != nil {
		t.Fatalf("error '%v' occured at 'B' decrement", err)
	}
	if v != "A" {
		t.Fatalf("cols count want %s, but got %s", "A", v)
	}

	c = "AA"
	v, err = decrementColumnNumber(c)
	if err != nil {
		t.Fatalf("error '%v' occured at 'AA' decrement", err)
	}
	if v != "Z" {
		t.Fatalf("cols count want %s, but got %s", "Z", v)
	}

	c = "A1"
	v, err = decrementColumnNumber(c)
	if err == nil {
		t.Fatalf("converted 'A1' but not error occured.")
	}

	c = "A"
	v, err = decrementColumnNumber(c)
	if err == nil {
		t.Fatalf("converted 'A' but not error occured.")
	}
}

func TestExistSheetName(t *testing.T) {
	sheet, _ := NewSpreadsheet("", "")

	v := sheet.existSheetName("Sheet1")
	if v != true {
		t.Fatalf("want true but got %v", v)
	}
	v = sheet.existSheetName("Sheet2")
	if v != false {
		t.Fatalf("want false but got %v", v)
	}
}

func TestSetSheetName(t *testing.T) {
	sheet, _ := NewSpreadsheet("", "")

	v := sheet.setSheetName("Sheet1", "Sheet2")
	if v != "Sheet2" {
		t.Fatalf("want 'Sheet2' but got %v", v)
	}
	v = sheet.setSheetName("Sheet1", "Sheet2")
	if v != "" {
		t.Fatalf("want '' but got %v", v)
	}
}

func TestCountSheet(t *testing.T) {
	sheet, _ := NewSpreadsheet("", "")

	v := sheet.countSheet()
	if v != 1 {
		t.Fatalf("countSheet() want 1 but got %d", v)
	}

	sheet.addSheet("Sheet2")

	v = sheet.countSheet()
	if v != 2 {
		t.Fatalf("countSheet() want 2 but got %d", v)
	}
}

func TestDeleteSheet(t *testing.T) {
	sheet, _ := NewSpreadsheet("", "")

	v := sheet.deleteSheet("Sheet1")
	if v != false {
		t.Fatalf("deleteSheet() return true on delete last sheet")
	}

	sheet.addSheet("Sheet2")

	v = sheet.deleteSheet("Sheet1")
	if v != true {
		t.Fatalf("deleteSheet() could not delete sheet")
	}

	n := sheet.countSheet()
	if n != 1 {
		t.Fatalf("deleteSheet() could not delete sheet")
	}
}

func TestCopySheet(t *testing.T) {
	sheet, _ := NewSpreadsheet("", "")

	v := sheet.copySheet("Sheet1", "Sheet1")
	if v != false {
		t.Fatalf("copySheet() return true on duplicate sheet")
	}

	sheet.setCellValue("A1", "Hello")
	v = sheet.copySheet("Sheet1", "Sheet2")
	if v != true {
		t.Fatalf("copySheet() could not copy sheet")
	}

	n := sheet.countSheet()
	if n != 2 {
		t.Fatalf("copySheet() could not copy sheet")
	}

	sheet.setActiveSheetByName("Sheet2")
	s := sheet.getCellValue("A1")
	if s != "Hello" {
		t.Fatalf("copySheet() could not copy value in sheet")
	}
}
