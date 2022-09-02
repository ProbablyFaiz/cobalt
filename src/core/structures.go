package core

import (
	"github.com/Workiva/go-datastructures/augmentedtree"
	"github.com/deckarep/golang-set/v2"
	"sync"
)

// TODO: Rename all the ReferenceId's that are called Uuid's, they are not uuid's anymore.

type ReferenceId uint64

type Spreadsheet struct {
	Sheets map[string]*Sheet

	CellMap  map[ReferenceId]*Cell
	DirtySet mapset.Set[ReferenceId]
	// Note that Children and Parents do not imply a nested structure, only dependencies.
	Parents  map[ReferenceId]mapset.Set[ReferenceId]
	Children map[ReferenceId]mapset.Set[ReferenceId]

	RangeMap                map[ReferenceId]*Range
	RangeDuplicateMap       map[string]*Range
	RangeDirtyParents       map[ReferenceId]mapset.Set[ReferenceId]
	RangesMarkedForDeletion []*Range

	Mutex  sync.Mutex
	NextId ReferenceId
}

type Sheet struct {
	Uuid        ReferenceId
	Spreadsheet *Spreadsheet
	Cells       [][]Cell
	RangeTree   augmentedtree.Tree
}

type Cell struct {
	Value interface{}
	Error error

	Uuid  ReferenceId
	Sheet *Sheet
	Row   int
	Col   int

	Formula    *FormulaNode
	RawContent string
}

type Range struct {
	Uuid     ReferenceId
	StartRow int
	EndRow   int
	StartCol int
	EndCol   int
	Sheet    *Sheet
	RefCount int
}

type ValueContainer interface {
	GetValue() interface{}
	GetUuid() ReferenceId
}

func NewSpreadsheet() *Spreadsheet {
	ss := &Spreadsheet{
		Sheets:            make(map[string]*Sheet),
		CellMap:           make(map[ReferenceId]*Cell),
		RangeMap:          make(map[ReferenceId]*Range),
		DirtySet:          mapset.NewThreadUnsafeSet[ReferenceId](),
		Parents:           make(map[ReferenceId]mapset.Set[ReferenceId]),
		Children:          make(map[ReferenceId]mapset.Set[ReferenceId]),
		RangeDuplicateMap: make(map[string]*Range),
		RangeDirtyParents: make(map[ReferenceId]mapset.Set[ReferenceId]),
		Mutex:             sync.Mutex{},
		NextId:            0,
	}
	// Add the default sheet.
	err := ss.AddSheet("Sheet1")
	if err != nil {
		return nil
	}
	return ss
}

func (ss *Spreadsheet) getNextId() ReferenceId {
	ss.NextId++
	return ss.NextId
}
