package src

import (
	"fmt"
	"github.com/Workiva/go-datastructures/augmentedtree"
	"github.com/deckarep/golang-set/v2"
	"github.com/google/uuid"
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
				Uuid:    ReferenceId(uuid.New().String()),
				Sheet:   newSheet,
				Formula: &formula,
				Value:   nil,
			}
			ss.CellMap[cells[i][j].Uuid] = &cells[i][j]
		}
	}
	return nil
}

func (cell *Cell) updateDependencies() {
	ss := cell.Sheet.Spreadsheet
	// Check if ss.Children[cell.Uuid] and ss.Parents[cell.Uuid] are nil and if so, initialize them.
	if ss.Children[cell.Uuid] == nil {
		ss.Children[cell.Uuid] = mapset.NewThreadUnsafeSet[ReferenceId]()
	}
	if ss.Parents[cell.Uuid] == nil {
		ss.Parents[cell.Uuid] = mapset.NewThreadUnsafeSet[ReferenceId]()
	}

	// Remove existing parents, as well as those parents' corresponding children.
	for parent := range ss.Parents[cell.Uuid].Iter() {
		ss.Children[parent].Remove(cell.Uuid)
	}
	ss.Parents[cell.Uuid].Clear()

	refs := (*cell.Formula).getRefs()
	for _, ref := range refs {
		if ref.Sheet == nil {
			ref.Sheet = cell.Sheet
		}
		// Resolve the reference.
		parent := ref.Sheet.Cells[ref.Row][ref.Col]
		ref.ResolvedUuid = parent.Uuid
		// Check if ss.Parents[parent.Uuid] is nil and initialize it if so.
		if ss.Parents[parent.Uuid] == nil {
			ss.Parents[parent.Uuid] = mapset.NewThreadUnsafeSet[ReferenceId]()
		}
		ss.Parents[parent.Uuid].Add(cell.Uuid)
		ss.Children[cell.Uuid].Add(parent.Uuid)
	}
}

func (cell *Cell) dirty(visited mapset.Set[ReferenceId]) error {
	if visited == nil {
		visited = mapset.NewThreadUnsafeSet[ReferenceId]()
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
			err := spreadsheet.CellMap[dependent].dirty(visited)
			if err != nil {
				return err
			}
		}
		visited.Remove(cell.Uuid)
	}
	return nil
}
