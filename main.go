package main

import (
	. "pasado/src"
)

func main() {
	// Initialize the spreadsheet.
	ss := NewSpreadsheet()

	// Add a formula to a cell.
	cell00, _ := ss.GetCell("Sheet1", 0, 0)
	_ = ss.UpdateCell(cell00.Uuid, `3`)
	cell01, _ := ss.GetCell("Sheet1", 0, 1)
	_ = ss.UpdateCell(cell01.Uuid, `9`)
	cell02, _ := ss.GetCell("Sheet1", 0, 2)
	_ = ss.UpdateCell(cell02.Uuid, `=Add(A1, B1)`)

	// Print the value of the cell.
	cell, _ := ss.GetCell("Sheet1", 0, 2)
	println(cell.Value.(int))
}
