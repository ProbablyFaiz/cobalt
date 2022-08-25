package src

import (
	"fmt"
	"github.com/deckarep/golang-set/v2"
	"github.com/google/uuid"
)

const DefaultNumCols = 26
const DefaultNumRows = 1000

func (ss *Spreadsheet) UpdateCell(sheetName string, cellId cellId, content string) error {
	ss.Mutex.Lock()

	cell := ss.CellMap[cellId]

	cell.RawContent = content

	newFormula, err := Parse(content)
	// TODO: Decorate all these errors with something useful.
	if err != nil {
		ss.Mutex.Unlock()
		return err
	}
	cell.Formula = &newFormula

	err = cell.Dirty(nil)
	if err != nil {
		ss.Mutex.Unlock()
		return err
	}

	for cellId := range ss.DirtySet.Iter() {
		currCell := ss.CellMap[cellId]
		res, err := (*currCell.Formula).Eval(nil)
		if err != nil {
			ss.Mutex.Unlock()
			return err
		}
		currCell.Value = res
	}

	ss.Mutex.Unlock()
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

func (cell *Cell) UpdateDependencies() error {
	// Remove existing parents, as well as those parents' corresponding children.
	for parent := range cell.Sheet.Spreadsheet.Parents[cell.Uuid].Iter() {
		cell.Sheet.Spreadsheet.Children[parent].Remove(cell.Uuid)
	}
	cell.Sheet.Spreadsheet.Parents[cell.Uuid].Clear()

	ss := cell.Sheet.Spreadsheet
	refs := (*cell.Formula).GetRefs()
	for _, ref := range refs {
		if ref.Sheet == nil {
			ref.Sheet = cell.Sheet
		}
		// Resolve the reference.
		parent := ref.Sheet.Cells[ref.Row][ref.Col]
		ref.ResolvedUuid = parent.Uuid
		// Check if ss.Parents[parent.Uuid] is nil and initialize it if so.
		if ss.Parents[parent.Uuid] == nil {
			ss.Parents[parent.Uuid] = mapset.NewThreadUnsafeSet[cellId]()
		}
		ss.Parents[parent.Uuid].Add(cell.Uuid)
		// Check if ss.Children[cell.Uuid] is nil and initialize it if so.
		if ss.Children[cell.Uuid] == nil {
			ss.Children[cell.Uuid] = mapset.NewThreadUnsafeSet[cellId]()
		}
		ss.Children[cell.Uuid].Add(parent.Uuid)
	}
	return nil
}

func (cell *Cell) Dirty(visited mapset.Set[cellId]) error {
	if visited == nil {
		visited = mapset.NewThreadUnsafeSet[cellId]()
	}

	if visited.Contains(cell.Uuid) {
		return fmt.Errorf("cycle detected")
	}

	spreadsheet := cell.Sheet.Spreadsheet
	spreadsheet.DirtySet.Add(cell.Uuid)
	// Dirty all dependent cells.
	if spreadsheet.Children[cell.Uuid] != nil {
		visited.Add(cell.Uuid)
		for dependent := range spreadsheet.Children[cell.Uuid].Iter() {
			err := spreadsheet.CellMap[dependent].Dirty(visited)
			if err != nil {
				return err
			}
		}
		visited.Remove(cell.Uuid)
	}
	return nil
}
