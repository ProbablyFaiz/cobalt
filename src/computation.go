package src

import (
	"fmt"
	"github.com/deckarep/golang-set/v2"
)

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
		child := ref.Sheet.Cells[ref.Row][ref.Col]
		ref.ResolvedUuid = child.Uuid
		// Check if ss.Parents[child.Uuid] is nil and initialize it if so.
		if ss.Parents[child.Uuid] == nil {
			ss.Parents[child.Uuid] = mapset.NewThreadUnsafeSet[ReferenceId]()
		}
		ss.Parents[child.Uuid].Add(cell.Uuid)
		ss.Children[cell.Uuid].Add(child.Uuid)
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

	visited.Add(cell.Uuid)
	// Query all ranges that contain the cell.
	ranges := cell.Sheet.RangeTree.Query(cell.toRange())
	for _, cr := range ranges {
		currRange := cr.(*Range)
		if visited.Contains(currRange.Uuid) {
			return fmt.Errorf("cycle detected")
		}
		visited.Add(currRange.Uuid)
	}

	// Dirty all dependent cells.
	if spreadsheet.Children[cell.Uuid] != nil {
		for dependent := range spreadsheet.Children[cell.Uuid].Iter() {
			err := spreadsheet.CellMap[dependent].dirty(visited)
			if err != nil {
				return err
			}
		}
	}
	visited.Remove(cell.Uuid)

	return nil
}
