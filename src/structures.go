package src

import (
	"github.com/deckarep/golang-set/v2"
	"sync"
)

type Spreadsheet struct {
	Sheets   map[string]*Sheet
	Mutex    sync.Mutex
	CellMap  map[cellId]*Cell
	DirtySet mapset.Set[cellId]
	// Note that Children and Parents do not imply a nested structure, only dependencies.
	Parents  map[cellId]mapset.Set[cellId]
	Children map[cellId]mapset.Set[cellId]
}

type cellId string

type Sheet struct {
	Spreadsheet *Spreadsheet
	Cells       [][]Cell
}

type Cell struct {
	Value interface{}

	Uuid  cellId
	Sheet *Sheet

	Formula    *FormulaNode
	RawContent string
}

func NewSpreadsheet() *Spreadsheet {
	ss := &Spreadsheet{
		Sheets:   make(map[string]*Sheet),
		CellMap:  make(map[cellId]*Cell),
		DirtySet: mapset.NewThreadUnsafeSet[cellId](),
		Parents:  make(map[cellId]mapset.Set[cellId]),
		Children: make(map[cellId]mapset.Set[cellId]),
	}
	// Add the default sheet.
	err := ss.AddSheet("Sheet1")
	if err != nil {
		return nil
	}
	return ss
}
