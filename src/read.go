package src

import (
	"fmt"
	"github.com/Workiva/go-datastructures/augmentedtree"
)

func (ss *Spreadsheet) GetCell(sheetName string, row int, col int) (*Cell, error) {
	// Check if the sheet exists and if the cell is in bounds.
	sheet, ok := ss.Sheets[sheetName]
	if !ok {
		return nil, fmt.Errorf("sheet %s does not exist", sheetName)
	}
	if row < 0 || row >= len(sheet.Cells) {
		return nil, fmt.Errorf("row %d is out of bounds", row)
	}
	if col < 0 || col >= len(sheet.Cells[row]) {
		return nil, fmt.Errorf("col %d is out of bounds", col)
	}

	return &sheet.Cells[row][col], nil
}

func (sheet *Sheet) GetRange(startRow int, startCol int, endRow int, endCol int) ([][]Cell, error) {
	// Check if the range is valid.
	if startRow < 0 || startRow >= len(sheet.Cells) {
		return nil, fmt.Errorf("startRow %d is out of bounds", startRow)
	}
	if startCol < 0 || startCol >= len(sheet.Cells[startRow]) {
		return nil, fmt.Errorf("startCol %d is out of bounds", startCol)
	}
	if endRow < 0 || endRow >= len(sheet.Cells) {
		return nil, fmt.Errorf("endRow %d is out of bounds", endRow)
	}
	if endCol < 0 || endCol >= len(sheet.Cells[endRow]) {
		return nil, fmt.Errorf("endCol %d is out of bounds", endCol)
	}
	if startRow > endRow {
		return nil, fmt.Errorf("startRow %d is greater than endRow %d", startRow, endRow)
	}
	if startCol > endCol {
		return nil, fmt.Errorf("startCol %d is greater than endCol %d", startCol, endCol)
	}

	// Get the range.
	rangeCells := make([][]Cell, endRow-startRow+1)
	for i := startRow; i <= endRow; i++ {
		rangeCells[i-startRow] = sheet.Cells[i][startCol : endCol+1]
	}
	return rangeCells, nil
}

func (cell *Cell) GetValue() interface{} {
	return cell.Value
}

func (cell *Cell) GetUuid() ReferenceId {
	return cell.Uuid
}

func (cr *Range) GetValue() interface{} {
	values, err := cr.Sheet.GetRange(cr.FromRow, cr.FromCol, cr.ToRow, cr.ToCol)
	if err != nil {
		// For a resolved range, it should be impossible to have an invalid range (I think).
		panic(err)
	}
	return values
}

func (cr *Range) GetUuid() ReferenceId {
	return cr.Uuid
}

func (cr *Range) LowAtDimension(dimension uint64) int64 {
	if dimension == 0 {
		return int64(cr.FromRow)
	} else {
		return int64(cr.FromCol)
	}
}

func (cr *Range) HighAtDimension(dimension uint64) int64 {
	if dimension == 0 {
		return int64(cr.ToRow)
	} else {
		return int64(cr.ToCol)
	}
}

func (cr *Range) OverlapsAtDimension(interval augmentedtree.Interval, dimension uint64) bool {
	if dimension == 0 {
		return int64(cr.FromRow) <= interval.HighAtDimension(uint64(dimension)) && int64(cr.ToRow) >= interval.LowAtDimension(uint64(dimension))
	} else {
		return int64(cr.FromCol) <= interval.HighAtDimension(uint64(dimension)) && int64(cr.ToCol) <= interval.LowAtDimension(uint64(dimension))
	}
}

func (cr *Range) ID() uint64 {
	return uint64(cr.Uuid)
}

func (cell *Cell) toRange() *Range {
	return &Range{
		Sheet:   cell.Sheet,
		FromRow: cell.Row,
		FromCol: cell.Col,
		ToRow:   cell.Row,
		ToCol:   cell.Col,
		Uuid:    cell.Uuid,
	}
}
