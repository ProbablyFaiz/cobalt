package src

type Spreadsheet struct {
	Sheets map[string]*Sheet
}

type Sheet struct {
	Spreadsheet *Spreadsheet
	Cells       [][]Cell
}

type Cell struct {
	Sheet *Sheet
	// Value can be a string, int, float, bool, or nil
	Value      interface{}
	Formula    FormulaNode
	RawContent string
	Dirty      bool
	Dependents []*Cell
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
