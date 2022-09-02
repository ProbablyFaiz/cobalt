package main

import (
	. "pasado/src/core"
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
	ss.UpdateCell(cell02.Uuid, `=Sum(A1:B1)`)

	// Print the value of the cell.
	cell, _ := ss.GetCell("Sheet1", 0, 2)
	println(cell.Value.(int))

	ss.UpdateCell(cell00.Uuid, `15`)
	println(`Cell00 value is now:`, cell00.Value.(int))
	println(`Cell value is now:`, cell.Value.(int))

	cell03, _ := ss.GetCell("Sheet1", 0, 3)
	ss.UpdateCell(cell03.Uuid, `=Sum(A1:C4)`)
	println(cell03.Value.(int))
}
