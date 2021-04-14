package main

import (
	"fmt"

	"github.com/tealeg/xlsx/v3"
)

type Spreadsheet struct {
	file   *xlsx.File
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
	s.file = xlsx.NewFile()
	_, err := s.file.AddSheet("Sheet1")
	if err != nil {
		return err
	}
	return nil
}

func (s *Spreadsheet) readSpreadsheet(frompath string) error {
	f, err := xlsx.OpenFile(frompath)
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

	err := s.file.Save(s.topath)
	if err != nil {
		return fmt.Errorf("on error spreadsheet writing. '%v'", err)
	}

	return nil
}
