package src

import (
	"github.com/Workiva/go-datastructures/augmentedtree"
	"github.com/deckarep/golang-set/v2"
	"sync"
)

type ReferenceId string

type Spreadsheet struct {
	Sheets   map[string]*Sheet
	Mutex    sync.Mutex
	CellMap  map[ReferenceId]*Cell
	RangeMap map[ReferenceId]*Range
	DirtySet mapset.Set[ReferenceId]
	// Note that Children and Parents do not imply a nested structure, only dependencies.
	Parents  map[ReferenceId]mapset.Set[ReferenceId]
	Children map[ReferenceId]mapset.Set[ReferenceId]
}

type Sheet struct {
	Spreadsheet *Spreadsheet
	Cells       [][]Cell
	RangeTree   augmentedtree.Tree
}

type Cell struct {
	Value interface{}
	Error error

	Uuid  ReferenceId
	Sheet *Sheet

	Formula    *FormulaNode
	RawContent string
}

type Range struct {
	Uuid     ReferenceId
	From     int
	To       int
	Sheet    *Sheet
	RefCount int
}

type ValueContainer interface {
	GetValue() interface{}
	GetUuid() ReferenceId
}

func NewSpreadsheet() *Spreadsheet {
	ss := &Spreadsheet{
		Sheets:   make(map[string]*Sheet),
		CellMap:  make(map[ReferenceId]*Cell),
		RangeMap: make(map[ReferenceId]*Range),
		DirtySet: mapset.NewThreadUnsafeSet[ReferenceId](),
		Parents:  make(map[ReferenceId]mapset.Set[ReferenceId]),
		Children: make(map[ReferenceId]mapset.Set[ReferenceId]),
	}
	// Add the default sheet.
	err := ss.AddSheet("Sheet1")
	if err != nil {
		return nil
	}
	return ss
}
