package src

import "fmt"

type FormulaNode interface {
	eval(ctx *EvalContext) (interface{}, error)
	getRefs() []*ReferenceNode
	toFormula() string
}

type LiteralNode struct {
	Value interface{}
}

type Referencer interface {
	Covers(r *ReferenceNode) bool
}

type ReferenceNode struct {
	Row          int
	Col          int
	Sheet        *Sheet      // If nil, then the cell is in the current sheet
	ResolvedUuid ReferenceId // If nil, then the cell has not been resolved yet
}

type RangeNode struct {
	From         *ReferenceNode
	To           *ReferenceNode
	ResolvedUuid ReferenceId // If nil, then the range has not been resolved yet
}

type FunctionNode struct {
	Name string
	Args []FormulaNode
}

type NilNode struct {
}

type EvalContext struct {
	Cell *Cell
	// TODO: Add more stuff here
}

// getRefs implementations

func (ln *LiteralNode) getRefs() []*ReferenceNode {
	return make([]*ReferenceNode, 0)
}

func (rn *ReferenceNode) getRefs() []*ReferenceNode {
	return []*ReferenceNode{rn}
}

func (rn *RangeNode) getRefs() []*ReferenceNode {
	// TODO: We probably want to set up global range tracking with segment
	//  trees, so we can efficiently dirty ranges etc.
	panic("not implemented")
}

func (fn *FunctionNode) getRefs() []*ReferenceNode {
	refs := make([]*ReferenceNode, 0)
	for _, arg := range fn.Args {
		refs = append(refs, arg.getRefs()...)
	}
	return refs
}

func (_ *NilNode) getRefs() []*ReferenceNode {
	return make([]*ReferenceNode, 0)
}

// toFormula implementations

func (ln *LiteralNode) toFormula() string {
	return fmt.Sprintf("%v", ln.Value)
}

func (rn *ReferenceNode) toFormula() string {
	return getA1Notation(rn.Row, rn.Col)
}

func (rn *RangeNode) toFormula() string {
	return fmt.Sprintf("%s:%s", rn.From.toFormula(), rn.To.toFormula())
}

func (fn *FunctionNode) toFormula() string {
	formula := fn.Name + "("
	for i, arg := range fn.Args {
		formula += arg.toFormula()
		if i != len(fn.Args)-1 {
			formula += ", "
		}
	}
	formula += ")"
	return formula
}

func (_ *NilNode) toFormula() string {
	return ""
}

// eval implementations in src/eval.go
