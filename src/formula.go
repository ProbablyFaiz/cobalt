package src

import "fmt"

type FormulaNode interface {
	Eval(ctx *EvalContext) (interface{}, error)
	GetRefs() []ReferenceNode
	ToFormula() string
}

type LiteralNode struct {
	Value interface{}
}

type ReferenceNode struct {
	Row          int
	Col          int
	Sheet        *Sheet // If nil, then the cell is in the current sheet
	ResolvedUuid cellId // If nil, then the cell has not been resolved yet
}

type FunctionNode struct {
	Name string
	Args []FormulaNode
}

type EvalContext struct {
	Cell *Cell
	// TODO: Add more stuff here
}

// GetRefs implementations

func (ln *LiteralNode) GetRefs() []ReferenceNode {
	return nil
}

func (rn *ReferenceNode) GetRefs() []ReferenceNode {
	return []ReferenceNode{*rn}
}

func (fn *FunctionNode) GetRefs() []ReferenceNode {
	var refs []ReferenceNode
	for _, arg := range fn.Args {
		refs = append(refs, arg.GetRefs()...)
	}
	return refs
}

// ToFormula implementations

func (ln *LiteralNode) ToFormula() string {
	return fmt.Sprintf("%v", ln.Value)
}

func (rn *ReferenceNode) ToFormula() string {
	// Convert rn.Col to a letter or letters (e.g. 0 -> A, 1 -> B, 26 -> AA, etc.)
	colLetters := ""
	for i := rn.Col; i >= 0; i = i/26 - 1 {
		colLetters = string(rune('A'+i%26)) + colLetters
	}
	return fmt.Sprintf("%s%d", colLetters, rn.Row+1)
}

func (fn *FunctionNode) ToFormula() string {
	formula := fn.Name + "("
	for i, arg := range fn.Args {
		formula += arg.ToFormula()
		if i != len(fn.Args)-1 {
			formula += ", "
		}
	}
	formula += ")"
	return formula
}

// Eval implementations in src/eval.go
