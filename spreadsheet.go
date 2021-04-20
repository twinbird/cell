package main

import (
	"fmt"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
)

type Spreadsheet struct {
	file        *excelize.File
	topath      string
	activeSheet string
}

func NewSpreadsheet(frompath string, topath string) (*Spreadsheet, error) {
	s := &Spreadsheet{
		topath: topath,
	}

	if frompath != "" {
		err := s.readSpreadsheet(frompath)
		if err != nil {
			return nil, fmt.Errorf("on error spreadsheet reading '%v'", err)
		}
	} else {
		err := s.createSpreadsheet()
		if err != nil {
			return nil, fmt.Errorf("on error spreadsheet creating '%v'", err)
		}
	}

	return s, nil
}

func (s *Spreadsheet) createSpreadsheet() error {
	s.file = excelize.NewFile()
	s.activeSheet = s.getActiveSheetName()
	return nil
}

func (s *Spreadsheet) readSpreadsheet(frompath string) error {
	f, err := excelize.OpenFile(frompath)
	if err != nil {
		return err
	}
	s.file = f
	s.activeSheet = s.getActiveSheetName()

	return nil
}

func (s *Spreadsheet) writeSpreadsheet() error {
	if s.topath == "" {
		return fmt.Errorf("on error spreadsheet writing: no specify write path.")
	}

	err := s.file.SaveAs(s.topath)
	if err != nil {
		return fmt.Errorf("on error spreadsheet writing. '%v'", err)
	}

	return nil
}

func (s *Spreadsheet) getCellValue(axis string) string {
	v, err := s.file.GetCellValue(s.activeSheet, axis)
	if err != nil {
		fatalError("cell '%s' refer failed", axis)
	}
	return v
}

func (s *Spreadsheet) setCellValue(axis string, v interface{}) {
	err := s.file.SetCellValue(s.activeSheet, axis, v)
	if err != nil {
		fatalError("cell '%s' set value failed", axis)
	}
}

func (s *Spreadsheet) getActiveSheetName() string {
	idx := s.file.GetActiveSheetIndex()
	name := s.file.GetSheetName(idx)
	return name
}

func (s *Spreadsheet) setActiveSheetByName(name string) error {
	idx := s.file.GetSheetIndex(name)
	if idx < 0 {
		return fmt.Errorf("sheet %s is not found.", name)
	}
	s.file.SetActiveSheet(idx)
	s.activeSheet = name

	return nil
}

func (s *Spreadsheet) getSheetList() []string {
	return s.file.GetSheetList()
}

func (s *Spreadsheet) addSheet(name string) error {
	s.file.NewSheet(name)
	return s.setActiveSheetByName(name)
}
