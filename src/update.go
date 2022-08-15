package src

import (
	"fmt"
	mapset "github.com/deckarep/golang-set/v2"
	"github.com/google/uuid"
)

const DefaultNumCols = 26
const DefaultNumRows = 1000

func (ss *Spreadsheet) UpdateCell(sheetName string, row int, col int, content string) error {
	ss.Mutex.Lock()

	sheet := ss.Sheets[sheetName]
	cell := sheet.Cells[row][col]
	cell.RawContent = content
	err := cell.Dirty(nil)

	ss.Mutex.Unlock()

	if err != nil {
		return err
	}
	return nil
}

func (ss *Spreadsheet) AddSheet(sheetName string) error {
	ss.Mutex.Lock()

	// Check if the sheet already exists.
	if _, ok := ss.Sheets[sheetName]; ok {
		ss.Mutex.Unlock()
		return fmt.Errorf("sheet %s already exists", sheetName)
	}

	cells := make([][]Cell, DefaultNumRows)
	newSheet := &Sheet{
		Spreadsheet: ss,
		Cells:       cells,
	}
	ss.Sheets[sheetName] = newSheet

	for i := 0; i < DefaultNumRows; i++ {
		cells[i] = make([]Cell, DefaultNumCols)
		for j := 0; j < DefaultNumCols; j++ {
			cells[i][j] = Cell{
				Uuid:  cellId(uuid.New().String()),
				Sheet: newSheet,
			}
		}
	}
	ss.Mutex.Unlock()
	return nil
}

func (cell *Cell) Dirty(visited mapset.Set[cellId]) error {
	// If visited is nil, initialize it.
	if visited == nil {
		visited = mapset.NewThreadUnsafeSet[cellId]()
	}

	// If the cell is already visited, cycle error.
	if visited.Contains(cell.Uuid) {
		return fmt.Errorf("cycle detected")
	}
	visited.Add(cell.Uuid)

	spreadsheet := cell.Sheet.Spreadsheet
	spreadsheet.DirtySet.Add(cell.Uuid)
	// Dirty all dependent cells.
	for dependent := range spreadsheet.Dependents[cell.Uuid].Iter() {
		err := spreadsheet.CellMap[dependent].Dirty(visited)
		if err != nil {
			return err
		}
	}
	visited.Remove(cell.Uuid)
	return nil
}
