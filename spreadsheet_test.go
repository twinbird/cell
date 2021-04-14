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
