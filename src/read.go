package src

import "fmt"

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
	// Check if the range is in bounds.
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

	// Check if the range is valid.
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
