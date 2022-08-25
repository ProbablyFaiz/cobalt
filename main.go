package main

import (
	"fmt"
	. "pasado/src"
)

func main() {
	// Initialize the spreadsheet.
	ss := NewSpreadsheet()

	// Add a formula to a cell.
	err := ss.UpdateCell("Sheet1", 0, 0, `=CONCAT("Hello, ", "world!")`)
	if err != nil {
		panic(err)
	}

	// Get the value of the cell and print it.
	cell, err := ss.GetCell("Sheet1", 0, 0)
	if err != nil {
		panic(err)
	}
	value, err := cell.GetValue()
	if err != nil {
		panic(err)
	}
	fmt.Println(value)
}
