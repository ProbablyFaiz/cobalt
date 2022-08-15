package main

import . "pasado/src"

func main() {
	// Initialize the spreadsheet.
	ss := Spreadsheet{
		Sheets: make(map[string]*Sheet),
	}
	// Add a sheet.
	ss.AddSheet("Sheet1")
	// UpdateCell a cell.
	ss.UpdateCell("Sheet1", 0, 0, "Hello, world!")
}
