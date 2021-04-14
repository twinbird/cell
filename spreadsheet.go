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
		err := s.readSpreadsheet()
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

func (s *Spreadsheet) readSpreadsheet() error {
	f, err := excelize.OpenFile(frompath)
	if err != nil {
		return err
	}
	s.file = f

	return nil
}

func (s *Spreadsheet) writeSpreadsheet() error {
	if topath != "" {
		err := s.file.SaveAs(topath)
		if err != nil {
			return fmt.Errorf("on error spreadsheet closing. '%v'", err)
		}
	}

	return nil
}
