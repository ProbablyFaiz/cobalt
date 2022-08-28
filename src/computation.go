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

	rangeRefs := (*cell.Formula).getRangeRefs()
	for _, ref := range rangeRefs {
		rangeSheet := ref.To.Sheet
		// Create a range object for the reference.
		currRange := &Range{
			Uuid:    ss.getNextId(),
			FromRow: ref.From.Row,
			ToRow:   ref.To.Row,
			FromCol: ref.From.Col,
			ToCol:   ref.To.Col,
			Sheet:   rangeSheet,
		}
		// Check if we already have this exact range.
		if ss.RangeDuplicateMap[currRange.getCompareKey()] != nil {
			currRange = ss.RangeDuplicateMap[currRange.getCompareKey()]
		} else {
			ss.RangeDuplicateMap[currRange.getCompareKey()] = currRange
			rangeSheet.RangeTree.Add(currRange)
			ss.RangeMap[currRange.Uuid] = currRange
			ss.Parents[currRange.Uuid] = mapset.NewThreadUnsafeSet[ReferenceId]()
			ss.Children[currRange.Uuid] = mapset.NewThreadUnsafeSet[ReferenceId]()
		}

		// Add the cell to the range's children.
		ss.Children[currRange.Uuid].Add(cell.Uuid)
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

		if spreadsheet.RangeDirtyParents[currRange.Uuid] == nil {
			spreadsheet.RangeDirtyParents[currRange.Uuid] = mapset.NewThreadUnsafeSet[ReferenceId]()
		}
		// Tracks the cells that need to be recomputed before recomputing the range.
		spreadsheet.RangeDirtyParents[cell.Uuid].Add(currRange.Uuid)
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
	for _, cr := range ranges {
		visited.Remove(cr.(*Range).Uuid)
	}

	return nil
}

func (ss *Spreadsheet) recomputeValues() {
	for cellId := range ss.DirtySet.Iter() {
		currCell := ss.CellMap[cellId]
		res, err := (*currCell.Formula).eval(&EvalContext{
			Cell: currCell,
		})
		currCell.Value, currCell.Error = res, err
		currCell.Value = res
	}
}
