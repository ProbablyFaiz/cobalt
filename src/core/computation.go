package core

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
		// If the parent is a range, decrement the ref count.
		if rangeParent, ok := ss.RangeMap[parent]; ok {
			rangeParent.RefCount--
			if rangeParent.RefCount == 0 {
				// We don't delete immediately, in case the range is re-added.
				ss.RangesMarkedForDeletion = append(ss.RangesMarkedForDeletion, rangeParent)
			}
		}
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
		rangeSheet := ref.Start.Sheet
		if rangeSheet == nil {
			rangeSheet = cell.Sheet
		}
		// Create a range object for the reference.
		currRange := &Range{
			Uuid:     ss.getNextId(),
			StartRow: ref.Start.Row,
			EndRow:   ref.End.Row,
			StartCol: ref.Start.Col,
			EndCol:   ref.End.Col,
			Sheet:    rangeSheet,
			RefCount: 1,
		}
		// Check if we already have this exact range.
		if ss.RangeDuplicateMap[currRange.getCompareKey()] != nil {
			currRange = ss.RangeDuplicateMap[currRange.getCompareKey()]
			currRange.RefCount++
		} else {
			ss.RangeDuplicateMap[currRange.getCompareKey()] = currRange
			rangeSheet.RangeTree.Add(currRange)
			ss.RangeMap[currRange.Uuid] = currRange
			ss.Children[currRange.Uuid] = mapset.NewThreadUnsafeSet[ReferenceId]()
		}
		ref.ResolvedUuid = currRange.Uuid

		ss.Parents[cell.Uuid].Add(currRange.Uuid)
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

	ss := cell.Sheet.Spreadsheet
	ss.DirtySet.Add(cell.Uuid)

	visited.Add(cell.Uuid)
	// Query all ranges that contain the cell.
	ranges := cell.Sheet.RangeTree.Query(cell.toRange())
	for _, cr := range ranges {
		currRange := cr.(*Range)
		err := currRange.dirty(visited)
		if err != nil {
			return err
		}

		if ss.RangeDirtyParents[currRange.Uuid] == nil {
			ss.RangeDirtyParents[currRange.Uuid] = mapset.NewThreadUnsafeSet[ReferenceId]()
		}
		// Tracks the cells that need to be recomputed before recomputing the range.
		ss.RangeDirtyParents[currRange.Uuid].Add(cell.Uuid)
	}

	// Dirty all dependent cells.
	if ss.Children[cell.Uuid] == nil {
		for dependent := range ss.Children[cell.Uuid].Iter() {
			err := ss.CellMap[dependent].dirty(visited)
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

func (cr *Range) dirty(visited mapset.Set[ReferenceId]) error {
	if visited.Contains(cr.Uuid) {
		return fmt.Errorf("cycle detected")
	}

	ss := cr.Sheet.Spreadsheet
	visited.Add(cr.Uuid)

	// Dirty all dependent cells.
	if ss.Children[cr.Uuid] != nil {
		for dependent := range ss.Children[cr.Uuid].Iter() {
			err := ss.CellMap[dependent].dirty(visited)
			if err != nil {
				return err
			}
		}
	}
	visited.Remove(cr.Uuid)
	return nil
}

func (ss *Spreadsheet) recomputeValues() {
	for cellId := range ss.DirtySet.Iter() {
		currCell := ss.CellMap[cellId]
		res, err := (*currCell.Formula).eval(&EvalContext{
			Cell: currCell,
		})
		currCell.Value, currCell.Error = res, err
	}
	ss.DirtySet.Clear()
}

func (ss *Spreadsheet) cleanupRanges() {
	for _, rangeObj := range ss.RangesMarkedForDeletion {
		if rangeObj.RefCount > 0 {
			continue
		}
		rangeObj.Sheet.RangeTree.Delete(rangeObj)
		delete(ss.RangeMap, rangeObj.Uuid)
		delete(ss.RangeDuplicateMap, rangeObj.getCompareKey())
	}
}
