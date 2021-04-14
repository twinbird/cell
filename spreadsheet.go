package main

import (
	"fmt"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
)

type Spreadsheet struct {
	file   *excelize.File
	topath string
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
	return nil
}

func (s *Spreadsheet) readSpreadsheet(frompath string) error {
	f, err := excelize.OpenFile(frompath)
	if err != nil {
		return err
	}
	s.file = f

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
	v, err := s.file.GetCellValue("Sheet1", axis)
	if err != nil {
		fatalError("cell '%s' refer failed", axis)
	}
	return v
}
