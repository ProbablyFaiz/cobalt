package src

import (
	"fmt"
	"github.com/Workiva/go-datastructures/augmentedtree"
)

const DefaultNumCols = 26
const DefaultNumRows = 1000

func (ss *Spreadsheet) UpdateCell(cellUuid ReferenceId, content string) {
	// TODO: This should be structured as a synchronous add to a queue, and a separate goroutine should
	//  handle the updates.

	ss.Mutex.Lock()
	defer ss.Mutex.Unlock()

	cell := ss.CellMap[cellUuid]

	cell.RawContent = content

	newFormula, err := Parse(content)
	if err != nil {
		cell.Formula, cell.Value, cell.Error = nil, nil, err
		return
	}

	cell.Formula = &newFormula
	cell.updateDependencies()

	err = cell.dirty(nil)
	if err != nil {
		cell.Value, cell.Error = nil, err
		return
	}

	for currCellId := range ss.DirtySet.Iter() {
		currCell := ss.CellMap[currCellId]
		res, err := (*currCell.Formula).eval(&EvalContext{
			Cell: currCell,
		})
		currCell.Value, currCell.Error = res, err
		currCell.Value = res
	}
}

func (ss *Spreadsheet) AddSheet(sheetName string) error {
	ss.Mutex.Lock()
	defer ss.Mutex.Unlock()

	// Check if the sheet already exists.
	if _, ok := ss.Sheets[sheetName]; ok {
		ss.Mutex.Unlock()
		return fmt.Errorf("sheet %s already exists", sheetName)
	}

	cells := make([][]Cell, DefaultNumRows)
	newSheet := &Sheet{
		Spreadsheet: ss,
		Cells:       cells,
		RangeTree:   augmentedtree.New(2),
	}
	ss.Sheets[sheetName] = newSheet

	for i := 0; i < DefaultNumRows; i++ {
		cells[i] = make([]Cell, DefaultNumCols)
		for j := 0; j < DefaultNumCols; j++ {
			var formula FormulaNode = &NilNode{}
			cells[i][j] = Cell{
				Uuid:    ss.GetNextId(),
				Sheet:   newSheet,
				Formula: &formula,
				Value:   nil,
			}
			ss.CellMap[cells[i][j].Uuid] = &cells[i][j]
		}
	}
	return nil
}
