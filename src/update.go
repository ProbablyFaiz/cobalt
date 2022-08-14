package src

func (ss *Spreadsheet) Update(sheetName string, row int, col int, content string) {
	sheet := ss.Sheets[sheetName]
	cell := sheet.Cells[row][col]
	cell.RawContent = content
	cell.Dirty = true
}
