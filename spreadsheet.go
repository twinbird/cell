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

func (s *Spreadsheet) setNextSheet() string {
	list := s.file.GetSheetList()
	current := s.file.GetSheetIndex(s.activeSheet)
	if current < 0 {
		fatalError("current worksheet not found in setNextSheet()")
	}
	idx := current + 1

	if idx < 0 || len(list) <= idx {
		return ""
	}

	s.file.SetActiveSheet(idx)
	name := s.file.GetSheetName(idx)
	s.activeSheet = name
	return name
}

func (s *Spreadsheet) setPrevSheet() string {
	list := s.file.GetSheetList()
	current := s.file.GetSheetIndex(s.activeSheet)
	if current < 0 {
		fatalError("current worksheet not found in setPrevSheet()")
	}
	idx := current - 1

	if idx < 0 || len(list) <= idx {
		return ""
	}

	s.file.SetActiveSheet(idx)
	name := s.file.GetSheetName(idx)
	s.activeSheet = name
	return name
}

func (s *Spreadsheet) setHeadSheet() string {
	current := s.file.GetSheetIndex(s.activeSheet)
	if current < 0 {
		fatalError("current worksheet not found in setPrevSheet()")
	}
	idx := 0

	s.file.SetActiveSheet(idx)
	name := s.file.GetSheetName(idx)
	s.activeSheet = name
	return name
}

func (s *Spreadsheet) setTailSheet() string {
	list := s.file.GetSheetList()
	current := s.file.GetSheetIndex(s.activeSheet)
	if current < 0 {
		fatalError("current worksheet not found in setPrevSheet()")
	}
	idx := len(list) - 1

	s.file.SetActiveSheet(idx)
	name := s.file.GetSheetName(idx)
	s.activeSheet = name
	return name
}

func (s *Spreadsheet) getColsCount() int {
	current := s.file.GetSheetIndex(s.activeSheet)
	if current < 0 {
		fatalError("current worksheet not found in getColsCount()")
	}
	cols, err := s.file.GetCols(s.activeSheet)
	if err != nil {
		fatalError("cols count error. %v", err)
	}
	return len(cols)
}

func (s *Spreadsheet) getAlphaColsCount() string {
	c := s.getColsCount()
	name, err := excelize.ColumnNumberToName(c)
	if err != nil {
		return ""
	}
	return name
}

func (s *Spreadsheet) getRowsCount() int {
	current := s.file.GetSheetIndex(s.activeSheet)
	if current < 0 {
		fatalError("current worksheet not found in getRowsCount()")
	}
	rows, err := s.file.GetRows(s.activeSheet)
	if err != nil {
		fatalError("rows count error. %v", err)
	}
	return len(rows)
}

func columnNameToNumber(name string) (int, error) {
	return excelize.ColumnNameToNumber(name)
}

func columnNumberToName(num int) (string, error) {
	return excelize.ColumnNumberToName(num)
}

func incrementColumnNumber(name string) (string, error) {
	n, err := columnNameToNumber(name)
	if err != nil {
		return "", err
	}
	n++
	return columnNumberToName(n)
}

func decrementColumnNumber(name string) (string, error) {
	n, err := columnNameToNumber(name)
	if err != nil {
		return "", err
	}
	n--
	return columnNumberToName(n)
}
