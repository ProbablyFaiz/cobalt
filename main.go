package main

import (
	. "pasado/src"
)

func main() {
	// Initialize the spreadsheet.
	ss := NewSpreadsheet()

	// Add a formula to a cell.
	_ = ss.UpdateCell("Sheet1", 0, 0, `="Hello, "`)
	_ = ss.UpdateCell("Sheet1", 0, 1, `="World!"`)
	_ = ss.UpdateCell("Sheet1", 0, 2, `=CONCAT(A1, B1)`)

	// Print the value of the cell.
	cell, _ := ss.GetCell("Sheet1", 0, 2)
	println(cell.Value.(string))
}
