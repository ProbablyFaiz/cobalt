package src

import (
	mapset "github.com/deckarep/golang-set/v2"
	"sync"
)

type Spreadsheet struct {
	Sheets     map[string]*Sheet
	Mutex      sync.Mutex
	CellMap    map[cellId]*Cell
	DirtySet   mapset.Set[cellId]
	Dependents map[cellId]mapset.Set[cellId]
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

type FormulaNode interface {
	Eval(ctx *EvalContext) (interface{}, error)
}

type LiteralNode struct {
	Value interface{}
}

type ReferenceNode struct {
	Row   int
	Col   int
	Sheet *Sheet // If nil, then the cell is in the current sheet
}

type FunctionNode struct {
	Name string
	Args []FormulaNode
}

type EvalContext struct {
	Cell *Cell
	// TODO: Add more stuff here
}
