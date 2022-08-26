package main

import (
	. "pasado/src"
)

func main() {
	// Initialize the spreadsheet.
	ss := NewSpreadsheet()

	// Add a formula to a cell.
	cell00, _ := ss.GetCell("Sheet1", 0, 0)
	ss.UpdateCell(cell00.Uuid, `9`)
	cell01, _ := ss.GetCell("Sheet1", 0, 1)
	ss.UpdateCell(cell01.Uuid, `8`)
	cell02, _ := ss.GetCell("Sheet1", 0, 2)
	ss.UpdateCell(cell02.Uuid, `=If(A1 - B1, 3, 4)`)

	// Print the value of the cell.
	cell, _ := ss.GetCell("Sheet1", 0, 2)
	println(cell.Value.(int))
}
