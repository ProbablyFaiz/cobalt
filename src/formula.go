package src

import (
	"fmt"
	"strconv"
)

type FormulaNode interface {
	eval(ctx *EvalContext) (interface{}, error)
	getRefs() []*ReferenceNode
	getRangeRefs() []*RangeNode
	toFormula() string
}

type LiteralNode struct {
	Value interface{}
}

type ReferenceNode struct {
	Row          int
	Col          int
	Sheet        *Sheet      // If nil, then the cell is in the current sheet
	ResolvedUuid ReferenceId // If nil, then the cell has not been resolved yet
}

type RangeNode struct {
	Start        *ReferenceNode
	End          *ReferenceNode
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
	return make([]*ReferenceNode, 0)
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

// getRangeRefs implementations

func (ln *LiteralNode) getRangeRefs() []*RangeNode {
	return make([]*RangeNode, 0)
}

func (rn *ReferenceNode) getRangeRefs() []*RangeNode {
	return make([]*RangeNode, 0)
}

func (rn *RangeNode) getRangeRefs() []*RangeNode {
	return []*RangeNode{rn}
}

func (fn *FunctionNode) getRangeRefs() []*RangeNode {
	refs := make([]*RangeNode, 0)
	for _, arg := range fn.Args {
		refs = append(refs, arg.getRangeRefs()...)
	}
	return refs
}

func (_ *NilNode) getRangeRefs() []*RangeNode {
	return make([]*RangeNode, 0)
}

// toFormula implementations

func (ln *LiteralNode) toFormula() string {
	return fmt.Sprintf("%v", ln.Value)
}

func (rn *ReferenceNode) toFormula() string {
	return getA1Notation(rn.Row, rn.Col)
}

func (rn *RangeNode) toFormula() string {
	return fmt.Sprintf("%s:%s", rn.Start.toFormula(), rn.End.toFormula())
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

// misc

func (r *Range) getCompareKey() string {
	return fmt.Sprintf("%s:%d:%d,%d:%d", strconv.FormatUint(uint64(r.Sheet.Uuid), 10), r.StartRow, r.StartCol, r.EndRow, r.EndCol)
}
